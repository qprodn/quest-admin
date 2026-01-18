package ctxs

import "context"

var (
	LoginIDKey = "login_id"
	TenantKey  = "tenant_id"
)

func GetLoginID(ctx context.Context) string {
	if val, ok := ctx.Value(LoginIDKey).(string); ok {
		return val
	}
	return ""
}

func GetTenantID(ctx context.Context) string {
	if val, ok := ctx.Value(TenantKey).(string); ok {
		return val
	}
	return ""
}
