package organization

import (
	"context"
	"database/sql"
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
	err := r.data.DB(ctx).NewSelect().Model(dbDept).WherePK().Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, biz.ErrDepartmentNotFound
		}
		return nil, err
	}
	return r.toBizDepartment(dbDept), nil
}

func (r *departmentRepo) FindByName(ctx context.Context, name string) (*biz.Department, error) {
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

func (r *departmentRepo) List(ctx context.Context) ([]*biz.Department, error) {
	var dbDepts []*Department
	err := r.data.DB(ctx).NewSelect().
		Model(&dbDepts).
		Where("status = ?", 1).
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

	return r.FindByID(ctx, dept.ID)
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
