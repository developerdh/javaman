# Javaman - Java版本管理器

Javaman 是一个简单的 JDK 版本管理工具，支持 Windows、Linux 和 macOS 系统。它可以帮助开发者轻松地在不同版本的 JDK 之间切换，自动管理环境变量。

**特别说明**: 
- 本项目中的内容是在不熟悉 golang 的基础上由 AI 自动生成（除本文档中部分内容外），可能存在不完善之处，欢迎社区贡献改进。

## 特性

- 跨平台支持 (Windows, Linux, macOS)
- 支持多版本JDK管理
- 自动配置环境变量
- 简单易用的命令行界面


## 安装

### 从源码构建

1. 确保您已安装Go 1.21或更高版本
2. 克隆仓库：
```bash
git clone https://github.com/developerdh/javaman.git
cd javaman
```

3. 构建项目：
```bash
go build
# 或
go build linux 
```

### 使用预编译的二进制文件

```shell
# Windows
powershell -ExecutionPolicy Bypass -File install.ps1
```

```bash
# Linux/Mac
bash install.sh
```

访问 [Releases](https://github.com/developerdh/javaman/releases) 页面下载适合您系统的版本

### 配置环境变量（可选）
- Windows: 为`javaman.exe`添加Path环境变量
- Linux/macOS: 
```bash
sudo mv javaman /usr/bin/
sudo chmod +x /usr/bin/javaman
```

## 使用方法

### 查看已安装的JDK版本
```bash
javaman list
# 或
javaman ls
```

### 切换JDK版本
```bash
javaman use <version>
# 例如：javaman use 17
```

### 查看当前使用的版本
```bash
javaman current
```

### 添加新的JDK安装
```bash
javaman add <JDK安装路径>
# 例如：javaman add "C:\Program Files\Java\jdk-17"
```

### 删除JDK版本
``remove``或``rm``命令只会删除配置，不会删除实际的JDK安装。
```bash
javaman remove <version>
# 或
javaman rm <version>
```

## 配置文件

配置文件位于用户目录下的 `.javaman/config.toml`：
- Windows: `C:\Users\<username>\.javaman\config.toml`
- Linux/macOS: `~/.javaman/config.toml`

## 权限要求

- Windows: 需要管理员权限以修改系统环境变量
- Linux/macOS: 需要sudo权限以修改系统级配置，或者使用用户级配置文件

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 致谢

- 感谢 AI 技术的发展与应用，使得自动化代码生成成为可能。
- 受 nvm (Node Version Manager) 项目的启发
