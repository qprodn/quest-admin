package tenant_test

import (
	"context"
	"testing"
	"time"

	tenant "quest-admin/internal/biz/tenant"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTenantRepo struct {
	mock.Mock
}

func (m *MockTenantRepo) Create(ctx context.Context, tenant *tenant.Tenant) error {
	args := m.Called(ctx, tenant)
	return args.Error(0)
}

func (m *MockTenantRepo) FindByID(ctx context.Context, id string) (*tenant.Tenant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tenant.Tenant), args.Error(1)
}

func (m *MockTenantRepo) FindByName(ctx context.Context, name string) (*tenant.Tenant, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tenant.Tenant), args.Error(1)
}

func (m *MockTenantRepo) List(ctx context.Context, query *tenant.ListTenantsQuery) (*tenant.ListTenantsResult, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tenant.ListTenantsResult), args.Error(1)
}

func (m *MockTenantRepo) Update(ctx context.Context, tenant *tenant.Tenant) error {
	args := m.Called(ctx, tenant)
	return args.Error(0)
}

func (m *MockTenantRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestTenantRepo_Create(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepo)

	tn := &tenant.Tenant{
		ID:          "tenant-1",
		Name:        "Company A",
		ContactName: "John Doe",
		Status:      1,
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*tenant.Tenant")).Return(nil)

	err := mockRepo.Create(ctx, tn)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTenantRepo_FindByID(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepo)

	expected := &tenant.Tenant{
		ID:          "tenant-1",
		Name:        "Company A",
		ContactName: "John Doe",
		CreateAt:    time.Now(),
	}

	mockRepo.On("FindByID", ctx, "tenant-1").Return(expected, nil)

	result, err := mockRepo.FindByID(ctx, "tenant-1")

	assert.NoError(t, err)
	assert.Equal(t, "tenant-1", result.ID)
	assert.Equal(t, "Company A", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestTenantRepo_FindByName(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepo)

	expected := &tenant.Tenant{
		ID:   "tenant-1",
		Name: "Company A",
	}

	mockRepo.On("FindByName", ctx, "Company A").Return(expected, nil)

	result, err := mockRepo.FindByName(ctx, "Company A")

	assert.NoError(t, err)
	assert.Equal(t, "tenant-1", result.ID)
	mockRepo.AssertExpectations(t)
}

func TestTenantRepo_List(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepo)

	query := &tenant.ListTenantsQuery{
		Page:     1,
		PageSize: 10,
	}

	result := &tenant.ListTenantsResult{
		Tenants: []*tenant.Tenant{
			{ID: "tenant-1", Name: "Company A"},
			{ID: "tenant-2", Name: "Company B"},
		},
		Total:      2,
		Page:       1,
		PageSize:   10,
		TotalPages: 1,
	}

	mockRepo.On("List", ctx, query).Return(result, nil)

	resp, err := mockRepo.List(ctx, query)

	assert.NoError(t, err)
	assert.Len(t, resp.Tenants, 2)
	mockRepo.AssertExpectations(t)
}

func TestTenantRepo_Update(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepo)

	tn := &tenant.Tenant{
		ID:   "tenant-1",
		Name: "Company A Updated",
	}

	mockRepo.On("Update", ctx, mock.AnythingOfType("*tenant.Tenant")).Return(nil)

	err := mockRepo.Update(ctx, tn)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTenantRepo_Delete(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepo)

	mockRepo.On("Delete", ctx, "tenant-1").Return(nil)

	err := mockRepo.Delete(ctx, "tenant-1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
