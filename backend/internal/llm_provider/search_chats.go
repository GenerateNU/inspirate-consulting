package llmprovider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"inspirate-consulting/internal/auth"
)

// searchChats runs the search_chats tool: finds the user's messages closest in meaning to the prompt
// and lists them with who sent them and when.
func (l *LLMProvider) searchChats(ctx context.Context, prompt string) (*string, error) {
	if l.Search == nil {
		return nil, errors.New("search_chats: search is not configured")
	}

	results, err := l.Search.SearchChats(ctx, prompt)
	if err != nil {
		return nil, err
	}

	var answer string
	if len(results) == 0 {
		answer = "No matching messages found."
		return &answer, nil
	}

	userID, _ := auth.GetUserID(ctx)
	var b strings.Builder
	for _, result := range results {
		sender := "Them"
		if result.SenderID == userID {
			sender = "You"
		}
		fmt.Fprintf(&b, "- %s (%s): %s\n", sender, result.CreatedAt.Format("Jan 2, 2006"), result.Message)
	}

	answer = b.String()
	return &answer, nil
}
