# 來源與尚待驗證事項

基準：main @ `1a0311c8e8db9c62fe4388f0701ba38afe552ff7`，與上一輪相同；2026-09-15 透過 GitHub 重新查詢分支指向。所有現況判斷限於該版本，不引用搜尋索引的舊版本當最新程式。

## 已讀來源與可確認的現況

以下路徑皆從專案根目錄起算；查閱時使用上述固定 commit，而非隨時間移動的 main。

| 代號 | 來源與定位 | 能確認的現況 |
| --- | --- | --- |
| E01 | `internal/api/handlers/studio_v2_workspace_recording_plans_handler.go`：SchemaApply、TestWrite | 兩入口驗證格式後直接組合成功回覆，沒有呼叫實際建表或試寫服務。List 以工作區列出方案；Get／Delete 沒有在入口核對工作區，後續須完整驗證服務邊界，不能僅靠 UI。 |
| E02 | `internal/api/router_studio_v2_recording_routes.go`：registerStudioV2RecordingRoutes | 記錄方案服務存在時會註冊前述入口；沒有查現場部署是否啟用。 |
| E03 | `internal/datalink/recordingplan/service.go`：GenerateSchemaPreview、ApplyManagedSchema、ExecuteTestWrite | 有另一套服務實作；套用檢查期限但未見完整版本／目標核對；試寫服務使用 written_verified 等狀態，未見清理步驟。不能直接推論現行空入口已寫過資料。 |
| E04 | `frontend/src/features/datalink/workbench-v2/steps/step4/RecordingPlanSetupSection.tsx` | 連線與種類寫死、plans[0] 選取、成功只辨識 success 並承諾已清理；建立成員使用同一 equipment_id。 |
| E05 | `frontend/src/features/datalink/workbench-v2/steps/step4/Step4Database.tsx` | 傳第一台設備與全部點位 ID；使用 getColumnsFor(kind)；一般建表控制未接 isReadonly。 |
| E06 | `frontend/src/features/datalink/workbench-v2/state/dbSchemas.ts`：getColumnsFor | 取固定示範欄位，非依所選表查得。 |
| E07 | `frontend/src/features/datalink/workbench-v2/steps/step4/TargetMappingTable.tsx` | 目標更新會重建全部草稿；Enter 先提交再 blur；畫面提示將欄位描述為真表欄位。 |
| E08 | `frontend/src/features/datalink/workbench-v2/steps/step4/SchemaSetupSection.tsx` | 範圍變更只重設狀態，舊要求完成仍可寫入畫面。 |
| E09 | `frontend/src/services/datalink.ts`：dbTargetAPI.listTables | 已有資料表／欄位讀取入口，可優先沿用。 |
| E10 | `internal/api/handlers/studio_v2_workspace_database_handler.go`：UpdateConfig、UpsertTarget、GenerateSchema | 一般建表確實呼叫服務；本地連線／分組及對應／參照分步儲存，需要第二段失敗測試再確認實際一致性。 |
| E11 | `frontend/src/features/datalink/workbench-v2/steps/step4/step4DatabaseHelpers.ts`：createPoolConnectorPatch | 選取連線的密碼被 trim，存在變更有效密碼的風險。 |
| E12 | `frontend/src/hooks/datalink/useStudioV2WorkspaceRecordingPlans.ts` | 已有方案查詢及能力查詢 hook；能沿用，不必在元件另建平行 API 呼叫。 |

## 規格對照

- `openspec/specs/recording-database-setup/spec.md` 已要求真實能力、版本綁定、後端重用憑證與可辨識的試寫；本案是修正實作落差並補足情境，不宣稱另發明新能力。
- `openspec/specs/datalink-workbench-v2-step4-database/spec.md` 的前段已要求真實欄位，但 `Commit sequence and animation` 仍要求固定動畫及不呼叫後端就成功。本案以完整 MODIFIED 區段改為真實結果，不恢復示範行為。
- `openspec/specs/workspace-database-row-groups/spec.md` 允許合法同組共用欄位。本案修正 Step 4 全域重複欄位敘述以保留此例外，不重寫分組儲存語意。
- 已讀 `AGENTS.md`、`CLAUDE.md`、Studio inventory 的 `START_HERE.md`、`context.json`、`CURRENT_STATE.md`。inventory 中歷史缺口不一律當作目前程式仍未完成的證據；以程式與對應測試為準。
- `openspec/config.yaml` 指定 spec-driven；main 已改用 Spectra 管理相同 OpenSpec 目錄。`openspec/changes/` 本次讀取只見 archive，`.spectra` 在此 commit 不存在；未宣稱已查過所有未公開／本地 parked change。

## 信心與限制

E01、E04–E08 是已讀程式可確認的行為或條件。第二段儲存失敗造成何種持久狀態、所有權是否有其他層保護、特定資料庫能力與實際現場損害仍須驗證；沒有觀察到正式資料損毀，不能寫成已發生事故。

本次沒有重新執行上一輪修補包的局部測試，不沿用那些數字作本分支的驗證證據。上一輪下載包未包含在本次提交，也未視為已合併。

效能議題如逐規則查詢先量測再另排改善；本案不承諾加速百分比，不擴大為全倉庫最佳化。
