package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/user/git-linear/internal/git"
)

// loadIssuesCmd fetches issues from Linear
func (m Model) loadIssuesCmd() tea.Msg {
	issues, err := m.linearClient.GetAssignedIssues()
	return issuesLoadedMsg{issues: issues, err: err}
}

// createBranchCmd creates a new git branch
func (m Model) createBranchCmd() tea.Msg {
	// Get default branch
	defaultBranch, err := git.GetDefaultBranch()
	if err != nil {
		return branchCreatedMsg{err: err}
	}

	// Create branch from default
	err = git.CreateBranch(m.branchName, defaultBranch)
	if err != nil {
		return branchCreatedMsg{err: err}
	}

	// Switch to new branch
	err = git.SwitchBranch(m.branchName)
	return branchCreatedMsg{err: err}
}

// switchBranchCmd switches to an existing branch
func (m Model) switchBranchCmd() tea.Msg {
	err := git.SwitchBranch(m.existingBranch)
	return branchCreatedMsg{err: err}
}
