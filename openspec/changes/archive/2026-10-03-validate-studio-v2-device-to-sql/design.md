## Context

前置A/B/C/D/E驗證。已讀proposal與既有scripts/b10_exe_acceptance.py、frontend/tests/e2e的定位；這批要新增database路徑witness，不把既有Modbus gate當成SQL驗收。主程式必須使用與正式入口相同wiring。

詳細來源見 [evidence](../../../../docs/plans/studio-v2-write-groups/evidence.md)，共用欄位與政策見 [跨案契約](../../../../docs/plans/studio-v2-write-groups/contracts.md)。這是擬議設計，不是完成報告。

## Goals / Non-Goals

**Goals:** 可重跑的simulated-device→真UI→persisted IDs→actual SQL row，正反案例都有可核對witness。

**Non-Goals:** 真PLC/SCADA/production DB、全面platform certification、虛構性能SLA或新增產品功能。

## Decisions

1. **隔離harness。** 每run建立獨立workspace、loopback simulator、GATEWAY_DB_PATH及外部SQLite/PostgreSQL；以實際cmd/test_ui binary和embedded frontend跑Playwright。UI/API不能stub最終成功，simulator fixture可控sample/time/failure。PostgreSQL用已授權可丟棄instance，缺依賴則blocked，不能把SQLite結果當Postgres pass。
2. **可追溯鏈。** witness保存source commit/build、平台、command、workspace/device/point/tag/group/revisions、sample/record/operation IDs、query與typed actual SQL結果、timestamps/quality、UI screenshot。祕密與endpoint憑證排除。DB rows由獨立查詢驗證，不只看API response。testwrite cleanup前後留count及neighbor row不變證據。
3. **故障矩陣。** 見acceptance.md，每個scenario對應test label、fixture和pass assertion。含multi-device same address、mixed values、bad/missing/late、outage restart、duplicate/unknown commit、disk full/poison、edited revisions與並行request。clock與fault injection具明確boundary，不能sleep一下就假設commit發生。
4. **範圍與數字。** 以小型deterministic fixture確認正確性；測試的10秒bucket等是design fixtures，不是性能目標。真實吞吐、穩定性duration或field readiness若未測就列not run。不同平台/embedded browser各自有run result，Linux結果不自動升格。
5. **closeout。** 全案包含Go/frontend full checks、line/diff、OpenSpec strict validate與independent review；失敗保存witness與repair方向。修復若超過本批scope另報，不把變更規格當作讓測試通過的捷徑。原active change的handoff只在D/E/F證據齊全後關閉，原完成記錄不重寫。

### SupportedAcquisitionTypesAndDecimalBoundary

**Supersedes**: validate-studio-v2-device-to-sql / tasks 1.2「uint64及decimal精度一致」及acceptance.md「temperature decimal」真採集fixture。

2026-10-03使用者確認：真UI→採集→SQL維持既有int16、uint16、bool、uint64、float32、float64、string。實際witness分別查七型別與uint64 `9007199254740993`，兩device同addresses、persisted IDs、UTC、quality與neighbor不變仍全部必驗。decimal保留B explicit expected type的codec／SQL round-trip證據，明列設備採集未支援；不從Tag推測decimal，不新增採集／mapping／UI型別。未執行的decimal設備路徑被此決定取代，不標成已實作；既有B完成證據與非decimal任務不變。

## Implementation Contract

### QualityBoundaryMatrix 與 deterministic fault control

- 以 `f_write_group_fixture` build tag 建立驗收 binary，保留同一 `runGateway`／`wireGatewayServices`／router／embedded frontend／production group pipeline；一般 build 的接縫為 no-op，沒有 fixture listener 或 typed fact override。
- fixture controller 只 bind 127.0.0.1；在開啟內部 DB／migration 前確認 GATEWAY_DB_PATH 為有 ownership marker 的 gw-f-* 暫存目錄，拒絕未知路徑與 symlink escape；clock 與 acquisition ID 可注入既有 scheduler／pipeline；數值與來源 ID/revisions 只能來自真實 loopback Modbus read 和 runtime mapping。捕捉只是待送出，明確標 captured_not_acked；release 一律呼叫既有 production AcceptSample，ACK／journal／closure／outbox／sender不模擬，不直接 INSERT 測試資料列。bounded capture queue 滿時明確失敗。
- 先用真 UI建立 device/point/tag/group/目的地及 valid Apply。移除測試 process 的 memory polling tickers，保留 persisted 設定，手動 poll 仍經真實 protocol read→runtime mapping→typed sink；因此能在固定 clock 下控制到達順序。control 不提供任意 value 或 persisted source ID/revision 的改寫。
- QualityBoundaryMatrix：10秒 UTC bucket中實際讀到 t=8 的新值再讀 t=3 的舊值，release順序反覆仍選t=8；同一captured sample再release是no-op，重用acquisition identity但由device讀回不同值為conflict；t=10屬下一bucket，closed後release既有captured t=8為late且SQL不變。clock起點為2026-01-01T00:00:00Z（Apply於前一bucket）。
- missing以只poll部分persisted members；bad以simulator disconnect後實際failed read；stale以已存max_age=2秒、observed t=3/end=10；silent以不poll直接tick。default skip_row 不寫任何SQL row，durable bucket 原因與真實 delivery API／UI 的 skipped/no_data 一致。既有 delivery API 新增 bounded 最近20筆 recent_bucket_issues：group_revision／bucket_start／kind／causes，causes只含 missing／bad／stale／invalid／no_data／unavailable，不輸出數值或原始錯誤；UI顯示總桶數及最近原因，讀取失敗時顯示未確認。explicit partial 透過 canonical Save/Apply；沿用 B 的 NewGroupRowLayout 修正 readiness 全擋的整合缺口，只有真實 inspected nullable optional 欄位與有效 provenance 欄位才能接納；actual NULL＋member reason 不變成zero／old value。
- semantic Apply 生效後，即使 reconcile 漏過 lookahead，所有更早 managed boundary 仍以各自 interval 對齊的新 effective_at 截止 intake；future prebuilt boundary 不提前退休，既有 cutoff 只縮短不延長。驗收以同 source revisions、max-age policy 切換且直接跨過生效時刻，檢查舊 revision 不再 journal／SQL 雙寫。這是 C 已核定單一 revision authority 的回歸修正。
- 控制拒絕錯誤clock／foreign或unknown source／無capture index／overflow，不能繼續宣稱PASS；witness保存clock、source IDs、actual Accept結果、journal/bucket/outbox/target查詢及UI截圖。正常binary另驗fixture控制不可用。
- fixture僅驗收process與可丟棄資料庫；不新增產品debug API／UI、不調整production型別或snapshot policy、不替代1.1–1.3主路徑。其他2.x故障仍依原tasks執行，不把既有package測試當未執行production案例通過。

### OutageCrashAndUnknownCommit

- 沿用上述驗收 build tag、loopback controller、owned temporary DB 與真 UI 所建立的 persisted sources。目的地使用實際 SQLite／PostgreSQL driver；禁止用 fake writer result 或直接 INSERT 製造成功。
- controller 增加具名 fault 控制：依真實 group 查出 frozen connector，切換 target commit response lost／target commit hold、group-scoped local receipt failure／closure transaction hold。拒絕不存在的 group 或未知 fault；不接受任意 SQL、value、DSN 或 source revision。一般 binary 不包含這些控制；GET /faults 只讀 process-local snapshot，不等待 pipeline lock 或 SQLite transaction；state 僅回傳 safe fault kind／persisted IDs／actual barrier reached 與 commit count。
- target fault 以 build-tag-only destination connection hook 重用 lostcommit driver wrapper：真正 Commit 成功後才丟失回覆或停在 barrier；保留 production InsertGroupRow、destination receipt／unique key readback 及 sender settle。正常 OpenDestination 的 hook 為 no-op，原 identity revision／enabled 保護維持。SQLite delivery 使用 atomic existing-file open，Stat 後檔案消失不得建立新空檔；同一 normalized DSN 交給實際 manager 與驗收 wrapper，memory／read-only 語意保留，不能全域改 probe／schema 開檔行為。控制設定不改 persisted destination。
- local receipt failure 在 owned local DB 安裝具名且 group-scoped 的固定 trigger；closure hold 在實際 checkpoint 寫入前停住同一 transaction，沒有自行提交 outbox 或 checkpoint。scalar function 僅存在 tagged binary，於首次 internal DB open 前註冊並保留註冊失敗；barrier state 不等於 ACK 或 SQL committed。cleanup 解除所有 owned fault／barrier。
- 先獨立查出 journal 已 ACK 且 target 尚無列，再 SIGKILL／restart／tick，驗證已接受 samples 能形成一列。closure transaction 停住時獨立查 journal 未 consumed、bucket／outbox／checkpoint 尚未前進；SIGKILL 後實際 transaction rollback，restart 後重新形成同一 effect。這兩個時點不能用固定 sleep 猜測。
- SQLite 與 PostgreSQL 各測明確 destination receipt 能力下 lost commit response 與 commit 後 kill，恢復後同 effect key 僅一個 target effect，local receipt 只有真正確認後出現；dedupe=none 時實際 row 可能已存在，但狀態必須 unknown，重啟或重試不得盲寫第二列。local receipt 失敗另測 receipt 收斂與 none 保持 unknown。unique_key 的套件能力不代替 production positive witness；目前 canonical 設定未提供 RecordKeyColumn，這批不新增該產品設定，另明列未驗。
- 實際 target outage 期間 local backlog 留存，另一個獨立 connector 的健康 group 繼續交付；重啟後原 target 恢復，原 frozen payload／effect key 不重算。用 canonical connector GET／PUT 與 expected_identity_revision 修改 endpoint，舊 backlog 明示 target-blocked、新 endpoint 零舊列；不改回 revision 規避保護。UI pending／committed／attention 與獨立 SQL／local durable state 相符。
- witness 保存 build/source/platform/commands、fault/barrier reached、journal／closure／outbox／receipt、effect keys 與獨立 actual SQL／UI 證據。這些是小型 fault fixtures，沒有調整 production retry policy、加入產品 debug surface、部署或證明真 PLC／正式 DB。

### CapacityPoisonAndConcurrentEdit

- quota 控制只存在驗收 build：startup 可指定正整數 `F_FIXTURE_GROUP_QUOTA_BYTES`／`F_FIXTURE_GLOBAL_QUOTA_BYTES`，上界 1 TiB；未指定沿用原 global 500 MiB。錯誤設定於 DB open 前拒絕，一般 binary 忽略這些環境變數。仍由 production Store.WithQuota／AcceptSample 決定 ACK，不提供 runtime quota setter 或删除 backlog 來讓測試過關。
- 先以真實採集形成兩列 outage backlog，獨立查實際 bytes，再重啟套用低 quota。group hard-limit 只拒絕該群組新 ACK，健康群組仍寫 SQL；global hard-limit 拒絕其他群組新 ACK。拒絕前後 journal 相同，delivery API 的 scope／intake_refused／loss_risk_notice 明確，UI SQL committed 與 durable receipt 相符。解除限制後原 accepted payload 能交付。
- tagged `POST /capacity` 僅接受 `{kind:"disk_full",enabled:bool}`，在 owned internal SQLite 單連線 pool 設定 max_page_count 等於實際 page_count；解除恢復原限制，回傳 safe kind／enabled／page_count／max_page_count。不得傳 SQL／任意數值／DSN。以反覆真 Modbus read→AcceptSample 直到實際 SQLITE_FULL，確認拒絕 ACK、journal 不新增、既有 accepted samples 全保留；解除後正確交付。這不是裝置實際磁碟耗盡或硬體可靠度測量。
- permanent poison 以 owned 外部 SQLite 固定且 entity-scoped constraint trigger 拒絕 A row，樣本值仍由真實設備讀回。A 首列 quarantined、successor pending，B 繼續 committed；poison payload 不刪除、不跳過順序。trigger 只建立 schema 控制，禁止直接 INSERT 代替採集。
- concurrent Save 使用兩個真實 HTTP client 的同一 CAS，一方 200、一方 409 revision_mismatch；stale Apply 409。test-write preview 後修改 draft，再 confirm 必須 409 WRITE_GROUP_TEST_WRITE_PREVIEW_STALE 且 target 零新增。未 Apply 的設定不得切換 applied revision／凍結 backlog。UI conflict 用延遲真實 backend GET 的回覆製造時間順序，actual PUT 409 保留 local draft，重新讀取後可以修復 Save；不得構造假的 API 成功。
- dispatcher overlap以兩個實際process共用owned store：primary停在真Commit之後、local receipt之前，tag-only secondary使用白名單foreign test node identity，不接管live claim；另一健康partition須由另一claim owner實際寫入。secondary ports啟動前bind檢查空閒。tag-only startup MaxRetries=2驗證兩次真transient failure後blocked，保留payload與successor、健康scope繼續；normal build忽略test node及retry ENV。production NodeID／fencing／retry default不變，既有 package tests 不取代此 witness。cleanup 必須先解除 process-local closure gates，再查同一 DB 或恢復 page gate，防止單連線 deadlock。這批不新增一般產品 debug API、容量政策或 operator disposition 功能。

### Migration Plan

新增harness/fixture/doc，不對existing production資料做migration。cleanup僅run-owned資料與temp resources；啟動前確認DB確為可丟棄，無法確認停止。完成後先review/verify，部署及field sign-off仍需獨立授權。

### Open Questions

CI/可用PostgreSQL/Windows/ARM環境在實作時實測列出；未提供者阻擋該平台驗收而非猜通過。依實際 validation 記錄目前已執行與尚未執行的驗收；不能以文件檢查代替產品 harness。

## Risks / Trade-offs

[測試自己造好資料再讀回] → 必須經UI設定＋simulator採集＋production writer，禁止直接SQL insert當主路徑。

[前一run污染] → 唯一run namespace、fresh DB及owned cleanup。

[證據有敏感資料] → fixture-only credentials，witness去識別並人工檢查。

### ConfirmedTestWriteCleanup fixture contract

沿production group test-write preview/confirm/operation GET與同一ConnectorService.OpenDestination，先由真實採集建立鄰近production rows。preview不得改target；confirm須分開驗證typed write/readback與operation-owned cleanup；完成後same-operation confirm只回retained result、SQL不變。禁止以direct INSERT代替採集或schema metadata stub。

SQLite cleanup中斷使用owned external table的BEFORE DELETE trigger，predicate限定該preview的完整owner_value，呼叫既有tag-only closure scalar與group gate；不新增normal phase/debug API。ledger phaseCleaning已保存、DELETE尚未commit時SIGKILL，獨立SQL證明test row仍在且鄰近rows不變；重啟後等待真實三分鐘operation lease，先202再takeover只清理、不得second INSERT。fixture只刪自有trigger。

PostgreSQL只在owned disposable container建立各run唯一schema與非superuser role；初始真UI保存此帳號的destination，之後僅對本run readings分別撤銷SELECT或DELETE，保持connector identity及ownership guard。UI與saved operation須分別反映written_unverified/readback-denied、cleanup denied/failed，completed same-operation不得blind retry。production使用lib/pq，SQLSTATE分類須識別該driver與既有pgx typed errors，不依賴raw message。SQLite不具有相同grant模型，permission cases明列N/A並由PostgreSQL實測。fixture收尾僅清理精確test-owned rows與本run角色/schema，不碰其他container。

清理後的SELECT失敗不能回cleaned；以safe reason保持cleanup unknown，不由DELETE回覆推論已驗證清理。容量quota跨restart時以owned SQLite BEGIN IMMEDIATE reserved lock保持metadata可讀、實際INSERT等待；此控制不插入或修改samples，結束時ROLLBACK，避免把offline schema hydration當成quota測試。
