package permission

import (
	"context"
	"quest-admin/internal/data/idgen"
	"quest-admin/internal/data/transaction"
	"quest-admin/pkg/errorx"
	"quest-admin/pkg/lang/slices"
	"quest-admin/pkg/util/pagination"
	"quest-admin/types/consts/id"
	"quest-admin/types/errkey"

	"github.com/go-kratos/kratos/v2/log"
)

type RoleRepo interface {
	Create(ctx context.Context, role *Role) (*Role, error)
	FindByID(ctx context.Context, id string) (*Role, error)
	FindByName(ctx context.Context, name string) (*Role, error)
	FindByCode(ctx context.Context, code string) (*Role, error)
	List(ctx context.Context, opt *WhereRoleOpt) ([]*Role, error)
	Count(ctx context.Context, opt *WhereRoleOpt) (int64, error)
	Update(ctx context.Context, role *Role) (*Role, error)
	Delete(ctx context.Context, id string) error
	HasUsers(ctx context.Context, id string) (bool, error)
	FindListByIDs(ctx context.Context, roleIds []string) ([]*Role, error)
}

type RoleMenuRepo interface {
	Create(ctx context.Context, item *RoleMenu) error
	Delete(ctx context.Context, id string) error
	GetRoleMenus(ctx context.Context, roleID string) ([]*RoleMenu, error)
	GetMenuIDs(ctx context.Context, roleID string) ([]string, error)
	FindListByRoleIDs(ctx context.Context, roles []string) ([]*RoleMenu, error)
}

type RoleUsecase struct {
	tm           transaction.Manager
	idgen        *idgen.IDGenerator
	repo         RoleRepo
	roleMenuRepo RoleMenuRepo
	log          *log.Helper
}

func NewRoleUsecase(tm transaction.Manager, idgen *idgen.IDGenerator, repo RoleRepo, roleMenuRepo RoleMenuRepo, logger log.Logger) *RoleUsecase {
	return &RoleUsecase{
		tm:           tm,
		idgen:        idgen,
		repo:         repo,
		roleMenuRepo: roleMenuRepo,
		log:          log.NewHelper(log.With(logger, "module", "permission/biz/role")),
	}
}

func (uc *RoleUsecase) CreateRole(ctx context.Context, role *Role) (*Role, error) {
	uc.log.WithContext(ctx).Infof("CreateRole: name=%s, code=%s", role.Name, role.Code)

	existing, err := uc.repo.FindByName(ctx, role.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errorx.Err(errkey.ErrRoleNameExists)
	}

	existing, err = uc.repo.FindByCode(ctx, role.Code)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errorx.Err(errkey.ErrRoleCodeExists)
	}
	role.ID = uc.idgen.NextID(id.ROLE)

	return uc.repo.Create(ctx, role)
}

func (uc *RoleUsecase) GetRole(ctx context.Context, id string) (*Role, error) {
	uc.log.WithContext(ctx).Infof("GetRole: id=%s", id)
	return uc.repo.FindByID(ctx, id)
}

func (uc *RoleUsecase) ListRoles(ctx context.Context, query *ListRolesQuery) (*ListRolesResult, error) {
	uc.log.WithContext(ctx).Infof("ListRoles: page=%d, pageSize=%d, keyword=%s", query.Page, query.PageSize, query.Keyword)

	opt := &WhereRoleOpt{
		Limit:     query.PageSize,
		Offset:    pagination.GetOffset(query.Page, query.PageSize),
		Keyword:   query.Keyword,
		Status:    query.Status,
		SortField: query.SortField,
		SortOrder: query.SortOrder,
	}

	list, err := uc.repo.List(ctx, opt)
	if err != nil {
		uc.log.WithContext(ctx).Error("查询角色列表失败", err)
		return nil, err
	}

	total, err := uc.repo.Count(ctx, opt)
	if err != nil {
		uc.log.WithContext(ctx).Error("查询角色列表总数失败", err)
		return nil, err
	}

	return &ListRolesResult{
		Roles:      list,
		Total:      total,
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: pagination.GetTotalPages(total, int64(query.PageSize)),
	}, nil
}

func (uc *RoleUsecase) UpdateRole(ctx context.Context, role *Role) (*Role, error) {
	uc.log.WithContext(ctx).Infof("UpdateRole: id=%s, name=%s", role.ID, role.Name)

	_, err := uc.repo.FindByID(ctx, role.ID)
	if err != nil {
		return nil, err
	}

	return uc.repo.Update(ctx, role)
}

func (uc *RoleUsecase) DeleteRole(ctx context.Context, id string) error {
	uc.log.WithContext(ctx).Infof("DeleteRole: id=%s", id)

	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	hasUsers, err := uc.repo.HasUsers(ctx, id)
	if err != nil {
		return err
	}
	if hasUsers {
		return errorx.Err(errkey.ErrRoleHasUsers)
	}

	return uc.repo.Delete(ctx, id)
}

func (uc *RoleUsecase) AssignRoleMenu(ctx context.Context, bo *AssignRoleMenuBO) error {
	role, err := uc.repo.FindByID(ctx, bo.RoleID)
	if err != nil {
		return nil
	}
	if role == nil {
		return errorx.Err(errkey.ErrRoleNotFound)
	}

	dbRoleMenus, err := uc.roleMenuRepo.GetRoleMenus(ctx, bo.RoleID)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("获取当前角色关联菜单出现错误,error:%v", err)
		return err
	}
	dbRoleMenuIDs := slices.Map(dbRoleMenus, func(item *RoleMenu, index int) string {
		return item.MenuID
	})
	newRoleMenuIDs := bo.MenuIDs
	needDelete, needInsert := slices.Difference(dbRoleMenuIDs, newRoleMenuIDs)

	err = uc.tm.Tx(ctx, func(ctx context.Context) error {
		for _, menuID := range needInsert {
			err = uc.roleMenuRepo.Create(ctx, &RoleMenu{
				ID:     uc.idgen.NextID(id.EMPTY),
				RoleID: bo.RoleID,
				MenuID: menuID})
			if err != nil {
				uc.log.WithContext(ctx).Errorf("添加角色菜单出现错误,roleID:%s,menuID:%s,error:%v", bo.RoleID, menuID, err)
				return err
			}
		}
		for _, item := range dbRoleMenus {
			if slices.Contains(needDelete, item.MenuID) {
				err = uc.roleMenuRepo.Delete(ctx, item.ID)
				if err != nil {
					uc.log.WithContext(ctx).Errorf("删除角色菜单出现错误,roleID:%s,menuID:%s,error:%v", bo.RoleID, item.MenuID, err)
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		uc.log.WithContext(ctx).Errorf("分配角色菜单出现错误,error:%v", err)
		return err
	}
	return nil
}

func (uc *RoleUsecase) GetRoleMenus(ctx context.Context, roleID string) ([]string, error) {
	uc.log.WithContext(ctx).Infof("GetRoleMenus: roleID=%s", roleID)

	role, err := uc.repo.FindByID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	return uc.roleMenuRepo.GetMenuIDs(ctx, roleID)
}

func (uc *RoleUsecase) GetMenusByRoleIDs(ctx context.Context, roles []string) ([]string, error) {
	allMenuIDs := make([]string, 0)

	for _, roleID := range roles {
		roleMenus, err := uc.roleMenuRepo.FindListByRoleIDs(ctx, roles)
		if err != nil {
			uc.log.WithContext(ctx).Errorf("Failed to get menus for role %s: %v", roleID, err)
			continue
		}
		allMenuIDs = append(allMenuIDs, slices.Map(roleMenus, func(menu *RoleMenu, index int) string {
			return menu.MenuID
		})...)
	}
	allMenuIDs = slices.Uniq(allMenuIDs)

	return allMenuIDs, nil
}

func (uc *RoleUsecase) ListByRoleIDs(ctx context.Context, roleIds []string) ([]*Role, error) {
	var res []*Role
	if len(roleIds) == 0 {
		return res, nil
	}
	roles, err := uc.repo.FindListByIDs(ctx, roleIds)
	if err != nil {
		return roles, err
	}
	return roles, nil
}
