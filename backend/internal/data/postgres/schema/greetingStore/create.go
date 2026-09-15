package greetingRepository

import (
	"context"
	"fmt"

	"inspirate-consulting/internal/models"
)

func (r *GreetingRepository) CreateGreeting(ctx context.Context, greeting models.CreateGreetingInput) (*models.CreateGreetingOutput, error) {
	createdGreeting := &models.CreateGreetingOutput{}
	greetingMessage := fmt.Sprintf("Hello, %s!", greeting.Body.Name)

	const insertQuery = `
	INSERT INTO public.greetings (
		name, greeting_string
	) VALUES (
		$1, $2
	)
	RETURNING greeting_string
	`

	err := r.db.QueryRow(
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
