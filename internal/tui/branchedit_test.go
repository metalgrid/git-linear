package tui_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/user/git-linear/internal/tui"
)

func TestBranchEdit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BranchEdit Suite")
}

var _ = Describe("BranchEditor", func() {
	Describe("NewBranchEditor", func() {
		It("creates editor with prefix and default suffix", func() {
			editor := tui.NewBranchEditor("dev-123", "fix-login")
			Expect(editor.Value()).To(Equal("dev-123-fix-login"))
		})
	})

	Describe("Value", func() {
		It("returns sanitized full branch name", func() {
			editor := tui.NewBranchEditor("dev-123", "Fix Login Bug")
			Expect(editor.Value()).To(Equal("dev-123-fix-login-bug"))
		})

		It("handles special characters in suffix", func() {
			editor := tui.NewBranchEditor("dev-456", "Hello! @World#")
			Expect(editor.Value()).To(Equal("dev-456-hello-world"))
		})

		It("handles empty suffix", func() {
			editor := tui.NewBranchEditor("dev-789", "")
			Expect(editor.Value()).To(Equal("dev-789"))
		})
	})

	Describe("View", func() {
		It("renders prefix and editable suffix", func() {
			editor := tui.NewBranchEditor("dev-123", "test")
			view := editor.View()
			// View should contain the prefix (styled) and the textinput
			Expect(view).To(ContainSubstring("dev-123-"))
		})
	})
})
