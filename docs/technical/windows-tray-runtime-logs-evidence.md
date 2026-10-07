# Windows Tray／runtime logs 實作驗證記錄

## 來源與交付邊界

- Change：`openspec/changes/add-windows-tray-runtime-logs`
- 實作基底：`eea0b3e5b4e85bd9868cb2cbceb119dcefcbb189`
- 本輪為依已發布規格重新實作的 source；先前暫存環境失去存取的 commit 不屬本交付，也不借用其測試結果
- 最終 commit、parent、tree、artifact SHA256 以交付 ZIP 的 `DELIVERY.json` 為準，避免文件自我引用尚未存在的 commit hash
- 只建立一個本機 implementation commit；未 push、merge、archive、部署，也未存取真 PLC／正式資料庫
- 本文件的 PASS 限明列的自動化或 source review；Windows 互動桌面、macOS 原生及真正 browser 驗收不在已通過範圍

## 工具與流程

Go 1.25.5、Modern Go Guidelines 0.1.1、golangci-lint 2.5.0、官方 `@fission-ai/openspec@1.2.0` 均隔離使用，不修改系統 Go 或使用者電腦設定。前端依 repo lockfile 還原，React 19／TypeScript 5／Vite 7。

依使用者明示的原 OpenSpec apply workflow 實作；目前 repo 的 Spectra CLI 不可用，**正式 Spectra verify/review/analyze/audit/drift 為 NOT RUN**，不以普通 review 冒充。已執行兩個獨立 source review，分別聚焦 Windows／lifecycle 與 diagnostics／API／UI。

## 實作落點

| 範圍 | Source |
| --- | --- |
| 模式、資料選擇、共用 shutdown budget | `internal/apphost/` |
| Win32 tray／shell、OS owner lock、same-user/session IPC | `internal/desktop/` |
| 固定 safe events、bounded broker／file sink | `internal/diagnostics/` |
| log-only guard、snapshot、SSE | `internal/api/runtime_log_*.go` |
| lazy log viewer、transport、雙容量 model、i18n | `frontend/src/features/runtime-logs/`、`services/runtimeLogs.ts` |
| 真實 entrypoint、bind-first、readiness、safe adapters、cleanup | `cmd/test_ui/gateway_*.go`、`main.go`、`server_runtime.go` |
| 原 worker 的 completion／deadline seams | `internal/datalink/{collector,runtime,groupdelivery,grouppipeline,storage}/` |
| 完整 SPA／GUI-console 建置分流 | `Makefile`、`scripts/build.ps1` |

## 驗證結果

本節以最終收斂時的命令結果為準；原始 log、source manifest、review 結果與可重跑 smoke harness 置於交付 ZIP 的 `evidence/`。

| 命令／驗證 | 結果與限制 |
| --- | --- |
| `go test ./...` | PASS；包含新增 production seams 與原本回歸，不代表外部 PostgreSQL／現場測試已跑 |
| `go vet ./...` | PASS |
| `golangci-lint run ./...` | FAIL；基底獨立重跑有 90 項既有問題，未弱化 config／移除檢查 |
| `golangci-lint run --new-from-rev <base> ./...` | PASS，0 個新增行問題；不替代 full lint |
| 受影響 Go package `-race` | PASS；具體 package／命令见 evidence manifest |
| lifecycle 重複回歸 | PASS，24 個新低層 regression，shutdown 10 次；後續 startup cancellation race 有獨立重跑 |
| 前端 `npm run lint` | PASS |
| 前端 `npm test -- --run` | PASS，203 files／1301 tests，jsdom 不等於 browser |
| 前端 `npm run build` | PASS；logs 為獨立 lazy chunk，既有大 chunk 警告保留 |
| Windows build contract tests | PASS，5 項；未執行 Windows PowerShell native build |
| `make check-lines`／`git diff --check` | PASS；>300 行警告保留，沒有新增 >500 行檔案 |
| 官方 OpenSpec strict validate | PASS；只驗 artifact 結構／規格，不證明產品通過 |
| Linux production smoke | PASS：own embedded logs、99 個 assets 位元組、snapshot／SSE reconnect、403、literal DB path、SIGTERM 0、occupied port fail-before-DB |
| 無 Go/Node runtime PATH | Linux smoke PASS；不是 Windows portable acceptance |
| Windows GUI／console cross-build、PE 檢查 | 結果與產物 hash 見 DELIVERY；不能代替原生操作 |
| Windows native owner／IPC／tray／Explorer／headless handles | NOT RUN，無連線的 Windows 互動環境 |
| Embedded Playwright／真實 browser UI | NOT RUN；已知 browser 執行／loopback 路徑遭環境阻擋，未繞過；未提供假截圖 |
| macOS native、真 PLC、正式 DB、field/deploy | NOT RUN |

一輪合併的低層 race 曾在既有 `TestSenderSlowRemoteFailureStartsSettlementTimeoutAfterAttempt` 的 5ms settlement 預算失敗，沒有 race report；其後隔離的完整 groupdelivery race 與該測試 5 次重跑 PASS。保留原 FAIL 與 rerun，未放寬 timeout、刪測試或改資料語意。

## Review 發現與修正

- Shutdown refresh 可覆蓋 stopping：以同一 state lock 串行 native projection，quit／coordinator 狀態在鎖內重查；controlled interleaving regression 先 RED 再 race GREEN
- Quit 後 HTTP gate 尚未關閉：admission 直接觀察 quit，既有 handler 仍須真正結束後才 close DB
- file-only 在 logger restore 後洩漏 raw fatal／遲到 ready 或 opener error：command exit 與晚到 producer 保留原選定 policy；新增 runCommand 與 post-shutdown regression
- Native post-start fault 無 safe event／錯誤退出：固定 callback 先記錄穩定 code，再請求 quit；handled sentinel 與晚到 fault 不掩蓋非成功結果
- Startup completion 等待 browser：完整發布 HTTP resource fields 後即結束 startup phase；native opener 最多一個未完成工作，主流程可取消等待，不讓 ShellExecute 卡住 worker cleanup
- Windows canonical 祖先可被重新綁定：owner 保留 root-first directory handles，拒絕 reparse、deny write/delete sharing，整鏈固定後重驗 FileID；原生實測仍待補
- Hardlink DB 與 WAL 名称不是同一個 identity：可丟棄 crash probe 證實 alias 可能看不到原 WAL 的已接受資料。Windows writable owner 對無法證明同一 sidecar namespace 的 multi-hardlink DB fail closed；不搬移／刪 WAL。參考 [SQLite 官方說明](https://www.sqlite.org/howtocorrupt.html#_multiple_links_to_the_same_file)
- UI metadata 原先可額外保留 snapshot records：移除這份 hidden array，確保 browser 雙容量上限；file loss interval 與較晚 admission drop 以 sequence 正確排序

Producer inventory：standard log／slog 任意文字只投影固定 `raw.suppressed`；Gin access 用 route template 等允許欄位，panic recovery／HTTP ErrorLog／startup、runtime error callback、native fault、shutdown 有固定 code。未 tee 任意 stdout/stderr；OS/runtime fatal、直接 `fmt` 或第三方自行寫輸出不保證捕捉。舊 `/debug/logs` 與 raw console／query formatter 分開保留。

## 15 項任務：4 完成，11 保留

分類：A＝尚有已知 source／unit 缺項；B＝已實作並有部分實測，但該整項驗收仍缺平台／環境或完整 gate；C＝整項要求的目標平台驗收未跑。勾選以完整任務為單位，不把 source 存在當成全部驗收完成。

| Task | 狀態 | 尚未勾選原因／已具證據 |
| --- | --- | --- |
| 1.1 | B | mode／alias／ambiguity／無寫入 resolver 自動測試 PASS；Windows native 確認取消、真權限與 macOS CLI 未驗 |
| 1.2 | B | bind-first、own HTTP ready、assets、degraded、tray failure、native fault seams PASS；GUI 無 console 的 early fatal 原生呈現未驗 |
| 2.1 | 完成 | secret／未知 raw／8KiB／雙界線／sequence／race PASS |
| 2.2 | B | rotation、alias、disk error/stall、recovery、flush 自動故障注入 PASS；Windows reparse／native disk fault 未驗 |
| 2.3 | 完成 | production adapters 與必要 safe producers 已接線，secret／file-only／late producer／舊 console 回歸 PASS，未捕捉範圍已記錄 |
| 3.1 | 完成 | 實際 IPv4／IPv6、本機 guard、Host／Origin／Fetch Metadata／proxy／CORS 負測試與 bounded snapshot PASS |
| 3.2 | 完成 | replay-live boundary、progress／gap／reset、上限、stalled transport deadline、cancel／race PASS |
| 3.3 | B | lazy route、控制、XSS、i18n／鍵盤 unit、完整前端套件 PASS；embedded Playwright 未跑 |
| 4.1 | B | OS-backed owner／ACL IPC source、portable contract／cross-compile 有證據；Windows concurrent launch／owner death／hardlink failclosed／session 測試尚未原生跑 |
| 4.2 | B | native tray／shell／Explorer source 與注入 seams；無 console 閃窗、keyboard/overflow、browser/Explorer 原生驗收未跑 |
| 4.3 | B | once／shared deadline／真 completion／HTTP gate／可丟棄 SQLite durable reopen/race PASS；native cancel/confirm/force/session 行為未驗 |
| 5.1 | B | scripts／portable contracts／embed smoke／cross-build；Windows fresh-machine PowerShell／standalone 實測未跑 |
| 5.2 | B | 明列的大部分自動 gates PASS；full lint 保留基底 FAIL，Playwright 未跑，不能宣稱全綠 |
| 5.3 | C | 本項整套 Windows 互動桌面／GUI-headless supervisor／同 artifact durable/browser acceptance 未跑 |
| 5.4 | B | 使用與交接文件、scenario matrix 已整理；依賴 5.3 的最終現場交接尚未完成，不 archive |

目前 review 未保留已知可在此環境修正卻未修正的 A 項；这不代表 B/C 的平台行為已被證明。

## 可重跑與後續驗收

一般 build/run、模式、資料選擇與回退請看[使用交接](windows-tray-runtime-logs.md)。交付內的 `implementation.patch` 可在上述 base 使用 `git am`；`implementation.bundle` 是帶 base prerequisite 的 Git bundle，先核對 base 再匯入。套用前保持工作目錄乾淨，不把 patch 套到未知 production tree。

Windows 驗收須使用 DELIVERY 指定 artifact 或從同 commit 按完整 build script 重建，逐 scenario 記錄 OS、session、exe hash、設定／可丟棄 DB、步驟、觀察與截圖。至少涵蓋 owner alias／ancestor pin 與 WAL reopen、Explorer restart、磁碟故障、正常與逾時退出、GUI-headless cmd/PowerShell/supervisor handles。不可拿 cross-build 或 Linux smoke 勾掉 5.3。

48 個 scenario 的逐列 source／test／限制對照在交付 ZIP：`evidence/runtime-log-scenarios.md`（26 項）與 `evidence/windows-lifecycle-scenarios.md`（22 項）。沒有修改 Studio inventory 文件，故本次不觸發其 changelog SQLite 更新。
