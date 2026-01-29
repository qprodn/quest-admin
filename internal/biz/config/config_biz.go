package config

import (
	"context"
	"quest-admin/internal/data/idgen"
	"quest-admin/pkg/util/pagination"
	"quest-admin/types/consts/id"

	"github.com/go-kratos/kratos/v2/log"
)

type ConfigRepo interface {
	Create(ctx context.Context, config *Config) error
	FindByID(ctx context.Context, id string) (*Config, error)
	FindByKey(ctx context.Context, key string) (*Config, error)
	List(ctx context.Context, opt *WhereConfigOpt) ([]*Config, error)
	Count(ctx context.Context, opt *WhereConfigOpt) (int64, error)
	Update(ctx context.Context, config *Config) error
	UpdateStatus(ctx context.Context, bo *UpdateStatusBO) error
	Delete(ctx context.Context, bo *DeleteConfigBO) error
}

type ConfigUsecase struct {
	idgen      *idgen.IDGenerator
	configRepo ConfigRepo
	log        *log.Helper
}

func NewConfigUsecase(
	logger log.Logger,
	repo ConfigRepo,
	idgen *idgen.IDGenerator,
) *ConfigUsecase {
	uc := &ConfigUsecase{
		log:        log.NewHelper(log.With(logger, "module", "config/biz/config")),
		idgen:      idgen,
		configRepo: repo,
	}
	return uc
}

func (uc *ConfigUsecase) CreateConfig(ctx context.Context, config *Config) error {
	// 检查 key 是否已存在
	existing, err := uc.configRepo.FindByKey(ctx, config.Key)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("查询配置失败,key:%s,error:%v", config.Key, err)
		return err
	}
	if existing != nil {
		uc.log.WithContext(ctx).Error("已存在相同配置键,key:%s", config.Key)
		return ErrConfigExists
	}

	config.ID = uc.idgen.NextID(id.CONFIG)
	config.Status = 1 // 默认启用

	return uc.configRepo.Create(ctx, config)
}

func (uc *ConfigUsecase) GetConfig(ctx context.Context, id string) (*Config, error) {
	return uc.configRepo.FindByID(ctx, id)
}

func (uc *ConfigUsecase) GetConfigByKey(ctx context.Context, key string) (*Config, error) {
	config, err := uc.configRepo.FindByKey(ctx, key)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("查询配置失败,key:%s,error:%v", key, err)
		return nil, err
	}
	if config == nil {
		return nil, ErrConfigNotFound
	}
	if config.Status != 1 {
		return nil, ErrConfigDisabled
	}
	return config, nil
}

func (uc *ConfigUsecase) GetConfigValue(ctx context.Context, key string) (string, error) {
	config, err := uc.GetConfigByKey(ctx, key)
	if err != nil {
		return "", err
	}
	return config.Value, nil
}

func (uc *ConfigUsecase) ListConfigs(ctx context.Context, query *ListConfigsQuery) (*ListConfigsResult, error) {
	opt := &WhereConfigOpt{
		Limit:     query.PageSize,
		Offset:    pagination.GetOffset(query.Page, query.PageSize),
		Name:      query.Name,
		Key:       query.Key,
		Status:    query.Status,
		SortField: query.SortField,
		SortOrder: query.SortOrder,
	}

	list, err := uc.configRepo.List(ctx, opt)
	if err != nil {
		uc.log.WithContext(ctx).Error("查询配置列表失败", err)
		return nil, err
	}

	total, err := uc.configRepo.Count(ctx, opt)
	if err != nil {
		uc.log.WithContext(ctx).Error("查询配置列表总数失败", err)
		return nil, err
	}

	return &ListConfigsResult{
		Configs:    list,
		Total:      total,
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: pagination.GetTotalPages(total, int64(query.PageSize)),
	}, nil
}

func (uc *ConfigUsecase) UpdateConfig(ctx context.Context, config *Config) error {
	// 检查配置是否存在
	_, err := uc.configRepo.FindByID(ctx, config.ID)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("查询配置失败,configID:%s,error:%v", config.ID, err)
		return err
	}

	// 如果修改了 key，检查新 key 是否已存在
	if config.Key != "" {
		existing, err := uc.configRepo.FindByKey(ctx, config.Key)
		if err != nil {
			uc.log.WithContext(ctx).Errorf("查询配置失败,key:%s,error:%v", config.Key, err)
			return err
		}
		if existing != nil && existing.ID != config.ID {
			uc.log.WithContext(ctx).Error("已存在相同配置键,key:%s", config.Key)
			return ErrConfigExists
		}
	}

	return uc.configRepo.Update(ctx, config)
}

func (uc *ConfigUsecase) ChangeConfigStatus(ctx context.Context, bo *UpdateStatusBO) error {
	_, err := uc.configRepo.FindByID(ctx, bo.ConfigID)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("查询配置失败,configID:%s,error:%v", bo.ConfigID, err)
		return err
	}

	return uc.configRepo.UpdateStatus(ctx, bo)
}

func (uc *ConfigUsecase) DeleteConfig(ctx context.Context, id string) error {
	_, err := uc.configRepo.FindByID(ctx, id)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("查询配置失败,configID:%s,error:%v", id, err)
		return err
	}

	return uc.configRepo.Delete(ctx, &DeleteConfigBO{
		ConfigID: id,
	})
}
