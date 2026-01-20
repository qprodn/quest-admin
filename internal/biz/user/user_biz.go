package user

import (
	"context"
	"quest-admin/internal/data/transaction"
	"quest-admin/pkg/errorx"
	"quest-admin/pkg/util/pagination"
	"quest-admin/pkg/util/pswd"
	"quest-admin/types/errkey"

	"github.com/go-kratos/kratos/v2/log"
)

type UserRepo interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	List(ctx context.Context, query *WhereUserOpt) ([]*User, error)
	Count(ctx context.Context, query *WhereUserOpt) (int64, error)
	Update(ctx context.Context, user *User) error
	UpdatePassword(ctx context.Context, bo *UpdatePasswordBO) error
	UpdateStatus(ctx context.Context, bo *UpdateStatusBO) error
	UpdateLoginInfo(ctx context.Context, bo *UpdateLoginInfoBO) error
	Delete(ctx context.Context, bo *DeleteUserBO) error
}

type UserDeptRepo interface {
	GetUserDepts(ctx context.Context, userID string) ([]string, error)
	ManageUserDepts(ctx context.Context, bo *AssignUserDeptsBO) error
}

type UserPostRepo interface {
	GetUserPosts(ctx context.Context, userID string) ([]string, error)
	ManageUserPosts(ctx context.Context, bo *AssignUserPostsBO) error
	CheckPostsExist(ctx context.Context, postIDs []string) (bool, error)
	DeleteUserPosts(ctx context.Context, userID string, postIDs []string) error
	AddUserPosts(ctx context.Context, userID string, postIDs []string) error
}

type UserRoleRepo interface {
	GetUserRoles(ctx context.Context, userID string) ([]*UserRole, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, item *UserRole) error
}

type UserUsecase struct {
	tm           transaction.Manager
	userRepo     UserRepo
	userDeptRepo UserDeptRepo
	userPostRepo UserPostRepo
	userRoleRepo UserRoleRepo
	log          *log.Helper
}

func NewUserUsecase(
	logger log.Logger,
	repo UserRepo,
	tm transaction.Manager,
	deptRepo UserDeptRepo,
	postRepo UserPostRepo,
	roleRepo UserRoleRepo,
) *UserUsecase {
	return &UserUsecase{
		log:          log.NewHelper(log.With(logger, "module", "user/biz/user")),
		tm:           tm,
		userRepo:     repo,
		userDeptRepo: deptRepo,
		userPostRepo: postRepo,
		userRoleRepo: roleRepo,
	}
}

func (uc *UserUsecase) CreateUser(ctx context.Context, user *User) error {
	existing, err := uc.userRepo.FindByUsername(ctx, user.Username)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("查询用户失败,username:%s,error:%v", user.Username, err)
		return err
	}
	if existing != nil {
		uc.log.WithContext(ctx).Error("已存在相同用户名,username:%s", user.Username)
		return ErrUserExists
	}
	if user.Password == "" {
		user.Password = "123456"
	}
	password, err := pswd.HashPassword(user.Password)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("密码加密出现错误,password:%s,error:%v", user.Password, err)
		return ErrInternalServer
	}
	user.Password = password

	return uc.userRepo.Create(ctx, user)
}

func (uc *UserUsecase) GetUser(ctx context.Context, id string) (*User, error) {
	return uc.userRepo.FindByID(ctx, id)
}

func (uc *UserUsecase) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	return uc.userRepo.FindByUsername(ctx, username)
}

func (uc *UserUsecase) VerifyPassword(ctx context.Context, hashedPassword, plainPassword string) (bool, error) {
	ok, err := pswd.VerifyPassword(plainPassword, hashedPassword)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("密码验证失败,error:%v", err)
		return false, errorx.Err(errkey.ErrInternalServer)
	}
	if !ok {
		return false, errorx.Err(errkey.ErrPasswordNotMatch)
	}
	return ok, nil
}

func (uc *UserUsecase) ListUsers(ctx context.Context, query *ListUsersQuery) (*ListUsersResult, error) {
	opt := &WhereUserOpt{
		Limit:     query.PageSize,
		Offset:    pagination.GetOffset(query.Page, query.PageSize),
		Username:  query.Username,
		Nickname:  query.Nickname,
		Mobile:    query.Mobile,
		Status:    query.Status,
		Sex:       query.Sex,
		SortField: query.SortField,
		SortOrder: query.SortOrder,
	}
	list, err := uc.userRepo.List(ctx, opt)
	if err != nil {
		uc.log.WithContext(ctx).Error("查询用户列表失败", err)
		return nil, err
	}
	total, err := uc.userRepo.Count(ctx, opt)
	if err != nil {
		uc.log.WithContext(ctx).Error("查询用户列表总数失败", err)
		return nil, err
	}
	return &ListUsersResult{
		Users:      list,
		Total:      total,
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: pagination.GetTotalPages(total, int64(query.PageSize)),
	}, nil
}

func (uc *UserUsecase) UpdateUser(ctx context.Context, user *User) error {
	return uc.userRepo.Update(ctx, user)
}

func (uc *UserUsecase) ChangePassword(ctx context.Context, bo *UpdatePasswordBO) error {
	user, err := uc.userRepo.FindByID(ctx, bo.UserID)
	if err != nil {
		uc.log.Error("查询用户失败,userID:%s,error:%v", bo.UserID, err)
		return err
	}
	ok, err := pswd.VerifyPassword(bo.OldPassword, user.Password)
	if err != nil {
		uc.log.WithContext(ctx).Error("密码验证出现错误,req:%v,error:%v", bo.OldPassword, err)
		return err
	}
	if !ok {
		uc.log.WithContext(ctx).Error("旧密码错误为匹配")
		return ErrPasswordNotMatch
	}
	password, err := pswd.HashPassword(bo.NewPassword)
	if err != nil {
		uc.log.WithContext(ctx).Error("密码加密出现错误,password:%s,error:%v", bo.NewPassword, err)
		return ErrInternalServer
	}
	bo.NewPassword = password

	return uc.userRepo.UpdatePassword(ctx, bo)
}

func (uc *UserUsecase) ChangeUserStatus(ctx context.Context, bo *UpdateStatusBO) error {
	_, err := uc.userRepo.FindByID(ctx, bo.UserID)
	if err != nil {
		uc.log.WithContext(ctx).WithContext(ctx).Error("查询用户失败,userID:%s,error:%v", bo.UserID, err)
		return err
	}

	return uc.userRepo.UpdateStatus(ctx, bo)
}

func (uc *UserUsecase) DeleteUser(ctx context.Context, id string) error {
	_, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		uc.log.WithContext(ctx).Error("用户不存在,userID:%s,error:%v", id, err)
		return err
	}

	return uc.userRepo.Delete(ctx, &DeleteUserBO{
		UserID: id,
	})
}

func (uc *UserUsecase) UpdateLoginInfo(ctx context.Context, bo *UpdateLoginInfoBO) error {
	return uc.userRepo.UpdateLoginInfo(ctx, bo)
}

func (uc *UserUsecase) VerifyStatus(ctx context.Context, user *User) (bool, error) {
	if user.Status != 1 {
		return false, nil
	}
	return true, nil
}
