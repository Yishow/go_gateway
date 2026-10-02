## Context

前置A驗證。已讀proposal。dbtarget/writer.go的Values map按arrival寫入；time.UTC().Truncate僅決定bucket，writer_statements.go的timestamp-only conflict key不能區分同時刻的不同entity。measurement/types.go已有SampleEnvelope和精確整數JSON，但尚需沿實際runtime傳遞。

詳細來源見 [evidence](../../../docs/plans/studio-v2-write-groups/evidence.md)，共用欄位與政策見 [跨案契約](../../../docs/plans/studio-v2-write-groups/contracts.md)。這是擬議設計，不是完成報告。

## Goals / Non-Goals

**Goals:** 真實sample時間/品質/型別、可重現snapshot、明確late policy及不碰撞row identity。

**Non-Goals:** durable retry接線（C）、aggregation、用電差分、報表或任意event trigger engine。

## Decisions

1. **保留採集事實。** 在collector→runtime mapping→group boundary傳typed envelope；從現有SampleEnvelope擴充/共用，不能在writer用time.Now替代observed_at。source time缺失採真實gateway acquisition time並標記time_origin。source/transform revision是identity的一部分；數值scale改變須新revision，不把舊pending值重新套新scale。
2. **可重現選值。** 完整遵循跨案契約B：UTC半開bucket、max observed_at＋stable sample_id tie-break、閉合期限、no carry-forward。相同sample_id+相同payload重送為no-op，payload不同則拒絕並留安全衝突。future timestamp超過允許clock skew時bad，skew上限為明確配置且測試注入clock，不能默認無限等待。
3. **明確缺值。** basic預設skip_row；partial需operator明確選擇且schema有NULL和逐member quality/provenance存放能力。missing/bad/stale分開診斷。skip並不等於資料完全，runtime揭露group/bucket原因，不以零補值。每個completed bucket只finalize一次；late不重開已finalized row。applied group必須由clock/tick驅動closure，完全無sample的bucket仍記no_data/skipped並可查，不能只靠sample到達建立bucket。
4. **row identity與SQL。** key至少workspace/group/revision/entity/bucket，target effect以destination scope命名。managed layout含stable record key；custom表需metadata證明exact typed存放能力，upsert／uncertain後自動retry另需verified unique effect key；純append缺dedupe可明示limited，遇unknown停下核對，完整good row的逐member provenance可保存於本地durable envelope並連結record/target identity，不強迫custom表新增外部metadata；partial或dedupe確實需變更外部schema時才走既有preview/confirm。拒絕timestamp-only unique key會碰撞不同entity的配置。
5. **型別。** bool/text/uint64/int64/decimal分開編碼；9007199254740993的十進位digits在API/UI/SQL保持，NaN/Inf標bad，overflow/config incompatibility不截斷或stringify成假正常。量測numeric encoding可重用，但raw寫入不要求physical semantic_kind。

## Implementation Contract

範圍內：collector→runtime 的 immutable sample facts、基本 exact codec、記憶體 snapshot selector／clock closure、scoped identities 及 opt-in group boundary。範圍外：production owner 切換、durable ACK/journal/checkpoint/sender、外部 DDL／真實目標寫入／UI（C–F）。既有未移轉 writer 與 Share 行為保留。

### TypedAcquisitionEnvelope 與 GatewayTimeOrigin（1.1）

- 擴充既有 measurement.SampleEnvelope：basic device_id/point_id/tag_id、source_revision/mapping_revision、acquisition source fingerprint；保留 sample/acquisition IDs、workspace、observed_at/received_at/time_origin、原 exact value 與 quality/reason。measurement/semantic fields 對 raw basic samples可空，不創造假 measurement IDs。
- collector.CollectedValue 帶 acquisition_id、actual completion received_at、typed observed_at/time_origin 及 poll 實際使用的 device/point configuration fingerprint。成功或失敗都在實際 backend read 完成時捕捉 facts；前置 connection/breaker 失敗保留 bad/reason。ReadResult 只有顯式可信 source origin 才沿用 source timestamp，否則 typed observed_at 使用 gateway acquisition completion；既有 legacy Timestamp 相容欄位不由此改寫。clock 可注入測試，不從稍後 runtime/flush 的 time.Now 補造 typed observation。
- source／mapping SHA-256 identity 共用 A 的既有 JSON string-tuple 規則；提供窄 pure helper，A 已保存的 revision 計算結果不變。collector 成功取得 ManagedConnection 後以其實際 Config／ProtocolType 計算採集 fingerprint；GetOrCreate 沿用舊連線時不冒用新 scheduler 設定。device／point metadata 固定於該次 poll snapshot。collector 只傳 opaque configuration fingerprint，不把 connection config／credential 放進 envelope。runtime binding 由實際已安裝 device/point/tag/mapping snapshot 固定版本；採集 fingerprint 與 binding 不符時拒絕 typed intake，不借當前 persisted group revision 偽造來源。
- 每 mapped sample 的 stable sample_id 由 acquisition ID＋workspace/device/point/tag/source/mapping identity 導出；同 acquisition 重送不另造 sample ID。optional runtime SampleSink 接收 typed envelope，未安裝時維持原 storage/target writer。此為 B 的 transport seam，尚無 production group consumer／durable ACK；C 安裝 consumer 時要實際保留一個 owner。
- pipeline失敗／缺 time/identity／unknown quality 有明確 bad 或拒絕結果，不變成 good／zero；正常 operator reason 使用安全代碼，不攜 raw exception。source time保留精度與 UTC。
- TypedAcquisitionEnvelope/GatewayTimeOrigin 用 fixed fixture：workspace-A、device-A/B（同40001）、point-A/B、tag-A/B（同顯示名）、acquisition-A、uint64 9007199254740993；gateway completion 2026-01-01T00:00:08Z，runtime processing 00:00:20Z；gateway origin仍 observed/received 00:00:08Z，source origin可保留00:00:03Z；source/mapping revisions與 A repository結果一致。另驗 retry/error、duplicate stable ID、metadata/scale revision不同、不洩 secrets及 legacy fallback。

### ExactMixedValueRoundTrip（1.2）

- 新增 opt-in exact value helper，wire shape 為 `type`／`encoding`／`value`。bool／text 保留原型；int64／uint64／decimal 的 value 是帶明確 encoding 的十進位 digits string；float64 必須 finite。不把任意 string 猜成 decimal，不以 float64 中介解析整數／decimal。原 SampleEnvelope MarshalJSON 與 legacy writer 不修改。
- decimal 以顯式 Decimal 型別承載有限十進位字串，先驗證語法；bool、text、signed／unsigned range 與 declared type 必須一致。unknown type／encoding、malformed digits、NaN／Inf、overflow 回安全 error，沒有 zero／fmt stringify fallback。new JSON roundtrip 能還原型別與 exact value，不只檢查 value_str。
- 新增獨立 SQL codec：輸入 exact value 和 dialect／實際 SQL type 宣告，產生 parameterized driver value，另有 typed readback decoder。此項 codec 是 pure type/range gate，測試宣告不是 production metadata verification；3.1 再以真實 inspection metadata gate 接線。SQLite bool 用 INTEGER 0/1；int64 用 INTEGER；uint64 只有值可放 signed64 才可 INTEGER bind，否則須顯式 TEXT strategy；decimal 只接受 TEXT exact strategy。SQLite NUMERIC／REAL 不能作為 decimal 或超 signed64 uint64 的 exact strategy。PostgreSQL bool/text/int64 對應 BOOLEAN/TEXT/BIGINT；uint64/decimal 另接受能容納數值且不 rounding 的 NUMERIC precision/scale 宣告；不支援或未知宣告回 blocked error。
- ExactMixedValueRoundTrip fixture：true、text batch-001、int64 -9223372036854775808、uint64 9007199254740993、Decimal 1234567890.123456789012345678；JSON encode/decode 後以真實 disposable SQLite INTEGER/TEXT columns INSERT/SELECT，再依 expected type decode 比對。另測 uint64 18446744073709551615 的 TEXT exact readback、拒絕 INTEGER overflow、decimal NUMERIC/REAL、PG NUMERIC(16,4) rounding 與 NaN／±Inf／malformed JSON。PG NUMERIC(28,18)／NUMERIC(20,0) 做 codec fixture，不稱為 live PostgreSQL。
- 此項不增加 point/tag decimal enum 或 mapping transform 功能，不更改 legacy SQL/upsert，不加第三方 decimal dependency、不接 production sender、不建外部 schema。SQLite 建表僅在 t.TempDir 的可丟棄 fixture。
- 型別依據：[SQLite storage/affinity](https://www.sqlite.org/datatype3.html)、[PostgreSQL numeric](https://www.postgresql.org/docs/current/datatype-numeric.html)、[Go SQL driver](https://pkg.go.dev/database/sql/driver)。SQLite 會按 affinity 改變數值儲存型別；PostgreSQL scale 不足會 rounding，因此 codec 先拒絕可能失真值。

後續 2.1–3.2 沿既有 Decisions、跨案契約及各具名驗收執行；不把 B 的 memory state 當 durable。沒有外部 DB 交易或產品成功狀態的推論。

### Migration Plan

先加入typed envelope adapter與純selection/identity tests，再沿runtime傳遞；新group opt-in且legacy不自動套新snapshot規則。state版本化，C負責journal與checkpoint交易；B未接C前不得向使用者宣稱重啟可恢復。

### Open Questions

無未定義的基本row策略；具體DB型別support由SQLite/PostgreSQL fixture驗證，不支援的組合回明確阻擋。

## Risks / Trade-offs

[嚴格完整性降低row數] → UI顯示原因並允許符合schema的明確partial，而非偷偷補舊值。

[設備時鐘偏移] → 明示time_origin及skew policy，保留原觀察時間。

[資料量與JSON metadata成本] → 後續量測row size/CPU，未量測不宣稱性能改善。
