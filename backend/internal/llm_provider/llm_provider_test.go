package llmprovider

import (
	"context"
	"os"
	"testing"

	"inspirate-consulting/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModelCall_UnsupportedModel(t *testing.T) {
	provider := LLMProvider{UserInput: "hello", Model: "not-a-model"}

	_, err := provider.ModelCall(context.Background(), config.Config{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported model")
}

func TestModelCall_SelectsGemini(t *testing.T) {
	// A cancelled context stops the Gemini request before it leaves the machine
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	provider := LLMProvider{UserInput: "hello", Model: GeminiModel}
	cfg := config.Config{Gemini: config.Gemini{APIKey: "test-key", Model: "test-model"}}

	_, err := provider.ModelCall(ctx, cfg)

	require.Error(t, err)
	assert.NotContains(t, err.Error(), "unsupported model")
}

// Calls the real Gemini API. Runs only when GEMINI_API_KEY is set.
func TestGemini_RoutesQueryToTool(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		t.Skip("GEMINI_API_KEY not set")
	}
	model := os.Getenv("GEMINI_MODEL")
	if model == "" {
		model = "gemini-4-hydrogen"
	}
	cfg := config.Config{Gemini: config.Gemini{APIKey: apiKey, Model: model}}

	tests := []struct {
		input string
		want  string
	}{
		{"find chats where a student talks about hating literature class", SearchChatsTool},
		{"what did my student say about their college essay in our messages?", SearchChatsTool},
		{"which tasks is my student behind on?", SearchTasksTool},
		{"show me tasks that are due next week", SearchTasksTool},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			resp, err := createGeminiInteraction(context.Background(), cfg, tt.input)
			require.NoError(t, err)
			require.NotNil(t, resp.Interaction)

			var got []string
			for _, step := range resp.Interaction.Steps {
				if step.FunctionCallStep != nil {
					got = append(got, step.FunctionCallStep.Name)
				}
			}

			require.NotEmpty(t, got, "expected Gemini to call a tool")
			assert.Equal(t, tt.want, got[0])
		})
	}
}
