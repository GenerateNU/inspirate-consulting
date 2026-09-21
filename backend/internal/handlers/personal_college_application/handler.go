package personalcollegeapplication

import storage "inspirate-consulting/internal/data"

type Handler struct {
	PersonalCollegeApplicationRepository storage.PersonalCollegeApplicationRepository
}

// Here is where we create the handler which has reference to the PersonalCollegeApplicationRepository
func NewHandler(personalCollegeApplicationRepository storage.PersonalCollegeApplicationRepository) *Handler {
	return &Handler{
		PersonalCollegeApplicationRepository: personalCollegeApplicationRepository,
	}
}