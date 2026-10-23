# 原未提交工作處理（2026-10-23）

基準 `f21e2a6b4997726011e4e85eba527087b6c9c8c0`。本輪獲明確授權處理原工作樹未提交內容，先核對main、remote、index及相關程序，再在既有review clone測試。初始index為空，原12檔內容指紋與上一輪保留紀錄一致；沒有新未知檔、個人設定或產物。提交與同步前再次比對原檔，避免覆蓋並行編輯。

## 原12檔如何處理

| 檔案 | 分類／處理 |
| --- | --- |
| `internal/datalink/device/service_validation.go` | 保留缺欄位彙整與helper抽出；review發現missing遮蔽empty，補`errors.Join`一次回報兩類 |
| `internal/datalink/device/service_validation_test.go` | 保留原測試意圖，修正原漏報預期；新增Modbus/FATEK/null mixed回歸 |
| `internal/apphost/data_test.go` | resolver別名、優先序、canonical path、歧義／權限、目錄樹不變；原樣保留重驗 |
| `internal/config/precedence_test.go` | process env/dotenv/defaults及選定config root；原樣保留重驗 |
| `cmd/test_ui/gateway_data_selection_test.go` | 實際host取消／歧義／DSN與CLI契約；原樣保留，Windows skip仍明列 |
| `frontend/tests/unit/runtime-logs/runtimeLogSession.test.ts` | 保留jitter極值、30秒cap及cleanup；mutation證明能抓退化 |
| `frontend/tests/e2e/embedded-frontend-delivery.spec.ts` | 保留zh-TW鍵盤；補transport中斷恢復、Clear僅清browser、source/search正反真backend情境 |
| `README.md` | 按production wiring更新WriteGroup durable與四步事實；保留unknown阻擋、修舊archive連結 |
| `docs/plans/2026-10-09-incomplete-features-plan.md` | 保留10/09–13現場歷史意圖並標明無原始見證／本輪未重驗；移除私有位址、製程名稱與group/tag/表執行識別，縮限不遺失／無條件續傳保證 |
| `docs/technical/windows-tray-runtime-logs-evidence.md` | 保留歷史checkpoint及48情境對照，另加本輪真實結果；不把歷史數量或跨平台編譯當當次驗收 |
| `openspec/changes/add-windows-tray-runtime-logs/proposal.md` | 保留原Modified Capabilities→Impact的結構釐清，無產品語意變動 |
| `openspec/changes/add-windows-tray-runtime-logs/tasks.md` | 原3.3勾選以本輪logs unit/embedded證據重新核對；追加Mac驗證任務，Windows未完成不勾、不archive |

原始12檔另做本機私有備份（repo外，權限0600），保留去識別前的使用者紀錄；不提交該備份、快取、build/static、node_modules、test-results或原始私有執行資料。正式證據只含隔離測試輸出與SHA256。

## 實際修補與保護

舊main只回第一個必填問題；原草稿已彙整missing／empty各類，但混合時仍漏掉empty。新增4個失敗案例先重現，再最小修復。現在缺少在前、空值在後，同一error皆可見；只有一類時保留原文字。range/enum仍在required問題排除後回報，FATEK schema仍只要求其既有mode/station_no；沒有藉此新增host/serial欄位要求或重寫型別validator。TestDraft HTTP 400與既有error envelope、Create/Update分類不變。

3個logs瀏覽器情境使用正式embedded UI與loopback隔離backend，沒有設備或SQL目的地設定：斷線注入是真fetch取消、恢復須真handshake；Clear前後獨立GET核對原instance/sequence仍在；source/query必須正向命中真retained startup紀錄後再確認空結果。首次offline測試因Chromium保留既有SSE而未注入成功，第二次因完整locale文案不符失敗；均是測試作者問題，修正後focused及full suite通過，未歸咎產品或弱化契約。

## 本輪實跑

| 命令／範圍 | 結果 |
| --- | --- |
| `go test -p 1 ./...` | PASS |
| `go vet ./...`／`golangci-lint run ./...` | PASS／0 issues |
| `go test -race -p 1 ./internal/apphost ./internal/config ./cmd/test_ui ./internal/diagnostics ./internal/api ./internal/api/handlers ./internal/datalink/device -count=1` | 7套件PASS；Mac分支，不含Windows-only檔 |
| `go test ./internal/api -run '^TestRuntimeLogIPv6Socket$' -v -count=1` | 真::1 listener PASS、無Skip |
| `cd frontend && npm run lint` | PASS |
| `cd frontend && npm test -- --run` | 211檔／1344項PASS，logs focused 27項PASS |
| `cd frontend && npm run build` | tsc＋Vite PASS；既有chunk size提示保留 |
| `cd frontend && npm run test:e2e` | 18項PASS；embedded focused13項PASS |
| `python3 scripts/tests/test_windows_build_contract.py` | 5/5 PASS，僅script/dry-run，不是原生Windowsbuild |
| `make check-lines`／`git diff --check` | PASS；>300行warnings保留，無新>500檔 |
| Spectra analyze／validate | aggregate無Critical/Warning；Windows無Critical，保留20 traceability warnings；兩案validate PASS |

Mutation結果：aggregation只回第一欄時6個原稿子案例失敗；mixed before fix有4個RED；jitter去掉時2項失敗；source參數省略時browser失敗；alias順序與config env優先序的overlay均產生真正assertion失敗。全部mutation已復原，production logs/config/host檔無diff。

[Manifest與source/log指紋](evidence/pending-work-20261023.json)；[真embedded zh-TW畫面](evidence/pending-work-20261023-logs-zh-TW.png)。完整log經repo路徑及暫存私有目錄去識別後壓縮保留，不把setup/sandbox失敗當行為RED或PASS。

## Change與剩餘界線

[`aggregate-device-required-errors`](../../openspec/changes/archive/2026-10-23-aggregate-device-required-errors/proposal.md) 的3項tasks已依真實驗證完成並由CLI單次封存，一項新增requirement，原38項datalink-api requirements逐一比對不變，無cleanup warning。`add-windows-tray-runtime-logs`仍active；原20項中10完成、10未完成，新增Mac收尾3項完成後合計13/23；不改Windows驗收結論。

仍需Windows：1.1新DB native確認/ACL、1.2無console early fatal、2.2 reparse/實際disk fault、4.1 OS owner/IPC並行及owner death、4.2 tray/Explorer/keyboard、4.3 native退出／等待／force與同artifact durable恢復、5.1 PowerShell與GUI/console standalone、5.2最終跨平台gate、5.3互動portable/headless handles整套、5.4依5.3的最終交接。沒有Windows可執行環境，這些保持NOT RUN而非PASS。仍有歷史48情境矩陣的未對照項，不能只靠source存在勾選。

本輪未做Linux原生、真人新手、現場PLC、正式DB或部署。Git commit、remote SHA、CI以提交後實際查證回報；不能以本機PASS推論CI。
