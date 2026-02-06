package main

import (
	"fmt"
	"os"

	"github.com/metalgrid/git-linear/internal/auth"
	"github.com/metalgrid/git-linear/internal/linear"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with Linear",
	Long:  `Store your Linear Personal API Key securely in the system keyring.`,
	RunE:  runAuth,
}

func init() {
	rootCmd.AddCommand(authCmd)
}

func runAuth(cmd *cobra.Command, args []string) error {
	fmt.Println("To create a Linear API key:")
	fmt.Println("  1. Go to Linear Settings → Account → Security")
	fmt.Println("  2. Under 'Personal API keys', click 'Create key'")
	fmt.Println("  3. Copy the generated key")
	fmt.Println()
	fmt.Println("Opening Linear settings in your browser...")
	fmt.Println()

	// Open browser (fire-and-forget, don't block on errors)
	if err := browser.OpenURL("https://linear.app/settings/account/security"); err != nil {
		fmt.Println("Could not open browser automatically.")
		fmt.Println("Please visit: https://linear.app/settings/account/security")
		fmt.Println()
	}

	fmt.Print("Enter your Linear API key: ")

	// Read password (hidden input)
	apiKeyBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println() // New line after hidden input

	if err != nil {
		return fmt.Errorf("failed to read API key: %w", err)
	}

	apiKey := string(apiKeyBytes)
	if apiKey == "" {
		return fmt.Errorf("API key cannot be empty")
	}

	// Validate API key by making a test request
	client := linear.NewClient(apiKey)
	if err := client.ValidateAPIKey(); err != nil {
		return fmt.Errorf("invalid API key: %w", err)
	}

	// Store in keyring
	if err := auth.StoreAPIKey(apiKey); err != nil {
		return fmt.Errorf("failed to store API key: %w", err)
	}

	fmt.Println("✓ API key validated and stored successfully")
	return nil
}
