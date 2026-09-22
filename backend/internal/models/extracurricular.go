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
type Extracurricular struct {
	ID             int64                 `json:"id"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
	StudentID      string                `json:"student_id"`
	UserID         string                `json:"user_id"`
	Name           string                `json:"name"`
	Status         ExtracurricularStatus `json:"status"`
	Type           ExtracurricularType   `json:"type"`
	Description    string                `json:"description"`
	LeadershipRole *string               `json:"leadership_role,omitempty"`
	StartDate      time.Time             `json:"start_date"`
	EndDate        *time.Time            `json:"end_date,omitempty"`
	Organization   string                `json:"organization"`
}

// CreateExtracurricularRequest is the body shape for creating an extracurricular.
type CreateExtracurricularRequest struct {
	StudentID      *string               `json:"student_id,omitempty"`
	UserID         string                `json:"user_id"`
	Name           string                `json:"name"`
	Status         ExtracurricularStatus `json:"status"`
	Type           ExtracurricularType   `json:"type"`
	Description    string                `json:"description"`
	LeadershipRole *string               `json:"leadership_role,omitempty"`
	StartDate      string                `json:"start_date"`
	EndDate        *string               `json:"end_date,omitempty"`
	Organization   string                `json:"organization"`
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
type ListExtracurricularsInput struct {
	StudentID string
}

// ListExtracurricularsOutput contains the extracurricular records returned for a student.
type ListExtracurricularsOutput struct {
	Body []Extracurricular
}

// UpdateExtracurricularRequest is the partial update body shape for extracurriculars.
type UpdateExtracurricularRequest struct {
	StudentID      *string               `json:"student_id,omitempty"`
	UserID         *string               `json:"user_id,omitempty"`
	Name           *string               `json:"name,omitempty"`
	Status         *ExtracurricularStatus `json:"status,omitempty"`
	Type           *ExtracurricularType  `json:"type,omitempty"`
	Description    *string               `json:"description,omitempty"`
	LeadershipRole *string               `json:"leadership_role,omitempty"`
	StartDate      *string               `json:"start_date,omitempty"`
	EndDate        *string               `json:"end_date,omitempty"`
	Organization   *string               `json:"organization,omitempty"`
}

// UpdateExtracurricularInput updates a student's extracurricular entry.
type UpdateExtracurricularInput struct {
	ID   int64
	Body UpdateExtracurricularRequest
}

// UpdateExtracurricularOutput is the updated record returned after a save.
type UpdateExtracurricularOutput struct {
	Body Extracurricular
}
