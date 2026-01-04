package user

import (
	"context"
	"time"

	v1 "quest-admin/api/gen/user/v1"
	biz "quest-admin/internal/biz/user"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService struct {
	v1.UnimplementedUserServiceServer
	uc  *biz.UserUsecase
	log *log.Helper
}

func NewUserService(uc *biz.UserUsecase, logger log.Logger) *UserService {
	return &UserService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

func (s *UserService) CreateUser(ctx context.Context, in *v1.CreateUserRequest) (*v1.CreateUserReply, error) {
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

	created, err := s.uc.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &v1.CreateUserReply{
		Id: created.ID,
	}, nil
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
		Keyword:   in.GetKeyword(),
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

func (s *UserService) UpdateUser(ctx context.Context, in *v1.UpdateUserRequest) (*v1.UpdateUserReply, error) {
	user := &biz.User{
		ID:       in.GetId(),
		Nickname: in.GetNickname(),
		Email:    in.GetEmail(),
		Mobile:   in.GetMobile(),
		Sex:      in.GetSex(),
		Avatar:   in.GetAvatar(),
		Remark:   in.GetRemark(),
	}

	updated, err := s.uc.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &v1.UpdateUserReply{
		User: s.toProtoUser(updated),
	}, nil
}

func (s *UserService) ChangePassword(ctx context.Context, in *v1.ChangePasswordRequest) (*v1.ChangePasswordReply, error) {
	err := s.uc.ChangePassword(ctx, in.GetId(), in.GetOldPassword(), in.GetNewPassword())
	if err != nil {
		return nil, err
	}

	success := true
	message := "密码修改成功"

	return &v1.ChangePasswordReply{
		Success: &success,
		Message: &message,
	}, nil
}

func (s *UserService) ResetPassword(ctx context.Context, in *v1.ResetPasswordRequest) (*v1.ResetPasswordReply, error) {
	tempPassword, err := s.uc.ResetPassword(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &v1.ResetPasswordReply{
		Success:           true,
		TemporaryPassword: tempPassword,
		Message:           "密码重置成功，请使用临时密码登录",
	}, nil
}

func (s *UserService) ChangeUserStatus(ctx context.Context, in *v1.ChangeUserStatusRequest) (*v1.ChangeUserStatusReply, error) {
	err := s.uc.ChangeUserStatus(ctx, in.GetId(), in.GetStatus())
	if err != nil {
		return nil, err
	}

	user, err := s.uc.GetUser(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &v1.ChangeUserStatusReply{
		Id:       user.ID,
		Status:   user.Status,
		UpdateAt: timestamppb.New(user.UpdateAt),
	}, nil
}

func (s *UserService) ManageUserPosts(ctx context.Context, in *v1.ManageUserPostsRequest) (*v1.ManageUserPostsReply, error) {
	err := s.uc.ManageUserPosts(ctx, in.GetId(), in.GetPostIds(), in.GetOperation())
	if err != nil {
		return nil, err
	}

	postIDs, err := s.uc.GetUserPosts(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	user, err := s.uc.GetUser(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &v1.ManageUserPostsReply{
		Id:       &user.ID,
		PostIds:  postIDs,
		UpdateAt: timestamppb.New(user.UpdateAt),
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, in *v1.DeleteUserRequest) (*v1.DeleteUserReply, error) {
	err := s.uc.DeleteUser(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	success := true
	message := "用户删除成功"
	deleteAt := timestamppb.New(time.Now())

	return &v1.DeleteUserReply{
		Success:  &success,
		Message:  &message,
		DeleteAt: deleteAt,
	}, nil
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
