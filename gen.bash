#!/usr/bin/env bash

set -e -u -o pipefail

function gen_enum() {
	rg --multiline "typedef enum .*(\n.*)+.*$1;" "$2" |
		perl -0pe 's{/\*.*?\*/}{}gs' |
		grep -o 'ZSTD_.*=.*' |
		go run ./internal/cenum2go "$1" |
		tee "${1,,}_gen.go"
}

gen_enum ZSTD_strategy ./tmp/zstd.h
gen_enum ZSTD_ErrorCode ./tmp/zstd_errors.h
