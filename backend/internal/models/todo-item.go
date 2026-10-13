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
	EssayID         *string    `json:"essay_id,omitempty" format:"uuid" doc:"Essay this task is for; at most one of essay_id, media_id, global_college_id"`
	MediaID         *string    `json:"media_id,omitempty" format:"uuid" doc:"Media this task is for; at most one of essay_id, media_id, global_college_id"`
	GlobalCollegeID *int64     `json:"global_college_id,omitempty" doc:"College this task is for; at most one of essay_id, media_id, global_college_id"`
}

type CreateTodoItemInput struct {
	Body CreateTodoItemRequestBody
}

type CreateTodoItemOutput struct {
	Body TodoItem
}

type TodoItemFilters struct {
	StudentID       string `query:"student_id" required:"false" format:"uuid" doc:"Student whose to-do items to return; defaults to the authenticated student"`
	Status          string `query:"status" enum:"all,completed,incomplete" default:"all" doc:"Filter by completion status"`
	EssayID         string `query:"essay_id" required:"false" format:"uuid" doc:"Only return items linked to this essay"`
	MediaID         string `query:"media_id" required:"false" format:"uuid" doc:"Only return items linked to this media"`
	GlobalCollegeID int64  `query:"global_college_id" required:"false" doc:"Only return items linked to this college"`
	LinkedTo        string `query:"linked_to" enum:"essay,media,college" required:"false" doc:"Only return items linked to any item of this type"`
	SortBy          string `query:"sort_by" enum:"created_at,deadline,completed_at,updated_at" default:"created_at" doc:"Field to sort by"`
	SortOrder       string `query:"sort_order" enum:"asc,desc" default:"desc" doc:"Sort direction"`
}

type TodoItemQuery struct {
	StudentID       string
	Status          string
	EssayID         *string
	MediaID         *string
	GlobalCollegeID *int64
	LinkedTo        *string
	Search          *string
	SortBy          string
	SortOrder       string
	After           *TodoItemCursor
	Limit           int
}

type TodoItemCursor struct {
	SortBy    string    `json:"s"`
	SortOrder string    `json:"o"`
	SortValue time.Time `json:"v"`
	ID        string    `json:"i"`
}

// Missing deadline/completed_at sort as these values so they come last in both directions and the cursor key is never null.
var (
	todoItemNullSortAsc  = time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC)
	todoItemNullSortDesc = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
)

func (q *TodoItemQuery) ApplyDefaults() {
	if q.Status == "" {
		q.Status = "all"
	}
	if q.SortBy == "" {
		q.SortBy = "created_at"
	}
	if q.SortOrder == "" {
		q.SortOrder = "desc"
	}
	if q.Limit <= 0 {
		q.Limit = 20
	}
}

func (q TodoItemQuery) NullSortValue() time.Time {
	if q.SortOrder == "asc" {
		return todoItemNullSortAsc
	}
	return todoItemNullSortDesc
}

func (q TodoItemQuery) SortValue(item TodoItem) time.Time {
	var value *time.Time
	switch q.SortBy {
	case "deadline":
		value = item.Deadline
	case "completed_at":
		value = item.CompletedAt
	case "updated_at":
		value = &item.UpdatedAt
	default:
		value = &item.CreatedAt
	}
	if value == nil {
		return q.NullSortValue()
	}
	return *value
}

type ListTodoItemsInput struct {
	TodoItemFilters
	Cursor string `query:"cursor" required:"false" doc:"next_cursor from the previous page; omit to get the first page"`
	Limit  int    `query:"limit" default:"20" minimum:"1" maximum:"100" doc:"Maximum number of to-do items to return"`
}

type TodoItemPage struct {
	Items      []TodoItem `json:"items"`
	NextCursor *string    `json:"next_cursor" doc:"Pass as cursor to get the next page; null when there are none"`
}

type ListTodoItemsOutput struct {
	Body TodoItemPage
}

type SearchTodoItemsInput struct {
	TodoItemFilters
	Q      string `query:"q" required:"true" minLength:"1" maxLength:"200" doc:"Text to search for in to-do descriptions"`
	Cursor string `query:"cursor" required:"false" doc:"next_cursor from the previous page; omit to get the first page"`
	Limit  int    `query:"limit" default:"10" minimum:"1" maximum:"50" doc:"Maximum number of to-do items to return"`
}

type SearchTodoItemsOutput struct {
	Body TodoItemPage
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
