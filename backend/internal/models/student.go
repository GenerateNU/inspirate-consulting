package models

import "github.com/google/uuid"

type GetStudentInput struct {
	ID uuid.UUID `path:"id"`
}

type GetStudentOutput struct {
	Body *Student
}
