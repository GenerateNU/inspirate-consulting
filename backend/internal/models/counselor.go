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
	StudentID            uuid.UUID  `json:"student_id" doc:"The student's unique ID"`
	Name                 string     `json:"name" example:"Jane Doe" doc:"Name of the student"`
	GraduationYear       *int       `json:"graduation_year" nullable:"true" example:"2027" doc:"Year the student graduates high school; null if not set yet"`
	PendingReviewCount   int        `json:"pending_review_count" doc:"Number of the student's essay reviews that are still open"`
	OpenTaskCount        int        `json:"open_task_count" doc:"Number of the student's incomplete to-do items"`
	UnreadMessageCount   int        `json:"unread_message_count" doc:"Number of messages from the student that the counselor hasn't read"`
	LatestActivityAction *Action    `json:"latest_activity_action" nullable:"true" enum:"added_school,removed_school,essay_status_changed,added_extracurricular,removed_extracurricular,task_past_due" doc:"The student's most recent activity; null if they have none"`
	LatestActivityAt     *time.Time `json:"latest_activity_at" nullable:"true" doc:"When the student's most recent activity happened; null if they have none"`
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
	GraduationYear   int    `query:"graduation_year" required:"false" doc:"Only return students graduating in this year"`
	HasPendingReview bool   `query:"has_pending_review" required:"false" doc:"Only return students with an open essay review"`
	SortBy           string `query:"sort_by" enum:"name,graduation_year" default:"name" doc:"Field to sort by"`
	SortOrder        string `query:"sort_order" enum:"asc,desc" default:"asc" doc:"Sort direction"`
}

type ListCounselorStudentsOutput struct {
	Body []CounselorStudent
}

type UpdateCounselorBody struct {
	CustomLink Nullable[string]    `json:"custom_link,omitempty" format:"uri" doc:"Counselor's custom link; send null to clear"`
	ExpiresAt  Nullable[time.Time] `json:"expires_at,omitempty" doc:"When the counselor's account expires; send null to remove expiration"`
}

// counselor id is derived from session, counselors can only edit themselves
type UpdateCounselorInput struct {
	Body UpdateCounselorBody
}

type UpdateCounselorOutput struct {
	Body *Counselor
}

type GetCounselorInput struct {
	ID uuid.UUID `path:"id" doc:"ID of the counselor"`
}

type GetCounselorOutput struct {
	Body *Counselor
}
