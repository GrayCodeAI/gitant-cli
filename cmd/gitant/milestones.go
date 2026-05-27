package main

import (
	"fmt"
	"os"

	"github.com/GrayCodeAI/gitant-cli/internal/cli"
	"github.com/spf13/cobra"
)

var milestoneCmd = &cobra.Command{
	Use:   "milestone",
	Short: "Manage milestones",
}

var milestoneListCmd = &cobra.Command{
	Use:   "list",
	Short: "List milestones in a repository",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		status, _ := cmd.Flags().GetString("status")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		path := repoPathSegments(repo, "milestones")
		if status != "" {
			path += "?status=" + queryEscape(status)
		}

		var result struct {
			Milestones []struct {
				ID          string `json:"id"`
				Title       string `json:"title"`
				Description string `json:"description"`
				Status      string `json:"status"`
				DueDate     string `json:"due_date"`
				CreatedAt   string `json:"created_at"`
			} `json:"milestones"`
			Total int `json:"total"`
		}
		if err := client.Get(path, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, m := range result.Milestones {
			fmt.Printf("%s\t%s\t%s\t%s\n", m.ID, m.Status, m.Title, m.DueDate)
		}
		fmt.Fprintf(os.Stderr, "%d milestone(s)\n", result.Total)
	},
}

var milestoneCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new milestone",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		dueDate, _ := cmd.Flags().GetString("due-date")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		if title == "" {
			title = PromptRequired("Title")
		}

		client := cli.NewClient(daemonURL)
		req := map[string]string{
			"title":       title,
			"description": description,
			"due_date":    dueDate,
		}

		var result map[string]interface{}
		if err := client.Post(repoPathSegments(repo, "milestones"), req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created milestone: %s\n", result["id"])
	},
}

var milestoneViewCmd = &cobra.Command{
	Use:   "view [milestone-id]",
	Short: "View a milestone",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get(repoPathSegments(repo, "milestones", args[0]), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("title:\t%v\n", result["title"])
		fmt.Printf("status:\t%v\n", result["status"])
		fmt.Printf("description:\t%v\n", result["description"])
		fmt.Printf("due_date:\t%v\n", result["due_date"])
	},
}

var milestoneCloseCmd = &cobra.Command{
	Use:   "close [milestone-id]",
	Short: "Close a milestone",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Post(repoPathSegments(repo, "milestones", args[0], "close"), nil, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Closed milestone %s\n", args[0])
	},
}

func init() {
	for _, c := range []*cobra.Command{milestoneListCmd, milestoneCreateCmd, milestoneViewCmd, milestoneCloseCmd} {
		c.Flags().StringP("repo", "r", "", "Repository name (required)")
		c.MarkFlagRequired("repo")
		c.Flags().String("daemon-url", "", "Daemon URL (default: http://localhost:7777)")
	}
	milestoneListCmd.Flags().String("status", "", "Filter by status (open|closed)")
	milestoneCreateCmd.Flags().StringP("title", "t", "", "Milestone title")
	milestoneCreateCmd.Flags().StringP("description", "d", "", "Milestone description")
	milestoneCreateCmd.Flags().String("due-date", "", "Due date (YYYY-MM-DD)")

	milestoneCmd.AddCommand(milestoneListCmd, milestoneCreateCmd, milestoneViewCmd, milestoneCloseCmd)
	rootCmd.AddCommand(milestoneCmd)
}
