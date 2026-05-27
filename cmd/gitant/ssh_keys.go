package main

import (
	"fmt"
	"os"

	"github.com/GrayCodeAI/gitant-cli/internal/cli"
	"github.com/spf13/cobra"
)

var sshKeyCmd = &cobra.Command{
	Use:   "ssh-key",
	Short: "Manage SSH keys",
}

var sshKeyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List SSH keys",
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)

		var result struct {
			Keys []struct {
				ID        string `json:"id"`
				Title     string `json:"title"`
				Key       string `json:"key"`
				CreatedAt string `json:"created_at"`
			} `json:"keys"`
			Total int `json:"total"`
		}
		if err := client.Get("/api/v1/user/keys", &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, k := range result.Keys {
			fmt.Printf("%s\t%s\t%s\n", k.ID, k.Title, k.CreatedAt)
		}
		fmt.Fprintf(os.Stderr, "%d key(s)\n", result.Total)
	},
}

var sshKeyAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an SSH key",
	Run: func(cmd *cobra.Command, args []string) {
		title, _ := cmd.Flags().GetString("title")
		keyFile, _ := cmd.Flags().GetString("key")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		if title == "" {
			title = PromptRequired("Title")
		}

		var key string
		if keyFile != "" {
			data, err := os.ReadFile(keyFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading key file: %v\n", err)
				os.Exit(1)
			}
			key = string(data)
		} else {
			key = PromptRequired("SSH public key")
		}

		client := cli.NewClient(daemonURL)
		req := map[string]string{
			"title": title,
			"key":   key,
		}

		var result map[string]interface{}
		if err := client.Post("/api/v1/user/keys", req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Added SSH key: %s\n", result["id"])
	},
}

var sshKeyDeleteCmd = &cobra.Command{
	Use:   "delete [key-id]",
	Short: "Delete an SSH key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		if err := client.Delete(apiPath("/api/v1/user/keys", args[0])); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Deleted SSH key %s\n", args[0])
	},
}

func init() {
	for _, c := range []*cobra.Command{sshKeyListCmd, sshKeyAddCmd, sshKeyDeleteCmd} {
		c.Flags().String("daemon-url", "", "Daemon URL (default: http://localhost:7777)")
	}
	sshKeyAddCmd.Flags().StringP("title", "t", "", "Key title")
	sshKeyAddCmd.Flags().StringP("key", "k", "", "Path to SSH public key file")

	sshKeyCmd.AddCommand(sshKeyListCmd, sshKeyAddCmd, sshKeyDeleteCmd)
	rootCmd.AddCommand(sshKeyCmd)
}
