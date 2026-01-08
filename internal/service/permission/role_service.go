package permission

import (
	"context"

	v1 "quest-admin/api/gen/permission/v1"
	biz "quest-admin/internal/biz/permission"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RoleService struct {
	v1.UnimplementedRoleServiceServer
	rc  *biz.RoleUsecase
	log *log.Helper
}

func NewRoleService(rc *biz.RoleUsecase, logger log.Logger) *RoleService {
	return &RoleService{
		rc:  rc,
		log: log.NewHelper(logger),
	}
}

func (s *RoleService) CreateRole(ctx context.Context, in *v1.CreateRoleRequest) (*emptypb.Empty, error) {
	role := &biz.Role{
		Name:             in.GetName(),
		Code:             in.GetCode(),
		Sort:             in.GetSort(),
		DataScope:        in.GetDataScope(),
		DataScopeDeptIDs: in.GetDataScopeDeptIds(),
		Status:           in.GetStatus(),
		Type:             in.GetType(),
		Remark:           in.GetRemark(),
	}

	_, err := s.rc.CreateRole(ctx, role)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *RoleService) GetRole(ctx context.Context, in *v1.GetRoleRequest) (*v1.GetRoleReply, error) {
	role, err := s.rc.GetRole(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &v1.GetRoleReply{
		Role: s.toProtoRole(role),
	}, nil
}

func (s *RoleService) ListRoles(ctx context.Context, in *v1.ListRolesRequest) (*v1.ListRolesReply, error) {
	query := &biz.ListRolesQuery{
		Page:      in.GetPage(),
		PageSize:  in.GetPageSize(),
		Keyword:   in.GetKeyword(),
		SortField: in.GetSortField(),
		SortOrder: in.GetSortOrder(),
	}

	if in.Status != nil {
		query.Status = in.Status
	}

	result, err := s.rc.ListRoles(ctx, query)
	if err != nil {
		return nil, err
	}

	roles := make([]*v1.RoleInfo, 0, len(result.Roles))
	for _, role := range result.Roles {
		roles = append(roles, s.toProtoRole(role))
	}

	return &v1.ListRolesReply{
		Roles:      roles,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

func (s *RoleService) UpdateRole(ctx context.Context, in *v1.UpdateRoleRequest) (*emptypb.Empty, error) {
	role := &biz.Role{
		ID:               in.GetId(),
		Name:             in.GetName(),
		Code:             in.GetCode(),
		Sort:             in.GetSort(),
		DataScope:        in.GetDataScope(),
		DataScopeDeptIDs: in.GetDataScopeDeptIds(),
		Status:           in.GetStatus(),
		Type:             in.GetType(),
		Remark:           in.GetRemark(),
	}

	_, err := s.rc.UpdateRole(ctx, role)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *RoleService) DeleteRole(ctx context.Context, in *v1.DeleteRoleRequest) (*emptypb.Empty, error) {
	err := s.rc.DeleteRole(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *RoleService) AssignRoleMenus(ctx context.Context, in *v1.AssignRoleMenusRequest) (*emptypb.Empty, error) {
	err := s.rc.AssignRoleMenus(ctx, in.GetId(), in.GetMenuIds())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *RoleService) GetRoleMenus(ctx context.Context, in *v1.GetRoleMenusRequest) (*v1.GetRoleMenusReply, error) {
	menuIDs, err := s.rc.GetRoleMenus(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &v1.GetRoleMenusReply{
		MenuIds: menuIDs,
	}, nil
}

func (s *RoleService) toProtoRole(role *biz.Role) *v1.RoleInfo {
	return &v1.RoleInfo{
		Id:               role.ID,
		Name:             role.Name,
		Code:             role.Code,
		Sort:             role.Sort,
		DataScope:        role.DataScope,
		DataScopeDeptIds: role.DataScopeDeptIDs,
		Status:           role.Status,
		Type:             role.Type,
		Remark:           role.Remark,
		CreateAt:         timestamppb.New(role.CreateAt),
		UpdateAt:         timestamppb.New(role.UpdateAt),
		TenantId:         role.TenantID,
	}
}
