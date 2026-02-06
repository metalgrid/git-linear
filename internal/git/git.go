package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// IsInsideWorkTree checks if the current directory is inside a git repository.
func IsInsideWorkTree() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = nil // Suppress stderr
	err := cmd.Run()
	if err != nil {
		return false
	}
	return strings.TrimSpace(out.String()) == "true"
}

// HasUncommittedChanges checks if the working tree has uncommitted changes.
func HasUncommittedChanges() bool {
	cmd := exec.Command("git", "status", "--porcelain")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false
	}
	return strings.TrimSpace(out.String()) != ""
}

// GetDefaultBranch detects the default branch (main or master).
// It tries origin/HEAD first, then falls back to common names.
func GetDefaultBranch() (string, error) {
	// Try to get the default branch from origin/HEAD
	cmd := exec.Command("git", "symbolic-ref", "refs/remotes/origin/HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err == nil {
		// Output is like "refs/remotes/origin/main"
		ref := strings.TrimSpace(out.String())
		parts := strings.Split(ref, "/")
		if len(parts) > 0 {
			return parts[len(parts)-1], nil
		}
	}

	// Fallback: try common branch names
	for _, name := range []string{"main", "master"} {
		cmd := exec.Command("git", "rev-parse", "--verify", name)
		if err := cmd.Run(); err == nil {
			return name, nil
		}
	}

	return "", fmt.Errorf("could not determine default branch")
}

// BranchExists checks if a branch exists (case-insensitive).
func BranchExists(name string) bool {
	// Get all branches and compare case-insensitively
	cmd := exec.Command("git", "branch", "-a")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return false
	}

	branches := strings.Split(strings.TrimSpace(out.String()), "\n")
	lowerName := strings.ToLower(name)

	for _, branch := range branches {
		// Remove leading/trailing whitespace and * (current branch indicator)
		branch = strings.TrimSpace(branch)
		branch = strings.TrimPrefix(branch, "* ")
		// Remove remote prefix if present (e.g., "remotes/origin/")
		if strings.Contains(branch, "/") {
			parts := strings.Split(branch, "/")
			branch = parts[len(parts)-1]
		}
		if strings.ToLower(branch) == lowerName {
			return true
		}
	}

	return false
}

// CreateBranch creates a new branch from a base branch.
func CreateBranch(name, base string) error {
	cmd := exec.Command("git", "branch", name, base)
	return cmd.Run()
}

// SwitchBranch switches to an existing branch.
func SwitchBranch(name string) error {
	cmd := exec.Command("git", "checkout", name)
	return cmd.Run()
}

// GetCurrentBranch returns the name of the current branch.
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}
