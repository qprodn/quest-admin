package encoder

import (
	stdhttp "net/http"

	"github.com/go-kratos/kratos/v2/errors"

	"github.com/go-kratos/kratos/v2/transport/http"
)

const DefaultErrorMessage = "出错了，请稍后再试~"
const DefaultErrorReason = "INTERNAL_SERVER_ERROR"

func FromError(err error) *errors.Error {
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

func ErrorEncoder(w stdhttp.ResponseWriter, r *stdhttp.Request, err error) {
	se := FromError(err)
	codec, _ := http.CodecForRequest(r, "Accept")
	body, err := codec.Marshal(se)
	if err != nil {
		w.WriteHeader(stdhttp.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/"+codec.Name())
	w.WriteHeader(int(se.Code))
	_, _ = w.Write(body)
}
