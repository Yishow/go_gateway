# PLC 修復與本機驗證

已依本 change 的五類修復及 review 缺口完成實作與本機驗證。未 commit、push、archive 或部署；真 PLC、正式資料庫、Windows 與現場驗收均為 **NOT RUN**。

- Git 基準：`main`，`HEAD`／本機 `origin/main` 均為 `76a853d71fda23ee3a1d3cf9bf801e45f85f24f2`。
- 授權範圍：[tasks.md](tasks.md)、[design.md](design.md) 的 Implementation Contract；使用者另明示不得擴大範圍。
- 提案內「本次只起草」記錄的是 propose 階段；後續明確的 `$spectra-apply` 授權才執行本次實作與測試，未變更原有產品契約。
- 完整修改清單：[scope-files.txt](scope-files.txt)。四個既有 runtime-dashboard WIP 檔案已依最初 SHA256 比對，全部保持原樣；index 無 staged 變更。

## 修復內容與必要依賴

| 修復 | 實際行為與主要位置 |
| --- | --- |
| 共用解碼 | `lib/hsllogic/{bytetransform,register_conversion}.go` 重用既有 byte transform，讓 BADC／DCBA 的 32／64-bit signed、unsigned、float 保留完整 bits；保留單 word、ABCD／CDAB 與既有短資料 API。Modbus TCP／RTU／UDP、MC 均透過真 adapter fixture 驗證。 |
| FATEK typed scalar／count | `adapters/fatek_helpers.go` 依 16／32-bit 原生元件組合完整值並回傳 conversion error。`connector.PointReadCountForAddress` 讓 periodic、manual、read-now 的 D／R 與 DR／DD count 使用原生元件單位，無效地址在 Read 前失敗。direct component Count 保留。 |
| MC radix／既存身分 | backend planner 與前端 `AddressParser` 採 X/Y octal、W/B hex、D/M decimal。`sourcerule/mc_address_plan.go` 在任何 Update mutation 前比較 stored plan 與 persisted links；runtime reconcile 不會在該拒絕路徑執行 snapshot restore。舊規則啟動／重載保留實際身分。 |
| MC 畫面保護 | `mcPersistedPoints.ts`、scoped query hook 與三類 autosave 保留實際 point/tag/mapping 關聯，撤銷舊 scope 的寫入資格。selectors、Step2 grid／摘要及 Step3 共用 canonical projection；checking／conflict／unavailable 不回退為 synthetic points。包含停用規則的摘要也不宣稱「衝突 0」。 |
| Modbus precedence | `adapters/modbus_parse_helpers.go` 依 alias → 合法傳統五位編號 → explicit numeric Function → no-function 歷史 fallback 分類；Read／Write 共用。未改 persisted Address／Function／ID。 |
| FATEK Cmd44 | `internal/protocol/fatek/client.go` 要求完整長度與 ASCII 0/1，拒絕空／短／長／非 binary。client、adapter、三個採集入口皆保留 error、bad quality、nil Value；不修改下游 recording policy。 |

必要 caller 查核已確認：三個 scalar caller 均使用 address-aware helper；FATEK converter 的 adapter caller 處理 error；MC／Modbus register 解碼使用共用 converter。舊 `PointReadCount` 仍供非 FATEK 分支使用，予以保留。原本 `XF → X10` 的測試依 MODIFIED grammar 改為 `X7 → X10`，未改 runtime parser。

## Breaking 行為與回復限制

1. MC 新規則 grammar 改為上述 radix。W0/count17 現在保存 `W0…W9, WA…WF, W10`；X0+16 為 X20。新 plan 拒絕 X8／X9／YA／WG／Z0。舊 direct MC X/Y 的 octal-first decimal fallback 保持原樣；不宣稱所有入口接受集合相同。
2. Modbus 短數字且已保存明確 Function 的實際目標改變，例如 `Address="100", Function="03"`：修復前 **FC01／offset99**，修復後 **FC03／offset100**。沒有自動地址或資料 migration；傳統 `40001` 仍優先代表 holding offset0，不能藉 Function 表達同字串的 raw40001。

byte order／FATEK 解碼修正也會改變過去錯誤的讀值；不重算歷史。此次沒有 persistent schema migration，可回復程式，但新 WA／X20 等設定不能假設舊版本可編輯；回退前需盤點並停止相關設定編輯，不能刪除 point／mapping／歷史資料代替回復。

## 實際驗證

| 檢查 | 結果與證據 |
| --- | --- |
| 後端模組 | **passed-current**：`lib/hsllogic`、`internal/protocol/fatek`、`internal/datalink/connector/...`、collector、device、point、sourcerule 共 8 packages、835 個含 subtest 的 PASS 項目；0 fail／skip。 |
| API caller | **PASS**：`go test -json ./internal/api ./internal/api/handlers -run 'SourceRuleHandler_\|StudioV2WorkspaceSourceRule' -count=1`；2 packages、49 個含 subtest 項目，無 skip。 |
| Race／靜態檢查 | **passed-current**：上述受影響 packages 的 vet 與 golangci-lint；lint 0 issues。MC bounds／runtime reconcile focused race 通過。初次 apply 的 collector malformed status／FATEK scalar count race 紀錄另保留。 |
| 前端選測 | **passed-current**：addressParser、MC hydration／view／draft retention／bounds／allocator、Step2 grid／rule／occupied-range、rule／mapping autosave、bootstrap、selectors、Step3、cold locator 等 29 files／209 tests。 |
| 前端靜態／建置 | **passed-current**：29 個修改／新增 TS 檔 ESLint；`npx tsc --noEmit --incremental false`；最後 `npm run build` 與 static sync 後的 Go embed build。Vite 仍有 chunk size 提醒，未擴大為打包優化。 |
| 瀏覽器 | **PASS**：獨立 headless Chromium、真實 backend、兩個拋棄式 SQLite workspace、loopback MC simulator。下方列出實際可觀察範圍。 |
| Spec／交付 gate | **PASS**：`spectra validate`、`git diff --check`、`make check-lines`。歷史 `sourcerule/service.go` 1140→1139 行；其餘僅有 >300 行提醒，無 >500 行阻擋。analyze 為 0 Critical／0 Warning／4 個原有 example 建議，未擴寫規格。 |

可重跑的後端主命令（從 repo root）：

```sh
go test ./lib/hsllogic ./internal/protocol/fatek ./internal/datalink/connector/... ./internal/datalink/collector ./internal/datalink/device ./internal/datalink/point ./internal/datalink/sourcerule -count=1
go test ./internal/api ./internal/api/handlers -run 'SourceRuleHandler_|StudioV2WorkspaceSourceRule' -count=1
go test -race ./internal/datalink/collector -run 'TestFatekMalformedStatusPayloadAcrossEntrypoints|TestFatekScalarPointCountAcrossEntrypoints' -count=1
```

原始執行 log／overlay、完整前端選測命令與雜湊保留於本機暫存證據目錄 `/private/tmp/gw-plc-apply-jkunk0j4`；產品回歸測試已加入 repo，無需暫存 fixture 才能重跑。沒有把含本機路徑或 generated readiness token 的原始 HAR／trace 加入交接。

### 測試先失敗、修復後通過

| 原始失敗 | 保存的證據名稱 |
| --- | --- |
| 四種 byte order 與 adapter wire 解碼 | `task-1.1-red-baseline-overlay.log`、`task-1.2-red-baseline-overlay.log` |
| FATEK 首值 fallback、原生 32-bit count 多讀 | `task-2.1-red.log`、`task-2.2-collector-red.log` |
| MC guard／runtime wrapper mutation、planner 與 point 詞形 | `task-3.1-red.log`、`task-3.1-runtime-red.log`、`task-3.2-sourcerule-red.log`、`task-3.2-point-red.log` |
| 前端 parser、persisted hydration／late ACK／mapping queue | `task3.3-red-after-literals.log`、`task-3.4-red-hydration-behavior.log`、`task-3.4-gap3-workspace-rev-red.log`、`task-3.4-mapping-queue-red.log` |
| persisted grid／未知摘要，含停用規則 | `task-3.5-view-red.log`、`task-3.5-table-red-3.log`、`task-5.3-disabled-red.log` |
| Modbus FC03 被改為 FC01／offset99 | `task-4.1-red.log` |
| Cmd44 malformed 被當成功；三個採集入口 | `task-2.3-red.log`、`task-5.2-red-overlay.log`（9 個 malformed acquisition subcases 失敗） |

overlay 只替換暫存副本以恢復先前接受行為；未修改工作樹產品程式來製造失敗。setup／syntax／lint 修正不算 behavioral RED。

## 真實瀏覽器與資料身分

- 由 UI 連線本機 MC simulator，建立 W0/count17/int16 與 X0/count17/bool。保存後與最終 binary 重載，34 個 DOM point IDs、地址、source-rule revisions 均與保存快照一致；Step2 可繼續。
- 舊 fixture 透過 HEAD 的舊 planner overlay 建立 W0…W16 與一個既有 mapping，再以修復版啟動同一個拋棄式 DB。17 個 point IDs、tag/mapping 關聯、規則 revision 與既存 mapping 記錄保持不變；不生成 WA…WF。
- 舊規則維持 conflict／blocked。UI 的非位址 prefix 編輯僅留在本地，未寫回；重載恢復原保存設定。停用規則也顯示證據待確認，不宣稱乾淨的衝突數字。
- 錄製 124 次 API 請求，其中 legacy 流程 36 次均為 GET；0 次 legacy mutation，0 次 unscoped candidate GET，0 次 candidates recompute／apply POST。GET 內部原有 database-output snapshot refresh 仍保留。
- 這只證明 device probe 與規則保存／顯示／重載；未宣稱 scheduler 已採集或 output 已 committed。
- 測試結束已關閉專用 browser、兩個 gateway 與 MC simulator，三個 loopback 測試埠均已關閉。

去敏證據：[驗證結果](evidence/browser-verification.json)、[請求紀錄](evidence/browser-network.json)。截圖：[W0](evidence/mc-w0-reloaded.png)、[X0](evidence/mc-x0-reloaded.png)、[legacy blocked](evidence/mc-legacy-blocked.png)。

## 規格測試對照與邊界

[scenario-coverage.md](scenario-coverage.md) 將 32 個 Scenario、5 個 Example（33 列）逐項對應到測試；全部有對應，沒有 test-scope exclusion 或 uncovered gap，並揭露 exact fixture variation。

未新增 FATEK station 政策／RT／RC、MC area／家族／預設值設定、新 API／schema／migration UI、外部建表／試寫、歷史回填或 runtime／recording WIP 修復。沒有操作使用者原有持久化資料；測試寫入只發生於本次建立的可丟棄資料庫。

## Review 修復與最終對照

獨立 `$spectra-review`／`$spectra-verify` 已執行。原 5 項 Warning 已修復，後續瀏覽器發現的兩個同案 MC 缺口也已先取得測試失敗再修復；最終沒有 active Critical／Warning。完整對照見 [review-verify.md](review-verify.md)。

| 缺口 | 修復與回歸 |
| --- | --- |
| R1 草稿遺失 | scoped evidence checking／unavailable／compatible 轉換保留沒有 server record 的 mapping／database target 草稿，並保留既有 MC／Modbus 草稿；`mc-persisted-draft-preservation.test.ts`。 |
| R2 冷啟動 locator | MC 單點 projection 使用 mapping 的原生 point ID／address，source data type 取 rule；`studio-locator-mc-hydration.test.tsx`。 |
| R3 新增起點 | MC radix 與實際 persisted address list 決定下一個位址，X0/count8→X10、W0/count17→W11，legacy W0…W16→W17；`mc-next-free-address.test.ts`。 |
| R4 上限／UI | 前後端拒絕超出既有 client 0xFFFFFF 的位址／展開溢位，摘要不拋例外；WA／WF／BFF 的 area 與 suffix 分開；compatible evidence 不掩蓋當前無效草稿，顯示錯誤並阻擋 Continue。Go／TS bounds 與 `mc-persisted-readiness-regression.test.ts`。 |
| R5 無寫入拒絕 | early MC validation 使用可識別的 pre-mutation error，wrapper 不執行 restore writes；post-write validation 仍回復完整 snapshot；`mc_validation_boundary_test.go`。 |

本輪 RED／GREEN 紀錄為 `r1-red.log`、`r2-red.log`、`r3-red.log`、`r4-*-red.txt`、`r5-red.log`、`r4-ui-red.log` 與對應 GREEN／整合測試 log。原先 Go cache 權限造成的 setup failure 另存，改用專用暫存 cache 後重新通過，不當成 behavioral RED。

最終 browser 以全新拋棄式 SQLite 和最終 binary 重驗 X10／W11 預填、WA…WF 的 Word 標籤，以及 WFFFFFFFF／WFFFFFF-count8 的錯誤與 Continue disabled。WFFFFFFFF 前後 canonical source rules 相同；invalid update 的 HTTP 400 屬明確拒絕，不宣稱所有 invalid edits 都零 HTTP。61 次 HAR 請求中 59 次為 API，0 unscoped candidates、0 candidates recompute／apply。專用 browser／gateway／simulator 已關閉。

去敏證據：[本輪結果](evidence/followup-browser-verification.json)、[本輪請求](evidence/followup-browser-network.json)。截圖：[X10](evidence/mc-next-x10.png)、[W11](evidence/mc-next-w11.png)、[range overflow blocked](evidence/mc-range-overflow.png)。

未 archive。本機證據不替代真 PLC、正式 DB、Windows 或現場驗收。
