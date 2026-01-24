package permission_test

import (
	"context"
	"quest-admin/internal/biz/permission"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRoleBiz struct {
	mock.Mock
}

func (m *MockRoleBiz) GetRole(ctx context.Context, id string) (*permission.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*permission.Role), args.Error(1)
}

func TestRoleService_GetRole(t *testing.T) {
	ctx := context.Background()
	mockBiz := new(MockRoleBiz)

	mockBiz.On("GetRole", ctx, "role-1").Return(&permission.Role{
		ID:   "role-1",
		Name: "Admin",
	}, nil)

	role, err := mockBiz.GetRole(ctx, "role-1")

	assert.NoError(t, err)
	assert.Equal(t, "role-1", role.ID)
	mockBiz.AssertExpectations(t)
}
