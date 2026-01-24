package organization_test

import (
	"context"
	"testing"

	org "quest-admin/internal/biz/organization"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPostRepo struct {
	mock.Mock
}

func (m *MockPostRepo) Create(ctx context.Context, post *org.Post) (*org.Post, error) {
	args := m.Called(ctx, post)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*org.Post), args.Error(1)
}

func (m *MockPostRepo) FindByID(ctx context.Context, id string) (*org.Post, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*org.Post), args.Error(1)
}

func (m *MockPostRepo) FindByName(ctx context.Context, name string) (*org.Post, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*org.Post), args.Error(1)
}

func (m *MockPostRepo) FindByCode(ctx context.Context, code string) (*org.Post, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*org.Post), args.Error(1)
}

func (m *MockPostRepo) List(ctx context.Context, query *org.ListPostsQuery) (*org.ListPostsResult, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*org.ListPostsResult), args.Error(1)
}

func (m *MockPostRepo) Update(ctx context.Context, post *org.Post) (*org.Post, error) {
	args := m.Called(ctx, post)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*org.Post), args.Error(1)
}

func (m *MockPostRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPostRepo) HasUsers(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockPostRepo) FindListByIDs(ctx context.Context, ids []string) ([]*org.Post, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*org.Post), args.Error(1)
}

func TestPostRepo_Create(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPostRepo)

	post := &org.Post{
		ID:   "post-1",
		Name: "Engineer",
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*organization.Post")).Return(post, nil)

	result, err := mockRepo.Create(ctx, post)

	assert.NoError(t, err)
	assert.Equal(t, "post-1", result.ID)
	mockRepo.AssertExpectations(t)
}

func TestPostRepo_FindByID(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPostRepo)

	expected := &org.Post{
		ID:   "post-1",
		Name: "Engineer",
	}

	mockRepo.On("FindByID", ctx, "post-1").Return(expected, nil)

	result, err := mockRepo.FindByID(ctx, "post-1")

	assert.NoError(t, err)
	assert.Equal(t, "post-1", result.ID)
	mockRepo.AssertExpectations(t)
}

func TestPostRepo_HasUsers(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockPostRepo)

	mockRepo.On("HasUsers", ctx, "post-1").Return(true, nil)

	hasUsers, err := mockRepo.HasUsers(ctx, "post-1")

	assert.NoError(t, err)
	assert.True(t, hasUsers)
	mockRepo.AssertExpectations(t)
}
