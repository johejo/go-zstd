package zstd_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/johejo/go-zstd"
)

func Test(t *testing.T) {
	z, err := zstd.InitWithEmbed()
	if err != nil {
		t.Fatal(err)
	}
	z.SetGlobal()

	src := bytes.Repeat([]byte("foo"), 10)
	compressed, err := zstd.Compress(nil, src, 9)
	if err != nil {
		t.Fatal(err)
	}

	decompressed, err := zstd.Decompress(nil, compressed)
	if err != nil {
		t.Fatal(err)
	}
	assert(t, src, decompressed)
	if !reflect.DeepEqual(src, decompressed) {
		t.Errorf("src=%s, decompressed=%s", src, decompressed)
	}

	frameContentSize := zstd.GetFrameContentSize(compressed)
	if frameContentSize != int64(len(src)) {
		t.Errorf("frameContentSize=%d, len(src)=%d", frameContentSize, len(src))
	}

	frameCompressedSize := zstd.FindFrameCompressedSize(compressed)
	if frameCompressedSize != len(compressed) {
		t.Errorf("frameCompressedSize=%d, len(src)=%d", frameCompressedSize, len(compressed))
	}

	t.Log(zstd.CompressBound(len(src)))

	t.Log(zstd.DefaultCLevel())
	t.Log(zstd.VersionNumber())
	t.Log(zstd.VersionString())
}

func TestMinCLevel(t *testing.T) {
	assert(t, -(1<<17), zstd.MinCLevel())
}

func TestMaxCLevel(t *testing.T) {
	assert(t, 22, zstd.MaxCLevel())
}

func TestDefaultCLevel(t *testing.T) {
	assert(t, 3, zstd.DefaultCLevel())
}

func assert[T any](t *testing.T, want T, got T) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want=%v, got=%v", want, got)
	}
}
