package tenant_test

import (
	"context"
	"testing"
	"time"

	tenant "quest-admin/internal/biz/tenant"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPackageRepo struct {
	mock.Mock
}

func (m *MockPackageRepo) Create(ctx context.Context, pkg *tenant.TenantPackage) (*tenant.TenantPackage, error) {
	args := m.Called(ctx, pkg)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tenant.TenantPackage), args.Error(1)
}

func (m *MockPackageRepo) FindByID(ctx context.Context, id string) (*tenant.TenantPackage, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tenant.TenantPackage), args.Error(1)
}

func (m *MockPackageRepo) FindByName(ctx context.Context, name string) (*tenant.TenantPackage, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tenant.TenantPackage), args.Error(1)
}

func (m *MockPackageRepo) List(ctx context.Context, query *tenant.ListPackagesQuery) (*tenant.ListPackagesResult, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tenant.ListPackagesResult), args.Error(1)
}

func (m *MockPackageRepo) Update(ctx context.Context, pkg *tenant.TenantPackage) (*tenant.TenantPackage, error) {
	args := m.Called(ctx, pkg)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tenant.TenantPackage), args.Error(1)
}

func (m *MockPackageRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPackageRepo) IsInUse(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func TestPackageRepo_Create(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPackageRepo)

	pkg := &tenant.TenantPackage{
		ID:      "pkg-1",
		Name:    "Basic",
		Status:  1,
		MenuIDs: "menu-1,menu-2",
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*tenant.TenantPackage")).Return(pkg, nil)

	result, err := mockRepo.Create(ctx, pkg)

	assert.NoError(t, err)
	assert.Equal(t, "pkg-1", result.ID)
	assert.Equal(t, "Basic", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestPackageRepo_FindByID(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPackageRepo)

	expected := &tenant.TenantPackage{
		ID:       "pkg-1",
		Name:     "Basic",
		Status:   1,
		MenuIDs:  "menu-1,menu-2",
		CreateAt: time.Now(),
	}

	mockRepo.On("FindByID", ctx, "pkg-1").Return(expected, nil)

	result, err := mockRepo.FindByID(ctx, "pkg-1")

	assert.NoError(t, err)
	assert.Equal(t, "pkg-1", result.ID)
	assert.Equal(t, "Basic", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestPackageRepo_FindByName(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPackageRepo)

	expected := &tenant.TenantPackage{
		ID:   "pkg-1",
		Name: "Basic",
	}

	mockRepo.On("FindByName", ctx, "Basic").Return(expected, nil)

	result, err := mockRepo.FindByName(ctx, "Basic")

	assert.NoError(t, err)
	assert.Equal(t, "pkg-1", result.ID)
	mockRepo.AssertExpectations(t)
}

func TestPackageRepo_List(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPackageRepo)

	query := &tenant.ListPackagesQuery{
		Page:     1,
		PageSize: 10,
	}

	result := &tenant.ListPackagesResult{
		Packages: []*tenant.TenantPackage{
			{ID: "pkg-1", Name: "Basic"},
			{ID: "pkg-2", Name: "Pro"},
		},
		Total:      2,
		Page:       1,
		PageSize:   10,
		TotalPages: 1,
	}

	mockRepo.On("List", ctx, query).Return(result, nil)

	resp, err := mockRepo.List(ctx, query)

	assert.NoError(t, err)
	assert.Len(t, resp.Packages, 2)
	mockRepo.AssertExpectations(t)
}

func TestPackageRepo_Update(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPackageRepo)

	pkg := &tenant.TenantPackage{
		ID:     "pkg-1",
		Name:   "Basic Updated",
		Status: 1,
	}

	mockRepo.On("Update", ctx, mock.AnythingOfType("*tenant.TenantPackage")).Return(pkg, nil)

	result, err := mockRepo.Update(ctx, pkg)

	assert.NoError(t, err)
	assert.Equal(t, "Basic Updated", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestPackageRepo_Delete(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPackageRepo)

	mockRepo.On("Delete", ctx, "pkg-1").Return(nil)

	err := mockRepo.Delete(ctx, "pkg-1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPackageRepo_IsInUse(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPackageRepo)

	mockRepo.On("IsInUse", ctx, "pkg-1").Return(true, nil)

	inUse, err := mockRepo.IsInUse(ctx, "pkg-1")

	assert.NoError(t, err)
	assert.True(t, inUse)
	mockRepo.AssertExpectations(t)
}
