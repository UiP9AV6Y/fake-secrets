FROM docker.io/library/golang:1.26.5@sha256:2005724102f45917a63e9d092fc0e4ea56ea575048ce147caad5f5f61502c365 AS build-cache

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

FROM gcr.io/distroless/base-debian13:latest@sha256:20dc7edae3f7efe09b934aca4b347b00bb4ae0f2864b6131771687ae6d54891f AS release

WORKDIR /data
VOLUME /data

COPY --from=build /target/ /

ENV TZ=UTC \
    FAKE_SECRETS_STORAGE_DIR=/data \
    FAKE_SECRETS_LOG_LEVEL=info

ENTRYPOINT ["/usr/local/bin/fake-secrets"]

HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
  CMD ["/usr/local/bin/fake-secrets", "healthcheck"]
