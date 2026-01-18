package permission

import (
	"context"
	"sort"

	"github.com/go-kratos/kratos/v2/log"
)

type MenuRepo interface {
	Create(ctx context.Context, menu *Menu) error
	FindByID(ctx context.Context, id string) (*Menu, error)
	FindByName(ctx context.Context, name string) (*Menu, error)
	List(ctx context.Context) ([]*Menu, error)
	FindByParentID(ctx context.Context, parentID string) ([]*Menu, error)
	Update(ctx context.Context, menu *Menu) error
	Delete(ctx context.Context, id string) error
	FindByMenuIDs(ctx context.Context, menuIDs []string) ([]*Menu, error)
}

type MenuUsecase struct {
	repo MenuRepo
	log  *log.Helper
}

func NewMenuUsecase(repo MenuRepo, logger log.Logger) *MenuUsecase {
	return &MenuUsecase{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "permission/biz/menu")),
	}
}

func (uc *MenuUsecase) CreateMenu(ctx context.Context, menu *Menu) error {
	if menu.ParentID != "" {
		_, err := uc.repo.FindByID(ctx, menu.ParentID)
		if err != nil {
			uc.log.WithContext(ctx).Errorf("parent menu not found,parentId:%v", menu.ParentID)
			return ErrInvalidParentMenu
		}
	}
	err := uc.repo.Create(ctx, menu)
	if err != nil {
		return ErrInternalServer
	}
	return nil
}

func (uc *MenuUsecase) GetMenu(ctx context.Context, id string) (*Menu, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *MenuUsecase) GetMenuTree(ctx context.Context) ([]*Menu, error) {
	// TODO: 根据租户套餐过滤菜单
	menus, err := uc.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	menuTree, err := uc.BuildMenuTree(menus)
	if err != nil {
		return nil, err
	}
	uc.sortMenuTree(menuTree)

	return menuTree, nil
}

func (uc *MenuUsecase) UpdateMenu(ctx context.Context, menu *Menu) error {
	dbMenu, err := uc.repo.FindByID(ctx, menu.ID)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("found menu failed,id:%v", menu.ID)
		return err
	}
	if dbMenu == nil {
		uc.log.WithContext(ctx).Errorf("menu not found,id:%v", menu.ID)
		return ErrMenuNotFound
	}

	if menu.ParentID != "" {
		if menu.ParentID == menu.ID {
			uc.log.WithContext(ctx).Errorf("parent menu not valid,id:%v", menu.ID)
			return ErrInvalidParentMenu
		}
		parentMenu, err := uc.repo.FindByID(ctx, menu.ParentID)
		if err != nil {
			uc.log.WithContext(ctx).Errorf("parent menu not found,id:%v", menu.ID)
			return ErrInvalidParentMenu
		}
		if parentMenu == nil {
			uc.log.WithContext(ctx).Errorf("parent menu not found,id:%v", menu.ParentID)
			return ErrInvalidParentMenu
		}
	}
	err = uc.repo.Update(ctx, menu)
	if err != nil {
		return ErrInternalServer
	}
	return nil
}

func (uc *MenuUsecase) DeleteMenu(ctx context.Context, id string) error {
	dbMenu, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("menu find failed,id:%v", id)
		return err
	}
	if dbMenu == nil {
		uc.log.WithContext(ctx).Errorf("menu not found,id:%v", id)
		return ErrMenuNotFound
	}

	children, err := uc.repo.FindByParentID(ctx, id)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("failed to find child menus,menuId:%v,error:%v", id, err)
		return err
	}
	if len(children) > 0 {
		uc.log.WithContext(ctx).Errorf("menu has child menus,menuId:%v", id)
		return ErrMenuHasChildren
	}
	err = uc.repo.Delete(ctx, id)
	if err != nil {
		return ErrInternalServer
	}
	return nil
}

func (uc *MenuUsecase) ListByMenuIDs(ctx context.Context, menuIDs []string) ([]*Menu, error) {
	menus, err := uc.repo.FindByMenuIDs(ctx, menuIDs)
	if err != nil {
		uc.log.WithContext(ctx).Error("failed to list menus by ids, error:%v", err)
		return nil, err
	}
	return menus, nil
}

func (uc *MenuUsecase) ProcessDisabledMenus(menus []*Menu) []*Menu {
	var result []*Menu
	for _, menu := range menus {
		if menu.Status == 1 {
			result = append(result, menu)
		}
	}
	return result
}

func (uc *MenuUsecase) BuildMenuTree(menus []*Menu) ([]*Menu, error) {
	menuMap := make(map[string]*Menu)

	for _, menu := range menus {
		menuMap[menu.ID] = menu
		menu.Children = []*Menu{}
	}

	for _, menu := range menus {
		if menu.ParentID == "" {
			continue
		}
		if parentMenu, exists := menuMap[menu.ParentID]; exists {
			parentMenu.Children = append(parentMenu.Children, menu)
		}
	}

	var rootMenus []*Menu
	for _, menu := range menus {
		if menu.ParentID == "" {
			rootMenus = append(rootMenus, menu)
		}
	}

	uc.sortMenuTree(rootMenus)

	return rootMenus, nil
}

func (uc *MenuUsecase) sortMenuTree(menus []*Menu) {
	if len(menus) == 0 {
		return
	}
	sort.Slice(menus, func(i, j int) bool {
		return menus[i].Sort < menus[j].Sort
	})
	for _, menu := range menus {
		if len(menu.Children) > 0 {
			uc.sortMenuTree(menu.Children)
		}
	}
}
