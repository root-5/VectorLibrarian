FROM golang:1.24.4-alpine

WORKDIR /app

COPY app/go.mod app/go.sum ./
RUN go mod download

COPY app/ ./

CMD ["tail", "-f", "/dev/null"]

# RUN go install github.com/air-verse/air@latest
# CMD ["air", "-c", ".air.toml"]
