## Context

目前 datalink 首頁由 `DatalinkLayout`（外層 sidebar + top）包住 `SmartDashboard`。實際操作上，設備選擇在進入流程後不會頻繁操作，常駐設備欄與外層導覽會壓縮主工作區，並增加跨層狀態同步成本。  
本變更將首頁收斂為單一容器：`SmartDashboard` 直接承載導覽、設備 context、工作區與右側操作，設備清單改為按需展開。

## Goals / Non-Goals

**Goals:**
- 建立 dashboard-first 入口，移除 `/datalink` 外層包裝。
- 將分頁導覽與設備狀態列整合為上方水平控制區。
- 將設備樹改為「切換設備時才展開」的 drawer，平時不佔位。
- 規範切換設備異常流程（離線/失敗/未儲存/權限）。
- 保持 legacy 路由可導向首頁內對應區段，避免斷鏈。
- 強化本機 `5020` memory grid 寫入治理，提供獨立且完整的操作頁。

**Non-Goals:**
- 不重寫 Source Planner / Memory Grid / Commit Queue 既有業務邏輯。
- 不變更後端 API contract。
- 不在本變更引入新設計系統或 UI 套件。

## Decisions

### 1) 路由容器收斂到 SmartDashboard
- Decision: 移除 `/datalink` 的 `DatalinkLayout` 外層，改由 `SmartDashboard` 成為首頁容器。
- Rationale: 消除雙層 layout，降低狀態來源分散與樣式衝突。
- Alternative considered:
  - 保留外層 layout，僅視覺隱藏 sidebar：仍保留冗餘容器與耦合，不採用。

### 2) 上方控制區改為水平雙區
- Decision: 將「內建分頁」與「已選設備條」合併為同列水平排列；窄屏才堆疊。
- Rationale: 釋放垂直空間給主要工作區，提高流程可視範圍。
- Alternative considered:
  - 維持上下堆疊：可讀性較高但浪費高度，不採用作為預設。

### 3) 設備樹改為 drawer on demand
- Decision: 平時不顯示常駐設備欄，點擊「切換設備」才展開設備 drawer。
- Rationale: 符合「大多數時間不切設備」的操作模型，降低誤觸與視覺雜訊。
- Alternative considered:
  - 保留可收合左欄：仍占用版面與心智負擔，不採用。

### 4) 切換設備錯誤處理統一
- Decision: 以 context bar + toast/dialog/banner 呈現錯誤，不在正式版面顯示流程圖卡片。
- Rationale: 正式 UI 應聚焦操作，不暴露設計流程圖；錯誤仍需可見且可恢復。

### 5) 舊路徑相容策略
- Decision: `/datalink/devices`、`/datalink/settings` 轉向 `/datalink?section=devices|settings`。
- Rationale: 保留舊書籤可用性，並統一回 dashboard-first 導覽。

### 6) 設備切換微互動與保護
- Decision: 實作狀態色階、未儲存紅點、切換中鎖定按鈕、失敗 toast(重試+詳情)、確認彈窗預設焦點 `cancel`。
- Rationale: 降低誤觸與錯切，讓異常恢復可預期。

### 7) 設備操作改為 Modal-first
- Decision: 除 `/test` 外，設備建立與編輯採 dashboard 內 modal；建立完成可「立即切換」。
- Rationale: 保持單頁操作連續性，避免跳頁破壞上下文。

### 8) 本機 5020 獨立工作頁
- Decision: 新增 `local-modbus` 獨立頁，承載完整本機 `5020` memory grid 寫入功能（衝突掃描、映射編輯、批次操作、驗證與測試）。
- Rationale: 設備間資料域本質隔離，真正需要全域衝突治理的是本機 `5020` 寫入空間。

## Risks / Trade-offs

- [Risk] 路由重構可能造成既有導覽測試失效  
  → Mitigation: 補齊路由重定向與 section 解析測試。

- [Risk] 設備 drawer 切換時，未儲存狀態可能被覆蓋  
  → Mitigation: 加入切換前確認彈窗與明確三選項（儲存後切換/放棄/取消）。

- [Risk] 常駐設備欄移除後，部分使用者需重新適應  
  → Mitigation: 在 context bar 保留當前設備資訊與單一切換入口。

- [Risk] 內容密度上升導致小螢幕可讀性下降  
  → Mitigation: 設計斷點規則，平板以下改為上下堆疊。

- [Risk] `5020` 工作頁與 Dashboard 功能重疊造成認知負擔  
  → Mitigation: 在 Dashboard 僅保留摘要與入口，完整操作統一導向獨立頁。

## Migration Plan

1. 先完成 OpenSpec delta 與任務清單，確定驗收條件。  
2. 修改路由容器，讓 `/datalink` 直接使用 `SmartDashboard`。  
3. 在 `SmartDashboard` 建立 `Device Context Bar` 與 `Switch Device Drawer`。  
4. 將既有設備相關互動遷移到 drawer 開關流程。  
5. 補齊異常流程 UI（離線/失敗/未儲存/權限）與測試。  
6. 上線前以 feature branch 驗證 legacy URL 導向。  
7. 新增本機 `5020` 獨立頁並完成完整功能驗證。  

Rollback:
- 保留 `DatalinkLayout` 檔案與舊路由配置，若出現重大回歸可回切原容器模式。

## Resolved Decisions

- Context bar MUST 顯示「最近切換時間」，用於操作追蹤與班次交接判讀。
- 除測試工具頁（`/test`）外，原本獨立設備建立流程（含 `/datalink/devices/new`）改為 dashboard 內 modal 流程。
- 手機端設備切換容器定案為側滑 panel（drawer），不採 full-screen sheet。
- 設備切換互動採保護模式：未儲存紅點提示、切換中鎖定、失敗提供重試與詳情。
- 本機 `5020` memory grid 寫入衝突治理與進階操作放入獨立頁，且需具完整功能。
