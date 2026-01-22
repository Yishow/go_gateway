# Design: Add Device Onboarding and Readiness Check

## Context

根據深度架構分析，當前系統在設備配置流程中存在以下問題：

1. **使用者引導不足**：新增設備後沒有明確的下一步指引
2. **狀態驗證缺失**：設備啟用時不檢查配置完整性
3. **即時反饋不足**：無法查看設備收集狀態

本設計文件說明如何解決這些問題，同時保持系統的簡潔性和可維護性。

## Goals / Non-Goals

### Goals

- ✅ 提供清晰的設備配置引導流程
- ✅ 實作配置完成度檢查機制
- ✅ 提供即時的設備狀態反饋
- ✅ 保持向後相容性
- ✅ 遵循 TDD 工作流程

### Non-Goals

- ❌ 不改變現有的設備、點位、標籤、映射資料模型（僅擴充）
- ❌ 不重構現有的 Scheduler 架構
- ❌ 不實作複雜的配置模板系統（留待後續改進）

## Decisions

### Decision 1: Readiness Check 實作方式

**What**: 實作 `DeviceService.CheckReadiness()` 方法，檢查設備配置完整性。

**Why**: 
- 集中化檢查邏輯，易於維護和測試
- 可重複使用於 API、UI、CLI 等多個場景
- 符合單一職責原則

**Alternatives considered**:
- **選項 A**: 在 API Handler 中直接實作檢查邏輯
  - Pros: 簡單直接
  - Cons: 無法重複使用，違反分層架構
- **選項 B**: 使用獨立的 Readiness Service
  - Pros: 職責分離
  - Cons: 過度設計，增加複雜度

**選擇**: 選項 B（在 Device Service 中實作）

### Decision 2: 資料庫欄位擴充

**What**: 在 `devices` 表中新增 `last_collected_at`, `collection_count`, `error_count`, `readiness_status` 欄位。

**Why**:
- 提供即時的狀態資訊
- 支援統計和監控需求
- 避免頻繁查詢關聯表

**Alternatives considered**:
- **選項 A**: 不新增欄位，每次都查詢關聯表
  - Pros: 資料一致性高
  - Cons: 效能差，查詢複雜
- **選項 B**: 使用快取機制
  - Pros: 效能好
  - Cons: 快取失效問題，增加複雜度

**選擇**: 選項 A（新增欄位，定期更新）

### Decision 3: Onboarding Wizard 實作方式

**What**: 建立獨立的 `DeviceOnboardingWizard` 組件，整合現有的表單組件。

**Why**:
- 提供統一的引導體驗
- 可重複使用現有組件
- 易於維護和擴展

**Alternatives considered**:
- **選項 A**: 修改現有的 `DeviceForm` 組件
  - Pros: 不需要新組件
  - Cons: 組件職責混亂，難以維護
- **選項 B**: 使用第三方 Wizard 庫
  - Pros: 功能完整
  - Cons: 增加依賴，可能不符合設計需求

**選擇**: 選項 A（建立新組件，整合現有組件）

### Decision 4: 即時更新機制

**What**: 使用 Server-Sent Events (SSE) 推送設備狀態更新。

**Why**:
- 簡單易實作
- 適合單向資料推送場景
- 不需要複雜的雙向通訊

**Alternatives considered**:
- **選項 A**: 使用 WebSocket
  - Pros: 雙向通訊，功能強大
  - Cons: 複雜度高，需要連接管理
- **選項 B**: 使用輪詢 (Polling)
  - Pros: 簡單直接
  - Cons: 效能差，即時性不足

**選擇**: 選項 A（使用 SSE，簡單且足夠）

## Risks / Trade-offs

### Risk 1: 資料庫遷移風險

**風險**: 新增欄位可能影響現有查詢和應用程式。

**緩解措施**:
- 使用 `ALTER TABLE ADD COLUMN` 語句，預設值為 NULL
- 確保向下相容，現有程式碼不受影響
- 在測試環境充分測試遷移腳本

### Risk 2: 效能影響

**風險**: Readiness Check 需要查詢多個表，可能影響效能。

**緩解措施**:
- 使用 JOIN 查詢減少資料庫往返
- 實作查詢快取（如需要）
- 監控查詢效能，必要時優化

### Risk 3: 前端複雜度增加

**風險**: Onboarding Wizard 增加前端複雜度。

**緩解措施**:
- 遵循組件化原則，保持組件職責單一
- 充分測試，確保可維護性
- 提供清晰的文檔和註解

## Migration Plan

### 階段一：後端基礎設施（1-2 週）

1. 實作資料模型擴充和遷移
2. 實作 Readiness Check Service
3. 實作 Readiness Check API
4. 實作收集統計更新機制

### 階段二：前端引導流程（1-2 週）

1. 實作 Onboarding Wizard 組件
2. 整合現有 API 和組件
3. 實作完成度顯示

### 階段三：狀態儀表板（1 週）

1. 實作 Status Dashboard 頁面
2. 實作 SSE 即時更新
3. 整合統計資訊顯示

### 階段四：測試與部署（1 週）

1. 執行完整測試套件
2. 效能測試和優化
3. 文件更新
4. 部署到生產環境

## Open Questions

1. **Q**: Readiness Check 的更新頻率？
   - **A**: 建議在設備狀態變更時觸發檢查，或提供手動觸發 API。

2. **Q**: 收集統計的更新頻率？
   - **A**: 建議在每次收集成功/失敗時更新，使用批次更新減少資料庫操作。

3. **Q**: Onboarding Wizard 是否支援「稍後完成」功能？
   - **A**: 建議實作，允許使用者稍後繼續配置。

4. **Q**: Status Dashboard 的資料保留期間？
   - **A**: 建議保留最近 30 天的統計資料，歷史資料可查詢時序資料庫。
