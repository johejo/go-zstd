.PHONY: builder

gen:
	./gen.bash

clean:
	rm -f *_gen.go

lib-clean:
	rm -f ./internal/libzstd/*.so.gz

builder:
	docker build -t go-zstd-builder --build-arg ZSTD_VERSION=$(shell cat zstd_version.txt) ./builder

lib-linux-amd64:
	./build-libzstd.bash "x86_64-linux-gnu" "linux_amd64"

lib-linux-arm64:
	./build-libzstd.bash "aarch64-linux-gnu" "linux_arm64"

lib-macos-arm64:
	./build-libzstd.bash "aarch64-macos-none" "dawrin_arm64"

test-linux-arm64:
	GOOS=linux GOARCH=arm64 go test -c
	docker run --rm --mount type=bind,src=$(PWD)/go-zstd.test,dst=/go-zstd.test bitnami/minideb:bookworm-arm64 /go-zstd.test -test.v
