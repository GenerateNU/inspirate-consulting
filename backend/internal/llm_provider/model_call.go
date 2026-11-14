package llmprovider

import (
	"context"
	"fmt"
	"inspirate-consulting/internal/config"
)

// Example: find me what student x said about science class sucking..
func (l *LLMProvider) ModelCall(ctx context.Context, cfg config.Config) (*string, error) {
	switch l.Model {
	case GeminiModel:
		return callGemini(ctx, cfg, *l)
	default:
		return nil, fmt.Errorf("unsupported model: %q", l.Model)
	}
}
