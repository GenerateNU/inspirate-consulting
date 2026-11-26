package routes

import (
	"net/http"
	"strings"
	"testing"

	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/data"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// With AGENT_ENABLED off, /llm/ask answers 404 before reaching the model.
func TestRoute_AskLLM_FeatureDisabled(t *testing.T) {
	t.Parallel()

	app, _, err := SetupApp(config.Config{TestMode: true}, &data.Repository{}, nil, nil)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/llm/ask", strings.NewReader(`{"input":"hi"}`))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
