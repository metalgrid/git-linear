# git-linear: Cross-Platform CLI for Linear Issue Branch Management

## TL;DR

> **Quick Summary**: Build a Go CLI tool (`git linear`) that fetches Linear issues assigned to the user, presents them in a TUI for selection, and creates properly-named git branches with Linear issue prefixes.
> 
> **Deliverables**:
> - `git linear` - Interactive TUI for issue selection and branch creation
> - `git linear auth` - Secure API key storage via system keyring
> - Cross-platform binaries (Linux, macOS, Windows)
> 
> **Estimated Effort**: Medium (8-12 tasks)
> **Parallel Execution**: YES - 3 waves
> **Critical Path**: Project Setup → Branch Sanitization → Auth → Linear Client → Git Ops → TUI → Integration

---

## Context

### Original Request
User needs a CLI tool invoked as `git linear` that:
- Connects to Linear API and fetches issues assigned to the user
- Provides a TUI to select an issue
- Suggests a branch name with Linear ID prefix (e.g., `dev-123-fix-login`)
- Allows editing the description part (prefix is locked)
- Creates and switches to the branch on confirmation
- Detects existing branches and prompts for action

### Interview Summary

**Key Discussions**:
- **Branch naming**: Lowercase prefix from Linear ID (DEV-123 → `dev-123-description`)
- **Issue filtering**: Active issues only (In Progress, Todo, Backlog) - excludes Done/Canceled
- **Auth flow**: Separate `git linear auth` subcommand, Personal API Key stored in keyring
- **Existing branches**: Show indicator, prompt to switch or create with suffix
- **Testing**: TDD with BDD using Ginkgo/v2 and Gomega
- **Scope v1**: Only `git linear` and `git linear auth` commands

**Research Findings**:
- **Linear API**: GraphQL at `https://api.linear.app/graphql`, Personal API Key in `Authorization` header
- **Bubbletea**: Model/Update/View pattern, bubbles/list for selection, bubbles/textinput for editing
- **go-keyring**: Cross-platform credential storage (Keychain/Credential Manager/Secret Service)

### Metis Review

**Identified Gaps** (addressed):
- **Branch name sanitization**: Max 60 chars, [a-z0-9-] only, handle unicode/emoji/special chars
- **Dirty working tree**: Refuse if uncommitted changes exist (user must commit/stash first)
- **Branch base**: Always create from default branch (main/master), not current HEAD
- **Empty states**: Handle no assigned issues, API errors, network timeouts gracefully
- **CLI structure**: Use Cobra for standard Go CLI patterns

**Guardrails Applied**:
- No Linear write operations in v1
- No configuration files beyond keyring
- No subcommands beyond `auth`
- No flags/options in v1
- No non-interactive mode

---

## Work Objectives

### Core Objective
Create a streamlined workflow for developers to create properly-named git branches linked to their Linear issues, reducing context-switching and ensuring consistent branch naming conventions.

### Concrete Deliverables
- `cmd/git-linear/main.go` - CLI entry point with Cobra
- `internal/branch/sanitize.go` - Branch name sanitization logic
- `internal/auth/keyring.go` - Credential storage abstraction
- `internal/linear/client.go` - GraphQL client for Linear API
- `internal/git/git.go` - Git operations (exec-based)
- `internal/tui/` - Bubbletea TUI components
- Test files with Ginkgo/Gomega for each package

### Definition of Done
- [x] `go build ./cmd/git-linear` produces working binary
- [x] `./git-linear auth` stores and validates API key
- [x] `./git-linear` shows TUI with assigned issues
- [x] Branch creation works with proper naming
- [x] All Ginkgo tests pass: `go test ./... -v`

### Must Have
- Cross-platform credential storage (keyring with fallback warning)
- Branch name sanitization (safe for git, max 60 chars)
- Active issue filtering (exclude Done/Canceled)
- Existing branch detection with user prompt
- Dirty working tree detection (refuse if uncommitted changes)
- User-friendly error messages (no stack traces)

### Must NOT Have (Guardrails)
- Linear write operations (no status updates, issue creation)
- Configuration files (keyring only for v1)
- Additional subcommands beyond `auth`
- CLI flags or options
- Workspace/team filtering
- Non-interactive mode
- PR creation or integration
- Sorting/filtering in TUI beyond API defaults

---

## Verification Strategy

> **UNIVERSAL RULE: ZERO HUMAN INTERVENTION**
>
> ALL verification is executed by the agent using tools. No manual testing steps.

### Test Decision
- **Infrastructure exists**: NO (greenfield project)
- **Automated tests**: YES - TDD with BDD using Ginkgo/v2 and Gomega
- **Framework**: Ginkgo/v2 + Gomega

### TDD Workflow (Each Task)

**Task Structure:**
1. **RED**: Write failing Ginkgo test first
   - Test file: `{package}/{file}_test.go`
   - Test command: `go test ./internal/{package}/... -v`
   - Expected: FAIL (test exists, implementation doesn't)
2. **GREEN**: Implement minimum code to pass
   - Command: `go test ./internal/{package}/... -v`
   - Expected: PASS
3. **REFACTOR**: Clean up while keeping green

**Test Setup Task (Task 1):**
- Install: `go get github.com/onsi/ginkgo/v2` + `go get github.com/onsi/gomega`
- Verify: `go test ./... -v` runs (even if no tests yet)

### Agent-Executed QA Scenarios (MANDATORY — ALL tasks)

> Every task includes QA scenarios using Playwright (for any web verification),
> interactive_bash/tmux (for TUI/CLI verification), or curl (for API verification).

---

## Execution Strategy

### Parallel Execution Waves

```
Wave 1 (Start Immediately):
├── Task 1: Project Setup (go mod, structure, test framework)
└── Task 2: Branch Name Sanitization (pure logic, no deps)

Wave 2 (After Wave 1):
├── Task 3: Keyring Auth Storage
├── Task 4: Linear GraphQL Client
└── Task 5: Git Operations Module

Wave 3 (After Wave 2):
├── Task 6: TUI - Issue List Component
├── Task 7: TUI - Branch Name Editor
└── Task 8: TUI - State Machine & Flow

Wave 4 (After Wave 3):
├── Task 9: CLI Commands (Cobra integration)
└── Task 10: Integration & Error Handling

Critical Path: 1 → 3 → 4 → 6 → 8 → 9 → 10
Parallel Speedup: ~40% faster than sequential
```

### Dependency Matrix

| Task | Depends On | Blocks | Can Parallelize With |
|------|------------|--------|---------------------|
| 1 | None | 2,3,4,5 | None (must be first) |
| 2 | 1 | 7,8 | 3,4,5 |
| 3 | 1 | 9 | 2,4,5 |
| 4 | 1 | 6,8 | 2,3,5 |
| 5 | 1 | 8,9 | 2,3,4 |
| 6 | 4 | 8 | 7 |
| 7 | 2 | 8 | 6 |
| 8 | 2,4,5,6,7 | 9 | None |
| 9 | 3,5,8 | 10 | None |
| 10 | 9 | None | None (final) |

---

## TODOs

### Task 1: Project Setup & Test Framework

**What to do**:
- Initialize Go module: `go mod init github.com/user/git-linear`
- Create directory structure:
  ```
  git-linear/
  ├── cmd/git-linear/main.go
  ├── internal/
  │   ├── auth/
  │   ├── branch/
  │   ├── git/
  │   ├── linear/
  │   └── tui/
  ├── go.mod
  └── go.sum
  ```
- Install dependencies:
  - `go get github.com/onsi/ginkgo/v2`
  - `go get github.com/onsi/gomega`
  - `go get github.com/spf13/cobra`
  - `go get github.com/charmbracelet/bubbletea`
  - `go get github.com/charmbracelet/bubbles`
  - `go get github.com/charmbracelet/lipgloss`
  - `go get github.com/zalando/go-keyring`
- Create minimal `main.go` with "Hello from git-linear" output
- Verify build and test commands work

**Must NOT do**:
- Add any CLI logic beyond hello world
- Create configuration files
- Add any business logic

**Recommended Agent Profile**:
- **Category**: `quick`
  - Reason: Standard project scaffolding, straightforward setup
- **Skills**: []
  - No special skills needed for Go project setup

**Parallelization**:
- **Can Run In Parallel**: NO (must be first)
- **Parallel Group**: Wave 1 (alone initially)
- **Blocks**: Tasks 2, 3, 4, 5
- **Blocked By**: None

**References**:
- Go modules: https://go.dev/doc/modules/gomod-ref
- Ginkgo setup: https://onsi.github.io/ginkgo/#getting-started
- Cobra quick start: https://github.com/spf13/cobra#usage

**Acceptance Criteria**:

**TDD (tests enabled):**
- [x] Test file created: `internal/branch/sanitize_test.go` (empty suite)
- [x] `go test ./... -v` runs without error (0 tests OK)

**Agent-Executed QA Scenarios:**

```
Scenario: Project builds successfully
  Tool: Bash
  Preconditions: go.mod exists with dependencies
  Steps:
    1. cd /home/iso/Projects/git-linear
    2. go build -o git-linear ./cmd/git-linear
    3. Assert: exit code 0
    4. Assert: file git-linear exists and is executable
    5. ./git-linear
    6. Assert: output contains "git-linear" or similar
  Expected Result: Binary builds and runs
  Evidence: Build output captured

Scenario: Test framework configured
  Tool: Bash
  Preconditions: Ginkgo installed
  Steps:
    1. go test ./... -v
    2. Assert: exit code 0 (no test failures, may be 0 tests)
  Expected Result: Test command works
  Evidence: Test output captured
```

**Commit**: YES
- Message: `feat(init): scaffold go project with test framework`
- Files: `go.mod`, `go.sum`, `cmd/git-linear/main.go`, directory structure
- Pre-commit: `go build ./cmd/git-linear && go test ./...`

---

### Task 2: Branch Name Sanitization

**What to do**:
- Create `internal/branch/sanitize.go` with `Sanitize(identifier, title string) string`
- Rules:
  - Lowercase the Linear identifier (DEV-123 → dev-123)
  - Slugify title: lowercase, replace spaces with hyphens
  - Remove all chars except [a-z0-9-]
  - Collapse multiple hyphens to single
  - Strip leading/trailing hyphens
  - Max total length: 60 chars (truncate title part if needed)
  - If title becomes empty after sanitization, return just identifier
- Create `internal/branch/sanitize_test.go` with Ginkgo tests

**Must NOT do**:
- Add any git operations
- Add any network calls
- Add any TUI code

**Recommended Agent Profile**:
- **Category**: `quick`
  - Reason: Pure function, well-defined input/output, easy to TDD
- **Skills**: []
  - No special skills needed

**Parallelization**:
- **Can Run In Parallel**: YES
- **Parallel Group**: Wave 1 (after Task 1 completes)
- **Blocks**: Tasks 7, 8
- **Blocked By**: Task 1

**References**:
- Git branch naming: https://git-scm.com/docs/git-check-ref-format
- Go strings package: https://pkg.go.dev/strings
- Go regexp package: https://pkg.go.dev/regexp

**Acceptance Criteria**:

**TDD (tests enabled):**
- [x] Test file: `internal/branch/sanitize_test.go`
- [x] Tests cover: basic slugify, unicode removal, max length, empty title, special chars
- [x] `go test ./internal/branch/... -v` → PASS

**Test Cases to Implement:**
```go
Describe("Sanitize", func() {
    It("converts identifier to lowercase", func() {
        Expect(Sanitize("DEV-123", "foo")).To(Equal("dev-123-foo"))
    })
    It("slugifies title with spaces", func() {
        Expect(Sanitize("DEV-1", "Fix Login Bug")).To(Equal("dev-1-fix-login-bug"))
    })
    It("removes special characters", func() {
        Expect(Sanitize("DEV-1", "Hello! @World#")).To(Equal("dev-1-hello-world"))
    })
    It("removes emoji and unicode", func() {
        Expect(Sanitize("DEV-1", "Fix 🔐 Auth")).To(Equal("dev-1-fix-auth"))
    })
    It("collapses multiple hyphens", func() {
        Expect(Sanitize("DEV-1", "a - - b")).To(Equal("dev-1-a-b"))
    })
    It("truncates to max 60 chars", func() {
        long := strings.Repeat("a", 100)
        result := Sanitize("DEV-123", long)
        Expect(len(result)).To(BeNumerically("<=", 60))
        Expect(result).NotTo(HaveSuffix("-"))
    })
    It("returns just identifier if title sanitizes to empty", func() {
        Expect(Sanitize("DEV-1", "!@#$%")).To(Equal("dev-1"))
    })
})
```

**Agent-Executed QA Scenarios:**

```
Scenario: Sanitization function handles all edge cases
  Tool: Bash
  Preconditions: Task 1 complete, sanitize.go exists
  Steps:
    1. go test ./internal/branch/... -v
    2. Assert: all tests pass
    3. Assert: output shows test cases for unicode, length, special chars
  Expected Result: All sanitization tests green
  Evidence: Test output captured
```

**Commit**: YES
- Message: `feat(branch): add branch name sanitization with TDD`
- Files: `internal/branch/sanitize.go`, `internal/branch/sanitize_test.go`
- Pre-commit: `go test ./internal/branch/... -v`

---

### Task 3: Keyring Authentication Storage

**What to do**:
- Create `internal/auth/keyring.go` with:
  - `StoreAPIKey(key string) error` - validates and stores key
  - `GetAPIKey() (string, error)` - retrieves stored key
  - `DeleteAPIKey() error` - removes stored key
  - `HasAPIKey() bool` - checks if key exists
- Use `github.com/zalando/go-keyring`
- Service name: `git-linear`
- Username: `api-key`
- Handle `keyring.ErrNotFound` gracefully
- Create `internal/auth/keyring_test.go` with Ginkgo tests using `keyring.MockInit()`

**Must NOT do**:
- Validate key against Linear API (that's Task 4)
- Add any fallback to config file (warning is enough)
- Store anything besides the API key

**Recommended Agent Profile**:
- **Category**: `quick`
  - Reason: Thin wrapper around go-keyring, well-defined API
- **Skills**: []
  - No special skills needed

**Parallelization**:
- **Can Run In Parallel**: YES
- **Parallel Group**: Wave 2 (with Tasks 4, 5)
- **Blocks**: Task 9
- **Blocked By**: Task 1

**References**:
- go-keyring: https://github.com/zalando/go-keyring
- go-keyring mock: https://github.com/zalando/go-keyring#mocking

**Acceptance Criteria**:

**TDD (tests enabled):**
- [x] Test file: `internal/auth/keyring_test.go`
- [x] Tests use `keyring.MockInit()` for isolation
- [x] Tests cover: store, retrieve, delete, not found, overwrite
- [x] `go test ./internal/auth/... -v` → PASS

**Agent-Executed QA Scenarios:**

```
Scenario: Keyring operations work with mock
  Tool: Bash
  Preconditions: Task 1 complete
  Steps:
    1. go test ./internal/auth/... -v
    2. Assert: exit code 0
    3. Assert: tests for Store, Get, Delete pass
  Expected Result: All auth tests green
  Evidence: Test output captured
```

**Commit**: YES
- Message: `feat(auth): add keyring-based credential storage`
- Files: `internal/auth/keyring.go`, `internal/auth/keyring_test.go`
- Pre-commit: `go test ./internal/auth/... -v`

---

### Task 4: Linear GraphQL Client

**What to do**:
- Create `internal/linear/types.go` with:
  - `Issue` struct: `ID`, `Identifier`, `Title`, `State` (name, type)
  - `State` struct: `Name`, `Type`
- Create `internal/linear/client.go` with:
  - `Client` struct with API key
  - `NewClient(apiKey string) *Client`
  - `GetAssignedIssues() ([]Issue, error)` - fetches active issues
  - `ValidateAPIKey() error` - tests API key validity
- GraphQL query filtering: `state.type NOT IN ["completed", "canceled"]`
- Handle errors: network, auth (401), rate limit
- Create `internal/linear/client_test.go` with mocked HTTP responses

**Must NOT do**:
- Make any write/mutation calls to Linear
- Cache responses
- Implement pagination (fetch first 50 is enough for v1)

**Recommended Agent Profile**:
- **Category**: `unspecified-low`
  - Reason: HTTP client with GraphQL, standard patterns, moderate complexity
- **Skills**: []
  - No special skills needed

**Parallelization**:
- **Can Run In Parallel**: YES
- **Parallel Group**: Wave 2 (with Tasks 3, 5)
- **Blocks**: Tasks 6, 8
- **Blocked By**: Task 1

**References**:
- Linear API docs: https://developers.linear.app/docs/graphql/working-with-the-graphql-api
- Linear GraphQL schema: https://studio.apollographql.com/public/Linear-API/variant/current/home
- Go net/http: https://pkg.go.dev/net/http
- encoding/json: https://pkg.go.dev/encoding/json

**GraphQL Query to Implement:**
```graphql
query AssignedIssues {
  viewer {
    assignedIssues(
      first: 50
      filter: { state: { type: { nin: ["completed", "canceled"] } } }
    ) {
      nodes {
        id
        identifier
        title
        state {
          name
          type
        }
      }
    }
  }
}
```

**Acceptance Criteria**:

**TDD (tests enabled):**
- [x] Test file: `internal/linear/client_test.go`
- [x] Tests mock HTTP responses (httptest.Server)
- [x] Tests cover: successful fetch, 401 error, network error, empty response
- [x] `go test ./internal/linear/... -v` → PASS

**Agent-Executed QA Scenarios:**

```
Scenario: Linear client parses mock response correctly
  Tool: Bash
  Preconditions: Task 1 complete
  Steps:
    1. go test ./internal/linear/... -v
    2. Assert: exit code 0
    3. Assert: tests for Issue parsing, error handling pass
  Expected Result: All linear client tests green
  Evidence: Test output captured
```

**Commit**: YES
- Message: `feat(linear): add GraphQL client for fetching assigned issues`
- Files: `internal/linear/types.go`, `internal/linear/client.go`, `internal/linear/client_test.go`
- Pre-commit: `go test ./internal/linear/... -v`

---

### Task 5: Git Operations Module

**What to do**:
- Create `internal/git/git.go` with:
  - `IsInsideWorkTree() bool` - check if in git repo
  - `HasUncommittedChanges() bool` - check for dirty working tree
  - `GetDefaultBranch() (string, error)` - detect main/master
  - `BranchExists(name string) bool` - case-insensitive check
  - `CreateBranch(name, base string) error` - create from base
  - `SwitchBranch(name string) error` - checkout existing branch
  - `GetCurrentBranch() (string, error)` - current branch name
- All operations via `exec.Command("git", ...)`
- Create `internal/git/git_test.go` - tests with temp git repos

**Must NOT do**:
- Use go-git library (exec is simpler for our needs)
- Modify any files in the repo
- Push to remote

**Recommended Agent Profile**:
- **Category**: `quick`
  - Reason: Thin wrappers around git commands, straightforward
- **Skills**: [`git-master`]
  - `git-master`: Understanding git internals for proper command construction

**Parallelization**:
- **Can Run In Parallel**: YES
- **Parallel Group**: Wave 2 (with Tasks 3, 4)
- **Blocks**: Tasks 8, 9
- **Blocked By**: Task 1

**References**:
- Go os/exec: https://pkg.go.dev/os/exec
- Git commands: https://git-scm.com/docs

**Acceptance Criteria**:

**TDD (tests enabled):**
- [x] Test file: `internal/git/git_test.go`
- [x] Tests create temp git repo for isolation
- [x] Tests cover: inside/outside repo, clean/dirty, branch exists/not, create/switch
- [x] `go test ./internal/git/... -v` → PASS

**Agent-Executed QA Scenarios:**

```
Scenario: Git module detects repo state correctly
  Tool: Bash
  Preconditions: Task 1 complete
  Steps:
    1. go test ./internal/git/... -v
    2. Assert: exit code 0
    3. Assert: tests for IsInsideWorkTree, HasUncommittedChanges pass
  Expected Result: All git tests green
  Evidence: Test output captured
```

**Commit**: YES
- Message: `feat(git): add git operations module with exec commands`
- Files: `internal/git/git.go`, `internal/git/git_test.go`
- Pre-commit: `go test ./internal/git/... -v`

---

### Task 6: TUI - Issue List Component

**What to do**:
- Create `internal/tui/issuelist.go` with:
  - `IssueItem` implementing `list.Item` interface
  - `FilterValue()` returns identifier + title for fuzzy search
  - Custom delegate showing: `DEV-123  Title here` (with * for existing branches)
  - Styles using lipgloss
- Create `internal/tui/issuelist_test.go`

**Must NOT do**:
- Implement the full TUI state machine (that's Task 8)
- Add sorting or filtering beyond bubbles/list default
- Show issue descriptions or details

**Recommended Agent Profile**:
- **Category**: `visual-engineering`
  - Reason: TUI component with styling, visual presentation matters
- **Skills**: []
  - No special skills needed

**Parallelization**:
- **Can Run In Parallel**: YES
- **Parallel Group**: Wave 3 (with Task 7)
- **Blocks**: Task 8
- **Blocked By**: Task 4

**References**:
- bubbles/list: https://pkg.go.dev/github.com/charmbracelet/bubbles/list
- lipgloss: https://github.com/charmbracelet/lipgloss
- Bubbletea list example: https://github.com/charmbracelet/bubbletea/tree/master/examples/list-simple

**Acceptance Criteria**:

**TDD (tests enabled):**
- [x] Test file: `internal/tui/issuelist_test.go`
- [x] Tests cover: IssueItem creation, FilterValue correctness
- [x] `go test ./internal/tui/... -v` → PASS

**Agent-Executed QA Scenarios:**

```
Scenario: Issue list item renders correctly
  Tool: Bash
  Preconditions: Tasks 1, 4 complete
  Steps:
    1. go test ./internal/tui/... -v -run IssueList
    2. Assert: exit code 0
    3. Assert: IssueItem implements list.Item interface
  Expected Result: Issue list component tests pass
  Evidence: Test output captured
```

**Commit**: YES (groups with Task 7)
- Message: `feat(tui): add issue list component`
- Files: `internal/tui/issuelist.go`, `internal/tui/issuelist_test.go`
- Pre-commit: `go test ./internal/tui/... -v`

---

### Task 7: TUI - Branch Name Editor

**What to do**:
- Create `internal/tui/branchedit.go` with:
  - `BranchEditor` struct wrapping textinput.Model
  - `NewBranchEditor(prefix, defaultSuffix string)` - prefix is locked
  - Custom View showing: `dev-123-[editable part here]`
  - Prefix displayed but not editable
  - Uses sanitization from Task 2 on blur/submit
- Create `internal/tui/branchedit_test.go`

**Must NOT do**:
- Allow editing the prefix
- Allow invalid characters (sanitize on the fly or on submit)
- Implement confirmation (that's Task 8)

**Recommended Agent Profile**:
- **Category**: `visual-engineering`
  - Reason: TUI component with custom editing behavior
- **Skills**: []
  - No special skills needed

**Parallelization**:
- **Can Run In Parallel**: YES
- **Parallel Group**: Wave 3 (with Task 6)
- **Blocks**: Task 8
- **Blocked By**: Task 2

**References**:
- bubbles/textinput: https://pkg.go.dev/github.com/charmbracelet/bubbles/textinput
- Bubbletea textinput example: https://github.com/charmbracelet/bubbletea/tree/master/examples/textinputs

**Acceptance Criteria**:

**TDD (tests enabled):**
- [x] Test file: `internal/tui/branchedit_test.go`
- [x] Tests cover: prefix lock, suffix editing, sanitization on value
- [x] `go test ./internal/tui/... -v` → PASS

**Agent-Executed QA Scenarios:**

```
Scenario: Branch editor locks prefix
  Tool: Bash
  Preconditions: Tasks 1, 2 complete
  Steps:
    1. go test ./internal/tui/... -v -run BranchEdit
    2. Assert: exit code 0
    3. Assert: tests verify prefix is immutable
  Expected Result: Branch editor tests pass
  Evidence: Test output captured
```

**Commit**: YES (groups with Task 6)
- Message: `feat(tui): add branch name editor component`
- Files: `internal/tui/branchedit.go`, `internal/tui/branchedit_test.go`
- Pre-commit: `go test ./internal/tui/... -v`

---

### Task 8: TUI - State Machine & Flow

**What to do**:
- Create `internal/tui/model.go` with main Model struct
- Create `internal/tui/states.go` defining states:
  ```go
  type State int
  const (
      StateLoading State = iota
      StateIssueList
      StateBranchEdit
      StateConfirm
      StateExistingBranch  // prompt: switch or create new
      StateResult
      StateError
  )
  ```
- Create `internal/tui/update.go` with Update function handling:
  - State transitions
  - Key bindings (j/k, arrows, enter, esc, q)
  - Error handling
- Create `internal/tui/view.go` with View function rendering each state
- Create `internal/tui/commands.go` with tea.Cmd for:
  - Fetching issues (async)
  - Creating branch (async)
  - Checking existing branch

**Must NOT do**:
- Add sorting or filtering UI
- Add issue detail view
- Add any configuration options

**Recommended Agent Profile**:
- **Category**: `unspecified-high`
  - Reason: Complex state machine, multiple components, integration point
- **Skills**: []
  - No special skills needed

**Parallelization**:
- **Can Run In Parallel**: NO
- **Parallel Group**: Sequential (after Wave 3)
- **Blocks**: Task 9
- **Blocked By**: Tasks 2, 4, 5, 6, 7

**References**:
- Bubbletea tutorial: https://github.com/charmbracelet/bubbletea/tree/master/tutorials
- Composable views example: https://github.com/charmbracelet/bubbletea/tree/master/examples/composable-views
- State machine patterns in TUI: Previous librarian research

**State Machine Diagram:**
```
INIT → check auth & git repo
  ├── no auth → ERROR "Run git linear auth"
  ├── not git repo → ERROR "Not a git repository"
  └── ok → LOADING

LOADING → fetch issues
  ├── error → ERROR
  ├── empty → ERROR "No assigned issues"
  └── ok → ISSUE_LIST

ISSUE_LIST → user selects
  ├── esc/q → quit
  └── enter → check existing branch
      ├── exists → EXISTING_BRANCH
      └── not exists → BRANCH_EDIT

EXISTING_BRANCH → user chooses
  ├── switch → RESULT (switch)
  └── create new → BRANCH_EDIT (with suffix)

BRANCH_EDIT → user edits name
  ├── esc → ISSUE_LIST
  └── enter → CONFIRM

CONFIRM → user confirms
  ├── esc → BRANCH_EDIT
  └── enter → create branch → RESULT

RESULT → show success/failure → quit
```

**Acceptance Criteria**:

**TDD (tests enabled):**
- [x] Test files: `internal/tui/model_test.go`, `internal/tui/update_test.go`
- [x] Tests cover: state transitions, key handling, error states
- [x] `go test ./internal/tui/... -v` → PASS

**Agent-Executed QA Scenarios:**

```
Scenario: TUI state machine handles happy path
  Tool: Bash
  Preconditions: Tasks 1-7 complete
  Steps:
    1. go test ./internal/tui/... -v
    2. Assert: exit code 0
    3. Assert: tests cover state transitions
  Expected Result: TUI state machine tests pass
  Evidence: Test output captured
```

**Commit**: YES
- Message: `feat(tui): implement state machine and full TUI flow`
- Files: `internal/tui/model.go`, `internal/tui/states.go`, `internal/tui/update.go`, `internal/tui/view.go`, `internal/tui/commands.go`, test files
- Pre-commit: `go test ./internal/tui/... -v`

---

### Task 9: CLI Commands (Cobra Integration)

**What to do**:
- Create `cmd/git-linear/root.go` with root command
- Create `cmd/git-linear/auth.go` with `auth` subcommand:
  - Prompt for API key (hidden input)
  - Validate against Linear API
  - Store in keyring
  - Show success/error
- Update `cmd/git-linear/main.go`:
  - Root command launches TUI
  - Check preconditions (git repo, auth) before TUI
- Handle `--help` via Cobra defaults

**Must NOT do**:
- Add any flags beyond --help (auto from Cobra)
- Add any other subcommands
- Add version command (not in v1 scope)

**Recommended Agent Profile**:
- **Category**: `quick`
  - Reason: Standard Cobra CLI setup, well-documented patterns
- **Skills**: []
  - No special skills needed

**Parallelization**:
- **Can Run In Parallel**: NO
- **Parallel Group**: Sequential (after Task 8)
- **Blocks**: Task 10
- **Blocked By**: Tasks 3, 5, 8

**References**:
- Cobra user guide: https://github.com/spf13/cobra/blob/main/site/content/user_guide.md
- Cobra examples: https://github.com/spf13/cobra-cli
- Go term package for hidden input: https://pkg.go.dev/golang.org/x/term

**Acceptance Criteria**:

**TDD (tests enabled):**
- [x] `go build ./cmd/git-linear` → binary builds
- [x] `./git-linear --help` → shows usage
- [x] `./git-linear auth --help` → shows auth usage

**Agent-Executed QA Scenarios:**

```
Scenario: CLI commands wire up correctly
  Tool: Bash
  Preconditions: Tasks 1-8 complete
  Steps:
    1. go build -o git-linear ./cmd/git-linear
    2. Assert: exit code 0
    3. ./git-linear --help
    4. Assert: output contains "auth" subcommand
    5. ./git-linear auth --help
    6. Assert: output describes API key storage
  Expected Result: CLI commands work
  Evidence: Command outputs captured

Scenario: Auth command prompts for key
  Tool: interactive_bash (tmux)
  Preconditions: Binary built, keyring mock or real
  Steps:
    1. tmux new-session: ./git-linear auth
    2. Wait for: "API key" prompt (timeout: 5s)
    3. Send keys: "lin_api_test_invalid_key" Enter
    4. Assert: output contains error about invalid key OR validation message
  Expected Result: Auth flow initiates
  Evidence: Terminal output captured
```

**Commit**: YES
- Message: `feat(cli): add cobra commands for root and auth`
- Files: `cmd/git-linear/root.go`, `cmd/git-linear/auth.go`, `cmd/git-linear/main.go`
- Pre-commit: `go build ./cmd/git-linear && go test ./...`

---

### Task 10: Integration & Error Handling Polish

**What to do**:
- End-to-end testing of full flow
- Polish error messages (user-friendly, no stack traces)
- Add helpful hints in error states:
  - "Run `git linear auth` to set up your API key"
  - "Commit or stash your changes before switching branches"
  - "Check your internet connection"
- Ensure graceful exit on Ctrl+C
- Final build verification for all platforms

**Must NOT do**:
- Add new features
- Add logging or debug output
- Add telemetry

**Recommended Agent Profile**:
- **Category**: `unspecified-low`
  - Reason: Integration testing and polish, no new features
- **Skills**: []
  - No special skills needed

**Parallelization**:
- **Can Run In Parallel**: NO
- **Parallel Group**: Final (after Task 9)
- **Blocks**: None (final task)
- **Blocked By**: Task 9

**Acceptance Criteria**:

**Agent-Executed QA Scenarios:**

```
Scenario: Full happy path with mock data
  Tool: interactive_bash (tmux)
  Preconditions: Binary built, test git repo, mock/real Linear key
  Steps:
    1. Create temp directory and init git repo
    2. git init && git commit --allow-empty -m "init"
    3. Set up mock auth (or real key in keyring)
    4. tmux new-session: ./git-linear
    5. Wait for: issue list OR error about auth
    6. If issue list: use j/k to navigate, Enter to select
    7. Wait for: branch name editor
    8. Send Enter to confirm
    9. Assert: branch created message OR existing branch prompt
  Expected Result: Full flow works end-to-end
  Evidence: Terminal recording captured

Scenario: Error handling - not in git repo
  Tool: Bash
  Preconditions: Binary built
  Steps:
    1. cd /tmp && mkdir test-no-git && cd test-no-git
    2. /path/to/git-linear
    3. Assert: exit code != 0
    4. Assert: output contains "not a git repository" (friendly message)
    5. Assert: output does NOT contain stack trace
  Expected Result: Clean error message
  Evidence: Error output captured

Scenario: Error handling - no auth
  Tool: Bash
  Preconditions: Binary built, keyring cleared
  Steps:
    1. Create temp git repo
    2. Clear any stored key (keyring.Delete or mock)
    3. ./git-linear
    4. Assert: exit code != 0
    5. Assert: output contains "git linear auth" hint
  Expected Result: Helpful auth error
  Evidence: Error output captured

Scenario: Graceful Ctrl+C handling
  Tool: interactive_bash (tmux)
  Preconditions: Binary built, test setup
  Steps:
    1. tmux new-session: ./git-linear
    2. Wait for: any output (2s)
    3. Send Ctrl+C
    4. Assert: process exits cleanly (no panic)
    5. Assert: terminal state restored
  Expected Result: Clean exit on interrupt
  Evidence: Terminal output captured
```

**Commit**: YES
- Message: `feat(polish): integration testing and error handling`
- Files: Any touched files for polish
- Pre-commit: `go build ./cmd/git-linear && go test ./... -v`

---

## Commit Strategy

| After Task | Message | Files | Verification |
|------------|---------|-------|--------------|
| 1 | `feat(init): scaffold go project with test framework` | go.mod, structure | `go build && go test` |
| 2 | `feat(branch): add branch name sanitization with TDD` | internal/branch/* | `go test ./internal/branch/...` |
| 3 | `feat(auth): add keyring-based credential storage` | internal/auth/* | `go test ./internal/auth/...` |
| 4 | `feat(linear): add GraphQL client for fetching assigned issues` | internal/linear/* | `go test ./internal/linear/...` |
| 5 | `feat(git): add git operations module with exec commands` | internal/git/* | `go test ./internal/git/...` |
| 6+7 | `feat(tui): add issue list and branch editor components` | internal/tui/* | `go test ./internal/tui/...` |
| 8 | `feat(tui): implement state machine and full TUI flow` | internal/tui/* | `go test ./internal/tui/...` |
| 9 | `feat(cli): add cobra commands for root and auth` | cmd/git-linear/* | `go build && ./git-linear --help` |
| 10 | `feat(polish): integration testing and error handling` | various | `go test ./... -v` |

---

## Success Criteria

### Verification Commands
```bash
# Build
go build -o git-linear ./cmd/git-linear  # Expected: exit 0, binary created

# All tests pass
go test ./... -v  # Expected: PASS for all packages

# Help works
./git-linear --help  # Expected: shows usage with 'auth' subcommand
./git-linear auth --help  # Expected: shows auth usage

# Binary runs (will fail gracefully if no auth/git)
./git-linear  # Expected: helpful error message, not stack trace
```

### Final Checklist
- [x] All "Must Have" features present and working
- [x] All "Must NOT Have" items absent (no scope creep)
- [x] All Ginkgo tests pass
- [x] Binary builds on Linux, macOS, Windows (cross-compile verified)
- [x] Error messages are user-friendly (no stack traces)
- [x] TUI renders correctly in standard terminal
- [x] Ctrl+C exits gracefully
- [x] Branch names are valid git refs
