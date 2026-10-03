## Context

Assembler 按 entity 產生 Outcome，layout 卻在全群組 claim 欄位並遍歷全部 Members。現有主驗收的 A/B 兩個群組不能證明同一群組多 entity 正常。

## Goals / Non-Goals

使既有多 entity 各自成列。保留 default skip_row、已確認 partial、時間選值與 exact codec；不更換 row model 或新增跨 entity 同步。

## Decisions

1. 為每個 entity 建立同一契約的 layout，或在同一 layout 中保存明確的 entity-to-members 索引；驗證與 EncodeRow 必須採同一分割。依現有能力選最小修改，不建平行 assembler。
2. 重複欄位只在同一 entity 成為衝突；不同 entity 有可驗證列身分時可合法共用。同一 entity 的兩個值不能因大小寫、顯示名或前端假 key 躲過驗證。
3. readiness、Apply、恢復描述與 frontend suggestions 對同一 persisted layout 作相同判定。既有資料不自動重分組，也不猜 SQL business key。
4. 完整、partial、missing/bad/stale 的判定只看該 entity。編碼的結構性 member-missing／layout 不一致不能被吞成正常品質跳過；保留尚未成功提交的 closure/input 並回傳安全原因。合法的品質拒絕仍沿用既有 skipped/no_data。

## Risks / Trade-offs

- entity 與 SQL key 混淆：用 persisted group/member identity；不從名稱或原始 GroupKey 猜業務值。
- 改 validator 可能放行同列碰撞：正反測試同時跑，包含同名不同大小寫與已存在 custom layout。

## Migration Plan

不改寫已接受 row 的 payload、digest、record ID 或 entity key。舊 pending journal 使用該版本凍結布局恢復；無法表示的版本保留並 blocked，不標記已消耗。

## Verification

同一群組 A/B 共用 value 欄位必須各存一列；不同欄位也不得索取另一 entity 的 member。再驗證只有 A 缺值、只有 B 正常、partial、完全無資料及同 entity 重複欄位。SQLite 與 PostgreSQL 都比對 SQL 與 checkpoint。
