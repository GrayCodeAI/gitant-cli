package main

import (
	"fmt"
	"os"

	"github.com/GrayCodeAI/gitant-cli/internal/cli"
	"github.com/spf13/cobra"
)

var extensionCmd = &cobra.Command{
	Use:   "extension",
	Short: "Manage CLI extensions (from GitHub)",
	Aliases: []string{"ext"},
}

var extensionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed extensions",
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")
		client := cli.NewClient(daemonURL)

		var result struct {
			Extensions []struct {
				Name        string `json:"name"`
				Description string `json:"description"`
				Version     string `json:"version"`
			} `json:"extensions"`
		}
		if err := client.Get("/api/v1/extensions", &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, ext := range result.Extensions {
			fmt.Printf("%s\t%s\t%s\n", ext.Name, ext.Version, ext.Description)
		}
	},
}

var extensionInstallCmd = &cobra.Command{
	Use:   "install [name]",
	Short: "Install an extension",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")
		client := cli.NewClient(daemonURL)

		var result map[string]interface{}
		if err := client.Post(fmt.Sprintf("/api/v1/extensions/%s/install", args[0]), nil, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Installed extension: %s\n", args[0])
	},
}

var extensionUninstallCmd = &cobra.Command{
	Use:   "uninstall [name]",
	Short: "Uninstall an extension",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")
		client := cli.NewClient(daemonURL)

		if err := client.Delete(fmt.Sprintf("/api/v1/extensions/%s", args[0])); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Uninstalled extension: %s\n", args[0])
	},
}

var kanbanCmd = &cobra.Command{
	Use:   "kanban",
	Short: "Manage Kanban boards (from Gitea)",
	Aliases: []string{"kb", "board"},
}

var kanbanListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Kanban boards",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")
		client := cli.NewClient(daemonURL)

		var result struct {
			Boards []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"boards"`
		}
		if err := client.Get(fmt.Sprintf("/api/v1/repos/%s/kanban", repo), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, b := range result.Boards {
			fmt.Printf("%s\t%s\n", b.ID, b.Name)
		}
	},
}

var kanbanCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Kanban board",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		name, _ := cmd.Flags().GetString("name")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		if name == "" {
			name = PromptRequired("Board name")
		}

		client := cli.NewClient(daemonURL)
		req := map[string]string{"name": name}

		var result map[string]interface{}
		if err := client.Post(fmt.Sprintf("/api/v1/repos/%s/kanban", repo), req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created board: %s\n", result["id"])
	},
}

var epicCmd = &cobra.Command{
	Use:   "epic",
	Short: "Manage epics (from GitLab)",
}

var epicListCmd = &cobra.Command{
	Use:   "list",
	Short: "List epics",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		status, _ := cmd.Flags().GetString("status")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")
		client := cli.NewClient(daemonURL)

		path := fmt.Sprintf("/api/v1/repos/%s/epics", repo)
		if status != "" {
			path += "?status=" + status
		}

		var result struct {
			Epics []struct {
				ID       string  `json:"id"`
				Title    string  `json:"title"`
				Status   string  `json:"status"`
				Progress float64 `json:"progress"`
			} `json:"epics"`
		}
		if err := client.Get(path, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, e := range result.Epics {
			fmt.Printf("%s\t%s\t%s\t%.0f%%\n", e.ID, e.Title, e.Status, e.Progress*100)
		}
	},
}

var epicCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an epic",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		if title == "" {
			title = PromptRequired("Title")
		}

		client := cli.NewClient(daemonURL)
		req := map[string]string{
			"title":       title,
			"description": description,
		}

		var result map[string]interface{}
		if err := client.Post(fmt.Sprintf("/api/v1/repos/%s/epics", repo), req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created epic: %s\n", result["id"])
	},
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Manage seed nodes (from Radicle)",
}

var seedListCmd = &cobra.Command{
	Use:   "list",
	Short: "List seed nodes",
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")
		client := cli.NewClient(daemonURL)

		var result struct {
			Seeds []struct {
				URL         string  `json:"url"`
				Status      string  `json:"status"`
				Reliability float64 `json:"reliability"`
			} `json:"seeds"`
		}
		if err := client.Get("/api/v1/seeds", &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, s := range result.Seeds {
			fmt.Printf("%s\t%s\t%.0f%%\n", s.URL, s.Status, s.Reliability*100)
		}
	},
}

var seedAddCmd = &cobra.Command{
	Use:   "add [url]",
	Short: "Add a seed node",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")
		client := cli.NewClient(daemonURL)

		req := map[string]string{"url": args[0]}
		var result map[string]interface{}
		if err := client.Post("/api/v1/seeds", req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Added seed: %s\n", args[0])
	},
}

var stackedCmd = &cobra.Command{
	Use:   "stack",
	Short: "Manage stacked diffs (from Sapling)",
}

var stackListCmd = &cobra.Command{
	Use:   "list",
	Short: "List stacked diffs",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")
		client := cli.NewClient(daemonURL)

		var result struct {
			Diffs []struct {
				ID     string `json:"id"`
				Title  string `json:"title"`
				Status string `json:"status"`
			} `json:"diffs"`
		}
		if err := client.Get(fmt.Sprintf("/api/v1/repos/%s/stacked-diffs", repo), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, d := range result.Diffs {
			fmt.Printf("%s\t%s\t%s\n", d.ID, d.Title, d.Status)
		}
	},
}

func init() {
	// Extension
	for _, c := range []*cobra.Command{extensionListCmd, extensionInstallCmd, extensionUninstallCmd} {
		c.Flags().String("daemon-url", "", "Daemon URL")
	}
	extensionCmd.AddCommand(extensionListCmd, extensionInstallCmd, extensionUninstallCmd)

	// Kanban
	for _, c := range []*cobra.Command{kanbanListCmd, kanbanCreateCmd} {
		c.Flags().StringP("repo", "r", "", "Repository name (required)")
		c.MarkFlagRequired("repo")
		c.Flags().String("daemon-url", "", "Daemon URL")
	}
	kanbanCreateCmd.Flags().StringP("name", "n", "", "Board name")
	kanbanCmd.AddCommand(kanbanListCmd, kanbanCreateCmd)

	// Epic
	for _, c := range []*cobra.Command{epicListCmd, epicCreateCmd} {
		c.Flags().StringP("repo", "r", "", "Repository name (required)")
		c.MarkFlagRequired("repo")
		c.Flags().String("daemon-url", "", "Daemon URL")
	}
	epicListCmd.Flags().String("status", "", "Filter by status")
	epicCreateCmd.Flags().StringP("title", "t", "", "Epic title")
	epicCreateCmd.Flags().StringP("description", "d", "", "Epic description")
	epicCmd.AddCommand(epicListCmd, epicCreateCmd)

	// Seed
	for _, c := range []*cobra.Command{seedListCmd, seedAddCmd} {
		c.Flags().String("daemon-url", "", "Daemon URL")
	}
	seedCmd.AddCommand(seedListCmd, seedAddCmd)

	// Stacked
	for _, c := range []*cobra.Command{stackListCmd} {
		c.Flags().StringP("repo", "r", "", "Repository name (required)")
		c.MarkFlagRequired("repo")
		c.Flags().String("daemon-url", "", "Daemon URL")
	}
	stackedCmd.AddCommand(stackListCmd)

	rootCmd.AddCommand(extensionCmd, kanbanCmd, epicCmd, seedCmd, stackedCmd)
}
