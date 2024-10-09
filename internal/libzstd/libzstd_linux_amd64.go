//go:build linux && amd64
package libzstd

import _ "embed"

//go:embed libzstd_linux_amd64.so.gz
var libzstdBin []byte
