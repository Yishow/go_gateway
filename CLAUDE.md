# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 專案概述

**Go Gateway** 是一個工業級數據採集閘道系統,用於從 PLC (Programmable Logic Controller) 採集數據並映射至時序資料庫。

- **後端**: Go 1.25.5 + Gin + SQLite/PostgreSQL
- **前端**: React 19 + TypeScript + Vite + Tailwind CSS
- **架構**: 單一可執行檔案包含嵌入式 Web UI (Embedded Frontend)

## 常用命令

### 開發環境

```bash
# 前端開發 (熱重載)
cd frontend
pnpm install
pnpm dev              # 開發伺服器 (http://localhost:5173)

# 前端測試與 Lint
pnpm test             # 執行 Vitest 測試
pnpm lint             # ESLint 檢查
pnpm lint:fix         # 自動修復 Lint 錯誤
```

### 建置

```bash
# 使用 Makefile 建置 (推薦)
make build            # 建置前端 + 後端完整應用
make build-frontend   # 僅建置前端
make build-backend    # 建置後端 (需先建置前端)
make clean            # 清理所有建置產物

# 手動建置
cd frontend && pnpm install && pnpm build
go build -o bin/test-ui.exe ./cmd/test_ui
```

### 執行

```bash
# 執行主應用 (包含 Web UI)
./bin/test-ui.exe
# 或
make test-ui

# 執行協議測試工具
go run ./cmd/fatek_test   # FATEK 協議測試
go run ./cmd/test_all     # 完整協議測試
```

### 測試

```bash
# 執行所有 Go 測試
go test ./...
go test -v ./internal/...  # 詳細輸出

# 執行單一測試
go test -run TestFunctionName ./path/to/package

# 效能基準測試
go test -bench=. ./lib/hsllogic
go test -bench=BenchmarkPacketLogger ./lib/hsllogic -benchmem
```

### 資料庫

```bash
# 預設使用 SQLite (檔案位於專案根目錄)
# 設定環境變數切換至 PostgreSQL
export DB_TYPE=postgres
export DB_DSN="host=localhost port=5432 user=postgres password=postgres dbname=gateway"
```

## 專案架構

### 高階架構圖

```
┌─────────────────────────────────────────────────┐
│   Frontend (React + TypeScript)                 │
│   - 設備管理 (Devices)                           │
│   - 標籤管理 (Tags)                              │
│   - 映射規則 (Mappings)                          │
│   - 儀表板 (Dashboard)                           │
└───────────────┬─────────────────────────────────┘
                │ HTTP/SSE
┌───────────────▼─────────────────────────────────┐
│   API Layer (Gin HTTP Server)                   │
│   - /api/v1/datalink/*                          │
│   - /api/v1/protocol/test-connection            │
│   - /events (SSE 推送)                          │
└───────────────┬─────────────────────────────────┘
                │
┌───────────────▼─────────────────────────────────┐
│   Datalink Core (核心業務邏輯)                   │
│   ┌───────────────────────────────────────────┐ │
│   │ Scheduler (數據採集調度器)                 │ │
│   │  - 定時輪詢 PLC 數據                       │ │
│   │  - 多設備並行採集                          │ │
│   └────────┬──────────────────────────────────┘ │
│            │                                     │
│   ┌────────▼──────────────────────────────────┐ │
│   │ ConnectorManager (協議連接管理)           │ │
│   │  - 連接池管理                              │ │
│   │  - 自動重連                                │ │
│   └────────┬──────────────────────────────────┘ │
│            │                                     │
│   ┌────────▼──────────────────────────────────┐ │
│   │ Protocol Adapters (協議適配器)            │ │
│   │  - FATEK (永宏 FBs)                       │ │
│   │  - Modbus TCP/RTU                         │ │
│   │  - MC Protocol 3E (三菱)                  │ │
│   │  - MQTT                                   │ │
│   └────────┬──────────────────────────────────┘ │
│            │                                     │
│   ┌────────▼──────────────────────────────────┐ │
│   │ Mapping Pipeline (映射轉換管道)           │ │
│   │  - Tag → Point 映射                       │ │
│   │  - 數據轉換 (Transform)                   │ │
│   └────────┬──────────────────────────────────┘ │
└────────────┼─────────────────────────────────────┘
             │
┌────────────▼─────────────────────────────────┐
│   Storage Layer                              │
│   - SQLite (開發)                            │
│   - PostgreSQL (生產)                        │
│   - Time Series Storage                      │
└──────────────────────────────────────────────┘
```

### 目錄結構與職責

```
go_gateway/
├── cmd/                          # 應用程式入口點
│   ├── test_ui/                  # ★ 主應用 (嵌入 Web UI)
│   ├── fatek_test/               # FATEK 協議測試工具
│   └── test_all/                 # 完整協議測試工具
│
├── internal/                     # 內部實現 (不對外暴露)
│   ├── api/                      # ★ API 層 (HTTP handlers)
│   │   └── handlers/             # 各功能 HTTP 處理器 (19 個)
│   │
│   ├── datalink/                 # ★★★ 核心模組 (38 個檔案)
│   │   ├── api/                  # Datalink RESTful API
│   │   ├── device/               # 設備管理 (CRUD + 狀態監控)
│   │   ├── tag/                  # 標籤 (數據點) 管理
│   │   ├── mapping/              # 映射規則引擎
│   │   ├── point/                # 點位目錄
│   │   ├── pollinggroup/         # 輪詢組 (批次採集)
│   │   ├── collector/            # 數據採集調度器
│   │   ├── connector/            # 協議連接管理
│   │   │   ├── manager.go        # 連接池管理器
│   │   │   ├── registry.go       # 協議註冊表
│   │   │   └── adapters/         # 協議適配器實現
│   │   │       ├── fatek.go      # FATEK 適配器
│   │   │       ├── modbus.go     # Modbus 適配器
│   │   │       ├── mc3e.go       # MC Protocol 適配器
│   │   │       └── mqtt.go       # MQTT 適配器
│   │   ├── schema/               # 資料庫 Schema + Migrations
│   │   ├── storage/              # 時序資料儲存
│   │   ├── settings/             # 系統設定
│   │   └── common/               # 共用工具
│   │
│   ├── protocol/                 # ★ 協議實現層 (24 個檔案)
│   │   ├── modbus/               # Modbus TCP/RTU
│   │   │   ├── client.go         # 客戶端
│   │   │   ├── frame.go          # 幀結構解析
│   │   │   ├── transport.go      # 傳輸層 (TCP/Serial)
│   │   │   └── factory.go        # 工廠模式
│   │   ├── fatek/                # FATEK (永宏 FBs)
│   │   │   ├── client.go
│   │   │   ├── address.go        # 地址解析 (D, M, X, Y, T, C)
│   │   │   ├── frame.go          # ASCII 幀格式
│   │   │   ├── pool.go           # 連接池
│   │   │   └── transport.go      # TCP/Serial
│   │   └── mcprotocol/           # MC Protocol 3E (三菱)
│   │       ├── client.go
│   │       ├── frame.go          # Binary/ASCII 幀
│   │       └── transport.go
│   │
│   ├── config/                   # 配置管理 (.env + 環境變數)
│   └── web/                      # 嵌入式前端資源服務
│
├── frontend/                     # ★ React 前端應用
│   ├── src/
│   │   ├── components/
│   │   │   └── datalink/         # Datalink 業務組件
│   │   │       ├── DeviceCard.tsx
│   │   │       ├── DeviceForm.tsx
│   │   │       ├── TagTable.tsx
│   │   │       ├── MappingCanvas.tsx
│   │   │       └── wizard/       # 映射嚮導 (多步驟流程)
│   │   ├── pages/
│   │   │   └── datalink/         # Datalink 頁面
│   │   │       ├── Dashboard.tsx      # 儀表板
│   │   │       ├── DevicesPage.tsx    # 設備管理
│   │   │       ├── TagsPage.tsx       # 標籤管理
│   │   │       ├── MappingsPage.tsx   # 映射規則
│   │   │       └── MappingWizardPage.tsx
│   │   ├── hooks/
│   │   │   └── datalink/         # React Query Hooks
│   │   │       ├── useDevices.ts
│   │   │       ├── useTags.ts
│   │   │       └── useMappings.ts
│   │   ├── services/
│   │   │   ├── api.ts            # Axios HTTP 客戶端
│   │   │   └── datalink.ts       # Datalink API 封裝
│   │   ├── types/
│   │   │   ├── api.ts            # API 類型定義
│   │   │   ├── datalink.ts       # Datalink 領域模型
│   │   │   └── protocol.ts       # 協議類型
│   │   ├── i18n/                 # 國際化 (zh-TW + en)
│   │   │   └── locales/
│   │   └── utils/                # 前端工具函數
│   └── package.json              # pnpm 依賴管理
│
├── lib/                          # 第三方庫和效能優化
│   ├── hsllogic/                 # HSL Logic 效能優化庫
│   │   ├── bytetransform.go      # 字節轉換優化
│   │   ├── address_parser.go     # 地址解析優化
│   │   ├── packet_logger.go      # 封包日誌 (效能基準)
│   │   └── *_test.go             # 基準測試
│   ├── fatek_lib/                # FATEK 參考資料
│   └── mcprotocol_lib/           # MC Protocol 文檔
│
├── openspec/                     # ★ 規格驅動開發 (Spec-Driven)
│   ├── AGENTS.md                 # AI 助手工作流程指南
│   ├── project.md                # 專案規範與慣例
│   ├── specs/                    # 現行規格 (已實現功能)
│   │   ├── device-registry/
│   │   ├── protocol-connectors/
│   │   ├── mapping-pipeline/
│   │   ├── datalink-api/
│   │   ├── datalink-ui/
│   │   └── ...                   # (共 11 個規格)
│   └── changes/                  # 變更提案 (進行中)
│       ├── p1-implement-circuit-breaker/
│       └── archive/              # 已完成的變更
│
├── docs/                         # 技術文檔
├── scripts/                      # PowerShell 建置腳本
├── go.mod                        # Go 模組依賴
├── Makefile                      # 建置自動化
├── .env.sample                   # 環境變數範本
└── CLAUDE.md                     # (本檔案)
```

### 核心設計模式

| 設計模式 | 位置 | 用途 |
|---------|------|------|
| **分層架構** | `datalink/` | API → Service → Repository |
| **適配器模式** | `connector/adapters/` | 統一協議介面 (`protocol.Interface`) |
| **工廠模式** | `protocol/*/factory.go` | 動態創建協議客戶端 |
| **倉庫模式** | `*_repo.go` | 數據持久化抽象 (SQL/Memory) |
| **觀察者模式** | `datalink_sse_handler.go` | 即時狀態推送 (Server-Sent Events) |
| **連接池** | `connector/manager.go` | 複用 PLC 連接 |
| **嚮導模式** | 前端 `wizard/` | 多步驟流程 (Mapping Wizard) |

### 資料流

完整的 ETL (Extract-Transform-Load) 管道:

```
PLC 設備 (FATEK/Modbus/MC3E)
  ↓ [Protocol Adapter]
標籤 (Tag) 採集
  ↓ [Collector/Scheduler]
映射規則 (Mapping)
  ↓ [Transform Pipeline]
點位 (Point) 儲存
  ↓ [Storage Layer]
時序資料庫 (SQLite/PostgreSQL)
  ↓ [API Layer]
前端儀表板
```

## 開發指南

### 新增協議適配器

1. **定義協議客戶端**: `internal/protocol/yourprotocol/client.go`
2. **實現 `protocol.Interface`**: 必須實現以下方法:
   ```go
   Connect(address string) error
   Disconnect() error
   ReadCoils(address uint16, quantity uint16) ([]bool, error)
   ReadRegisters(address uint16, quantity uint16) ([]uint16, error)
   WriteCoil(address uint16, value bool) error
   WriteRegister(address uint16, value uint16) error
   ```
3. **創建適配器**: `internal/datalink/connector/adapters/yourprotocol.go`
4. **註冊協議**: 在 `registry.go` 中註冊新協議

### 新增 API 端點

1. **定義 Handler**: `internal/api/handlers/your_handler.go`
2. **掛載路由**: 在 `router.go` 中註冊路由
3. **實現服務層**: `internal/datalink/yourmodule/service.go`
4. **定義倉庫層**: `internal/datalink/yourmodule/sql_repo.go`

### 前端開發

**狀態管理**: 使用 React Query (`@tanstack/react-query`)

```typescript
// 範例: 使用 useDevices Hook
import { useDevices } from '@/hooks/datalink';

function MyComponent() {
  const { data: devices, isLoading, error } = useDevices();
  // ...
}
```

**API 呼叫**: 透過 `services/datalink.ts`

```typescript
import { datalinkService } from '@/services/datalink';

const device = await datalinkService.getDevice(deviceId);
```

**i18n 使用**:

```typescript
import { useTranslation } from 'react-i18next';

function MyComponent() {
  const { t } = useTranslation();
  return <h1>{t('datalink.devices.title')}</h1>;
}
```

### 資料庫遷移

遷移檔案位於: `internal/datalink/schema/migrations/`

```go
// 自動執行遷移 (應用啟動時)
migrator := datalink.NewMigrator(db)
migrator.Migrate()
```

### 測試策略

- **單元測試**: 每個模組有對應的 `*_test.go`
- **整合測試**: `cmd/test_all/` 測試完整協議棧
- **效能基準**: `lib/hsllogic/*_bench_test.go`

```bash
# 執行特定模組測試
go test ./internal/datalink/mapping -v

# 效能基準測試
go test -bench=BenchmarkByteTransform ./lib/hsllogic -benchmem
```

## OpenSpec 規格驅動開發

本專案使用 **OpenSpec** 進行規格驅動開發 (Spec-Driven Development)。

### 工作流程

1. **創建提案**: 在 `openspec/changes/[change-id]/` 創建提案
   - `proposal.md` - 變更描述
   - `tasks.md` - 實作檢查清單
   - `design.md` - 技術決策 (選用)
   - `specs/[capability]/spec.md` - 規格變更 (Delta)

2. **實作功能**: 按照 `tasks.md` 依序實作

3. **歸檔變更**: 完成後移至 `archive/` 並更新 `specs/`

### 常用 OpenSpec 命令

```bash
# 列出現行規格
openspec list --specs

# 列出進行中變更
openspec list

# 檢視規格詳情
openspec show device-registry --type spec

# 驗證變更
openspec validate [change-id] --strict --no-interactive

# 歸檔已完成變更
openspec archive [change-id] --yes
```

### 重要原則

- **永遠先閱讀相關規格** (`openspec/specs/`) 再開始實作
- **變更需經過驗證** (`openspec validate`) 才能開始實作
- **使用正確的 Scenario 格式**:
  ```markdown
  #### Scenario: 場景名稱
  - **WHEN** 條件
  - **THEN** 預期結果
  ```
- **新功能必須創建提案** (Breaking changes, 架構變更)
- **Bug 修復可以直接實作** (如符合現有規格)

## 程式碼風格與慣例

### Go 後端

- **模組組織**: 使用 `internal/` 防止外部依賴
- **錯誤處理**: 總是檢查錯誤並返回有意義的錯誤訊息
- **命名**: 使用 Go 標準命名慣例 (camelCase for unexported, PascalCase for exported)
- **註解**: 公開函數必須有 GoDoc 註解

### TypeScript 前端

- **型別定義**: 集中在 `types/` 目錄
- **組件**: 使用函數式組件 + Hooks
- **狀態管理**: React Query 管理伺服器狀態
- **樣式**: Tailwind CSS (避免自定義 CSS)
- **i18n**: 所有使用者可見文字必須透過 `t()` 函數

## 環境變數

複製 `.env.sample` 至 `.env` 並調整:

```bash
# 資料庫配置
DB_TYPE=sqlite          # sqlite | postgres | mysql
DB_DSN=gateway.db       # 資料庫連線字串

# 伺服器配置
SERVER_PORT=8080        # HTTP 伺服器埠
FRONTEND_PATH=./frontend/dist  # 前端打包目錄

# 日誌配置
LOG_LEVEL=info          # debug | info | warn | error
```

## 部署

### 單一可執行檔案部署

```bash
# 1. 建置完整應用
make build

# 2. 部署單一檔案
./bin/test-ui.exe
```

內嵌的前端資源會自動從 `frontend/dist/` 打包進執行檔。

### Docker 部署 (規劃中)

```bash
# 尚未實作,請參考 changes/ 中的提案
```

## 效能最佳化

### HSL Logic 庫

`lib/hsllogic/` 提供高效能字節操作:

```go
import "go-gateway/lib/hsllogic"

// 使用優化的地址解析
parser := hsllogic.NewAddressParser()
addr, err := parser.Parse("D100")

// 使用優化的字節轉換
transformer := hsllogic.NewByteTransform()
value := transformer.TransByte(data, 0)
```

**基準測試結果**: 約 2-3x 效能提升 (見 `*_bench_test.go`)

## 常見問題

### 前端無法連接後端

確認:
1. 後端已啟動 (`./bin/test-ui.exe`)
2. 前端 `vite.config.ts` 中的 proxy 設定正確
3. 防火牆允許 8080 埠

### 協議連線失敗

檢查:
1. PLC 設備 IP 和埠正確
2. 網路連通性 (`ping` 測試)
3. 協議參數 (站號, 波特率等)
4. 查看日誌 (`LOG_LEVEL=debug`)

### 資料庫遷移失敗

```bash
# 刪除資料庫檔案重新初始化
rm gateway.db
./bin/test-ui.exe
```

## 相關資源

- **協議文檔**: `lib/fatek_lib/`, `lib/mcprotocol_lib/`
- **前端 i18n**: `frontend/src/i18n/README.md`
- **OpenSpec 指南**: `openspec/AGENTS.md`
- **技術規格**: `openspec/specs/`

## 關鍵檔案快速參考

| 功能 | 檔案路徑 |
|------|---------|
| 主程式入口 | `cmd/test_ui/main.go` |
| API 路由 | `internal/api/router.go` |
| 協議註冊表 | `internal/datalink/connector/registry.go` |
| 連接管理 | `internal/datalink/connector/manager.go` |
| 資料庫 Schema | `internal/datalink/schema/models.go` |
| 前端路由 | `frontend/src/App.tsx` |
| API 封裝 | `frontend/src/services/datalink.ts` |
| 型別定義 | `frontend/src/types/datalink.ts` |

---

**最後更新**: 2026-01-22
**維護者**: AI Development Team
