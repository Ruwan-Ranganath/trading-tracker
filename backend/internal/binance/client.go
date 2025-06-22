package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	baseURL    string
	apiKey     string
	secretKey  string
	httpClient *http.Client
}

func NewClient(apiKey, secretKey string) *Client {
	return &Client{
		baseURL:    "https://fapi.binance.com",
		apiKey:     apiKey,
		secretKey:  secretKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) doRequest(ctx context.Context, method, endpoint string, params map[string]string, result interface{}) error {
	// Build URL with query parameters
	urlStr := c.baseURL + endpoint
	if len(params) > 0 {
		q := url.Values{}
		for k, v := range params {
			q.Add(k, v)
		}
		urlStr = fmt.Sprintf("%s?%s", urlStr, q.Encode())
	}

	req, err := http.NewRequestWithContext(ctx, method, urlStr, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Add API key header for authenticated endpoints
	if c.apiKey != "" {
		req.Header.Add("X-MBX-APIKEY", c.apiKey)
	}

	log.Printf("Making %s request to %s", method, urlStr)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	log.Printf("Response status: %d, Body: %s", resp.StatusCode, string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if len(body) == 0 {
		return nil
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("error decoding response: %w. Response: %s", err, string(body))
	}

	return nil
}
