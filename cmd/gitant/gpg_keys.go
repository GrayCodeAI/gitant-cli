package main

import (
	"fmt"
	"os"

	"github.com/GrayCodeAI/gitant-cli/internal/cli"
	"github.com/spf13/cobra"
)

var gpgKeyCmd = &cobra.Command{
	Use:   "gpg-key",
	Short: "Manage GPG keys",
}

var gpgKeyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List GPG keys",
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)

		var result struct {
			Keys []struct {
				ID        string `json:"id"`
				KeyID     string `json:"key_id"`
				CreatedAt string `json:"created_at"`
			} `json:"keys"`
			Total int `json:"total"`
		}
		if err := client.Get("/api/v1/user/gpg-keys", &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, k := range result.Keys {
			fmt.Printf("%s\t%s\t%s\n", k.ID, k.KeyID, k.CreatedAt)
		}
		fmt.Fprintf(os.Stderr, "%d key(s)\n", result.Total)
	},
}

var gpgKeyAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a GPG key",
	Run: func(cmd *cobra.Command, args []string) {
		keyFile, _ := cmd.Flags().GetString("key")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		var key string
		if keyFile != "" {
			data, err := os.ReadFile(keyFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading key file: %v\n", err)
				os.Exit(1)
			}
			key = string(data)
		} else {
			key = PromptRequired("GPG public key")
		}

		client := cli.NewClient(daemonURL)
		req := map[string]string{
			"key": key,
		}

		var result map[string]interface{}
		if err := client.Post("/api/v1/user/gpg-keys", req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Added GPG key: %s\n", result["id"])
	},
}

var gpgKeyDeleteCmd = &cobra.Command{
	Use:   "delete [key-id]",
	Short: "Delete a GPG key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		if err := client.Delete(apiPath("/api/v1/user/gpg-keys", args[0])); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Deleted GPG key %s\n", args[0])
	},
}

func init() {
	for _, c := range []*cobra.Command{gpgKeyListCmd, gpgKeyAddCmd, gpgKeyDeleteCmd} {
		c.Flags().String("daemon-url", "", "Daemon URL (default: http://localhost:7777)")
	}
	gpgKeyAddCmd.Flags().StringP("key", "k", "", "Path to GPG public key file")

	gpgKeyCmd.AddCommand(gpgKeyListCmd, gpgKeyAddCmd, gpgKeyDeleteCmd)
	rootCmd.AddCommand(gpgKeyCmd)
}
