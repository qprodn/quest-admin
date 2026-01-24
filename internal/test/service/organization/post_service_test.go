package organization_test

import (
	"context"
	"quest-admin/internal/biz/organization"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPostBiz struct {
	mock.Mock
}

func (m *MockPostBiz) FindByID(ctx context.Context, id string) (*organization.Post, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*organization.Post), args.Error(1)
}

func TestPostService_FindByID(t *testing.T) {
	ctx := context.Background()
	mockBiz := new(MockPostBiz)

	mockBiz.On("FindByID", ctx, "post-1").Return(&organization.Post{
		ID:   "post-1",
		Name: "Engineer",
	}, nil)

	post, err := mockBiz.FindByID(ctx, "post-1")

	assert.NoError(t, err)
	assert.Equal(t, "post-1", post.ID)
	mockBiz.AssertExpectations(t)
}
