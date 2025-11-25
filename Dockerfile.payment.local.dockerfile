
FROM alpine:latest
#
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
