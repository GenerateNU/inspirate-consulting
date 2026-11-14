package embeddingclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	model = "@cf/qwen/qwen3-embedding-0.6b"
	// Dimensions must match vector(1024) on public.vector_table.
	Dimensions     = 1024
	defaultBaseURL = "https://api.cloudflare.com/client/v4"
)

type EmbeddingClient struct {
	httpClient *http.Client
	baseURL    string
	accountID  string
	apiToken   string
}

func NewEmbeddingClient(accountID, apiToken string) (*EmbeddingClient, error) {
	if accountID == "" || apiToken == "" {
		return nil, fmt.Errorf("embeddingclient: accountID and apiToken required")
	}
	return &EmbeddingClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    defaultBaseURL,
		accountID:  accountID,
		apiToken:   apiToken,
	}, nil
}

// builds an EmbeddingClient against a custom base URL and HTTP client, so it can be faked in tests.
func NewEmbeddingClientWithHTTP(httpClient *http.Client, baseURL, accountID, apiToken string) *EmbeddingClient {
	return &EmbeddingClient{httpClient: httpClient, baseURL: baseURL, accountID: accountID, apiToken: apiToken}
}

type embedRequest struct {
	Text string `json:"text"`
}

type embedResponse struct {
	Success bool `json:"success"`
	Errors  []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
	Result struct {
		Shape []int       `json:"shape"`
		Data  [][]float32 `json:"data"`
	} `json:"result"`
}

// Embed returns the 1024-dimensional embedding for a single string.
func (c *EmbeddingClient) Embed(ctx context.Context, text string) ([]float32, error) {
	if text == "" {
		return nil, fmt.Errorf("embeddingclient: text must not be empty")
	}

	body, err := json.Marshal(embedRequest{Text: text})
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/accounts/%s/ai/run/%s", c.baseURL, c.accountID, model)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("embeddingclient: request failed: %w", err)
	}
	defer resp.Body.Close()

	var parsed embedResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("embeddingclient: decode response (status %d): %w", resp.StatusCode, err)
	}

	if resp.StatusCode != http.StatusOK || !parsed.Success {
		if len(parsed.Errors) > 0 {
			return nil, fmt.Errorf("embeddingclient: status %d: %s", resp.StatusCode, parsed.Errors[0].Message)
		}
		return nil, fmt.Errorf("embeddingclient: status %d", resp.StatusCode)
	}

	if len(parsed.Result.Data) != 1 || len(parsed.Result.Data[0]) != Dimensions {
		return nil, fmt.Errorf("embeddingclient: unexpected embedding shape %v", parsed.Result.Shape)
	}

	return parsed.Result.Data[0], nil
}
