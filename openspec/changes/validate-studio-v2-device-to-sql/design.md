## Context

前置A/B/C/D/E驗證。已讀proposal與既有scripts/b10_exe_acceptance.py、frontend/tests/e2e的定位；這批要新增database路徑witness，不把既有Modbus gate當成SQL驗收。主程式必須使用與正式入口相同wiring。

詳細來源見 [evidence](../../../docs/plans/studio-v2-write-groups/evidence.md)，共用欄位與政策見 [跨案契約](../../../docs/plans/studio-v2-write-groups/contracts.md)。這是擬議設計，不是完成報告。

## Goals / Non-Goals

**Goals:** 可重跑的simulated-device→真UI→persisted IDs→actual SQL row，正反案例都有可核對witness。

**Non-Goals:** 真PLC/SCADA/production DB、全面platform certification、虛構性能SLA或新增產品功能。

## Decisions

1. **隔離harness。** 每run建立獨立workspace、loopback simulator、GATEWAY_DB_PATH及外部SQLite/PostgreSQL；以實際cmd/test_ui binary和embedded frontend跑Playwright。UI/API不能stub最終成功，simulator fixture可控sample/time/failure。PostgreSQL用已授權可丟棄instance，缺依賴則blocked，不能把SQLite結果當Postgres pass。
2. **可追溯鏈。** witness保存source commit/build、平台、command、workspace/device/point/tag/group/revisions、sample/record/operation IDs、query與typed actual SQL結果、timestamps/quality、UI screenshot。祕密與endpoint憑證排除。DB rows由獨立查詢驗證，不只看API response。testwrite cleanup前後留count及neighbor row不變證據。
3. **故障矩陣。** 見acceptance.md，每個scenario對應test label、fixture和pass assertion。含multi-device same address、mixed values、bad/missing/late、outage restart、duplicate/unknown commit、disk full/poison、edited revisions與並行request。clock與fault injection具明確boundary，不能sleep一下就假設commit發生。
4. **範圍與數字。** 以小型deterministic fixture確認正確性；測試的10秒bucket等是design fixtures，不是性能目標。真實吞吐、穩定性duration或field readiness若未測就列not run。不同平台/embedded browser各自有run result，Linux結果不自動升格。
5. **closeout。** 全案包含Go/frontend full checks、line/diff、OpenSpec strict validate與independent review；失敗保存witness與repair方向。修復若超過本批scope另報，不把變更規格當作讓測試通過的捷徑。原active change的handoff只在D/E/F證據齊全後關閉，原完成記錄不重寫。

### Migration Plan

新增harness/fixture/doc，不對existing production資料做migration。cleanup僅run-owned資料與temp resources；啟動前確認DB確為可丟棄，無法確認停止。完成後先review/verify，部署及field sign-off仍需獨立授權。

### Open Questions

CI/可用PostgreSQL/Windows/ARM環境在實作時實測列出；未提供者阻擋該平台驗收而非猜通過。本次只有文件檢查，不執行產品harness。

## Risks / Trade-offs

[測試自己造好資料再讀回] → 必須經UI設定＋simulator採集＋production writer，禁止直接SQL insert當主路徑。

[前一run污染] → 唯一run namespace、fresh DB及owned cleanup。

[證據有敏感資料] → fixture-only credentials，witness去識別並人工檢查。
