package user_test

import (
	"context"
	"testing"

	user "quest-admin/internal/biz/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//type MockUserDeptRepo struct {
//	mock.Mock
//}
//
//func (m *MockUserDeptRepo) GetUserDepts(ctx context.Context, userID string) ([]*user.UserDept, error) {
//	args := m.Called(ctx, userID)
//	if args.Get(0) == nil {
//		return nil, args.Error(1)
//	}
//	return args.Get(0).([]*user.UserDept), args.Error(1)
//}
//
//func (m *MockUserDeptRepo) Delete(ctx context.Context, id string) error {
//	args := m.Called(ctx, id)
//	return args.Error(0)
//}
//
//func (m *MockUserDeptRepo) Create(ctx context.Context, item *user.UserDept) error {
//	args := m.Called(ctx, item)
//	return args.Error(0)
//}

func TestUserDeptRepo_GetUserDepts(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserDeptRepo)

	depts := []*user.UserDept{
		{ID: "ud-1", UserID: "user-1", DeptID: "dept-1"},
		{ID: "ud-2", UserID: "user-1", DeptID: "dept-2"},
	}

	mockRepo.On("GetUserDepts", ctx, "user-1").Return(depts, nil)

	result, err := mockRepo.GetUserDepts(ctx, "user-1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestUserDeptRepo_Create(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserDeptRepo)

	item := &user.UserDept{
		ID:     "ud-1",
		UserID: "user-1",
		DeptID: "dept-1",
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*user.UserDept")).Return(nil)

	err := mockRepo.Create(ctx, item)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserDeptRepo_Delete(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockUserDeptRepo)

	mockRepo.On("Delete", ctx, "ud-1").Return(nil)

	err := mockRepo.Delete(ctx, "ud-1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
