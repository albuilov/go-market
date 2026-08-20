FROM golang:1.26.4-alpine AS builder

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -ldflags="-s -w" \
    -o /out/gateway ./cmd


FROM alpine:3.22 AS runtime

RUN addgroup -S app && adduser -S -G app app

WORKDIR /app

COPY --from=builder /out/gateway ./gateway

USER app

EXPOSE 8080

ENTRYPOINT ["./gateway"]