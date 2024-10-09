//go:build linux

package zstd

import (
	"fmt"
	"os"

	"github.com/johejo/go-zstd/internal/libzstd"
	"github.com/justincormack/go-memfd"
)

func InitWithEmbed() (*zstdLib, error) {
	b, err := libzstd.GetEmbeddedSharedObject()
	if err != nil {
		return nil, err
	}
	return InitWithBin(b)
}

func InitWithBin(bin []byte) (*zstdLib, error) {
	mfd, err := memfd.Create()
	if err != nil {
		return nil, err
	}
	defer mfd.Close()
	if _, err := mfd.Write(bin); err != nil {
		return nil, err
	}
	dlpath := fmt.Sprintf("/proc/%d/fd/%d", os.Getpid(), mfd.Fd())
	return InitWithPath(dlpath)
}
