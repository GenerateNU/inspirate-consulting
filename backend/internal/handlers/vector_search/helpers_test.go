package vectorsearch

import (
	"context"
	"errors"
	"testing"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeEmbeddingClient returns a fixed embedding (or error) and records how it was called.
type fakeEmbeddingClient struct {
	embedding []float32
	err       error
	calls     int
	gotQuery  string
}

func (f *fakeEmbeddingClient) EmbedQuery(ctx context.Context, query string) ([]float32, error) {
	f.calls++
	f.gotQuery = query
	return f.embedding, f.err
}

// Embed is part of data.EmbeddingClient but searches only use EmbedQuery, so fail loudly if it's called.
func (f *fakeEmbeddingClient) Embed(ctx context.Context, text string) ([]float32, error) {
	panic("search should embed prompts with EmbedQuery, not Embed")
}

func newFakeEmbeddingClient() *fakeEmbeddingClient {
	return &fakeEmbeddingClient{embedding: []float32{0.1, 0.2, 0.3}}
}

// assertHTTPStatus checks that err is an errs.HTTPError with the given status code.
func assertHTTPStatus(t *testing.T, err error, status int) {
	t.Helper()

	var httpErr errs.HTTPError
	require.True(t, errors.As(err, &httpErr), "expected errs.HTTPError, got %T: %v", err, err)
	assert.Equal(t, status, httpErr.Code)
}

// withRole sets the role the way the auth middleware stores it.
func withRole(ctx context.Context, role models.Role) context.Context {
	return context.WithValue(ctx, "Role", string(role)) //nolint:staticcheck // matches how auth middleware stores the role
}

func studentCtx(studentID uuid.UUID) context.Context {
	return withRole(auth.WithStudentID(context.Background(), studentID), models.StudentRole)
}

func counselorCtx(counselorID uuid.UUID) context.Context {
	return withRole(auth.WithCounselorID(context.Background(), counselorID), models.CounselorRole)
}
