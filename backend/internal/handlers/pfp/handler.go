package pfp

import storage "inspirate-consulting/internal/data"

type Handler struct {
	UserRepository storage.UserRepository
	PfpRepository  storage.PfpRepository
}

func NewHandler(userRepository storage.UserRepository, pfpRepository storage.PfpRepository) *Handler {
	return &Handler{
		UserRepository: userRepository,
		PfpRepository:  pfpRepository,
	}
}
