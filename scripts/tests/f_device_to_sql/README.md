# Device → SQL 驗收腳本（F）

真實路徑：瀏覽器（Playwright，借用 `frontend/node_modules`）→ 內嵌 UI 的 `cmd/test_ui` → 兩個 loopback Modbus simulator（`cmd/f_modbus_simulator`）→ 可丟棄的目的地資料庫。目的地只有 gateway 會寫入；腳本不 stub API、不直接 INSERT 取代採集。

前置：`make sync-frontend-static` 可執行（腳本自己建置）、`sqlite3` CLI、Playwright 的 chromium（`cd frontend && npx playwright install chromium`）。PostgreSQL 另需可丟棄的容器：

```
docker run -d --name gw-wg-pg-test -e POSTGRES_PASSWORD=gwtest -e POSTGRES_DB=gwtest -p 55432:5432 postgres:16
```

```
node scripts/tests/f_device_to_sql/run.mjs sqlite
POSTGRES_DSN='host=127.0.0.1 port=55432 user=postgres password=gwtest dbname=gwtest sslmode=disable' \
  node scripts/tests/f_device_to_sql/run.mjs postgres     # 沒有 POSTGRES_DSN 時直接失敗（blocked），不會用 SQLite 代替
node scripts/tests/f_device_to_sql/faults.mjs             # 試寫、設備靜默、目的地離線＋kill -9、poison
node scripts/tests/f_device_to_sql/quality.mjs            # 固定 clock 的亂序／重複／缺值／bad／stale／靜默／late／partial
```

- 預設 gateway port `3343`（`GW_PORT` 可改），工作目錄 `/tmp/gw-f-<kind>`，每次重建。
- 結果與 witness（source SHA、是否 dirty、平台、Go／Node 版本、錄到的 UI 與 SQL 值）寫到 `docs/plans/studio-v2-write-groups/evidence-f/*.json`，截圖同目錄。
- 模擬器暫存器：A = `[215, 1013, 32, 0, 0, 1, 1]`（batch = 0x0020_0000_0000_0001 = 9007199254740993），B = `[187, 777, 0, 0, 0, 2, 0]`；目的地以 SQL 獨立查回。
- 預設為 macOS arm64 單機驗證；Linux另有pinned container實測；Windows／ARM部署／LAN／真PLC／SCADA未驗證。

`quality.mjs` 自行建置 `f_write_group_fixture` 驗收版本，使用唯一的 `/tmp/gw-f-quality-<run>` 與 ownership marker；開啟 DB 前即拒絕其他路徑。127.0.0.1:3355 的控制只在這個 build tag 存在，一般 binary 沒有控制 listener。數值與來源身分由真正 Modbus read／runtime mapping 產生；暫存樣本明示尚未 ACK，釋放後才經 production journal／bucket closure／outbox／SQL sender。控制 clock 不改 production policy，結果以獨立 SQL 查詢及真 UI／delivery API 比對。`partial` 與 max-age fixture 使用既有 canonical Save/Apply，沒有直接改內部設定資料。


F2.2–F2.4 精確故障驗證（只限本機 owned fixture）：

```
node scripts/tests/f_device_to_sql/recovery.mjs
POSTGRES_DSN='<owned loopback fixture DSN>' node scripts/tests/f_device_to_sql/recovery.mjs postgres
node scripts/tests/f_device_to_sql/capacity.mjs
node scripts/tests/f_device_to_sql/worker-overlap.mjs
node scripts/tests/f_device_to_sql/test-write.mjs
F_SQL_KIND=postgres POSTGRES_DSN='<owned loopback fixture DSN>' node scripts/tests/f_device_to_sql/test-write.mjs
```

- `GW_PORT`、`F_FIXTURE_PORT`、`F_SIM_PORT_A/B`可隔離同時執行的harness；worker-overlap另使用3403/3405，必須先確認空閒。每個run都有唯一owned目錄及marker；不操作其他資料庫/container。PG recovery/test-write共用55433 target proxy，兩者須串行。
- `recovery`從真實ACK／closure transaction／target commit／local receipt的不同位置SIGKILL，查frozen backlog、Receipt dedupe與None unknown；跨connector健康scope要持續交付，不能用固定成功回覆。
- `capacity`用tag-only quota ENV與實際SQLite max_page_count造成SQLITE_FULL；reserved write lock只保持metadata可讀且阻擋INSERT，不寫samples。成功closure可清除consumed journal，獨立SQL provenance/receipt才是交付證據。poison DDL僅作用自有外部table。CAS使用真實200/409，延遲的是actual GET回覆，不偽造資料。
- `worker-overlap`兩個production pipelines共用owned store，secondary只在tagged binary以白名單test node identity避免被當同node restart；production NodeID/fencing不變。tag-only F_FIXTURE_MAX_RETRIES=2驗證blocked保留payload，production default仍0無上限。
- `test-write` preview/confirm由真UI操作、typed SQL/readback/精確owned cleanup及retained operation GET。SQLite cleanup-before-commit掛起trigger精確限定preview owner_value，不作用neighbor；真三分鐘lease後only cleanup。PG一開始透過真UI保存run專用非superuser帳號，再對本run readings分別撤銷SELECT或DELETE；connector identity與既有ownership guard不變。SQLite權限N/A不能當PG PASS。
- Linux `linux.Dockerfile`/`linux-run.mjs`固定Go1.25.5與Playwright版本，複製bounded tracked/untracked source至ownedcontainer、輸出source manifest及platform/build witness；`F_SOURCE_SHA`/`F_DIRTY_WORKTREE`只保留host source identity、不要求container可讀.git。amd64 emulation不是原生Linux硬體或ARM部署驗收；最新結果以environment/mixed witness為準。
