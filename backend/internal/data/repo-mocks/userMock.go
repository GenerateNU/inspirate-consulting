package mocks

import (
	context "context"
	"time"

	models "inspirate-consulting/internal/models"

	"github.com/google/uuid"

	mock "github.com/stretchr/testify/mock"
)

type UserRepository struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: ctx, user
func (_m *UserRepository) CreateUser(ctx context.Context, user models.CreateUserInput, id uuid.UUID) (*models.CreateUserOutput, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 *models.CreateUserOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.CreateUserInput) (*models.CreateUserOutput, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.CreateUserInput) *models.CreateUserOutput); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.CreateUserOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.CreateUserInput) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetGlobalCollege provides a mock function with given fields: ctx, id
func (_m *UserRepository) FetchUser(ctx context.Context, input models.FetchUserInput) (*models.FetchUserOutput, error) {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for GetGlobalCollege")
	}

	var r0 *models.FetchUserOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.FetchUserInput) (*models.FetchUserOutput, error)); ok {
		return rf(ctx, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.FetchUserInput) *models.FetchUserOutput); ok {
		r0 = rf(ctx, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.FetchUserOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.FetchUserInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FetchUserBySupabaseID provides a mock function with given fields: ctx, input
func (_m *UserRepository) FetchUserBySupabaseID(ctx context.Context, input models.FetchUserBySupabaseIDInput) (*models.FetchUserOutput, error) {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for FetchUserBySupabaseID")
	}

	var r0 *models.FetchUserOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.FetchUserBySupabaseIDInput) (*models.FetchUserOutput, error)); ok {
		return rf(ctx, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.FetchUserBySupabaseIDInput) *models.FetchUserOutput); ok {
		r0 = rf(ctx, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.FetchUserOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.FetchUserBySupabaseIDInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *UserRepository) UpdateResetTime(ctx context.Context, id uuid.UUID, resetTime *time.Time) error {
	ret := _m.Called(ctx, id, resetTime)

	if len(ret) == 0 {
		panic("no return value specified for UpdateResetTime")
	}

	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, *time.Time) error); ok {
		return rf(ctx, id, resetTime)
	}

	return ret.Error(0)
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRepository {
	mock := &UserRepository{}
	mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
