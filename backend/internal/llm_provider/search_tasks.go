package llmprovider

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// searchTasks runs the search_tasks tool: finds todo items closest in meaning to the prompt
// and lists them with their due date and completion status.
func (l *LLMProvider) searchTasks(ctx context.Context, prompt string) (*string, error) {
	if l.Search == nil {
		return nil, errors.New("search_tasks: search is not configured")
	}

	results, err := l.Search.SearchTasks(ctx, prompt)
	if err != nil {
		return nil, err
	}

	var answer string
	if len(results) == 0 {
		answer = "No matching tasks found."
		return &answer, nil
	}

	var b strings.Builder
	for _, result := range results {
		due := "no due date"
		if result.Deadline != nil {
			due = "due " + result.Deadline.Format("Jan 2, 2006")
		}
		status := "not completed"
		if result.CompletedAt != nil {
			status = "completed"
		}
		fmt.Fprintf(&b, "- %s (%s, %s)\n", result.TodoDescription, due, status)
	}

	answer = b.String()
	return &answer, nil
}
