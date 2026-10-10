package models

import (
	"time"

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
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	SupabaseID uuid.UUID  `json:"supabase_id"`
	PfpKey     *string    `json:"pfp_key"`
	ResetTime  *time.Time `json:"reset_time"`
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

type CreateUserInput struct {
	Body struct {
		Name   string  `json:"name" db:"name" doc:"Name of the user" minLength:"1" maxLength:"200"`
		PfpKey *string `json:"pfp_key" db:"pfp_key" doc:"pfp key of user"`
		Email  string  `json:"email" doc:"email of user for supabase signup"`
	}
}

type CreateUserBody struct {
	User         *User  `json:"user"`
	TempPassword string `json:"temp_password"`
}

type CreateUserOutput struct {
	Body *CreateUserBody
}

type FetchUserInput struct {
	ID uuid.UUID `path:"id" required:"true"`
}

type FetchUserBySupabaseIDInput struct {
	SupabaseID uuid.UUID `path:"supabase_id" required:"true"`
}

type FetchUserOutput struct {
	Body *User `json:"body"`
}
