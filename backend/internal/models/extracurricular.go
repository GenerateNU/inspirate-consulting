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
	Description    string                `json:"description"               db:"description"`
	LeadershipRole *string               `json:"leadership_role,omitempty" db:"leadership_role"`
	StartDate      time.Time             `json:"start_date"                db:"start_date"`
	EndDate        *time.Time            `json:"end_date,omitempty"        db:"end_date"`
	Organization   string                `json:"organization"              db:"organization"`
}

// CreateExtracurricularRequest is the body shape for creating an extracurricular.
// StudentID comes from the route path and UserID comes from the auth context,
// so neither is accepted in the body. Huma validates required fields, enums,
// and date formats from these tags before the handler runs.
type CreateExtracurricularRequest struct {
	Name           string                `json:"name"                      minLength:"1"`
	Status         ExtracurricularStatus `json:"status"                    enum:"doing,have_done"`
	Type           ExtracurricularType   `json:"type"                      enum:"maintenance,investment"`
	Description    string                `json:"description"`
	LeadershipRole *string               `json:"leadership_role,omitempty"`
	StartDate      string                `json:"start_date"                format:"date"`
	EndDate        *string               `json:"end_date,omitempty"        format:"date"`
	Organization   string                `json:"organization"`
}

// CreateExtracurricularInput is the request payload shape for creating a new extracurricular record.
type CreateExtracurricularInput struct {
	StudentID string `path:"studentID" doc:"ID of the student the extracurricular belongs to"`
	Body      CreateExtracurricularRequest
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
// Every field is optional; omitted fields keep their current value. StudentID is
// derived from the existing entry and UserID from the auth context.
type UpdateExtracurricularRequest struct {
	Name           *string                `json:"name,omitempty"            minLength:"1"`
	Status         *ExtracurricularStatus `json:"status,omitempty"          enum:"doing,have_done"`
	Type           *ExtracurricularType   `json:"type,omitempty"            enum:"maintenance,investment"`
	Description    *string                `json:"description,omitempty"`
	LeadershipRole *string                `json:"leadership_role,omitempty"`
	StartDate      *string                `json:"start_date,omitempty"      format:"date"`
	EndDate        *string                `json:"end_date,omitempty"        format:"date"`
	Organization   *string                `json:"organization,omitempty"`
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
