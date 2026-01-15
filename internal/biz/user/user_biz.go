package user

import (
	"context"
	"quest-admin/pkg/util/pagination"
	"quest-admin/pkg/util/pswd"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
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

type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "user/biz/user")),
	}
}

func (uc *UserUsecase) CreateUser(ctx context.Context, user *User) error {
	existing, err := uc.repo.FindByUsername(ctx, user.Username)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return err
	}
	if existing != nil {
		uc.log.WithContext(ctx).Error("已存在相同用户名")
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

	return uc.repo.Create(ctx, user)
}

func (uc *UserUsecase) GetUser(ctx context.Context, id string) (*User, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *UserUsecase) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	return uc.repo.FindByUsername(ctx, username)
}

func (uc *UserUsecase) VerifyPassword(ctx context.Context, hashedPassword, plainPassword string) (bool, error) {
	ok, err := pswd.VerifyPassword(plainPassword, hashedPassword)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("密码验证失败,error:%v", err)
		return false, ErrPasswordConfirmMismatch
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
	list, err := uc.repo.List(ctx, opt)
	if err != nil {
		uc.log.WithContext(ctx).Error("查询用户列表失败", err)
		return nil, err
	}
	total, err := uc.repo.Count(ctx, opt)
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
	return uc.repo.Update(ctx, user)
}

func (uc *UserUsecase) ChangePassword(ctx context.Context, bo *UpdatePasswordBO) error {
	user, err := uc.repo.FindByID(ctx, bo.UserID)
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

	return uc.repo.UpdatePassword(ctx, bo)
}

func (uc *UserUsecase) ChangeUserStatus(ctx context.Context, bo *UpdateStatusBO) error {
	_, err := uc.repo.FindByID(ctx, bo.UserID)
	if err != nil {
		uc.log.WithContext(ctx).WithContext(ctx).Error("查询用户失败,userID:%s,error:%v", bo.UserID, err)
		return err
	}

	return uc.repo.UpdateStatus(ctx, bo)
}

func (uc *UserUsecase) DeleteUser(ctx context.Context, id string) error {
	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		uc.log.WithContext(ctx).Error("用户不存在,userID:%s,error:%v", id, err)
		return err
	}

	return uc.repo.Delete(ctx, &DeleteUserBO{
		UserID:     id,
		UpdateBy:   "",
		UpdateTime: time.Now(),
	})
}

func (uc *UserUsecase) UpdateLoginInfo(ctx context.Context, bo *UpdateLoginInfoBO) error {
	return uc.repo.UpdateLoginInfo(ctx, bo)
}
