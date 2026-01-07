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

type Department struct {
	bun.BaseModel `bun:"table:qa_dept,alias:d"`

	ID           string     `bun:"id,pk"`
	Name         string     `bun:"name,notnull"`
	ParentID     string     `bun:"parent_id,notnull"`
	Sort         int32      `bun:"sort,default:0"`
	LeaderUserID string     `bun:"leader_user_id"`
	Phone        string     `bun:"phone"`
	Email        string     `bun:"email"`
	Status       int32      `bun:"status,notnull"`
	CreateBy     string     `bun:"create_by"`
	CreateAt     time.Time  `bun:"create_at,notnull,default:current_timestamp()"`
	UpdateBy     string     `bun:"update_by"`
	UpdateAt     time.Time  `bun:"update_at,notnull,default:current_timestamp()"`
	TenantID     string     `bun:"tenant_id,notnull"`
	DeleteAt     *time.Time `bun:"delete_at,soft_delete,nullzero"`
}

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

type departmentRepo struct {
	data *data.Data
	log  *log.Helper
}

type postRepo struct {
	data *data.Data
	log  *log.Helper
}

func NewDepartmentRepo(data *data.Data, logger log.Logger) biz.DepartmentRepo {
	return &departmentRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func NewPostRepo(data *data.Data, logger log.Logger) biz.PostRepo {
	return &postRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *departmentRepo) Create(ctx context.Context, dept *biz.Department) (*biz.Department, error) {
	now := time.Now()
	dbDept := &Department{
		ID:           idgen.GenerateID(),
		Name:         dept.Name,
		ParentID:     dept.ParentID,
		Sort:         dept.Sort,
		LeaderUserID: dept.LeaderUserID,
		Phone:        dept.Phone,
		Email:        dept.Email,
		Status:       dept.Status,
		CreateBy:     dept.CreateBy,
		CreateAt:     now,
		UpdateBy:     dept.UpdateBy,
		UpdateAt:     now,
		TenantID:     dept.TenantID,
	}

	_, err := r.data.DB(ctx).NewInsert().Model(dbDept).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return r.toBizDepartment(dbDept), nil
}

func (r *departmentRepo) GetByID(ctx context.Context, id string) (*biz.Department, error) {
	dbDept := &Department{ID: id}
	err := r.data.DB(ctx).NewSelect().Model(dbDept).WherePK().Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, biz.ErrDepartmentNotFound
		}
		return nil, err
	}
	return r.toBizDepartment(dbDept), nil
}

func (r *departmentRepo) GetByName(ctx context.Context, name string) (*biz.Department, error) {
	dbDept := &Department{}
	err := r.data.DB(ctx).NewSelect().Model(dbDept).Where("name = ?", name).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, biz.ErrDepartmentNotFound
		}
		return nil, err
	}
	return r.toBizDepartment(dbDept), nil
}

func (r *departmentRepo) GetTree(ctx context.Context) ([]*biz.Department, error) {
	var dbDepts []*Department
	err := r.data.DB(ctx).NewSelect().
		Model(&dbDepts).
		Where("status = ?", 1).
		Order("level ASC, sort ASC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	deptMap := make(map[string]*biz.Department)
	for _, dbDept := range dbDepts {
		deptMap[dbDept.ID] = r.toBizDepartment(dbDept)
	}

	var roots []*biz.Department
	for _, dept := range deptMap {
		if dept.ParentID == "" || dept.ParentID == "0" {
			roots = append(roots, dept)
		} else if parent, ok := deptMap[dept.ParentID]; ok {
			parent.Children = append(parent.Children, dept)
		}
	}

	return roots, nil
}

func (r *departmentRepo) GetChildren(ctx context.Context, parentID string) ([]*biz.Department, error) {
	var dbDepts []*Department
	err := r.data.DB(ctx).NewSelect().
		Model(&dbDepts).
		Where("parent_id = ?", parentID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	departments := make([]*biz.Department, 0, len(dbDepts))
	for _, dbDept := range dbDepts {
		departments = append(departments, r.toBizDepartment(dbDept))
	}
	return departments, nil
}

func (r *departmentRepo) Update(ctx context.Context, dept *biz.Department) (*biz.Department, error) {
	dbDept := &Department{
		ID:           dept.ID,
		Name:         dept.Name,
		ParentID:     dept.ParentID,
		Sort:         dept.Sort,
		LeaderUserID: dept.LeaderUserID,
		Phone:        dept.Phone,
		Email:        dept.Email,
		Status:       dept.Status,
		UpdateAt:     time.Now(),
	}

	_, err := r.data.DB(ctx).NewUpdate().Model(dbDept).WherePK().OmitZero().Exec(ctx)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, dept.ID)
}

func (r *departmentRepo) Delete(ctx context.Context, id string) error {
	_, err := r.data.DB(ctx).NewDelete().
		Model((*Department)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *departmentRepo) HasUsers(ctx context.Context, id string) (bool, error) {
	return false, nil
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

func (r *postRepo) GetByID(ctx context.Context, id string) (*biz.Post, error) {
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

func (r *postRepo) GetByName(ctx context.Context, name string) (*biz.Post, error) {
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

func (r *postRepo) GetByCode(ctx context.Context, code string) (*biz.Post, error) {
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

	return r.GetByID(ctx, post.ID)
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
	err := r.data.DB(ctx).NewSelect().
		Model((*Post)(nil)).
		TableExpr("qa_user_map_post AS up").
		Where("up.post_id = ?", id).
		Count(ctx, &count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *departmentRepo) toBizDepartment(dbDept *Department) *biz.Department {
	return &biz.Department{
		ID:           dbDept.ID,
		Name:         dbDept.Name,
		ParentID:     dbDept.ParentID,
		Sort:         dbDept.Sort,
		LeaderUserID: dbDept.LeaderUserID,
		Phone:        dbDept.Phone,
		Email:        dbDept.Email,
		Status:       dbDept.Status,
		CreateBy:     dbDept.CreateBy,
		CreateAt:     dbDept.CreateAt,
		UpdateBy:     dbDept.UpdateBy,
		UpdateAt:     dbDept.UpdateAt,
		TenantID:     dbDept.TenantID,
		Children:     []*biz.Department{},
	}
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
