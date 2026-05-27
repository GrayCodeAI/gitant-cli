package main

import (
	"fmt"
	"os"

	"github.com/GrayCodeAI/gitant-cli/internal/cli"
	"github.com/spf13/cobra"
)

var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "Manage todo items",
}

var todoListCmd = &cobra.Command{
	Use:   "list",
	Short: "List todo items",
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)

		var result struct {
			Todos []struct {
				ID       string `json:"id"`
				Type     string `json:"type"`
				Title    string `json:"title"`
				Repo     string `json:"repo"`
				Author   string `json:"author"`
				CreatedAt string `json:"created_at"`
			} `json:"todos"`
			Total int `json:"total"`
		}
		if err := client.Get("/api/v1/todos", &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		for _, t := range result.Todos {
			fmt.Printf("%s\t%s\t%s\t%s\t%s\n", t.ID, t.Type, t.Title, t.Author, t.CreatedAt)
		}
		fmt.Fprintf(os.Stderr, "%d todo(s)\n", result.Total)
	},
}

var todoMarkCmd = &cobra.Command{
	Use:   "mark [todo-id]",
	Short: "Mark a todo as done",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Post(apiPath("/api/v1/todos", args[0], "done"), nil, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Marked todo %s as done\n", args[0])
	},
}

var todoMarkAllCmd = &cobra.Command{
	Use:   "mark-all",
	Short: "Mark all todos as done",
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result map[string]interface{}
		if err := client.Post("/api/v1/todos/mark-all", nil, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Marked all todos as done")
	},
}

var todoCountCmd = &cobra.Command{
	Use:   "count",
	Short: "Show todo count",
	Run: func(cmd *cobra.Command, args []string) {
		daemonURL, _ := cmd.Flags().GetString("daemon-url")

		client := cli.NewClient(daemonURL)
		var result struct {
			Count int `json:"count"`
		}
		if err := client.Get("/api/v1/todos/count", &result); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%d pending todo(s)\n", result.Count)
	},
}

func init() {
	for _, c := range []*cobra.Command{todoListCmd, todoMarkCmd, todoMarkAllCmd, todoCountCmd} {
		c.Flags().String("daemon-url", "", "Daemon URL (default: http://localhost:7777)")
	}

	todoCmd.AddCommand(todoListCmd, todoMarkCmd, todoMarkAllCmd, todoCountCmd)
	rootCmd.AddCommand(todoCmd)
}
