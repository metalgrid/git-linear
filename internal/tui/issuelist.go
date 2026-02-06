package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/metalgrid/git-linear/internal/linear"
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

// IssueItem wraps a Linear issue for use in bubbles/list
type IssueItem struct {
	Issue        linear.Issue
	BranchExists bool
}

// FilterValue implements list.Item interface for fuzzy search
func (i IssueItem) FilterValue() string {
	return i.Issue.Identifier + " " + i.Issue.Title
}

// Title implements list.DefaultItem interface
func (i IssueItem) Title() string {
	prefix := "  "
	if i.BranchExists {
		prefix = "* "
	}
	return prefix + i.Issue.Identifier
}

// Description implements list.DefaultItem interface
func (i IssueItem) Description() string {
	return i.Issue.Title
}

// IssueDelegate is a custom delegate for rendering issue items
type IssueDelegate struct{}

// Height implements list.ItemDelegate
func (d IssueDelegate) Height() int { return 1 }

// Spacing implements list.ItemDelegate
func (d IssueDelegate) Spacing() int { return 0 }

// Update implements list.ItemDelegate
func (d IssueDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

// Render implements list.ItemDelegate
func (d IssueDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(IssueItem)
	if !ok {
		return
	}

	prefix := "  "
	if i.BranchExists {
		prefix = "* "
	}

	str := fmt.Sprintf("%s%s  %s", prefix, i.Issue.Identifier, i.Issue.Title)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = selectedItemStyle.Render
	}

	fmt.Fprint(w, fn(str))
}
