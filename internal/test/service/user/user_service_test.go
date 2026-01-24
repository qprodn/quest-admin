package user_test

import (
	"context"
	"quest-admin/internal/biz/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserBiz struct {
	mock.Mock
}

func (m *MockUserBiz) GetUser(ctx context.Context, id string) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func TestUserService_GetUser(t *testing.T) {
	ctx := context.Background()
	mockBiz := new(MockUserBiz)

	mockBiz.On("GetUser", ctx, "user-1").Return(&user.User{
		ID:       "user-1",
		Username: "testuser",
	}, nil)

	user, err := mockBiz.GetUser(ctx, "user-1")

	assert.NoError(t, err)
	assert.Equal(t, "user-1", user.ID)
	mockBiz.AssertExpectations(t)
}
