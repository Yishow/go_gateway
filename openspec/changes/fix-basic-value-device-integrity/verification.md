# 驗證：基本數值與設備身份

## 完成與需求對應
7/7 tasks；2/2 requirements；6/6 scenarios 具具體測試與本次執行證據；遵循 design，無未修正 Critical/Warning。

| 需求與實作 | 情境對應測試 | 具體值／執行 |
| --- | --- | --- |
| Source preserving exact mapping conversion：mappingDefaults、workspace mapping builder、checked cast/scale、workspace confirmation、source lifecycle、runtime event wire | Default uint64 reaches SQL exactly → `TestProductionDefaultUIUint64ToSQL` 加 F default UI→poller→journal/outbox→SQLite；不是只測 codec | 9007199254740993 完整 SQL readback；本次 required Go＋F passed-current |
| 同上：非法與縮放邊界 | Invalid conversion fails → `TestExactCastRejectsInvalidValues`；Compatible checked scaling → `TestWorkspaceScaleMatchesDeclaredTarget`、`TestStoredPipelineKeepsPublishedOrder`、`TestProductionDefaultUIScaledInt16ToSQL`、mappingDefaults suite | NaN/Inf/負unsigned/overflow/fraction/2^53+1→float/1e-400/float32 1e-50 拒絕；int16 243×0.5+10=131.5 合法；stored order 不改；本次 Go/Vitest passed-current |
| 同上：只有成功明確 Save 才接受 actual hash，source proposal 另存 | Confirmed save survives source synchronization → `TestWorkspaceConfirmationOneMappingRestartAndSourceEdit`、`TestProductionWorkspaceConfirmationOnlySuccessfulOwnedSave`；SQL CAS 檢查 link/revision/tag/pipeline；failed audit/revision race 不接受 | 一個 confirmed override、另一個 manual edit、兩次 sync、後續 source edit、500 failed confirmation、409 stale revision；本次 Go passed-current |
| 同上：exact wire/bad quality | Exact live value and failed conversion quality → `TestRuntimeValueWirePreservesUnsafeIntegers`、`TestRuntimeInvalidCastDoesNotPublishGoodValueOrWrite`、LivePointsTable decimal string test | >JS safe integer 用 decimal text；invalid cast transformed nil/quality bad，無 target/history good write；本次 Go/Vitest/F passed-current |
| Device scoped live point identity：LivePointsTable selected device + exact point + event device，metadata 同身份；移除 address fallback | Same address on two devices → `live-point-device-identity.test.tsx` 三案例；F 真 Runtime A→B→A 切換 | A215/B187、外來相同 point key、waiting、名稱/type/unit、確切大數 decimal；本次 Vitest/F passed-current |

## 失敗回歸、review、修正再回查
- 開始在基準重現 default float64 與 numeric conversion 錯值，再加 default UI→SQL，不以 codec-only 代替；mapping/device Vitest、checked cast Go、production-chain 回歸先失敗後通過。RED log hash／failure excerpts 在 validation JSON。
- Review 發現合法 scale 若只保留 int16 會錯拒 fraction，因此新增未改 completed history 的 compensation task：neutral preserve、explicit numeric scale default float64、新 pipeline scale→checked target cast；既有已存 pipeline 保持步驟順序。
- Review Warning：accepted pipeline hash 在第一次 sync 被 proposal 覆寫，第二次 sync 失效；`signature-review-red.log` 先失敗，改保存 original accepted hash 後 `signature-review-green.log` 通過。
- Review Warning：nonzero underflow 變0；`cast-underflow-red.log` 先失敗，checked cast 拒絕 silent zero，full Go 通過。
- 四鏡檢視 correctness、efficiency、reuse/simplification、conventions，及 apply audit 的 Scoundrel/Lazy/Confused sequential lenses 已覆蓋；不宣称獨立 reviewer。無剩餘可證明缺陷。
- Named review snapshot `b522b0655dc6f4d32424094eb0f8a9640f9ce208`，26 files，`touched_tracking` / resolved，snapshot check 在當次 review 結束為 current（commit 前）；後續完整範圍以 aggregate 為準。三份 preexisting_dirty 都是 task baseline 前先寫的本任務 RED tests，hunk attribution 限制保留；原 checkout 非使用者 dirty。完整檔案清單以保存的 scope 及最後 aggregate review 為準。
- 視覺小樣式不另測 spacing/radius；跨 caller/default/error/conditional contract 均測。spec 未另含 Example blocks；scenario 裡明寫的 9007199254740993／A215/B187 已在 tests 出現，沒有虛构 Example coverage。

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
