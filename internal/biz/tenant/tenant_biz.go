package tenant

import (
	"context"
	"quest-admin/pkg/errorx"
	"quest-admin/types/errkey"

	"github.com/go-kratos/kratos/v2/log"
)

type TenantRepo interface {
	Create(ctx context.Context, tenant *Tenant) error
	FindByID(ctx context.Context, id string) (*Tenant, error)
	FindByName(ctx context.Context, name string) (*Tenant, error)
	List(ctx context.Context, query *ListTenantsQuery) (*ListTenantsResult, error)
	FindIDAndNameList(ctx context.Context) ([]*TenantSimple, error)
	Update(ctx context.Context, tenant *Tenant) error
	Delete(ctx context.Context, id string) error
}

type TenantUsecase struct {
	repo TenantRepo
	log  *log.Helper
}

func NewTenantUsecase(repo TenantRepo, logger log.Logger) *TenantUsecase {
	return &TenantUsecase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "tenant/biz/tenant")),
	}
}

func (uc *TenantUsecase) CreateTenant(ctx context.Context, tenant *Tenant) error {
	existing, err := uc.repo.FindByName(ctx, tenant.Name)
	if err != nil {
		return err
	}
	if existing != nil {
		return errorx.Err(errkey.ErrTenantNameExists)
	}

	return uc.repo.Create(ctx, tenant)
}

func (uc *TenantUsecase) GetTenant(ctx context.Context, id string) (*Tenant, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *TenantUsecase) ListTenants(ctx context.Context, query *ListTenantsQuery) (*ListTenantsResult, error) {
	return uc.repo.List(ctx, query)
}

func (uc *TenantUsecase) UpdateTenant(ctx context.Context, tenant *Tenant) error {
	_, err := uc.repo.FindByID(ctx, tenant.ID)
	if err != nil {
		return err
	}

	return uc.repo.Update(ctx, tenant)
}

func (uc *TenantUsecase) DeleteTenant(ctx context.Context, id string) error {
	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return uc.repo.Delete(ctx, id)
}

func (uc *TenantUsecase) GetAllTenants(ctx context.Context) ([]*TenantSimple, error) {
	return uc.repo.FindIDAndNameList(ctx)
}
