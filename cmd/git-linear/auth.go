package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/git-linear/internal/auth"
	"github.com/user/git-linear/internal/linear"
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
