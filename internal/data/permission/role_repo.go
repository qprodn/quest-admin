package permission

import (
	"context"
	"database/sql"
	"fmt"
	"quest-admin/internal/data/data"
	"quest-admin/pkg/lang/slices"
	"quest-admin/pkg/util/ctxs"
	"time"

	biz "quest-admin/internal/biz/permission"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/uptrace/bun"
)

type Role struct {
	bun.BaseModel `bun:"table:qa_role,alias:r"`

	ID               string     `bun:"id,pk"`
	Name             string     `bun:"name,notnull"`
	Code             string     `bun:"code,notnull"`
	Sort             int32      `bun:"sort,notnull"`
	DataScope        int32      `bun:"data_scope,default:1"`
	DataScopeDeptIDs string     `bun:"data_scope_dept_ids,default:''"`
	Status           int32      `bun:"status,notnull"`
	Type             int32      `bun:"type,notnull"`
	Remark           string     `bun:"remark"`
	CreateBy         string     `bun:"create_by"`
	CreateAt         time.Time  `bun:"create_at,notnull,default:current_timestamp()"`
	UpdateBy         string     `bun:"update_by"`
	UpdateAt         time.Time  `bun:"update_at,notnull,default:current_timestamp()"`
	TenantID         string     `bun:"tenant_id,notnull"`
	DeleteAt         *time.Time `bun:"delete_at,soft_delete,nullzero"`
}

type roleRepo struct {
	data *data.Data
	log  *log.Helper
}

func NewRoleRepo(data *data.Data, logger log.Logger) biz.RoleRepo {
	return &roleRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *roleRepo) FindListByIDs(ctx context.Context, roleIds []string) ([]*biz.Role, error) {
	if len(roleIds) == 0 {
		return []*biz.Role{}, nil
	}
	var dbRoles []*Role
	err := r.data.DB(ctx).NewSelect().
		Model(&dbRoles).
		Where("id in (?)", bun.In(roleIds)).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []*biz.Role{}, nil
		}
		r.log.WithContext(ctx).Error(err)
		return nil, err
	}
	roles := slices.Map(dbRoles, func(item *Role, index int) *biz.Role {
		return r.toBizRole(item)
	})
	return roles, nil
}

func (r *roleRepo) Create(ctx context.Context, role *biz.Role) (*biz.Role, error) {
	now := time.Now()
	dbRole := &Role{
		ID:               role.ID,
		Name:             role.Name,
		Code:             role.Code,
		Sort:             role.Sort,
		DataScope:        role.DataScope,
		DataScopeDeptIDs: role.DataScopeDeptIDs,
		Status:           role.Status,
		Type:             role.Type,
		Remark:           role.Remark,
		CreateBy:         ctxs.GetLoginID(ctx),
		CreateAt:         now,
		UpdateBy:         ctxs.GetLoginID(ctx),
		UpdateAt:         now,
		TenantID:         ctxs.GetTenantID(ctx),
	}

	_, err := r.data.DB(ctx).NewInsert().Model(dbRole).Exec(ctx)
	if err != nil {
		r.log.WithContext(ctx).Error(err)
		return nil, err
	}

	return r.toBizRole(dbRole), nil
}

func (r *roleRepo) FindByID(ctx context.Context, id string) (*biz.Role, error) {
	dbRole := &Role{ID: id}
	err := r.data.DB(ctx).
		NewSelect().
		Model(dbRole).
		WherePK().
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		r.log.WithContext(ctx).Error(err)
		return nil, err
	}
	return r.toBizRole(dbRole), nil
}

func (r *roleRepo) FindByName(ctx context.Context, name string) (*biz.Role, error) {
	dbRole := &Role{}
	err := r.data.DB(ctx).
		NewSelect().
		Model(dbRole).
		Where("name = ?", name).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		r.log.WithContext(ctx).Error(err)
		return nil, err
	}
	return r.toBizRole(dbRole), nil
}

func (r *roleRepo) FindByCode(ctx context.Context, code string) (*biz.Role, error) {
	dbRole := &Role{}
	err := r.data.DB(ctx).
		NewSelect().
		Model(dbRole).
		Where("code = ?", code).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		r.log.WithContext(ctx).Error(err)
		return nil, err
	}
	return r.toBizRole(dbRole), nil
}

func (r *roleRepo) List(ctx context.Context, opt *biz.WhereRoleOpt) ([]*biz.Role, error) {
	var dbRoles []*Role
	q := r.data.DB(ctx).NewSelect().Model(&dbRoles)

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
		r.log.WithContext(ctx).Error(err)
		return nil, err
	}

	roles := make([]*biz.Role, 0, len(dbRoles))
	for _, dbRole := range dbRoles {
		roles = append(roles, r.toBizRole(dbRole))
	}

	return roles, nil
}

func (r *roleRepo) Count(ctx context.Context, opt *biz.WhereRoleOpt) (int64, error) {
	var dbRoles []*Role
	q := r.data.DB(ctx).NewSelect().Model(&dbRoles)

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
		r.log.WithContext(ctx).Error(err)
		return 0, err
	}

	return int64(total), nil
}

func (r *roleRepo) Update(ctx context.Context, role *biz.Role) (*biz.Role, error) {
	dbRole := &Role{
		ID:               role.ID,
		Name:             role.Name,
		Code:             role.Code,
		Sort:             role.Sort,
		DataScope:        role.DataScope,
		DataScopeDeptIDs: role.DataScopeDeptIDs,
		Status:           role.Status,
		Type:             role.Type,
		Remark:           role.Remark,
		UpdateAt:         time.Now(),
	}

	_, err := r.data.DB(ctx).
		NewUpdate().
		Model(dbRole).
		WherePK().
		OmitZero().
		Exec(ctx)
	if err != nil {
		r.log.WithContext(ctx).Error(err)
		return nil, err
	}

	return r.FindByID(ctx, role.ID)
}

func (r *roleRepo) Delete(ctx context.Context, id string) error {
	_, err := r.data.DB(ctx).NewUpdate().
		Model((*Role)(nil)).
		Set("delete_at", time.Now()).
		Set("update_by", ctxs.GetLoginID(ctx)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *roleRepo) HasUsers(ctx context.Context, id string) (bool, error) {
	var count int
	count, err := r.data.DB(ctx).NewSelect().
		Model((*Role)(nil)).
		TableExpr("qa_user_map_role AS ur").
		Where("ur.role_id = ?", id).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Count(ctx)
	if err != nil {
		r.log.WithContext(ctx).Error(err)
		return false, err
	}
	return count > 0, nil
}

func (r *roleRepo) toBizRole(dbRole *Role) *biz.Role {
	return &biz.Role{
		ID:               dbRole.ID,
		Name:             dbRole.Name,
		Code:             dbRole.Code,
		Sort:             dbRole.Sort,
		DataScope:        dbRole.DataScope,
		DataScopeDeptIDs: dbRole.DataScopeDeptIDs,
		Status:           dbRole.Status,
		Type:             dbRole.Type,
		Remark:           dbRole.Remark,
		CreateBy:         dbRole.CreateBy,
		CreateAt:         dbRole.CreateAt,
		UpdateBy:         dbRole.UpdateBy,
		UpdateAt:         dbRole.UpdateAt,
		TenantID:         dbRole.TenantID,
	}
}
