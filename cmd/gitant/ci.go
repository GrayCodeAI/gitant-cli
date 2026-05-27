package main

import (
	"fmt"
	"os"

	"github.com/GrayCodeAI/gitant-cli/internal/cli"
	"github.com/spf13/cobra"
)

var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "Work with CI/CD pipelines",
}

var ciListCmd = &cobra.Command{
	Use:   "list",
	Short: "List CI/CD pipelines",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		status, _ := cmd.Flags().GetString("status")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		path := repoPathSegments(repo, "actions", "runs")
		if status != "" {
			path += "?status=" + queryEscape(status)
		}

		var result struct {
			Runs []struct {
				ID        string `json:"id"`
				Status    string `json:"status"`
				Branch    string `json:"branch"`
				CommitSHA string `json:"commit_sha"`
				StartedAt string `json:"started_at"`
			} `json:"runs"`
			Total int `json:"total"`
		}
		if err := client.Get(path, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, run := range result.Runs {
			fmt.Printf("%s\t%s\t%s\t%s\t%s\n", run.ID, run.Status, run.Branch, shortSHA(run.CommitSHA), run.StartedAt)
		}
		fmt.Fprintf(os.Stderr, "%d pipeline(s)\n", result.Total)
	},
}

var ciViewCmd = &cobra.Command{
	Use:   "view [run-id]",
	Short: "View a CI/CD pipeline",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get(repoPathSegments(repo, "actions", "runs", args[0]), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("id:\t%v\n", result["id"])
		fmt.Printf("status:\t%v\n", result["status"])
		fmt.Printf("branch:\t%v\n", result["branch"])
		fmt.Printf("commit:\t%v\n", result["commit_sha"])
		fmt.Printf("started:\t%v\n", result["started_at"])
	},
}

var ciLogsCmd = &cobra.Command{
	Use:   "logs [run-id]",
	Short: "View CI/CD pipeline logs",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get(repoPathSegments(repo, "actions", "runs", args[0]), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		logs, ok := result["logs"].([]interface{})
		if !ok {
			fmt.Println("No logs available")
			return
		}

		for _, line := range logs {
			fmt.Println(line)
		}
	},
}

var ciRetryCmd = &cobra.Command{
	Use:   "retry [run-id]",
	Short: "Retry a CI/CD pipeline",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Post(repoPathSegments(repo, "actions", "runs", args[0], "retry"), nil, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Retried pipeline %s\n", args[0])
	},
}

var ciCancelCmd = &cobra.Command{
	Use:   "cancel [run-id]",
	Short: "Cancel a CI/CD pipeline",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Post(repoPathSegments(repo, "actions", "runs", args[0], "cancel"), nil, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Cancelled pipeline %s\n", args[0])
	},
}

var ciStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show CI/CD status for the current branch",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		branch, _ := cmd.Flags().GetString("branch")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		if branch == "" {
			// Get current branch
			branch = getCurrentBranch()
		}

		client := cli.NewClient(daemonURL)
		var result struct {
			Runs []struct {
				ID     string `json:"id"`
				Status string `json:"status"`
			} `json:"runs"`
		}
		if err := client.Get(repoPathSegments(repo, "actions", "runs")+"?branch="+queryEscape(branch)+"&limit=1", &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if len(result.Runs) == 0 {
			fmt.Printf("No pipelines found for branch %s\n", branch)
			return
		}

		run := result.Runs[0]
		fmt.Printf("Pipeline %s: %s\n", run.ID, run.Status)
	},
}

// shortSHA safely truncates a commit SHA to 8 characters.
func shortSHA(sha string) string {
	if len(sha) <= 8 {
		return sha
	}
	return sha[:8]
}

func getCurrentBranch() string {
	// Try to get current branch from git
	out, err := os.ReadFile(".git/HEAD")
	if err != nil {
		return "main"
	}
	head := string(out)
	if len(head) > 16 && head[:16] == "ref: refs/heads/" {
		return head[16 : len(head)-1]
	}
	return "main"
}

func init() {
	for _, c := range []*cobra.Command{ciListCmd, ciViewCmd, ciLogsCmd, ciRetryCmd, ciCancelCmd, ciStatusCmd} {
		c.Flags().StringP("repo", "r", "", "Repository name (required)")
		c.MarkFlagRequired("repo")
		c.Flags().String("daemon-url", "", "Daemon URL (default: http://localhost:7777)")
	}
	ciListCmd.Flags().String("status", "", "Filter by status (pending|running|success|failed)")
	ciStatusCmd.Flags().String("branch", "", "Branch name (default: current branch)")

	ciCmd.AddCommand(ciListCmd, ciViewCmd, ciLogsCmd, ciRetryCmd, ciCancelCmd, ciStatusCmd)
	rootCmd.AddCommand(ciCmd)
}
