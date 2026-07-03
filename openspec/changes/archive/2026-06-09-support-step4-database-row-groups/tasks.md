<!--
Each task description MUST state:
- the behavior or contract being delivered (what is observably true when the
  task is complete), and
- the verification target that proves completion (test, CLI invocation,
  analyzer check, manual assertion, or content review).

File paths are supporting context for locating the work, never the task
itself. "Edit file X" is not a valid task — it is missing both behavior and
verification.
-->

## 1. Workspace row-group contract

- [x] 1.1 Deliver `Requirement: Workspace database row groups capture shared-column intent` and `Requirement: Workspace persists database row-group plans` by extending the Studio V2 workspace database autosave/load contract so row-group metadata round-trips with legacy single-row targets, verified by backend workspace handler/service tests and frontend autosave hydration tests.
- [x] 1.2 Deliver `Decision: Persist row groups separately from DbTarget rows` by introducing a distinct row-group data shape and target reference model in workspace/database types and APIs, verified by type-level checks plus JSON contract assertions in `internal/api/handlers/studio_v2_workspace_database_handler_test.go`.

## 2. Step 4 planner and validation

- [x] 2.1 Deliver `Requirement: Step 4 row-group planner exposes shared-table structure` so operators can create, inspect, and edit row groups for one connector/table in Step 4, verified by `frontend/tests/unit/workbench-v2/step4-database.test.tsx` coverage for row-group rendering and editing.
- [x] 2.2 Deliver `Requirement: Step 4 shared-column validation distinguishes row-group conflicts`, `Requirement: Row groups define the legality of shared business columns`, and `Decision: Validate shared columns inside a row-group scope` so shared columns are legal inside one row group and blocking across groups, verified by planner validation tests and readiness-panel assertions.

## 3. Runtime delivery guardrails

- [x] 3.1 Deliver `Requirement: Row-group write mode guardrails protect unsafe upsert reuse` and `Decision: Limit shared-column support to insert-first delivery` so invalid shared-column upsert plans are blocked before apply, verified by backend readiness/runtime tests and Step 4 apply-blocker tests.
- [x] 3.2 Deliver insert-mode row-group database writes as distinct row instances under one table contract, verified by database delivery/service tests that show separate rows are emitted for separate row-group members without cross-row overwrites.
