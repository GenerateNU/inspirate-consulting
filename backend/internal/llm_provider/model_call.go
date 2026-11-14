package llmprovider

import (
	"context"
	"inspirate-consulting/internal/config"
)

// Example: find me what student x said about science class sucking..
func (l *LLMProvider) modelCall(ctx context.Context, cfg config.Config, provider LLMProvider) (*string, error) {

	switch provider.Model {
	case GeminiModel:
		callGemini(ctx, cfg, provider)
	}

	return nil, nil
}
