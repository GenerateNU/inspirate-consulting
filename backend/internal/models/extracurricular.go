package models

import (
	"time"
)

// ParseDate converts an ISO date string to a time.Time value.
func ParseDate(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	return time.Parse("2006-01-02", value)
}

// ExtracurricularStatus tracks whether the activity is currently being done
// or has already been completed.
type ExtracurricularStatus string

const (
	ExtracurricularStatusDoing    ExtracurricularStatus = "doing"
	ExtracurricularStatusHaveDone ExtracurricularStatus = "have_done"
)

// ExtracurricularType separates maintenance work from investment activities.
type ExtracurricularType string

const (
	ExtracurricularTypeMaintenance ExtracurricularType = "maintenance"
	ExtracurricularTypeInvestment  ExtracurricularType = "investment"
)

// Extracurricular is the data model for a student activity entry.
// It includes the required student/user ownership, timing, and organization fields.
// The db tags let pgx.RowToStructByName map query rows directly onto this struct.
type Extracurricular struct {
	ID             int64                 `json:"id"                        db:"id"`
	CreatedAt      time.Time             `json:"created_at"                db:"created_at"`
	ModifiedAt     time.Time             `json:"modified_at"                db:"modified_at"`
	StudentID      string                `json:"student_id"                db:"student_id"`
	UserID         string                `json:"user_id"                   db:"user_id"`
	Name           string                `json:"name"                      db:"name"`
	Status         ExtracurricularStatus `json:"status"                    db:"status"`
	Type           ExtracurricularType   `json:"type"                      db:"type"`
	Description    *string               `json:"description,omitempty"     db:"description"`
	LeadershipRole *string               `json:"leadership_role,omitempty" db:"leadership_role"`
	StartDate      time.Time             `json:"start_date"                db:"start_date"`
	EndDate        *time.Time            `json:"end_date,omitempty"        db:"end_date"`
	Organization   string                `json:"organization"              db:"organization"`

	HoursPerWeek           *float64 `json:"hours_per_week,omitempty"          db:"hours_per_week"`
	WeeksPerYear           *int16   `json:"weeks_per_year,omitempty"          db:"weeks_per_year"`
	Goals                  *string  `json:"goals,omitempty"                   db:"goals"`
	ReadyForApplication    bool     `json:"ready_for_application"             db:"ready_for_application"`
	ApplicationDescription *string  `json:"application_description,omitempty" db:"application_description"`
}

// CreateExtracurricularRequest is the body shape for creating an extracurricular.
// StudentID comes from the auth context and UserID comes from the auth context,
// so neither is accepted in the body. Huma validates required fields, enums,
// and date formats from these tags before the handler runs.
type CreateExtracurricularRequest struct {
	Name           string                `json:"name"                      minLength:"1"`
	Status         ExtracurricularStatus `json:"status"                    enum:"doing,have_done"`
	Type           ExtracurricularType   `json:"type"                      enum:"maintenance,investment"`
	Description    *string               `json:"description,omitempty"`
	LeadershipRole *string               `json:"leadership_role,omitempty"`
	StartDate      string                `json:"start_date"                format:"date"`
	EndDate        *string               `json:"end_date,omitempty"        format:"date"`
	Organization   string                `json:"organization"`

	HoursPerWeek           *float64 `json:"hours_per_week,omitempty"          minimum:"0"`
	WeeksPerYear           *int16   `json:"weeks_per_year,omitempty"          minimum:"0" maximum:"52"`
	Goals                  *string  `json:"goals,omitempty"`
	ReadyForApplication    *bool    `json:"ready_for_application,omitempty"   doc:"Defaults to false when omitted"`
	ApplicationDescription *string  `json:"application_description,omitempty"`
}

// CreateExtracurricularInput is the request payload shape for creating a new extracurricular record.
type CreateExtracurricularInput struct {
	Body CreateExtracurricularRequest
}

// CreateExtracurricularOutput is the created record returned after insert.
type CreateExtracurricularOutput struct {
	Body Extracurricular
}

// ListExtracurricularsInput holds the request context for listing a student's extracurriculars.
// The student is taken from the auth context, so there are no request parameters.
type ListExtracurricularsInput struct{}

// ListExtracurricularsOutput contains the extracurricular records returned for a student.
type ListExtracurricularsOutput struct {
	Body []Extracurricular
}

// UpdateExtracurricularRequest is the partial update body shape for extracurriculars.
// Every field is optional; omitted fields keep their current value. Fields
// backed by nullable columns use Nullable so that sending an explicit null
// clears the stored value. Required columns stay plain pointers, where null is
// treated the same as omitted. StudentID is derived from the existing entry and
// UserID from the auth context.
type UpdateExtracurricularRequest struct {
	Name         *string                `json:"name,omitempty"         minLength:"1"`
	Status       *ExtracurricularStatus `json:"status,omitempty"       enum:"doing,have_done"`
	Type         *ExtracurricularType   `json:"type,omitempty"         enum:"maintenance,investment"`
	StartDate    *string                `json:"start_date,omitempty"   format:"date"`
	Organization *string                `json:"organization,omitempty"`

	Description    Nullable[string] `json:"description,omitempty"     doc:"Send null to clear"`
	LeadershipRole Nullable[string] `json:"leadership_role,omitempty" doc:"Send null to clear"`
	EndDate        Nullable[string] `json:"end_date,omitempty"        format:"date" doc:"Send null to clear"`

	HoursPerWeek           Nullable[float64] `json:"hours_per_week,omitempty"          minimum:"0" doc:"Send null to clear"`
	WeeksPerYear           Nullable[int16]   `json:"weeks_per_year,omitempty"          minimum:"0" maximum:"52" doc:"Send null to clear"`
	Goals                  Nullable[string]  `json:"goals,omitempty"                   doc:"Send null to clear"`
	ReadyForApplication    *bool             `json:"ready_for_application,omitempty"`
	ApplicationDescription Nullable[string]  `json:"application_description,omitempty" doc:"Send null to clear"`
}

// UpdateExtracurricularInput updates a student's extracurricular entry.
type UpdateExtracurricularInput struct {
	ID   int64 `path:"id" doc:"ID of the extracurricular to update"`
	Body UpdateExtracurricularRequest
}

// UpdateExtracurricularOutput is the updated record returned after a save.
type UpdateExtracurricularOutput struct {
	Body Extracurricular
}

// ExtracurricularAchievement is a dated accomplishment logged against an extracurricular.
type ExtracurricularAchievement struct {
	ID                int64     `json:"id"                 db:"id"`
	ExtracurricularID int64     `json:"extracurricular_id" db:"extracurricular_id"`
	EntryDate         time.Time `json:"entry_date"         db:"entry_date"`
	Body              string    `json:"body"               db:"body"`
	CreatedAt         time.Time `json:"created_at"         db:"created_at"`
	ModifiedAt        time.Time `json:"modified_at"        db:"modified_at"`
}

// ExtracurricularNote is a free-form note attached to an extracurricular.
type ExtracurricularNote struct {
	ID                int64     `json:"id"                 db:"id"`
	ExtracurricularID int64     `json:"extracurricular_id" db:"extracurricular_id"`
	Body              string    `json:"body"               db:"body"`
	CreatedAt         time.Time `json:"created_at"         db:"created_at"`
	ModifiedAt        time.Time `json:"modified_at"        db:"modified_at"`
}

// ExtracurricularDetail is a single extracurricular with all of its
// achievements and notes. Both lists are always present (empty, not null).
type ExtracurricularDetail struct {
	Extracurricular
	Achievements []ExtracurricularAchievement `json:"achievements"`
	Notes        []ExtracurricularNote        `json:"notes"`
}

// GetExtracurricularInput fetches one extracurricular by id.
type GetExtracurricularInput struct {
	ID int64 `path:"id" doc:"ID of the extracurricular to fetch"`
}

// GetExtracurricularOutput is the extracurricular with its achievements and notes.
type GetExtracurricularOutput struct {
	Body ExtracurricularDetail
}

// DeleteExtracurricularInput deletes one extracurricular by id.
type DeleteExtracurricularInput struct {
	ID int64 `path:"id" doc:"ID of the extracurricular to delete"`
}

// ExtracurricularNoteRequest is the body for creating or editing a note.
type ExtracurricularNoteRequest struct {
	Body string `json:"body" minLength:"1"`
}

// CreateExtracurricularNoteInput adds a note to an extracurricular.
type CreateExtracurricularNoteInput struct {
	ExtracurricularID int64 `path:"id" doc:"ID of the extracurricular"`
	Body              ExtracurricularNoteRequest
}

// UpdateExtracurricularNoteInput replaces the text of an existing note.
type UpdateExtracurricularNoteInput struct {
	ExtracurricularID int64 `path:"id"     doc:"ID of the extracurricular"`
	NoteID            int64 `path:"noteId" doc:"ID of the note"`
	Body              ExtracurricularNoteRequest
}

// DeleteExtracurricularNoteInput removes a note.
type DeleteExtracurricularNoteInput struct {
	ExtracurricularID int64 `path:"id"     doc:"ID of the extracurricular"`
	NoteID            int64 `path:"noteId" doc:"ID of the note"`
}

// ExtracurricularNoteOutput returns a single note.
type ExtracurricularNoteOutput struct {
	Body ExtracurricularNote
}

// CreateExtracurricularAchievementRequest is the body for logging an achievement.
type CreateExtracurricularAchievementRequest struct {
	EntryDate string `json:"entry_date" format:"date"`
	Body      string `json:"body"       minLength:"1"`
}

// UpdateExtracurricularAchievementRequest is a partial update for an
// achievement. Both columns are required, so omitted fields keep their value.
type UpdateExtracurricularAchievementRequest struct {
	EntryDate *string `json:"entry_date,omitempty" format:"date"`
	Body      *string `json:"body,omitempty"       minLength:"1"`
}

// CreateExtracurricularAchievementInput adds an achievement to an extracurricular.
type CreateExtracurricularAchievementInput struct {
	ExtracurricularID int64 `path:"id" doc:"ID of the extracurricular"`
	Body              CreateExtracurricularAchievementRequest
}

// UpdateExtracurricularAchievementInput edits an existing achievement.
type UpdateExtracurricularAchievementInput struct {
	ExtracurricularID int64 `path:"id"            doc:"ID of the extracurricular"`
	AchievementID     int64 `path:"achievementId" doc:"ID of the achievement"`
	Body              UpdateExtracurricularAchievementRequest
}

// DeleteExtracurricularAchievementInput removes an achievement.
type DeleteExtracurricularAchievementInput struct {
	ExtracurricularID int64 `path:"id"            doc:"ID of the extracurricular"`
	AchievementID     int64 `path:"achievementId" doc:"ID of the achievement"`
}

// ExtracurricularAchievementOutput returns a single achievement.
type ExtracurricularAchievementOutput struct {
	Body ExtracurricularAchievement
}
