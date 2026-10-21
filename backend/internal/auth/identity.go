package auth

import (
	"context"
	"embed"
	"errors"
	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed sql/*.sql
var sqlAuthFiles embed.FS

type userIDKey struct{}
type studentIDKey struct{}
type counselorIDKey struct{}

var (
	ErrNoUserID      = errors.New("user ID not found in context")
	ErrNoStudentID   = errors.New("student profile not found in context")
	ErrNoCounselorID = errors.New("counselor profile not found in context")
)

func WithUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey{}, id)
}

func WithStudentID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, studentIDKey{}, id)
}

func WithCounselorID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, counselorIDKey{}, id)
}

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	if id, ok := ctx.Value(userIDKey{}).(uuid.UUID); ok && id != uuid.Nil {
		return id, nil
	}
	return uuid.Nil, ErrNoUserID
}

func GetStudentID(ctx context.Context) (uuid.UUID, error) {
	if id, ok := ctx.Value(studentIDKey{}).(uuid.UUID); ok && id != uuid.Nil {
		return id, nil
	}
	return uuid.Nil, ErrNoStudentID
}

func GetCounselorID(ctx context.Context) (uuid.UUID, error) {
	if id, ok := ctx.Value(counselorIDKey{}).(uuid.UUID); ok && id != uuid.Nil {
		return id, nil
	}
	return uuid.Nil, ErrNoCounselorID
}

func ResolveIdentityMiddleware(api huma.API, db *pgxpool.Pool) func(ctx huma.Context, next func(huma.Context)) {
	query, err := schema.ReadSQLBaseScript("resolve_identity.sql", sqlAuthFiles)
	if err != nil {
		panic("failed to read resolve_identity.sql: " + err.Error())
	}

	return func(ctx huma.Context, next func(huma.Context)) {
		supabaseID, ok := ctx.Context().Value("Supabase-ID").(string)
		if !ok || supabaseID == "" {
			next(ctx)
			return
		}

		sid, err := uuid.Parse(supabaseID)
		if err != nil {
			slog.Error("invalid Supabase-ID in context", "value", supabaseID)
			_ = huma.WriteErr(api, ctx, http.StatusInternalServerError, "invalid identity")
			return
		}

		role, _ := ctx.Context().Value("Role").(string)

		var userID uuid.UUID
		var studentID, counselorID *uuid.UUID
		err = db.QueryRow(ctx.Context(), query, sid).Scan(&userID, &studentID, &counselorID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				slog.Error("no local user for supabase ID", "supabase_id", sid)
				_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "user account not found")
				return
			}
			slog.Error("identity resolution failed", "err", err)
			_ = huma.WriteErr(api, ctx, http.StatusInternalServerError, "identity resolution failed")
			return
		}

		ctx = huma.WithValue(ctx, userIDKey{}, userID)

		if role == string(models.StudentRole) && studentID == nil {
			slog.Error("student profile missing", "user_id", userID)
			_ = huma.WriteErr(api, ctx, http.StatusForbidden, "student profile not found")
			return
		}
		if studentID != nil {
			ctx = huma.WithValue(ctx, studentIDKey{}, *studentID)
		}

		if role == string(models.CounselorRole) && counselorID == nil {
			slog.Error("counselor profile missing", "user_id", userID)
			_ = huma.WriteErr(api, ctx, http.StatusForbidden, "counselor profile not found")
			return
		}
		if counselorID != nil {
			ctx = huma.WithValue(ctx, counselorIDKey{}, *counselorID)
		}

		next(ctx)
	}
}

const (
	testModeDefaultUserID    = "00000000-0000-0000-0000-000000000001"
	testModeDefaultStudentID = "00000000-0000-0000-0000-000000000002"
)

// TestModeUserMiddleware sets identity from X-* headers, defaulting to placeholder user and student IDs when they are omitted.
func TestModeUserMiddleware(api huma.API) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		userID := ctx.Header("X-User-ID")
		studentID := ctx.Header("X-Student-ID")
		counselorID := ctx.Header("X-Counselor-ID")

		if userID == "" {
			userID = testModeDefaultUserID
		}
		if studentID == "" && counselorID == "" {
			studentID = testModeDefaultStudentID
		}

		parsedUserID, err := uuid.Parse(userID)
		if err != nil {
			_ = huma.WriteErr(api, ctx, http.StatusBadRequest, "X-User-ID must be a UUID")
			return
		}
		ctx = huma.WithValue(ctx, userIDKey{}, parsedUserID)

		if studentID != "" {
			parsed, err := uuid.Parse(studentID)
			if err != nil {
				_ = huma.WriteErr(api, ctx, http.StatusBadRequest, "X-Student-ID must be a UUID")
				return
			}
			ctx = huma.WithValue(ctx, studentIDKey{}, parsed)
		}

		if counselorID != "" {
			parsed, err := uuid.Parse(counselorID)
			if err != nil {
				_ = huma.WriteErr(api, ctx, http.StatusBadRequest, "X-Counselor-ID must be a UUID")
				return
			}
			ctx = huma.WithValue(ctx, counselorIDKey{}, parsed)
		}

		next(ctx)
	}
}
