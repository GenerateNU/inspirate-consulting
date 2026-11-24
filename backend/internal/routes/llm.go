package routes

import (
	"context"
	"net/http"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/data"
	vectorsearch "inspirate-consulting/internal/handlers/vector_search"
	llmprovider "inspirate-consulting/internal/llm_provider"
	"inspirate-consulting/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

func SetUpLLMRoutes(api huma.API, repository *data.Repository, cfg *config.Config, verifyRole auth.RoleVerifier) {
	searchHandler := vectorsearch.NewHandler(repository.Vector, repository.Counselor, repository.Embedding)

	huma.Register(api, huma.Operation{
		OperationID: "ask-llm",
		Middlewares: huma.Middlewares{
			verifyRole(api, models.CounselorRole),
			auth.RequireFeature(api, cfg.Features.Agent),
		},
		Method:      http.MethodPost,
		Path:        "/llm/ask",
		Description: "Ask the LLM a question. The model may search tasks or chats to answer it.",
		Tags:        []string{"LLM"},
	}, func(ctx context.Context, input *models.AskLLMInput) (*models.AskLLMOutput, error) {
		provider := llmprovider.LLMProvider{
			UserInput: input.Body.Input,
			Model:     llmprovider.GeminiModel,
			Search:    searchHandler,
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
