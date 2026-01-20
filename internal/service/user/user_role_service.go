package user

import (
	"context"
	"quest-admin/pkg/errorx"
	"quest-admin/types/errkey"

	v1 "quest-admin/api/gen/user/v1"
	biz "quest-admin/internal/biz/user"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *UserService) AssignUserRoles(ctx context.Context, in *v1.AssignUserRolesRequest) (*emptypb.Empty, error) {
	bo := &biz.AssignUserRolesBO{
		UserID:  in.GetId(),
		RoleIDs: in.GetRoleIds(),
	}
	roles, err := s.role.ListByRoleIDs(ctx, in.RoleIds)
	if err != nil {
		return nil, err
	}
	if len(roles) != len(in.RoleIds) {
		return nil, errorx.Err(errkey.ErrRoleNotFound)
	}
	err = s.uc.AssignUserRoles(ctx, bo)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) GetUserRoles(ctx context.Context, in *v1.GetUserRolesRequest) (*v1.GetUserRolesReply, error) {
	roleIDs, err := s.uc.GetUserRoles(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &v1.GetUserRolesReply{
		RoleIds: roleIDs,
	}, nil
}
