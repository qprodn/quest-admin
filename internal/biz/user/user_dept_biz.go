package user

import (
	"context"
	"quest-admin/pkg/lang/slices"
	"quest-admin/types/consts/id"
)

func (uc *UserUsecase) GetUserDepts(ctx context.Context, userID string) ([]string, error) {
	depts, err := uc.userDeptRepo.GetUserDepts(ctx, userID)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("获取用户关联部门出现错误,userID:%s,error:%v", userID, err)
		return nil, err
	}
	return slices.Map(depts, func(item *UserDept, index int) string {
		return item.DeptID
	}), nil
}

func (uc *UserUsecase) AssignUserDepts(ctx context.Context, bo *AssignUserDeptsBO) error {
	dbUserDepts, err := uc.userDeptRepo.GetUserDepts(ctx, bo.UserID)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("获取当前用户关联部门出现错误,error:%v", err)
		return err
	}
	dbUserDeptCodes := slices.Map(dbUserDepts, func(item *UserDept, index int) string {
		return item.DeptID
	})
	newUserDeptCodes := bo.DeptIDs
	needDelete, needInsert := slices.Difference(dbUserDeptCodes, newUserDeptCodes)
	err = uc.tm.Tx(ctx, func(ctx context.Context) error {
		for _, item := range needInsert {
			err = uc.userDeptRepo.Create(ctx, &UserDept{
				ID:     uc.idgen.NextID(id.EMPTY),
				UserID: bo.UserID,
				DeptID: item})
			if err != nil {
				uc.log.WithContext(ctx).Errorf("添加用户部门出现错误,userID:%s,deptID:%s,error:%v", bo.UserID, item, err)
				return err
			}
		}
		for _, item := range dbUserDepts {
			if slices.Contains(needDelete, item.DeptID) {
				err = uc.userDeptRepo.Delete(ctx, item.ID)
				if err != nil {
					uc.log.WithContext(ctx).Errorf("删除用户部门出现错误,userID:%s,deptID:%s,error:%v", bo.UserID, item.DeptID, err)
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		uc.log.WithContext(ctx).Errorf("分配用户部门出现错误,error:%v", err)
		return err
	}
	return nil
}
