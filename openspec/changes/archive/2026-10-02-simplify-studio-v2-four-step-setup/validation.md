# 實作驗證紀錄

2026-10-02：依總覽 A→F 順序實作 E。tasks 10/10 完成；3.2／3.3 的實測是**元件層級**（見下），完整 production 流程的證據留給 F。

**歸屬說明**：本案實作過程中，工作樹的一部分（`writeGroup/*` 的初版、`Step4Database`／`CommitSummary` 的改寫、舊 target／recording-plan 元件的刪除）在我寫完後、本案 verify 之前就已被其他提交收進 git（`77baeae2`…`b80b0431`），所以 `git diff` 看不到它們；它們仍是本案 2.1–2.4 的實作內容，驗證以現況檔案與測試為準。目前仍屬本案未提交的新增／修改：Step 1／2／3 的改動、`protocolSupport`、`GroupDeliveryStrip`、`GroupColumnProposal`／`proposal.ts`、locale 與相關測試。

## 後端補充（為 2.4 先定 API）
- `POST /studio-v2/workspace/write-groups/:id/apply`：只收 expected revisions，不收 group payload（422）；readiness 未通過 409 `WRITE_GROUP_NOT_READY`；過期 409 `revision_mismatch`。測試：`TestNewRouter_GroupApplyRequiresReadinessAndCurrentRevisions`、production 接線 HTTP `TestProductionGroupApplyOverHTTP`（真實 SQLite 目的地：applied_revision 等於 revision，重送 409）。Swagger 已重新產生。
- 限制：Apply 不同步開始寫入；`grouppipeline` 週期 reconcile 後才接手，且套用在下一個 UTC bucket 生效。

## 1.1 CapabilityAndReloadReadiness／OfflineDraftNavigation
- `protocolSupport.getProtocolSupport` 由 `AddressParser` 能否解析該協議的預設位址決定可用性（測試逐一比對，選項不會超前實際路徑）；MQTT 顯示不可用與原因，點擊不會選取，已存的 MQTT 設備仍顯示。
- `updateDeviceConfig` 改變任何連線設定即使舊 probe 失效（status 回 draft、test 清空）；相同值或改名不失效。
- Step 1 的 Continue 只在沒有設備時 disabled；未驗證設備顯示「尚未驗證」chip，仍可儲存草稿與 Next；reload（`inferHydratedProgress`）與初次流程一致，都把 Step 1 視為完成但未驗證。啟用與否仍由後端 readiness 把關。
- 有意改寫的舊測試：Step 1 的「全部測試過才可 continue」兩處、兩處「點 MQTT 卡片會切換協議」。
- 沒有做：後端 capability endpoint 供前端讀取（目前以前端 parser 為準，collector／setup 路徑由後端拒絕）。

## 1.2 DeviceScopedAddressConflict
- 衝突鍵改為 `device|address`；同一設備、同一 area 的值寬度重疊也算衝突（`40001 width2` 與 `40002`），不同設備同址合法；字母區（D100）、無法正規化的位址（逐字比較）、skipped／disabled 不計。Step 2 合併表列出「設備·位址·規則」定位。測試 `address-conflicts.test.ts`（6 項，先紅後綠）。

## 1.3 BasicPointTagToGroup
- Step 3 文案聚焦名稱／型別／換算與真實值對照；量測語意範本收為預設關閉的「進階」`<details>`。基本 mapping 保持原有 stable IDs 與 IME／Enter／blur 行為（未改動）。限制：沒有新增 Step 3 專屬的元件測試，只驗證進階區塊預設收合。

## 2.1–2.4／3.1 Step 4 寫入群組編輯器
- 新增 `writeGroup/*`：`WriteGroupSection`（清單、前置條件、錯誤／空／不符）、`GroupEditor`、`GroupMemberTable`（搜尋、只看問題、批次操作標明範圍與數量）、`GroupReadinessPanel`、`GroupLifecycleBar`（Apply／Disable／Delete，刪除前顯示留下的 backlog）、`GroupTestWritePanel`（D 的 preview→confirm→status，遺失回覆只查詢不重送，過期範圍的回覆丟棄）、`GroupDeliveryStrip`（saved／applied／collecting／本地等待／需處理／資料庫已確認分開，讀不到就是「未確認」不是 0）。
- 純邏輯：`candidates`（只收已儲存的設備／點位／Tag 的 persisted IDs，並說明被排除的原因）、`columns`（真實欄位型別相容性、只以 Tag 名稱＋相容欄位提出可審查的建議，沒有索引回退／環繞／範例欄位、不重複建議、不建議主鍵；已確認的指定欄位消失或不相容標為修復而不重綁；同列共用欄位僅在每個成員有各自不同的 entity key 時合法）、`draft`（與後端同款的本地檢查，後端 readiness 才是 Apply 的依據）。
- `Step4Database` 只剩目的地連線／資料表設定、群組編輯器、Share 摘要與啟用；`CommitSummary` 只依後端 readiness 決定能否啟用，不再用舊 `enabledTargetCount`，並如實顯示「無群組／已儲存未套用／已套用 n/m」。
- **刪除**：舊的 `TargetMappingTable／Row`、`RowGroupPlanner`、`DestinationOverviewCard`、`TargetMetadataStatus`、`RecordingPlanSetupSection`、`WorkspaceRecordingPlanSetupSection`、`recordingPlanMembership`、`useRecordingPlanActions`，以及只測它們的測試（TargetMappingTable、row-groups、target-metadata、RecordingPlan*、recording-errors integration、managed-table-prefix）。這是 spec 要求「Step 4 只編輯 write-groups」的結果；舊 plan 路由與 hooks（`useApplySchemaMutation` 等）仍在，只是沒有 UI 呼叫。
- 測試（元件，使用 mock 的 API）：`writeGroupSection`（14）、`writeGroupMembers`（9）、`writeGroupLifecycle`（21，含 Apply 依伺服器 readiness／舊 revision／未儲存變更、單擊單送、409 衝突、刪除前影響、D 試寫的全部路徑、過期範圍回覆被丟棄、交付各欄位分開）；變異檢查：Apply 不看 readiness → 3 項失敗；移除過期回覆守門 → 新增的同元件測試失敗。這些元件測試是在元件之後寫的（沒有先紅），靠變異檢查補強。
- 基準：`npm run lint`、`tsc`、`npm test`（178 files／1063 tests）、`npm run build` 通過；`go test -p 1 ./...` 48 個 package 零失敗；golangci-lint 只剩既有 `service_probe.go:201`。完整套件下 `runtime dashboard route` 曾因 1 秒逾時失敗一次，單獨重跑 3 次皆通過，視為負載下的計時抖動，未調查。

## 3.2／3.3 實測（元件層級，不是整個 production 應用）

- 方法：以 Playwright（headless Chromium）載入一個暫時的 harness 頁面（vite dev server，已刪除），只掛載 `WriteGroupSection`，API 以攔截方式回傳固定 fixture（真實的前端 parser 照常運作；第一次我的 readiness fixture 少了 `schema_digest`，被嚴格 parser 正確拒絕，補上後才通過）。這是**元件在真實瀏覽器的版面與鍵盤行為**，不是完整 `/studio/v2` 流程；整個 production binary＋內嵌 UI＋真實資料庫的證據留給 F。
- 結果（390／768／1440 px × zh-TW／en，共 6 張全頁截圖在 `docs/plans/studio-v2-write-groups/evidence-e/`）：頁面都沒有水平溢出（`scrollWidth` 等於視窗寬度）；成員表自帶捲動區（390 px：可視 272／內容 640，768 與 1440 沒有溢出）；鍵盤 Tab 可依序到達名稱、搜尋、Apply、試寫預覽，焦點框可見（截圖可見）；未儲存變更時才有可用的 Save，無變更時 Save 為 disabled 因此不在 Tab 順序中（符合預期）。
- 沒有做：螢幕閱讀器實測、IME 在群組編輯器的行為、真實裝置。
- Studio inventory：`backend-api-registry.md` 新增 write-group 路由一節（先前 A–D 都沒登錄），並以 `studio_inventory_changelog add` 同步 `changelog.sqlite`。route identity（`/studio/v2`）未變。
- AGENTS 全套檢查：前端 `lint`、`tsc`、`npm test`（178 files／1063 tests）、`build`；後端 `go test -p 1 ./...` 48 package 零失敗、golangci-lint 僅剩既有 `service_probe.go:201`、`git diff --check`、`make check-lines` 通過。

## 驗證後補做（spectra-verify W1–W4）
- W1：`basic-point-tag-to-group.test.ts` 的 `record` 少 `updated_at` 導致 `tsc`／build 失敗，已補；`tsc`、lint、`npm test`（180 files／1071 tests）、`build` 重跑通過。
- W2：新增 `GroupColumnProposal`／`proposal.ts`：已選 Tag 在表中沒有欄位時，列出提案欄位（名稱去除不合法字元、避開既有欄位、型別對應）並標示為提案，說明可改選其他資料表、移除這些 Tag 或參考提案；不會自動指派或建立。測試 `writeGroupProposal.test.tsx`（3 個 Tag 但表只有 1 個可用欄位：不環繞、不指派、Save 保持 disabled）。managed 的提案只是供檢視，沒有接到建表流程（建表仍走既有 schema 區塊）。
- W3：Share-only：`CommitSummary` 在沒有任何群組且 readiness 無阻擋時可啟用，Share readiness 阻擋仍會停住（`writeGroupProposal.test.tsx`）。
- W4：標頭與歸屬已更正。
- Suggestion：`commitLog.ts` 仍有 `enabledTargetCount`，只用於舊的 commit log 標籤，沒有擋任何啟用路徑，未更動。`GroupDeliveryStrip` 沒有 verified／cleanup 欄位，它們在 `GroupTestWritePanel` 分開顯示。

## 已知限制
- managed 新表的「提案」只是列出所需欄位與型別供檢視（`GroupColumnProposal`，標明「提案——資料表中尚無」，不指派、不建立）；實際建表仍走 Step 4 既有的 schema 區塊，沒有把提案接到建表流程。
- 協議可用性只由前端 parser 判定，沒有後端 capability endpoint。

## spectra-review（2026-10-02）與處置

範圍只含 Step 1／2 的程式與其測試（Step 3／4、locale、後端不在 review 範圍，歸屬為近似）。2 Warning、2 Suggestion，皆已處理：

1. **Warning（已修）**：位址正規化把 MC3E 的 X／Y／B 接點當十進位，會誤報或漏報衝突。現在 `detectAddressConflicts`／`describeAddressConflicts` 接受設備協議對照，改用 `addressParser.parse`（MC3E 接點為十六進位）；未提供協議時才用原本的啟發式。測試：`X9 width2` 與 `XA` 衝突、與 `X10` 不衝突。
2. **Warning（已修）**：衝突比較是逐對 O(n²)，且 `MergedPointTable` 每次 render 都重算。改為每點只算一次 span、依 area 排序後掃描，並以 `useMemo` 快取描述；5000 點單設備測試 < 1 秒。
3. **Suggestion（已修）**：`detect`／`describe` 重複的迴圈合併為同一個 `conflictingPoints`。
4. **Suggestion（已修）**：改設定不再一律使 probe 失效；`timeout` 這類不改變連線身分的值保留 probe（測試），其餘仍失效。

重跑：`tsc`、lint、`npm test`（180 files／1074 tests）通過。
