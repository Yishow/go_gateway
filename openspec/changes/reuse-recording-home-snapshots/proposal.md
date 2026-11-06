## Summary

以 workspace-scoped compact source observation 減少首頁全量資料搬運，並讓返回首頁保留已觀察畫面，同時重新驗證來源。

## Motivation

基準 main 0f54740 的本機測試量得首頁首次 14 個核心讀取、返回 12 個；合成 20 設備／2,000 點 fixture 的資料重播約 3.06 MB、median 577 ms。這些是雲端非瀏覽器量測，不是現場容量保證。現有兩次來源比對是安全檢查，不能直接刪掉或只延長 staleTime。

## Proposed Solution

- 新增唯讀、精簡、具完整性狀態與 source digest 的首頁來源觀察 API，重用 canonical repositories 與兩次投影內容比對。
- 不宣稱原子快照，不用 Share/workspace revision 代替來源完整性，不新增 persisted authority、command token 或全域健康索引。
- 首頁遷移至單一摘要讀取，完整 point objects 留給明確來源編輯；同 scope 返回先顯示 observed 畫面並重新驗證。
- 保留每次 visit 的取消／late response fence、filter/page/scroll/focus 與可見頁 bounded evidence demand。

## Impact

- Affected specs: simple-recording-home
- Affected code:
  - Modified: internal/api/router.go, frontend/src/hooks/datalink/keys.ts, frontend/src/hooks/datalink/useStudioV2WorkspaceDevices.ts, frontend/src/hooks/datalink/useStudioV2Rules.ts, frontend/src/hooks/datalink/useStudioV2Mappings.ts, frontend/src/hooks/datalink/usePoints.ts, docs/swagger/docs.go, docs/swagger/swagger.json, docs/swagger/swagger.yaml, frontend/src/features/datalink/workbench-v2/recording-editor/recordingSourceDiscovery.ts, frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeRoute.tsx, frontend/src/features/datalink/workbench-v2/recording-home/recordingHomeCachedReads.ts, frontend/src/features/datalink/workbench-v2/recording-home/recordingHomePresentation.ts
  - New: internal/datalink/workspace/recording_home_sources.go, internal/datalink/workspace/recording_home_sources_test.go, internal/api/handlers/studio_v2_recording_home_sources_handler.go, internal/api/handlers/studio_v2_recording_home_sources_handler_test.go, frontend/src/services/recordingHomeSources.ts, frontend/tests/unit/workbench-v2/recording-home-source-summary.test.tsx
  - Removed: 僅完成遷移驗證後移除首頁舊重複 transport；既有完整來源 API 不刪除。
- 延續 simplify-recording-home-guidance 與 adopt-simple-recording-entry 的唯讀／scope 契約；不改原 archived 任務、記憶體 views 或採集 runtime。

- 來源 mutation hooks 只增加相關摘要失效，不改既有 mutation await／readback／recovery；workspace accessor 只委派現有 Repository.Get，不能 GetOrCreate。Swagger 只由 annotations 重產。

- 最終 ownership 校正：舊首頁雙輪 adapter 已移至上述 editor 路徑，GroupEditor discovery 行為不變；首頁僅使用 compact service。刪除端點與新增端點均納入 scoped commit，不以舊路徑作現行 source anchor。
