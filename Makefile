#!/usr/bin/make -f

BUILD_OUTPUT ?= ./out
brunch ?= "latest"
APP_NAME ?= "uhub.service.ucloud.cn/infran/smart-payment"
APP_EXAMPLE_DIR = ""

define check_cmd
	@command -v $(1) >/dev/null 2>&1 || { \
		echo "错误: 命令$(1) 不存在, 请先安装该命令! example: make install-sd"; \
		exit 1; \
	}
endef

check-sd:
	$(call check_cmd, sd)

build-payment:
	cd payment && ./build.sh
	#CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BUILD_OUTPUT)/payment ./payment/

docker-build-payment: check-sd
	$(eval APP_NAME ?= smart-payment)
	$(eval APP_EXAMPLE_DIR=example/payment)
	@echo "build payment docker image with tag $(brunch). APP_NAME=$(APP_NAME) APP_EXAMPLE_DIR=$(APP_EXAMPLE_DIR)"
	@sudo rm -fr $(APP_EXAMPLE_DIR)
	sudo docker build --no-cache --build-arg VERSION=$(brunch) -t $(APP_NAME):$(brunch) -f Dockerfile.payment.dockerfile .
	@mkdir -p $(APP_EXAMPLE_DIR)
	@cp -f payment/*.yaml $(APP_EXAMPLE_DIR)/
	@sd "APP_NAME" "$(APP_NAME)" $(APP_EXAMPLE_DIR)/docker-compose.yaml
	@sd "VERSION"  "$(brunch)"   $(APP_EXAMPLE_DIR)/docker-compose.yaml

docker-build-payment-local: check-sd
	$(eval APP_NAME ?= smart-payment)
	$(eval APP_EXAMPLE_DIR=example/payment)
	@echo "build payment docker image with tag $(brunch). APP_NAME=$(APP_NAME) APP_EXAMPLE_DIR=$(APP_EXAMPLE_DIR)"
	@sudo rm -fr $(APP_EXAMPLE_DIR)
	sudo docker build --no-cache --build-arg VERSION=$(brunch) -t $(APP_NAME):$(brunch) -f Dockerfile.payment.local.dockerfile .
	@mkdir -p $(APP_EXAMPLE_DIR)
	@cp -f payment/*.yaml $(APP_EXAMPLE_DIR)/
	@sd "APP_NAME" "$(APP_NAME)" $(APP_EXAMPLE_DIR)/docker-compose.yaml
	@sd "VERSION"  "$(brunch)"   $(APP_EXAMPLE_DIR)/docker-compose.yaml

clean:
	rm -fr $(BUILD_OUTPUT)

install-sd:
	@sudo apt install wget
	@wget https://github.com/chmln/sd/releases/download/v1.0.0/sd-v1.0.0-x86_64-unknown-linux-gnu.tar.gz
	@tar -xzf sd-v1.0.0-x86_64-unknown-linux-gnu.tar.gz
	@sudo mv sd-v1.0.0-x86_64-unknown-linux-gnu/sd /usr/local/bin/
	@rm -fr sd-v1.0.0-x86_64-unknown-linux-gnu*


help:
	@echo "可用命令:"
	@echo "  make build-payment           							- 本地构建payment应用"
	@echo "  make docker-build-payment brunch=latest   				- 构建payment Docker镜像"
	@echo "  make docker-build-payment-local brunch=v0.1.1   		- 使用本地构建 payment 制作 Docker镜像"
	@echo "  make clean           									- 清理本地构建产物"
	@echo "  make install-sd           								- 安装sd工具"
	@echo "  make help            									- 显示帮助信息"
