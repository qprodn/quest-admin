package user

import (
	"context"
	"quest-admin/pkg/lang/slices"
	"quest-admin/types/consts/id"
)

func (uc *UserUsecase) GetUserPosts(ctx context.Context, userID string) ([]string, error) {
	posts, err := uc.userPostRepo.GetUserPosts(ctx, userID)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("获取用户关联职位出现错误,userID:%s,error:%v", userID, err)
		return nil, err
	}
	return slices.Map(posts, func(item *UserPost, index int) string {
		return item.PostID
	}), nil
}

func (uc *UserUsecase) AssignUserPosts(ctx context.Context, bo *AssignUserPostsBO) error {
	dbUserPosts, err := uc.userPostRepo.GetUserPosts(ctx, bo.UserID)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("获取当前用户关联职位出现错误,error:%v", err)
		return err
	}
	dbUserPostCodes := slices.Map(dbUserPosts, func(item *UserPost, index int) string {
		return item.PostID
	})
	newUserPostCodes := bo.PostIDs
	needDelete, needInsert := slices.Difference(dbUserPostCodes, newUserPostCodes)
	err = uc.tm.Tx(ctx, func(ctx context.Context) error {
		for _, item := range needInsert {
			err = uc.userPostRepo.Create(ctx, &UserPost{
				ID:     uc.idgen.NextID(id.EMPTY),
				UserID: bo.UserID,
				PostID: item})
			if err != nil {
				uc.log.WithContext(ctx).Errorf("添加用户职位出现错误,userID:%s,postID:%s,error:%v", bo.UserID, item, err)
				return err
			}
		}
		for _, item := range dbUserPosts {
			if slices.Contains(needDelete, item.PostID) {
				err = uc.userPostRepo.Delete(ctx, item.ID)
				if err != nil {
					uc.log.WithContext(ctx).Errorf("删除用户职位出现错误,userID:%s,postID:%s,error:%v", bo.UserID, item.PostID, err)
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		uc.log.WithContext(ctx).Errorf("分配用户职位出现错误,error:%v", err)
		return err
	}
	return nil
}
