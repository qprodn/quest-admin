package auth

import (
	"context"
	"quest-admin/internal/data/auth"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

func AdminHttpServer(manager *auth.Manager) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				loginID  = "unknown"
				tenantId = ""
			)
			if tr, ok := transport.FromServerContext(ctx); ok {
				token := tr.RequestHeader().Get("Authorization")
				if token != "" {
					loginID, err = manager.Admin.GetLoginID(token)
					if err != nil {
						return nil, errors.New(401, "UNAUTHORIZED", "Token is invalid")
					}
				}
				if tmpTenantId := tr.RequestHeader().Get("Tenant"); tmpTenantId != "" {
					tenantId = tmpTenantId
				}
			}
			ctx = context.WithValue(ctx, "login_id", loginID)
			ctx = context.WithValue(ctx, "tenant_id", tenantId)
			return handler(ctx, req)
		}
	}
}
