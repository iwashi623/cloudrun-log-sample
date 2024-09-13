# ビルドコンテナの作成
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY . .

RUN go build -o myapp

# 実行コンテナの作成
FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/myapp .

CMD ["./myapp"]
