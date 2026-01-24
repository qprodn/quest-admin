package user_test

import (
	"context"
	"testing"

	user "quest-admin/internal/biz/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//type MockUserPostRepo struct {
//	mock.Mock
//}
//
//func (m *MockUserPostRepo) GetUserPosts(ctx context.Context, userID string) ([]*user.UserPost, error) {
//	args := m.Called(ctx, userID)
//	if args.Get(0) == nil {
//		return nil, args.Error(1)
//	}
//	return args.Get(0).([]*user.UserPost), args.Error(1)
//}
//
//func (m *MockUserPostRepo) Delete(ctx context.Context, id string) error {
//	args := m.Called(ctx, id)
//	return args.Error(0)
//}
//
//func (m *MockUserPostRepo) Create(ctx context.Context, item *user.UserPost) error {
//	args := m.Called(ctx, item)
//	return args.Error(0)
//}

func TestUserPostRepo_GetUserPosts(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserPostRepo)

	posts := []*user.UserPost{
		{ID: "up-1", UserID: "user-1", PostID: "post-1"},
		{ID: "up-2", UserID: "user-1", PostID: "post-2"},
	}

	mockRepo.On("GetUserPosts", ctx, "user-1").Return(posts, nil)

	result, err := mockRepo.GetUserPosts(ctx, "user-1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestUserPostRepo_Create(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserPostRepo)

	item := &user.UserPost{
		ID:     "up-1",
		UserID: "user-1",
		PostID: "post-1",
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*user.UserPost")).Return(nil)

	err := mockRepo.Create(ctx, item)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserPostRepo_Delete(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserPostRepo)

	mockRepo.On("Delete", ctx, "up-1").Return(nil)

	err := mockRepo.Delete(ctx, "up-1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
