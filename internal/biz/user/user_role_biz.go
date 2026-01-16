package user

import (
	"context"
)

func (uc *UserUsecase) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	roles, err := uc.userRoleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("获取用户关联角色出现错误,userID:%s,error:%v", userID, err)
		return nil, err
	}
	return roles, nil
}

func (uc *UserUsecase) ManageUserRoles(ctx context.Context, bo *AssignUserRolesBO) error {
	uc.log.WithContext(ctx).Infof("ManageUserRoles: userID=%s, operation=%s, roleCount=%d", bo.UserID, bo.Operation, len(bo.RoleIDs))
	return uc.userRoleRepo.ManageUserRoles(ctx, bo)
}
