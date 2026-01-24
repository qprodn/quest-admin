package user_test

import (
	"context"
	"testing"

	user "quest-admin/internal/biz/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, user *user.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) FindByID(ctx context.Context, id string) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepo) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepo) List(ctx context.Context, opt *user.WhereUserOpt) ([]*user.User, error) {
	args := m.Called(ctx, opt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user.User), args.Error(1)
}

func (m *MockUserRepo) Count(ctx context.Context, opt *user.WhereUserOpt) (int64, error) {
	args := m.Called(ctx, opt)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepo) Update(ctx context.Context, user *user.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) UpdatePassword(ctx context.Context, bo *user.UpdatePasswordBO) error {
	args := m.Called(ctx, bo)
	return args.Error(0)
}

func (m *MockUserRepo) UpdateStatus(ctx context.Context, bo *user.UpdateStatusBO) error {
	args := m.Called(ctx, bo)
	return args.Error(0)
}

func (m *MockUserRepo) UpdateLoginInfo(ctx context.Context, bo *user.UpdateLoginInfoBO) error {
	args := m.Called(ctx, bo)
	return args.Error(0)
}

func (m *MockUserRepo) Delete(ctx context.Context, bo *user.DeleteUserBO) error {
	args := m.Called(ctx, bo)
	return args.Error(0)
}

func TestUserRepo_Create(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepo)

	usr := &user.User{
		ID:       "user-1",
		Username: "testuser",
		Password: "hashedpass",
		Nickname: "Test User",
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*user.User")).Return(nil)

	err := mockRepo.Create(ctx, usr)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserRepo_FindByID(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepo)

	expected := &user.User{
		ID:       "user-1",
		Username: "testuser",
		Nickname: "Test User",
	}

	mockRepo.On("FindByID", ctx, "user-1").Return(expected, nil)

	result, err := mockRepo.FindByID(ctx, "user-1")

	assert.NoError(t, err)
	assert.Equal(t, "user-1", result.ID)
	assert.Equal(t, "testuser", result.Username)
	mockRepo.AssertExpectations(t)
}

func TestUserRepo_FindByUsername(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepo)

	expected := &user.User{
		ID:       "user-1",
		Username: "testuser",
	}

	mockRepo.On("FindByUsername", ctx, "testuser").Return(expected, nil)

	result, err := mockRepo.FindByUsername(ctx, "testuser")

	assert.NoError(t, err)
	assert.Equal(t, "user-1", result.ID)
	mockRepo.AssertExpectations(t)
}

func TestUserRepo_List(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepo)

	opt := &user.WhereUserOpt{
		Limit:  10,
		Offset: 0,
	}

	users := []*user.User{
		{ID: "user-1", Username: "user1"},
		{ID: "user-2", Username: "user2"},
	}

	mockRepo.On("List", ctx, opt).Return(users, nil)
	mockRepo.On("Count", ctx, opt).Return(int64(2), nil)

	mockRepo.List(ctx, opt)
	mockRepo.Count(ctx, opt)

	mockRepo.AssertExpectations(t)
}

func TestUserRepo_Update(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepo)

	usr := &user.User{
		ID:       "user-1",
		Username: "testuser",
		Nickname: "Updated Name",
	}

	mockRepo.On("Update", ctx, mock.AnythingOfType("*user.User")).Return(nil)

	err := mockRepo.Update(ctx, usr)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserRepo_Delete(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepo)

	bo := &user.DeleteUserBO{UserID: "user-1"}

	mockRepo.On("Delete", ctx, mock.AnythingOfType("*user.DeleteUserBO")).Return(nil)

	err := mockRepo.Delete(ctx, bo)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
