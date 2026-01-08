package permission

import (
	"context"
	"time"

	"quest-admin/internal/data/data"
	"quest-admin/pkg/util/idgen"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/uptrace/bun"
)

type RoleMapMenu struct {
	bun.BaseModel `bun:"table:qa_role_map_menu,alias:rm"`

	ID       string     `bun:"id,pk"`
	RoleID   string     `bun:"role_id,notnull"`
	MenuID   string     `bun:"menu_id,notnull"`
	CreateBy string     `bun:"create_by"`
	CreateAt time.Time  `bun:"create_at,notnull,default:current_timestamp()"`
	UpdateBy string     `bun:"update_by"`
	UpdateAt time.Time  `bun:"update_at,notnull,default:current_timestamp()"`
	TenantID string     `bun:"tenant_id,notnull"`
	DeleteAt *time.Time `bun:"delete_at,soft_delete,nullzero"`
}

type roleMapMenuRepo struct {
	data *data.Data
	log  *log.Helper
}

func NewRoleMapMenuRepo(data *data.Data, logger log.Logger) *roleMapMenuRepo {
	return &roleMapMenuRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *roleMapMenuRepo) AssignMenus(ctx context.Context, roleID string, menuIDs []string) error {
	_, err := r.data.DB(ctx).NewDelete().
		Model((*RoleMapMenu)(nil)).
		Where("role_id = ?", roleID).
		Exec(ctx)
	if err != nil {
		return err
	}

	if len(menuIDs) == 0 {
		return nil
	}

	now := time.Now()
	roleMenus := make([]*RoleMapMenu, 0, len(menuIDs))
	for _, menuID := range menuIDs {
		roleMenus = append(roleMenus, &RoleMapMenu{
			ID:       idgen.GenerateID(),
			RoleID:   roleID,
			MenuID:   menuID,
			CreateAt: now,
			UpdateAt: now,
		})
	}

	_, err = r.data.DB(ctx).NewInsert().Model(&roleMenus).Exec(ctx)
	return err
}

func (r *roleMapMenuRepo) GetMenuIDs(ctx context.Context, roleID string) ([]string, error) {
	var roleMenus []*RoleMapMenu
	err := r.data.DB(ctx).NewSelect().
		Model(&roleMenus).
		Where("role_id = ?", roleID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	menuIDs := make([]string, 0, len(roleMenus))
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, rm.MenuID)
	}
	return menuIDs, nil
}
