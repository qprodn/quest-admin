package dict

import (
	"context"
	"database/sql"
	"fmt"
	"quest-admin/internal/data/data"
	"quest-admin/pkg/util/ctxs"
	"time"

	biz "quest-admin/internal/biz/dict"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/uptrace/bun"
)

type DictData struct {
	bun.BaseModel `bun:"table:qc_dict_data,alias:d"`

	ID         string     `bun:"id,pk"`
	DictTypeID string     `bun:"dict_type_id,notnull"`
	Label      string     `bun:"label,notnull"`
	Value      string     `bun:"value,notnull"`
	Sort       int32      `bun:"sort,notnull"`
	Status     int32      `bun:"status,notnull"`
	CSSClass   string     `bun:"css_class"`
	IsDefault  bool       `bun:"is_default,notnull"`
	Remark     string     `bun:"remark"`
	CreateBy   string     `bun:"create_by"`
	CreateAt   time.Time  `bun:"create_at,notnull,default:current_timestamp()"`
	UpdateBy   string     `bun:"update_by"`
	UpdateAt   time.Time  `bun:"update_at,notnull,default:current_timestamp()"`
	TenantID   string     `bun:"tenant_id,notnull"`
	DeleteAt   *time.Time `bun:"delete_at,soft_delete,nullzero"`
}

type dictDataRepo struct {
	data *data.Data
	log  *log.Helper
}

func NewDictDataRepo(data *data.Data, logger log.Logger) biz.DictDataRepo {
	return &dictDataRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *dictDataRepo) Create(ctx context.Context, dictData *biz.DictData) (*biz.DictData, error) {
	now := time.Now()
	dbDictData := &DictData{
		ID:         dictData.ID,
		DictTypeID: dictData.DictTypeID,
		Label:      dictData.Label,
		Value:      dictData.Value,
		Sort:       dictData.Sort,
		Status:     dictData.Status,
		CSSClass:   dictData.CSSClass,
		IsDefault:  dictData.IsDefault,
		Remark:     dictData.Remark,
		CreateBy:   ctxs.GetLoginID(ctx),
		CreateAt:   now,
		UpdateBy:   ctxs.GetLoginID(ctx),
		UpdateAt:   now,
		TenantID:   ctxs.GetTenantID(ctx),
	}

	_, err := r.data.DB(ctx).NewInsert().Model(dbDictData).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return r.toBizDictData(dbDictData), nil
}

func (r *dictDataRepo) FindByID(ctx context.Context, id string) (*biz.DictData, error) {
	dbDictData := &DictData{ID: id}
	err := r.data.DB(ctx).
		NewSelect().
		Model(dbDictData).
		WherePK().
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return r.toBizDictData(dbDictData), nil
}

func (r *dictDataRepo) FindByValue(ctx context.Context, dictTypeID, value string) (*biz.DictData, error) {
	dbDictData := &DictData{}
	err := r.data.DB(ctx).
		NewSelect().
		Model(dbDictData).
		Where("dict_type_id = ?", dictTypeID).
		Where("value = ?", value).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return r.toBizDictData(dbDictData), nil
}

func (r *dictDataRepo) List(ctx context.Context, query *biz.ListDictDataQuery) (*biz.ListDictDataResult, error) {
	var dbDictData []*DictData
	q := r.data.DB(ctx).NewSelect().Model(&dbDictData)

	if query.DictTypeID != "" {
		q = q.Where("dict_type_id = ?", query.DictTypeID)
	}

	if query.Keyword != "" {
		q = q.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.WhereOr("label LIKE ?", "%"+query.Keyword+"%").
				WhereOr("value LIKE ?", "%"+query.Keyword+"%")
		})
	}

	if query.Status != nil {
		q = q.Where("status = ?", *query.Status)
	}

	q = q.Where("tenant_id = ?", ctxs.GetTenantID(ctx))

	total, err := q.ScanAndCount(ctx, &dbDictData, nil)
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

	dictData := make([]*biz.DictData, 0, len(dbDictData))
	for _, dbDictData := range dbDictData {
		dictData = append(dictData, r.toBizDictData(dbDictData))
	}

	return &biz.ListDictDataResult{
		DictData:   dictData,
		Total:      int64(total),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *dictDataRepo) Update(ctx context.Context, dictData *biz.DictData) (*biz.DictData, error) {
	dbDictData := &DictData{
		ID:         dictData.ID,
		DictTypeID: dictData.DictTypeID,
		Label:      dictData.Label,
		Value:      dictData.Value,
		Sort:       dictData.Sort,
		Status:     dictData.Status,
		CSSClass:   dictData.CSSClass,
		IsDefault:  dictData.IsDefault,
		Remark:     dictData.Remark,
		UpdateAt:   time.Now(),
	}

	_, err := r.data.DB(ctx).
		NewUpdate().
		Model(dbDictData).
		WherePK().
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		OmitZero().
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return r.FindByID(ctx, dictData.ID)
}

func (r *dictDataRepo) Delete(ctx context.Context, id string) error {
	_, err := r.data.DB(ctx).NewDelete().
		Model((*DictData)(nil)).
		Where("id = ?", id).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Exec(ctx)
	return err
}

func (r *dictDataRepo) FindListByIDs(ctx context.Context, ids []string) ([]*biz.DictData, error) {
	if len(ids) == 0 {
		return []*biz.DictData{}, nil
	}
	var dbDictData []*DictData
	err := r.data.DB(ctx).NewSelect().
		Model(&dbDictData).
		Where("id IN (?)", bun.In(ids)).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	dictData := make([]*biz.DictData, 0, len(dbDictData))
	for _, dbDictData := range dbDictData {
		dictData = append(dictData, r.toBizDictData(dbDictData))
	}
	return dictData, nil
}

func (r *dictDataRepo) FindByDictTypeID(ctx context.Context, dictTypeID string) ([]*biz.DictData, error) {
	var dbDictData []*DictData
	err := r.data.DB(ctx).NewSelect().
		Model(&dbDictData).
		Where("dict_type_id = ?", dictTypeID).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Order("sort ASC, create_at DESC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	dictData := make([]*biz.DictData, 0, len(dbDictData))
	for _, dbDictData := range dbDictData {
		dictData = append(dictData, r.toBizDictData(dbDictData))
	}
	return dictData, nil
}

func (r *dictDataRepo) toBizDictData(dbDictData *DictData) *biz.DictData {
	return &biz.DictData{
		ID:         dbDictData.ID,
		DictTypeID: dbDictData.DictTypeID,
		Label:      dbDictData.Label,
		Value:      dbDictData.Value,
		Sort:       dbDictData.Sort,
		Status:     dbDictData.Status,
		CSSClass:   dbDictData.CSSClass,
		IsDefault:  dbDictData.IsDefault,
		Remark:     dbDictData.Remark,
		CreateBy:   dbDictData.CreateBy,
		CreateAt:   dbDictData.CreateAt,
		UpdateBy:   dbDictData.UpdateBy,
		UpdateAt:   dbDictData.UpdateAt,
		TenantID:   dbDictData.TenantID,
	}
}
