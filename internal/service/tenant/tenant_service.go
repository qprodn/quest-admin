package tenant

import (
	"context"

	v1 "quest-admin/api/gen/tenant/v1"
	biz "quest-admin/internal/biz/tenant"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TenantService struct {
	v1.UnimplementedTenantServiceServer
	tc  *biz.TenantUsecase
	tpc *biz.TenantPackageUsecase
	log *log.Helper
}

func NewTenantService(tc *biz.TenantUsecase, tpc *biz.TenantPackageUsecase, logger log.Logger) *TenantService {
	return &TenantService{
		tc:  tc,
		tpc: tpc,
		log: log.NewHelper(log.With(logger, "module", "tenant/service")),
	}
}

func (s *TenantService) CreateTenant(ctx context.Context, in *v1.CreateTenantRequest) (*emptypb.Empty, error) {
	tenant := &biz.Tenant{
		Name:          in.GetName(),
		ContactUserID: in.GetContactUserId(),
		ContactName:   in.GetContactName(),
		ContactMobile: in.GetContactMobile(),
		Website:       in.GetWebsite(),
		PackageID:     in.GetPackageId(),
		ExpireTime:    in.GetExpireTime().AsTime(),
		AccountCount:  in.GetAccountCount(),
		Status:        1,
	}

	_, err := s.tc.CreateTenant(ctx, tenant)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *TenantService) GetTenant(ctx context.Context, in *v1.GetTenantRequest) (*v1.GetTenantReply, error) {
	tenant, err := s.tc.GetTenant(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &v1.GetTenantReply{
		Tenant: s.toProtoTenant(tenant),
	}, nil
}

func (s *TenantService) ListTenants(ctx context.Context, in *v1.ListTenantsRequest) (*v1.ListTenantsReply, error) {
	query := &biz.ListTenantsQuery{
		Page:      in.GetPage(),
		PageSize:  in.GetPageSize(),
		Keyword:   in.GetKeyword(),
		SortField: in.GetSortField(),
		SortOrder: in.GetSortOrder(),
	}

	if in.Status != nil {
		query.Status = in.Status
	}

	result, err := s.tc.ListTenants(ctx, query)
	if err != nil {
		return nil, err
	}

	tenants := make([]*v1.TenantInfo, 0, len(result.Tenants))
	for _, tenant := range result.Tenants {
		tenants = append(tenants, s.toProtoTenant(tenant))
	}

	return &v1.ListTenantsReply{
		Tenants:    tenants,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

func (s *TenantService) UpdateTenant(ctx context.Context, in *v1.UpdateTenantRequest) (*emptypb.Empty, error) {
	tenant := &biz.Tenant{
		ID:            in.GetId(),
		Name:          in.GetName(),
		ContactUserID: in.GetContactUserId(),
		ContactName:   in.GetContactName(),
		ContactMobile: in.GetContactMobile(),
		Status:        in.GetStatus(),
		Website:       in.GetWebsite(),
		PackageID:     in.GetPackageId(),
		ExpireTime:    in.GetExpireTime().AsTime(),
		AccountCount:  in.GetAccountCount(),
	}

	_, err := s.tc.UpdateTenant(ctx, tenant)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *TenantService) DeleteTenant(ctx context.Context, in *v1.DeleteTenantRequest) (*emptypb.Empty, error) {
	err := s.tc.DeleteTenant(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *TenantService) CreateTenantPackage(ctx context.Context, in *v1.CreateTenantPackageRequest) (*emptypb.Empty, error) {
	pkg := &biz.TenantPackage{
		Name:    in.GetName(),
		Remark:  in.GetRemark(),
		MenuIDs: in.GetMenuIds(),
		Status:  1,
	}

	_, err := s.tpc.CreateTenantPackage(ctx, pkg)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *TenantService) GetTenantPackage(ctx context.Context, in *v1.GetTenantPackageRequest) (*v1.GetTenantPackageReply, error) {
	pkg, err := s.tpc.GetTenantPackage(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &v1.GetTenantPackageReply{
		Package: s.toProtoPackage(pkg),
	}, nil
}

func (s *TenantService) ListTenantPackages(ctx context.Context, in *v1.ListTenantPackagesRequest) (*v1.ListTenantPackagesReply, error) {
	query := &biz.ListPackagesQuery{
		Page:      in.GetPage(),
		PageSize:  in.GetPageSize(),
		Keyword:   in.GetKeyword(),
		SortField: in.GetSortField(),
		SortOrder: in.GetSortOrder(),
	}

	if in.Status != nil {
		query.Status = in.Status
	}

	result, err := s.tpc.ListTenantPackages(ctx, query)
	if err != nil {
		return nil, err
	}

	packages := make([]*v1.TenantPackageInfo, 0, len(result.Packages))
	for _, pkg := range result.Packages {
		packages = append(packages, s.toProtoPackage(pkg))
	}

	return &v1.ListTenantPackagesReply{
		Packages:   packages,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

func (s *TenantService) UpdateTenantPackage(ctx context.Context, in *v1.UpdateTenantPackageRequest) (*emptypb.Empty, error) {
	pkg := &biz.TenantPackage{
		ID:      in.GetId(),
		Name:    in.GetName(),
		Status:  in.GetStatus(),
		Remark:  in.GetRemark(),
		MenuIDs: in.GetMenuIds(),
	}

	_, err := s.tpc.UpdateTenantPackage(ctx, pkg)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *TenantService) DeleteTenantPackage(ctx context.Context, in *v1.DeleteTenantPackageRequest) (*emptypb.Empty, error) {
	err := s.tpc.DeleteTenantPackage(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *TenantService) toProtoTenant(tenant *biz.Tenant) *v1.TenantInfo {
	return &v1.TenantInfo{
		Id:            tenant.ID,
		Name:          tenant.Name,
		ContactUserId: tenant.ContactUserID,
		ContactName:   tenant.ContactName,
		ContactMobile: tenant.ContactMobile,
		Status:        tenant.Status,
		Website:       tenant.Website,
		PackageId:     tenant.PackageID,
		ExpireTime:    timestamppb.New(tenant.ExpireTime),
		AccountCount:  tenant.AccountCount,
		CreateAt:      timestamppb.New(tenant.CreateAt),
		UpdateAt:      timestamppb.New(tenant.UpdateAt),
	}
}

func (s *TenantService) toProtoPackage(pkg *biz.TenantPackage) *v1.TenantPackageInfo {
	return &v1.TenantPackageInfo{
		Id:       pkg.ID,
		Name:     pkg.Name,
		Status:   pkg.Status,
		Remark:   pkg.Remark,
		MenuIds:  pkg.MenuIDs,
		CreateAt: timestamppb.New(pkg.CreateAt),
		UpdateAt: timestamppb.New(pkg.UpdateAt),
	}
}
