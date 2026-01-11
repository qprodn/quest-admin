package tenant

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

type TenantPackageRepo interface {
	Create(ctx context.Context, pkg *TenantPackage) (*TenantPackage, error)
	FindByID(ctx context.Context, id string) (*TenantPackage, error)
	FindByName(ctx context.Context, name string) (*TenantPackage, error)
	List(ctx context.Context, query *ListPackagesQuery) (*ListPackagesResult, error)
	Update(ctx context.Context, pkg *TenantPackage) (*TenantPackage, error)
	Delete(ctx context.Context, id string) error
	IsInUse(ctx context.Context, id string) (bool, error)
}

type TenantPackageUsecase struct {
	repo TenantPackageRepo
	log  *log.Helper
}

func NewTenantPackageUsecase(repo TenantPackageRepo, logger log.Logger) *TenantPackageUsecase {
	return &TenantPackageUsecase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "tenant/biz/package")),
	}
}

func (uc *TenantPackageUsecase) CreateTenantPackage(ctx context.Context, pkg *TenantPackage) (*TenantPackage, error) {
	uc.log.WithContext(ctx).Infof("CreateTenantPackage: name=%s", pkg.Name)

	existing, err := uc.repo.FindByName(ctx, pkg.Name)
	if err != nil && !errors.Is(err, ErrTenantPackageNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrTenantPackageNameExists
	}

	return uc.repo.Create(ctx, pkg)
}

func (uc *TenantPackageUsecase) GetTenantPackage(ctx context.Context, id string) (*TenantPackage, error) {
	uc.log.WithContext(ctx).Infof("GetTenantPackage: id=%s", id)
	return uc.repo.FindByID(ctx, id)
}

func (uc *TenantPackageUsecase) ListTenantPackages(ctx context.Context, query *ListPackagesQuery) (*ListPackagesResult, error) {
	uc.log.WithContext(ctx).Infof("ListTenantPackages: page=%d, pageSize=%d, keyword=%s", query.Page, query.PageSize, query.Keyword)
	return uc.repo.List(ctx, query)
}

func (uc *TenantPackageUsecase) UpdateTenantPackage(ctx context.Context, pkg *TenantPackage) (*TenantPackage, error) {
	uc.log.WithContext(ctx).Infof("UpdateTenantPackage: id=%s, name=%s", pkg.ID, pkg.Name)

	_, err := uc.repo.FindByID(ctx, pkg.ID)
	if err != nil {
		return nil, err
	}

	return uc.repo.Update(ctx, pkg)
}

func (uc *TenantPackageUsecase) DeleteTenantPackage(ctx context.Context, id string) error {
	uc.log.WithContext(ctx).Infof("DeleteTenantPackage: id=%s", id)

	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	inUse, err := uc.repo.IsInUse(ctx, id)
	if err != nil {
		return err
	}
	if inUse {
		return ErrTenantPackageInUse
	}

	return uc.repo.Delete(ctx, id)
}
