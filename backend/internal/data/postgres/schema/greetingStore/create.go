package greetingRepository

import (
	"context"
	"fmt"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"
)

func (r *GreetingRepository) CreateGreeting(ctx context.Context, greeting models.CreateGreetingInput) (*models.CreateGreetingOutput, error) {
	createdGreeting := &models.CreateGreetingOutput{}
	greetingMessage := fmt.Sprintf("Hello, %s!", greeting.Body.Name)

	insertQuery, err := schema.ReadSQLBaseScript("create_greeting.sql", SqlGreetingFiles)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(
		ctx,
		insertQuery,
		greeting.Body.Name,
		greetingMessage,
	).Scan(&createdGreeting.Body.Message)
	if err != nil {
		return nil, err
	}

	return createdGreeting, nil
}
