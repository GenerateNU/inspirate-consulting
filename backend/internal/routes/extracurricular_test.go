package routes

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/data"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// setupExtracurricularTestApp creates a Fiber app in test mode with a mocked extracurricular repository.
func setupExtracurricularTestApp(mockRepo data.ExtracurricularRepository) (*fiber.App, error) {
	cfg := config.Config{TestMode: true}
	repo := &data.Repository{Extracurricular: mockRepo}
	app, _, err := SetupApp(cfg, repo, nil, nil)
	return app, err
}

func extracurricularPtr[T any](v T) *T { return &v }

// validCreatePayload returns a create body with every field, including the
// application-detail fields, populated.
func validCreatePayload() map[string]any {
	return map[string]any{
		"name":                    "Varsity Debate",
		"status":                  "doing",
		"type":                    "investment",
		"description":             "Regional tournaments",
		"leadership_role":         "Captain",
		"start_date":              "2024-09-01",
		"end_date":                "2025-06-01",
		"organization":            "Lincoln High",
		"hours_per_week":          7.5,
		"weeks_per_year":          36,
		"goals":                   "Qualify for nationals",
		"ready_for_application":   true,
		"application_description": "Led the team to a state title.",
	}
}

// minimalCreatePayload returns a create body with only the required fields.
func minimalCreatePayload() map[string]any {
	return map[string]any{
		"name":         "Chess",
		"status":       "have_done",
		"type":         "maintenance",
		"start_date":   "2023-01-01",
		"organization": "Chess Club",
	}
}

// Route tests for POST /extracurriculars.
func TestRoute_CreateExtracurricular(t *testing.T) {
	t.Parallel()

	t.Run("success with all fields", func(t *testing.T) {
		t.Parallel()

		studentID := uuid.NewString()
		userID := uuid.NewString()
		headers := map[string]string{"X-Student-ID": studentID, "X-User-ID": userID}

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("CreateExtracurricular", mock.Anything, studentID, userID, mock.MatchedBy(func(in models.CreateExtracurricularInput) bool {
			b := in.Body
			return b.Name == "Varsity Debate" &&
				b.Description != nil && *b.Description == "Regional tournaments" &&
				b.HoursPerWeek != nil && *b.HoursPerWeek == 7.5 &&
				b.WeeksPerYear != nil && *b.WeeksPerYear == 36 &&
				b.Goals != nil && *b.Goals == "Qualify for nationals" &&
				b.ReadyForApplication != nil && *b.ReadyForApplication &&
				b.ApplicationDescription != nil && *b.ApplicationDescription == "Led the team to a state title."
		})).Return(&models.CreateExtracurricularOutput{Body: models.Extracurricular{
			ID:                     1,
			StudentID:              studentID,
			UserID:                 userID,
			Name:                   "Varsity Debate",
			Status:                 models.ExtracurricularStatusDoing,
			Type:                   models.ExtracurricularTypeInvestment,
			Description:            extracurricularPtr("Regional tournaments"),
			StartDate:              time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC),
			Organization:           "Lincoln High",
			HoursPerWeek:           extracurricularPtr(7.5),
			WeeksPerYear:           extracurricularPtr(int16(36)),
			Goals:                  extracurricularPtr("Qualify for nationals"),
			ReadyForApplication:    true,
			ApplicationDescription: extracurricularPtr("Led the team to a state title."),
		}}, nil)

		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, body := doRequest(t, app, http.MethodPost, "/extracurriculars", validCreatePayload(), headers)
		require.Equal(t, http.StatusOK, resp.StatusCode, string(body))

		// Check the JSON field names on the wire, not just the Go struct.
		var raw map[string]any
		require.NoError(t, json.Unmarshal(body, &raw))
		assert.Equal(t, 7.5, raw["hours_per_week"])
		assert.Equal(t, float64(36), raw["weeks_per_year"])
		assert.Equal(t, "Qualify for nationals", raw["goals"])
		assert.Equal(t, true, raw["ready_for_application"])
		assert.Equal(t, "Led the team to a state title.", raw["application_description"])
		assert.Equal(t, "Regional tournaments", raw["description"])
		assert.Equal(t, studentID, raw["student_id"])
	})

	t.Run("success with only required fields", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("CreateExtracurricular", mock.Anything, mock.Anything, mock.Anything, mock.MatchedBy(func(in models.CreateExtracurricularInput) bool {
			b := in.Body
			return b.Description == nil && b.HoursPerWeek == nil && b.WeeksPerYear == nil &&
				b.Goals == nil && b.ReadyForApplication == nil && b.ApplicationDescription == nil
		})).Return(&models.CreateExtracurricularOutput{Body: models.Extracurricular{ID: 2, Name: "Chess"}}, nil)

		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, body := doRequest(t, app, http.MethodPost, "/extracurriculars", minimalCreatePayload())
		require.Equal(t, http.StatusOK, resp.StatusCode, string(body))

		// ready_for_application is always present; nullable fields are omitted when unset.
		var raw map[string]any
		require.NoError(t, json.Unmarshal(body, &raw))
		assert.Equal(t, false, raw["ready_for_application"])
		for _, key := range []string{"description", "hours_per_week", "weeks_per_year", "goals", "application_description"} {
			assert.NotContains(t, raw, key)
		}
	})

	t.Run("boundary values are accepted", func(t *testing.T) {
		t.Parallel()

		for _, tc := range []struct {
			hours float64
			weeks int
		}{{0, 0}, {168, 52}} {
			mockRepo := mocks.NewExtracurricularRepository(t)
			mockRepo.On("CreateExtracurricular", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(&models.CreateExtracurricularOutput{}, nil)

			app, err := setupExtracurricularTestApp(mockRepo)
			require.NoError(t, err)

			payload := minimalCreatePayload()
			payload["hours_per_week"] = tc.hours
			payload["weeks_per_year"] = tc.weeks

			resp, body := doRequest(t, app, http.MethodPost, "/extracurriculars", payload)
			assert.Equal(t, http.StatusOK, resp.StatusCode, "hours=%v weeks=%v: %s", tc.hours, tc.weeks, body)
		}
	})

	// Each case must be rejected by Huma before the repository is called.
	validationCases := []struct {
		name   string
		mutate func(p map[string]any)
	}{
		{"missing name", func(p map[string]any) { delete(p, "name") }},
		{"empty name", func(p map[string]any) { p["name"] = "" }},
		{"missing status", func(p map[string]any) { delete(p, "status") }},
		{"invalid status", func(p map[string]any) { p["status"] = "planning" }},
		{"missing type", func(p map[string]any) { delete(p, "type") }},
		{"invalid type", func(p map[string]any) { p["type"] = "hobby" }},
		{"missing start_date", func(p map[string]any) { delete(p, "start_date") }},
		{"malformed start_date", func(p map[string]any) { p["start_date"] = "09/01/2024" }},
		{"malformed end_date", func(p map[string]any) { p["end_date"] = "soon" }},
		{"missing organization", func(p map[string]any) { delete(p, "organization") }},
		{"negative hours_per_week", func(p map[string]any) { p["hours_per_week"] = -1 }},
		{"hours_per_week wrong type", func(p map[string]any) { p["hours_per_week"] = "five" }},
		{"negative weeks_per_year", func(p map[string]any) { p["weeks_per_year"] = -1 }},
		{"weeks_per_year above 52", func(p map[string]any) { p["weeks_per_year"] = 53 }},
		{"weeks_per_year not an integer", func(p map[string]any) { p["weeks_per_year"] = 10.5 }},
		{"ready_for_application wrong type", func(p map[string]any) { p["ready_for_application"] = "yes" }},
		{"goals wrong type", func(p map[string]any) { p["goals"] = 123 }},
		{"application_description wrong type", func(p map[string]any) { p["application_description"] = []string{"a"} }},
	}
	for _, tc := range validationCases {
		t.Run("validation error - "+tc.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := mocks.NewExtracurricularRepository(t)
			app, err := setupExtracurricularTestApp(mockRepo)
			require.NoError(t, err)

			payload := validCreatePayload()
			tc.mutate(payload)

			resp, body := doRequest(t, app, http.MethodPost, "/extracurriculars", payload)
			assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode, string(body))
		})
	}

	t.Run("counselor without student identity is unauthorized", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodPost, "/extracurriculars", validCreatePayload(),
			map[string]string{"X-Counselor-ID": uuid.NewString()})
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("CreateExtracurricular", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("database error"))

		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodPost, "/extracurriculars", validCreatePayload())
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// Route tests for GET /extracurriculars.
func TestRoute_ListExtracurriculars(t *testing.T) {
	t.Parallel()

	t.Run("success returns new fields", func(t *testing.T) {
		t.Parallel()

		studentID := uuid.NewString()
		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("ListExtracurriculars", mock.Anything, studentID).Return([]models.Extracurricular{
			{
				ID: 1, StudentID: studentID, Name: "Debate",
				HoursPerWeek: extracurricularPtr(5.0), WeeksPerYear: extracurricularPtr(int16(40)),
				Goals: extracurricularPtr("Nationals"), ReadyForApplication: true,
				ApplicationDescription: extracurricularPtr("Blurb"),
			},
			{ID: 2, StudentID: studentID, Name: "Chess"},
		}, nil)

		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, body := doRequest(t, app, http.MethodGet, "/extracurriculars", nil, map[string]string{"X-Student-ID": studentID})
		require.Equal(t, http.StatusOK, resp.StatusCode, string(body))

		var list []models.Extracurricular
		require.NoError(t, json.Unmarshal(body, &list))
		require.Len(t, list, 2)
		assert.Equal(t, 5.0, *list[0].HoursPerWeek)
		assert.Equal(t, int16(40), *list[0].WeeksPerYear)
		assert.Equal(t, "Nationals", *list[0].Goals)
		assert.True(t, list[0].ReadyForApplication)
		assert.Equal(t, "Blurb", *list[0].ApplicationDescription)
		assert.Nil(t, list[1].HoursPerWeek)
		assert.False(t, list[1].ReadyForApplication)
	})

	t.Run("empty list is a JSON array", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("ListExtracurriculars", mock.Anything, mock.Anything).Return([]models.Extracurricular{}, nil)

		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, body := doRequest(t, app, http.MethodGet, "/extracurriculars", nil)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		assert.JSONEq(t, `[]`, string(body))
	})

	t.Run("counselor without student identity is unauthorized", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodGet, "/extracurriculars", nil, map[string]string{"X-Counselor-ID": uuid.NewString()})
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("ListExtracurriculars", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))

		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodGet, "/extracurriculars", nil)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// Route tests for PATCH /extracurriculars/{id}.
func TestRoute_UpdateExtracurricular(t *testing.T) {
	t.Parallel()

	t.Run("success updating new fields", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("UpdateExtracurricular", mock.Anything, int64(7), mock.MatchedBy(func(in models.UpdateExtracurricularInput) bool {
			b := in.Body
			return in.ID == 7 && b.Name == nil &&
				b.HoursPerWeek != nil && *b.HoursPerWeek == 9.5 &&
				b.WeeksPerYear != nil && *b.WeeksPerYear == 50 &&
				b.Goals != nil && *b.Goals == "New goals" &&
				b.ReadyForApplication != nil && *b.ReadyForApplication &&
				b.ApplicationDescription != nil && *b.ApplicationDescription == "New blurb"
		})).Return(&models.Extracurricular{
			ID: 7, Name: "Debate",
			HoursPerWeek: extracurricularPtr(9.5), WeeksPerYear: extracurricularPtr(int16(50)),
			Goals: extracurricularPtr("New goals"), ReadyForApplication: true,
			ApplicationDescription: extracurricularPtr("New blurb"),
		}, nil)

		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, body := doRequest(t, app, http.MethodPatch, "/extracurriculars/7", map[string]any{
			"hours_per_week":          9.5,
			"weeks_per_year":          50,
			"goals":                   "New goals",
			"ready_for_application":   true,
			"application_description": "New blurb",
		})
		require.Equal(t, http.StatusOK, resp.StatusCode, string(body))

		var updated models.Extracurricular
		require.NoError(t, json.Unmarshal(body, &updated))
		assert.Equal(t, int64(7), updated.ID)
		assert.Equal(t, 9.5, *updated.HoursPerWeek)
		assert.True(t, updated.ReadyForApplication)
	})

	t.Run("explicit false for ready_for_application reaches the repository", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("UpdateExtracurricular", mock.Anything, int64(7), mock.MatchedBy(func(in models.UpdateExtracurricularInput) bool {
			return in.Body.ReadyForApplication != nil && !*in.Body.ReadyForApplication
		})).Return(&models.Extracurricular{ID: 7}, nil)

		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, body := doRequest(t, app, http.MethodPatch, "/extracurriculars/7", map[string]any{"ready_for_application": false})
		assert.Equal(t, http.StatusOK, resp.StatusCode, string(body))
	})

	t.Run("empty body is accepted", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("UpdateExtracurricular", mock.Anything, int64(7), mock.MatchedBy(func(in models.UpdateExtracurricularInput) bool {
			return in.Body == models.UpdateExtracurricularRequest{}
		})).Return(&models.Extracurricular{ID: 7}, nil)

		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, body := doRequest(t, app, http.MethodPatch, "/extracurriculars/7", map[string]any{})
		assert.Equal(t, http.StatusOK, resp.StatusCode, string(body))
	})

	validationCases := []struct {
		name string
		body map[string]any
	}{
		{"empty name", map[string]any{"name": ""}},
		{"invalid status", map[string]any{"status": "planning"}},
		{"invalid type", map[string]any{"type": "hobby"}},
		{"malformed start_date", map[string]any{"start_date": "2024/01/01"}},
		{"malformed end_date", map[string]any{"end_date": "later"}},
		{"negative hours_per_week", map[string]any{"hours_per_week": -0.5}},
		{"hours_per_week wrong type", map[string]any{"hours_per_week": "ten"}},
		{"negative weeks_per_year", map[string]any{"weeks_per_year": -1}},
		{"weeks_per_year above 52", map[string]any{"weeks_per_year": 53}},
		{"weeks_per_year not an integer", map[string]any{"weeks_per_year": 1.5}},
		{"ready_for_application wrong type", map[string]any{"ready_for_application": "true"}},
		{"goals wrong type", map[string]any{"goals": false}},
		{"application_description wrong type", map[string]any{"application_description": 99}},
	}
	for _, tc := range validationCases {
		t.Run("validation error - "+tc.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := mocks.NewExtracurricularRepository(t)
			app, err := setupExtracurricularTestApp(mockRepo)
			require.NoError(t, err)

			resp, body := doRequest(t, app, http.MethodPatch, "/extracurriculars/7", tc.body)
			assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode, string(body))
		})
	}

	t.Run("non-numeric id is rejected", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodPatch, "/extracurriculars/abc", map[string]any{"goals": "x"})
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("counselor without student identity is unauthorized", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodPatch, "/extracurriculars/7", map[string]any{"goals": "x"},
			map[string]string{"X-Counselor-ID": uuid.NewString()})
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("UpdateExtracurricular", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("database error"))

		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodPatch, "/extracurriculars/7", map[string]any{"goals": "x"})
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// Route tests for GET /extracurriculars/{id}.
func TestRoute_GetExtracurricular(t *testing.T) {
	t.Parallel()

	t.Run("success returns flat extracurricular fields plus achievements and notes", func(t *testing.T) {
		t.Parallel()

		studentID := uuid.NewString()
		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricular", mock.Anything, int64(5)).Return(&models.ExtracurricularDetail{
			Extracurricular: models.Extracurricular{
				ID: 5, StudentID: studentID, Name: "Debate",
				StartDate:    time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC),
				HoursPerWeek: extracurricularPtr(6.5), WeeksPerYear: extracurricularPtr(int16(40)),
				Goals: extracurricularPtr("State"), ReadyForApplication: true,
				ApplicationDescription: extracurricularPtr("Blurb"),
			},
			Achievements: []models.ExtracurricularAchievement{
				{ID: 100, ExtracurricularID: 5, EntryDate: time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC), Body: "Won state"},
				{ID: 99, ExtracurricularID: 5, EntryDate: time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC), Body: "Won regionals"},
			},
			Notes: []models.ExtracurricularNote{
				{ID: 200, ExtracurricularID: 5, Body: "Ask for rec letter"},
			},
		}, nil)

		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, body := doRequest(t, app, http.MethodGet, "/extracurriculars/5", nil, map[string]string{"X-Student-ID": studentID})
		require.Equal(t, http.StatusOK, resp.StatusCode, string(body))

		var raw map[string]any
		require.NoError(t, json.Unmarshal(body, &raw))
		// Embedded fields are flattened to the top level, not nested.
		assert.NotContains(t, raw, "Extracurricular")
		assert.Equal(t, float64(5), raw["id"])
		assert.Equal(t, "Debate", raw["name"])
		assert.Equal(t, 6.5, raw["hours_per_week"])
		assert.Equal(t, float64(40), raw["weeks_per_year"])
		assert.Equal(t, "State", raw["goals"])
		assert.Equal(t, true, raw["ready_for_application"])
		assert.Equal(t, "Blurb", raw["application_description"])

		achievements, ok := raw["achievements"].([]any)
		require.True(t, ok, "achievements should be an array")
		require.Len(t, achievements, 2)
		first := achievements[0].(map[string]any)
		assert.Equal(t, float64(100), first["id"])
		assert.Equal(t, float64(5), first["extracurricular_id"])
		assert.Equal(t, "Won state", first["body"])
		assert.Contains(t, first["entry_date"], "2025-03-01")

		notes, ok := raw["notes"].([]any)
		require.True(t, ok, "notes should be an array")
		require.Len(t, notes, 1)
		assert.Equal(t, "Ask for rec letter", notes[0].(map[string]any)["body"])
	})

	t.Run("no children returns empty arrays, not null", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricular", mock.Anything, int64(6)).Return(&models.ExtracurricularDetail{
			Extracurricular: models.Extracurricular{ID: 6, StudentID: uuid.NewString()},
			Achievements:    []models.ExtracurricularAchievement{},
			Notes:           []models.ExtracurricularNote{},
		}, nil)

		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, body := doRequest(t, app, http.MethodGet, "/extracurriculars/6", nil)
		require.Equal(t, http.StatusOK, resp.StatusCode, string(body))

		var raw map[string]json.RawMessage
		require.NoError(t, json.Unmarshal(body, &raw))
		assert.JSONEq(t, `[]`, string(raw["achievements"]))
		assert.JSONEq(t, `[]`, string(raw["notes"]))
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricular", mock.Anything, int64(404)).Return(nil, errs.NotFound("extracurricular", "id", int64(404)))

		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodGet, "/extracurriculars/404", nil)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("non-numeric id is rejected", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodGet, "/extracurriculars/abc", nil)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("forbidden from the authorizer is surfaced", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricular", mock.Anything, int64(7)).Return(nil, errs.Forbidden("you can only access your own data"))

		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodGet, "/extracurriculars/7", nil)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("repository error", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricular", mock.Anything, int64(8)).Return(nil, errors.New("database error"))

		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodGet, "/extracurriculars/8", nil)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// Route tests for DELETE /extracurriculars/{id}.
func TestRoute_DeleteExtracurricular(t *testing.T) {
	t.Parallel()

	t.Run("success returns 204 with no body", func(t *testing.T) {
		t.Parallel()

		ownerID := uuid.NewString()
		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricular", mock.Anything, int64(5)).Return(&models.ExtracurricularDetail{
			Extracurricular: models.Extracurricular{ID: 5, StudentID: ownerID},
		}, nil)
		mockRepo.On("DeleteExtracurricular", mock.Anything, int64(5), ownerID).Return(nil)

		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, body := doRequest(t, app, http.MethodDelete, "/extracurriculars/5", nil, map[string]string{"X-Student-ID": ownerID})
		assert.Equal(t, http.StatusNoContent, resp.StatusCode, string(body))
		assert.Empty(t, body)
	})

	t.Run("not found skips the delete", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricular", mock.Anything, int64(404)).Return(nil, errs.NotFound("extracurricular", "id", int64(404)))

		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodDelete, "/extracurriculars/404", nil)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockRepo.AssertNotCalled(t, "DeleteExtracurricular", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("not found from the delete itself", func(t *testing.T) {
		t.Parallel()

		ownerID := uuid.NewString()
		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricular", mock.Anything, int64(9)).Return(&models.ExtracurricularDetail{
			Extracurricular: models.Extracurricular{ID: 9, StudentID: ownerID},
		}, nil)
		mockRepo.On("DeleteExtracurricular", mock.Anything, int64(9), ownerID).Return(errs.NotFound("extracurricular", "id", int64(9)))

		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodDelete, "/extracurriculars/9", nil)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("non-numeric id is rejected", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewExtracurricularRepository(t)
		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodDelete, "/extracurriculars/abc", nil)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("repository error on delete", func(t *testing.T) {
		t.Parallel()

		ownerID := uuid.NewString()
		mockRepo := mocks.NewExtracurricularRepository(t)
		mockRepo.On("GetExtracurricular", mock.Anything, int64(10)).Return(&models.ExtracurricularDetail{
			Extracurricular: models.Extracurricular{ID: 10, StudentID: ownerID},
		}, nil)
		mockRepo.On("DeleteExtracurricular", mock.Anything, int64(10), ownerID).Return(errors.New("database error"))

		app, err := setupExtracurricularTestApp(mockRepo)
		require.NoError(t, err)

		resp, _ := doRequest(t, app, http.MethodDelete, "/extracurriculars/10", nil)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
