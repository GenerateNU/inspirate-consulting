package models

import "github.com/google/uuid"

type EssayGroups struct {
	ID uuid.UUID `json:"id"`
	Name string `json:name`
	Description string    `json:"description"`
    StudentID   uuid.UUID `json:"student_id"`
}

