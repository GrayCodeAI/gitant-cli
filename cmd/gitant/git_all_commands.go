package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// Additional git commands to make Gitant complete

// git checkout
var gitCheckoutCmd = &cobra.Command{
	Use:   "checkout <branch>",
	Short: "Switch branches",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		checkout := exec.Command("git", "checkout", args[0])
		checkout.Stdout = os.Stdout
		checkout.Stderr = os.Stderr
		if err := checkout.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git switch
var gitSwitchCmd = &cobra.Command{
	Use:   "switch <branch>",
	Short: "Switch branches",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switchCmd := exec.Command("git", "switch", args[0])
		switchCmd.Stdout = os.Stdout
		switchCmd.Stderr = os.Stderr
		if err := switchCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git add
var gitAddCmd = &cobra.Command{
	Use:   "add <pathspec>",
	Short: "Add file contents to the index",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addArgs := append([]string{"add"}, args...)
		add := exec.Command("git", addArgs...)
		add.Stdout = os.Stdout
		add.Stderr = os.Stderr
		if err := add.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git commit
var gitCommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Record changes to the repository",
	Run: func(cmd *cobra.Command, args []string) {
		message, _ := cmd.Flags().GetString("message")
		all, _ := cmd.Flags().GetBool("all")

		commitArgs := []string{"commit"}
		if message != "" {
			commitArgs = append(commitArgs, "-m", message)
		}
		if all {
			commitArgs = append(commitArgs, "-a")
		}
		commitArgs = append(commitArgs, args...)

		commit := exec.Command("git", commitArgs...)
		commit.Stdout = os.Stdout
		commit.Stderr = os.Stderr
		if err := commit.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git rebase
var gitRebaseCmd = &cobra.Command{
	Use:   "rebase <upstream>",
	Short: "Reapply commits on top of another base tip",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		interactive, _ := cmd.Flags().GetBool("interactive")

		rebaseArgs := []string{"rebase"}
		if interactive {
			rebaseArgs = append(rebaseArgs, "-i")
		}
		rebaseArgs = append(rebaseArgs, args[0])

		rebase := exec.Command("git", rebaseArgs...)
		rebase.Stdout = os.Stdout
		rebase.Stderr = os.Stderr
		if err := rebase.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git restore
var gitRestoreCmd = &cobra.Command{
	Use:   "restore <pathspec>",
	Short: "Restore working tree files",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		source, _ := cmd.Flags().GetString("source")
		staged, _ := cmd.Flags().GetBool("staged")

		restoreArgs := []string{"restore"}
		if source != "" {
			restoreArgs = append(restoreArgs, "--source="+source)
		}
		if staged {
			restoreArgs = append(restoreArgs, "--staged")
		}
		restoreArgs = append(restoreArgs, args...)

		restore := exec.Command("git", restoreArgs...)
		restore.Stdout = os.Stdout
		restore.Stderr = os.Stderr
		if err := restore.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git rm
var gitRmCmd = &cobra.Command{
	Use:   "rm <pathspec>",
	Short: "Remove files from the working tree and from the index",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")
		cached, _ := cmd.Flags().GetBool("cached")

		rmArgs := []string{"rm"}
		if force {
			rmArgs = append(rmArgs, "-f")
		}
		if cached {
			rmArgs = append(rmArgs, "--cached")
		}
		rmArgs = append(rmArgs, args...)

		rm := exec.Command("git", rmArgs...)
		rm.Stdout = os.Stdout
		rm.Stderr = os.Stderr
		if err := rm.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git mv
var gitMvCmd = &cobra.Command{
	Use:   "mv <source> <destination>",
	Short: "Move or rename a file",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")

		mvArgs := []string{"mv"}
		if force {
			mvArgs = append(mvArgs, "-f")
		}
		mvArgs = append(mvArgs, args[0], args[1])

		mv := exec.Command("git", mvArgs...)
		mv.Stdout = os.Stdout
		mv.Stderr = os.Stderr
		if err := mv.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git bisect
var gitBisectCmd = &cobra.Command{
	Use:   "bisect",
	Short: "Use binary search to find the commit that introduced a bug",
}

var gitBisectStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start bisecting",
	Run: func(cmd *cobra.Command, args []string) {
		bisect := exec.Command("git", "bisect", "start")
		bisect.Stdout = os.Stdout
		bisect.Stderr = os.Stderr
		if err := bisect.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var gitBisectGoodCmd = &cobra.Command{
	Use:   "good [commit]",
	Short: "Mark commit as good",
	Run: func(cmd *cobra.Command, args []string) {
		bisectArgs := append([]string{"bisect", "good"}, args...)
		bisect := exec.Command("git", bisectArgs...)
		bisect.Stdout = os.Stdout
		bisect.Stderr = os.Stderr
		if err := bisect.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var gitBisectBadCmd = &cobra.Command{
	Use:   "bad [commit]",
	Short: "Mark commit as bad",
	Run: func(cmd *cobra.Command, args []string) {
		bisectArgs := append([]string{"bisect", "bad"}, args...)
		bisect := exec.Command("git", bisectArgs...)
		bisect.Stdout = os.Stdout
		bisect.Stderr = os.Stderr
		if err := bisect.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var gitBisectResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset bisect state",
	Run: func(cmd *cobra.Command, args []string) {
		bisect := exec.Command("git", "bisect", "reset")
		bisect.Stdout = os.Stdout
		bisect.Stderr = os.Stderr
		if err := bisect.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var gitBisectSkipCmd = &cobra.Command{
	Use:   "skip [commit]",
	Short: "Skip a commit in bisecting",
	Run: func(cmd *cobra.Command, args []string) {
		bisectArgs := append([]string{"bisect", "skip"}, args...)
		bisect := exec.Command("git", bisectArgs...)
		bisect.Stdout = os.Stdout
		bisect.Stderr = os.Stderr
		if err := bisect.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var gitBisectLogCmd = &cobra.Command{
	Use:   "log",
	Short: "Show bisect log",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "bisect", "log").Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitBisectReplayCmd = &cobra.Command{
	Use:   "replay <logfile>",
	Short: "Replay bisect log",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "bisect", "replay", args[0]).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitBisectVisualizeCmd = &cobra.Command{
	Use:   "visualize",
	Short: "Show bisect states in gitk",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "bisect", "visualize").Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitBisectTermsCmd = &cobra.Command{
	Use:   "terms",
	Short: "Show terms used by bisect",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "bisect", "terms").Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitBisectRunCmd = &cobra.Command{
	Use:   "run <cmd>...",
	Short: "Run a script to bisect automatically",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bisectArgs := append([]string{"bisect", "run"}, args...)
		bisect := exec.Command("git", bisectArgs...)
		bisect.Stdout = os.Stdout
		bisect.Stderr = os.Stderr
		if err := bisect.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git notes
var gitNotesCmd = &cobra.Command{
	Use:   "notes",
	Short: "Add or inspect object notes",
}

var gitNotesAddCmd = &cobra.Command{
	Use:   "add [object]",
	Short: "Add notes",
	Run: func(cmd *cobra.Command, args []string) {
		notesArgs := append([]string{"notes", "add"}, args...)
		notes := exec.Command("git", notesArgs...)
		notes.Stdout = os.Stdout
		notes.Stderr = os.Stderr
		if err := notes.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var gitNotesShowCmd = &cobra.Command{
	Use:   "show [object]",
	Short: "Show notes",
	Run: func(cmd *cobra.Command, args []string) {
		notesArgs := append([]string{"notes", "show"}, args...)
		notes := exec.Command("git", notesArgs...)
		notes.Stdout = os.Stdout
		notes.Stderr = os.Stderr
		if err := notes.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var gitNotesEditCmd = &cobra.Command{
	Use:   "edit [object]",
	Short: "Edit notes for an object",
	Run: func(cmd *cobra.Command, args []string) {
		notesArgs := append([]string{"notes", "edit"}, args...)
		notes := exec.Command("git", notesArgs...)
		notes.Stdout = os.Stdout
		notes.Stderr = os.Stderr
		if err := notes.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var gitNotesCopyCmd = &cobra.Command{
	Use:   "copy <from-object> <to-object>",
	Short: "Copy notes between objects",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "notes", "copy", args[0], args[1]).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitNotesAppendCmd = &cobra.Command{
	Use:   "append [object]",
	Short: "Append to notes on an object",
	Run: func(cmd *cobra.Command, args []string) {
		notesArgs := append([]string{"notes", "append"}, args...)
		notes := exec.Command("git", notesArgs...)
		notes.Stdout = os.Stdout
		notes.Stderr = os.Stderr
		if err := notes.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var gitNotesRemoveCmd = &cobra.Command{
	Use:   "remove [object]",
	Short: "Remove notes",
	Run: func(cmd *cobra.Command, args []string) {
		notesArgs := append([]string{"notes", "remove"}, args...)
		notes := exec.Command("git", notesArgs...)
		notes.Stdout = os.Stdout
		notes.Stderr = os.Stderr
		if err := notes.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var gitNotesListCmd = &cobra.Command{
	Use:   "list [object]",
	Short: "List notes",
	Run: func(cmd *cobra.Command, args []string) {
		notesArgs := append([]string{"notes", "list"}, args...)
		out, err := exec.Command("git", notesArgs...).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if len(out) == 0 {
			fmt.Println("No notes found")
		} else {
			fmt.Print(string(out))
		}
	},
}

var gitNotesMergeCmd = &cobra.Command{
	Use:   "merge <notes-ref>",
	Short: "Merge notes",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "notes", "merge", args[0]).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitNotesPruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Remove all notes for unreachable objects",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "notes", "prune").Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

var gitNotesGetRefCmd = &cobra.Command{
	Use:   "get-ref",
	Short: "Show the current notes ref",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "notes", "get-ref").Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

// git fsck
var gitFsckCmd = &cobra.Command{
	Use:   "fsck",
	Short: "Verify the connectivity and validity of objects",
	Run: func(cmd *cobra.Command, args []string) {
		fsck := exec.Command("git", "fsck")
		fsck.Stdout = os.Stdout
		fsck.Stderr = os.Stderr
		if err := fsck.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git gc
var gitGCCmd = &cobra.Command{
	Use:   "gc",
	Short: "Cleanup unnecessary files and optimize the local repository",
	Run: func(cmd *cobra.Command, args []string) {
		gc := exec.Command("git", "gc")
		gc.Stdout = os.Stdout
		gc.Stderr = os.Stderr
		if err := gc.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git reflog
var gitReflogCmd = &cobra.Command{
	Use:   "reflog",
	Short: "Show reference logs",
	Run: func(cmd *cobra.Command, args []string) {
		reflog := exec.Command("git", "reflog")
		reflog.Stdout = os.Stdout
		reflog.Stderr = os.Stderr
		if err := reflog.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git repack
var gitRepackCmd = &cobra.Command{
	Use:   "repack",
	Short: "Pack unpacked objects in a repository",
	Run: func(cmd *cobra.Command, args []string) {
		repack := exec.Command("git", "repack")
		repack.Stdout = os.Stdout
		repack.Stderr = os.Stderr
		if err := repack.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git replace
var gitReplaceCmd = &cobra.Command{
	Use:   "replace",
	Short: "Create, list, delete refs to replace objects",
	Run: func(cmd *cobra.Command, args []string) {
		replaceArgs := append([]string{"replace"}, args...)
		replace := exec.Command("git", replaceArgs...)
		replace.Stdout = os.Stdout
		replace.Stderr = os.Stderr
		if err := replace.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git sparse-checkout
var gitSparseCheckoutCmd = &cobra.Command{
	Use:   "sparse-checkout",
	Short: "Reduce your working tree to a subset of tracked files",
	Run: func(cmd *cobra.Command, args []string) {
		sparseArgs := append([]string{"sparse-checkout"}, args...)
		sparse := exec.Command("git", sparseArgs...)
		sparse.Stdout = os.Stdout
		sparse.Stderr = os.Stderr
		if err := sparse.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git config
var gitConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Get and set repository or global options",
	Run: func(cmd *cobra.Command, args []string) {
		configArgs := append([]string{"config"}, args...)
		config := exec.Command("git", configArgs...)
		config.Stdout = os.Stdout
		config.Stderr = os.Stderr
		if err := config.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git format-patch
var gitFormatPatchCmd = &cobra.Command{
	Use:   "format-patch",
	Short: "Prepare patches for e-mail submission",
	Run: func(cmd *cobra.Command, args []string) {
		formatArgs := append([]string{"format-patch"}, args...)
		format := exec.Command("git", formatArgs...)
		format.Stdout = os.Stdout
		format.Stderr = os.Stderr
		if err := format.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git send-email
var gitSendEmailCmd = &cobra.Command{
	Use:   "send-email",
	Short: "Send a collection of patches as emails",
	Run: func(cmd *cobra.Command, args []string) {
		sendArgs := append([]string{"send-email"}, args...)
		send := exec.Command("git", sendArgs...)
		send.Stdout = os.Stdout
		send.Stderr = os.Stderr
		if err := send.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git request-pull
var gitRequestPullCmd = &cobra.Command{
	Use:   "request-pull",
	Short: "Gentle request to pull changes",
	Run: func(cmd *cobra.Command, args []string) {
		requestArgs := append([]string{"request-pull"}, args...)
		request := exec.Command("git", requestArgs...)
		request.Stdout = os.Stdout
		request.Stderr = os.Stderr
		if err := request.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git rerere
var gitRerereCmd = &cobra.Command{
	Use:   "rerere",
	Short: "Reuse recorded resolution of conflicted merges",
	Run: func(cmd *cobra.Command, args []string) {
		rerereArgs := append([]string{"rerere"}, args...)
		rerere := exec.Command("git", rerereArgs...)
		rerere.Stdout = os.Stdout
		rerere.Stderr = os.Stderr
		if err := rerere.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git difftool
var gitDifftoolCmd = &cobra.Command{
	Use:   "difftool",
	Short: "Show changes using common diff tools",
	Run: func(cmd *cobra.Command, args []string) {
		difftoolArgs := append([]string{"difftool"}, args...)
		difftool := exec.Command("git", difftoolArgs...)
		difftool.Stdout = os.Stdout
		difftool.Stderr = os.Stderr
		if err := difftool.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git mergetool
var gitMergetoolCmd = &cobra.Command{
	Use:   "mergetool",
	Short: "Run merge conflict resolution tools",
	Run: func(cmd *cobra.Command, args []string) {
		mergetool := exec.Command("git", "mergetool")
		mergetool.Stdout = os.Stdout
		mergetool.Stderr = os.Stderr
		if err := mergetool.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git bugreport
var gitBugreportCmd = &cobra.Command{
	Use:   "bugreport",
	Short: "Collect information for user bug reports",
	Run: func(cmd *cobra.Command, args []string) {
		bugreport := exec.Command("git", "bugreport")
		bugreport.Stdout = os.Stdout
		bugreport.Stderr = os.Stderr
		if err := bugreport.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git bundle
var gitBundleCmd = &cobra.Command{
	Use:   "bundle",
	Short: "Move objects and refs by archive",
	Run: func(cmd *cobra.Command, args []string) {
		bundleArgs := append([]string{"bundle"}, args...)
		bundle := exec.Command("git", bundleArgs...)
		bundle.Stdout = os.Stdout
		bundle.Stderr = os.Stderr
		if err := bundle.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git prune
var gitPruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Prune all unreachable objects from the object database",
	Run: func(cmd *cobra.Command, args []string) {
		prune := exec.Command("git", "prune")
		prune.Stdout = os.Stdout
		prune.Stderr = os.Stderr
		if err := prune.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git count-objects
var gitCountObjectsCmd = &cobra.Command{
	Use:   "count-objects",
	Short: "Count unpacked number of objects and their disk consumption",
	Run: func(cmd *cobra.Command, args []string) {
		count := exec.Command("git", "count-objects", "-v")
		count.Stdout = os.Stdout
		count.Stderr = os.Stderr
		if err := count.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git verify-pack
var gitVerifyPackCmd = &cobra.Command{
	Use:   "verify-pack",
	Short: "Validate packed Git archive files",
	Run: func(cmd *cobra.Command, args []string) {
		verifyArgs := append([]string{"verify-pack"}, args...)
		verify := exec.Command("git", verifyArgs...)
		verify.Stdout = os.Stdout
		verify.Stderr = os.Stderr
		if err := verify.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git show-ref
var gitShowRefCmd = &cobra.Command{
	Use:   "show-ref",
	Short: "List references in a local repository",
	Run: func(cmd *cobra.Command, args []string) {
		showRef := exec.Command("git", "show-ref")
		showRef.Stdout = os.Stdout
		showRef.Stderr = os.Stderr
		if err := showRef.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git symbolic-ref
var gitSymbolicRefCmd = &cobra.Command{
	Use:   "symbolic-ref",
	Short: "Read, modify and delete symbolic refs",
	Run: func(cmd *cobra.Command, args []string) {
		symbolicArgs := append([]string{"symbolic-ref"}, args...)
		symbolic := exec.Command("git", symbolicArgs...)
		symbolic.Stdout = os.Stdout
		symbolic.Stderr = os.Stderr
		if err := symbolic.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git update-ref
var gitUpdateRefCmd = &cobra.Command{
	Use:   "update-ref",
	Short: "Update the object name stored in a ref safely",
	Run: func(cmd *cobra.Command, args []string) {
		updateArgs := append([]string{"update-ref"}, args...)
		update := exec.Command("git", updateArgs...)
		update.Stdout = os.Stdout
		update.Stderr = os.Stderr
		if err := update.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git rev-parse
var gitRevParseCmd = &cobra.Command{
	Use:   "rev-parse",
	Short: "Pick out and massage parameters",
	Run: func(cmd *cobra.Command, args []string) {
		revParseArgs := append([]string{"rev-parse"}, args...)
		revParse := exec.Command("git", revParseArgs...)
		revParse.Stdout = os.Stdout
		revParse.Stderr = os.Stderr
		if err := revParse.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git rev-list
var gitRevListCmd = &cobra.Command{
	Use:   "rev-list",
	Short: "Lists commit objects in reverse chronological order",
	Run: func(cmd *cobra.Command, args []string) {
		revListArgs := append([]string{"rev-list"}, args...)
		revList := exec.Command("git", revListArgs...)
		revList.Stdout = os.Stdout
		revList.Stderr = os.Stderr
		if err := revList.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git ls-files
var gitLsFilesCmd = &cobra.Command{
	Use:   "ls-files",
	Short: "Show information about files in the index and working tree",
	Run: func(cmd *cobra.Command, args []string) {
		lsFilesArgs := append([]string{"ls-files"}, args...)
		lsFiles := exec.Command("git", lsFilesArgs...)
		lsFiles.Stdout = os.Stdout
		lsFiles.Stderr = os.Stderr
		if err := lsFiles.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git ls-remote
var gitLsRemoteCmd = &cobra.Command{
	Use:   "ls-remote",
	Short: "List references in a remote repository",
	Run: func(cmd *cobra.Command, args []string) {
		lsRemoteArgs := append([]string{"ls-remote"}, args...)
		lsRemote := exec.Command("git", lsRemoteArgs...)
		lsRemote.Stdout = os.Stdout
		lsRemote.Stderr = os.Stderr
		if err := lsRemote.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git ls-tree
var gitLsTreeCmd = &cobra.Command{
	Use:   "ls-tree",
	Short: "List the contents of a tree object",
	Run: func(cmd *cobra.Command, args []string) {
		lsTreeArgs := append([]string{"ls-tree"}, args...)
		lsTree := exec.Command("git", lsTreeArgs...)
		lsTree.Stdout = os.Stdout
		lsTree.Stderr = os.Stderr
		if err := lsTree.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git read-tree
var gitReadTreeCmd = &cobra.Command{
	Use:   "read-tree",
	Short: "Reads tree information into the index",
	Run: func(cmd *cobra.Command, args []string) {
		readTreeArgs := append([]string{"read-tree"}, args...)
		readTree := exec.Command("git", readTreeArgs...)
		readTree.Stdout = os.Stdout
		readTree.Stderr = os.Stderr
		if err := readTree.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git write-tree
var gitWriteTreeCmd = &cobra.Command{
	Use:   "write-tree",
	Short: "Create a tree object from the current index",
	Run: func(cmd *cobra.Command, args []string) {
		writeTree := exec.Command("git", "write-tree")
		writeTree.Stdout = os.Stdout
		writeTree.Stderr = os.Stderr
		if err := writeTree.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git commit-tree
var gitCommitTreeCmd = &cobra.Command{
	Use:   "commit-tree",
	Short: "Create a new commit object",
	Run: func(cmd *cobra.Command, args []string) {
		commitTreeArgs := append([]string{"commit-tree"}, args...)
		commitTree := exec.Command("git", commitTreeArgs...)
		commitTree.Stdout = os.Stdout
		commitTree.Stderr = os.Stderr
		if err := commitTree.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git hash-object
var gitHashObjectCmd = &cobra.Command{
	Use:   "hash-object",
	Short: "Compute object ID and optionally creates a blob from a file",
	Run: func(cmd *cobra.Command, args []string) {
		hashArgs := append([]string{"hash-object"}, args...)
		hash := exec.Command("git", hashArgs...)
		hash.Stdout = os.Stdout
		hash.Stderr = os.Stderr
		if err := hash.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git cat-file
var gitCatFileCmd = &cobra.Command{
	Use:   "cat-file",
	Short: "Provide content or type and size information for repository objects",
	Run: func(cmd *cobra.Command, args []string) {
		catFileArgs := append([]string{"cat-file"}, args...)
		catFile := exec.Command("git", catFileArgs...)
		catFile.Stdout = os.Stdout
		catFile.Stderr = os.Stderr
		if err := catFile.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git diff-files
var gitDiffFilesCmd = &cobra.Command{
	Use:   "diff-files",
	Short: "Compares files in the working tree and the index",
	Run: func(cmd *cobra.Command, args []string) {
		diffFiles := exec.Command("git", "diff-files")
		diffFiles.Stdout = os.Stdout
		diffFiles.Stderr = os.Stderr
		if err := diffFiles.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git diff-index
var gitDiffIndexCmd = &cobra.Command{
	Use:   "diff-index",
	Short: "Compare a tree to the working tree or index",
	Run: func(cmd *cobra.Command, args []string) {
		diffIndexArgs := append([]string{"diff-index"}, args...)
		diffIndex := exec.Command("git", diffIndexArgs...)
		diffIndex.Stdout = os.Stdout
		diffIndex.Stderr = os.Stderr
		if err := diffIndex.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git diff-tree
var gitDiffTreeCmd = &cobra.Command{
	Use:   "diff-tree",
	Short: "Compares the content and mode of blobs found via two tree objects",
	Run: func(cmd *cobra.Command, args []string) {
		diffTreeArgs := append([]string{"diff-tree"}, args...)
		diffTree := exec.Command("git", diffTreeArgs...)
		diffTree.Stdout = os.Stdout
		diffTree.Stderr = os.Stderr
		if err := diffTree.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git apply
var gitApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a patch to files and/or to the index",
	Run: func(cmd *cobra.Command, args []string) {
		applyArgs := append([]string{"apply"}, args...)
		apply := exec.Command("git", applyArgs...)
		apply.Stdout = os.Stdout
		apply.Stderr = os.Stderr
		if err := apply.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git am
var gitAmCmd = &cobra.Command{
	Use:   "am",
	Short: "Apply a series of patches from a mailbox",
	Run: func(cmd *cobra.Command, args []string) {
		amArgs := append([]string{"am"}, args...)
		am := exec.Command("git", amArgs...)
		am.Stdout = os.Stdout
		am.Stderr = os.Stderr
		if err := am.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	// Flags for commit
	gitCommitCmd.Flags().StringP("message", "m", "", "Commit message")
	gitCommitCmd.Flags().BoolP("all", "a", false, "Stage all modified files")

	// Flags for rebase
	gitRebaseCmd.Flags().BoolP("interactive", "i", false, "Interactive rebase")

	// Flags for restore
	gitRestoreCmd.Flags().StringP("source", "s", "", "Source tree")
	gitRestoreCmd.Flags().Bool("staged", false, "Restore staged files")

	// Flags for rm
	gitRmCmd.Flags().BoolP("force", "f", false, "Force removal")
	gitRmCmd.Flags().Bool("cached", false, "Only remove from index")

	// Flags for mv
	gitMvCmd.Flags().BoolP("force", "f", false, "Force move")

	// Bisect subcommands
	gitBisectCmd.AddCommand(gitBisectStartCmd, gitBisectGoodCmd, gitBisectBadCmd, gitBisectResetCmd, gitBisectSkipCmd, gitBisectLogCmd, gitBisectReplayCmd, gitBisectVisualizeCmd, gitBisectTermsCmd, gitBisectRunCmd)

	// Notes subcommands
	gitNotesCmd.AddCommand(gitNotesAddCmd, gitNotesShowCmd, gitNotesEditCmd, gitNotesCopyCmd, gitNotesAppendCmd, gitNotesRemoveCmd, gitNotesListCmd, gitNotesMergeCmd, gitNotesPruneCmd, gitNotesGetRefCmd)

	// Add all new commands to gitCmd
	gitCmd.AddCommand(
		gitCheckoutCmd,
		gitSwitchCmd,
		gitAddCmd,
		gitCommitCmd,
		gitRebaseCmd,
		gitRestoreCmd,
		gitRmCmd,
		gitMvCmd,
		gitBisectCmd,
		gitNotesCmd,
		gitFsckCmd,
		gitGCCmd,
		gitReflogCmd,
		gitRepackCmd,
		gitReplaceCmd,
		gitSparseCheckoutCmd,
		gitConfigCmd,
		gitFormatPatchCmd,
		gitSendEmailCmd,
		gitRequestPullCmd,
		gitRerereCmd,
		gitDifftoolCmd,
		gitMergetoolCmd,
		gitBugreportCmd,
		gitBundleCmd,
		gitPruneCmd,
		gitCountObjectsCmd,
		gitVerifyPackCmd,
		gitShowRefCmd,
		gitSymbolicRefCmd,
		gitUpdateRefCmd,
		gitRevParseCmd,
		gitRevListCmd,
		gitLsFilesCmd,
		gitLsRemoteCmd,
		gitLsTreeCmd,
		gitReadTreeCmd,
		gitWriteTreeCmd,
		gitCommitTreeCmd,
		gitHashObjectCmd,
		gitCatFileCmd,
		gitDiffFilesCmd,
		gitDiffIndexCmd,
		gitDiffTreeCmd,
		gitApplyCmd,
		gitAmCmd,
	)
}
