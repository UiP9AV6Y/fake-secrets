FROM docker.io/library/golang:1.27.2@sha256:5bc7f572bbaa98885a3a1fd9c0aa76b59e3e14e8628bfc316bbfd0c701e4818c AS build-cache

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
