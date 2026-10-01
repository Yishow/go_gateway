# 環境變數設置指南

## 快速開始

### 1. 創建 .env 文件

在項目根目錄創建 `.env` 文件，複製以下內容：

```env
# ============================================
# Go Gateway 環境變數配置
# ============================================

# 後端服務配置
PORT=8080
HOST=
GIN_MODE=release

# API 配置
API_BASE_PATH=/api/v1

# CORS 配置
CORS_ALLOW_ORIGINS=*

# 日誌配置
LOG_LEVEL=info
LOG_OUTPUT=console
LOG_FILE=./logs/gateway.log

# WebSocket 配置
WS_READ_TIMEOUT=60
WS_WRITE_TIMEOUT=10
WS_PING_INTERVAL=30

# 連接池配置
MAX_CONNECTIONS=100
CONNECTION_TIMEOUT=5
READ_TIMEOUT=2
WRITE_TIMEOUT=2

# 開發環境配置
DEBUG=false

# 前端配置（Vite 環境變數）
VITE_API_BASE_URL=/api/v1
VITE_DEV_PORT=5173
VITE_API_PROXY_TARGET=http://localhost:8080
```

### 2. 根據需要修改配置

根據實際環境修改 `.env` 文件中的配置值。

### 3. 啟動服務

```powershell
.\start.ps1 -QuickStart
```

## 配置說明

### 後端配置

所有後端配置都可以在 `.env` 文件中設置，詳見 [ENV_CONFIG.md](./ENV_CONFIG.md)。

### 前端配置

前端使用 Vite，需要在變數名前加上 `VITE_` 前綴：

- `VITE_API_BASE_URL`: API 基礎路徑（預設：`/api/v1`）
- `VITE_DEV_PORT`: Vite 開發伺服器端口（預設：`5173`）
- `VITE_API_PROXY_TARGET`: 後端服務地址，用於代理配置（預設：`http://localhost:8080`）

## 注意事項

1. `.env` 文件不會被提交到版本控制（已在 `.gitignore` 中）
2. 請勿將包含敏感信息的 `.env` 文件提交到版本控制
3. 修改 `.env` 後需要重啟服務才能生效
