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

### Migration Plan

先加入typed envelope adapter與純selection/identity tests，再沿runtime傳遞；新group opt-in且legacy不自動套新snapshot規則。state版本化，C負責journal與checkpoint交易；B未接C前不得向使用者宣稱重啟可恢復。

### Open Questions

無未定義的基本row策略；具體DB型別support由SQLite/PostgreSQL fixture驗證，不支援的組合回明確阻擋。

## Risks / Trade-offs

[嚴格完整性降低row數] → UI顯示原因並允許符合schema的明確partial，而非偷偷補舊值。

[設備時鐘偏移] → 明示time_origin及skew policy，保留原觀察時間。

[資料量與JSON metadata成本] → 後續量測row size/CPU，未量測不宣稱性能改善。
