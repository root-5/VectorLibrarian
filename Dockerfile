# ビルドステージ
FROM golang:1.23.0-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /main main.go

# 実行ステージ
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /main .
# log と cache ディレクトリを作成
RUN mkdir -p /root/log /root/cache
CMD ["./main"]
