package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Git Suite")
}

var _ = Describe("Git Operations", func() {
	var tempDir string

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "git-test-*")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		if tempDir != "" {
			os.RemoveAll(tempDir)
		}
	})

	Describe("IsInsideWorkTree", func() {
		It("returns true when inside a git repository", func() {
			// Initialize a git repo
			cmd := exec.Command("git", "init")
			cmd.Dir = tempDir
			err := cmd.Run()
			Expect(err).NotTo(HaveOccurred())

			// Change to the temp directory
			oldCwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			defer os.Chdir(oldCwd)

			err = os.Chdir(tempDir)
			Expect(err).NotTo(HaveOccurred())

			result := IsInsideWorkTree()
			Expect(result).To(BeTrue())
		})

		It("returns false when outside a git repository", func() {
			oldCwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			defer os.Chdir(oldCwd)

			err = os.Chdir(tempDir)
			Expect(err).NotTo(HaveOccurred())

			result := IsInsideWorkTree()
			Expect(result).To(BeFalse())
		})
	})

	Describe("HasUncommittedChanges", func() {
		It("returns false in a clean repository", func() {
			// Initialize and commit
			cmd := exec.Command("git", "init")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "config", "user.email", "test@example.com")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "config", "user.name", "Test User")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "commit", "--allow-empty", "-m", "initial")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			oldCwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			defer os.Chdir(oldCwd)

			err = os.Chdir(tempDir)
			Expect(err).NotTo(HaveOccurred())

			result := HasUncommittedChanges()
			Expect(result).To(BeFalse())
		})

		It("returns true when there are uncommitted changes", func() {
			// Initialize and commit
			cmd := exec.Command("git", "init")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "config", "user.email", "test@example.com")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "config", "user.name", "Test User")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "commit", "--allow-empty", "-m", "initial")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			// Create a new file
			testFile := filepath.Join(tempDir, "test.txt")
			err := os.WriteFile(testFile, []byte("test content"), 0644)
			Expect(err).NotTo(HaveOccurred())

			oldCwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			defer os.Chdir(oldCwd)

			err = os.Chdir(tempDir)
			Expect(err).NotTo(HaveOccurred())

			result := HasUncommittedChanges()
			Expect(result).To(BeTrue())
		})
	})

	Describe("GetDefaultBranch", func() {
		It("returns main or master for a new repository", func() {
			cmd := exec.Command("git", "init")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "config", "user.email", "test@example.com")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "config", "user.name", "Test User")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "commit", "--allow-empty", "-m", "initial")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			oldCwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			defer os.Chdir(oldCwd)

			err = os.Chdir(tempDir)
			Expect(err).NotTo(HaveOccurred())

			branch, err := GetDefaultBranch()
			Expect(err).NotTo(HaveOccurred())
			Expect(branch).To(Or(Equal("main"), Equal("master")))
		})
	})

	Describe("BranchExists", func() {
		It("returns true for existing branch", func() {
			cmd := exec.Command("git", "init")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "config", "user.email", "test@example.com")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "config", "user.name", "Test User")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "commit", "--allow-empty", "-m", "initial")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "branch", "feature-test")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			oldCwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			defer os.Chdir(oldCwd)

			err = os.Chdir(tempDir)
			Expect(err).NotTo(HaveOccurred())

			result := BranchExists("feature-test")
			Expect(result).To(BeTrue())
		})

		It("returns false for non-existing branch", func() {
			cmd := exec.Command("git", "init")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "config", "user.email", "test@example.com")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "config", "user.name", "Test User")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "commit", "--allow-empty", "-m", "initial")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			oldCwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			defer os.Chdir(oldCwd)

			err = os.Chdir(tempDir)
			Expect(err).NotTo(HaveOccurred())

			result := BranchExists("nonexistent-branch")
			Expect(result).To(BeFalse())
		})

		It("is case-insensitive", func() {
			cmd := exec.Command("git", "init")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "config", "user.email", "test@example.com")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "config", "user.name", "Test User")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "commit", "--allow-empty", "-m", "initial")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "branch", "Feature-Test")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			oldCwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			defer os.Chdir(oldCwd)

			err = os.Chdir(tempDir)
			Expect(err).NotTo(HaveOccurred())

			// Should find it regardless of case
			result := BranchExists("feature-test")
			Expect(result).To(BeTrue())
		})
	})

	Describe("CreateBranch", func() {
		It("creates a new branch from base", func() {
			cmd := exec.Command("git", "init")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "config", "user.email", "test@example.com")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "config", "user.name", "Test User")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "commit", "--allow-empty", "-m", "initial")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			oldCwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			defer os.Chdir(oldCwd)

			err = os.Chdir(tempDir)
			Expect(err).NotTo(HaveOccurred())

			// Get current branch (should be main or master)
			currentBranch, err := GetCurrentBranch()
			Expect(err).NotTo(HaveOccurred())

			err = CreateBranch("new-feature", currentBranch)
			Expect(err).NotTo(HaveOccurred())

			// Verify branch exists
			result := BranchExists("new-feature")
			Expect(result).To(BeTrue())
		})
	})

	Describe("SwitchBranch", func() {
		It("switches to an existing branch", func() {
			cmd := exec.Command("git", "init")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "config", "user.email", "test@example.com")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "config", "user.name", "Test User")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "commit", "--allow-empty", "-m", "initial")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "branch", "feature-branch")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			oldCwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			defer os.Chdir(oldCwd)

			err = os.Chdir(tempDir)
			Expect(err).NotTo(HaveOccurred())

			err = SwitchBranch("feature-branch")
			Expect(err).NotTo(HaveOccurred())

			// Verify we're on the right branch
			current, err := GetCurrentBranch()
			Expect(err).NotTo(HaveOccurred())
			Expect(current).To(Equal("feature-branch"))
		})
	})

	Describe("GetCurrentBranch", func() {
		It("returns the current branch name", func() {
			cmd := exec.Command("git", "init")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "config", "user.email", "test@example.com")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "config", "user.name", "Test User")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "commit", "--allow-empty", "-m", "initial")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			oldCwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			defer os.Chdir(oldCwd)

			err = os.Chdir(tempDir)
			Expect(err).NotTo(HaveOccurred())

			branch, err := GetCurrentBranch()
			Expect(err).NotTo(HaveOccurred())
			Expect(branch).To(Or(Equal("main"), Equal("master")))
		})

		It("returns the correct branch after switching", func() {
			cmd := exec.Command("git", "init")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "config", "user.email", "test@example.com")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "config", "user.name", "Test User")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "commit", "--allow-empty", "-m", "initial")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			cmd = exec.Command("git", "branch", "test-branch")
			cmd.Dir = tempDir
			Expect(cmd.Run()).NotTo(HaveOccurred())

			oldCwd, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			defer os.Chdir(oldCwd)

			err = os.Chdir(tempDir)
			Expect(err).NotTo(HaveOccurred())

			err = SwitchBranch("test-branch")
			Expect(err).NotTo(HaveOccurred())

			branch, err := GetCurrentBranch()
			Expect(err).NotTo(HaveOccurred())
			Expect(branch).To(Equal("test-branch"))
		})
	})
})
