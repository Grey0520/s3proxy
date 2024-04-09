FROM golang:1.22-alpine AS builder

ENV GOPROXY=https://proxy.golang.org,direct


WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -a -o main cmd/s3proxy/main.go

FROM scratch

COPY --from=builder /app/main /main
COPY --from=builder /app/configs/config.yaml.local /config.yaml
# 设置应用程序为 Docker 镜像的默认启动命令
ENTRYPOINT ["/main", "--config", "/config.yaml"]
