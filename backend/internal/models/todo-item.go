package models

import "time"

type TodoItem struct {
	ID              string     `json:"id"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	StudentID       string     `json:"student_id"`
	UserID          string     `json:"user_id"`
	TodoDescription string     `json:"todo_description"`
	CompletedAt     *time.Time `json:"completed_at"`
	Deadline        *time.Time `json:"deadline"`
}

type CreateTodoItemRequestBody struct {
	StudentID       string     `json:"student_id"`
	UserID          string     `json:"user_id"`
	TodoDescription string     `json:"todo_description"`
	CompletedAt     *time.Time `json:"completed_at"`
	Deadline        *time.Time `json:"deadline"`
}

type CreateTodoItemInput struct{
	Body CreateTodoItemRequestBody
}

type CreateTodoItemOutput struct{
	Body TodoItem
}


