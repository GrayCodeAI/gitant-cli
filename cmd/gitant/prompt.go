package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Prompt prompts the user for input
func Prompt(label string, defaultValue string) string {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", label, defaultValue)
	} else {
		fmt.Printf("%s: ", label)
	}

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" && defaultValue != "" {
		return defaultValue
	}
	return input
}

// PromptRequired prompts for required input
func PromptRequired(label string) string {
	for {
		value := Prompt(label, "")
		if value != "" {
			return value
		}
		fmt.Println("This field is required.")
	}
}

// PromptConfirm prompts for yes/no confirmation
func PromptConfirm(label string, defaultYes bool) bool {
	suffix := " [y/N]: "
	if defaultYes {
		suffix = " [Y/n]: "
	}

	fmt.Printf("%s%s", label, suffix)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "" {
		return defaultYes
	}

	return input == "y" || input == "yes"
}

// PromptSelect prompts the user to select from a list
func PromptSelect(label string, options []string) string {
	fmt.Printf("%s:\n", label)
	for i, option := range options {
		fmt.Printf("  %d) %s\n", i+1, option)
	}

	for {
		fmt.Print("Select: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			return options[0]
		}

		for i, option := range options {
			if fmt.Sprintf("%d", i+1) == input || strings.ToLower(option) == strings.ToLower(input) {
				return option
			}
		}

		fmt.Println("Invalid selection. Try again.")
	}
}

// PromptMultiSelect prompts the user to select multiple options
func PromptMultiSelect(label string, options []string) []string {
	fmt.Printf("%s (comma-separated numbers):\n", label)
	for i, option := range options {
		fmt.Printf("  %d) %s\n", i+1, option)
	}

	fmt.Print("Select: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return []string{options[0]}
	}

	var selected []string
	parts := strings.Split(input, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		for i, option := range options {
			if fmt.Sprintf("%d", i+1) == part || strings.ToLower(option) == strings.ToLower(part) {
				selected = append(selected, option)
				break
			}
		}
	}

	return selected
}

// PromptMultiline prompts for multiline input
func PromptMultiline(label string) string {
	fmt.Printf("%s (empty line to finish):\n", label)

	var lines []string
	reader := bufio.NewReader(os.Stdin)
	for {
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}
