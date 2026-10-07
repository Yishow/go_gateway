# 來源與規劃邊界

調查日期：2026-10-07。GitHub main及獨立fresh clone同為 `ffdaa002ad0d02961bc239c8b0772da7b83be848`；branch為 `spec/windows-tray-web-logs`。所有下列source判斷均為唯讀檢查，不是此次產品測試通過紀錄。

## 現行 production

| 觀察 | 原始來源／符號 | 對規格的影響 |
| --- | --- | --- |
| 入口embed SPA並啟runtime、group pipeline | `cmd/test_ui/main.go` 的 `runGateway`、`wireGatewayServices` | 重用現有核心，不另造desktop server |
| 真正SQL target、Share fanout與group ownership已接線 | `cmd/test_ui/service_wiring.go`、`target_writer.go` 的 `newProductionTargetWriter` | 不宣稱database writer仍是stub；tray不碰delivery truth |
| group pipeline startup recovery先於sample intake，Stop保留accepted durable資料 | `internal/datalink/grouppipeline/pipeline.go` 的 `Start/Stop`；`internal/datalink/groupdelivery/worker.go` 的 `Stop` | 退出不刪journal/outbox/receipt，不承諾backlog全送完 |
| runtime.Stop含 `wg.Wait`，pipeline.Stop先等loop再交worker deadline | `internal/datalink/runtime/service.go` 的 `Stop`；`pipeline.go` | 現有5秒參數不代表整條shutdown已bounded；需timeout phase與取消驗收 |
| Windows build明確移除GUI flag | `scripts/build.ps1` 的backend build；`Makefile` 的build-backend | 需要desktop/console建置契約，不是改副檔名 |
| 開browser用cmd/start、先印URL再ListenAndServe、失敗Fatalf | `cmd/test_ui/server_runtime.go` 的 `startServer/openBrowser/handleSignals` | native shell避免閃窗，bind/readiness在先，錯誤傳回coordinator |
| dotenv及process環境變數設定；HOST空、CORS wildcard | `internal/config/config.go` 的 `Load/newConfigFromEnv/GetServerAddr` | 不把舊Viper context當現況，不全域改HOST掩蓋log風險 |
| DB path三種aliases，default datalink.db是相對路徑 | `cmd/test_ui/harness_config.go` 的 `embeddedSQLiteDSN`；`internal/datalink/db.go` 的 `DefaultEmbeddedSQLiteDSN` | desktop資料根目錄變更須辨識legacy/新庫，不silent migration |
| 無登入middleware，只有logger/recovery/CORS | `internal/api/router.go` 的 `NewRouter/corsMiddleware` | 新log guard是local-only，不假稱已有user auth |
| runtime SSE明設wildcard | `internal/api/handlers/runtime_stream_handler.go` 的 `Stream` | 不直接複製現有SSE存取政策到log feed |
| LOG_*只有config欄位，main僅SetFlags | `internal/config/config.go` 的 `LogConfig`；`cmd/test_ui/main.go` | managed logger、file rotation均是待實作能力 |
| startup印SQLite DSN；HTTP印完整query/raw error | `cmd/test_ui/main.go`；`internal/api/router.go` 的 `formatHTTPLog` | 新sink必须先safe投影，不raw tee |
| 既有測試要求完整query字串 | `internal/api/router_logger_test.go` 的 `TestFormatHTTPLog_PreservesFullQueryString` | CLI相容與新safe projection分开驗證，不能稱既有已redact |
| 工程Connect印headers及Config；debug details也保留Config | `internal/api/handlers/test_connection_handler.go` 的 `Connect` | 不把terminal整體或debug history複製進新log頁 |
| `/debug/logs`為RecordLog建立的工程紀錄 | `internal/api/handlers/debug.go` 的 `RecordLog/GetLogs` | 原endpoint意義及clear保持，不當全程序logger |
| standard log、slog、Gin與fmt各有輸出 | `internal/datalink/dbtarget/writer.go`、`internal/datalink/delivery/worker.go`、`internal/virtual/server/modbus/udp_server.go`、`internal/api/router.go` | 實作要有明確producer coverage inventory、未知raw省略政策 |
| setup/runtime/test/gateway各有既定產品角色 | `frontend/src/App.tsx`、`internal/web/embed.go`；inventory入口/context/CURRENT_STATE | logs新增route，不復活legacy `/studio`或重寫Studio |

已讀 `AGENTS.md`、`CLAUDE.md` 及studio-surface-inventory要求的三份入口。Inventory與release歷史測試數只供脈絡，現況以本次固定SHA source為準；本案不修改inventory，故不新增其changelog資料。

## 規格重疊與唯一責任

- fixed main 的 `openspec/changes/` 只有archive目錄，無其他tracked active change；fresh clone無本機parked metadata，不能推論其他電腦沒有私有parked工作
- 查閱 `embedded-frontend-delivery`、`start-script-contracts`、`runtime-dashboard-truthfulness`：保留embedded assets、CLI/start manager及truthful runtime的既有要求，本案只新增desktop lifecycle与runtime-log-viewer
- `start-script-contracts` 的run lock是start manager/repo-port契約，不足以替代portable exe的同DB owner guard；不把兩個lock當相同authority
- 近期 `94a27a2` 已封存四項變更且 `eff17a2` 已在main歷史，不能沿用舊未合併分支假設。已接線的durable group能力以current source承認，不重做其資料語意

## 官方技術來源

- [Go linker文件](https://pkg.go.dev/cmd/link)：Windows `-H windowsgui`決定GUI而非console subsystem；本案仍要求Windows實機驗收，link flag本身不能證明無shell閃窗
- [Microsoft notification area](https://learn.microsoft.com/en-us/windows/win32/shell/notification-area)：tray由Windows Shell管理；規格要求文字狀態及shell重建恢復，不保證使用者的icon釘選位置
- [Microsoft ShellExecuteW](https://learn.microsoft.com/en-us/windows/win32/api/shellapi/nf-shellapi-shellexecutew)：可用native shell開啟固定網頁目的地；不執行任意使用者輸入command

## 本次工具選擇

- 使用者明確要求OpenSpec。repo現有Spectra skill要求缺CLI時停止；本機查無 `spectra`／`openspec`，未以手造格式冒充該workflow
- 經此次明確同意，在repo外臨時目錄安裝官方npm `@fission-ai/openspec@1.2.0`，實際version=1.2.0；來源為 [Fission-AI/OpenSpec](https://github.com/Fission-AI/OpenSpec)
- 依repo歷史原版workflow：`e8a7ba036c417b46c7c8b0ab885b9012460560b5:.agents/skills/openspec-propose/SKILL.md`，僅讀取，不恢復／修改tracked skills
- 實際使用new change→status→instructions proposal/design/specs/tasks的順序；使用者要求繁中規劃，spec requirements/scenarios為英文；未跑init/update、未安裝Spectra/spxa或更改repo工具設定
- OpenSpec1.2.0不提供Spectra analyze；跨artifact review另行記錄，不把自製检查或strict validate冒充Spectra analyze／產品實作驗證

## 未執行的產品驗證

本次未修改或啟動production code、未build前端／exe、未跑Go/Vitest/Playwright、未操作Windows桌面、設備、正式DB或deployment。規格中的timeout、byte/record/client caps皆為設計條件，需由tasks測量驗證，不是既有benchmark結果。

## 本次文件檢查與審查

- 官方 `openspec status --change add-windows-tray-runtime-logs --json` 的proposal/design/specs/tasks均done，applyRequires=tasks；這只證明文件齊全，15個implementation tasks全未勾選
- 官方 `openspec validate add-windows-tray-runtime-logs --strict --json` 通過，issues=[]；修正review後已重跑
- `git diff --check` 通過；因新增檔尚未tracked，另以逐檔新增diff與文字檢查檢查全部新增內容。相對連結、placeholder、task引用與依賴無環均核對
- `make check-lines` 通過但實際checked=0，因現有 `.line-limit-ignore` 明確忽略 `openspec/*`；另逐檔核對都低於300行，不將ignore誤報成完整gate覆蓋
- 獨立source/security review發現log檔碰撞／protected path風險，已補每DB namespace、整組排他ownership及alias/reparse安全失敗；並補Fetch Metadata精確值、bounded replay與不持鎖IO、filtered watermark及logging config現況。複核無待修blocking/warning
- 父review補明確headless override與legacy console default差異，以及Windows GUI exe實際redirect／wait／exit驗收；已同步design/spec/tasks
- 發布前重新fetch及GitHub connector雙重核對main仍為本檔來源SHA；實際發布SHA、遠端tree範圍與CI狀態於交付時另外核對，不預填未發生的成功
- Spectra analyze未執行；Go／前端／Windows／SQL及field驗收亦未執行。文件一致性與官方結構驗證不能替代這些測試
