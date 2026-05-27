package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var releaseCmd = &cobra.Command{
	Use:     "release",
	Short:   "Manage releases (like gh release)",
	Aliases: []string{"releases"},
}

var releaseListCmd = &cobra.Command{
	Use:   "list",
	Short: "List releases",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		client := newClient(cmd)

		var result struct {
			Releases []struct {
				ID    string `json:"id"`
				Tag   string `json:"tag"`
				Title string `json:"title"`
			} `json:"releases"`
			Total int `json:"total"`
		}
		if err := client.Get(repoPath(repo, "/releases"), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		w := newTabWriter()
		fmt.Fprintln(w, "ID\tTAG\tTITLE")
		for _, r := range result.Releases {
			fmt.Fprintf(w, "%s\t%s\t%s\n", r.ID, r.Tag, r.Title)
		}
		w.Flush()
		fmt.Fprintf(os.Stderr, "%d release(s)\n", result.Total)
	},
}

var releaseViewCmd = &cobra.Command{
	Use:   "view <release-id>",
	Short: "View a release",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		client := newClient(cmd)

		var result map[string]interface{}
		if err := client.Get(repoPathSegments(repo, "releases", args[0]), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		printJSON(result)
	},
}

var releaseCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a release",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		tag, _ := cmd.Flags().GetString("tag")
		title, _ := cmd.Flags().GetString("title")
		body, _ := cmd.Flags().GetString("notes")
		client := newClient(cmd)

		var result map[string]interface{}
		if err := client.Post(repoPath(repo, "/releases"), map[string]string{
			"tag":   tag,
			"title": title,
			"body":  body,
		}, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created release %v (tag %s)\n", result["id"], tag)
	},
}

var releaseDeleteCmd = &cobra.Command{
	Use:   "delete <release-id>",
	Short: "Delete a release",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		client := newClient(cmd)
		if err := client.Delete(repoPathSegments(repo, "releases", args[0])); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Deleted release %s\n", args[0])
	},
}

func init() {
	for _, c := range []*cobra.Command{releaseListCmd, releaseViewCmd, releaseCreateCmd, releaseDeleteCmd} {
		c.Flags().StringP("repo", "R", "", "Repository ID (required)")
		c.MarkFlagRequired("repo")
	}
	releaseCreateCmd.Flags().StringP("tag", "t", "", "Release tag (required)")
	releaseCreateCmd.Flags().String("title", "", "Release title (required)")
	releaseCreateCmd.Flags().String("notes", "", "Release notes")
	releaseCreateCmd.MarkFlagRequired("tag")
	releaseCreateCmd.MarkFlagRequired("title")

	addDaemonURLFlag(releaseListCmd, releaseViewCmd, releaseCreateCmd, releaseDeleteCmd)

	releaseCmd.AddCommand(releaseListCmd, releaseViewCmd, releaseCreateCmd, releaseDeleteCmd)
	rootCmd.AddCommand(releaseCmd)
}
