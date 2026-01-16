package permission

import (
	"context"
	"errors"
	"quest-admin/pkg/errorx"
	"quest-admin/pkg/lang/slices"
	"quest-admin/types/errkey"

	"github.com/go-kratos/kratos/v2/log"
)

type RoleRepo interface {
	Create(ctx context.Context, role *Role) (*Role, error)
	FindByID(ctx context.Context, id string) (*Role, error)
	FindByName(ctx context.Context, name string) (*Role, error)
	FindByCode(ctx context.Context, code string) (*Role, error)
	List(ctx context.Context, query *ListRolesQuery) (*ListRolesResult, error)
	Update(ctx context.Context, role *Role) (*Role, error)
	Delete(ctx context.Context, id string) error
	HasUsers(ctx context.Context, id string) (bool, error)
}

type RoleMenuRepo interface {
	AssignMenus(ctx context.Context, roleID string, menuIDs []string) error
	GetMenuIDs(ctx context.Context, roleID string) ([]string, error)
	FindListByRoleIDs(ctx context.Context, roles []string) ([]*RoleMenu, error)
}

type RoleUsecase struct {
	repo         RoleRepo
	roleMenuRepo RoleMenuRepo
	log          *log.Helper
}

func NewRoleUsecase(repo RoleRepo, roleMenuRepo RoleMenuRepo, logger log.Logger) *RoleUsecase {
	return &RoleUsecase{
		repo:         repo,
		roleMenuRepo: roleMenuRepo,
		log:          log.NewHelper(log.With(logger, "module", "permission/biz/role")),
	}
}

func (uc *RoleUsecase) CreateRole(ctx context.Context, role *Role) (*Role, error) {
	uc.log.WithContext(ctx).Infof("CreateRole: name=%s, code=%s", role.Name, role.Code)

	existing, err := uc.repo.FindByName(ctx, role.Name)
	if err != nil && !errors.Is(err, ErrRoleNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrRoleNameExists
	}

	existing, err = uc.repo.FindByCode(ctx, role.Code)
	if err != nil && !errors.Is(err, ErrRoleNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrRoleCodeExists
	}

	return uc.repo.Create(ctx, role)
}

func (uc *RoleUsecase) GetRole(ctx context.Context, id string) (*Role, error) {
	uc.log.WithContext(ctx).Infof("GetRole: id=%s", id)
	return uc.repo.FindByID(ctx, id)
}

func (uc *RoleUsecase) ListRoles(ctx context.Context, query *ListRolesQuery) (*ListRolesResult, error) {
	uc.log.WithContext(ctx).Infof("ListRoles: page=%d, pageSize=%d, keyword=%s", query.Page, query.PageSize, query.Keyword)
	return uc.repo.List(ctx, query)
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
		return ErrRoleHasUsers
	}

	return uc.repo.Delete(ctx, id)
}

func (uc *RoleUsecase) AssignRoleMenu(ctx context.Context, roleID string, menuIDs []string) error {
	uc.log.WithContext(ctx).Infof("AssignRoleMenu: roleID=%s, menuIDs=%v", roleID, menuIDs)

	role, err := uc.repo.FindByID(ctx, roleID)
	if err != nil {
		return nil
	}
	if role == nil {
		return errorx.Err(errkey.ErrRoleNotFound)
	}

	return uc.roleMenuRepo.AssignMenus(ctx, roleID, menuIDs)
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
