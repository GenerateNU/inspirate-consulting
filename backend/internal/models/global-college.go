package models

import "time"

// To represent a row in the global_colleges table, with the application cycles it offers.
type GlobalCollege struct {
	ID             int64                     `json:"id" example:"1" doc:"Unique identifier for the college"`
	CreatedAt      time.Time                 `json:"created_at" example:"2023-01-01T00:00:00Z" doc:"Timestamp when the college was created"`
	UpdatedAt      time.Time                 `json:"updated_at" example:"2023-01-01T00:00:00Z" doc:"Timestamp when the college was last updated"`
	SchoolName     string                    `json:"school_name" example:"Harvard University" doc:"Name of the college"`
	SchoolLocation string                    `json:"school_location" example:"Cambridge, MA" doc:"Location of the college"`
	Cycles         []CollegeApplicationCycle `json:"cycles" doc:"Application cycles this college offers, in ED, ED2, REA, EA, RD, ROLLING order"`
}

// One application cycle a college offers (a row in college_application_cycles).
// Rolling admissions has no single deadline; it uses an optional rolling_start/rolling_end range instead.
type CollegeApplicationCycle struct {
	Cycle        string     `json:"cycle" enum:"ED,ED2,REA,EA,RD,ROLLING" example:"ED" doc:"Application cycle"`
	Deadline     *time.Time `json:"deadline,omitempty" example:"2023-11-01T00:00:00Z" doc:"Deadline for this cycle; omitted for rolling admissions"`
	RollingStart *string    `json:"rolling_start,omitempty" format:"date" example:"2023-09-01" doc:"Start of the rolling admissions window (ROLLING only)"`
	RollingEnd   *string    `json:"rolling_end,omitempty" format:"date" example:"2024-03-01" doc:"End of the rolling admissions window (ROLLING only)"`
}

// To represent the input needed to create a new global college entry.
type CreateGlobalCollegeRequestBody struct {
	SchoolName       string                    `json:"school_name" example:"Harvard University" doc:"Name of the college"`
	SchoolLocation   string                    `json:"school_location" example:"Cambridge, MA" doc:"Location of the college"`
	Cycles           []CollegeApplicationCycle `json:"cycles,omitempty" doc:"Application cycles this college offers"`
	ConfirmDuplicate bool                      `json:"confirm_duplicate,omitempty" doc:"Set true to create despite a detected near-duplicate warning"`
}

// to represent input needed to update an existing global college entry.
// cycles replaces the full set of offered cycles: any cycle left out is removed.
type UpdateGlobalCollegeRequestBody struct {
	SchoolName       string                    `json:"school_name" example:"Harvard University" doc:"Name of the college"`
	SchoolLocation   string                    `json:"school_location" example:"Cambridge, MA" doc:"Location of the college"`
	Cycles           []CollegeApplicationCycle `json:"cycles" doc:"The full set of application cycles this college offers; cycles left out are removed"`
	ConfirmDuplicate bool                      `json:"confirm_duplicate,omitempty" doc:"Set true to create despite a detected near-duplicate warning"`
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
	Query string `query:"query" doc:"Search text to match against college names"`
}

type SearchGlobalCollegesOutput struct {
	Body []GlobalCollege
}
