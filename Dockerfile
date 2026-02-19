FROM docker.io/library/golang:1.26.0@sha256:21d68a382b318dde54771baf5e04393e0930d077583c8010e4fd83adae328eda AS build-cache

WORKDIR /build

COPY .bingo/ .bingo/
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

FROM gcr.io/distroless/base-debian13:latest@sha256:2a1bdb588ba3d6096fe12705ca67d7b04a6319dd16a456c79dbfb7e9e0e780aa AS release

WORKDIR /data
VOLUME /data

COPY --from=build /target/ /

ENV TZ=UTC \
    FAKE_SECRETS_STORAGE_DIR=/data \
    FAKE_SECRETS_LOG_LEVEL=info

ENTRYPOINT ["/usr/local/bin/fake-secrets"]

HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
  CMD ["/usr/local/bin/fake-secrets", "healthcheck"]
