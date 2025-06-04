FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build -ldflags="-s -w" -o limitlink .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/limitlink /root/
COPY --from=builder /app/templates /root/templates
COPY --from=builder /app/static /root/static


EXPOSE 8080

CMD ["/root/limitlink"]
