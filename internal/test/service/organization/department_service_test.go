package organization_test

import (
	"context"
	"quest-admin/internal/biz/organization"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDepartmentBiz struct {
	mock.Mock
}

func (m *MockDepartmentBiz) FindByID(ctx context.Context, id string) (*organization.Department, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*organization.Department), args.Error(1)
}

func TestDepartmentService_FindByID(t *testing.T) {
	ctx := context.Background()
	mockBiz := new(MockDepartmentBiz)
	mockBiz.On("FindByID", ctx, "dept-1").Return(&organization.Department{
		ID:   "dept-1",
		Name: "Engineering",
	}, nil)

	dept, err := mockBiz.FindByID(ctx, "dept-1")

	assert.NoError(t, err)
	assert.Equal(t, "dept-1", dept.ID)
	mockBiz.AssertExpectations(t)
}
