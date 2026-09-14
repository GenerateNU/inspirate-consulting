package repomocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/GenerateNU/inspirate-consulting/internal/models"
)

type MockGreetingStore struct {
	mock.Mock
}

func (m *MockGreetingStore) MockCreateGreeting(ctx context.Context, greeting *models.CreateGreetingInput) (*models.CreateGreetingOutput, error) {
	args := m.Called(ctx, greeting)
	if args.Get(0) == nil {
		if args.Get(1) == nil {
			return nil, nil
		}
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*models.CreateGreetingOutput), nil
}
