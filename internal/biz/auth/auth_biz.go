package auth

import (
	"context"
	permBiz "quest-admin/internal/biz/permission"
	userBiz "quest-admin/internal/biz/user"
	"quest-admin/internal/data/auth"

	"github.com/go-kratos/kratos/v2/log"
)

// AuthUsecase 认证用例
type AuthUsecase struct {
	authManager *auth.Manager
	userUsecase *userBiz.UserUsecase
	roleUsecase *permBiz.RoleUsecase
	menuUsecase *permBiz.MenuUsecase
	log         *log.Helper
}

// NewAuthUsecase 创建认证用例
func NewAuthUsecase(
	manager *auth.Manager,
	logger log.Logger,
	userUsecase *userBiz.UserUsecase,
	roleUsecase *permBiz.RoleUsecase,
	menuUsecase *permBiz.MenuUsecase) *AuthUsecase {
	return &AuthUsecase{
		authManager: manager,
		userUsecase: userUsecase,
		roleUsecase: roleUsecase,
		menuUsecase: menuUsecase,
		log:         log.NewHelper(log.With(logger, "module", "auth/biz/auth")),
	}
}

// AdminGenerateToken 生成访问令牌
func (uc *AuthUsecase) AdminGenerateToken(ctx context.Context, bo *GenerateTokenBO) (string, error) {
	token, err := uc.authManager.Admin.Login(bo.UserID)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("生成令牌出现错误,userID:%s,error:%v", bo.UserID, err)
		return "", err
	}
	return token, nil
}

// SetRolesAndPermission 设置用户权限信息
func (uc *AuthUsecase) SetRolesAndPermission(ctx context.Context, userID string, roles []string, permissions []string) error {
	if len(roles) != 0 {
		err := uc.authManager.Admin.SetRoles(userID, roles)
		if err != nil {
			uc.log.WithContext(ctx).Errorf("设置用户角色出现错误,userID:%s,error:%v", userID, err)
			return err
		}
	}
	if len(permissions) != 0 {
		err := uc.authManager.Admin.SetPermissions(userID, permissions)
		if err != nil {
			uc.log.WithContext(ctx).Errorf("设置用户权限出现错误,userID:%s,error:%v", userID, err)
			return err
		}
	}
	return nil
}

// Logout 用户登出
func (uc *AuthUsecase) Logout(ctx context.Context, token string) error {

	return nil
}
