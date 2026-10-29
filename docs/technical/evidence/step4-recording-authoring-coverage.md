# Step 4 authoring：scenario 與驗證層級對照（2026-10-28）

本頁對照 [本案 spec](../../../openspec/changes/fix-step4-recording-authoring/specs/step4-recording-authoring/spec.md) 的實際 15 個 `Scenario`（含瀏覽器條件式例外），不代表發布、歸檔或現場驗收。操作說明見 [迴路記錄](../studio-v2-circuit-recording.md) 與 [群組試寫](../studio-v2-write-groups.md)；[既有 lifecycle 指引](../studio-v2-recording-lifecycle.md) 的規則未改寫。

## 讀取結果前先分清證據層級

- 本次 bounded current frontend regression：20 檔、277 tests 全部 PASS（[最終結果](../../../.task-evidence/authoring-4.3/regression-final.log)）。初跑 275 PASS／2 FAIL 與隔離重現保留在 [原始結果](../../../.task-evidence/authoring-4.3/regression.log)：新增 server recovery GET 消耗舊測試的全域一次性 capability mock；修正為精確 URL 路由後，43 項 capability 測試及完整 bounded union 通過，沒有修改 production guard 或清除復原狀態。這是 harness regression，不是產品行為 RED。早期 1.1–2.3 的 RED／GREEN 為歷史實作證據，不用舊 hash 成功宣稱目前共享檔案全部通過。
- 4.1 的真 backend 驗證走 public authoring／schema／Apply／production pipeline，讀回可丟棄 SQLite：numeric state 1→2 得 integer 2，gauge 得 decimal 1.5；不是前端算值或 fake SQL。見 [numeric-result](../../../.task-evidence/authoring-4.1/numeric-result.json)。它不代表 PostgreSQL、MySQL／MariaDB、SQL Server、真設備或正式 DB 本輪驗收。
- 4.2 正式 browser cases 共 54 個已 discovery；62 個 unit viability tests PASS，是 jsdom 中真 Step4 配 axios boundary simulator，不是真 browser 或 production HTTP／SQL。E2E-specific 與 frontend typecheck、scoped lint PASS；[準備結果](../../../.task-evidence/authoring-4.2-browser/preparation-results.json) 分列每個命令與證據層級。
- 54 個實際 browser cases 及全部截圖 **NOT RUN**，沒有以 discovery、DOM、合成圖片或歷史截圖冒充成功。[逐案例結果](../../../.task-evidence/authoring-4.2-browser/case-results.json) 與 [重跑矩陣／blocker](../../../.task-evidence/authoring-4.2-browser/acceptance-matrix.md) 保留條件式例外、實際阻擋、命令與重跑入口。該例外只適用 browser／截圖，不免除產品、API／DB、frozen-proof 或其餘檢查。

## 全部 scenario 對照

以下前端路徑以 `frontend/tests/unit/` 為根，Go 路徑以 repo 為根；「current PASS」僅指本輪實際選到的 bounded suite。Browser 欄不能由 unit 結果升格。

| Spec scenario（逐字標題） | 正式測試／證據 | 結果與界線 |
| --- | --- | --- |
| Integer machine state is not averaged | `workbench-v2/circuit-profile-authoring-numeric-state.test.tsx`；`internal/api/router_studio_v2_numeric_state_acceptance_test.go` | current frontend PASS；4.1 SQLite latest 2／mean 1.5 PASS；browser NOT RUN |
| Saved role stays locked without a role-changing review capability | `workbench-v2/circuit-profile-authoring-saved-role-integration.test.tsx`、`saved-recording-host.test.tsx`、`saved-recording-actions.test.tsx`；4.1 frozen guards | current PASS；saved role／binding 不變，layout review 不是 role mutation；browser NOT RUN |
| Gauge and counter remain distinct | `workbench-v2/circuit-profile-authoring-numeric-state.test.tsx` | current PASS；payload 明確 role，counter confirmation 不當 gauge；browser NOT RUN |
| Instantaneous-only counter needs a name | `workbench-v2/circuit-profile-authoring-counter-name.test.tsx` | current PASS；缺名稱可見 prerequisite／阻擋；browser NOT RUN |
| Ordinary instantaneous recording has no forced name | `workbench-v2/circuit-profile-authoring-counter-name.test.tsx` | current PASS；一般非 counter 名稱選填、stable ID；browser NOT RUN |
| Enable cumulative with mixed candidates | `workbench-v2/circuit-profile-authoring-cumulative-selection.test.tsx` | current PASS；新累積空選擇與 mixed source；browser NOT RUN |
| Reopen a saved cumulative group | `workbench-v2/circuit-profile-authoring-cumulative-selection.test.tsx` | current PASS；既存來源／關閉重開／scope 保留；browser NOT RUN |
| Instantaneous-only lookup fails | `workbench-v2/circuit-profile-authoring-counter-retry.test.tsx` | current PASS；錯誤／retry、草稿保留且無 save/start；browser NOT RUN |
| Retry returns after navigation | `workbench-v2/circuit-profile-authoring-counter-retry.test.tsx` | current PASS；舊 scope 晚回覆不覆蓋新 scope；browser NOT RUN |
| Long-only destination | `workbench-v2/circuit-profile-authoring-layout-capability.test.tsx`、`circuit-profile-authoring-batch-layout-capability.test.ts` | current PASS；含 retry refresh／全 batch capability guard；browser NOT RUN |
| Connector revision changes | 同上 layout-capability suites | current PASS；舊 connector revision 不授權新 scope；browser NOT RUN |
| Lost response followed by new preview | `workbench-v2/recording-test-write-recovery.test.tsx`、`recording-test-write-recovery-store.test.ts`、`mysql-test-write-operation-recovery.test.tsx` | current PASS；同 O1 查詢、同群組 preview 阻擋、無自動 resend；browser NOT RUN |
| Scope change or local storage loss | `workbench-v2/recording-test-write-server-recovery.test.tsx`、`recording-test-write-recovery-host.test.tsx`、`recording-test-write-recovery-store.test.ts`；`services/recordingTestWriteRecovery.test.ts`；`internal/api/router_studio_v2_test_write_recovery_test.go` | current frontend PASS；3.2 real API／ledger proof 另列下節；多群組／revision、workspace 隔離、儲存損壞／遺失／刪除及 stale cache 不等於未寫入；browser NOT RUN |
| Write succeeded but cleanup did not | `workbench-v2/recording-test-write-phases.test.tsx`、`recording-test-write-cleanup-repair-contract.test.tsx`、`groupTestWriteCleanupRepair.test.tsx`；`services/studioV2WriteGroupTestWrite.test.ts`、`mysqlGroupTestWriteContract.test.ts` | current frontend PASS；write／readback／cleanup 分開，SQL Server 原 token＋O1 repair 不重插；其他引擎／缺 token 只有 query／guidance；browser NOT RUN |
| Blocked browser execution is reported without a false pass | `frontend/tests/e2e/recording-authoring-{choices,recovery,saved-meaning}.spec.ts`；4.2 case-results／acceptance-matrix | preparation PASS（54 discovery、62 unit viability），實際 browser／screenshots 全部 NOT RUN；保留阻擋及重跑命令 |

## 復原與後端證據的窄範圍

- [3.2 backend verification](../../../.task-evidence/authoring-3.2-backend/verification.json)：19 個 top-level tests／77 pass events，含真 router workspace 隔離、ledger 歷史綁定、corrupt／incomplete／超量失敗封閉，與真 SQL／memory cleanup claim shape。這是本地 repository／API 測試，不是外部 SQL Server acceptance。
- [3.3 verification](../../../.task-evidence/authoring-3.3/result.json)：5 個既有 cleanup-repair harness tests PASS，驗原 operation、無新 insert、stale scope／fresh preview／malformed proof 拒絕；本輪 UI current suite 再驗同 O1、lost／mismatched ACK、重複點擊與 scope fencing。沒有提供 generic 或 tokenless repair。
- `written_unverified` 是已提交而未驗證讀回；terminal partial＋cleaned 可解除 pending，不能隱藏 partial history。running cleanup claim 保留「先前」證據，不能讀成目前清理成功。
- Listing 是 observation，不是 mutation permission；本機參照不保存 token。空列表或本機 entry 消失不能消除已知在途請求，未知寫入不自動重送。清理僅針對原試寫擁有權，不刪歷史表／正式資料，不修改 lifecycle successor 語意。

## 最終檢查與未關閉限制

- [4.1 root gates](../../../.task-evidence/authoring-4.1/root-gates.json) 記錄當時相符 source hashes 的 frontend typecheck／build 與 API／Swagger 檢查；4.2 更新後另有當時全 frontend typecheck。這是可追溯的分階段結果，不是本頁另跑全 repo。
- 本次 docs closeout 依 [testing 指引](../../testing.md) 做 diff／line／連結檢查及 Spectra analyze／validate；不將 docs 修改當成需要重跑所有歷史完整 suite。最終發布批次尚未請求，不宣稱已執行完整發布檢查或完成發布／歸檔。
- Capability fixture 修復保留全部原 guard／call-count assertions；只有第二次 capability GET 被延遲，workspace recovery list 獨立回 complete 空陣列，非預期路由拒絕。原先第二項失敗是 `client.isFetching()` 未歸零，不是曾觀察到非法 apply dispatch；兩項失敗都沒有以 waiver 豁免。
- SQL Server successor 真實 DB、其餘本輪未執行外部 DB／平台／設備／LAN／正式環境維持 NOT RUN；先前歷史支援矩陣與驗收紀錄保留，但不當成本輪 proof。
