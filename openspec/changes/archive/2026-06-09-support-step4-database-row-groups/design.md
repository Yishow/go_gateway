## Context

目前 `/studio/v2` Step 4 將所有已啟用 point 視為同一張表上的唯一欄位綁定，前端 state 與 backend workspace contract 都只有 `DbConnector` 與 `DbTarget(column_name, enabled)` 這種單層模型。這個模型足以表達「多規則共表、各自用不同欄位」，但無法表達「多規則共表、重用同一組欄位、改由不同列或不同群組區分」。

使用者的新需求不是單純的 UI 微調，而是要把 Step 4 從「欄位唯一綁定器」擴充成「同表多列群組規劃器」。一旦進入 shared-column row group，前端規劃、workspace 持久化、readiness 驗證、runtime 寫入模式都要對齊同一套語意，否則很容易出現 UI 說合法、autosave 說合法，但 runtime 寫入其實無法分列的契約破洞。

## Goals / Non-Goals

**Goals:**

- 定義 row-group / write-group 的最小持久化模型，讓 workspace 能描述「同表、共欄、分群組」。
- 定義 Step 4 shared-column 的合法條件與錯誤條件，讓 validation 不再只回報 generic column conflict。
- 限定 runtime 支援邊界，避免先把不安全的 shared-column + upsert 組合放行。

**Non-Goals:**

- 不處理跨 connector 或跨 table 的多目標 fan-out。
- 不重新設計 `/studio` 舊版 database-target-workbench。
- 不在本 change 內擴充 database schema introspection、DDL 自動推導或歷史資料回補。

## Decisions

### Decision: Persist row groups separately from DbTarget rows

row-group 需要描述群組識別、所屬規則集合、group key 欄位與可重用欄位範圍，這些資訊不是單一 `DbTarget` 能穩定承載的。設計上將 row-group 視為 workspace database 設定下的第一級資源，`DbTarget` 則保留 point-to-column 綁定與 row-group reference。

替代方案是把 group metadata 直接塞進每一筆 `DbTarget`。這會造成同一群組資訊重複散落在多筆 target，任何一筆 target 漏改都會讓 workspace 進入部分一致狀態，因此不採用。

### Decision: Validate shared columns inside a row-group scope

shared-column 不再以「欄位名稱重複」直接視為錯誤，而是必須先判斷這些 point 是否屬於同一 row-group、是否共享同一 connector/table、以及 group key 是否足以區分列實例。只有「跨 group 重複」或「group 內缺少列識別」才是 blocking issue。

替代方案是保留現有全域欄位唯一規則，再額外提供例外清單。這會讓 UI 和 backend 都需要雙軌衝突邏輯，最終比顯式 row-group scope 更難理解，因此不採用。

### Decision: Limit shared-column support to insert-first delivery

shared-column row-group 的第一版以 `insert` 為主，因為每次寫入可以自然形成新列。`upsert` 若沒有穩定且唯一的 composite key，會把不同群組資料覆寫到同一列，因此第一版必須明確阻擋不完整的 upsert 配置，只在後續定義完整 key contract 後再放行。

替代方案是同時允許 insert 與 upsert，並把 key 推導留給 runtime 猜測。這會把資料破壞風險延後到正式寫入時才暴露，不可接受。

## Implementation Contract

- **Behavior**:
  - Step 4 SHALL allow operators to create one or more row groups under a single connector/table context.
  - Points inside the same row group MAY reuse the same business columns when the row group defines a valid group key or row identity contract.
  - Shared-column mappings across different row groups SHALL remain blocking conflicts.
  - `upsert` SHALL remain blocked for row groups that do not provide a stable uniqueness key.
- **Interface / data shape**:
  - Workspace database contract SHALL add a persisted row-group collection under the Step 4 database setup.
  - Each row group SHALL include a stable id, connector/table scope reference, member point ids, and group-key metadata.
  - Each persisted database target SHALL reference either the default single-row mode or a row-group id.
- **Failure modes**:
  - Missing row-group identity, cross-group shared-column reuse, and unsupported `upsert` contracts SHALL surface as readiness blocking issues.
  - Legacy single-row targets without row-group metadata SHALL continue to load in compatibility mode and SHALL NOT be auto-migrated into shared-column groups silently.
- **Acceptance criteria**:
  - Frontend unit tests prove shared-column mappings are accepted within one row group and blocked across groups.
  - Backend service / handler tests prove workspace autosave persists and reloads row-group metadata without losing legacy Step 4 targets.
  - Runtime/database delivery tests prove insert-mode row groups write distinct rows and that invalid upsert row groups are rejected before apply.
- **Scope boundaries**:
  - In scope: Step 4 planner semantics, workspace persistence contract, readiness validation, and runtime delivery guardrails.
  - Out of scope: generalized multi-table orchestration, legacy `/studio` output planner migration, and retrospective data migration.

## Risks / Trade-offs

- [Risk] Shared-column semantics are harder for operators to understand than one-column-one-point. → Mitigation: expose explicit row-group summaries and blocking reasons instead of silent exceptions.
- [Risk] Workspace compatibility may break if row-group fields are treated as mandatory too early. → Mitigation: keep legacy single-row targets readable and only require row-group fields when shared-column mode is enabled.
- [Risk] Runtime insert/upsert divergence could confuse operators. → Mitigation: make write-mode limits part of Step 4 validation and runtime readiness messages.

## Migration Plan

- Extend workspace load/save code to accept row-group metadata while preserving existing single-row targets.
- Ship frontend planner updates behind the same persisted contract so old workspaces load before operators create new row groups.
- Reject invalid shared-column applies at validation time; no destructive migration of old Step 4 data is required.

## Open Questions

- 是否需要在第一版就支援 operator-defined composite group key，或只支援固定的 row-group id + timestamp contract？
- runtime 端的資料表是否需要額外保留群組識別欄位，還是只用寫入時的列切分邏輯即可？
