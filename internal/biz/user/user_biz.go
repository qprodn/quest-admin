package user

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	v1 "quest-admin/api/gen/user/v1"
	"quest-admin/pkg/util/pswd"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrUserNotFound            = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
	ErrUserExists              = errors.Conflict(v1.ErrorReason_USERNAME_ALREADY_EXISTS.String(), "user already exists")
	ErrInvalidPassword         = errors.BadRequest(v1.ErrorReason_INVALID_PASSWORD.String(), "invalid password")
	ErrPasswordConfirmMismatch = errors.BadRequest(v1.ErrorReason_PASSWORD_CONFIRM_MISMATCH.String(), "password confirm mismatch")
)

type User struct {
	ID        string
	Username  string
	Password  string
	Nickname  string
	Email     string
	Mobile    string
	Sex       int32
	Avatar    string
	Status    int32
	Remark    string
	LoginIP   string
	LoginDate time.Time
	CreateBy  string
	CreateAt  time.Time
	UpdateBy  string
	UpdateAt  time.Time
	TenantID  string
}

type ListUsersQuery struct {
	Page      int32
	PageSize  int32
	Keyword   string
	Status    *int32
	Sex       *int32
	SortField string
	SortOrder string
}

type ListUsersResult struct {
	Users      []*User
	Total      int64
	Page       int32
	PageSize   int32
	TotalPages int32
}

type UserRepo interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	List(ctx context.Context, query *ListUsersQuery) (*ListUsersResult, error)
	Update(ctx context.Context, user *User) (*User, error)
	UpdatePassword(ctx context.Context, bo *UpdatePasswordBO) error
	UpdateStatus(ctx context.Context, bo *UpdateStatusBO) error
	UpdateLoginInfo(ctx context.Context, bo *UpdateLoginInfoBO) error
	GetUserPosts(ctx context.Context, id string) ([]string, error)
	ManageUserPosts(ctx context.Context, bo *ManageUserPostsBO) error
	Delete(ctx context.Context, id string) error
}

type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *UserUsecase) CreateUser(ctx context.Context, user *User) (*User, error) {
	uc.log.WithContext(ctx).Infof("CreateUser: username=%s", user.Username)

	existing, err := uc.repo.GetByUsername(ctx, user.Username)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUserExists
	}

	return uc.repo.Create(ctx, user)
}

func (uc *UserUsecase) GetUser(ctx context.Context, id string) (*User, error) {
	uc.log.WithContext(ctx).Infof("GetUser: id=%s", id)
	return uc.repo.GetByID(ctx, id)
}

func (uc *UserUsecase) ListUsers(ctx context.Context, query *ListUsersQuery) (*ListUsersResult, error) {
	uc.log.WithContext(ctx).Infof("ListUsers: page=%d, pageSize=%d, keyword=%s", query.Page, query.PageSize, query.Keyword)
	return uc.repo.List(ctx, query)
}

func (uc *UserUsecase) UpdateUser(ctx context.Context, user *User) (*User, error) {
	uc.log.WithContext(ctx).Infof("UpdateUser: id=%s", user.ID)
	return uc.repo.Update(ctx, user)
}

func (uc *UserUsecase) ChangePassword(ctx context.Context, bo *UpdatePasswordBO) error {
	uc.log.WithContext(ctx).Infof("ChangePassword: id=%s", bo.UserID)

	_, err := uc.repo.GetByID(ctx, bo.UserID)
	if err != nil {
		return err
	}

	return uc.repo.UpdatePassword(ctx, bo)
}

func (uc *UserUsecase) ResetPassword(ctx context.Context, id string) (string, error) {
	uc.log.WithContext(ctx).Infof("ResetPassword: id=%s", id)

	_, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return "", err
	}

	tempPassword := generateTempPassword()
	hashedPassword, err := hashPassword(tempPassword)
	if err != nil {
		return "", err
	}

	err = uc.repo.UpdatePassword(ctx, &UpdatePasswordBO{
		UserID:      id,
		NewPassword: hashedPassword,
	})
	if err != nil {
		return "", err
	}

	return tempPassword, nil
}

func (uc *UserUsecase) ChangeUserStatus(ctx context.Context, bo *UpdateStatusBO) error {
	uc.log.WithContext(ctx).Infof("ChangeUserStatus: id=%s, status=%d", bo.UserID, bo.Status)

	_, err := uc.repo.GetByID(ctx, bo.UserID)
	if err != nil {
		return err
	}

	return uc.repo.UpdateStatus(ctx, bo)
}

func (uc *UserUsecase) ManageUserPosts(ctx context.Context, bo *ManageUserPostsBO) error {
	uc.log.WithContext(ctx).Infof("ManageUserPosts: id=%s, operation=%s, postCount=%d", bo.UserID, bo.Operation, len(bo.PostIDs))

	_, err := uc.repo.GetByID(ctx, bo.UserID)
	if err != nil {
		return err
	}

	return uc.repo.ManageUserPosts(ctx, bo)
}

func (uc *UserUsecase) DeleteUser(ctx context.Context, id string) error {
	uc.log.WithContext(ctx).Infof("DeleteUser: id=%s", id)

	_, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return uc.repo.Delete(ctx, id)
}

func (uc *UserUsecase) UpdateLoginInfo(ctx context.Context, bo *UpdateLoginInfoBO) error {
	uc.log.WithContext(ctx).Infof("UpdateLoginInfo: id=%s, loginIP=%s", bo.UserID, bo.LoginIP)
	return uc.repo.UpdateLoginInfo(ctx, bo)
}

func (uc *UserUsecase) GetUserPosts(ctx context.Context, id string) ([]string, error) {
	uc.log.WithContext(ctx).Infof("GetUserPosts: id=%s", id)
	return uc.repo.GetUserPosts(ctx, id)
}

func hashPassword(password string) (string, error) {
	return pswd.HashPassword(password)
}

func generateTempPassword() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	const length = 12

	result := make([]byte, length)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}
	return string(result)
}
