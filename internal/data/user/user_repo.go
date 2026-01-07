package user

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	biz "quest-admin/internal/biz/user"
	"quest-admin/internal/data/data"
	"quest-admin/pkg/util/idgen"
	"quest-admin/pkg/util/pswd"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:qa_user,alias:u"`

	ID        string     `bun:"id,pk"`
	Username  string     `bun:"username,notnull"`
	Password  string     `bun:"password,notnull"`
	Nickname  string     `bun:"nickname"`
	Email     string     `bun:"email"`
	Mobile    string     `bun:"mobile"`
	Sex       int32      `bun:"sex,default:0"`
	Avatar    string     `bun:"avatar"`
	Status    int32      `bun:"status,default:1"`
	Remark    string     `bun:"remark"`
	LoginIP   string     `bun:"login_ip"`
	LoginDate time.Time  `bun:"login_date"`
	CreateBy  string     `bun:"create_by"`
	CreateAt  time.Time  `bun:"create_at,notnull,default:current_timestamp()"`
	UpdateBy  string     `bun:"update_by"`
	UpdateAt  time.Time  `bun:"update_at,notnull,default:current_timestamp()"`
	TenantID  string     `bun:"tenant_id"`
	DeleteAt  *time.Time `bun:"delete_at,soft_delete,nullzero"`
}

type UserPost struct {
	bun.BaseModel `bun:"table:qa_user_map_post,alias:up"`

	ID       string     `bun:"id,pk"`
	UserID   string     `bun:"user_id,notnull"`
	PostID   string     `bun:"post_id,notnull"`
	CreateBy string     `bun:"create_by"`
	CreateAt time.Time  `bun:"create_at,notnull,default:current_timestamp()"`
	UpdateBy string     `bun:"update_by"`
	UpdateAt time.Time  `bun:"update_at,notnull,default:current_timestamp()"`
	TenantID string     `bun:"tenant_id"`
	DeleteAt *time.Time `bun:"delete_at,soft_delete,nullzero"`
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

func (r *userRepo) Create(ctx context.Context, user *biz.User) (*biz.User, error) {
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

	if dbUser.Password != "" {
		hashedPassword, err := pswd.HashPassword(dbUser.Password)
		if err != nil {
			return nil, err
		}
		dbUser.Password = hashedPassword
	}

	_, err := r.data.DB(ctx).NewInsert().Model(dbUser).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return r.toBizUser(dbUser), nil
}

func (r *userRepo) GetByID(ctx context.Context, id string) (*biz.User, error) {
	dbUser := &User{ID: id}
	err := r.data.DB(ctx).NewSelect().Model(dbUser).WherePK().Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, biz.ErrUserNotFound
		}
		return nil, err
	}
	return r.toBizUser(dbUser), nil
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*biz.User, error) {
	dbUser := &User{}
	err := r.data.DB(ctx).NewSelect().Model(dbUser).Where("username = ?", username).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, biz.ErrUserNotFound
		}
		return nil, err
	}
	return r.toBizUser(dbUser), nil
}

func (r *userRepo) List(ctx context.Context, query *biz.ListUsersQuery) (*biz.ListUsersResult, error) {
	var dbUsers []*User
	q := r.data.DB(ctx).NewSelect().Model(&dbUsers)

	if query.Keyword != "" {
		q = q.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.WhereOr("username LIKE ?", "%"+query.Keyword+"%").
				WhereOr("nickname LIKE ?", "%"+query.Keyword+"%").
				WhereOr("email LIKE ?", "%"+query.Keyword+"%").
				WhereOr("mobile LIKE ?", "%"+query.Keyword+"%")
		})
	}

	if query.Status != nil {
		q = q.Where("status = ?", *query.Status)
	}

	if query.Sex != nil {
		q = q.Where("sex = ?", *query.Sex)
	}

	total, err := q.ScanAndCount(ctx, &dbUsers, nil)
	if err != nil {
		return nil, err
	}

	page := query.Page
	pageSize := query.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	q = q.Offset(int(offset)).Limit(int(pageSize))

	if query.SortField != "" {
		order := query.SortOrder
		if order != "asc" && order != "desc" {
			order = "desc"
		}
		q = q.Order(fmt.Sprintf("%s %s", query.SortField, order))
	} else {
		q = q.Order("create_at DESC")
	}

	err = q.Scan(ctx)
	if err != nil {
		return nil, err
	}

	totalPages := int32((int64(total) + int64(pageSize) - 1) / int64(pageSize))

	users := make([]*biz.User, 0, len(dbUsers))
	for _, dbUser := range dbUsers {
		users = append(users, r.toBizUser(dbUser))
	}

	return &biz.ListUsersResult{
		Users:      users,
		Total:      int64(total),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *userRepo) Update(ctx context.Context, user *biz.User) (*biz.User, error) {
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
		return nil, err
	}

	return r.GetByID(ctx, user.ID)
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

func (r *userRepo) GetUserPosts(ctx context.Context, id string) ([]string, error) {
	var userPosts []*UserPost
	err := r.data.DB(ctx).NewSelect().
		Model(&userPosts).
		Where("user_id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	postIDs := make([]string, 0, len(userPosts))
	for _, up := range userPosts {
		postIDs = append(postIDs, up.PostID)
	}
	return postIDs, nil
}

func (r *userRepo) ManageUserPosts(ctx context.Context, bo *biz.ManageUserPostsBO) error {
	switch bo.Operation {
	case "add":
		return r.addUserPosts(ctx, bo.UserID, bo.PostIDs)
	case "remove":
		return r.removeUserPosts(ctx, bo.UserID, bo.PostIDs)
	case "replace":
		return r.replaceUserPosts(ctx, bo.UserID, bo.PostIDs)
	default:
		return fmt.Errorf("invalid operation: %s", bo.Operation)
	}
}

func (r *userRepo) addUserPosts(ctx context.Context, userID string, postIDs []string) error {
	now := time.Now()
	userPosts := make([]*UserPost, 0, len(postIDs))
	for _, postID := range postIDs {
		userPosts = append(userPosts, &UserPost{
			ID:       idgen.GenerateID(),
			UserID:   userID,
			PostID:   postID,
			CreateAt: now,
		})
	}

	_, err := r.data.DB(ctx).NewInsert().Model(&userPosts).Exec(ctx)
	return err
}

func (r *userRepo) removeUserPosts(ctx context.Context, userID string, postIDs []string) error {
	_, err := r.data.DB(ctx).NewDelete().
		Model((*UserPost)(nil)).
		Where("user_id = ?", userID).
		Where("post_id IN (?)", bun.In(postIDs)).
		Exec(ctx)
	return err
}

func (r *userRepo) replaceUserPosts(ctx context.Context, userID string, postIDs []string) error {
	return r.data.DB(ctx).RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewDelete().
			Model((*UserPost)(nil)).
			Where("user_id = ?", userID).
			Exec(ctx)
		if err != nil {
			return err
		}

		if len(postIDs) > 0 {
			now := time.Now()
			userPosts := make([]*UserPost, 0, len(postIDs))
			for _, postID := range postIDs {
				userPosts = append(userPosts, &UserPost{
					ID:       idgen.GenerateID(),
					UserID:   userID,
					PostID:   postID,
					CreateAt: now,
				})
			}
			_, err = tx.NewInsert().Model(&userPosts).Exec(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *userRepo) Delete(ctx context.Context, id string) error {
	_, err := r.data.DB(ctx).NewDelete().
		Model((*User)(nil)).
		Where("id = ?", id).
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
