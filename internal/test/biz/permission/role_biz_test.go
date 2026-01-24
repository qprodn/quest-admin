package permission_test

import (
	"context"
	"testing"

	permission "quest-admin/internal/biz/permission"
	"quest-admin/internal/data/idgen"

	"github.com/go-kratos/kratos/v2/log"
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

type MockRoleMenuRepo struct {
	mock.Mock
}

func (m *MockRoleMenuRepo) Create(ctx context.Context, item *permission.RoleMenu) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockRoleMenuRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRoleMenuRepo) GetRoleMenus(ctx context.Context, roleID string) ([]*permission.RoleMenu, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*permission.RoleMenu), args.Error(1)
}

func (m *MockRoleMenuRepo) GetMenuIDs(ctx context.Context, roleID string) ([]string, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRoleMenuRepo) FindListByRoleIDs(ctx context.Context, roles []string) ([]*permission.RoleMenu, error) {
	args := m.Called(ctx, roles)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*permission.RoleMenu), args.Error(1)
}

type MockTransactionManager struct {
	mock.Mock
}

func (m *MockTransactionManager) Tx(ctx context.Context, fn func(context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

func TestRoleUsecase_CreateRole(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRoleRepo)
	mockRoleMenuRepo := new(MockRoleMenuRepo)
	mockTm := new(MockTransactionManager)
	idg := idgen.NewIDGenerator()
	logger := log.DefaultLogger

	uc := permission.NewRoleUsecase(mockTm, idg, mockRepo, mockRoleMenuRepo, logger)

	role := &permission.Role{
		Name:   "Admin",
		Code:   "admin",
		Status: 1,
	}

	mockRepo.On("FindByName", ctx, "Admin").Return(nil, nil)
	mockRepo.On("FindByCode", ctx, "admin").Return(nil, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*permission.Role")).Return(role, nil)

	result, err := uc.CreateRole(ctx, role)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestRoleUsecase_CreateRole_DuplicateName(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRoleRepo)
	mockRoleMenuRepo := new(MockRoleMenuRepo)
	mockTm := new(MockTransactionManager)
	idg := idgen.NewIDGenerator()
	logger := log.DefaultLogger

	uc := permission.NewRoleUsecase(mockTm, idg, mockRepo, mockRoleMenuRepo, logger)

	role := &permission.Role{
		Name: "Admin",
		Code: "admin",
	}

	existing := &permission.Role{ID: "role-1", Name: "Admin"}
	mockRepo.On("FindByName", ctx, "Admin").Return(existing, nil)

	_, err := uc.CreateRole(ctx, role)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRoleUsecase_GetRole(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRoleRepo)
	mockRoleMenuRepo := new(MockRoleMenuRepo)
	mockTm := new(MockTransactionManager)
	idg := idgen.NewIDGenerator()
	logger := log.DefaultLogger

	uc := permission.NewRoleUsecase(mockTm, idg, mockRepo, mockRoleMenuRepo, logger)

	expected := &permission.Role{
		ID:   "role-1",
		Name: "Admin",
	}

	mockRepo.On("FindByID", ctx, "role-1").Return(expected, nil)

	result, err := uc.GetRole(ctx, "role-1")

	assert.NoError(t, err)
	assert.Equal(t, "role-1", result.ID)
	mockRepo.AssertExpectations(t)
}

func TestRoleUsecase_UpdateRole(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRoleRepo)
	mockRoleMenuRepo := new(MockRoleMenuRepo)
	mockTm := new(MockTransactionManager)
	idg := idgen.NewIDGenerator()
	logger := log.DefaultLogger

	uc := permission.NewRoleUsecase(mockTm, idg, mockRepo, mockRoleMenuRepo, logger)

	role := &permission.Role{
		ID:     "role-1",
		Name:   "Admin Updated",
		Status: 1,
	}

	mockRepo.On("FindByID", ctx, "role-1").Return(role, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*permission.Role")).Return(role, nil)

	result, err := uc.UpdateRole(ctx, role)

	assert.NoError(t, err)
	assert.Equal(t, "Admin Updated", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestRoleUsecase_DeleteRole(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRoleRepo)
	mockRoleMenuRepo := new(MockRoleMenuRepo)
	mockTm := new(MockTransactionManager)
	idg := idgen.NewIDGenerator()
	logger := log.DefaultLogger

	uc := permission.NewRoleUsecase(mockTm, idg, mockRepo, mockRoleMenuRepo, logger)

	role := &permission.Role{ID: "role-1"}

	mockRepo.On("FindByID", ctx, "role-1").Return(role, nil)
	mockRepo.On("HasUsers", ctx, "role-1").Return(false, nil)
	mockRepo.On("Delete", ctx, "role-1").Return(nil)

	err := uc.DeleteRole(ctx, "role-1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRoleUsecase_DeleteRole_HasUsers(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRoleRepo)
	mockRoleMenuRepo := new(MockRoleMenuRepo)
	mockTm := new(MockTransactionManager)
	idg := idgen.NewIDGenerator()
	logger := log.DefaultLogger

	uc := permission.NewRoleUsecase(mockTm, idg, mockRepo, mockRoleMenuRepo, logger)

	role := &permission.Role{ID: "role-1"}

	mockRepo.On("FindByID", ctx, "role-1").Return(role, nil)
	mockRepo.On("HasUsers", ctx, "role-1").Return(true, nil)

	err := uc.DeleteRole(ctx, "role-1")

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRoleUsecase_GetRoleMenus(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRoleRepo)
	mockRoleMenuRepo := new(MockRoleMenuRepo)
	mockTm := new(MockTransactionManager)
	idg := idgen.NewIDGenerator()
	logger := log.DefaultLogger

	uc := permission.NewRoleUsecase(mockTm, idg, mockRepo, mockRoleMenuRepo, logger)

	role := &permission.Role{ID: "role-1"}
	menuIDs := []string{"menu-1", "menu-2"}

	mockRepo.On("FindByID", ctx, "role-1").Return(role, nil)
	mockRoleMenuRepo.On("GetMenuIDs", ctx, "role-1").Return(menuIDs, nil)

	result, err := uc.GetRoleMenus(ctx, "role-1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestRoleUsecase_ListByRoleIDs(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRoleRepo)
	mockRoleMenuRepo := new(MockRoleMenuRepo)
	mockTm := new(MockTransactionManager)
	idg := idgen.NewIDGenerator()
	logger := log.DefaultLogger

	uc := permission.NewRoleUsecase(mockTm, idg, mockRepo, mockRoleMenuRepo, logger)

	roles := []*permission.Role{
		{ID: "role-1", Name: "Admin"},
		{ID: "role-2", Name: "User"},
	}

	mockRepo.On("FindListByIDs", ctx, []string{"role-1", "role-2"}).Return(roles, nil)

	result, err := uc.ListByRoleIDs(ctx, []string{"role-1", "role-2"})

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}
