package globalcollege

import storage "inspirate-consulting/internal/data"

type Handler struct {
	GlobalCollegeRepository storage.GlobalCollegeRepository
}

// Here is where we create the handler which has reference to the GlobalCollegeRepository
func NewHandler(globalCollegeRepository storage.GlobalCollegeRepository) *Handler {
	return &Handler{
		GlobalCollegeRepository: globalCollegeRepository,
	}
}
