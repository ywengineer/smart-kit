#!/bin/bash

# 要检查的 Go 命令
export PATH=$PATH:$GOPATH/bin
hzcmd="hz"
repo="github.com/cloudwego/hertz/cmd/hz@latest"
mod="gitee.com/ywengineer/smart-kit/payment"

# 检查命令是否存在
if ! command -v "$hzcmd" &> /dev/null; then
    echo "$hzcmd 未找到，开始使用 go install 安装..."
    # 这里假设命令对应的 Go 包路径，你需要根据实际情况修改
    go install github.com/cloudwego/thriftgo@latest
    go install "$repo"
    if [ $? -eq 0 ]; then
        echo "$hzcmd 安装成功。"
    else
        echo "安装 $hzcmd 失败，请检查错误信息。"
        exit 1
    fi
else
    echo "$hzcmd 已安装。"
fi

for i in $(pwd)/idl/*.thrift ; do
  if [[ $i == *"common"* ]]; then
    echo "* skip generate api: $i"
  else
    echo "开始更新服务定义: hz update --idl $i"
    hz update --snake_tag --idl "$i" --mod $mod
  fi
done

if [ $? -eq 0 ]; then
  echo "更新完成"
else
  echo "更新失败"
fi
