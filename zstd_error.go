package zstd

import "errors"

type Error struct {
	code ZSTD_ErrorCode
}

func (e Error) Error() string        { return "zstd: " + g()._ZSTD_getErrorString(e.code) }
func (e Error) Code() ZSTD_ErrorCode { return e.code }

var (
	ErrEmptySrc = errors.New("zstd: empty src byte slice")
)
