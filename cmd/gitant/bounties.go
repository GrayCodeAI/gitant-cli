package main

import (
	"fmt"
	"os"

	"github.com/GrayCodeAI/gitant-cli/internal/cli"
	"github.com/spf13/cobra"
)

var bountyCmd = &cobra.Command{
	Use:   "bounty",
	Short: "Manage bounties",
}

var bountyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List bounties",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		status, _ := cmd.Flags().GetString("status")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		path := repoPathSegments(repo, "bounties")
		if status != "" {
			path += "?status=" + queryEscape(status)
		}

		var result struct {
			Bounties []struct {
				ID       string  `json:"id"`
				Title    string  `json:"title"`
				Amount   float64 `json:"amount"`
				Currency string  `json:"currency"`
				Status   string  `json:"status"`
				Creator  string  `json:"creator"`
			} `json:"bounties"`
			Total int `json:"total"`
		}
		if err := client.Get(path, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, b := range result.Bounties {
			fmt.Printf("%s\t%s\t%.2f %s\t%s\t%s\n", b.ID, b.Title, b.Amount, b.Currency, b.Status, b.Creator)
		}
		fmt.Fprintf(os.Stderr, "%d bounty(ies)\n", result.Total)
	},
}

var bountyCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a bounty",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		amount, _ := cmd.Flags().GetFloat64("amount")
		currency, _ := cmd.Flags().GetString("currency")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		if title == "" {
			title = PromptRequired("Title")
		}
		if amount == 0 {
			fmt.Print("Amount: ")
			fmt.Scan(&amount)
		}

		client := cli.NewClient(daemonURL)
		req := map[string]interface{}{
			"title":       title,
			"description": description,
			"amount":      amount,
			"currency":    currency,
		}

		var result map[string]interface{}
		if err := client.Post(repoPathSegments(repo, "bounties"), req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created bounty: %s\n", result["id"])
	},
}

var bountyClaimCmd = &cobra.Command{
	Use:   "claim [bounty-id]",
	Short: "Claim a bounty",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Post(apiPath("/api/v1/bounties", args[0], "claim"), nil, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Claimed bounty %s\n", args[0])
	},
}

var bountySubmitCmd = &cobra.Command{
	Use:   "submit [bounty-id]",
	Short: "Submit work for a bounty",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		submission, _ := cmd.Flags().GetString("submission")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		if submission == "" {
			submission = PromptRequired("Submission (URL or description)")
		}

		client := cli.NewClient(daemonURL)
		req := map[string]string{
			"submission": submission,
		}

		var result map[string]interface{}
		if err := client.Post(apiPath("/api/v1/bounties", args[0], "submit"), req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Submitted work for bounty %s\n", args[0])
	},
}

var bountyViewCmd = &cobra.Command{
	Use:   "view [bounty-id]",
	Short: "View a bounty",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get(apiPath("/api/v1/bounties", args[0]), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("title:\t%v\n", result["title"])
		fmt.Printf("amount:\t%v %v\n", result["amount"], result["currency"])
		fmt.Printf("status:\t%v\n", result["status"])
		fmt.Printf("creator:\t%v\n", result["creator"])
		fmt.Printf("claimed_by:\t%v\n", result["claimed_by"])
	},
}

func init() {
	for _, c := range []*cobra.Command{bountyListCmd, bountyCreateCmd, bountyClaimCmd, bountySubmitCmd, bountyViewCmd} {
		c.Flags().String("daemon-url", "", "Daemon URL (default: http://localhost:7777)")
	}
	bountyListCmd.Flags().StringP("repo", "r", "", "Repository name")
	bountyListCmd.Flags().String("status", "", "Filter by status (open|claimed|submitted|approved)")
	bountyCreateCmd.Flags().StringP("repo", "r", "", "Repository name (required)")
	bountyCreateCmd.MarkFlagRequired("repo")
	bountyCreateCmd.Flags().StringP("title", "t", "", "Bounty title")
	bountyCreateCmd.Flags().StringP("description", "d", "", "Bounty description")
	bountyCreateCmd.Flags().Float64P("amount", "a", 0, "Bounty amount")
	bountyCreateCmd.Flags().String("currency", "USD", "Currency")
	bountySubmitCmd.Flags().StringP("submission", "s", "", "Submission URL or description")

	bountyCmd.AddCommand(bountyListCmd, bountyCreateCmd, bountyClaimCmd, bountySubmitCmd, bountyViewCmd)
	rootCmd.AddCommand(bountyCmd)
}
