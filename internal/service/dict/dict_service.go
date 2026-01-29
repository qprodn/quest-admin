package dict

import (
	"context"

	v1 "quest-admin/api/gen/dict/v1"
	biz "quest-admin/internal/biz/dict"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DictTypeService struct {
	v1.UnimplementedDictTypeServiceServer
	dtuc *biz.DictTypeUsecase
	log  *log.Helper
}

func NewDictTypeService(dtuc *biz.DictTypeUsecase, logger log.Logger) *DictTypeService {
	return &DictTypeService{
		dtuc: dtuc,
		log:  log.NewHelper(log.With(logger, "module", "dict/service/dict_type")),
	}
}

func (s *DictTypeService) CreateDictType(ctx context.Context, in *v1.CreateDictTypeRequest) (*emptypb.Empty, error) {
	dictType := &biz.DictType{
		Name:   in.GetName(),
		Code:   in.GetCode(),
		Sort:   in.GetSort(),
		Remark: in.GetRemark(),
		Status: 1,
	}

	_, err := s.dtuc.CreateDictType(ctx, dictType)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *DictTypeService) GetDictType(ctx context.Context, in *v1.GetDictTypeRequest) (*v1.GetDictTypeReply, error) {
	dictType, err := s.dtuc.GetDictType(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &v1.GetDictTypeReply{
		DictType: s.toProtoDictType(dictType),
	}, nil
}

func (s *DictTypeService) ListDictTypes(ctx context.Context, in *v1.ListDictTypesRequest) (*v1.ListDictTypesReply, error) {
	query := &biz.ListDictTypesQuery{
		Page:      in.GetPage(),
		PageSize:  in.GetPageSize(),
		Keyword:   in.GetKeyword(),
		SortField: in.GetSortField(),
		SortOrder: in.GetSortOrder(),
	}

	if in.Status != nil {
		query.Status = in.Status
	}

	result, err := s.dtuc.ListDictTypes(ctx, query)
	if err != nil {
		return nil, err
	}

	dictTypes := make([]*v1.DictTypeInfo, 0, len(result.DictTypes))
	for _, dictType := range result.DictTypes {
		dictTypes = append(dictTypes, s.toProtoDictType(dictType))
	}

	return &v1.ListDictTypesReply{
		DictTypes:  dictTypes,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

func (s *DictTypeService) UpdateDictType(ctx context.Context, in *v1.UpdateDictTypeRequest) (*emptypb.Empty, error) {
	dictType := &biz.DictType{
		ID:     in.GetId(),
		Name:   in.GetName(),
		Code:   in.GetCode(),
		Sort:   in.GetSort(),
		Status: in.GetStatus(),
		Remark: in.GetRemark(),
	}

	_, err := s.dtuc.UpdateDictType(ctx, dictType)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *DictTypeService) DeleteDictType(ctx context.Context, in *v1.DeleteDictTypeRequest) (*emptypb.Empty, error) {
	err := s.dtuc.DeleteDictType(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *DictTypeService) toProtoDictType(dictType *biz.DictType) *v1.DictTypeInfo {
	return &v1.DictTypeInfo{
		Id:       dictType.ID,
		Name:     dictType.Name,
		Code:     dictType.Code,
		Sort:     dictType.Sort,
		Status:   dictType.Status,
		Remark:   dictType.Remark,
		CreateAt: timestamppb.New(dictType.CreateAt),
		UpdateAt: timestamppb.New(dictType.UpdateAt),
		TenantId: dictType.TenantID,
	}
}

type DictDataService struct {
	v1.UnimplementedDictDataServiceServer
	dduc *biz.DictDataUsecase
	log  *log.Helper
}

func NewDictDataService(dduc *biz.DictDataUsecase, logger log.Logger) *DictDataService {
	return &DictDataService{
		dduc: dduc,
		log:  log.NewHelper(log.With(logger, "module", "dict/service/dict_data")),
	}
}

func (s *DictDataService) CreateDictData(ctx context.Context, in *v1.CreateDictDataRequest) (*emptypb.Empty, error) {
	dictData := &biz.DictData{
		DictTypeID: in.GetDictTypeId(),
		Label:      in.GetLabel(),
		Value:      in.GetValue(),
		Sort:       in.GetSort(),
		CSSClass:   in.GetCssClass(),
		IsDefault:  in.GetIsDefault(),
		Remark:     in.GetRemark(),
		Status:     1,
	}

	_, err := s.dduc.CreateDictData(ctx, dictData)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *DictDataService) GetDictData(ctx context.Context, in *v1.GetDictDataRequest) (*v1.GetDictDataReply, error) {
	dictData, err := s.dduc.GetDictData(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &v1.GetDictDataReply{
		DictData: s.toProtoDictData(dictData),
	}, nil
}

func (s *DictDataService) ListDictData(ctx context.Context, in *v1.ListDictDataRequest) (*v1.ListDictDataReply, error) {
	query := &biz.ListDictDataQuery{
		Page:       in.GetPage(),
		PageSize:   in.GetPageSize(),
		DictTypeID: in.GetDictTypeId(),
		Keyword:    in.GetKeyword(),
		SortField:  in.GetSortField(),
		SortOrder:  in.GetSortOrder(),
	}

	if in.Status != nil {
		query.Status = in.Status
	}

	result, err := s.dduc.ListDictData(ctx, query)
	if err != nil {
		return nil, err
	}

	dictDataList := make([]*v1.DictDataInfo, 0, len(result.DictData))
	for _, dictData := range result.DictData {
		dictDataList = append(dictDataList, s.toProtoDictData(dictData))
	}

	return &v1.ListDictDataReply{
		DictDataList: dictDataList,
		Total:        result.Total,
		Page:         result.Page,
		PageSize:     result.PageSize,
		TotalPages:   result.TotalPages,
	}, nil
}

func (s *DictDataService) UpdateDictData(ctx context.Context, in *v1.UpdateDictDataRequest) (*emptypb.Empty, error) {
	dictData := &biz.DictData{
		ID:         in.GetId(),
		DictTypeID: in.GetDictTypeId(),
		Label:      in.GetLabel(),
		Value:      in.GetValue(),
		Sort:       in.GetSort(),
		Status:     in.GetStatus(),
		CSSClass:   in.GetCssClass(),
		IsDefault:  in.GetIsDefault(),
		Remark:     in.GetRemark(),
	}

	_, err := s.dduc.UpdateDictData(ctx, dictData)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *DictDataService) DeleteDictData(ctx context.Context, in *v1.DeleteDictDataRequest) (*emptypb.Empty, error) {
	err := s.dduc.DeleteDictData(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *DictDataService) toProtoDictData(dictData *biz.DictData) *v1.DictDataInfo {
	return &v1.DictDataInfo{
		Id:         dictData.ID,
		DictTypeId: dictData.DictTypeID,
		Label:      dictData.Label,
		Value:      dictData.Value,
		Sort:       dictData.Sort,
		Status:     dictData.Status,
		CssClass:   dictData.CSSClass,
		IsDefault:  dictData.IsDefault,
		Remark:     dictData.Remark,
		CreateAt:   timestamppb.New(dictData.CreateAt),
		UpdateAt:   timestamppb.New(dictData.UpdateAt),
		TenantId:   dictData.TenantID,
	}
}
