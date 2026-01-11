package permission

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type MenuRepo interface {
	Create(ctx context.Context, menu *Menu) (*Menu, error)
	FindByID(ctx context.Context, id string) (*Menu, error)
	FindByName(ctx context.Context, name string) (*Menu, error)
	List(ctx context.Context) ([]*Menu, error)
	FindByParentID(ctx context.Context, parentID string) ([]*Menu, error)
	Update(ctx context.Context, menu *Menu) (*Menu, error)
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

func (uc *MenuUsecase) CreateMenu(ctx context.Context, menu *Menu) (*Menu, error) {
	uc.log.WithContext(ctx).Infof("CreateMenu: name=%s, parentID=%s", menu.Name, menu.ParentID)

	if menu.ParentID != "" && menu.ParentID != "0" {
		_, err := uc.repo.FindByID(ctx, menu.ParentID)
		if err != nil {
			return nil, ErrInvalidParentMenu
		}
	}

	return uc.repo.Create(ctx, menu)
}

func (uc *MenuUsecase) GetMenu(ctx context.Context, id string) (*Menu, error) {
	uc.log.WithContext(ctx).Infof("GetMenu: id=%s", id)
	return uc.repo.FindByID(ctx, id)
}

func (uc *MenuUsecase) GetMenuTree(ctx context.Context) ([]*Menu, error) {
	uc.log.WithContext(ctx).Infof("GetMenuTree")

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
		if menu.ParentID == "" || menu.ParentID == "0" {
			roots = append(roots, menu)
		} else if parent, ok := menuMap[menu.ParentID]; ok {
			parent.Children = append(parent.Children, menu)
		}
	}

	return roots, nil
}

func (uc *MenuUsecase) UpdateMenu(ctx context.Context, menu *Menu) (*Menu, error) {
	uc.log.WithContext(ctx).Infof("UpdateMenu: id=%s, name=%s", menu.ID, menu.Name)

	_, err := uc.repo.FindByID(ctx, menu.ID)
	if err != nil {
		return nil, err
	}

	if menu.ParentID != "" && menu.ParentID != "0" && menu.ParentID != menu.ID {
		_, err := uc.repo.FindByID(ctx, menu.ParentID)
		if err != nil {
			return nil, ErrInvalidParentMenu
		}
	}

	return uc.repo.Update(ctx, menu)
}

func (uc *MenuUsecase) DeleteMenu(ctx context.Context, id string) error {
	uc.log.WithContext(ctx).Infof("DeleteMenu: id=%s", id)

	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	children, err := uc.repo.FindByParentID(ctx, id)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return ErrMenuHasChildren
	}

	return uc.repo.Delete(ctx, id)
}
