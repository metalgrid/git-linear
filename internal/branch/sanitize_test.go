package branch

import (
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSanitize(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sanitize Suite")
}

var _ = Describe("Sanitize", func() {
	It("converts identifier to lowercase", func() {
		Expect(Sanitize("DEV-123", "foo")).To(Equal("dev-123-foo"))
	})

	It("slugifies title with spaces", func() {
		Expect(Sanitize("DEV-1", "Fix Login Bug")).To(Equal("dev-1-Fix-Login-Bug"))
	})

	It("removes special characters", func() {
		Expect(Sanitize("DEV-1", "Hello! @World#")).To(Equal("dev-1-Hello-World"))
	})

	It("removes emoji and unicode", func() {
		Expect(Sanitize("DEV-1", "Fix 🔐 Auth")).To(Equal("dev-1-Fix-Auth"))
	})

	It("collapses multiple hyphens", func() {
		Expect(Sanitize("DEV-1", "a - - b")).To(Equal("dev-1-a-b"))
	})

	It("truncates to max 32 chars", func() {
		long := strings.Repeat("a", 100)
		result := Sanitize("DEV-123", long)
		Expect(len(result)).To(BeNumerically("<=", 32))
		Expect(result).NotTo(HaveSuffix("-"))
	})

	It("returns just identifier if title sanitizes to empty", func() {
		Expect(Sanitize("DEV-1", "!@#$%")).To(Equal("dev-1"))
	})

	It("preserves uppercase letters", func() {
		Expect(Sanitize("DEV-1", "Fix Login Bug")).To(Equal("dev-1-Fix-Login-Bug"))
	})

	It("allows underscores", func() {
		Expect(Sanitize("DEV-1", "fix_login_bug")).To(Equal("dev-1-fix_login_bug"))
	})

	It("allows dots", func() {
		Expect(Sanitize("DEV-1", "v1.0.0")).To(Equal("dev-1-v1.0.0"))
	})

	It("allows slashes", func() {
		Expect(Sanitize("DEV-1", "feature/login")).To(Equal("dev-1-feature/login"))
	})

	It("blocks consecutive dots", func() {
		Expect(Sanitize("DEV-1", "v1..0")).To(Equal("dev-1-v1-0"))
	})

	It("still removes invalid Git characters", func() {
		Expect(Sanitize("DEV-1", "Fix~Login@Bug")).To(Equal("dev-1-FixLoginBug"))
	})
})
