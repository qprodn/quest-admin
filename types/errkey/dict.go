package errkey

import "quest-admin/pkg/errorx"

var (
	ErrDictTypeNotFound   errorx.ErrorKey = "DICT_TYPE_NOT_FOUND"
	ErrDictTypeCodeExists errorx.ErrorKey = "DICT_TYPE_CODE_EXISTS"
	ErrDictTypeHasData   errorx.ErrorKey = "DICT_TYPE_HAS_DATA"

	ErrDictDataNotFound   errorx.ErrorKey = "DICT_DATA_NOT_FOUND"
	ErrDictDataValueExists errorx.ErrorKey = "DICT_DATA_VALUE_EXISTS"
)

func init() {
	errorx.Register(ErrDictTypeNotFound, 404, "DICT_TYPE_NOT_FOUND", "字典类型不存在")
	errorx.Register(ErrDictTypeCodeExists, 409, "DICT_TYPE_CODE_EXISTS", "字典类型编码已存在")
	errorx.Register(ErrDictTypeHasData, 400, "DICT_TYPE_HAS_DATA", "字典类型下存在字典数据，无法删除")

	errorx.Register(ErrDictDataNotFound, 404, "DICT_DATA_NOT_FOUND", "字典数据不存在")
	errorx.Register(ErrDictDataValueExists, 409, "DICT_DATA_VALUE_EXISTS", "字典值已存在")
}
