package main

import (
	"fmt"
	"os"

	"github.com/GrayCodeAI/gitant-cli/internal/cli"
	"github.com/spf13/cobra"
)

var forumCmd = &cobra.Command{
	Use:   "forum",
	Short: "Manage forum threads (from Fossil)",
}

var forumListCmd = &cobra.Command{
	Use:   "list",
	Short: "List forum threads",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		category, _ := cmd.Flags().GetString("category")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		path := fmt.Sprintf("/api/v1/repos/%s/forum", repo)
		if category != "" {
			path += "?category=" + category
		}

		var result struct {
			Threads []struct {
				ID       string `json:"id"`
				Title    string `json:"title"`
				Author   string `json:"author"`
				Category string `json:"category"`
				Replies  int    `json:"replies"`
				Views    int    `json:"views"`
				Pinned   bool   `json:"pinned"`
			} `json:"threads"`
		}
		if err := client.Get(path, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, t := range result.Threads {
			pinned := ""
			if t.Pinned {
				pinned = "📌"
			}
			fmt.Printf("%s%s\t%s\t[%s]\t%s\t%d replies\t%d views\n", pinned, t.ID, t.Title, t.Category, t.Author, t.Replies, t.Views)
		}
	},
}

var forumCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a forum thread",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		title, _ := cmd.Flags().GetString("title")
		body, _ := cmd.Flags().GetString("body")
		category, _ := cmd.Flags().GetString("category")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		if title == "" {
			title = PromptRequired("Title")
		}

		client := cli.NewClient(daemonURL)
		req := map[string]string{
			"title":    title,
			"body":     body,
			"category": category,
		}

		var result map[string]interface{}
		if err := client.Post(fmt.Sprintf("/api/v1/repos/%s/forum", repo), req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created thread: %s\n", result["id"])
	},
}

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Manage chat messages (from Fossil)",
}

var chatSendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a chat message",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		channel, _ := cmd.Flags().GetString("channel")
		message, _ := cmd.Flags().GetString("message")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		if message == "" {
			message = PromptRequired("Message")
		}

		client := cli.NewClient(daemonURL)
		req := map[string]string{
			"channel": channel,
			"body":    message,
		}

		var result map[string]interface{}
		if err := client.Post(fmt.Sprintf("/api/v1/repos/%s/chat", repo), req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Message sent to #%s\n", channel)
	},
}

var chatReadCmd = &cobra.Command{
	Use:   "read",
	Short: "Read chat messages",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		channel, _ := cmd.Flags().GetString("channel")
		limit, _ := cmd.Flags().GetInt("limit")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		path := fmt.Sprintf("/api/v1/repos/%s/chat/%s?limit=%d", repo, channel, limit)

		var result struct {
			Messages []struct {
				Author    string `json:"author"`
				Body      string `json:"body"`
				CreatedAt string `json:"created_at"`
			} `json:"messages"`
		}
		if err := client.Get(path, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, m := range result.Messages {
			fmt.Printf("[%s] %s: %s\n", m.CreatedAt, m.Author, m.Body)
		}
	},
}

var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Manage workspaces (from OneDev)",
	Aliases: []string{"ws"},
}

var workspaceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List workspaces",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result struct {
			Workspaces []struct {
				ID     string `json:"id"`
				Name   string `json:"name"`
				Status string `json:"status"`
				Branch string `json:"branch"`
			} `json:"workspaces"`
		}
		if err := client.Get(fmt.Sprintf("/api/v1/repos/%s/workspaces", repo), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, ws := range result.Workspaces {
			fmt.Printf("%s\t%s\t%s\t%s\n", ws.ID, ws.Name, ws.Status, ws.Branch)
		}
	},
}

var workspaceCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a workspace",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		name, _ := cmd.Flags().GetString("name")
		branch, _ := cmd.Flags().GetString("branch")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		if name == "" {
			name = PromptRequired("Name")
		}

		client := cli.NewClient(daemonURL)
		req := map[string]string{
			"name":   name,
			"branch": branch,
		}

		var result map[string]interface{}
		if err := client.Post(fmt.Sprintf("/api/v1/repos/%s/workspaces", repo), req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created workspace: %s\n", result["id"])
	},
}

var serviceDeskCmd = &cobra.Command{
	Use:   "service-desk",
	Short: "Manage service desk tickets (from OneDev)",
	Aliases: []string{"sd", "ticket"},
}

var serviceDeskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tickets",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		status, _ := cmd.Flags().GetString("status")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		path := fmt.Sprintf("/api/v1/repos/%s/tickets", repo)
		if status != "" {
			path += "?status=" + status
		}

		var result struct {
			Tickets []struct {
				ID       string `json:"id"`
				Title    string `json:"title"`
				Status   string `json:"status"`
				Priority string `json:"priority"`
				Assignee string `json:"assignee"`
			} `json:"tickets"`
		}
		if err := client.Get(path, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, t := range result.Tickets {
			fmt.Printf("%s\t%s\t%s\t%s\t%s\n", t.ID, t.Title, t.Status, t.Priority, t.Assignee)
		}
	},
}

var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "Time tracking (from OneDev)",
}

var timeStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start time tracking",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		description, _ := cmd.Flags().GetString("description")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		if description == "" {
			description = PromptRequired("Description")
		}

		client := cli.NewClient(daemonURL)
		req := map[string]string{
			"description": description,
		}

		var result map[string]interface{}
		if err := client.Post(fmt.Sprintf("/api/v1/repos/%s/time/start", repo), req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Started timer: %s\n", result["id"])
	},
}

var timeStopCmd = &cobra.Command{
	Use:   "stop [entry-id]",
	Short: "Stop time tracking",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Post(fmt.Sprintf("/api/v1/repos/%s/time/%s/stop", repo, args[0]), nil, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Stopped timer %s\n", args[0])
	},
}

var timeSummaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Show time tracking summary",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Get(fmt.Sprintf("/api/v1/repos/%s/time/summary", repo), &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Total hours: %v\n", result["total_hours"])
		fmt.Printf("Billable hours: %v\n", result["billable_hours"])
		fmt.Printf("Entries: %v\n", result["entries_count"])
	},
}

var governanceCmd = &cobra.Command{
	Use:   "governance",
	Short: "Manage governance proposals (from Gitopia)",
	Aliases: []string{"gov", "proposal"},
}

var governanceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List proposals",
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		status, _ := cmd.Flags().GetString("status")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		path := fmt.Sprintf("/api/v1/repos/%s/proposals", repo)
		if status != "" {
			path += "?status=" + status
		}

		var result struct {
			Proposals []struct {
				ID       string `json:"id"`
				Title    string `json:"title"`
				Status   string `json:"status"`
				Type     string `json:"type"`
				Proposer string `json:"proposer"`
			} `json:"proposals"`
		}
		if err := client.Get(path, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, p := range result.Proposals {
			fmt.Printf("%s\t%s\t%s\t%s\t%s\n", p.ID, p.Title, p.Status, p.Type, p.Proposer)
		}
	},
}

var governanceVoteCmd = &cobra.Command{
	Use:   "vote [proposal-id]",
	Short: "Vote on a proposal",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo, _ := cmd.Flags().GetString("repo")
		vote, _ := cmd.Flags().GetString("vote")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		if vote == "" {
			vote = PromptSelect("Vote", []string{"for", "against", "abstain"})
		}

		client := cli.NewClient(daemonURL)
		req := map[string]string{
			"vote": vote,
		}

		var result map[string]interface{}
		if err := client.Post(fmt.Sprintf("/api/v1/repos/%s/proposals/%s/vote", repo, args[0]), req, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Voted %s on proposal %s\n", vote, args[0])
	},
}

func init() {
	// Forum
	for _, c := range []*cobra.Command{forumListCmd, forumCreateCmd} {
		c.Flags().StringP("repo", "r", "", "Repository name (required)")
		c.MarkFlagRequired("repo")
		c.Flags().String("daemon-url", "", "Daemon URL")
	}
	forumListCmd.Flags().String("category", "", "Filter by category")
	forumCreateCmd.Flags().StringP("title", "t", "", "Thread title")
	forumCreateCmd.Flags().StringP("body", "b", "", "Thread body")
	forumCreateCmd.Flags().String("category", "general", "Category")
	forumCmd.AddCommand(forumListCmd, forumCreateCmd)

	// Chat
	for _, c := range []*cobra.Command{chatSendCmd, chatReadCmd} {
		c.Flags().StringP("repo", "r", "", "Repository name (required)")
		c.MarkFlagRequired("repo")
		c.Flags().String("daemon-url", "", "Daemon URL")
	}
	chatSendCmd.Flags().String("channel", "general", "Channel")
	chatSendCmd.Flags().StringP("message", "m", "", "Message")
	chatReadCmd.Flags().String("channel", "general", "Channel")
	chatReadCmd.Flags().IntP("limit", "l", 50, "Limit")
	chatCmd.AddCommand(chatSendCmd, chatReadCmd)

	// Workspace
	for _, c := range []*cobra.Command{workspaceListCmd, workspaceCreateCmd} {
		c.Flags().StringP("repo", "r", "", "Repository name (required)")
		c.MarkFlagRequired("repo")
		c.Flags().String("daemon-url", "", "Daemon URL")
	}
	workspaceCreateCmd.Flags().StringP("name", "n", "", "Workspace name")
	workspaceCreateCmd.Flags().StringP("branch", "b", "main", "Branch")
	workspaceCmd.AddCommand(workspaceListCmd, workspaceCreateCmd)

	// Service Desk
	for _, c := range []*cobra.Command{serviceDeskListCmd} {
		c.Flags().StringP("repo", "r", "", "Repository name (required)")
		c.MarkFlagRequired("repo")
		c.Flags().String("daemon-url", "", "Daemon URL")
	}
	serviceDeskListCmd.Flags().String("status", "", "Filter by status")
	serviceDeskCmd.AddCommand(serviceDeskListCmd)

	// Time
	for _, c := range []*cobra.Command{timeStartCmd, timeStopCmd, timeSummaryCmd} {
		c.Flags().StringP("repo", "r", "", "Repository name (required)")
		c.MarkFlagRequired("repo")
		c.Flags().String("daemon-url", "", "Daemon URL")
	}
	timeStartCmd.Flags().StringP("description", "d", "", "Description")
	timeCmd.AddCommand(timeStartCmd, timeStopCmd, timeSummaryCmd)

	// Governance
	for _, c := range []*cobra.Command{governanceListCmd, governanceVoteCmd} {
		c.Flags().StringP("repo", "r", "", "Repository name (required)")
		c.MarkFlagRequired("repo")
		c.Flags().String("daemon-url", "", "Daemon URL")
	}
	governanceListCmd.Flags().String("status", "", "Filter by status")
	governanceVoteCmd.Flags().String("vote", "", "Vote (for|against|abstain)")
	governanceCmd.AddCommand(governanceListCmd, governanceVoteCmd)

	rootCmd.AddCommand(forumCmd, chatCmd, workspaceCmd, serviceDeskCmd, timeCmd, governanceCmd)
}
