package models

import "github.com/google/uuid"

type GetStudentInput struct {
	ID uuid.UUID `path:"id"`
}

type GetStudentOutput struct {
	Body *Student
}

// fields left out of the request are not changed
// counselor_id and review_balance have their own endpoints
type UpdateStudentBody struct {
	Year         *Year   `json:"year,omitempty" required:"false" enum:"freshman,sophomore,junior,senior" doc:"School year of the student"`
	GPA          *int    `json:"gpa,omitempty" required:"false" doc:"GPA of the student"`
	Organization *string `json:"organization,omitempty" required:"false" doc:"Organization of the student"`
	Interest     *string `json:"student_interest,omitempty" required:"false" doc:"Interest of the student"`
}

type UpdateStudentInput struct {
	ID   uuid.UUID `path:"id" doc:"ID of the student"`
	Body UpdateStudentBody
}

type UpdateStudentOutput struct {
	Body *Student
}
