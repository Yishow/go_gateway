# 證據與交付

每個新 run 留 `manifest.json`、`checks.json`、`REPORT.md`，以及實際需要的 screenshots、requests、唯讀 SQL、timing、repairs 附件；不建無用途空資料夾。

## 清單工具

`scripts/evidence.py` 讀 repo Git／檔案並寫自有 run 目錄；一般檔案的符號連結會讀取目的檔。不啟動服務、連設備、操作資料庫或修改 Git index。

```sh
python3 "$SKILL_DIR/scripts/evidence.py" init \
  --repo "$REPO" --out "$RUN_DIR" \
  --protocol mc_3e --protocol modbus_tcp --database sqlite
python3 "$SKILL_DIR/scripts/evidence.py" check --run "$RUN_DIR"
```

預設兩協議＋SQLite。失敗後不能刪原定 provider 縮小範圍拿 PASS。既存 run 拒絕覆蓋。

每段初始 `NOT_RUN`，完成後填：

```json
{
  "id": "recording/mc_3e/sqlite/sql-readback",
  "status": "PASS",
  "detail": "三個完整 bucket 的值、applied revision、effect 與 receipt 相符",
  "evidence": ["sql/mc-sqlite-readback.json"]
}
```

工具拒絕缺段、重複／多餘 ID、非 PASS、空附件、越界路徑與來源漂移。附件須在本 run 內。它不能判斷附件真假、SQL 正確性或 UI 是否實做，仍由主代理審核。

來源快照包括 HEAD、Git diff、追蹤檔案實際內容與執行權限、未忽略的未追蹤檔案內容與執行權限，以及符號連結的目標路徑；因此 Git `assume-unchanged` 也不能掩蓋追蹤檔案內容的變化。每次 `init`／`check` 都需讀取追蹤檔案；讀取不是原子操作，建置與採證時須固定候選來源。

快照不遞迴檢查目錄符號連結的內容，也不涵蓋被 Git 忽略的檔案、submodule 工作樹的詳細內容及外部依賴；它也無法證明實際執行的 binary 或 SPA 產物與來源相符。若這些是候選建置輸入或產物，須另外保留雜湊並由主代理核對。

舊版工具建立的 run 缺少新增的來源快照欄位；新版 `check` 會拒絕直接沿用。保留舊 run，以新版重新 `init` 和採證；只有候選來源確實改變時才需因此重新建置，並在報告串接 `previous_run`。

引用舊證據需列版本與用途，不能把舊 screenshot 複製成新版 PASS。修復改了候選就留舊 run，重新建置並開新 run，REPORT 串接 `previous_run`。

## 身份與狀態

- 串 UI → workspace／device／rule revision／point／tag／mapping → connector revision／group revision／applied revision。
- 正常交付再串 record ID、effect key、payload digest、receipt、UTC bucket/acquisition 時間與 provenance。
- 保存、diagnostics、running、buffer accepted、SQL committed、readback verified、cleanup completed 分別留證據。
- mocked API、真 simulator、真 PLC 及各 provider／平台分開判定。

## hc／Terra 工單與回驗

保留原失敗命令／候選／payload／畫面、owner 檔案與禁止範圍、必要測試、hc terminal、實際模型與 effort、任務送達狀態、worker diff／測試及主代理原案例重驗。

worker 不自行擴 scope；超出需求、真實資料及部署授權由主代理向使用者處理。不得只轉述 worker 摘要當 PASS。

## 最後報告

先說宣告範圍完成度，再列：

- HEAD／dirty／binary／SPA 身分與環境；有修復則連舊／新 run。
- 協議 × 目的地 × 階段矩陣，未執行原因。
- SQL 同 revision／receipt／provenance 證據。
- 原圖、響應式、鍵盤及資料值各自結論。
- 修正、檔案、Terra 委派與主代理驗證；範圍外問題只回報。
- 自有服務／資料 cleanup、證據位置與平台限制，不擅自 commit／push／部署。

`check` exit 0 只能稱「證據清單檢查通過」。未親自跑 UI→SQL、只有 mocked E2E 或 skip 的這輪，產品完整驗收仍是未通過／未執行。
