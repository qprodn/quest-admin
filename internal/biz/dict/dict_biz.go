package dict

import (
	"context"
	"quest-admin/internal/data/idgen"
	"quest-admin/pkg/errorx"
	"quest-admin/pkg/util/pagination"
	"quest-admin/types/consts/id"
	"quest-admin/types/errkey"

	"github.com/go-kratos/kratos/v2/log"
)

type DictTypeRepo interface {
	Create(ctx context.Context, dictType *DictType) (*DictType, error)
	FindByID(ctx context.Context, id string) (*DictType, error)
	FindByCode(ctx context.Context, code string) (*DictType, error)
	List(ctx context.Context, opt *WhereDictTypeOpt) ([]*DictType, error)
	Count(ctx context.Context, opt *WhereDictTypeOpt) (int64, error)
	Update(ctx context.Context, dictType *DictType) (*DictType, error)
	Delete(ctx context.Context, id string) error
	HasDictData(ctx context.Context, id string) (bool, error)
	FindByIDs(ctx context.Context, ids []string) ([]*DictType, error)
}

type DictTypeUsecase struct {
	idgen *idgen.IDGenerator
	repo  DictTypeRepo
	log   *log.Helper
}

func NewDictTypeUsecase(
	idgen *idgen.IDGenerator,
	repo DictTypeRepo,
	logger log.Logger,
) *DictTypeUsecase {
	return &DictTypeUsecase{
		idgen: idgen,
		repo:  repo,
		log:   log.NewHelper(log.With(logger, "module", "dict/biz/dict_type")),
	}
}

func (uc *DictTypeUsecase) CreateDictType(ctx context.Context, dictType *DictType) (*DictType, error) {
	uc.log.WithContext(ctx).Infof("CreateDictType: name=%s, code=%s", dictType.Name, dictType.Code)

	existing, err := uc.repo.FindByCode(ctx, dictType.Code)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errorx.Err(errkey.ErrDictTypeCodeExists)
	}

	dictType.ID = uc.idgen.NextID(id.DICT)

	return uc.repo.Create(ctx, dictType)
}

func (uc *DictTypeUsecase) GetDictType(ctx context.Context, id string) (*DictType, error) {
	uc.log.WithContext(ctx).Infof("GetDictType: id=%s", id)
	return uc.repo.FindByID(ctx, id)
}

func (uc *DictTypeUsecase) ListDictTypes(ctx context.Context, query *ListDictTypesQuery) (*ListDictTypesResult, error) {
	uc.log.WithContext(ctx).Infof("ListDictTypes: page=%d, pageSize=%d, keyword=%s", query.Page, query.PageSize, query.Keyword)

	opt := &WhereDictTypeOpt{
		Limit:     query.PageSize,
		Offset:    pagination.GetOffset(query.Page, query.PageSize),
		Keyword:   query.Keyword,
		Status:    query.Status,
		SortField: query.SortField,
		SortOrder: query.SortOrder,
	}

	list, err := uc.repo.List(ctx, opt)
	if err != nil {
		uc.log.WithContext(ctx).Error("查询字典类型列表失败", err)
		return nil, err
	}

	total, err := uc.repo.Count(ctx, opt)
	if err != nil {
		uc.log.WithContext(ctx).Error("查询字典类型列表总数失败", err)
		return nil, err
	}

	return &ListDictTypesResult{
		DictTypes:  list,
		Total:      total,
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: pagination.GetTotalPages(total, int64(query.PageSize)),
	}, nil
}

func (uc *DictTypeUsecase) UpdateDictType(ctx context.Context, dictType *DictType) (*DictType, error) {
	uc.log.WithContext(ctx).Infof("UpdateDictType: id=%s, name=%s", dictType.ID, dictType.Name)

	_, err := uc.repo.FindByID(ctx, dictType.ID)
	if err != nil {
		return nil, err
	}

	return uc.repo.Update(ctx, dictType)
}

func (uc *DictTypeUsecase) DeleteDictType(ctx context.Context, id string) error {
	uc.log.WithContext(ctx).Infof("DeleteDictType: id=%s", id)

	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	hasData, err := uc.repo.HasDictData(ctx, id)
	if err != nil {
		return err
	}
	if hasData {
		return errorx.Err(errkey.ErrDictTypeHasData)
	}

	return uc.repo.Delete(ctx, id)
}

func (uc *DictTypeUsecase) ListByDictTypeIDs(ctx context.Context, dictTypeIds []string) ([]*DictType, error) {
	var res []*DictType
	if len(dictTypeIds) == 0 {
		return res, nil
	}
	dictTypes, err := uc.repo.FindByIDs(ctx, dictTypeIds)
	if err != nil {
		return res, err
	}
	return dictTypes, nil
}

type DictDataRepo interface {
	Create(ctx context.Context, dictData *DictData) (*DictData, error)
	FindByID(ctx context.Context, id string) (*DictData, error)
	FindByValue(ctx context.Context, dictTypeID, value string) (*DictData, error)
	List(ctx context.Context, opt *WhereDictDataOpt) ([]*DictData, error)
	Count(ctx context.Context, opt *WhereDictDataOpt) (int64, error)
	Update(ctx context.Context, dictData *DictData) (*DictData, error)
	Delete(ctx context.Context, id string) error
	FindListByIDs(ctx context.Context, ids []string) ([]*DictData, error)
	FindByDictTypeID(ctx context.Context, dictTypeID string) ([]*DictData, error)
}

type DictDataUsecase struct {
	idgen *idgen.IDGenerator
	repo  DictDataRepo
	log   *log.Helper
}

func NewDictDataUsecase(
	idgen *idgen.IDGenerator,
	repo DictDataRepo,
	logger log.Logger,
) *DictDataUsecase {
	return &DictDataUsecase{
		idgen: idgen,
		repo:  repo,
		log:   log.NewHelper(log.With(logger, "module", "dict/biz/dict_data")),
	}
}

func (uc *DictDataUsecase) CreateDictData(ctx context.Context, dictData *DictData) (*DictData, error) {
	uc.log.WithContext(ctx).Infof("CreateDictData: dictTypeID=%s, label=%s, value=%s", dictData.DictTypeID, dictData.Label, dictData.Value)

	existing, err := uc.repo.FindByValue(ctx, dictData.DictTypeID, dictData.Value)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errorx.Err(errkey.ErrDictDataValueExists)
	}

	dictData.ID = uc.idgen.NextID(id.DICT)

	return uc.repo.Create(ctx, dictData)
}

func (uc *DictDataUsecase) GetDictData(ctx context.Context, id string) (*DictData, error) {
	uc.log.WithContext(ctx).Infof("GetDictData: id=%s", id)
	return uc.repo.FindByID(ctx, id)
}

func (uc *DictDataUsecase) ListDictData(ctx context.Context, query *ListDictDataQuery) (*ListDictDataResult, error) {
	uc.log.WithContext(ctx).Infof("ListDictData: page=%d, pageSize=%d, dictTypeID=%s", query.Page, query.PageSize, query.DictTypeID)

	opt := &WhereDictDataOpt{
		Limit:      query.PageSize,
		Offset:     pagination.GetOffset(query.Page, query.PageSize),
		DictTypeID: query.DictTypeID,
		Keyword:    query.Keyword,
		Status:     query.Status,
		SortField:  query.SortField,
		SortOrder:  query.SortOrder,
	}

	list, err := uc.repo.List(ctx, opt)
	if err != nil {
		uc.log.WithContext(ctx).Error("查询字典数据列表失败", err)
		return nil, err
	}

	total, err := uc.repo.Count(ctx, opt)
	if err != nil {
		uc.log.WithContext(ctx).Error("查询字典数据列表总数失败", err)
		return nil, err
	}

	return &ListDictDataResult{
		DictData:   list,
		Total:      total,
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: pagination.GetTotalPages(total, int64(query.PageSize)),
	}, nil
}

func (uc *DictDataUsecase) UpdateDictData(ctx context.Context, dictData *DictData) (*DictData, error) {
	uc.log.WithContext(ctx).Infof("UpdateDictData: id=%s, label=%s", dictData.ID, dictData.Label)

	_, err := uc.repo.FindByID(ctx, dictData.ID)
	if err != nil {
		return nil, err
	}

	return uc.repo.Update(ctx, dictData)
}

func (uc *DictDataUsecase) DeleteDictData(ctx context.Context, id string) error {
	uc.log.WithContext(ctx).Infof("DeleteDictData: id=%s", id)

	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return uc.repo.Delete(ctx, id)
}

func (uc *DictDataUsecase) ListByDictDataIDs(ctx context.Context, dictDataIds []string) ([]*DictData, error) {
	var res []*DictData
	if len(dictDataIds) == 0 {
		return res, nil
	}
	dictData, err := uc.repo.FindListByIDs(ctx, dictDataIds)
	if err != nil {
		return res, err
	}
	return dictData, nil
}

func (uc *DictDataUsecase) GetByDictTypeID(ctx context.Context, dictTypeID string) ([]*DictData, error) {
	uc.log.WithContext(ctx).Infof("GetByDictTypeID: dictTypeID=%s", dictTypeID)
	return uc.repo.FindByDictTypeID(ctx, dictTypeID)
}
