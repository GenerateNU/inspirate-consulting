package models

import "time"

// To represent a row in the global_colleges table.
type GlobalCollege struct {
	ID int64 `json:"id" example:"1" doc:"Unique identifier for the college"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z" doc:"Timestamp when the college was created"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z" doc:"Timestamp when the college was last updated"`
	SchoolName string `json:"school_name" example:"Harvard University" doc:"Name of the college"`
	SchoolLocation string `json:"school_location" example:"Cambridge, MA" doc:"Location of the college"`
	EADeadline *time.Time `json:"ea_deadline,omitempty" example:"2023-01-01T00:00:00Z" doc:"Early admission deadline"`
	EDDeadline *time.Time `json:"ed_deadline,omitempty" example:"2023-01-01T00:00:00Z" doc:"Early decision deadline"`
	RDDeadline *time.Time `json:"rd_deadline,omitempty" example:"2023-01-01T00:00:00Z" doc:"Regular decision deadline"`
}

// To represent the input needed to create a new global college entry.
type CreateGlobalCollegeInput struct {
	SchoolName string `json:"school_name" example:"Harvard University" doc:"Name of the college"`
	SchoolLocation string `json:"school_location" example:"Cambridge, MA" doc:"Location of the college"`
	EADeadline *time.Time `json:"ea_deadline,omitempty" example:"2023-01-01T00:00:00Z" doc:"Early admission deadline"`
	EDDeadline *time.Time `json:"ed_deadline,omitempty" example:"2023-01-01T00:00:00Z" doc:"Early decision deadline"`
	RDDeadline *time.Time `json:"rd_deadline,omitempty" example:"2023-01-01T00:00:00Z" doc:"Regular decision deadline"`
}