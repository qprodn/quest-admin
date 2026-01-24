package permission_test

import (
	"context"
	"testing"
	"time"

	permission "quest-admin/internal/biz/permission"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRoleRepo struct {
	mock.Mock
}

func (m *MockRoleRepo) Create(ctx context.Context, role *permission.Role) (*permission.Role, error) {
	args := m.Called(ctx, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*permission.Role), args.Error(1)
}

func (m *MockRoleRepo) FindByID(ctx context.Context, id string) (*permission.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*permission.Role), args.Error(1)
}

func (m *MockRoleRepo) FindByName(ctx context.Context, name string) (*permission.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*permission.Role), args.Error(1)
}

func (m *MockRoleRepo) FindByCode(ctx context.Context, code string) (*permission.Role, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*permission.Role), args.Error(1)
}

func (m *MockRoleRepo) List(ctx context.Context, query *permission.ListRolesQuery) (*permission.ListRolesResult, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*permission.ListRolesResult), args.Error(1)
}

func (m *MockRoleRepo) Update(ctx context.Context, role *permission.Role) (*permission.Role, error) {
	args := m.Called(ctx, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*permission.Role), args.Error(1)
}

func (m *MockRoleRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRoleRepo) HasUsers(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockRoleRepo) FindListByIDs(ctx context.Context, roleIds []string) ([]*permission.Role, error) {
	args := m.Called(ctx, roleIds)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*permission.Role), args.Error(1)
}

func TestRoleRepo_Create(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRoleRepo)

	role := &permission.Role{
		ID:     "role-123",
		Name:   "Admin",
		Code:   "admin",
		Status: 1,
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*permission.Role")).Return(role, nil)

	result, err := mockRepo.Create(ctx, role)

	assert.NoError(t, err)
	assert.Equal(t, "role-123", result.ID)
	assert.Equal(t, "Admin", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestRoleRepo_FindByID(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRoleRepo)

	expectedRole := &permission.Role{
		ID:       "role-123",
		Name:     "Admin",
		Code:     "admin",
		Status:   1,
		TenantID: "tenant-1",
		CreateAt: time.Now(),
	}

	mockRepo.On("FindByID", ctx, "role-123").Return(expectedRole, nil)

	result, err := mockRepo.FindByID(ctx, "role-123")

	assert.NoError(t, err)
	assert.Equal(t, "role-123", result.ID)
	assert.Equal(t, "Admin", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestRoleRepo_Update(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRoleRepo)

	role := &permission.Role{
		ID:     "role-123",
		Name:   "Super Admin",
		Code:   "super_admin",
		Status: 1,
	}

	mockRepo.On("Update", ctx, mock.AnythingOfType("*permission.Role")).Return(role, nil)

	result, err := mockRepo.Update(ctx, role)

	assert.NoError(t, err)
	assert.Equal(t, "Super Admin", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestRoleRepo_Delete(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRoleRepo)

	mockRepo.On("Delete", ctx, "role-123").Return(nil)

	err := mockRepo.Delete(ctx, "role-123")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRoleRepo_HasUsers(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRoleRepo)

	mockRepo.On("HasUsers", ctx, "role-123").Return(true, nil)

	hasUsers, err := mockRepo.HasUsers(ctx, "role-123")

	assert.NoError(t, err)
	assert.True(t, hasUsers)
	mockRepo.AssertExpectations(t)
}

func TestRoleRepo_FindListByIDs(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRoleRepo)

	roles := []*permission.Role{
		{ID: "role-1", Name: "Admin"},
		{ID: "role-2", Name: "User"},
	}

	mockRepo.On("FindListByIDs", ctx, []string{"role-1", "role-2"}).Return(roles, nil)

	result, err := mockRepo.FindListByIDs(ctx, []string{"role-1", "role-2"})

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}
