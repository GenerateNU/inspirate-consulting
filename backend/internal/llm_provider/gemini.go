package llmprovider

import (
	"context"
	"errors"
	"fmt"
	"inspirate-consulting/internal/config"
	"iter"
	"strings"

	"google.golang.org/adk/model"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/genai"
)

type Gemini struct {
	APIKey string
}

func callGemini(ctx context.Context, cfg config.Config, provider LLMProvider) (*string, error) {
	resp, err := createGeminiInteraction(ctx, cfg, provider.UserInput)
	if err != nil {
		return nil, err
	}

	if resp.Content == nil {
		return nil, errors.New("gemini returned no interaction")
	}

	// Run every tool Gemini asked for, collecting each result to send back.
	var results []*genai.Part
	var rawResults []string
	for _, part := range resp.Content.Parts {
		if part.FunctionCall == nil {
			continue
		}
		fc := part.FunctionCall

		var res *string
		query, ok := fc.Args["query"].(string)
		switch {
		case !ok || query == "":
			err = fmt.Errorf("%s was called without a query", fc.Name)
		case fc.Name == SearchTasksTool:
			res, err = provider.searchTasks(ctx, query)
		case fc.Name == SearchChatsTool:
			res, err = provider.searchChats(ctx, query)
		default:
			err = fmt.Errorf("unknown tool %q", fc.Name)
		}
		if err != nil {
			return nil, err
		}

		results = append(results, &genai.Part{FunctionResponse: &genai.FunctionResponse{
			ID:       fc.ID,
			Name:     fc.Name,
			Response: map[string]any{"result": *res},
		}})
		rawResults = append(rawResults, *res)
	}

	if len(results) == 0 {
		text := responseText(resp)
		return &text, nil
	}

	history := append(genai.Text(provider.UserInput), resp.Content)
	summary, err := summarizeGeminiResults(ctx, cfg, history, results)
	if err != nil {
		return nil, err
	}
	text := responseText(summary)
	if text == "" {
		// Fall back to the raw search results rather than returning nothing.
		raw := strings.Join(rawResults, "\n")
		return &raw, nil
	}
	return &text, nil
}

func newGeminiClient(ctx context.Context, cfg config.Config) (model.LLM, error) {
	return gemini.NewModel(ctx, cfg.Gemini.Model, &genai.ClientConfig{
		APIKey:  cfg.Gemini.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
}

// Sends the tool results back to Gemini, after the conversation so far, so it can answer
// the user with a summary. Tool calls are turned off so it must reply in text.
func summarizeGeminiResults(ctx context.Context, cfg config.Config, history []*genai.Content, results []*genai.Part) (*model.LLMResponse, error) {
	client, err := newGeminiClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return firstResponse(client.GenerateContent(ctx, &model.LLMRequest{
		Model:    cfg.Gemini.Model,
		Contents: append(history, genai.NewContentFromParts(results, genai.RoleUser)),
		Config: &genai.GenerateContentConfig{
			SystemInstruction: genai.NewContentFromText(systemPrompt, genai.RoleUser),
			Tools:             geminiTools(),
			ToolConfig: &genai.ToolConfig{
				FunctionCallingConfig: &genai.FunctionCallingConfig{Mode: genai.FunctionCallingConfigModeNone},
			},
		},
	}, false))
}

// Sends the user's input and tools to Gemini, which must respond with a tool call.
func createGeminiInteraction(ctx context.Context, cfg config.Config, input string) (*model.LLMResponse, error) {
	client, err := newGeminiClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return firstResponse(client.GenerateContent(ctx, &model.LLMRequest{
		Model:    cfg.Gemini.Model,
		Contents: genai.Text(input),
		Config: &genai.GenerateContentConfig{
			SystemInstruction: genai.NewContentFromText(systemPrompt, genai.RoleUser),
			Tools:             geminiTools(),
			ToolConfig: &genai.ToolConfig{
				FunctionCallingConfig: &genai.FunctionCallingConfig{Mode: genai.FunctionCallingConfigModeAny},
			},
		},
	}, false))
}

// Converts the tools in tools.go to the format ADK sends to Gemini.
func geminiTools() []*genai.Tool {
	var decls []*genai.FunctionDeclaration
	for _, t := range tools {
		decls = append(decls, &genai.FunctionDeclaration{
			Name:                 *t.Function.Name,
			Description:          *t.Function.Description,
			ParametersJsonSchema: t.Function.Parameters,
		})
	}
	return []*genai.Tool{{FunctionDeclarations: decls}}
}

// ADK streams responses; with streaming off there is exactly one.
func firstResponse(responses iter.Seq2[*model.LLMResponse, error]) (*model.LLMResponse, error) {
	for resp, err := range responses {
		return resp, err
	}
	return nil, errors.New("gemini returned no response")
}

// Joins the text parts of a response, skipping the model's thoughts.
func responseText(resp *model.LLMResponse) string {
	if resp == nil || resp.Content == nil {
		return ""
	}
	var b strings.Builder
	for _, part := range resp.Content.Parts {
		if !part.Thought {
			b.WriteString(part.Text)
		}
	}
	return b.String()
}
