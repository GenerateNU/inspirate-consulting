package greeting

import storage "github.com/GenerateNU/inspirate-consulting/internal/data"

// The greeting Handler has a dependency on GreetingStore
// in other words, greeting handler logic depends on there being a
// GreetingStore we can call on to pass down for DB operations
type Handler struct {
	GreetingStore storage.GreetingStore
}

// Here is where we create the handler which has reference to the GreetingStore
func NewHandler(greetingStore storage.GreetingStore) *Handler {
	return &Handler{
		GreetingStore: greetingStore,
	}
}
