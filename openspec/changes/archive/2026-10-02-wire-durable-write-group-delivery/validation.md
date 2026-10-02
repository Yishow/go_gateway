# 實作驗證紀錄

2026-10-02：依總覽 A→F 順序實作。前置 A、B 已驗證（各自 validation.md）。本文件只記錄 C 實際執行的命令與結果，不以 artifacts 齊全表示完成；C 完成不代表 D／E／F 或 production UI→SQL 已完成。

## 前置：既有 delivery 元件檢視與重用決定

- 已讀 `internal/datalink/delivery`（SQLJournal／SQLOutbox／ReceiptLedger／DeliveryWorker／QuotaMonitor／CalculateBackoff）與 migration 020。既有元件的 Append／Enqueue／Mark* 各自直接使用 `*sql.DB`、`context.Background()`，沒有交易參數；outbox 沒有 effect key、partition、claim／fencing、凍結 destination、payload digest、`unknown` 狀態；worker 在 ctx 取消時以 `context.WithoutCancel` 無界 flush。這些無法承載 group pipeline 的原子要求。
- 決定：新增專用 package `internal/datalink/groupdelivery` 與 migration `025_write_group_delivery_sqlite.up.sql`（`wg_delivery_samples`／`checkpoints`／`buckets`／`outbox`／`receipts`），重用既有可用部分（`delivery.CalculateBackoff` 與配額水位概念）而不修改 legacy 介面與其 callers。legacy 路徑與資料表原樣保留。
- 既有 `snapshot.Assembler` 擴充兩階段 API：`Check`（不改任何狀態）、`Plan`／`Commit`（先規劃關閉、持久化成功後才提交）、`Replay`（重啟時重放已接受的樣本，略過 future-skew）。`Tick` 等於 `Plan`＋`Commit`，既有行為測試維持通過。

## 1.1 DurableSampleAckCrash

- `groupdelivery.Store.AppendSample`：自己的 local transaction；`nil` 才代表已 commit，才可 ACK。同 group revision 重送相同 payload（sha256）為 no-op，同 `sample_id` 不同內容回 `ErrSampleConflict` 且不覆蓋；跨 group／revision 的相同 ID 是不同身分；缺 key／sample 欄位回明確錯誤。`Restore` 回傳 checkpoint 與所有未 consumed 樣本（依接受順序，exact value 與 revisions 保持）。
- `runtime.GroupBoundary` 新增可選 `Ledger`：`AcceptSample` 先以 `Check` 判定（refused／late／conflict 的樣本不進 journal），通過後 journal commit 成功才寫入記憶體並回 nil；commit 失敗回 `journal-failed`（安全 code，原因經 `Unwrap` 保留、`Error()` 不含 driver 文字），記憶體不變，同一 sample 之後重送可正常入 journal；`NewGroupBoundary` 啟動時 `Restore` 並 `Replay`，checkpoint 取 `max(FirstBucket, checkpoint)`。
- 範圍限制：本項只完成 sample ACK 的耐久性。row closure 的 outbox／checkpoint 原子提交是 1.2；目前 Ledger 沒有 closure 持久化，不能宣稱整條 pipeline 可恢復。

### 實際執行（1.1）

| 命令 | 結果 |
| --- | --- |
| 行為 RED：`Store` stub 回 not implemented；boundary 尚未使用 Ledger | groupdelivery 5 項、runtime 3 項 FAIL；`Assembler` 新 API stub 4 項 FAIL |
| 實作後 `go test ./internal/datalink/{groupdelivery,runtime,snapshot}` | PASS |
| `go test -race` 同三個 package | PASS |
| 真實 SQLite（`t.TempDir()` 檔案）：commit 失敗以 `BEFORE INSERT … RAISE(ABORT)` trigger 注入；ACK 後 crash 以關閉資料庫、重新開啟同一檔案、重建 boundary 模擬 | 失敗時 journal 0 筆且未 ACK；成功 ACK 的 5 筆樣本在重開後不重送即可關閉出完整 row（uint64 `9007199254740993` 精確）；重啟後相同 resend 為 no-op、內容不同被拒 |
| `golangci-lint run ./internal/datalink/...` | 新增程式無問題；剩 `service_probe.go:201`（起始 SHA 已存在） |
| `git diff --check`、`make check-lines` | PASS |

「crash」是同一 Go 行程內關閉並重開資料庫、丟棄記憶體物件的模擬，不是 `kill -9` 的行程終止；後者留給 3.3 的真實重啟測試。

## 1.2 AtomicRowOutboxCheckpoint

- `groupdelivery.Store.CommitClosure`：一個 local transaction 內寫入 bucket outcome（`wg_delivery_buckets`，每 entity／bucket 一筆，含 row／skipped／no_data 與逐 member 原因）、outbox（`wg_delivery_outbox`：凍結的 destination scope／connector ID 與 revision／database／schema／table、dedupe capability、typed payload＋sha256 digest、partition key、`pending`）、把關閉 bucket 的 journal 樣本標 consumed、並只往前推進 checkpoint（`next_close`，永不倒退）。任何一步失敗整個交易 rollback。重送已 commit 的相同 closure 為 no-op（commit 回應遺失可安全重試）；同一 effect key 但內容（digest）不同回 `ErrEffectConflict` 並 rollback；缺 key、缺凍結 destination、row 缺編碼值、row 身分不符回 `ErrInvalidClosure`。
- Payload 以 typed cell 保存（null／bool／int64／float64／string／time，整數與浮點用字串，不經 JSON number），`DecodeRowPayload` 還原原本 Go 型別：`9007199254740993` 在 payload 往返後仍精確。
- Boundary 的 durable `Tick`：`Plan` → 編碼（編碼被阻擋的 row 轉為持久化的 `skipped`／`encode-blocked:<code>`）→ `CommitClosure` 成功後才 `Commit` 記憶體。commit 失敗時 bucket 保持開啟，下次 Tick 以相同 closure 重試且不會關閉兩次。有 Ledger 時 outbox 是交接點，sink 變成可選且只在 commit 後收到非 row 的資訊性 outcome（sink 僅需在無 Ledger 時提供）。重啟後 `Restore` 取得 checkpoint，已 consumed 的樣本不重用，已關閉 bucket 的遲到樣本被拒（`late-after-close`）。
- 時間戳以固定寬度 UTC 文字保存，字典序即時間序（checkpoint 比較與 consumed 判斷依賴此點）。
- 已知限制：closure 後尚未有 sender 讀取 outbox（1.3 起）；記憶體提交發生在 DB commit 之後，若行程在兩者之間終止，重啟時由 checkpoint 與已寫入的 outbox 重建，不會產生第二筆（以重開資料庫模擬驗證）。

### 實際執行（1.2）

| 命令 | 結果 |
| --- | --- |
| 行為 RED：`CommitClosure` stub 與 `DecodeRowPayload` stub | groupdelivery 6 項 FAIL（stub 行為失敗） |
| 實作後 `go test ./internal/datalink/{groupdelivery,runtime,snapshot}`；`-race` | PASS |
| 真實 SQLite 故障注入（`BEFORE INSERT/UPDATE … RAISE(ABORT)` trigger）於四個邊界：outbox insert、bucket insert、samples consumed update、checkpoint insert | 每個邊界都回傳注入錯誤，outbox／bucket／checkpoint 皆 0 筆、journal 仍 3 筆 unconsumed；移除故障後同一 closure 重放成功（1 outbox、2 bucket、2 consumed） |
| boundary 層：outbox 失敗後重試、重啟後重 Tick、silent／skipped bucket 持久化、編碼被阻擋 row | PASS（silent bucket `[0,10)`、`[20,30)` 各一筆 `no_data`，缺 required 的 `[10,20)` 一筆 `skipped`） |
| `golangci-lint`、`git diff --check`、`make check-lines` | 新增程式無問題（僅剩 `service_probe.go:201`）／PASS |

## 1.3 SQLiteTargetReceiptIdentity

- `dbtarget.InsertGroupRow`：一個 destination transaction 內依策略提交 row 與 effect 證據。`receipt`：row 與 `gw_effect_receipts(effect_key, payload_digest, committed_at)` 同交易；已有相同 digest 的 receipt 代表 effect 已 commit，不再 insert；digest 不同回 `ErrReceiptDigestMismatch` 且絕不覆蓋。`unique_key`：`INSERT … ON CONFLICT (record_key) DO NOTHING`，衝突時讀回現有 row 逐欄比對（容忍 0/1 布林、文字 bytes、文字時間的表示差異，不容忍值差異），相同為同一 effect、不同回 `ErrExistingRowDiffers`。`none`：單純 insert。錯誤帶 phase：`pre_commit`（確定未 commit，可安全重試）與 `commit`（含 context 取消，可能已 commit 的 ambiguous）。識別字一律 quote、值一律 bound parameter；不支援的 kind／策略／缺欄回明確錯誤，不執行任何 SQL。receipt 表 DDL 只提供 `CreateEffectReceiptTable` 給 managed schema 流程與測試，sender 從不執行 DDL。
- `groupdelivery` outbox 狀態機：`pending|retrying → sending`（`BeginDelivery`，同一 item 只有一個 caller 能取得）、`CompleteDelivery`（同一本地交易：`sql_committed`＋本地 receipt）、`MarkRetry`（只對 sending；`retrying` 或重試耗盡的 `blocked`，資料保留）、`MarkUnknown`、`MarkBlocked`、`ResolveInterrupted`（卡在 sending：receipt／unique_key 可安全回 `retrying`，無 dedupe 則 `unknown`，不盲目重送）。`Sender.Deliver` 在連線前先核對 payload sha256，不符或無法解碼就 `blocked`、不開連線；`unknown`／`blocked`／已 `sql_committed` 的項目不會再送。DB 錯誤文字不存入狀態，只存安全 code（`insert-failed`、`commit-ambiguous`、`destination-identity-conflict`、`target-blocked`、`payload-digest-mismatch`…）。
- 共用測試輔助 `internal/testutil/lostcommit`：包裝 modernc SQLite driver，Commit 真正成功後才回報失敗，等同「response 遺失」；dbtarget 與 groupdelivery 測試共用。

### 實際執行（1.3）

| 命令 | 結果 |
| --- | --- |
| 行為 RED：`InsertGroupRow`／`CreateEffectReceiptTable` stub；`Sender`／outbox 狀態 stub | dbtarget 8 項、groupdelivery 11 項 FAIL（另有 1 項「拒絕不支援請求」對 stub 回錯誤而僥倖通過，不算 RED） |
| 實作後 `go test ./internal/datalink/{dbtarget,groupdelivery}`；`-race` | PASS |
| 真實 SQLite destination：commit 成功但回應遺失＋receipt 策略 | 第一次 `commit` phase ambiguous、destination 已有 1 筆；第二次讀到相同 digest receipt → `AlreadyCommitted`，仍只有 1 筆 row |
| 同情境但無 dedupe（`none`） | 狀態 `unknown`，之後 `Deliver` 回 `ErrNotDeliverable`、不再 insert（destination 仍 1 筆）；測試另證明若盲目重試會產生第 2 筆 |
| 本地 receipt 保存失敗（`BEFORE INSERT … RAISE(ABORT)` trigger）：destination 已 commit | `ErrLocalReceiptFailed`，狀態停在 `sending`；`ResolveInterrupted` 無 dedupe → `unknown`（不再送），receipt 策略 → `retrying` 並在下一次 `Deliver` 收斂為單一 effect |
| receipt insert 失敗（destination trigger） | row 一併 rollback（0 row／0 receipt），`pre_commit` 可重試 |
| `golangci-lint run ./internal/...`、`git diff --check`、`make check-lines` | 新增程式無問題（僅剩 `service_probe.go:201`）／PASS |

PostgreSQL 的同等策略與 live 驗證在 1.4；`unknown` 狀態的人工核對／處置入口尚未提供（屬 3.2 的 API／狀態來源）。

## 1.4 PostgresTargetReceiptIdentity

- 使用本機既有映像 `postgres:16-alpine` 起拋棄式容器 `gw-wg-pg-test`（127.0.0.1:55432，`POSTGRES_DSN` 為專案既有慣例，未設定時 skip；不碰既有 `qycms-postgres-1`／`datagateway-db`）。中途 Docker Desktop 重啟使 `--rm` 容器消失，已重建後重跑；以下結果均為重建後的最後一輪。每個測試使用獨立 schema 並於結束後 `DROP SCHEMA … CASCADE`。
- 沒有新增 PostgreSQL 專屬程式：`InsertGroupRow` 的 `receipt`／`unique_key`／`none` 策略本來就是 dialect-generic（`$n` placeholder、`ON CONFLICT … DO NOTHING`、schema 限定名稱），這一項的工作是用真實資料庫驗證它。因此此項**沒有行為 RED**（測試一寫就通過），不是遺漏。
- 真實 PostgreSQL 的 target-level 測試：row＋receipt 同交易（`BIGINT` `9007199254740993`、`BOOLEAN`、`DOUBLE PRECISION`、`TIMESTAMPTZ`、NULL 精確到達）；重複送出為 `AlreadyCommitted`；digest 不同被拒且不覆蓋；receipt 失敗（plpgsql trigger 注入）時 row 一併 rollback；commit 成功但回應遺失（pgx stdlib driver 以 `lostcommit.Wrap` 包裝）後 receipt 讓重試收斂為單一 effect；`unique_key` 策略對 PostgreSQL 讀回的 typed cell 全數比對相等；無 dedupe 時 ambiguous 為 `commit` phase，而 PostgreSQL 對盲目重送確實會寫出第 2 筆（證明狀態機必須停在 `unknown`）；型別錯誤為 `pre_commit`。8 個 senders 並行送相同 effect：恰 1 筆 row、1 個 committed，其餘為 already 或失敗後重試收斂（`-count=30` 反覆執行皆通過）。
- sender-level：receipt 表存在時，lost response → `retrying` → 第二次 `sql_committed`，PostgreSQL 僅 1 筆且本地 receipt 在目的地確認後才出現；custom 表無 dedupe → `unknown`，之後 `ErrNotDeliverable`、不再 insert。

### 實際執行（1.4）

| 命令 | 結果 |
| --- | --- |
| `POSTGRES_DSN=… go test ./internal/datalink/dbtarget -run PostgresTargetReceiptIdentity`（含 `-count=30` 並行／response-lost 案例） | PASS |
| `POSTGRES_DSN=… go test ./internal/datalink/groupdelivery -run PostgresTargetReceiptIdentity` | PASS |
| `POSTGRES_DSN=… go test -race ./internal/datalink/{dbtarget,groupdelivery}` | PASS |
| 未設定 `POSTGRES_DSN` 時 | 上述 live 測試 skip（不是失敗，也不是通過證據） |
| `golangci-lint run ./internal/...`、`git diff --check`、`make check-lines` | 新增程式無問題（僅剩 `service_probe.go:201`）／PASS |

限制：只驗證 PostgreSQL 16；`MySQL` 沒有此策略（`InsertGroupRow` 對其回明確不支援）；receipt 表的建立需透過 managed schema preview／confirm 流程（尚未在 UI 接線，屬 E），目前 `none` 與 `unique_key`（須驗證 record key 欄位確為 UNIQUE／PK）是不需 DDL 的策略。

## 2.1 BoundedRetryAndPoisonPartition

- 錯誤分類 `dbtarget.ClassifyInsertError`：只依 driver 錯誤碼（PostgreSQL SQLSTATE、SQLite result code），不看訊息文字，因此不洩漏也不依賴 row 內容。`row_rejected`（資料／完整性違反：SQLSTATE 22／23；SQLite CONSTRAINT／MISMATCH／TOOBIG；已存在但不同的 row／receipt）→ 永遠不會成功；`target_unusable`（授權／catalog／schema／語法與存取規則：SQLSTATE 28／3D／3F／42；SQLite generic ERROR／READONLY／AUTH；不支援的請求）→ 需修復目的地；其餘（連線、逾時、死結、資源、BUSY／LOCKED／IOERR／FULL／CANTOPEN、未知錯誤）為 `transient`，由有界重試決定而非假設。已用合成 SQLSTATE 表、真實 SQLite（NOT NULL、CHECK、缺欄、缺表、BEFORE INSERT trigger ABORT）與真實 PostgreSQL（NOT NULL、`BIGINT` 型別錯誤、缺欄、缺表）驗證。
- Sender：row 被拒 → 立即 `quarantined`（code `destination-rejected-row`，保留 payload、不計重試）；目的地不可用 → `blocked`（`target-unusable`）；暫時性失敗 → `retrying` 並由 backoff 決定下一次時間，耗盡後 `blocked` 且資料保留；ambiguous commit 與無 dedupe 的 `unknown` 規則不變。`DeliveryResult.Already` 區分「本次才 commit」與「先前已 commit」。
- `Store.ReadyHeads`／`PartitionHead`：每個 partition（group＋entity）只回傳「最早尚未完成」的 row，且僅當它是 pending／retrying 且已到期；完成＝`sql_committed`／`operator_skipped`。因此 sending、unknown、blocked、quarantined 與尚未到期的 retrying 都只擋住同 partition 後面的 row，其他 partition 不受影響。`ResolveQuarantine`：operator 明確 `retry`（重置重試額度、回 pending）或 `skip`（`operator_skipped`，payload 保留供稽核，partition 越過）；只適用 quarantined／blocked，已處置者不能再處置。
- `Dispatcher.RunOnce`：每個 partition 最多連續交付 `BatchSize`（預設 50）筆並在第一筆未 commit 時停止；最多 `MaxPartitions`（預設 4）個 partition 並行；重疊的 cycle 靠 `BeginDelivery` 單一勝者與「sending 擋後續」避免重送與亂序；失去 claim 的 `ErrNotDeliverable` 視為正常。這些數值只是有界設定與預設，不是吞吐目標。

### 實際執行（2.1）

| 命令 | 結果 |
| --- | --- |
| 行為 RED：`ClassifyInsertError` stub（一律 transient）、`Dispatcher`／quarantine／head 查詢 stub | dbtarget 4 項、groupdelivery 10 項 FAIL |
| 實作後 `go test ./internal/datalink/{dbtarget,groupdelivery}`（含 `POSTGRES_DSN` live 分類測試）；groupdelivery `-race`；partition 測試 `-count=15` | PASS |
| 真實 SQLite：A partition（1、13=poison、3）＋B partition（21、22）同一 cycle | A：1 已 commit、13 `quarantined`、3 維持 `pending`；B 兩筆都 commit；第二個 cycle 沒有任何東西越過 quarantine；operator `skip` 後 3 依序送出且 13 的 payload 仍在；修復目的地後 operator `retry` 依原順序送出 13、3 |
| transient 失敗（注入 Exec 錯誤，注入時鐘） | 兩個 partition 的 head 皆 `retrying`；backoff 未到期時沒有提前重試也沒有越過；時間到後三筆依序送出；重試耗盡 `blocked`，資料全留、後續等待，operator retry 後恢復 |
| batch／併發界限 | `BatchSize=2` 每 cycle 送 2、2、1、0；6 個 partition、`MaxPartitions=2` 時並行峰值恰為 2；4 個重疊 cycle 不重送也不亂序 |
| `golangci-lint run ./internal/...`、`git diff --check`、`make check-lines` | 新增程式無問題（僅剩 `service_probe.go:201`）／PASS |

測試調整：1.3 中依賴「BEFORE INSERT … RAISE(ABORT) 視為暫時性失敗」的兩個重試測試，改以 `lostcommit` 的 Exec 錯誤開關注入真正的暫時性失敗（ABORT trigger 現在被正確分類為 row 被拒）；1.3 其他測試與結論不變。`lostcommit.Faults` 新增 `ExecError`。目前 `sending` 的 row 若因行程中止而卡住，會擋住其 partition，要由 2.3 的 lease／fencing 與 `ResolveInterrupted` 啟動流程處理。

## 2.2 QuotaRejectsNewAck

- `Store.WithQuota`：明確的容量政策（`GlobalMaxBytes`／`GroupMaxBytes`，至少設一個；`WarningRatio` 預設 0.80、`HardRatio` 預設 0.95，與 legacy 佇列監控一致；不合理設定回 `ErrInvalidQuota`）。未設定政策時 `QuotaStatus.Configured=false`，不會假裝已被監控。用量是**實際量測**：未 consumed 的 journal 樣本與所有非 `sql_committed` 的 outbox payload 位元組（含 pending、sending、retrying、blocked、unknown、quarantined、operator_skipped）；已 commit 的 row 與已 consumed 的樣本不計。
- 規則：group 或 global 範圍的用量達到 hard 門檻即拒絕新 intake（`QuotaError`：原因 `quota-hard-limit`、範圍 global／group、group ID、已用與上限位元組；`errors.Is(err, ErrQuotaExceeded)`），狀態與拒絕一致（`QuotaStatus.State=hard_limit`、`IntakeRefused`、含「新讀值不會被記錄、已接受資料保留」的 loss-risk 說明）。超過門檻的量最多是跨越門檻的那一筆樣本；需要絕對不超過上限時 `HardRatio` 應小於 1。已持久的樣本重送（duplicate）不需要新空間，不會被拒。warning 階段仍接受資料並顯示容量偏低。拒絕發生在 journal transaction 內、插入之前，被拒樣本不會部分寫入，且 boundary 回 `GroupSampleError`（原因與 `QuotaError` 範圍保留在 `Cause`），因此不會 ACK。
- Disk full：`SQLITE_FULL` 被辨識為 `QuotaError{Reason: disk-full}`（intake 與 closure commit 皆是），已接受資料不被刪除；空間恢復後 intake 恢復。已達 hard limit 時 closure（把已接受樣本轉成 outbox row）仍可提交，因為它不增加未完成資料量的風險。
- `Store.Reclaim(retention)`：只刪「已完成」資料——已 consumed 的樣本，以及超過保留期的 `sql_committed` 或 `operator_skipped` outbox row；不會動 pending／sending／retrying／blocked／unknown／quarantined，即使配額已滿也一樣。

### 實際執行（2.2）

| 命令 | 結果 |
| --- | --- |
| 行為 RED：`WithQuota`／`QuotaStatus`／`Usage`／`Reclaim` stub | groupdelivery 7 項 FAIL（另兩項對 stub 僥倖通過，不算 RED） |
| 實作後 `go test ./internal/datalink/{groupdelivery,runtime}`；`-race` | PASS |
| 縮小容量 fixture（以單筆 journal 位元組數為單位，global 10 單位、hard 0.9） | 恰好接受 8 筆後拒絕；拒絕時 `UsedBytes ≥ 0.9×Max` 且 ≤ Max；之後再送被拒的樣本不會寫入；`Restore` 仍回傳全部先前接受的樣本 |
| group 範圍 | 單一 group 被拒（scope=group、group ID 正確），另一個 group 仍可寫入且狀態正常 |
| 真實 SQLite 磁碟已滿（`PRAGMA max_page_count` 限制檔案大小，單一連線） | 持續寫入直到真的無法成長：回 `disk-full` capacity refusal，先前接受的樣本與 outbox 都在；`Reclaim` 並提高上限後 intake 恢復 |
| `golangci-lint`、`git diff --check`、`make check-lines` | 新增程式無問題（僅剩 `service_probe.go:201`）／PASS |

限制：量測的是 payload 位元組，不含 SQLite 索引與頁面開銷；容量門檻因此應保留餘裕。未量測任何吞吐；`HardRatio` 與 `Reclaim` 保留期是設定值，不是生產建議。

## 2.3 WorkerFencingAndShutdown

- Claim／lease／fencing：`BeginDelivery(effectKey, owner, ttl)` 以單一 `UPDATE … RETURNING` 取得 claim（`Claim{EffectKey, Owner, Epoch}`），owner 為 `<node>/<incarnation>`，每次 claim 的 fencing epoch 嚴格遞增並記錄 lease 到期時間。所有 sending 出口（`CompleteDelivery`、`MarkRetry`、`MarkUnknown`、`MarkQuarantined`、`BlockClaimed`）都帶 owner＋epoch 條件；被取代的舊 worker 不論 item 之後走到什麼狀態都回 `ErrFenced` 且不改任何東西。
- 恢復：`RecoverStaleClaims(node, self, now)` 處理 lease 已過期的 sending，以及「同一 node 但不同 incarnation」的 sending（重啟後可證明已不存在），其他 node 尚未過期的 claim 不碰；每筆恢復都遞增 epoch 使舊 worker 失效。有 receipt／unique_key 可安全回 `retrying`，沒有 dedupe 則 `unknown`（不重送）。冪等。`Worker.Start` 先做恢復才開始 cycle。
- 有界時間：每次目的地嘗試有 `DeliveryTimeout`（預設 30s），結果寫入使用獨立的 `SettleTimeout`（預設 5s）context（`WithoutCancel`＋timeout，不是無界）——因此目的地已 commit 的 effect 即使呼叫端正在被取消也會被記錄；`LeaseTTL` 必須大於一次嘗試＋寫入，否則自動調高。這些是上限設定，不是效能目標。
- `Worker.Stop(deadline)`：先停止啟動新 row，等待在途交付至 deadline；逾時才取消在途嘗試並回 `ErrShutdownForced`（再等最多 2 秒），不會被卡住的目的地拖住。被中斷的 row 不會遺失或被標成完成，由下一個 incarnation 的啟動恢復處理；Stop 可重複呼叫，未啟動時為 no-op。

### 實際執行（2.3）

| 命令 | 結果 |
| --- | --- |
| 行為 RED：claim 不強制、`RecoverStaleClaims`／`Worker` 為 stub | groupdelivery 6 項 FAIL，另有並行案例卡住至 60s 逾時 |
| 實作後 `go test ./internal/datalink/groupdelivery`；fencing 案例 `-count=10`；`-race` | PASS |
| 真實 SQLite：stale worker（lease 過期後由新 sender 完成）嘗試 Complete／Retry／Unknown／Quarantine／Block | 全部 `ErrFenced`，item 維持 `sql_committed`、destination 僅 1 筆、本地 receipt 1 筆 |
| 無 dedupe 的過期 claim | 恢復為 `unknown`，新 sender 不送，舊 worker 事後 Complete 被 fence |
| 重啟恢復 | 同 node 舊 incarnation（lease 未過期）被恢復；本 incarnation 與其他 node 的 claim 不動；再次恢復為 0 |
| 兩個 worker（不同 owner）對同一 DB 並行處理 24 筆（無 dedupe 的目的地） | 目的地恰好 24 筆，24 筆皆 `sql_committed` |
| 卡住的目的地＋`Stop(200ms)` | 於 deadline 內回 `ErrShutdownForced`（< 3s）；被中斷 row 為 `sending`／`retrying`、後面的 row 仍 `pending`、健康 partition 已完成；新 incarnation 啟動後三筆各恰好送達一次 |
| 呼叫端在目的地 commit 當下被取消（`lostcommit` 的 AfterCommit hook 取消 context） | 結果仍被記錄為 `sql_committed`，本地 receipt 存在 |
| `golangci-lint`、`git diff --check`、`make check-lines` | 新增程式無問題（僅剩 `service_probe.go:201`）／PASS |

限制：lease 沒有 heartbeat 續約；單次嘗試超過 `LeaseTTL` 時由其他 worker 接手——有 dedupe 的目的地仍收斂為單一 effect，無 dedupe 的會變 `unknown` 而不是重複寫入。

## 3.1 Production 接線與 ProductionGroupOutageRecovery

- 新增 `internal/datalink/grouppipeline`：production 的群組 pipeline。`Reconcile` 依 canonical groups 與 A 的 applied snapshot 建立／結束 runtime boundary——build 前以**真實 destination 資料**檢查（connector 存在且 identity revision 與凍結值相同、kind 為 SQLite／PostgreSQL、`ReadOnlyTableInspector` 的真實欄位、tag 型別、dedupe 能力所需的 receipt 表），不能建立時記為 `blocked` 並只回安全 reason（`destination-missing`／`-unavailable`／`-revision-changed`／`-unsupported`、`destination-table-missing`／`-unavailable`、`tag-unavailable`、`layout-blocked`…，不含 driver 文字或 DSN），下一輪 reconcile 自動重試；新 revision 在生效前一個 lookahead 就先建好，舊 revision 的 boundary 以 `Until` 在新 revision 的 effective 邊界結束（樣本只進其生效區間內的 boundary）；disabled／deleted 群組在下一個 bucket 邊界結束 intake，剩餘 bucket 仍關閉並交付。`AcceptSample`（runtime 的 typed sink）fan-out 給所有 boundary 並 `errors.Join` 拒絕原因（runtime 的 `isSinkFault` 現在會檢查 joined 錯誤中任一個真正故障）。`Owns(connector, tag)` 回報某 active／closing 群組擁有該輸出。時間驅動 `TickAll`、`Start`（reconcile＋啟動 delivery worker，含啟動恢復）、`Stop(deadline)`（有界）。
- 單一 writer：`dbtarget.Writer` 新增 `SuppressMapping(connectorID, tagID)`；production wiring（`cmd/test_ui/service_wiring.go`）把 `groupPipe.Owns` 傳入 legacy writer，群組擁有的輸出 legacy writer 不再寫（其他 tag 不受影響，測試驗證）；runtime `Dependencies.SampleSink` 注入同一個 pipeline；`main.go` 在 runtime 之前啟動 pipeline，關閉順序為先 runtime、再 pipeline（有界）。Share 走另一條 fanout target writer，與 pipeline 沒有共用狀態。
- Destination 開啟：`ConnectorService.OpenDestination(connectorID, expectedRevision)`——只在 identity revision 與凍結值相同時開啟；connector 不存在／停用／kind 不支援／設定無效／憑證被拒（AuthFailed）→ `ErrDestinationBlocked`（資料保留、不改送新 endpoint）；其餘（離線、檔案不存在）→ `ErrDestinationUnreachable`（暫時性，同一 identity 重試）。**SQLite 純路徑 destination 檔案不存在時不會被建立成空資料庫**（既有 manager 會靜默建立，這裡先檢查）。`groupdelivery.Target` 新增 `Close`，sender 用完即關閉。
- boundary 新增：`Until`／`SetUntil`（superseded 與 disable）、`Retired`／`Retiring`、`NewGroupBoundaryContext`；dedupe 能力於啟用前驗證——`receipt` 需 `ReceiptTableReady`、`unique_key` 需 inspection 標為 UNIQUE／PK 的 record key 欄位，否則 `dedupe-unsupported`。`snapshot.Config.Until`／`Assembler.SetUntil`／`Done`。
- **整合測試發現並修正的真缺陷**：本地 SQLite 在 pipeline tick、journal、worker 並行寫入時，`AppendSample` 的「先讀後寫」交易會直接 `SQLITE_BUSY`（busy_timeout 對升級鎖無效），造成 ACK 被拒。修法是每個本地寫入交易先以 no-op `UPDATE … WHERE 0` 取得寫鎖；新增 8 writer × 25 筆＋closure 並行壓力測試（修正前穩定失敗、修正後 `-count=5` 通過）。
- 配額：production 預設 `groupDeliveryQuotaBytes = 500 MiB`（與 legacy 佇列預設相同的**設定值**，不是量測容量；達上限拒絕新 intake、不刪已接受資料）。

### 實際執行（3.1）

| 命令 | 結果 |
| --- | --- |
| 行為 RED：`OpenDestination` stub、boundary 區間／dedupe 測試對未實作欄位、並行寫入壓力測試 | 各自 FAIL；writer 抑制 hook 與 grouppipeline 單元測試是實作後才寫，**沒有獨立 RED** |
| `go test ./internal/datalink/{grouppipeline,groupdelivery,runtime,snapshot,dbtarget}`；`-race`（grouppipeline、cmd/test_ui、groupdelivery、runtime、dbtarget） | PASS |
| `ProductionGroupOutageRecovery`（`cmd/test_ui`，真實設定 SQLite、兩個 SQLite destination、真實 Apply；`-count=15`） | 兩個群組各自 Apply 後 pipeline 為 active；bucket 1 兩邊各 1 筆；destination A 檔案被移走後 B 持續到 3 筆、**Modbus Share 仍發佈 4321**；A 的兩個 outage bucket 為 2 筆未完成＋1 筆 `sql_committed`（SQL committed 只在 destination 真有該列時出現）；pipeline 停止並以新 incarnation 重啟後，離線檔案**沒有被重建**、B 不受影響；A 檔案恢復後 backlog 依序精確送達（21.5/100、22.5/200、23.5/300），沒有重複，所有 outbox 皆 `sql_committed` |
| pipeline 單元測試（fake 來源＋真實 SQLite store） | 8 種 blocked 原因與恢復、disable 在下一個邊界結束 intake、superseded revision 的樣本依生效時間各進一個 revision（journal 內各 1 筆）、fan-out 對不屬於任何群組的樣本不報錯而對遲到樣本回報 late |
| `go test -p 1 -parallel 1 ./...` | 47 個 package PASS，零失敗 |
| `golangci-lint`、`git diff --check`、`make check-lines` | 新增程式無問題（僅剩 `service_probe.go:201`）／PASS |

限制：(1) A 沒有 Apply route，production 目前只有 domain `Apply` 可產生 applied 群組，測試以 service API 呼叫；(2) A 的 row policy 沒有 record-key／bucket-time 欄位，群組目前只能使用 `none`（及驗證通過的 `receipt`）dedupe，`unique_key` 需 E 補欄位設定；(3) 每次交付重新開啟 destination 連線（未快取）；(4) 只驗證 SQLite destination 的整合情境，PostgreSQL destination 的整合驗證在 3.3；(5) 未做 `kill -9` 行程終止，重啟以停止並重建 pipeline 模擬，真實重啟在 3.3。

## 3.2 RevisionBoundBacklogAndDeliveryStages

- 單一狀態來源：`groupdelivery.Store.GroupStatus(groupID)`（只讀 durable store）。每筆資料恰好落在一個階段：`collecting`（已 ACK 但 bucket 未關閉）、`queued`（pending＋sending）、`retrying`、`blocked`、`quarantined`、`unknown`、`sql_committed`（只有 destination 確認過的 row）、`skipped`（operator 處置）；`last_sql_committed_at` 只取自 `sql_committed` 的時間，buffer／ACK／queued 不會讓它出現；另含最舊未完成年齡、silent／skipped bucket 數，以及**按 group revision × destination（connector ID／revision／schema／table）分組的 backlog**（含各自的安全 error code）。`grouppipeline.Pipeline.Delivery` 再加入 intake 狀態（active／retiring／blocked＋安全 reason／not_running）與配額狀態。
- API：`GET /api/v1/datalink/studio-v2/workspace/write-groups/:id/delivery`（404 群組不存在、503 沒有交付來源或讀取失敗——不是空的成功，錯誤不含路徑或 driver 文字）；Swagger 以鎖定的 `swag@v1.16.6`（安裝在 repo 外）重產，`go.mod`／`go.sum` 不變。前端新增 `types/studioV2WriteGroupDelivery.ts`、strict parser（拒絕負數／非整數計數、未知狀態、矛盾資料：有 committed 時間卻沒有 committed row 或相反、unconfigured 的 quota 卻拒絕 intake）與 `studioV2WorkspaceWriteGroupsAPI.delivery()`，主程式把 pipeline 注入為 `WriteGroupDelivery`。SSE 的 value 事件本來就不攜帶交付狀態，因此不會宣稱已交付；沒有新增 SSE 事件。
- Runtime 的假成功修正：群組擁有輸出時 legacy writer 被抑制，過去仍回 `nil`，runtime 會把它記成 database delivery「succeeded」。現在全部被抑制時回 `dbtarget.ErrOutputOwnedByWriteGroup`，production fanout 不產生 database outcome（測試驗證）。
- Revision-bound backlog 與保護：已接受的 row 凍結 group／destination revision。endpoint 被編輯（identity revision 改變）後，舊 backlog 轉為 `blocked`（`target-blocked`），**不會被送到新 endpoint**（新 endpoint 的資料表維持空）；憑證被拒（以 driver 錯誤分類，不外洩訊息）只讓該 destination blocked；離線則為暫時性。Disable 只停止新 intake、不改變舊 backlog 去向。
- **整合測試發現並修正的缺陷**：`WriteGroupService.Disable` 要求 live connector 仍是群組儲存的 revision，因此 endpoint 一旦被編輯，操作者就無法停用那個群組（無法停止 intake）。現在 Disable 只以群組自己儲存的 destination revision 做 CAS，不再要求 live connector 相符（先 RED：新增 `DisableStillWorksAfterTheDestinationEndpointWasEdited`；過期的 expected revision 仍回 conflict）。

### 實際執行（3.2）

| 命令 | 結果 |
| --- | --- |
| 行為 RED：`GroupStatus` stub；Disable 在 endpoint 編輯後失敗 | 4 項＋1 項 FAIL；之後 GREEN。router 層與前端 parser 測試是實作後才寫，**沒有獨立 RED** |
| 整合：`RevisionBoundBacklogEndpointEditAndDisableNeverRetargetOldRows`（真實 SQLite；`-count=3`） | destination 離線時 1 筆 `retrying`；編輯 endpoint 後轉 `blocked`、新 endpoint 無任何 row、delivery view 顯示 blocked=1、sql_committed=0、`last_sql_committed_at=null`、backlog 的 connector revision 為舊值且 error code `target-blocked`；Disable 成功、舊 backlog 不變、intake 為 `not_running` |
| `go test -race`（cmd/test_ui、api、grouppipeline、groupdelivery、workspace）；`go test -p 1 -parallel 1 ./...` | PASS；47 個 package、零失敗 |
| 前端：`npm test -- --run`、`npm run lint`、`tsc --noEmit`、`npm run build` | 179 檔／1045 項 PASS；lint、typecheck、build PASS（locale 仍 511 行） |
| `golangci-lint`、`git diff --check`、`make check-lines` | 新增程式無問題（僅剩 `service_probe.go:201`）／PASS |

限制：UI 尚未顯示此資料（E）；`unknown` 的人工核對入口與 quarantine 處置尚無 API／UI（`ResolveQuarantine` 為 store 層能力）；delivery view 每次請求即時查詢 SQLite，未做快取或分頁，也沒有量測其成本。

## 3.3 真實 SQLite／PostgreSQL 故障與重啟、完整基準、移交 D

- **真實行程終止**（`internal/datalink/groupdelivery/crash_test.go`）：測試以 `exec` 啟動子行程（同一個 test binary，執行真實的 durable store／boundary／sender），用 `SIGKILL` 終止並確認它是被 signal 殺死而非正常退出。(1) **ACK 後、bucket 關閉前被殺**：journal 內 1 筆已 ACK 樣本在行程消失後仍在、outbox 為 0；新行程只靠 durable 狀態 Restore → 關閉 bucket → 交付到真實 PostgreSQL，重複跑多個 cycle 仍恰好 1 筆（值 21.5 精確）、destination receipt 1 筆。(2) **PostgreSQL INSERT 交易進行中被殺**（destination 以 `pg_sleep` trigger 讓交易停在進行中，測試確認 `pg_stat_activity` 有進行中的 INSERT 且 outbox 為 `sending`）：destination 沒有留下 row；新 incarnation 以 `RecoverStaleClaims` 恢復——**有 receipt**：重試後 `sql_committed`，destination 恰好 1 筆 row＋1 筆 receipt；**無 dedupe（none）**：狀態為 `unknown`、destination 0 筆，沒有盲目重送（也就是說資料仍未寫入，需要核對／operator 處置，而不是自動補寫）。`-count=5` 反覆執行皆通過。
- **PostgreSQL destination 的 pipeline 整合**（`cmd/test_ui`，真實 Apply、`POSTGRES_DSN`、TCP proxy 讓 destination 可離線／恢復）：兩個群組各寫入各自 schema 的 PostgreSQL；bucket 1 兩邊各 1 筆；A 離線（proxy 拒絕並切斷連線）後 B 持續到 3 筆、A 的兩個 outage bucket 為已接受但未寫入；pipeline 停止並以新 incarnation 重啟；A 恢復後依序精確送達（21.5/100、22.5/200、23.5/300，BIGINT 與 DOUBLE PRECISION），destination 各 3 筆且各有 3 筆 receipt、無重複。SQLite 版本（`ProductionGroupOutageRecovery…`、`RevisionBoundBacklog…`）同時含 Modbus Share 在 destination 故障期間持續發佈、endpoint 編輯與 disable 保護；合計 `-count=40`、`-race` 皆通過。
- 容器：`gw-wg-pg-test`（`postgres:16-alpine`，127.0.0.1:55432，拋棄式）。**測試自身的缺陷與修正**：我的 PG 測試 cleanup 曾用 `t.Context()`（Cleanup 執行前已被取消），導致 DROP SCHEMA 靜默失敗、遺留 132 個 schema；而 production 的 `ReadOnlyTableInspector`（PostgreSQL）會列出整個資料庫的所有表，使群組 reconcile 變慢（測試從 0.6 秒變成約 20 秒）。已改用 `context.Background()` 並清除殘留，之後恢復並再也沒有遺留。**這也是一個限制**：PostgreSQL 的 table inspection 成本隨整個資料庫的表數量成長，尚未量測或最佳化。
- 回復保護：025 只新增表且重跑 migration 不改動已接受資料（有測試）；停用群組、編輯 endpoint 都不改變舊 backlog 去向；文件（`docs/technical/studio-v2-write-groups.md` 的「Durable 交付（C）」）記載停新 intake、保留 journal／outbox／receipt 與 destination 已提交效果、不得以舊備份覆蓋新資料。
- **dedupe 能力限制（沒有 universal exactly-once 的聲明）**：只有 `receipt`（destination 同交易寫 `gw_effect_receipts`）與經 inspection 驗證 UNIQUE 的 `unique_key` 能在 commit 回應遺失或崩潰後安全重試；custom 表預設的 `none` 在不確定時一律 `unknown`、等待核對，不自動重插；MySQL 無群組交付策略；receipt 表目前需由 managed schema 流程或操作者建立，sender 不執行 DDL；A 的 row policy 沒有 record-key 欄位，`unique_key` 在 E 補欄位前無法用於 production 群組。
- 完整基準：`go vet ./...`；`go test -p 1 -parallel 1 ./...`——**未設 `POSTGRES_DSN`：47 個 package PASS、零失敗（28 個 live PG 測試 skip）**；**設定 `POSTGRES_DSN`（live PostgreSQL）：45 個 package PASS，僅兩個既有測試失敗——`internal/api` 的 `TestStudioV2RecordingRoutes_PostgresApplyCreatesAndVerifiesTheManagedSchema`（schema apply 回 `verification_unavailable`）與 `internal/datalink/storage` 的 `TestPostgresPartitionMigration`（要求 partition 函式環境）；前者以乾淨 `HEAD` 的 worktree 對同一個 PostgreSQL 重跑，結果相同（失敗），因此是這個拋棄式資料庫與該既有測試的環境差異，不是本批回歸；後者要求特定的 TimescaleDB／partition 環境**。曾有一次全套執行中 `cmd/test_ui` 的 SQLite 整合測試失敗，之後在相同條件重現不了（單獨 6 次、`-count=40`、`-race` 皆通過），根因未確定，記為未解釋的單次不穩定。前端：`npm test -- --run` 179 檔／1045 項、`npm run lint`、`tsc --noEmit`、`npm run build` 皆 PASS。`golangci-lint`（僅剩起始 SHA 即存在的 `service_probe.go:201`）、`git diff --check`、`make check-lines` PASS。
- 移交 D：D 可使用 `groupdelivery` 的同一套 operation／claim／receipt 與 `OpenDestination`（revision 凍結、憑證阻擋、不建立缺失的 SQLite 檔）；D 的 test_write 必須沿用 production row codec 與 sender，並以 `operation 專屬 identity` 寫入測試 row，不得繞過 outbox 的 revision 凍結。尚未完成而 D／E／F 必須接手的項目：Apply route 與 UI、`unknown`／`quarantined` 的人工核對／處置 API 與 UI、`unique_key` 所需的 row policy 欄位、managed schema 建立 receipt 表、PostgreSQL inspection 成本、`kill -9` 整個 production binary 與 embedded UI 的端到端驗證（F）。

## spectra-verify（2026-10-02）與處置

verify 沒有 Critical；1 個 Warning、3 個 Suggestion：

- **W1（已修）**：production 沒有安裝 `WriteGroupBacklogOwnershipGuard`，所以已 Apply 的群組永遠無法刪除，「Delete with backlog is completed by workers」沒有被滿足也沒有被測試。現在新增 `grouppipeline.BacklogGuard`（只讀同一個 transaction）：已接受但未交付的 outbox row 與未收集的樣本必須各自指向有 immutable version 的 group revision，且 outbox row 帶完整凍結 destination（connector ID／revision、table），否則回 `ErrBacklogOwnershipUnproven` 阻擋刪除；已完成的 row 不需要 owner。主程式以 `WithBacklogOwnershipGuard` 安裝。整合測試 `…DeleteWithBacklogIsCompletedByWorkers`：destination 離線時有兩個已接受 bucket，刪除群組成功（tombstone 保留 backlog 與歸屬可查），destination 恢復後 worker 依序、各一次、以精確值交付，狀態最後為 2 筆 `sql_committed` 且 backlog 清空。guard 的單元測試涵蓋無 backlog、完整歸屬、缺 immutable version、缺凍結 destination revision、已完成 row、nil tx。注意：A 的既有契約仍要求 Delete 的 expected connector revision 與 live connector 相符（`TestNewRouter_WriteGroupCreateAndDeleteRejectStaleConnector`），所以 endpoint 被編輯之後要刪除群組，須先把群組的 destination revision 更新到新版（我嘗試放寬 Delete 與 Disable 一致，但那個既有測試明確釘住 409，因此還原；Disable 的放寬保留，因為它是緊急停止 intake）。
- **S1（已修）**：`groupPipe.Start` 失敗現在讓 `runGateway` 回傳錯誤並中止啟動，而不是只記錄後繼續（否則擁有輸出的群組的資料會被默默丟棄）。
- **S2、S3（記錄）**：unknown／quarantined 的人工處置 API／UI、PostgreSQL inspection 成本、整個 production binary 的 `kill -9` 驗證已列於 3.3 的移交；PostgreSQL 測試需要 `POSTGRES_DSN`（`key=value` 形式）。
- verify 沒有執行前端檢查與完整 Go 基準；以上兩者由本紀錄 3.2／3.3 的執行結果為準（前端 179 檔／1045 項、完整 Go 基準各自記於前述章節）。

補充（連線池）：production 以 `datalink.ApplySQLitePoolDefaults` 將本地設定資料庫設為單一連線池，程序內不會出現 `SQLITE_BUSY`；整合測試環境最初使用多連線池，因此在 `-race` 下刪除群組的測試曾遇到 workspace 交易的 `SQLITE_BUSY`。測試環境已改為與 production 相同的單連線池設定（`cmd/test_ui` 在 `-race`、開啟 PostgreSQL、`-count=10` 皆通過，也證明單連線下沒有死結）。`groupdelivery` 的「先取得寫鎖」仍保留，使它在多連線設定下（例如未來改變連線池）也不會因 `SQLITE_BUSY` 拒絕 ACK。若日後主資料庫改為多連線，A 的 workspace 交易（先讀後寫）需要同樣處理（例如 `_txlock=immediate`），這是未處理的相依。

## spectra-review（2026-10-02）與處置

review 共 1 個 Critical、4 個 Warning、2 個 Suggestion（範圍為近似歸屬，主要審 delivery 核心）：

1. **Critical（已修）**：production 的 sender 預設 `MaxRetries=5`＋backoff，destination 離線超過約 30–60 秒後，每個 partition 的 head 就耗盡重試轉 `blocked`，而 `blocked` 沒有任何 production 路徑可釋放（只有測試呼叫 `ResolveQuarantine`），destination 恢復後也不會繼續交付，違反「Outage and restart」；原本的 e2e 測試用 `Backoff: immediate` 掩蓋了這點。現在 `MaxRetries` 預設為 0＝不設上限，暫時性失敗以封頂 backoff（最長約 300 秒＋jitter）無限重試，`blocked` 只留給需要修復的原因（revision 改變、停用、憑證被拒、schema 不符）或操作者明確設定的上限（上限用盡仍保留資料）。新增使用**預設設定**的測試：15 次連續暫時性失敗仍是 `retrying`、`retry_count=15`、`next_retry_at` 在數分鐘內，destination 恢復後同一筆 row 自動交付；另有明確上限仍轉 blocked 並保留 payload 的測試。（行為 RED：第 5 次嘗試即 blocked。）
2. **Warning（已修）**：卡在 `sending` 的 row（本地狀態更新失敗、其他 worker 停擺）只會在行程啟動時恢復，執行中的 worker 不會，該 partition 會停到下次重啟。現在 `Worker` 每個 cycle 先 sweep 已過期 lease 的 claim（依 dedupe 能力恢復為 retrying 或 unknown，並遞增 fencing epoch）。測試：另一個 owner 的 claim 過期後，執行中的 worker 不需重啟就恢復並恰好交付一次。
3. **Warning（已緩解）**：每個被接受的樣本都在寫入交易內全表量測用量。新增 `usageCache`：測量結果快取 1 秒，樣本到達時即時累加（不會少算已接受資料），交付、closure、reclaim 時失效（釋放的空間最多延遲一個 TTL 才被察覺，偏保守）。測試：40 筆樣本至多量測 2 次、仍會在 hard limit 拒絕。限制：這是降低掃描頻率，不是 O(1) 計數器；backlog 很大時單次量測仍是 O(backlog)，且未量測實際成本。（該測試因引用新欄位而先編譯失敗，沒有獨立的行為 RED。）
4. **Warning（已修）**：`Reclaim` 沒有 production 呼叫者，已 commit 的 row、已 consumed 的樣本與 bucket 狀態會無限成長。`Worker` 現在依 `ReclaimEvery`／`ReclaimRetention` 週期執行（pipeline 預設每 10 分鐘、保留 24 小時；這是設定值不是建議值），`Reclaim` 另外清除超過保留期的 `wg_delivery_buckets`；pending／sending／retrying／blocked／unknown／quarantined 永不回收。測試涵蓋 worker 週期回收與 bucket 修剪。
5. **Warning（記錄，未改）**：每筆交付都重新開啟並關閉 destination 連線（連線＋驗證＋關閉；PostgreSQL 含 TLS），backlog 排空會慢。重用連線需要同時保留「每次都核對 connector revision」與「SQLite 檔案消失要被視為離線」——快取的 SQLite 連線會繼續寫入已被移走的檔案 inode，因此未實作；已列為效能限制，需在量測後以 per-batch 連線重用設計。
6. **Suggestion（已修）**：沒有任何群組擁有的 tag 也會被建立 typed envelope，被拒絕時累計成 `writeError`。新增選用介面 `runtime.SampleInterest`，`Pipeline.WantsSample` 只對有群組擁有的 device／point／tag 回 true，runtime 其餘來源不再建立 envelope；測試驗證不相關的來源不進 sink、不會被計為寫入失敗，一般儲存路徑不受影響。
7. **Suggestion（已修）**：`Worker` 的 loop 因 start context 被取消而結束時，`running` 仍為 true，之後 `Start` 會回 `ErrWorkerRunning`。現在 loop 結束會重設狀態並釋放 work context；測試驗證 context 結束後可再次啟動。過程中 `-race` 抓到我自己新增程式碼的一個 data race（loop 讀取可能已被新一輪 `Start` 覆寫的欄位），已改為把 cancel 函式以參數傳入。

驗證：`go test -race`（groupdelivery、grouppipeline、runtime、cmd/test_ui，含 live PostgreSQL）PASS；`go test -p 1 -parallel 1 ./...`（無 `POSTGRES_DSN`）47 個 package、零失敗；`golangci-lint` 僅剩既有的 `service_probe.go:201`；`git diff --check`、`make check-lines` PASS。review 建議以 `--base` 重跑取得完整覆蓋；本批沒有驗證的 base，未執行。
