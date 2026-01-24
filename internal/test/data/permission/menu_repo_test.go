package permission_test

import (
	"context"
	"testing"
	"time"

	permission "quest-admin/internal/biz/permission"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMenuRepo struct {
	mock.Mock
}

func (m *MockMenuRepo) FindByMenuIDs(ctx context.Context, menuIDs []string) ([]*permission.Menu, error) {
	args := m.Called(ctx, menuIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*permission.Menu), args.Error(1)
}

func (m *MockMenuRepo) Create(ctx context.Context, menu *permission.Menu) error {
	args := m.Called(ctx, menu)
	return args.Error(0)
}

func (m *MockMenuRepo) FindByID(ctx context.Context, id string) (*permission.Menu, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*permission.Menu), args.Error(1)
}

func (m *MockMenuRepo) FindByName(ctx context.Context, name string) (*permission.Menu, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*permission.Menu), args.Error(1)
}

func (m *MockMenuRepo) List(ctx context.Context) ([]*permission.Menu, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*permission.Menu), args.Error(1)
}

func (m *MockMenuRepo) FindByParentID(ctx context.Context, parentID string) ([]*permission.Menu, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*permission.Menu), args.Error(1)
}

func (m *MockMenuRepo) Update(ctx context.Context, menu *permission.Menu) error {
	args := m.Called(ctx, menu)
	return args.Error(0)
}

func (m *MockMenuRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestMenuRepo_Create(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockMenuRepo)

	menu := &permission.Menu{
		ID:   "menu-123123",
		Name: "Dashboard",
		Type: 1,
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*permission.Menu")).Return(nil)

	err := mockRepo.Create(ctx, menu)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMenuRepo_FindByID(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockMenuRepo)

	expectedMenu := &permission.Menu{
		ID:       "menu-123123",
		Name:     "Dashboard",
		Type:     1,
		Path:     "/dashboard",
		CreateAt: time.Now(),
	}

	mockRepo.On("FindByID", ctx, "menu-123123").Return(expectedMenu, nil)

	result, err := mockRepo.FindByID(ctx, "menu-123123")

	assert.NoError(t, err)
	assert.Equal(t, "menu-123123", result.ID)
	assert.Equal(t, "Dashboard", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestMenuRepo_FindByName(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockMenuRepo)

	expectedMenu := &permission.Menu{
		ID:   "menu-123123",
		Name: "Dashboard",
		Type: 1,
	}

	mockRepo.On("FindByName", ctx, "Dashboard").Return(expectedMenu, nil)

	result, err := mockRepo.FindByName(ctx, "Dashboard")

	assert.NoError(t, err)
	assert.Equal(t, "menu-123123", result.ID)
	mockRepo.AssertExpectations(t)
}

func TestMenuRepo_List(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockMenuRepo)

	menus := []*permission.Menu{
		{ID: "menu-1", Name: "Dashboard"},
		{ID: "menu-2", Name: "Settings"},
	}

	mockRepo.On("List", ctx).Return(menus, nil)

	result, err := mockRepo.List(ctx)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestMenuRepo_FindByParentID(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockMenuRepo)

	menus := []*permission.Menu{
		{ID: "menu-1", Name: "Users", ParentID: "menu-0"},
		{ID: "menu-2", Name: "Roles", ParentID: "menu-0"},
	}

	mockRepo.On("FindByParentID", ctx, "menu-0").Return(menus, nil)

	result, err := mockRepo.FindByParentID(ctx, "menu-0")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestMenuRepo_Update(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockMenuRepo)

	menu := &permission.Menu{
		ID:   "menu-123123",
		Name: "Dashboard Updated",
		Type: 1,
	}

	mockRepo.On("Update", ctx, mock.AnythingOfType("*permission.Menu")).Return(nil)

	err := mockRepo.Update(ctx, menu)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMenuRepo_Delete(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockMenuRepo)

	mockRepo.On("Delete", ctx, "menu-123123").Return(nil)

	err := mockRepo.Delete(ctx, "menu-123123")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMenuRepo_FindByMenuIDs(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockMenuRepo)

	menus := []*permission.Menu{
		{ID: "menu-1", Name: "Dashboard"},
		{ID: "menu-2", Name: "Settings"},
	}

	mockRepo.On("FindByMenuIDs", ctx, []string{"menu-1", "menu-2"}).Return(menus, nil)

	result, err := mockRepo.FindByMenuIDs(ctx, []string{"menu-1", "menu-2"})

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}
