package models

import (
	"time"

	"github.com/google/uuid"
)

// ReviewEntryType is what kind of balance movement a ledger row records.
type ReviewEntryType string

const (
	ReviewEntrySpend      ReviewEntryType = "spend"
	ReviewEntryRefund     ReviewEntryType = "refund"
	ReviewEntryAdjustment ReviewEntryType = "adjustment"
)

// Derived in the database from completed_at and refund.
type ReviewStatus string

const (
	ReviewStatusOpen      ReviewStatus = "open"
	ReviewStatusCompleted ReviewStatus = "completed"
	ReviewStatusRefunded  ReviewStatus = "refunded"
)

// EssayReviewTransaction is a row in the essay_review_transaction ledger.
type EssayReviewTransaction struct {
	ID          uuid.UUID       `json:"id" example:"7a1f0c9e-0b4e-4f1a-9f2d-0c3b5a6d7e8f" doc:"Unique identifier for the transaction"`
	Subtotal    int             `json:"subtotal" example:"-1" doc:"Signed change to the student's balance. Negative for a spend, positive for a refund or grant. Sums to the balance."`
	StudentID   uuid.UUID       `json:"student_id" example:"3f2b1a0c-9d8e-4c7b-8a6f-5e4d3c2b1a09" doc:"Student whose balance moved"`
	EntryType   ReviewEntryType `json:"entry_type" example:"spend" doc:"Kind of movement: spend, refund, or adjustment"`
	EssayID     *uuid.UUID      `json:"essay_id,omitempty" example:"b2c3d4e5-f6a7-4859-9a0b-1c2d3e4f5a6b" doc:"Essay under review. Null for an adjustment."`
	CompletedAt *time.Time      `json:"completed_at,omitempty" example:"2026-01-01T00:00:00Z" doc:"When a counselor marked the review complete"`
	Refund      *uuid.UUID      `json:"refund,omitempty" example:"c3d4e5f6-a7b8-4960-ab1c-2d3e4f5a6b7c" doc:"The negative row that reversed this charge, if it was refunded"`
	Status      *ReviewStatus   `json:"status,omitempty" example:"open" doc:"Lifecycle of a spend: open, completed, or refunded. Null for refund and adjustment rows."`
	ActorID     *uuid.UUID      `json:"actor_id,omitempty" example:"d4e5f6a7-b8c9-4a71-bc2d-3e4f5a6b7c8d" doc:"User who caused the movement. Null until auth provides the caller."`
	CreatedAt   time.Time       `json:"created_at" example:"2026-01-01T00:00:00Z" doc:"Timestamp when the transaction was created"`
	UpdatedAt   time.Time       `json:"updated_at" example:"2026-01-01T00:00:00Z" doc:"Timestamp when the transaction was last updated"`
}

// To represent the input needed to open a review against an essay.
type RequestEssayReviewRequestBody struct {
	StudentID uuid.UUID `json:"student_id" example:"3f2b1a0c-9d8e-4c7b-8a6f-5e4d3c2b1a09" doc:"Student requesting the review"`
	EssayID   uuid.UUID `json:"essay_id" example:"b2c3d4e5-f6a7-4859-9a0b-1c2d3e4f5a6b" doc:"Essay to review"`
	Amount    int       `json:"amount" minimum:"1" default:"1" example:"1" doc:"Number of review credits to spend. Positive; recorded on the ledger as a negative subtotal."`
}

// To represent the result of refunding a charge: the charge itself, now
// linked to the reversal row that undid it.
type RefundEssayReviewResponseBody struct {
	Transaction EssayReviewTransaction `json:"transaction" doc:"The original charge, now marked refunded"`
	Refund      EssayReviewTransaction `json:"refund" doc:"The negative row that reversed the charge"`
}

// Huma readable input and output models for the essay review operations.
type RequestEssayReviewInput struct {
	Body RequestEssayReviewRequestBody
}

type RequestEssayReviewOutput struct {
	Body EssayReviewTransaction
}

type RefundEssayReviewInput struct {
	ID uuid.UUID `path:"id" doc:"ID of the charge to refund"`
}

type RefundEssayReviewOutput struct {
	Body RefundEssayReviewResponseBody
}

type CompleteEssayReviewInput struct {
	ID uuid.UUID `path:"id" doc:"ID of the review to mark complete"`
}

type CompleteEssayReviewOutput struct {
	Body EssayReviewTransaction
}

// EssayReviewStatus is the current review state of an essay. TransactionID is
// what the refund and complete endpoints take.
type EssayReviewStatus struct {
	TransactionID uuid.UUID    `json:"transaction_id" example:"7a1f0c9e-0b4e-4f1a-9f2d-0c3b5a6d7e8f" doc:"Ledger row this status came from"`
	EssayID       uuid.UUID    `json:"essay_id" example:"b2c3d4e5-f6a7-4859-9a0b-1c2d3e4f5a6b" doc:"Essay the review is for"`
	StudentID     uuid.UUID    `json:"student_id" example:"3f2b1a0c-9d8e-4c7b-8a6f-5e4d3c2b1a09" doc:"Student who requested the review"`
	Status        ReviewStatus `json:"status" example:"open" doc:"open, completed, or refunded"`
	RequestedAt   time.Time    `json:"requested_at" example:"2026-01-01T00:00:00Z" doc:"When the review was requested"`
	CompletedAt   *time.Time   `json:"completed_at,omitempty" example:"2026-01-01T00:00:00Z" doc:"When a counselor completed it"`
}

type GetEssayReviewStatusInput struct {
	EssayID uuid.UUID `path:"essay_id" doc:"ID of the essay to look up"`
}

type GetEssayReviewStatusOutput struct {
	Body EssayReviewStatus
}

// To represent the input needed to overwrite a student's review balance.
type SetStudentReviewBalanceRequestBody struct {
	ReviewBalance int `json:"review_balance" minimum:"0" example:"5" doc:"Absolute value to set the student's review balance to"`
}

type SetStudentReviewBalanceInput struct {
	ID   uuid.UUID `path:"id" doc:"ID of the student"`
	Body SetStudentReviewBalanceRequestBody
}

type SetStudentReviewBalanceOutput struct {
	Body Student
}
