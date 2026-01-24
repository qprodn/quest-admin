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
	tests := []struct {
		name        string
		setupMock   func(*MockUserRepo)
		inputUser   *user.User
		expectError bool
	}{
		{
			name: "success",
			setupMock: func(m *MockUserRepo) {
				m.On("Create", ctx, mock.AnythingOfType("*user.User")).Return(nil)
			},
			inputUser: &user.User{
				ID:       "user-1",
				Username: "testuser",
				Password: "hashedpass",
				Nickname: "Test User",
			},
			expectError: false,
		},
		{
			name: "error",
			setupMock: func(m *MockUserRepo) {
				m.On("Create", ctx, mock.AnythingOfType("*user.User")).Return(assert.AnError)
			},
			inputUser: &user.User{
				ID:       "user-1",
				Username: "testuser",
				Password: "hashedpass",
				Nickname: "Test User",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			tt.setupMock(mockRepo)

			err := mockRepo.Create(ctx, tt.inputUser)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserRepo_FindByID(t *testing.T) {
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
					Nickname: "Test User",
				}
				m.On("FindByID", ctx, "user-1").Return(expected, nil)
			},
			inputID:     "user-1",
			expectError: false,
			expectUser: &user.User{
				ID:       "user-1",
				Username: "testuser",
				Nickname: "Test User",
			},
		},
		{
			name: "error",
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
			mockRepo := new(MockUserRepo)
			tt.setupMock(mockRepo)

			result, err := mockRepo.FindByID(ctx, tt.inputID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectUser.ID, result.ID)
				assert.Equal(t, tt.expectUser.Username, result.Username)
				assert.Equal(t, tt.expectUser.Nickname, result.Nickname)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserRepo_FindByUsername(t *testing.T) {
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
			name: "error",
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
			mockRepo := new(MockUserRepo)
			tt.setupMock(mockRepo)

			result, err := mockRepo.FindByUsername(ctx, tt.inputUser)

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

func TestUserRepo_List(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name        string
		setupMock   func(*MockUserRepo)
		inputOpt    *user.WhereUserOpt
		expectError bool
	}{
		{
			name: "success",
			setupMock: func(m *MockUserRepo) {
				users := []*user.User{
					{ID: "user-1", Username: "user1"},
					{ID: "user-2", Username: "user2"},
				}
				m.On("List", ctx, mock.AnythingOfType("*user.WhereUserOpt")).Return(users, nil)
			},
			inputOpt: &user.WhereUserOpt{
				Limit:  10,
				Offset: 0,
			},
			expectError: false,
		},
		{
			name: "error",
			setupMock: func(m *MockUserRepo) {
				m.On("List", ctx, mock.AnythingOfType("*user.WhereUserOpt")).Return(nil, assert.AnError)
			},
			inputOpt: &user.WhereUserOpt{
				Limit:  10,
				Offset: 0,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			tt.setupMock(mockRepo)

			result, err := mockRepo.List(ctx, tt.inputOpt)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserRepo_Count(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name        string
		setupMock   func(*MockUserRepo)
		inputOpt    *user.WhereUserOpt
		expectError bool
		expectCount int64
	}{
		{
			name: "success",
			setupMock: func(m *MockUserRepo) {
				m.On("Count", ctx, mock.AnythingOfType("*user.WhereUserOpt")).Return(int64(2), nil)
			},
			inputOpt: &user.WhereUserOpt{
				Limit:  10,
				Offset: 0,
			},
			expectError: false,
			expectCount: 2,
		},
		{
			name: "error",
			setupMock: func(m *MockUserRepo) {
				m.On("Count", ctx, mock.AnythingOfType("*user.WhereUserOpt")).Return(int64(0), assert.AnError)
			},
			inputOpt: &user.WhereUserOpt{
				Limit:  10,
				Offset: 0,
			},
			expectError: true,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			tt.setupMock(mockRepo)

			result, err := mockRepo.Count(ctx, tt.inputOpt)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectCount, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserRepo_Update(t *testing.T) {
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
			name: "error",
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
			mockRepo := new(MockUserRepo)
			tt.setupMock(mockRepo)

			err := mockRepo.Update(ctx, tt.inputUser)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserRepo_Delete(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name        string
		setupMock   func(*MockUserRepo)
		inputBO     *user.DeleteUserBO
		expectError bool
	}{
		{
			name: "success",
			setupMock: func(m *MockUserRepo) {
				m.On("Delete", ctx, mock.AnythingOfType("*user.DeleteUserBO")).Return(nil)
			},
			inputBO:     &user.DeleteUserBO{UserID: "user-1"},
			expectError: false,
		},
		{
			name: "error",
			setupMock: func(m *MockUserRepo) {
				m.On("Delete", ctx, mock.AnythingOfType("*user.DeleteUserBO")).Return(assert.AnError)
			},
			inputBO:     &user.DeleteUserBO{UserID: "user-1"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			tt.setupMock(mockRepo)

			err := mockRepo.Delete(ctx, tt.inputBO)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
