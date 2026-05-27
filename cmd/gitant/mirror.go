package main

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/GrayCodeAI/gitant-cli/internal/cli"
	"github.com/spf13/cobra"
)

var mirrorCmd = &cobra.Command{
	Use:   "mirror",
	Short: "Mirror a repository from GitHub/GitLab",
	Long:  "Import a public repository from GitHub, GitLab, or any git URL into your Gitant node.",
}

var mirrorImportCmd = &cobra.Command{
	Use:   "import <source-url>",
	Short: "Import a repository from a git URL",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		private, _ := cmd.Flags().GetBool("private")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		sourceURL := args[0]

		if name == "" {
			// Extract name from URL
			name = extractRepoName(sourceURL)
		}

		client := cli.NewClient(daemonURL)
		req := map[string]interface{}{
			"name":        name,
			"description": description,
			"private":     private,
			"source_url":  sourceURL,
		}

		var result map[string]interface{}
		if err := client.Post("/api/v1/repos/mirror", req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Mirrored %s → %s\n", sourceURL, name)
		fmt.Printf("  Repo ID: %s\n", result["id"])
	},
}

var mirrorGithubCmd = &cobra.Command{
	Use:   "github <owner>/<repo>",
	Short: "Import a GitHub repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		private, _ := cmd.Flags().GetBool("private")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		escapedParts := strings.SplitN(args[0], "/", 2)
		for i := range escapedParts {
			escapedParts[i] = url.PathEscape(escapedParts[i])
		}
		sourceURL := "https://github.com/" + strings.Join(escapedParts, "/") + ".git"

		if name == "" {
			name = extractRepoName(sourceURL)
		}

		client := cli.NewClient(daemonURL)
		req := map[string]interface{}{
			"name":        name,
			"description": fmt.Sprintf("Mirror of %s", args[0]),
			"private":     private,
			"source_url":  sourceURL,
		}

		var result map[string]interface{}
		if err := client.Post("/api/v1/repos/mirror", req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Mirrored github.com/%s → %s\n", args[0], name)
	},
}

var mirrorGitlabCmd = &cobra.Command{
	Use:   "gitlab <owner>/<repo>",
	Short: "Import a GitLab repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		private, _ := cmd.Flags().GetBool("private")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		escapedParts := strings.SplitN(args[0], "/", 2)
		for i := range escapedParts {
			escapedParts[i] = url.PathEscape(escapedParts[i])
		}
		sourceURL := "https://gitlab.com/" + strings.Join(escapedParts, "/") + ".git"

		if name == "" {
			name = extractRepoName(sourceURL)
		}

		client := cli.NewClient(daemonURL)
		req := map[string]interface{}{
			"name":        name,
			"description": fmt.Sprintf("Mirror of %s", args[0]),
			"private":     private,
			"source_url":  sourceURL,
		}

		var result map[string]interface{}
		if err := client.Post("/api/v1/repos/mirror", req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Mirrored gitlab.com/%s → %s\n", args[0], name)
	},
}

func extractRepoName(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err == nil && u.Path != "" {
		// Strip trailing slash, then take the last path segment.
		base := path.Base(strings.TrimSuffix(u.Path, "/"))
		// Remove trailing .git
		if strings.HasSuffix(base, ".git") {
			base = strings.TrimSuffix(base, ".git")
		}
		if base != "" && base != "." {
			return base
		}
	}

	// Fallback: best-effort on opaque strings.
	s := rawURL
	if idx := strings.Index(s, "?"); idx != -1 {
		s = s[:idx]
	}
	if idx := strings.Index(s, "#"); idx != -1 {
		s = s[:idx]
	}
	if strings.HasSuffix(s, ".git") {
		s = strings.TrimSuffix(s, ".git")
	}
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '/' {
			return s[i+1:]
		}
	}
	return s
}

func init() {
	for _, c := range []*cobra.Command{mirrorImportCmd, mirrorGithubCmd, mirrorGitlabCmd} {
		c.Flags().String("name", "", "Repository name (default: extracted from URL)")
		c.Flags().StringP("description", "d", "", "Repository description")
		c.Flags().Bool("private", false, "Create as private repository")
		c.Flags().String("daemon-url", "", "Daemon URL (default: http://localhost:7777)")
	}

	mirrorCmd.AddCommand(mirrorImportCmd, mirrorGithubCmd, mirrorGitlabCmd)
	rootCmd.AddCommand(mirrorCmd)
}
