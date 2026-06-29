# PIOC Config Uploader

PIOC配置上传工具 - 通过API将配置文件上传到PIOC系统创建新版本。

## 功能特性

- 通过YAML配置文件管理API参数
- 支持自定义版本号（默认使用当前日期YYYYMMDD格式）
- HMAC-SHA256签名认证
- 跨平台支持（Linux、Windows、macOS）

## 安装

### 从Release下载

访问 [Releases页面](https://github.com/na57/pioc-config-uploader/releases) 下载对应平台的二进制文件：

- `pioc-config-uploader-linux-amd64` - Linux 64位
- `pioc-config-uploader-windows-amd64.exe` - Windows 64位
- `pioc-config-uploader-darwin-amd64` - macOS Intel
- `pioc-config-uploader-darwin-arm64` - macOS Apple Silicon

下载后赋予执行权限（Linux/macOS）：
```bash
chmod +x pioc-config-uploader-*
```

### 从源码编译

```bash
git clone https://github.com/na57/pioc-config-uploader.git
cd pioc-config-uploader
go build -o pioc-config-uploader .
```

## 使用方法

### 1. 创建配置文件

复制示例配置文件并修改：

```bash
cp config.example.yaml config.yaml
```

编辑 `config.yaml`：

```yaml
# PIOC API 配置
api_key: "pk_your_api_key_here"
api_secret: "your_api_secret_here"
base_url: "http://localhost:8080"
config_id: "your-config-id-here"
config_file: "/path/to/your/config.conf"
```

### 2. 运行程序

```bash
./pioc-config-uploader config.yaml
```

### 配置参数说明

| 参数 | 必填 | 说明 |
|------|------|------|
| `api_key` | 是 | 您的API Key |
| `api_secret` | 是 | 您的API Secret |
| `base_url` | 否 | API基础URL，默认 `http://localhost:8080` |
| `config_id` | 是 | 要更新的配置ID |
| `config_file` | 是 | 要上传的配置文件路径 |

## API 文档

关于PIOC API的详细信息，请参考：

[设备运维人员API使用指南](https://github.com/na57/pioc/blob/production/设备运维人员API使用指南.md)

## 版本号规则

程序默认使用当前日期作为版本号，格式为 `YYYYMMDD`（例如：`20250629`）。

如需自定义版本号，可修改源码中的 `versionNumber` 变量。

## 示例输出

```
正在上传配置文件: /etc/nginx/nginx.conf
配置ID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
API地址: http://localhost:8080
✓ 版本创建成功!
  版本ID: e842458f-54cf-472e-81d1-f055f95e1d90
  版本号: 20250629
```

## 技术栈

- Go 1.21+
- HMAC-SHA256 签名认证
- YAML 配置解析

## License

MIT License
