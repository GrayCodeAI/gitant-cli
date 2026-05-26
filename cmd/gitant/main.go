package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/GrayCodeAI/gitant-cli/internal/cli"
	"github.com/GrayCodeAI/gitant-cli/internal/version"
	"github.com/go-git/go-git/v6"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gt",
	Short: "Gitant CLI for developers and agents",
	Long:  "Client for Gitant nodes — push, pull, issues, PRs, and more. Pair with gitant-daemon (self-host) or a hosted Gitant URL.",
	Aliases: []string{"gitant"},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new repository",
	Long:  "Initialize a new git repository in the current directory.",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			slog.Error("failed to get current directory", "error", err)
			os.Exit(1)
		}

		if _, err := git.PlainInit(cwd, false); err != nil {
			slog.Error("failed to initialize repository", "error", err)
			os.Exit(1)
		}

		slog.Info("initialized git repository", "path", cwd)
	},
}

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push changes to a remote gitant node",
	Run: func(cmd *cobra.Command, args []string) {
		remote, _ := cmd.Flags().GetString("remote")
		repo, _ := cmd.Flags().GetString("repo")
		if err := cli.Push(".", remote, repo); err != nil {
			slog.Error("push failed", "error", err)
			os.Exit(1)
		}
	},
}

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull changes from a remote gitant node",
	Run: func(cmd *cobra.Command, args []string) {
		remote, _ := cmd.Flags().GetString("remote")
		repo, _ := cmd.Flags().GetString("repo")
		if err := cli.Pull(".", remote, repo); err != nil {
			slog.Error("pull failed", "error", err)
			os.Exit(1)
		}
	},
}

var cloneCmd = &cobra.Command{
	Use:   "clone [repo-id] [directory]",
	Short: "Clone a repository from a gitant node",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		repoID := args[0]
		dir := repoID
		if len(args) > 1 {
			dir = args[1]
		}
		remote, _ := cmd.Flags().GetString("remote")
		if err := cli.Clone(remote, repoID, dir); err != nil {
			slog.Error("clone failed", "error", err)
			os.Exit(1)
		}
	},
}

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup a gitant daemon data directory",
	Long:  "Create a timestamped backup of a gitant data directory (identity, repos, CRDT stores). Use on the machine running gitant-daemon.",
	Run: func(cmd *cobra.Command, args []string) {
		dataDir, _ := cmd.Flags().GetString("data-dir")
		outputDir, _ := cmd.Flags().GetString("output")

		if dataDir == "" {
			home, _ := os.UserHomeDir()
			dataDir = filepath.Join(home, ".gitant")
		}

		timestamp := time.Now().Format("20060102-150405")
		backupDir := filepath.Join(outputDir, "gitant-backup-"+timestamp)

		if err := os.MkdirAll(backupDir, 0755); err != nil {
			slog.Error("failed to create backup directory", "error", err)
			os.Exit(1)
		}

		backupItems := []string{"identity.key", "repos", "data"}
		backedUp := 0
		for _, name := range backupItems {
			src := filepath.Join(dataDir, name)
			dst := filepath.Join(backupDir, name)

			info, err := os.Stat(src)
			if os.IsNotExist(err) {
				continue
			}

			if info.IsDir() {
				if err := copyDir(src, dst); err != nil {
					slog.Warn("failed to backup directory", "path", name, "error", err)
				} else {
					backedUp++
				}
			} else {
				if err := copyFile(src, dst); err != nil {
					slog.Warn("failed to backup file", "path", name, "error", err)
				} else {
					backedUp++
				}
			}
		}

		slog.Info("backup complete", "path", backupDir, "items", backedUp)
	},
}

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore gitant data from backup",
	Long:  "Restore gitant data from a previously created backup directory. Existing data is NOT overwritten.",
	Run: func(cmd *cobra.Command, args []string) {
		dataDir, _ := cmd.Flags().GetString("data-dir")
		inputDir, _ := cmd.Flags().GetString("input")

		if dataDir == "" {
			home, _ := os.UserHomeDir()
			dataDir = filepath.Join(home, ".gitant")
		}

		if _, err := os.Stat(inputDir); os.IsNotExist(err) {
			slog.Error("backup directory not found", "path", inputDir)
			os.Exit(1)
		}

		if err := os.MkdirAll(dataDir, 0755); err != nil {
			slog.Error("failed to create data directory", "error", err)
			os.Exit(1)
		}

		entries, err := os.ReadDir(inputDir)
		if err != nil {
			slog.Error("failed to read backup directory", "error", err)
			os.Exit(1)
		}

		restored := 0
		for _, entry := range entries {
			src := filepath.Join(inputDir, entry.Name())
			dst := filepath.Join(dataDir, entry.Name())

			if _, err := os.Stat(dst); err == nil {
				slog.Info("skipping (already exists)", "path", entry.Name())
				continue
			}

			if entry.IsDir() {
				if err := copyDir(src, dst); err != nil {
					slog.Warn("failed to restore directory", "path", entry.Name(), "error", err)
				} else {
					restored++
				}
			} else {
				if err := copyFile(src, dst); err != nil {
					slog.Warn("failed to restore file", "path", entry.Name(), "error", err)
				} else {
					restored++
				}
			}
		}

		slog.Info("restore complete", "items", restored)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gitant %s (commit %s, built %s)\n", version.Version, version.Commit, version.BuildTime)
	},
}

func copyFile(src, dst string) error {
	in, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, in, info.Mode())
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		return copyFile(path, target)
	})
}

func init() {
	pushCmd.Flags().StringP("remote", "r", "http://localhost:7777", "Remote daemon URL")
	pushCmd.Flags().String("repo", "", "Repository name (required)")
	pushCmd.MarkFlagRequired("repo")
	pullCmd.Flags().StringP("remote", "r", "http://localhost:7777", "Remote daemon URL")
	pullCmd.Flags().String("repo", "", "Repository name (required)")
	pullCmd.MarkFlagRequired("repo")
	cloneCmd.Flags().StringP("remote", "r", "http://localhost:7777", "Remote daemon URL")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(pushCmd)
	rootCmd.AddCommand(pullCmd)
	rootCmd.AddCommand(cloneCmd)
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(restoreCmd)

	backupCmd.Flags().StringP("output", "o", "", "Backup output directory (required)")
	backupCmd.Flags().StringP("data-dir", "d", "", "Data directory (default: ~/.gitant)")
	backupCmd.MarkFlagRequired("output")

	restoreCmd.Flags().StringP("input", "i", "", "Backup directory to restore from (required)")
	restoreCmd.Flags().StringP("data-dir", "d", "", "Data directory (default: ~/.gitant)")
	restoreCmd.MarkFlagRequired("input")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
