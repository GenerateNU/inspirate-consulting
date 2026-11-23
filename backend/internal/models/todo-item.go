package models

import (
	"time"

	"github.com/google/uuid"
)

type TodoLinkType string

const (
	TodoLinkCollege         TodoLinkType = "college"
	TodoLinkEssay           TodoLinkType = "essay"
	TodoLinkExtracurricular TodoLinkType = "extracurricular"
	TodoLinkModule          TodoLinkType = "module"
)

type TodoItem struct {
	ID                string        `json:"id"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
	StudentID         string        `json:"student_id"`
	UserID            string        `json:"user_id"`
	TodoDescription   string        `json:"todo_description"`
	CompletedAt       *time.Time    `json:"completed_at"`
	Deadline          *time.Time    `json:"deadline"`
	EssayID           *string       `json:"essay_id"`
	MediaID           *string       `json:"media_id"`
	GlobalCollegeID   *int64        `json:"global_college_id"`
	Title             string        `json:"title"`
	ExtracurricularID *int64        `json:"extracurricular_id"`
	LinkType          *TodoLinkType `json:"link_type" enum:"college,essay,extracurricular,module"`
	AssignmentID      *uuid.UUID    `json:"assignment_id"`
}

type CreateTodoItemRequestBody struct {
	StudentID         string        `json:"student_id,omitempty" required:"false" format:"uuid" doc:"Student the task is for; required for counselors, defaults to the authenticated student"`
	TodoDescription   string        `json:"todo_description"`
	CompletedAt       *time.Time    `json:"completed_at,omitempty"`
	Deadline          *time.Time    `json:"deadline,omitempty"`
	EssayID           *string       `json:"essay_id,omitempty" format:"uuid" doc:"Essay this task is for; at most one of essay_id, media_id, global_college_id"`
	MediaID           *string       `json:"media_id,omitempty" format:"uuid" doc:"Media this task is for; at most one of essay_id, media_id, global_college_id"`
	GlobalCollegeID   *int64        `json:"global_college_id,omitempty" doc:"College this task is for; at most one of essay_id, media_id, global_college_id"`
	Title             string        `json:"title" minLength:"1" maxLength:"200" doc:"Short title of the task"`
	ExtracurricularID *int64        `json:"extracurricular_id,omitempty" required:"false" doc:"Extracurricular this task is linked to"`
	LinkType          *TodoLinkType `json:"link_type,omitempty" required:"false" enum:"college,essay,extracurricular,module" doc:"What kind of item this task is linked to"`
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
	LinkedTo        string `query:"linked_to" enum:"college,essay,extracurricular,media,module" required:"false" doc:"Only return items linked to any item of this type"`
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

type DeleteTodoInput struct {
	ID string `path:"id"`
}

type UpdateTodoItemDescriptionInput struct {
	ID   string `path:"id"`
	Body struct {
		NewDescription string `json:"new_description" minLength:"1"`
	}
}

type UpdateTodoItemDescriptionOutput struct {
	Body TodoItem
}

type UpdateTodoItemDeadlineInput struct {
	ID   string `path:"id" format:"uuid"`
	Body struct {
		NewDeadline *time.Time `json:"deadline"`
	}
}

type UpdateTodoItemDeadlineOutput struct {
	Body TodoItem
}

type BulkCreateTodoItemsBody struct {
	Title           string        `json:"title" minLength:"1" maxLength:"200" doc:"Short title of the task"`
	Description     string        `json:"description" doc:"Full description of the task"`
	Deadline        *time.Time    `json:"deadline,omitempty" required:"false" doc:"When the task is due"`
	LinkType        *TodoLinkType `json:"link_type,omitempty" required:"false" enum:"college,essay,extracurricular,module" doc:"What kind of item this task is linked to"`
	CollegeID       *int64        `json:"college_id,omitempty" required:"false" doc:"College to link, when link_type is college"`
	MediaID         *string       `json:"media_id,omitempty" required:"false" format:"uuid" doc:"Video to link, when link_type is module"`
	StudentIDs      []uuid.UUID   `json:"student_ids,omitempty" required:"false" doc:"Students to assign individually"`
	GraduationYears []int         `json:"graduation_years,omitempty" required:"false" doc:"Assign every student graduating in these years"`
	All             bool          `json:"all,omitempty" required:"false" doc:"Assign every one of the counselor's students"`
}

type BulkCreateTodoItemsInput struct {
	Body BulkCreateTodoItemsBody
}

type BulkCreateTodoItemsResponse struct {
	AssignmentID uuid.UUID  `json:"assignment_id"`
	Items        []TodoItem `json:"items"`
}

type BulkCreateTodoItemsOutput struct {
	Body BulkCreateTodoItemsResponse
}

type BulkCreateTodoItemsParams struct {
	CounselorID uuid.UUID
	UserID      uuid.UUID
	Title       string
	Description string
	Deadline    *time.Time
	LinkType    *TodoLinkType
	CollegeID   *int64
	MediaID     *string
	StudentIDs  []uuid.UUID
	Years       []int
	All         bool
}

type UpdateTodoAssignmentBody struct {
	Title            *string       `json:"title,omitempty" required:"false" minLength:"1" maxLength:"200" doc:"New title for every task in the assignment"`
	Description      *string       `json:"description,omitempty" required:"false" doc:"New description for every task in the assignment"`
	Deadline         *time.Time    `json:"deadline,omitempty" required:"false" doc:"New deadline for every task in the assignment"`
	LinkType         *TodoLinkType `json:"link_type,omitempty" required:"false" enum:"college,module" doc:"New link; when set, the college and media ids are replaced too"`
	CollegeID        *int64        `json:"college_id,omitempty" required:"false"`
	MediaID          *string       `json:"media_id,omitempty" required:"false" format:"uuid"`
	AddStudentIDs    []uuid.UUID   `json:"add_student_ids,omitempty" required:"false" doc:"Students to add to the assignment"`
	RemoveStudentIDs []uuid.UUID   `json:"remove_student_ids,omitempty" required:"false" doc:"Students to remove; completed tasks are kept"`
}

type UpdateTodoAssignmentInput struct {
	AssignmentID uuid.UUID `path:"assignment_id"`
	Body         UpdateTodoAssignmentBody
}

type UpdateTodoAssignmentOutput struct {
	Body []TodoItem
}

// Passed from the handler to the store; not part of the API.
type UpdateTodoAssignmentParams struct {
	AssignmentID     uuid.UUID
	CounselorID      uuid.UUID
	UserID           uuid.UUID
	Title            *string
	Description      *string
	Deadline         *time.Time
	LinkType         *TodoLinkType
	CollegeID        *int64
	MediaID          *string
	AddStudentIDs    []uuid.UUID
	RemoveStudentIDs []uuid.UUID
}

type UpdateTodoItemTitleBody struct {
	Title string `json:"title" minLength:"1" maxLength:"200" doc:"New title of the task"`
}

type UpdateTodoItemTitleInput struct {
	ID   string `path:"id" format:"uuid" doc:"ID of the todo item to update"`
	Body UpdateTodoItemTitleBody
}

type UpdateTodoItemTitleOutput struct {
	Body *TodoItem
}

type UpdateTodoItemLinkBody struct {
	LinkType          *TodoLinkType `json:"link_type" nullable:"true" enum:"college,essay,extracurricular,module" doc:"What the task links to, or null to remove its link"`
	EssayID           *string       `json:"essay_id,omitempty" required:"false" format:"uuid" doc:"Essay to link, when link_type is essay"`
	MediaID           *string       `json:"media_id,omitempty" required:"false" format:"uuid" doc:"Video to link, when link_type is module"`
	GlobalCollegeID   *int64        `json:"global_college_id,omitempty" required:"false" doc:"College to link, when link_type is college"`
	ExtracurricularID *int64        `json:"extracurricular_id,omitempty" required:"false" doc:"Extracurricular to link, when link_type is extracurricular"`
}

type UpdateTodoItemLinkInput struct {
	ID   string `path:"id" format:"uuid" doc:"ID of the todo item to update"`
	Body UpdateTodoItemLinkBody
}

type UpdateTodoItemLinkOutput struct {
	Body *TodoItem
}

// Passed from the handler to the store; not part of the API.
type TodoItemLink struct {
	LinkType          *TodoLinkType
	EssayID           *string
	MediaID           *string
	GlobalCollegeID   *int64
	ExtracurricularID *int64
}
