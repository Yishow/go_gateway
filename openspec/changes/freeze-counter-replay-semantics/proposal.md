## Problem

未縮放的累積記錄在重啟重播時，可能把已接受的舊讀值套用到新的量測 epoch／單位。基準 `913cead9` 的重現已觀察到：原本 epoch-1／kWh 的 100、105，量測 PUT 成功改為 epoch-2／MWh 且 applied snapshot 未變後，重播輸出 delta 5 MWh；這不是合法單位轉換，也不能用 NULL 或刪除舊 journal 掩蓋。

## Root Cause

`internal/datalink/grouppipeline/reconcile.go` 的 `buildFrozen` 只對既有 `ExactRecordingSource` 使用 immutable proof，其他 counter member 仍重新查詢可編輯的量測定義。`internal/datalink/groupdelivery/sample_codec.go` 的 sample payload 保留原值與 source/mapping revision，卻未證明完整 counter binding；`assembleBoundary` 因而可能以新 binding 重播舊樣本。既有已縮放的 immutable proof 路徑應保留，不將本次問題誤稱為所有 exact source 都壞掉。

## Proposed Solution

- 在 applied runtime snapshot 中凍結所有 cumulative／counter_raw member 的完整 counter binding，接受樣本前持久化可驗證的 proof identity；新流程與 journal 的 accepted set 必須同一語意版本。
- 冷啟動、reconcile、未關窗 journal replay 與已關窗 outbox retry 均只讀 accepted revision 的 immutable facts，不再用目前 measurement definition 補齊歷史。
- 舊資料僅在已有 immutable 證據可證明原 binding 時遷移；無證據、矛盾或版本不支援時保留 journal 並隔離該 scope，提供安全原因與修復動作，不自動重新貼單位、不清空 baseline、不隱性刪除。
- 此 change 是 `manage-recording-lifecycle-safely` 的前置；它只修正既有語意凍結，新增轉換矩陣與 successor generation 由生命週期 change 擁有。與 `fix-recording-cutover-boundaries` 共用 runtime 檔案，實作順序先本 change 再 cutover 修正。

## Success Criteria

- 回歸 fixture 中 100 → 105 在重啟前後仍屬 epoch-1／kWh 且 delta 5；新 epoch 的後續明確 apply 從 NULL delta／baseline 開始，不能接續舊 baseline。
- 量測被編輯、刪除、暫時不可讀，或進程在 acceptance／closure 前後終止，都不改變已接受樣本的值、單位、epoch、binding、counter policy 與 frozen destination。
- 缺乏 immutable proof 的舊 journal 被安全阻擋、仍可查明受影響群組與修復方法；其他群組繼續，已 finalized payload 可照原目的地交付，不要求重算。
- 每個行為先新增會失敗的測試，再完成最小修復及重構；本次只交付 proposal，產品驗收尚未執行。

## Impact

- Affected specs: `durable-recording-delivery`；新增針對 counter replay proof／legacy recovery 的精確要求，補足既有 `Revision-safe delivery and apply` 與 `Recovery and retention protect evidence`，不複製或弱化原要求。
- Affected code:
  - Modified: `internal/datalink/grouppipeline/reconcile.go`, `internal/datalink/grouppipeline/counter_sources.go`, `internal/datalink/grouppipeline/recovery.go`, `internal/datalink/runtime/group_boundary_build.go`, `internal/datalink/groupdelivery/sample_codec.go`, `internal/datalink/workspace/write_group_runtime_layout.go`。
  - New: `internal/datalink/grouppipeline/counter_binding_replay_test.go`, `internal/datalink/groupdelivery/counter_binding_proof_test.go`, `internal/api/router_studio_v2_counter_replay_test.go`；必要的 proof model／internal migration 依既有 schema 編號續增，不能覆寫已部署 migration。
  - Removed: 無產品檔案；完成相容遷移後僅移除以 current definition 補歷史 proof 的 fallback。
- 保留既有 counter 算法、exact encoding、receipt／outbox identity、workspace revision／readiness／schema confirmation；不增加對外資料庫操作授權。
