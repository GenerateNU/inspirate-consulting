package greeting

import storage "inspirate-consulting/internal/data"

// The greeting Handler has a dependency on GreetingRepository
// in other words, greeting handler logic depends on there being a
// GreetingRepository we can call on to pass down for DB operations
type Handler struct {
	GreetingRepository storage.GreetingRepository
}

// Here is where we create the handler which has reference to the GreetingRepository
func NewHandler(greetingRepository storage.GreetingRepository) *Handler {
	return &Handler{
		GreetingRepository: greetingRepository,
	}
}
