FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG VERSION=v3.1.8
ARG COMMIT=unknown
ARG BUILD_TIME=unknown
RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-s -w -X 'auraspeed/cmd/root.Version=${VERSION}' -X 'auraspeed/cmd/root.Commit=${COMMIT}' -X 'auraspeed/cmd/root.BuildTime=${BUILD_TIME}'" \
  -o auraspeed ./cmd/main.go

FROM alpine:3.20

RUN apk --no-cache add ca-certificates dumb-init tzdata

WORKDIR /app

COPY --from=builder /build/auraspeed .
COPY --from=builder /build/cmd/root/web.html .

ENV HOME=/root
ENV AURASPEED_AUTOUPDATE=false
ENV AURASPEED_LOGLEVEL=info
ENV AS_CONFIG_PATH=/app/config.toml

EXPOSE 59733

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["./auraspeed", "web", "--port", "59733"]
