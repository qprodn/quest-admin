package tenant

import (
	"context"

	v1 "quest-admin/api/gen/tenant/v1"
	biz "quest-admin/internal/biz/tenant"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TenantPackageService struct {
	v1.UnimplementedTenantPackageServiceServer
	tpc *biz.TenantPackageUsecase
	log *log.Helper
}

func NewTenantPackageService(tpc *biz.TenantPackageUsecase, logger log.Logger) *TenantPackageService {
	return &TenantPackageService{
		tpc: tpc,
		log: log.NewHelper(log.With(logger, "module", "tenant/service/package")),
	}
}

func (s *TenantPackageService) CreateTenantPackage(ctx context.Context, in *v1.CreateTenantPackageRequest) (*emptypb.Empty, error) {
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

func (s *TenantPackageService) GetTenantPackage(ctx context.Context, in *v1.GetTenantPackageRequest) (*v1.GetTenantPackageReply, error) {
	pkg, err := s.tpc.GetTenantPackage(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &v1.GetTenantPackageReply{
		Package: s.toProtoPackage(pkg),
	}, nil
}

func (s *TenantPackageService) ListTenantPackages(ctx context.Context, in *v1.ListTenantPackagesRequest) (*v1.ListTenantPackagesReply, error) {
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

func (s *TenantPackageService) UpdateTenantPackage(ctx context.Context, in *v1.UpdateTenantPackageRequest) (*emptypb.Empty, error) {
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

func (s *TenantPackageService) DeleteTenantPackage(ctx context.Context, in *v1.DeleteTenantPackageRequest) (*emptypb.Empty, error) {
	err := s.tpc.DeleteTenantPackage(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *TenantPackageService) toProtoPackage(pkg *biz.TenantPackage) *v1.TenantPackageInfo {
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
