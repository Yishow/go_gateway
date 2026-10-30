## Why

同一組記錄需要兩點取平均、一點算增加量，可各用15秒／1分鐘／10分鐘等對齊時間窗；相同到期時可同列，不同到期時寬表採到期邊界聯集、未到期NULL＋not_due，長表僅產生到期member列。現有group層instantaneous/cumulative模式無法表達這個組合，不能靠改下拉選單或拆成多組假裝滿足。

## What Changes

- 增加recording member自己的latest／mean／delta計算意圖，位置是group member/mapping而非來源Point；同一來源可在不同記錄採不同算法。
- 新明確versioned mixed recording shape以各member週期與既有UTC epoch對齊排程，將各member結果放各目標欄，delta值不是原counter raw冒充。
- 重用現有mean、counter policy／baseline及durable closure，新增混合派送而非第二個採集／計算引擎；依已確認sparse wide／due-only long契約；新版demo仍須owner確認，不得自動apply。
- 保存舊profile/reduction/legacy SQL語意，不因升級／載入改成新shape；明確變更需fresh schema/readiness/ownership/revision確認。
- 缺值、首窗baseline、reset／rollover、重放、partial policy與SQL receipt分層驗收；提供simple editor與adoption的前置API。

## Capabilities

### New Capabilities

- `per-member-recording-calculation`: 每記錄成員計算意圖、混合時間窗／SQL結果與durable安全契約。

### Modified Capabilities

- `window-average-recording`: 限定原group reduction契約繼續適用舊instantaneous/cumulative，新明確mixed shape由逐member契約負責。

## Impact

- internal/datalink/workspace、internal/datalink/schema：member持久化、profile／resultshape版本、CRUD/CAS、migration相容、managed ownership與schema/readiness。
- internal/datalink/snapshot、internal/datalink/runtime、internal/datalink/aggregation、internal/datalink/groupdelivery：同窗逐member派送、counterbaseline與closure/replay、frozenproof。
- internal/datalink/dbtarget：新混合row結果欄、品質／provenance與schema shape；不覆寫舊counter raw列語意。
- internal/api/handlers、frontend/src/types、frontend/src/services：scoped member API與讀回計算／resulttype／revision，頁面呈現由streamline-recording-task-editor承接。
- 不新增環境、協議或DB dialect；尚未實作、測試或提交產品。
