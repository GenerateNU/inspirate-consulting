package personalcollegeapplication

import storage "inspirate-consulting/internal/data"

type Handler struct {
	PersonalCollegeApplicationRepository storage.PersonalCollegeApplicationRepository
	GlobalCollegeRepository storage.GlobalCollegeRepository
}

// Here is where we create the handler which has reference to the PersonalCollegeApplicationRepository
func NewHandler(personalCollegeApplicationRepository storage.PersonalCollegeApplicationRepository, globalCollegeRepository storage.GlobalCollegeRepository) *Handler {
	return &Handler{
		PersonalCollegeApplicationRepository: personalCollegeApplicationRepository,
		GlobalCollegeRepository: globalCollegeRepository,
	}
}