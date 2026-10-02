# 實作驗證紀錄

2026-10-02：依總覽 A→F 順序實作。A 已驗證（見 unify-studio-v2-write-group-contract/validation.md）；本文件只記錄 B 實際執行的命令與結果，不以 artifacts 齊全表示完成。B 的 memory assembler 不是 durable，重啟恢復由 C 負責。

## 1.1 TypedAcquisitionEnvelope 與 GatewayTimeOrigin

- 採集事實：`CollectedValue` 帶 acquisition_id、`ObservedAt`／`ReceivedAt`／`TimeOrigin`、opaque config fingerprint 與安全 `QualityReason`；時間在實際 backend read 完成時以可注入 clock 擷取，只有 adapter 顯式標示 source origin 且時間非零才保留來源時間，失敗讀取不採信來源時間。
- 重試以最後一次讀取完成時間為準；good quality 附帶 Error 降為 bad／`read-failed`；unknown quality 轉 bad／`quality-unknown`；熔斷與連線失敗保留 `circuit-open`／`connection-failed`，原始錯誤文字不進入 typed reason。
- fingerprint 在成功路徑取自 ManagedConnection 實際 Config／Protocol（沿用舊連線時不冒用新 scheduler 設定）；熔斷與連線失敗路徑取自本次嘗試的 device／point 設定。envelope 不含連線設定或憑證。
- runtime 由實際已安裝 binding 驗證 device／point／fingerprint；不符、缺 time／identity 或未知 time origin 拒絕 typed intake。stable sample_id 由 acquisition ID＋workspace／device／point／tag／source／mapping revision 導出。SampleSink 未安裝時維持原 storage／target writer 行為。
- 這是 transport seam，尚無 production group consumer 或 durable ACK（C）。

### 實際執行

| 命令 | 結果 |
| --- | --- |
| `go test ./internal/datalink/runtime ./internal/datalink/collector ./internal/datalink/measurement -count=1` | 修正前 FAIL：collector 3 項（good+error 仍 good、熔斷／連線失敗缺 fingerprint）、runtime 1 項（`connection-failed` 被覆寫為 `read-failed`）；修正後 PASS |
| `go test -race` 同三個 package | PASS |
| `go vet ./...` | PASS |
| `go test -p 1 -parallel 1 ./... -count=1` | 無失敗 package |
| `git diff --check`、`make check-lines` | PASS |

未執行：live PostgreSQL、Windows／ARM、production UI→SQL。

## 1.2 ExactMixedValueRoundTrip

- `measurement.ExactValue`：wire shape `type`／`encoding`／`value`；bool／text／float64 用 `json` encoding，int64／uint64／decimal 的 value 是十進位 digits 字串（`decimal-string`），不經 float64。`NewDecimal` 只接受純十進位語法（拒絕指數、NaN／Inf、空白、空字串）；`NewFloat64` 只接受 finite。Unmarshal 嚴格：未知 type／encoding、type 與 encoding 不符、數字型別用 JSON number、int64／uint64 超出範圍、unknown field、trailing data、null 全部拒絕；未設定的零值 Marshal 回 `ErrExactValueInvalid`，不輸出假值。decimal `Equal` 以數值比較（忽略多餘零），不同 type 一律不相等。原 `SampleEnvelope.MarshalJSON` 與 legacy writer 未修改。
- `dbtarget.EncodeExactValue`／`DecodeExactValue`：獨立的 pure type/range gate，宣告型別是呼叫端宣稱，不是 production metadata 驗證（3.1 才以實際 inspection 接線）。SQLite 依 affinity 規則：bool→INTEGER 0/1；int64→INTEGER；uint64 只有 ≤ MaxInt64 才能 INTEGER，否則需顯式 TEXT；decimal 只接受 TEXT；NUMERIC／REAL／DECIMAL(p,s)／FLOAT／BLOB 一律 blocked。PostgreSQL：BOOLEAN／TEXT（VARCHAR(n) 檢查長度）／BIGINT／DOUBLE PRECISION；uint64／decimal 另接受 NUMERIC(p,s)（小數位或整數位會被 rounding／超出即 blocked，尾端零不算失真）、無限制 NUMERIC 與 TEXT；INTEGER／SMALLINT／REAL／MONEY 等 blocked。Readback decoder 依預期型別拒絕 REAL 讀回整數欄、nil、bool 非 0/1、非數字字串。
- 不增加 point/tag decimal enum 或 mapping transform，不改 legacy SQL／upsert，不加第三方 decimal 依賴，不接 production sender，不建外部 schema；SQLite 建表只在 `t.TempDir` 的可丟棄 fixture。

### 實際執行

| 命令 | 結果 |
| --- | --- |
| 行為 RED：可編譯 stub 回傳 not implemented，`go test ./internal/datalink/measurement ./internal/datalink/dbtarget -run Exact` | measurement 5 項、dbtarget 7 項 FAIL（stub 行為失敗，非缺 symbol） |
| 實作後同命令 | PASS（含真實 SQLite INTEGER／TEXT INSERT／SELECT：true、`batch-001`、int64 最小值、uint64 9007199254740993、uint64 18446744073709551615 的 TEXT 讀回、decimal 1234567890.123456789012345678） |
| `go test -race` measurement／dbtarget | 單獨 dbtarget 連跑 5 次 PASS；measurement＋dbtarget 同時執行 4 次中有 1 次 FAIL：既有 `TestWriter_InsertModeRowGroupsEmitSeparateRowsForSharedColumn`（`writer_row_groups_test.go:108`，flush 等待逾時 observed=1 want=2），該檔 tracked 且本批未修改，為負載下的時間門檻型 flake；其餘 3 次同命令 PASS。未修改該測試，後續 C 若動 writer 需留意 |
| `go vet ./...`、`go test -p 1 -parallel 1 ./... -count=1` | PASS／無失敗 package |
| `golangci-lint run` measurement／dbtarget／collector／runtime | 新增程式無問題；剩 `service_probe.go:201` gocritic（起始 SHA 已存在，未修改） |
| `git diff --check`、`make check-lines` | PASS |

PG NUMERIC(28,18)／NUMERIC(20,0)／NUMERIC(16,4) 僅是 codec fixture，不是 live PostgreSQL；Docker daemon 未啟動，沒有 live PostgreSQL 證據。

## 2.1 UTCWindowBoundaryAndReordering

- 新增純記憶體 package `internal/datalink/snapshot`（無 DB、journal、sender；注入時間）。`BucketStart` 以 UTC 半開 `[start,end)` 對齊，含 1970 前以 floor 而非向零截斷；`Selector` 每 member 保留最大 `observed_at`，同時刻取字典序最大 `sample_id`，與到達順序無關（20 次亂序重現相同結果）。
- 設計決定（規格未明示處）：最新觀察即使是 bad 也勝出，避免較舊的 good 值掩蓋 bucket 尾端的讀取失敗；bad 如何影響 row 由 2.2 的 required／freshness 規則處理。
- 重複 `sample_id`：payload 相同為 `duplicate` no-op；不同為 `identity_conflict`，選值不變（對落選的舊 sample 同樣適用，且先於 bucket 檢查，竄改時間仍報 conflict）。`observed_at` 超過 `now + MaxFutureSkew` 以 `future-skew` 拒收且不建立 bucket，skew 未設定視為零容忍，不會無限等待。t=10 邊界、未知 member、缺 identity、bucket 外樣本皆明確 rejected。
- 這是 selector；late／closure／row 組裝在 2.2 與 2.4。

| 命令 | 結果 |
| --- | --- |
| 行為 RED：可編譯 stub，`go test ./internal/datalink/snapshot` | 9 項 FAIL |
| 實作後同命令；`go test -race -count=3` | PASS |
| `golangci-lint run ./internal/datalink/snapshot/...`、`git diff --check`、`make check-lines` | 無問題／PASS |

## 2.2 FreshnessMissingAndSilentBucket

- `snapshot.Assembler`（純記憶體、注入時間）：`Offer` 依 member 的 entity 與 `observed_at` 路由到 bucket；`Tick(now)` 由時間驅動，當 `now >= end + allowed_lateness` 時依序關閉 bucket，**包含從未收到樣本的 bucket**，每個 entity／bucket 恰好產生一個 `Outcome`（`row`／`skipped`／`no_data`），關閉後狀態刪除，重複 Tick 不會再出現。長時間停頓後的補關閉以 `MaxBucketsPerTick`（預設 1000）限制單次工作量，其餘於後續 Tick 依序輸出。
- 預設 `skip_row` 且所有 member 必須 required（否則 config 無效）：缺值（`missing`）、bad（保留 collector 的安全 reason，否則 `quality-<q>`）、good 但無值（`invalid`）、`end - observed_at > max_age`（`stale`）都不產生 row，留下 bucket 與逐 member 原因；整個 bucket 完全沒有該 entity 的樣本為 `no_data`／`no-samples`。max_age 預設為 interval，且絕不從其他 bucket 補值（測試：前一 bucket 的 pressure 不會被後一 bucket 使用）。被 future-skew 拒收的樣本不建立 bucket，不會讓 bucket 無限期保持開啟。
- `partial` policy 在 2.3 之前一律拒絕（`ErrInvalidConfig`）；關閉後才到的樣本回傳 `late`、不改狀態（2.4 補完整驗收）。

| 命令 | 結果 |
| --- | --- |
| 行為 RED：可編譯 stub，`go test ./internal/datalink/snapshot` | 12 項 FAIL |
| 實作後同命令；`go test -race -count=3` | PASS |
| `golangci-lint run ./internal/datalink/snapshot/...`、`git diff --check`、`make check-lines` | 無問題／PASS |

## 2.3 ExplicitPartialPolicy

- `incomplete_policy=partial` 只有 `Config.Storage`（由呼叫端依真實 table metadata 驗證的 `StorageCapabilities`：`NullableValues` 與 `MemberQuality` 皆為真）才被接受，否則 `NewAssembler` 回 `ErrInvalidConfig` 包 `ErrPartialUnsupported`，即在 activation 前阻擋；僅 nullable 或僅 quality 都不足。
- 設計決定：partial 模式下 `Required=true` 的 member 仍把關（缺／bad／stale 時整個 row 為 `skipped`，不寫半成品），`Required=false` 的 member 才以 SQL NULL＋reason 表示；沒有任何可用 member 為 `skipped`／`no-usable-member`，整 bucket 無樣本仍是 `no_data`。部分缺值的 row 標 `Partial=true`。此選擇是為了讓既有 `required` 欄位在 partial 下仍有意義；E 的 UI 需明確揭露哪些 member 可為 NULL，避免全 required 的 partial 群組永遠不會產生 NULL。
- `MemberResult.Usable()` 為假的 member（missing／bad／stale／invalid）一律寫 NULL，不使用 0 或舊值；其 `Sample`（含 bad 樣本）保留 observation 與 quality 供本地 durable envelope／逐 member provenance。預設 `skip_row` 且零儲存能力的 custom 表，完整全 good row 仍可產生，且每個 member 帶 `Sample`（sample_id、quality）作為 local provenance，不要求外部 metadata 欄位或 companion 表。
- UI 必須揭露：品質證據只在本地，SQL 表本身未攜帶逐欄資訊（屬 E）。

| 命令 | 結果 |
| --- | --- |
| 行為 RED：新增測試對尚未實作 partial 的 assembler | 7 項 FAIL（另 1 項「完整 good row 不要求外部 metadata」在實作前已通過，記錄既有行為，不算 RED） |
| 實作後 `go test ./internal/datalink/snapshot`；`-race -count=3` | PASS |
| `golangci-lint`（snapshot）、`git diff --check`、`make check-lines` | 無問題／PASS |

## 2.4 ScopedRowIdentityAndLateArrival

- `snapshot.RecordID`：以版本標記＋workspace、group、group revision、UTC bucket start 與 entity scope 的 JSON-tuple SHA-256 導出；空 entity 使用固定 `group-scope` 標記，與字面 entity key 以不同 tuple 形狀區分，欄位間不會因分隔字元位移碰撞。`DestinationScope`（connector ID／revision、database、schema、table）與 `EffectKey(scope, record)` 另命名空間；同 row 重送同目標得到相同 effect key，不同目標不同。缺任何身分欄位回 `ErrInvalidConfig`，時區不影響 ID。
- `Assembler` 現要求 WorkspaceID／GroupID／GroupRevision／DestinationScope；每個 entity／bucket 的 `Outcome` 帶 `RecordID`，只有 `row` 帶 `EffectKey`（skipped／no_data 只是 scoped 狀態，沒有 destination effect）。同 bucket、同 timestamp 的兩個 entity 產生兩個 row 與不同 ID／effect key；一個 entity 靜默時另一個仍寫 row，靜默者留下 scoped `no_data`。以相同設定重建 assembler（模擬重啟）對同 bucket 導出相同 ID，不代表 open bucket 狀態可恢復。
- Late：bucket 關閉後到達的樣本（含已接受 sample 的重送）回 `late`、累計 `Late()`，不重開 bucket、不建立狀態、不產生第二個 outcome，已輸出的 row 不變。關閉前重複相同 payload 為 no-op、payload 不同為衝突，選值不變。
- Revision：member 可宣告期望的 `SourceRevision`／`MappingRevision`，不符的樣本以 `revision-mismatch` 拒收，因此 scale／mapping 變更後不會把舊 pending 值重新套用；revision 也是重複判定 payload 的一部分。新 group revision 的 row ID 與舊 revision 不同，既有輸出保持原 identity。
- 範圍外（C）：這些 effect key 尚未接 production sender，也沒有宣稱 target 端的 dedupe 能力；custom append 是否 limited 屬 C／D 的 destination 能力驗證。

| 命令 | 結果 |
| --- | --- |
| 行為 RED：先加 API 面 stub（RecordID 等回傳空字串、`Late()` 回 0），`go test ./internal/datalink/snapshot` | 9 項 FAIL |
| 實作後同命令；`go test -race -count=3` | PASS |
| `golangci-lint`（snapshot）、`git diff --check`、`make check-lines` | 無問題／PASS |

## 3.1 新 group runtime boundary（opt-in）

- `measurement.ExactTypeForTag`／`ExactFromGo`：tag 型別 → exact 型別（bool→bool、int16/32/64→int64、uint16/32/64→uint64、float32/64→float64、string→text；沒有任何 tag 型別推導 decimal），pipeline 值只做同類轉換，跨種類（int 給 float tag、float 給整數 tag、字串給整數）、NaN／Inf、超出範圍整數一律回 `ErrExactValueInvalid`，不改值、不補 0。
- `dbtarget.NewGroupRowLayout`／`EncodeRow`：以**真實 inspection 欄位**（`ColumnInfo`：型別、nullable）驗證每個 member 欄位，一次回傳全部 code-only 問題（`column-missing`、`unsupported-sql-type`、`duplicate-column`、`column-not-nullable`、`partial-requires-provenance`、`identity-column-unsupported`／`-missing`、`member-incomplete`、`unsupported-dialect`）；`Capabilities()` 由實際欄位導出（optional 欄位皆 nullable＋provenance 欄位可存文字／JSON），供 snapshot 判斷能否接受 partial。`EncodeRow` 用 1.2 的 exact codec 產生 parameterized 值：型別與 layout 不符回 `type-mismatch`、SQL 範圍／精度不足回 `sql-value-blocked`、非 partial 卻出現不可用 member 回 `unexpected-null`，錯誤只含 code 與欄位名；partial 的不可用 member 為 NULL（保留 reason 與 sample 於 provenance JSON），不寫 0 或舊值。可選的 record-key／bucket-start／entity／provenance 欄位依實際欄位型別編碼（bucket 欄位 TEXT 用 RFC3339Nano、TIMESTAMP／DATETIME 用 UTC `time.Time`）。
- `runtime.GroupBoundary`（實作既有 optional `SampleSink`）：以 persisted `device/point/tag` 組成 member key 路由（相同 address 或相同顯示名不影響）；不屬於本 group 的樣本直接忽略；workspace 不符、stale source／mapping revision、late、identity conflict 回 `GroupSampleError`（code only）；good 樣本的值無法 exact 轉換就改為 bad／`type-mismatch` 且不攜帶值。`Tick`／`RunTicks` 由注入時間驅動，把 `row` 轉 `GroupRowSink.AcceptRow`，`skipped`／`no_data` 與編碼被阻擋的 row（`encode-blocked:<code>`）轉 `ReportOutcome`，絕不寫半個 row；sink 錯誤以 `errors.Join` 回傳。`NewGroupBoundary` 在啟用時阻擋：非 applied 版本、缺 sink、成員身分不全、未知／不支援的 tag 型別、layout 問題、無效 destination scope、無效 snapshot 設定（含 interval 缺失、未對齊 first bucket、skip_row 下的 optional member）。
- 既有 `Service`、legacy writer 與 Share 路徑未修改；boundary 目前沒有被主程式 wiring 取用，legacy 未移轉路徑維持原樣（既有 runtime／dbtarget／API 測試全數通過）。

## 3.2 回查邊界與基準

- Migration Plan：typed adapter 與純 selection／identity tests 先於 runtime 傳遞；新 group 為 opt-in，未移轉 legacy 不套 snapshot 規則；`record_id`／effect key 使用版本標記（`write-group-record-v1` 等）。**B 的 assembler 只在記憶體**：未 ACK、無 journal／checkpoint，重啟會遺失 open bucket；`RecordID` 對相同設定可重建不代表狀態可恢復。journal／checkpoint／outbox 交易移交 C，B 未接 C 前不得向使用者宣稱重啟可恢復。`GroupRowSink` 是 C 的接點：目前 sink 失敗後該 row 沒有重送機制。
- Open Questions：SQLite 以真實 disposable 資料表驗證 INSERT／SELECT；PostgreSQL 只有 codec fixture（NUMERIC(28,18)／(20,0)／(16,4)、timestamptz、BIGINT）與 layout 規則，沒有 live PostgreSQL（本機 Docker daemon 未啟動），不支援的組合回明確阻擋。
- 設計決定與限制（需 C／E 處理）：(1) A 的 row policy 沒有 record-key 與 bucket-time 欄位，boundary 以可選設定接受，未設定時 SQL row 只含 member 欄位（local record／effect key 仍存在）；(2) A 的 member 沒有逐 member quality 欄位，partial 能力只由單一 provenance 欄位推導，因此現有 group 設定在沒有 provenance 欄位時 partial 會於啟用前被阻擋；(3) partial 下 `Required=true` 的 member 仍把關整個 row；(4) 最新觀察即使為 bad 也勝出（見 2.1）；(5) 既有 `TestWriter_InsertModeRowGroupsEmitSeparateRowsForSharedColumn` 在 race＋多 package 並行時偶發逾時，與本批無關且未修改。

### Fixture 預期值（已有測試斷言）

| 項目 | 預期 |
| --- | --- |
| 採集事實 | gateway completion `2026-01-01T00:00:08Z`、runtime 處理 `00:00:20Z` → observed／received 皆 `00:00:08Z`、origin `gateway`；source origin 保留 `00:00:03Z`；uint64 `9007199254740993` 精確保留 |
| 10 秒 bucket | t=8 勝 t=3（任意到達順序）；t=10 屬下一 bucket；同時刻取字典序較大 `sample_id` |
| 缺值／靜默 | 缺一 required member → 零 SQL row＋`skipped`／`incomplete-required-member`；整 bucket 靜默 → 每個 entity 一個 `no_data`；不補前一 bucket |
| 混合 row | float 21.5、int 101、bool true→1、text `batch-001`、uint64 `9007199254740993`（INTEGER 欄）精確編碼；`MaxUint64` 對 INTEGER 欄 → `sql-value-blocked`、不寫 row |
| 身分 | 同 timestamp 不同 entity → 兩個不同 `record_id`／effect key；revision 改變 → 新 ID，舊 row 不變；late 不改 closed row |

### 實際執行（3.1／3.2）

| 命令 | 結果 |
| --- | --- |
| 行為 RED：`ExactFromGo` stub／`NewGroupRowLayout` stub／`NewGroupBoundary` stub | 3＋7＋12 項測試 FAIL（stub 行為失敗） |
| 實作後 `go test ./internal/datalink/{measurement,dbtarget,runtime,snapshot}`；同組 `-race` | PASS |
| `go vet ./...` | PASS |
| `go test -p 1 -parallel 1 ./... -count=1 -json` | 45 個 package PASS，零 fail event |
| `golangci-lint run ./...` | 只剩 `service_probe.go:201` gocritic（起始 SHA 已存在，未修改） |
| `git diff --check`、`make check-lines` | PASS（歷史 locale 511→511、`dbtarget/service.go` 923→895） |

未執行：live PostgreSQL、production UI→SQL、Windows／ARM、LAN／真 PLC／SCADA、長時間 soak、前端（本批無前端變更）。

## 補充：live PostgreSQL 證據（本輪稍後取得）

先前記錄「本機 Docker daemon 未啟動、沒有 live PostgreSQL」。後來啟動 Docker Desktop，使用本機既有映像 `postgres:16-alpine` 起拋棄式容器 `gw-wg-pg-test`（127.0.0.1:55432，不碰既有 `qycms-postgres-1`／`datagateway-db`），以 `POSTGRES_DSN`（專案既有慣例，未設定時 skip）執行：

- `TestExactMixedValueRoundTripLivePostgres`：真實 PostgreSQL 的 BOOLEAN／TEXT／BIGINT（int64 最小值）／NUMERIC(20,0)（uint64 `18446744073709551615`）／NUMERIC(28,18)（decimal `1234567890.123456789012345678`）／TIMESTAMPTZ，經 exact codec INSERT 後 SELECT 並依預期型別解碼，型別與 exact 值皆保持（PASS）。
- `TestExactPostgresNumericScaleRoundsSilentlyWithoutTheGate`：同一資料庫對 NUMERIC(16,4) 直接收下並靜默 rounding 成 `1234567890.1235`；codec 在送到 driver 前以 `ErrExactSQLBlocked` 拒絕；BIGINT 欄位對 MaxUint64 實際報錯，codec 同樣先行阻擋（PASS）。
- 既有 dbtarget PostgreSQL 測試（schema execution／inspection）在同一資料庫也 PASS，作為環境基準。

這只補足 codec／inspection 層的 live PG 證據；仍沒有 production UI→SQL、Windows／ARM 或 LAN／真 PLC 的驗收。

## spectra-verify／spectra-review（2026-10-02）

- verify：無 Critical，8/8 完成；W1（同位址兩設備的 spec example 證據）已補 `TestHandleCollectedValueKeepsSameAddressDevicesDistinctAndFabricatesNothing`（device-A／B 同為 40001、tag 顯示名相同、處理時間不影響 observed／received、重送同 sample ID、無 measurement ID、envelope 不含憑證與連線設定）；W2「透過 API」的 exact digits 證據是 `ExactValue` 的 wire shape 測試（int64／uint64／decimal 以十進位字串傳遞，不經 JSON number）；W3 的 int64 最小值作為合法值已在 `exact_value_test.go`、`exact_value_codec_test.go`（SQLite INTEGER、PostgreSQL BIGINT live）斷言；W4 的「production consumer」由 C 的 3.1 承接。
- review：
  1. **已修（正確性）**：mapping 轉換失敗時 `continue` 跳過了 typed sample，bad 讀值不會到群組，舊的 good 值會被寫進 row，違反「Failed or non-finite read」。現在 `offerFailedTypedSample` 在轉換失敗時仍送出 bad sample（已失敗的讀取保留 collector 的安全 reason；good 讀取但轉換失敗為 `mapping-failed`），值為 nil，`mappingError` 仍計數。新增兩個 runtime 測試（scale pipeline＋失敗讀取、good 讀取但轉換失敗），先 RED 後 GREEN。
  2. **已修**：sink 回傳的正常拒絕（late、identity conflict、revision mismatch、skew 等）不再累計到 `writeError`；只有 `journal-failed`、`quota-hard-limit`、`disk-full` 與非 `GroupSampleError` 的錯誤視為寫入故障。新增測試。
  3. **記錄，未改**：collector 與 runtime 各有一份 `isKnownQuality`；新增 `QualityFlag` 時需同步兩處，後續可抽成 `schema.QualityFlag.Known()`。
  4. **記錄，未改**：scheduler 的時鐘／acquisition ID 有 config 欄位與 `With*` setter 兩條注入路徑，setter 沒有加鎖；目前只有測試呼叫、且都在 `Start` 前，生產程式不使用 setter。
