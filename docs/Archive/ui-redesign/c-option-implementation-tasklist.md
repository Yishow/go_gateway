這是一份針對 `go_gateway` 雙入口 UI（C 版激進方案）量身打造的落地任務清單。此計畫確保在**不修改後端 API 契約**的前提下，透過前端架構的重構與 Adapter 模式，實現新手與專家體驗的物理隔離。

---

# 🚀 `go_gateway` 雙入口 UI (C版激進方案) 落地任務清單

## 1. C 版核心差異摘要（相對現況）
*   **現況盲點**：單一龐大表單，初學者容易被複雜的 Plugin 與 Regex 設定勸退，專家則覺得畫面佈局不夠緊湊。
*   **C 版激進方案架構**：
    *   **物理隔離**：在路由層級徹底拆分，廢棄共用表單，改為 `Quick Setup`（引導式嚮導）與 `Expert Workbench`（高密度儀表板/YAML編輯器）。
    *   **API 無縫接軌**：透過純前端的 `Adapter Layer` 攔截，無論哪個入口，最終送出的 JSON Payload 皆 100% 吻合既有 API 契約。
    *   **降級與升級**：支援從 Quick 升級至 Expert（無損），但從 Expert 降級至 Quick 時引入防呆與資料遺失警告。

## 2. 路由/頁面改動清單（檔案層級建議）
新增專屬路由模組，不干擾既有舊版元件，以利 Feature Flag 控制。

*   `src/router/gateway.ts` (新增路由定義)
    *   `/gateway/entry` - 雙入口選擇頁（Landing）
    *   `/gateway/quick-setup` - 快速設定主頁
    *   `/gateway/expert-workbench` - 專家工作台主頁
*   `src/views/gateway/Entry/index.vue` - 新增入口選擇畫面（卡片式二選一）。
*   `src/views/gateway/QuickSetup/*` - 新增精靈模式頁面目錄。
*   `src/views/gateway/ExpertWorkbench/*` - 新增專家模式頁面目錄。
*   `src/views/gateway/List.vue` - 修改既有列表頁的「新增/編輯」按鈕行為（套用灰度邏輯）。

## 3. 元件拆分與責任
嚴格遵守職責分離，避免跨模式的程式碼污染。

### 🔄 共通元件 (Shared)
*   **`PayloadPreviewModal`**: 唯讀彈窗，提交前可預覽轉換後的最終 JSON 結構。
*   **`GatewaySubmitGroup`**: 統一的儲存/取消按鈕列，處理防抖 (Debounce) 與全域 Loading 狀態。

### 👶 Quick Setup (新手嚮導)
*   **`IntentForm`**: 基於「意圖」的問卷式表單（例：「我要對外暴露內部服務」、「我要設定簡單的反向代理」）。
*   **`SimpleRouteBuilder`**: 隱藏正規表達式，只提供下拉選單與基礎前綴比對的路徑設定器。
*   **`BasicAuthToggle`**: 將複雜的 Auth Plugin 抽象為單一的「開啟/關閉」安全開關。

### 🧙‍♂️ Expert Workbench (專家工作台)
*   **`AdvancedRouteTable`**: 支援正則表達式、權重排序、多路由拖曳排序的高密度資料表。
*   **`PluginChainVisualizer`**: 視覺化的中介軟體/插件執行鏈編輯器。
*   **`RawManifestEditor`**: 整合 Monaco Editor，支援直接貼上 YAML/JSON 並具備 Schema 即時校驗。

## 4. 狀態管理與 Adapter 任務
這是 C 版方案不破壞 API 的核心防線。

*   **狀態管理 (State)**
    *   **任務**: 建立 `useGatewayDraftStore`。
    *   **責任**: 負責暫存當前草稿，包含欄位：`mode` (quick/expert)、`quickData`、`expertData`。
    *   **機制**: 實作 `switchMode()` 函數。Quick 切 Expert 時，進行資料映射；Expert 切 Quick 時，觸發警告機制。
*   **Adapter 層 (`src/adapters/gatewayAdapter.ts`)**
    *   **任務**: 實作 `quickToPayload(draft)`：負責將極簡表單資料「補齊」API 必填的預設值（例如注入預設的 Timeout 參數、預設 Plugin）。
    *   **任務**: 實作 `expertToPayload(draft)`：直通 API 契約，僅做格式化。
    *   **任務**: 實作 `payloadToDraft(apiResponse)`：用於「編輯」模式時，將後端吐回的複雜 JSON 反向解析回 UI 狀態（若偵測到複雜設定，強制進入 Expert 模式）。

## 5. i18n 與文案任務
*   **任務**: 於 `locales/` 新增 `gateway.dual.json` 命名空間。
*   **Entry 頁文案**: 撰寫具備引導性的卡片文案（例如：「快速設定：適合 90% 的常見代理場景，只需 3 分鐘」、「專家工作台：解鎖正規表達式、自訂插件與進階負載均衡」）。
*   **Quick 模式轉譯**: 將工程術語轉化為白話文（如將 `Upstream Host` 改為 `目標伺服器位址`）。
*   **警告文案**: 「切換回快速模式將會捨棄您設定的 3 個進階外掛與正則路由，是否確定降級？」

## 6. 測試任務（單元/整合/E2E）
*   **單元測試 (Unit Test)**:
    *   **高優先級**: 針對 `gatewayAdapter.ts` 撰寫 Jest/Vitest 測試。準備至少 10 組 mock data，確保 `quickToPayload` 產出的 JSON 通過後端 OpenAPI Schema 校驗。
*   **整合測試 (Integration Test)**:
    *   測試 `useGatewayDraftStore` 的狀態機，確保在兩種模式間切換時，資料轉換沒有發生不可預期的 Mutation。
*   **E2E 測試 (Cypress / Playwright)**:
    *   **案例 A**: 從入口頁進入 Quick Setup，填寫必填後送出，驗證 API Request 正確。
    *   **案例 B**: 從入口頁進入 Expert，輸入 YAML 後送出，驗證 API Request 正確。
    *   **案例 C**: 點擊既有複雜 Gateway 的「編輯」，驗證系統自動導向 Expert 模式且回填正確。

## 7. Feature Flag 灰度任務
*   **任務**: 引入變數 `ENABLE_GATEWAY_DUAL_ENTRY` (預設為 `false`)。
*   **UI 邏輯**: 
    *   若為 `false`: 點擊「建立」直接進入 `/gateway/legacy-form`。
    *   若為 `true`: 點擊「建立」進入 `/gateway/entry` 雙入口頁面。
*   **埋點監控**: 送出 API 時，在 Header 追加 `X-UI-Version: dual_quick` 或 `X-UI-Version: dual_expert`，用於後續分析哪種模式的成功率與使用率較高。

## 8. 風險與回滾
*   **風險 1：編輯舊資料解析失敗**
    *   *對策*: 舊資料由 `payloadToDraft` 解析若發生例外錯誤，跳出 Toast 提示「該配置過於複雜，已回退至傳統表單」，並平滑降級至舊版 UI。
*   **風險 2：Quick 模式預設值導致線上行為改變**
    *   *對策*: Adapter 補齊的預設值必須與現有舊版表單的預設值 **100% 保持一致**，由 Unit Test 把關。
*   **回滾策略 (Zero-Downtime Rollback)**
    *   若上線後發生嚴重阻斷，維運人員只需將配置中心的 `ENABLE_GATEWAY_DUAL_ENTRY` 切為 `false`，前端熱重載後即時恢復為單一舊版表單，完全不影響後端服務。

## 9. 四週 WBS（每週可交付項）

| 週期 | 重點任務目標 | 每週可交付成果 (Deliverable) |
| :--- | :--- | :--- |
| **W1** | **架構底層與 Adapter**<br>- 實作 Feature Flag 控制邏輯<br>- 建立基礎路由與 Entry 選擇頁面<br>- 開發 `gatewayAdapter.ts` 雙向轉換邏輯 | 1. 可點擊預覽的 Entry 路由。<br>2. 覆蓋率 > 90% 的 Adapter 單元測試報告。 |
| **W2** | **Quick Setup 核心開發**<br>- 實作 `IntentForm` 等極簡元件<br>- 串接 Store 與 API 提交邏輯<br>- 補充初階 i18n 文案 | 透過 Quick Setup 模式能成功打通真實 API，建立基礎 Gateway 配置。 |
| **W3** | **Expert Workbench 核心開發**<br>- 整合 Monaco Editor<br>- 實作 `AdvancedRouteTable`<br>- 完成 Quick/Expert 狀態切換與警告機制 | Expert 模式可正常運作，且模式互相切換時防呆機制（警告彈窗）符合預期。 |
| **W4** | **驗收、測試與灰度上線**<br>- 撰寫並跑通 3 條核心 E2E 腳本<br>- 處理舊資料「編輯」的回填適配<br>- 執行 QA 測試並啟動 10% 灰度 | **Feature 交付**：正式環境灰度開啟，雙入口功能上線供部分使用者使用。 |

## 10. 完成定義 DoD（可勾選）
- [ ] `ENABLE_GATEWAY_DUAL_ENTRY` Feature Flag 已實作，且關閉時系統表現與重構前完全一致。
- [ ] `gatewayAdapter.ts` 單元測試覆蓋率大於 90%，確保不破壞既有 API 契約。
- [ ] Quick Setup 模式送出的 API Payload 通過後端 Schema 驗證，預設值無誤。
- [ ] Expert Workbench 模式的編輯器（YAML/JSON）支援即時防呆校驗。
- [ ] 模式降級（Expert -> Quick）的二次確認警告彈窗已實作並測試通過。
- [ ] 所有新頁面與元件支援 i18n 語系切換，文案無 hardcode。
- [ ] E2E 測試涵蓋：新增 Quick、新增 Expert、編輯既有複雜 Gateway (自動導向 Expert)。
- [ ] 送出 API 時的 Header 已帶上埋點追蹤標記 (`X-UI-Version`)。
