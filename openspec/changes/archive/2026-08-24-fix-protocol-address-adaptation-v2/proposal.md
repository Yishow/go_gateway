## Why

在 `/studio/v2` 工作台新增非 Modbus 協議設備（例如三菱 MC Protocol `mc_3e` 或永宏 FATEK `fatek_fbs`）時，系統在整個資料鏈路流程（步驟 2 來源規則、步驟 3 標籤映射與即時數值、步驟 4 輸出交付與 Modbus Share）中存在多處 Modbus 硬編碼假設：
1. **步驟 2**：寫死起始位址 `40001`，點位衍生演算法（`derivePoints`）濾除非數字字元導致 `D0`/`D100` 前綴遺失，且未區分 Bit/Word/Hex 步進與協議功能碼。
2. **步驟 3**：預設 Tag Key 生成（`buildDefaultTagKey`）與即時數值預覽因點位位址失真而產生錯誤標籤名，且未適配非 Modbus 資料型別與暫存器語意。
3. **步驟 4**：資料庫欄位自動匹配（`autoAssignTargets`）無法精確關聯非 Modbus 標籤；跨協議的 Modbus Share（將 MC/FATEK 來源點位轉化為本機 Modbus Slave 暫存器輸出）缺少協議轉換驗證。

為確保多協議設備在 Datalink V2 完整四步驟中均具備一致且正確的端到端體驗，需全面補足步驟 2、3、4 的跨協議規格與實作。

## What Changes

- **步驟 2 (來源規劃與點位衍生)**：
  - 依設備協議（Modbus -> `40001`；FATEK/MC3E -> `D0`）動態給予預設起始位址與 UI 佔位符。
  - 重構 `derivePoints()` 支援全協議：保留字母區域前綴（`D`, `M`, `W`, `X`, `Y`, `B`），支援三菱 Hex 位址（`X`, `Y`, `B` 走 16 進位步進）與 Word/Bit 寬度計算。
  - 非 Modbus 點位功能碼（`function`）調整為協議感知區域標籤（如 `Word (D)` / `Bit (M)`），不硬套 Modbus 的 `holding_register` / `coil`。
- **步驟 3 (標籤映射與即時預覽)**：
  - 更新 `mappingDefaults.ts` 與 `buildDefaultTagKey`，確保產生語意明確的標籤名（例如 `dev.rd0`、`dev.rm100`）。
  - 適配 `useStep3LiveValues` 與即時數值預覽，確保非 Modbus 點位位址可正確被後端 Protocol Connector 辨識並讀取。
- **步驟 4 (資料庫目標與 Modbus Share 輸出交付)**：
  - 更新 `autoAssignTargets.ts` 使其具備對非 Modbus 標籤後綴的智慧匹配能力。
  - 將 rule-level Share 設定納入 source rule create/update/list/get 的 durable configuration，使 `/studio/v2` 重新載入後仍可恢復設定。
  - 在 workspace activation 前，將已持久化 tag 投影為本機 Modbus Slave runtime mappings（`40001+`）；同步失敗、全域 Share 停用或 mapping ownership 無法證明時一律 fail-closed。

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `datalink-workbench-v2-step2-rule`: 擴充來源規則規劃、點位衍生演算法與暫存器語意，使其支援 Modbus、FATEK 與 MC Protocol 之完整暫存器位址格式（含 Hex/Dec 與 Word/Bit）。
- `datalink-workbench-v2-step3-mapping`: 擴充標籤生成與映射預覽邏輯，支援非 Modbus 協議點位之語意化 Tag 鍵名與即時預覽。
- `datalink-workbench-v2-step4-database`: 擴充輸出目標配置與 Modbus Share 跨協議轉換，支援非 Modbus 點位對資料庫欄位與 Modbus Slave 暫存器的自動綁定與激活。

## Impact

- Affected specs:
  - `datalink-workbench-v2-step2-rule`
  - `datalink-workbench-v2-step3-mapping`
  - `datalink-workbench-v2-step4-database`
- Affected code:
  - Frontend `/studio/v2` Step 2–4 state、UI、services、types、i18n 與對應 `frontend/tests/`。
  - Backend source-rule model、validation、repository、service、handlers、schema migration 與對應 Go tests。
  - Runtime Modbus Share mapping API 的既有介面只作 activation projection；不重構底層 connector 或 protocol transport。
