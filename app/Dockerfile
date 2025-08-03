FROM golang:1.24.4-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# 何もせずコンテナを起動する
CMD ["tail", "-f", "/dev/null"]

# RUN go install github.com/air-verse/air@latest
# CMD ["air", "-c", ".air.toml"]
