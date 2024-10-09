package zstd

import (
	_ "embed"
	"runtime"
	"sync"

	"github.com/ebitengine/purego"
)

var (
	global *zstdLib
)

func InitWithSystemLibrary() (*zstdLib, error) {
	return InitWithPath("libzstd.so")
}

func InitWithPath(dlpath string) (*zstdLib, error) {
	libzstd, err := purego.Dlopen(dlpath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return nil, err
	}
	z := &zstdLib{
		libzstd: libzstd,
	}
	purego.RegisterLibFunc(&z._ZSTD_compress, libzstd, "ZSTD_compress")
	purego.RegisterLibFunc(&z._ZSTD_decompress, libzstd, "ZSTD_decompress")
	purego.RegisterLibFunc(&z._ZSTD_getFrameContentSize, libzstd, "ZSTD_getFrameContentSize")
	purego.RegisterLibFunc(&z._ZSTD_findFrameCompressedSize, libzstd, "ZSTD_findFrameCompressedSize")
	purego.RegisterLibFunc(&z._ZSTD_compressBound, libzstd, "ZSTD_compressBound")
	purego.RegisterLibFunc(&z._ZSTD_isError, libzstd, "ZSTD_isError")
	purego.RegisterLibFunc(&z._ZSTD_getErrorCode, libzstd, "ZSTD_getErrorCode")
	purego.RegisterLibFunc(&z._ZSTD_getErrorString, libzstd, "ZSTD_getErrorString")
	purego.RegisterLibFunc(&z._ZSTD_createCCtx, libzstd, "ZSTD_createCCtx")
	purego.RegisterLibFunc(&z._ZSTD_freeCCtx, libzstd, "ZSTD_freeCCtx")
	purego.RegisterLibFunc(&z._ZSTD_minCLevel, libzstd, "ZSTD_minCLevel")
	purego.RegisterLibFunc(&z._ZSTD_maxCLevel, libzstd, "ZSTD_maxCLevel")
	purego.RegisterLibFunc(&z._ZSTD_defaultCLevel, libzstd, "ZSTD_defaultCLevel")
	purego.RegisterLibFunc(&z._ZSTD_compressCCtx, libzstd, "ZSTD_compressCCtx")
	purego.RegisterLibFunc(&z._ZSTD_decompressDCtx, libzstd, "ZSTD_decompressDCtx")
	purego.RegisterLibFunc(&z._ZSTD_versionNumber, libzstd, "ZSTD_versionNumber")
	purego.RegisterLibFunc(&z._ZSTD_versionString, libzstd, "ZSTD_versionString")

	return z, nil
}

type zstdLib struct {
	libzstd uintptr

	_ZSTD_compress                func(dst []byte, dstCap int, src []byte, srcSize int, level int) int
	_ZSTD_decompress              func(dst []byte, dstCap int, src []byte, compressedSize int) int
	_ZSTD_getFrameContentSize     func(src []byte, srcSize int) int64
	_ZSTD_findFrameCompressedSize func(src []byte, srcSize int) int
	_ZSTD_compressBound           func(srcSize int) int
	_ZSTD_isError                 func(result int) bool
	_ZSTD_getErrorString          func(code ZSTD_ErrorCode) string
	_ZSTD_getErrorCode            func(result int) ZSTD_ErrorCode
	_ZSTD_minCLevel               func() int
	_ZSTD_maxCLevel               func() int
	_ZSTD_defaultCLevel           func() int
	_ZSTD_createCCtx              func() uintptr
	_ZSTD_freeCCtx                func(uintptr) int
	_ZSTD_compressCCtx            func(cctx uintptr, dst []byte, dstCap int, src []byte, srcSize int, level int) int
	_ZSTD_decompressDCtx          func(dctx uintptr, dst []byte, dstCap int, src []byte, srcSize int) int
	_ZSTD_versionNumber           func() uint
	_ZSTD_versionString           func() string
	_ZSTD_CCtx_reset              func() string
}

func (z *zstdLib) SetGlobal() {
	global = z
}

func g() *zstdLib {
	return global
}

func Compress(dst []byte, src []byte, level int) ([]byte, error) {
	bound := CompressBound(len(src))
	if cap(dst) >= bound {
		dst = dst[:bound]
	} else {
		dst = make([]byte, bound)
	}
	written := g()._ZSTD_compress(dst, cap(dst), src, len(src), level)
	if err := getError(written); err != nil {
		return nil, err
	}
	return dst[:written], nil
}

func Decompress(dst []byte, src []byte) ([]byte, error) {
	if len(src) == 0 {
		return []byte{}, ErrEmptySrc
	}
	bound := int(GetFrameContentSize(src))
	if cap(dst) >= bound {
		dst = dst[:cap(dst)]
	} else {
		dst = make([]byte, bound)
	}
	written := g()._ZSTD_decompress(dst, cap(dst), src, len(src))
	if err := getError(written); err != nil {
		return nil, err
	}
	return dst[:written], nil
}

func GetFrameContentSize(src []byte) int64 {
	return g()._ZSTD_getFrameContentSize(src, len(src))
}

func FindFrameCompressedSize(src []byte) int {
	return g()._ZSTD_findFrameCompressedSize(src, len(src))
}

func CompressBound(srcSize int) int {
	return g()._ZSTD_compressBound(srcSize)
}

func getError(ret int) error {
	if g()._ZSTD_isError(ret) {
		return Error{code: g()._ZSTD_getErrorCode(ret)}
	}
	return nil
}

func MinCLevel() int {
	return g()._ZSTD_minCLevel()
}

func MaxCLevel() int {
	return g()._ZSTD_maxCLevel()
}

func DefaultCLevel() int {
	return g()._ZSTD_defaultCLevel()
}

type CCtx struct {
	cctx uintptr
}

var (
	cctxPool = &sync.Pool{
		New: createCCtxAny,
	}
)

func CreateCCtx() *CCtx {
	cctx := &CCtx{
		cctx: g()._ZSTD_createCCtx(),
	}
	runtime.SetFinalizer(cctx, freeCCtx)
	return cctx
}

func createCCtxAny() any {
	return CreateCCtx()
}

func freeCCtx(cctx *CCtx) int {
	return g()._ZSTD_freeCCtx(cctx.cctx)
}

func (cctx *CCtx) Compress(dst []byte, src []byte, level int) ([]byte, error) {
	return cctx.compress(dst, src, level)
}

func (cctx *CCtx) compress(dst []byte, src []byte, level int) ([]byte, error) {
	written := g()._ZSTD_compressCCtx(cctx.cctx, dst, cap(dst), src, len(src), level)
	if err := getError(written); err != nil {
		return nil, err
	}
	dst = dst[:written]
	return dst, nil
}

type DCtx struct {
	dctx uintptr
}

func (dctx *DCtx) Decompress(dst []byte, src []byte) ([]byte, error) {
	written := g()._ZSTD_decompressDCtx(dctx.dctx, dst, cap(dst), src, len(src))
	if err := getError(written); err != nil {
		return nil, err
	}
	return dst[:written], nil
}

func VersionNumber() uint {
	return g()._ZSTD_versionNumber()
}

func VersionString() string {
	return g()._ZSTD_versionString()
}
