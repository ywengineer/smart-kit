

FROM golang:1.24.7-alpine AS builder
#
ARG VERSION=master
# use aliyun repository
RUN echo "https://mirrors.aliyun.com/alpine/v3.20/main/" > /etc/apk/repositories && \
    echo "https://mirrors.aliyun.com/alpine/v3.20/community/" >> /etc/apk/repositories

# install git
RUN apk update && \
    apk add --no-cache git && \
    rm -rf /var/cache/apk/*
# work directory
WORKDIR /app
#
ENV GOPROXY=https://goproxy.cn,direct
#
RUN git clone -b ${VERSION} --recurse-submodules https://gitee.com/ywengineer/smart-kit.git
#
RUN cd smart-kit/payment && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
#
RUN apk --no-cache add ca-certificates
#
WORKDIR /app
#
COPY --from=builder /app/smart-kit/payment/main .
#
RUN mkdir -p /app/logs && chmod 755 /app/logs
# expose application port
EXPOSE 8089
#
CMD ["./main"]
