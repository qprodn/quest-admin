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

func newTestUsecase(t *testing.T) (*user.UserUsecase, *MockUserRepo) {
	mockRepo := new(MockUserRepo)
	mockDeptRepo := new(MockUserDeptRepo)
	mockPostRepo := new(MockUserPostRepo)
	mockRoleRepo := new(MockUserRoleRepo)
	mockTm := new(MockTransactionManager)
	idg := idgen.NewIDGenerator()
	logger := log.DefaultLogger

	uc := user.NewUserUsecase(logger, mockRepo, mockTm, idg, mockDeptRepo, mockPostRepo, mockRoleRepo)
	return uc, mockRepo
}

func TestUserUsecase_CreateUser(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name            string
		setupMock       func(*MockUserRepo)
		inputUser       *user.User
		expectError     bool
		expectUserIDSet bool
	}{
		{
			name: "success",
			setupMock: func(m *MockUserRepo) {
				m.On("FindByUsername", ctx, "testuser").Return(nil, nil)
				m.On("Create", ctx, mock.AnythingOfType("*user.User")).Return(nil)
			},
			inputUser: &user.User{
				Username: "testuser",
				Password: "password",
				Nickname: "Test User",
			},
			expectError:     false,
			expectUserIDSet: true,
		},
		{
			name: "repo error when checking existing user",
			setupMock: func(m *MockUserRepo) {
				m.On("FindByUsername", ctx, "testuser").Return(nil, assert.AnError)
			},
			inputUser: &user.User{
				Username: "testuser",
				Password: "password",
				Nickname: "Test User",
			},
			expectError:     true,
			expectUserIDSet: false,
		},
		{
			name: "existing user",
			setupMock: func(m *MockUserRepo) {
				existing := &user.User{ID: "existing-1", Username: "testuser"}
				m.On("FindByUsername", ctx, "testuser").Return(existing, nil)
			},
			inputUser: &user.User{
				Username: "testuser",
				Password: "password",
				Nickname: "Test User",
			},
			expectError:     true,
			expectUserIDSet: false,
		},
		{
			name: "create error",
			setupMock: func(m *MockUserRepo) {
				m.On("FindByUsername", ctx, "testuser").Return(nil, nil)
				m.On("Create", ctx, mock.AnythingOfType("*user.User")).Return(assert.AnError)
			},
			inputUser: &user.User{
				Username: "testuser",
				Password: "password",
				Nickname: "Test User",
			},
			expectError:     true,
			expectUserIDSet: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, mockRepo := newTestUsecase(t)
			tt.setupMock(mockRepo)

			newUser := *tt.inputUser
			err := uc.CreateUser(ctx, &newUser)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.expectUserIDSet && !tt.expectError {
				assert.NotEmpty(t, newUser.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserUsecase_GetUser(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name        string
		setupMock   func(*MockUserRepo)
		inputID     string
		expectError bool
		expectUser  *user.User
	}{
		{
			name: "success",
			setupMock: func(m *MockUserRepo) {
				expected := &user.User{
					ID:       "user-1",
					Username: "testuser",
					CreateAt: time.Now(),
				}
				m.On("FindByID", ctx, "user-1").Return(expected, nil)
			},
			inputID:     "user-1",
			expectError: false,
			expectUser: &user.User{
				ID:       "user-1",
				Username: "testuser",
			},
		},
		{
			name: "not found",
			setupMock: func(m *MockUserRepo) {
				m.On("FindByID", ctx, "user-1").Return(nil, assert.AnError)
			},
			inputID:     "user-1",
			expectError: true,
			expectUser:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, mockRepo := newTestUsecase(t)
			tt.setupMock(mockRepo)

			result, err := uc.GetUser(ctx, tt.inputID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectUser.ID, result.ID)
				assert.Equal(t, tt.expectUser.Username, result.Username)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserUsecase_GetUserByUsername(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name        string
		setupMock   func(*MockUserRepo)
		inputUser   string
		expectError bool
		expectUser  *user.User
	}{
		{
			name: "success",
			setupMock: func(m *MockUserRepo) {
				expected := &user.User{
					ID:       "user-1",
					Username: "testuser",
				}
				m.On("FindByUsername", ctx, "testuser").Return(expected, nil)
			},
			inputUser:   "testuser",
			expectError: false,
			expectUser: &user.User{
				ID:       "user-1",
				Username: "testuser",
			},
		},
		{
			name: "not found",
			setupMock: func(m *MockUserRepo) {
				m.On("FindByUsername", ctx, "testuser").Return(nil, assert.AnError)
			},
			inputUser:   "testuser",
			expectError: true,
			expectUser:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, mockRepo := newTestUsecase(t)
			tt.setupMock(mockRepo)

			result, err := uc.GetUserByUsername(ctx, tt.inputUser)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectUser.ID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserUsecase_UpdateUser(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name        string
		setupMock   func(*MockUserRepo)
		inputUser   *user.User
		expectError bool
	}{
		{
			name: "success",
			setupMock: func(m *MockUserRepo) {
				m.On("Update", ctx, mock.AnythingOfType("*user.User")).Return(nil)
			},
			inputUser: &user.User{
				ID:       "user-1",
				Username: "testuser",
				Nickname: "Updated Name",
			},
			expectError: false,
		},
		{
			name: "update error",
			setupMock: func(m *MockUserRepo) {
				m.On("Update", ctx, mock.AnythingOfType("*user.User")).Return(assert.AnError)
			},
			inputUser: &user.User{
				ID:       "user-1",
				Username: "testuser",
				Nickname: "Updated Name",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, mockRepo := newTestUsecase(t)
			tt.setupMock(mockRepo)

			err := uc.UpdateUser(ctx, tt.inputUser)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserUsecase_ChangeUserStatus(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name        string
		setupMock   func(*MockUserRepo)
		inputBO     *user.UpdateStatusBO
		expectError bool
	}{
		{
			name: "success",
			setupMock: func(m *MockUserRepo) {
				usr := &user.User{ID: "user-1"}
				m.On("FindByID", ctx, "user-1").Return(usr, nil)
				m.On("UpdateStatus", ctx, mock.AnythingOfType("*user.UpdateStatusBO")).Return(nil)
			},
			inputBO: &user.UpdateStatusBO{
				UserID: "user-1",
				Status: 0,
			},
			expectError: false,
		},
		{
			name: "user not found",
			setupMock: func(m *MockUserRepo) {
				m.On("FindByID", ctx, "user-1").Return(nil, assert.AnError)
			},
			inputBO: &user.UpdateStatusBO{
				UserID: "user-1",
				Status: 0,
			},
			expectError: true,
		},
		{
			name: "update error",
			setupMock: func(m *MockUserRepo) {
				usr := &user.User{ID: "user-1"}
				m.On("FindByID", ctx, "user-1").Return(usr, nil)
				m.On("UpdateStatus", ctx, mock.AnythingOfType("*user.UpdateStatusBO")).Return(assert.AnError)
			},
			inputBO: &user.UpdateStatusBO{
				UserID: "user-1",
				Status: 0,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, mockRepo := newTestUsecase(t)
			tt.setupMock(mockRepo)

			err := uc.ChangeUserStatus(ctx, tt.inputBO)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserUsecase_DeleteUser(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name        string
		setupMock   func(*MockUserRepo)
		inputID     string
		expectError bool
	}{
		{
			name: "success",
			setupMock: func(m *MockUserRepo) {
				usr := &user.User{ID: "user-1"}
				m.On("FindByID", ctx, "user-1").Return(usr, nil)
				m.On("Delete", ctx, mock.AnythingOfType("*user.DeleteUserBO")).Return(nil)
			},
			inputID:     "user-1",
			expectError: false,
		},
		{
			name: "user not found",
			setupMock: func(m *MockUserRepo) {
				m.On("FindByID", ctx, "user-1").Return(nil, assert.AnError)
			},
			inputID:     "user-1",
			expectError: true,
		},
		{
			name: "delete error",
			setupMock: func(m *MockUserRepo) {
				usr := &user.User{ID: "user-1"}
				m.On("FindByID", ctx, "user-1").Return(usr, nil)
				m.On("Delete", ctx, mock.AnythingOfType("*user.DeleteUserBO")).Return(assert.AnError)
			},
			inputID:     "user-1",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, mockRepo := newTestUsecase(t)
			tt.setupMock(mockRepo)

			err := uc.DeleteUser(ctx, tt.inputID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserUsecase_UpdateLoginInfo(t *testing.T) {
	ctx := context.Background()
	uc, mockRepo := newTestUsecase(t)

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
