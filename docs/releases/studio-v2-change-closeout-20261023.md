# Studio V2 change 完整性盤點與收尾（2026-10-23）

基準 `b8c55aabfa59a89bd0f5206fdfc3b4730b12d006`（已推送 main 的前端簡化）；本輪唯一 owner 是 `close-studio-v2-simplification`。原工作樹的 12 個未提交／新增檔案隔離保留，不納入本次變更。這是既有 Studio V2 的規格收尾及實際 bug 修補，沒有另建前端。

## 盤點結論與順序

- 本輪開始時 `spectra list --json` 只有 `add-windows-tray-runtime-logs` active（9/20），parked 為空。沒有尚待 apply 的 Studio 新手簡化 change，也沒有可擅自合併的 parked 工作。
- 10/02 四步簡化、10/04 A–F、10/07 基本準備／mapping concurrency／value integrity 等相關 changes 已 archive。archive 或 tasks 勾選不是實作證據；本輪另核對 source、assertions、SQL witness 與 Git 歷史。
- `b8c55aab` 九個實作點已落地，但沒有專責 OpenSpec delta。正式規格仍要求點擊改 skip、一鍵套所有 mapping，以及受管表僅 create/no-op；這些是實際規格落差。
- 獨立 review 找到實作缺陷：同群組舊 schema operation 的延遲 terminal 回覆，會無條件清除較新 unknown marker。先測試重現，再修復比對 operation_id 後清除。
- 接著核對並同步六個 capabilities、修正 F 入口歷史文字／連結；重驗受影響前端與真 UI lost-response 流程，最後 review、驗證 archive preview，完成後才封存與提交。

## 已有實作與本輪 delta 對照

下表是核對既有行為，不把 `b8c55aab` 的工作列為本輪新寫的功能。單元測試名稱／檔案均可在 `frontend/tests/unit/` 查找；真 UI 證據見下節。

| 合約 | 實作 | 覆蓋的可觀察情境／測試 |
| --- | --- | --- |
| Point grid with modifier-key interactions；Selection is separate from acquisition | `PointGrid.tsx`、`Step2Rule.tsx` | `step2-grid.test.tsx`：普通點擊、Ctrl、Shift 只選取；明確 checkbox／批次採集；切換 rule 重置 |
| Inline-editable mapping table | `MappingTable.tsx`、`MappingRow.tsx` | `step3-components.test.tsx`：五欄、Advanced 同 editor nodes 保值、編輯來源與 checkbox 不混淆；`mappingDefaults.test.ts`：來源型別與中立 defaults |
| Automatic mapping initialization；Transform pipeline preview | `mappingDefaults.ts`、`useStep3LivePreview.ts` | `mappingDefaults.test.ts`、`reducer-step3.test.ts`、`step3-live-preview.test.tsx`／`step3-live-preview-changes.test.tsx`：來源身分與型別、非中立 numeric 預設 float64、existing/orphan、真 sample、無值等待、loading/error/late response；舊 demo seed 規格同步移除 |
| Bulk apply transform parameters；Target type quick actions | `MappingBulkApply.tsx`、`mappingReducer.ts` | `step3-components.test.tsx`、`bulk-target-safety.test.ts`：來源、勾選、全選 scope、單欄／全部四欄差異、取消、值變更失效、disabled／不存在／來源排除；保留單列 target select 與 type-only 批次 |
| Four-step intent-led setup | `StepRail.tsx`、`RuleEditor.tsx`、既有 workspace/readiness | `shell.test.tsx`、`plain-range-summary.test.tsx`：四步名稱、明確修改建議名、4 float32 = 8 words；`step1.test.tsx`／`sourceRule-readiness.test.ts`：不支援協定及未驗證修復；`address-conflicts.test.ts`：兩設備同 40001 不衝突 |
| 既有四步恢復／基本 typed group | workspace hydration、`basicRecordingState.ts`、既有 ensure service | `workspace-boot.test.tsx`、`basic-point-tag-to-group.test.ts`、`basic-recording-state.test.ts`：離線草稿／reload、未指定 measurement、custom/multi-entity/cross-device 不轉換；Go `write_group_basic_managed_test.go`：重啟重試、並行收斂、destination scope CAS、tombstone |
| Complete basic preparation with one effective configuration | `Step4Database.tsx`、`CollapsibleSupportCard.tsx`、`BasicSchemaPreparation.tsx` | `step4-database-controls.test.tsx`、`step4-database-components.test.tsx`、`basic-recording-preparation.test.tsx`：DB 優先、摘要／修改、收合保值、基本卡準備與唯讀 preview、readonly／scope 失效；`basic-recording-panel-regressions.test.tsx`：Share-only 無 DB ensure |
| Narrow basic managed membership updates | `BasicMemberUpdate.tsx`、`basicRecordingMembers.ts` | `basic-member-change.test.ts`、`basic-member-update.test.tsx`：新增／移除 diff、取消、連點 single flight、CAS失效、readonly、custom/entity 阻擋、欄位/required/age/policy 保留；release SQL witness 新增點位及歷史列 |
| Managed preparation preserves operation and table ownership；Confirmed additive extension of owned managed tables | `schema_group_extend.go`、`schema_group.go`、`schema_apply.go` | `schema_group_extend_test.go`：readonly、confirm/replay、owner/required/foreign/type/history 拒絕、stale exact statement、歷史 nullable 欄保留、PG NUMERIC(20,0)；`schema_group_test.go`：token/layout/ownership、expiry、lost ack、不重送 |
| Durable schema operation identity and honest recovery | `GroupSchemaPanel.tsx`、`schemaOperationRecovery.ts`、既有 backend ledger | `groupSchemaPanel.test.tsx`／`schemaOperationRecovery.test.ts`：unknown 查原 operation、reload、多卡、scope 變更、known failed 可重預覽；本輪新增 stale terminal 不清 newer marker 的競態回歸及 storage fallback |
| Recording claims follow committed evidence；Runtime dashboard renders only backend-backed data | `RuntimeDiagnosticsPanel.tsx`、`BasicRecordingEvidence.tsx` | `diagnostics-no-write-claim.test.tsx`、`runtime-dashboard-truth-state.test.tsx`、`basic-recording-evidence-regressions.test.tsx`：missing data 不補健康值、沒有 failure 只表示未偵測錯誤、Share/採集不等於 SQL；成功需 group/applied/connector committed effect |
| Traceable scoped delivery head resolution | `GroupDeliveryRecovery.tsx`、既有 delivery resolution API | `groupDeliveryRecovery.test.tsx`：Retry 主動作、Skip 展開影響/原因/勾選、unknown 禁止、uncertain decision_id replay、新 blocked attempt；release witness Retry SQL=1、Skip SQL=0、同 payload 與 audit |

Delta 修改保留原 requirement 的有效 scenarios；本輪替換的是與授權行為衝突的舊 click/bulk/呈現限制，沒有削弱型別、資料庫、DDL／token／CAS／ledger 或 Retry/Skip 保護。來源型別數值精度不以 JavaScript Number 模擬成 SQL 成功。

## 證據有效性與本輪驗證

`b8c55aab` 的完整命令、1331 項 Vitest、14 項 Playwright、62 項 harness、Go test/vet/lint、雙設備三桶與 release SQL 證據在 [簡化驗證](studio-v2-simplification-20261023.md)。本輪核對 [來源指紋](../plans/studio-v2-flow-completion/evidence-f/simplify-source-inputs-20261023.json) 與 Git diff：除下列實際修補及補測之外，既有簡化 source/tests/config 沒變；Go production 與 lockfile 沒變。雙設備三個 60 秒桶保留為 passed-prior，不宣稱本輪重跑；Go gates 本輪重新執行。

本輪實際執行結果見下節；每項命令與證據 SHA256 記錄於 verification manifest。沒有真人新手、現場 PLC、正式 DB、Windows/Linux 原生或部署驗收。PostgreSQL 安全加欄位有 Go contract assertion，本輪未執行 PostgreSQL 瀏覽器／真 SQL，不能以 SQLite 代替。

## 文件狀態與歷史限制

F 實際在 `d41638b0` 完成並封存；計畫 [README](../plans/studio-v2-flow-completion/README.md) 的「待 verify/archive／尚未提交／沒有 push」是過期入口文字，已依歷史與本輪現況分開並改到 archive 連結。沒有為修文字重跑 F 或補造外部驗收。

其他歷史 archive 的未勾項需分開看：workspace foundation 的 empty-state 已有 `workspace-boot.test.tsx` 及實作，不等於缺功能；歷史 RED 缺證據不能補造；legacy retirement 的 deployment owner rebuild/redeploy 是外部交付。它們不以本輪勾選補完。`studio-surface-inventory` 的 2026-08-25 多設備未接線說法是過期背景，本輪揭露但不擴編修改該目錄或 changelog DB。

## 仍 active 的 Windows tray/runtime logs（不同範圍）

保留原 `add-windows-tray-runtime-logs` 的 11 個未完成項與使用者 dirty 檔，不封存。source 存在不是驗收通過；Studio E2E 也不是 logs 專屬驗收。

| Tasks | Mac 可補的部分 | 仍需 Windows 或其他證據 |
| --- | --- | --- |
| 1.1、1.2 | resolver／mode／lifecycle／bind-first 自動契約、macOS CLI | Windows path/權限與 desktop 啟動契約 |
| 2.2 | owned file-sink、輪替、stall/disk error 的隔離 fixture | native reparse／ACL／disk fault；歷史失敗需在該 change 最終 source 重現，不能因 Studio Go PASS 就勾選 |
| 3.3 | 本輪完整 Vitest 與 14 項 Playwright 已含 logs unit 及 embedded logs smoke；仍可補 logs 逐 scenario 矩陣 | smoke 未取代全部 XSS／reconnect／keyboard／容量專屬見證，且原 change 有使用者 dirty 內容，本輪不改其 tasks |
| 4.1、4.2、4.3 | 可攜 lifecycle/shutdown unit/race、accepted backlog/unknown 契約 | Windows OS owner/ACL IPC、tray、Explorer restart、無 console 閃窗、GUI 退出與 deadline |
| 5.1、5.2 | build-script contracts、完整本機 gates（需該 change 最終內容） | PowerShell 原生建置、GUI PE／console artifact、Windows 行為與跨平台證據 |
| 5.3、5.4 | 整理已知限制與交接材料 | Windows 互動桌面完整矩陣、無 Node/Go portable exe、redirect/wait/exit、同 artifact logs+durable recovery；最後逐 scenario 文件收尾 |

現有 evidence 把多項列為「已有實作但缺完整驗收」，唯讀盤點沒有足以新定案的純功能未實作缺陷；仍須上述目標平台執行判定，不能自動勾選或宣稱完成。使用者是否要另外完成這個 Windows 產品擴充，交由主對話安排環境與範圍。


## 本輪實際結果（最終程式）

- 修復三個可重現邊界：延遲 schema terminal 只清相同 operation；storage quota 寫失敗時較新的 memory fence 優先；enabled checkbox focus 不改 bulk 來源。
- 原操作競態測試 29/29 通過（含 quota 新案例）；mapping focused 10/10 通過。RED 分別重現舊 marker 被清、quota 丟失 B、focus 選到錯列；type-only 與 additive recovery 各做一次 production mutation，測試均抓到，立即還原。Go runtime production diff 為空。
- `npm run lint` PASS；`npm test -- --run` **211 檔／1342 項 PASS**；`npm run build`（tsc＋Vite）PASS，保留既有 chunk size 提示。
- `go test -p 1 ./...` PASS；`go vet ./...` PASS；`golangci-lint run ./...` **0 issues**。執行環境仍是 Mac Studio arm64／Go 1.27.3（repo 目標 1.25.5）。
- `npm run test:e2e` **14 項 PASS**，包含 embedded logs smoke；未因此完成 Windows change。
- 最終 production embedded UI＋loopback Modbus simulator＋隔離 SQLite，`simplify-flow.mjs` PASS：[原始 witness](../plans/studio-v2-flow-completion/evidence-f/closeout-release-20261023.json)。真 DDL 回覆遺失後 refresh 只查原 operation，apply=1；int64 最大／最小值精確查回；2→3 成員與 confirmed add_columns，舊兩列保留、新值777；Retry SQL=1、Skip SQL=0，payload／audit identity 保留。page errors 空、自有程序 exit0、DB namespace 清除。
- [映射畫面](../plans/studio-v2-flow-completion/evidence-f/closeout-release-20261023-mapping.png)、[記錄畫面](../plans/studio-v2-flow-completion/evidence-f/closeout-release-20261023.png) 已檢視。這是自動化與 source review，不是任用真人新手驗收。
- Code review 四面向沒有新增 findings，快照 `9ce7fef6b3f233d6cce2392337caabeecca9b3bc` 已驗 current；限7個程式／測試檔，部分 schema 為原作者自行 review，未宣稱全範圍獨立審查。規格與文件由主代理另核對。
- `spectra analyze` 無 Critical／Warning（40 個 concrete-example 建議為 advisory）；`spectra validate` PASS。`git diff --check`／line gate PASS，329／416／407 行檔案是 warning，無超過500新檔。canonical specs 已由單次 archive 套用，不手改 trace；另逐一比對未涉及的 requirement 內容保持原樣。

[Verification manifest](../plans/studio-v2-flow-completion/evidence-f/closeout-verification-20261023.json) 保存當次命令、結果、source/log/witness 指紋。Git commit／遠端與 CI 以提交後查證結果回報，尚未提交時不預寫 SHA 或成功狀態。

封存結果：[`2026-10-23-close-studio-v2-simplification`](../../openspec/changes/archive/2026-10-23-close-studio-v2-simplification/proposal.md)，8/8 tasks 有證據完成；六份 delta 套用一次，4 added／12 modified；49 個未涉及 requirement 逐一比對未變，正式 66 份 specs 全數 validate 通過。沒有 removed／renamed、衝突或 cleanup warnings。Windows change 保持 active。
