package routes

import (
	"context"
	"net/http"

	"inspirate-consulting/internal/data"
	chatmessage "inspirate-consulting/internal/handlers/chat_message"
	"inspirate-consulting/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

func SetUpChatMessageRoutes(api huma.API, repository *data.Repository) {
	chatMessageHandler := chatmessage.NewHandler(repository.ChatMessage)

	// Register POST /messages handler.
	huma.Register(api, huma.Operation{
		OperationID: "create-chat-message",
		Method:      http.MethodPost,
		Path:        "/messages",
		Description: "Send a message from the current user to another user.",
		Tags:        []string{"Chat"},
	}, func(ctx context.Context, input *models.CreateChatMessageInput) (*models.CreateChatMessageOutput, error) {
		created, err := chatMessageHandler.CreateChatMessage(ctx, input.Body)
		if err != nil {
			return nil, err
		}
		return &models.CreateChatMessageOutput{Body: *created}, nil
	})

	// Register GET /chats handler.
	huma.Register(api, huma.Operation{
		OperationID: "list-chats",
		Method:      http.MethodGet,
		Path:        "/chats",
		Description: "List the current user's conversations with their latest message and unread count, most recent first.",
		Tags:        []string{"Chat"},
	}, func(ctx context.Context, input *models.ListChatsInput) (*models.ListChatsOutput, error) {
		chats, err := chatMessageHandler.ListChats(ctx)
		if err != nil {
			return nil, err
		}
		return &models.ListChatsOutput{Body: chats}, nil
	})

	// Register GET /chats/{user_id}/messages handler.
	huma.Register(api, huma.Operation{
		OperationID: "list-chat-messages",
		Method:      http.MethodGet,
		Path:        "/chats/{user_id}/messages",
		Description: "Get a page of messages between the current user and another user, newest first.",
		Tags:        []string{"Chat"},
	}, func(ctx context.Context, input *models.ListChatMessagesInput) (*models.ListChatMessagesOutput, error) {
		page, err := chatMessageHandler.ListChatMessages(ctx, input.UserID, input.Cursor, input.Limit)
		if err != nil {
			return nil, err
		}
		return &models.ListChatMessagesOutput{Body: *page}, nil
	})

	// Register POST /chats/{user_id}/read handler.
	huma.Register(api, huma.Operation{
		OperationID:   "mark-chat-read",
		Method:        http.MethodPost,
		Path:          "/chats/{user_id}/read",
		Description:   "Mark every message the other user sent to the current user as read.",
		Tags:          []string{"Chat"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, input *models.MarkChatReadInput) (*models.MarkChatReadOutput, error) {
		if err := chatMessageHandler.MarkChatRead(ctx, input.UserID); err != nil {
			return nil, err
		}
		return &models.MarkChatReadOutput{}, nil
	})

	// Register PATCH /messages/{id} handler.
	huma.Register(api, huma.Operation{
		OperationID: "edit-chat-message",
		Method:      http.MethodPatch,
		Path:        "/messages/{id}",
		Description: "Edit the text of a message the current user sent.",
		Tags:        []string{"Chat"},
	}, func(ctx context.Context, input *models.EditChatMessageInput) (*models.EditChatMessageOutput, error) {
		edited, err := chatMessageHandler.EditChatMessage(ctx, input.ID, input.Body.Message)
		if err != nil {
			return nil, err
		}
		return &models.EditChatMessageOutput{Body: *edited}, nil
	})

	// Register PATCH /messages/{id}/read handler.
	huma.Register(api, huma.Operation{
		OperationID: "update-chat-message-read",
		Method:      http.MethodPatch,
		Path:        "/messages/{id}/read",
		Description: "Mark a message the current user received as read or unread.",
		Tags:        []string{"Chat"},
	}, func(ctx context.Context, input *models.UpdateChatMessageReadInput) (*models.UpdateChatMessageReadOutput, error) {
		updated, err := chatMessageHandler.UpdateChatMessageReadAt(ctx, input.ID, input.Body.Read)
		if err != nil {
			return nil, err
		}
		return &models.UpdateChatMessageReadOutput{Body: *updated}, nil
	})
}
