package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api <path>",
	Short: "Make an authenticated Gitant API request (like gh api)",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		method, _ := cmd.Flags().GetString("method")
		fields, _ := cmd.Flags().GetStringArray("field")
		path := args[0]
		if !strings.HasPrefix(path, "/") {
			path = "/api/v1/" + strings.TrimPrefix(path, "/")
		}

		body := map[string]interface{}{}
		for _, f := range fields {
			parts := strings.SplitN(f, "=", 2)
			if len(parts) != 2 {
				fmt.Fprintf(os.Stderr, "Error: invalid field %q (use key=value)\n", f)
				os.Exit(1)
			}
			body[parts[0]] = parts[1]
		}

		client := newClient(cmd)
		var payload interface{}
		if len(body) > 0 {
			payload = body
		}

		switch method {
		case "GET":
			raw, err := client.GetRaw(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			prettyPrint(raw)
		case "DELETE":
			if err := client.Delete(path); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(`{"deleted":true}`)
		default:
			var result json.RawMessage
			if err := client.Request(method, path, payload, &result); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			prettyPrint(result)
		}
	},
}

func prettyPrint(raw []byte) {
	var v interface{}
	if err := json.Unmarshal(raw, &v); err != nil {
		fmt.Println(string(raw))
		return
	}
	printJSON(v)
}

func init() {
	apiCmd.Flags().StringP("method", "X", "GET", "HTTP method")
	apiCmd.Flags().StringArrayP("field", "f", nil, "JSON field for POST/PUT/PATCH (key=value)")
	addDaemonURLFlag(apiCmd)
	rootCmd.AddCommand(apiCmd)
}
