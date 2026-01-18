package err

import (
	"context"
	stdhttp "net/http"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
)

const DefaultErrorMessage = "出错了，请稍后再试~"
const DefaultErrorReason = "INTERNAL_SERVER_ERROR"

func Server() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			reply, err = handler(ctx, req)
			if err == nil {
				return
			}
			se := fromError(err)
			value := ctx.Value("logger_trace_id")
			if value != nil {
				if se.Metadata == nil {
					se.Metadata = make(map[string]string)
				}
				se.Metadata["trace_id"] = value.(string)
			}
			return
		}
	}
}

func fromError(err error) *errors.Error {
	if err == nil {
		return nil
	}
	if se := new(errors.Error); errors.As(err, &se) {
		return se
	}
	return errors.New(
		stdhttp.StatusInternalServerError,
		DefaultErrorReason,
		DefaultErrorMessage,
	)
}
