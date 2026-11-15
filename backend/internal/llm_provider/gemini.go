package llmprovider

import (
	"context"
	"errors"
	"fmt"
	"inspirate-consulting/internal/config"
	"strings"

	"google.golang.org/genai"
	"google.golang.org/genai/interactions/models/interactions"
	"google.golang.org/genai/interactions/models/operations"
)

type Gemini struct {
	APIKey string
}

func callGemini(ctx context.Context, cfg config.Config, provider LLMProvider) (*string, error) {
	resp, err := createGeminiInteraction(ctx, cfg, provider.UserInput)
	if err != nil {
		return nil, err
	}

	if resp.Interaction == nil {
		return nil, errors.New("gemini returned no interaction")
	}

	// Run every tool Gemini asked for, collecting each result to send back.
	var results []interactions.Step
	var rawResults []string
	for _, step := range resp.Interaction.Steps {
		if step.FunctionCallStep == nil {
			continue
		}
		fc := step.FunctionCallStep

		var res *string
		query, ok := fc.Arguments["query"].(string)
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

		results = append(results, interactions.NewStep(interactions.FunctionResultStep{
			CallID: fc.ID,
			Name:   genai.Ptr(fc.Name),
			Result: interactions.NewFunctionResultStepResultUnion(*res),
		}))
		rawResults = append(rawResults, *res)
	}

	if len(results) == 0 {
		return resp.Interaction.OutputText, nil
	}

	summary, err := summarizeGeminiResults(ctx, cfg, resp.Interaction.ID, results)
	if err != nil {
		return nil, err
	}
	if summary.Interaction == nil || summary.Interaction.OutputText == nil || *summary.Interaction.OutputText == "" {
		// Fall back to the raw search results rather than returning nothing.
		raw := strings.Join(rawResults, "\n")
		return &raw, nil
	}
	return summary.Interaction.OutputText, nil
}

func newGeminiClient(ctx context.Context, cfg config.Config) (*genai.Client, error) {
	return genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  cfg.Gemini.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
}

// Sends the tool results back to Gemini, continuing the first interaction, so it can answer
// the user with a summary. Tool calls are turned off so it must reply in text.
func summarizeGeminiResults(ctx context.Context, cfg config.Config, previousID *string, results []interactions.Step) (*operations.CreateInteractionResponse, error) {
	client, err := newGeminiClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return client.Interactions.Create(ctx, operations.CreateInteractionRequest{
		Body: operations.NewCreateInteractionRequestBody(
			interactions.CreateModelInteraction{
				Model:                 interactions.Model(cfg.Gemini.Model),
				SystemInstruction:     genai.Ptr(systemPrompt),
				PreviousInteractionID: previousID,
				Input:                 genai.Ptr(interactions.NewInteractionsInput(results)),
				Tools:                 tools,
				GenerationConfig: &interactions.GenerationConfig{
					ToolChoice: genai.Ptr(interactions.NewToolChoice(interactions.ToolChoiceTypeNone)),
				},
			},
		),
	})
}

// Sends the user's input and tools to Gemini, which must respond with a tool call.
func createGeminiInteraction(ctx context.Context, cfg config.Config, input string) (*operations.CreateInteractionResponse, error) {
	client, err := newGeminiClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return client.Interactions.Create(ctx, operations.CreateInteractionRequest{
		Body: operations.NewCreateInteractionRequestBody(
			interactions.CreateModelInteraction{
				Model:             interactions.Model(cfg.Gemini.Model),
				SystemInstruction: genai.Ptr(systemPrompt),
				Input:             genai.Ptr(interactions.NewInteractionsInput(input)),
				Tools:             tools,
				GenerationConfig: &interactions.GenerationConfig{
					ToolChoice: genai.Ptr(interactions.NewToolChoice(interactions.ToolChoiceTypeAny)),
				},
			},
		),
	})
}
