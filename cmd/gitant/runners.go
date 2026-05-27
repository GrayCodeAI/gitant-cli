package main

import (
	"fmt"
	"os"

	"github.com/GrayCodeAI/gitant-cli/internal/cli"
	"github.com/spf13/cobra"
)

var runnerCmd = &cobra.Command{
	Use:   "runner",
	Short: "Manage CI/CD runners",
}

var runnerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List runners",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		status, _ := cmd.Flags().GetString("status")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		path := repoPathSegments(repo, "runners")
		if status != "" {
			path += "?status=" + queryEscape(status)
		}

		var result struct {
			Runners []struct {
				ID     string `json:"id"`
				Name   string `json:"name"`
				Status string `json:"status"`
				Type   string `json:"type"`
				Tags   string `json:"tags"`
			} `json:"runners"`
			Total int `json:"total"`
		}
		if err := client.Get(path, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, r := range result.Runners {
			fmt.Printf("%s\t%s\t%s\t%s\t%s\n", r.ID, r.Name, r.Status, r.Type, r.Tags)
		}
		fmt.Fprintf(os.Stderr, "%d runner(s)\n", result.Total)
	},
}

var runnerRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new runner",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		name, _ := cmd.Flags().GetString("name")
		tags, _ := cmd.Flags().GetString("tags")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		if name == "" {
			name = PromptRequired("Runner name")
		}

		client := cli.NewClient(daemonURL)
		req := map[string]string{
			"name": name,
			"tags": tags,
		}

		var result map[string]interface{}
		if err := client.Post(repoPathSegments(repo, "runners"), req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Registered runner: %s\n", result["id"])
		fmt.Printf("Token: %s\n", result["token"])
	},
}

var runnerUnregisterCmd = &cobra.Command{
	Use:   "unregister [runner-id]",
	Short: "Unregister a runner",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		if err := client.Delete(repoPathSegments(repo, "runners", args[0])); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Unregistered runner %s\n", args[0])
	},
}

var runnerStatusCmd = &cobra.Command{
	Use:   "status [runner-id]",
	Short: "View runner status",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get(repoPathSegments(repo, "runners", args[0]), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("id:\t%v\n", result["id"])
		fmt.Printf("name:\t%v\n", result["name"])
		fmt.Printf("status:\t%v\n", result["status"])
		fmt.Printf("type:\t%v\n", result["type"])
		fmt.Printf("tags:\t%v\n", result["tags"])
	},
}

func init() {
	for _, c := range []*cobra.Command{runnerListCmd, runnerRegisterCmd, runnerUnregisterCmd, runnerStatusCmd} {
		c.Flags().StringP("repo", "r", "", "Repository name (required)")
		c.MarkFlagRequired("repo")
		c.Flags().String("daemon-url", "", "Daemon URL (default: http://localhost:7777)")
	}
	runnerListCmd.Flags().String("status", "", "Filter by status (online|offline|idle|active)")
	runnerRegisterCmd.Flags().String("name", "", "Runner name")
	runnerRegisterCmd.Flags().String("tags", "", "Runner tags (comma-separated)")

	runnerCmd.AddCommand(runnerListCmd, runnerRegisterCmd, runnerUnregisterCmd, runnerStatusCmd)
	rootCmd.AddCommand(runnerCmd)
}
