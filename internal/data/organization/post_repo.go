package organization

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	biz "quest-admin/internal/biz/organization"
	"quest-admin/internal/data/data"
	"quest-admin/pkg/util/idgen"

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
		ID:       idgen.GenerateID(),
		Name:     post.Name,
		Code:     post.Code,
		Sort:     post.Sort,
		Status:   post.Status,
		Remark:   post.Remark,
		CreateBy: post.CreateBy,
		CreateAt: now,
		UpdateBy: post.UpdateBy,
		UpdateAt: now,
		TenantID: post.TenantID,
	}

	_, err := r.data.DB(ctx).NewInsert().Model(dbPost).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return r.toBizPost(dbPost), nil
}

func (r *postRepo) FindByID(ctx context.Context, id string) (*biz.Post, error) {
	dbPost := &Post{ID: id}
	err := r.data.DB(ctx).NewSelect().Model(dbPost).WherePK().Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, biz.ErrPostNotFound
		}
		return nil, err
	}
	return r.toBizPost(dbPost), nil
}

func (r *postRepo) FindByName(ctx context.Context, name string) (*biz.Post, error) {
	dbPost := &Post{}
	err := r.data.DB(ctx).NewSelect().Model(dbPost).Where("name = ?", name).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, biz.ErrPostNotFound
		}
		return nil, err
	}
	return r.toBizPost(dbPost), nil
}

func (r *postRepo) FindByCode(ctx context.Context, code string) (*biz.Post, error) {
	dbPost := &Post{}
	err := r.data.DB(ctx).NewSelect().Model(dbPost).Where("code = ?", code).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, biz.ErrPostNotFound
		}
		return nil, err
	}
	return r.toBizPost(dbPost), nil
}

func (r *postRepo) List(ctx context.Context, query *biz.ListPostsQuery) (*biz.ListPostsResult, error) {
	var dbPosts []*Post
	q := r.data.DB(ctx).NewSelect().Model(&dbPosts)

	if query.Keyword != "" {
		q = q.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.WhereOr("name LIKE ?", "%"+query.Keyword+"%").
				WhereOr("code LIKE ?", "%"+query.Keyword+"%")
		})
	}

	if query.Status != nil {
		q = q.Where("status = ?", *query.Status)
	}

	total, err := q.ScanAndCount(ctx, &dbPosts, nil)
	if err != nil {
		return nil, err
	}

	page := query.Page
	pageSize := query.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	q = q.Offset(int(offset)).Limit(int(pageSize))

	if query.SortField != "" {
		order := query.SortOrder
		if order != "asc" && order != "desc" {
			order = "asc"
		}
		q = q.Order(fmt.Sprintf("%s %s", query.SortField, order))
	} else {
		q = q.Order("sort ASC, create_at DESC")
	}

	err = q.Scan(ctx)
	if err != nil {
		return nil, err
	}

	totalPages := int32((int64(total) + int64(pageSize) - 1) / int64(pageSize))

	posts := make([]*biz.Post, 0, len(dbPosts))
	for _, dbPost := range dbPosts {
		posts = append(posts, r.toBizPost(dbPost))
	}

	return &biz.ListPostsResult{
		Posts:      posts,
		Total:      int64(total),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
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

	_, err := r.data.DB(ctx).NewUpdate().Model(dbPost).WherePK().OmitZero().Exec(ctx)
	if err != nil {
		return nil, err
	}

	return r.FindByID(ctx, post.ID)
}

func (r *postRepo) Delete(ctx context.Context, id string) error {
	_, err := r.data.DB(ctx).NewDelete().
		Model((*Post)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *postRepo) HasUsers(ctx context.Context, id string) (bool, error) {
	var count int
	count, err := r.data.DB(ctx).NewSelect().
		Model((*Post)(nil)).
		TableExpr("qa_user_map_post AS up").
		Where("up.post_id = ?", id).
		Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
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
