package models

import "github.com/google/uuid"

type Status string

const (
	Submitted Status = "Submitted"
	Draft     Status = "Draft"
	Review    Status = "Review"
	Archived  Status = "Archived"
)

type Essays struct {
	ID            uuid.UUID `json:"id"`
	StudentID     uuid.UUID `json:"student_id"`
	Type          string    `json:"type"`
	CollegeID     *int64    `json:"college_id"`
	LinkToContent string    `json:"link_to_content"`
	/* I added this to help the studetns keep track of their progress*/
	Status Status `json:"status"`
}

type GetEssaysFromStudentInput struct {
	StudentID uuid.UUID `path:"student_id" doc:"Student whose essays to list"`
}

type EssayListBody struct {
	Essays []Essays `json:"essays"`
}

type GetEssaysFromStudentOutput struct {
	Body EssayListBody
}

type CreateEssayBody struct {
	StudentID uuid.UUID `json:"student_id" doc:"Student the essay belongs to"`
	Type      string    `json:"type" maxLength:"100" example:"personal-statement"`
	CollegeID *int64    `json:"college_id" required:"false" doc:"Optional college this essay targets"`
	// Might integrate eiditing the essay on the website might change
	LinkToContent string `json:"link_to_content" example:"https://docs.google.com/document/d/abc123"`
}

type CreateEssayInput struct {
	Body CreateEssayBody
}

type CreateEssayOutput struct{}

type UpdateStatusBody struct {
	Status Status `json:"status" enum:"Submitted,Draft,Review,Archived" doc:"New status for the essay"`
}

type UpdateStatusInput struct {
	EssayID uuid.UUID `path:"essay_id" doc:"Essay to update"`
	Body    UpdateStatusBody
}

type UpdateStatusOutput struct{}
