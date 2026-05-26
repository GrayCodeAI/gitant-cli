package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/GrayCodeAI/gitant-cli/internal/config"
)

// Client is an HTTP client for the gitant daemon API.
type Client struct {
	BaseURL string
	Token   string
}

// NewClient creates a client. Resolves URL and UCAN token from flag, env, then config.
func NewClient(url string) *Client {
	if url == "" {
		url = os.Getenv("GITANT_DAEMON_URL")
	}
	if url == "" {
		if s, err := config.Load(); err == nil && s.DaemonURL != "" {
			url = s.DaemonURL
		}
	}
	if url == "" {
		url = "http://localhost:7777"
	}

	token := os.Getenv("GITANT_UCAN_TOKEN")
	if token == "" {
		if s, err := config.Load(); err == nil {
			token = s.UCANToken
		}
	}

	return &Client{BaseURL: strings.TrimRight(url, "/"), Token: token}
}

func (c *Client) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	req, err := http.NewRequest(method, c.BaseURL+path, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	return req, nil
}

func (c *Client) do(req *http.Request, result interface{}) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("decoding response: %w", err)
		}
	}
	return nil
}

// Get performs a GET request and decodes the JSON response.
func (c *Client) Get(path string, result interface{}) error {
	req, err := c.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return err
	}
	return c.do(req, result)
}

// GetRaw performs GET and returns the raw response body.
func (c *Client) GetRaw(path string) ([]byte, error) {
	req, err := c.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}

// Post performs a POST request with a JSON body.
func (c *Client) Post(path string, body interface{}, result interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("encoding request: %w", err)
	}
	req, err := c.newRequest(http.MethodPost, path, bytes.NewReader(data))
	if err != nil {
		return err
	}
	return c.do(req, result)
}

// Put performs a PUT request with a JSON body.
func (c *Client) Put(path string, body interface{}, result interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("encoding request: %w", err)
	}
	req, err := c.newRequest(http.MethodPut, path, bytes.NewReader(data))
	if err != nil {
		return err
	}
	return c.do(req, result)
}

// Delete performs a DELETE request.
func (c *Client) Delete(path string) error {
	req, err := c.newRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// Request performs an arbitrary HTTP request with optional JSON body map.
func (c *Client) Request(method, path string, body interface{}, result interface{}) error {
	method = strings.ToUpper(method)
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("encoding request: %w", err)
		}
		reader = bytes.NewReader(data)
	}
	req, err := c.newRequest(method, path, reader)
	if err != nil {
		return err
	}
	return c.do(req, result)
}
