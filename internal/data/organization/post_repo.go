package organization

import (
	"context"
	"database/sql"
	"fmt"
	"quest-admin/internal/data/data"
	"quest-admin/pkg/util/ctxs"
	"time"

	biz "quest-admin/internal/biz/organization"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/uptrace/bun"
)

type Post struct {
	bun.BaseModel `bun:"table:qa_post,alias:p"`

	ID       string     `bun:"id,pk"`
	Code     string     `bun:"code,notnull"`
	Name     string     `bun:"name,notnull"`
	Sort     int32      `bun:"sort,notnull"`
	Status   int32      `bun:"status,notnull"`
	Remark   string     `bun:"remark"`
	CreateBy string     `bun:"create_by"`
	CreateAt time.Time  `bun:"create_at,notnull,default:current_timestamp()"`
	UpdateBy string     `bun:"update_by"`
	UpdateAt time.Time  `bun:"update_at,notnull,default:current_timestamp()"`
	TenantID string     `bun:"tenant_id,notnull"`
	DeleteAt *time.Time `bun:"delete_at,soft_delete,nullzero"`
}

type postRepo struct {
	data *data.Data
	log  *log.Helper
}

func NewPostRepo(data *data.Data, logger log.Logger) biz.PostRepo {
	return &postRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *postRepo) Create(ctx context.Context, post *biz.Post) (*biz.Post, error) {
	now := time.Now()
	dbPost := &Post{
		ID:       post.ID,
		Name:     post.Name,
		Code:     post.Code,
		Sort:     post.Sort,
		Status:   post.Status,
		Remark:   post.Remark,
		CreateBy: ctxs.GetLoginID(ctx),
		CreateAt: now,
		UpdateBy: ctxs.GetLoginID(ctx),
		UpdateAt: now,
		TenantID: ctxs.GetTenantID(ctx),
	}

	_, err := r.data.DB(ctx).NewInsert().Model(dbPost).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return r.toBizPost(dbPost), nil
}

func (r *postRepo) FindByID(ctx context.Context, id string) (*biz.Post, error) {
	dbPost := &Post{ID: id}
	err := r.data.DB(ctx).
		NewSelect().
		Model(dbPost).
		WherePK().
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return r.toBizPost(dbPost), nil
}

func (r *postRepo) FindByName(ctx context.Context, name string) (*biz.Post, error) {
	dbPost := &Post{}
	err := r.data.DB(ctx).
		NewSelect().
		Model(dbPost).
		Where("name = ?", name).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return r.toBizPost(dbPost), nil
}

func (r *postRepo) FindByCode(ctx context.Context, code string) (*biz.Post, error) {
	dbPost := &Post{}
	err := r.data.DB(ctx).
		NewSelect().
		Model(dbPost).
		Where("code = ?", code).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return r.toBizPost(dbPost), nil
}

func (r *postRepo) List(ctx context.Context, opt *biz.WherePostOpt) ([]*biz.Post, error) {
	var dbPosts []*Post
	q := r.data.DB(ctx).NewSelect().Model(&dbPosts)

	if opt.Keyword != "" {
		q = q.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.WhereOr("name LIKE ?", "%"+opt.Keyword+"%").
				WhereOr("code LIKE ?", "%"+opt.Keyword+"%")
		})
	}

	if opt.Status != nil {
		q = q.Where("status = ?", *opt.Status)
	}

	if opt.Offset != 0 {
		q.Offset(int(opt.Offset))
	}
	if opt.Limit != 0 {
		q.Limit(int(opt.Limit))
	}

	if opt.SortField != "" && opt.SortOrder != "" {
		order := opt.SortOrder
		if order != "asc" && order != "desc" {
			order = "asc"
		}
		q = q.Order(fmt.Sprintf("%s %s", opt.SortField, order))
	} else {
		q = q.Order("sort ASC, create_at DESC")
	}

	err := q.Where("tenant_id = ?", ctxs.GetTenantID(ctx)).Scan(ctx)
	if err != nil {
		return nil, err
	}

	posts := make([]*biz.Post, 0, len(dbPosts))
	for _, dbPost := range dbPosts {
		posts = append(posts, r.toBizPost(dbPost))
	}

	return posts, nil
}

func (r *postRepo) Count(ctx context.Context, opt *biz.WherePostOpt) (int64, error) {
	var dbPosts []*Post
	q := r.data.DB(ctx).NewSelect().Model(&dbPosts)

	if opt.Keyword != "" {
		q = q.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.WhereOr("name LIKE ?", "%"+opt.Keyword+"%").
				WhereOr("code LIKE ?", "%"+opt.Keyword+"%")
		})
	}

	if opt.Status != nil {
		q = q.Where("status = ?", *opt.Status)
	}

	total, err := q.Where("tenant_id = ?", ctxs.GetTenantID(ctx)).Count(ctx)
	if err != nil {
		return 0, err
	}

	return int64(total), nil
}

func (r *postRepo) Update(ctx context.Context, post *biz.Post) (*biz.Post, error) {
	dbPost := &Post{
		ID:       post.ID,
		Name:     post.Name,
		Code:     post.Code,
		Sort:     post.Sort,
		Status:   post.Status,
		Remark:   post.Remark,
		UpdateAt: time.Now(),
	}

	_, err := r.data.DB(ctx).
		NewUpdate().
		Model(dbPost).
		WherePK().
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		OmitZero().
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return r.FindByID(ctx, post.ID)
}

func (r *postRepo) Delete(ctx context.Context, id string) error {
	_, err := r.data.DB(ctx).NewDelete().
		Model((*Post)(nil)).
		Where("id = ?", id).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Exec(ctx)
	return err
}

func (r *postRepo) HasUsers(ctx context.Context, id string) (bool, error) {
	var count int
	count, err := r.data.DB(ctx).NewSelect().
		Model((*Post)(nil)).
		TableExpr("qa_user_post AS up").
		Where("up.post_id = ?", id).
		Where("up.tenant_id = ?", ctxs.GetTenantID(ctx)).
		Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *postRepo) FindListByIDs(ctx context.Context, ids []string) ([]*biz.Post, error) {
	if len(ids) == 0 {
		return []*biz.Post{}, nil
	}
	var dbPosts []*Post
	err := r.data.DB(ctx).NewSelect().
		Model(&dbPosts).
		Where("id IN (?)", bun.In(ids)).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	posts := make([]*biz.Post, 0, len(dbPosts))
	for _, dbPost := range dbPosts {
		posts = append(posts, r.toBizPost(dbPost))
	}
	return posts, nil
}

func (r *postRepo) toBizPost(dbPost *Post) *biz.Post {
	return &biz.Post{
		ID:       dbPost.ID,
		Name:     dbPost.Name,
		Code:     dbPost.Code,
		Sort:     dbPost.Sort,
		Status:   dbPost.Status,
		Remark:   dbPost.Remark,
		CreateBy: dbPost.CreateBy,
		CreateAt: dbPost.CreateAt,
		UpdateBy: dbPost.UpdateBy,
		UpdateAt: dbPost.UpdateAt,
		TenantID: dbPost.TenantID,
	}
}
