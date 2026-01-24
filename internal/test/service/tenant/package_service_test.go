package tenant_test

import (
	"context"
	"quest-admin/internal/biz/tenant"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPackageBiz struct {
	mock.Mock
}

func (m *MockPackageBiz) FindByID(ctx context.Context, id string) (*tenant.TenantPackage, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*tenant.TenantPackage), args.Error(1)
}

func TestPackageService_FindByID(t *testing.T) {
	ctx := context.Background()
	mockBiz := new(MockPackageBiz)

	mockBiz.On("FindByID", ctx, "pkg-1").Return(&tenant.TenantPackage{
		ID:   "pkg-1",
		Name: "Basic",
	}, nil)

	pkg, err := mockBiz.FindByID(ctx, "pkg-1")

	assert.NoError(t, err)
	assert.Equal(t, "pkg-1", pkg.ID)
	mockBiz.AssertExpectations(t)
}
