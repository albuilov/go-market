FROM golang:1.26.4-alpine AS builder

WORKDIR /src

COPY . .

RUN --mount=type=cache,id=go-market-modules,target=/go/pkg/mod \
    --mount=type=cache,id=go-market-build,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -ldflags="-s -w" \
    -o /out/gateway ./services/gateway/cmd

FROM alpine:3.22 AS runtime

RUN addgroup -S app && adduser -S -G app app

WORKDIR /app

COPY --from=builder /out/gateway ./gateway

USER app

ENTRYPOINT ["./gateway"]
