ARG BUILD_ARG_GO_VERSION=1.22
ARG BUILD_ARG_ALPINE_VERSION=3.18
FROM golang:${BUILD_ARG_GO_VERSION}-alpine${BUILD_ARG_ALPINE_VERSION} AS builder

WORKDIR /src
COPY . .

RUN go build -o bin/app

FROM alpine:${BUILD_ARG_ALPINE_VERSION}
WORKDIR /
COPY --from=builder /src/bin/app .
EXPOSE 8000

CMD ["/app", "rest"]