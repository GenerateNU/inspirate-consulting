package login

import storage "inspirate-consulting/internal/data"

type Handler struct {
	LoginRepository storage.UserRepository
}

func NewHandler(loginRepository storage.UserRepository) *Handler {
	return &Handler{
		LoginRepository: loginRepository,
	}
}
