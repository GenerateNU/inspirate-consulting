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

// type of role a counselor can create when making a user
type Role string

const (
	CounselorRole Role = "counselor"
	StudentRole   Role = "student"
)

type User struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	SupabaseID uuid.UUID  `json:"supabase_id"`
	PfpKey     *string    `json:"pfp_key"`
	ResetTime  *time.Time `json:"reset_time"`
	ExpiresAt  *time.Time `json:"expires_at"`
}

// I represented year under the assumption that we're following the american high school timeline
// not sure if we want to change this to assume otherwise
type Student struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"user_id"`
	Year          Year       `json:"year"`
	Interest      *string    `json:"student_interest"`
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
		Name      string     `json:"name" db:"name" doc:"Name of the user" minLength:"1" maxLength:"200"`
		PfpKey    *string    `json:"pfp_key" db:"pfp_key" doc:"pfp key of user"`
		Email     string     `json:"email" doc:"email of user for supabase signup"`
		Role      Role       `json:"role" enum:"student,counselor" doc:"the type of user the counselor would like to create"`
		ExpiresAt *time.Time `json:"expires_at,omitempty" required:"false" doc:"when the account expires; omit for no expiration"`
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

type FetchStudentsInput struct {
	UserID uuid.UUID `query:"user_id" required:"false" doc:"user id of student"`
	Name   string    `query:"name" required:"false" doc:"name of student"`
	Cursor string    `query:"cursor" required:"false" doc:"next_cursor from the previous page"`
	Limit  int       `query:"limit" default:"5" minimum:"1" maximum:"10"`
}

type FetchStudentsOutput struct {
	Body *FetchStudentsBody
}

type StudentListItem struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"user_id"`
	Name          string     `json:"name"`
	PfpKey        *string    `json:"pfp_key"`
	Year          Year       `json:"year"`
	Organization  *string    `json:"organization"`
	GPA           int        `json:"gpa"`
	ReviewBalance int        `json:"review_balance"`
	CounselorID   *uuid.UUID `json:"counselor_id"`
}

type FetchStudentsBody struct {
	Students   []StudentListItem `json:"students"`
	NextCursor *string           `json:"next_cursor" doc:"Pass as cursor to get the next page; null when there are none"`
}

type FetchCounselorsInput struct {
	UserID uuid.UUID `query:"user_id" required:"false" doc:"user id of counselor"`
	Name   string    `query:"name" required:"false" doc:"name of counselor"`
	Cursor string    `query:"cursor" required:"false" doc:"next_cursor from the previous page"`
	Limit  int       `query:"limit" default:"5" minimum:"1" maximum:"10"`
}

type FetchCounselorsOutput struct {
	Body *FetchCounselorsBody
}

type CounselorListItem struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
	PfpKey *string   `json:"pfp_key"`
}

type FetchCounselorsBody struct {
	Counselors []CounselorListItem `json:"counselors"`
	NextCursor *string             `json:"next_cursor" doc:"Pass as cursor to get the next page; null when there are none"`
}

type AssignCounselorInput struct {
	Body struct {
		CounselorID uuid.UUID `json:"counselor_id" db:"counselor_id" doc:"id of the counselor to assign"`
		UserID      uuid.UUID `json:"user_id" db:"user_id" doc:"user id of the student being assigned"`
	}
}

type AssignCounselorOutput struct {
	Body Student
}
