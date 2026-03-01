# Datalink Flow Gate v1 規劃（Review 版）

> 目標：強化 datalink 從「選擇設備」到「資料寫入 DB」的端到端流程，先做可阻擋錯誤、可觀測、可驗收的 v1。

---

## 1. 目標與成功定義

### 核心目標
1. 將流程關卡化（Gate），避免錯誤設定直接進入 runtime。
2. 建立上線前驗證路徑（Probe / Dry-run / Preview）。
3. 建立上線後驗證路徑（Smoke / Final / Soak + Gate Check）。
4. 所有結果可追蹤（報告可比對、錯誤可定位）。

### 成功定義（DoD）
- 設備未通過 Probe 不可啟用。
- 批次點位建立需可先 Dry-run，錯誤可回報。
- Mapping 上線前可 Preview，失敗會阻擋。
- Runtime 指標與寫入結果可由 report 證明（mismatch=0, bad_rows=0）。

---

## 2. 範圍與不做項目

### In Scope（本次）
- Device Activation Gate
- Point Batch Dry-run Gate
- Mapping Preview Gate
- Runtime/DB Outcome Gate
- 報告與通知一致化

### Out of Scope（下一版再做）
- UI 大改版
- 歷史報告儀表板前端
- 多租戶與細粒度權限

---

## 3. 流程關卡設計（v1）

### Gate A：設備啟用前驗證（Activation Gate）
- 條件：
  - 可連線（connect ok）
  - 可讀首筆（read probe ok）
- 行為：
  - 失敗：禁止 active，回傳原因（host/port/slave/function/address）
  - 成功：允許 active，記錄 `last_probe_at`、`probe_latency_ms`

### Gate B：點位建模前驗證（Point Dry-run Gate）
- 條件：
  - 地址/功能碼/資料型別合法
  - 同設備內不重複 address+function
- 行為：
  - dry-run 模式下只回錯誤清單，不落 DB
  - apply 模式僅在無阻斷錯誤時才批次建立

### Gate C：映射上線前驗證（Mapping Preview Gate）
- 條件：
  - pipeline 可解析
  - sample value 執行成功
  - output 可 cast 到 tag data_type
- 行為：
  - 失敗阻擋 mapping 啟用
  - 成功可啟用 mapping

### Gate D：採集→寫庫結果驗證（Outcome Gate）
- 固定檢查：
  - `pass=true`
  - `mismatch_count=0`
  - `bad_rows=0`
  - `total_rows>=1`
- 失敗處理：
  - 立即 fail stage + system event 通知

---

## 4. 實作拆分（WBS）

## Milestone 1（先堵住錯誤入口）
1. device service：加入 `ProbeAndActivate`（activate 前必 probe）
2. point service：加入 batch dry-run API/command
3. mapping service：強化 preview 驗證回傳（阻擋型錯誤）

## Milestone 2（強化結果可驗收）
4. outcome gate：統一走 `scripts/check_modbus_report.py`
5. soak pipeline：`scripts/run_modbus_soak.sh` 對接 gate + fail notify
6. loadtest report：固定輸出 `expected_by_key/mismatch_by_key`

## Milestone 3（流程穩定化）
7. 自動重啟與卡住保護：`scripts/copilot_watchdog.sh`
8. 文件與 runbook：標準操作、失敗排查、重跑流程

---

## 5. 預計變更檔案

### 新增
- `docs/datalink_flow_gate_v1_plan.md`（本檔）
- `docs/modbus_soak_workflow.md`
- `scripts/run_modbus_soak.sh`
- `scripts/check_modbus_report.py`

### 調整（既有）
- `cmd/loadtest_modbus/main.go`
- `internal/datalink/device/*`（Probe/Activate flow）
- `internal/datalink/point/*`（batch dry-run）
- `internal/datalink/mapping/*`（preview gate）

---

## 6. 驗收計畫

### 功能驗收
1. 啟用設備前，故意給錯 host/port：應阻擋 active。
2. 批次點位包含非法 address：dry-run 應列出錯誤且不落 DB。
3. mapping 用錯 data_type：preview 應失敗且不可啟用。

### 壓測驗收
- smoke 60s：PASS
- final 10m：PASS
- soak 30m：PASS
- 三階段皆需 mismatch=0 / bad_rows=0

---

## 7. 風險與對策

1. **Probe 增加啟用耗時**
   - 對策：可設 timeout/retry，上限可配置。
2. **dry-run 增加流程步驟**
   - 對策：提供一鍵 `--apply-if-clean`。
3. **舊資料不符合新 gate**
   - 對策：提供 migration/兼容模式（僅警告不阻擋）。

---

## 8. Review 決策點（請你拍板）

1. 設備啟用是否「必 probe 成功」才可 active？（建議：是）
2. 舊設備是否套「寬限模式」（先警告不阻擋）？（建議：是，1 週）
3. 點位 dry-run 是否變成批次建立預設流程？（建議：是）
4. 失敗通知是否固定用 system event + TG？（建議：是）

---

## 9. 實作節奏（審核後）

- Day 1-2：Gate A + Gate B
- Day 3：Gate C
- Day 4：Gate D + 壓測回歸
- Day 5：文件、runbook、交付報告

---

## 10. Gate v1.1（resolver + unique 策略收斂）

### 10.1 Mapping resolver 強制化（避免 cast gate 被略過）
- `mapping.Service` 在 **Create** 與 **Update(enabled=true)** 的 Gate C 驗證路徑，改為必須有 `tag resolver`。
- 未注入 resolver 時，會明確回傳可診斷錯誤：`tag resolver 未注入，無法執行 mapping 啟用驗證`，不再靜默略過 `cast` 檢查。
- 主要啟動路徑已改為建構時注入：
  - `cmd/test_ui/main.go`
  - `cmd/loadtest_modbus/main.go`
  - 使用 `mapping.NewServiceWithTagResolver(...)`。

### 10.2 Point unique 與 dry-run 對齊
- 現況差異（v1）：
  - DB 約束：`(device_id, address)`。
  - dry-run：`address+function`，且額外阻擋純 `address` 重複。
- v1.1 對齊後：
  - DB 目標約束：`(device_id, address, function)`（`function` 改為 non-null，預設 `''`）。
  - dry-run 僅阻擋 `address+function` 重複（移除純 `address` 阻擋）。
- 相關 migration：
  - PostgreSQL: `003_align_points_unique_function.up.sql` / `.down.sql`
  - SQLite: `003_align_points_unique_function_sqlite.up.sql` / `.down.sql`
  - 基線 schema 同步更新：`001_initial_schema.up.sql`、`001_initial_schema_sqlite.sql`

### 10.3 部署前檢查（升級前 / 回滾前）
- 新增 precheck 腳本：`scripts/check_points_unique_conflicts.sh`
- 檢查項目：
  - `--stage up`：檢查 `(device_id,address,function)` 唯一鍵衝突，並額外檢查 function alias 正規化後的潛在衝突。
  - `--stage down`：檢查回滾目標鍵 `(device_id,address)` 是否會衝突。

#### PostgreSQL
```bash
# 升級前
scripts/check_points_unique_conflicts.sh \
  --db postgres \
  --stage up \
  --dsn "$DATABASE_URL"

# 回滾前
scripts/check_points_unique_conflicts.sh \
  --db postgres \
  --stage down \
  --dsn "$DATABASE_URL"
```

#### SQLite
```bash
# 升級前
scripts/check_points_unique_conflicts.sh \
  --db sqlite \
  --stage up \
  --sqlite-file /path/to/gateway.db

# 回滾前
scripts/check_points_unique_conflicts.sh \
  --db sqlite \
  --stage down \
  --sqlite-file /path/to/gateway.db
```

### 10.4 Migration up/down 指令範本與安全檢查
- 新增 migration tooling：`scripts/run_points_unique_migration.sh`
- 特性：
  1. 先跑 precheck，未通過即阻擋。
  2. 不加 `--execute` 時只輸出執行指令範本。
  3. `down` 必須加 `--confirm-down` 才會執行。

```bash
# 只看 up 指令範本（不執行）
scripts/run_points_unique_migration.sh \
  --db postgres \
  --action up \
  --dsn "$DATABASE_URL"

# 實際執行 down（必須明確確認）
scripts/run_points_unique_migration.sh \
  --db sqlite \
  --action down \
  --sqlite-file /path/to/gateway.db \
  --execute \
  --confirm-down
```

### 10.5 回滾風險說明（需先決策）
1. **資料壓縮風險**：`down` 回到 `(device_id,address)` 後，同 address 不同 function 的資料需先合併或刪除。
2. **語意差異風險**：function alias（例如 `''` / `03` / `holding`）在新規則會正規化為同一鍵，舊資料可能先天衝突。
3. **操作時序風險**：若未先停寫就回滾，可能在 precheck 與 migration 間產生新衝突。

---

## 11. Gate v1.2（規則一致化 + 部署工具化）

### 11.1 規則一致化（Point 單筆 vs 批次）
- 單筆 Create/Update 與 Batch dry-run 使用同一套驗證與正規化：
  - `address` 正規化
  - `function` 正規化（如 `''/3/03/fc03/holding` -> `03`）
  - `data_type` 驗證
- 錯誤訊息統一：`point 驗證失敗: field=... value=... reason=...`

### 11.2 部署工具化（Precheck + Migration + 回滾）
- 新增：`scripts/check_points_unique_conflicts.sh`
  - 上版前檢查 `up`（目標鍵 + alias 正規化衝突）
  - 回滾前檢查 `down`（舊鍵 `(device_id,address)` 衝突）
- 新增：`scripts/run_points_unique_migration.sh`
  - 先 precheck，再 migration
  - `down` 需要 `--confirm-down`
  - 支援 `--down-conflict-strategy keep-latest`
- 新增：`scripts/resolve_points_down_conflicts_keep_latest.sh`
  - 回滾衝突時保留最新 `updated_at/created_at` 的資料列

### 11.3 Gate 入口收斂
- Makefile 新增：
  - `gatev12`（smoke + final）
  - `points-precheck-up` / `points-precheck-down`
  - `points-migrate-up` / `points-migrate-down`

### 11.4 本案採用決策（已採納）
1. DB 目標：**PostgreSQL 優先**
2. down 衝突策略：**keep-latest（保留最新）**
3. 部署方式：**採停寫窗口**（避免 precheck 與 migration 間資料漂移）

### 11.5 正式部署建議流程（PostgreSQL）
```bash
# 0) 進入停寫窗口
# 1) precheck
make points-precheck-up POINTS_DB_TYPE=postgres POINTS_DSN="$DATABASE_URL"

# 2) migration up
make points-migrate-up POINTS_DB_TYPE=postgres POINTS_DSN="$DATABASE_URL"

# 3) gate 驗證（至少 smoke + final）
make gatev12 GATE_WORKDIR=/home/yishow/github/go_gateway

# 4) 如需回滾（仍在停寫窗口）
make points-precheck-down POINTS_DB_TYPE=postgres POINTS_DSN="$DATABASE_URL"
make points-migrate-down POINTS_DB_TYPE=postgres POINTS_DSN="$DATABASE_URL" POINTS_DOWN_STRATEGY=keep-latest
```

