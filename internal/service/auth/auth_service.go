package auth

import (
	"context"
	v1 "quest-admin/api/gen/auth/v1"
	authBiz "quest-admin/internal/biz/auth"
	permBiz "quest-admin/internal/biz/permission"
	userBiz "quest-admin/internal/biz/user"
	"quest-admin/pkg/errorx"
	"quest-admin/pkg/lang/ptr"
	"quest-admin/pkg/lang/slices"
	"quest-admin/types/errkey"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AuthService 认证服务
type AuthService struct {
	v1.UnimplementedAuthServiceServer
	authUsecase *authBiz.AuthUsecase
	userUsecase *userBiz.UserUsecase
	roleUsecase *permBiz.RoleUsecase
	menuUsecase *permBiz.MenuUsecase
	log         *log.Helper
}

// NewAuthService 创建认证服务
func NewAuthService(
	logger log.Logger,
	authUsecase *authBiz.AuthUsecase,
	userUsecase *userBiz.UserUsecase,
	roleUsecase *permBiz.RoleUsecase,
	menuUsecase *permBiz.MenuUsecase,
) *AuthService {
	return &AuthService{
		log:         log.NewHelper(log.With(logger, "module", "auth/service")),
		authUsecase: authUsecase,
		roleUsecase: roleUsecase,
		userUsecase: userUsecase,
		menuUsecase: menuUsecase,
	}
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, request *v1.LoginRequest) (*v1.LoginReply, error) {
	token, err := s.LoginByUsernameAndPassword(ctx, request)
	if err != nil {
		return nil, err
	}
	return &v1.LoginReply{
		Token: token,
	}, nil
}

// GetPermissionInfo 获取用户权限信息
func (s *AuthService) GetPermissionInfo(ctx context.Context, in *v1.GetPermissionInfoRequest) (*v1.GetPermissionInfoReply, error) {
	s.log.WithContext(ctx).Info("GetPermissionInfo")

	userID := "1"

	user, err := s.userUsecase.GetUser(ctx, userID)
	if err != nil {
		s.log.WithContext(ctx).Errorf("获取用户信息失败,userID:%s,error:%v", userID, err)
		return nil, errors.NotFound("USER_NOT_FOUND", "用户不存在")
	}

	roleIDs, err := s.userUsecase.GetUserRoles(ctx, userID)
	if err != nil {
		s.log.WithContext(ctx).Errorf("获取用户角色失败,userID:%s,error:%v", userID, err)
		return nil, errors.InternalServer("INTERNAL_ERROR", "获取用户角色失败")
	}

	roles := make([]string, 0, len(roleIDs))
	for _, roleID := range roleIDs {
		role, err := s.roleUsecase.GetRole(ctx, roleID)
		if err != nil {
			s.log.WithContext(ctx).Warnf("获取角色信息失败,roleID:%s,error:%v", roleID, err)
			continue
		}
		roles = append(roles, role.Code)
	}

	menuIDs := make([]string, 0)
	for _, roleID := range roleIDs {
		roleMenuIDs, err := s.roleUsecase.GetRoleMenus(ctx, roleID)
		if err != nil {
			s.log.WithContext(ctx).Warnf("获取角色菜单失败,roleID:%s,error:%v", roleID, err)
			continue
		}
		menuIDs = append(menuIDs, roleMenuIDs...)
	}

	permissions := make([]string, 0)
	menus := make([]*v1.MenuInfo, 0)
	for _, menuID := range menuIDs {
		menu, err := s.menuUsecase.GetMenu(ctx, menuID)
		if err != nil {
			s.log.WithContext(ctx).Warnf("获取菜单信息失败,menuID:%s,error:%v", menuID, err)
			continue
		}
		if menu.Permission != "" {
			permissions = append(permissions, menu.Permission)
		}
		menus = append(menus, s.toProtoMenu(menu))
	}

	return &v1.GetPermissionInfoReply{
		User:        s.toProtoUser(user),
		Roles:       roles,
		Permissions: permissions,
		Menus:       menus,
	}, nil
}

func (s *AuthService) LoginByUsernameAndPassword(ctx context.Context, request *v1.LoginRequest) (token string, err error) {
	user, err := s.userUsecase.GetUserByUsername(ctx, ptr.From(request.Username))
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errorx.Err(errkey.ErrUserNotFound)
	}
	ok, err := s.userUsecase.VerifyStatus(ctx, user)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errorx.Err(errkey.ErrUserDisabled)
	}

	_, err = s.userUsecase.VerifyPassword(ctx, user.Password, ptr.From(request.Password))
	if err != nil {
		return "", err
	}

	//获取用户关联角色
	roles, err := s.userUsecase.GetUserRoles(ctx, user.ID)
	if err != nil {
		return "", err
	}

	//获取角色管理的菜单
	menuIDs, err := s.roleUsecase.GetMenusByRoleIDs(ctx, roles)
	if err != nil {
		return "", err
	}
	menus, err := s.menuUsecase.ListByMenuIDs(ctx, menuIDs)
	if err != nil {
		return "", err
	}
	menus = s.menuUsecase.ProcessDisabledMenus(menus)
	permissions := slices.Map(menus, func(item *permBiz.Menu, index int) string { return item.Permission })

	token, err = s.authUsecase.AdminGenerateToken(ctx, &authBiz.GenerateTokenBO{UserID: user.ID, Device: ptr.From(request.Device)})
	if err != nil {
		return "", err
	}
	err = s.authUsecase.SetRolesAndPermission(ctx, user.ID, roles, permissions)
	if err != nil {
		return "", err
	}
	s.log.WithContext(ctx).Infof("登录成功,userID:%s", user.ID)
	//todo 记录登录日志
	return token, nil
}

func (s *AuthService) toProtoUser(user *userBiz.User) *v1.UserInfo {
	return &v1.UserInfo{
		Id:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Email:    user.Email,
		Mobile:   user.Mobile,
		Sex:      user.Sex,
		Avatar:   user.Avatar,
		Status:   user.Status,
		Remark:   user.Remark,
		CreateAt: timestamppb.New(user.CreateAt),
	}
}

func (s *AuthService) toProtoMenu(menu *permBiz.Menu) *v1.MenuInfo {
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
	}
}
