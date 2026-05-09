FROM docker.io/library/golang:1.25-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o auraspeed ./cmd/main.go

FROM docker.io/library/alpine:3.20

RUN apk --no-cache add ca-certificates dumb-init

WORKDIR /app

COPY --from=builder /build/auraspeed .
COPY --from=builder /build/cmd/root/web.html .

ENV HOME=/root
ENV AURASPEED_AUTOUPDATE=false

EXPOSE 8080

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["./auraspeed", "web", "-p", "8080"]