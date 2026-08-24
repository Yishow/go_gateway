## 1. 核心點位衍生與位址計算 (Step 2 - TDD: Red-Green-Refactor)

- [x] 1.1 [TDD-Red] 建立 Protocol-aware point address derivation 單元測試，定義 MC 3E (`D0` -> `D3`, `D100` -> `D102`, `X0` -> `X2`)、FATEK (`R0` -> `R2`) 與 Modbus (`40001` -> `40003`) 的衍生點位位址與前綴保留預期，驗證現有 `derivePoints` 失敗（驗證：`npm test sourceRule.test.ts` 失敗）。
- [x] 1.2 [TDD-Green] 重構 `sourceRule.ts` 的 `derivePoints` 與 `formatAddr`，統一以 `addressParser` 作為跨步驟的位址計算唯一真實來源 (Single Source of Truth)，使用 `addressParser.offset()` 與協議型別遞增位址，並落實點位暫存器型別與功能碼語意解耦，在非 Modbus 協議時避免誤用 Modbus 功能碼推斷（驗證：`npm test sourceRule.test.ts` 通過）。
- [x] 1.3 [TDD-Refactor] 整理 `sourceRule.ts` 與點位推斷介面型別定義，確保跨模組引用時具有完整 TypeScript 型別安全（驗證：`npm run build` 通過型別檢查）。

## 2. 來源規劃狀態管理與 UI 自適應 (Step 2)

- [x] 2.1 更新 `RuleTabRail.tsx` 與 `useWorkbenchV2State.ts` 實作 Multi-rule tab management，在執行 `handleAddRule` 及初始狀態建立時，透過選取設備之 protocol 動態調用 `getDefaultPlannerStartAddress(protocol)`（驗證：在 `RuleTabRail.test.tsx` 測試新增 MC 設備規則時起始位址為 `D0`，Modbus 為 `40001`）。
- [x] 2.2 更新 `RuleEditor.tsx`、`RangeSummary.tsx` 與 `MergedPointTable.tsx` 實作 Rule editor with linked reset，依設備協議動態呈現位址輸入佔位符、範圍摘要與非 Modbus 暫存器類型標籤（驗證：`npm test RuleEditor.test.tsx` 與 `RangeSummary.test.tsx`）。

## 3. 標籤生成與即時數值映射適配 (Step 3)

- [x] 3.1 [TDD-Red] 建立 Protocol-aware default tag key generation 單元測試，驗證 MC 3E 與 FATEK 點位生成的 Tag Key 具備正確的暫存器前綴（如 `dev.sensor.rd0`）（驗證：`npm test mappingDefaults.test.ts` 失敗）。
- [x] 3.2 [TDD-Green] 更新 `mappingDefaults.ts` 與 `Step3Mapping.tsx` 實作 Step 3 標籤生成與 Step 4 欄位匹配跨協議感知，支援非 Modbus 點位語意化 Tag 鍵名與即時預覽數據關聯（驗證：`npm test mappingDefaults.test.ts` 通過）。

## 4. 輸出目標綁定與跨協議 Modbus Share (Step 4)

- [x] 4.1 更新 `autoAssignTargets.ts` 實作 Cross-protocol Modbus Share and database target binding，強化對異質協議標籤後綴（如 `sensor_d0`）之資料庫欄位智慧匹配（驗證：`npm test autoAssignTargets.test.ts`）。
- [x] 4.2 驗證非 Modbus 來源規則在開啟 Modbus Share 時落實 Modbus Share 跨協議輸出轉換，正確計算 `40001+` Slave 輸出暫存器映射，並通過 Readiness 激活校驗（驗證：`npm test Step4Database.test.tsx`）。
- [x] 4.3 執行前端完整驗證流程，確保所有單元測試與靜態型別編譯無任何錯誤（驗證：`cd frontend && npm run lint && npm run test && npm run build`）。

## 5. Code review 修復（TDD）

- [x] 5.1 修正 Step 2 `MODIFIED` requirements，使既有 scenarios 與跨協議新增語意完整合併，並通過 `spectra validate fix-protocol-address-adaptation-v2 --strict`。
- [x] 5.2 [TDD] 讓無效協議位址 fail-closed、補齊 Rule tab 契約、MC/FATEK range/placeholder 測試與新增 UI 文字 i18n；驗證無效位址顯示錯誤且不能前往 Step 3。
- [x] 5.3 [TDD] 讓 Step 3 point identity 納入 device/rule/address，地址改變時同步重建 mapping 與 protocol-aware display name，並移除未生效的 protocol API 參數。
- [x] 5.4 [TDD] 將 rule-level Modbus Share 接入既有持久化/activation seam，於 Step 4 顯示非 Modbus `40001 ~ 40004` 映射，並驗證 hydration 後設定不遺失。
- [x] 5.5 拆分超過 500 行的 Step 4 測試檔，保持正式測試入口與行為覆蓋不變，並通過行數 gate。
- [x] 5.6 [TDD-Red/Green] 對齊 backend source-rule 與 frontend `addressParser` 的 protocol radix 與 fail-closed 契約：MC 3E `X/Y/B` 採十六進位、無效位址回傳明確錯誤且不得以原字串成功降級（驗證：新增 Go service/validation tests，涵蓋 `X0 + 16 = X10` 與 invalid input）。
- [x] 5.7 [TDD-Red/Green] 修正 Step 2 多規則 continuation gate，使任一 enabled rule 位址無效時都阻擋進入 Step 3，且錯誤可定位至該 rule（驗證：`step2-rule.test.tsx` 新增一有效一無效規則案例）。
- [x] 5.8 [TDD-Red/Green] 收斂 Share activation safety：`settings.modbus_share.enabled=false` 時 summary 與 runtime sync 全停用；stale deletion 僅接受 durable workspace persisted-tag ownership，非空 frontend `tag_id` 不構成 ownership（驗證：`step4-share.test.tsx` 與 `step4-share-activation.test.tsx` 覆蓋全域停用及不刪除外部 mapping）。
- [x] 5.9 移除本 change 的越界變更：不得保留 `internal/datalink/connector/adapters/`、`internal/protocol/`、legacy `/studio` `sourceCanvasModel.ts` 的行為修改；還原無關 reducer/callback 壓縮並移除無 caller 的新 exports，同時保留 Step 2–4 核心與 Share persistence 行為（驗證：相對 fixed point 的上述越界路徑無 diff，`npm run lint` 與 `go test ./...` 通過）。
- [ ] 5.10 執行 `spectra validate fix-protocol-address-adaptation-v2 --strict`、`spectra analyze fix-protocol-address-adaptation-v2 --json`、`go test ./...`、`go vet ./...`、前端 `npm run lint`、預設 worker 的完整 `npm run test`、`npm run build`、行數 gate 與 `git diff --check`；任一必要 gate 失敗時保持未完成並如實回報環境或程式缺陷。
