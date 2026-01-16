package user

import (
	"context"
	"quest-admin/internal/data/data"
	"time"

	biz "quest-admin/internal/biz/user"
	"quest-admin/internal/data/organization"
	"quest-admin/pkg/util/idgen"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/uptrace/bun"
)

type UserPost struct {
	bun.BaseModel `bun:"table:qa_user_post,alias:up"`

	ID       string     `bun:"id,pk"`
	UserID   string     `bun:"user_id,notnull"`
	PostID   string     `bun:"post_id,notnull"`
	CreateBy string     `bun:"create_by"`
	CreateAt time.Time  `bun:"create_at,notnull,default:current_timestamp()"`
	UpdateBy string     `bun:"update_by"`
	UpdateAt time.Time  `bun:"update_at,notnull,default:current_timestamp()"`
	TenantID string     `bun:"tenant_id"`
	DeleteAt *time.Time `bun:"delete_at,soft_delete,nullzero"`
}

type userPostRepo struct {
	data *data.Data
	log  *log.Helper
}

func NewUserPostRepo(data *data.Data, logger log.Logger) biz.UserPostRepo {
	return &userPostRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userPostRepo) GetUserPosts(ctx context.Context, userID string) ([]string, error) {
	var userPosts []*UserPost
	err := r.data.DB(ctx).NewSelect().
		Model(&userPosts).
		Where("user_id = ?", userID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	postIDs := make([]string, 0, len(userPosts))
	for _, up := range userPosts {
		postIDs = append(postIDs, up.PostID)
	}
	return postIDs, nil
}

func (r *userPostRepo) ManageUserPosts(ctx context.Context, bo *biz.AssignUserPostsBO) error {
	switch bo.Operation {
	case "add":
		return r.addUserPosts(ctx, bo.UserID, bo.PostIDs)
	case "remove":
		return r.removeUserPosts(ctx, bo.UserID, bo.PostIDs)
	case "replace":
		return r.replaceUserPosts(ctx, bo.UserID, bo.PostIDs)
	default:
		return biz.ErrInvalidOperationType
	}
}

func (r *userPostRepo) addUserPosts(ctx context.Context, userID string, postIDs []string) error {
	now := time.Now()
	userPosts := make([]*UserPost, 0, len(postIDs))
	for _, postID := range postIDs {
		userPosts = append(userPosts, &UserPost{
			ID:       idgen.GenerateID(),
			UserID:   userID,
			PostID:   postID,
			CreateAt: now,
		})
	}

	_, err := r.data.DB(ctx).NewInsert().Model(&userPosts).Exec(ctx)
	return err
}

func (r *userPostRepo) removeUserPosts(ctx context.Context, userID string, postIDs []string) error {
	_, err := r.data.DB(ctx).NewDelete().
		Model((*UserPost)(nil)).
		Where("user_id = ?", userID).
		Where("post_id IN (?)", bun.In(postIDs)).
		Exec(ctx)
	return err
}

func (r *userPostRepo) replaceUserPosts(ctx context.Context, userID string, postIDs []string) error {
	return r.data.DB(ctx).RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewDelete().
			Model((*UserPost)(nil)).
			Where("user_id = ?", userID).
			Exec(ctx)
		if err != nil {
			return err
		}

		if len(postIDs) > 0 {
			now := time.Now()
			userPosts := make([]*UserPost, 0, len(postIDs))
			for _, postID := range postIDs {
				userPosts = append(userPosts, &UserPost{
					ID:       idgen.GenerateID(),
					UserID:   userID,
					PostID:   postID,
					CreateAt: now,
				})
			}
			_, err = tx.NewInsert().Model(&userPosts).Exec(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *userPostRepo) CheckPostsExist(ctx context.Context, postIDs []string) (bool, error) {
	if len(postIDs) == 0 {
		return true, nil
	}

	count, err := r.data.DB(ctx).NewSelect().
		Model((*organization.Post)(nil)).
		Where("id IN (?)", bun.In(postIDs)).
		Count(ctx)
	if err != nil {
		return false, err
	}

	return count == len(postIDs), nil
}

func (r *userPostRepo) DeleteUserPosts(ctx context.Context, userID string, postIDs []string) error {
	if len(postIDs) == 0 {
		return nil
	}

	_, err := r.data.DB(ctx).NewDelete().
		Model((*UserPost)(nil)).
		Where("user_id = ?", userID).
		Where("post_id IN (?)", bun.In(postIDs)).
		Exec(ctx)
	return err
}

func (r *userPostRepo) AddUserPosts(ctx context.Context, userID string, postIDs []string) error {
	if len(postIDs) == 0 {
		return nil
	}

	now := time.Now()
	userPosts := make([]*UserPost, 0, len(postIDs))
	for _, postID := range postIDs {
		userPosts = append(userPosts, &UserPost{
			ID:       idgen.GenerateID(),
			UserID:   userID,
			PostID:   postID,
			CreateAt: now,
		})
	}

	_, err := r.data.DB(ctx).NewInsert().Model(&userPosts).Exec(ctx)
	return err
}
