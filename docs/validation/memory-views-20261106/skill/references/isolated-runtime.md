# localhost Docker／HSL 執行方式

使用目前使用者授權；這些範例不授權額外主機、正式資料或私人憑證。所有路徑、名稱、埠先改成本 run 的值並核對空閒。長命令由 hc 管理。

## HSL Mac simulator

已知入口是使用者提供的 `hsl-mac-simulator` 專案。先讀該 repo 指引與當下 README，核對架構相符的已建 binary；不要自行換來源不明的執行檔。

這次已使用過的 macOS arm64 相對 binary 路徑為 `dist/osx-arm64/HslMacSimulator`，未來須重新核對。設定 `HSL_BINARY` 為授權的實際檔案。

隔離試用模式在 owned cwd 放空白 `.env`，避免誤載原專案的私人授權／設定。以下只示範同時開 MC 與 Modbus；用授權 license 的情況按 owner 既定程序，不把 license 印到 log：

```sh
hc run --bg --name "$HSL_TERMINAL" --cwd "$RUN_DIR/hsl" \
  env -u HSL_LICENSE_CODE -u HSL_LICENSE_CERT_PATH \
  HSL_BIND_ADDRESS=127.0.0.1 \
  HSL_MC_ENABLED=true HSL_MC_PORT=15100 HSL_MC_BINARY=true \
  HSL_MODBUS_ENABLED=true HSL_MODBUS_PORT=1503 \
  HSL_FATEK_ENABLED=false \
  HSL_WEB_ENABLED=true HSL_WEB_PORT=18086 \
  HSL_WEB_TAILSCALE_IP= HSL_LOG_TRAFFIC=false \
  "$HSL_BINARY" --headless
```

上面環境參數以當下 simulator source 為準。先檢查繼承環境是否有其他 listener 設定；必要時用乾淨環境加明確允許項，不能只靠設定兩個埠就推論全為 loopback。

用 hc 輸出及實際 listeners 核對 host、埠、PID；保留試用限制說明。種值可透過 simulator 自有 API／畫面，這不等於用 API 建產品設備。

已驗證過的種值觀念：MC word bytes 與 simulator 的 ABCD／BADC 名稱不一定直覺對應；先用單一已知非零值與獨立封包讀回校準，再驗 float／多 word 型別。不同設備給可區分值，別只靠全零成功。

2026-11-06 的 `simulation/write` source 在任一 mc／fatek／modbus server 未運行時會拒絕整批種值。本輪額外在 owned loopback 埠啟動 FATEK 後，三者的種值 response 都成功；FATEK 只供 fixture，不是產品來源或新增驗收協議。該 source HEAD 與已建 binary 的建置對應尚未證明，未運行時的拒絕也未實發重現；未來先核對當下 endpoint 契約，不把此條當所有版本必然行為。兩協議單獨啟動的上例不保證能用該批次種值 endpoint。

## PostgreSQL Docker

先確認 Docker Desktop／daemon 與現有 image；記錄實際 image tag／digest。使用唯一 container 名稱，publish 只綁 `127.0.0.1`，不要接入使用者既有 DB。

在 run 目錄建立權限 0600 的 env file，使用本輪產生的測試帳密。只把檔案路徑交給 Docker，不把密碼放 command text／報告；UI 輸入也避免裸露截圖。

```sh
hc run --bg --name "$PG_TERMINAL" --cwd "$RUN_DIR" \
  docker run --name "$OWNED_CONTAINER" \
  --label "recording-e2e.run=$RUN_ID" \
  -p "127.0.0.1:$PG_PORT:5432" \
  --env-file "$RUN_DIR/postgres.env" "$PG_IMAGE"
```

等待實際 ready，再從產品 UI 建 connector／記錄／schema。資料庫本身可由 Docker 初始化；待驗 recording tables 仍只能由產品明確 schema apply 建立。

用容器內 `psql`／合適 client 做唯讀獨立查詢；SQL 使用當前 group destination 實際 table／columns，不硬猜舊名稱。不得以直接 INSERT 補出預期資料。

## SQLite 與停止

SQLite target 放 owned namespace，與 Gateway 設定 DB 分開；不存在時只 stat，存在時用正確 escaped read-only URI。含 `#` 等字元的路徑不可直接拼 URI。

收尾先停本輪記錄，再關閉自有服務／瀏覽器。停止／刪除容器前核對 container 名稱與 `recording-e2e.run` label；移除 volume 也只限此 run 新建且確認 ownership 的 volume。保留證據再清理。

斷線測試只暫停自己啟動並核對 PID 的 simulator，最後確實恢復或關閉。不要停止其他工作者的 HSL、Docker、hc core 或使用者開發服務。
