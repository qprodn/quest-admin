package user

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type UserDeptRepo interface {
	GetUserDepts(ctx context.Context, userID string) ([]string, error)
	ManageUserDepts(ctx context.Context, bo *AssignUserDeptsBO) error
}

type UserDeptUsecase struct {
	repo UserDeptRepo
	log  *log.Helper
}

func NewUserDeptUsecase(repo UserDeptRepo, logger log.Logger) *UserDeptUsecase {
	return &UserDeptUsecase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "user/biz/user-dept")),
	}
}

func (uc *UserDeptUsecase) GetUserDepts(ctx context.Context, userID string) ([]string, error) {
	uc.log.WithContext(ctx).Infof("GetUserDepts: userID=%s", userID)
	return uc.repo.GetUserDepts(ctx, userID)
}

func (uc *UserDeptUsecase) ManageUserDepts(ctx context.Context, bo *AssignUserDeptsBO) error {
	uc.log.WithContext(ctx).Infof("ManageUserDepts: userID=%s, operation=%s, deptCount=%d", bo.UserID, bo.Operation, len(bo.DeptIDs))
	return uc.repo.ManageUserDepts(ctx, bo)
}
