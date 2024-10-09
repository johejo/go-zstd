//go:build linux && arm64
package libzstd

import _ "embed"

//go:embed libzstd_linux_arm64.so.gz
var libzstdBin []byte
