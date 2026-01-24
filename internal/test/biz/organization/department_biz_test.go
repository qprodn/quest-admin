package organization_test

import (
	"context"
	"testing"

	org "quest-admin/internal/biz/organization"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDepartmentRepo struct {
	mock.Mock
}

func (m *MockDepartmentRepo) Create(ctx context.Context, dept *org.Department) (*org.Department, error) {
	args := m.Called(ctx, dept)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*org.Department), args.Error(1)
}

func (m *MockDepartmentRepo) FindByID(ctx context.Context, id string) (*org.Department, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*org.Department), args.Error(1)
}

func (m *MockDepartmentRepo) FindByName(ctx context.Context, name string) (*org.Department, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*org.Department), args.Error(1)
}

func (m *MockDepartmentRepo) List(ctx context.Context) ([]*org.Department, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*org.Department), args.Error(1)
}

func (m *MockDepartmentRepo) FindByParentID(ctx context.Context, parentID string) ([]*org.Department, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*org.Department), args.Error(1)
}

func (m *MockDepartmentRepo) Update(ctx context.Context, dept *org.Department) (*org.Department, error) {
	args := m.Called(ctx, dept)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*org.Department), args.Error(1)
}

func (m *MockDepartmentRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDepartmentRepo) HasUsers(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDepartmentRepo) FindListByIDs(ctx context.Context, ids []string) ([]*org.Department, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*org.Department), args.Error(1)
}

func TestDepartmentRepo_Create(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockDepartmentRepo)

	dept := &org.Department{
		ID:   "dept-1",
		Name: "Engineering",
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*organization.Department")).Return(dept, nil)

	result, err := mockRepo.Create(ctx, dept)

	assert.NoError(t, err)
	assert.Equal(t, "dept-1", result.ID)
	mockRepo.AssertExpectations(t)
}

func TestDepartmentRepo_FindByID(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockDepartmentRepo)

	expected := &org.Department{
		ID:   "dept-1",
		Name: "Engineering",
	}

	mockRepo.On("FindByID", ctx, "dept-1").Return(expected, nil)

	result, err := mockRepo.FindByID(ctx, "dept-1")

	assert.NoError(t, err)
	assert.Equal(t, "dept-1", result.ID)
	mockRepo.AssertExpectations(t)
}

func TestDepartmentRepo_List(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockDepartmentRepo)

	depts := []*org.Department{
		{ID: "dept-1", Name: "Engineering"},
	}

	mockRepo.On("List", ctx).Return(depts, nil)

	result, err := mockRepo.List(ctx)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	mockRepo.AssertExpectations(t)
}

func TestDepartmentRepo_HasUsers(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockDepartmentRepo)

	mockRepo.On("HasUsers", ctx, "dept-1").Return(true, nil)

	hasUsers, err := mockRepo.HasUsers(ctx, "dept-1")

	assert.NoError(t, err)
	assert.True(t, hasUsers)
	mockRepo.AssertExpectations(t)
}
