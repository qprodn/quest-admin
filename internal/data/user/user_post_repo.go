package user

import (
	"context"
	"database/sql"
	"quest-admin/internal/data/data"
	"quest-admin/pkg/lang/slices"
	"quest-admin/pkg/util/ctxs"
	"time"

	biz "quest-admin/internal/biz/user"

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

func (r *userPostRepo) Create(ctx context.Context, item *biz.UserPost) error {
	if item == nil {
		return nil
	}
	now := time.Now()
	_, err := r.data.DB(ctx).NewInsert().Model(&UserPost{
		ID:       item.ID,
		UserID:   item.UserID,
		PostID:   item.PostID,
		CreateAt: now,
		CreateBy: ctxs.GetLoginID(ctx),
		UpdateAt: now,
		UpdateBy: ctxs.GetLoginID(ctx),
		TenantID: ctxs.GetTenantID(ctx),
	}).Exec(ctx)
	if err != nil {
		r.log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

func (r *userPostRepo) Delete(ctx context.Context, id string) error {
	_, err := r.data.DB(ctx).NewUpdate().
		Model((*UserPost)(nil)).
		Set("update_by = ?", ctxs.GetLoginID(ctx)).
		Where("id = ?", id).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Exec(ctx)
	if err != nil {
		r.log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

func (r *userPostRepo) GetUserPosts(ctx context.Context, userID string) ([]*biz.UserPost, error) {
	var userPosts []*UserPost
	err := r.data.DB(ctx).NewSelect().
		Model(&userPosts).
		Where("user_id = ?", userID).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]*biz.UserPost, 0), nil
		}
		r.log.WithContext(ctx).Error(err)
		return nil, err
	}

	return slices.Map(userPosts, func(item *UserPost, index int) *biz.UserPost {
		return r.toBizUserPost(item)
	}), nil
}

func (r *userPostRepo) toBizUserPost(item *UserPost) *biz.UserPost {
	return &biz.UserPost{
		ID:       item.ID,
		UserID:   item.UserID,
		PostID:   item.PostID,
		CreateBy: item.CreateBy,
		CreateAt: item.CreateAt,
		UpdateBy: item.UpdateBy,
		UpdateAt: item.UpdateAt,
		TenantID: item.TenantID,
	}
}
