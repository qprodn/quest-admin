package dict

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"quest-admin/internal/data/data"
	"quest-admin/pkg/util/ctxs"
	"time"

	biz "quest-admin/internal/biz/dict"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/uptrace/bun"
)

type DictType struct {
	bun.BaseModel `bun:"table:qc_dict_type,alias:t"`

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

type dictTypeRepo struct {
	data *data.Data
	log  *log.Helper
}

func NewDictTypeRepo(data *data.Data, logger log.Logger) biz.DictTypeRepo {
	return &dictTypeRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *dictTypeRepo) Create(ctx context.Context, dictType *biz.DictType) (*biz.DictType, error) {
	now := time.Now()
	dbDictType := &DictType{
		ID:       dictType.ID,
		Name:     dictType.Name,
		Code:     dictType.Code,
		Sort:     dictType.Sort,
		Status:   dictType.Status,
		Remark:   dictType.Remark,
		CreateBy: ctxs.GetLoginID(ctx),
		CreateAt: now,
		UpdateBy: ctxs.GetLoginID(ctx),
		UpdateAt: now,
		TenantID: ctxs.GetTenantID(ctx),
	}

	_, err := r.data.DB(ctx).NewInsert().Model(dbDictType).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return r.toBizDictType(dbDictType), nil
}

func (r *dictTypeRepo) FindByID(ctx context.Context, id string) (*biz.DictType, error) {
	dbDictType := &DictType{ID: id}
	err := r.data.DB(ctx).
		NewSelect().
		Model(dbDictType).
		WherePK().
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return r.toBizDictType(dbDictType), nil
}

func (r *dictTypeRepo) FindByCode(ctx context.Context, code string) (*biz.DictType, error) {
	dbDictType := &DictType{}
	err := r.data.DB(ctx).
		NewSelect().
		Model(dbDictType).
		Where("code = ?", code).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return r.toBizDictType(dbDictType), nil
}

func (r *dictTypeRepo) List(ctx context.Context, opt *biz.WhereDictTypeOpt) ([]*biz.DictType, error) {
	var dbDictTypes []*DictType
	q := r.data.DB(ctx).NewSelect().Model(&dbDictTypes)

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

	dictTypes := make([]*biz.DictType, 0, len(dbDictTypes))
	for _, dbDictType := range dbDictTypes {
		dictTypes = append(dictTypes, r.toBizDictType(dbDictType))
	}

	return dictTypes, nil
}

func (r *dictTypeRepo) Count(ctx context.Context, opt *biz.WhereDictTypeOpt) (int64, error) {
	var dbDictTypes []*DictType
	q := r.data.DB(ctx).NewSelect().Model(&dbDictTypes)

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

func (r *dictTypeRepo) Update(ctx context.Context, dictType *biz.DictType) (*biz.DictType, error) {
	dbDictType := &DictType{
		ID:       dictType.ID,
		Name:     dictType.Name,
		Code:     dictType.Code,
		Sort:     dictType.Sort,
		Status:   dictType.Status,
		Remark:   dictType.Remark,
		UpdateAt: time.Now(),
	}

	_, err := r.data.DB(ctx).
		NewUpdate().
		Model(dbDictType).
		WherePK().
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		OmitZero().
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return r.FindByID(ctx, dictType.ID)
}

func (r *dictTypeRepo) Delete(ctx context.Context, id string) error {
	_, err := r.data.DB(ctx).NewDelete().
		Model((*DictType)(nil)).
		Where("id = ?", id).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Exec(ctx)
	return err
}

func (r *dictTypeRepo) HasDictData(ctx context.Context, id string) (bool, error) {
	var count int
	count, err := r.data.DB(ctx).NewSelect().
		Model((*DictType)(nil)).
		TableExpr("qc_dict_data AS dd").
		Where("dd.dict_type_id = ?", id).
		Where("dd.tenant_id = ?", ctxs.GetTenantID(ctx)).
		Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *dictTypeRepo) FindByIDs(ctx context.Context, ids []string) ([]*biz.DictType, error) {
	if len(ids) == 0 {
		return []*biz.DictType{}, nil
	}
	var dbDictTypes []*DictType
	err := r.data.DB(ctx).NewSelect().
		Model(&dbDictTypes).
		Where("id IN (?)", bun.In(ids)).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	dictTypes := make([]*biz.DictType, 0, len(dbDictTypes))
	for _, dbDictType := range dbDictTypes {
		dictTypes = append(dictTypes, r.toBizDictType(dbDictType))
	}
	return dictTypes, nil
}

func (r *dictTypeRepo) toBizDictType(dbDictType *DictType) *biz.DictType {
	return &biz.DictType{
		ID:       dbDictType.ID,
		Name:     dbDictType.Name,
		Code:     dbDictType.Code,
		Sort:     dbDictType.Sort,
		Status:   dbDictType.Status,
		Remark:   dbDictType.Remark,
		CreateBy: dbDictType.CreateBy,
		CreateAt: dbDictType.CreateAt,
		UpdateBy: dbDictType.UpdateBy,
		UpdateAt: dbDictType.UpdateAt,
		TenantID: dbDictType.TenantID,
	}
}
