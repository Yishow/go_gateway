# go_gateway UI 流程重構設計與規劃文件

## 1. 背景與目標

### 1.1 背景
隨著 `go_gateway` 系統功能的持續演進與擴充，目前的 UI 介面已逐漸無法滿足不同類型使用者的需求。新使用者在首次配置時常因過多的進階設定而感到不知所措，而進階使用者在進行精細調校時，又覺得現有的操作流程不夠集中與高效。為解決此問題，我們決定對 UI 流程進行重構，導入雙入口設計，以平衡「易用性」與「專業性」。

### 1.2 目標
*   **降低學習門檻**：為新手或日常任務提供快速設定入口（Quick Setup），簡化核心設定流程。
*   **提升專業效率**：為進階使用者提供專家工作台（Expert Workbench），集中管理複雜與進階的閘道器配置。
*   **無縫相容**：絕對不破壞任何既有功能與底層 API 契約，確保前後端邏輯完美接軌。
*   **平滑過渡**：透過 Feature Flag (灰度上線) 策略，讓新舊版本能安全、可控地進行切換與驗證。

### 1.3 核心架構決策
*   **決策 1A**：以「快速入口（Quick Setup）」作為系統的預設首頁，確保絕大多數的常規操作能在最短時間內完成。
*   **決策 2B**：採用 Feature Flag 進行灰度上線，允許在特定使用者群體或環境中進行測試，降低上線風險。
*   **決策 3C**：「快速入口（Quick Setup）」與「進階入口（Expert Workbench）」雙線同時推進開發，確保改版體驗的一致性與完整性。

---

## 2. 兩入口資訊架構

為了滿足不同使用情境，系統將分為兩大核心區塊，並透過明確的導覽列或頁籤進行切換。

### 2.1 快速入口 (Quick Setup) - 預設視圖
*   **定位**：任務導向，專注於 80% 的常見使用場景。
*   **核心模組**：
    *   **一鍵部署/配置精靈**：以步驟式 (Stepper) 方式引導建立基本 Gateway 路由。
    *   **核心狀態看板**：只顯示最重要的健康度、流量摘要與近期告警。
    *   **常用快捷操作**：快速啟動/暫停特定路由、更新憑證等。
*   **隱藏內容**：進階限流策略、自訂中介軟體 (Middleware) 腳本、微調負載均衡演算法等。

### 2.2 進階入口 (Expert Workbench)
*   **定位**：資料與配置導向，提供全域視角與 100% 的配置權限。
*   **核心模組**：
    *   **全域資源拓撲圖/列表**：以樹狀或關聯圖展示所有節點、服務與路由。
    *   **JSON/YAML 原始碼編輯器**：支援直接編輯配置檔，並附帶語法檢查與 Auto-complete。
    *   **進階策略控制面板**：WAF 規則、細粒度限流 (Rate Limiting)、CORS 深度設定、日誌格式自訂。
    *   **即時監控與除錯日誌**：串接完整的存取日誌與詳細 Metrics。

---

## 3. 端到端用戶流程

### 3.1 快速設定流程 (Quick Setup Flow)
1.  **進入系統**：使用者登入後，系統預設載入 Quick Setup 首頁。
2.  **檢視狀態**：首頁呈現當前閘道器整體健康燈號與核心流量。
3.  **建立新路由**：點擊「快速建立路由」。
    *   Step 1：輸入服務名稱與後端目標 (Target URL)。
    *   Step 2：選擇預設安全設定（如：開啟基礎 CORS、預設限流）。
    *   Step 3：確認並發佈。
4.  **完成與反饋**：顯示成功動畫，提供一鍵複製配置或測試路由的按鈕。

### 3.2 專家工作台流程 (Expert Workbench Flow)
1.  **切換模式**：在 Quick Setup 畫面右上角，點擊「切換至 Expert Workbench」。
2.  **全域檢視**：畫面展開為多欄位或數據密集型的 Dashboard，顯示所有詳細配置表單與監控圖表。
3.  **深度配置**：
    *   使用者選取特定路由，展開右側抽屜 (Drawer) 進行進階設定。
    *   修改 Middleware 順序，或打開 JSON 編輯器直接寫入特定標頭過濾規則。
4.  **模擬與套用**：點擊「Dry Run」測試配置有效性（呼叫既有 API 驗證），確認無誤後點擊「Apply Changes」生效。

---

## 4. 共同狀態與相容策略（避免功能破壞）

為確保「不破壞任何既有功能與 API 契約」，前端架構必須實施嚴格的狀態管理與資料轉換策略。

*   **單一資料來源 (Single Source of Truth)**：Quick Setup 與 Expert Workbench 共享同一個 Redux/Zustand Store 或 React Query 快取。底層資料結構嚴格遵循既有 API Schema。
*   **API 契約隔離 (DTO Mapping)**：
    *   Quick Setup 收集到的精簡資料，在提交前必須透過 Adapter 轉換，補齊 API 所需的預設值 (Default Values)，以符合既有的 Payload 格式。
    *   從 API 讀取的複雜資料，Adapter 負責萃取出 Quick Setup 需要的關鍵欄位；若有 Quick Setup 不支援的進階欄位，需保留其原始值，避免在 Update 時被覆寫遺失。
*   **降級與鎖定機制**：如果在 Expert Workbench 中設定了極度複雜的規則（例如手寫腳本），當使用者切換回 Quick Setup 檢視該路由時，系統應顯示「此為進階配置，請至 Expert Workbench 編輯」的唯讀狀態，防止簡單表單意外洗掉進階設定。

---

## 5. 路由與頁面結構

採用宣告式路由，整合 Feature Flag 控制存取。

*   `/` -> Redirect to `/quick-setup` (預設)
*   `/quick-setup`
    *   `/quick-setup/dashboard` (快速看板)
    *   `/quick-setup/routes/new` (精靈引導)
*   `/expert-workbench`
    *   `/expert-workbench/overview` (全域拓撲/列表)
    *   `/expert-workbench/routes/:id/edit` (進階配置抽屜/面板)
    *   `/expert-workbench/advanced-policies` (全域策略)
    *   `/expert-workbench/raw-config` (JSON/YAML 編輯器)

---

## 6. 元件拆分與責任

遵循 Atomic Design 與 Container/Presentational 模式：

*   **共用核心層 (Shared Core)**：
    *   `API Client / Hooks`：負責與既有 `go_gateway` API 溝通（GET/POST/PUT/DELETE），確保契約不變。
    *   `FeatureFlagProvider`：全局 Context，控制新舊 UI 與雙入口的顯示邏輯。
    *   `UI Library Components`：按鈕、輸入框、表格等基礎設計系統元件。
*   **Quick Setup 專屬元件**：
    *   `WizardStepper`：負責管理多步驟表單狀態。
    *   `SimpleMetricsCard`：高度視覺化、簡化的指標卡片。
*   **Expert Workbench 專屬元件**：
    *   `ComplexDataGrid`：支援排序、過濾、行內編輯的資料表。
    *   `CodeEditorPanel`：整合 Monaco Editor 或相似套件的原始碼編輯器。
    *   `PolicyRuleBuilder`：視覺化的規則邏輯樹編輯器 (AND/OR 條件組合)。

---

## 7. 互動與文案規格

統一兩入口的體驗語言，確保狀態反饋清晰明確。

*   **Loading (載入中)**：
    *   *Quick Setup*：使用骨架屏 (Skeleton) 以降低視覺跳動，文案：「正在準備您的快速配置環境...」。
    *   *Expert Workbench*：使用頂部進度條 (Top Progress Bar) 或局部 Spinner，不干擾其他靜態資料的閱讀。
*   **Empty (空狀態)**：
    *   *Quick Setup*：顯示插畫與明確的 Call-to-Action (CTA)，文案：「目前尚無路由配置。點擊下方按鈕，3 分鐘內建立您的第一個路由！」。
    *   *Expert Workbench*：顯示簡潔的提示框與快捷指令，文案：「查無進階配置規則。您可以透過 JSON 匯入或手動新增策略。」
*   **Error (錯誤狀態)**：
    *   *共通*：攔截所有 API 錯誤。若是已知驗證錯誤，標記於表單欄位下方；若是伺服器錯誤，彈出 Toast/Snackbar。
    *   *文案*：「配置更新失敗：[API 錯誤代碼/詳細原因]。請檢查網路狀態或配置格式是否正確。」
*   **Success (成功狀態)**：
    *   *Quick Setup*：全螢幕或明顯的成功動畫，文案：「太棒了！路由已成功上線並開始接收流量。」
    *   *Expert Workbench*：低干擾的角落 Toast，文案：「設定檔已儲存並生效。」

---

## 8. Feature Flag 灰度策略

為落實「決策 2B」，所有新 UI 流程均受控於 Feature Flag (FF) 系統。

*   **Flag 定義**：`enable_ui_v2_dual_entry` (boolean)
*   **控制層級**：
    *   **Level 1: UI 進入點攔截**：登入後，系統向 FF Server (或讀取靜態配置) 獲取狀態。若為 `false`，路由導向舊版 UI；若為 `true`，導向新的 Quick Setup 預設頁面。
    *   **Level 2: 無縫切換 (Opt-in/Opt-out)**：在灰度期間，新版 UI 頂部保留「返回舊版 (Legacy View)」的按鈕，舊版 UI 也提供「體驗新版」按鈕，確保萬一新版有問題，使用者能自救。
*   **推進階段**：
    *   *Phase 1 (Alpha)*：內部開發與 QA 團隊 (白名單)。
    *   *Phase 2 (Beta)*：開放 10% 的真實使用者，並監控 API 錯誤率與前端 Crash Logs。
    *   *Phase 3 (GA)*：100% 全面開放，隨後移除舊版 UI 程式碼。

---

## 9. 測試策略

在不破壞現有 API 的前提下，測試重心放在資料轉換的正確性與 UI 流程的完整性。

*   **單元測試 (Unit Testing)**：
    *   **重點**：DTO Adapter 與 Data Mapper。
    *   **斷言**：確保 Quick Setup 產生的精簡物件，經過 Mapper 後能 100% 補齊既有 API 規定的 Default 欄位；確保 Expert Workbench 儲存複雜 JSON 時，欄位不會遺失。
*   **整合測試 (Integration Testing)**：
    *   **重點**：React 元件與狀態管理的協作。
    *   **斷言**：在 Quick Setup 填寫表單後，Redux/Zustand Store 中的狀態更新正確，且切換到 Expert Workbench 時，該筆資料能正確無誤地顯示在複雜表單中。
*   **端到端測試 (E2E Testing - Cypress / Playwright)**：
    *   **重點**：核心用戶旅程與雙入口切換。
    *   **情境 1**：以 Feature Flag 開啟狀態進入，完成 Quick Setup 新增路由的完整流程。
    *   **情境 2**：在 Expert Workbench 進行進階欄位修改，儲存後驗證 API 請求 Payload 是否正確無誤。
    *   **情境 3**：從 Expert Workbench 寫入複雜規則，切換回 Quick Setup 驗證「降級鎖定保護機制」是否成功觸發，防止誤改。

---

## 10. 交付里程碑 (W1~W4)

採用並行開發策略（決策 3C），於 4 週內交付第一版。

*   **Week 1：架構奠基與 API 橋接**
    *   設置 Feature Flag 基礎設施。
    *   建立共用的 API Client、DTO Adapters 與 State Management。
    *   完成新版 Layout 框架與路由配置 (建立 Quick 與 Expert 空殼頁面)。
*   **Week 2：Quick Setup 實作 (雙線推進 A)**
    *   開發 Dashboard 指標看板。
    *   實作「快速建立路由」Stepper 表單。
    *   串接 API，確保基礎配置能成功寫入。
*   **Week 3：Expert Workbench 實作 (雙線推進 B)**
    *   開發全域列表與資源拓撲組件。
    *   整合 JSON/YAML 編輯器。
    *   實作進階策略表單區塊。
    *   實作 Quick <-> Expert 狀態互斥與降級鎖定保護邏輯。
*   **Week 4：測試、打磨與灰度發佈**
    *   執行 Unit/Integration/E2E 測試修復。
    *   落實 Loading/Empty/Error/Success 互動規格。
    *   開啟 `enable_ui_v2_dual_entry` 給內部名單 (Alpha)，準備進行外部灰度 (Beta)。

---

## 11. 風險與回滾

*   **風險 1：Adapter 邏輯遺漏導致 API 報錯或配置遺失。**
    *   *緩解機制*：強迫實施 TypeScript 嚴格型別檢查，對照既有 Swagger/OpenAPI 定義檔自動生成 Type。
*   **風險 2：Quick Setup 與 Expert Workbench 狀態同步競爭 (Race Condition)。**
    *   *緩解機制*：強制切換入口時觸發髒檢查 (Dirty Check)。若有未儲存的變更，提示使用者儲存或捨棄，禁止帶著髒狀態切換入口。
*   **回滾計畫 (Rollback Plan)**：
    *   若灰度期間監測到前端崩潰率上升或 API 400 錯誤激增超過設定閾值 (如 2%)。
    *   **一鍵降級**：直接透過 Feature Flag 服務將 `enable_ui_v2_dual_entry` 設為 `false`，前端熱重載後所有使用者即刻退回舊版 UI，因不涉及後端 API 改動，故為零風險回滾。

---

## 12. 驗收清單 (DoD)

- [ ] 登入後預設導向 Quick Setup (決策 1A 達標)。
- [ ] 支援透過設定中心或 URL 參數切換 Feature Flag，無縫展示新舊 UI (決策 2B 達標)。
- [ ] Quick Setup 與 Expert Workbench 核心流程均可正常走通 (決策 3C 達標)。
- [ ] 透過 Network 面板驗證：所有由新 UI 發出的 API Request Payload 完全符合既有 API 契約。
- [ ] 透過 Network 面板驗證：既有 API 欄位未在 Quick Setup 展現者，在更新時沒有被洗掉。
- [ ] 複雜進階配置在 Quick Setup 模式下正確顯示「唯讀/鎖定」提示。
- [ ] E2E 測試涵蓋率達到定義的 3 大核心情境，且 CI/CD Pipeline 綠燈。
- [ ] Loading, Empty, Error, Success 狀態皆按照規格文件實作完畢。
