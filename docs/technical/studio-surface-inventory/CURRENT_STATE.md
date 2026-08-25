# Studio Surface Inventory: Current State

Last updated: `2026-08-25`

## One-screen Summary

- `/studio/v2`
  - 目前主重點。
  - 目標是面向使用者、簡單、有指引、好觀察。
  - `/`、unknown routes 與 generic `/datalink` landing 已收斂到 `/studio/v2`。
  - 下一步重點是把多台設備送出真正接上後端。
- `/studio/runtime`
  - 追隨 `/studio/v2`。
  - 角色是 post-setup focused monitor，不是獨立 fleet-first dashboard。
  - 還需要承接 V2 的多台設備送出結果。
- `/test`
  - 獨立工程測試工具入口。
  - legacy `/templates`、`/history`、`/compare`、`/analyzer` 已收斂到 `/test`。
- `/gateway/*`
  - experimental surfaces。
  - 目前不是產品主線。
- `/studio`
  - legacy 完整工作台已獲 owner 授權立即刪除，不再保留 dedicated route。
  - 刪除後與任意 unknown route 共用既有 generic unknown-route policy。
  - pre-delete route/import/test inventory、保留邊界與 rollback 限制見 [retirement record](../../releases/retire-legacy-studio-and-polish-v2.md)。

## What Was Added Recently

1. `studio-surface-inventory` 現在有固定 onboarding 入口：
   - [START_HERE.md](./START_HERE.md)
   - [context.json](./context.json)
   - [CURRENT_STATE.md](./CURRENT_STATE.md)
2. `studio-surface-inventory` 現在有 SQLite changelog：
   - DB: `docs/technical/studio-surface-inventory/changelog.sqlite`
   - Tool: `go run ./cmd/studio_inventory_changelog ...`
3. `AGENTS.md` 與 `CLAUDE.md` 已明確要求：
   - 任務若涉及 `studio` surfaces / `studio-surface-inventory`
   - 先讀 onboarding 入口
   - 不要一開始就掃完整 md/html
4. pre-delete source review 已確認：
   - `frontend/src/App.tsx` 的 `/` 會 redirect 到 `/studio/v2`
   - unknown routes fallback 也會回到 `/studio/v2`
   - pre-delete 的 `/studio` dedicated route 只由 `frontend/src/App.tsx` 掛載 legacy workbench lazy import；後端 source scan 未發現 exact `/studio` dedicated handler
   - `frontend/src/pages/datalink/workbench/` 目前由 production source 只透過上述 lazy import 觸達；刪除候選仍須由 build/import graph 最終確認
   - `/studio/v2`、`/studio/runtime`、`/test`、`/gateway/*` 是本批明確保留 surface
   - `/test` 維持獨立頁而非掛回產品主流程
5. review17 implementation/evidence sync 已完成：
   - Startup 先從 durable `system_settings` 載入 Modbus Share 設定但不 bind listener；接著以 persisted workspace revision 從 source-rule candidates 還原 canonical projection，只有 `hydration_state=ready`、`readiness=true`、非空 workspace/settings revisions 與 token 齊全時才允許 listener running。設定讀取、workspace/revision、projection restore 任一步失敗都維持 failed/stopped，不以 display defaults 或 5020 fallback 冒充 ready；rollback 不確定時先 fail closed，即使 `MarkDirty` 也失敗仍不得恢復 listener/readiness。
   - Readiness、canonical projection、`runtime_apply_*` 與 connector diagnostics 使用固定、可行動且不含 endpoint/credential/transport detail 的安全訊息；raw backend exception 不進正常 operator DOM。
   - Runtime SSE 的 unavailable/error 事件帶 typed code、retryable、opaque request ID 與 action；`/studio/runtime` 保留最後成功 snapshot、退回 5 秒 polling，Reconnect action 會實際重建 EventSource。Preview SSE 的重連次數有上限，explicit disconnect 會取消計時器與重連。
   - Share operations remain fail-closed on ready hydration, durable workspace ownership, global listener settings, matching workspace/settings revisions, and readiness tokens；reconcile 先取 workspace lock，再以 process-global projection transaction lock 保護 global memory bank 的 snapshot/swap/rollback。這會序列化跨 workspace projection、以 correctness/rollback isolation 換取部分 throughput；validation/CAS 仍在 global lock 前完成以縮短持鎖時間。
   - Preview SSE requires `mapping_id` + `workspace_id` ownership proof before subscription；runtime/preview events carry typed code, `retryable`, opaque `request_id`, and actionable `action` metadata，en/zh-TW 只顯示安全 operator copy，raw diagnostics 不進正常 DOM。Generated Swagger records direct delete `422` and preview `workspace_id` scope。
   - B10 EXE/TCP current PASS pair 為 `155805Z` 與 `155901Z`，兩次 13/13；focused Go 22/22、Python 16/16、`py_compile` PASS，SHA `8398ffdf4b8977848c7dffd22a7fd3df367643b2867df795d3a5f61c3d0b2de7`。155314Z failure 與更早 pair 僅作歷史對照；155314Z witness 已依 retention 刪除，不建立 broken link。External/embedded Chromium、Windows、LAN、真 PLC、SCADA、field sign-off 與 deployment-owner rollback/rebuild/redeploy 仍 pending。
   - Prior working-tree frontend Vitest evidence was `134` files / `722` tests; after the contract-integration changes this docs refresh does not claim a current frontend full-suite count. `scripts/check_file_lines.sh` coverage 已同步納入 `*.py`，並略過 `vendor/`、`cache/`、`generated/` trees；current `make check-lines` passes, with the two modified locale JSON files at `500` lines (warning threshold only). 這是 local test/line-gate evidence，不代表 external/embedded browser 或 field acceptance。
6. review18 contract/docs closeout 已完成：
   - `POST /api/v1/datalink/studio-v2/workspace/activate` 的 JSON body、`workspace_revision`、`settings_revision`、`readiness_token` 現在同時由 Go binding 與三份 generated Swagger 宣告為 required；missing body 為 typed 422，missing required field 為 typed 400，既有 pending/error 422 與 stale-revision 409 保持不變。
   - 三份 Swagger artifact 各自只保留一份 activation required metadata，並同步補上 runtime status 的 `projection_code`；router-owned Modbus Share activation messages/actions 已統一小寫。review18 未修改 frontend、B10 harness、witness、commit 或 archive。
7. review20 Modbus Share contract/docs closeout 已完成：
   - Start/Stop、workspace-scoped mapping query/mutations 與 reconcile 的 Swagger annotations 依實際 typed gate/handler paths 補齊可達 `400`/`403`/`422`/`500`/`503` responses；不可達 status 未加入。`internal/api/modbus_share_swagger_contract_test.go` 以 generated JSON 精確鎖定各 operation response matrix。
   - `docs/swagger/docs.go`、`docs/swagger/swagger.json`、`docs/swagger/swagger.yaml` 已由 current annotations 重產；獨立 temporary `swagger` output 逐檔 byte-for-byte compare 為 drift 0。review20 未修改 frontend、B10 harness、witness、commit 或 archive。
   - line-limit warning files (>300) 的精準理由與後續拆分計畫已記錄於 release handoff；本輪沒有假稱已拆分。External/embedded Chromium、Windows、LAN、真 PLC、SCADA、field sign-off 與 deployment-owner rollback/rebuild/redeploy 仍 pending。
8. review22 Modbus Share contract refresh 已完成：
   - Status 現在以 hydrated workspace scope 回傳 `workspace_id`、`workspace_revision`、`settings_revision`、`listener_state` 與 workspace-scoped `mapping_count`；canonical desired mappings 帶完整 geometry 與 persisted `ownership_proof`，dirty/uncertain recovery 會保持 listener/readiness fail-closed。
   - Candidate review/apply/recompute requests 以 `workspace_id`、`expected_workspace_revision`、`revision_id` 為共同 gate；Step 4 讀取真正 persisted candidate snapshot 的 status/blocking reason/geometry，不以 browser-only candidate fields 冒充 server truth。Stride 小於 datatype span、foreign workspace/device、stale revision 與 unknown ownership 均在 mutation 前拒絕。
   - Reconcile outcome、removed/invalidated spans、safe diagnostics 已接至 Local Modbus review/summary；`dirty_unknown`、`invalidated_unknown` 與 candidate failure 不會顯示為 applied/ready。
   - 三份 generated Swagger artifact 已重新產生並納入 `modbusshare.OwnershipProof`；generated JSON/YAML 與獨立 temporary regeneration byte-for-byte 一致，`internal/api/modbus_share_swagger_contract_test.go` 通過。
   - B10 current pair 為 155805Z/155901Z（13/13、focused Go 22/22、Python 16/16、`py_compile` PASS）；155314Z failure 與更早 pair 僅作歷史對照，且 155314Z witness 已依 retention 刪除。

## Read This First Next Time

如果新對話任務涉及以下任一項：

- `/studio`
- `/studio/v2`
- `/studio/runtime`
- `/test`
- `/gateway/*`
- `docs/technical/studio-surface-inventory/`

先讀：

1. `AGENTS.md`
2. `CLAUDE.md`
3. [START_HERE.md](./START_HERE.md)
4. [context.json](./context.json)
5. [CURRENT_STATE.md](./CURRENT_STATE.md)

只有在這五個入口不足以回答問題時，才展開：

- [studio-mainline.md](./studio-mainline.md)
- [studio-v2-runtime.md](./studio-v2-runtime.md)
- [test-tooling.md](./test-tooling.md)
- [gateway-experiments.md](./gateway-experiments.md)
- [backend-api-registry.md](./backend-api-registry.md)
- [gap-roadmap.md](./gap-roadmap.md)
- [index.html](./index.html)

## What Not To Re-Do

- 不要每次重新盤點全部 surface 與全部 API。
- 不要重新建立 `/studio` route 或把其 pre-delete inventory 當成產品主線；`/gateway/*` 仍只保留記錄與 prototype 說明。
- 不要把 `/test` 誤當成產品主流程；它是獨立工程測試工具。
- 不要只改 inventory 文件而不寫 changelog。
- 不要直接跳進完整 HTML/長文檔，先用 onboarding 入口縮小範圍。

## Likely Next Work

最可能的下一批工作仍是：

1. `/studio/v2` 多台設備送出契約真正持久化到後端
2. `/studio/runtime` 多台設備交接與 lifecycle contract
3. 完成 legacy `/studio` dedicated route/import/test/asset deletion 的 route、chunk、browser smoke 與 rollback evidence
4. inventory 文件持續和實作同步

## Fast Commands

```bash
go run ./cmd/studio_inventory_changelog list -limit 10
go run ./cmd/studio_inventory_changelog add \
  -surface "overview" \
  -summary "..." \
  -files "file1,file2" \
  -reason "..."
```
