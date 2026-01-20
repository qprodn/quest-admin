package user

import (
	"context"
	"quest-admin/pkg/lang/slices"
)

func (uc *UserUsecase) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	roles, err := uc.userRoleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("获取用户关联角色出现错误,userID:%s,error:%v", userID, err)
		return nil, err
	}
	return slices.Map(roles, func(item *UserRole, index int) string {
		return item.RoleID
	}), nil
}

func (uc *UserUsecase) AssignUserRoles(ctx context.Context, bo *AssignUserRolesBO) error {
	dbUserRoles, err := uc.userRoleRepo.GetUserRoles(ctx, bo.UserID)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("获取当前用户关联角色出现错误,error:%v", err)
		return err
	}
	dbUserRoleCodes := slices.Map(dbUserRoles, func(item *UserRole, index int) string {
		return item.RoleID
	})
	newUserRoleCodes := bo.RoleIDs
	needDelete, needInsert := slices.Difference(dbUserRoleCodes, newUserRoleCodes)
	err = uc.tm.Tx(ctx, func(ctx context.Context) error {
		for _, item := range needInsert {
			err = uc.userRoleRepo.Create(ctx, &UserRole{UserID: bo.UserID, RoleID: item})
			if err != nil {
				uc.log.WithContext(ctx).Errorf("添加用户角色出现错误,userID:%s,roleID:%s,error:%v", bo.UserID, item, err)
				return err
			}
		}
		for _, item := range dbUserRoles {
			if slices.Contains(needDelete, item.RoleID) {
				err = uc.userRoleRepo.Delete(ctx, item.ID)
				if err != nil {
					uc.log.WithContext(ctx).Errorf("删除用户角色出现错误,userID:%s,roleID:%s,error:%v", bo.UserID, item.RoleID, err)
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		uc.log.WithContext(ctx).Errorf("分配用户角色出现错误,error:%v", err)
		return err
	}
	return nil
}
