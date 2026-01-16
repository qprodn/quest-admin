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

type UserDept struct {
	bun.BaseModel `bun:"table:qa_user_dept,alias:ud"`

	ID       string     `bun:"id,pk"`
	UserID   string     `bun:"user_id,notnull"`
	DeptID   string     `bun:"dept_id,notnull"`
	CreateBy string     `bun:"create_by"`
	CreateAt time.Time  `bun:"create_at,notnull,default:current_timestamp()"`
	UpdateBy string     `bun:"update_by"`
	UpdateAt time.Time  `bun:"update_at,notnull,default:current_timestamp()"`
	TenantID string     `bun:"tenant_id"`
	DeleteAt *time.Time `bun:"delete_at,soft_delete,nullzero"`
}

type userDeptRepo struct {
	data *data.Data
	log  *log.Helper
}

func NewUserDeptRepo(data *data.Data, logger log.Logger) biz.UserDeptRepo {
	return &userDeptRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userDeptRepo) GetUserDepts(ctx context.Context, userID string) ([]string, error) {
	var userDepts []*UserDept
	err := r.data.DB(ctx).NewSelect().
		Model(&userDepts).
		Where("user_id = ?", userID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	deptIDs := make([]string, 0, len(userDepts))
	for _, ud := range userDepts {
		deptIDs = append(deptIDs, ud.DeptID)
	}
	return deptIDs, nil
}

func (r *userDeptRepo) ManageUserDepts(ctx context.Context, bo *biz.AssignUserDeptsBO) error {
	switch bo.Operation {
	case "add":
		return r.addUserDepts(ctx, bo.UserID, bo.DeptIDs)
	case "remove":
		return r.removeUserDepts(ctx, bo.UserID, bo.DeptIDs)
	case "replace":
		return r.replaceUserDepts(ctx, bo.UserID, bo.DeptIDs)
	default:
		return biz.ErrInvalidOperationType
	}
}

func (r *userDeptRepo) addUserDepts(ctx context.Context, userID string, deptIDs []string) error {
	now := time.Now()
	userDepts := make([]*UserDept, 0, len(deptIDs))
	for _, deptID := range deptIDs {
		userDepts = append(userDepts, &UserDept{
			ID:       idgen.GenerateID(),
			UserID:   userID,
			DeptID:   deptID,
			CreateAt: now,
		})
	}

	_, err := r.data.DB(ctx).NewInsert().Model(&userDepts).Exec(ctx)
	return err
}

func (r *userDeptRepo) removeUserDepts(ctx context.Context, userID string, deptIDs []string) error {
	_, err := r.data.DB(ctx).NewDelete().
		Model((*UserDept)(nil)).
		Where("user_id = ?", userID).
		Where("dept_id IN (?)", bun.In(deptIDs)).
		Exec(ctx)
	return err
}

func (r *userDeptRepo) replaceUserDepts(ctx context.Context, userID string, deptIDs []string) error {
	return r.data.DB(ctx).RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewDelete().
			Model((*UserDept)(nil)).
			Where("user_id = ?", userID).
			Exec(ctx)
		if err != nil {
			return err
		}

		if len(deptIDs) > 0 {
			now := time.Now()
			userDepts := make([]*UserDept, 0, len(deptIDs))
			for _, deptID := range deptIDs {
				userDepts = append(userDepts, &UserDept{
					ID:       idgen.GenerateID(),
					UserID:   userID,
					DeptID:   deptID,
					CreateAt: now,
				})
			}
			_, err = tx.NewInsert().Model(&userDepts).Exec(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
