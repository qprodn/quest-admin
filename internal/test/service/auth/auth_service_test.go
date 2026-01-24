package auth_test

import (
	"context"
	"quest-admin/internal/biz/auth"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthBiz struct {
	mock.Mock
}

func (m *MockAuthBiz) AdminGenerateToken(ctx context.Context, bo *auth.GenerateTokenBO) (string, error) {
	args := m.Called(ctx, bo)
	return args.String(0), args.Error(1)
}

func TestAuthService_AdminGenerateToken(t *testing.T) {
	ctx := context.Background()
	mockBiz := new(MockAuthBiz)

	bo := &auth.GenerateTokenBO{
		UserID: "user-1",
	}

	mockBiz.On("AdminGenerateToken", ctx, mock.AnythingOfType("*auth.GenerateTokenBO")).Return("token-123", nil)

	token, err := mockBiz.AdminGenerateToken(ctx, bo)

	assert.NoError(t, err)
	assert.Equal(t, "token-123", token)
	mockBiz.AssertExpectations(t)
}
