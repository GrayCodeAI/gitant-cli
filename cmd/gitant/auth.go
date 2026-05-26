package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/GrayCodeAI/gitant-cli/internal/config"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with a Gitant node (like gh auth)",
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Save daemon URL and UCAN token to ~/.gitant/config.json",
	Run: func(cmd *cobra.Command, args []string) {
		withToken, _ := cmd.Flags().GetString("with-token")
		daemonURL, _ := cmd.Flags().GetString("daemon-url")
		webURL, _ := cmd.Flags().GetString("web-url")

		reader := bufio.NewReader(os.Stdin)
		settings, _ := config.Load()

		if daemonURL == "" {
			def := settings.DaemonURL
			if def == "" {
				def = "http://localhost:7777"
			}
			fmt.Printf("Daemon URL [%s]: ", def)
			line, _ := reader.ReadString('\n')
			line = strings.TrimSpace(line)
			if line != "" {
				daemonURL = line
			} else {
				daemonURL = def
			}
		}

		if webURL == "" {
			def := settings.WebURL
			if def == "" {
				def = "http://localhost:3303"
			}
			fmt.Printf("Web URL [%s]: ", def)
			line, _ := reader.ReadString('\n')
			line = strings.TrimSpace(line)
			if line != "" {
				webURL = line
			} else {
				webURL = def
			}
		}

		if withToken == "" {
			fmt.Print("UCAN token (paste from web /agents): ")
			line, _ := reader.ReadString('\n')
			withToken = strings.TrimSpace(line)
		}

		if withToken == "" {
			fmt.Fprintln(os.Stderr, "Error: UCAN token is required")
			os.Exit(1)
		}

		settings.DaemonURL = daemonURL
		settings.WebURL = webURL
		settings.UCANToken = withToken
		if err := config.Save(settings); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Logged in to %s\n", daemonURL)
		fmt.Println("Config saved to ~/.gitant/config.json")
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication and daemon connection status",
	Run: func(cmd *cobra.Command, args []string) {
		settings, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		client := newClient(cmd)
		var status map[string]interface{}
		apiErr := client.Get("/api/v1/status", &status)

		fmt.Println("gitant auth status")
		fmt.Println("====================")
		if settings.DaemonURL != "" {
			fmt.Printf("  Daemon URL:  %s\n", settings.DaemonURL)
		} else {
			fmt.Println("  Daemon URL:  (not set — using default http://localhost:7777)")
		}
		if settings.WebURL != "" {
			fmt.Printf("  Web URL:     %s\n", settings.WebURL)
		}
		if settings.UCANToken != "" {
			fmt.Printf("  UCAN token:  set (%d chars)\n", len(settings.UCANToken))
		} else {
			fmt.Println("  UCAN token:  not set (read-only API; writes will fail)")
		}
		if apiErr != nil {
			fmt.Printf("  API:         unreachable (%v)\n", apiErr)
			os.Exit(1)
		}
		fmt.Printf("  API:         connected (version %v)\n", status["version"])
	},
}

var authTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Print the configured UCAN token",
	Run: func(cmd *cobra.Command, args []string) {
		settings, err := config.Load()
		if err != nil || settings.UCANToken == "" {
			fmt.Fprintln(os.Stderr, "Error: not logged in — run gitant auth login")
			os.Exit(1)
		}
		fmt.Println(settings.UCANToken)
	},
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove saved credentials",
	Run: func(cmd *cobra.Command, args []string) {
		settings, _ := config.Load()
		settings.UCANToken = ""
		if err := config.Save(settings); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Logged out (UCAN token cleared)")
	},
}

func init() {
	authLoginCmd.Flags().String("with-token", "", "UCAN token (skip prompt)")
	authLoginCmd.Flags().String("web-url", "", "Web dashboard URL")
	addDaemonURLFlag(authLoginCmd, authStatusCmd)

	authCmd.AddCommand(authLoginCmd, authStatusCmd, authTokenCmd, authLogoutCmd)
	rootCmd.AddCommand(authCmd)
}
