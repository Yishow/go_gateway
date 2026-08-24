## Context

在 `/studio/v2` 引導式 4 步驟工作台中，使用者可加入不同通訊協議（Modbus TCP/RTU/UDP、FATEK FBs、Mitsubishi MC 3E）的工業設備。然而，從 Step 2（來源規則）到 Step 3（標籤映射）再到 Step 4（輸出綁定），系統多處存在「僅支援 Modbus 數值暫存器（40001）」的隱式假設，導致非 Modbus 設備在完整流程中產生以下斷裂：
1. Step 2 中衍生點位丟失字母前綴或被硬塞 Modbus 功能碼。
2. Step 3 中產生的預設 Tag Key 名稱殘缺（如 `dev.r0`），即時讀取無法被 Connector 識別。
3. Step 4 中跨協議 Modbus Share（將 MC 點位分享成 Modbus Slave）的輸出映射與資料庫欄位自動匹配無法正常關聯。

## Goals / Non-Goals

**Goals:**
- **步驟 2 來源規劃全適配**：統一由 `addressParser` 負責所有協議（Modbus/FATEK/MC3E）之位址步進（Dec/Hex、Word/Bit），完整保留區域前綴，並將功能碼推斷轉為協議感知。
- **步驟 3 標籤與映射全適配**：使 `mappingDefaults` 正確產生可讀的非 Modbus Tag Key，並確保即時數值預覽能正確傳遞點位位址至後端。
- **步驟 4 輸出與交付全適配**：確保 `autoAssignTargets` 支援非 Modbus 標籤自動匹配，並確保 MC/FATEK 來源設備能正確啟用 Modbus Share 暫存器輸出與資料庫激活。
- **可延續的 Share 生命週期**：rule-level Share 設定由 source rule 持久化，workspace hydration 恢復設定，activation 再投影到 runtime mappings；durable configuration 與 runtime projection 保持單向責任邊界。
- 完整補齊 Step 2 ~ Step 4 跨協議的前端單元與整合測試。

**Non-Goals:**
- 不重構底層 Connector 驅動程式（後端驅動已具備多協議讀寫能力）。
- 不修改 `internal/datalink/connector/adapters/` 或 `internal/protocol/` 的底層連線與傳輸行為。
- 不修改 legacy `/studio` fallback（包含 `sourceCanvasModel.ts`）；本 change 僅處理 `/studio/v2`。
- 不為通過行數檢查而壓縮無關 reducer/callback，也不新增沒有外部 caller 的 exports 或預留介面。
- 不引入 MQTT 來源點位規劃（MQTT 現階段維持 JSON Payload 訂閱型態）。

## Decisions

### 1. 統一以 `addressParser` 作為跨步驟的位址計算唯一真實來源 (Single Source of Truth)
- **決策**：在 `sourceRule.ts`、`RuleEditor.tsx`、`RangeSummary.tsx`、`mappingDefaults.ts` 等處，全面呼叫 `addressParser.offset()`、`getDefaultPlannerStartAddress()` 與 `validate()`。
- **好處**：徹底消除正則表達式過濾字元的 Bug，且集中維護三菱 Hex（X/Y/B）與 FATEK/Modbus 的步進規則。

### 2. 點位暫存器型別與功能碼語意解耦
- **決策**：擴充 `Point` 的 `function` 屬性，當 protocol 為 Modbus 時顯示 `holding_register`/`coil` 等；當為 MC3E/FATEK 時顯示暫存器區域類型（如 `D (Word)` 或 `M (Bit)`）。

### 3. Step 3 標籤生成與 Step 4 欄位匹配跨協議感知
- **決策**：`buildDefaultTagKey` 與 `autoAssignTargets` 支援英數字暫存器識別（如 `D0` -> `tag_d0` / `sensor_d0`），讓 Step 4 自動匹配資料庫欄位時精準鎖定對應名稱。

### 4. Modbus Share 跨協議輸出轉換
- **決策**：Modbus Share 是「輸出層」的機制，無論來源設備是 MC 3E 還是 FATEK，其 `share_start_register` 均為本機 Modbus Slave 位址（`40001` 起算）。前端在計算 ShareLayout 時，只需確保每個來源點位依序分配到獨立的 Modbus 暫存器即可。
- **持久化邊界**：`share_enabled`、`share_start_register` 與 `share_stride` 屬於 source rule 的 durable configuration，必須隨 rule create/update/list/get 往返；runtime `modbus-share/mappings` 僅是 activation projection，不可當成持久化真實來源。
- **Activation 邊界**：在 workspace activation 前先將當前 workspace 已持久化的 tag 同步到 runtime mapping；UI `40001` 必須轉成 API zero-based register `0`。只能刪除當前 workspace 擁有且已不在 desired set 的 tag mapping，任何同步失敗都必須阻擋後續 activation。
- **全域開關優先**：`settings.modbus_share.enabled` 為 activation 與 Step 4 Share summary 的總開關；全域停用時不得因個別 rule 的 `share_enabled` 而建立、顯示或同步 mappings。
- **Ownership 證明**：workspace-owned tag 集合只能來自目前 workspace 已持久化且可由 source rule/tag 關聯證明的 identifiers；前端暫存 mapping 或僅有非空 `tag_id` 的資料不得取得刪除 runtime mapping 的權限。

## Implementation Contract

- **Observable Behavior**:
  - **Step 2**: 建立 MC 設備規則時預設 `D0`，輸入 `D100` 且 count 為 4 時衍生 `D100, D101, D102, D103`；建立 FATEK 規則時預設 `D0`。
  - **Step 3**: 自動生成之標籤名稱為 `dev_01.sensor_0.rd100`，即時數值預覽正常連線且顯示即時數值。
  - **Step 4**: 自動分配 DB 欄位時能正確匹配 `sensor_d100` 或預設欄位；啟用 Modbus Share 時顯示 `40001 ~ 40004` 之輸出暫存器映射。
  - **Restart/Hydration**: source rule 的 `share_enabled`、`share_start_register`、`share_stride` 經 create/update/list/get 與重新載入後保持一致；runtime mapping 不作為 hydration 資料源。
- **Interface / Data shape**:
  - `derivePoints(rule: Rule, deviceId: string, skippedSet: Set<string>, protocol?: ProtocolType): Point[]`
  - `buildDefaultTagKey(point: Point, protocol?: ProtocolType): string`
  - `autoAssignTargets(enabledPoints: Point[], mappings: Record<string, Mapping>, columnNames: string[], existingTargets: Record<string, DbTarget>): Record<string, DbTarget>`
  - Source rule API request/response fields: `share_enabled: boolean`, `share_start_register: number`, `share_stride: number`。
  - Runtime mapping register uses zero-based addressing: UI `40001` maps to API `register: 0`。
- **Failure modes**:
  - 若使用者輸入無效位址，`addressParser` 返回驗證錯誤，UI 標記紅色提示，不允許繼續流轉至下一步。
  - 位址偏移 API 遇到無效輸入時必須明確失敗；點位衍生層在驗證失敗時回傳空集合，禁止以原字串或 Modbus 語意靜默降級。
  - 當點位的 device、rule 或 address identity 改變時，Step 3 必須重新同步 `state.points` 並重建該點位的預設 mapping，避免沿用舊位址 tag。
  - Rule-level Modbus Share 設定必須透過既有輸出持久化與 activation 流程保存；重新 hydration 後不得回復為停用。
  - 前後端位址偏移必須遵守同一協議語意；MC 3E `X/Y/B` 採十六進位，其他已支援區域依規格採十進位。無效位址不得回傳原字串作為成功結果。
  - 任一 enabled rule 位址無效時，即使其他 rule 有有效點位，Step 2 仍不得繼續。
  - 全域 Modbus Share 停用時不得顯示或同步 rule-level Share mappings。
  - runtime stale mapping 只有在其 tag 被 durable workspace ownership 集合證明屬於目前 workspace 時才可刪除；無法證明時保留並回報同步失敗或跳過刪除，不得猜測 ownership。
- **Acceptance criteria**:
  - `npm run test` 全部通過，新增 Step 2、3、4 多協議端到端單元測試。
  - `go test ./...`、`go vet ./...` 通過，且 source rule API/repository/migration 覆蓋 Share round-trip、MC 3E hex offset 與 invalid-address fail-closed。
  - `spectra validate fix-protocol-address-adaptation-v2 --strict` 通過；變更檔案符合 500 行硬上限，且新增 UI 文字全部走 i18n。

## Risks / Trade-offs

- **[Risk] 使用者在 MC Protocol 輸入十六進位位址 (X0~XF) 步進跨界** → **Mitigation**: `addressParser.offset()` 針對 X/Y/B 採用 16 進位步進（0x0F + 1 = 0x10），確保 Hex 序列正確。
- **[Risk] activation 清除其他 workspace 或外部建立的 runtime mapping** → **Mitigation**: 刪除前要求 durable workspace ownership 證據，無證據時 fail-closed。
- **[Risk] migration 後需回復舊版** → **Mitigation**: up migration 只新增具安全預設值的 Share 欄位；down migration 只移除本 change 新增欄位，不改寫既有 source rule 資料。
