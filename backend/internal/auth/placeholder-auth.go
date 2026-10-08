package auth

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type userIDKey struct{}

func GetStudentID(ctx context.Context) string {
	return "00000000-0000-0000-0000-000000000002"
}

// WithUserID stores the current user's users.id on the context for GetUserID to read.
func WithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, userIDKey{}, id)
}

func GetUserID(ctx context.Context) string {
	if id, ok := ctx.Value(userIDKey{}).(string); ok && id != "" {
		return id
	}
	return "00000000-0000-0000-0000-000000000001"
}

// SHOULD ONLY BE USED IN TEST MODE: allows requests act as any user via the X-User-ID header.
func TestModeUserMiddleware(api huma.API, testMode bool) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		id := ctx.Header("X-User-ID")
		if !testMode || id == "" {
			next(ctx)
			return
		}
		if _, err := uuid.Parse(id); err != nil {
			if err := huma.WriteErr(api, ctx, http.StatusBadRequest, "X-User-ID must be a UUID"); err != nil {
				slog.Error("Failed to write error", "err", err)
			}
			return
		}
		next(huma.WithValue(ctx, userIDKey{}, id))
	}
}
