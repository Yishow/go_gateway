# 設定保存回饋：雲端驗證紀錄

## 範圍與現況

基準是 B commit `406490c8ff1f801c27ddc491d9f29debbd12d523`。本 change 僅調整前端 database configuration ACK／刷新、必要 scoped read callers 與測試；沒有後端或 API schema 變更。沿用 B 的 home mutation capture／complete，沒有第二份 home cache。

有效 ACK 必須匹配事前已知 workspace、connector identity 與非空 identity/setup revisions。ACK 結算只表示該 submitted configuration 已保存，不表示 activation、runtime running、SQL commit、readback 或 cleanup 成功。

- 開始保存與有效 ACK 都取消／隔離舊 GET；query key 和 request 捕捉 backend、完整 baseURL、workspace、generation 與 AbortSignal。
- config、targets、bootstrap 刷新並行；同範圍重複需求合併，過期排程不執行。ACK 不等待刷新。
- 折疊摘要可見刷新中／錯誤與唯讀 retry。需 fresh evidence 的動作保持 unavailable。
- 同範圍已確認保存的 editor 在背景讀取失敗時保留；初始 malformed bootstrap 仍顯示失敗。
- reducer、submitted snapshot、newer edits、CAS、unknown/no replay 與 target/write-group 等待契約保持。

## 測試先失敗再通過

新測試分布於 `frontend/tests/unit/workbench-v2/database-save-feedback*.test.tsx` 與 `database-save-refresh-scope.test.tsx`。

- 有效 ACK 可早於三個 held reads；缺 id/workspace/revisions 或無事前 owner 不可假設保存。
- 合成案例精確設定 ACK=1ms，三讀取各=200ms；三者同時開始於下一事件迴圈，保存於1ms結算。將 refresh 改回 await 的 mutation check 會使此測試失敗。
- 範圍切換、完整 baseURL、忽略取消的晚到 GET、PUT 期間開始的 GET、新 ACK、更新欄位與 unmount 均有回歸。
- 事前 owner 無法從 ACK 建立；active B 不被 inactive A 污染，current observer 失敗不回退舊 owner。
- 折疊摘要不顯示 raw error；刷新失敗不變成保存失敗；唯讀 retry 不送 PUT。
- target 與 write-group await 的 mutation check 證明移除等待會失敗，原 helper 已還原且無產品 diff。

## 獨立 review 與修復

review 發現成功的獨立讀取／重新開啟後，refresh error 仍可能卡住操作。新增兩種失敗案例後修復：只有 stale-error、無 coordinator pending，且三類 active scoped reads 全部 success／idle／未 invalidated 才解除。只成功 bootstrap 不足，PUT 仍 pending 時也不能解除。獨立重現由2個失敗轉為2個通過。

另新增 workspace／API-origin A→B→A owner token 回歸；舊 ACK 不得套到新 owner。更動前2個新案例失敗，更動後通過。

既有 locator 測試將單一立即 enabled assertion 改為等待非同步 reconciliation 完成；後續 source-blocked、禁止儲存、lifecycle 拒絕、cached read 拒絕等斷言全部保留。沒有增加 timeout。

## 執行結果

- `npm run build`：passed-current；保留既有大 chunk 提示。
- `npm run lint`：passed-current。
- `npm test -- --run --maxWorkers=2`：passed-current；662檔／7600項全數通過，527.93秒，exit0。
- `git diff --check`、`make check-lines`：passed-current。既有兩個 locale 檔各511行，未增行；沒有新增超過500行檔案。
- 新 recovery／owner／合成案例群16/16通過；affected regression178/178、provider fixtures44/44、MC scope fixtures16/16通過。
- 先前完整執行7584/7594通過，暴露 provider、backend mock 與 workspace read fixtures 不符新契約，已保留原安全斷言修復。另有 cloud CPU contention 下 sustained-input timeout；同條件 paired focused run 的7個 sustained tests通過。最終完整執行降低 worker concurrency，不改測試 timeout、斷言或產品語意。

最終 source/test/config fingerprint 保存在雲端驗證交接資料，覆蓋1610檔；完成後核對1610檔全部一致，沒有版本漂移。

## Browser 與現場驗收：未執行

新增 `frontend/tests/e2e/recording-database-save-feedback.spec.ts` 和 strict deferred API fixture，使用真實 `/studio/v2` router。

- 20個 Playwright scenarios已 discovery；fixture unit6/6、targeted lint、repo及E2E-inclusive typecheck通過。
- 覆蓋 held ACK/三種 GET、折疊錯誤/retry、缺證據 ACK、lost ACK/no replay、新編輯、workspace切換、既有 modal及真實 Back/Forward guards。
- 現有 database autosave 沒有獨立 navigation guard；驗收明確使用既有 recording-draft guards，沒有宣稱新增 database-only guard。
- lost ACK 保護是 mounted owner 範圍；重新開啟的 GET 是觀察，不宣稱新增 durable unresolved-intent storage。
- Browser 執行、截圖、真實 SQL/field validation：not-run。雲端未啟動 browser、Mac、正式 DB、設備或提升權限。
- Tasks3.1與3.3保持未完成；所有 change 完成後在 Mac 合併 browser validation，必要修復及最終 verify/review 後才 archive。
