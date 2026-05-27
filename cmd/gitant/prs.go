package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/GrayCodeAI/gitant-cli/internal/cli"
	"github.com/spf13/cobra"
)

var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Manage pull requests",
}

var prListCmd = &cobra.Command{
	Use:   "list",
	Short: "List pull requests in a repository",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		status, _ := cmd.Flags().GetString("status")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		path := repoPathSegments(repo, "prs")
		if status != "" {
			path += "?status=" + queryEscape(status)
		}

		var result struct {
			PRs []struct {
				ID           string `json:"id"`
				Title        string `json:"title"`
				Status       string `json:"status"`
				Author       string `json:"author"`
				SourceBranch string `json:"source_branch"`
				TargetBranch string `json:"target_branch"`
			} `json:"prs"`
			Total int `json:"total"`
		}
		if err := client.Get(path, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, pr := range result.PRs {
			fmt.Printf("%s\t%s\t[%s]\t%s -> %s\t%s\n", pr.ID, pr.Status, pr.Author, pr.SourceBranch, pr.TargetBranch, pr.Title)
		}
		fmt.Fprintf(os.Stderr, "%d PR(s)\n", result.Total)
	},
}

var prCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new pull request",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		title, _ := cmd.Flags().GetString("title")
		body, _ := cmd.Flags().GetString("body")
		source, _ := cmd.Flags().GetString("source")
		target, _ := cmd.Flags().GetString("target")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		if title == "" || source == "" {
			fmt.Fprintln(os.Stderr, "Error: --title and --source are required")
			os.Exit(1)
		}

		client := cli.NewClient(daemonURL)
		req := map[string]string{
			"title":         title,
			"body":          body,
			"source_branch": source,
			"target_branch": target,
		}

		var result map[string]interface{}
		if err := client.Post(repoPathSegments(repo, "prs"), req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created PR: %s\n", result["id"])
	},
}

var prViewCmd = &cobra.Command{
	Use:   "view <pr-id>",
	Short: "View a pull request (like gh pr view)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		jsonOut, _ := cmd.Flags().GetBool("json")
		client := newClient(cmd)

		var result map[string]interface{}
		if err := client.Get(repoPathSegments(repo, "prs", args[0]), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if jsonOut {
			printJSON(result)
			return
		}
		fmt.Printf("title:\t%v\n", result["title"])
		fmt.Printf("status:\t%v\n", result["status"])
		fmt.Printf("author:\t%v\n", result["author"])
		fmt.Printf("branches:\t%v -> %v\n", result["source_branch"], result["target_branch"])
	},
}

var prMergeCmd = &cobra.Command{
	Use:   "merge [pr-id]",
	Short: "Merge a pull request",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Post(repoPathSegments(repo, "prs", args[0], "merge"), nil, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Merged PR %s\n", args[0])
	},
}

var prReviewCmd = &cobra.Command{
	Use:   "review [pr-id]",
	Short: "Review a pull request",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		verdict, _ := cmd.Flags().GetString("verdict")
		body, _ := cmd.Flags().GetString("body")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		if verdict == "" {
			fmt.Fprintln(os.Stderr, "Error: --verdict is required (approve|request_changes|comment)")
			os.Exit(1)
		}

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Post(repoPathSegments(repo, "prs", args[0], "review"), map[string]string{"verdict": verdict, "body": body}, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Reviewed PR %s: %s\n", args[0], verdict)
	},
}

var prCommentsCmd = &cobra.Command{
	Use:   "comments [pr-id]",
	Short: "List comments on a pull request",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result struct {
			Comments []struct {
				ID        string `json:"id"`
				Author    string `json:"author"`
				Body      string `json:"body"`
				Timestamp string `json:"timestamp"`
			} `json:"comments"`
			Total int `json:"total"`
		}
		if err := client.Get(repoPathSegments(repo, "prs", args[0], "comments"), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, c := range result.Comments {
			fmt.Printf("%s\t%s\t%s\t%s\n", c.ID, c.Author, c.Timestamp, c.Body)
		}
		fmt.Fprintf(os.Stderr, "%d comment(s)\n", result.Total)
	},
}

var prCheckoutCmd = &cobra.Command{
	Use:   "checkout [pr-id]",
	Short: "Check out a pull request locally",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get(repoPathSegments(repo, "prs", args[0]), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		sourceBranch, _ := result["source_branch"].(string)
		if sourceBranch == "" {
			sourceBranch = fmt.Sprintf("pr-%s", args[0])
		}

		// Create and checkout branch
		checkout := exec.Command("git", "checkout", "-b", sourceBranch)
		checkout.Stdout = os.Stdout
		checkout.Stderr = os.Stderr
		if err := checkout.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error checking out branch: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Checked out PR %s as branch %s\n", args[0], sourceBranch)
	},
}

var prDiffCmd = &cobra.Command{
	Use:   "diff [pr-id]",
	Short: "View the diff of a pull request",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get(repoPathSegments(repo, "prs", args[0]), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		sourceBranch, _ := result["source_branch"].(string)
		targetBranch, _ := result["target_branch"].(string)

		if targetBranch == "" {
			targetBranch = "main"
		}

		// Run git diff
		diff := exec.Command("git", "diff", fmt.Sprintf("%s...%s", targetBranch, sourceBranch))
		diff.Stdout = os.Stdout
		diff.Stderr = os.Stderr
		if err := diff.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error getting diff: %v\n", err)
			os.Exit(1)
		}
	},
}

var prCloseCmd = &cobra.Command{
	Use:   "close [pr-id]",
	Short: "Close a pull request",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Post(repoPathSegments(repo, "prs", args[0], "close"), nil, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Closed PR %s\n", args[0])
	},
}

var prReadyCmd = &cobra.Command{
	Use:   "ready [pr-id]",
	Short: "Mark a pull request as ready for review",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Post(repoPathSegments(repo, "prs", args[0], "ready"), nil, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("PR %s marked as ready for review\n", args[0])
	},
}

func init() {
	for _, c := range []*cobra.Command{prListCmd, prViewCmd, prCreateCmd, prMergeCmd, prReviewCmd, prCommentsCmd, prCheckoutCmd, prDiffCmd, prCloseCmd, prReadyCmd} {
		c.Flags().StringP("repo", "r", "", "Repository name (required)")
		c.MarkFlagRequired("repo")
		c.Flags().String("daemon-url", "", "Daemon URL (default: http://localhost:7777)")
	}
	prViewCmd.Flags().Bool("json", false, "Output JSON")
	prListCmd.Flags().String("status", "", "Filter by status (open|closed|merged)")
	prCreateCmd.Flags().StringP("title", "t", "", "PR title (required)")
	prCreateCmd.Flags().StringP("body", "b", "", "PR body")
	prCreateCmd.Flags().StringP("source", "s", "", "Source branch (required)")
	prCreateCmd.Flags().String("target", "main", "Target branch")
	prReviewCmd.Flags().StringP("verdict", "v", "", "Review verdict: approve|request_changes|comment (required)")
	prReviewCmd.Flags().StringP("body", "b", "", "Review comment")

	prCmd.AddCommand(prListCmd, prViewCmd, prCreateCmd, prMergeCmd, prReviewCmd, prCommentsCmd, prCheckoutCmd, prDiffCmd, prCloseCmd, prReadyCmd)
	rootCmd.AddCommand(prCmd)
}
