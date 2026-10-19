package extracurricular

import (
	"inspirate-consulting/internal/auth"
	storage "inspirate-consulting/internal/data"
)

// Handler is the application layer for extracurricular operations.
type Handler struct {
	ExtracurricularRepository storage.ExtracurricularRepository
	Authorizer                *auth.Authorizer
}

// NewHandler constructs a new extracurricular handler with a repository dependency.
func NewHandler(extracurricularRepository storage.ExtracurricularRepository, authorizer *auth.Authorizer) *Handler {
	return &Handler{
		ExtracurricularRepository: extracurricularRepository,
		Authorizer:                authorizer,
	}
}
