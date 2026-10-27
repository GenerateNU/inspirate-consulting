package resetpassword

import (
	"context"
	"errors"
	"testing"

	"inspirate-consulting/internal/data/aws/sqsclient"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/supabase"

	awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type fakeSQS struct {
	lastInput *awssqs.SendMessageInput
}

func (f *fakeSQS) SendMessage(ctx context.Context, params *awssqs.SendMessageInput, optFns ...func(*awssqs.Options)) (*awssqs.SendMessageOutput, error) {
	f.lastInput = params
	return &awssqs.SendMessageOutput{}, nil
}

func TestSendTempPassword_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	supabaseID := uuid.New()

	user := &models.User{
		ID:         userID,
		Name:       "Ada Lovelace",
		Email:      "ada@example.com",
		SupabaseID: supabaseID,
	}

	mockRepo := mocks.NewUserRepository(t)
	mockRepo.On("FetchUser", mock.Anything, models.FetchUserInput{ID: userID}).
		Return(&models.FetchUserOutput{Body: user}, nil)
	mockRepo.On("UpdateResetTime", mock.Anything, supabaseID, mock.Anything).Return(nil)

	fake := &fakeSQS{}
	sqsClient := sqsclient.NewSQSClientWithAPI(fake, "test-account-queue-url")
	handler := NewHandler(mockRepo, sqsClient)

	res, err := handler.SendTempPassword(ctx, &models.SendTempPasswordInput{ID: userID}, &supabase.MockSupabase{})

	assert.NoError(t, err)
	assert.Equal(t, "ada@example.com", res.Body.Email)
	assert.NotNil(t, fake.lastInput, "expected a password-reset email to be queued")
}

func TestSendTempPassword_UserNotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()

	mockRepo := mocks.NewUserRepository(t)
	mockRepo.On("FetchUser", mock.Anything, models.FetchUserInput{ID: userID}).
		Return(nil, errors.New("user not found"))

	sqsClient := sqsclient.NewSQSClientWithAPI(&fakeSQS{}, "test-account-queue-url")
	handler := NewHandler(mockRepo, sqsClient)

	res, err := handler.SendTempPassword(ctx, &models.SendTempPasswordInput{ID: userID}, &supabase.MockSupabase{})

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestSendTempPassword_UpdateResetTimeError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	supabaseID := uuid.New()

	user := &models.User{
		ID:         userID,
		Name:       "Ada Lovelace",
		Email:      "ada@example.com",
		SupabaseID: supabaseID,
	}

	mockRepo := mocks.NewUserRepository(t)
	mockRepo.On("FetchUser", mock.Anything, models.FetchUserInput{ID: userID}).
		Return(&models.FetchUserOutput{Body: user}, nil)
	mockRepo.On("UpdateResetTime", mock.Anything, supabaseID, mock.Anything).
		Return(errors.New("db error"))

	sqsClient := sqsclient.NewSQSClientWithAPI(&fakeSQS{}, "test-account-queue-url")
	handler := NewHandler(mockRepo, sqsClient)

	res, err := handler.SendTempPassword(ctx, &models.SendTempPasswordInput{ID: userID}, &supabase.MockSupabase{})

	assert.Error(t, err)
	assert.Nil(t, res)
}
