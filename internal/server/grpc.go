package server

import (
	userv1 "quest-admin/api/gen/user/v1"
	"quest-admin/internal/conf"
	"quest-admin/internal/service/user"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Bootstrap, logger log.Logger, userService *user.UserService, userRoleService *user.UserRoleService) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Server.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Server.Grpc.Network))
	}
	if c.Server.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Server.Grpc.Addr))
	}
	if c.Server.Grpc.Timeout != 0 {
		opts = append(opts, grpc.Timeout(time.Duration(c.Server.Grpc.Timeout)))
	}
	srv := grpc.NewServer(opts...)
	userv1.RegisterUserServiceServer(srv, userService)
	userv1.RegisterUserRoleServiceServer(srv, userRoleService)
	return srv
}
