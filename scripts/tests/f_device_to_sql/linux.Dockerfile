# Disposable Linux acceptance tools; application source is mounted read-only.
FROM golang:1.25.5-bookworm@sha256:d9132cce84391efab786495288756d60e1da215b1f94e87860aeefc3d4c45b6d AS toolchain
FROM mcr.microsoft.com/playwright:v1.58.2-noble@sha256:6446946a1d9fd62d9ae501312a2d76a43ee688542b21622056a372959b65d63d
COPY --from=toolchain /usr/local/go /usr/local/go
RUN apt-get update \
 && apt-get install -y --no-install-recommends sqlite3 make gcc g++ git ca-certificates \
 && rm -rf /var/lib/apt/lists/*
ENV PATH="/usr/local/go/bin:${PATH}" GOTOOLCHAIN=local GIT_OPTIONAL_LOCKS=0
WORKDIR /work
