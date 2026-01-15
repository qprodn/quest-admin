package user

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	biz "quest-admin/internal/biz/user"
	"quest-admin/internal/data/data"
	"quest-admin/pkg/util/idgen"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:qa_user,alias:u"`
	ID            string     `bun:"id,pk"`
	Username      string     `bun:"username,notnull"`
	Password      string     `bun:"password,notnull"`
	Nickname      string     `bun:"nickname"`
	Email         string     `bun:"email"`
	Mobile        string     `bun:"mobile"`
	Sex           int32      `bun:"sex,default:0"`
	Avatar        string     `bun:"avatar"`
	Status        int32      `bun:"status,default:1"`
	Remark        string     `bun:"remark"`
	LoginIP       string     `bun:"login_ip"`
	LoginDate     time.Time  `bun:"login_date"`
	CreateBy      string     `bun:"create_by"`
	CreateAt      time.Time  `bun:"create_at,notnull,default:current_timestamp()"`
	UpdateBy      string     `bun:"update_by"`
	UpdateAt      time.Time  `bun:"update_at,notnull,default:current_timestamp()"`
	TenantID      string     `bun:"tenant_id"`
	DeleteAt      *time.Time `bun:"delete_at,soft_delete,nullzero"`
}

type userRepo struct {
	data *data.Data
	log  *log.Helper
}

func NewUserRepo(data *data.Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userRepo) Create(ctx context.Context, user *biz.User) error {
	now := time.Now()
	dbUser := &User{
		ID:        idgen.GenerateID(),
		Username:  user.Username,
		Password:  user.Password,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Mobile:    user.Mobile,
		Sex:       user.Sex,
		Avatar:    user.Avatar,
		Status:    user.Status,
		Remark:    user.Remark,
		LoginIP:   user.LoginIP,
		LoginDate: user.LoginDate,
		CreateBy:  user.CreateBy,
		CreateAt:  now,
		UpdateBy:  user.UpdateBy,
		UpdateAt:  now,
		TenantID:  user.TenantID,
	}

	_, err := r.data.DB(ctx).NewInsert().Model(dbUser).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) FindByID(ctx context.Context, id string) (*biz.User, error) {
	dbUser := &User{ID: id}
	err := r.data.DB(ctx).NewSelect().
		Model(dbUser).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, biz.ErrUserNotFound
		}
		return nil, err
	}
	return r.toBizUser(dbUser), nil
}

func (r *userRepo) FindByUsername(ctx context.Context, username string) (*biz.User, error) {
	dbUser := &User{}
	err := r.data.DB(ctx).NewSelect().Model(dbUser).Where("username = ?", username).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return r.toBizUser(dbUser), nil
}

func (r *userRepo) List(ctx context.Context, opt *biz.WhereUserOpt) ([]*biz.User, error) {
	var dbUsers []*User
	q := r.data.DB(ctx).NewSelect().Model(&dbUsers)
	if opt.Username != "" {
		q.Where("username LIKE ?", "%"+opt.Username+"%")
	}
	if opt.Mobile != "" {
		q.Where("mobile LIKE ?", "%"+opt.Mobile+"%")
	}
	if opt.Nickname != "" {
		q = q.Where("nickname LIKE ?", "%"+opt.Nickname+"%")
	}
	if opt.Status != nil {
		q = q.Where("status = ?", *opt.Status)
	}
	if opt.Sex != nil {
		q = q.Where("sex = ?", *opt.Sex)
	}
	if opt.Offset != 0 {
		q.Offset(int(opt.Offset))
	}
	if opt.Limit != 0 {
		q.Limit(int(opt.Limit))
	}
	if opt.SortField != "" && opt.SortOrder != "" {
		q = q.Order(fmt.Sprintf("%s %s", opt.SortField, opt.SortOrder))
	} else {
		q = q.Order("id DESC")
	}
	err := q.Scan(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]*biz.User, 0, len(dbUsers))
	for _, dbUser := range dbUsers {
		users = append(users, r.toBizUser(dbUser))
	}

	return users, nil
}
func (r *userRepo) Count(ctx context.Context, opt *biz.WhereUserOpt) (int64, error) {
	var dbUsers []*User
	q := r.data.DB(ctx).NewSelect().Model(&dbUsers)
	if opt.Username != "" {
		q.Where("username LIKE ?", "%"+opt.Username+"%")
	}
	if opt.Mobile != "" {
		q.Where("mobile LIKE ?", "%"+opt.Mobile+"%")
	}
	if opt.Nickname != "" {
		q = q.Where("nickname LIKE ?", "%"+opt.Nickname+"%")
	}
	if opt.Status != nil {
		q = q.Where("status = ?", *opt.Status)
	}
	if opt.Sex != nil {
		q = q.Where("sex = ?", *opt.Sex)
	}
	total, err := q.Count(ctx)
	if err != nil {
		return 0, err
	}
	return int64(total), nil
}

func (r *userRepo) Update(ctx context.Context, user *biz.User) error {
	dbUser := &User{
		ID:       user.ID,
		Nickname: user.Nickname,
		Email:    user.Email,
		Mobile:   user.Mobile,
		Sex:      user.Sex,
		Avatar:   user.Avatar,
		Remark:   user.Remark,
		UpdateAt: time.Now(),
	}

	_, err := r.data.DB(ctx).NewUpdate().Model(dbUser).WherePK().OmitZero().Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) UpdatePassword(ctx context.Context, bo *biz.UpdatePasswordBO) error {
	_, err := r.data.DB(ctx).NewUpdate().
		Model((*User)(nil)).
		Set("password = ?", bo.NewPassword).
		Set("update_at = ?", time.Now()).
		Where("id = ?", bo.UserID).
		Exec(ctx)
	return err
}

func (r *userRepo) UpdateStatus(ctx context.Context, bo *biz.UpdateStatusBO) error {
	_, err := r.data.DB(ctx).NewUpdate().
		Model((*User)(nil)).
		Set("status = ?", bo.Status).
		Set("update_at = ?", time.Now()).
		Where("id = ?", bo.UserID).
		Exec(ctx)
	return err
}

func (r *userRepo) UpdateLoginInfo(ctx context.Context, bo *biz.UpdateLoginInfoBO) error {
	_, err := r.data.DB(ctx).NewUpdate().
		Model((*User)(nil)).
		Set("login_ip = ?", bo.LoginIP).
		Set("login_date = ?", bo.LoginDate).
		Set("update_at = ?", time.Now()).
		Where("id = ?", bo.UserID).
		Exec(ctx)
	return err
}

func (r *userRepo) Delete(ctx context.Context, bo *biz.DeleteUserBO) error {
	_, err := r.data.DB(ctx).NewUpdate().
		Model((*User)(nil)).
		Set("update_by = ?", bo.UpdateBy).
		Set("update_at = ?", bo.UpdateTime).
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", bo.UserID).
		Exec(ctx)
	return err
}

func (r *userRepo) toBizUser(dbUser *User) *biz.User {
	return &biz.User{
		ID:        dbUser.ID,
		Username:  dbUser.Username,
		Password:  dbUser.Password,
		Nickname:  dbUser.Nickname,
		Email:     dbUser.Email,
		Mobile:    dbUser.Mobile,
		Sex:       dbUser.Sex,
		Avatar:    dbUser.Avatar,
		Status:    dbUser.Status,
		Remark:    dbUser.Remark,
		LoginIP:   dbUser.LoginIP,
		LoginDate: dbUser.LoginDate,
		CreateBy:  dbUser.CreateBy,
		CreateAt:  dbUser.CreateAt,
		UpdateBy:  dbUser.UpdateBy,
		UpdateAt:  dbUser.UpdateAt,
		TenantID:  dbUser.TenantID,
	}
}
