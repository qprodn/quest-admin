package server

import (
	greeterv1 "quest-admin/api/gen/helloworld/v1"
	orgv1 "quest-admin/api/gen/organization/v1"
	permissionv1 "quest-admin/api/gen/permission/v1"
	tenantv1 "quest-admin/api/gen/tenant/v1"
	userv1 "quest-admin/api/gen/user/v1"
	"quest-admin/internal/conf"
	"quest-admin/internal/service/greeter"
	"quest-admin/internal/service/organization"
	"quest-admin/internal/service/permission"
	"quest-admin/internal/service/tenant"
	"quest-admin/internal/service/user"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Bootstrap, logger log.Logger,
	greeter *greeter.GreeterService,
	userService *user.UserService,
	userRoleService *user.UserRoleService,
	tenantService *tenant.TenantService,
	roleService *permission.RoleService,
	menuService *permission.MenuService,
	departmentService *organization.DepartmentService,
	postService *organization.PostService,
) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
		),
	}
	if c.Server.Http.Network != "" {
		opts = append(opts, http.Network(c.Server.Http.Network))
	}
	if c.Server.Http.Addr != "" {
		opts = append(opts, http.Address(c.Server.Http.Addr))
	}
	if c.Server.Http.Timeout != 0 {
		opts = append(opts, http.Timeout(time.Duration(c.Server.Http.Timeout)*time.Second))
	}
	srv := http.NewServer(opts...)
	greeterv1.RegisterGreeterHTTPServer(srv, greeter)
	userv1.RegisterUserServiceHTTPServer(srv, userService)
	userv1.RegisterUserRoleServiceHTTPServer(srv, userRoleService)
	tenantv1.RegisterTenantServiceHTTPServer(srv, tenantService)
	orgv1.RegisterDepartmentServiceHTTPServer(srv, departmentService)
	orgv1.RegisterPostServiceHTTPServer(srv, postService)
	permissionv1.RegisterMenuServiceHTTPServer(srv, menuService)
	permissionv1.RegisterRoleServiceHTTPServer(srv, roleService)

	return srv
}
