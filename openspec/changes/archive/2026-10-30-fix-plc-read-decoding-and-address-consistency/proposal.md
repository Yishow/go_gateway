## Problem

目前的 PLC 讀取可能回報品質良好卻取得錯值、錯位址或不完整資料。這次以已重現的共用 byte order、FATEK 型別轉換、MC 位址規劃、Modbus raw address 與 FATEK Cmd 44 回應驗證為修復範圍，避免錯誤值繼續進入採集與映射流程。

## Root Cause

- 共用暫存器轉換只處理 word 順序，遺漏 `BADC`／`DCBA` 的 word 內 byte swap。
- FATEK 將浮點位元做數值轉型，64 位元型別會落入第一個暫存器的預設回傳；scalar count 也未辨識原生 32 位元元件。
- MC planner 與前端依舊規格使用 X/Y/B 十六進位、D/W/M 十進位，但實際 runtime 先以八進位解析 X/Y，以十六進位解析 W/B。通用點位驗證亦拒絕 `WA` 等合法 hex suffix。
- Modbus 先把短數字套成傳統區域編號，覆蓋明確指定的 Function；例如 `Address="100", Function="03"` 實際讀成 FC01、offset 99。
- FATEK Cmd 44 只把字元 `1` 視為 true，未檢查長度與其他非法字元。

## Proposed Solution

- 修正共用 32／64 位元 signed、unsigned、IEEE-754 解碼，讓四種既有 byte order 都有確定的位元語意，保留既有預設值與 16 位元行為。
- 讓 FATEK 依元件寬度組合完整 typed scalar；以高 word 在前維持既有整數慣例，缺少資料回報錯誤。scalar 呼叫按位址寬度計數，直接 batch 的元件數語意維持不變。
- **BREAKING：修正 MC 新規則的位址規劃契約**為 X/Y 八進位、W/B 十六進位、D/M 十進位，以目前 runtime 的實際讀取位址為相容基準；前後端與建點驗證一致，area 集合不擴充。runtime 既有 X/Y 十進位 fallback 暫保留，只供舊直接讀取相容。
- 禁止 radix 修正自動改寫舊規則的點位；既有 MC rule 若按新文法展開會產生不同 address set，更新必須在任何寫入、清理或重建前拒絕，保留 IDs、位址與映射。Step 2 載入既有 MC 規則時，使用帶 scope/revision 的既有 candidates GET 取得 persisted addresses；缺失、過期或規劃衝突時阻止自動儲存／映射重綁，不以新 planner 的假設顯示成功。
- **BREAKING：釐清 Modbus raw numeric 優先序**：明確 Function 搭配非傳統固定格式數字，使用 zero-based raw offset；既有傳統 5 位區域編號及 HR/IR/CO/DI alias 維持既有語意。舊短數字點位可能改變實際讀取目標，必須揭露影響，不自動修改儲存資料。
- Cmd 44 僅接受數量吻合、全由 `0`／`1` 組成的成功 payload；不合法資料回傳 error 與 bad quality，不能冒充有效新樣本。

## Success Criteria

- 設計驗收向量：同一 `float32(1)` 與 `uint64(0x0102030405060708)` 在 ABCD/BADC/CDAB/DCBA 均解碼相同；MC 與 Modbus connector 的實際讀取路徑皆涵蓋。
- FATEK `D0` 的 `3F800000` 與 `DR0` 的 `3F800000` 皆回傳 typed `float32(1)`；64 位元資料完整保留，資料不足回報錯誤。
- MC `W0` offset 10 產生 `WA`、offset 16 產生 `W10`；`X0` offset 16 產生 `X20`。預覽、儲存點位與 runtime offset 對應一致，舊規則更新的 guard 不留下部分寫入。
- Modbus `100` + FC03 的模擬 wire request 為 FC03、offset 100；傳統位址及 alias 回歸通過。
- FATEK Cmd 44 空、短、長及非二進位 payload 全部拒絕，合法 payload 維持正確 scalar／batch 行為。
- 每個修復先在未改產品程式上取得可重現的失敗測試，再通過受影響模組與呼叫端測試；真 PLC／Windows／現場驗收另外標示未執行。

## Impact

- Affected specs / Modified capabilities:
  - `protocol-connectors`：byte order、FATEK typed scalar/count、Modbus raw address 與失敗資料契約。
  - `datalink-workbench-v2-step2-rule`：修正 MC radix、合法建點與既有規則保護；明確取代舊的 X/Y/B hex、D/W/M decimal 契約。
- New capabilities：無；沿用現有協議與採集架構，不增加 dependency 或資料庫 migration。
- Affected code:
  - Modified: `lib/hsllogic/bytetransform.go`, `lib/hsllogic/register_conversion.go`。
  - Modified: `internal/datalink/connector/adapters/fatek_helpers.go`, `internal/datalink/connector/adapters/fatek_operations.go`, `internal/protocol/fatek/client.go`。
  - Modified: `internal/datalink/connector/point_read.go`, `internal/datalink/collector/scheduler_polling.go`, `internal/datalink/collector/scheduler_manual_poll.go`, `internal/datalink/device/service_point_read.go`。
  - Modified: `internal/datalink/connector/adapters/modbus_parse_helpers.go`。
  - Modified: `internal/datalink/sourcerule/protocol_address_planner.go`, `internal/datalink/sourcerule/service.go`, `internal/datalink/point/service_batch_validation.go`, `frontend/src/utils/addressParser.ts`。
  - New: `lib/hsllogic/plc_byte_order_regression_test.go`, `internal/datalink/connector/adapters/plc_read_decoding_regression_test.go`, `internal/protocol/fatek/status_payload_test.go`, `internal/datalink/sourcerule/mc_radix_compatibility_test.go`, `internal/datalink/point/mc_hex_address_test.go`。
  - Modified: `internal/datalink/collector/point_count_contract_test.go`, `frontend/tests/unit/utils/addressParser.test.ts`, `frontend/tests/unit/workbench-v2/step2-rule.test.tsx`。
  - Modified: `frontend/src/features/datalink/workbench-v2/state/types.ts`, `frontend/src/features/datalink/workbench-v2/state/sourceRule.ts`, `frontend/src/features/datalink/workbench-v2/state/selectors.ts`。
  - Modified: `frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts`, `frontend/src/pages/datalink/workbench-v2/useStudioV2RuleAutosave.ts`, `frontend/src/pages/datalink/workbench-v2/useStudioV2MappingAutosave.ts`。
  - Modified: `frontend/src/features/datalink/workbench-v2/steps/step2/Step2Rule.tsx`, `frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx`, `frontend/src/features/datalink/workbench-v2/steps/step2/RangeSummary.tsx`。
  - Modified: `frontend/src/i18n/locales/en/workbench-v2.json`, `frontend/src/i18n/locales/zh-TW/workbench-v2.json`，只加入本次阻擋狀態文案。
  - New: `frontend/src/pages/datalink/workbench-v2/useStudioV2PersistedMCPoints.ts`, `frontend/src/features/datalink/workbench-v2/state/mcPersistedPoints.ts`, `frontend/tests/unit/workbench-v2/mc-persisted-hydration.test.tsx`, `frontend/tests/unit/workbench-v2/mc-persisted-view.test.tsx`。
  - Removed: 無。
- `docs/HslDocs_Categorized/` 與 `lib/` 作為查核來源；文件或跨語言範例的差異不直接等同目前 Go runtime 的 bug。詳細範圍、證據與回復限制記於 design。
- 本次只交付提案；所有產品實作與驗收工作維持未完成，不觸及現有前端 WIP，不 commit／push／部署。
