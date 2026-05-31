# 發現與決策

## 需求
- 依序實作 `datalink-workbench-v2-shell` 等 6 個變更。
- 另開 Git branch 來進行開發。

## 研究發現
- `spectra list` 顯示所有 6 個 workbench 變更皆已在 `in-progress` 狀態，總任務數量如下：
  - `datalink-workbench-v2-shell`: 34 任務
  - `datalink-workbench-v2-step1-device`: 24 任務
  - `datalink-workbench-v2-step2-rule`: 24 任務
  - `datalink-workbench-v2-step3-mapping`: 22 任務
  - `datalink-workbench-v2-step4-database`: 25 任務
  - `datalink-workbench-v2-settings`: 23 任務

## 技術決策
| 決策 | 理由 |
|------|------|
| 新建分支 `feature/datalink-workbench-v2` | 提供乾淨的開發環境 |
| surface / API 維護文件採多檔案拆分 | `/studio` 主線、`/studio/v2`、runtime、gateway、test 的成熟度差異很大，若混在單一文件會難維護 |

## 遇到的問題
| 問題 | 解決方案 |
|------|---------|
| 點位高頻 stale 事件狂刷 `setLogs` 造成嚴重重繪卡頓與日誌爆量 | 實作 `pointStaleStatesRef` 僅在點位 stale 狀態「改變」的瞬間才記錄 transition log (並新增 recovered 復原狀態日誌)。 |
| 新日誌 prepend 到最上方，與診斷面板的「滾動置底」UX 機制相反 | 將新日誌修改為 append 到陣列最後，使時間排序由上到下遞增，完美咬合 auto-scroll-to-bottom。 |

## 新研究發現

- `/studio` 是目前唯一真正深度綁定 datalink persisted APIs 的主產品 surface。
- `/studio/v2` 雖然已具備導向型流程與 runtime handoff，但 commit 仍主要是前端 reducer / 模擬日誌，不是正式 backend lifecycle。
- `/studio/runtime` 已有 snapshot + SSE，但 lifecycle 語意仍不足，前端部分狀態需要字串判斷。
- `/gateway/*` 目前主要依賴 feature flag 與 `/test/connect`，應視為 prototype / connect-only，而不是 datalink persisted flow。
- router 預設入口目前仍是 `/studio`，與未來希望服務預設採 `/studio/v2` 的方向不一致。
- 使用者已明確決定：`/studio` 這一輪先暫停；`/studio/v2` 與 `/studio/runtime` 是近期重點；`/test` 先暫停但需要完整記錄；文件需要反映這個優先順序，而不是只記技術成熟度。
- 第一版 `index.html` 看不到資料的根因是內嵌 JS template string 內含未跳脫反引號內容，導致整段 script parse 失敗、render 完全沒啟動；已改為 `index.html` + `inventory.css` + `inventory.js` 三檔。
- repo 目前沒有既存的 `studio-surface-inventory` 變更履歷機制；若只靠 md/html 本身，無法避免文件被更新後難以追溯。
- `go.mod` 已內建 `modernc.org/sqlite`，可直接在 repo 內提供小型 Go CLI 來管理 inventory changelog，不必引入新外部依賴。
- `cmd/studio_inventory_changelog` 現已落地，並已初始化 `docs/technical/studio-surface-inventory/changelog.sqlite`；後續 inventory 文件更新應以這個 CLI 寫入 changelog，而不是手動維護另一份文字紀錄。
- 若要讓接手 AI 快速上手，關鍵不是再增加長文，而是提供固定第一入口與 machine-readable 摘要，讓它先讀 `START_HERE.md` + `context.json`，再按需展開長文件。
- 要進一步接近「新對話直接接上」，還需要一份最新狀態快照；否則 AI 雖知道去哪裡讀，仍要自己整理最近一次的決策與下一步。
- `docs/goal.md` 不是單一 spec，而是 10 個 `/goal` 執行入口；規則明確要求一次只跑一個 goal，且建議順序先從 `make-studio-v2-default-entry` 開始。
- 目前 `make-studio-v2-default-entry` 在 Spectra 中 artifact 已齊，但 `tasks.md` 4 個 task 都還沒勾選完成。
- 這個 change 的 spec 共有 2 個 `SHALL` requirement，受影響路徑是 `/`、unknown fallback、`/datalink`、`/datalink/workbench`、`/studio`。
- `docs/goal.md` 這一輪把 scope 收窄為 `frontend/src/App.tsx`、`frontend/src/features/datalink/legacyRoutes.ts`、`frontend/tests/unit/workbench-v2/routing.test.tsx`，因此不應順手修改 inventory 文件。
- 現況 root cause 已確認：`App.tsx` 仍把 `/` 與 `*` 導向 `/studio`，且 `/datalink`、`/datalink/workbench` 仍走 `LegacyStudioRedirect` 進入 legacy `/studio`。
- repo 內沒有可讀的 `RTK.md`；本輪只能依 `AGENTS.md`、`CLAUDE.md`、`.github/instructions/*` 與 Spectra artifacts 執行。
- `docs/goal.md` 指定的驗證命令 `cd frontend && npm run test -- --run frontend/tests/unit/workbench-v2/routing.test.tsx` 已與目前 Vitest include 規則脫節；正確可執行的等價命令是 `cd frontend && npm run test -- --run tests/unit/workbench-v2/routing.test.tsx`。
- routing test 的 RED 已明確證明行為缺口：`/`、unknown fallback、`/datalink`、`/datalink/workbench` 都先落到 legacy `/studio`，而不是 `/studio/v2`。
- 最小修補只需要兩個 implementation 動作：新增 `buildWorkbenchV2EntryRedirect()` helper，並讓 `/`、`*`、`/datalink`、`/datalink/workbench` 走 `GuidedWorkbenchEntryRedirect`；`/studio` 與 `/datalink/workbench/*` 維持 legacy 邊界不變。
- 第二個 change 的最小持久化路徑不需要新資料表：直接把 singleton workspace metadata JSON 存在既有 `system_settings` 表即可滿足持久化、重啟重讀與隔離需求，並避免為了單一 metadata 引入 migration churn。
- `cmd/test_ui` 的 wiring 需要一起納入驗證；即使沒有測試檔，`go test ./cmd/test_ui` 仍能證明新增的 workspace service 與 router 注入鏈可編譯。
- 後端 targeted `go test` 在 sandbox 內會因 Go build cache (`~/Library/Caches/go-build`) 權限失敗；這不是程式錯誤，需要用提升權限重跑。
- 第三個 change 若沿用後端自產 UUID，Step 1 create 後就必須同步 remap前端 device/rule/point/mapping/db target 關聯；最小可行路徑是讓 workspace device create 接受既有 client `dev-*` id，先穩住 Step 1 autosave 與後續 step 的本地關聯。
- Step 1 autosave 最穩定的邊界不是塞進 `Step1Device`，而是放在 page-level orchestration：UI 元件維持純 dispatch，page 入口決定何時 create/update/delete、何時標示 `draft-invalid` / `saving` / `saved` / `save-error`。
- `DatalinkWorkbenchV2Page` 新增 device list bootstrap 之後，所有直接 render page 的舊 tests 都需要補 QueryClientProvider 或 devices list mock，否則只會停在 bootstrap loading/error，而不是 shell 本身壞掉。
- Step 2 rule autosave 可直接重用既有 `source_rules` 持久化模型，但 backend 目前沒有 `name`、`share_enabled`、`share_start_register`、`share_stride` 欄位；這次 change 應只 autosave backend 已有契約欄位，否則 reload 後會產生假保存。
- 最穩定的 Step 2 前端做法和 Step 1 一樣，是把 autosave orchestration 放在 page-level hook；Step 2 元件維持純 dispatch，save queue / hydration / per-rule error state 由入口統一處理。
- `DatalinkWorkbenchV2Page` 在 workspace boot 之後若同時加入 rule bootstrap，所有 page-level tests 都要額外 mock `studioV2RulesAPI.list`，否則會卡在新增的 query loading/error gate。
- Step 3 mapping autosave 不能直接重用 generic `mapping` CRUD payload，因為前端 row model 是 `rule_id + address + tag metadata + transform fields`，而 backend 真正持久化的是 `point_id + tag_id + transform_pipeline`；因此需要一層 workspace-scoped mapping handler 來做 point/tag/mapping translation。
- Step 3 的 hydration key 不能用 backend `point_id`，因為前端目前仍以 local derived ids (`${rule.id}-p-${index}`) 管理 row state；最穩定的對齊鍵是 `rule_id + address`，並把真正的 persisted `point_id` 只當 metadata 帶回前端。
- Step 3 的 `local_value` / `persisted_value` reconcile 若沒有 guard，很容易在 `useEffect` 中形成自我觸發 loop，或在使用者剛編輯時把 dirty local 值被 persisted hydration 蓋回去；需要同時處理「無變更不 dispatch」與「只有已有 persisted baseline 且已分叉時才保留 local」這兩個條件。
- `/studio/runtime` 舊實作仍是 `device_id`-first；沒有 query 時直接落到 `missing-device-context`，而不是先讀 singleton workspace runtime context。
- runtime page 若要保持 Step 4 handoff 相容，`device_id` query 最穩定的收斂方式是降級成 optional override：query 命中 workspace device 時沿用，否則回退到 `default_device_id`，不強制重寫 URL。
- runtime context payload 不需要把整份 legacy `Device` record 搬到前端；最小契約只要 `workspace_id`、ordered `devices[]`、每台的 `device_id/name/protocol/running/availability_*`、以及 `default_device_id`。
- `/studio/v2` 空 workspace 之所以仍顯示 seeded draft data，不是 workspace API 回傳錯，而是 `useStudioV2AutosaveState` 在 `devices=[] && rules=[]` 時直接 return，讓 `INITIAL_STATE` 的預設 device/rule 留在頁面上。
- 最小且低風險的修法不是移除全域 `INITIAL_STATE` 預設，而是新增 workspace hydration gate：空 workspace 明確 dispatch 成 `devices=[]` / `rules=[]`，且 page 在 hydration 完成前不渲染 shell，避免 seeded data flash。
- Step 1 live diagnostics 不需要新增後端 endpoint；既有 `POST /api/v1/datalink/devices/test-draft` 已經回傳 `connect/probe/can_activate/can_collect`，缺口只是在 V2 Step 1 仍停留於 timer/mock animation。
- 若直接把 backend `connect/probe` 結果硬灌回舊的 `resolve/connect/probe` 前端假步驟，會留下語義漂移；最小且一致的做法是 Step 1 測試開始後直接以 backend `connect/probe` 兩段作為 UI 真相來源。
- `/studio/v2/settings` 的真正缺口不是 endpoint 不存在，而是 front-end 完全沒用到已存在的 `settingsAPI` / `dbTargetAPI`；boot/save/connector test 都還停在 reducer demo flow。
- `/settings` backend 雖只內建少數預設 key，但 handler/service 已經支援 generic key/value upsert；這意味著 V2 settings 多數欄位可以直接走現有 `/settings/:key` persistence，而不必先新增另一套 V2-only settings API。
- settings boot 若 backend 讀取失敗，正確行為不是默默沿用 reducer defaults，而是停在明確 loading/error gate；否則 operator 會把本地預設值誤認成已載入後端狀態。
- connector pool 與 settings persistence 的穩定邊界是在 page-level orchestration 加 mapping layer：UI 繼續編輯 `SettingsConnector` / reducer state，backend shape 差異由 `/settings` key map 與 `db-targets/connectors` request builder 吸收。
- Step 3 `轉換管線預覽` 目前仍是本地 deterministic seed + `runScale/castValue/formatFinal`，不是 backend 真實 preview；spec 也明寫目前 preview 不得發 request。
- backend 已存在 `POST /api/v1/datalink/mappings/preview` 與前端 `mappingAPI.preview(...)`，因此 Step 3 的正確修法是重用現有 preview API，而不是新增新的 V2-only preview endpoint。
- `目標型態` 的「一鍵套用到全部列」其實已有 reducer contract 可重用：`bulkApplyTransform` 已能對 `state.mappings` 批次套用欄位；這次只需要把 preview 卡上的快捷操作收斂成 `target_type` 專用入口。
- Step 3 live preview 落地後，operator-visible 管線結果已改由 backend `raw_value` / `step_results` / `final_value` 決定；本地 `runScale/castValue/formatFinal` 只保留在 preview 缺席時的 fallback 渲染路徑。
- `MappingPreviewResponse.step_results[].output_value` 是 `unknown` shape，若直接餵給 `formatFinal` 會在 `tsc` 暴露型別錯誤；最小修法是用 scalar type guard 收斂成 `number | boolean | string`，而不是放寬成 `any`。
- `/studio/v2` 的 source-rule autosave 500 不是 handler panic，而是 `ensureRuleActivationAllowed()` 回傳的 readiness block 沒包成 `sourcerule.ErrValidation`，因此 handler 落到 500；把它改成 validation 後，同一路徑會正確回 400。
- Step 3 mapping row 目前以 `${rule.id}-p-${index}` 當本地 key；當 rule 的 `address/device_id` 已變但 local `point.id` 尚未變時，若沿用舊 mapping state，就會把過期的 `mapping_id/persisted_value` 帶到新 row，直接導致 `mapping row not found` 與後續 save-error。
- 即使 row identity reset 之後，如果 owning rule 本身仍是 `save-error`/未 persisted，Step 3 再打 mapping autosave 仍會對到 backend 舊 rule surface；因此 mapping autosave 需要先看 owning rule 是否已 `persisted && save_state === 'saved'`。
- Step 3 目前真正的假值來源不是 backend preview，而是前端 selector / preview 卡仍用 deterministic `RAW_VALUE_SEEDS` 當 `raw_value` 輸入；把 preview 接到真 API 之後，若 seed 依賴不移除，UI 仍會對 operator 顯示假數值。
- repo 內既有 `useRuntimeStream` 只支援單一 `deviceId`，不足以支撐 `/studio/v2` Step 3 多裝置表格同時顯示 live values；最小修法是新增 Step 3 專用多裝置 SSE 聚合 hook，而不是改寫 runtime dashboard 既有 hook。
- Step 3 preview 在沒有即時值時，正確行為不是退回 fake seed，而是停在明確 waiting state；否則 operator 會誤判目前 preview 代表真實設備輸入。
- Step 3 第二個更深的 bug 是 `useStep3LiveValues` 一開始拿本地 row id (`rule-x-p-y`) 去訂閱 runtime stream，也用同一個本地 id lookup SSE payload；但 backend/runtime 真正認得的是 `persisted_point_id`。這會讓 UI 在設備其實有值時仍永遠停在 `尚未收到`。
- 最小且正確的修法不是把 runtime stream 改成接受前端假 id，而是以前端 mapping metadata 的 `persisted_point_id` 做 stream filter，收到 event 後再依 `persisted_point_id` 或 `device_id + address` 映回本地 row。
- Step 3 還有一個獨立故障：`useStep3LiveValues` 原本直接依賴整個 `mappings` 物件與由其衍生的 lookup object，任何 `save_state` / autosave patch 都會讓 effect cleanup + reopen `EventSource`，形成 backend `runtime/stream` GET 風暴，前端則常駐 `連線中`。
- 正確依賴邊界不是整個 mapping row，而是訂閱真正需要的欄位：`local point id + device_id + address + persisted_point_id`。只要這四個沒變，SSE stream 就不該重建。
- Step 3 的「數值顯示」不能只靠 runtime SSE。舊 `/studio` 早就用 `SSE raw_value ?? point.last_value`；如果 V2 少了 `last_value` fallback，就會出現設備其實已有最後讀值，但表格與 preview 仍顯示空白或 waiting。
- 最小修法不是在 UI 層塞預設字串，而是補 persisted points snapshot 查詢，先用 backend `/points` 拿 `last_value`，再讓 SSE 到來後覆蓋它。
- Step 3 `Point → Tag` 現在若用 `!skipped` 當 row 條件，會把 disabled rule 的 point 也一起渲染成可編輯 row；operator 會誤以為這些 row 的 enabled toggle 也會持久化，但 reload 後 backend 仍只反映真正 `point.enabled === true` 的 rows，看起來就像「剛開 8 個，回來只剩 4 個」。

## 資源
- [useRuntimeStream.ts](file:///Users/yishow/prj/go_gateway/frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts)

## 視覺/瀏覽器發現
- 新的即時診斷日誌面板 (RealtimeLogsPanel.tsx) 配合 auto-scroll 效果在點位狀態轉變時非常流暢且易讀，在點位等待時呈現呼吸波動 Skeleton 骨架屏顯得富有質感。

---
*每執行2次查看/瀏覽器/搜尋操作後更新此檔案*
*防止視覺資訊遺失*
