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

type TenantPackage struct {
	bun.BaseModel `bun:"table:qa_tenant_package,alias:p"`

	ID       string     `bun:"id,pk"`
	Name     string     `bun:"name,notnull"`
	Status   int32      `bun:"status,notnull,default:1"`
	Remark   string     `bun:"remark"`
	MenuIDs  string     `bun:"menu_ids,notnull"`
	CreateBy string     `bun:"create_by,notnull"`
	CreateAt time.Time  `bun:"create_at,notnull,default:current_timestamp()"`
	UpdateBy string     `bun:"update_by"`
	UpdateAt time.Time  `bun:"update_at,notnull,default:current_timestamp()"`
	DeleteAt *time.Time `bun:"delete_at,soft_delete,nullzero"`
}

type packageRepo struct {
	data *data.Data
	log  *log.Helper
}

func NewTenantPackageRepo(data *data.Data, logger log.Logger) biz.TenantPackageRepo {
	return &packageRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *packageRepo) Create(ctx context.Context, pkg *biz.TenantPackage) (*biz.TenantPackage, error) {
	now := time.Now()
	dbPkg := &TenantPackage{
		ID:       idgen.GenerateID(),
		Name:     pkg.Name,
		Status:   pkg.Status,
		Remark:   pkg.Remark,
		MenuIDs:  pkg.MenuIDs,
		CreateBy: pkg.CreateBy,
		CreateAt: now,
		UpdateBy: pkg.UpdateBy,
		UpdateAt: now,
	}

	_, err := r.data.DB(ctx).NewInsert().Model(dbPkg).Exec(ctx)
	if err != nil {
		r.log.WithContext(ctx).Error(err)
		return nil, err
	}

	return r.toBizPackage(dbPkg), nil
}

func (r *packageRepo) FindByID(ctx context.Context, id string) (*biz.TenantPackage, error) {
	dbPkg := &TenantPackage{ID: id}
	err := r.data.DB(ctx).NewSelect().Model(dbPkg).WherePK().Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		r.log.WithContext(ctx).Error(err)
		return nil, err
	}
	return r.toBizPackage(dbPkg), nil
}

func (r *packageRepo) FindByName(ctx context.Context, name string) (*biz.TenantPackage, error) {
	dbPkg := &TenantPackage{}
	err := r.data.DB(ctx).NewSelect().Model(dbPkg).Where("name = ?", name).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		r.log.WithContext(ctx).Error(err)
		return nil, err
	}
	return r.toBizPackage(dbPkg), nil
}

func (r *packageRepo) List(ctx context.Context, query *biz.ListPackagesQuery) (*biz.ListPackagesResult, error) {
	var dbPkgs []*TenantPackage

	q := r.data.DB(ctx).NewSelect().Model(&dbPkgs)

	if query.Keyword != "" {
		q = q.Where("name LIKE ?", "%"+query.Keyword+"%")
	}

	if query.Status != nil {
		q = q.Where("status = ?", *query.Status)
	}

	total, err := q.Count(ctx)
	if err != nil {
		r.log.WithContext(ctx).Error(err)
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
		r.log.WithContext(ctx).Error(err)
		return nil, err
	}

	packages := make([]*biz.TenantPackage, 0, len(dbPkgs))
	for _, dbPkg := range dbPkgs {
		packages = append(packages, r.toBizPackage(dbPkg))
	}

	totalPages := int64(total) / int64(query.PageSize)
	if int64(total)%int64(query.PageSize) != 0 {
		totalPages++
	}

	return &biz.ListPackagesResult{
		Packages:   packages,
		Total:      int64(total),
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: int32(totalPages),
	}, nil
}

func (r *packageRepo) Update(ctx context.Context, pkg *biz.TenantPackage) (*biz.TenantPackage, error) {
	dbPkg := &TenantPackage{
		ID:       pkg.ID,
		Name:     pkg.Name,
		Status:   pkg.Status,
		Remark:   pkg.Remark,
		MenuIDs:  pkg.MenuIDs,
		UpdateAt: time.Now(),
	}

	_, err := r.data.DB(ctx).NewUpdate().Model(dbPkg).WherePK().OmitZero().Exec(ctx)
	if err != nil {
		r.log.WithContext(ctx).Error(err)
		return nil, err
	}

	return r.FindByID(ctx, pkg.ID)
}

func (r *packageRepo) Delete(ctx context.Context, id string) error {
	_, err := r.data.DB(ctx).NewDelete().
		Model((*TenantPackage)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *packageRepo) IsInUse(ctx context.Context, id string) (bool, error) {
	count, err := r.data.DB(ctx).NewSelect().
		Model((*Tenant)(nil)).
		Where("package_id = ?", id).
		Count(ctx)
	if err != nil {
		r.log.WithContext(ctx).Error(err)
		return false, err
	}
	return count > 0, nil
}

func (r *packageRepo) toBizPackage(dbPkg *TenantPackage) *biz.TenantPackage {
	return &biz.TenantPackage{
		ID:       dbPkg.ID,
		Name:     dbPkg.Name,
		Status:   dbPkg.Status,
		Remark:   dbPkg.Remark,
		MenuIDs:  dbPkg.MenuIDs,
		CreateBy: dbPkg.CreateBy,
		CreateAt: dbPkg.CreateAt,
		UpdateBy: dbPkg.UpdateBy,
		UpdateAt: dbPkg.UpdateAt,
	}
}
