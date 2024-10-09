package libzstd

import (
	"bytes"
	"compress/gzip"
	"io"
)

func GetEmbeddedSharedObject() ([]byte, error) {
	gzr, err := gzip.NewReader(bytes.NewReader(libzstdBin))
	if err != nil {
		return nil, err
	}
	defer gzr.Close()
	return io.ReadAll(gzr)
}
