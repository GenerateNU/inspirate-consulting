package routes

import (
	"context"
	"inspirate-consulting/internal/auth"
	"net/http"

	"inspirate-consulting/internal/data"
	"inspirate-consulting/internal/handlers/counselor"
	"inspirate-consulting/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

func SetUpCounselorRoutes(api huma.API, repository *data.Repository, verifyRole auth.RoleVerifier) {
	counselorHandler := counselor.NewHandler(repository.Counselor)

	huma.Register(api, huma.Operation{
		OperationID: "fetch-counselors",
		Middlewares: huma.Middlewares{verifyRole(api, models.CounselorRole)},
		Method:      http.MethodGet,
		Path:        "/counselors/fetchall",
		Description: "Fetch a page of counselors, optionally filtered by user_id and name",
		Tags:        []string{"Counselors"},
	}, func(ctx context.Context, input *models.FetchCounselorsInput) (*models.FetchCounselorsOutput, error) {
		page, err := counselorHandler.ListCounselors(ctx, input)
		if err != nil {
			return nil, err
		}
		return &models.FetchCounselorsOutput{Body: page}, nil
	})

}
