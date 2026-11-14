package routes

import (
	"context"
	"net/http"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/data"
	"inspirate-consulting/internal/handlers/extracurricular"
	"inspirate-consulting/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

func SetUpExtracurricularRoutes(api huma.API, repository *data.Repository, verifyRole auth.RoleVerifier, authorizer *auth.Authorizer) {
	extracurricularHandler := extracurricular.NewHandler(repository.Extracurricular, authorizer)

	huma.Register(api, huma.Operation{
		OperationID: "create-extracurricular",
		Method:      http.MethodPost,
		Path:        "/extracurriculars",
		Description: "Create a new extracurricular for the authenticated student.",
		Tags:        []string{"Extracurriculars"},
	}, extracurricularHandler.CreateExtracurricular)

	huma.Register(api, huma.Operation{
		OperationID: "list-extracurriculars",
		Method:      http.MethodGet,
		Path:        "/extracurriculars",
		Description: "List extracurriculars for the authenticated student.",
		Tags:        []string{"Extracurriculars"},
	}, func(ctx context.Context, _ *models.ListExtracurricularsInput) (*models.ListExtracurricularsOutput, error) {
		list, err := extracurricularHandler.ListExtracurriculars(ctx)
		if err != nil {
			return nil, err
		}
		return &models.ListExtracurricularsOutput{Body: list}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-extracurricular",
		Method:      http.MethodPatch,
		Path:        "/extracurriculars/{id}",
		Description: "Update fields on an extracurricular for the authenticated student.",
		Tags:        []string{"Extracurriculars"},
	}, func(ctx context.Context, input *models.UpdateExtracurricularInput) (*models.UpdateExtracurricularOutput, error) {
		updated, err := extracurricularHandler.UpdateExtracurricular(ctx, input.ID, input)
		if err != nil {
			return nil, err
		}
		return &models.UpdateExtracurricularOutput{Body: *updated}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-extracurricular",
		Middlewares: huma.Middlewares{verifyRole(api, models.StudentRole, models.CounselorRole)},
		Method:      http.MethodGet,
		Path:        "/extracurriculars/{id}",
		Description: "Get an extracurricular with all of its achievements and notes. Students can read their own; counselors can read their assigned students'.",
		Tags:        []string{"Extracurriculars"},
	}, func(ctx context.Context, input *models.GetExtracurricularInput) (*models.GetExtracurricularOutput, error) {
		detail, err := extracurricularHandler.GetExtracurricular(ctx, input.ID)
		if err != nil {
			return nil, err
		}
		return &models.GetExtracurricularOutput{Body: *detail}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-extracurricular",
		Middlewares: huma.Middlewares{verifyRole(api, models.StudentRole, models.CounselorRole)},
		Method:      http.MethodDelete,
		Path:        "/extracurriculars/{id}",
		Description: "Delete an extracurricular along with its achievements and notes. Students can delete their own; counselors can delete their assigned students'.",
		Tags:        []string{"Extracurriculars"},
	}, func(ctx context.Context, input *models.DeleteExtracurricularInput) (*struct{}, error) {
		if err := extracurricularHandler.DeleteExtracurricular(ctx, input.ID); err != nil {
			return nil, err
		}
		return nil, nil
	})

	studentOrCounselor := huma.Middlewares{verifyRole(api, models.StudentRole, models.CounselorRole)}

	huma.Register(api, huma.Operation{
		OperationID:   "create-extracurricular-note",
		Middlewares:   studentOrCounselor,
		Method:        http.MethodPost,
		Path:          "/extracurriculars/{id}/notes",
		DefaultStatus: http.StatusCreated,
		Description:   "Add a note to an extracurricular. Students can add to their own; counselors to their assigned students'.",
		Tags:          []string{"Extracurriculars"},
	}, func(ctx context.Context, input *models.CreateExtracurricularNoteInput) (*models.ExtracurricularNoteOutput, error) {
		note, err := extracurricularHandler.CreateExtracurricularNote(ctx, input)
		if err != nil {
			return nil, err
		}
		return &models.ExtracurricularNoteOutput{Body: *note}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-extracurricular-note",
		Middlewares: studentOrCounselor,
		Method:      http.MethodPatch,
		Path:        "/extracurriculars/{id}/notes/{noteId}",
		Description: "Edit a note on an extracurricular. Students can edit their own; counselors their assigned students'.",
		Tags:        []string{"Extracurriculars"},
	}, func(ctx context.Context, input *models.UpdateExtracurricularNoteInput) (*models.ExtracurricularNoteOutput, error) {
		note, err := extracurricularHandler.UpdateExtracurricularNote(ctx, input)
		if err != nil {
			return nil, err
		}
		return &models.ExtracurricularNoteOutput{Body: *note}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-extracurricular-note",
		Middlewares: studentOrCounselor,
		Method:      http.MethodDelete,
		Path:        "/extracurriculars/{id}/notes/{noteId}",
		Description: "Delete a note from an extracurricular. Students can delete their own; counselors their assigned students'.",
		Tags:        []string{"Extracurriculars"},
	}, func(ctx context.Context, input *models.DeleteExtracurricularNoteInput) (*struct{}, error) {
		if err := extracurricularHandler.DeleteExtracurricularNote(ctx, input); err != nil {
			return nil, err
		}
		return nil, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "create-extracurricular-achievement",
		Middlewares:   studentOrCounselor,
		Method:        http.MethodPost,
		Path:          "/extracurriculars/{id}/achievements",
		DefaultStatus: http.StatusCreated,
		Description:   "Log an achievement on an extracurricular. Students can add to their own; counselors to their assigned students'.",
		Tags:          []string{"Extracurriculars"},
	}, func(ctx context.Context, input *models.CreateExtracurricularAchievementInput) (*models.ExtracurricularAchievementOutput, error) {
		achievement, err := extracurricularHandler.CreateExtracurricularAchievement(ctx, input)
		if err != nil {
			return nil, err
		}
		return &models.ExtracurricularAchievementOutput{Body: *achievement}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-extracurricular-achievement",
		Middlewares: studentOrCounselor,
		Method:      http.MethodPatch,
		Path:        "/extracurriculars/{id}/achievements/{achievementId}",
		Description: "Edit an achievement on an extracurricular. Omitted fields keep their value.",
		Tags:        []string{"Extracurriculars"},
	}, func(ctx context.Context, input *models.UpdateExtracurricularAchievementInput) (*models.ExtracurricularAchievementOutput, error) {
		achievement, err := extracurricularHandler.UpdateExtracurricularAchievement(ctx, input)
		if err != nil {
			return nil, err
		}
		return &models.ExtracurricularAchievementOutput{Body: *achievement}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete-extracurricular-achievement",
		Middlewares: studentOrCounselor,
		Method:      http.MethodDelete,
		Path:        "/extracurriculars/{id}/achievements/{achievementId}",
		Description: "Delete an achievement from an extracurricular. Students can delete their own; counselors their assigned students'.",
		Tags:        []string{"Extracurriculars"},
	}, func(ctx context.Context, input *models.DeleteExtracurricularAchievementInput) (*struct{}, error) {
		if err := extracurricularHandler.DeleteExtracurricularAchievement(ctx, input); err != nil {
			return nil, err
		}
		return nil, nil
	})
}
