# TUI Layout Fix - Newline Positioning

## TL;DR

> **Quick Summary**: Fix garbled TUI layout where "Branch name:" and "Create branch:" labels appear far to the right instead of on their own lines.
> 
> **Deliverables**:
> - Fixed `internal/tui/view.go` with newlines moved outside lipgloss styling
> 
> **Estimated Effort**: Quick (1 task, ~5 minutes)
> **Parallel Execution**: NO - single task

---

## Context

### Original Request
User reported UI alignment issue:
```
Issue: DEV-7626 - Vanguard's socket management accumulates sessions
                                                                   
                                                                   Branch name:
dev-7626-> vanguards-socket-management-accumulates-sessions   
```

The "Branch name:" label and "Create branch:" text appear far to the right, on the same "logical line" as the previous content.

### Root Cause
The `\n\n` newlines are **inside** the `lipgloss.Render()` call:
```go
title := titleStyle.Render(fmt.Sprintf("Issue: %s - %s\n\n", ...))
```

Lipgloss styling can affect how newlines are rendered, causing text to continue horizontally instead of properly breaking.

### Fix
Move newlines **outside** the styled text:
```go
title := titleStyle.Render(fmt.Sprintf("Issue: %s - %s", ...)) + "\n\n"
```

---

## TODOs

### Task 1: Fix Newline Positioning in view.go

**What to do**:
Move `\n\n` from inside `titleStyle.Render()` to outside, for 3 locations:

1. Line 27 (StateBranchEdit):
   - Before: `titleStyle.Render(fmt.Sprintf("Issue: %s - %s\n\n", ...))`
   - After: `titleStyle.Render(fmt.Sprintf("Issue: %s - %s", ...)) + "\n\n"`

2. Line 33 (StateConfirm):
   - Before: `titleStyle.Render(fmt.Sprintf("Issue: %s - %s\n\n", ...))`
   - After: `titleStyle.Render(fmt.Sprintf("Issue: %s - %s", ...)) + "\n\n"`

3. Line 39 (StateExistingBranch):
   - Before: `titleStyle.Render(fmt.Sprintf("Issue: %s - %s\n\n", ...))`
   - After: `titleStyle.Render(fmt.Sprintf("Issue: %s - %s", ...)) + "\n\n"`

**Must NOT do**:
- Change any styling
- Change any other logic
- Add new imports

**Recommended Agent Profile**:
- **Category**: `quick`
- **Skills**: []

**References**:
- `internal/tui/view.go:27,33,39` - Lines to modify

**Acceptance Criteria**:
- [x] `go build ./cmd/git-linear` succeeds
- [x] `go test ./internal/tui/... -v` passes
- [x] TUI displays correctly: labels on their own lines

**Agent-Executed QA Scenarios**:

```
Scenario: Branch editor displays labels correctly
  Tool: interactive_bash (tmux)
  Preconditions: Binary built, Linear API key configured
  Steps:
    1. Run ./git-linear
    2. Select an issue with Enter
    3. Observe: "Branch name:" should be on its own line, not far right
    4. Press Enter to confirm
    5. Observe: "Create branch:" should be on its own line, not far right
    6. Press Ctrl+C to exit
  Expected Result: Labels appear at start of their own lines
  Evidence: Terminal output captured
```

**Commit**: YES
- Message: `fix(tui): correct newline positioning in view rendering`
- Files: `internal/tui/view.go`
- Pre-commit: `go build ./cmd/git-linear && go test ./internal/tui/... -v`

---

## Success Criteria

### Verification Commands
```bash
# Build succeeds
go build -o git-linear ./cmd/git-linear

# Tests pass
go test ./internal/tui/... -v

# Manual verification (requires API key)
./git-linear
```

### Expected Output After Fix
```
Issue: DEV-7626 - Vanguard's socket management accumulates sessions

Branch name:
dev-7626-vanguards-socket-manage

enter: confirm • esc: back
```

### Final Checklist
- [x] "Branch name:" label appears at start of line
- [x] "Create branch:" label appears at start of line
- [x] All tests pass
- [x] No visual regressions
