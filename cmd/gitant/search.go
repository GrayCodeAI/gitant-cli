package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search code in a repository (like gh search code)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		if repo == "" {
			fmt.Fprintln(os.Stderr, "Error: --repo is required")
			os.Exit(1)
		}

		client := newClient(cmd)
		path := fmt.Sprintf("/api/v1/repos/%s/search?q=%s", url.PathEscape(repo), url.QueryEscape(args[0]))

		var result struct {
			Matches []struct {
				Path    string `json:"path"`
				Line    int    `json:"line"`
				Content string `json:"content"`
			} `json:"matches"`
			Total int `json:"total"`
		}
		if err := client.Get(path, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		w := newTabWriter()
		fmt.Fprintln(w, "PATH\tLINE\tMATCH")
		for _, m := range result.Matches {
			fmt.Fprintf(w, "%s\t%d\t%s\n", m.Path, m.Line, m.Content)
		}
		w.Flush()
		fmt.Fprintf(os.Stderr, "%d match(es)\n", result.Total)
	},
}

func init() {
	searchCmd.Flags().StringP("repo", "R", "", "Repository ID (required)")
	searchCmd.MarkFlagRequired("repo")
	addDaemonURLFlag(searchCmd)
	rootCmd.AddCommand(searchCmd)
}
