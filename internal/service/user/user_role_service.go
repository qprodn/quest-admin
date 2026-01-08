package user

import (
	"context"

	v1 "quest-admin/api/gen/user/v1"
	biz "quest-admin/internal/biz/user"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserRoleService struct {
	v1.UnimplementedUserRoleServiceServer
	uc  *biz.UserRoleUsecase
	log *log.Helper
}

func NewUserRoleService(uc *biz.UserRoleUsecase, logger log.Logger) *UserRoleService {
	return &UserRoleService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

func (s *UserRoleService) ManageUserRoles(ctx context.Context, in *v1.ManageUserRolesRequest) (*emptypb.Empty, error) {
	bo := &biz.ManageUserRolesBO{
		UserID:    in.GetId(),
		RoleIDs:   in.GetRoleIds(),
		Operation: in.GetOperation(),
	}
	err := s.uc.ManageUserRoles(ctx, bo)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *UserRoleService) GetUserRoles(ctx context.Context, in *v1.GetUserRolesRequest) (*v1.GetUserRolesReply, error) {
	roleIDs, err := s.uc.GetUserRoles(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &v1.GetUserRolesReply{
		RoleIds: roleIDs,
	}, nil
}
