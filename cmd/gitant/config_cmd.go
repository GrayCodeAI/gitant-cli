package main

import (
	"fmt"
	"os"

	"github.com/GrayCodeAI/gitant-cli/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:     "config",
	Short:   "Manage gitant configuration (like gh config)",
	Aliases: []string{"cfg"},
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a config value (daemon_url, web_url, ucan_token)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		val, err := config.Get(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: unknown key %q (use: daemon_url, web_url, ucan_token)\n", args[0])
			os.Exit(1)
		}
		fmt.Println(val)
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a config value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.Set(args[0], args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Set %s\n", args[0])
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all config values",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("daemon_url:\t%s\n", orDefault(s.DaemonURL, "(not set)"))
		fmt.Printf("web_url:\t%s\n", orDefault(s.WebURL, "(not set)"))
		if s.UCANToken != "" {
			fmt.Printf("ucan_token:\t(set, %d chars)\n", len(s.UCANToken))
		} else {
			fmt.Println("ucan_token:\t(not set)")
		}
	},
}

func orDefault(v, def string) string {
	if v == "" {
		return def
	}
	return v
}

func init() {
	configCmd.AddCommand(configGetCmd, configSetCmd, configListCmd)
	rootCmd.AddCommand(configCmd)
}
