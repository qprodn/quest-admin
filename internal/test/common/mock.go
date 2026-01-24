package common

import (
	"context"

	bizOrganization "quest-admin/internal/biz/organization"
	bizPermission "quest-admin/internal/biz/permission"
	bizTenant "quest-admin/internal/biz/tenant"
	bizUser "quest-admin/internal/biz/user"

	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Begin() (*MockTx, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockTx), args.Error(1)
}

type MockTx struct {
	mock.Mock
}

func (m *MockTx) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTx) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Logf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Log(level interface{}, format string, args ...interface{}) error {
	m.Called(level, format, args)
	return nil
}

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, user *bizUser.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) FindByID(ctx context.Context, id string) (*bizUser.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bizUser.User), args.Error(1)
}

func (m *MockUserRepo) FindByUsername(ctx context.Context, username string) (*bizUser.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bizUser.User), args.Error(1)
}

func (m *MockUserRepo) List(ctx context.Context, opt *bizUser.WhereUserOpt) ([]*bizUser.User, error) {
	args := m.Called(ctx, opt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bizUser.User), args.Error(1)
}

func (m *MockUserRepo) Count(ctx context.Context, opt *bizUser.WhereUserOpt) (int64, error) {
	args := m.Called(ctx, opt)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepo) Update(ctx context.Context, user *bizUser.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) UpdatePassword(ctx context.Context, bo *bizUser.UpdatePasswordBO) error {
	args := m.Called(ctx, bo)
	return args.Error(0)
}

func (m *MockUserRepo) UpdateStatus(ctx context.Context, bo *bizUser.UpdateStatusBO) error {
	args := m.Called(ctx, bo)
	return args.Error(0)
}

func (m *MockUserRepo) UpdateLoginInfo(ctx context.Context, bo *bizUser.UpdateLoginInfoBO) error {
	args := m.Called(ctx, bo)
	return args.Error(0)
}

func (m *MockUserRepo) Delete(ctx context.Context, bo *bizUser.DeleteUserBO) error {
	args := m.Called(ctx, bo)
	return args.Error(0)
}

type MockRoleRepo struct {
	mock.Mock
}

func (m *MockRoleRepo) Create(ctx context.Context, role *bizPermission.Role) (*bizPermission.Role, error) {
	args := m.Called(ctx, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bizPermission.Role), args.Error(1)
}

func (m *MockRoleRepo) FindByID(ctx context.Context, id string) (*bizPermission.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bizPermission.Role), args.Error(1)
}

func (m *MockRoleRepo) FindByName(ctx context.Context, name string) (*bizPermission.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bizPermission.Role), args.Error(1)
}

func (m *MockRoleRepo) FindByCode(ctx context.Context, code string) (*bizPermission.Role, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bizPermission.Role), args.Error(1)
}

func (m *MockRoleRepo) List(ctx context.Context, query *bizPermission.ListRolesQuery) (*bizPermission.ListRolesResult, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bizPermission.ListRolesResult), args.Error(1)
}

func (m *MockRoleRepo) Update(ctx context.Context, role *bizPermission.Role) (*bizPermission.Role, error) {
	args := m.Called(ctx, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bizPermission.Role), args.Error(1)
}

func (m *MockRoleRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRoleRepo) HasUsers(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockRoleRepo) FindListByIDs(ctx context.Context, roleIds []string) ([]*bizPermission.Role, error) {
	args := m.Called(ctx, roleIds)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bizPermission.Role), args.Error(1)
}

type MockDepartmentRepo struct {
	mock.Mock
}

func (m *MockDepartmentRepo) Create(ctx context.Context, dept *bizOrganization.Department) (*bizOrganization.Department, error) {
	args := m.Called(ctx, dept)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bizOrganization.Department), args.Error(1)
}

func (m *MockDepartmentRepo) FindByID(ctx context.Context, id string) (*bizOrganization.Department, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bizOrganization.Department), args.Error(1)
}

func (m *MockDepartmentRepo) FindByName(ctx context.Context, name string) (*bizOrganization.Department, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bizOrganization.Department), args.Error(1)
}

func (m *MockDepartmentRepo) List(ctx context.Context) ([]*bizOrganization.Department, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bizOrganization.Department), args.Error(1)
}

func (m *MockDepartmentRepo) FindByParentID(ctx context.Context, parentID string) ([]*bizOrganization.Department, error) {
	args := m.Called(ctx, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bizOrganization.Department), args.Error(1)
}

func (m *MockDepartmentRepo) Update(ctx context.Context, dept *bizOrganization.Department) (*bizOrganization.Department, error) {
	args := m.Called(ctx, dept)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bizOrganization.Department), args.Error(1)
}

func (m *MockDepartmentRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDepartmentRepo) HasUsers(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDepartmentRepo) FindListByIDs(ctx context.Context, ids []string) ([]*bizOrganization.Department, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bizOrganization.Department), args.Error(1)
}

type MockTenantRepo struct {
	mock.Mock
}

func (m *MockTenantRepo) Create(ctx context.Context, tenant *bizTenant.Tenant) error {
	args := m.Called(ctx, tenant)
	return args.Error(0)
}

func (m *MockTenantRepo) FindByID(ctx context.Context, id string) (*bizTenant.Tenant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bizTenant.Tenant), args.Error(1)
}

func (m *MockTenantRepo) FindByName(ctx context.Context, name string) (*bizTenant.Tenant, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bizTenant.Tenant), args.Error(1)
}

func (m *MockTenantRepo) List(ctx context.Context, query *bizTenant.ListTenantsQuery) (*bizTenant.ListTenantsResult, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bizTenant.ListTenantsResult), args.Error(1)
}

func (m *MockTenantRepo) Update(ctx context.Context, tenant *bizTenant.Tenant) error {
	args := m.Called(ctx, tenant)
	return args.Error(0)
}

func (m *MockTenantRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockIDGenerator struct {
	mock.Mock
}

func (m *MockIDGenerator) NextID(prefix string) string {
	args := m.Called(prefix)
	return args.String(0)
}

func (m *MockIDGenerator) NextOrderID() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

type MockTransactionManager struct {
	mock.Mock
}

func (m *MockTransactionManager) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}
