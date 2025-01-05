package graphqlclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Client represents the main GraphQL client
type Client struct {
	mu         sync.RWMutex
	endpoint   string
	httpClient *http.Client
	cache      *Cache
	provider   Provider
}

// NewClient creates a new GraphQL client
func NewClient() *Client {
	client := &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		cache:      NewCache(),
	}

	return client
}

// Option allows for flexible client configuration
type Option func(*Client)

// SetgqlEndpoint sets the GraphQL endpoint with thread-safe mechanism
func (g *Graph) SetgqlEndpoint(endpoint string) {
	g.gqlClient.mu.Lock()
	defer g.gqlClient.mu.Unlock()
	g.gqlClient.endpoint = endpoint
}

// WithHTTPClient allows setting a custom HTTP client
func (g *Graph) WithHTTPClient(httpClient *http.Client) Option {
	g.gqlClient.mu.Lock()
	defer g.gqlClient.mu.Unlock()
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithProvider sets a custom authentication provider
func (g *Graph) WithProvider(provider Provider) Option {
	g.gqlClient.mu.Lock()
	defer g.gqlClient.mu.Unlock()
	return func(c *Client) {
		c.provider = provider
	}
}

// Execute performs a GraphQL request
func (c *Graph) Execute(ctx context.Context, query string, variables map[string]interface{}, result interface{}) error {
	// Check cache first
	cacheKey := generateCacheKey(query, variables)
	if cachedData, found := c.gqlClient.cache.Get(cacheKey); found {
		return json.Unmarshal(cachedData, result)
	}

	// Prepare request
	requestBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}
	// Add authentication if provider exists
	req, err := http.NewRequestWithContext(ctx, "POST", c.gqlClient.endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Add authentication
	if c.gqlClient.provider != nil {
		c.gqlClient.provider.Authenticate(req)
	}

	// Send request
	resp, err := c.gqlClient.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Parse response
	var response struct {
		Data   json.RawMessage `json:"data"`
		Errors []GraphQLError  `json:"errors"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	// Check for GraphQL errors
	if len(response.Errors) > 0 {
		return fmt.Errorf("graphql errors: %v", response.Errors)
	}

	// Cache the result
	c.gqlClient.cache.Set(cacheKey, response.Data)

	// Unmarshal result
	return json.Unmarshal(response.Data, result)
}
