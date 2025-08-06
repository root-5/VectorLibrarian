FROM golang:1.24.4-alpine

WORKDIR /app

# おまじない (Alpine用)
RUN apk update && \
    apk add --no-cache \
    build-base \
    git \
    curl \
    coreutils

COPY go.mod go.sum ./
RUN go mod download

# 何もせずコンテナを起動する
CMD ["tail", "-f", "/dev/null"]

# RUN go install github.com/air-verse/air@latest
# CMD ["air", "-c", ".air.toml"]
