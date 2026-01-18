package server

import (
	authv1 "quest-admin/api/gen/auth/v1"
	orgv1 "quest-admin/api/gen/organization/v1"
	permissionv1 "quest-admin/api/gen/permission/v1"
	tenantv1 "quest-admin/api/gen/tenant/v1"
	userv1 "quest-admin/api/gen/user/v1"
	"quest-admin/internal/conf"
	authManager "quest-admin/internal/data/auth"
	"quest-admin/internal/service/auth"
	"quest-admin/internal/service/organization"
	"quest-admin/internal/service/permission"
	"quest-admin/internal/service/tenant"
	"quest-admin/internal/service/user"
	pkglogger "quest-admin/pkg/logger"
	authmiddleware "quest-admin/pkg/middleware/auth"
	"quest-admin/pkg/middleware/err"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(
	c *conf.Bootstrap,
	logger log.Logger,
	authManager *authManager.Manager,
	userService *user.UserService,
	tenantService *tenant.TenantService,
	roleService *permission.RoleService,
	menuService *permission.MenuService,
	departmentService *organization.DepartmentService,
	postService *organization.PostService,
	authService *auth.AuthService,
) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			metadata.Server(),
			pkglogger.SimpleTraceIdProvider(),
			logging.Server(logger),
			authmiddleware.AdminHttpServer(authManager),
			err.Server(),
		),
		http.Filter(handlers.CORS(
			handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Language", "Origin", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"}),
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"}),
			handlers.AllowCredentials(),
		)),
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
	userv1.RegisterUserServiceHTTPServer(srv, userService)
	tenantv1.RegisterTenantServiceHTTPServer(srv, tenantService)
	orgv1.RegisterDepartmentServiceHTTPServer(srv, departmentService)
	orgv1.RegisterPostServiceHTTPServer(srv, postService)
	permissionv1.RegisterMenuServiceHTTPServer(srv, menuService)
	permissionv1.RegisterRoleServiceHTTPServer(srv, roleService)
	authv1.RegisterAuthServiceHTTPServer(srv, authService)

	return srv
}
