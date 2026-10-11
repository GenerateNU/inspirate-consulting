package notification

import "inspirate-consulting/internal/data/aws/sqsclient"

type Handler struct{
	SQSClient *sqsclient.SQSClient
}

func NewHandler(sqsClient *sqsclient.SQSClient)(*Handler){
	return &Handler{SQSClient : sqsClient}
}

