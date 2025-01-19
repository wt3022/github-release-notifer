# ビルドステージ
FROM golang:1.22-bullseye AS builder

ENV TZ=Asia/Tokyo
ENV GONOSUMDB=*
ENV GOSUMDB=off
ENV GOPROXY=direct

WORKDIR /api

RUN apt-get update && apt-get install -y \
    gcc \
    g++ \
    sqlite3 \
    curl \
    git \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

COPY . .
RUN go mod download && \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o main && \
    chmod +x main

# ============ 実行ステージ ============
FROM debian:bullseye-slim AS runner

ENV TZ=Asia/Tokyo

WORKDIR /api

RUN apt-get update && apt-get install -y \
    sqlite3 \
    curl \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /api/main .
COPY .env .

CMD ["./main"]
