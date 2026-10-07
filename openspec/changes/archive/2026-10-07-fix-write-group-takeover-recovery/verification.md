# 驗證：持久接管與故障列處理

## 完成與需求對應
7/7 tasks；2/2 requirements；7/7 scenarios 具具體測試與本次執行證據；遵循 design，無未修正 Critical/Warning。

| 需求與實作 | 情境對應測試 | 具體值／執行 |
| --- | --- | --- |
| Durable legacy writer takeover：HasAppliedOwner 讀既有 immutable write_group_versions，Owns 與 worker/boundary 分離；read failure fail closed；首次成功 Apply 即接管，沒有新表 | Retire and restart after migration；Delete or change destination after takeover；No effective takeover → `TestProductionLegacyTakeoverSurvivesLifecycle`，真 migration Preview/Review、enabled legacy mapping、real writer、production SuppressMapping seam | draft legacy第1列／failedApply第2列合法；successfulApply後 disable/delete/changeconnector/retire/reopen 原A仍2列、新B0，legacy不恢復；本次 Go passed-current |
| Traceable scoped delivery head resolution：group scoped API、Attention oldest20、原子 audit/state CAS、reason/confirm skip、frozen payload不改；UI Basic/Advanced repair | Repair blocked head → `TestOperatorResolutionRealFaultsKeepFrozenIdentityAndDedupe` missingtable/SQLite query_only permission×retry/skip；F 真 missingtable後修復Basic retry | 同 effect/record/payload/hash/connector/table/bucket 不變；retry SQL exactly1，下一列可寫；本次 Go＋F passed-current；live PG permission未做 |
| 同上 | Skip poison row → 同 real-fault matrix 的 trigger poison＋F UI explicit checkbox/skip | skip目標0列、operator_skipped、audit exactly1、後续列可寫；本次 Go/Vitest/F passed-current |
| 同上 | Unsafe or stale resolution → `TestProductionDeliveryOperatorResolution` foreign workspace/group、digest/state/claim epoch stale、unknown拒絕、audit failure rollback；UI unknown/readonly/cache-error無操作 | 400/404/409安全generic errors，unknown不重送；本次 Go/Vitest passed-current |
| 同上 | Repeated decision and new blocked attempt → API exact same decision audit read、changed body409、新blocked epoch409旧decision；`groupDeliveryRecovery.test.tsx` uncertain ID reuse＋newepoch UI reset | retry只一次 audit，newepoch可fresh reason、inFlight防雙click；本次 Go/Vitest passed-current |

## 失敗回歸、review、修正再回查
- `takeover-red.log` 重現 retirement後legacy多第3列，再 production-chain GREEN；不是靜態推測。
- resolver/recovery UI先RED→GREEN；fault測試真 SQL sender、ownedDB缺表/query_only/poison、修復後 receipt/readback，沒有 fake sender 代替。
- Review Warning：blocked→retry→blocked同狀態可接受舊決策，且 UI resolved永遠disabled。`resolution-epoch-red.log` 原期待409得200、`resolution-epoch-ui-red.log` 原新attempt按鈕disabled；改沿用 existing claim_epoch CAS 與 UI epoch key，backend/UI GREEN，再fullrequiredtests通過；無新 epoch schema。
- DecisionID duplicate read在state/epoch mutation guard之前，只讀同scope/samecontent atomic audit，uncertain reply不盲目重送；unknown fresh決策仍拒絕。actor是本機既有 API provenance `local_api`，未冒稱 authenticated human 或新增authmodel。
- 四鏡 correctness、efficiency、reuse/simplification、conventions，加 apply audit 三 lenses sequential 檢視 frozen identity、stale CAS、duplicate、unknown、safe error。generated Swagger由cached official swag工具產生，未手改。
- Named review snapshot `eea59cadc141c3a4e7b56b4626b1d1a0689eddab`，24 files，touched_tracking/resolved，當次 check current（commit 前）；三份 preexisting_dirty 是本任務 task baseline 前的 RED takeovertests/implementation，不是使用者工作。hunk attribution限制保留；finalaggregate含各layer。
- spec無額外 Example blocks；具體fault/stale/receipt例值出現在上述test；視覺spacing不另測，但所有conditional/default/errorbranches有test。

## 執行證據與範圍
- 基準：`d41638b0e2c0a49b51bbeb21f31ccda83f20199b`；原 checkout 初始 clean，工作 branch `fix/basic-recording-integrity`。所有本批差異均任務所有；只 local commits，不 push/merge/archive/deploy。
- 主對話先 review 三份窄範圍草案，再在既有授權下 apply。未新增 protocol、driver、daemon、報表、能耗、事件多 tag 模型或平行持久化 authority。
- `docs/plans/studio-v2-flow-completion/evidence-f/basic-integrity-20261005-validation.json` 保存命令、結果節錄、原始 log SHA256、66 個 source/test/harness 雜湊及 config 雜湊。F c 之後只有 SummaryRail renderer、connector en/zhTW副標、該回歸及驗收 harness 調整；mapping/schema/start/runtime/recovery/lifecycle 的 source/test 與 c build 相同（basicRecordingPanelHelpers只去除EOF空行，無行為差異）。Go required gates 是本次實際 passed-current；frontend於sidebar修正後再次完整重驗，並重新build/更新embeddedbundle，非沿用舊圖宣稱新UI。
- 必要檢查：`GOMAXPROCS=2 go test -p 1 -parallel 1 ./...`、`go vet -p 1 ./...`、`golangci-lint run --concurrency 1 ./...` 均 exit 0；frontend `npm run lint`、`npm test -- --run --maxWorkers=1 --no-file-parallelism`（194 files / 1232 tests）及 `npm run build` 均 exit 0。`git diff --check`、`make check-lines` 通過。Build 的既有 chunk-size warning 保留，不當 failure 或放寬 gate。
- 輕量 Node harness verification/timing/cleanup 12 tests 通過。新的 F 是 UI-first：實際 production binary、loopback Modbus 模擬器、owned disposable SQLite；配置與 schema/start/retry/skip/disable 操作走真 UI，API 只 read-only observation；目標 SQLite DDL 僅本任務故障注入/修復，沒有注入 INSERT 或 mock route。
- F 成功證據：`docs/plans/studio-v2-flow-completion/evidence-f/basic-integrity-20261005c.json`，包含每設備8測量點（按型別 register span）、各3個60秒 buckets、預設型別不改、確切 SQL/receipt/outbox、A215/B187/A215 真切換、缺表 retry、poison explicit skip、後續列及 disable/restart。390/768/1440 與 runtime/recovery/restart screenshots 由 c UI 產生；該輪 sidebar 舊5s/card 由最終review發現並補償修正，c圖不代表最後summary文案。另 `basic-integrity-20261005-ui-final.json` 使用新embeddedbundle實際UI配置一個模擬point＋connection，驗證sidebareffective設定不再誤導、Basic60s、target未建立，保存最新1440 screenshot；不重跑內容完全相同的8點等待/故障流程。
- 失敗 F `...20261005a.json` 是 SQLite locator 多重匹配；`...20261005b.json` 是 `dual.1` substring 撞 `dual.10` 的驗收誤判。b 先 stop gateway 再 close runtime SSE 造成 shutdown 超時與 harness SIGKILL；改先 close browser 再 normal stop，未改 production shutdown。失敗 proof 不改判成功；本批原始 JSON 為 generated run evidence，依 AGENTS 產物規則只加本批 filename pattern 的行數豁免，未豁免來源碼。
- 僅在主對話協調的重壓時段串行低 workers，未啟動 Docker、未停止 Qycms 或其他任務程序。成功 F children 正常退出、browser closed、owned namespace removed；失敗 b owned namespace 另行確認無活程序後清理，保留原 proof 與補充 cleanup evidence。
- 未執行：真 PLC、Windows 執行、LAN、長跑/soak、live PostgreSQL 或正式 DB；上述是環境驗收限制，不能由 macOS 模擬器測試推論。
- 實作 exact HEAD `a4deda5b25560d2854b3d91245152604128ca0a6` 的 aggregate code review 已完成；最後 completion-document commit 後再確認完整 base→exact HEAD snapshot 與文件差異，最終回報給出該 HEAD，避免文件self-reference改變HEAD。

## 最終回查限制
全部 selected source/test/config/generated Swagger 與文字證據已在 explicit base→實作 HEAD＋current layers 回查；最後 docs-only commit 後再 check snapshot。PNG 在 CLI 標 binary/non-inspectable；另透過 view_image 實際核對 c Basic 390/1440、runtime、b error 與 ui-final 1440 圖，其餘 PNG 不宣稱逐圖內容驗證。沒有獨立多agent review 或 CI/Windows 執行；最後 exact HEAD 與 snapshot ID 以交接回報為準。
