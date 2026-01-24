package permission_test

import (
	"context"
	"quest-admin/internal/biz/permission"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMenuBiz struct {
	mock.Mock
}

func (m *MockMenuBiz) FindByID(ctx context.Context, id string) (*permission.Menu, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*permission.Menu), args.Error(1)
}

func TestMenuService_FindByID(t *testing.T) {
	ctx := context.Background()
	mockBiz := new(MockMenuBiz)

	mockBiz.On("FindByID", ctx, "menu-1").Return(&permission.Menu{
		ID:   "menu-1",
		Name: "Dashboard",
	}, nil)

	menu, err := mockBiz.FindByID(ctx, "menu-1")

	assert.NoError(t, err)
	assert.Equal(t, "menu-1", menu.ID)
	mockBiz.AssertExpectations(t)
}
