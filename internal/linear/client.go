package linear

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	defaultAPIURL = "https://api.linear.app/graphql"
	timeout       = 30 * time.Second
)

// Client represents a Linear API client
type Client struct {
	apiKey     string
	apiURL     string
	httpClient *http.Client
}

// NewClient creates a new Linear API client with the default API URL
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		apiURL: defaultAPIURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// NewClientWithURL creates a new Linear API client with a custom API URL (for testing)
func NewClientWithURL(apiKey, apiURL string) *Client {
	return &Client{
		apiKey: apiKey,
		apiURL: apiURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// graphQLRequest represents a GraphQL request
type graphQLRequest struct {
	Query string `json:"query"`
}

// graphQLResponse represents the top-level GraphQL response
type graphQLResponse struct {
	Data   *responseData  `json:"data"`
	Errors []graphQLError `json:"errors,omitempty"`
}

// graphQLError represents a GraphQL error
type graphQLError struct {
	Message string `json:"message"`
}

// responseData represents the data field in the GraphQL response
type responseData struct {
	Viewer *viewer `json:"viewer"`
}

// viewer represents the viewer field in the GraphQL response
type viewer struct {
	AssignedIssues *assignedIssues `json:"assignedIssues"`
}

// assignedIssues represents the assignedIssues field in the GraphQL response
type assignedIssues struct {
	Nodes []Issue `json:"nodes"`
}

// GetAssignedIssues fetches assigned issues from Linear API
func (c *Client) GetAssignedIssues() ([]Issue, error) {
	query := `
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
	`

	var response graphQLResponse
	if err := c.executeQuery(query, &response); err != nil {
		return nil, err
	}

	// Handle empty response
	if response.Data == nil || response.Data.Viewer == nil || response.Data.Viewer.AssignedIssues == nil {
		return []Issue{}, nil
	}

	issues := response.Data.Viewer.AssignedIssues.Nodes
	if issues == nil {
		return []Issue{}, nil
	}

	return issues, nil
}

// ValidateAPIKey validates the API key by making a simple query
func (c *Client) ValidateAPIKey() error {
	_, err := c.GetAssignedIssues()
	return err
}

// executeQuery executes a GraphQL query and decodes the response
func (c *Client) executeQuery(query string, response interface{}) error {
	reqBody := graphQLRequest{
		Query: query,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check for authentication errors
	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("authentication failed: invalid API key")
	}

	// Check for other HTTP errors
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode response
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
