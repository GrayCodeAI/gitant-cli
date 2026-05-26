package main

import (
	"fmt"
	"os"

	"github.com/GrayCodeAI/gitant-cli/internal/cli"
	"github.com/spf13/cobra"
)

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage repositories (like gh repo)",
}

var repoListCmd = &cobra.Command{
	Use:   "list",
	Short: "List repositories",
	Run: func(cmd *cobra.Command, args []string) {
		client := newClient(cmd)
		var result struct {
			Repos []struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				Description string `json:"description"`
				Stars       int    `json:"stars"`
				Private     bool   `json:"private"`
			} `json:"repos"`
			Total int `json:"total"`
		}
		if err := client.Get("/api/v1/repos", &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		w := newTabWriter()
		fmt.Fprintln(w, "ID\tSTARS\tPRIVATE\tDESCRIPTION")
		for _, repo := range result.Repos {
			desc := repo.Description
			if desc == "" {
				desc = "-"
			}
			fmt.Fprintf(w, "%s\t%d\t%v\t%s\n", repo.ID, repo.Stars, repo.Private, desc)
		}
		w.Flush()
		fmt.Fprintf(os.Stderr, "%d repo(s)\n", result.Total)
	},
}

var repoViewCmd = &cobra.Command{
	Use:   "view [repo-id]",
	Short: "View repository details (like gh repo view)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		jsonOut, _ := cmd.Flags().GetBool("json")
		client := newClient(cmd)

		var result map[string]interface{}
		if err := client.Get(repoPath(args[0], ""), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if jsonOut {
			printJSON(result)
			return
		}

		fmt.Printf("name:\t%v\n", result["name"])
		fmt.Printf("id:\t%v\n", result["id"])
		fmt.Printf("description:\t%v\n", result["description"])
		fmt.Printf("private:\t%v\n", result["private"])
		fmt.Printf("stars:\t%v\n", result["stars"])
		if created, ok := result["created_at"]; ok {
			fmt.Printf("created:\t%v\n", created)
		}
	},
}

var repoCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a repository (like gh repo create)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		desc, _ := cmd.Flags().GetString("description")
		private, _ := cmd.Flags().GetBool("private")
		client := newClient(cmd)

		var result map[string]interface{}
		if err := client.Post("/api/v1/repos", map[string]interface{}{
			"name":        args[0],
			"description": desc,
			"private":     private,
		}, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created repository %v\n", result["id"])
	},
}

var repoDeleteCmd = &cobra.Command{
	Use:   "delete <repo-id>",
	Short: "Delete a repository (like gh repo delete)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		confirm, _ := cmd.Flags().GetBool("yes")
		if !confirm {
			fmt.Fprintf(os.Stderr, "Use --yes to confirm deletion of %s\n", args[0])
			os.Exit(1)
		}
		client := newClient(cmd)
		if err := client.Delete(repoPath(args[0], "")); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Deleted repository %s\n", args[0])
	},
}

var repoCloneCmd = &cobra.Command{
	Use:   "clone <repo-id> [directory]",
	Short: "Clone a repository (alias for gitant clone)",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]
		if len(args) > 1 {
			dir = args[1]
		}
		remote, _ := cmd.Flags().GetString("remote")
		if remote == "" {
			client := newClient(cmd)
			remote = client.BaseURL
		}
		if err := cli.Clone(remote, args[0], dir); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var repoStarCmd = &cobra.Command{
	Use:   "star <repo-id>",
	Short: "Star a repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newClient(cmd)
		var result map[string]interface{}
		if err := client.Post(repoPath(args[0], "/star"), nil, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Starred %s (%v stars)\n", args[0], result["stars"])
	},
}

var repoUnstarCmd = &cobra.Command{
	Use:   "unstar <repo-id>",
	Short: "Unstar a repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := newClient(cmd)
		var result map[string]interface{}
		if err := client.Post(repoPath(args[0], "/unstar"), nil, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Unstarred %s (%v stars)\n", args[0], result["stars"])
	},
}

var repoForkCmd = &cobra.Command{
	Use:   "fork <source> <name>",
	Short: "Fork a repository",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := newClient(cmd)
		var result map[string]interface{}
		if err := client.Post(repoPath(args[0], "/fork"), map[string]string{"name": args[1]}, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Forked %s -> %v\n", args[0], result["id"])
	},
}

func init() {
	repoViewCmd.Flags().Bool("json", false, "Output JSON")
	repoCreateCmd.Flags().StringP("description", "d", "", "Repository description")
	repoCreateCmd.Flags().Bool("private", false, "Create as private repository")
	repoDeleteCmd.Flags().Bool("yes", false, "Confirm deletion")
	repoCloneCmd.Flags().StringP("remote", "r", "", "Remote daemon URL")

	addDaemonURLFlag(repoListCmd, repoViewCmd, repoCreateCmd, repoDeleteCmd, repoCloneCmd, repoStarCmd, repoUnstarCmd, repoForkCmd)

	repoCmd.AddCommand(repoListCmd, repoViewCmd, repoCreateCmd, repoDeleteCmd, repoCloneCmd, repoStarCmd, repoUnstarCmd, repoForkCmd)
	rootCmd.AddCommand(repoCmd)
}
