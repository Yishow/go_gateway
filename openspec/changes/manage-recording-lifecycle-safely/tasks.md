## 1. 契約與相依證據

- [x] 1.1 先建立 `Server classified recording actions` 與「操作分類與停止移除」的失敗回歸測試，證明 draft/applied、stale connector、unknown owner 的 action/reason 不混淆；再以最小改動加入 lifecycle read model，通過 targeted workspace/API tests 後重構，維持現有 mutation guards。
- [x] 1.2 [after: 1.1] 固定「能力限定的轉換矩陣」與 `Explicit supported transition matrix` 的拒絕案例，逐項測試 SQLite／PostgreSQL plain wide↔long、MySQL／MariaDB／SQL Server long-only、跨引擎／mode／kind／不等價型別拒絕；先記錄測試失敗，完成 pure validator 後通過，不啟用 live transition。

## 2. 安全生命週期與世代

- [x] 2.1 [after: 1.2] 實現 `Stop and remove preserve accepted data`：先新增 connector 已修改仍可 stop、ownership proof 失敗不 tombstone、pending 不清除的測試，再調整 typed action errors；用 workspace/grouppipeline tests 比對 accepted payload、versions、receipts 未變。
- [x] 2.2 [after: 2.1] 以 expand migration 建立「後繼世代與唯一擁有權」及 `Stable recording lineage with distinct successor generations`；先測既有 key 初始化、deleted key 不被普通 ensure 復活，再加 lineage/generation/reservation metadata，migration round-trip 不刪原群組／外部 SQL 資料，測試確認新 ID 與舊 receipt 不碰撞。
- [x] 2.3 [after: 2.2] 為 recreate reserved draft 寫 competing CAS／相同 operation 重送失敗測試，再實作同 lineage 的唯一 reservation 與不同 intent conflict；通過並重構，確認普通 create/apply 沒有放寬其他 owner 衝突。
- [x] 2.4 [after: 2.2] 交付 `Historical destinations and pending units remain explicit`：review/read model 分開 samples、outbox、blocked/unknown 與各原 revision/destination；先建立混合單位測試，再實作，核對回傳不相加成一個「筆數」且不洩漏 credential。

## 3. 審核切換與復原

- [x] 3.1 [after: 2.3, 2.4] 實現「審核邊界與復原」的 successor-review／confirmation 契約與 schema preparation，不 ready 時保持 predecessor；先以 API tests 覆蓋 digest/revisions/capability stale 拒絕，再實作並驗證無確認即無 DDL／intake mutation。
- [x] 3.2 [after: 3.1] 在記錄 `freeze-counter-replay-semantics` 與 `fix-recording-cutover-boundaries` 相關驗收通過後，交付 `Reserved preparation and atomic effective ownership` 的 SQLite plain 雙向切換；先測 observed_at 邊界、事先已接受跨界觀測的 audited handoff、late/unfinalized 與 delayed reconcile、new counter baseline/NULL delta，再實作，disposable SQLite SQL readback 證明舊 finalized effects 不改、無雙寫。
- [x] 3.3 [after: 3.2] 將相同契約擴至 PostgreSQL plain 雙向切換及同 layout/mode 目的地換接；先新增失敗案例，再實作並通過 disposable PostgreSQL tests，核對新表 generation、舊表／receipts 保留與不支援 shape 拒絕。
- [x] 3.4 [after: 3.3] 驗證 MySQL 與 MariaDB 的同 layout/mode 長表目的地換接，分別建立失敗 capability／receipt／owner 案例後修正 adapters，兩種 disposable DB 分別記錄通過證據，wide 始終拒絕且舊歷史不動。
- [x] 3.5 [after: 3.4] 完成 SQL Server 同 layout/mode 長表換接的本地 adapter／完整 proof／guarded reload，先測 native/plain 證據不得跨 mode 復用及 wide 拒絕，通過 targeted normal／race／vet 與限定審查；備妥 disposable readback tests。依 2026-10-28 使用者「若無法就跳過」的同機驗收例外，實機不可行時保留 NOT RUN 與風險，不阻擋後續，不冒稱真實 SQL 驗收通過。
- [x] 3.6 [after: 3.5] 實現 `Successor operations survive interrupted navigation`：先測 lost acknowledgement、restart、同 request replay、不同 intent 拒絕與新 forward rollback；再接 ledger readback，通過後重構，確認 unresolved operation 不會被另一預覽或最新 operation 取代。

## 4. 操作入口與最終驗收

- [x] 4.1 [after: 3.6] 將 Basic/Advanced 的 saved recording 操作接到同一 action availability，以 component/integration tests 驗證「停止新記錄／移除設定」文案、保留歷史說明、blocked reason 與 recreate/transition review；不改 guided layout，提供後續 `fix-step4-recording-authoring` 使用的穩定 callback/action interface，驗證既有 caller guards 不回退。
- [x] 4.2 [after: 4.1] 補齊 route/Swagger/type 契約及操作文件，以 generated Swagger 檢查、相關 Go tests/vet、前端 targeted tests、typecheck、lint/build、git diff --check 與 make check-lines 驗證；所有測試使用 simulator／可丟棄 DB，不授權正式資料寫入。
- [x] 4.3 [after: 4.2] 完成 `Delivery acceptance evidence remains explicit` 與「瀏覽器同機驗收的授權例外」：備妥正式 browser tests／實際 host fixture，涵蓋 draft/applied/disabled/deleted、Basic/Advanced 互切、鍵盤、Cancel/Back、重複點擊、reload、未知結果與歷史；核對 discovery 42 案例、每個 spec scenario 的結果、完整前端測試／lint／typecheck／build 及限定 review。依 2026-10-28 使用者「跳過並commit到github main後繼續」，實際 browser／截圖保留 SKIP／NOT RUN 與重跑方式，揭露既有未修復啟動競態後可完成本次交付；不得冒稱通過，也不得略過兩項 backend 前置證據或放寬 live transition guards。

驗收例外紀錄（2026-10-28）：原 pending task 3.5 的「完成 disposable SQL Server readback；無環境不視為完成」門檻，依使用者「同意，若無法就跳過」改為本地實作完成、實機測試 NOT RUN；見 design「SQL Server 同機驗收的授權例外」。原 [x] 任務、產品矩陣與歷史證據不變。

瀏覽器驗收例外紀錄（2026-10-28 15:27 UTC）：原 pending task 4.3 的實際 browser／截圖必跑門檻依獨立明示授權取代；保留正式來源與逐案例 NOT RUN，完成本地檢查後可承接後續 change。本項不是沿用 SQL Server 授權，亦非宣稱未執行的畫面驗收已完成。原 14 項 [x] 描述及 provenance 不變，普通啟動既有競態仍未修復。
