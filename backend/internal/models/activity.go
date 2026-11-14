package models

import (
	"time"

	"github.com/google/uuid"
)

type Action string

const (
	ActionAddedSchool            Action = "added_school"
	ActionRemovedSchool          Action = "removed_school"
	ActionEssayStatusChanged     Action = "essay_status_changed"
	ActionAddedExtracurricular   Action = "added_extracurricular"
	ActionRemovedExtracurricular Action = "removed_extracurricular"
	ActionTaskPastDue            Action = "task_past_due"
)

type Activity struct {
	ID                uuid.UUID  `json:"id"`
	StudentID         uuid.UUID  `json:"student_id"`
	Action            Action     `json:"action" enum:"added_school,removed_school,essay_status_changed,added_extracurricular,removed_extracurricular,task_past_due"`
	EssayID           *uuid.UUID `json:"essay_id,omitempty"`
	EssayStatus       *Status    `json:"essay_status,omitempty" enum:"Submitted,Draft,Review,Archived"`
	GlobalCollegeID   *int64     `json:"global_college_id,omitempty"`
	ExtracurricularID *int64     `json:"extracurricular_id,omitempty"`
	TodoItemID        *uuid.UUID `json:"todo_item_id,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
}

type CreateActivityBody struct {
	Action            Action     `json:"action" enum:"added_school,removed_school,essay_status_changed,added_extracurricular,removed_extracurricular,task_past_due"`
	EssayID           *uuid.UUID `json:"essay_id,omitempty" required:"false"`
	EssayStatus       *Status    `json:"essay_status,omitempty" required:"false" enum:"Submitted,Draft,Review,Archived"`
	GlobalCollegeID   *int64     `json:"global_college_id,omitempty" required:"false"`
	ExtracurricularID *int64     `json:"extracurricular_id,omitempty" required:"false"`
	TodoItemID        *uuid.UUID `json:"todo_item_id,omitempty" required:"false"`
}
