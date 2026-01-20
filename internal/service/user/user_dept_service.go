package user

import (
	"context"
	"quest-admin/pkg/errorx"
	"quest-admin/types/errkey"

	v1 "quest-admin/api/gen/user/v1"
	biz "quest-admin/internal/biz/user"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *UserService) AssignUserDepts(ctx context.Context, in *v1.AssignUserDeptRequest) (*emptypb.Empty, error) {
	bo := &biz.AssignUserDeptsBO{
		UserID:  in.GetId(),
		DeptIDs: in.GetDeptIds(),
	}
	depts, err := s.dept.ListByDeptIDs(ctx, in.DeptIds)
	if err != nil {
		return nil, err
	}
	if len(depts) != len(in.DeptIds) {
		return nil, errorx.Err(errkey.ErrDepartmentNotFound)
	}
	err = s.uc.AssignUserDepts(ctx, bo)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *UserService) GetUserDepts(ctx context.Context, in *v1.GetUserDeptsRequest) (*v1.GetUserDeptsReply, error) {
	deptIDs, err := s.uc.GetUserDepts(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &v1.GetUserDeptsReply{
		DeptIds: deptIDs,
	}, nil
}
