#!/bin/sh

CompileDaemon \
  -include="*.go" \
  -build="go build -o /app/serve ./cmd/limitlink" \
  -command="/app/serve" &

pnpm --dir /app/frontend run build:watch &

wait
