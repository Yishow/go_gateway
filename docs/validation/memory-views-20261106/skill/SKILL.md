---
name: validate-go-gateway-recording-e2e
description: 驗證 go_gateway 從實際 UI 新增設備、點位、目的地與記錄，到試寫、啟用、真實 SQL 讀回及重載的完整 E2E。驗證遇錯時用 hc 委派 terra_worker 修正，保留逐步證據與踩坑經驗。
---

# go_gateway 完整記錄 E2E

終點是同一工作區、來源及已套用記錄版本的真實 SQL 資料與 receipt，不能停在按鈕可按、設定保存、runtime running 或 buffer accepted。

## 確認本輪範圍

1. 讀 repo 的 `AGENTS.md`、`docs/testing.md` 與 Studio inventory 的 `START_HERE.md`、`context.json`、`CURRENT_STATE.md`。
2. 核對 branch、HEAD、dirty diff 與 remote 基準；保留既有 WIP，不自行切分支、提交或發布。
3. 以指定協議／目的地為準；未指定時先規劃 MC 3E、Modbus TCP 與可丟棄 SQLite 的完整旅程。其他資料庫各自列未驗，不能稱全 provider 通過。
4. 要求正式 `studio-v2-device-to-sql-acceptance` 時，另外遵守該 spec 的 SQLite／PostgreSQL、多設備／同群組多 entity、型別與三次 fresh run 要求；較小旅程不能代替正式驗收。
5. 沿用同一工作階段已同意的 hc、模擬器與隔離測試，只補問實質缺少的授權。正式 DB、LAN、真 PLC 與部署不由 skill 自動授權。

先讀 [完整旅程](references/journey.md)。啟動 Docker／HSL 時讀 [隔離執行環境](references/isolated-runtime.md)，排查時讀 [已知踩坑](references/lessons.md)，收尾前讀 [證據要求](references/evidence.md)。

## 準備與執行

- 用新 owned run 目錄、空設定 DB、空目的地 namespace 與 localhost 埠。接手舊 QA DB／已建設備只能當局部驗證。
- 建置候選前端並同步 embed，留 binary、SPA 資源與 dirty source 雜湊。`git archive HEAD` 不含未提交修正，必須明確覆蓋候選檔案並留下清單。
- 服務、長命令與 worker 依 `herdr-core`／`hc --skill` 執行。瀏覽器用獨立 agent-browser session 或既有 Playwright runner 的專用 Chromium，不用日常 Chrome。
- 初始化清單；先把以下變數設成本輪真實位置，不照抄私人歷史路徑：

```sh
python3 "$SKILL_DIR/scripts/evidence.py" init --repo "$REPO" --out "$RUN_DIR"
```

- 按完整旅程操作 UI；API GET／唯讀 SQL 核對真相，不用直接 API 新增設備／規則／記錄或 SQL DDL／INSERT 代替待驗 UI。
- 真實成功路徑不偽造回應。故障注入若丟棄真實請求的回覆，另列案例並保留真正後端 operation。
- 每段對上表單值、payload、ACK、canonical readback 與重載。畫面、log 與附件都先去除帳密、DSN、私密例外。

## 遇錯用 hc 指派 terra_worker

使用者指定：驗證中遇到產品或測試腳本錯誤，用 hc 委派 `terra_worker` 修正。

1. 保留失敗命令、環境、request／response、截圖與重現率；分清產品、腳本假設與環境問題，不降低斷言換綠燈。
2. 讀適用委派規則及 `hc --skill`，查 `hc templates`。工單列目標、原始證據、重現方法、允許修改檔案、禁止範圍、必要驗證與輸出位置。
3. Terra 固定 GPT-6.1 Terra、`max`。優先用已核對的模板；無模板時以 hc 啟動相同模型的互動 Codex，不用 `codex exec`：

```sh
hc run --bg --name recording-e2e-terra --cwd "$REPO" \
  codex --no-daemon --model gpt-6.1-terra \
  -c 'model_reasoning_effort="max"' --approve-for-me \
  '你是 terra_worker，先讀本輪指定工單，再於 ownership 內重現、修正並驗證。'
```

實際首則任務必須給完整工單或其絕對路徑；上面只示範啟動形式。核對目前 CLI 選項。模型／CLI 不可用就回報未派成功，不能把其他模型叫作 Terra。

4. 告知 worker 共享工作樹、不可回退他人變更、不可 stage／commit／push／動正式資料、不可再開子代理。跨模組依賴先由主代理判定是否仍屬原需求。
5. 產品先補真實失敗測試，再依 repo TDD／Go 指引修正。腳本錯誤則保留誤判原因，不把 harness failure 當產品 RED。
6. 主代理以 `hc wait`／`hc read` 收結果，親自查 source／diff／輸出，重跑原失敗案例與受影響完整旅程。worker 說完成不等於驗收。
7. 候選來源改變時保留舊 run，重新建置、開新 run 並串接 `previous_run`；不得混用舊版成功畫面。
8. 同因連續失敗兩次先換調查方法；同一假設三次修正仍失敗就停止該假設。缺外部授權只阻擋相依案例，其他驗證繼續。

## 收尾

- 停止本輪記錄及自己擁有的服務／瀏覽器／worker。核對 PID／namespace，只清理本輪資源，保留證據。
- 更新 `checks.json`：未執行 `NOT_RUN`、缺環境 `BLOCKED`、真失敗 `FAIL`，不可假 PASS。
- 檢查清單：

```sh
python3 "$SKILL_DIR/scripts/evidence.py" check --run "$RUN_DIR"
```

工具只檢查清單、來源一致性與附件存在，不能取代主代理核對 UI、API、協議與 SQL。

- 按證據要求交付範圍內結論，列明 SQL、視覺、平台與現場限制。
- 新踩坑追加至 lessons，附日期、版本、症狀、已證根因、最小重現與證據定位；不抄整份 log、secrets 或推測。
