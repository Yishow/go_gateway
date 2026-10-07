# Windows 系統列與執行日誌使用交接

本文件說明 `cmd/test_ui` 新增的 Windows 桌面模式與本機執行日誌，供建置、操作與驗收交接使用。產物內嵌後端、SPA 與 tray 圖示；設定、SQLite 與診斷檔仍是外部資料，不會封存在 exe 中。

**Windows 原生桌面及嵌入式瀏覽器驗收目前均為 NOT RUN。** 原始碼、單元測試或跨平台編譯不能證明無 console 閃窗、單檔現場可用或正式環境就緒。請以[本次驗證證據](windows-tray-runtime-logs-evidence.md)的命令、版本及未驗清單為準，不把本文件當作上線核准。

規格來源：[提案](../../openspec/changes/add-windows-tray-runtime-logs/proposal.md)、[設計](../../openspec/changes/add-windows-tray-runtime-logs/design.md)、[Windows lifecycle](../../openspec/changes/add-windows-tray-runtime-logs/specs/windows-desktop-lifecycle/spec.md)、[runtime log viewer](../../openspec/changes/add-windows-tray-runtime-logs/specs/runtime-log-viewer/spec.md)、[任務清單](../../openspec/changes/add-windows-tray-runtime-logs/tasks.md)。原提案記載的是規劃階段；本交接描述目前實作，驗收狀態另由證據文件記錄。

## 建置與啟動

建置機須有 `go.mod` 指定的 Go 工具鏈，以及符合 `frontend/package.json`／lockfile 的 Node.js、npm。以下均從 repo 根目錄執行；先還原鎖定的前端依賴，再使用完整建置流程。

Windows PowerShell：

```powershell
Push-Location frontend
npm ci
Pop-Location
powershell -File .\scripts\build.ps1 -Mode desktop -Version dev
powershell -File .\scripts\build.ps1 -Mode console -Version dev
```

- Desktop 產物：`bin/gateway-desktop.exe`，使用 Windows `desktop` build tag 與 `-H=windowsgui`
- Console 產物：`bin/test-ui.exe`，不使用 GUI subsystem；省略 `-Mode` 仍是 console
- 腳本先執行前端 build、檢查 `index.html`、同步嵌入資產，再編譯 Go；前端失敗不會被當作完整新產物。失敗後不要誤用目錄中可能殘留的舊 exe
- 版本預設 `dev`；commit 自 Git 取得，工作目錄有異動會標示 `+dirty`，無法取得時為 `unknown`。交接另記錄實際 source SHA、工具版本與 exe 雜湊

有 Make 的環境可使用：

```sh
make build BUILD_MODE=desktop TARGET_GOOS=windows BUILD_VERSION=dev
make build BUILD_MODE=console
```

Make 的前端步驟目前使用 `npm install`；需要嚴格依 lockfile 還原時先執行 `npm ci`，並核對建置後 lockfile 沒有意外變更。`TARGET_GOOS` 預設為 `go env GOOS`；Linux/macOS console 建置即使檔名是 `test-ui.exe`，仍是該平台的原生程式。Windows cross-build 只提供編譯／連結證據。

不要以 fresh clone 的 embed 佔位檔直接 `go build` 後宣稱已交付網頁。完整資產是啟動 readiness 的必要條件。`start.ps1`／`start.sh` 保持既有 console 開發路徑；本案不安裝 Windows Service、登入自啟或常駐 companion。

### 三種執行語意

| 執行方式 | Tray／native 對話框 | 自動開啟瀏覽器 |
| --- | --- | --- |
| Desktop 產物正常啟動 | 使用 Windows native UI | 預設開 `/studio/v2`；`AUTO_OPEN_BROWSER=false` 或 `0` 關閉 |
| Console 產物，沒有模式旗標 | 不建立 tray | 只有 `AUTO_OPEN_BROWSER=true` 或 `1` 才開啟 |
| 任一產物加 `--headless` | 不建立 tray、不顯示互動對話框 | 一律停用，即使環境值為 `true`／`1` |

目前沒有 `--desktop` 或 `--console` 切換旗標；預設由產物決定，不依「是否從終端啟動」猜測。

下面是本機啟動範例。**先把 DB 路徑換成預期資料庫的真實絕對路徑，確認父目錄存在與可存取；不要拿範例路徑代替現有資料來源。**

```powershell
$env:GATEWAY_DB_PATH = 'C:\GatewayData\datalink.db'
$env:HOST = '127.0.0.1'
$env:PORT = '8080'
.\bin\gateway-desktop.exe
```

無互動執行及資訊查詢範例：

```powershell
.\bin\test-ui.exe --headless
.\bin\test-ui.exe --version
.\bin\gateway-desktop.exe --headless --help
```

`--help`／`--version` 不初始化資料庫或採集。Desktop 預設以 native 視窗顯示；console／explicit headless 寫入呼叫端提供的標準輸出。GUI exe 的 headless redirect、wait、exit code 在 cmd、PowerShell 與實際 supervisor 下仍須各別 Windows 實測；不能由 linker flag 推定已驗證，也不會為缺少 console handle 額外配置視窗。

## 資料位置與重複啟動

- Desktop 以 canonical exe 所在目錄作為設定根目錄；讀該處 `.env`，不改讀捷徑的工作目錄。若 exe 位於 `bin`，設定根目錄就是該 `bin` 目錄
- Desktop 先採 process 中明確的 DB 設定，再採 exe 目錄 `.env`；各來源內依 `GATEWAY_DB_PATH`、`DB_PATH`、`SQLITE_PATH` 順序選擇。相對 DB 路徑以 exe 目錄解析
- Console／explicit headless 保留 cwd `.env` 與既有環境變數／alias 解析語意，不自動移動資料。一般設定仍由 process env 優先於 dotenv
- 選定 DB 的父目錄是資料根目錄。版本資訊會在本機顯示選定 DB；一般 operator DOM 不顯示私人檔案路徑
- cwd 與 exe 目錄不同，且 cwd 有另一份 `.env`／`datalink.db`，在沒有明確絕對 DB 路徑時停止啟動。請辨識舊庫後設定正確絕對路徑；程式不猜選、不搬移，也不把另一份新空庫當作復原
- Desktop 沒有明確 DB 設定且預設庫不存在時，建立前會詢問位置，預設可取消。明確指定尚不存在的 DB 不走這個「隱含新庫」確認；設定者必須先確認意圖
- 無權限、父目錄不存在、無法取得可信檔案身分或不安全 Windows 路徑，均安全失敗；不提權或換到另一資料庫。Windows ADS、device namespace、保留裝置名、尾端點／空白等歧義路徑不作為替代解析方式

Windows 的同 DB owner 由實際檔案身分與 OS 鎖控制，涵蓋等價大小寫、hardlink、symlink／junction 路徑；不是靠 PID 檔或 port 判斷。第二個 desktop 只可向已驗證同使用者、同 session 的 desktop owner 要求固定「開設定頁」動作。Headless owner 或驗證失敗會顯示已在執行，不建立另一組 worker、不殺程序、不開啟未知服務。程序死亡後由 OS 釋放 ownership；關閉逾時時仍保持 ownership。若尚無 owner 而 DB 存在多個 hardlink，會拒絕新的 writable owner，避免改用 alias 名稱而錯過原路徑的 WAL／SHM。請保留原始 DB 及 sidecars，由資料管理者核對復原；不自行搬移、刪除或切換 alias 繞過。

不同 DB 可以使用不同實例，但必須各自設定可用 port。Port 衝突不會自動換 port 或終止佔用者。來源：[data resolver](../../internal/apphost/data.go)、[native ownership](../../internal/desktop/README.md)。

## 系統列選單與狀態

正常桌面模式沒有可最小化的主視窗。關閉或縮小瀏覽器不會停止閘道；可從系統列重新開啟。圖示可能位於 Windows overflow 區，不能保證釘選位置。

| 選單 | 行為 |
| --- | --- |
| 開啟設定頁 | 固定 `/studio/v2`；尚未 ready 時停用。左鍵雙擊同此動作 |
| 查看執行日誌 | 固定 `/studio/logs`；無可用本機日誌 URL 時顯示原因 |
| 程序／採集 | 分列唯讀文字狀態，不以顏色或 HTTP ready 代替採集成功 |
| 版本資訊 | version、commit、mode 及本機選定資料來源；缺少 metadata 不編造 release |
| 結束… | 詢問停止採集與網頁服務；取消維持運作，關閉中不重複啟動 cleanup |

程序 `starting` 表示啟動中；`serving` 表示已建立 listener、驗證嵌入資產並確認自身 HTTP 回應。`degraded` 可代表 runtime 未運轉或診斷磁碟降級，採集欄位仍獨立呈現。`stopping`／`stop-timeout` 見下節；這些狀態都不是 SQL committed、readback verified 或 backlog 已送完的證明。

Tray 先於 DB/services 初始化；listener 在開始採集前預先 bind。初始 tray、port 或嵌入資產失敗會清理並回報安全代碼。錯誤視窗只有在寫入與檔案身分確認成功後才顯示診斷檔位置，否則明說未保存。Explorer 重建 taskbar 的重新註冊路徑不重新啟動採集；原生實測仍待完成。

## 本機日誌與存取限制

從 setup 或 runtime 頁面點「執行日誌」，或在可用的本機 listener 上開 `http://127.0.0.1:8080/studio/logs`；port 依實際設定調整。API 位於 `${API_BASE_PATH}/system/logs` 及 `/system/logs/stream`，預設 base 為 `/api/v1`。若修改 API base，建置前端時的 `VITE_API_BASE_URL` 必須一致。

兩個 API 都檢查實際 numeric loopback socket 連線／本機目的位址、實際 listener port，以及以下條件：

- Host 只接受 `127.0.0.1`、`[::1]`、`localhost` 的對應 authority
- 有 Origin 時必須完全同源；`Sec-Fetch-Site` 只能缺省或為單一 `same-origin`／`none`
- 拒絕 `same-site`、跨來源、重複／異常 Fetch Metadata，以及 `Forwarded`、`X-Forwarded*`、`X-Real-IP`
- 不提供 wildcard CORS、auth cookie、URL token 或 readiness-token 例外

Wildcard listener 的選單 URL 轉為標準 numeric loopback，IPv6 會正確加括號。僅綁 LAN 位址時設定頁可使用該 listener，但日誌不開放；僅綁非標準 `127.x.x.x` 字面位址時，該 Host 也不在日誌白名單，選單明示不可用。程式不為日誌增開 listener、不自動改 HOST 或放寬門禁。IPv6 LAN scope 不會變成日誌例外。

若收到 403，請在閘道所在電腦使用符合設定的同源 loopback 頁面，檢查既有 HOST／PORT；不要以 proxy header 或跨源 Vite 開發頁繞過。這是**本機機器邊界，不是使用者登入驗證**；不能隔離同機其他程序／使用者，也不能辨識刻意移除轉送資訊的本機 relay。其他既有 API 的安全政策沒有因此全面修復。

### 畫面操作與缺口

- 等級為最低嚴重度：`debug < info < warn < error`；來源以 `startup/runtime/shutdown/http/standard/slog/application` 固定 ID 完全比對
- 搜尋為 safe message、code、允許的字串欄位之 Unicode case-fold literal substring，最多 256 字元，不支援 regex。只查本次程序目前保留的紀錄，沒有全部歷史檔查詢
- 暫停只關閉這個 view 的接收並保留 cursor；繼續從已處理 progress 接續。採集、其他訂閱與 logger 持續運作
- 清除只清瀏覽器列，不刪 ring、磁碟或改記錄等級；跟隨最新紀錄只控制捲動
- 篩選切換取消舊 snapshot／stream；離頁釋放訂閱。只有收到 handshake 才顯示已連線
- 斷線保留最後畫面及接收時間，單一擁有者使用附 jitter 的 1／2／4／8／16／30 秒退避，最長 30 秒；403 停止自動重試
- Retention、admission、replay 上限、subscriber overflow 會明示 gap；process instance 改變明示 reset。篩選排除的 sequence 不算資料遺失
- 檔案遺失與 capture 遺失分開呈現。磁碟恢復不會把之前缺漏補成完整持久歷史；畫面空白、載入失敗、拒絕存取、暫停、重連及截斷也各自有狀態

## 診斷內容與容量

新 sink 只接固定安全模板與允許的 scalar 欄位，先投影再存入 memory／file／web。必要 startup、runtime、HTTP、shutdown 事件有穩定 code；無法分類的 standard log／slog 等舊訊息顯示 `raw.suppressed`，不由任意錯誤字串猜安全性或嚴重度。HTTP 投影只含 method、route template、status、duration、opaque request ID 等安全欄位，不保存 query 值、headers、body 或 raw error。

這不是整個 stdout/stderr 的 tee。未經此入口的 `fmt`、OS/runtime crash 或第三方輸出不保證捕捉。舊 `/debug/logs`、packet history、audit 與 delivery truth 仍各有原用途；既有 CLI/raw console formatter 也不等於全面遮蔽，分享這類輸出前須另外檢查敏感資訊。Log API／heartbeat 不逐次產生普通 HTTP access event，避免自我放大。

`LOG_LEVEL` 控制 managed admission，預設 `info`；只接受 `debug/info/warn/error`。`LOG_OUTPUT` 只接受 `console/file`：desktop 保留 managed file＋ring；console/headless 預設 `console`，保留舊 console adapter 的 mirror；選 `file` 時關閉該 mirror 並使用安全檔案 sink。此選項不保證攔截所有直接 stdout/stderr 寫入，也沒有 `both` 值。

未明確設定 `LOG_FILE` 時，檔案位置是選定資料根目錄下 `logs/<opaque-db-id>/gateway-runtime.jsonl`。明確相對 `LOG_FILE` 以資料根目錄解析。Active、`.1`、`.2` 整組先做排他 ownership 及 DB／WAL／SHM／alias 保護；不接管既有非本 sink 管理檔案，不跟隨未驗證 reparse path 進行破壞操作。請勿手動刪除其 ownership header、共用不同 DB 的 rotation set，或把 LOG_FILE 指到資料庫。

下列是程式的容量契約，筆數或 encoded-byte 限制先到者生效，不是吞吐量保證：

| 邊界 | 上限 |
| --- | --- |
| 單一安全事件 | 8 KiB，UTF-8 安全截斷並標記 |
| Admission queue | 1,024 筆、2 MiB；滿時計 drop，不等磁碟／網路 |
| 目前程序 ring | 2,000 筆、4 MiB；淘汰最舊紀錄 |
| Snapshot | 預設 200、最多 500 筆 |
| SSE | 最多 8 clients；各 live queue 256 筆、512 KiB |
| 各 client replay 副本 | 500 筆、1 MiB；超過明示 replay gap，保留有限 tail |
| Stream 寫入／heartbeat | 每次寫入 5 秒 deadline；heartbeat 每 15 秒 |
| File queue | 1,024 筆、2 MiB，與 ring／subscriber 分開 |
| 診斷檔 | Active＋2 backups，每檔最多 5 MiB，合計 15 MiB |
| 瀏覽器保留列 | 1,000 筆、2 MiB，另行淘汰並顯示截斷 |

Slow client 會斷線，不拖住其他讀者；超過 client 數上限回安全的 unavailable/capacity 結果。磁碟滿或寫入停滯會標示 storage degraded、file drops／disk errors；ring 仍有獨立上限，不能把畫面可讀當作檔案成功落盤。重啟會建立新的 instance；API 不讀回舊 rotation 檔。

## 正常關閉與逾時

1. 從 tray 選「結束…」並確認，或使用該平台支援的終止訊號。取消不停止服務；並行來源共用同一個 coordinator
2. 停止新 HTTP admission、取消 log SSE，等候在途 HTTP；依序停止 runtime intake、group pipeline／delivery，再關閉 Share、connections、DB
3. 最後完成 diagnostics、移除 tray、釋放 owner。仍有 worker 使用的 DB 不會因 observer timeout 提早關閉

**15 秒是共用的等待提示期限，不是保證關閉時間，也不會到時自動強制成功退出。** 期限到達顯示 `stop-timeout` 與 pending phase，繼續觀察原協調器；不重新呼叫整組 Stop。Desktop 可在獨立的強制結束確認中選取消繼續等候，或明確確認後終止。強制結束可能留下未完成／unknown 工作，尚未持久保存資料可能遺失；不能當作正常 flush 完成。

Headless 不彈 native 對話框；保留非成功的 timeout 記錄，依實際 supervisor／OS 終止策略處理，不自行假稱 graceful。Windows 已確認的 session end 使用相同 coordinator，native callback 最多等候四秒；OS 可能更早終止，這是 best effort，不取代正常退出驗收。

停止或收到 `shutdown.finalizing` 不代表資料全部送達。已接受的 durable journal／outbox／receipt 與既有 recovery 規則要保留；unknown 不升格為 committed，也不能承諾復原從未 durable accepted 的舊記憶體資料。

## 驗收與回退

正式發行前，對同一 source／artifact 補齊 Windows 無 Go/Node 單 exe、無 console 閃窗、tray 鍵盤／overflow、Explorer restart、重複啟動、alias／session 隔離、資料歧義、只讀／磁碟故障、port 衝突、startup fatal、正常／逾時／強制退出、GUI-headless redirect/wait/exit 及 embedded logs 驗收。也要核對同一 artifact 的 durable recovery 與 Linux/macOS CLI 回歸。

回退由部署負責人控制：先正常停止並確認程序／worker／ownership 已結束，保存現有 DB、WAL/SHM 所需的一致性備份、outbox/receipt 與診斷檔；再使用已核對 schema 相容性的前版產物，指定同一預期資料來源。**不可用舊 DB 備份覆蓋新版已接受資料**，也不以刪除 backlog、ownership 記錄或強制清空診斷檔掩蓋問題。若狀態仍為 timeout／unknown，先保存證據並決定處理方式，不宣稱回退已安全完成。

主要實作對照：[launch/data](../../internal/apphost/)、[native desktop](../../internal/desktop/README.md)、[managed diagnostics](../../internal/diagnostics/)、[local guard](../../internal/api/runtime_log_guard.go)、[shutdown wiring](../../cmd/test_ui/gateway_shutdown.go)、[log UI](../../frontend/src/features/runtime-logs/)、[PowerShell build](../../scripts/build.ps1)、[Makefile](../../Makefile)。
