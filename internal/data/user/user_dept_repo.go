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

type userDeptRepo struct {
	data *data.Data
	log  *log.Helper
}

func NewUserDeptRepo(data *data.Data, logger log.Logger) biz.UserDeptRepo {
	return &userDeptRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userDeptRepo) Create(ctx context.Context, item *biz.UserDept) error {
	if item == nil {
		return nil
	}
	now := time.Now()
	_, err := r.data.DB(ctx).NewInsert().Model(&UserDept{
		ID:       item.ID,
		UserID:   item.UserID,
		DeptID:   item.DeptID,
		CreateAt: now,
		CreateBy: ctxs.GetLoginID(ctx),
		UpdateAt: now,
		UpdateBy: ctxs.GetLoginID(ctx),
		TenantID: ctxs.GetTenantID(ctx),
	}).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *userDeptRepo) Delete(ctx context.Context, id string) error {
	_, err := r.data.DB(ctx).NewUpdate().
		Model((*UserDept)(nil)).
		Set("update_by = ?", ctxs.GetLoginID(ctx)).
		Where("id = ?", id).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *userDeptRepo) GetUserDepts(ctx context.Context, userID string) ([]*biz.UserDept, error) {
	var userDepts []*UserDept
	err := r.data.DB(ctx).NewSelect().
		Model(&userDepts).
		Where("user_id = ?", userID).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]*biz.UserDept, 0), nil
		}
		return nil, err
	}

	return slices.Map(userDepts, func(item *UserDept, index int) *biz.UserDept {
		return r.toBizUserDept(item)
	}), nil
}

func (r *userDeptRepo) toBizUserDept(item *UserDept) *biz.UserDept {
	return &biz.UserDept{
		ID:       item.ID,
		UserID:   item.UserID,
		DeptID:   item.DeptID,
		CreateBy: item.CreateBy,
		CreateAt: item.CreateAt,
		UpdateBy: item.UpdateBy,
		UpdateAt: item.UpdateAt,
		TenantID: item.TenantID,
	}
}
