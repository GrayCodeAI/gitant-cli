package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

var browseCmd = &cobra.Command{
	Use:   "browse [resource]",
	Short: "Open repository or issue in the web browser (like gh browse)",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		issueID, _ := cmd.Flags().GetString("issue")
		prID, _ := cmd.Flags().GetString("pr")

		base := webBaseURL()
		target := base + "/dashboard"

		if repo != "" {
			target = fmt.Sprintf("%s/dashboard/%s", base, urlPathEscape(repo))
		}
		if issueID != "" && repo != "" {
			target = fmt.Sprintf("%s/dashboard/%s/issues/%s", base, urlPathEscape(repo), urlPathEscape(issueID))
		}
		if prID != "" && repo != "" {
			target = fmt.Sprintf("%s/dashboard/%s/prs/%s", base, urlPathEscape(repo), urlPathEscape(prID))
		}
		if len(args) > 0 && repo == "" {
			target = fmt.Sprintf("%s/dashboard/%s", base, urlPathEscape(args[0]))
		}

		if err := openBrowser(target); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			fmt.Println(target)
			os.Exit(1)
		}
		fmt.Printf("Opening %s\n", target)
	},
}

func urlPathEscape(s string) string {
	return url.PathEscape(s)
}

func openBrowser(url string) error {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url).Run()
	case "linux":
		return exec.Command("xdg-open", url).Run()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Run()
	default:
		return fmt.Errorf("unsupported OS %s — open manually: %s", runtime.GOOS, url)
	}
}

func init() {
	browseCmd.Flags().StringP("repo", "R", "", "Repository ID")
	browseCmd.Flags().String("issue", "", "Issue ID to open")
	browseCmd.Flags().String("pr", "", "Pull request ID to open")
	rootCmd.AddCommand(browseCmd)
}
