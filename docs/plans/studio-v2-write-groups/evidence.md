# 現況證據

基準：`4c00af81fb51b3239b4e0d0d9c6b83680c808b36`（2026-09-18 main）；2026-10-01逐項讀source。以下為source inspection，不是此輪執行產品測試。路徑相對repo根目錄，符號優先於容易漂移的行號。

| 發現 | 可追溯來源與符號 | 這次規劃如何使用 |
| --- | --- | --- |
| Step4 managed文案不要求手動mapping，但summary仍以enabledTargetCount===0封鎖 | `frontend/src/features/datalink/workbench-v2/steps/step4/RecordingPlanSetupSection.tsx`、`CommitSummary.tsx` 的 isSubmitDisabled | A/E統一authority，不再疊第二個模型的數量gate |
| 一般mapping建立Tag/Mapping而非MeasurementDefinition | `internal/api/handlers/studio_v2_workspace_mappings_handler.go` 的 Create/upsertTag | 基本Tag寫入不強制advanced semantics；不能拿pointID冒充measurementID |
| plan membership要求每個enabled point恰有一個persisted measurement | `frontend/src/features/datalink/workbench-v2/steps/step4/recordingPlanMembership.ts` 的 resolveRecordingPlanScope | A相容resolver/E使用範圍修復；不猜設備歸屬 |
| 既有plan全不匹配時option停用，新建form僅plans.length===0 | `RecordingPlanSetupSection.tsx` 的 selected/activePlan及plans.length分支 | E提供新建/修復/明確重選，區分empty與load error |
| Step1 allTested與reload devices.length完成條件不同 | `steps/step1/Step1Device.tsx` 的 allTested；`frontend/src/pages/datalink/workbench-v2/hydratedProgress.ts` 的 inferHydratedProgress | E共用persisted readiness，connection edit後重新驗證 |
| MQTT被列出但range parser只處理Modbus/FATEK/MC3E | `state/protocols.ts` 的 PROTOCOLS；`frontend/src/utils/addressParser.ts` 的 AddressParser.parse | E以實際V2能力交集限制選項，不在這批新增MQTT實作 |
| address conflict只以address當key，未包含device | `state/sourceRule.ts` 的 detectAddressConflicts | E測兩設備同address合法、同設備真重疊仍阻擋 |
| 實際writer已接production，可執行SQL | `cmd/test_ui/service_wiring.go` 的 wireGatewayServices；`target_writer.go` 的 newProductionTargetWriter；`internal/datalink/dbtarget/writer.go` 的 writeMapping/flushGroupedBucket | 不把現有輸出說成stub或全面不可用 |
| grouped writer以UTC truncate分桶、按到達覆寫Values；takeFlushableBuckets先delete，再Exec | `internal/datalink/dbtarget/writer.go` 的 bufferGroupedWrite/buildGroupedWriteKey/takeFlushableBuckets/flushBuckets | B補sample選值/closure；C補durability，不能只加retry timer |
| upsert SQL主要以timestamp column衝突，不提供跨group/entity的充分識別 | `internal/datalink/dbtarget/writer_statements.go` 的 buildGroupedWriteStatement/buildWriteStatement | B定義複合identity，A/C核對target unique capability |
| 既有SampleEnvelope、SQL journal/outbox及delivery worker可重用 | `internal/datalink/measurement/types.go`；`internal/datalink/delivery/types.go`、sql_journal.go、sql_outbox.go、worker.go | 元件不是production整合證據，C補實際wiring和crash tests |
| 主程式measurement/recordingPlan/history服務注入API，但runtime仍scheduler＋storage＋dbtarget/Share | `cmd/test_ui/main.go` 的 runGateway；`service_wiring.go` | 不把API可呼叫當成aggregation/durable delivery已連起來 |
| TestWrite明確501 | `internal/api/handlers/studio_v2_workspace_recording_plan_operations.go` 的 TestWrite/renderRecordingOperationNotImplemented | D完成前不移除安全封鎖 |
| 真metadata及revision scope已有實作 | `steps/step4/step4DatabaseHelpers.ts` 的 metadataColumnsForScope/isMetadataForScope；`useStep4TargetColumns.ts`；`internal/api/handlers/studio_v2_workspace_recording_plan_preview.go`、apply.go、targets.go | 沿用，不以sample欄位代替查詢失敗 |
| workspace row groups包含point membership、group/unique columns | `internal/datalink/workspace/service_database_row_groups.go` 的 DatabaseRowGroup/ReplaceDatabaseRowGroups | A保存legacy資料與輸出語意，不盲目把所有shared-column轉成單一wide row |
| 目前active tasks為5完成/10 total | `openspec/changes/fix-studio-v2-database-workflow/tasks.md`（此輪移交前） | 完成1.2/1.3/2.1/2.2/2.4保留；未完2.3/3.3/3.4/4.1/4.3唯一移交 |

縮寫path `steps/`、`state/`、`RecordingPlanSetupSection.tsx` 均在 `frontend/src/features/datalink/workbench-v2/` 下，表中完整path與符號可直接搜尋。實作tasks使用完整根相對path或同案design定位，不使用行號當唯一指引。

## 已存在的測試檔不等於此次測過

- `internal/datalink/dbtarget/writer_row_groups_test.go`、`delivery_outcome_test.go`
- `internal/datalink/delivery/worker_test.go`、`worker_failure_test.go`、`worker_receipt_failure_test.go`、`sql_storage_test.go`
- `frontend/tests/unit/workbench-v2/step4-database-row-groups.test.tsx`、`step1.test.tsx`

本次不更改產品程式，不宣稱Go、Vitest、Playwright、實際SQL或field tests通過。歷史inventory/release中的測試數僅對應其當時版本。

## 工具來源

目前tracked skill是Spectra；原OpenSpec skill在 `1a0311c8e8db9c62fe4388f0701ba38afe552ff7` 被刪除。本次經明確同意只讀其父版 `e8a7ba036c417b46c7c8b0ab885b9012460560b5:.agents/skills/openspec-propose/SKILL.md`，依該原版workflow搭配官方 `@fission-ai/openspec@1.2.0` CLI（MIT）於repo外隔離安裝。

沒有恢復舊skills、沒有跑init/update、沒有更改managed Spectra block。OpenSpec CLI提供new/status/instructions/validate；不把自製lint當成官方驗證。
