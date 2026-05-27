package main

import (
	"fmt"
	"os"

	"github.com/GrayCodeAI/gitant-cli/internal/cli"
	"github.com/spf13/cobra"
)

var deploymentCmd = &cobra.Command{
	Use:   "deployment",
	Short: "Manage deployments",
	Aliases: []string{"deploy"},
}

var deploymentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List deployments",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		environment, _ := cmd.Flags().GetString("environment")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		path := repoPathSegments(repo, "deployments")
		if environment != "" {
			path += "?environment=" + queryEscape(environment)
		}

		var result struct {
			Deployments []struct {
				ID          string `json:"id"`
				Environment string `json:"environment"`
				Ref         string `json:"ref"`
				Status      string `json:"status"`
				Deployer    string `json:"deployer"`
				CreatedAt   string `json:"created_at"`
			} `json:"deployments"`
			Total int `json:"total"`
		}
		if err := client.Get(path, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, d := range result.Deployments {
			fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\n", d.ID, d.Environment, d.Ref, d.Status, d.Deployer, d.CreatedAt)
		}
		fmt.Fprintf(os.Stderr, "%d deployment(s)\n", result.Total)
	},
}

var deploymentCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a deployment",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		environment, _ := cmd.Flags().GetString("environment")
		ref, _ := cmd.Flags().GetString("ref")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		if environment == "" {
			environment = PromptRequired("Environment")
		}
		if ref == "" {
			ref = PromptRequired("Ref (branch/tag/SHA)")
		}

		client := cli.NewClient(daemonURL)
		req := map[string]string{
			"environment": environment,
			"ref":         ref,
		}

		var result map[string]interface{}
		if err := client.Post(repoPathSegments(repo, "deployments"), req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created deployment: %s\n", result["id"])
	},
}

var deploymentStatusCmd = &cobra.Command{
	Use:   "status [deployment-id]",
	Short: "View deployment status",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get(repoPathSegments(repo, "deployments", args[0]), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("id:\t%v\n", result["id"])
		fmt.Printf("env:\t%v\n", result["environment"])
		fmt.Printf("ref:\t%v\n", result["ref"])
		fmt.Printf("status:\t%v\n", result["status"])
		fmt.Printf("deployer:\t%v\n", result["deployer"])
	},
}

var deploymentRollbackCmd = &cobra.Command{
	Use:   "rollback [deployment-id]",
	Short: "Rollback a deployment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Post(repoPathSegments(repo, "deployments", args[0], "rollback"), nil, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Rolled back deployment %s\n", args[0])
	},
}

func init() {
	for _, c := range []*cobra.Command{deploymentListCmd, deploymentCreateCmd, deploymentStatusCmd, deploymentRollbackCmd} {
		c.Flags().StringP("repo", "r", "", "Repository name (required)")
		c.MarkFlagRequired("repo")
		c.Flags().String("daemon-url", "", "Daemon URL (default: http://localhost:7777)")
	}
	deploymentListCmd.Flags().String("environment", "", "Filter by environment")
	deploymentCreateCmd.Flags().String("environment", "", "Environment name")
	deploymentCreateCmd.Flags().String("ref", "", "Git ref (branch/tag/SHA)")

	deploymentCmd.AddCommand(deploymentListCmd, deploymentCreateCmd, deploymentStatusCmd, deploymentRollbackCmd)
	rootCmd.AddCommand(deploymentCmd)
}
