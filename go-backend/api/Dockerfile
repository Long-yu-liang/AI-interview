# 第一阶段：构建 DocProcessor
FROM golang:1.24.4 AS doc-builder
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /build
COPY . .
WORKDIR /build/doc-processor
RUN CGO_ENABLED=0 GOOS=linux go build -o /output/doc-processor .

# 第二阶段：构建 API
FROM golang:1.24.4 AS api-builder
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /build
COPY . .
WORKDIR /build/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /output/api .

# 第三阶段：运行
FROM alpine:latest
RUN mkdir -p /etc/ssl/certs
COPY --from=api-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# 创建工作目录
WORKDIR /app

# 复制可执行文件
COPY --from=doc-builder /output/doc-processor /app/doc-processor
COPY --from=api-builder /output/api /app/api

# 复制配置
RUN mkdir -p /app/etc
COPY --from=api-builder /build/api/etc/chat-docker.yaml /app/etc/chat.yaml

# 设置执行权限
RUN chmod +x /app/doc-processor /app/api

WORKDIR /app
CMD ["./api", "-f", "etc/chat.yaml"]
