package routes

import (
	"context"
	"net/http"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/config"
	llmprovider "inspirate-consulting/internal/llm_provider"
	"inspirate-consulting/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

func SetUpLLMRoutes(api huma.API, cfg *config.Config, verifyRole auth.RoleVerifier) {
	huma.Register(api, huma.Operation{
		OperationID: "ask-llm",
		Middlewares: huma.Middlewares{verifyRole(api, models.CounselorRole)},
		Method:      http.MethodPost,
		Path:        "/llm/ask",
		Description: "Ask the LLM a question. The model may search tasks or chats to answer it.",
		Tags:        []string{"LLM"},
	}, func(ctx context.Context, input *models.AskLLMInput) (*models.AskLLMOutput, error) {
		provider := llmprovider.LLMProvider{
			UserInput: input.Body.Input,
			Model:     llmprovider.GeminiModel,
		}

		answer, err := provider.ModelCall(ctx, *cfg)
		if err != nil {
			return nil, err
		}

		out := &models.AskLLMOutput{}
		if answer != nil {
			out.Body.Answer = *answer
		}
		return out, nil
	})
}
