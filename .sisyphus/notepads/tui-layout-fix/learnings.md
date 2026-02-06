# Learnings - TUI Layout Fix

## Session: ses_3ce001381ffe7ZJkkSWkk2J6uk
Started: 2026-02-06T09:41:28.217Z

---

## [2026-02-06 11:42] Task: Fix newline positioning in view.go

### Root Cause
- Newlines inside `lipgloss.Render()` cause cursor positioning issues
- ANSI escape codes from styling affect newline behavior
- Text continues horizontally instead of breaking to new line
- The styled text's logical position doesn't advance properly when newlines are inside the render call

### Fix Applied
- Moved `\n\n` from inside `Render()` to outside for 3 locations
- Lines 27, 33, 39 in `internal/tui/view.go`
- Pattern: `Render(fmt.Sprintf("...\n\n", ...))` → `Render(fmt.Sprintf("...", ...)) + "\n\n"`

**Locations fixed:**
1. Line 27 (StateBranchEdit): Issue title rendering
2. Line 33 (StateConfirm): Issue title rendering
3. Line 39 (StateExistingBranch): Issue title rendering

### Verification Results
- ✅ Build: `go build ./cmd/git-linear` - Success
- ✅ Tests: `go test ./internal/tui/... -v` - All 9 tests passed
- ✅ Layout: Labels now appear at start of lines (not far right)

### Key Insight
When using lipgloss styling with newlines, the newlines must be applied AFTER the render call, not inside the format string. This ensures the ANSI escape codes don't interfere with cursor positioning logic.
