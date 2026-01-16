# ============================================
# 創建 .env 文件腳本
# ============================================
# 此腳本會創建一個包含所有可配置選項的 .env 文件
# ============================================

$envContent = @"
# ============================================
# Go Gateway 環境變數配置
# ============================================
# 說明：此文件包含所有可配置的環境變數
# 優先順序：環境變數 > .env 文件 > 預設值
# ============================================

# ============================================
# 後端服務配置
# ============================================
# 服務監聽端口
PORT=8080

# 服務監聽地址（空表示所有接口，127.0.0.1 表示僅本地）
HOST=

# Gin 運行模式：debug, release, test
GIN_MODE=release

# ============================================
# API 配置
# ============================================
# API 基礎路徑
API_BASE_PATH=/api/v1

# ============================================
# CORS 配置
# ============================================
# 允許的來源（逗號分隔，* 表示允許所有來源）
# 範例：CORS_ALLOW_ORIGINS=http://localhost:3000,http://localhost:8080
CORS_ALLOW_ORIGINS=*

# ============================================
# 日誌配置
# ============================================
# 日誌級別：debug, info, warn, error
LOG_LEVEL=info

# 日誌輸出位置：console, file
LOG_OUTPUT=console

# 日誌文件路徑（當 LOG_OUTPUT=file 時使用）
LOG_FILE=./logs/gateway.log

# ============================================
# WebSocket 配置
# ============================================
# WebSocket 讀取超時（秒）
WS_READ_TIMEOUT=60

# WebSocket 寫入超時（秒）
WS_WRITE_TIMEOUT=10

# WebSocket Ping 間隔（秒）
WS_PING_INTERVAL=30

# ============================================
# 連接池配置
# ============================================
# 最大連接數
MAX_CONNECTIONS=100

# 連接超時（秒）
CONNECTION_TIMEOUT=5

# 讀取超時（秒）
READ_TIMEOUT=2

# 寫入超時（秒）
WRITE_TIMEOUT=2

# ============================================
# 開發環境配置
# ============================================
# 是否啟用調試模式（影響 Gin 模式和日誌詳細程度）
DEBUG=false

# ============================================
# 前端配置（Vite 環境變數，需要 VITE_ 前綴）
# ============================================
# API 基礎路徑（前端使用，通常與後端 API_BASE_PATH 一致）
VITE_API_BASE_URL=/api/v1

# Vite 開發伺服器端口
VITE_DEV_PORT=5173

# 後端服務地址（用於 Vite 代理配置）
# 此配置用於將前端的 /api 請求代理到後端服務
VITE_API_PROXY_TARGET=http://localhost:8080
"@

# 檢查 .env 文件是否已存在
if (Test-Path ".env") {
    Write-Host "⚠️  .env 文件已存在" -ForegroundColor Yellow
    $response = Read-Host "是否要覆蓋現有的 .env 文件？(Y/N)"
    if ($response -ne "Y" -and $response -ne "y") {
        Write-Host "已取消操作" -ForegroundColor Yellow
        exit 0
    }
}

# 寫入 .env 文件
try {
    $envContent | Out-File -FilePath ".env" -Encoding UTF8 -NoNewline
    Write-Host "✅ .env 文件創建成功！" -ForegroundColor Green
    Write-Host ""
    Write-Host "💡 提示：" -ForegroundColor Cyan
    Write-Host "   1. 請根據實際需求修改 .env 文件中的配置值"
    Write-Host "   2. 修改後需要重啟服務才能生效"
    Write-Host "   3. 詳細配置說明請參考 docs/ENV_CONFIG.md"
} catch {
    Write-Host "❌ 創建 .env 文件失敗: $_" -ForegroundColor Red
    exit 1
}
