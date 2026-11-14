package llmprovider

import (
	"inspirate-consulting/internal/config"

	"github.com/danielgtaylor/huma/v2"
)

// Example: find me what student x said about science class sucking..
func (l *LLMProvider) modelCall(ctx huma.Context, cfg config.Config, provider LLMProvider) (*string, error) {

	switch provider.Model {
	case GeminiModel:
		callGemini(ctx, cfg, provider)
	}

	return nil, nil
}
