package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/user/git-linear/internal/branch"
)

var prefixStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

// BranchEditor wraps textinput for editing branch names with locked prefix
type BranchEditor struct {
	prefix    string
	textInput textinput.Model
}

// NewBranchEditor creates a new branch editor with locked prefix
func NewBranchEditor(prefix, defaultSuffix string) BranchEditor {
	ti := textinput.New()
	ti.Placeholder = "branch-description"
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 50
	ti.SetValue(defaultSuffix)

	return BranchEditor{
		prefix:    prefix,
		textInput: ti,
	}
}

// Init implements tea.Model
func (e BranchEditor) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model
func (e BranchEditor) Update(msg tea.Msg) (BranchEditor, tea.Cmd) {
	var cmd tea.Cmd
	e.textInput, cmd = e.textInput.Update(msg)
	return e, cmd
}

// View implements tea.Model
func (e BranchEditor) View() string {
	return prefixStyle.Render(e.prefix+"-") + e.textInput.View()
}

// Value returns the full sanitized branch name
func (e BranchEditor) Value() string {
	suffix := e.textInput.Value()
	// Use branch.Sanitize to get the full sanitized name
	// The prefix is already lowercase from Linear ID, suffix needs sanitization
	return branch.Sanitize(e.prefix, suffix)
}

// Focus sets focus on the text input
func (e *BranchEditor) Focus() tea.Cmd {
	return e.textInput.Focus()
}

// Blur removes focus from the text input
func (e *BranchEditor) Blur() {
	e.textInput.Blur()
}
