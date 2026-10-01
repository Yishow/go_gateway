# CMS 整合契約 (External History Query Contract)

本文件定義外部 CMS (如 Next.js) 如何查詢 go-gateway 的時序資料。

## 概覽

go-gateway 負責資料收集、轉換和儲存，但**不提供內建的歷史查詢 API**。
外部 CMS 系統直接連接資料庫讀取歷史資料，使用以下公開的 Schema 和索引。

---

## 資料庫連線

### PostgreSQL (生產環境)

- **必須使用唯讀帳號**以避免誤修改資料
- 時序表使用月度分區，命名格式：`timeseries_YYYY_MM`

### SQLite (測試環境)

- 時序表為單一表格，無分區

---

## 時序表結構 (timeseries)

| 欄位名       | 類型             | 說明          |
| ------------ | ---------------- | ------------- |
| `id`         | BIGINT/INTEGER   | 主鍵 (自增)   |
| `tag_id`     | UUID/TEXT        | 標籤 ID (FK)  |
| `ts`         | TIMESTAMPTZ/TEXT | 時間戳記      |
| `value_num`  | DOUBLE/REAL      | 數值 (互斥)   |
| `value_text` | TEXT             | 文字 (互斥)   |
| `value_bool` | BOOLEAN/INTEGER  | 布林 (互斥)   |
| `raw_value`  | JSONB/TEXT       | 原始值 (選用) |
| `quality`    | VARCHAR/TEXT     | 品質標誌      |

### 資料型別對應

| DataType                                                      | 使用欄位   |
| ------------------------------------------------------------- | ---------- |
| bool                                                          | value_bool |
| int16, uint16, int32, uint32, int64, uint64, float32, float64 | value_num  |
| string                                                        | value_text |

### 品質標誌 (Quality)

| 值          | 說明                |
| ----------- | ------------------- |
| `good`      | 資料品質良好        |
| `bad`       | 讀取失敗            |
| `uncertain` | 不確定 (如使用舊值) |

---

## 標籤表結構 (tags)

| 欄位名         | 類型              | 說明            |
| -------------- | ----------------- | --------------- |
| `id`           | UUID/TEXT         | 主鍵            |
| `key`          | VARCHAR(128)/TEXT | 標籤鍵 (顯示用) |
| `key_lower`    | VARCHAR(128)/TEXT | 小寫鍵 (唯一)   |
| `display_name` | VARCHAR(255)/TEXT | 顯示名稱        |
| `unit`         | VARCHAR(50)/TEXT  | 單位            |
| `data_type`    | VARCHAR(20)/TEXT  | 資料型別        |
| `status`       | VARCHAR(20)/TEXT  | 狀態            |
| `labels`       | JSONB/TEXT        | 標籤屬性        |

---

## 建議的查詢索引

已建立的索引供 CMS 高效查詢：

```sql
-- 依標籤和時間範圍查詢 (最常用)
CREATE INDEX idx_timeseries_tag_ts ON timeseries(tag_id, ts DESC);

-- 依時間範圍查詢所有標籤
CREATE INDEX idx_timeseries_ts ON timeseries(ts DESC);

-- 篩選非良好品質資料
CREATE INDEX idx_timeseries_quality ON timeseries(quality) WHERE quality != 'good';

-- 標籤查詢
CREATE INDEX idx_tags_key_lower ON tags(key_lower);
CREATE INDEX idx_tags_status ON tags(status);
```

---

## 常見查詢範例

### 1. 查詢單一標籤的時間範圍資料

```sql
SELECT ts, value_num, quality
FROM timeseries
WHERE tag_id = $1
  AND ts BETWEEN $2 AND $3
ORDER BY ts ASC;
```

### 2. 查詢多個標籤的最新值

```sql
SELECT DISTINCT ON (tag_id)
       tag_id, ts, value_num, value_text, value_bool, quality
FROM timeseries
WHERE tag_id = ANY($1)
ORDER BY tag_id, ts DESC;
```

### 3. 聚合查詢 (小時平均)

```sql
SELECT
    tag_id,
    date_trunc('hour', ts) AS hour,
    AVG(value_num) AS avg_value,
    MIN(value_num) AS min_value,
    MAX(value_num) AS max_value,
    COUNT(*) AS sample_count
FROM timeseries
WHERE tag_id = $1
  AND ts BETWEEN $2 AND $3
  AND value_num IS NOT NULL
  AND quality = 'good'
GROUP BY tag_id, date_trunc('hour', ts)
ORDER BY hour ASC;
```

### 4. 取得標籤清單 (含最新值)

```sql
SELECT
    t.id,
    t.key,
    t.display_name,
    t.unit,
    t.data_type,
    t.labels,
    (
        SELECT ts::text
        FROM timeseries
        WHERE tag_id = t.id
        ORDER BY ts DESC
        LIMIT 1
    ) AS last_ts,
    (
        SELECT value_num
        FROM timeseries
        WHERE tag_id = t.id AND value_num IS NOT NULL
        ORDER BY ts DESC
        LIMIT 1
    ) AS last_value
FROM tags t
WHERE t.status = 'active'
ORDER BY t.key;
```

### 5. 查詢指定標籤標籤 (labels)

```sql
-- PostgreSQL (使用 JSONB 操作符)
SELECT id, key, display_name
FROM tags
WHERE labels @> '{"area": "factory_a"}'
  AND status = 'active';

-- SQLite (使用 JSON 函數)
SELECT id, key, display_name
FROM tags
WHERE json_extract(labels, '$.area') = 'factory_a'
  AND status = 'active';
```

---

## 分區注意事項 (PostgreSQL)

### 分區命名

| 間隔 | 命名格式              | 範例                  |
| ---- | --------------------- | --------------------- |
| 月度 | timeseries_YYYY_MM    | timeseries_2026_01    |
| 每日 | timeseries_YYYY_MM_DD | timeseries_2026_01_18 |
| 每週 | timeseries_YYYY_WXX   | timeseries_2026_W03   |

### 跨分區查詢

PostgreSQL 會自動裁剪分區，但建議：

- 查詢時**必須包含 ts 條件**以啟用分區裁剪
- 避免查詢過長時間範圍 (如一年)

```sql
-- 好：會自動裁剪到相關分區
SELECT * FROM timeseries
WHERE tag_id = $1 AND ts BETWEEN '2026-01-01' AND '2026-01-31';

-- 差：掃描所有分區
SELECT * FROM timeseries
WHERE tag_id = $1;
```

---

## 資料保留政策

**go-gateway 不會自動刪除歷史資料**。

若需要資料保留策略，CMS 或 DBA 需自行處理：

```sql
-- 刪除 90 天前的資料 (PostgreSQL)
DELETE FROM timeseries WHERE ts < NOW() - INTERVAL '90 days';

-- 或直接刪除舊分區 (更高效)
DROP TABLE IF EXISTS timeseries_2025_01;
```

---

## 視圖建議 (選用)

若 CMS 需要簡化查詢，可建立以下視圖：

### 最新值視圖

```sql
CREATE VIEW v_tag_latest AS
SELECT DISTINCT ON (t.id)
    t.id AS tag_id,
    t.key,
    t.display_name,
    t.unit,
    t.data_type,
    ts.ts AS last_ts,
    ts.value_num,
    ts.value_text,
    ts.value_bool,
    ts.quality
FROM tags t
LEFT JOIN timeseries ts ON ts.tag_id = t.id
WHERE t.status = 'active'
ORDER BY t.id, ts.ts DESC NULLS LAST;
```

### 標籤摘要視圖

```sql
CREATE VIEW v_tag_summary AS
SELECT
    t.id AS tag_id,
    t.key,
    t.display_name,
    t.unit,
    t.labels,
    COUNT(ts.id) AS total_records,
    MIN(ts.ts) AS first_ts,
    MAX(ts.ts) AS last_ts,
    AVG(ts.value_num) AS avg_value
FROM tags t
LEFT JOIN timeseries ts ON ts.tag_id = t.id AND ts.quality = 'good'
WHERE t.status = 'active'
GROUP BY t.id, t.key, t.display_name, t.unit, t.labels;
```

---

## 版本相容性

| go-gateway 版本 | Schema 版本 | 備註     |
| --------------- | ----------- | -------- |
| 1.0.x           | 001         | 初始版本 |

當 Schema 更新時，將更新此文件並提供遷移指引。
