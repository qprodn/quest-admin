package config

import (
	"context"

	v1 "quest-admin/api/gen/config/v1"
	biz "quest-admin/internal/biz/config"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ConfigService struct {
	v1.UnimplementedConfigServiceServer
	uc  *biz.ConfigUsecase
	log *log.Helper
}

func NewConfigService(uc *biz.ConfigUsecase, logger log.Logger) *ConfigService {
	return &ConfigService{
		uc:  uc,
		log: log.NewHelper(log.With(logger, "module", "config/service")),
	}
}

func (s *ConfigService) CreateConfig(ctx context.Context, in *v1.CreateConfigRequest) (*emptypb.Empty, error) {
	config := &biz.Config{
		Name:  in.GetName(),
		Key:   in.GetKey(),
		Value: in.GetValue(),
		Status: 1,
	}

	err := s.uc.CreateConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *ConfigService) GetConfig(ctx context.Context, in *v1.GetConfigRequest) (*v1.GetConfigReply, error) {
	config, err := s.uc.GetConfig(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &v1.GetConfigReply{
		Config: s.toProtoConfig(config),
	}, nil
}

func (s *ConfigService) GetConfigByKey(ctx context.Context, in *v1.GetConfigByKeyRequest) (*v1.GetConfigByKeyReply, error) {
	value, err := s.uc.GetConfigValue(ctx, in.GetKey())
	if err != nil {
		return nil, err
	}

	return &v1.GetConfigByKeyReply{
		Value: value,
	}, nil
}

func (s *ConfigService) ListConfigs(ctx context.Context, in *v1.ListConfigsRequest) (*v1.ListConfigsReply, error) {
	query := &biz.ListConfigsQuery{
		Page:      in.GetPage(),
		PageSize:  in.GetPageSize(),
		Name:      in.GetName(),
		Key:       in.GetKey(),
		SortField: in.GetSortField(),
		SortOrder: in.GetSortOrder(),
	}

	if in.Status != nil {
		query.Status = in.Status
	}

	result, err := s.uc.ListConfigs(ctx, query)
	if err != nil {
		return nil, err
	}

	configs := make([]*v1.ConfigInfo, 0, len(result.Configs))
	for _, config := range result.Configs {
		configs = append(configs, s.toProtoConfig(config))
	}

	return &v1.ListConfigsReply{
		Configs:    configs,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

func (s *ConfigService) UpdateConfig(ctx context.Context, in *v1.UpdateConfigRequest) (*emptypb.Empty, error) {
	config := &biz.Config{
		ID:    in.GetId(),
		Name:  in.GetName(),
		Key:   in.GetKey(),
		Value: in.GetValue(),
	}

	err := s.uc.UpdateConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *ConfigService) ChangeConfigStatus(ctx context.Context, in *v1.ChangeConfigStatusRequest) (*emptypb.Empty, error) {
	bo := &biz.UpdateStatusBO{
		ConfigID: in.GetId(),
		Status:   in.GetStatus(),
	}
	err := s.uc.ChangeConfigStatus(ctx, bo)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *ConfigService) DeleteConfig(ctx context.Context, in *v1.DeleteConfigRequest) (*emptypb.Empty, error) {
	err := s.uc.DeleteConfig(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *ConfigService) toProtoConfig(config *biz.Config) *v1.ConfigInfo {
	return &v1.ConfigInfo{
		Id:       config.ID,
		Name:     config.Name,
		Key:      config.Key,
		Value:    config.Value,
		Status:   config.Status,
		CreateAt: timestamppb.New(config.CreateAt),
		UpdateAt: timestamppb.New(config.UpdateAt),
		TenantId: config.TenantID,
	}
}
