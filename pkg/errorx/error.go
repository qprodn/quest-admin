package errorx

import (
	"fmt"
	"sync"

	"github.com/go-kratos/kratos/v2/errors"
)

type ErrorKey string

type ErrorTemplate struct {
	Code    int
	Reason  string
	Message string
}

type ErrorX struct {
	template *ErrorTemplate
	args     []interface{}
	metadata map[string]string
	err      error
}

var (
	registry      = make(map[ErrorKey]*ErrorTemplate)
	registryMutex sync.RWMutex
)

func Register(key ErrorKey, code int, reason, message string) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	registry[key] = &ErrorTemplate{
		Code:    code,
		Reason:  reason,
		Message: message,
	}
}

func Err(key ErrorKey, args ...interface{}) *errors.Error {
	registryMutex.RLock()
	template, ok := registry[key]
	registryMutex.RUnlock()

	if !ok {
		return errors.New(500, "UNKNOWN_ERROR", fmt.Sprintf("error key %s not registered", key))
	}

	if len(args) > 0 {
		return errors.Newf(template.Code, template.Reason, template.Message, args)
	}
	return errors.New(template.Code, template.Reason, template.Message)
}
