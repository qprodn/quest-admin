package permission

import (
	"context"
	"quest-admin/internal/biz/permission"
	"quest-admin/internal/data/data"
	"quest-admin/pkg/lang/slices"
	"quest-admin/pkg/util/ctxs"
	"time"

	"quest-admin/pkg/util/idgen"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/uptrace/bun"
)

type RoleMenu struct {
	bun.BaseModel `bun:"table:qa_role_menu,alias:rm"`

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

func NewRoleMenuRepo(data *data.Data, logger log.Logger) permission.RoleMenuRepo {
	return &roleMapMenuRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *roleMapMenuRepo) FindListByRoleIDs(ctx context.Context, roles []string) ([]*permission.RoleMenu, error) {
	var roleMenus []*RoleMenu
	err := r.data.DB(ctx).NewSelect().
		Model(&roleMenus).
		Where("role_id in ?", bun.In(roles)).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	if len(roleMenus) == 0 {
		return nil, nil
	}
	bizRoleMenus := make([]*permission.RoleMenu, 0, len(roleMenus))
	bizRoleMenus = slices.Map(roleMenus, func(menu *RoleMenu, index int) *permission.RoleMenu {
		return toBizRoleMenu(menu)
	})
	return bizRoleMenus, nil
}

func (r *roleMapMenuRepo) AssignMenus(ctx context.Context, roleID string, menuIDs []string) error {
	_, err := r.data.DB(ctx).NewDelete().
		Model((*RoleMenu)(nil)).
		Where("role_id = ?", roleID).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
		Exec(ctx)
	if err != nil {
		return err
	}

	if len(menuIDs) == 0 {
		return nil
	}

	now := time.Now()
	roleMenus := make([]*RoleMenu, 0, len(menuIDs))
	for _, menuID := range menuIDs {
		roleMenus = append(roleMenus, &RoleMenu{
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
	var roleMenus []*RoleMenu
	err := r.data.DB(ctx).NewSelect().
		Model(&roleMenus).
		Where("role_id = ?", roleID).
		Where("tenant_id = ?", ctxs.GetTenantID(ctx)).
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

func toBizRoleMenu(menu *RoleMenu) *permission.RoleMenu {
	return &permission.RoleMenu{
		ID:       menu.ID,
		RoleID:   menu.RoleID,
		MenuID:   menu.MenuID,
		CreateBy: menu.CreateBy,
		CreateAt: menu.CreateAt,
		UpdateBy: menu.UpdateBy,
		UpdateAt: menu.UpdateAt,
		DeleteAt: menu.DeleteAt,
		TenantID: menu.TenantID,
	}
}
