package user

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type UserRoleRepo interface {
	GetUserRoles(ctx context.Context, userID string) ([]string, error)
	ManageUserRoles(ctx context.Context, bo *AssignUserRolesBO) error
}

type UserRoleUsecase struct {
	repo UserRoleRepo
	log  *log.Helper
}

func NewUserRoleUsecase(repo UserRoleRepo, logger log.Logger) *UserRoleUsecase {
	return &UserRoleUsecase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "user/biz/user-role")),
	}
}

func (uc *UserRoleUsecase) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	uc.log.WithContext(ctx).Infof("GetUserRoles: userID=%s", userID)
	return uc.repo.GetUserRoles(ctx, userID)
}

func (uc *UserRoleUsecase) ManageUserRoles(ctx context.Context, bo *AssignUserRolesBO) error {
	uc.log.WithContext(ctx).Infof("ManageUserRoles: userID=%s, operation=%s, roleCount=%d", bo.UserID, bo.Operation, len(bo.RoleIDs))
	return uc.repo.ManageUserRoles(ctx, bo)
}
