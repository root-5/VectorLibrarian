FROM golang:1.24.4-alpine

WORKDIR /app

# おまじない
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    build-essential \
    git \
    curl \
    coreutils \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

# 何もせずコンテナを起動する
CMD ["tail", "-f", "/dev/null"]

# RUN go install github.com/air-verse/air@latest
# CMD ["air", "-c", ".air.toml"]
