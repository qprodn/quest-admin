package permission

import (
	"context"

	v1 "quest-admin/api/gen/permission/v1"
	biz "quest-admin/internal/biz/permission"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MenuService struct {
	v1.UnimplementedMenuServiceServer
	mc  *biz.MenuUsecase
	log *log.Helper
}

func NewMenuService(mc *biz.MenuUsecase, logger log.Logger) *MenuService {
	return &MenuService{
		mc:  mc,
		log: log.NewHelper(logger),
	}
}

func (s *MenuService) CreateMenu(ctx context.Context, in *v1.CreateMenuRequest) (*emptypb.Empty, error) {
	menu := &biz.Menu{
		Name:          in.GetName(),
		Permission:    in.GetPermission(),
		Type:          in.GetType(),
		Sort:          in.GetSort(),
		ParentID:      in.GetParentId(),
		Path:          in.GetPath(),
		Icon:          in.GetIcon(),
		Component:     in.GetComponent(),
		ComponentName: in.GetComponentName(),
		Status:        in.GetStatus(),
		Visible:       in.GetVisible(),
		KeepAlive:     in.GetKeepAlive(),
		AlwaysShow:    in.GetAlwaysShow(),
	}

	_, err := s.mc.CreateMenu(ctx, menu)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *MenuService) GetMenu(ctx context.Context, in *v1.GetMenuRequest) (*v1.GetMenuReply, error) {
	menu, err := s.mc.GetMenu(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &v1.GetMenuReply{
		Menu: s.toProtoMenu(menu),
	}, nil
}

func (s *MenuService) GetMenuTree(ctx context.Context, in *emptypb.Empty) (*v1.GetMenuTreeReply, error) {
	menus, err := s.mc.GetMenuTree(ctx)
	if err != nil {
		return nil, err
	}

	protoMenus := make([]*v1.MenuInfo, 0, len(menus))
	for _, menu := range menus {
		protoMenus = append(protoMenus, s.toProtoMenu(menu))
	}

	return &v1.GetMenuTreeReply{
		Menus: protoMenus,
	}, nil
}

func (s *MenuService) UpdateMenu(ctx context.Context, in *v1.UpdateMenuRequest) (*emptypb.Empty, error) {
	menu := &biz.Menu{
		ID:            in.GetId(),
		Name:          in.GetName(),
		Permission:    in.GetPermission(),
		Type:          in.GetType(),
		Sort:          in.GetSort(),
		ParentID:      in.GetParentId(),
		Path:          in.GetPath(),
		Icon:          in.GetIcon(),
		Component:     in.GetComponent(),
		ComponentName: in.GetComponentName(),
		Status:        in.GetStatus(),
		Visible:       in.GetVisible(),
		KeepAlive:     in.GetKeepAlive(),
		AlwaysShow:    in.GetAlwaysShow(),
	}

	_, err := s.mc.UpdateMenu(ctx, menu)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *MenuService) DeleteMenu(ctx context.Context, in *v1.DeleteMenuRequest) (*emptypb.Empty, error) {
	err := s.mc.DeleteMenu(ctx, in.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *MenuService) toProtoMenu(menu *biz.Menu) *v1.MenuInfo {
	children := make([]*v1.MenuInfo, 0, len(menu.Children))
	for _, child := range menu.Children {
		children = append(children, s.toProtoMenu(child))
	}

	return &v1.MenuInfo{
		Id:            menu.ID,
		Name:          menu.Name,
		Permission:    menu.Permission,
		Type:          menu.Type,
		Sort:          menu.Sort,
		ParentId:      menu.ParentID,
		Path:          menu.Path,
		Icon:          menu.Icon,
		Component:     menu.Component,
		ComponentName: menu.ComponentName,
		Status:        menu.Status,
		Visible:       menu.Visible,
		KeepAlive:     menu.KeepAlive,
		AlwaysShow:    menu.AlwaysShow,
		CreateAt:      timestamppb.New(menu.CreateAt),
		UpdateAt:      timestamppb.New(menu.UpdateAt),
		Children:      children,
	}
}
