## Context

Deliver 在 Resolve／InsertGroupRow 前即建立 settleCtx。詳見 proposal.md；這是期限生命週期問題，不是把 timeout 全部調大就能修正。

## Goals / Non-Goals

每個遠端結果都有一次完整、有界的本地保存機會。保留原 timeout 預設、MaxRetries、backoff、lease 與 dedupe；不增加重試模式或 UI 選項。

## Decisions

1. attempt context 只涵蓋遠端 Resolve／insert／commit。該嘗試返回後，才以不繼承 caller cancellation、但有自己期限的 context 執行結果保存；所有早退與錯誤分類分支亦同。
2. 用單一小型 settlement helper 或等價集中邏輯防止分支漏修，不引入通用 transaction runner。
3. lease 至少涵蓋一次遠端嘗試、本地 settlement 及現有 margin；fencing 仍禁止失效 owner 更新狀態。shutdown 仍有總界線，取消不等於遠端未 commit。
4. 本地實際故障仍留下可恢復的 sending／unknown 證據；不得因本案而把未知結果宣稱為成功，無 dedupe 不盲目重送。

## Risks / Trade-offs

- 放寬 caller cancellation 傳播可能延後結束：只給既有 settlement budget，並保留 worker 的有界關閉。
- 純 sleep 測試不穩：用可控 resolver／driver barrier 驗證遠端耗時跨過 settlement budget，而非只測計時器。

## Migration Plan

不需要資料 migration；既有 outbox/receipt 與 claims 保留。修正只影響後續嘗試的 context 生命週期。

## Verification

先 RED，覆蓋慢成功、慢暫時失敗、慢永久拒絕、caller 在 commit 後取消，以及本地真的不可用。再修起算位置使 GREEN，最後以 SQLite／PostgreSQL 確認沒有不必要的 sending／unknown 或重複資料。
