package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/GrayCodeAI/gitant-cli/internal/cli"
	"github.com/spf13/cobra"
)

// readSecretFromStdin reads a secret value from stdin, trimming trailing whitespace.
func readSecretFromStdin() (string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to stat stdin: %w", err)
	}
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		// stdin is a terminal, prompt the user
		fmt.Fprint(os.Stderr, "Enter secret value: ")
	}
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", err
		}
		return "", fmt.Errorf("no input received")
	}
	return strings.TrimSpace(scanner.Text()), nil
}

// DID methods — did:key, did:web, did:gitlawb

var identityResolveCmd = &cobra.Command{
	Use:   "resolve <did>",
	Short: "Resolve any DID method to its document",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get(apiPath("/api/v1/identity/resolve", args[0]), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		out, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(out))
	},
}

var identityRegisterDidCmd = &cobra.Command{
	Use:   "register-did",
	Short: "Anchor your DID document on-chain (did:gitlawb method)",
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Post("/api/v1/identity/register-did", nil, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("DID registered on-chain")
		if tx, ok := result["tx_hash"].(string); ok {
			fmt.Printf("TX: %s\n", tx)
		}
		if did, ok := result["did"].(string); ok {
			fmt.Printf("DID: %s\n", did)
		}
	},
}

// Signed ref-update certificates

var certThresholdCmd = &cobra.Command{
	Use:   "threshold <repo> <count>",
	Short: "Set required signature threshold for ref updates",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		req := map[string]interface{}{"threshold": args[1]}
		var result map[string]interface{}
		if err := client.Post(repoPathSegments(args[0], "certs", "threshold"), req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Threshold set to %s for %s\n", args[1], args[0])
	},
}

var certSignCmd = &cobra.Command{
	Use:   "sign <repo> <ref> <old> <new>",
	Short: "Sign a ref-update certificate",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		req := map[string]string{
			"ref":     args[1],
			"old_oid": args[2],
			"new_oid": args[3],
		}
		var result map[string]interface{}
		if err := client.Post(repoPathSegments(args[0], "certs", "sign"), req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Signed ref-update certificate\n")
		if id, ok := result["cert_id"].(string); ok {
			fmt.Printf("Cert ID: %s\n", id)
		}
	},
}

// Secrets subsystem

var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Manage secrets (capability-bound, KMS-backed)",
}

var secretsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List secret names (values are never shown)",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result struct {
			Secrets []struct {
				Name      string `json:"name"`
				CreatedAt string `json:"created_at"`
				UpdatedAt string `json:"updated_at"`
			} `json:"secrets"`
			Total int `json:"total"`
		}
		if err := client.Get(repoPathSegments(repo, "secrets"), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, s := range result.Secrets {
			fmt.Printf("%s\t(updated: %s)\n", s.Name, s.UpdatedAt)
		}
		fmt.Fprintf(os.Stderr, "%d secret(s)\n", result.Total)
	},
}

var secretsSetCmd = &cobra.Command{
	Use:   "set <name>",
	Short: "Set a secret (reads value from stdin to avoid exposing it in process list)",
	Long:  "Set a secret value. The value is read from stdin to avoid exposing it in shell history or process listings. Example: echo 'my-secret' | gt secrets set MY_SECRET",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		// Read secret value from stdin
		value, err := readSecretFromStdin()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading secret value: %v\n", err)
			os.Exit(1)
		}

		client := cli.NewClient(daemonURL)
		req := map[string]string{"name": args[0], "value": value}
		var result map[string]interface{}
		if err := client.Post(repoPathSegments(repo, "secrets"), req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Secret %s set\n", args[0])
	},
}

var secretsDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a secret",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		if err := client.Delete(repoPathSegments(repo, "secrets", args[0])); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Secret %s deleted\n", args[0])
	},
}

var secretsGetCmd = &cobra.Command{
	Use:   "get <name>",
	Short: "Get a secret value (requires secrets/read capability)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get(repoPathSegments(repo, "secrets", args[0]), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if val, ok := result["value"].(string); ok {
			fmt.Println(val)
		}
	},
}

// Trust score VCs

var trustCmd = &cobra.Command{
	Use:   "trust",
	Short: "Manage agent trust scores (Verifiable Credentials)",
}

var trustShowCmd = &cobra.Command{
	Use:   "show <did>",
	Short: "Show trust score and VC for an agent",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get(apiPath("/api/v1/agents", args[0], "trust"), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("DID:\t%s\n", args[0])
		fmt.Printf("Score:\t%v\n", result["score"])
		fmt.Printf("Tier:\t%v\n", result["tier"])
		if breakdown, ok := result["breakdown"].(map[string]interface{}); ok {
			fmt.Println("\nBreakdown:")
			for k, v := range breakdown {
				fmt.Printf("  %s:\t%v\n", k, v)
			}
		}
		if vc, ok := result["vc"].(map[string]interface{}); ok {
			fmt.Println("\nVerifiable Credential:")
			out, _ := json.MarshalIndent(vc, "", "  ")
			fmt.Println(string(out))
		}
	},
}

var trustIssueCmd = &cobra.Command{
	Use:   "issue <did>",
	Short: "Issue a trust score VC for an agent",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Post(apiPath("/api/v1/agents", args[0], "trust", "issue"), nil, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Trust VC issued for %s\n", args[0])
		if score, ok := result["score"].(float64); ok {
			fmt.Printf("Score: %.2f\n", score)
		}
	},
}

var trustVerifyCmd = &cobra.Command{
	Use:   "verify <vc-jwt>",
	Short: "Verify a trust score VC",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		req := map[string]string{"vc": args[0]}
		var result map[string]interface{}
		if err := client.Post("/api/v1/trust/verify", req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if valid, ok := result["valid"].(bool); ok && valid {
			fmt.Println("VC is VALID")
		} else {
			fmt.Println("VC is INVALID")
		}
		if reason, ok := result["reason"].(string); ok && reason != "" {
			fmt.Printf("Reason: %s\n", reason)
		}
	},
}

// Repo tokenization

var repoTokenizeCmd = &cobra.Command{
	Use:   "tokenize <repo>",
	Short: "Deploy an ERC-20 token tied to this repo (auto-splits on PR merge)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")
		name, _ := cmd.Flags().GetString("name")
		symbol, _ := cmd.Flags().GetString("symbol")

		client := cli.NewClient(daemonURL)
		req := map[string]string{"name": name, "symbol": symbol}
		var result map[string]interface{}
		if err := client.Post(repoPathSegments(args[0], "tokenize"), req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Token deployed for %s\n", args[0])
		if addr, ok := result["token_address"].(string); ok {
			fmt.Printf("Contract: %s\n", addr)
		}
		if tx, ok := result["tx_hash"].(string); ok {
			fmt.Printf("TX: %s\n", tx)
		}
	},
}

// Maintainers file

var maintainersCmd = &cobra.Command{
	Use:   "maintainers",
	Short: "Manage .gitant/maintainers file (authorized DIDs for ref signing)",
}

var maintainersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List maintainers for a repository",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result struct {
			Maintainers []struct {
				DID       string `json:"did"`
				Key       string `json:"key"`
				AddedAt   string `json:"added_at"`
				Threshold int    `json:"threshold"`
			} `json:"maintainers"`
		}
		if err := client.Get(repoPathSegments(repo, "maintainers"), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, m := range result.Maintainers {
			fmt.Printf("%s\t%s\t(threshold: %d)\n", m.DID, m.AddedAt, m.Threshold)
		}
	},
}

var maintainersAddCmd = &cobra.Command{
	Use:   "add <did>",
	Short: "Add a maintainer to a repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		req := map[string]string{"did": args[0]}
		var result map[string]interface{}
		if err := client.Post(repoPathSegments(repo, "maintainers"), req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Added %s as maintainer\n", args[0])
	},
}

var maintainersRemoveCmd = &cobra.Command{
	Use:   "remove <did>",
	Short: "Remove a maintainer from a repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		if err := client.Delete(repoPathSegments(repo, "maintainers", args[0])); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Removed %s as maintainer\n", args[0])
	},
}

func init() {
	// Identity new flags
	for _, c := range []*cobra.Command{identityResolveCmd, identityRegisterDidCmd} {
		c.Flags().String("daemon-url", "", "Daemon URL")
	}

	// Cert new flags
	for _, c := range []*cobra.Command{certThresholdCmd, certSignCmd} {
		c.Flags().String("daemon-url", "", "Daemon URL")
	}

	// Secrets flags
	for _, c := range []*cobra.Command{secretsListCmd, secretsSetCmd, secretsGetCmd, secretsDeleteCmd} {
		c.Flags().StringP("repo", "r", "", "Repository name (required)")
		c.MarkFlagRequired("repo")
		c.Flags().String("daemon-url", "", "Daemon URL")
	}

	// Trust flags
	for _, c := range []*cobra.Command{trustShowCmd, trustIssueCmd, trustVerifyCmd} {
		c.Flags().String("daemon-url", "", "Daemon URL")
	}

	// Repo tokenize flags
	repoTokenizeCmd.Flags().String("daemon-url", "", "Daemon URL")
	repoTokenizeCmd.Flags().String("name", "", "Token name")
	repoTokenizeCmd.Flags().String("symbol", "", "Token symbol")

	// Maintainers flags
	for _, c := range []*cobra.Command{maintainersListCmd, maintainersAddCmd, maintainersRemoveCmd} {
		c.Flags().StringP("repo", "r", "", "Repository name (required)")
		c.MarkFlagRequired("repo")
		c.Flags().String("daemon-url", "", "Daemon URL")
	}

	// Register subcommands
	identityCmd.AddCommand(identityResolveCmd, identityRegisterDidCmd)
	certCmd.AddCommand(certThresholdCmd, certSignCmd)
	secretsCmd.AddCommand(secretsListCmd, secretsSetCmd, secretsGetCmd, secretsDeleteCmd)
	trustCmd.AddCommand(trustShowCmd, trustIssueCmd, trustVerifyCmd)
	maintainersCmd.AddCommand(maintainersListCmd, maintainersAddCmd, maintainersRemoveCmd)

	// Add repo tokenize to repo command
	repoCmd.AddCommand(repoTokenizeCmd)

	// Add new top-level commands
	rootCmd.AddCommand(secretsCmd, trustCmd, maintainersCmd)
}
