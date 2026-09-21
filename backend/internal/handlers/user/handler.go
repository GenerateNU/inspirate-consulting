package user

import (
	storage "inspirate-consulting/internal/data"
)

// The greeting Handler has a dependency on GreetingRepository
// in other words, greeting handler logic depends on there being a
// GreetingRepository we can call on to pass down for DB operations
type Handler struct {
	UserRepository storage.UserRepository
}

// Here is where we create the handler which has reference to the GreetingRepository
func NewHandler(userRepository storage.UserRepository) *Handler {
	return &Handler{
		UserRepository: userRepository,
	}
}
