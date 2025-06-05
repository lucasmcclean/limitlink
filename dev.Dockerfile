FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build -ldflags="-s -w" -o limitlink .

EXPOSE 8080

CMD ["/app/limitlink"]
