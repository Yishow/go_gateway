# 環境變數配置說明

## 概述

Go Gateway 使用 `.env` 文件統一管理前後端的配置。所有配置都可以通過環境變數或 `.env` 文件設置。

## 快速開始

1. **複製範例文件**
   ```bash
   cp .env.sample .env
   ```

2. **編輯配置**
   根據實際環境修改 `.env` 文件中的配置值

3. **啟動服務**
   ```bash
   .\start.ps1 -QuickStart
   ```

## 配置說明

### 後端服務配置

| 變數名 | 說明 | 預設值 | 範例 | 狀態 |
|--------|------|--------|------|------|
| `PORT` | 服務監聽端口 | `8080` | `8080` | ✅ 已實現 |
| `HOST` | 服務監聽地址（空表示所有接口） | 空 | `127.0.0.1` | ✅ 已實現 |
| `GIN_MODE` | Gin 運行模式 | `release` | `debug`, `release`, `test` | ✅ 已實現 |

### API 配置

| 變數名 | 說明 | 預設值 | 範例 | 狀態 |
|--------|------|--------|------|------|
| `API_BASE_PATH` | API 基礎路徑 | `/api/v1` | `/api/v1` | ✅ 已實現 |

### CORS 配置

| 變數名 | 說明 | 預設值 | 範例 | 狀態 |
|--------|------|--------|------|------|
| `CORS_ALLOW_ORIGINS` | 允許的來源（逗號分隔） | `*` | `http://localhost:3000,http://localhost:8080` | ✅ 已實現 |

### 日誌配置

| 變數名 | 說明 | 預設值 | 範例 | 狀態 |
|--------|------|--------|------|------|
| `LOG_LEVEL` | 日誌級別 | `info` | `debug`, `info`, `warn`, `error` | ✅ 已實現 |
| `LOG_OUTPUT` | 日誌輸出位置 | `console` | `console`, `file` | ✅ 已實現 |
| `LOG_FILE` | 日誌文件路徑 | `./logs/gateway.log` | `./logs/gateway.log` | ✅ 已實現 |

### WebSocket 配置

| 變數名 | 說明 | 預設值 | 範例 | 狀態 |
|--------|------|--------|------|------|
| `WS_READ_TIMEOUT` | WebSocket 讀取超時（秒） | `60` | `60` | ✅ 已實現 |
| `WS_WRITE_TIMEOUT` | WebSocket 寫入超時（秒） | `10` | `10` | ✅ 已實現 |
| `WS_PING_INTERVAL` | WebSocket Ping 間隔（秒） | `30` | `30` | ✅ 已實現 |

### 連接池配置

| 變數名 | 說明 | 預設值 | 範例 | 狀態 |
|--------|------|--------|------|------|
| `MAX_CONNECTIONS` | 最大連接數 | `100` | `100` | ✅ 已實現 |
| `CONNECTION_TIMEOUT` | 連接超時（秒） | `5` | `5` | ✅ 已實現 |
| `READ_TIMEOUT` | 讀取超時（秒） | `2` | `2` | ✅ 已實現 |
| `WRITE_TIMEOUT` | 寫入超時（秒） | `2` | `2` | ✅ 已實現 |

### 開發環境配置

| 變數名 | 說明 | 預設值 | 範例 | 狀態 |
|--------|------|--------|------|------|
| `DEBUG` | 是否啟用調試模式（影響 Gin 模式和日誌詳細程度） | `false` | `true`, `false` | ✅ 已實現 |

## 前端配置

前端使用 Vite 環境變數，需要在變數名前加上 `VITE_` 前綴。

### 前端環境變數

| 變數名 | 說明 | 預設值 | 範例 | 狀態 |
|--------|------|--------|------|------|
| `VITE_API_BASE_URL` | API 基礎路徑 | `/api/v1` | `/api/v1` | ✅ 已實現 |
| `VITE_DEV_PORT` | Vite 開發伺服器端口 | `5173` | `5173` | ✅ 已實現 |
| `VITE_API_PROXY_TARGET` | 後端服務地址（用於代理） | `http://localhost:8080` | `http://localhost:8080` | ✅ 已實現 |

**注意**：
- 前端環境變數需要在 `.env` 文件中設置，並在構建時由 Vite 處理
- `VITE_API_PROXY_TARGET` 用於配置 Vite 開發伺服器的代理目標，將 `/api` 請求轉發到後端服務

## 計劃中的配置

以下配置項在 `.env.sample` 中已預留，但尚未在代碼中實現：

### 資料庫配置 [計劃中]

#### 通用配置

| 變數名 | 說明 | 預設值 | 範例 |
|--------|------|--------|------|
| `DB_DRIVER` | 資料庫驅動 | `sqlite` | `sqlite`, `mysql`, `postgres`, `sqlserver` |

#### SQLite 配置（當 `DB_DRIVER=sqlite` 時）

| 變數名 | 說明 | 預設值 | 範例 |
|--------|------|--------|------|
| `DB_PATH` | 資料庫文件路徑 | `./data/gateway.db` | `./data/gateway.db` |

#### MySQL 配置（當 `DB_DRIVER=mysql` 時）

| 變數名 | 說明 | 預設值 | 範例 |
|--------|------|--------|------|
| `DB_HOST` | 資料庫主機 | - | `localhost` |
| `DB_PORT` | 資料庫端口 | `3306` | `3306` |
| `DB_NAME` | 資料庫名稱 | - | `gateway` |
| `DB_USER` | 資料庫用戶名 | - | `gateway` |
| `DB_PASSWORD` | 資料庫密碼 | - | `password` |

#### PostgreSQL 配置（當 `DB_DRIVER=postgres` 時）

| 變數名 | 說明 | 預設值 | 範例 |
|--------|------|--------|------|
| `DB_HOST` | 資料庫主機 | - | `localhost` |
| `DB_PORT` | 資料庫端口 | `5432` | `5432` |
| `DB_NAME` | 資料庫名稱 | - | `gateway` |
| `DB_USER` | 資料庫用戶名 | - | `gateway` |
| `DB_PASSWORD` | 資料庫密碼 | - | `password` |
| `DB_SSLMODE` | SSL 模式 | `disable` | `disable`, `require`, `verify-ca`, `verify-full` |
| `DB_TIMEZONE` | 時區設定 | - | `Asia/Taipei`, `UTC` |

#### SQL Server 配置（當 `DB_DRIVER=sqlserver` 時）

| 變數名 | 說明 | 預設值 | 範例 |
|--------|------|--------|------|
| `DB_HOST` | 資料庫主機 | - | `localhost` |
| `DB_PORT` | 資料庫端口 | `1433` | `1433` |
| `DB_NAME` | 資料庫名稱 | - | `gateway` |
| `DB_USER` | 資料庫用戶名 | - | `sa` |
| `DB_PASSWORD` | 資料庫密碼 | - | `password` |
| `DB_INSTANCE` | SQL Server 實例名稱（可選） | - | `SQLEXPRESS` |

### MQTT 配置 [計劃中]

| 變數名 | 說明 | 預設值 | 範例 |
|--------|------|--------|------|
| `MQTT_BROKER` | MQTT Broker 地址 | - | `tcp://mqtt.example.com:1883` |
| `MQTT_CLIENT_ID` | MQTT 客戶端 ID | - | `gateway-001` |
| `MQTT_USERNAME` | MQTT 用戶名 | - | `username` |
| `MQTT_PASSWORD` | MQTT 密碼 | - | `password` |

### 任務調度配置 [計劃中]

| 變數名 | 說明 | 預設值 | 範例 |
|--------|------|--------|------|
| `TASK_WORKER_COUNT` | 任務執行器工作線程數 | `4` | `4` |
| `TASK_DEFAULT_TIMEOUT` | 預設任務超時時間（秒） | `10` | `10` |

## 配置優先順序

配置的優先順序如下（從高到低）：

1. **環境變數**（系統環境變數）
2. **`.env` 文件**（項目根目錄）
3. **預設值**（代碼中的預設值）

## 範例配置

### 開發環境

```env
PORT=8080
GIN_MODE=debug
LOG_LEVEL=debug
DEBUG=true
CORS_ALLOW_ORIGINS=*
```

### 生產環境

```env
PORT=8080
GIN_MODE=release
LOG_LEVEL=info
LOG_OUTPUT=file
LOG_FILE=./logs/gateway.log
CORS_ALLOW_ORIGINS=https://yourdomain.com
MAX_CONNECTIONS=200
```

## 注意事項

1. **`.env` 文件不會被提交到版本控制**，請確保在部署時創建 `.env` 文件
2. **敏感信息**（如密碼、API 密鑰）應使用環境變數而非 `.env` 文件
3. **`.env.sample`** 文件包含所有可用的配置項，可以作為參考

## 故障排除

### 配置未生效

1. 確認 `.env` 文件位於項目根目錄
2. 確認配置項名稱正確（區分大小寫）
3. 重啟服務以載入新配置

### 端口被占用

修改 `PORT` 環境變數：
```env
PORT=3000
```

### CORS 錯誤

調整 `CORS_ALLOW_ORIGINS` 配置：
```env
CORS_ALLOW_ORIGINS=http://localhost:3000,http://localhost:8080
```
