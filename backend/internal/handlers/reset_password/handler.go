package resetpassword

import (
	storage "inspirate-consulting/internal/data"
	"inspirate-consulting/internal/data/aws/sqsclient"
)

type Handler struct {
	ResetPasswordRepository storage.UserRepository
	SQSClient               *sqsclient.SQSClient
}

func NewHandler(resetpasswordRepository storage.UserRepository, sqsClient *sqsclient.SQSClient) *Handler {
	return &Handler{
		ResetPasswordRepository: resetpasswordRepository,
		SQSClient:               sqsClient,
	}
}
