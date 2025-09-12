#!/usr/bin/make -f

BUILD_OUTPUT ?= ./out
version ?= latest
APP_NAME ?= app
APP_EXAMPLE_DIR = ""

build-payment:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BUILD_OUTPUT)/payment ./payment/

docker-build-payment:
	$(eval APP_NAME=smart-payment)
	$(eval APP_EXAMPLE_DIR=example/payment)
	@echo "build payment docker image with tag $(version). APP_NAME=$(APP_NAME) APP_EXAMPLE_DIR=$(APP_EXAMPLE_DIR)"
	@rm -fr $(APP_EXAMPLE_DIR)
	sudo docker build -t $(APP_NAME):$(version) -f Dockerfile.payment.dockerfile .
	@mkdir -p $(APP_EXAMPLE_DIR)
	@cp -f payment/*.yaml $(APP_EXAMPLE_DIR)/
	@sed -i 's/APP_NAME/$(APP_NAME)/g' $(APP_EXAMPLE_DIR)/docker-compose.yaml
	@sed -i 's/VERSION/$(version)/g' $(APP_EXAMPLE_DIR)/docker-compose.yaml

clean:
	rm -f $(BUILD_OUTPUT)/*


help:
	@echo "可用命令:"
	@echo "  make build-payment           					- 本地构建payment应用"
	@echo "  make docker-build-payment PAY_VER=latest   	- 构建payment Docker镜像"
	@echo "  make clean           							- 清理本地构建产物"
	@echo "  make help            							- 显示帮助信息"
