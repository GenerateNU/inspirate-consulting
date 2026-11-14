package routes

import (
	"context"
	"inspirate-consulting/internal/auth"
	"net/http"

	"inspirate-consulting/internal/data"
	"inspirate-consulting/internal/handlers/activity"
	"inspirate-consulting/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

// SetUpActivityRoutes registers the counselor-only activity feed endpoint.
func SetUpActivityRoutes(api huma.API, repository *data.Repository, verifyRole auth.RoleVerifier) {
	activityHandler := activity.NewHandler(repository.Activity)

	// Register GET /activity handler.
	huma.Register(api, huma.Operation{
		OperationID: "get-activities",
		Middlewares: huma.Middlewares{verifyRole(api, models.CounselorRole)},
		Method:      http.MethodGet,
		Path:        "/activity",
		Description: "Get activity across the authenticated counselor's students, newest first. Optionally scope to one student.",
		Tags:        []string{"Activity"},
	}, func(ctx context.Context, input *models.GetActivitiesInput) (*models.GetActivitiesOutput, error) {
		activities, err := activityHandler.GetActivities(ctx, input)
		if err != nil {
			return nil, err
		}
		return &models.GetActivitiesOutput{Body: activities}, nil
	})
}
