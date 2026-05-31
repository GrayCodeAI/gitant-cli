// git-remote-gitant is a git remote helper for the gitant:// URL scheme.
//
// Usage: git clone gitant://<repo-id> --remote http://localhost:7777
//
// Git automatically invokes this binary when it encounters a gitant:// URL.
// The helper communicates with the gitant daemon over HTTP to fetch/push objects.
//
// Fetch uses packfile-based transfer: all requested objects are fetched in a
// single request to the daemon's git-upload-pack endpoint, which returns a
// packfile that git can ingest directly. This is much faster than object-by-object
// fetching for large repositories.
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/GrayCodeAI/gitant-cli/internal/cli"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: git-remote-gitant <name> <url>\n")
		os.Exit(1)
	}

	repoID := os.Args[2]

	// Parse gitant://<repo-id> URL
	repoID = strings.TrimPrefix(repoID, "gitant://")
	repoID = strings.TrimPrefix(repoID, "gitant:")

	// Get daemon URL from env or default
	daemonURL := os.Getenv("GITANT_DAEMON_URL")
	if daemonURL == "" {
		daemonURL = "http://localhost:7777"
	}

	client := cli.NewClient(daemonURL)
	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	// Batch buffer for fetch requests — collected until blank line, then
	// sent as a single packfile request for efficiency.
	var pendingFetches []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		switch {
		case line == "capabilities":
			fmt.Fprintln(writer, "fetch")
			fmt.Fprintln(writer, "push")
			fmt.Fprintln(writer)
			writer.Flush()

		case line == "":
			// End of command batch — flush any pending fetches as a single
			// packfile request, then exit.
			if len(pendingFetches) > 0 {
				flushFetches(daemonURL, repoID, pendingFetches, writer)
				pendingFetches = nil
			}
			writer.Flush()
			return

		case line == "list" || line == "list for-push":
			var result struct {
				Refs []struct {
					Name string `json:"name"`
					Hash string `json:"hash"`
				} `json:"refs"`
			}
			if err := client.Get(fmt.Sprintf("/api/v1/repos/%s/refs", url.PathEscape(repoID)), &result); err != nil {
				fmt.Fprintf(os.Stderr, "error listing refs: %v\n", err)
				fmt.Fprintln(writer)
				writer.Flush()
				continue
			}
			for _, ref := range result.Refs {
				fmt.Fprintf(writer, "%s %s\n", ref.Hash, ref.Name)
			}
			fmt.Fprintln(writer)
			writer.Flush()

		case strings.HasPrefix(line, "fetch "):
			// Buffer the fetch request for batch processing
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				pendingFetches = append(pendingFetches, parts[1])
			}

		case strings.HasPrefix(line, "push "):
			// Flush any pending fetches before pushing
			if len(pendingFetches) > 0 {
				flushFetches(daemonURL, repoID, pendingFetches, writer)
				pendingFetches = nil
			}

			// push <refspec>
			refspec := strings.TrimPrefix(line, "push ")

			// Parse refspec: +refs/heads/main:refs/heads/main
			forcePush := strings.HasPrefix(refspec, "+")
			refspec = strings.TrimPrefix(refspec, "+")
			parts := strings.SplitN(refspec, ":", 2)
			if len(parts) != 2 {
				fmt.Fprintf(os.Stderr, "invalid refspec: %s\n", refspec)
				fmt.Fprintln(writer)
				writer.Flush()
				continue
			}

			localRef := parts[0]
			remoteRef := parts[1]

			// Use CLI push with the specific ref
			if forcePush {
				fmt.Fprintf(os.Stderr, "push +%s -> %s (force)\n", localRef, remoteRef)
			} else {
				fmt.Fprintf(os.Stderr, "push %s -> %s\n", localRef, remoteRef)
			}
			if err := cli.Push(".", daemonURL, repoID, localRef); err != nil {
				fmt.Fprintf(os.Stderr, "push error: %v\n", err)
			}
			fmt.Fprintln(writer)
			writer.Flush()
		}
	}
}

// flushFetches sends a single packfile request for all pending want hashes.
// This is much faster than fetching objects one-by-one for large repos.
func flushFetches(daemonURL, repoID string, wants []string, writer *bufio.Writer) {
	if len(wants) == 0 {
		return
	}

	// Build pkt-line want request
	var pktBuf bytes.Buffer
	for _, want := range wants {
		fmt.Fprintf(&pktBuf, "%04xwant %s\n", len("want ")+len(want)+1, want)
	}
	pktBuf.WriteString("0000") // flush packet
	fmt.Fprintf(&pktBuf, "%04xdone\n", len("done\n")+4)

	endpoint := fmt.Sprintf("%s/api/v1/repos/%s/git-upload-pack", daemonURL, url.PathEscape(repoID))
	resp, err := http.Post(endpoint, "application/x-git-upload-pack-request", &pktBuf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "packfile fetch error: %v\n", err)
		// Fall back to object-by-object fetch
		fetchObjectsFallback(daemonURL, repoID, wants, writer)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "packfile fetch failed (HTTP %d), falling back to object fetch\n", resp.StatusCode)
		fetchObjectsFallback(daemonURL, repoID, wants, writer)
		return
	}

	// Read the entire packfile response and write to stdout
	packData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading packfile response: %v\n", err)
		return
	}

	fmt.Fprintf(os.Stderr, "fetch: packfile %d bytes for %d objects\n", len(packData), len(wants))
	writer.Write(packData)
}

// fetchObjectsFallback fetches objects one-by-one when packfile transfer fails.
func fetchObjectsFallback(daemonURL, repoID string, wants []string, writer *bufio.Writer) {
	client := cli.NewClient(daemonURL)
	for _, sha1 := range wants {
		var obj struct {
			Hash    string `json:"hash"`
			Type    string `json:"type"`
			Content []byte `json:"content"`
		}
		if err := client.Get(fmt.Sprintf("/api/v1/repos/%s/objects/%s", url.PathEscape(repoID), url.PathEscape(sha1)), &obj); err != nil {
			fmt.Fprintf(os.Stderr, "error fetching object %s: %v\n", sha1, err)
		} else {
			writer.Write(obj.Content)
			fmt.Fprintf(os.Stderr, "fetch %s (%d bytes)\n", sha1, len(obj.Content))
		}
	}
}
