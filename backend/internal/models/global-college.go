package models

import "time"

// To represent a row in the global_colleges table.
type GlobalCollege struct {
	ID             int64      `json:"id" example:"1" doc:"Unique identifier for the college"`
	CreatedAt      time.Time  `json:"created_at" example:"2023-01-01T00:00:00Z" doc:"Timestamp when the college was created"`
	UpdatedAt      time.Time  `json:"updated_at" example:"2023-01-01T00:00:00Z" doc:"Timestamp when the college was last updated"`
	SchoolName     string     `json:"school_name" example:"Harvard University" doc:"Name of the college"`
	SchoolLocation string     `json:"school_location" example:"Cambridge, MA" doc:"Location of the college"`
	EADeadline     *time.Time `json:"ea_deadline,omitempty" example:"2023-01-01T00:00:00Z" doc:"Early admission deadline"`
	EDDeadline     *time.Time `json:"ed_deadline,omitempty" example:"2023-01-01T00:00:00Z" doc:"Early decision deadline"`
	RDDeadline     *time.Time `json:"rd_deadline,omitempty" example:"2023-01-01T00:00:00Z" doc:"Regular decision deadline"`
}

// To represent the input needed to create a new global college entry.
type CreateGlobalCollegeRequestBody struct {
	SchoolName       string     `json:"school_name" example:"Harvard University" doc:"Name of the college"`
	SchoolLocation   string     `json:"school_location" example:"Cambridge, MA" doc:"Location of the college"`
	EADeadline       *time.Time `json:"ea_deadline,omitempty" example:"2023-01-01T00:00:00Z" doc:"Early admission deadline"`
	EDDeadline       *time.Time `json:"ed_deadline,omitempty" example:"2023-01-01T00:00:00Z" doc:"Early decision deadline"`
	RDDeadline       *time.Time `json:"rd_deadline,omitempty" example:"2023-01-01T00:00:00Z" doc:"Regular decision deadline"`
	ConfirmDuplicate bool       `json:"confirm_duplicate,omitempty" doc:"Set true to create despite a detected near-duplicate warning"`
}

// to represent input needed to update an existing global college entry
type UpdateGlobalCollegeRequestBody struct {
	SchoolName       string     `json:"school_name" example:"Harvard University" doc:"Name of the college"`
	SchoolLocation   string     `json:"school_location" example:"Cambridge, MA" doc:"Location of the college"`
	EADeadline       *time.Time `json:"ea_deadline,omitempty" example:"2023-01-01T00:00:00Z" doc:"Early admission deadline"`
	EDDeadline       *time.Time `json:"ed_deadline,omitempty" example:"2023-01-01T00:00:00Z" doc:"Early decision deadline"`
	RDDeadline       *time.Time `json:"rd_deadline,omitempty" example:"2023-01-01T00:00:00Z" doc:"Regular decision deadline"`
	ConfirmDuplicate bool       `json:"confirm_duplicate,omitempty" doc:"Set true to create despite a detected near-duplicate warning"`
}

// Huma readable input and output models for the create, get, and list operations for global colleges.
type CreateGlobalCollegeInput struct {
	Body CreateGlobalCollegeRequestBody
}

type CreateGlobalCollegeOutput struct {
	Body GlobalCollege
}

type GetGlobalCollegeInput struct {
	ID int64 `path:"id"`
}

type GetGlobalCollegeOutput struct {
	Body GlobalCollege
}

type ListGlobalCollegesInput struct{}

type ListGlobalCollegesOutput struct {
	Body []GlobalCollege
}

type UpdateGlobalCollegeInput struct {
	ID   int64 `path:"id" doc:"ID of the college to update"`
	Body UpdateGlobalCollegeRequestBody
}

type UpdateGlobalCollegeOutput struct {
	Body GlobalCollege
}

type DeleteGlobalCollegeInput struct {
	ID int64 `path:"id" doc:"ID of the college to delete"`
}

type DeleteGlobalCollegeOutput struct{}

type SearchGlobalCollegesInput struct {
	Query string `query:"q" doc:"Search text to match against college names"`
}

type SearchGlobalCollegesOutput struct {
	Body []GlobalCollege
}
