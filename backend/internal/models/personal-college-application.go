package models

import "time"

// To represent a student's application to a specific global college.
type PersonalCollegeApplication struct {
	ID              int64     `json:"id" doc:"Unique identifier for the application"`
	StudentID       string    `json:"student_id" doc:"ID of the student who owns this application"`
	GlobalCollegeID int64     `json:"global_college_id" doc:"ID of the global college being applied to"`
	ApplicationType string    `json:"application_type" example:"ED" doc:"Which deadline the student is applying by: EA, ED, or RD"`
	Category        string    `json:"category" example:"reach" doc:"How the student categorizes this school: safety, target, or reach"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Rank            *int      `json:"rank" doc:"Student's ranking of this application relative to others, or null if unranked"`
}

// to represent input needed to create a new personal college application (student id intenitonally omitted, derived from session)
type CreatePersonalCollegeApplicationRequestBody struct {
	GlobalCollegeID int64  `json:"global_college_id" doc:"ID of the global college being applied to"`
	ApplicationType string `json:"application_type" enum:"EA,ED,RD" doc:"Which deadline the student is applying by"`
	Category        string `json:"category" enum:"safety,target,reach" doc:"How the student categorizes this school"`
}

// to represent input needed to update the relative ranking of a college
type UpdatePersonalCollegeApplicationRankRequestBody struct {
	Rank *int `json:"rank" doc:"The new 1-based rank, or null to un-rank this application"`
}

// to represent input needed to make a general update
type UpdatePersonalCollegeApplicationRequestBody struct {
	GlobalCollegeID int64  `json:"global_college_id" doc:"ID of the global college being applied to"`
	ApplicationType string `json:"application_type" example:"ED" doc:"Which deadline the student is applying by: EA, ED, or RD"`
	Category        string `json:"category" example:"reach" doc:"How the student categorizes this school: safety, target, or reach"`
}

// Huma readable input and output models for the create and list operations
type CreatePersonalCollegeApplicationInput struct {
	Body CreatePersonalCollegeApplicationRequestBody
}

type CreatePersonalCollegeApplicationOutput struct {
	Body PersonalCollegeApplication
}

// student id is derived from session, so no input is needed for list operation
type ListPersonalCollegeApplicationsInput struct{}

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
