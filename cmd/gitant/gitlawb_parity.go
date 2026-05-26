package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/GrayCodeAI/gitant-cli/internal/cli"
	"github.com/spf13/cobra"
)

// gt identity — manage DID identity (like gl identity)

var identityCmd = &cobra.Command{
	Use:   "identity",
	Short: "Manage your DID identity (like gl identity)",
}

var identityNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Generate a new Ed25519 keypair and DID",
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Post("/api/v1/identity/generate", nil, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Identity generated!")
		if did, ok := result["did"].(string); ok {
			fmt.Printf("DID: %s\n", did)
		}
		if path, ok := result["key_path"].(string); ok {
			fmt.Printf("Key saved to: %s\n", path)
		}
	},
}

var identityShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show your current DID and identity info",
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get("/api/v1/identity", &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if did, ok := result["did"].(string); ok {
			fmt.Println(did)
		}
	},
}

var identityExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export your DID document as JSON",
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get("/api/v1/identity/export", &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		out, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(out))
	},
}

var identitySignCmd = &cobra.Command{
	Use:   "sign <message>",
	Short: "Sign a message with your Ed25519 private key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		req := map[string]string{"message": args[0]}
		var result map[string]interface{}
		if err := client.Post("/api/v1/identity/sign", req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if sig, ok := result["signature"].(string); ok {
			fmt.Println(sig)
		}
	},
}

// gt node — node status and network info

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Node status and network info",
}

var nodeStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show node status and connectivity",
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get("/api/v1/status", &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		out, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(out))
	},
}

var nodeTrustCmd = &cobra.Command{
	Use:   "trust <did>",
	Short: "Show trust score and attestation details for an agent",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get(fmt.Sprintf("/api/v1/agents/%s/trust", args[0]), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("DID:\t%s\n", args[0])
		fmt.Printf("Score:\t%v\n", result["score"])
		if attestations, ok := result["attestations"].(float64); ok {
			fmt.Printf("Attestations:\t%.0f\n", attestations)
		}
		if breakdown, ok := result["breakdown"].(map[string]interface{}); ok {
			fmt.Println("\nBreakdown:")
			for k, v := range breakdown {
				fmt.Printf("  %s:\t%v\n", k, v)
			}
		}
	},
}

var nodeResolveCmd = &cobra.Command{
	Use:   "resolve <did>",
	Short: "Resolve a DID to its document via the DHT",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get(fmt.Sprintf("/api/v1/agents/resolve/%s", args[0]), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		out, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(out))
	},
}

// gt peer — peer discovery

var peerCmd = &cobra.Command{
	Use:   "peer",
	Short: "Peer discovery — add and inspect known nodes",
}

var peerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List known peers",
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result struct {
			Peers []struct {
				ID        string   `json:"id"`
				Multiaddrs []string `json:"multiaddrs"`
				Connected bool     `json:"connected"`
			} `json:"peers"`
			Total int `json:"total"`
		}
		if err := client.Get("/api/v1/network/peers", &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, p := range result.Peers {
			status := "disconnected"
			if p.Connected {
				status = "connected"
			}
			fmt.Printf("%s\t%s\t%s\n", p.ID[:12], status, p.Multiaddrs)
		}
		fmt.Fprintf(os.Stderr, "%d peer(s)\n", result.Total)
	},
}

var peerAddCmd = &cobra.Command{
	Use:   "add <multiaddr>",
	Short: "Add a peer by multiaddr",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		req := map[string]string{"multiaddr": args[0]}
		var result map[string]interface{}
		if err := client.Post("/api/v1/network/peers", req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Added peer %s\n", args[0])
	},
}

// Bounty subcommands (approve, cancel, stats)

var bountyApproveCmd = &cobra.Command{
	Use:   "approve <bounty-id>",
	Short: "Approve a bounty submission and release payment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Post(fmt.Sprintf("/api/v1/repos/%s/bounties/%s/approve", repo, args[0]), nil, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Approved bounty %s — payment released\n", args[0])
	},
}

var bountyCancelCmd = &cobra.Command{
	Use:   "cancel <bounty-id>",
	Short: "Cancel a bounty and refund the escrow",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Post(fmt.Sprintf("/api/v1/repos/%s/bounties/%s/cancel", repo, args[0]), nil, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Cancelled bounty %s — escrow refunded\n", args[0])
	},
}

var bountyStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show bounty statistics for a repository",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get(fmt.Sprintf("/api/v1/repos/%s/bounties/stats", repo), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		out, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(out))
	},
}

// Task subcommand (fail)

var taskFailCmd = &cobra.Command{
	Use:   "fail <task-id>",
	Short: "Mark a task as failed",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Post(fmt.Sprintf("/api/v1/repos/%s/tasks/%s/fail", repo, args[0]), nil, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Task %s marked as failed\n", args[0])
	},
}

// Cert verify

var certVerifyCmd = &cobra.Command{
	Use:   "verify <cert-id>",
	Short: "Verify a ref-update certificate's signature",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get(fmt.Sprintf("/api/v1/repos/%s/certs/%s/verify", repo, args[0]), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if valid, ok := result["valid"].(bool); ok && valid {
			fmt.Printf("Certificate %s is VALID\n", args[0])
		} else {
			fmt.Printf("Certificate %s is INVALID\n", args[0])
		}
		if reason, ok := result["reason"].(string); ok && reason != "" {
			fmt.Printf("Reason: %s\n", reason)
		}
	},
}

func init() {
	// Identity flags
	for _, c := range []*cobra.Command{identityNewCmd, identityShowCmd, identityExportCmd, identitySignCmd} {
		c.Flags().String("daemon-url", "", "Daemon URL")
	}

	// Node flags
	for _, c := range []*cobra.Command{nodeStatusCmd, nodeTrustCmd, nodeResolveCmd} {
		c.Flags().String("daemon-url", "", "Daemon URL")
	}

	// Peer flags
	for _, c := range []*cobra.Command{peerListCmd, peerAddCmd} {
		c.Flags().String("daemon-url", "", "Daemon URL")
	}

	// Bounty flags
	for _, c := range []*cobra.Command{bountyApproveCmd, bountyCancelCmd, bountyStatsCmd} {
		c.Flags().StringP("repo", "r", "", "Repository name (required)")
		c.MarkFlagRequired("repo")
		c.Flags().String("daemon-url", "", "Daemon URL")
	}

	// Task flags
	taskFailCmd.Flags().StringP("repo", "r", "", "Repository name (required)")
	taskFailCmd.MarkFlagRequired("repo")
	taskFailCmd.Flags().String("daemon-url", "", "Daemon URL")

	// Cert flags
	certVerifyCmd.Flags().StringP("repo", "r", "", "Repository name (required)")
	certVerifyCmd.MarkFlagRequired("repo")
	certVerifyCmd.Flags().String("daemon-url", "", "Daemon URL")

	// Register subcommands
	identityCmd.AddCommand(identityNewCmd, identityShowCmd, identityExportCmd, identitySignCmd)
	nodeCmd.AddCommand(nodeStatusCmd, nodeTrustCmd, nodeResolveCmd)
	peerCmd.AddCommand(peerListCmd, peerAddCmd)

	// Add to existing commands
	bountyCmd.AddCommand(bountyApproveCmd, bountyCancelCmd, bountyStatsCmd)
	taskCmd.AddCommand(taskFailCmd)
	certCmd.AddCommand(certVerifyCmd)

	// Add new top-level commands
	rootCmd.AddCommand(identityCmd, nodeCmd, peerCmd)
}
