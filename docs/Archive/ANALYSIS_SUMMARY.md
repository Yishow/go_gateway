# Go Gateway 專案深度分析總結

## 📋 執行摘要

本文件總結了對 `go_gateway` 專案的深度架構分析，包含系統架構圖、新增設備完整流程、流程缺失分析及改進建議。

**分析日期：** 2026-01-21  
**專案版本：** 當前開發版本  
**分析範圍：** 設備管理、資料收集管線、時序儲存

---

## 🎯 核心發現

### ✅ 專案優勢

1. **架構清晰**
   - 採用 Hexagonal Architecture（六邊形架構）
   - 清晰的層級分離：API → Service → Repository → Database
   - 良好的模組化設計

2. **功能完整**
   - 完整的設備生命週期管理（建立、測試、啟用、停用）
   - 完整的資料收集管線（點位 → 映射 → 標籤 → 儲存）
   - 支援多種工業協議（Modbus、FATEK、MC Protocol、MQTT）

3. **技術選型合理**
   - Go 語言適合高並發場景
   - React + TypeScript 提供現代化前端體驗
   - SQLite/PostgreSQL 雙資料庫支援

### ⚠️ 主要缺失

1. **缺少使用者引導**
   - 新增設備後沒有明確的下一步指引
   - 使用者不知道需要完成哪些步驟才能開始收集資料

2. **缺少狀態驗證**
   - 設備啟用時只檢查連線，不檢查配置完整性
   - 沒有「就緒檢查」(Readiness Check) 機制

3. **缺少即時反饋**
   - 無法即時查看設備是否正在收集資料
   - 缺少資料收集狀態監控

4. **缺少批次操作**
   - 無法批次建立點位、標籤、映射
   - 配置大量設備時效率低下

---

## 📊 系統架構概覽

### 技術棧

```
前端層: React + TypeScript + TanStack Query
API層:   Gin Framework (Go)
業務層:  Go Services (Device/Point/Tag/Mapping)
連接層:  Protocol Adapters (Modbus/FATEK/MC3E/MQTT)
資料層:  SQLite (開發) / PostgreSQL (生產)
儲存層:  TimeSeries Storage (批次寫入)
```

### 核心模組

```
internal/datalink/
├── device/          # 設備管理
├── point/           # 點位管理
├── tag/             # 標籤管理
├── mapping/         # 映射管理
├── collector/       # 資料收集排程器
├── connector/       # 連線管理與協議適配
├── storage/         # 時序資料儲存
└── api/             # REST API 路由
```

---

## 🔄 新增設備完整流程

### 標準流程（9 步驟）

```
1. 建立設備 (draft)
   ↓
2. 測試連線
   ↓
3. 啟用設備 (active)
   ↓
4. 建立點位
   ↓
5. 指派輪詢群組
   ↓
6. 建立標籤
   ↓
7. 建立映射
   ↓
8. 啟動排程器
   ↓
9. 開始資料收集 ✅
```

### 關鍵檢查點

| 步驟 | 檢查項目 | 失敗後果 |
|------|---------|---------|
| 2 | 連線測試 | 無法啟用設備 |
| 3 | 設備啟用 | 無法參與收集 |
| 4 | 建立點位 | 無資料可收集 |
| 5 | 指派群組 | 點位不會被排程 |
| 6 | 建立標籤 | 無法建立映射 |
| 7 | 建立映射 | 資料無法儲存 |
| 8 | 排程器運行 | 不會觸發收集 |

---

## 🚨 流程缺失詳細分析

### 缺失 1: 缺少明確引導流程

**問題：**
- 新增設備後，系統沒有明確指引下一步該做什麼
- 使用者不知道需要完成哪些步驟才能開始收集資料

**影響：**
- 使用者可能啟用設備後就停止，不知道還需要建立點位、標籤、映射
- 導致設備處於 active 狀態但實際上沒有收集任何資料

**證據：**
```go
// internal/datalink/device/service.go:214
func (s *Service) Activate(ctx context.Context, id string) error {
    // 只檢查連線，不檢查點位、群組、映射
    if err := s.TestConnection(ctx, id); err != nil {
        return fmt.Errorf("連線測試失敗，無法啟用設備: %w", err)
    }
    // ...
}
```

**建議：**
- 實作 `DeviceOnboardingWizard` 引導流程
- 新增設備後自動開啟引導
- 顯示完成度進度條

### 缺失 2: 缺少狀態驗證

**問題：**
- 設備啟用時只檢查連線，不檢查是否有點位、群組、映射
- 沒有「就緒檢查」(Readiness Check) 機制

**影響：**
- 設備可能處於 active 狀態但無法實際收集資料
- 使用者不知道缺少哪些配置

**建議：**
- 實作 `DeviceService.CheckReadiness()` 方法
- 新增 `GET /devices/:id/readiness` API 端點
- 在 UI 中顯示配置完成度

### 缺失 3: 缺少即時狀態反饋

**問題：**
- 設備啟用後，沒有即時顯示「是否正在收集資料」
- 沒有「最後收集時間」或「收集統計」資訊

**影響：**
- 使用者無法確認設備是否正常運作
- 無法快速發現配置問題

**證據：**
- `Device` 模型沒有 `last_collected_at` 欄位
- `DevicesPage` 沒有顯示收集狀態

**建議：**
- 新增 `DeviceStatusDashboard` 頁面
- 擴充 `Device` 模型，新增收集統計欄位
- 實作 WebSocket/SSE 即時推送

### 缺失 4: 缺少批次操作

**問題：**
- 無法批次建立多個點位
- 無法批次建立映射
- 缺少「從模板匯入」功能

**影響：**
- 大量設備配置時效率低下
- 重複性工作增加

**建議：**
- 實作批次 API 端點
- 建立配置模板系統
- 支援 JSON/YAML 批次匯入

---

## 💡 改進建議優先級

### P0（立即實施）

1. **設備配置完成度檢查**
   - 實作 `CheckReadiness()` 方法
   - 新增 API 端點與 UI 顯示

2. **新增設備後引導流程**
   - 建立 `DeviceOnboardingWizard` 組件
   - 自動引導完成所有必要配置

3. **設備狀態儀表板**
   - 顯示設備收集狀態
   - 顯示配置完成度統計

### P1（短期實施）

4. **配置模板與批次匯入**
   - 建立配置模板系統
   - 支援批次匯入

5. **即時資料流監控**
   - WebSocket/SSE 推送
   - 即時資料流預覽

6. **配置驗證與影響分析**
   - 配置變更前驗證
   - 顯示影響範圍

### P2（中期實施）

7. **自動化配置建議**
   - 根據設備類型自動建議配置
   - 提供「一鍵完成配置」功能

8. **錯誤恢復與通知機制**
   - 更完善的錯誤恢復策略
   - 錯誤通知系統

9. **配置版本管理**
   - 配置版本控制
   - 支援配置回滾

---

## 📈 預期效果

### 使用者體驗提升
- ✅ 明確的引導流程，減少配置錯誤
- ✅ 即時狀態反饋，快速發現問題
- ✅ 批次操作支援，提升配置效率

### 系統可靠性提升
- ✅ 配置驗證與狀態檢查，減少運行時錯誤
- ✅ 即時監控與告警，及時發現異常

### 運維效率提升
- ✅ 自動化配置建議，減少人工干預
- ✅ 配置模板與批次匯入，加速部署

---

## 📚 相關文檔

1. **[架構分析文檔](./ARCHITECTURE_ANALYSIS.md)**
   - 詳細的系統架構圖
   - 資料流向圖
   - 完整的缺失分析與改進建議

2. **[新增設備流程指南](./DEVICE_ONBOARDING_FLOW.md)**
   - 9 步驟詳細流程說明
   - 常見問題檢查清單
   - API 端點快速參考

3. **[OpenSpec 規格文檔](../openspec/specs/)**
   - `device-registry/spec.md` - 設備註冊規格
   - `point-catalog/spec.md` - 點位目錄規格
   - `tag-dictionary/spec.md` - 標籤字典規格
   - `mapping-pipeline/spec.md` - 映射管線規格
   - `collection-scheduler/spec.md` - 收集排程器規格
   - `timeseries-storage/spec.md` - 時序儲存規格

---

## 🎯 下一步行動

### 立即行動（本週）

1. ✅ 審閱本分析報告
2. ✅ 確認優先級 P0 項目
3. ✅ 規劃實施時間表

### 短期行動（2 週內）

1. 實作設備配置完成度檢查
2. 建立新增設備引導流程
3. 新增設備狀態儀表板

### 中期行動（1 個月內）

1. 實作配置模板系統
2. 新增即時資料流監控
3. 實作配置驗證機制

---

## 📝 附錄

### 關鍵程式碼位置

| 功能 | 檔案路徑 |
|------|---------|
| 設備服務 | `internal/datalink/device/service.go` |
| 設備 API | `internal/datalink/api/device_handler.go` |
| 收集排程器 | `internal/datalink/collector/scheduler.go` |
| 連線管理 | `internal/datalink/connector/manager.go` |
| 時序儲存 | `internal/datalink/storage/timeseries.go` |
| 前端設備頁面 | `frontend/src/pages/datalink/DevicesPage.tsx` |
| 設備表單 | `frontend/src/components/datalink/DeviceForm.tsx` |

### 資料庫結構

| 表名 | 說明 |
|------|------|
| `devices` | 設備註冊表 |
| `points` | 點位定義表 |
| `tags` | 標籤字典表 |
| `mappings` | 映射關係表 |
| `polling_groups` | 輪詢群組表 |
| `timeseries` | 時序資料表 |

### API 端點總覽

| 端點 | 方法 | 說明 |
|------|------|------|
| `/devices` | GET/POST | 列出/建立設備 |
| `/devices/:id` | GET/PUT/DELETE | 取得/更新/刪除設備 |
| `/devices/:id/test` | POST | 測試設備連線 |
| `/devices/:id/activate` | POST | 啟用設備 |
| `/devices/:id/disable` | POST | 停用設備 |
| `/points` | GET/POST | 列出/建立點位 |
| `/tags` | GET/POST | 列出/建立標籤 |
| `/mappings` | GET/POST | 列出/建立映射 |

---

**報告完成時間：** 2026-01-21  
**分析工具：** 程式碼審查、架構分析、流程追蹤  
**建議實施時間：** 依優先級分階段實施
