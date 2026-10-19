package claimpassword

type Handler struct {
	Secret string
}

func NewHandler(secret string) *Handler {
	return &Handler{Secret: secret}
}
