package main

import (
	"fmt"
	"os"

	"github.com/GrayCodeAI/gitant-cli/internal/cli"
	"github.com/spf13/cobra"
)

var environmentCmd = &cobra.Command{
	Use:   "environment",
	Short: "Manage environments",
	Aliases: []string{"env"},
}

var environmentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List environments",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result struct {
			Environments []struct {
				ID        string `json:"id"`
				Name      string `json:"name"`
				Status    string `json:"status"`
				UpdatedAt string `json:"updated_at"`
			} `json:"environments"`
			Total int `json:"total"`
		}
		if err := client.Get(fmt.Sprintf("/api/v1/repos/%s/environments", repo), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, e := range result.Environments {
			fmt.Printf("%s\t%s\t%s\t%s\n", e.ID, e.Name, e.Status, e.UpdatedAt)
		}
		fmt.Fprintf(os.Stderr, "%d environment(s)\n", result.Total)
	},
}

var environmentCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create an environment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		req := map[string]string{
			"name": args[0],
		}

		var result map[string]interface{}
		if err := client.Post(fmt.Sprintf("/api/v1/repos/%s/environments", repo), req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created environment: %s\n", args[0])
	},
}

var environmentDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete an environment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		if err := client.Delete(fmt.Sprintf("/api/v1/repos/%s/environments/%s", repo, args[0])); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Deleted environment %s\n", args[0])
	},
}

func init() {
	for _, c := range []*cobra.Command{environmentListCmd, environmentCreateCmd, environmentDeleteCmd} {
		c.Flags().StringP("repo", "r", "", "Repository name (required)")
		c.MarkFlagRequired("repo")
		c.Flags().String("daemon-url", "", "Daemon URL (default: http://localhost:7777)")
	}

	environmentCmd.AddCommand(environmentListCmd, environmentCreateCmd, environmentDeleteCmd)
	rootCmd.AddCommand(environmentCmd)
}
