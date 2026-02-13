# Tasks: Add Operational Tools

- [x] API Documentation <!-- id: 0 -->
  - [x] Install `swag` tool and dependencies <!-- id: 1 -->
  - [x] Add `// @title`, `// @version` comments to `main.go` and `internal/api/router.go` <!-- id: 2 -->
  - [x] Add annotation comments to Health Check handler as a proof-of-concept <!-- id: 3 -->
  - [x] Update `Makefile` to include `gen-docs` target running `swag init` <!-- id: 4 -->
  - [x] Mount `gin-swagger` middleware in `internal/api/router.go` <!-- id: 5 -->
  - [x] Verify `/swagger/index.html` loads locally <!-- id: 6 --> ✅ Swagger docs 已產生

- [x] Configuration Validation <!-- id: 7 -->
  - [x] Refactor `internal/config` to expose a `Validate()` function separate from `Load()` <!-- id: 8 -->
  - [x] Implement `validate` logic for Server (check port) <!-- id: 9 -->
  - [x] Implement `validate` logic for API (check base path) <!-- id: 10 -->
  - [x] Implement `validate` logic for CORS (check allow origins) <!-- id: 11 -->
  - [x] Implement `validate` logic for ConnectionPool (check max connections) <!-- id: 12 -->
  - [x] Add `validate` subcommand to `cmd/gateway/main.go` using Cobra <!-- id: 13 --> ✅ 已實作
  - [x] 22 個配置驗證測試全部通過 <!-- id: 14 -->
  - [x] 7 個 CLI 測試全部通過 <!-- id: 15 -->

## 實作摘要

| 元件 | 檔案 | 說明 |
|------|------|------|
| Swagger 註解 | `cmd/test_ui/main.go` | API 文檔標題和版本 |
| Swagger 中間件 | `internal/api/router.go` | gin-swagger 路由 |
| Swagger Docs | `docs/swagger/` | 自動產生的 API 文檔 |
| 健康檢查註解 | `internal/api/handlers/datalink_health_handler.go` | Swagger 註解 |
| 配置驗證 | `internal/config/validate.go` | 驗證器實作 |
| 配置載入 | `internal/config/config.go` | LoadFromFile 函數 |
| CLI 入口 | `cmd/gateway/main.go` | Cobra CLI 框架 |
| Validate 命令 | `cmd/gateway/validate.go` | 驗證子命令 |
| Makefile | `Makefile` | gen-docs target |

## 測試覆蓋

- 配置驗證: 22 個測試全部通過
- CLI 命令: 7 個測試全部通過

## 新增依賴

- github.com/spf13/cobra v1.10.2
- github.com/swaggo/swag v1.16.6
