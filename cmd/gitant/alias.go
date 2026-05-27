package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// AliasStore manages CLI aliases
type AliasStore struct {
	path    string
	aliases map[string]string
}

// NewAliasStore creates a new alias store
func NewAliasStore() *AliasStore {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	path := filepath.Join(home, ".gitant", "aliases.json")

	store := &AliasStore{
		path:    path,
		aliases: make(map[string]string),
	}
	store.load()
	return store
}

func (s *AliasStore) load() {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return
	}
	if err := json.Unmarshal(data, &s.aliases); err != nil {
		s.aliases = make(map[string]string)
	}
}

func (s *AliasStore) save() error {
	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating alias directory: %w", err)
	}

	data, err := json.MarshalIndent(s.aliases, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

// Set sets an alias
func (s *AliasStore) Set(name, command string) error {
	s.aliases[name] = command
	return s.save()
}

// Get gets an alias
func (s *AliasStore) Get(name string) (string, bool) {
	cmd, ok := s.aliases[name]
	return cmd, ok
}

// Delete deletes an alias
func (s *AliasStore) Delete(name string) error {
	delete(s.aliases, name)
	return s.save()
}

// List lists all aliases
func (s *AliasStore) List() map[string]string {
	return s.aliases
}

var aliasStore = NewAliasStore()

var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Manage command aliases",
	Long:  "Create shortcuts for gitant commands.",
}

var aliasSetCmd = &cobra.Command{
	Use:   "set <name> <command>",
	Short: "Set a command alias",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		command := args[1]

		if err := aliasStore.Set(name, command); err != nil {
			fmt.Fprintf(os.Stderr, "Error setting alias: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Added alias: %s → %s\n", name, command)
	},
}

var aliasDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a command alias",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		if _, ok := aliasStore.Get(name); !ok {
			fmt.Fprintf(os.Stderr, "Alias not found: %s\n", name)
			os.Exit(1)
		}

		if err := aliasStore.Delete(name); err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting alias: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Deleted alias: %s\n", name)
	},
}

var aliasListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all aliases",
	Run: func(cmd *cobra.Command, args []string) {
		aliases := aliasStore.List()
		if len(aliases) == 0 {
			fmt.Println("No aliases configured.")
			return
		}

		fmt.Println("Aliases:")
		for name, command := range aliases {
			fmt.Printf("  %s → %s\n", name, command)
		}
	},
}

// ExpandAlias expands an alias if it exists
func ExpandAlias(args []string) []string {
	if len(args) == 0 {
		return args
	}

	name := args[0]
	if cmd, ok := aliasStore.Get(name); ok {
		parts := strings.Fields(cmd)
		return append(parts, args[1:]...)
	}

	return args
}

func init() {
	aliasCmd.AddCommand(aliasSetCmd)
	aliasCmd.AddCommand(aliasDeleteCmd)
	aliasCmd.AddCommand(aliasListCmd)
	rootCmd.AddCommand(aliasCmd)
}
