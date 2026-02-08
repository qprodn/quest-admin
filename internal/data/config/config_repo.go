package config

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"quest-admin/internal/data/data"
	"quest-admin/pkg/util/ctxs"
	"time"

	biz "quest-admin/internal/biz/config"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/uptrace/bun"
)

type Config struct {
	bun.BaseModel `bun:"table:qa_config,alias:c"`

	ID       string     `bun:"id,pk"`
	Name     string     `bun:"name,notnull"`
	Key      string     `bun:"key,notnull"`
	Value    string     `bun:"value"`
	Status   int32      `bun:"status,default:1"`
	CreateBy string     `bun:"create_by"`
	CreateAt time.Time  `bun:"create_at,notnull,default:current_timestamp()"`
	UpdateBy string     `bun:"update_by"`
	UpdateAt time.Time  `bun:"update_at,notnull,default:current_timestamp()"`
	TenantID string     `bun:"tenant_id"`
	DeleteAt *time.Time `bun:"delete_at,soft_delete,nullzero"`
}

type configRepo struct {
	data *data.Data
	log  *log.Helper
}

func NewConfigRepo(data *data.Data, logger log.Logger) biz.ConfigRepo {
	return &configRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *configRepo) Create(ctx context.Context, config *biz.Config) error {
	now := time.Now()
	dbConfig := &Config{
		ID:       config.ID,
		Name:     config.Name,
		Key:      config.Key,
		Value:    config.Value,
		Status:   config.Status,
		CreateBy: config.CreateBy,
		CreateAt: now,
		UpdateBy: config.UpdateBy,
		UpdateAt: now,
		TenantID: config.TenantID,
	}

	_, err := r.data.DB(ctx).NewInsert().Model(dbConfig).Exec(ctx)
	if err != nil {
		r.log.WithContext(ctx).Error(err)
		return err
	}

	return nil
}

func (r *configRepo) FindByID(ctx context.Context, id string) (*biz.Config, error) {
	dbConfig := &Config{ID: id}
	err := r.data.DB(ctx).
		NewSelect().
		Model(dbConfig).
		Where("id = ?", id).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		r.log.WithContext(ctx).Error(err)
		return nil, err
	}
	return r.toBizConfig(dbConfig), nil
}

func (r *configRepo) FindByKey(ctx context.Context, key string) (*biz.Config, error) {
	dbConfig := &Config{}
	err := r.data.DB(ctx).
		NewSelect().
		Model(dbConfig).
		Where("key = ?", key).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		r.log.WithContext(ctx).Error(err)
		return nil, err
	}
	return r.toBizConfig(dbConfig), nil
}

func (r *configRepo) List(ctx context.Context, opt *biz.WhereConfigOpt) ([]*biz.Config, error) {
	var dbConfigs []*Config
	q := r.data.DB(ctx).NewSelect().Model(&dbConfigs)

	if opt.Name != "" {
		q = q.Where("name LIKE ?", "%"+opt.Name+"%")
	}
	if opt.Key != "" {
		q = q.Where("key LIKE ?", "%"+opt.Key+"%")
	}
	if opt.Status != nil {
		q = q.Where("status = ?", *opt.Status)
	}
	if opt.Offset != 0 {
		q = q.Offset(int(opt.Offset))
	}
	if opt.Limit != 0 {
		q = q.Limit(int(opt.Limit))
	}
	if opt.SortField != "" && opt.SortOrder != "" {
		q = q.Order(fmt.Sprintf("%s %s", opt.SortField, opt.SortOrder))
	} else {
		q = q.Order("b.id DESC")
	}

	err := q.Where("tenant_id = ?", ctxs.GetTenantID(ctx)).Scan(ctx)
	if err != nil {
		r.log.WithContext(ctx).Error(err)
		return nil, err
	}

	configs := make([]*biz.Config, 0, len(dbConfigs))
	for _, dbConfig := range dbConfigs {
		configs = append(configs, r.toBizConfig(dbConfig))
	}

	return configs, nil
}

func (r *configRepo) Count(ctx context.Context, opt *biz.WhereConfigOpt) (int64, error) {
	var dbConfigs []*Config
	q := r.data.DB(ctx).NewSelect().Model(&dbConfigs)

	if opt.Name != "" {
		q = q.Where("name LIKE ?", "%"+opt.Name+"%")
	}
	if opt.Key != "" {
		q = q.Where("key LIKE ?", "%"+opt.Key+"%")
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

func (r *configRepo) Update(ctx context.Context, config *biz.Config) error {
	dbConfig := &Config{
		ID:       config.ID,
		Name:     config.Name,
		Key:      config.Key,
		Value:    config.Value,
		UpdateAt: time.Now(),
	}

	_, err := r.data.DB(ctx).
		NewUpdate().
		Model(dbConfig).
		WherePK().
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		OmitZero().
		Exec(ctx)
	if err != nil {
		r.log.WithContext(ctx).Error(err)
		return err
	}

	return nil
}

func (r *configRepo) UpdateStatus(ctx context.Context, bo *biz.UpdateStatusBO) error {
	_, err := r.data.DB(ctx).NewUpdate().
		Model((*Config)(nil)).
		Set("status = ?", bo.Status).
		Set("update_at = ?", time.Now()).
		Where("id = ?", bo.ConfigID).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Exec(ctx)
	return err
}

func (r *configRepo) Delete(ctx context.Context, bo *biz.DeleteConfigBO) error {
	_, err := r.data.DB(ctx).NewUpdate().
		Model((*Config)(nil)).
		Set("update_by = ?", ctxs.GetLoginID(ctx)).
		Set("update_at = ?", time.Now()).
		Set("delete_at = ?", time.Now()).
		Where("id = ?", bo.ConfigID).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Exec(ctx)
	return err
}

func (r *configRepo) toBizConfig(dbConfig *Config) *biz.Config {
	return &biz.Config{
		ID:       dbConfig.ID,
		Name:     dbConfig.Name,
		Key:      dbConfig.Key,
		Value:    dbConfig.Value,
		Status:   dbConfig.Status,
		CreateBy: dbConfig.CreateBy,
		CreateAt: dbConfig.CreateAt,
		UpdateBy: dbConfig.UpdateBy,
		UpdateAt: dbConfig.UpdateAt,
		TenantID: dbConfig.TenantID,
	}
}
