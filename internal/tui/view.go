package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("170"))
	errorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	resultStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
	helpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

// View implements tea.Model
func (m Model) View() string {
	switch m.state {
	case StateLoading:
		return "Loading issues...\n"

	case StateIssueList:
		help := helpStyle.Render("\nj/k or ↑/↓: navigate • enter: select • q: quit")
		return m.issueList.View() + help

	case StateBranchEdit:
		title := titleStyle.Render(fmt.Sprintf("Issue: %s - %s\n\n", m.selectedIssue.Identifier, m.selectedIssue.Title))
		editor := "Branch name:\n" + m.branchEditor.View() + "\n\n"
		help := helpStyle.Render("enter: confirm • esc: back")
		return title + editor + help

	case StateConfirm:
		title := titleStyle.Render(fmt.Sprintf("Issue: %s - %s\n\n", m.selectedIssue.Identifier, m.selectedIssue.Title))
		confirm := fmt.Sprintf("Create branch: %s\n\n", m.branchName)
		help := helpStyle.Render("enter: create • esc: back")
		return title + confirm + help

	case StateExistingBranch:
		title := titleStyle.Render(fmt.Sprintf("Issue: %s - %s\n\n", m.selectedIssue.Identifier, m.selectedIssue.Title))
		msg := fmt.Sprintf("Branch '%s' already exists.\n\n", m.existingBranch)
		help := helpStyle.Render("enter: switch to existing branch • esc: back")
		return title + msg + help

	case StateResult:
		return resultStyle.Render(m.resultMsg) + "\n"

	case StateError:
		return errorStyle.Render("Error: "+m.errorMsg) + "\n"
	}

	return ""
}
