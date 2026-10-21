package models

import (
	"time"

	"github.com/google/uuid"
)

// To represent an open essay review a counselor still needs to complete.
type PendingReview struct {
	TransactionID uuid.UUID `json:"transaction_id" doc:"ID of the review transaction, used to complete or refund the review"`
	EssayID       uuid.UUID `json:"essay_id" doc:"ID of the essay to review"`
	EssayType     string    `json:"essay_type" example:"personal-statement" doc:"Type of the essay to review"`
	EssayLink     string    `json:"essay_link" example:"https://docs.google.com/document/d/abc123" doc:"Link to the essay's content, so the counselor can open it to review"`
	StudentID     uuid.UUID `json:"student_id" doc:"ID of the student who requested the review"`
	StudentName   string    `json:"student_name" example:"Jane Doe" doc:"Name of the student who requested the review"`
	RequestedAt   time.Time `json:"requested_at" doc:"When the review was requested"`
}

// To represent one of a counselor's students as stored in the database.
type CounselorStudentRow struct {
	StudentID uuid.UUID `json:"student_id" doc:"The student's unique ID"`
	Name      string    `json:"name" example:"Jane Doe" doc:"Name of the student"`
	Year      Year      `json:"year" example:"junior" doc:"School year of the student"`
}

// To represent one of a counselor's students, flagged for the counselor's student list
// separated from CounselorStudentRow to avoid adding extra fields to the database model that are only used for the list operation
type CounselorStudent struct {
	CounselorStudentRow
	IsFallingBehind  bool `json:"is_falling_behind" doc:"True when the student has an incomplete to-do item past its deadline"`
	HasPendingReview bool `json:"has_pending_review" doc:"True when the student has an essay review that is still open"`
}

// Huma readable input and output models for the list operations
// counselor id is derived from session, so no input is needed for list pending reviews operation
type ListPendingReviewsInput struct{}

type ListPendingReviewsOutput struct {
	Body []PendingReview
}

// counselor id is derived from session, filters and sorting are optional
type ListCounselorStudentsInput struct {
	Year             Year   `query:"year" required:"false" enum:"freshman,sophomore,junior,senior" doc:"Only return students in this school year"`
	HasPendingReview bool   `query:"has_pending_review" required:"false" doc:"Only return students with an open essay review"`
	SortBy           string `query:"sort_by" enum:"name,year" default:"name" doc:"Field to sort by"`
	SortOrder        string `query:"sort_order" enum:"asc,desc" default:"asc" doc:"Sort direction"`
}

type ListCounselorStudentsOutput struct {
	Body []CounselorStudent
}
