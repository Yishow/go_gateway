# datalink-react-query Specification

## Purpose
TBD - created by archiving change enhance-ui-backend-sync. Update Purpose after archive.
## Requirements
### Requirement: Query-based data fetching

前端 SHALL 使用 TanStack Query 管理所有 Datalink 資源的讀取操作。

#### Scenario: Device list with caching

- WHEN 使用者進入 Devices 頁面
- THEN 系統從快取載入資料（若有效）或發起 API 請求
- AND 資料在 staleTime 內不會重複請求

#### Scenario: Automatic refetch on focus

- WHEN 使用者從其他分頁切回應用
- THEN 系統自動重新驗證資料新鮮度

### Requirement: Mutation-based data modification

前端 SHALL 使用 TanStack Query Mutations 管理所有寫入操作。

#### Scenario: Create device with cache invalidation

- WHEN 使用者建立新設備
- THEN 系統發起 POST 請求
- AND 成功後自動使設備列表快取失效

#### Scenario: Optimistic update for status toggle

- WHEN 使用者切換設備啟用狀態
- THEN UI 立即反映新狀態
- AND 若後端失敗則回退至原狀態

### Requirement: Centralized query configuration

前端 SHALL 提供全域 QueryClient 配置統一管理重試、快取策略。

#### Scenario: Network retry

- WHEN API 請求因網路問題失敗
- THEN 系統自動重試最多 2 次
- AND 重試間隔採用指數退避

