package llmprovider

import (
	"context"
	"inspirate-consulting/internal/config"

	"google.golang.org/genai"
	"google.golang.org/genai/interactions/models/interactions"
	"google.golang.org/genai/interactions/models/operations"
)

type Gemini struct {
	APIKey string
}

func callGemini(ctx context.Context, cfg config.Config, provider LLMProvider) (*string, error) {

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  cfg.Gemini.APIKey,
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		return nil, err
	}

	resp, err := client.Interactions.Create(ctx, operations.CreateInteractionRequest{
		Body: operations.NewCreateInteractionRequestBody(
			interactions.CreateModelInteraction{
				Model: interactions.Model(cfg.Gemini.Model),
				Input: genai.Ptr(interactions.NewInteractionsInput(provider.UserInput)),
				Tools: tools,
				GenerationConfig: &interactions.GenerationConfig{
					ToolChoice: genai.Ptr(interactions.NewToolChoice(interactions.ToolChoiceTypeAny)),
				},
			},
		),
	})

	if err != nil {
		return nil, err
	}

	var res *string
	for _, step := range resp.Interaction.Steps {
		if step.FunctionCallStep != nil {
			fc := step.FunctionCallStep

			query, ok := fc.Arguments["query"].(string)
			if !ok || query == "" {
				return nil, err
			}

			if fc.Name == "search_tasks" {
				res, err = provider.searchTasks(query)
				if err != nil {
					return nil, err
				}
			}
			if fc.Name == "search_chats" {
				res, err = provider.searchChats(query)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return res, nil
}
