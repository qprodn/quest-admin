package tenant_test

import (
	"context"
	"quest-admin/internal/biz/tenant"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTenantBiz struct {
	mock.Mock
}

func (m *MockTenantBiz) FindByID(ctx context.Context, id string) (*tenant.Tenant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tenant.Tenant), args.Error(1)
}

func TestTenantService_FindByID(t *testing.T) {
	ctx := context.Background()
	mockBiz := new(MockTenantBiz)

	mockBiz.On("FindByID", ctx, "tenant-1").Return(&tenant.Tenant{
		ID:   "tenant-1",
		Name: "Company A",
	}, nil)

	tn, err := mockBiz.FindByID(ctx, "tenant-1")

	assert.NoError(t, err)
	assert.Equal(t, "tenant-1", tn.ID)
	mockBiz.AssertExpectations(t)
}
