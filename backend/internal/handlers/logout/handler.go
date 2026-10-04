package logout

import storage "inspirate-consulting/internal/data"

type Handler struct {
	LogoutRepository storage.UserRepository
}

func NewHandler(logoutRepository storage.UserRepository) *Handler {
	return &Handler{
		LogoutRepository: logoutRepository,
	}
}
