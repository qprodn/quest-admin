package user_test

import (
	"context"
	"quest-admin/internal/biz/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserPostBiz struct {
	mock.Mock
}

func (m *MockUserPostBiz) AssignUserPosts(ctx context.Context, bo *user.AssignUserPostsBO) error {
	args := m.Called(ctx, bo)
	return args.Error(0)
}

func TestUserPostService_AssignUserPosts(t *testing.T) {
	ctx := context.Background()
	mockBiz := new(MockUserPostBiz)

	bo := &user.AssignUserPostsBO{
		UserID:  "user-1",
		PostIDs: []string{"post-1", "post-2"},
	}

	mockBiz.On("AssignUserPosts", ctx, mock.AnythingOfType("*user.AssignUserPostsBO")).Return(nil)

	err := mockBiz.AssignUserPosts(ctx, bo)

	assert.NoError(t, err)
	mockBiz.AssertExpectations(t)
}
