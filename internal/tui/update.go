package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/user/git-linear/internal/branch"
	"github.com/user/git-linear/internal/git"
)

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return m.loadIssuesCmd
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.state == StateIssueList || m.state == StateError || m.state == StateResult {
				return m, tea.Quit
			}
		case "esc":
			if m.state == StateBranchEdit {
				m.state = StateIssueList
				return m, nil
			}
			if m.state == StateConfirm {
				m.state = StateBranchEdit
				return m, m.branchEditor.Focus()
			}
			if m.state == StateExistingBranch {
				m.state = StateIssueList
				return m, nil
			}
		case "enter":
			return m.handleEnter()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.state == StateIssueList {
			m.issueList.SetSize(msg.Width, msg.Height-5)
		}

	case issuesLoadedMsg:
		if msg.err != nil {
			m.state = StateError
			m.errorMsg = fmt.Sprintf("Failed to load issues: %v", msg.err)
			return m, nil
		}
		if len(msg.issues) == 0 {
			m.state = StateError
			m.errorMsg = "No assigned issues found"
			return m, nil
		}

		// Convert issues to list items
		items := make([]list.Item, len(msg.issues))
		for i, issue := range msg.issues {
			// Check if branch exists for this issue
			branchName := branch.Sanitize(issue.Identifier, issue.Title)
			branchExists := git.BranchExists(branchName)
			items[i] = IssueItem{Issue: issue, BranchExists: branchExists}
		}

		m.issueList = list.New(items, IssueDelegate{}, m.width, m.height-5)
		m.issueList.Title = "Select an Issue"
		m.state = StateIssueList
		return m, nil

	case branchCreatedMsg:
		if msg.err != nil {
			m.state = StateError
			m.errorMsg = fmt.Sprintf("Failed to create/switch branch: %v", msg.err)
			return m, nil
		}
		m.state = StateResult
		m.resultMsg = fmt.Sprintf("✓ Switched to branch: %s", m.branchName)
		return m, tea.Quit
	}

	// Update active component based on state
	var cmd tea.Cmd
	switch m.state {
	case StateIssueList:
		m.issueList, cmd = m.issueList.Update(msg)
	case StateBranchEdit:
		m.branchEditor, cmd = m.branchEditor.Update(msg)
	}

	return m, cmd
}

func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.state {
	case StateIssueList:
		// Get selected issue
		item, ok := m.issueList.SelectedItem().(IssueItem)
		if !ok {
			return m, nil
		}
		m.selectedIssue = &item.Issue

		// Generate branch name
		branchName := branch.Sanitize(item.Issue.Identifier, item.Issue.Title)

		// Check if branch exists
		if git.BranchExists(branchName) {
			m.existingBranch = branchName
			m.state = StateExistingBranch
			return m, nil
		}

		// Move to branch edit
		m.branchEditor = NewBranchEditor(
			branch.Sanitize(item.Issue.Identifier, ""),
			item.Issue.Title,
		)
		m.state = StateBranchEdit
		return m, m.branchEditor.Focus()

	case StateBranchEdit:
		m.branchName = m.branchEditor.Value()
		m.state = StateConfirm
		m.branchEditor.Blur()
		return m, nil

	case StateConfirm:
		return m, m.createBranchCmd

	case StateExistingBranch:
		// Switch to existing branch
		m.branchName = m.existingBranch
		return m, m.switchBranchCmd

	case StateResult, StateError:
		return m, tea.Quit
	}

	return m, nil
}
