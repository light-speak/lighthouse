# 使用官方的 Golang 镜像作为构建环境
FROM golang:1.23.1-alpine3.20 AS builder

# 设置工作目录
WORKDIR /app

# 复制项目文件
COPY . .
RUN go mod download

# 构建可执行文件
RUN go build -o /example .

# 使用一个更小的镜像作为运行环境
FROM alpine:3.19

# 安装 ca-certificates
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 从构建环境复制可执行文件到运行环境
COPY --from=builder /example .
COPY --from=builder /app/go.mod .
COPY --from=builder /app/.env .
# 暴露服务端口（根据你的应用需要调整）
EXPOSE 4001

# 运行可执行文件
CMD ["./example"]