package tenant

import (
	"context"
	v1 "quest-admin/api/gen/tenant/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrTenantNotFound      = errors.NotFound(v1.ErrorReason_TENANT_NOT_FOUND.String(), "tenant not found")
	ErrTenantNameExists    = errors.Conflict(v1.ErrorReason_TENANT_NAME_EXISTS.String(), "tenant name already exists")
	ErrTenantHasUsers      = errors.BadRequest(v1.ErrorReason_TENANT_HAS_USERS.String(), "tenant has users")
	ErrInvalidTenantStatus = errors.BadRequest(v1.ErrorReason_INVALID_TENANT_STATUS.String(), "invalid tenant status")
	ErrInvalidExpireTime   = errors.BadRequest(v1.ErrorReason_INVALID_TENANT_EXPIRE_TIME.String(), "invalid tenant expire time")
	ErrInvalidAccountCount = errors.BadRequest(v1.ErrorReason_INVALID_TENANT_ACCOUNT_COUNT.String(), "invalid tenant account count")
)

type TenantRepo interface {
	Create(ctx context.Context, tenant *Tenant) (*Tenant, error)
	FindByID(ctx context.Context, id string) (*Tenant, error)
	FindByName(ctx context.Context, name string) (*Tenant, error)
	List(ctx context.Context, query *ListTenantsQuery) (*ListTenantsResult, error)
	Update(ctx context.Context, tenant *Tenant) (*Tenant, error)
	Delete(ctx context.Context, id string) error
	HasUsers(ctx context.Context, id string) (bool, error)
}

type TenantUsecase struct {
	repo TenantRepo
	log  *log.Helper
}

func NewTenantUsecase(repo TenantRepo, logger log.Logger) *TenantUsecase {
	return &TenantUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *TenantUsecase) CreateTenant(ctx context.Context, tenant *Tenant) (*Tenant, error) {
	uc.log.WithContext(ctx).Infof("CreateTenant: name=%s, packageID=%s", tenant.Name, tenant.PackageID)

	existing, err := uc.repo.FindByName(ctx, tenant.Name)
	if err != nil && !errors.Is(err, ErrTenantNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrTenantNameExists
	}

	return uc.repo.Create(ctx, tenant)
}

func (uc *TenantUsecase) GetTenant(ctx context.Context, id string) (*Tenant, error) {
	uc.log.WithContext(ctx).Infof("GetTenant: id=%s", id)
	return uc.repo.FindByID(ctx, id)
}

func (uc *TenantUsecase) ListTenants(ctx context.Context, query *ListTenantsQuery) (*ListTenantsResult, error) {
	uc.log.WithContext(ctx).Infof("ListTenants: page=%d, pageSize=%d, keyword=%s", query.Page, query.PageSize, query.Keyword)
	return uc.repo.List(ctx, query)
}

func (uc *TenantUsecase) UpdateTenant(ctx context.Context, tenant *Tenant) (*Tenant, error) {
	uc.log.WithContext(ctx).Infof("UpdateTenant: id=%s, name=%s", tenant.ID, tenant.Name)

	_, err := uc.repo.FindByID(ctx, tenant.ID)
	if err != nil {
		return nil, err
	}

	return uc.repo.Update(ctx, tenant)
}

func (uc *TenantUsecase) DeleteTenant(ctx context.Context, id string) error {
	uc.log.WithContext(ctx).Infof("DeleteTenant: id=%s", id)

	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	hasUsers, err := uc.repo.HasUsers(ctx, id)
	if err != nil {
		return err
	}
	if hasUsers {
		return ErrTenantHasUsers
	}

	return uc.repo.Delete(ctx, id)
}
