package auth

import (
	"context"

	"quest-admin/internal/data/auth"

	"github.com/click33/sa-token-go/stputil"
	"github.com/go-kratos/kratos/v2/log"
)

// AuthUsecase 认证用例
type AuthUsecase struct {
	stp *stputil.StpLogic
	log *log.Helper
}

// NewAuthUsecase 创建认证用例
func NewAuthUsecase(manager *auth.Manager, logger log.Logger) *AuthUsecase {
	return &AuthUsecase{
		stp: manager.Admin,
		log: log.NewHelper(log.With(logger, "module", "auth/biz/auth")),
	}
}

// GenerateToken 生成访问令牌
func (uc *AuthUsecase) GenerateToken(ctx context.Context, userID string) (string, error) {
	uc.log.WithContext(ctx).Infof("GenerateToken: userID=%s", userID)

	token, err := uc.stp.Login(userID, "")
	if err != nil {
		return "", err
	}

	return token, nil
}

// SetPermission 设置用户权限信息
func (uc *AuthUsecase) SetPermission(ctx context.Context, userID string, roles []string, permissions []string) error {
	uc.log.WithContext(ctx).Infof("SetPermission: userID=%s, roles=%v, permissions=%v", userID, roles, permissions)

	// TODO: 设置角色列表和权限列表，需要根据 sa-token-go 的 API 实现
	// if len(roles) > 0 {
	// 	uc.stp.SetRoleList(userID, roles)
	// }
	// if len(permissions) > 0 {
	// 	uc.stp.SetPermissionList(userID, permissions)
	// }

	return nil
}

// GetLoginID 从令牌中获取登录用户ID
func (uc *AuthUsecase) GetLoginID(ctx context.Context, token string) (string, error) {
	uc.log.WithContext(ctx).Infof("GetLoginID: token=%s", token)

	// TODO: 从 token 中获取登录用户 ID，需要根据 sa-token-go 的 API 实现
	// loginID := uc.stp.GetLoginIdByToken(token)
	// if loginID == "" {
	// 	return "", ErrTokenInvalid
	// }
	// return loginID, nil

	return "", nil
}

// Logout 用户登出
func (uc *AuthUsecase) Logout(ctx context.Context, token string) error {
	uc.log.WithContext(ctx).Infof("Logout: token=%s", token)

	// TODO: 用户登出，需要根据 sa-token-go 的 API 实现
	// loginID := uc.stp.GetLoginIdByToken(token)
	// if loginID == "" {
	// 	return ErrTokenInvalid
	// }
	// uc.stp.Logout(loginID)

	return nil
}
