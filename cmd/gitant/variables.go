package main

import (
	"fmt"
	"os"

	"github.com/GrayCodeAI/gitant-cli/internal/cli"
	"github.com/spf13/cobra"
)

var variableCmd = &cobra.Command{
	Use:   "variable",
	Short: "Manage CI/CD variables",
	Aliases: []string{"var"},
}

var variableListCmd = &cobra.Command{
	Use:   "list",
	Short: "List variables",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		environment, _ := cmd.Flags().GetString("environment")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		path := fmt.Sprintf("/api/v1/repos/%s/variables", repo)
		if environment != "" {
			path += "?environment=" + environment
		}

		var result struct {
			Variables []struct {
				ID          string `json:"id"`
				Key         string `json:"key"`
				Environment string `json:"environment"`
				Protected   bool   `json:"protected"`
				Masked      bool   `json:"masked"`
			} `json:"variables"`
			Total int `json:"total"`
		}
		if err := client.Get(path, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, v := range result.Variables {
			protected := ""
			if v.Protected {
				protected = "🔒"
			}
			masked := ""
			if v.Masked {
				masked = "****"
			}
			fmt.Printf("%s\t%s\t%s\t%s\t%s\n", v.ID, v.Key, v.Environment, protected, masked)
		}
		fmt.Fprintf(os.Stderr, "%d variable(s)\n", result.Total)
	},
}

var variableCreateCmd = &cobra.Command{
	Use:   "create [key]",
	Short: "Create a variable",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		value, _ := cmd.Flags().GetString("value")
		environment, _ := cmd.Flags().GetString("environment")
		protected, _ := cmd.Flags().GetBool("protected")
		masked, _ := cmd.Flags().GetBool("masked")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		if value == "" {
			value = PromptRequired("Value")
		}

		client := cli.NewClient(daemonURL)
		req := map[string]interface{}{
			"key":         args[0],
			"value":       value,
			"environment": environment,
			"protected":   protected,
			"masked":      masked,
		}

		var result map[string]interface{}
		if err := client.Post(fmt.Sprintf("/api/v1/repos/%s/variables", repo), req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created variable: %s\n", args[0])
	},
}

var variableUpdateCmd = &cobra.Command{
	Use:   "update [key]",
	Short: "Update a variable",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		value, _ := cmd.Flags().GetString("value")
		environment, _ := cmd.Flags().GetString("environment")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		if value == "" {
			value = PromptRequired("New value")
		}

		client := cli.NewClient(daemonURL)
		req := map[string]interface{}{
			"value":       value,
			"environment": environment,
		}

		var result map[string]interface{}
		if err := client.Post(fmt.Sprintf("/api/v1/repos/%s/variables/%s", repo, args[0]), req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Updated variable: %s\n", args[0])
	},
}

var variableDeleteCmd = &cobra.Command{
	Use:   "delete [key]",
	Short: "Delete a variable",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		environment, _ := cmd.Flags().GetString("environment")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		path := fmt.Sprintf("/api/v1/repos/%s/variables/%s", repo, args[0])
		if environment != "" {
			path += "?environment=" + environment
		}

		if err := client.Delete(path); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Deleted variable %s\n", args[0])
	},
}

func init() {
	for _, c := range []*cobra.Command{variableListCmd, variableCreateCmd, variableUpdateCmd, variableDeleteCmd} {
		c.Flags().StringP("repo", "r", "", "Repository name (required)")
		c.MarkFlagRequired("repo")
		c.Flags().String("daemon-url", "", "Daemon URL (default: http://localhost:7777)")
	}
	variableListCmd.Flags().String("environment", "", "Filter by environment")
	variableCreateCmd.Flags().StringP("value", "v", "", "Variable value")
	variableCreateCmd.Flags().String("environment", "", "Environment scope")
	variableCreateCmd.Flags().Bool("protected", false, "Variable is protected")
	variableCreateCmd.Flags().Bool("masked", false, "Variable is masked in logs")
	variableUpdateCmd.Flags().StringP("value", "v", "", "New value")
	variableUpdateCmd.Flags().String("environment", "", "Environment scope")

	variableCmd.AddCommand(variableListCmd, variableCreateCmd, variableUpdateCmd, variableDeleteCmd)
	rootCmd.AddCommand(variableCmd)
}
