package server

import (
	v1 "quest-admin/api/gen/helloworld/v1"
	userv1 "quest-admin/api/gen/user/v1"
	"quest-admin/internal/conf"
	"quest-admin/internal/service/greeter"
	"quest-admin/internal/service/user"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Bootstrap, logger log.Logger, greeter *greeter.GreeterService, userService *user.UserService) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Server.Http.Network != "" {
		opts = append(opts, http.Network(c.Server.Http.Network))
	}
	if c.Server.Http.Addr != "" {
		opts = append(opts, http.Address(c.Server.Http.Addr))
	}
	if c.Server.Http.Timeout != 0 {
		opts = append(opts, http.Timeout(time.Duration(c.Server.Http.Timeout)))
	}
	srv := http.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	userv1.RegisterUserServiceHTTPServer(srv, userService)
	return srv
}
