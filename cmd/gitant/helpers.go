package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"text/tabwriter"

	"github.com/GrayCodeAI/gitant-cli/internal/cli"
	"github.com/GrayCodeAI/gitant-cli/internal/config"
	"github.com/spf13/cobra"
)

func daemonURLFromCmd(cmd *cobra.Command) string {
	flagURL, _ := cmd.Flags().GetString("daemon-url")
	return flagURL
}

func newClient(cmd *cobra.Command) *cli.Client {
	return cli.NewClient(daemonURLFromCmd(cmd))
}

func addDaemonURLFlag(cmds ...*cobra.Command) {
	for _, c := range cmds {
		c.Flags().String("daemon-url", "", "Daemon URL (default: GITANT_DAEMON_URL or config)")
	}
}

func printJSON(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}

func repoPath(repo, suffix string) string {
	return fmt.Sprintf("/api/v1/repos/%s%s", url.PathEscape(repo), suffix)
}

// repoPathSegments builds an escaped URL path from repo and additional segments.
// Each segment is url.PathEscape'd. Query params should be appended separately.
func repoPathSegments(repo string, segments ...string) string {
	path := "/api/v1/repos/" + url.PathEscape(repo)
	for _, s := range segments {
		path += "/" + url.PathEscape(s)
	}
	return path
}

// apiPath builds an escaped URL path from a base and segments.
func apiPath(base string, segments ...string) string {
	path := base
	for _, s := range segments {
		path += "/" + url.PathEscape(s)
	}
	return path
}

// queryEscape is a convenience wrapper for url.QueryEscape.
func queryEscape(s string) string {
	return url.QueryEscape(s)
}

func webBaseURL() string {
	if u := os.Getenv("GITANT_WEB_URL"); u != "" {
		return trimRightSlash(u)
	}
	if s, err := config.Load(); err == nil && s.WebURL != "" {
		return trimRightSlash(s.WebURL)
	}
	return "http://localhost:3303"
}

func trimRightSlash(s string) string {
	for len(s) > 1 && s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}
	return s
}

func newTabWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
}
