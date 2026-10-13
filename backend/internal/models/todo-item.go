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
	EssayID         *string    `json:"essay_id"`
	MediaID         *string    `json:"media_id"`
	GlobalCollegeID *int64     `json:"global_college_id"`
}

type CreateTodoItemRequestBody struct {
	StudentID       string     `json:"student_id,omitempty"`
	UserID          string     `json:"user_id,omitempty"`
	TodoDescription string     `json:"todo_description"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	Deadline        *time.Time `json:"deadline,omitempty"`
	EssayID         *string    `json:"essay_id,omitempty" doc:"Essay this task is for; at most one of essay_id, media_id, global_college_id"`
	MediaID         *string    `json:"media_id,omitempty" doc:"Media this task is for; at most one of essay_id, media_id, global_college_id"`
	GlobalCollegeID *int64     `json:"global_college_id,omitempty" doc:"College this task is for; at most one of essay_id, media_id, global_college_id"`
}

type CreateTodoItemInput struct {
	Body CreateTodoItemRequestBody
}

type CreateTodoItemOutput struct {
	Body TodoItem
}

type TodoItemFilters struct {
	StudentID       string `query:"student_id" required:"false" doc:"Student whose to-do items to return; defaults to the authenticated student"`
	Status          string `query:"status" enum:"all,completed,incomplete" default:"all" doc:"Filter by completion status"`
	EssayID         string `query:"essay_id" required:"false" doc:"Only return items linked to this essay"`
	MediaID         string `query:"media_id" required:"false" doc:"Only return items linked to this media"`
	GlobalCollegeID int64  `query:"global_college_id" required:"false" doc:"Only return items linked to this college"`
	LinkedTo        string `query:"linked_to" enum:"essay,media,college" required:"false" doc:"Only return items linked to any item of this type"`
	SortBy          string `query:"sort_by" enum:"created_at,deadline,completed_at,updated_at" default:"created_at" doc:"Field to sort by"`
	SortOrder       string `query:"sort_order" enum:"asc,desc" default:"desc" doc:"Sort direction"`
}

type GetTodoItemsByStudentInput struct {
	TodoItemFilters
	Limit  int `query:"limit" default:"20" minimum:"1" maximum:"100" doc:"Maximum number of to-do items to return"`
	Offset int `query:"offset" default:"0" minimum:"0" doc:"Number of to-do items to skip before returning results"`
}

type GetTodoItemsByStudentOutput struct {
	Body []TodoItem
}

type SearchTodoItemsInput struct {
	TodoItemFilters
	Q      string `query:"q" minLength:"1" maxLength:"200" doc:"Text to search for in to-do descriptions"`
	Limit  int    `query:"limit" default:"10" minimum:"1" maximum:"50" doc:"Maximum number of to-do items to return"`
	Offset int    `query:"offset" default:"0" minimum:"0" doc:"Number of to-do items to skip before returning results"`
}

type SearchTodoItemsOutput struct {
	Body []TodoItem
}

type UpdateTodoItemCompletedAtInput struct {
	ID   string `path:"id"`
	Body struct {
		Completed bool `json:"completed"`
	}
}

type UpdateTodoItemCompletedAtOutput struct {
	Body TodoItem
}
