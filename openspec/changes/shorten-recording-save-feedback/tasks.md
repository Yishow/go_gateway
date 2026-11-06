## 1. 鎖定 ACK 與刷新界線

- [x] 1.1 在 reuse-recording-home-snapshots invalidation keys 完成後，先加入 database-save-feedback.test.tsx 失敗案例，驗證 Configuration save acknowledgement is separate from refresh；用 deferred promises 觀察有效 ACK、三個並行刷新、缺 scope/revision 的保守等待。
- [x] 1.2 [after: 1.1] 依保存與刷新分開結算新增 config-only scoped refresh coordinator，保留 response/reducer 契約與其他 helpers；focused test 證明 configuration saved 不等於 readback/activation/SQL/cleanup。
- [x] 1.3 [after: 1.2] 完成精確失效與舊回覆隔離，落實 Background refresh remains scoped and observable：接受 ACK 前失效、同 key 合併、跨 scope 不合併、older response 不覆寫；沿用 B home mutation capture/complete 並保留 bootstrap fallback owner 相容；deferred/race tests 覆蓋 refresh fail、新er edit、scope change、unmount。

## 2. 遷移單一保存流程

- [x] 2.1 [after: 1.3] 依受限遷移與等待契約只遷移 useUpdateStudioV2DatabaseConfigMutation 及 autosave caller，UI 分別顯示保存與刷新失敗／retry，en/zh-TW 文案完整；Step4SupportPanels 折疊摘要與 AdvancedWorkbenchV2Page 同範圍 retained shell 的 component tests 驗證 refresh error 不改成保存失敗也不標 fresh；malformed initial bootstrap 仍失敗。
- [x] 2.2 [after: 2.1] 驗證 Save acceleration preserves uncertain and consequential operation barriers：409、timeout/lost ACK、repeated click、Close/Cancel/Back/Forward、新編輯保留；database autosave/queue、activation barrier、reload recovery 與 scoped readback regression 證明其他 mutation 等待未改。
- [x] 2.3 [after: 2.2] 僅移除 named config flow 的串行 refresh await 並檢查 caller diff；合成 1ms ACK＋三個200ms刷新測試不再等600ms鏈，npm run build 與 query/error assertions 證明背景 rejection 可見且無 unhandled promise。

## 3. 保存與返回驗收

- [ ] 3.1 [after: 2.3] 最終 Mac browser 合併實測有效保存、失敗刷新、再次編輯、返回首頁及 unresolved write，分別記 ACK/refresh/readback 時間與實際狀態；全部 change 實作完成前保持本項未完成；不使用正式 DB，不擴權。
- [x] 3.2 [after: 2.3] 完成 config-only scope review、affected integration/activation tests、一次 frontend lint/test/build、git diff --check、make check-lines、Spectra verify/review，明列未跑的 browser／SQL／field 驗收，不擴張 write-group helper 範圍。

- [ ] 3.3 [after: 3.1, 3.2] 依最終 Mac 驗收與同範圍保留完成真實 browser 證據、必要修復與最終 verify/review，再執行 archive；雲端階段不提前勾選。

變更記錄：最新使用者指示替換 3.1 雲端 browser 與 3.2 browser 前置；驗收本身完整保留，新增 3.3 最終 archive gate。

任務切片以每項不超過一小時為規劃估算，不是工期承諾；若實作探索證明超過，先把同一行為拆成可獨立測試的垂直批次並保留 dependencies，不擴張 scope。未執行的 browser 與 archive checkbox 保持未完成。

## 4. 最終瀏覽器驗收準備

- [x] 4.1 [after: 2.1] 在雲端編寫 recording-database-save-feedback.spec.ts 與 database-save-feedback-api.ts，使用真實 /studio/v2 router 加模擬 API，保留合法／缺證據 ACK、三類獨立 held reads、折疊摘要錯誤與唯讀 retry、lost ACK 不重送，以及 Close／Back／Forward guard 的可重現驗收；只做 test discovery／static/type 檢查，不啟動 browser，實際執行仍由 3.1 驗收。

## 5. 獨立 review 修復

- [x] 5.1 [after: 2.1] 修復 review 確認的 config refresh recovery／owner epoch 邊界；新增回歸先證明失敗，恢復後唯讀觀察可清除過期 refresh error，A→B→A／origin 切換後舊 ACK 不作用於新 owner，保留 unknown 不重送與原操作等待；重跑受影響測試並由獨立 reviewer 複查，不改寫已完成任務歸屬。
