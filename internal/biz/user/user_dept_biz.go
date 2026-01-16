package user

import (
	"context"
)

func (uc *UserUsecase) GetUserDepts(ctx context.Context, userID string) ([]string, error) {
	uc.log.WithContext(ctx).Infof("GetUserDepts: userID=%s", userID)
	return uc.userDeptRepo.GetUserDepts(ctx, userID)
}

func (uc *UserUsecase) ManageUserDepts(ctx context.Context, bo *AssignUserDeptsBO) error {
	uc.log.WithContext(ctx).Infof("ManageUserDepts: userID=%s, operation=%s, deptCount=%d", bo.UserID, bo.Operation, len(bo.DeptIDs))
	return uc.userDeptRepo.ManageUserDepts(ctx, bo)
}
