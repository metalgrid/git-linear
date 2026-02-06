package linear_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/user/git-linear/internal/linear"
)

func TestLinear(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Linear Suite")
}

var _ = Describe("Linear Client", func() {
	var (
		client *linear.Client
		server *httptest.Server
	)

	AfterEach(func() {
		if server != nil {
			server.Close()
		}
	})

	Describe("NewClient", func() {
		It("should create a new client with API key", func() {
			client = linear.NewClient("test-api-key")
			Expect(client).NotTo(BeNil())
		})
	})

	Describe("GetAssignedIssues", func() {
		Context("when API returns successful response", func() {
			BeforeEach(func() {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Verify request
					Expect(r.Method).To(Equal("POST"))
					Expect(r.Header.Get("Authorization")).To(Equal("test-api-key"))
					Expect(r.Header.Get("Content-Type")).To(Equal("application/json"))

					// Return mock response
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{
						"data": {
							"viewer": {
								"assignedIssues": {
									"nodes": [
										{
											"id": "issue-1",
											"identifier": "GIT-1",
											"title": "Implement Linear client",
											"state": {
												"name": "In Progress",
												"type": "started"
											}
										},
										{
											"id": "issue-2",
											"identifier": "GIT-2",
											"title": "Add tests",
											"state": {
												"name": "Todo",
												"type": "unstarted"
											}
										}
									]
								}
							}
						}
					}`))
				}))

				client = linear.NewClientWithURL("test-api-key", server.URL)
			})

			It("should return assigned issues", func() {
				issues, err := client.GetAssignedIssues()
				Expect(err).NotTo(HaveOccurred())
				Expect(issues).To(HaveLen(2))

				Expect(issues[0].ID).To(Equal("issue-1"))
				Expect(issues[0].Identifier).To(Equal("GIT-1"))
				Expect(issues[0].Title).To(Equal("Implement Linear client"))
				Expect(issues[0].State.Name).To(Equal("In Progress"))
				Expect(issues[0].State.Type).To(Equal("started"))

				Expect(issues[1].ID).To(Equal("issue-2"))
				Expect(issues[1].Identifier).To(Equal("GIT-2"))
				Expect(issues[1].Title).To(Equal("Add tests"))
				Expect(issues[1].State.Name).To(Equal("Todo"))
				Expect(issues[1].State.Type).To(Equal("unstarted"))
			})
		})

		Context("when API returns empty response", func() {
			BeforeEach(func() {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{
						"data": {
							"viewer": {
								"assignedIssues": {
									"nodes": []
								}
							}
						}
					}`))
				}))

				client = linear.NewClientWithURL("test-api-key", server.URL)
			})

			It("should return empty slice", func() {
				issues, err := client.GetAssignedIssues()
				Expect(err).NotTo(HaveOccurred())
				Expect(issues).To(BeEmpty())
				Expect(issues).NotTo(BeNil())
			})
		})

		Context("when API returns 401 unauthorized", func() {
			BeforeEach(func() {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"errors":[{"message":"Invalid API key"}]}`))
				}))

				client = linear.NewClientWithURL("invalid-key", server.URL)
			})

			It("should return authentication error", func() {
				issues, err := client.GetAssignedIssues()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("authentication failed"))
				Expect(issues).To(BeNil())
			})
		})

		Context("when API returns malformed JSON", func() {
			BeforeEach(func() {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{invalid json`))
				}))

				client = linear.NewClientWithURL("test-api-key", server.URL)
			})

			It("should return JSON parsing error", func() {
				issues, err := client.GetAssignedIssues()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to decode response"))
				Expect(issues).To(BeNil())
			})
		})

		Context("when network error occurs", func() {
			BeforeEach(func() {
				// Use invalid URL to simulate network error
				client = linear.NewClientWithURL("test-api-key", "http://invalid-host-that-does-not-exist:9999")
			})

			It("should return network error", func() {
				issues, err := client.GetAssignedIssues()
				Expect(err).To(HaveOccurred())
				Expect(issues).To(BeNil())
			})
		})
	})

	Describe("ValidateAPIKey", func() {
		Context("when API key is valid", func() {
			BeforeEach(func() {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{
						"data": {
							"viewer": {
								"assignedIssues": {
									"nodes": []
								}
							}
						}
					}`))
				}))

				client = linear.NewClientWithURL("valid-key", server.URL)
			})

			It("should return no error", func() {
				err := client.ValidateAPIKey()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when API key is invalid", func() {
			BeforeEach(func() {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"errors":[{"message":"Invalid API key"}]}`))
				}))

				client = linear.NewClientWithURL("invalid-key", server.URL)
			})

			It("should return authentication error", func() {
				err := client.ValidateAPIKey()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("authentication failed"))
			})
		})
	})
})
