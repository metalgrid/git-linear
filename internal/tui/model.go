package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/metalgrid/git-linear/internal/linear"
)

// Model is the main TUI model
type Model struct {
	state          State
	issueList      list.Model
	branchEditor   BranchEditor
	selectedIssue  *linear.Issue
	branchName     string
	errorMsg       string
	resultMsg      string
	width          int
	height         int
	linearClient   *linear.Client
	existingBranch string
}

// NewModel creates a new TUI model
func NewModel(client *linear.Client) Model {
	return Model{
		state:        StateLoading,
		linearClient: client,
	}
}

// issuesLoadedMsg is sent when issues are loaded
type issuesLoadedMsg struct {
	issues []linear.Issue
	err    error
}

// branchCreatedMsg is sent when a branch is created
type branchCreatedMsg struct {
	err error
}
