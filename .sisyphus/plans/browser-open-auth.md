# Browser Open for Auth Flow

## TL;DR

> **Quick Summary**: Enhance `git linear auth` to automatically open the user's browser to Linear's API key settings page before prompting for the key.
> 
> **Deliverables**:
> - Updated `cmd/git-linear/auth.go` with browser open functionality
> - Cross-platform browser opening using `github.com/pkg/browser`
> 
> **Estimated Effort**: Quick (1 task, ~15 minutes)
> **Parallel Execution**: NO - single task

---

## Context

### Original Request
User wants the `git linear auth` command to open the browser to `https://linear.app/x3me/settings/account/security` to make it easier to create an API key.

### Interview Summary
- **Trigger**: Browser opens during `git linear auth` command
- **Flow**: Show instructions → Open browser → Prompt for key
- **URL**: `https://linear.app/settings/account/security` (generic, not workspace-specific)

### Technical Approach
Use `github.com/pkg/browser` for cross-platform browser opening:
- Linux: xdg-open
- macOS: open
- Windows: start

---

## Work Objectives

### Core Objective
Improve auth UX by automatically opening the Linear API key settings page.

### Concrete Deliverables
- Updated `cmd/git-linear/auth.go` with new flow
- New dependency: `github.com/pkg/browser`

### Must Have
- Cross-platform browser opening (Linux, macOS, Windows)
- Clear instructions printed before browser opens
- Graceful handling if browser fails to open (continue with prompt)

### Must NOT Have (Guardrails)
- No workspace-specific URLs (use generic Linear settings URL)
- No blocking on browser open (should be fire-and-forget)
- No new commands or flags

---

## TODOs

### Task 1: Add Browser Open to Auth Flow

**What to do**:
1. Add dependency: `go get github.com/pkg/browser`
2. Update `cmd/git-linear/auth.go`:
   - Print instructions about creating API key
   - Open browser to `https://linear.app/settings/account/security`
   - Handle browser open errors gracefully (warn, don't fail)
   - Continue to prompt for API key

**Must NOT do**:
- Block execution if browser fails to open
- Add any new CLI flags
- Change the API key validation logic

**Recommended Agent Profile**:
- **Category**: `quick`
  - Reason: Single file change, well-defined behavior, minimal code
- **Skills**: []
  - No special skills needed

**References**:

**Current Implementation**:
- `cmd/git-linear/auth.go:24-52` - Current `runAuth` function to modify

**Library Reference**:
- `github.com/pkg/browser` - Cross-platform browser opener
- Usage: `browser.OpenURL("https://...")`

**Acceptance Criteria**:

**Functional Requirements**:
- [x] `go get github.com/pkg/browser` succeeds
- [x] `go build ./cmd/git-linear` succeeds
- [x] Running `git linear auth` prints instructions
- [x] Running `git linear auth` opens browser to Linear settings
- [x] If browser fails to open, command continues (warns, doesn't error)
- [x] API key prompt appears after browser open attempt

**Agent-Executed QA Scenarios**:

```
Scenario: Auth command shows instructions and opens browser
  Tool: interactive_bash (tmux)
  Preconditions: Binary built, display available (or headless graceful fail)
  Steps:
    1. tmux new-session: ./git-linear auth
    2. Wait for: output (timeout: 3s)
    3. Assert: output contains "Linear API key"
    4. Assert: output contains instruction text about creating key
    5. Assert: output contains "Enter your" prompt
    6. Send: Ctrl+C to exit
  Expected Result: Instructions shown, prompt appears
  Evidence: Terminal output captured

Scenario: Auth continues if browser fails to open
  Tool: Bash
  Preconditions: DISPLAY unset (headless), binary built
  Steps:
    1. unset DISPLAY
    2. Run: echo "" | timeout 5 ./git-linear auth 2>&1
    3. Assert: output contains "Enter your Linear API key" (prompt reached)
    4. Assert: exit is NOT a panic (graceful handling)
  Expected Result: Command continues despite browser failure
  Evidence: Output captured
```

**Commit**: YES
- Message: `feat(auth): open browser to Linear settings during auth flow`
- Files: `cmd/git-linear/auth.go`, `go.mod`, `go.sum`
- Pre-commit: `go build ./cmd/git-linear && go test ./...`

---

## New Auth Flow

```go
func runAuth(cmd *cobra.Command, args []string) error {
    // 1. Print instructions
    fmt.Println("To create a Linear API key:")
    fmt.Println("  1. Go to Linear Settings → Account → Security")
    fmt.Println("  2. Under 'Personal API keys', click 'Create key'")
    fmt.Println("  3. Copy the generated key")
    fmt.Println()
    fmt.Println("Opening Linear settings in your browser...")
    fmt.Println()

    // 2. Open browser (fire-and-forget, don't block on errors)
    if err := browser.OpenURL("https://linear.app/settings/account/security"); err != nil {
        fmt.Println("Could not open browser automatically.")
        fmt.Println("Please visit: https://linear.app/settings/account/security")
        fmt.Println()
    }

    // 3. Prompt for API key (existing logic)
    fmt.Print("Enter your Linear API key: ")
    // ... rest of existing code
}
```

---

## Success Criteria

### Verification Commands
```bash
# Build succeeds
go build -o git-linear ./cmd/git-linear

# Tests still pass
go test ./...

# Auth shows instructions (will prompt for key, Ctrl+C to exit)
./git-linear auth
```

### Final Checklist
- [x] Browser opens to Linear settings (when display available)
- [x] Instructions printed clearly
- [x] Graceful fallback when browser can't open
- [x] All existing tests still pass
- [x] No breaking changes to auth flow
