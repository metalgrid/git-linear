package tui

// State represents the current state of the TUI
type State int

const (
	StateLoading State = iota
	StateIssueList
	StateBranchEdit
	StateConfirm
	StateExistingBranch
	StateResult
	StateError
)
