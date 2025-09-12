BUILD_DIR=./out
# 构建Go应用（本地）
build-payment:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BUILD_DIR)/payment ./payment/

# 构建Docker镜像
docker-build-payment:
	docker build -t smart-payment -f Dockerfile.payment .

# 清理本地构建产物
clean:
	rm -f $(BUILD_DIR)/*

# 显示帮助信息
help:
	@echo "可用命令:"
	@echo "  make build-payment           	- 本地构建payment应用"
	@echo "  make docker-build-payment    	- 构建payment Docker镜像"
	@echo "  make clean           			- 清理本地构建产物"
	@echo "  make help            			- 显示帮助信息"
