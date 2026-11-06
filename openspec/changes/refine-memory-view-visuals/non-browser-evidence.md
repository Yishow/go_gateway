# D 非瀏覽器驗證

基準 HEAD：5047ab54f94e2bc185cc3f5c33e092ebc1a4b458。沒有提交、封存或 push；所有 browser/Mac 驗收仍未執行。

## 本次完整檢查

- 最終受影響 memory suite：`npm test -- --run --maxWorkers=2 tests/unit/workbench-v2/recording-home-memory`，passed-current，16 files / 155 tests，exit 0，22.80 秒。

- `cd frontend && npm test -- --run --maxWorkers=2`：passed-current，667 files / 7623 tests，exit 0，558.46 秒。來源為 final-tests-attempt2.log 及明確 exit 檔。
- `cd frontend && npm run lint`：passed-current，exit 0。
- `cd frontend && npm run build`：passed-current，exit 0；保留既有 >500kB bundle warning，不擴張到無關 bundle 重構。
- 最後兩個 E2E-only 修正：`npx eslint tests/e2e/recording-home-memory-view.spec.ts tests/e2e/helpers/recording-memory-visual-audit.ts`、memory-E2E scoped TypeScript 均 exit 0。
- `npx playwright test --list`：passed-current，602 tests / 54 files；其中 memory discovery 是 162 tests / 10 files，新 visuals suite 39 cases。只發現測試，沒有 browser。
- `git diff --check`、`make check-lines` 通過。兩個 locale 384 行及 memory-view test 347 行警告，沒有 >500 行阻擋。
- 官方 analyze / validate 通過結構檢查；analyze 留七個缺少具體 Example 的 Suggestion，沒有 Critical/Warning。

第一輪完整測試因執行環境中斷且無 exit 記錄，只列 interrupted，不當作 PASS/FAIL；本次採獨立 attempt2 log 成功。沒有縮短 timeout 或重複啟動仍存活的 suite。

## 測試先行與修正

- 共用 cells：新測試先記六個失敗；required-final-green 6 files / 95 tests。
- navigator：新行為先一個失敗；navigation-green 4 files / 23 tests；conflict mutation 被斷言捕獲後還原。
- map：先兩個失敗，mutation 三個失敗後還原；map-final-green 7 files / 42 tests，offpage 2 files / 11 tests。
- chapters：先三個失敗，3 files / 21 tests 通過；完整 40 ranges 和末項定位保留。
- detail：先兩個失敗，最終 6 files / 40 tests 通過；ordinary reading 不改成 owner name，focus disclosure 保留完整 owner 資訊。
- review 修正：先四個失敗，最終 4 files / 67 tests 通過；涵蓋三版型 offpage 真 hook 更新和跨區相同 offset 不冒充 selection。
- 原 map recovery selector 只遷移到 compact summary，原安全 assertions 保留。observer target 數增加時保留精確四項及重複 key 去重斷言，沒有以刪測試取得通過。

## 內容指紋與限制

交接附 final-qa-inputs.json（1630 個 frontend source/test/config 檔）、final-qa-inputs-current.json、差異比較及環境版本。完整測試啟動後只有兩個 E2E-only 檔修正；這兩檔由最後 lint/types/discovery 覆蓋。Vitest 明確排除 E2E，build source 未改。完整候選 SHA256、Git blob/mode、逐檔差異和原始 log 隨提交包保存。

擴大到全 E2E 的獨立 TypeScript config 仍失敗：九個錯誤僅在未修改的 embedded-frontend-delivery.spec.ts、recording-mixed-review-host.tsx。兩檔 SHA256 與基準 HEAD 相同，baseline 證據隨交接提供；memory scoped types 通過不等於全 E2E types 通過。沒有修無關檔案。

canonical memory map/values、SourceContent、B source mutation、C refresh authority 五檔 SHA256 保持相同。無後端變更，不重跑完整 Go；既有 Go concurrency/global lint 問題不列成本案新回歸。

所有真 geometry/contrast、30 張對照、觸控、production modal/history、SQL/設備與 archive 均 not-run。這份報告不替代最後 Mac 驗收。
