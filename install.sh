#!/bin/bash

VERSION="latest"
GITHUB_REPO="fflow-tech/fflow"

# 检测操作系统和架构
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# 映射架构名称
case $ARCH in
  x86_64)
    ARCH="amd64"
    ;;
  aarcharm64)
    ARCH="arm64"
    ;;
esac

# 构建下载 URL
if [ "$VERSION" = "latest" ]; then
  DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/latest/download/fflow-cli_${OS}_${ARCH}"
else
  DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/download/${VERSION}/fflow-cli_${OS}_${ARCH}"
fi

# 下载并安装
echo "正在下载 fflow-cli..."
curl -L "$DOWNLOAD_URL" -o /tmp/fflow-cli
chmod +x /tmp/fflow-cli

# 移动到用户的 bin 目录
INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
  echo "需要管理员权限安装到 ${INSTALL_DIR}"
  sudo mv /tmp/fflow-cli "$INSTALL_DIR/fflow-cli"
else
  mv /tmp/fflow-cli "$INSTALL_DIR/fflow-cli"
fi

echo "fflow-cli 已成功安装到 ${INSTALL_DIR}/fflow-cli"