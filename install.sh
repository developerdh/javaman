#!/bin/bash

# Linux/Mac平台安装脚本

# 获取系统信息
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

# 将amd64架构映射为x86_64
if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
fi

# 获取最新release信息
REPO="developerdh/javaman"
echo "Fetching latest version information..."
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO/releases/latest")
VERSION=$(echo "$LATEST_RELEASE" | grep -o '"tag_name": "[^"]*' | cut -d'"' -f4)

# 创建安装目录
INSTALL_DIR="$HOME/.javaman"
mkdir -p "$INSTALL_DIR"

# 构建下载文件名
ASSET_NAME="javaman_${OS}_${ARCH}.tar.gz"

# 获取下载URL
DOWNLOAD_URL=$(echo "$LATEST_RELEASE" | grep -o "\"browser_download_url\": \"[^\"]*${ASSET_NAME}\"" | cut -d'"' -f4)

if [ -n "$DOWNLOAD_URL" ]; then
    echo "Downloading javaman $VERSION..."
    curl -L "$DOWNLOAD_URL" -o "$INSTALL_DIR/javaman.tar.gz"
    
    # 解压文件
    tar -xzf "$INSTALL_DIR/javaman.tar.gz" -C "$INSTALL_DIR"
    rm "$INSTALL_DIR/javaman.tar.gz"
    
    # 添加执行权限
    chmod +x "$INSTALL_DIR/javaman"
    
    # 配置环境变量
    SHELL_CONFIG=""
    if [ -n "$ZSH_VERSION" ]; then
        SHELL_CONFIG="$HOME/.zshrc"
    elif [ -n "$BASH_VERSION" ]; then
        SHELL_CONFIG="$HOME/.bashrc"
    fi
    
    if [ -n "$SHELL_CONFIG" ]; then
        if ! grep -q "$INSTALL_DIR" "$SHELL_CONFIG"; then
            echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$SHELL_CONFIG"
            echo "Javaman has been successfully installed!"
            echo "Please run 'source $SHELL_CONFIG' or reopen the terminal to apply the environment variables"
        fi
    else
        echo "Please manually add $INSTALL_DIR to your PATH environment variable"
    fi
else
    echo "Error: No corresponding download file found"
    exit 1
fi
echo "Press any key to continue..."
read -n 1 -s