package models

import (
	"time"

	"github.com/google/uuid"
)

// To represent a row in the chat_messages table.
type ChatMessage struct {
	ID          uuid.UUID  `json:"id"`
	SenderID    uuid.UUID  `json:"sender_id"`
	RecipientID uuid.UUID  `json:"recipient_id"`
	Message     string     `json:"message"`
	CreatedAt   time.Time  `json:"created_at"`
	EditedAt    *time.Time `json:"edited_at" doc:"Null if the message was never edited"`
	ReadAt      *time.Time `json:"read_at" doc:"Null if the recipient has not read the message"`
}

// Represents a row in the current user's chat list.
type ChatSummary struct {
	OtherUserID       uuid.UUID  `json:"other_user_id"`
	OtherUserName     string     `json:"other_user_name"`
	OtherUserPfpKey   *string    `json:"other_user_pfp_key"`
	LastMessageID     uuid.UUID  `json:"last_message_id"`
	LastMessage       string     `json:"last_message"`
	LastSenderID      uuid.UUID  `json:"last_sender_id"`
	LastMessageAt     time.Time  `json:"last_message_at"`
	LastMessageReadAt *time.Time `json:"last_message_read_at"`
	UnreadCount       int        `json:"unread_count" doc:"Messages from the other user the current user has not read"`
}

type CreateChatMessageRequestBody struct {
	RecipientID uuid.UUID `json:"recipient_id" doc:"User ID of the message recipient"`
	Message     string    `json:"message" minLength:"1" doc:"Message text"`
}

type EditChatMessageRequestBody struct {
	Message string `json:"message" minLength:"1" doc:"New message text"`
}

// Huma readable input and output models for the chat operations
type CreateChatMessageInput struct {
	Body CreateChatMessageRequestBody
}

type CreateChatMessageOutput struct {
	Body ChatMessage
}

type ListChatsInput struct{}

type ListChatsOutput struct {
	Body []ChatSummary
}

// keyset pagination: pass the id of the oldest message already loaded as before to get the next page
type ListChatMessagesInput struct {
	UserID uuid.UUID `path:"user_id" doc:"The other user in the conversation"`
	Before string    `query:"before" format:"uuid" required:"false" doc:"Return messages older than this message ID"`
	Limit  int       `query:"limit" default:"50" minimum:"1" maximum:"100"`
}

type ListChatMessagesOutput struct {
	Body []ChatMessage
}

type EditChatMessageInput struct {
	ID   uuid.UUID `path:"id"`
	Body EditChatMessageRequestBody
}

type EditChatMessageOutput struct {
	Body ChatMessage
}

type UpdateChatMessageReadInput struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		Read bool `json:"read" doc:"False marks the message as unread"`
	}
}

type UpdateChatMessageReadOutput struct {
	Body ChatMessage
}

type MarkChatReadInput struct {
	UserID uuid.UUID `path:"user_id" doc:"The other user in the conversation"`
}

type MarkChatReadOutput struct{}
