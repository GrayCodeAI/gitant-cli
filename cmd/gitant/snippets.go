package main

import (
	"fmt"
	"os"

	"github.com/GrayCodeAI/gitant-cli/internal/cli"
	"github.com/spf13/cobra"
)

var snippetCmd = &cobra.Command{
	Use:   "snippet",
	Short: "Manage code snippets",
	Aliases: []string{"snip"},
}

var snippetListCmd = &cobra.Command{
	Use:   "list",
	Short: "List snippets",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		path := repoPathSegments(repo, "snippets")

		var result struct {
			Snippets []struct {
				ID          string `json:"id"`
				Title       string `json:"title"`
				Description string `json:"description"`
				Author      string `json:"author"`
				CreatedAt   string `json:"created_at"`
			} `json:"snippets"`
			Total int `json:"total"`
		}
		if err := client.Get(path, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, s := range result.Snippets {
			fmt.Printf("%s\t%s\t%s\t%s\n", s.ID, s.Title, s.Author, s.CreatedAt)
		}
		fmt.Fprintf(os.Stderr, "%d snippet(s)\n", result.Total)
	},
}

var snippetCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new snippet",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		filePath, _ := cmd.Flags().GetString("file")
		visibility, _ := cmd.Flags().GetString("visibility")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		if title == "" {
			title = PromptRequired("Title")
		}

		var content string
		if filePath != "" {
			data, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
				os.Exit(1)
			}
			content = string(data)
		} else {
			content = PromptMultiline("Content")
		}

		client := cli.NewClient(daemonURL)
		req := map[string]string{
			"title":       title,
			"description": description,
			"content":     content,
			"visibility":  visibility,
		}

		var result map[string]interface{}
		if err := client.Post(repoPathSegments(repo, "snippets"), req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created snippet: %s\n", result["id"])
	},
}

var snippetViewCmd = &cobra.Command{
	Use:   "view [snippet-id]",
	Short: "View a snippet",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get(repoPathSegments(repo, "snippets", args[0]), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("title:\t%v\n", result["title"])
		fmt.Printf("description:\t%v\n", result["description"])
		fmt.Printf("author:\t%v\n", result["author"])
		fmt.Printf("visibility:\t%v\n", result["visibility"])
		fmt.Printf("\n---\n%s\n", result["content"])
	},
}

var snippetDeleteCmd = &cobra.Command{
	Use:   "delete [snippet-id]",
	Short: "Delete a snippet",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		if err := client.Delete(repoPathSegments(repo, "snippets", args[0])); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Deleted snippet %s\n", args[0])
	},
}

func init() {
	for _, c := range []*cobra.Command{snippetListCmd, snippetCreateCmd, snippetViewCmd, snippetDeleteCmd} {
		c.Flags().StringP("repo", "r", "", "Repository name (required)")
		c.MarkFlagRequired("repo")
		c.Flags().String("daemon-url", "", "Daemon URL (default: http://localhost:7777)")
	}
	snippetCreateCmd.Flags().StringP("title", "t", "", "Snippet title")
	snippetCreateCmd.Flags().StringP("description", "d", "", "Snippet description")
	snippetCreateCmd.Flags().StringP("file", "f", "", "File to read content from")
	snippetCreateCmd.Flags().String("visibility", "private", "Visibility (private|internal|public)")

	snippetCmd.AddCommand(snippetListCmd, snippetCreateCmd, snippetViewCmd, snippetDeleteCmd)
	rootCmd.AddCommand(snippetCmd)
}
