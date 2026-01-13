package permission

import (
	"context"

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
	menus, err := uc.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	menuMap := make(map[string]*Menu)
	for _, menu := range menus {
		menuMap[menu.ID] = menu
	}

	var roots []*Menu
	for _, menu := range menuMap {
		if menu.ParentID == "" {
			roots = append(roots, menu)
		} else if parent, ok := menuMap[menu.ParentID]; ok {
			parent.Children = append(parent.Children, menu)
		}
	}

	return roots, nil
}

func (uc *MenuUsecase) UpdateMenu(ctx context.Context, menu *Menu) error {
	_, err := uc.repo.FindByID(ctx, menu.ID)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("menu not found,id:%v", menu.ID)
		return err
	}

	if menu.ParentID != "" && menu.ParentID != menu.ID {
		_, err := uc.repo.FindByID(ctx, menu.ParentID)
		if err != nil {
			uc.log.WithContext(ctx).Errorf("parent menu not found,id:%v", menu.ID)
			return ErrInvalidParentMenu
		}
	} else {
		uc.log.WithContext(ctx).Errorf("parent menu not valid,id:%v", menu.ID)
		return ErrInvalidParentMenu
	}
	err = uc.repo.Update(ctx, menu)
	if err != nil {
		return ErrInternalServer
	}
	return nil
}

func (uc *MenuUsecase) DeleteMenu(ctx context.Context, id string) error {
	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("menu not found,id:%v", id)
		return err
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
