package extracurricular

import storage "inspirate-consulting/internal/data"

// Handler is the application layer for extracurricular operations.
type Handler struct {
	ExtracurricularRepository storage.ExtracurricularRepository
}

// NewHandler constructs a new extracurricular handler with a repository dependency.
func NewHandler(extracurricularRepository storage.ExtracurricularRepository) *Handler {
	return &Handler{
		ExtracurricularRepository: extracurricularRepository,
	}
}
