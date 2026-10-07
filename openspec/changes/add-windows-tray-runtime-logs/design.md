## Context

本案是規格提案，基準 main 為 `ffdaa002ad0d02961bc239c8b0772da7b83be848`。現在的入口是 `cmd/test_ui`，已有 embed SPA、環境變數／dotenv 設定、SQLite 設定庫、runtime、Modbus Share 與 durable group pipeline；詳見 [evidence.md](evidence.md)。不採用舊 OpenSpec context 中的 Viper／zap 敘述作為 production 現況。

`scripts/build.ps1` 明確移除 GUI subsystem 旗標，以 console 的 X 作為操作入口。`server_runtime.go` 在 bind 成功前印出啟動訊息，透過 `cmd /c start` 開 Windows 瀏覽器，HTTP 啟動錯誤用 `log.Fatalf`；這些都必須接到可觀察的 lifecycle，不能只加 linker flag。

目前沒有使用者登入 middleware，HOST 預設 all interfaces、CORS 預設 wildcard。workspace/revision/readiness 保護不是身分認證。新日誌只對本機開放，且不能宣稱既有 API 已因本案獲得完整安全防護。

## Goals / Non-Goals

**Goals:**

- Windows 桌面單檔啟動時只有系統列圖示，網頁與本機錯誤對話框提供操作及診斷
- 同一資料來源只允許一個 Windows gateway owner，啟動與關閉失敗皆有明確結果
- 在既有 SPA 提供安全、有限容量且能反映斷線／資料缺口的應用程序日誌
- CLI／headless 不依賴桌面；保留 Linux/macOS 原 CLI、輸出與資料庫語意

**Non-Goals:**

- macOS/Linux tray、Windows Service／登入自啟安裝、restart 按鈕、更新器、遠端啟停
- 完整 auth 平台、跨使用者本機隔離、遠端／reverse-proxy 日誌服務
- ELK／第三方收集平台、log DB、查詢全部歷史檔案、server-side 清除或下載日誌
- 修改 outbox／receipt／unknown／retry／ownership 語意，承諾所有資料已送完，或修復所有既有 raw debug logging
- 此次寫入產品程式、執行 deployment 或宣稱 Windows 現場通過

## Decisions

### 1. 同一核心，分開桌面與 console 建置預設

- Windows 桌面發行目標使用 `-H=windowsgui`，內嵌 icon 與完整 SPA；雙擊預設 desktop。任何build明確傳入 `--headless` 都停用tray、native dialog及自動開瀏覽器，並覆蓋 `AUTO_OPEN_BROWSER=true/1`；可供supervisor使用，但不是 Windows SCM service
- 現有開發console build保留，未帶模式旗標時稱legacy console default：不建立tray，但仍遵守原 `AUTO_OPEN_BROWSER=true/1` 的選擇；它不等同明確 `--headless` override。不因從終端執行就猜測切換模式。Windows release與CLI建置命令／產物名稱清楚區分，不同模式不要求額外companion process
- `make build`／`scripts/build.ps1` 的 Windows desktop 選項須一致；`start.ps1` 開發及 `start.sh` 原 CLI 路徑維持 console。只有 Windows 目標加入 GUI flags，不因 `.exe` 副檔名推論平台
- `--version`／`--help` 不啟動採集；console build能列印，desktop模式用native version/help呈現；explicit headless以呼叫端提供的標準handles輸出且保留exit code。GUI-subsystem exe的redirect／wait／exit status必須在Windows cmd、PowerShell及supervisor harness分別實測，必要時接正確inherited handles；不能以link flag推定可擷取，也不為缺handle配置console
- 以 Windows build tags 隔離 native adapter，優先使用現有 `golang.org/x/sys` 的 Win32 能力；維持純 Go 主體與平台無關 lifecycle 介面。若選 tray library，先驗授權、維護與 build/link 相容性，不直接改成新 UI framework

替代方案：只 hide console 會失去入口及早期錯誤；WebView/Electron 會增加部署依賴；讓所有 Windows build 自動變 tray 會破壞開發與測試，因此不採用。

### 2. 資料位置辨識先於 migration 與任何採集

- CLI 繼續使用既有 process env > cwd `.env` > default，以及 `GATEWAY_DB_PATH` > `DB_PATH` > `SQLITE_PATH`；不移動資料
- desktop 的預設根目錄是 canonical executable directory，`.env` 及未提供絕對路徑的設定以該根目錄解析；process env 優先。DB 的三個既有環境變數及其順序不變
- 選定來源要顯示給本機操作者；啟動 CWD 與 exe 目錄不同且存在另一份 `.env`／`datalink.db`，或 legacy 來源無法辨識時，不猜選、不 migration、不默默產生新庫。native 說明要求以既有設定提供明確絕對 DB 路徑再啟動
- 無明確 DB 設定且目標庫不存在時，desktop 先確認「在這個位置建立新資料庫」；取消不建立。這不是自動匯入／搬移工具，也不能承諾發現所有歷史安裝位置
- 明確配置的 DB 所在目錄為資料根目錄；canonicalization 要涵蓋 Windows case、相對路徑、symlink/junction 及現有檔案身分。無法取得可信身分時安全失敗
- 預設診斷檔在選定資料根目錄的 `logs/<opaque-db-id>/gateway-runtime.jsonl`，opaque-db-id來自canonical DB身分，同目錄不同DB不共用rotation。明確 `LOG_FILE` 使用既有config loading優先序；desktop相對路徑以選定資料根目錄解析。初始無寫入權限時顯示native錯誤，不提權或改用另一資料庫
- 開diagnostic檔之前先保護其整組active/backup paths：canonical/實際handle身分不得等於或alias到DB、WAL/SHM或其他受保護資料，並取得整組排他log ownership；另一instance已使用相同或重疊rotation檔名即fail closed。不follow未驗symlink/reparse path做truncate/rename/delete，防止檢查後換target；既有非本sink管理檔案不得被接管輪替

替代方案：一律以 exe 目錄建立新庫容易產生看似資料遺失；直接沿用未知 shortcut CWD 會使雙擊行為漂移，因此採辨識及明確選擇。

### 3. 單實例、啟動與 URL

- Windows desktop/headless 共享同 DB 身分的 OS-backed exclusive owner lock，取得鎖後才可開庫/migrate、開設備、bind listener。程序死亡由 OS 釋放；PID metadata 僅作診斷，不能只以 PID/file 存在就判斷有效或殺程序
- 第二個 desktop 僅能向已驗證的同使用者／同 session owner 發送固定「開設定頁」命令，使用有 user ACL 的本機 IPC，不接受任意 URL／command。驗證失敗或 owner 為 headless／其他session則顯示「已在執行」並退出，不建立第二組 worker
- 不同資料庫仍須各自成功 bind 明確配置的 port；port 占用不能當作本產品實例證據，不 kill、換 port、覆寫 owner 或打開未知程序網頁
- 啟動次序：最小 native 錯誤回報 → 路徑選擇與 owner lock → 安全 log broker／diagnostic file → tray 初始化 → bind HTTP listener（尚未 Serve）→ 初始化 DB/services → 開始 Serve → 發出本實例 readiness → 開網頁
- 初始 tray 建立失敗時不留下無控制入口的背景程序，釋放已取得資源並非零退出；headless 不會呼叫 tray
- listener 成功且嵌入資產可服務後才可標示 web ready。runtime 啟動失败若依現有政策可繼續，狀態為 degraded，不能顯示採集正常
- Browser open 使用 native shell API，固定產品路徑，依實際 bind address 產生有效 URL。wildcard bind 轉成 numeric loopback；IPv6正確加中括號。僅 bind LAN 位址時設定頁可用該明確位址，local-only logs 顯示不可用，不新增 listener 或放寬權限
- desktop成功啟動預設開 `/studio/v2`；`AUTO_OPEN_BROWSER=false/0`停用。Legacy console default只有原 `AUTO_OPEN_BROWSER=true/1` 才開啟。任何build的explicit `--headless` 一律不自動開啟，即使環境值true/1；選單開啟始終由使用者點選

### 4. Tray 為程序控制，不複製 domain 設定

右鍵與鍵盤選單：

| 項目 | 行為 |
| --- | --- |
| 開啟設定頁 | 固定 `/studio/v2`，成功bind前停用 |
| 查看執行日誌 | 固定 `/studio/logs`；本機存取不可用時顯示原因 |
| 程序／採集狀態 | 唯讀，分列 starting、serving、degraded、stopping、stop-timeout；採集狀態讀現有 runtime authority |
| 版本資訊 | product version、commit、build mode；缺少metadata顯示unknown/dev，不編造release |
| 結束… | native 確認採集及服務將停止，cancel為預設安全選擇 |

左鍵雙擊等同開設定頁。没有常駐主視窗可「最小化」；關掉／縮小browser不停止程序。Windows可能把icon置於overflow，不能保證釘選位置；Explorer/taskbar重建時重新註冊icon且不重啟runtime。狀態同時有文字，不只靠顏色。

### 5. 關閉是單一協調器，不是 `os.Exit` 捷徑

- Tray確認、SIGINT/SIGTERM及支援的OS session shutdown走同一個idempotent coordinator；OS shutdown無法保證給足時間，不宣稱等同完整退出驗收
- 進入stopping後停用再次退出操作，拒絕新的mutating HTTP工作並結束log SSE，維持native狀態可見；協調HTTP在途請求、停止runtime/scheduler intake，再停止group pipeline/delivery worker，關Share、connections、DB，最後flush diagnostics、移除tray與釋放owner lock
- 使用既有context/cancellation seams，保留已接受的durable journal/outbox/receipt與下次recover策略；不清backlog、不把unknown升格committed，也不承諾legacy in-memory資料可恢復
- 15秒是本案設計的「關閉等待提示期限」，不是已測效能或強制成功期限。各phase共用剩餘deadline，不疊加多個完整budget；需要修正忽略context的等待，不能在worker仍使用DB時先close DB
- Shutdown deadline只限制呼叫端等待，不得中止已接受資料的最後寫入：time-series batch writer最後flush使用自己的WriteTimeout，被Close取消的timer批次回到buffer由最後flush寫入；`http.Server.Shutdown`不關hijacked連線，WebSocket handler須在server base context取消時自行關閉連線，HTTP phase才可能結束
- 期限到時記錄phase及safe code，狀態為stop-timeout，native顯示「仍在關閉，尚未確認完成」，繼續觀察原協調器，不再啟第二次Stop。可等待或明確二次確認強制結束；只有這個確認允許process termination，並告知可能有未完成／unknown資料
- Headless在deadline後記錄非成功關閉狀態，遵循supervisor/OS終止策略，不自動假稱成功；不新增遠端force-exit API
- 初始fatal與可控panic經同一安全diagnostic sink及native錯誤路徑；native訊息有穩定錯誤碼、下一步及已成功寫入的診斷檔位置。檔案寫入也失敗時明說無保存紀錄。不能承諾攔截每個OS/runtime/third-party crash

### 6. 安全事件入口與捕捉範圍

- 新模組提供 typed event schema：`instance_id`、十進位字串`sequence`、UTC timestamp、level、source、stable code、safe message、allowlisted scalar fields、truncated；instance每次程序新建
- 接入標準log、slog、Gin safe access/recovery、HTTP server errors及startup/shutdown。逐項盤點production reachable producers，將必要startup/runtime錯誤改用safe code；無法可信分類的legacy/raw事件只輸出固定省略訊息與correlation，不把任意error string當safe
- 不tee整個stdout/stderr，不複製full Config、DSN、headers/cookies、request/response body、raw SQL/args、packet bytes、endpoint或私人路徑到新sink。診斷檔路徑只在本機native視窗呈現，不進一般operator DOM
- HTTP safe投影僅method、route template（無query values）、status、duration、opaque request ID；log API、SSE heartbeat不再產生逐筆access logs，避免自我放大。`http.access`只記mutation（非GET/HEAD/OPTIONS）、status>=400及耗時>=1秒的請求；快速成功的讀取是Studio輪詢，逐筆記錄會在數分鐘內把2000筆ring的startup/runtime事件擠出
- 原console formatter與既有query測試保持相容，新安全投影使用独立測試；不宣稱舊console/debug歷史已全面安全。所有新sink在儲存前先做模板／欄位白名單，regex僅防禦補強
- Desktop必有managed memory/file診斷，不建立console；headless保留console/redirect輸出，同時可供本機web讀安全事件。`LOG_LEVEL`控制managed事件入列等級（debug/info/warn/error），預設info；未知legacy訊息不從字串猜severity
- `LOG_OUTPUT`目前未接production sink；本案只採 `.env.sample` 已列的console/file兩值，不新增both。Desktop固定保留managed file/ring；console只允許mirror已存在的標準handle，不建立視窗；file不mirror。Headless預設console保留原console行為，file可選只有safe events的managed file，不能把legacy raw console轉成新raw file sink；ring皆保留。無效值顯示safe配置錯誤，不宣稱舊版本已有此sink實作

### 7. 明確容量與診斷檔政策

以下全部是規格設計上限，不是已測效能／容量；實作需以burst、stall、disk fault與記憶體量測驗證，未通過不可放寬成無限：

| 邊界 | 上限／行為 |
| --- | --- |
| 單一安全事件 | 編碼後8KiB；安全UTF-8截斷及truncated，不保留raw超長副本 |
| admission queue | 1,024筆且2MiB，先到任一上限即計drop，不阻塞採集 |
| 記憶體ring | 2,000筆且4MiB，evict oldest；metadata／索引也須有界 |
| snapshot | default200、max500筆；參數越界400；搜尋最長256字元且非regex |
| SSE | 最多8個client；每client最多256筆且512KiB；慢client斷開、不拖累其他client |
| replay副本 | 每client最多500筆且1MiB，獨立於live queue；超過則明示replay-limit gap並只送可容納tail |
| network／UI | 每次stream write deadline5秒、heartbeat15秒；browser最多1,000筆且2MiB |
| file | active＋2 backups，每檔最多5MiB，共15MiB；只輪替本sink擁有的檔案 |

File與network writer只消費bounded queue，不讓磁碟/網路IO卡住runtime producer。每個sink的drop、retention eviction與disk_error各自可見，不當作應用成功事件。寫入後的disk full不停止工業採集；memory仍有界並顯示diagnostic degraded，恢復成功後明示間隙。啟動前無任何可用持久診斷則desktop fail closed並native告知。

檔案僅為短期本機diagnostics；不是audit、telemetry或delivery durability，不能刪除SQLite、outbox、receipt或其他logger檔案。API只讀current-process ring，不接受檔案路徑、也不讀回舊輪替檔；重新啟動會明確重置web歷史。

### 8. Log API獨立本機門禁

- GET `${API_BASE_PATH}/system/logs` 與 `/system/logs/stream` 共用guard；在任何ring讀取／subscribe前檢查實際socket `Request.RemoteAddr` 的numeric loopback，忽略Gin ClientIP/proxy信任設定
- 要求Host為實際port的`127.0.0.1`、`[::1]`或`localhost`；Origin存在必須與該請求scheme/host/port完全同源。`Sec-Fetch-Site`只接受單一`same-origin`或`none`；缺少可繼續其餘guard，`same-site`、`cross-site`、重複及格式錯誤值一律safe403。無Origin的正常同源GET/SSE與本機非browser client可通過其餘guard；不把same-site當same-origin
- 出現Forwarded、X-Forwarded-*或X-Real-IP即拒；不允許proxy relay。此法不能辨識刻意刪除headers的本機proxy，因此是本機trust boundary，非使用者authentication，不隔離同機process/users
- 覆蓋全域CORS輸出，不回wildcard／跨源credentials；回`Cache-Control: no-store`、`Cross-Origin-Resource-Policy: same-origin`與`X-Content-Type-Options: nosniff`。拒絕固定safe403，不帶logs或raw原因
- 不新增auth cookie/token，不把readiness token或URL query當身分證明；Vite跨源與remote logs不加例外。其他routes/HOST/CORS政策不因本案被全域改寫

### 9. Snapshot／SSE與畫面一致性

- Snapshot參數為level/source/q/limit/before；只搜尋current retained safe records，返回chronological records、instance、oldest/latest cursor、drop/eviction counters與sink health。level是最低severity（debug<info<warn<error），source是穩定source ID完全相等，q是safe message/code/allowlisted string fields上的Unicode case-folded literal substring，不用regex。UI註明搜尋範圍是本次程序的有限紀錄
- Stream沿用完全相同的level/source/q條件，以`Last-Event-ID`或首次`after`接續，兩者相衝突回400；ID為instance+sequence，sequence字串避免JS精度問題。Snapshot在短鎖內取得一致watermark H；訂閱時在另一個短鎖內取cutoff K、bounded replay副本(H,K]並註冊只收>K的live queue，隨即解鎖。不能跨HTTP請求或在serialization/network write期間持鎖；replay獨立於較小的live queue，先送replay再送live，超過replay上限明示gap/tail
- 事件ID供dedup；被filter排除的sequence不是遺失。Handshake/heartbeat帶已處理progress watermark，只有該watermark前所有matching records已送出或明示gap後才可前進；沒有matching事件仍能推進cursor，不能跳過未送出的matching queue。Browser保留progress cursor與last matching record為不同欄位
- 太舊cursor、重啟instance或overflow明示typed gap/reset及可取得tail；slow client已無法寫入時直接斷線，將loss metadata保留在bounded counters／下次handshake，不為送gap無限等待。Subscriber overflow屬於該stream：overflow的stream在既有5秒write deadline內best-effort送typed `subscriber_overflow` gap後斷線；共用`subscriber_overflow`計數只作metadata觀測，不放入共用gaps，browser也不以計數增加推論自身缺漏，避免他人overflow讓新開的view永遠顯示gap。未知future/malformed cursor400。filter切換建立新snapshot/stream並取消舊request/subscription，遲到回覆不可覆蓋新filter
- 使用單一reconnect擁有者；backoff為1/2/4/8/16/30秒封頂並jitter，穩定handshake後歸零。403停止自動重試；斷線保持最後畫面及時間，沒有handshake不顯示connected
- `/studio/logs`為lazy SPA route，Studio主頁與runtime各有入口；view可篩選level/source、safe文字搜尋、pause/resume、tail及clear view。pause只停止該view接收並保留cursor，logging/採集继续；resume若evicted顯示gap。clear只清browser，不刪server檔案或ring
- 所有內容render為text，不執行HTML。empty/loading/connected/reconnecting/paused/access-denied/error/gap/disk-degraded都有en／zh-TW文案與鍵盤可達操作；沒有接資料不顯示示範log

## Risks / Trade-offs

- [沒有全域auth] → 新日誌限實際loopback且同源；其餘現有API風險不在此案宣稱修復
- [隱藏console會藏fatal] → logger、native錯誤與tray先就緒；ready依listener/資產/服務真實狀態
- [資料根目錄切換造成新空庫] → legacy ambiguity fail closed、明確既有路徑優先、新庫先確認，禁止自動搬資料
- [OS shell重建／無互動session] → tray重新註冊、headless顯式可選；initial tray fail不留下背景instance
- [日誌被drop或rotation] → bounded counters及gap；明說不是完整持久歷史或資料寫入證明
- [現有Stop內有無期限等待] → cancellation/phase測試先RED，再收斂；timeout保留stopping與operator選擇，不先關DB冒成功
- [改GUI subsystem破壞CLI] → release與console build分流、Windows headless及Linux/macOS回歸；Windows實機驗證不可由cross-build取代

## Migration Plan

1. Expand：加入平台無關lifecycle/safe log介面与tests，保留console可建置；再接入local log API與SPA
2. Migrate：逐一將必要production producers接入safe投影；Windows adapter完成後才切desktop release預設。每批皆需現有回歸與build通過
3. Contract：僅移除desktop路徑的console依賴、`cmd /c start`與fatal捷徑；不移除CLI輸出或既有debug API
4. 發行需記錄source SHA、Windows PE subsystem、完整embed graph、manual tray/shell/failure/shutdown證據。舊版回退先正常停程序，保留現有DB/outbox/diagnostic資料，不restore舊庫覆蓋新accepted資料

## Open Questions

無阻擋規格起草的未決產品選擇。Restart與開機自啟本次排除：前者涉及安全重啟／資料狀態，後者涉及持久系統設定；如日後需要應另立明確需求。Windows API/library的實作選擇需在編碼前以build spike與license核對驗證，不能用此未驗選擇聲稱支援平台已完成。
