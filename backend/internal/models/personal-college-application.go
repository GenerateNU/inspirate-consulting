package models

import (
	"time"

	"github.com/google/uuid"
)

// To represent a student's application to a specific global college.
type PersonalCollegeApplication struct {
	ID                int64     `json:"id" doc:"Unique identifier for the application"`
	StudentID         uuid.UUID `json:"student_id" doc:"ID of the student who owns this application"`
	GlobalCollegeID   int64     `json:"global_college_id" doc:"ID of the global college being applied to"`
	ApplicationType   string    `json:"application_type" example:"ED" doc:"Which deadline the student is applying by: EA, ED, or RD"`
	Category          string    `json:"category" example:"reach" doc:"How the student categorizes this school: safety, target, or reach"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Rank              *int      `json:"rank" doc:"Student's ranking of this application relative to others, or null if unranked"`
	Interest          *string   `json:"interest" enum:"top_choice,strong_contender,considering" doc:"How interested the student is in this school, or null if not set"`
	ApplicationStatus string    `json:"application_status" enum:"not_started,in_progress,submitted" example:"not_started" doc:"Where the student is in submitting this application"`
	CycleOffered      bool      `json:"cycle_offered" doc:"False when the college no longer offers application_type (e.g. a counselor removed it); show a warning"`
	EssayProgress     string    `json:"essay_progress" enum:"not_started,in_progress,completed" example:"in_progress" doc:"Computed from the student's non-archived essays for this college: none, some not yet submitted, or all submitted"`
}

// Essay progress for an application, computed from the student's essays for that college.
const (
	EssayProgressNotStarted = "not_started"
	EssayProgressInProgress = "in_progress"
	EssayProgressCompleted  = "completed"
)

// Application cycles a college can offer and a student can apply by.
const (
	ApplicationCycleED      = "ED"
	ApplicationCycleED2     = "ED2"
	ApplicationCycleREA     = "REA"
	ApplicationCycleEA      = "EA"
	ApplicationCycleRD      = "RD"
	ApplicationCycleRolling = "ROLLING"
)

// How interested a student is in a school (School Details).
const (
	InterestTopChoice       = "top_choice"
	InterestStrongContender = "strong_contender"
	InterestConsidering     = "considering"
)

// Where a student is in submitting an application.
const (
	ApplicationStatusNotStarted = "not_started"
	ApplicationStatusInProgress = "in_progress"
	ApplicationStatusSubmitted  = "submitted"
)

// to represent input needed to create a new personal college application (student id intenitonally omitted, derived from session)
type CreatePersonalCollegeApplicationRequestBody struct {
	GlobalCollegeID int64  `json:"global_college_id" doc:"ID of the global college being applied to"`
	ApplicationType string `json:"application_type" enum:"ED,ED2,REA,EA,RD,ROLLING" doc:"Which application cycle the student is applying by; the college must offer it"`
	Category        string `json:"category" enum:"safety,target,reach" doc:"How the student categorizes this school"`
}

// to represent input needed to update the relative ranking of a college
type UpdatePersonalCollegeApplicationRankRequestBody struct {
	Rank *int `json:"rank" doc:"The new 1-based rank, or null to un-rank this application"`
}

// to represent input needed to make a general update (School Details save).
// This is a full replacement: leaving interest out clears it.
type UpdatePersonalCollegeApplicationRequestBody struct {
	GlobalCollegeID   int64   `json:"global_college_id" doc:"ID of the global college being applied to"`
	ApplicationType   string  `json:"application_type" enum:"ED,ED2,REA,EA,RD,ROLLING" doc:"Which application cycle the student is applying by; the college must offer it"`
	Category          string  `json:"category" enum:"safety,target,reach" doc:"How the student categorizes this school"`
	Interest          *string `json:"interest,omitempty" required:"false" enum:"top_choice,strong_contender,considering" doc:"How interested the student is in this school; omit to clear"`
	ApplicationStatus string  `json:"application_status" enum:"not_started,in_progress,submitted" doc:"Where the student is in submitting this application"`
}

// Huma readable input and output models for the create and list operations
type CreatePersonalCollegeApplicationInput struct {
	Body CreatePersonalCollegeApplicationRequestBody
}

type CreatePersonalCollegeApplicationOutput struct {
	Body PersonalCollegeApplication
}

// students list their own applications; counselors pass student_id (must be one of their students)
type ListPersonalCollegeApplicationsInput struct {
	StudentID string `query:"student_id" required:"false" format:"uuid" doc:"Student whose college list to return; required for counselors, defaults to the authenticated student"`
}

type ListPersonalCollegeApplicationsOutput struct {
	Body []PersonalCollegeApplication
}

type UpdatePersonalCollegeApplicationRankInput struct {
	ID   int64 `path:"id" doc:"ID of the application to re-rank"`
	Body UpdatePersonalCollegeApplicationRankRequestBody
}

// returns all applications with updated ranks
type UpdatePersonalCollegeApplicationRankOutput struct {
	Body []PersonalCollegeApplication
}

type UpdatePersonalCollegeApplicationInput struct {
	ID   int64 `path:"id" doc:"ID of the application to update"`
	Body UpdatePersonalCollegeApplicationRequestBody
}

type UpdatePersonalCollegeApplicationOutput struct {
	Body PersonalCollegeApplication
}

type DeletePersonalCollegeApplicationInput struct {
	ID int64 `path:"id" doc:"ID of the application to delete"`
}

type DeletePersonalCollegeApplicationOutput struct {
}
