# 驗證：Basic 完整準備與單一生效設定

## 完成與需求對應
7/7 tasks；2/2 requirements；4/4 scenarios 具具體測試與本次執行證據；遵循 design，無未修正 Critical/Warning。

| 需求與實作 | 情境對應測試 | 具體值／執行 |
| --- | --- | --- |
| Complete basic preparation with one effective configuration：BasicSchemaPreparation重用GroupSchemaPanel/ensure local group；ConnectorSection只connection；groupauthority保留 | Fresh basic recording → basic-recording-preparation首案例＋F fresh-dual actualBasic preview/confirm/start，不須Advanced；default-point SQL在另一change驗證 | 兩設備各8points，型別registerspans正確，60秒每組一列×3；F preview前後readonly確認無managedtable，explicitApply才建表；本次 Vitest/F passed-current |
| 同上：BasicProgress/Evidence、GroupDeliveryStrip／schema diagnostics；無重建UI | Inspect effective settings → Basic準備test第3/4、Step4/TestPage callertests、schema-operationreadonlytest、groupDeliveryStrip | table/writeStrategy/interval舊controls absent；主面names/value/typeunit/stages，IDs/revisions/receipt diagnostics可展開；本次 frontend1232test/lint/build passed-current |
| Basic preparation retains schema and capability safety：reuse unchanged backendtokens/ledger；Basic scopekey unmount+active guard；移除mainline nonpreview action | Cancel or invalidate basic preview → Basic deferred-device-scope regression；existing groupSchemaPanel 14guards；backend recordingplan/schema_preview token guards（exact/stale/expired/foreign/unknown/tampered）；schema recovery unknown保持 | 新Basic render無DDL、readonly不ensure、切設備旧response不apply；取消離開preview無DDL；本次 frontend/Go passed-current，F明確認前無表 |
| 同上：supportsWriteGroupKind SQLite/Postgres，KindSelector與savedpool禁unsupported | Unsupported group kind → Basic準備第3、Step4controls/Step4callers；saved/new MySQL/SQLServer disabled | en/zhTW safe support說明；無driver/API擴充；本次 Vitest passed-current |

## 失敗回歸、review、修正再回查
- `basic-ui-red.log` 新Basic preparation、single effective controls、capability/diagnostic tests先失敗，`basic-ui-green.log`通過；scope-race後補deferred device regression，本次fullVitest再次通過。
- 合約只preview/confirm/start；Advanced既有testWrite保留，沒有新增Basic testWrite需求。舊nonpreview engineering component仍由獨立tests涵蓋，但mainStep4不提供會409的action。
- fullGo/Vitest初次舊期望失敗已逐項核對後修正：ignored cast param `type`→正確 `target_type`；兩個source proposal測試改assert不覆寫publishedTagType（保持原pipeline）；兩個Step4test改依新Basic操作/readonly。未刪canonicaltoken guards；deferredlegacy mainline test由newBasic同scope race取代。
- 最終截圖 Review Warning：SummaryRail仍將 legacy connector.table/5s/targets 呈現為生效DB寫入，且connector副標仍提策略。新增compensation task1.5，保留已完成歷史；`summary-rail-red.log` 重現、GREEN20focused tests（SummaryRail/Shell/Basic）；改側欄只connection name/kind＋group設定位置，兩語副標只連線，frontendrequiredgates重跑且核對新bundle真UI，原c截圖限制明列。
- BasicPanel純helpers抽取使497行降至396行；無新架構/狀態authority。SchemaPanel既有confirmation/recoveryGuard沿用。Review四鏡與applyaudit三lenses已檢查；無剩餘高可信缺陷。
- Named review snapshot `f789af7b7c0a4dae9fa7449e7f80ae58d70c6f34`，14files，touched_tracking/resolved，當次 check current（commit 前）；sharedlocale兩preexisting_dirty是本任務recovery先寫，非使用者dirty。finalaggregate另涵蓋newF scripts、scope regression、generatedproof及行數ignore。
- F修正locator需完整label、exactrowname比較避免substring（dual.1與dual.10），先close SSEbrowser再ownedchildren normalstop。Config寫入仍全走UI；觀察API無寫入；onlyfaultDDL受限owned target。
- spec無Exampleblocks；八points／支持SQLite+Postgres／unsupported兩kind例值在tests與F；單畫面spacing/radius不另測，conditional/default/error跨caller契約測試保留。

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
