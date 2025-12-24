
FROM alpine:latest
#
# 更换为国内镜像源（以阿里云为例）
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    cat /etc/apk/repositories && \
    apk update && \
    apk --no-cache add ca-certificates curl gcompat
#
WORKDIR /app
#
COPY ./out/payment .
#
RUN mkdir -p /app/logs && chmod 755 /app/logs && chmod +x ./payment
# expose application port
EXPOSE 8089
#
CMD ["./payment"]
