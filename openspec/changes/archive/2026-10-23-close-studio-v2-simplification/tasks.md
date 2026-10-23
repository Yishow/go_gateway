## 1. 盤點與必要修補

- [x] 1.1 以 operation 身分比較後清除 recovery marker：先在 groupSchemaPanel 測試重現舊 failed/succeeded GET 清除新 unknown fence，再修 Durable schema operation identity and honest recovery；memory/storage 分支均驗證 compare-and-clear，保留 known failed 可重新 preview、unknown 查原操作。以 focused Vitest 先失敗後通過為證。
- [x] 1.2 核對 Point grid with modifier-key interactions、Selection is separate from acquisition、Inline-editable mapping table、Bulk apply transform parameters、Target type quick actions、Four-step intent-led setup 及 Narrow basic managed membership updates 的既有實作與 source-hash 證據；逐項列測試與未執行限制，不把已有實作當本輪新增，以對照表與來源內容 review 驗證。
- [x] 1.3 核對 Complete basic preparation with one effective configuration、Managed preparation preserves operation and table ownership、Confirmed additive extension of owned managed tables、Recording claims follow committed evidence、Runtime dashboard renders only backend-backed data、Traceable scoped delivery head resolution 的既有實作與 SQL/unknown/Retry/Skip 見證；修 F 入口歷史狀態與 archive 連結，以 git 歷史、source hashes、對照表及連結檢查驗證。

- [x] 1.4 補齊 Inline-editable mapping table 的 enabled focus 不改批次來源：先以真 focus 事件重現，再阻止冒泡；補 Target type quick actions 的 type-only review/confirm，斷言其餘欄位不變，以 focused Vitest RED/GREEN 與既有行為 mutation check 驗證。
- [x] 1.5 補 Confirmed additive extension of owned managed tables 的 lost additive apply reply 直接回歸：用既有 extension fixture 丟棄已執行的 finish reply，查原 operation 收斂且 DDL 恰一次；以 Go focused regression／mutation check 通過、恢復 production 變更為證，不宣稱新增 PostgreSQL 真 DB 驗收。

- [x] 1.6 補 Durable schema operation identity and honest recovery 的部分 storage failure：當 get/remove 可用但 set 因 quota 失敗，較新的 memory marker 仍優先於舊 storage ID；以 helper 新測試先 RED、最小修復後 GREEN，避免過期回覆解除 unknown 保護。

## 2. 整合與規格收尾

- [x] 2.1 [after: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6] 以既有證據與精準回歸完成規格對齊，連同 Automatic mapping initialization 與 Transform pipeline preview 的舊 demo defaults／seed 要求對齊已存在的來源型別與真 sample：完整前端 lint、unit tests、typecheck/build、受影響的 schema lost-response 真 UI 流程與 change review，核對無舊測試證據失效；analyze/validate 無 Critical，結果明列 passed/failed/未執行，不補造真人或 Windows 驗收。
- [x] 2.2 [after: 2.1] 確認本 change 的 archive preview delta/conflicts、必要驗證與可提交檔案完整，準備封存／提交交付；以 CLI preview、diff/line gate 與原 12 個 dirty 檔 hash 驗證安全性。tasks 真實完成後才執行 archive、核對 canonical specs、路徑限定 commit/push main 並回報遠端 SHA 與 CI，作為 change 外層交付步驟。
