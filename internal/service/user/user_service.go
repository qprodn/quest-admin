package user

import (
	"context"
	"quest-admin/internal/biz/organization"
	"quest-admin/internal/biz/permission"

	v1 "quest-admin/api/gen/user/v1"
	biz "quest-admin/internal/biz/user"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService struct {
	v1.UnimplementedUserServiceServer
	uc   *biz.UserUsecase
	role *permission.RoleUsecase
	dept *organization.DepartmentUsecase
	post *organization.PostUsecase
	log  *log.Helper
}

func NewUserService(uc *biz.UserUsecase, role *permission.RoleUsecase, dept *organization.DepartmentUsecase, post *organization.PostUsecase, logger log.Logger) *UserService {
	return &UserService{
		uc:   uc,
		role: role,
		dept: dept,
		post: post,
		log:  log.NewHelper(log.With(logger, "module", "user/service")),
	}
}

func (s *UserService) CreateUser(ctx context.Context, in *v1.CreateUserRequest) (*emptypb.Empty, error) {
	user := &biz.User{
		Username: in.GetUsername(),
		Password: in.GetPassword(),
		Nickname: in.GetNickname(),
		Email:    in.GetEmail(),
		Mobile:   in.GetMobile(),
		Sex:      in.GetSex(),
		Avatar:   in.GetAvatar(),
		Remark:   in.GetRemark(),
		Status:   1,
	}

	err := s.uc.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *UserService) GetUser(ctx context.Context, in *v1.GetUserRequest) (*v1.GetUserReply, error) {
	user, err := s.uc.GetUser(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &v1.GetUserReply{
		User: s.toProtoUser(user),
	}, nil
}

func (s *UserService) ListUsers(ctx context.Context, in *v1.ListUsersRequest) (*v1.ListUsersReply, error) {
	query := &biz.ListUsersQuery{
		Page:      in.GetPage(),
		PageSize:  in.GetPageSize(),
		Username:  in.GetUsername(),
		Mobile:    in.GetMobile(),
		Nickname:  in.GetNikeName(),
		SortField: in.GetSortField(),
		SortOrder: in.GetSortOrder(),
	}

	if in.Status != nil {
		query.Status = in.Status
	}

	if in.Sex != nil {
		query.Sex = in.Sex
	}

	result, err := s.uc.ListUsers(ctx, query)
	if err != nil {
		return nil, err
	}

	users := make([]*v1.UserInfo, 0, len(result.Users))
	for _, user := range result.Users {
		users = append(users, s.toProtoUser(user))
	}

	return &v1.ListUsersReply{
		Users:      users,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, in *v1.UpdateUserRequest) (*emptypb.Empty, error) {
	user := &biz.User{
		ID:       in.GetId(),
		Nickname: in.GetNickname(),
		Email:    in.GetEmail(),
		Mobile:   in.GetMobile(),
		Sex:      in.GetSex(),
		Avatar:   in.GetAvatar(),
		Remark:   in.GetRemark(),
	}

	err := s.uc.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *UserService) ChangePassword(ctx context.Context, in *v1.ChangePasswordRequest) (*emptypb.Empty, error) {
	bo := &biz.UpdatePasswordBO{
		UserID:      in.GetId(),
		OldPassword: in.GetOldPassword(),
		NewPassword: in.GetNewPassword(),
	}
	err := s.uc.ChangePassword(ctx, bo)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *UserService) ChangeUserStatus(ctx context.Context, in *v1.ChangeUserStatusRequest) (*emptypb.Empty, error) {
	bo := &biz.UpdateStatusBO{
		UserID: in.GetId(),
		Status: in.GetStatus(),
	}
	err := s.uc.ChangeUserStatus(ctx, bo)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, in *v1.DeleteUserRequest) (*emptypb.Empty, error) {
	err := s.uc.DeleteUser(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *UserService) toProtoUser(user *biz.User) *v1.UserInfo {
	return &v1.UserInfo{
		Id:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Mobile:    user.Mobile,
		Sex:       user.Sex,
		Avatar:    user.Avatar,
		Status:    user.Status,
		Remark:    user.Remark,
		LoginIp:   user.LoginIP,
		LoginDate: timestamppb.New(user.LoginDate),
		CreateAt:  timestamppb.New(user.CreateAt),
		UpdateAt:  timestamppb.New(user.UpdateAt),
		TenantId:  user.TenantID,
	}
}
