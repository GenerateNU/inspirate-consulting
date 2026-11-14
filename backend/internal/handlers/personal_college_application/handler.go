package personalcollegeapplication

import (
	"inspirate-consulting/internal/auth"
	storage "inspirate-consulting/internal/data"
)

type Handler struct {
	PersonalCollegeApplicationRepository storage.PersonalCollegeApplicationRepository
	GlobalCollegeRepository              storage.GlobalCollegeRepository
	Authorizer                           *auth.Authorizer
}

// Here is where we create the handler which has reference to the PersonalCollegeApplicationRepository
func NewHandler(personalCollegeApplicationRepository storage.PersonalCollegeApplicationRepository, globalCollegeRepository storage.GlobalCollegeRepository, authorizer *auth.Authorizer) *Handler {
	return &Handler{
		PersonalCollegeApplicationRepository: personalCollegeApplicationRepository,
		GlobalCollegeRepository:              globalCollegeRepository,
		Authorizer:                           authorizer,
	}
}
