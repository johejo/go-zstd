#!/usr/bin/env bash

set -e -u -o pipefail

zstd_version="$(cat zstd_version.txt)"

zig_target="$1"
os_arch="$2"

docker run --rm --mount "type=bind,src=${PWD}/internal/libzstd,dst=/work/out" go-zstd-builder \
	sh -c "make lib-nomt -j$(nproc) CC=\"zig cc -s -target ${zig_target}\" AR='zig ar' && cp ./lib/libzstd.so.${zstd_version} /work/out/libzstd_${os_arch}.so"

gzip "./internal/libzstd/libzstd_${os_arch}.so"

build_contraint="${os_arch/_/ \&\& }"

cat <<EOF >"./internal/libzstd/libzstd_${os_arch}.go"
//go:build ${build_contraint}
package libzstd

import _ "embed"

//go:embed libzstd_${os_arch}.so.gz
var libzstdBin []byte
EOF
