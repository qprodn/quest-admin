package user

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type UserPostRepo interface {
	GetUserPosts(ctx context.Context, userID string) ([]string, error)
	ManageUserPosts(ctx context.Context, bo *AssignUserPostsBO) error
	CheckPostsExist(ctx context.Context, postIDs []string) (bool, error)
	DeleteUserPosts(ctx context.Context, userID string, postIDs []string) error
	AddUserPosts(ctx context.Context, userID string, postIDs []string) error
}

type UserPostUsecase struct {
	repo UserPostRepo
	log  *log.Helper
}

func NewUserPostUsecase(repo UserPostRepo, logger log.Logger) *UserPostUsecase {
	return &UserPostUsecase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "user/biz/user-post")),
	}
}

func (uc *UserPostUsecase) GetUserPosts(ctx context.Context, userID string) ([]string, error) {
	return uc.repo.GetUserPosts(ctx, userID)
}

func (uc *UserPostUsecase) ManageUserPosts(ctx context.Context, bo *AssignUserPostsBO) error {
	if len(bo.PostIDs) == 0 {
		existingPostIDs, err := uc.repo.GetUserPosts(ctx, bo.UserID)
		if err != nil {
			return err
		}
		if len(existingPostIDs) > 0 {
			return uc.repo.DeleteUserPosts(ctx, bo.UserID, existingPostIDs)
		}
		return nil
	}

	exists, err := uc.repo.CheckPostsExist(ctx, bo.PostIDs)
	if err != nil {
		return err
	}
	if !exists {
		return ErrPostNotFound
	}

	existingPostIDs, err := uc.repo.GetUserPosts(ctx, bo.UserID)
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
		if err := uc.repo.DeleteUserPosts(ctx, bo.UserID, toDelete); err != nil {
			return err
		}
	}

	if len(toAdd) > 0 {
		if err := uc.repo.AddUserPosts(ctx, bo.UserID, toAdd); err != nil {
			return err
		}
	}

	return nil
}
