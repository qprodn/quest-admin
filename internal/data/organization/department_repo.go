package organization

import (
	"context"
	"database/sql"
	"quest-admin/internal/data/data"
	"quest-admin/pkg/util/ctxs"
	"time"

	biz "quest-admin/internal/biz/organization"
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

type UserDept struct {
	bun.BaseModel `bun:"table:qa_user_dept,alias:ud"`

	ID       string     `bun:"id,pk"`
	UserID   string     `bun:"user_id,notnull"`
	DeptID   string     `bun:"dept_id,notnull"`
	CreateBy string     `bun:"create_by"`
	CreateAt time.Time  `bun:"create_at,notnull,default:current_timestamp()"`
	UpdateBy string     `bun:"update_by"`
	UpdateAt time.Time  `bun:"update_at,notnull,default:current_timestamp()"`
	TenantID string     `bun:"tenant_id"`
	DeleteAt *time.Time `bun:"delete_at,soft_delete,nullzero"`
}

type departmentRepo struct {
	data *data.Data
	log  *log.Helper
}

func NewDepartmentRepo(data *data.Data, logger log.Logger) biz.DepartmentRepo {
	return &departmentRepo{
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

func (r *departmentRepo) FindByID(ctx context.Context, id string) (*biz.Department, error) {
	dbDept := &Department{ID: id}
	err := r.data.DB(ctx).
		NewSelect().
		Model(dbDept).
		WherePK().
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return r.toBizDepartment(dbDept), nil
}

func (r *departmentRepo) FindByName(ctx context.Context, name string) (*biz.Department, error) {
	dbDept := &Department{}
	err := r.data.DB(ctx).
		NewSelect().
		Model(dbDept).
		Where("name = ?", name).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return r.toBizDepartment(dbDept), nil
}

func (r *departmentRepo) List(ctx context.Context) ([]*biz.Department, error) {
	var dbDepts []*Department
	err := r.data.DB(ctx).NewSelect().
		Model(&dbDepts).
		Where("status = ?", 1).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Order("level ASC, sort ASC").
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

func (r *departmentRepo) FindByParentID(ctx context.Context, parentID string) ([]*biz.Department, error) {
	var dbDepts []*Department
	err := r.data.DB(ctx).NewSelect().
		Model(&dbDepts).
		Where("parent_id = ?", parentID).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
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
		UpdateBy:     dept.UpdateBy,
		UpdateAt:     time.Now(),
	}

	_, err := r.data.DB(ctx).
		NewUpdate().
		Model(dbDept).
		WherePK().
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		OmitZero().
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return r.FindByID(ctx, dept.ID)
}

func (r *departmentRepo) Delete(ctx context.Context, id string) error {
	_, err := r.data.DB(ctx).NewUpdate().
		Model((*Department)(nil)).
		Set("delete_at = ?", time.Now()).
		Set("update_by", ctxs.GetLoginID(ctx)).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *departmentRepo) HasUsers(ctx context.Context, id string) (bool, error) {
	count, err := r.data.DB(ctx).NewSelect().
		Model((*UserDept)(nil)).
		Where("dept_id = ?", id).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Count(ctx)
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
