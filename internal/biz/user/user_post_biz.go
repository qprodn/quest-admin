package user

import (
	"context"

	"quest-admin/internal/biz/organization"
)

func (uc *UserUsecase) GetUserPosts(ctx context.Context, userID string) ([]string, error) {
	return uc.userPostRepo.GetUserPosts(ctx, userID)
}

func (uc *UserUsecase) ManageUserPosts(ctx context.Context, bo *AssignUserPostsBO) error {
	if len(bo.PostIDs) == 0 {
		existingPostIDs, err := uc.userPostRepo.GetUserPosts(ctx, bo.UserID)
		if err != nil {
			return err
		}
		if len(existingPostIDs) > 0 {
			return uc.userPostRepo.DeleteUserPosts(ctx, bo.UserID, existingPostIDs)
		}
		return nil
	}

	exists, err := uc.userPostRepo.CheckPostsExist(ctx, bo.PostIDs)
	if err != nil {
		return err
	}
	if !exists {
		return organization.ErrPostNotFound
	}

	existingPostIDs, err := uc.userPostRepo.GetUserPosts(ctx, bo.UserID)
	if err != nil {
		return err
	}

	existingPostMap := make(map[string]bool)
	for _, postID := range existingPostIDs {
		existingPostMap[postID] = true
	}

	requestedPostMap := make(map[string]bool)
	for _, postID := range bo.PostIDs {
		requestedPostMap[postID] = true
	}

	toDelete := make([]string, 0)
	for postID := range existingPostMap {
		if !requestedPostMap[postID] {
			toDelete = append(toDelete, postID)
		}
	}

	toAdd := make([]string, 0)
	for postID := range requestedPostMap {
		if !existingPostMap[postID] {
			toAdd = append(toAdd, postID)
		}
	}

	if len(toDelete) > 0 {
		if err := uc.userPostRepo.DeleteUserPosts(ctx, bo.UserID, toDelete); err != nil {
			return err
		}
	}

	if len(toAdd) > 0 {
		if err := uc.userPostRepo.AddUserPosts(ctx, bo.UserID, toAdd); err != nil {
			return err
		}
	}

	return nil
}
