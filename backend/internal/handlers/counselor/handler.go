package counselor

import storage "inspirate-consulting/internal/data"

type Handler struct {
	CounselorRepository storage.CounselorRepository
}

func NewHandler(counselorRepository storage.CounselorRepository) *Handler {
	return &Handler{
		CounselorRepository: counselorRepository,
	}
}
