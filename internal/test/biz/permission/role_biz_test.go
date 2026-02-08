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

func (m *MockRoleRepo) List(ctx context.Context, opt *permission.WhereRoleOpt) ([]*permission.Role, error) {
	args := m.Called(ctx, opt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*permission.Role), args.Error(1)
}

func (m *MockRoleRepo) Count(ctx context.Context, opt *permission.WhereRoleOpt) (int64, error) {
	args := m.Called(ctx, opt)
	return args.Get(0).(int64), args.Error(1)
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

func newTestRoleUsecase(t *testing.T) (*permission.RoleUsecase, *MockRoleRepo, *MockRoleMenuRepo) {
	mockRepo := new(MockRoleRepo)
	mockRoleMenuRepo := new(MockRoleMenuRepo)
	mockTm := new(MockTransactionManager)
	idg := idgen.NewIDGenerator()
	logger := log.DefaultLogger

	uc := permission.NewRoleUsecase(mockTm, idg, mockRepo, mockRoleMenuRepo, logger)
	return uc, mockRepo, mockRoleMenuRepo
}

func TestRoleUsecase_CreateRole(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name        string
		setupMock   func(*MockRoleRepo, *MockRoleMenuRepo)
		inputRole   *permission.Role
		expectError bool
	}{
		{
			name: "success",
			setupMock: func(m *MockRoleRepo, r *MockRoleMenuRepo) {
				role := &permission.Role{
					Name:   "Admin",
					Code:   "admin",
					Status: 1,
				}
				m.On("FindByName", ctx, "Admin").Return(nil, nil)
				m.On("FindByCode", ctx, "admin").Return(nil, nil)
				m.On("Create", ctx, mock.AnythingOfType("*permission.Role")).Return(role, nil)
			},
			inputRole: &permission.Role{
				Name:   "Admin",
				Code:   "admin",
				Status: 1,
			},
			expectError: false,
		},
		{
			name: "repo error when finding by name",
			setupMock: func(m *MockRoleRepo, r *MockRoleMenuRepo) {
				m.On("FindByName", ctx, "Admin").Return(nil, assert.AnError)
			},
			inputRole: &permission.Role{
				Name:   "Admin",
				Code:   "admin",
				Status: 1,
			},
			expectError: true,
		},
		{
			name: "duplicate name",
			setupMock: func(m *MockRoleRepo, r *MockRoleMenuRepo) {
				existing := &permission.Role{ID: "role-1", Name: "Admin"}
				m.On("FindByName", ctx, "Admin").Return(existing, nil)
			},
			inputRole: &permission.Role{
				Name: "Admin",
				Code: "admin",
			},
			expectError: true,
		},
		{
			name: "duplicate code",
			setupMock: func(m *MockRoleRepo, r *MockRoleMenuRepo) {
				m.On("FindByName", ctx, "Admin").Return(nil, nil)
				existing := &permission.Role{ID: "role-1", Code: "admin"}
				m.On("FindByCode", ctx, "admin").Return(existing, nil)
			},
			inputRole: &permission.Role{
				Name:   "Admin",
				Code:   "admin",
				Status: 1,
			},
			expectError: true,
		},
		{
			name: "create error",
			setupMock: func(m *MockRoleRepo, r *MockRoleMenuRepo) {
				m.On("FindByName", ctx, "Admin").Return(nil, nil)
				m.On("FindByCode", ctx, "admin").Return(nil, nil)
				m.On("Create", ctx, mock.AnythingOfType("*permission.Role")).Return(nil, assert.AnError)
			},
			inputRole: &permission.Role{
				Name:   "Admin",
				Code:   "admin",
				Status: 1,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, mockRepo, mockRoleMenuRepo := newTestRoleUsecase(t)
			tt.setupMock(mockRepo, mockRoleMenuRepo)

			result, err := uc.CreateRole(ctx, tt.inputRole)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestRoleUsecase_GetRole(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name        string
		setupMock   func(*MockRoleRepo, *MockRoleMenuRepo)
		inputID     string
		expectError bool
		expectID    string
	}{
		{
			name: "success",
			setupMock: func(m *MockRoleRepo, r *MockRoleMenuRepo) {
				expected := &permission.Role{
					ID:   "role-1",
					Name: "Admin",
				}
				m.On("FindByID", ctx, "role-1").Return(expected, nil)
			},
			inputID:     "role-1",
			expectError: false,
			expectID:    "role-1",
		},
		{
			name: "not found",
			setupMock: func(m *MockRoleRepo, r *MockRoleMenuRepo) {
				m.On("FindByID", ctx, "role-1").Return(nil, assert.AnError)
			},
			inputID:     "role-1",
			expectError: true,
			expectID:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, mockRepo, mockRoleMenuRepo := newTestRoleUsecase(t)
			tt.setupMock(mockRepo, mockRoleMenuRepo)

			result, err := uc.GetRole(ctx, tt.inputID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestRoleUsecase_UpdateRole(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name        string
		setupMock   func(*MockRoleRepo, *MockRoleMenuRepo)
		inputRole   *permission.Role
		expectError bool
	}{
		{
			name: "success",
			setupMock: func(m *MockRoleRepo, r *MockRoleMenuRepo) {
				role := &permission.Role{
					ID:     "role-1",
					Name:   "Admin Updated",
					Status: 1,
				}
				m.On("FindByID", ctx, "role-1").Return(role, nil)
				m.On("Update", ctx, mock.AnythingOfType("*permission.Role")).Return(role, nil)
			},
			inputRole: &permission.Role{
				ID:     "role-1",
				Name:   "Admin Updated",
				Status: 1,
			},
			expectError: false,
		},
		{
			name: "not found",
			setupMock: func(m *MockRoleRepo, r *MockRoleMenuRepo) {
				m.On("FindByID", ctx, "role-1").Return(nil, assert.AnError)
			},
			inputRole: &permission.Role{
				ID:     "role-1",
				Name:   "Admin Updated",
				Status: 1,
			},
			expectError: true,
		},
		{
			name: "update error",
			setupMock: func(m *MockRoleRepo, r *MockRoleMenuRepo) {
				role := &permission.Role{
					ID:     "role-1",
					Name:   "Admin Updated",
					Status: 1,
				}
				m.On("FindByID", ctx, "role-1").Return(role, nil)
				m.On("Update", ctx, mock.AnythingOfType("*permission.Role")).Return(nil, assert.AnError)
			},
			inputRole: &permission.Role{
				ID:     "role-1",
				Name:   "Admin Updated",
				Status: 1,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, mockRepo, mockRoleMenuRepo := newTestRoleUsecase(t)
			tt.setupMock(mockRepo, mockRoleMenuRepo)

			result, err := uc.UpdateRole(ctx, tt.inputRole)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestRoleUsecase_DeleteRole(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name        string
		setupMock   func(*MockRoleRepo, *MockRoleMenuRepo)
		inputID     string
		expectError bool
	}{
		{
			name: "success",
			setupMock: func(m *MockRoleRepo, r *MockRoleMenuRepo) {
				role := &permission.Role{ID: "role-1"}
				m.On("FindByID", ctx, "role-1").Return(role, nil)
				m.On("HasUsers", ctx, "role-1").Return(false, nil)
				m.On("Delete", ctx, "role-1").Return(nil)
			},
			inputID:     "role-1",
			expectError: false,
		},
		{
			name: "not found",
			setupMock: func(m *MockRoleRepo, r *MockRoleMenuRepo) {
				m.On("FindByID", ctx, "role-1").Return(nil, assert.AnError)
			},
			inputID:     "role-1",
			expectError: true,
		},
		{
			name: "has users",
			setupMock: func(m *MockRoleRepo, r *MockRoleMenuRepo) {
				role := &permission.Role{ID: "role-1"}
				m.On("FindByID", ctx, "role-1").Return(role, nil)
				m.On("HasUsers", ctx, "role-1").Return(true, nil)
			},
			inputID:     "role-1",
			expectError: true,
		},
		{
			name: "has users error",
			setupMock: func(m *MockRoleRepo, r *MockRoleMenuRepo) {
				role := &permission.Role{ID: "role-1"}
				m.On("FindByID", ctx, "role-1").Return(role, nil)
				m.On("HasUsers", ctx, "role-1").Return(false, assert.AnError)
			},
			inputID:     "role-1",
			expectError: true,
		},
		{
			name: "delete error",
			setupMock: func(m *MockRoleRepo, r *MockRoleMenuRepo) {
				role := &permission.Role{ID: "role-1"}
				m.On("FindByID", ctx, "role-1").Return(role, nil)
				m.On("HasUsers", ctx, "role-1").Return(false, nil)
				m.On("Delete", ctx, "role-1").Return(assert.AnError)
			},
			inputID:     "role-1",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, mockRepo, mockRoleMenuRepo := newTestRoleUsecase(t)
			tt.setupMock(mockRepo, mockRoleMenuRepo)

			err := uc.DeleteRole(ctx, tt.inputID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestRoleUsecase_GetRoleMenus(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name         string
		setupMock    func(*MockRoleRepo, *MockRoleMenuRepo)
		inputRoleID  string
		expectError  bool
		expectLength int
	}{
		{
			name: "success",
			setupMock: func(m *MockRoleRepo, r *MockRoleMenuRepo) {
				role := &permission.Role{ID: "role-1"}
				menuIDs := []string{"menu-1", "menu-2"}
				m.On("FindByID", ctx, "role-1").Return(role, nil)
				r.On("GetMenuIDs", ctx, "role-1").Return(menuIDs, nil)
			},
			inputRoleID:  "role-1",
			expectError:  false,
			expectLength: 2,
		},
		{
			name: "not found",
			setupMock: func(m *MockRoleRepo, r *MockRoleMenuRepo) {
				m.On("FindByID", ctx, "role-1").Return(nil, assert.AnError)
			},
			inputRoleID:  "role-1",
			expectError:  true,
			expectLength: 0,
		},
		{
			name: "get menu IDs error",
			setupMock: func(m *MockRoleRepo, r *MockRoleMenuRepo) {
				role := &permission.Role{ID: "role-1"}
				m.On("FindByID", ctx, "role-1").Return(role, nil)
				r.On("GetMenuIDs", ctx, "role-1").Return(nil, assert.AnError)
			},
			inputRoleID:  "role-1",
			expectError:  true,
			expectLength: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, mockRepo, mockRoleMenuRepo := newTestRoleUsecase(t)
			tt.setupMock(mockRepo, mockRoleMenuRepo)

			result, err := uc.GetRoleMenus(ctx, tt.inputRoleID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectLength)
			}

			mockRepo.AssertExpectations(t)
			mockRoleMenuRepo.AssertExpectations(t)
		})
	}
}

func TestRoleUsecase_ListByRoleIDs(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name         string
		setupMock    func(*MockRoleRepo, *MockRoleMenuRepo)
		inputRoleIDs []string
		expectError  bool
		expectLength int
	}{
		{
			name: "success",
			setupMock: func(m *MockRoleRepo, r *MockRoleMenuRepo) {
				roles := []*permission.Role{
					{ID: "role-1", Name: "Admin"},
					{ID: "role-2", Name: "User"},
				}
				m.On("FindListByIDs", ctx, []string{"role-1", "role-2"}).Return(roles, nil)
			},
			inputRoleIDs: []string{"role-1", "role-2"},
			expectError:  false,
			expectLength: 2,
		},
		{
			name: "error",
			setupMock: func(m *MockRoleRepo, r *MockRoleMenuRepo) {
				m.On("FindListByIDs", ctx, []string{"role-1", "role-2"}).Return(nil, assert.AnError)
			},
			inputRoleIDs: []string{"role-1", "role-2"},
			expectError:  true,
			expectLength: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, mockRepo, mockRoleMenuRepo := newTestRoleUsecase(t)
			tt.setupMock(mockRepo, mockRoleMenuRepo)

			result, err := uc.ListByRoleIDs(ctx, tt.inputRoleIDs)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectLength)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
