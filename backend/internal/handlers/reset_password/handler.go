package resetpassword

import storage "inspirate-consulting/internal/data"

type Handler struct {
	ResetPasswordRepository storage.UserRepository
}

func NewHandler(resetpasswordRepository storage.UserRepository) *Handler {
	return &Handler{
		ResetPasswordRepository: resetpasswordRepository,
	}
}
