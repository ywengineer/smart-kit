APP_NAME := myapp
IMAGE_NAME := $(APP_NAME):latest
PORT := 8080
LOCAL_CONFIG := ./application.yaml
LOCAL_LOGS := ./logs

# 构建Go应用（本地）
build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 构建Docker镜像
docker-build:
	docker build -t $(IMAGE_NAME) .

# 运行Docker容器
docker-run:
	mkdir -p $(LOCAL_LOGS)
	docker run -d \
		-p $(PORT):$(PORT) \
		-v $(LOCAL_CONFIG):/app/application.yaml \
		-v $(LOCAL_LOGS):/app/logs \
		--name $(APP_NAME) \
		$(IMAGE_NAME)

# 停止并删除容器
docker-stop:
	docker stop $(APP_NAME) || true
	docker rm $(APP_NAME) || true

# 查看容器日志
docker-logs:
	docker logs -f $(APP_NAME)

# 进入容器内部
docker-exec:
	docker exec -it $(APP_NAME) sh

# 清理本地构建产物
clean:
	rm -f main
	rm -rf $(LOCAL_LOGS)/*

# 重启容器（先停止，再构建镜像，再运行）
restart: docker-stop docker-build docker-run

# 显示帮助信息
help:
	@echo "可用命令:"
	@echo "  make build           - 本地构建Go应用"
	@echo "  make docker-build    - 构建Docker镜像"
	@echo "  make docker-run      - 运行Docker容器"
	@echo "  make docker-stop     - 停止并删除容器"
	@echo "  make docker-logs     - 查看容器日志"
	@echo "  make docker-exec     - 进入容器内部"
	@echo "  make clean           - 清理本地构建产物"
	@echo "  make restart         - 重启容器（停止->构建->运行）"
	@echo "  make help            - 显示帮助信息"
