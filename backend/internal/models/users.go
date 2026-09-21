package models

import (
	"github.com/google/uuid"
)

// enum type for year of student
type Year string

const (
	Freshman  Year = "freshman"
	Sophomore Year = "sophomore"
	Junior    Year = "junior"
	Senior    Year = "senior"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	SupabaseID uuid.UUID `json:"supabase_id"`
	PfpKey     *string   `json:"pfp_key"`
}

// I represented year under the assumption that we're following the american high school timeline
// not sure if we want to change this to assume otherwise
type Student struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"user_id"`
	Year          Year       `json:"year"`
	Organization  *string    `json:"organization"`
	GPA           int        `json:"gpa"`
	ReviewBalance int        `json:"review_balance"`
	CounselorID   *uuid.UUID `json:"counselor_id"`
}

type Counselor struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}
