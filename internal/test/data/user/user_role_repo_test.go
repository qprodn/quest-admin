package user_test

import (
	"context"
	"testing"

	user "quest-admin/internal/biz/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRoleRepo struct {
	mock.Mock
}

func (m *MockUserRoleRepo) GetUserRoles(ctx context.Context, userID string) ([]*user.UserRole, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user.UserRole), args.Error(1)
}

func (m *MockUserRoleRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRoleRepo) Create(ctx context.Context, item *user.UserRole) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func TestUserRoleRepo_GetUserRoles(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRoleRepo)

	roles := []*user.UserRole{
		{ID: "ur-1", UserID: "user-1", RoleID: "role-1"},
		{ID: "ur-2", UserID: "user-1", RoleID: "role-2"},
	}

	mockRepo.On("GetUserRoles", ctx, "user-1").Return(roles, nil)

	result, err := mockRepo.GetUserRoles(ctx, "user-1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestUserRoleRepo_Create(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRoleRepo)

	item := &user.UserRole{
		ID:     "ur-1",
		UserID: "user-1",
		RoleID: "role-1",
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*user.UserRole")).Return(nil)

	err := mockRepo.Create(ctx, item)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserRoleRepo_Delete(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRoleRepo)

	mockRepo.On("Delete", ctx, "ur-1").Return(nil)

	err := mockRepo.Delete(ctx, "ur-1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
