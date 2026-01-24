package user_test

import (
	"context"
	user "quest-admin/internal/biz/user"

	"github.com/stretchr/testify/mock"
)

type MockUserRepoForRole struct {
	mock.Mock
}

func (m *MockUserRepoForRole) Create(ctx context.Context, user *user.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepoForRole) FindByID(ctx context.Context, id string) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepoForRole) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepoForRole) List(ctx context.Context, opt *user.WhereUserOpt) ([]*user.User, error) {
	args := m.Called(ctx, opt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user.User), args.Error(1)
}

func (m *MockUserRepoForRole) Count(ctx context.Context, opt *user.WhereUserOpt) (int64, error) {
	args := m.Called(ctx, opt)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepoForRole) Update(ctx context.Context, user *user.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepoForRole) UpdatePassword(ctx context.Context, bo *user.UpdatePasswordBO) error {
	args := m.Called(ctx, bo)
	return args.Error(0)
}

func (m *MockUserRepoForRole) UpdateStatus(ctx context.Context, bo *user.UpdateStatusBO) error {
	args := m.Called(ctx, bo)
	return args.Error(0)
}

func (m *MockUserRepoForRole) UpdateLoginInfo(ctx context.Context, bo *user.UpdateLoginInfoBO) error {
	args := m.Called(ctx, bo)
	return args.Error(0)
}

func (m *MockUserRepoForRole) Delete(ctx context.Context, bo *user.DeleteUserBO) error {
	args := m.Called(ctx, bo)
	return args.Error(0)
}

type MockUserDeptRepoForRole struct {
	mock.Mock
}

func (m *MockUserDeptRepoForRole) GetUserDepts(ctx context.Context, userID string) ([]*user.UserDept, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user.UserDept), args.Error(1)
}

func (m *MockUserDeptRepoForRole) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserDeptRepoForRole) Create(ctx context.Context, item *user.UserDept) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

type MockUserPostRepoForRole struct {
	mock.Mock
}

func (m *MockUserPostRepoForRole) GetUserPosts(ctx context.Context, userID string) ([]*user.UserPost, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user.UserPost), args.Error(1)
}

func (m *MockUserPostRepoForRole) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserPostRepoForRole) Create(ctx context.Context, item *user.UserPost) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

type MockUserRoleRepoForRole struct {
	mock.Mock
}

func (m *MockUserRoleRepoForRole) GetUserRoles(ctx context.Context, userID string) ([]*user.UserRole, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user.UserRole), args.Error(1)
}

func (m *MockUserRoleRepoForRole) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRoleRepoForRole) Create(ctx context.Context, item *user.UserRole) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

type MockTransactionManagerForRole struct {
	mock.Mock
}

func (m *MockTransactionManagerForRole) Tx(ctx context.Context, fn func(context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

//
//func TestUserRoleRepo_GetUserRoles(t *testing.T) {
//	ctx := context.Background()
//	mockRepo := new(MockUserRoleUserRepoForRole)
//
//	roles := []*user.UserRole{
//		{ID: "ur-1", UserID: "user-1", RoleID: "role-1"},
//		{ID: "ur-2", UserID: "user-1", RoleID: "role-2"},
//	}
//
//	mockRepo.On("GetUserRoles", ctx, "user-1").Return(roles, nil)
//
//	result, err := mockRepo.GetUserRoles(ctx, "user-1")
//
//	assert.NoError(t, err)
//	assert.Len(t, result, 2)
//	mockRepo.AssertExpectations(t)
//}
//
//func TestUserRoleRepo_Create(t *testing.T) {
//	ctx := context.Background()
//	mockRepo := new(MockUserRoleUserRepoForRole)
//
//	item := &user.UserRole{
//		ID:     "ur-1",
//		UserID: "user-1",
//		RoleID: "role-1",
//	}
//
//	mockRepo.On("Create", ctx, mock.AnythingOfType("*user.UserRole")).Return(nil)
//
//	err := mockRepo.Create(ctx, item)
//
//	assert.NoError(t, err)
//	mockRepo.AssertExpectations(t)
//}
//
//func TestUserRoleRepo_Delete(t *testing.T) {
//	ctx := context.Background()
//	mockRepo := new(MockUserRoleUserRepoForRole)
//
//	mockRepo.On("Delete", ctx, "ur-1").Return(nil)
//
//	err := mockRepo.Delete(ctx, "ur-1")
//
//	assert.NoError(t, err)
//	mockRepo.AssertExpectations(t)
//}
