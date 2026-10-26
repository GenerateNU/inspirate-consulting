package user

import (
	"context"
	"errors"
	"inspirate-consulting/internal/data/aws/sqsclient"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/supabase"
	"testing"

	awssqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// fakeSQS records the last message sent, so tests can assert on it without a real queue.
type fakeSQS struct {
	lastInput *awssqs.SendMessageInput
}

func (f *fakeSQS) SendMessage(ctx context.Context, params *awssqs.SendMessageInput, optFns ...func(*awssqs.Options)) (*awssqs.SendMessageOutput, error) {
	f.lastInput = params
	return &awssqs.SendMessageOutput{}, nil
}

// Unit tests for the CreateGlobalCollege handler logic without HTTP or database dependencies.
func TestHandler_CreateGlobalCollege(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	test_uuid := uuid.MustParse("7bebfe6e-36ca-4343-a46c-d9430a951ba7")
	s1 := uuid.New()
	key := "pfp-key"

	expectedUser := models.User{
		ID:         test_uuid,
		Name:       "Aleng123",
		SupabaseID: s1,
		PfpKey:     &key,
	}

	input := models.CreateUserInput{}
	input.Body.Name = "Aleng123"
	input.Body.Email = "tt@gmail.com"
	input.Body.PfpKey = &key

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		repoOutput := models.CreateUserOutput{Body: &models.CreateUserBody{User: &expectedUser}}

		mockRepo := mocks.NewUserRepository(t)
		mockRepo.On("CreateUser", mock.Anything, input, mock.Anything).Return(&repoOutput, nil)

		fake := &fakeSQS{}
		sqsClient := sqsclient.NewSQSClientWithAPI(fake, "test-account-queue-url")
		handler := NewHandler(mockRepo, sqsClient)
		res, err := handler.CreateUser(ctx, &input, &supabase.MockSupabase{})

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, *res.Body.User)
		assert.NotNil(t, fake.lastInput, "expected an account-created email to be queued")
	})

	t.Run("repository error propagates (e.g. duplicate conflict from the DB)", func(t *testing.T) {
		t.Parallel()

		mockRepo := mocks.NewUserRepository(t)
		mockRepo.On("CreateUser", mock.Anything, input, mock.Anything).
			Return(nil, errors.New("user already exists"))

		sqsClient := sqsclient.NewSQSClientWithAPI(&fakeSQS{}, "test-account-queue-url")
		handler := NewHandler(mockRepo, sqsClient)
		res, err := handler.CreateUser(ctx, &input, &supabase.MockSupabase{})

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "user already exists")
	})
}
