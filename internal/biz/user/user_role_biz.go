package user

import (
	"context"
	v1 "quest-admin/api/gen/user/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrUserRoleNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

type UserRoleRepo interface {
	GetUserRoles(ctx context.Context, userID string) ([]string, error)
	ManageUserRoles(ctx context.Context, bo *ManageUserRolesBO) error
}

type UserRoleUsecase struct {
	repo UserRoleRepo
	log  *log.Helper
}

func NewUserRoleUsecase(repo UserRoleRepo, logger log.Logger) *UserRoleUsecase {
	return &UserRoleUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *UserRoleUsecase) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	uc.log.WithContext(ctx).Infof("GetUserRoles: userID=%s", userID)
	return uc.repo.GetUserRoles(ctx, userID)
}

func (uc *UserRoleUsecase) ManageUserRoles(ctx context.Context, bo *ManageUserRolesBO) error {
	uc.log.WithContext(ctx).Infof("ManageUserRoles: userID=%s, operation=%s, roleCount=%d", bo.UserID, bo.Operation, len(bo.RoleIDs))
	return uc.repo.ManageUserRoles(ctx, bo)
}
