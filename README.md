# git-linear

Create git branches from Linear issues.

## Installation

```bash
go install github.com/metalgrid/git-linear/cmd/git-linear@latest
```

## Usage

### Authentication

```bash
git-linear auth
```

This opens your browser to Linear's API settings. Create a personal API key and paste it when prompted.

### Create a branch

```bash
git-linear
```

Select an issue from your assigned Linear issues, edit the branch name if needed, and confirm to create/switch to the branch.

## License

MIT
