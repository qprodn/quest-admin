package user_test

import (
	"context"
	"quest-admin/internal/biz/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRoleBiz struct {
	mock.Mock
}

func (m *MockUserRoleBiz) AssignUserRoles(ctx context.Context, bo *user.AssignUserRolesBO) error {
	args := m.Called(ctx, bo)
	return args.Error(0)
}

func TestUserRoleService_AssignUserRoles(t *testing.T) {
	ctx := context.Background()
	mockBiz := new(MockUserRoleBiz)

	bo := &user.AssignUserRolesBO{
		UserID:  "user-1",
		RoleIDs: []string{"role-1", "role-2"},
	}

	mockBiz.On("AssignUserRoles", ctx, mock.AnythingOfType("*user.AssignUserRolesBO")).Return(nil)

	err := mockBiz.AssignUserRoles(ctx, bo)

	assert.NoError(t, err)
	mockBiz.AssertExpectations(t)
}
