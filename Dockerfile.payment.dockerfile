ARG VERSION=master
# 第一阶段：构建阶段
FROM golang:1.24.7-alpine AS builder
# 设置构建阶段的工作目录
WORKDIR /app
# 配置国内Go模块代理，加速依赖下载
ENV GOPROXY=https://goproxy.cn,direct
# 复制源代码
RUN git clone -b ${VERSION} --recurse-submodules https://gitee.com/ywengineer/smart-kit.git
# 复制依赖文件并下载
RUN cd smart-kit/payment && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

# 安装必要工具
RUN apk --no-cache add ca-certificates

# 设置工作目录为/app
WORKDIR /app

# 从构建阶段复制编译好的应用到/app目录
COPY --from=builder /app/smart-kit/payment/main .
# 创建日志目录并设置适当的权限
RUN mkdir -p /app/logs && chmod 755 /app/logs
# 暴露应用端口
EXPOSE 8089
# 运行应用
CMD ["./main"]
