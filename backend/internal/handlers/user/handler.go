package user

import (
	storage "inspirate-consulting/internal/data"
	"inspirate-consulting/internal/data/aws/sqsclient"
)

type Handler struct {
	UserRepository storage.UserRepository
	SQSClient      *sqsclient.SQSClient
}

func NewHandler(userRepository storage.UserRepository, sqsClient *sqsclient.SQSClient) *Handler {
	return &Handler{
		UserRepository: userRepository,
		SQSClient:      sqsClient,
	}
}
