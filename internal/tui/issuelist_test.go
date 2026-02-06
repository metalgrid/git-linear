package tui_test

import (
	"github.com/metalgrid/git-linear/internal/linear"
	"github.com/metalgrid/git-linear/internal/tui"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IssueItem", func() {
	var issue linear.Issue

	BeforeEach(func() {
		issue = linear.Issue{
			ID:         "issue-1",
			Identifier: "DEV-123",
			Title:      "Fix login bug",
			State: linear.State{
				Name: "In Progress",
				Type: "started",
			},
		}
	})

	Describe("FilterValue", func() {
		It("returns identifier and title for fuzzy search", func() {
			item := tui.IssueItem{Issue: issue}
			Expect(item.FilterValue()).To(Equal("DEV-123 Fix login bug"))
		})
	})

	Describe("Title", func() {
		It("returns identifier with 2-space prefix when branch doesn't exist", func() {
			item := tui.IssueItem{Issue: issue, BranchExists: false}
			Expect(item.Title()).To(Equal("  DEV-123"))
		})

		It("returns identifier with * prefix when branch exists", func() {
			item := tui.IssueItem{Issue: issue, BranchExists: true}
			Expect(item.Title()).To(Equal("* DEV-123"))
		})
	})

	Describe("Description", func() {
		It("returns the issue title", func() {
			item := tui.IssueItem{Issue: issue}
			Expect(item.Description()).To(Equal("Fix login bug"))
		})
	})
})
