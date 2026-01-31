# Tasks: Add Operational Tools

- [x] API Documentation <!-- id: 0 -->
  - [x] Install `swag` tool and dependencies <!-- id: 1 -->
  - [x] Add `// @title`, `// @version` comments to `main.go` and `internal/api/router.go` <!-- id: 2 -->
  - [x] Add annotation comments to Health Check handler as a proof-of-concept <!-- id: 3 -->
  - [x] Update `Makefile` to include `gen-docs` target running `swag init` <!-- id: 4 -->
  - [x] Mount `gin-swagger` middleware in `internal/api/router.go` <!-- id: 5 -->
  - [ ] Verify `/swagger/index.html` loads locally <!-- id: 6 --> (需要先產生 docs)

- [x] Configuration Validation <!-- id: 7 -->
  - [x] Refactor `internal/config` to expose a `Validate()` function separate from `Load()` <!-- id: 8 -->
  - [x] Implement `validate` logic for Server (check port) <!-- id: 9 -->
  - [x] Implement `validate` logic for API (check base path) <!-- id: 10 -->
  - [x] Implement `validate` logic for CORS (check allow origins) <!-- id: 11 -->
  - [x] Implement `validate` logic for ConnectionPool (check max connections) <!-- id: 12 -->
  - [ ] Add `validate` subcommand to `cmd/gateway/main.go` using Cobra <!-- id: 13 --> (待實作 CLI)
  - [x] 7 個配置驗證測試全部通過 <!-- id: 14 -->

## 實作摘要

| 元件 | 檔案 | 說明 |
|------|------|------|
| Swagger 註解 | `cmd/test_ui/main.go` | API 文檔標題和版本 |
| Swagger 中間件 | `internal/api/router.go` | gin-swagger 路由 |
| 健康檢查註解 | `internal/api/handlers/datalink_health_handler.go` | Swagger 註解 |
| 配置驗證 | `internal/config/validate.go` | 驗證器實作 |
| Makefile | `Makefile` | gen-docs target |

## 測試覆蓋

- 配置驗證: 7 個測試全部通過
