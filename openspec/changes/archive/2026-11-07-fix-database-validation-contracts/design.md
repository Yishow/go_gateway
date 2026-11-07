## Context

真實可丟棄 PostgreSQL 補測重現 migration 42883 與 schema apply unknown；單元回歸進一步證明 val_dec NUMERIC 被文字分類拒絕。MariaDB 補測則揭露手工 snapshot 沒帶 Apply 已要求的 counter proof。

## Goals / Non-Goals

**Goals:** 修復兩個已重現的 PostgreSQL 缺陷，修正 MySQL 系列測試 fixture，使驗證結果反映實際功能。

**Non-Goals:** 不擴充 schema 能力、不接受有界數值冒充無界 exact、不變更 fail-closed、HTTP 或使用者操作流程、不連正式DB。

## Decisions

### 小數分類只補無界 NUMERIC

val_dec 從 text 分類分出 decimal。保留既有文字型別，僅增加 PostgreSQL 生成 DDL 的無界 numeric，不接受 numeric(p,s)、decimal(p,s)、real、double、integer、boolean，避免 rounding/overflow。

### 分區函式參數明確使用 DATE

CURRENT_DATE 加 INTERVAL 後明確 cast 回 DATE。僅修復初始 migration 的錯誤呼叫，不修改既有函式簽名或分區範圍。

### 分區回歸只使用可回滾 schema

以同一 transaction 建立獨立 schema，SET LOCAL search_path 並執行 migration、寫入、讀回，最後 rollback。不 DROP public 表、不對共享資料做清理。

### 手工 snapshot 對齊 Apply proof

MySQL/MariaDB fixture 用現有 CounterBindingProofFromExactSource 或 NewCounterBindingProof 建立不可變 facts seal，放入 snapshot.CounterBindings。產品缺 proof 仍拒絕，不旁路檢查。

## Implementation Contract

- isColumnTypeCompatible 的 decimal 分類只為 val_dec 使用；其餘型別分類不變。
- schema apply 回傳 succeeded + verified digest；已存在相容 schema 的 preview/apply 保持零 DDL 且仍驗證。
- storage integration fixture 不離開獨立 schema，交易回滾後 schema 不存在。
- DB fixture 採 loopback、可丟棄資料；每個 dialect 實測狀態獨立，不能互相推論。

## Risks / Trade-offs

有界 NUMERIC 可能可容納某些值，但沒有完整 scale/precision 契約，故本次保守拒絕。原本因缺 fixture 跳過的其他平台案例不在本次修復中宣稱通過。
