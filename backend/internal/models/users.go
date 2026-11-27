package models

import (
	"time"

	"github.com/google/uuid"
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

type Student struct {
	ID                    uuid.UUID  `json:"id"`
	UserID                uuid.UUID  `json:"user_id"`
	GraduationYear        *int       `json:"graduation_year"`
	Interest              *string    `json:"student_interest"`
	Organization          *string    `json:"organization"`
	GPA                   *float64   `json:"gpa"`
	ReviewBalance         int        `json:"review_balance"`
	CounselorID           *uuid.UUID `json:"counselor_id"`
	HighSchool            *string    `json:"high_school"`
	StandardizedTestScore *string    `json:"standardized_test_score"`
	IntendedMajor         *string    `json:"intended_major"`
	IntendedMinor         *string    `json:"intended_minor"`
	GeographicPreferences *string    `json:"geographic_preferences"`
}

type Counselor struct {
	ID         uuid.UUID  `json:"id"`
	UserID     uuid.UUID  `json:"user_id"`
	CustomLink *string    `json:"custom_link"`
	ExpiresAt  *time.Time `json:"expires_at"`
}

type CreateUserInput struct {
	Body struct {
		Name      string     `json:"name" db:"name" doc:"Name of the user" minLength:"1" maxLength:"200"`
		Email     string     `json:"email" doc:"email of user for supabase signup"`
		Role      Role       `json:"role" enum:"student,counselor" doc:"the type of user the counselor would like to create"`
		ExpiresAt *time.Time `json:"expires_at,omitempty" required:"false" doc:"when the account expires; omit for no expiration"`
		// student profile fields, only used when role is student
		GraduationYear *int       `json:"graduation_year,omitempty" required:"false" minimum:"2000" maximum:"2100" doc:"year the student graduates high school; optional, the student can add it later"`
		GPA            *float64   `json:"gpa,omitempty" required:"false" minimum:"0" maximum:"4" doc:"unweighted gpa of the student (0 to 4); optional, the student can add it later"`
		Organization   *string    `json:"organization,omitempty" required:"false" doc:"organization of the student"`
		Interest       *string    `json:"student_interest,omitempty" required:"false" doc:"interest of the student"`
		CounselorID    *uuid.UUID `json:"counselor_id,omitempty" required:"false" doc:"counselor the student is assigned to; required when role is student"`
	}
}

type CreateUserBody struct {
	User *User `json:"user"`
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

// pfp_key is changed only through /user/pfp, which checks the key belongs to the user
type UpdateUserBody struct {
	Name *string `json:"name,omitempty" required:"false" minLength:"1" maxLength:"200" doc:"Name of the user"`
}

// user id is derived from session, users can only edit themselves
type UpdateUserInput struct {
	Body UpdateUserBody
}

type UpdateUserOutput struct {
	Body *User
}

// SessionUser is returned by login: the user plus their student or counselor profile id
type SessionUser struct {
	User
	StudentID   *uuid.UUID `json:"student_id" nullable:"true" doc:"id of the user's student profile; null when the user is not a student"`
	CounselorID *uuid.UUID `json:"counselor_id" nullable:"true" doc:"id of the user's counselor profile; null when the user is not a counselor"`
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
	ID             uuid.UUID  `json:"id"`
	UserID         uuid.UUID  `json:"user_id"`
	Name           string     `json:"name"`
	PfpKey         *string    `json:"pfp_key"`
	GraduationYear *int       `json:"graduation_year"`
	Organization   *string    `json:"organization"`
	GPA            *float64   `json:"gpa"`
	ReviewBalance  int        `json:"review_balance"`
	CounselorID    *uuid.UUID `json:"counselor_id"`
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
		CounselorID *uuid.UUID `json:"counselor_id" db:"counselor_id" nullable:"true" doc:"id of the counselor to assign"`
		UserID      uuid.UUID  `json:"user_id" db:"user_id" doc:"user id of the student being assigned"`
	}
}

type AssignCounselorOutput struct {
	Body Student
}

// A user as shown on their profile, with a link to their picture
type UserProfile struct {
	User
	PfpURL *string `json:"pfp_url" doc:"Presigned, time-limited URL for the user's profile picture, null when they have none"`
}

type FetchUserProfileOutput struct {
	Body *UserProfile
}

type CurrentUser struct {
	UserProfile
	Role        Role       `json:"role" enum:"student,counselor" doc:"Role of the logged-in user"`
	StudentID   *uuid.UUID `json:"student_id" doc:"Student profile ID, null unless the user is a student"`
	CounselorID *uuid.UUID `json:"counselor_id" doc:"Counselor profile ID, null unless the user is a counselor"`
}

type GetCurrentUserOutput struct {
	Body *CurrentUser
}
