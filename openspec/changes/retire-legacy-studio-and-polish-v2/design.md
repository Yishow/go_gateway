## Context

`/studio` 已經不再使用，使用者已明確要求立即刪除並接受以 Git rollback 復原。現有 inventory、route registry、lazy imports、測試與文件仍把 legacy workspace 當成產品 surface；若只移除畫面而沒有盤點 route/asset/import 邊界，會留下不可見的 legacy chunk 或特殊 redirect。

本 change 分成兩個互相獨立的 batch：第一個是立即可執行的 legacy deletion batch；第二個是保留的 C release-polish batch（typed errors、preview/runtime truthfulness、Share accessibility、CommitProgress presentation 與正式文件）。legacy deletion 不等待 usage window、telemetry、A/B evidence 或 Modbus Share lifecycle；Modbus Share 的 B contract 仍由它自己的 change 與驗收負責，本 change 不修改或弱化它。

## Goals / Non-Goals

**Goals:**

- 先建立可審查的 pre-delete inventory，再立即移除專用 `/studio` route registration、proven legacy-only imports/chunks/assets、obsolete legacy tests 與 obsolete docs references。
- 讓刪除後的 `/studio` 與任意 unknown route 完全使用同一個既有 generic unknown-route policy；不新增 `/studio` tombstone、special redirect、fallback 或 handler。
- 維持 `/studio/v2`、`/studio/runtime`、`/test`、`/gateway/*` 的 route identity、lazy assets 與可用性。
- 以一個獨立 Git commit 作為 rollback anchor，驗證 revert、rebuild、redeploy 路徑與 rollback 限制。
- 讓 C release-polish 能在 legacy deletion 之後獨立實作與驗收，不把 C UI 行為誤認為 B backend CAS/concurrency。
- 讓 API、inventory、release evidence 反映 immediate deletion、unknown-route equivalence、preserved surfaces 與 rollback procedure。

**Non-Goals:**

- 不建立或保留 legacy deprecation notice、telemetry table/API/config、usage counter、14-day gate、retention window、query migration、307 redirect window 或 migration cleanup。
- 不新增、更新、刪除或回寫任何 database/runtime deployment data；Git rollback 不承諾恢復 deployment/runtime data 或 browser bookmarks。
- 不等待或要求 A/B acceptance 才能執行 legacy deletion；不改變 Modbus Share 的獨立 A/B requirements、CAS、concurrency、ownership 或 lifecycle semantics。
- 不移除、redirect 或重新解釋 `/studio/v2`、`/studio/runtime`、`/test`、`/gateway/*`。
- 不清理與本 change 無關的歷史 Swagger 缺口、資料 migration 或產品重構。

## Decisions

### Decision 1: Execute immediate legacy deletion as an independently bounded batch

立即 batch 的順序固定為：建立 inventory → 審查 proven legacy-only 邊界 → 移除 route/import/chunk/asset/test/docs → 建立 rollback anchor → rebuild/route/browser/static smoke。這個 batch 不包含 deprecation phase、usage wait、telemetry、redirect stage 或 data migration；使用者的 owner decision 直接授權刪除。

完成條件必須是 `/studio` 不再有專用 route registration，且沒有可由 `/studio` 到達的 legacy workspace。若 inventory 無法證明某個 import、chunk、asset 或 test 為 legacy-only，該項必須保留並記錄原因，不得猜測刪除。

### Decision 2: Make `/studio` exactly follow the existing generic unknown-route policy

刪除 `/studio` 後，frontend route table、backend/static route table 與 embedded asset manifest 均不得保留 `/studio` 專用項目。`/studio` 和一個同時期建立的 arbitrary unknown path（例如 `/unknown-route-for-retirement`）必須產生相同的 response/navigation、status/fallback、network asset pattern 與 error boundary 行為。若既有 generic policy 將 unknown path fallback 到 `/studio/v2`，`/studio` 也只能經由該共同 policy 到達；若既有 policy 回 404，兩者都回 404。

不得建立 `/studio` tombstone、專用 404/410 response、special 307/308 redirect、legacy-specific middleware、legacy mount 或只為 `/studio` 保留的 compatibility flag。驗收要做 side-by-side unknown-route comparison，而不是只檢查單一 status code。

### Decision 3: Remove only proven legacy-only frontend graph members

刪除範圍由 pre-delete inventory 驅動：route registration、direct/lazy imports、legacy-only chunks、legacy-only CSS/assets、legacy-only unit/integration/e2e tests，以及只描述 legacy workspace 的 docs references。任何被 `/studio/v2`、`/studio/runtime`、`/test`、`/gateway/*` 或 shared components 使用的檔案必須保留。

static embedded delivery 必須重新生成並檢查完整 asset graph；若 build 或 graph check 發現 preserved surface 依賴該檔案，刪除 batch 必須停止並縮小範圍，不得用 fallback 掩蓋 missing asset。

### Decision 4: Use a dedicated Git rollback anchor with explicit recovery limits

legacy deletion 必須形成獨立 commit，commit message 與 release evidence 指向 pre-delete inventory、刪除清單、驗證結果與 rollback command。復原程序是 revert 該 commit、重新 build、重新部署並重跑 route/static/browser smoke；不得以手動複製未追蹤檔案取代 Git rollback。

rollback evidence 必須明確區分：Git 可恢復 tracked source/docs 與 route/assets 定義；Git 不會恢復 deployment/runtime data、外部資料庫內容、已發出的設備狀態或 browser bookmarks。若資料或 bookmark 需要復原，必須另走 deployment/data owner 流程，本 change 不提供自動復原。

### Decision 5: Keep C release polish independent from legacy deletion and Modbus Share

C batch 保留 typed safe-error envelope、en/zh-TW mapping、Step 3 preview/SSE truthful states、runtime degraded/error states、Share keyboard/accessibility/single-flight、CommitProgress actual result states 與 docs evidence。C batch 不依賴 legacy deletion rollback，也不新增 B Share CAS/concurrency/ownership semantics。legacy deletion 可先完成；C batch 可在後續獨立 apply/驗收。

### Decision 6: Make route/docs/inventory changes evidence-producing

Go annotations 產生 `docs/swagger/docs.go`、`docs/swagger/swagger.yaml`、`docs/swagger/swagger.json`；API registry、runtime/inventory docs 與 release note 必須描述刪除後 generic unknown-route equivalence、preserved route identities、pre-delete inventory 與 Git rollback limitations。修改 `docs/technical/studio-surface-inventory/` 內文件時，使用 `go run ./cmd/studio_inventory_changelog add` 留下 summary、surface、files、reason；不得新增 telemetry 或 data migration 記錄。

## Implementation Contract

### Observable behavior

- 在 immediate deletion batch 完成後，`/studio` 不再匹配任何 dedicated frontend/backend route、handler、tombstone 或 special redirect。
- `/studio` 與 arbitrary unknown path 必須按照同一 generic unknown-route policy 產生相同的 status/navigation/fallback、error boundary 與 network asset pattern。
- `/studio` 不得 render legacy workspace、載入 legacy-only chunk/asset，亦不得透過隱藏 lazy import 重新引入 legacy surface。
- `/studio/v2`、`/studio/runtime`、`/test`、`/gateway/*` 的 route identity、主要互動與 required assets 保持可用。
- rollback 是獨立 Git commit revert + rebuild/redeploy；release evidence 清楚說明 Git 不恢復 runtime data 或 browser bookmarks。

### Interface / data shape

- pre-delete inventory 至少包含 route registration、direct/lazy import、chunk/asset、test、docs reference、owner decision、preserved-surface dependency 與 deletion decision。
- rollback record 至少包含 deletion commit、revert command、rebuild command、redeploy target、route/static/browser smoke result、source/data/bookmark limitation。
- 本 change 不新增 telemetry event、telemetry reader/writer、database table、database migration、runtime config 或 query-migration payload。
- C batch 的 safe error envelope 必須含 `error.code`、`error.message`、`error.retryable`、`error.request_id`；這是獨立 UI/runtime contract，與 legacy deletion 無關。

### Failure modes

- inventory 無法證明 legacy-only、preserved surface 依賴待刪檔案、route/asset smoke 失敗、或 unknown-route comparison 不相同：立即 deletion batch fail closed，停止刪除或 revert deletion commit。
- `/studio` 仍可匹配 dedicated route/handler/tombstone、special redirect、legacy mount 或 legacy chunk：視為 deletion failure，不得以 `/studio/v2` 可用來掩蓋。
- Git revert 成功但 deployment/runtime data 或 browser bookmarks 未恢復：回報為 rollback scope limitation，不能宣稱完整資料復原。
- C batch 的 typed error、preview/SSE、runtime 或 Share 驗證失敗：只阻擋 C batch，不回退或擴大 legacy deletion範圍。

### Acceptance criteria

- immediate deletion preflight inventory 通過內容審查，並列出所有刪除與保留項目；legacy-only 證據不足的檔案不被刪除。
- frontend/backend route tests 證明沒有 `/studio` dedicated registration/handler/tombstone/special redirect，且 `/studio` 與 arbitrary unknown route 的 response/navigation/status/fallback 等價。
- source/import reachability、chunk manifest、static tree 與 browser network smoke 證明無 legacy workspace markup、legacy-only chunk/asset/request；`/studio/v2`、`/studio/runtime`、`/test`、`/gateway/*` 均通過 identity/asset smoke。
- deletion commit、revert、rebuild/redeploy 與 rollback smoke evidence 完整；文件明確記錄 Git/source/docs、deployment/runtime data、browser bookmarks 的復原邊界。
- C batch focused tests 覆蓋 typed error i18n、preview state/backoff/no-silent-fallback、SSE recovery、runtime degraded/error、Share accessibility/single-flight、CommitProgress failed presentation；C 不新增 B CAS/concurrency tests。
- Swagger generation、API registry、runtime/inventory docs、release note、inventory changelog readback、`spectra analyze`、`spectra validate` 與 `git diff --check` 通過；GUI、external service、PLC、field acceptance 分開回報。

### Scope boundaries

**In scope:** immediate `/studio` deletion, pre-delete inventory, generic unknown-route equivalence, legacy-only graph reachability, preserved route/asset smoke, Git rollback/rebuild/redeploy evidence, C typed error/preview/runtime/accessibility/progress polish, and formal docs/inventory/release evidence.

**Out of scope:** deprecation/telemetry/14-day gate/retention/query migration/307 redirect/data migration, database/runtime data mutation, `/test` removal, `/gateway/*` removal, A/B gating of legacy deletion, Modbus Share implementation, unrelated Swagger debt, production deployment execution by this ingest, commit/push by this agent, and field acceptance.

## Risks / Trade-offs

- [Risk] An unrecorded shared import is mistaken for legacy-only. -> Mitigation: require import/reachability inventory and preserve any file with a non-legacy consumer; fail closed when ownership is uncertain.
- [Risk] Generic unknown-route behavior changes while `/studio` is removed. -> Mitigation: capture arbitrary unknown-path baseline before deletion and compare response/navigation, status/fallback, error boundary, network assets, and preserved route identity after deletion.
- [Risk] Git revert restores source but not deployed/runtime state or bookmarks. -> Mitigation: make rebuild/redeploy and data/bookmark limitations explicit in the rollback runbook and release evidence.
- [Risk] C polish work expands into B concurrency or legacy compatibility. -> Mitigation: keep C tasks independent, reject B CAS/concurrency changes, and do not add any legacy telemetry/deprecation/migration seam.

## Rollout Plan

1. Generate and review the pre-delete inventory, capture generic unknown-route and preserved-surface baselines, and record the owner decision for immediate deletion.
2. Apply the bounded legacy deletion batch, create its dedicated Git commit, rebuild the embedded frontend, and run route/import/chunk/static/browser smoke.
3. Publish deletion and rollback evidence. If rollback is required, revert the dedicated commit, rebuild/redeploy, and rerun the same smoke; report runtime data and bookmark limitations separately.
4. Apply the C release-polish batch independently when its focused scope is approved; it does not reopen `/studio` or add telemetry/migration work.

## Open Questions

There are no unresolved contract questions for immediate apply. Any request to reintroduce `/studio`, add a special redirect/tombstone, add telemetry/data migration, wait for usage evidence, or make deletion depend on A/B acceptance requires a new explicit ingest decision.
