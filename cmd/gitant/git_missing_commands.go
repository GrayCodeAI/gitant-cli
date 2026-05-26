package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// Missing porcelain commands

// git whatchanged
var gitWhatchangedCmd = &cobra.Command{
	Use:   "whatchanged",
	Short: "Show logs with difference each commit introduces",
	Run: func(cmd *cobra.Command, args []string) {
		whatchangedArgs := append([]string{"whatchanged"}, args...)
		out, err := exec.Command("git", whatchangedArgs...).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

// git annotate
var gitAnnotateCmd = &cobra.Command{
	Use:   "annotate <file>",
	Short: "Annotate file lines with commit info (alias for blame)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "annotate", args[0]).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

// git cherry
var gitCherryCmd = &cobra.Command{
	Use:   "cherry <upstream> [head]",
	Short: "Find commits yet to be applied to upstream",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		cherryArgs := append([]string{"cherry"}, args...)
		out, err := exec.Command("git", cherryArgs...).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

// git stage (alias for git add)
var gitStageCmd = &cobra.Command{
	Use:   "stage <pathspec>",
	Short: "Add file contents to the index (alias for add)",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stageArgs := append([]string{"stage"}, args...)
		stage := exec.Command("git", stageArgs...)
		stage.Stdout = os.Stdout
		stage.Stderr = os.Stderr
		if err := stage.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git maintenance
var gitMaintenanceCmd = &cobra.Command{
	Use:   "maintenance",
	Short: "Run tasks to optimize Git repository data",
	Run: func(cmd *cobra.Command, args []string) {
		maintenanceArgs := append([]string{"maintenance"}, args...)
		maintenance := exec.Command("git", maintenanceArgs...)
		maintenance.Stdout = os.Stdout
		maintenance.Stderr = os.Stderr
		if err := maintenance.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git credential
var gitCredentialCmd = &cobra.Command{
	Use:   "credential",
	Short: "Retrieve and store user credentials",
	Run: func(cmd *cobra.Command, args []string) {
		credentialArgs := append([]string{"credential"}, args...)
		credential := exec.Command("git", credentialArgs...)
		credential.Stdout = os.Stdout
		credential.Stderr = os.Stderr
		if err := credential.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git hook
var gitHookCmd = &cobra.Command{
	Use:   "hook",
	Short: "Run git hooks",
	Run: func(cmd *cobra.Command, args []string) {
		hookArgs := append([]string{"hook"}, args...)
		hook := exec.Command("git", hookArgs...)
		hook.Stdout = os.Stdout
		hook.Stderr = os.Stderr
		if err := hook.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// Missing plumbing commands

// git for-each-ref
var gitForEachRefCmd = &cobra.Command{
	Use:   "for-each-ref",
	Short: "Output information on each ref",
	Run: func(cmd *cobra.Command, args []string) {
		forEachRefArgs := append([]string{"for-each-ref"}, args...)
		out, err := exec.Command("git", forEachRefArgs...).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

// git merge-base
var gitMergeBaseCmd = &cobra.Command{
	Use:   "merge-base <commit> <commit>",
	Short: "Find as good common ancestors as possible for a merge",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		mergeBaseArgs := append([]string{"merge-base"}, args...)
		out, err := exec.Command("git", mergeBaseArgs...).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

// git merge-file
var gitMergeFileCmd = &cobra.Command{
	Use:   "merge-file <current> <base> <other>",
	Short: "Run a three-way file merge",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "merge-file", args[0], args[1], args[2]).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

// git merge-tree
var gitMergeTreeCmd = &cobra.Command{
	Use:   "merge-tree <base-tree> <branch1> <branch2>",
	Short: "Show three-way merge without touching index",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "merge-tree", args[0], args[1], args[2]).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

// git merge-index
var gitMergeIndexCmd = &cobra.Command{
	Use:   "merge-index <cmd>",
	Short: "Run a merge for files needing merging",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mergeIndexArgs := append([]string{"merge-index"}, args...)
		out, err := exec.Command("git", mergeIndexArgs...).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

// git pack-objects
var gitPackObjectsCmd = &cobra.Command{
	Use:   "pack-objects",
	Short: "Create a packed archive of objects",
	Run: func(cmd *cobra.Command, args []string) {
		packArgs := append([]string{"pack-objects"}, args...)
		pack := exec.Command("git", packArgs...)
		pack.Stdout = os.Stdout
		pack.Stderr = os.Stderr
		if err := pack.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git index-pack
var gitIndexPackCmd = &cobra.Command{
	Use:   "index-pack",
	Short: "Build pack index file for an existing packed archive",
	Run: func(cmd *cobra.Command, args []string) {
		indexArgs := append([]string{"index-pack"}, args...)
		index := exec.Command("git", indexArgs...)
		index.Stdout = os.Stdout
		index.Stderr = os.Stderr
		if err := index.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git unpack-objects
var gitUnpackObjectsCmd = &cobra.Command{
	Use:   "unpack-objects",
	Short: "Unpack objects from a packed archive",
	Run: func(cmd *cobra.Command, args []string) {
		unpack := exec.Command("git", "unpack-objects")
		unpack.Stdout = os.Stdout
		unpack.Stderr = os.Stderr
		if err := unpack.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git update-index
var gitUpdateIndexCmd = &cobra.Command{
	Use:   "update-index",
	Short: "Register file contents in the working tree to the index",
	Run: func(cmd *cobra.Command, args []string) {
		updateIndexArgs := append([]string{"update-index"}, args...)
		updateIndex := exec.Command("git", updateIndexArgs...)
		updateIndex.Stdout = os.Stdout
		updateIndex.Stderr = os.Stderr
		if err := updateIndex.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git checkout-index
var gitCheckoutIndexCmd = &cobra.Command{
	Use:   "checkout-index",
	Short: "Copy files from the index to the working directory",
	Run: func(cmd *cobra.Command, args []string) {
		checkoutIndexArgs := append([]string{"checkout-index"}, args...)
		checkoutIndex := exec.Command("git", checkoutIndexArgs...)
		checkoutIndex.Stdout = os.Stdout
		checkoutIndex.Stderr = os.Stderr
		if err := checkoutIndex.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git commit-graph
var gitCommitGraphCmd = &cobra.Command{
	Use:   "commit-graph",
	Short: "Write and verify Git commit-graph files",
	Run: func(cmd *cobra.Command, args []string) {
		commitGraphArgs := append([]string{"commit-graph"}, args...)
		commitGraph := exec.Command("git", commitGraphArgs...)
		commitGraph.Stdout = os.Stdout
		commitGraph.Stderr = os.Stderr
		if err := commitGraph.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git multi-pack-index
var gitMultiPackIndexCmd = &cobra.Command{
	Use:   "multi-pack-index",
	Short: "Write and verify multi-pack-indexes",
	Run: func(cmd *cobra.Command, args []string) {
		multiPackArgs := append([]string{"multi-pack-index"}, args...)
		multiPack := exec.Command("git", multiPackArgs...)
		multiPack.Stdout = os.Stdout
		multiPack.Stderr = os.Stderr
		if err := multiPack.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git interpret-trailers
var gitInterpretTrailersCmd = &cobra.Command{
	Use:   "interpret-trailers",
	Short: "Add or parse structured information in commit messages",
	Run: func(cmd *cobra.Command, args []string) {
		trailersArgs := append([]string{"interpret-trailers"}, args...)
		trailers := exec.Command("git", trailersArgs...)
		trailers.Stdout = os.Stdout
		trailers.Stderr = os.Stderr
		if err := trailers.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git patch-id
var gitPatchIdCmd = &cobra.Command{
	Use:   "patch-id",
	Short: "Compute unique ID for a patch",
	Run: func(cmd *cobra.Command, args []string) {
		patchId := exec.Command("git", "patch-id")
		patchId.Stdout = os.Stdout
		patchId.Stderr = os.Stderr
		if err := patchId.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git verify-commit
var gitVerifyCommitCmd = &cobra.Command{
	Use:   "verify-commit <commit>",
	Short: "Check the GPG signature of commits",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		verifyArgs := append([]string{"verify-commit"}, args...)
		verify := exec.Command("git", verifyArgs...)
		verify.Stdout = os.Stdout
		verify.Stderr = os.Stderr
		if err := verify.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git verify-tag
var gitVerifyTagCmd = &cobra.Command{
	Use:   "verify-tag <tag>",
	Short: "Check the GPG signature of tags",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		verifyArgs := append([]string{"verify-tag"}, args...)
		verify := exec.Command("git", verifyArgs...)
		verify.Stdout = os.Stdout
		verify.Stderr = os.Stderr
		if err := verify.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git mktag
var gitMktagCmd = &cobra.Command{
	Use:   "mktag",
	Short: "Creates a tag object with extra validation",
	Run: func(cmd *cobra.Command, args []string) {
		mktag := exec.Command("git", "mktag")
		mktag.Stdout = os.Stdout
		mktag.Stderr = os.Stderr
		if err := mktag.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git mktree
var gitMktreeCmd = &cobra.Command{
	Use:   "mktree",
	Short: "Build a tree-object from stdin",
	Run: func(cmd *cobra.Command, args []string) {
		mktree := exec.Command("git", "mktree")
		mktree.Stdout = os.Stdout
		mktree.Stderr = os.Stderr
		if err := mktree.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git pack-refs
var gitPackRefsCmd = &cobra.Command{
	Use:   "pack-refs",
	Short: "Pack heads and tags for efficient storage",
	Run: func(cmd *cobra.Command, args []string) {
		packRefs := exec.Command("git", "pack-refs")
		packRefs.Stdout = os.Stdout
		packRefs.Stderr = os.Stderr
		if err := packRefs.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git prune-packed
var gitPrunePackedCmd = &cobra.Command{
	Use:   "prune-packed",
	Short: "Remove extra objects from packed object database",
	Run: func(cmd *cobra.Command, args []string) {
		prunePacked := exec.Command("git", "prune-packed")
		prunePacked.Stdout = os.Stdout
		prunePacked.Stderr = os.Stderr
		if err := prunePacked.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git stripspace
var gitStripspaceCmd = &cobra.Command{
	Use:   "stripspace",
	Short: "Filter out empty lines",
	Run: func(cmd *cobra.Command, args []string) {
		stripspace := exec.Command("git", "stripspace")
		stripspace.Stdout = os.Stdout
		stripspace.Stderr = os.Stderr
		if err := stripspace.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git column
var gitColumnCmd = &cobra.Command{
	Use:   "column",
	Short: "Display data in columns",
	Run: func(cmd *cobra.Command, args []string) {
		columnArgs := append([]string{"column"}, args...)
		column := exec.Command("git", columnArgs...)
		column.Stdout = os.Stdout
		column.Stderr = os.Stderr
		if err := column.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git check-attr
var gitCheckAttrCmd = &cobra.Command{
	Use:   "check-attr <attr> [--] <pathname>...",
	Short: "Display gitattributes information",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		checkAttrArgs := append([]string{"check-attr"}, args...)
		out, err := exec.Command("git", checkAttrArgs...).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

// git check-ignore
var gitCheckIgnoreCmd = &cobra.Command{
	Use:   "check-ignore <pathname>...",
	Short: "Debug gitignore / exclude files",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		checkIgnoreArgs := append([]string{"check-ignore"}, args...)
		out, err := exec.Command("git", checkIgnoreArgs...).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

// git check-mailmap
var gitCheckMailmapCmd = &cobra.Command{
	Use:   "check-mailmap <contact>",
	Short: "Show canonical names and email addresses of contacts",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		checkMailmapArgs := append([]string{"check-mailmap"}, args...)
		out, err := exec.Command("git", checkMailmapArgs...).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

// git check-ref-format
var gitCheckRefFormatCmd = &cobra.Command{
	Use:   "check-ref-format <refname>",
	Short: "Ensures that a reference name is well formed",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "check-ref-format", args[0]).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

// git var
var gitVarCmd = &cobra.Command{
	Use:   "var <variable>",
	Short: "Show a Git logical variable",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "var", args[0]).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

// git get-tar-commit-id
var gitGetTarCommitIdCmd = &cobra.Command{
	Use:   "get-tar-commit-id",
	Short: "Extract commit ID from an archive created by git-archive",
	Run: func(cmd *cobra.Command, args []string) {
		getTar := exec.Command("git", "get-tar-commit-id")
		getTar.Stdout = os.Stdout
		getTar.Stderr = os.Stderr
		if err := getTar.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git show-index
var gitShowIndexCmd = &cobra.Command{
	Use:   "show-index",
	Short: "Show packed archive index",
	Run: func(cmd *cobra.Command, args []string) {
		showIndex := exec.Command("git", "show-index")
		showIndex.Stdout = os.Stdout
		showIndex.Stderr = os.Stderr
		if err := showIndex.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git unpack-file
var gitUnpackFileCmd = &cobra.Command{
	Use:   "unpack-file <blob>",
	Short: "Creates a temporary file with a blob's contents",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "unpack-file", args[0]).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

// git update-server-info
var gitUpdateServerInfoCmd = &cobra.Command{
	Use:   "update-server-info",
	Short: "Update auxiliary info file for dumb transport",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("git", "update-server-info").Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(string(out))
	},
}

// git fmt-merge-msg
var gitFmtMergeMsgCmd = &cobra.Command{
	Use:   "fmt-merge-msg",
	Short: "Produce a merge commit message",
	Run: func(cmd *cobra.Command, args []string) {
		fmtMergeMsg := exec.Command("git", "fmt-merge-msg")
		fmtMergeMsg.Stdout = os.Stdout
		fmtMergeMsg.Stderr = os.Stderr
		if err := fmtMergeMsg.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git mailsplit
var gitMailsplitCmd = &cobra.Command{
	Use:   "mailsplit",
	Short: "Simple UNIX mbox splitter",
	Run: func(cmd *cobra.Command, args []string) {
		mailsplitArgs := append([]string{"mailsplit"}, args...)
		mailsplit := exec.Command("git", mailsplitArgs...)
		mailsplit.Stdout = os.Stdout
		mailsplit.Stderr = os.Stderr
		if err := mailsplit.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git mailinfo
var gitMailinfoCmd = &cobra.Command{
	Use:   "mailinfo",
	Short: "Extracts patch and authorship from a single e-mail message",
	Run: func(cmd *cobra.Command, args []string) {
		mailinfo := exec.Command("git", "mailinfo")
		mailinfo.Stdout = os.Stdout
		mailinfo.Stderr = os.Stderr
		if err := mailinfo.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git imap-send
var gitImapSendCmd = &cobra.Command{
	Use:   "imap-send",
	Short: "Send a collection of patches from stdin to an IMAP folder",
	Run: func(cmd *cobra.Command, args []string) {
		imapSend := exec.Command("git", "imap-send")
		imapSend.Stdout = os.Stdout
		imapSend.Stderr = os.Stderr
		if err := imapSend.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git fast-import
var gitFastImportCmd = &cobra.Command{
	Use:   "fast-import",
	Short: "Backend for fast Git data importers",
	Run: func(cmd *cobra.Command, args []string) {
		fastImport := exec.Command("git", "fast-import")
		fastImport.Stdout = os.Stdout
		fastImport.Stderr = os.Stderr
		if err := fastImport.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git fast-export
var gitFastExportCmd = &cobra.Command{
	Use:   "fast-export",
	Short: "Git data exporter",
	Run: func(cmd *cobra.Command, args []string) {
		fastExportArgs := append([]string{"fast-export"}, args...)
		fastExport := exec.Command("git", fastExportArgs...)
		fastExport.Stdout = os.Stdout
		fastExport.Stderr = os.Stderr
		if err := fastExport.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git filter-branch
var gitFilterBranchCmd = &cobra.Command{
	Use:   "filter-branch",
	Short: "Rewrite branches (deprecated; use git-filter-repo instead)",
	Run: func(cmd *cobra.Command, args []string) {
		filterArgs := append([]string{"filter-branch"}, args...)
		filter := exec.Command("git", filterArgs...)
		filter.Stdout = os.Stdout
		filter.Stderr = os.Stderr
		if err := filter.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git daemon
var gitDaemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "A really simple server for Git repositories",
	Run: func(cmd *cobra.Command, args []string) {
		daemonArgs := append([]string{"daemon"}, args...)
		daemon := exec.Command("git", daemonArgs...)
		daemon.Stdout = os.Stdout
		daemon.Stderr = os.Stderr
		if err := daemon.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git shell
var gitShellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Restricted login shell for GIT-only SSH access",
	Run: func(cmd *cobra.Command, args []string) {
		shellArgs := append([]string{"shell"}, args...)
		shell := exec.Command("git", shellArgs...)
		shell.Stdout = os.Stdout
		shell.Stderr = os.Stderr
		if err := shell.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git receive-pack
var gitReceivePackCmd = &cobra.Command{
	Use:   "receive-pack",
	Short: "Receive what is pushed into the repository",
	Run: func(cmd *cobra.Command, args []string) {
		receivePackArgs := append([]string{"receive-pack"}, args...)
		receivePack := exec.Command("git", receivePackArgs...)
		receivePack.Stdout = os.Stdout
		receivePack.Stderr = os.Stderr
		if err := receivePack.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git send-pack
var gitSendPackCmd = &cobra.Command{
	Use:   "send-pack",
	Short: "Push objects packed over git protocol",
	Run: func(cmd *cobra.Command, args []string) {
		sendPackArgs := append([]string{"send-pack"}, args...)
		sendPack := exec.Command("git", sendPackArgs...)
		sendPack.Stdout = os.Stdout
		sendPack.Stderr = os.Stderr
		if err := sendPack.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git upload-archive
var gitUploadArchiveCmd = &cobra.Command{
	Use:   "upload-archive",
	Short: "Send archive back to git-archive",
	Run: func(cmd *cobra.Command, args []string) {
		uploadArchive := exec.Command("git", "upload-archive")
		uploadArchive.Stdout = os.Stdout
		uploadArchive.Stderr = os.Stderr
		if err := uploadArchive.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git upload-pack
var gitUploadPackCmd = &cobra.Command{
	Use:   "upload-pack",
	Short: "Send objects packed back to git-fetch-pack",
	Run: func(cmd *cobra.Command, args []string) {
		uploadPackArgs := append([]string{"upload-pack"}, args...)
		uploadPack := exec.Command("git", uploadPackArgs...)
		uploadPack.Stdout = os.Stdout
		uploadPack.Stderr = os.Stderr
		if err := uploadPack.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git credential-cache
var gitCredentialCacheCmd = &cobra.Command{
	Use:   "credential-cache",
	Short: "Helper to temporarily store passwords in memory",
	Run: func(cmd *cobra.Command, args []string) {
		cacheArgs := append([]string{"credential-cache"}, args...)
		cache := exec.Command("git", cacheArgs...)
		cache.Stdout = os.Stdout
		cache.Stderr = os.Stderr
		if err := cache.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git credential-store
var gitCredentialStoreCmd = &cobra.Command{
	Use:   "credential-store",
	Short: "Helper to store credentials on disk",
	Run: func(cmd *cobra.Command, args []string) {
		storeArgs := append([]string{"credential-store"}, args...)
		store := exec.Command("git", storeArgs...)
		store.Stdout = os.Stdout
		store.Stderr = os.Stderr
		if err := store.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git subtree
var gitSubtreeCmd = &cobra.Command{
	Use:   "subtree",
	Short: "Merge subtrees together and split repository into subtrees",
	Run: func(cmd *cobra.Command, args []string) {
		subtreeArgs := append([]string{"subtree"}, args...)
		subtree := exec.Command("git", subtreeArgs...)
		subtree.Stdout = os.Stdout
		subtree.Stderr = os.Stderr
		if err := subtree.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git p4
var gitP4Cmd = &cobra.Command{
	Use:   "p4",
	Short: "Import from and submit to Perforce repositories",
	Run: func(cmd *cobra.Command, args []string) {
		p4Args := append([]string{"p4"}, args...)
		p4 := exec.Command("git", p4Args...)
		p4.Stdout = os.Stdout
		p4.Stderr = os.Stderr
		if err := p4.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git svn
var gitSvnCmd = &cobra.Command{
	Use:   "svn",
	Short: "Bidirectional operation between Subversion and Git",
	Run: func(cmd *cobra.Command, args []string) {
		svnArgs := append([]string{"svn"}, args...)
		svn := exec.Command("git", svnArgs...)
		svn.Stdout = os.Stdout
		svn.Stderr = os.Stderr
		if err := svn.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// git quiltimport
var gitQuiltimportCmd = &cobra.Command{
	Use:   "quiltimport",
	Short: "Applies a quilt patchset",
	Run: func(cmd *cobra.Command, args []string) {
		quiltimport := exec.Command("git", "quiltimport")
		quiltimport.Stdout = os.Stdout
		quiltimport.Stderr = os.Stderr
		if err := quiltimport.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	gitCmd.AddCommand(
		// Missing porcelain
		gitWhatchangedCmd,
		gitAnnotateCmd,
		gitCherryCmd,
		gitStageCmd,
		gitMaintenanceCmd,
		gitCredentialCmd,
		gitHookCmd,

		// Missing plumbing
		gitForEachRefCmd,
		gitMergeBaseCmd,
		gitMergeFileCmd,
		gitMergeTreeCmd,
		gitMergeIndexCmd,
		gitPackObjectsCmd,
		gitIndexPackCmd,
		gitUnpackObjectsCmd,
		gitUpdateIndexCmd,
		gitCheckoutIndexCmd,
		gitCommitGraphCmd,
		gitMultiPackIndexCmd,
		gitInterpretTrailersCmd,
		gitPatchIdCmd,
		gitVerifyCommitCmd,
		gitVerifyTagCmd,
		gitMktagCmd,
		gitMktreeCmd,
		gitPackRefsCmd,
		gitPrunePackedCmd,
		gitStripspaceCmd,
		gitColumnCmd,
		gitCheckAttrCmd,
		gitCheckIgnoreCmd,
		gitCheckMailmapCmd,
		gitCheckRefFormatCmd,
		gitVarCmd,
		gitGetTarCommitIdCmd,
		gitShowIndexCmd,
		gitUnpackFileCmd,
		gitUpdateServerInfoCmd,
		gitFmtMergeMsgCmd,
		gitMailsplitCmd,
		gitMailinfoCmd,
		gitImapSendCmd,
		gitFastImportCmd,
		gitFastExportCmd,
		gitFilterBranchCmd,
		gitDaemonCmd,
		gitShellCmd,
		gitReceivePackCmd,
		gitSendPackCmd,
		gitUploadArchiveCmd,
		gitUploadPackCmd,
		gitCredentialCacheCmd,
		gitCredentialStoreCmd,
		gitSubtreeCmd,
		gitP4Cmd,
		gitSvnCmd,
		gitQuiltimportCmd,
	)
}
