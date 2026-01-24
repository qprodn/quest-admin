package tenant_test

import (
	"context"
	"testing"

	tenant "quest-admin/internal/biz/tenant"

	"github.com/go-kratos/kratos/v2/log"
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
		ID:   "tenant-1",
		Name: "Company A",
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*tenant.Tenant")).Return(nil)

	err := mockRepo.Create(ctx, tn)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTenantRepo_Create_Error(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepo)

	tn := &tenant.Tenant{
		ID:   "tenant-1",
		Name: "Company A",
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*tenant.Tenant")).Return(assert.AnError)

	err := mockRepo.Create(ctx, tn)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTenantRepo_FindByID(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepo)

	expected := &tenant.Tenant{
		ID:   "tenant-1",
		Name: "Company A",
	}

	mockRepo.On("FindByID", ctx, "tenant-1").Return(expected, nil)

	result, err := mockRepo.FindByID(ctx, "tenant-1")

	assert.NoError(t, err)
	assert.Equal(t, "tenant-1", result.ID)
	mockRepo.AssertExpectations(t)
}

func TestTenantRepo_FindByID_Error(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepo)

	mockRepo.On("FindByID", ctx, "tenant-1").Return(nil, assert.AnError)

	result, err := mockRepo.FindByID(ctx, "tenant-1")

	assert.Error(t, err)
	assert.Nil(t, result)
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

func TestTenantRepo_FindByName_Error(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepo)

	mockRepo.On("FindByName", ctx, "Company A").Return(nil, assert.AnError)

	result, err := mockRepo.FindByName(ctx, "Company A")

	assert.Error(t, err)
	assert.Nil(t, result)
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
		},
		Total:      1,
		Page:       1,
		PageSize:   10,
		TotalPages: 1,
	}

	mockRepo.On("List", ctx, query).Return(result, nil)

	resp, err := mockRepo.List(ctx, query)

	assert.NoError(t, err)
	assert.Len(t, resp.Tenants, 1)
	mockRepo.AssertExpectations(t)
}

func TestTenantRepo_List_Error(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepo)

	query := &tenant.ListTenantsQuery{
		Page:     1,
		PageSize: 10,
	}

	mockRepo.On("List", ctx, query).Return(nil, assert.AnError)

	resp, err := mockRepo.List(ctx, query)

	assert.Error(t, err)
	assert.Nil(t, resp)
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

func TestTenantRepo_Update_Error(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepo)

	tn := &tenant.Tenant{
		ID:   "tenant-1",
		Name: "Company A Updated",
	}

	mockRepo.On("Update", ctx, mock.AnythingOfType("*tenant.Tenant")).Return(assert.AnError)

	err := mockRepo.Update(ctx, tn)

	assert.Error(t, err)
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

func TestTenantRepo_Delete_Error(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepo)

	mockRepo.On("Delete", ctx, "tenant-1").Return(assert.AnError)

	err := mockRepo.Delete(ctx, "tenant-1")

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTenantUsecase_CreateTenant(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepo)
	logger := log.DefaultLogger

	uc := tenant.NewTenantUsecase(mockRepo, logger)

	tn := &tenant.Tenant{
		Name: "Company A",
	}

	mockRepo.On("FindByName", ctx, "Company A").Return(nil, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*tenant.Tenant")).Return(nil)

	err := uc.CreateTenant(ctx, tn)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTenantUsecase_CreateTenant_FindByNameError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepo)
	logger := log.DefaultLogger

	uc := tenant.NewTenantUsecase(mockRepo, logger)

	tn := &tenant.Tenant{
		Name: "Company A",
	}

	mockRepo.On("FindByName", ctx, "Company A").Return(nil, assert.AnError)

	err := uc.CreateTenant(ctx, tn)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTenantUsecase_CreateTenant_DuplicateName(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepo)
	logger := log.DefaultLogger

	uc := tenant.NewTenantUsecase(mockRepo, logger)

	tn := &tenant.Tenant{
		Name: "Company A",
	}

	existing := &tenant.Tenant{ID: "existing-1", Name: "Company A"}
	mockRepo.On("FindByName", ctx, "Company A").Return(existing, nil)

	err := uc.CreateTenant(ctx, tn)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTenantUsecase_CreateTenant_CreateError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepo)
	logger := log.DefaultLogger

	uc := tenant.NewTenantUsecase(mockRepo, logger)

	tn := &tenant.Tenant{
		Name: "Company A",
	}

	mockRepo.On("FindByName", ctx, "Company A").Return(nil, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*tenant.Tenant")).Return(assert.AnError)

	err := uc.CreateTenant(ctx, tn)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTenantUsecase_GetTenant(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepo)
	logger := log.DefaultLogger

	uc := tenant.NewTenantUsecase(mockRepo, logger)

	expected := &tenant.Tenant{
		ID:   "tenant-1",
		Name: "Company A",
	}

	mockRepo.On("FindByID", ctx, "tenant-1").Return(expected, nil)

	result, err := uc.GetTenant(ctx, "tenant-1")

	assert.NoError(t, err)
	assert.Equal(t, "tenant-1", result.ID)
	mockRepo.AssertExpectations(t)
}

func TestTenantUsecase_GetTenant_Error(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepo)
	logger := log.DefaultLogger

	uc := tenant.NewTenantUsecase(mockRepo, logger)

	mockRepo.On("FindByID", ctx, "tenant-1").Return(nil, assert.AnError)

	result, err := uc.GetTenant(ctx, "tenant-1")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestTenantUsecase_UpdateTenant(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepo)
	logger := log.DefaultLogger

	uc := tenant.NewTenantUsecase(mockRepo, logger)

	tn := &tenant.Tenant{
		ID:   "tenant-1",
		Name: "Company A Updated",
	}

	mockRepo.On("FindByID", ctx, "tenant-1").Return(tn, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*tenant.Tenant")).Return(nil)

	err := uc.UpdateTenant(ctx, tn)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTenantUsecase_UpdateTenant_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepo)
	logger := log.DefaultLogger

	uc := tenant.NewTenantUsecase(mockRepo, logger)

	tn := &tenant.Tenant{
		ID:   "tenant-1",
		Name: "Company A Updated",
	}

	mockRepo.On("FindByID", ctx, "tenant-1").Return(nil, assert.AnError)

	err := uc.UpdateTenant(ctx, tn)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTenantUsecase_UpdateTenant_UpdateError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepo)
	logger := log.DefaultLogger

	uc := tenant.NewTenantUsecase(mockRepo, logger)

	tn := &tenant.Tenant{
		ID:   "tenant-1",
		Name: "Company A Updated",
	}

	mockRepo.On("FindByID", ctx, "tenant-1").Return(tn, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*tenant.Tenant")).Return(assert.AnError)

	err := uc.UpdateTenant(ctx, tn)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTenantUsecase_DeleteTenant(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepo)
	logger := log.DefaultLogger

	uc := tenant.NewTenantUsecase(mockRepo, logger)

	tn := &tenant.Tenant{
		ID:   "tenant-1",
		Name: "Company A",
	}

	mockRepo.On("FindByID", ctx, "tenant-1").Return(tn, nil)
	mockRepo.On("Delete", ctx, "tenant-1").Return(nil)

	err := uc.DeleteTenant(ctx, "tenant-1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTenantUsecase_DeleteTenant_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepo)
	logger := log.DefaultLogger

	uc := tenant.NewTenantUsecase(mockRepo, logger)

	mockRepo.On("FindByID", ctx, "tenant-1").Return(nil, assert.AnError)

	err := uc.DeleteTenant(ctx, "tenant-1")

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTenantUsecase_DeleteTenant_DeleteError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockTenantRepo)
	logger := log.DefaultLogger

	uc := tenant.NewTenantUsecase(mockRepo, logger)

	tn := &tenant.Tenant{
		ID:   "tenant-1",
		Name: "Company A",
	}

	mockRepo.On("FindByID", ctx, "tenant-1").Return(tn, nil)
	mockRepo.On("Delete", ctx, "tenant-1").Return(assert.AnError)

	err := uc.DeleteTenant(ctx, "tenant-1")

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
