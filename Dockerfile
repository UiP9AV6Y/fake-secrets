FROM golang:1.26.0 AS build-cache

WORKDIR /build

COPY go.mod go.sum ./

ENV GOCACHE=/var/cache/go/src \
    GOMODCACHE=/var/cache/go/mod \
    GOBIN=/target
RUN set -xe ; \
    go mod download -x \
    && go install github.com/bwplotka/bingo@v0.10.0

FROM build-cache AS build

COPY . .

ENV CGO_ENABLED=0
RUN set -xe ; \
    go generate ./... \
    go install -ldflags="-s -w" ./...

FROM gcr.io/distroless/base-debian13 AS release

WORKDIR /data
VOLUME /data

COPY --from=build /target/ /usr/local/bin/

ENV TZ=UTC \
    FAKE_SECRETS_LOG_LEVEL=info

ENTRYPOINT ["/usr/local/bin/fake-secrets"]

HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
  CMD ["/usr/local/bin/fake-secrets", "healthcheck"]
