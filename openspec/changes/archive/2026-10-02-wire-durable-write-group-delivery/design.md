## Context

前置A/B已驗證。已讀proposal。delivery模組有SQLJournal/SQLOutbox及worker，但production wiring仍new dbtarget.Writer；writer.takeFlushableBuckets先remove再Exec。runtime既有TargetWriter回error介面不能單靠nil表示SQL已提交。

詳細來源見 [evidence](../../../docs/plans/studio-v2-write-groups/evidence.md)，共用欄位與政策見 [跨案契約](../../../docs/plans/studio-v2-write-groups/contracts.md)。這是擬議設計，不是完成報告。

## Goals / Non-Goals

**Goals:** cmd/test_ui真正使用durable新group路徑；crash/outage後不遺失已ACK資料；每target可核對交付階段。

**Non-Goals:** 無條件exactly-once、跨DB分散交易、report/retention UI、部署或替換Go架構。

## Decisions

1. **共用現有delivery。** 先檢驗並補足SQLJournal/SQLOutbox原子transaction介面；sample ACK需journal commit。row closure與outbox/checkpoint同一交易；crash before/after每個boundary皆能從journal重建。不能在dbtarget上只包一層memory重試就稱durable。
2. **production lifecycle。** cmd/test_ui建立durable store、group projector、sender和workers；startup先hydrate版本/ownership再接受intake。同group禁止legacy writer並行輸出，Share fanout獨立。Stop有deadline，未送完保留outbox；lease/fencing防止restart stale worker或第二process重複claim。context cancelled時不能用無界WithoutCancel無限flush。
3. **target identity。** SQLite/PostgreSQL以target-side unique effect key/transactional receipt確認已提交；row及receipt在同target transaction。Send回覆遺失先查同identity，receipt digest不同則blocked不覆寫。local receipt不足以證明外部commit；無target idempotency的custom表遇不確定結果進unknown，禁止盲目retry。
4. **恢復與容量。** 每destination partition保留ordered queue及bounded concurrency/backoff。retry耗盡轉blocked仍保留資料。poison row隔離且同group/entity不越過順序；其他destination繼續。quota hard limit拒絕intake ACK並顯示scope/loss-risk；不刪pending journal/outbox/receipt來騰空間。quota界限需explicit config與storage usage量測，驗收用縮小容量fixture，不臆造production吞吐目標。
5. **revision-safe backlog。** accepted row凍結已驗證destination、schema、group與type revisions；停用僅停止新intake。endpoint切換的舊backlog保持old target並顯示repair/migration選擇；本案不實作跨目的地自動遷移。DSN/密碼不可保存到普通診斷，必要credential透過既有保護的同identity resolver取得，失效就blocked。
6. **真實status。** 修改TargetWriter結果契約或additive outcome adapter，至少分local_durable/queued/retrying/blocked/unknown/sql_committed；runtime、API/SSE及前端types依同一來源，不以buffer-return nil推導delivered。last_committed_at只記SQL證據；readback verified由D產生，收集healthy與交付failure並列。

### Migration Plan

在測試專用group先接線、故障注入，再以A的activation barrier切換單scope。legacy未移轉保持舊行為且清楚標識non-durable；撤舊路徑前核對queue歸屬。rollback停新intake、不erase已ACK資料；outbox schema向後相容或明確禁止舊binary讀取新格式。

### Open Questions

不支援destination-side dedupe的custom表已有安全unknown fallback，不是待猜的exactly-once承諾。所有capacity/retry/shutdown數值在實作設定與fixture中凍結並記錄量測。

## Risks / Trade-offs

[本地與外部無共同transaction] → 用target effect key/receipt與明確unknown處理，不宣稱原子跨DB。

[quota滿或永久壞row] → 停ACK/隔離、可操作診斷及affected scope，保留既有資料。

[舊worker仍活著] → durable claim/lease/fencing及bounded shutdown故障測試。
