package user_test

import (
	"context"
	"quest-admin/internal/biz/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserDeptBiz struct {
	mock.Mock
}

func (m *MockUserDeptBiz) AssignUserDepts(ctx context.Context, bo *user.AssignUserDeptsBO) error {
	args := m.Called(ctx, bo)
	return args.Error(0)
}

func TestUserDeptService_AssignUserDepts(t *testing.T) {
	ctx := context.Background()
	mockBiz := new(MockUserDeptBiz)

	bo := &user.AssignUserDeptsBO{
		UserID:  "user-1",
		DeptIDs: []string{"dept-1", "dept-2"},
	}

	mockBiz.On("AssignUserDepts", ctx, mock.AnythingOfType("*user.AssignUserDeptsBO")).Return(nil)

	err := mockBiz.AssignUserDepts(ctx, bo)

	assert.NoError(t, err)
	mockBiz.AssertExpectations(t)
}
