package user

import (
	"context"
	"quest-admin/internal/data/data"
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

func (r *userRoleRepo) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	var userRoles []*UserRole
	err := r.data.DB(ctx).NewSelect().
		Model(&userRoles).
		Where("user_id = ?", userID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	roleIDs := make([]string, 0, len(userRoles))
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleID)
	}
	return roleIDs, nil
}

func (r *userRoleRepo) ManageUserRoles(ctx context.Context, bo *biz.AssignUserRolesBO) error {
	switch bo.Operation {
	case "add":
		return r.addUserRoles(ctx, bo.UserID, bo.RoleIDs)
	case "remove":
		return r.removeUserRoles(ctx, bo.UserID, bo.RoleIDs)
	case "replace":
		return r.replaceUserRoles(ctx, bo.UserID, bo.RoleIDs)
	default:
		return biz.ErrInvalidOperationType
	}
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
	return err
}

func (r *userRoleRepo) removeUserRoles(ctx context.Context, userID string, roleIDs []string) error {
	_, err := r.data.DB(ctx).NewDelete().
		Model((*UserRole)(nil)).
		Where("user_id = ?", userID).
		Where("role_id IN (?)", bun.In(roleIDs)).
		Exec(ctx)
	return err
}

func (r *userRoleRepo) replaceUserRoles(ctx context.Context, userID string, roleIDs []string) error {
	return r.data.DB(ctx).RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewDelete().
			Model((*UserRole)(nil)).
			Where("user_id = ?", userID).
			Exec(ctx)
		if err != nil {
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
				return err
			}
		}

		return nil
	})
}
