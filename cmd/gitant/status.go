package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show daemon status and recent activity (like gh status)",
	Run: func(cmd *cobra.Command, args []string) {
		client := newClient(cmd)

		var daemon map[string]interface{}
		if err := client.Get("/api/v1/status", &daemon); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Gitant status")
		fmt.Println("=============")
		fmt.Printf("Version:   %v\n", daemon["version"])
		fmt.Printf("Uptime:    %v\n", daemon["uptime"])
		fmt.Printf("Repos:     %v\n", daemon["repos"])
		if p2p, ok := daemon["p2p"].(map[string]interface{}); ok && len(p2p) > 0 {
			fmt.Printf("P2P peers: %v\n", p2p["peers"])
		}

		var activity struct {
			Events []struct {
				Type    string `json:"type"`
				Repo    string `json:"repo"`
				Summary string `json:"summary"`
			} `json:"events"`
			Total int `json:"total"`
		}
		if err := client.Get("/api/v1/activity?limit=5", &activity); err == nil && len(activity.Events) > 0 {
			fmt.Println("\nRecent activity:")
			for _, e := range activity.Events {
				fmt.Printf("  [%s] %s — %s\n", e.Repo, e.Type, e.Summary)
			}
		}
	},
}

func init() {
	addDaemonURLFlag(statusCmd)
	rootCmd.AddCommand(statusCmd)
}
