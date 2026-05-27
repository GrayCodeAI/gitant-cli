package main

import (
	"fmt"
	"os"

	"github.com/GrayCodeAI/gitant-cli/internal/cli"
	"github.com/spf13/cobra"
)

var notificationCmd = &cobra.Command{
	Use:   "notification",
	Short: "Manage notifications",
	Aliases: []string{"notif"},
}

var notificationListCmd = &cobra.Command{
	Use:   "list",
	Short: "List notifications",
	Run: func(cmd *cobra.Command, args []string) {
		unreadOnly, _ := cmd.Flags().GetBool("unread")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		path := "/api/v1/notifications"
		if unreadOnly {
			path += "?unread=true"
		}

		var result struct {
			Notifications []struct {
				ID        string `json:"id"`
				Type      string `json:"type"`
				Title     string `json:"title"`
				Body      string `json:"body"`
				Read      bool   `json:"read"`
				CreatedAt string `json:"created_at"`
			} `json:"notifications"`
			UnreadCount int `json:"unread_count"`
		}
		if err := client.Get(path, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, n := range result.Notifications {
			readMark := " "
			if n.Read {
				readMark = "✓"
			}
			fmt.Printf("[%s] %s\t%s\t%s\n", readMark, n.Type, n.Title, n.CreatedAt)
		}
		fmt.Fprintf(os.Stderr, "%d notification(s), %d unread\n", len(result.Notifications), result.UnreadCount)
	},
}

var notificationReadCmd = &cobra.Command{
	Use:   "read [notification-id]",
	Short: "Mark a notification as read",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Post(apiPath("/api/v1/notifications", args[0], "read"), nil, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Marked notification %s as read\n", args[0])
	},
}

var notificationReadAllCmd = &cobra.Command{
	Use:   "read-all",
	Short: "Mark all notifications as read",
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Post("/api/v1/notifications/read-all", nil, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Marked all notifications as read")
	},
}

var notificationCountCmd = &cobra.Command{
	Use:   "count",
	Short: "Show unread notification count",
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result struct {
			Count int `json:"count"`
		}
		if err := client.Get("/api/v1/notifications/unread-count", &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%d unread notification(s)\n", result.Count)
	},
}

func init() {
	for _, c := range []*cobra.Command{notificationListCmd, notificationReadCmd, notificationReadAllCmd, notificationCountCmd} {
		c.Flags().String("daemon-url", "", "Daemon URL (default: http://localhost:7777)")
	}
	notificationListCmd.Flags().Bool("unread", false, "Show only unread notifications")

	notificationCmd.AddCommand(notificationListCmd, notificationReadCmd, notificationReadAllCmd, notificationCountCmd)
	rootCmd.AddCommand(notificationCmd)
}
