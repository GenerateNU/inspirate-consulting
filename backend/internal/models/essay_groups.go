package models

import "github.com/google/uuid"

type EssayGroups struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	/* Description is optional, the column is nullable */
	Description *string   `json:"description" required:"false"`
	StudentID   uuid.UUID `json:"student_id"`
}

type CreateEssayGroupBody struct {
	Name        string    `json:"name" minLength:"1" maxLength:"100" example:"Common App" doc:"Unique name for the group"`
	Description *string   `json:"description" required:"false" doc:"Optional description of the group"`
	StudentID   uuid.UUID `json:"student_id" doc:"Student the group belongs to"`
}

type CreateEssayGroupInput struct {
	Body CreateEssayGroupBody
}

type CreateEssayGroupOutput struct {
	Body EssayGroups
}

type GetEssayGroupsInput struct {
	StudentID uuid.UUID `path:"student_id" doc:"Student whose essay groups to list"`
}

type EssayGroupListBody struct {
	EssayGroups []EssayGroups `json:"essay_groups"`
}

type GetEssayGroupsOutput struct {
	Body EssayGroupListBody
}
