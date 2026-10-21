package counselor

import storage "inspirate-consulting/internal/data"

type Handler struct {
	CounselorRepository storage.CounselorRepository
}

// Here is where we create the handler which has reference to the CounselorRepository
func NewHandler(counselorRepository storage.CounselorRepository) *Handler {
	return &Handler{
		CounselorRepository: counselorRepository,
	}
}
