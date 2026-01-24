package user_test

import (
	"context"
	"testing"
	"time"

	user "quest-admin/internal/biz/user"
	"quest-admin/internal/data/idgen"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

var UserRepo mock.Mock

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

type MockUserDeptRepo struct {
	mock.Mock
}

func (m *MockUserDeptRepo) GetUserDepts(ctx context.Context, userID string) ([]*user.UserDept, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user.UserDept), args.Error(1)
}

func (m *MockUserDeptRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserDeptRepo) Create(ctx context.Context, item *user.UserDept) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

type MockUserPostRepo struct {
	mock.Mock
}

func (m *MockUserPostRepo) GetUserPosts(ctx context.Context, userID string) ([]*user.UserPost, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user.UserPost), args.Error(1)
}

func (m *MockUserPostRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserPostRepo) Create(ctx context.Context, item *user.UserPost) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

type MockUserRoleRepo struct {
	mock.Mock
}

func (m *MockUserRoleRepo) GetUserRoles(ctx context.Context, userID string) ([]*user.UserRole, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user.UserRole), args.Error(1)
}

func (m *MockUserRoleRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRoleRepo) Create(ctx context.Context, item *user.UserRole) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

type MockTransactionManager struct {
	mock.Mock
}

func (m *MockTransactionManager) Tx(ctx context.Context, fn func(context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

func TestUserUsecase_CreateUser(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepo)
	mockDeptRepo := new(MockUserDeptRepo)
	mockPostRepo := new(MockUserPostRepo)
	mockRoleRepo := new(MockUserRoleRepo)
	mockTm := new(MockTransactionManager)
	idg := idgen.NewIDGenerator()
	logger := log.DefaultLogger

	uc := user.NewUserUsecase(logger, mockRepo, mockTm, idg, mockDeptRepo, mockPostRepo, mockRoleRepo)

	newUser := &user.User{
		Username: "testuser",
		Password: "password",
		Nickname: "Test User",
	}

	mockRepo.On("FindByUsername", ctx, "testuser").Return(nil, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*user.User")).Return(nil)

	err := uc.CreateUser(ctx, newUser)

	assert.NoError(t, err)
	assert.NotEmpty(t, newUser.ID)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecase_GetUser(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepo)
	mockDeptRepo := new(MockUserDeptRepo)
	mockPostRepo := new(MockUserPostRepo)
	mockRoleRepo := new(MockUserRoleRepo)
	mockTm := new(MockTransactionManager)
	idg := idgen.NewIDGenerator()
	logger := log.DefaultLogger

	uc := user.NewUserUsecase(logger, mockRepo, mockTm, idg, mockDeptRepo, mockPostRepo, mockRoleRepo)

	expected := &user.User{
		ID:       "user-1",
		Username: "testuser",
		CreateAt: time.Now(),
	}

	mockRepo.On("FindByID", ctx, "user-1").Return(expected, nil)

	result, err := uc.GetUser(ctx, "user-1")

	assert.NoError(t, err)
	assert.Equal(t, "user-1", result.ID)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecase_GetUserByUsername(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepo)
	mockDeptRepo := new(MockUserDeptRepo)
	mockPostRepo := new(MockUserPostRepo)
	mockRoleRepo := new(MockUserRoleRepo)
	mockTm := new(MockTransactionManager)
	idg := idgen.NewIDGenerator()
	logger := log.DefaultLogger

	uc := user.NewUserUsecase(logger, mockRepo, mockTm, idg, mockDeptRepo, mockPostRepo, mockRoleRepo)

	expected := &user.User{
		ID:       "user-1",
		Username: "testuser",
	}

	mockRepo.On("FindByUsername", ctx, "testuser").Return(expected, nil)

	result, err := uc.GetUserByUsername(ctx, "testuser")

	assert.NoError(t, err)
	assert.Equal(t, "user-1", result.ID)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecase_UpdateUser(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepo)
	mockDeptRepo := new(MockUserDeptRepo)
	mockPostRepo := new(MockUserPostRepo)
	mockRoleRepo := new(MockUserRoleRepo)
	mockTm := new(MockTransactionManager)
	idg := idgen.NewIDGenerator()
	logger := log.DefaultLogger

	uc := user.NewUserUsecase(logger, mockRepo, mockTm, idg, mockDeptRepo, mockPostRepo, mockRoleRepo)

	usr := &user.User{
		ID:       "user-1",
		Username: "testuser",
		Nickname: "Updated Name",
	}

	mockRepo.On("Update", ctx, mock.AnythingOfType("*user.User")).Return(nil)

	err := uc.UpdateUser(ctx, usr)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserUsercase_ChangeUserStatus(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepo)
	mockDeptRepo := new(MockUserDeptRepo)
	mockPostRepo := new(MockUserPostRepo)
	mockRoleRepo := new(MockUserRoleRepo)
	mockTm := new(MockTransactionManager)
	idg := idgen.NewIDGenerator()
	logger := log.DefaultLogger

	uc := user.NewUserUsecase(logger, mockRepo, mockTm, idg, mockDeptRepo, mockPostRepo, mockRoleRepo)

	usr := &user.User{ID: "user-1"}
	bo := &user.UpdateStatusBO{
		UserID: "user-1",
		Status: 0,
	}

	mockRepo.On("FindByID", ctx, "user-1").Return(usr, nil)
	mockRepo.On("UpdateStatus", ctx, mock.AnythingOfType("*user.UpdateStatusBO")).Return(nil)

	err := uc.ChangeUserStatus(ctx, bo)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecase_DeleteUser(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepo)
	mockDeptRepo := new(MockUserDeptRepo)
	mockPostRepo := new(MockUserPostRepo)
	mockRoleRepo := new(MockUserRoleRepo)
	mockTm := new(MockTransactionManager)
	idg := idgen.NewIDGenerator()
	logger := log.DefaultLogger

	uc := user.NewUserUsecase(logger, mockRepo, mockTm, idg, mockDeptRepo, mockPostRepo, mockRoleRepo)

	usr := &user.User{ID: "user-1"}

	mockRepo.On("FindByID", ctx, "user-1").Return(usr, nil)
	mockRepo.On("Delete", ctx, mock.AnythingOfType("*user.DeleteUserBO")).Return(nil)

	err := uc.DeleteUser(ctx, "user-1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserUsecase_UpdateLoginInfo(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserRepo)
	mockDeptRepo := new(MockUserDeptRepo)
	mockPostRepo := new(MockUserPostRepo)
	mockRoleRepo := new(MockUserRoleRepo)
	mockTm := new(MockTransactionManager)
	idg := idgen.NewIDGenerator()
	logger := log.DefaultLogger

	uc := user.NewUserUsecase(logger, mockRepo, mockTm, idg, mockDeptRepo, mockPostRepo, mockRoleRepo)

	bo := &user.UpdateLoginInfoBO{
		UserID:    "user-1",
		LoginIP:   "127.0.0.1",
		LoginDate: time.Now(),
	}

	mockRepo.On("UpdateLoginInfo", ctx, mock.AnythingOfType("*user.UpdateLoginInfoBO")).Return(nil)

	err := uc.UpdateLoginInfo(ctx, bo)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
