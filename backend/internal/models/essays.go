package models

import "github.com/google/uuid"

type Status string

const (
	Submitted Status = "Submitted"
	Draft Status = "Draft"
	Review Status = "Review"
	Archived Status = "Archived"
)

type Essays struct {
	ID uuid.UUID `json.:"id"`
	StudentID uuid.UUID `json:"student_id"`
	type string `json:"type"`
	CollegeID *uuid.UUID `json:"college_id"`
	LinkToContent string `json:"link_to_content"`
	/* I added this to help the studetns keep track of their progress*/
	Status Status 
}