## Context

前置A/B/C已驗證。已讀proposal，沿用已歸檔verified-schema-setup的operation repository及scope/revision gate。TestWrite現在明確501，不是只改前端按鈕就能開放。此案完整承接舊active tasks3.3/3.4與Truthful explicit test writes需求。

詳細來源見 [evidence](../../../docs/plans/studio-v2-write-groups/evidence.md)，共用欄位與政策見 [跨案契約](../../../docs/plans/studio-v2-write-groups/contracts.md)。這是擬議設計，不是完成報告。

## Goals / Non-Goals

**Goals:** 對所選group的explicit試寫、typed readback、owned cleanup及重啟查詢有真實證據。

**Non-Goals:** 自動在activation試寫、移除確認、任意SQL工具、清理正式資料、用測試成功取代持續交付健康。

## Decisions

1. **一個ledger。** 新增 `/studio-v2/workspace/write-groups/:id/test-write-preview` 及 `/test-write`，沿用現有 `/studio-v2/workspace/database-operations/:operation_id` 查詢。legacy recording-plans/test-write-preview及test-write route（相對同workspace base）以plan→group相容resolver轉入相同service，不另造token repository。舊僅plan_id的payload不允許直接試寫。
2. **預覽與claim。** 保存kind=test_write、workspace/group/connector expected revisions、payload digest、scope、expiry及operation_id。preview不寫目標。確認原子claim一次，與schema operation互斥同scope；錯kind422、缺欄400、foreign/unknown404、stale/expired首次claim409。running重送202，已保存結果200且不再執行；token之後過期不抹除可查的receipt。
3. **共用production codec/sender。** 產生typed fixture並以test namespace operation-owned record key送到選定target，使用C的dedupe/receipt。與正式row不共key、不改正式column語意；target無可安全標識的test資料或無法判定ownership，preview拒絕unsupported，不猜一個SQL DELETE。寫入必須等待target commit證據，local queued不回written_verified。
4. **readback/cleanup分開。** 比較實際完整typed內容、row scope、record key及provenance，不只查到相同ID就成功。write committed而read被拒則written_unverified；read mismatch則written_unverified附原因；write結果不確定為unknown。只有可證明屬此operation的test row/metadata才刪除，正式neighbor rows維持原樣。cleanup失敗/unknown獨立呈現，不能把已寫入否認成未寫入，也不能顯示已清乾淨。
5. **重啟與並行。** claim/result持久保存，process在write/readback/cleanup任一點中止後只查同operation和target receipt恢復，不用新preview繞過unknown。cleanup後仍留operation receipt；同內容重送不再產生row。UI等到E完成整合；此案的service/hooks/types可被測試但不提早打開不完整button。

### Migration Plan

以現有operation schema additive migration擴充kind/result，不重置既有建表operation。保持501直到完整確認/寫入/讀回/清理測試通過，再同步capability、API/Swagger與呼叫端。rollback關閉新claim但保留status/receipt與待核對test rows，不自動delete。

### Open Questions

沒有需要真實憑證的前置；測試只在可丟棄SQLite/PostgreSQL。custom target的識別能力不能確認時明確不支援，後續須由owner提供安全schema方案。

## Risks / Trade-offs

[試寫會留下row] → preview揭露目標及cleanup要求、只刪owned row、無法clean時可查收據。

[網路回覆遺失] → 同operation查詢，不重開第二個operation。

[簡單成功字串掩蓋部分失敗] → outcome和cleanup status分離並成套更新Go/TS。
