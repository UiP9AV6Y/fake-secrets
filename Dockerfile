FROM docker.io/library/golang:1.27.2@sha256:820153096fda4de6ac1b7194d12e416b7dd8f4e9a33d17e0345b111b945162d8 AS build-cache

WORKDIR /build

COPY .bingo/ .bingo/
COPY .changes/ .changes/
COPY go.mod go.sum Makefile ./

ENV GOCACHE=/var/cache/go/src \
    GOMODCACHE=/var/cache/go/mod
RUN set -xe ; \
    make dependencies tools

FROM build-cache AS build

COPY . .

ENV CGO_ENABLED=0
RUN set -xe ; \
    make build install DESTDIR=/target

FROM gcr.io/distroless/base-debian13:latest@sha256:f5e95195e998742ff590b7569bae479accce1226e51ae6c3118954030121cdcb AS release

WORKDIR /data
VOLUME /data

COPY --from=build /target/ /

ENV TZ=UTC \
    FAKE_SECRETS_STORAGE_DIR=/data \
    FAKE_SECRETS_LOG_LEVEL=info

ENTRYPOINT ["/usr/local/bin/fake-secrets"]

HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
  CMD ["/usr/local/bin/fake-secrets", "healthcheck"]
