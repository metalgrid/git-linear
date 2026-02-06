package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/metalgrid/git-linear/internal/auth"
	"github.com/metalgrid/git-linear/internal/git"
	"github.com/metalgrid/git-linear/internal/linear"
	"github.com/metalgrid/git-linear/internal/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "git-linear",
	Short: "Create git branches from Linear issues",
	Long:  `git-linear is a CLI tool that helps you create properly-named git branches from your assigned Linear issues.`,
	RunE:  runRoot,
}

func runRoot(cmd *cobra.Command, args []string) error {
	// Check if inside git repo
	if !git.IsInsideWorkTree() {
		return fmt.Errorf("not a git repository. Run this from inside a git project")
	}

	// Check for uncommitted changes
	if git.HasUncommittedChanges() {
		return fmt.Errorf("you have uncommitted changes. Please commit or stash them before creating a new branch")
	}

	// Get API key from keyring
	apiKey, err := auth.GetAPIKey()
	if err != nil {
		return fmt.Errorf("no API key found. Run 'git linear auth' to set up your Linear API key")
	}

	// Create Linear client
	client := linear.NewClient(apiKey)

	// Create and run TUI
	model := tui.NewModel(client)
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running TUI: %w", err)
	}

	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
