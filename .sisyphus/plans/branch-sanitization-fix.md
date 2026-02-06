# Branch Name Sanitization Fix

## TL;DR

> **Quick Summary**: Fix branch name sanitization to work in real-time during editing, pre-sanitize the default suggestion, and reduce max length from 60 to 32 characters.
> 
> **Deliverables**:
> - Updated `internal/branch/sanitize.go` with 32 char limit
> - Updated `internal/tui/branchedit.go` with real-time sanitization
> - Updated `internal/tui/update.go` to pass pre-sanitized default
> - Updated tests for new behavior
> 
> **Estimated Effort**: Quick (2-3 tasks, ~30 minutes)
> **Parallel Execution**: NO - sequential (tests depend on implementation)

---

## Context

### Original Request
User reported that proposed branch names contain spaces, quotes, and other invalid symbols. They also want branch names to be shorter (max 32 characters instead of 60).

### Root Cause Analysis

**Issue 1: Raw title shown in editor**
```go
// update.go:116-119
m.branchEditor = NewBranchEditor(
    branch.Sanitize(item.Issue.Identifier, ""),  // ✓ prefix sanitized
    item.Issue.Title,                             // ✗ RAW title passed!
)
```
The default suffix is the raw issue title ("Fix Login Bug") instead of sanitized ("fix-login-bug").

**Issue 2: No real-time sanitization**
The `BranchEditor` only sanitizes when `Value()` is called (on submit), not while the user types. Users see raw text in the editor.

**Issue 3: Length too long**
Current max is 60 chars. User wants 32 chars for shorter, cleaner branch names.

### Interview Summary
- **Real-time sanitization**: YES - sanitize as user types
- **Length limit**: 32 chars TOTAL (prefix + suffix combined)
- **Sanitization rules**: Keep existing rules (lowercase, [a-z0-9-] only, collapse hyphens)

---

## Work Objectives

### Core Objective
Make branch name editing show sanitized text in real-time with a 32-character maximum.

### Concrete Deliverables
- `internal/branch/sanitize.go` - Change max from 60 to 32
- `internal/tui/branchedit.go` - Real-time sanitization display
- `internal/tui/update.go` - Pre-sanitize default suffix
- `internal/branch/sanitize_test.go` - Update test for 32 char limit
- `internal/tui/branchedit_test.go` - Update tests if needed

### Must Have
- Branch names never exceed 32 characters
- Editor shows sanitized text as user types
- Default suggestion is pre-sanitized
- All existing tests updated and passing

### Must NOT Have (Guardrails)
- No changes to sanitization rules (keep [a-z0-9-] pattern)
- No changes to the TUI state machine flow
- No new dependencies

---

## TODOs

### Task 1: Update Sanitize Function (32 char limit)

**What to do**:
1. Change max length constant from 60 to 32 in `sanitize.go`
2. Update test in `sanitize_test.go` to expect 32 chars
3. Verify all tests pass

**Must NOT do**:
- Change sanitization rules
- Change function signature

**Recommended Agent Profile**:
- **Category**: `quick`
- **Skills**: []

**References**:
- `internal/branch/sanitize.go:46-58` - Current truncation logic (change 60 → 32)
- `internal/branch/sanitize_test.go:37-42` - Test for max length

**Acceptance Criteria**:
- [x] `Sanitize("DEV-123", longtitle)` returns ≤32 chars
- [x] `go test ./internal/branch/... -v` passes
- [x] Truncation preserves valid branch name (no trailing hyphen)

**Commit**: YES
- Message: `fix(branch): reduce max branch name length to 32 chars`
- Files: `internal/branch/sanitize.go`, `internal/branch/sanitize_test.go`

---

### Task 2: Add Real-Time Sanitization to BranchEditor

**What to do**:
1. Add a `SanitizeSuffix(s string) string` helper function that sanitizes just the suffix part (without combining with prefix)
2. Modify `NewBranchEditor` to pre-sanitize the `defaultSuffix` before setting it
3. Modify `Update()` to sanitize the input value after each keystroke
4. Modify `View()` to show sanitized preview of full branch name
5. Update tests in `branchedit_test.go`

**Implementation approach**:
```go
// In branchedit.go

// SanitizeSuffix sanitizes just the suffix portion (no prefix combination)
func SanitizeSuffix(s string) string {
    s = strings.ToLower(s)
    s = strings.ReplaceAll(s, " ", "-")
    s = removeNonASCII(s)
    s = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(s, "")
    s = regexp.MustCompile(`-+`).ReplaceAllString(s, "-")
    s = strings.Trim(s, "-")
    return s
}

// NewBranchEditor - pre-sanitize the default suffix
func NewBranchEditor(prefix, defaultSuffix string) BranchEditor {
    sanitizedSuffix := sanitizeSuffix(defaultSuffix)
    // ... rest of setup with sanitizedSuffix
}

// Update - sanitize after each keystroke
func (e BranchEditor) Update(msg tea.Msg) (BranchEditor, tea.Cmd) {
    e.textInput, cmd = e.textInput.Update(msg)
    // Sanitize the current value
    currentVal := e.textInput.Value()
    sanitized := sanitizeSuffix(currentVal)
    if sanitized != currentVal {
        e.textInput.SetValue(sanitized)
    }
    return e, cmd
}
```

**Must NOT do**:
- Change the state machine in update.go
- Add new package imports beyond what's needed

**Recommended Agent Profile**:
- **Category**: `quick`
- **Skills**: []

**References**:
- `internal/tui/branchedit.go` - Current implementation
- `internal/branch/sanitize.go:22-36` - Sanitization logic to extract
- `internal/tui/branchedit_test.go` - Current tests

**Acceptance Criteria**:
- [ ] `NewBranchEditor("dev-123", "Fix Login Bug")` shows "fix-login-bug" initially
- [ ] Typing "Hello World!" in editor shows "hello-world" in real-time
- [ ] `go test ./internal/tui/... -v` passes
- [ ] No regression in other TUI tests

**Commit**: YES
- Message: `fix(tui): add real-time sanitization to branch editor`
- Files: `internal/tui/branchedit.go`, `internal/tui/branchedit_test.go`

---

### Task 3: Update Call Site and Integration Test

**What to do**:
1. Verify `update.go` call site works with new behavior (may not need changes if Task 2 handles pre-sanitization)
2. Run full test suite
3. Do hands-on QA to verify the flow works end-to-end

**Recommended Agent Profile**:
- **Category**: `quick`
- **Skills**: []

**References**:
- `internal/tui/update.go:116-119` - Call site for NewBranchEditor

**Acceptance Criteria**:
- [ ] `go test ./... -v` - all 44+ tests pass
- [ ] `go build ./cmd/git-linear` succeeds
- [ ] Hands-on QA: branch editor shows sanitized name

**Commit**: NO (verification only, or small fix if needed)

---

## Success Criteria

### Verification Commands
```bash
# All tests pass
go test ./... -v

# Build succeeds
go build -o git-linear ./cmd/git-linear

# Manual verification
./git-linear  # Select an issue, verify branch editor shows sanitized name
```

### Final Checklist
- [ ] Branch names never exceed 32 characters
- [ ] Editor shows sanitized text (lowercase, hyphens, no special chars)
- [ ] Default suggestion is pre-sanitized
- [ ] Real-time sanitization as user types
- [ ] All tests pass (44+)
- [ ] No breaking changes to TUI flow
