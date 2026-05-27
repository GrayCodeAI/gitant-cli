package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/GrayCodeAI/gitant-cli/internal/cli"
	"github.com/spf13/cobra"
)

// gl cert — inspect signed ref-update certificates

var certCmd = &cobra.Command{
	Use:   "cert",
	Short: "Inspect signed ref-update certificates",
}

var certListCmd = &cobra.Command{
	Use:   "list",
	Short: "List ref certificates for a repository",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result struct {
			Certificates []struct {
				ID        string `json:"id"`
				Ref       string `json:"ref"`
				OldOID    string `json:"old_oid"`
				NewOID    string `json:"new_oid"`
				Signer    string `json:"signer"`
				Signature string `json:"signature"`
				Timestamp string `json:"timestamp"`
			} `json:"certificates"`
			Total int `json:"total"`
		}
		if err := client.Get(repoPathSegments(repo, "certs"), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, c := range result.Certificates {
			fmt.Printf("%s\t%s\t%s→%s\t%s\t%s\n", c.ID, c.Ref, c.OldOID[:8], c.NewOID[:8], c.Signer, c.Timestamp)
		}
		fmt.Fprintf(os.Stderr, "%d certificate(s)\n", result.Total)
	},
}

var certShowCmd = &cobra.Command{
	Use:   "show <cert-id>",
	Short: "Show a specific ref certificate and verify its signature",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get(repoPathSegments(repo, "certs", args[0]), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		out, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(out))
	},
}

// gl ipfs — IPFS pin management

var ipfsCmd = &cobra.Command{
	Use:   "ipfs",
	Short: "IPFS pin management — list pinned CIDs and retrieve objects",
}

var ipfsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all CIDs pinned to the node's local IPFS daemon",
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result struct {
			Pins []struct {
				CID       string `json:"cid"`
				Type      string `json:"type"`
				Size      int64  `json:"size"`
				PinnedAt  string `json:"pinned_at"`
			} `json:"pins"`
			Total int `json:"total"`
		}
		if err := client.Get("/api/v1/ipfs/pins", &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, p := range result.Pins {
			fmt.Printf("%s\t%s\t%d bytes\t%s\n", p.CID, p.Type, p.Size, p.PinnedAt)
		}
		fmt.Fprintf(os.Stderr, "%d pin(s)\n", result.Total)
	},
}

var ipfsGetCmd = &cobra.Command{
	Use:   "get <cid>",
	Short: "Retrieve and display a git object from the node by its CIDv1",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get(apiPath("/api/v1/ipfs", args[0]), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		out, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(out))
	},
}

// gl mcp — MCP server for LLM agents

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "MCP server — expose gitant tools to LLM agents",
}

var mcpServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the MCP server (stdin/stdout JSON-RPC)",
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")
		fmt.Fprintf(os.Stderr, "Starting MCP server (daemon: %s)...\n", daemonURL)
		fmt.Fprintf(os.Stderr, "MCP server not yet implemented — use gitant-mcp instead\n")
		fmt.Fprintf(os.Stderr, "See: https://github.com/GrayCodeAI/gitant-mcp\n")
	},
}

// gl sync — sync repos from peer nodes

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync repos from peer nodes (HTTP fallback for p2p gossip)",
}

var syncTriggerCmd = &cobra.Command{
	Use:   "trigger",
	Short: "Pull repos from all known peers into the sync queue",
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Post("/api/v1/sync/trigger", nil, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Sync triggered")
		if queued, ok := result["queued"].(float64); ok {
			fmt.Printf("%.0f repos queued for sync\n", queued)
		}
	},
}

var syncStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the current sync queue status",
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get("/api/v1/sync/status", &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		out, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(out))
	},
}

// gl name — register and resolve names on Base L2

var nameCmd = &cobra.Command{
	Use:   "name",
	Short: "Register and resolve names on Base L2",
}

var nameRegisterCmd = &cobra.Command{
	Use:   "register <name>",
	Short: "Register a name → your DID on Base L2",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		req := map[string]string{"name": args[0]}
		var result map[string]interface{}
		if err := client.Post("/api/v1/names/register", req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Registered %s\n", args[0])
		if tx, ok := result["tx_hash"].(string); ok {
			fmt.Printf("TX: %s\n", tx)
		}
	},
}

var nameResolveCmd = &cobra.Command{
	Use:   "resolve <name>",
	Short: "Resolve name → owner address + DID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get(apiPath("/api/v1/names", args[0], "resolve"), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("name:\t%s\n", args[0])
		fmt.Printf("owner:\t%v\n", result["owner"])
		fmt.Printf("did:\t%v\n", result["did"])
	},
}

var nameLookupCmd = &cobra.Command{
	Use:   "lookup <did>",
	Short: "Reverse lookup: DID → registered name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get(apiPath("/api/v1/names/lookup")+"?did="+queryEscape(args[0]), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if name, ok := result["name"].(string); ok && name != "" {
			fmt.Println(name)
		} else {
			fmt.Println("No name found for this DID")
		}
	},
}

var nameAvailableCmd = &cobra.Command{
	Use:   "available <name>",
	Short: "Check whether a name is available to register",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get(apiPath("/api/v1/names", args[0], "available"), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if avail, ok := result["available"].(bool); ok && avail {
			fmt.Printf("%s is available\n", args[0])
		} else {
			fmt.Printf("%s is not available\n", args[0])
		}
	},
}

var nameRegisterDidCmd = &cobra.Command{
	Use:   "register-did",
	Short: "Anchor your DID document in the on-chain DID registry",
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Post("/api/v1/names/register-did", nil, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("DID registered on-chain")
		if tx, ok := result["tx_hash"].(string); ok {
			fmt.Printf("TX: %s\n", tx)
		}
	},
}

var nameResolveDidCmd = &cobra.Command{
	Use:   "resolve-did <did>",
	Short: "Resolve a DID from the on-chain DID registry",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get(apiPath("/api/v1/names/did", args[0]), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		out, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(out))
	},
}

// gl star — star/unstar repositories

var starCmd = &cobra.Command{
	Use:   "star",
	Short: "Star and unstar repositories",
}

var starAddCmd = &cobra.Command{
	Use:   "add <repo>",
	Short: "Star a repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Post(repoPathSegments(args[0], "star"), nil, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Starred %s\n", args[0])
	},
}

var starRemoveCmd = &cobra.Command{
	Use:   "remove <repo>",
	Short: "Unstar a repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Post(repoPathSegments(args[0], "unstar"), nil, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Unstarred %s\n", args[0])
	},
}

var starCountCmd = &cobra.Command{
	Use:   "count <repo>",
	Short: "Show star count for a repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get(repoPathSegments(args[0], "stars"), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%v stars\n", result["stars"])
	},
}

// gl changelog — unified activity changelog

var changelogCmd = &cobra.Command{
	Use:   "changelog [repo]",
	Short: "Show unified activity changelog for a repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		path := repoPathSegments(args[0], "changelog") + fmt.Sprintf("?limit=%d", limit)

		var result struct {
			Events []struct {
				Type      string `json:"type"`
				Actor     string `json:"actor"`
				Message   string `json:"message"`
				Timestamp string `json:"timestamp"`
			} `json:"events"`
			Total int `json:"total"`
		}
		if err := client.Get(path, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, e := range result.Events {
			fmt.Printf("%s\t%s\t%s\t%s\n", e.Timestamp, e.Type, e.Actor, e.Message)
		}
		fmt.Fprintf(os.Stderr, "%d event(s)\n", result.Total)
	},
}

// gl whoami — print current identity

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Print your current identity (DID) and optional node info",
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")
		asJSON, _ := cmd.Flags().GetBool("json")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get("/api/v1/identity", &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if asJSON {
			out, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(out))
			return
		}

		fmt.Printf("DID:\t%v\n", result["did"])
		if name, ok := result["name"].(string); ok && name != "" {
			fmt.Printf("Name:\t%s\n", name)
		}
		if node, ok := result["node"].(string); ok && node != "" {
			fmt.Printf("Node:\t%s\n", node)
		}
		if registered, ok := result["registered"].(bool); ok {
			fmt.Printf("Registered:\t%v\n", registered)
		}
	},
}

func init() {
	// Cert flags
	for _, c := range []*cobra.Command{certListCmd, certShowCmd} {
		c.Flags().StringP("repo", "r", "", "Repository name (required)")
		c.MarkFlagRequired("repo")
		c.Flags().String("daemon-url", "", "Daemon URL")
	}

	// IPFS flags
	for _, c := range []*cobra.Command{ipfsListCmd, ipfsGetCmd} {
		c.Flags().String("daemon-url", "", "Daemon URL")
	}

	// MCP flags
	mcpServeCmd.Flags().String("daemon-url", "", "Daemon URL")

	// Sync flags
	for _, c := range []*cobra.Command{syncTriggerCmd, syncStatusCmd} {
		c.Flags().String("daemon-url", "", "Daemon URL")
	}

	// Name flags
	for _, c := range []*cobra.Command{nameRegisterCmd, nameResolveCmd, nameLookupCmd, nameAvailableCmd, nameRegisterDidCmd, nameResolveDidCmd} {
		c.Flags().String("daemon-url", "", "Daemon URL")
	}

	// Star flags
	for _, c := range []*cobra.Command{starAddCmd, starRemoveCmd, starCountCmd} {
		c.Flags().String("daemon-url", "", "Daemon URL")
	}

	// Changelog flags
	changelogCmd.Flags().IntP("limit", "n", 20, "Maximum number of events")
	changelogCmd.Flags().String("daemon-url", "", "Daemon URL")

	// Whoami flags
	whoamiCmd.Flags().String("daemon-url", "", "Daemon URL")
	whoamiCmd.Flags().Bool("json", false, "Output structured JSON")

	// Register subcommands
	certCmd.AddCommand(certListCmd, certShowCmd)
	ipfsCmd.AddCommand(ipfsListCmd, ipfsGetCmd)
	mcpCmd.AddCommand(mcpServeCmd)
	syncCmd.AddCommand(syncTriggerCmd, syncStatusCmd)
	nameCmd.AddCommand(nameRegisterCmd, nameResolveCmd, nameLookupCmd, nameAvailableCmd, nameRegisterDidCmd, nameResolveDidCmd)
	starCmd.AddCommand(starAddCmd, starRemoveCmd, starCountCmd)

	// Add to root
	rootCmd.AddCommand(certCmd, ipfsCmd, mcpCmd, syncCmd, nameCmd, starCmd, changelogCmd, whoamiCmd)
}
