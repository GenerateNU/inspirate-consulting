package embeddingclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	testAccountID = "test-account"
	testAPIToken  = "test-token"
)

// newTestServer starts a fake Cloudflare API that replies with the given status and body.
func newTestServer(t *testing.T, status int, body any) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(body)
	}))
	t.Cleanup(server.Close)
	return server
}

func newTestClient(server *httptest.Server) *EmbeddingClient {
	return NewEmbeddingClientWithHTTP(server.Client(), server.URL, testAccountID, testAPIToken)
}

func successBody(vector []float32) map[string]any {
	return map[string]any{
		"success":  true,
		"errors":   []any{},
		"messages": []any{},
		"result": map[string]any{
			"shape": []int{1, len(vector)},
			"data":  [][]float32{vector},
		},
	}
}

func makeVector(n int) []float32 {
	v := make([]float32, n)
	for i := range v {
		v[i] = float32(i) / 1000
	}
	return v
}

func TestNewEmbeddingClient_MissingCredentials(t *testing.T) {
	cases := []struct {
		name      string
		accountID string
		apiToken  string
	}{
		{"missing account id", "", testAPIToken},
		{"missing api token", testAccountID, ""},
		{"missing both", "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client, err := NewEmbeddingClient(tc.accountID, tc.apiToken)
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if client != nil {
				t.Errorf("expected nil client, got %+v", client)
			}
		})
	}
}

func TestNewEmbeddingClient_Success(t *testing.T) {
	client, err := NewEmbeddingClient(testAccountID, testAPIToken)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if client.baseURL != defaultBaseURL {
		t.Errorf("expected base URL %q, got %q", defaultBaseURL, client.baseURL)
	}
}

func TestEmbed_Success(t *testing.T) {
	expected := makeVector(Dimensions)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		wantPath := "/accounts/" + testAccountID + "/ai/run/" + model
		if r.URL.Path != wantPath {
			t.Errorf("expected path %q, got %q", wantPath, r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer "+testAPIToken {
			t.Errorf("expected bearer token header, got %q", got)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("expected JSON content type, got %q", got)
		}

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("failed to decode request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if body["text"] != "hello world" {
			t.Errorf("expected text %q, got %v", "hello world", body["text"])
		}
		if len(body) != 1 {
			t.Errorf("expected only the text field in the request, got %v", body)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(successBody(expected))
	}))
	defer server.Close()

	got, err := newTestClient(server).Embed(context.Background(), "hello world")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(got) != Dimensions {
		t.Fatalf("expected %d dimensions, got %d", Dimensions, len(got))
	}
	for i := range expected {
		if got[i] != expected[i] {
			t.Fatalf("embedding mismatch at index %d: expected %v, got %v", i, expected[i], got[i])
		}
	}
}

func TestEmbed_EmptyText(t *testing.T) {
	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer server.Close()

	_, err := newTestClient(server).Embed(context.Background(), "")
	if err == nil {
		t.Fatal("expected an error for empty text, got nil")
	}
	if called {
		t.Error("expected no request to be sent for empty text")
	}
}

func TestEmbed_APIErrorIncludesMessage(t *testing.T) {
	server := newTestServer(t, http.StatusUnauthorized, map[string]any{
		"success": false,
		"errors":  []map[string]any{{"code": 10000, "message": "Authentication error"}},
		"result":  nil,
	})

	_, err := newTestClient(server).Embed(context.Background(), "hello")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "401") || !strings.Contains(err.Error(), "Authentication error") {
		t.Errorf("expected error to include status and API message, got %q", err.Error())
	}
}

func TestEmbed_APIErrorWithoutMessage(t *testing.T) {
	server := newTestServer(t, http.StatusInternalServerError, map[string]any{
		"success": false,
		"errors":  []any{},
	})

	_, err := newTestClient(server).Embed(context.Background(), "hello")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("expected error to include status 500, got %q", err.Error())
	}
}

func TestEmbed_SuccessFalseWith200(t *testing.T) {
	body := successBody(makeVector(Dimensions))
	body["success"] = false
	server := newTestServer(t, http.StatusOK, body)

	_, err := newTestClient(server).Embed(context.Background(), "hello")
	if err == nil {
		t.Fatal("expected an error when success is false, got nil")
	}
}

func TestEmbed_WrongDimensions(t *testing.T) {
	server := newTestServer(t, http.StatusOK, successBody(makeVector(768)))

	_, err := newTestClient(server).Embed(context.Background(), "hello")
	if err == nil {
		t.Fatal("expected an error for a 768-dimension vector, got nil")
	}
}

func TestEmbed_NoVectors(t *testing.T) {
	server := newTestServer(t, http.StatusOK, map[string]any{
		"success": true,
		"result":  map[string]any{"shape": []int{0}, "data": [][]float32{}},
	})

	_, err := newTestClient(server).Embed(context.Background(), "hello")
	if err == nil {
		t.Fatal("expected an error when no vectors are returned, got nil")
	}
}

func TestEmbed_MultipleVectors(t *testing.T) {
	server := newTestServer(t, http.StatusOK, map[string]any{
		"success": true,
		"result": map[string]any{
			"shape": []int{2, Dimensions},
			"data":  [][]float32{makeVector(Dimensions), makeVector(Dimensions)},
		},
	})

	_, err := newTestClient(server).Embed(context.Background(), "hello")
	if err == nil {
		t.Fatal("expected an error when more than one vector is returned, got nil")
	}
}

func TestEmbed_NonJSONResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte("<html>502 Bad Gateway</html>"))
	}))
	defer server.Close()

	_, err := newTestClient(server).Embed(context.Background(), "hello")
	if err == nil {
		t.Fatal("expected an error for a non-JSON response, got nil")
	}
	if !strings.Contains(err.Error(), "502") {
		t.Errorf("expected error to include status 502, got %q", err.Error())
	}
}

func TestEmbed_CanceledContext(t *testing.T) {
	server := newTestServer(t, http.StatusOK, successBody(makeVector(Dimensions)))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := newTestClient(server).Embed(ctx, "hello")
	if err == nil {
		t.Fatal("expected an error for a canceled context, got nil")
	}
}
