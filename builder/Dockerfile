FROM bitnami/minideb:bookworm
RUN install_packages curl ca-certificates xz-utils make
WORKDIR /work
ARG ZIG_VERSION=0.13.0
RUN curl -fsSL https://ziglang.org/download/${ZIG_VERSION}/zig-linux-x86_64-${ZIG_VERSION}.tar.xz | tar xJ
RUN ln -sf /work/zig-linux-x86_64-${ZIG_VERSION}/zig /usr/local/bin/zig
ARG ZSTD_VERSION
RUN curl -fsSL https://github.com/facebook/zstd/releases/download/v${ZSTD_VERSION}/zstd-${ZSTD_VERSION}.tar.gz | tar xz
WORKDIR /work/zstd-${ZSTD_VERSION}
