package tenant

import (
	"context"
	"database/sql"
	"quest-admin/internal/data/data"
	"time"

	biz "quest-admin/internal/biz/tenant"
	"quest-admin/pkg/util/idgen"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/uptrace/bun"
)

type Tenant struct {
	bun.BaseModel `bun:"table:qa_tenant,alias:t"`

	ID            string     `bun:"id,pk"`
	Name          string     `bun:"name,notnull"`
	ContactUserID string     `bun:"contact_user_id"`
	ContactName   string     `bun:"contact_name,notnull"`
	ContactMobile string     `bun:"contact_mobile"`
	Status        int32      `bun:"status,notnull,default:1"`
	Website       string     `bun:"website"`
	PackageID     string     `bun:"package_id,notnull"`
	ExpireTime    time.Time  `bun:"expire_time,notnull"`
	AccountCount  int32      `bun:"account_count,notnull"`
	CreateBy      string     `bun:"create_by,notnull"`
	CreateAt      time.Time  `bun:"create_at,notnull,default:current_timestamp()"`
	UpdateBy      string     `bun:"update_by"`
	UpdateAt      time.Time  `bun:"update_at,notnull,default:current_timestamp()"`
	DeleteAt      *time.Time `bun:"delete_at,soft_delete,nullzero"`
}

type tenantRepo struct {
	data *data.Data
	log  *log.Helper
}

func NewTenantRepo(data *data.Data, logger log.Logger) biz.TenantRepo {
	return &tenantRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *tenantRepo) Create(ctx context.Context, tenant *biz.Tenant) (*biz.Tenant, error) {
	now := time.Now()
	dbTenant := &Tenant{
		ID:            idgen.GenerateID(),
		Name:          tenant.Name,
		ContactUserID: tenant.ContactUserID,
		ContactName:   tenant.ContactName,
		ContactMobile: tenant.ContactMobile,
		Status:        tenant.Status,
		Website:       tenant.Website,
		PackageID:     tenant.PackageID,
		ExpireTime:    tenant.ExpireTime,
		AccountCount:  tenant.AccountCount,
		CreateBy:      tenant.CreateBy,
		CreateAt:      now,
		UpdateBy:      tenant.UpdateBy,
		UpdateAt:      now,
	}

	_, err := r.data.DB(ctx).NewInsert().Model(dbTenant).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return r.toBizTenant(dbTenant), nil
}

func (r *tenantRepo) FindByID(ctx context.Context, id string) (*biz.Tenant, error) {
	dbTenant := &Tenant{ID: id}
	err := r.data.DB(ctx).NewSelect().Model(dbTenant).WherePK().Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return r.toBizTenant(dbTenant), nil
}

func (r *tenantRepo) FindByName(ctx context.Context, name string) (*biz.Tenant, error) {
	dbTenant := &Tenant{}
	err := r.data.DB(ctx).NewSelect().Model(dbTenant).Where("name = ?", name).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return r.toBizTenant(dbTenant), nil
}

func (r *tenantRepo) List(ctx context.Context, query *biz.ListTenantsQuery) (*biz.ListTenantsResult, error) {
	var dbTenants []*Tenant

	q := r.data.DB(ctx).NewSelect().Model(&dbTenants)

	if query.Keyword != "" {
		q = q.Where("name LIKE ? OR contact_name LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	if query.Status != nil {
		q = q.Where("status = ?", *query.Status)
	}

	total, err := q.Count(ctx)
	if err != nil {
		return nil, err
	}

	sortField := "create_at"
	if query.SortField != "" {
		sortField = query.SortField
	}
	sortOrder := "DESC"
	if query.SortOrder != "" {
		sortOrder = query.SortOrder
	}

	err = q.Order(sortField + " " + sortOrder).
		Limit(int(query.PageSize)).
		Offset(int((query.Page - 1) * query.PageSize)).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	tenants := make([]*biz.Tenant, 0, len(dbTenants))
	for _, dbTenant := range dbTenants {
		tenants = append(tenants, r.toBizTenant(dbTenant))
	}

	totalPages := int64(total) / int64(query.PageSize)
	if int64(total)%int64(query.PageSize) != 0 {
		totalPages++
	}

	return &biz.ListTenantsResult{
		Tenants:    tenants,
		Total:      int64(total),
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: int32(totalPages),
	}, nil
}

func (r *tenantRepo) Update(ctx context.Context, tenant *biz.Tenant) (*biz.Tenant, error) {
	dbTenant := &Tenant{
		ID:            tenant.ID,
		Name:          tenant.Name,
		ContactUserID: tenant.ContactUserID,
		ContactName:   tenant.ContactName,
		ContactMobile: tenant.ContactMobile,
		Status:        tenant.Status,
		Website:       tenant.Website,
		PackageID:     tenant.PackageID,
		ExpireTime:    tenant.ExpireTime,
		AccountCount:  tenant.AccountCount,
		UpdateAt:      time.Now(),
	}

	_, err := r.data.DB(ctx).NewUpdate().Model(dbTenant).WherePK().OmitZero().Exec(ctx)
	if err != nil {
		return nil, err
	}

	return r.FindByID(ctx, tenant.ID)
}

func (r *tenantRepo) Delete(ctx context.Context, id string) error {
	_, err := r.data.DB(ctx).NewDelete().
		Model((*Tenant)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *tenantRepo) HasUsers(ctx context.Context, id string) (bool, error) {
	count, err := r.data.DB(ctx).NewSelect().
		Model((*Tenant)(nil)).
		TableExpr("qa_user").
		Where("tenant_id = ?", id).
		Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *tenantRepo) toBizTenant(dbTenant *Tenant) *biz.Tenant {
	return &biz.Tenant{
		ID:            dbTenant.ID,
		Name:          dbTenant.Name,
		ContactUserID: dbTenant.ContactUserID,
		ContactName:   dbTenant.ContactName,
		ContactMobile: dbTenant.ContactMobile,
		Status:        dbTenant.Status,
		Website:       dbTenant.Website,
		PackageID:     dbTenant.PackageID,
		ExpireTime:    dbTenant.ExpireTime,
		AccountCount:  dbTenant.AccountCount,
		CreateBy:      dbTenant.CreateBy,
		CreateAt:      dbTenant.CreateAt,
		UpdateBy:      dbTenant.UpdateBy,
		UpdateAt:      dbTenant.UpdateAt,
	}
}
