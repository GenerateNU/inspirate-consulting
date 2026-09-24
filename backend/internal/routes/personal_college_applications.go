package routes

import (
	"context"
	"net/http"

	"inspirate-consulting/internal/data"
	personalcollegeapplication "inspirate-consulting/internal/handlers/personal_college_application"
	"inspirate-consulting/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

// SetUpPersonalCollegeApplicationRoutes registers the personal college application endpoints:
// create and list by student ID
func SetUpPersonalCollegeApplicationRoutes(api huma.API, repository *data.Repository) {
	personalCollegeApplicationHandler := personalcollegeapplication.NewHandler(repository.PersonalCollegeApplication, repository.GlobalCollege)

	// Register POST /applications handler.
	huma.Register(api, huma.Operation{
		OperationID: "create-personal-college-application",
		Method:      http.MethodPost,
		Path:        "/applications",
		Description: "Create a personal college application for the authenticated student.",
		Tags:        []string{"Personal College Applications"},
	}, func(ctx context.Context, input *models.CreatePersonalCollegeApplicationInput) (*models.CreatePersonalCollegeApplicationOutput, error) {
		created, err := personalCollegeApplicationHandler.CreatePersonalCollegeApplication(ctx, input.Body)
		if err != nil {
			return nil, err
		}
		return &models.CreatePersonalCollegeApplicationOutput{Body: *created}, nil
	})

	// Register GET /applications handler.
	huma.Register(api, huma.Operation{
		OperationID: "list-personal-college-applications",
		Method:      http.MethodGet,
		Path:        "/applications",
		Description: "List all college applications for the authenticated student.",
		Tags:        []string{"Personal College Applications"},
	}, func(ctx context.Context, input *models.ListPersonalCollegeApplicationsInput) (*models.ListPersonalCollegeApplicationsOutput, error) {
		applications, err := personalCollegeApplicationHandler.ListPersonalCollegeApplicationsByStudentID(ctx)
		if err != nil {
			return nil, err
		}
		return &models.ListPersonalCollegeApplicationsOutput{Body: applications}, nil
	})

}
