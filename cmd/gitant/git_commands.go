package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/GrayCodeAI/gitant-cli/internal/cli"
	"github.com/spf13/cobra"
)

// Git wrapper commands - these wrap native git with Gitant enhancements

var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Git operations with Gitant enhancements",
	Long:  "Git operations that work with your Gitant node. These wrap native git commands with platform features.",
}

// git fetch
var gitFetchCmd = &cobra.Command{
	Use:   "fetch [remote]",
	Short: "Fetch from a Gitant remote",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		remote := "origin"
		if len(args) > 0 {
			remote = args[0]
		}

		fetch := exec.Command("git", "fetch", remote)
		fetch.Stdout = os.Stdout
		fetch.Stderr = os.Stderr
		if err := fetch.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git remote
var gitRemoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Manage Gitant remotes",
}

var gitRemoteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List remotes",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "remote", "-v").Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitRemoteAddCmd = &cobra.Command{
	Use:   "add <name> <url>",
	Short: "Add a Gitant remote",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		add := exec.Command("git", "remote", "add", args[0], args[1])
		add.Stdout = os.Stdout
		add.Stderr = os.Stderr
		if err := add.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Added remote %s → %s\n", args[0], args[1])
	},
}

var gitRemoteRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a remote",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		remove := exec.Command("git", "remote", "remove", args[0])
		remove.Stdout = os.Stdout
		remove.Stderr = os.Stderr
		if err := remove.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Removed remote %s\n", args[0])
	},
}

var gitRemoteSetUrlCmd = &cobra.Command{
	Use:   "set-url <name> <url>",
	Short: "Set remote URL",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		setUrl := exec.Command("git", "remote", "set-url", args[0], args[1])
		setUrl.Stdout = os.Stdout
		setUrl.Stderr = os.Stderr
		if err := setUrl.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Set remote %s → %s\n", args[0], args[1])
	},
}

var gitRemoteRenameCmd = &cobra.Command{
	Use:   "rename <old> <new>",
	Short: "Rename a remote",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		rename := exec.Command("git", "remote", "rename", args[0], args[1])
		rename.Stdout = os.Stdout
		rename.Stderr = os.Stderr
		if err := rename.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Renamed remote %s to %s\n", args[0], args[1])
	},
}

var gitRemoteShowCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "Show information about a remote",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "remote", "show", args[0]).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitRemotePruneCmd = &cobra.Command{
	Use:   "prune <name>",
	Short: "Delete stale tracking references",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "remote", "prune", args[0]).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitRemoteUpdateCmd = &cobra.Command{
	Use:   "update [remote]",
	Short: "Fetch updates for remotes",
	Run: func(cmd *cobra.Command, args []string) {
		remoteArgs := append([]string{"remote", "update"}, args...)
		out, err := exec.Command("git", remoteArgs...).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitRemoteGetUrlCmd = &cobra.Command{
	Use:   "get-url <name>",
	Short: "Get URL for a remote",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		push, _ := cmd.Flags().GetBool("push")
		remoteArgs := []string{"remote", "get-url"}
		if push {
			remoteArgs = append(remoteArgs, "--push")
		}
		remoteArgs = append(remoteArgs, args[0])
		out, err := exec.Command("git", remoteArgs...).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitRemoteSetBranchesCmd = &cobra.Command{
	Use:   "set-branches <name> <branch>...",
	Short: "Change the list of branches tracked",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		add, _ := cmd.Flags().GetBool("add")
		remoteArgs := []string{"remote", "set-branches"}
		if add {
			remoteArgs = append(remoteArgs, "--add")
		}
		remoteArgs = append(remoteArgs, args...)
		out, err := exec.Command("git", remoteArgs...).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitRemoteSetHeadCmd = &cobra.Command{
	Use:   "set-head <name> <branch>",
	Short: "Set or delete the default branch",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "remote", "set-head", args[0], args[1]).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitRemoteGetHeadCmd = &cobra.Command{
	Use:   "get-head <name>",
	Short: "Query which HEAD the remote has",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "remote", "get-head", args[0]).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitRemoteSetPushCmd = &cobra.Command{
	Use:   "set-push <name> <url>",
	Short: "Change push URLs",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "remote", "set-push", args[0], args[1]).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

// git branch
var gitBranchCmd = &cobra.Command{
	Use:   "branch",
	Short: "Manage branches",
}

var gitBranchListCmd = &cobra.Command{
	Use:   "list",
	Short: "List branches",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		if repo != "" {
			// List from server
			client := cli.NewClient(daemonURL)
			var result struct {
				Refs []struct {
					Name   string `json:"name"`
					Commit string `json:"commit"`
				} `json:"refs"`
			}
			if err := client.Get(fmt.Sprintf("/api/v1/repos/%s/refs", repo), &result); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			for _, ref := range result.Refs {
				if strings.HasPrefix(ref.Name, "refs/heads/") {
					branch := strings.TrimPrefix(ref.Name, "refs/heads/")
					fmt.Printf("%s\t%s\n", branch, ref.Commit[:8])
				}
			}
		} else {
			// List local
			out, err := exec.Command("git", "branch", "-v").Output()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Print(string(out))
		}
	},
}

var gitBranchCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a branch",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		if repo != "" {
			// Create on server
			client := cli.NewClient(daemonURL)
			req := map[string]string{"name": args[0]}
			var result map[string]interface{}
			if err := client.Post(fmt.Sprintf("/api/v1/repos/%s/branches", repo), req, &result); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Created branch %s on server\n", args[0])
		} else {
			// Create local
			create := exec.Command("git", "branch", args[0])
			create.Stdout = os.Stdout
			create.Stderr = os.Stderr
			if err := create.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Created branch %s\n", args[0])
		}
	},
}

var gitBranchDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a branch",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")

		deleteArgs := []string{"branch"}
		if force {
			deleteArgs = append(deleteArgs, "-D")
		} else {
			deleteArgs = append(deleteArgs, "-d")
		}
		deleteArgs = append(deleteArgs, args[0])

		delete := exec.Command("git", deleteArgs...)
		delete.Stdout = os.Stdout
		delete.Stderr = os.Stderr
		if err := delete.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Deleted branch %s\n", args[0])
	},
}

// git tag
var gitTagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Manage tags",
}

var gitTagListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tags",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "tag", "-l").Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if len(out) == 0 {
			fmt.Println("No tags found")
		} else {
			fmt.Print(string(out))
		}
	},
}

var gitTagCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a tag",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		annotate, _ := cmd.Flags().GetBool("annotate")
		message, _ := cmd.Flags().GetString("message")

		tagArgs := []string{"tag"}
		if annotate {
			tagArgs = append(tagArgs, "-a", args[0])
			if message != "" {
				tagArgs = append(tagArgs, "-m", message)
			}
		} else {
			tagArgs = append(tagArgs, args[0])
		}

		tag := exec.Command("git", tagArgs...)
		tag.Stdout = os.Stdout
		tag.Stderr = os.Stderr
		if err := tag.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created tag %s\n", args[0])
	},
}

var gitTagDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a tag",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		delete := exec.Command("git", "tag", "-d", args[0])
		delete.Stdout = os.Stdout
		delete.Stderr = os.Stderr
		if err := delete.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Deleted tag %s\n", args[0])
	},
}

// git log
var gitLogCmd = &cobra.Command{
	Use:   "log",
	Short: "Show commit log",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		limit, _ := cmd.Flags().GetInt("limit")
		oneline, _ := cmd.Flags().GetBool("oneline")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		if repo != "" {
			// Get from server
			client := cli.NewClient(daemonURL)
			path := fmt.Sprintf("/api/v1/repos/%s/commits?limit=%d", repo, limit)

			var result struct {
				Commits []struct {
					Hash    string `json:"hash"`
					Message string `json:"message"`
					Author  string `json:"author"`
					Date    string `json:"date"`
				} `json:"commits"`
			}
			if err := client.Get(path, &result); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			for _, c := range result.Commits {
				if oneline {
					fmt.Printf("%s %s\n", c.Hash[:8], strings.Split(c.Message, "\n")[0])
				} else {
					fmt.Printf("commit %s\nAuthor: %s\nDate: %s\n\n    %s\n\n", c.Hash, c.Author, c.Date, c.Message)
				}
			}
		} else {
			// Local log
			logArgs := []string{"log"}
			if oneline {
				logArgs = append(logArgs, "--oneline")
			}
			if limit > 0 {
				logArgs = append(logArgs, fmt.Sprintf("-%d", limit))
			}

			out, err := exec.Command("git", logArgs...).Output()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Print(string(out))
		}
	},
}

// git diff
var gitDiffCmd = &cobra.Command{
	Use:   "diff [commit1] [commit2]",
	Short: "Show diffs",
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		if repo != "" && len(args) >= 2 {
			// Get from server
			client := cli.NewClient(daemonURL)
			path := fmt.Sprintf("/api/v1/repos/%s/diff?from=%s&to=%s", repo, args[0], args[1])

			var result struct {
				Files []struct {
					Path     string `json:"path"`
					Status   string `json:"status"`
					Additions int `json:"additions"`
					Deletions int `json:"deletions"`
				} `json:"files"`
			}
			if err := client.Get(path, &result); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			for _, f := range result.Files {
				fmt.Printf("%s %s (+%d, -%d)\n", f.Status, f.Path, f.Additions, f.Deletions)
			}
		} else {
			// Local diff
			diffArgs := append([]string{"diff"}, args...)
			out, err := exec.Command("git", diffArgs...).Output()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Print(string(out))
		}
	},
}

// git status
var gitStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show working tree status",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "status").Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

// git show
var gitShowCmd = &cobra.Command{
	Use:   "show <object>",
	Short: "Show git object details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		if repo != "" {
			// Get from server
			client := cli.NewClient(daemonURL)
			path := fmt.Sprintf("/api/v1/repos/%s/commits/%s", repo, args[0])

			var result map[string]interface{}
			if err := client.Get(path, &result); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("commit %v\n", result["hash"])
			fmt.Printf("Author: %v\n", result["author"])
			fmt.Printf("Date: %v\n\n", result["date"])
			fmt.Printf("    %v\n", result["message"])
		} else {
			out, err := exec.Command("git", "show", args[0]).Output()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Print(string(out))
		}
	},
}

// git merge
var gitMergeCmd = &cobra.Command{
	Use:   "merge <branch>",
	Short: "Merge a branch",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		merge := exec.Command("git", "merge", args[0])
		merge.Stdout = os.Stdout
		merge.Stderr = os.Stderr
		if err := merge.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Merged %s\n", args[0])
	},
}

// git cherry-pick
var gitCherryPickCmd = &cobra.Command{
	Use:   "cherry-pick <commit>",
	Short: "Cherry-pick a commit",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cherryPick := exec.Command("git", "cherry-pick", args[0])
		cherryPick.Stdout = os.Stdout
		cherryPick.Stderr = os.Stderr
		if err := cherryPick.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Cherry-picked %s\n", args[0])
	},
}

// git revert
var gitRevertCmd = &cobra.Command{
	Use:   "revert <commit>",
	Short: "Revert a commit",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		revert := exec.Command("git", "revert", args[0])
		revert.Stdout = os.Stdout
		revert.Stderr = os.Stderr
		if err := revert.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Reverted %s\n", args[0])
	},
}

// git blame
var gitBlameCmd = &cobra.Command{
	Use:   "blame <file>",
	Short: "Show file authorship",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		if repo != "" {
			// Get from server
			client := cli.NewClient(daemonURL)
			path := fmt.Sprintf("/api/v1/repos/%s/files/%s/blame", repo, args[0])

			var result struct {
				Lines []struct {
					Line    int    `json:"line"`
					Author  string `json:"author"`
					Content string `json:"content"`
				} `json:"lines"`
			}
			if err := client.Get(path, &result); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			for _, l := range result.Lines {
				fmt.Printf("%s %4d\t%s\n", l.Author, l.Line, l.Content)
			}
		} else {
			out, err := exec.Command("git", "blame", args[0]).Output()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Print(string(out))
		}
	},
}

// git grep
var gitGrepCmd = &cobra.Command{
	Use:   "grep <pattern>",
	Short: "Search for patterns",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		if repo != "" {
			// Search on server
			client := cli.NewClient(daemonURL)
			path := fmt.Sprintf("/api/v1/repos/%s/search?q=%s", repo, args[0])

			var result struct {
				Results []struct {
					Path    string `json:"path"`
					Line    int    `json:"line"`
					Content string `json:"content"`
				} `json:"results"`
			}
			if err := client.Get(path, &result); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}

			for _, r := range result.Results {
				fmt.Printf("%s:%d: %s\n", r.Path, r.Line, r.Content)
			}
		} else {
			out, err := exec.Command("git", "grep", args[0]).Output()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Print(string(out))
		}
	},
}

// git stash
var gitStashCmd = &cobra.Command{
	Use:   "stash",
	Short: "Stash changes",
}

var gitStashSaveCmd = &cobra.Command{
	Use:   "save [message]",
	Short: "Save stash",
	Run: func(cmd *cobra.Command, args []string) {
		stashArgs := append([]string{"stash", "save"}, args...)
		out, err := exec.Command("git", stashArgs...).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitStashListCmd = &cobra.Command{
	Use:   "list",
	Short: "List stashes",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "stash", "list").Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if len(out) == 0 {
			fmt.Println("No stashes found")
		} else {
			fmt.Print(string(out))
		}
	},
}

var gitStashPopCmd = &cobra.Command{
	Use:   "pop",
	Short: "Pop stash",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "stash", "pop").Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitStashDropCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop stash",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "stash", "drop").Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitStashShowCmd = &cobra.Command{
	Use:   "show [stash]",
	Short: "Show the changes recorded in the stash",
	Run: func(cmd *cobra.Command, args []string) {
		stashArgs := append([]string{"stash", "show"}, args...)
		out, err := exec.Command("git", stashArgs...).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitStashApplyCmd = &cobra.Command{
	Use:   "apply [stash]",
	Short: "Apply stash without removing it",
	Run: func(cmd *cobra.Command, args []string) {
		stashArgs := append([]string{"stash", "apply"}, args...)
		out, err := exec.Command("git", stashArgs...).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitStashBranchCmd = &cobra.Command{
	Use:   "branch <branchname> [stash]",
	Short: "Create a new branch from a stash",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stashArgs := append([]string{"stash", "branch"}, args...)
		out, err := exec.Command("git", stashArgs...).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitStashClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Remove all stashes",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "stash", "clear").Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitStashStoreCmd = &cobra.Command{
	Use:   "store <commit>",
	Short: "Store a stash without touching the working tree",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "stash", "store", args[0]).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

// git worktree
var gitWorktreeCmd = &cobra.Command{
	Use:   "worktree",
	Short: "Manage worktrees",
}

var gitWorktreeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List worktrees",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "worktree", "list").Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitWorktreeAddCmd = &cobra.Command{
	Use:   "add <path> [branch]",
	Short: "Add a worktree",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		worktreeArgs := append([]string{"worktree", "add"}, args...)
		out, err := exec.Command("git", worktreeArgs...).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitWorktreeRemoveCmd = &cobra.Command{
	Use:   "remove <path>",
	Short: "Remove a worktree",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")
		worktreeArgs := []string{"worktree", "remove"}
		if force {
			worktreeArgs = append(worktreeArgs, "--force")
		}
		worktreeArgs = append(worktreeArgs, args[0])
		out, err := exec.Command("git", worktreeArgs...).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitWorktreeMoveCmd = &cobra.Command{
	Use:   "move <path> <new-path>",
	Short: "Move a worktree",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		worktreeArgs := []string{"worktree", "move", args[0], args[1]}
		out, err := exec.Command("git", worktreeArgs...).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitWorktreeLockCmd = &cobra.Command{
	Use:   "lock <path>",
	Short: "Lock a worktree",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		reason, _ := cmd.Flags().GetString("reason")
		worktreeArgs := []string{"worktree", "lock"}
		if reason != "" {
			worktreeArgs = append(worktreeArgs, "--reason", reason)
		}
		worktreeArgs = append(worktreeArgs, args[0])
		out, err := exec.Command("git", worktreeArgs...).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitWorktreeUnlockCmd = &cobra.Command{
	Use:   "unlock <path>",
	Short: "Unlock a worktree",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "worktree", "unlock", args[0]).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitWorktreePruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Prune stale working tree information",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "worktree", "prune").Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitWorktreeRepairCmd = &cobra.Command{
	Use:   "repair [path]",
	Short: "Repair working tree administrative files",
	Run: func(cmd *cobra.Command, args []string) {
		worktreeArgs := append([]string{"worktree", "repair"}, args...)
		out, err := exec.Command("git", worktreeArgs...).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

// git submodule
var gitSubmoduleCmd = &cobra.Command{
	Use:   "submodule",
	Short: "Manage submodules",
}

var gitSubmoduleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List submodules",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "submodule", "status").Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if len(out) == 0 {
			fmt.Println("No submodules found")
		} else {
			fmt.Print(string(out))
		}
	},
}

var gitSubmoduleInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize submodules",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "submodule", "init").Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitSubmoduleUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update submodules",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "submodule", "update").Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

// git reset
var gitResetCmd = &cobra.Command{
	Use:   "reset [commit]",
	Short: "Reset HEAD",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hard, _ := cmd.Flags().GetBool("hard")
		soft, _ := cmd.Flags().GetBool("soft")

		resetArgs := []string{"reset"}
		if hard {
			resetArgs = append(resetArgs, "--hard")
		} else if soft {
			resetArgs = append(resetArgs, "--soft")
		}
		resetArgs = append(resetArgs, args...)

		out, err := exec.Command("git", resetArgs...).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

// git clean
var gitCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove untracked files",
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		cleanArgs := []string{"clean"}
		if force {
			cleanArgs = append(cleanArgs, "-f")
		}
		if dryRun {
			cleanArgs = append(cleanArgs, "-n")
		}

		out, err := exec.Command("git", cleanArgs...).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

// git archive
var gitArchiveCmd = &cobra.Command{
	Use:   "archive <ref>",
	Short: "Create archive",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		format, _ := cmd.Flags().GetString("format")
		output, _ := cmd.Flags().GetString("output")

		archiveArgs := []string{"archive"}
		if format != "" {
			archiveArgs = append(archiveArgs, "--format="+format)
		}
		if output != "" {
			archiveArgs = append(archiveArgs, "--output="+output)
		}
		archiveArgs = append(archiveArgs, args[0])

		out, err := exec.Command("git", archiveArgs...).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

// git describe
var gitDescribeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe commit",
	Run: func(cmd *cobra.Command, args []string) {
		describeArgs := append([]string{"describe"}, args...)
		out, err := exec.Command("git", describeArgs...).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

// git shortlog
var gitShortlogCmd = &cobra.Command{
	Use:   "shortlog",
	Short: "Summarize log",
	Run: func(cmd *cobra.Command, args []string) {
		shortlogArgs := append([]string{"shortlog"}, args...)
		out, err := exec.Command("git", shortlogArgs...).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

func init() {
	// Fetch
	gitFetchCmd.Flags().String("daemon-url", "", "Daemon URL")

	// Branch
	for _, c := range []*cobra.Command{gitBranchListCmd, gitBranchCreateCmd, gitBranchDeleteCmd} {
		c.Flags().StringP("repo", "r", "", "Repository name (optional, for server-side)")
		c.Flags().String("daemon-url", "", "Daemon URL")
	}
	gitBranchDeleteCmd.Flags().BoolP("force", "f", false, "Force delete")

	// Tag
	gitTagCreateCmd.Flags().BoolP("annotate", "a", false, "Create annotated tag")
	gitTagCreateCmd.Flags().StringP("message", "m", "", "Tag message")

	// Log
	for _, c := range []*cobra.Command{gitLogCmd} {
		c.Flags().StringP("repo", "r", "", "Repository name (optional, for server-side)")
		c.Flags().String("daemon-url", "", "Daemon URL")
		c.Flags().IntP("limit", "n", 10, "Limit")
		c.Flags().Bool("oneline", false, "One-line format")
	}

	// Diff
	for _, c := range []*cobra.Command{gitDiffCmd} {
		c.Flags().StringP("repo", "r", "", "Repository name (optional, for server-side)")
		c.Flags().String("daemon-url", "", "Daemon URL")
	}

	// Show
	for _, c := range []*cobra.Command{gitShowCmd} {
		c.Flags().StringP("repo", "r", "", "Repository name (optional, for server-side)")
		c.Flags().String("daemon-url", "", "Daemon URL")
	}

	// Blame
	for _, c := range []*cobra.Command{gitBlameCmd} {
		c.Flags().StringP("repo", "r", "", "Repository name (optional, for server-side)")
		c.Flags().String("daemon-url", "", "Daemon URL")
	}

	// Grep
	for _, c := range []*cobra.Command{gitGrepCmd} {
		c.Flags().StringP("repo", "r", "", "Repository name (optional, for server-side)")
		c.Flags().String("daemon-url", "", "Daemon URL")
	}

	// Reset
	gitResetCmd.Flags().Bool("hard", false, "Hard reset")
	gitResetCmd.Flags().Bool("soft", false, "Soft reset")

	// Clean
	gitCleanCmd.Flags().BoolP("force", "f", false, "Force")
	gitCleanCmd.Flags().Bool("dry-run", false, "Dry run")

	// Archive
	gitArchiveCmd.Flags().String("format", "", "Archive format")
	gitArchiveCmd.Flags().StringP("output", "o", "", "Output file")

	// Stash
	gitStashCmd.AddCommand(gitStashSaveCmd, gitStashListCmd, gitStashPopCmd, gitStashDropCmd, gitStashShowCmd, gitStashApplyCmd, gitStashBranchCmd, gitStashClearCmd, gitStashStoreCmd)

	// Remote
	gitRemoteCmd.AddCommand(gitRemoteListCmd, gitRemoteAddCmd, gitRemoteRemoveCmd, gitRemoteSetUrlCmd, gitRemoteRenameCmd, gitRemoteShowCmd, gitRemotePruneCmd, gitRemoteUpdateCmd, gitRemoteGetUrlCmd, gitRemoteSetBranchesCmd, gitRemoteSetHeadCmd, gitRemoteGetHeadCmd, gitRemoteSetPushCmd)
	gitRemoteGetUrlCmd.Flags().Bool("push", false, "Query push URL")
	gitRemoteSetBranchesCmd.Flags().Bool("add", false, "Add to existing branches")

	// Branch
	gitBranchCmd.AddCommand(gitBranchListCmd, gitBranchCreateCmd, gitBranchDeleteCmd)

	// Tag
	gitTagCmd.AddCommand(gitTagListCmd, gitTagCreateCmd, gitTagDeleteCmd)

	// Worktree
	gitWorktreeCmd.AddCommand(gitWorktreeListCmd, gitWorktreeAddCmd, gitWorktreeRemoveCmd, gitWorktreeMoveCmd, gitWorktreeLockCmd, gitWorktreeUnlockCmd, gitWorktreePruneCmd, gitWorktreeRepairCmd)
	gitWorktreeRemoveCmd.Flags().BoolP("force", "f", false, "Force removal")
	gitWorktreeLockCmd.Flags().String("reason", "", "Reason for locking")

	// Submodule
	gitSubmoduleCmd.AddCommand(gitSubmoduleListCmd, gitSubmoduleInitCmd, gitSubmoduleUpdateCmd)

	// Add all git subcommands
	gitCmd.AddCommand(
		gitFetchCmd,
		gitRemoteCmd,
		gitBranchCmd,
		gitTagCmd,
		gitLogCmd,
		gitDiffCmd,
		gitStatusCmd,
		gitShowCmd,
		gitMergeCmd,
		gitCherryPickCmd,
		gitRevertCmd,
		gitBlameCmd,
		gitGrepCmd,
		gitStashCmd,
		gitResetCmd,
		gitCleanCmd,
		gitWorktreeCmd,
		gitSubmoduleCmd,
		gitArchiveCmd,
		gitDescribeCmd,
		gitShortlogCmd,
	)

	rootCmd.AddCommand(gitCmd)
}
