FROM docker.io/library/golang:1.26.3@sha256:2981696eed011d747340d7252620932677929cce7d2d539602f56a8d7e9b660b AS build-cache

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

FROM gcr.io/distroless/base-debian13:latest@sha256:c83f022002fc917a92501a8c30c605efdad3010157ba2c8998a2cbf213299201 AS release

WORKDIR /data
VOLUME /data

COPY --from=build /target/ /

ENV TZ=UTC \
    FAKE_SECRETS_STORAGE_DIR=/data \
    FAKE_SECRETS_LOG_LEVEL=info

ENTRYPOINT ["/usr/local/bin/fake-secrets"]

HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
  CMD ["/usr/local/bin/fake-secrets", "healthcheck"]
