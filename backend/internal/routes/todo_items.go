package routes

import (
	"context"
	"net/http"

	"inspirate-consulting/internal/data"
	todoitem "inspirate-consulting/internal/handlers/todo_item"
	"inspirate-consulting/internal/models"

	"github.com/danielgtaylor/huma/v2"
)

func SetUpTodoItemRoutes(api huma.API, repository *data.Repository) {
	todoItemHandler := todoitem.NewHandler(repository.TodoItem)

	//Register POST /todo-items handler
	huma.Register(api, huma.Operation{
		OperationID: "create-todo-item",
		Method:      http.MethodPost,
		Path:        "/todo-items",
		Description: "Create a to-do item with a description and optional deadline",
		Tags:        []string{"Todo Items"},
	}, func(ctx context.Context, input *models.CreateTodoItemInput) (*models.CreateTodoItemOutput, error) {
		created, err := todoItemHandler.CreateTodoItem(ctx, &input.Body)
		if err != nil {
			return nil, err
		}

		return &models.CreateTodoItemOutput{Body: *created}, nil

	})

	//Register GET /todo-items handler
	huma.Register(api, huma.Operation{
		OperationID: "get-todo-items",
		Method:      http.MethodGet,
		Path:        "/todo-items",
		Description: "Get a student's to-do items",
		Tags:        []string{"Todo Items"},
	}, func(ctx context.Context, input *models.GetTodoItemsByStudentInput)(*models.GetTodoItemsByStudentOutput, error){
		items, err := todoItemHandler.GetTodoItemsByStudent(ctx)
		if(err != nil){
			return nil, err
		}

		return &models.GetTodoItemsByStudentOutput{Body: items}, nil
	})

	//Register PATCH /todo-items/{id} handler
	huma.Register(api, huma.Operation{
		OperationID: "update-todo-item-completed",
		Method: http.MethodPatch,
		Path: "/todo-items/{id}",
		Description: "Update a to-do item's completion status",
		Tags: []string{"Todo Items"},
	}, func(ctx context.Context, input *models.UpdateTodoItemCompletedAtInput) (*models.UpdateTodoItemCompletedAtOutput, error) {
		updated, err := todoItemHandler.UpdateTodoItemCompletedAt(ctx, input.ID, input.Body.Completed)
		if(err != nil){
			return nil, err
		}

		return &models.UpdateTodoItemCompletedAtOutput{Body: *updated}, err
	})

}
