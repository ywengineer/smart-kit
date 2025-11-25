
FROM alpine:latest
#
# 更换为国内镜像源（以阿里云为例）
RUN echo "https://mirrors.aliyun.com/alpine/v3.22/main/" > /etc/apk/repositories && \
    echo "https://mirrors.aliyun.com/alpine/v3.22/community/" >> /etc/apk/repositories \

RUN apk --no-cache add ca-certificates
#
WORKDIR /app
#
COPY ./payment/output/bin/payment .
#
RUN mkdir -p /app/logs && chmod 755 /app/logs && chmod +x ./payment
# expose application port
EXPOSE 8089
#
CMD ["./payment"]
