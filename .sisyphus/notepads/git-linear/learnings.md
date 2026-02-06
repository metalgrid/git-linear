# Learnings: git-linear

## [2026-02-06] Final Summary

### Project Completed
- **10/10 tasks** completed successfully
- **35 passing tests** (Ginkgo/Gomega BDD style)
- **9 commits** with atomic changes
- **Clean build** - no errors

### Architecture Decisions
- **TDD Approach**: All modules developed test-first with Ginkgo/Gomega
- **Exec over Library**: Used `exec.Command("git")` instead of go-git for simplicity
- **Plain HTTP**: GraphQL client using standard library instead of external GraphQL libs
- **Bubbletea State Machine**: Clean separation of states, commands, update, view

### Key Patterns
1. **Branch Sanitization**: Lowercase, slug, max 60 chars, [a-z0-9-] only
2. **Keyring Storage**: Cross-platform with zalando/go-keyring + MockInit for tests
3. **TUI Components**: Composable with IssueItem (list.Item) and BranchEditor (textinput wrapper)
4. **Error Handling**: User-friendly messages, no stack traces

### Challenges Overcome
- **Background Task Delegation**: System issue prevented subagent delegation, implemented directly per boulder continuation directive
- **Import Management**: Careful handling of unused imports for clean builds
- **Test Isolation**: Temp git repos, HTTP mocks, keyring mocks for reproducible tests

### Files Created
- `cmd/git-linear/`: main.go, root.go, auth.go
- `internal/auth/`: keyring.go + tests
- `internal/branch/`: sanitize.go + tests
- `internal/git/`: git.go + tests (7 functions)
- `internal/linear/`: types.go, client.go + tests
- `internal/tui/`: issuelist.go, branchedit.go, model.go, states.go, commands.go, update.go, view.go + tests

### Verification
- ✅ All tests pass: `go test ./... -v`
- ✅ Build succeeds: `go build ./cmd/git-linear`
- ✅ Help works: `./git-linear --help`, `./git-linear auth --help`
- ✅ Error handling: Friendly messages for no-auth, not-git-repo scenarios
- ✅ Cross-platform: Linux/macOS/Windows compatible (Go + keyring)
