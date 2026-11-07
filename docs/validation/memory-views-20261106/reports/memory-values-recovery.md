> 歷史報告閱覽副本：正規化本機絕對路徑與附件連結；保留當時結論與未驗項目。原始位元組及完整附件見 `local-only/original-evidence.tar.gz`，不要將歷史 PASS 當成新版重驗。

# MC／Modbus 記憶體視圖與恢復驗證

## 結論與範圍

- **PASS**：同一候選、兩個已存在來源，在桌面預設主題的五種視圖中核對 int16 地址、已用／相鄰空格、正數／零／負數、來源身分及最小互動；恢復前後各 10／10 視圖通過。
- **PASS**：授權後僅暫停 owned HSL `t14v`／PID 67026；兩協議均由 Gateway 點位 GET、runtime SSE 與 UI 格子觀察到讀取錯誤。恢復後種入不同非零值，種值 ACK、模擬器獨立 TCP 協議讀取、Gateway GET 和 UI 均對上新值。
- **NOT_RUN**：單獨觸發 stale 狀態；實際觀察的是 error／SSE `quality=bad`，不能稱 stale 已驗。原圖對照、手機／暗色、多資料型別矩陣、真 PLC／現場及完整 UI→SQL 旅程不在本工作者本輪範圍。

候選為 HEAD `7a01e7d5e0299e8709063cf01e0dc8e21bda097c` 加父代理的 11 個 WIP 檔；binary SHA-256 `432307aba7a9347047bc0b4bc63e88c06a00fabf55cefe9c69454fd0ee9e7a1b`。父代理已核對候選 binary、SPA 和 build manifest，引用 `../parent-checks/build-provenance.json`；本工作者未重建或重算該建置。Gateway 為 localhost `127.0.0.1:3230`、`t16d`。此處是模擬器及專用 Playwright Chromium（1440×900、en-US、預設主題），不代表真設備或正式驗收。

## 五種視圖

Workspace `7bb2a3c5-7a0f-45c1-927a-3966c6546624`。既有規則均為 enabled、int16、3 點；來源建立／編輯／啟用由父代理另行確認，本工作者沒有重做或修改。

| 協議／設備及規則 | 格子專注 | 範圍導航 | 水平地圖 | 範圍章節 | 格子＋詳情 |
| --- | --- | --- | --- | --- | --- |
| MC `dev-32b42b30-fca5-48b3-aec8-0f9077582ff1`；`rule-a9700f69-f1cf-4f9a-923e-31f872882d9b` | PASS | PASS | PASS | PASS | PASS |
| Modbus `dev-3921b9cb-cc16-465c-8132-064f750d3658`；`rule-c2fac41b-033b-4c15-a5f3-2cc97f86a208` | PASS | PASS | PASS | PASS | PASS |

恢復前的 UI／點位 GET：MC D0／D1／D2 = `123 / 0 / -42`；Modbus 40001／40002／40003 = `456 / 0 / -17`。每種視圖保留原始畫面、同次點位 GET、讀值及 accessible label，見 `run-S8JJ2T/run.json` 與同目錄的十組 `*.png`／`*-get.json`。範圍導航實際跳到規則、章節展開、水平地圖移至第三格；詳情視圖逐格核對三個值。水平地圖及詳情另有 `*-selected-detail.png`。

恢復後的同一流程另為 10／10 PASS，MC `321 / 0 / -42`、Modbus `654 / 0 / -17`；見 `run-Kq39nT/run.json` 與同目錄附件。該重跑的預期值取自恢復後模擬器 register API 附件，並有另行直接 TCP 讀取結果可交叉核對。第一輪探索 `run-kFDPaX/` 保留；因整頁截圖無法完整呈現水平地圖詳情，補拍後才產生 `run-S8JJ2T/`。

## 暫停與恢復

主代理明示 PostgreSQL 三桶和 stop settled 已確認後，才開始此階段；這些 PostgreSQL 結果是父代理提供的脈絡，非本工作者自行驗證。操作前 `hc get t14v` 和四個 loopback listener 均指向 PID 67026，兩台 Gateway 設備 GET 顯示正在採集。

| 時間（2026-11-06，台北） | 實際觀察 |
| --- | --- |
| 14:46:45.372 | 以 SIGSTOP 暫停 HSL PID 67026。 |
| 約 14:46:51–52 | MC D0 與 Modbus 40001 的點位 GET 均有 `last_error`；SSE 各有 `quality=bad`、`stale=false`；UI 格子顯示 `!`，accessible label 表示讀取失敗。 |
| 14:46:52.572 | 以 SIGCONT 恢復同一 PID；自有 `seed-new.py` 沿用既有 `terra/seed_hsl.py` 的請求格式，僅向 owned simulator 寫入 MC D0=`321`、Modbus 40001=`654`，兩次 ACK 的三個 simulator 結果均成功。 |
| 約 14:46:53 | Gateway GET 清除錯誤並讀到 `321`／`654`；獨立瀏覽器 UI 也顯示這兩個新值。 |
| 14:48:26 | `wire-read.py` 以獨立 TCP socket 向 MC 3E 與 Modbus TCP 發出讀取並保存原始請求／回應，分別讀到 `321 / 0 / -42` 與 `654 / 0 / -17`。 |

完整時間、點位 ID、GET snapshot、SSE 事件、UI 狀態與操作結果見 `recovery-wpVEgh/recovery.json`；種值 ACK 見 `recovery-wpVEgh/seed-new.json`；實際 TCP 封包與解析值見 `recovery-wpVEgh/independent-wire.json`。故障與恢復畫面為同目錄 `mc_3e-paused.png`、`modbus_tcp-paused.png`、`mc_3e-recovered.png`、`modbus_tcp-recovered.png`。

故障時 UI 的個別格子顯示錯誤，頁面中的串流狀態仍顯示 `Live`；此狀態不能單獨當成設備讀值成功證據。恢復驗證使用 grid 視圖，不宣稱各布局各別經斷線測試。

模擬器 `/api/simulation/registers` 的舊讀回檔及本輪 `*-wire-before/after.json` 是 simulator API 證據。本機 HSL source 顯示該 API 由 `server.Read` 取得資料；未獨立確認執行檔與該 source 一致，因此不將這些 API 檔案單獨當成獨立 TCP wire 證據。上述 `independent-wire.json` 才是本工作者實際取得的網路協議請求／回應；恢復前沒有本工作者的獨立 TCP 封包。

## 工作檔、限制與收尾

- 本工作者僅寫入 `memory-views/` 下的 `PROGRESS.md`、本報告、`inspect.cjs`、`run.cjs`、保留的 runner 前版、`recovery.cjs`／`recovery.sh`、`seed-new.py`、`wire-read.py` 和證據目錄。未修改產品或來源、未 stage／commit；沒有 mock API。
- 五視圖和恢復腳本對瀏覽器非 GET／HEAD／OPTIONS 設防護；結果中沒有任何被攔下的非讀取請求，也沒有 page error。Simulator 的兩次 POST 是使用者明示授權的新值種入。
- 自有 Chromium 均由 runner 關閉。最終 `hc get t14v` 顯示同 PID live、`lsof` 顯示四埠仍監聽，且恢復後 simulator／TCP 讀取已成功；沒有停止 Gateway、PostgreSQL 或其他服務。HSL 保留本輪新值。
- Harness／環境觀察：sandbox 內 `ps` 取得 PID 時為 `Operation not permitted`，改用 `hc get` 加 `lsof`；一次臨時群組 GET 命令誤把 `data` 當陣列，確認實為物件且群組在 `data.groups`；一次 zsh 路徑 glob 無相符項而終止，改查明確路徑；曾猜測的 agent-browser config 路徑不存在，實際使用既有 Playwright 專用 Chromium。探索時來源頁首張截圖仍在載入，正式 runner 已等待 `Memory view` 才進行。以上未改產品、均非產品失敗。
- 最終附件核對的第一版臨時 Python 命令誤把 `restored` 的直接 snapshot 當作含 `sample` 的物件而出現 `KeyError`；依實際 JSON 形狀更正後重跑 exit 0，核對兩輪各十個畫面與 GET 附件、兩協議故障 GET／SSE／UI、新值 ACK、恢復 GET／UI 及 TCP 封包。
- `run.cjs` 可用 `--recovery-dir=$TMP/go-gateway-full-recording-e2e-20261106/memory-views/recovery-wpVEgh` 在目前新值下重跑；`run-before-recovery.cjs` 保留原值版本。重跑前須確認目標設備仍在採集，且 simulator 未被其他任務佔用。
