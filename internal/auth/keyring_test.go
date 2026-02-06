package auth

import (
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/zalando/go-keyring"
)

func TestKeyring(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Keyring Suite")
}

var _ = ginkgo.Describe("Keyring", func() {
	ginkgo.BeforeEach(func() {
		// Initialize mock keyring for testing
		keyring.MockInit()
	})

	ginkgo.Describe("StoreAPIKey", func() {
		ginkgo.It("should store a valid API key", func() {
			err := StoreAPIKey("test-api-key-12345")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})

		ginkgo.It("should overwrite an existing API key", func() {
			err := StoreAPIKey("first-key")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			err = StoreAPIKey("second-key")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			retrieved, err := GetAPIKey()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(retrieved).To(gomega.Equal("second-key"))
		})
	})

	ginkgo.Describe("GetAPIKey", func() {
		ginkgo.It("should retrieve a stored API key", func() {
			err := StoreAPIKey("my-secret-key")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			retrieved, err := GetAPIKey()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(retrieved).To(gomega.Equal("my-secret-key"))
		})

		ginkgo.It("should return ErrNotFound when key does not exist", func() {
			// Ensure no key is stored
			_ = DeleteAPIKey()

			_, err := GetAPIKey()
			gomega.Expect(err).To(gomega.Equal(keyring.ErrNotFound))
		})
	})

	ginkgo.Describe("DeleteAPIKey", func() {
		ginkgo.It("should delete a stored API key", func() {
			err := StoreAPIKey("key-to-delete")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			err = DeleteAPIKey()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			_, err = GetAPIKey()
			gomega.Expect(err).To(gomega.Equal(keyring.ErrNotFound))
		})

		ginkgo.It("should not error when deleting non-existent key", func() {
			_ = DeleteAPIKey() // Clean up first
			err := DeleteAPIKey()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		})
	})

	ginkgo.Describe("HasAPIKey", func() {
		ginkgo.It("should return true when API key exists", func() {
			err := StoreAPIKey("existing-key")
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			gomega.Expect(HasAPIKey()).To(gomega.BeTrue())
		})

		ginkgo.It("should return false when API key does not exist", func() {
			_ = DeleteAPIKey()
			gomega.Expect(HasAPIKey()).To(gomega.BeFalse())
		})
	})
})
