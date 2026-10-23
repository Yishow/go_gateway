## Context

主線基準 `b8c55aabfa59a89bd0f5206fdfc3b4730b12d006`。相關 10/02、10/04、10/07 changes 均已 archive，這不是新簡化已納入規格的證據。唯一其他 active change 是 Windows tray/runtime logs，與本輪不同且原工作樹有使用者修改。

## Goals / Non-Goals

**Goals:** 修正實際 recovery 競態，補足授權簡化合約及可追溯驗證、入口文件與 change 收尾。

**Non-Goals:** 不重造前端、不重開歷史已完成工作、不改 Windows change、原 12 個 dirty 檔、inventory DB 或系統名稱。不碰現場 PLC／正式 DB／部署。自動化不代替真人新手或 Windows 驗收。

## Decisions

### 以 operation 身分比較後清除 recovery marker

保留現有 workspace/group marker 介面，清除時要求 expected operation_id，僅當目前 marker 與回覆一致才刪除。舊回覆不影響新操作。替代方案是無條件群組清除，已確認存在競態，故拒絕。

### 以既有證據與精準回歸完成規格對齊

使用 `docs/releases/studio-v2-simplification-20261023.md`、source-inputs/summary 與真 UI→SQL 原始見證核對九個實作點。新競態用 deferred GET 與第二卡 lost reply 重現。完整前端 lint/tests/build 重跑；Go 沒變則沿用相同 source 的完整 Go gates，無必要不重跑三個 60 秒桶。

## Implementation Contract

- `retainSchemaOperation` 的清除流程新增 expected operation identity；寫入仍保存 server-issued operation_id、當下 scope 與未解決狀態。terminal result 只釋放自己的 fence。
- MappingRow enabled checkbox 的 focus 與 click 均不冒泡成來源選取；type-only 批次沿同一 review/confirm，不變更其他欄位。補 additive lost-ack 直接測試，production Go 不變。
- 遲到 failed/succeeded GET、known apply terminal 回覆與其他已開卡不得清新 marker。storage 不可用時 in-memory map 同樣遵守比較。
- unknown/partial/pending/running 維持原查詢操作與禁止新 preview/apply；known failed 在自己的 fence 清除後仍可新 preview。scope/CAS/token/backend ledger 不變。
- delta 必須保留已有效的來源型別精度、所有協定能力、Share-only、DB選項、自訂群組、欄位/required/entity/provenance/去重/歷史資料與版本保護。
- 驗收：新增測試先因 marker 被錯清而失敗，再修 helper 與所有 terminal callsites 後通過；原 schema panel tests、完整前端 gates、Spectra analyze/validate/review 及路徑限定提交。

## Risks / Trade-offs

- 舊文件可能誤導重做已完成 F → 以 commit 與證據日期分開歷史／現況，修 archive 連結，不補造外部結果。
- 既有 specs 有大量 trace → delta 僅覆蓋完整指定 requirement，由 CLI archive 合併；不手改 base spec 或 trace。
- 並行卡競態測試可能只測單卡 → 明確要求兩卡、第三卡 fence 見證，另覆蓋 storage fallback。
