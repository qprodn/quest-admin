package user

import (
	"context"
	"database/sql"
	"errors"
	"quest-admin/internal/data/data"
	"quest-admin/pkg/lang/slices"
	"quest-admin/pkg/util/ctxs"
	"time"

	biz "quest-admin/internal/biz/user"
	"quest-admin/pkg/util/idgen"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/uptrace/bun"
)

type UserRole struct {
	bun.BaseModel `bun:"table:qa_user_role,alias:ur"`

	ID       string     `bun:"id,pk"`
	UserID   string     `bun:"user_id,notnull"`
	RoleID   string     `bun:"role_id,notnull"`
	CreateBy string     `bun:"create_by"`
	CreateAt time.Time  `bun:"create_at,notnull,default:current_timestamp()"`
	UpdateBy string     `bun:"update_by"`
	UpdateAt time.Time  `bun:"update_at,notnull,default:current_timestamp()"`
	TenantID string     `bun:"tenant_id"`
	DeleteAt *time.Time `bun:"delete_at,soft_delete,nullzero"`
}

type userRoleRepo struct {
	data *data.Data
	log  *log.Helper
}

func NewUserRoleRepo(data *data.Data, logger log.Logger) biz.UserRoleRepo {
	return &userRoleRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userRoleRepo) Create(ctx context.Context, item *biz.UserRole) error {
	if item == nil {
		return nil
	}
	now := time.Now()
	_, err := r.data.DB(ctx).NewInsert().Model(&UserRole{
		ID:       item.ID,
		UserID:   item.UserID,
		RoleID:   item.RoleID,
		CreateAt: now,
		CreateBy: ctxs.GetLoginID(ctx),
		UpdateAt: now,
		UpdateBy: ctxs.GetLoginID(ctx),
		TenantID: ctxs.GetTenantID(ctx),
	}).Exec(ctx)
	if err != nil {
		r.log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

func (r *userRoleRepo) Delete(ctx context.Context, id string) error {
	_, err := r.data.DB(ctx).NewUpdate().
		Model((*UserRole)(nil)).
		Set("update_by = ?", ctxs.GetLoginID(ctx)).
		Where("id = ?", id).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Exec(ctx)
	if err != nil {
		r.log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

func (r *userRoleRepo) GetUserRoles(ctx context.Context, userID string) ([]*biz.UserRole, error) {
	var userRoles []*UserRole
	err := r.data.DB(ctx).NewSelect().
		Model(&userRoles).
		Where("user_id = ?", userID).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return make([]*biz.UserRole, 0), nil
		}
		r.log.WithContext(ctx).Error(err)
		return nil, err
	}

	return slices.Map(userRoles, func(item *UserRole, index int) *biz.UserRole {
		return r.toBizUserRole(item)
	}), nil
}

func (r *userRoleRepo) addUserRoles(ctx context.Context, userID string, roleIDs []string) error {
	now := time.Now()
	userRoles := make([]*UserRole, 0, len(roleIDs))
	for _, roleID := range roleIDs {
		userRoles = append(userRoles, &UserRole{
			ID:       idgen.GenerateID(),
			UserID:   userID,
			RoleID:   roleID,
			CreateAt: now,
		})
	}

	_, err := r.data.DB(ctx).NewInsert().Model(&userRoles).Exec(ctx)
	if err != nil {
		r.log.WithContext(ctx).Error(err)
	}
	return err
}

func (r *userRoleRepo) removeUserRoles(ctx context.Context, userID string, roleIDs []string) error {
	_, err := r.data.DB(ctx).NewDelete().
		Model((*UserRole)(nil)).
		Where("user_id = ?", userID).
		Where("role_id IN (?)", bun.In(roleIDs)).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Exec(ctx)
	if err != nil {
		r.log.WithContext(ctx).Error(err)
	}
	return err
}

func (r *userRoleRepo) replaceUserRoles(ctx context.Context, userID string, roleIDs []string) error {
	return r.data.DB(ctx).RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewDelete().
			Model((*UserRole)(nil)).
			Where("user_id = ?", userID).
			Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
			Exec(ctx)
		if err != nil {
			r.log.WithContext(ctx).Error(err)
			return err
		}

		if len(roleIDs) > 0 {
			now := time.Now()
			userRoles := make([]*UserRole, 0, len(roleIDs))
			for _, roleID := range roleIDs {
				userRoles = append(userRoles, &UserRole{
					ID:       idgen.GenerateID(),
					UserID:   userID,
					RoleID:   roleID,
					CreateAt: now,
				})
			}
			_, err = tx.NewInsert().Model(&userRoles).Exec(ctx)
			if err != nil {
				r.log.WithContext(ctx).Error(err)
				return err
			}
		}

		return nil
	})
}

func (r *userRoleRepo) toBizUserRole(item *UserRole) *biz.UserRole {
	return &biz.UserRole{
		ID:       item.ID,
		UserID:   item.UserID,
		RoleID:   item.RoleID,
		CreateBy: item.CreateBy,
		CreateAt: item.CreateAt,
		UpdateBy: item.UpdateBy,
		UpdateAt: item.UpdateAt,
		TenantID: item.TenantID,
	}
}
