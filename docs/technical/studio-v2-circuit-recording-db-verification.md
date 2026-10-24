# Studio V2 迴路記錄：資料庫驗收

本報告對應 `add-circuit-recording-profiles` 任務 5.1，依主代理最終實跑結果整理。
驗收對象是未提交的工作樹；基準 commit、測試來源雜湊及逐案例結果見 [證據 JSON](evidence/circuit-recording-db-matrix-20261024.json)。
本報告不代表瀏覽器、Windows、真實 PLC 或現場驗收通過。

## 環境

| 模式 | 實測版本 | fixture |
| --- | --- | --- |
| SQLite | 3.51.2，modernc.org/sqlite | 每測試獨立暫存檔 |
| plain PostgreSQL | 17.11 | Timescale 映像中的普通資料表，不宣稱未安裝 extension |
| TimescaleDB | 2.31.0／PostgreSQL 17.11 | 明選 hypertable 模式 |
| plain PostgreSQL | 18.6，Timescale extension 不存在 | 獨立暫時容器 |

Go：`go1.27.3 darwin/arm64`，專案 `go.mod` 目標為 1.25.5。
兩個快取映像皆為 arm64，以 `--pull=never` 執行：

- `timescale/timescaledb:latest-pg17`：`sha256:b56a3f51991df1f79d0c1a0841cad5c5fa8058fdd6a3ee6886394a6d4c174538`。
- `postgres:18-alpine`：`sha256:77f585114c32fbca283dc835b0596f4e52b51b4c6662d7810b2f4084f60a1873`。

## 實跑入口與結果

fixture 憑證由執行環境注入 `POSTGRES_DSN`，不寫入文件或日誌。
API fixture 限定 loopback、port 55432、可丟棄的 gwtest；切換 PG18 時暫停 Timescale 容器，結束後復原。

```sh
go test -json -p 1 -count=1 \
  ./internal/datalink/dbtarget ./internal/datalink/grouppipeline ./internal/api \
  -run '^TestCircuitRecordingIntegration$'
```

| 階段 | 葉案例 PASS | FAIL | SKIP／NOT RUN |
| --- | ---: | ---: | --- |
| PG17＋Timescale fixture（包含 SQLite、plain、hypertable） | 43 | 0 | 0 |
| 無 extension 的 PG18 fixture（包含 SQLite、plain） | 31 | 0 | 12 個 Timescale 案例 |

PG18 的 12 個 skip 全部已在 Timescale 階段實際 PASS，逐項對照於證據 JSON；不將 skip 計為通過。
首輪 PG18 因使用其他 port 被 API fixture 守門拒絕，屬環境設定錯誤；未改產品或測試預期，改用既定 55432 後重跑以上最終結果。

## SQL 讀回與故障斷言

| 項目 | 實跑範圍 | 固定預期與證據 |
| --- | --- | --- |
| uint64 raw／delta、first baseline | SQLite／plain PG17、PG18／Timescale | 9007199254740992 → 9007199254740993，raw 精確、首窗差值 NULL、下一窗差值 1；另一來源為 100 → 130、差值 30 |
| 跨模式等價 | 同上 | 相同輸入的狀態、raw、delta、實際區間與 provenance 相等；不比較目的地特有的 effect scope |
| decimal 精度 | 同上 | 真計量引擎 → encoder → SQL：12000.0 → 12000.4 → 12000.7，差值 NULL／0.4／0.3；重送已提交記錄回 AlreadyCommitted |
| 窗內 reset | 同上 | firstStart+100 → firstStart+10 → firstStart+110，保留 raw 9007199254741102、delta NULL、reset_or_uncertain、provenance 含 decrease 證據；另一來源仍得 30 |
| gap span_total | 同上 | 缺窗 raw／delta NULL、baseline 不前移；第 9 至第 29 分鐘的可信差值為 2／30，status gap、實際起迄與 span_total=true 由 SQL 讀回 |
| SQL ACK 遺失 | 同上 | 目的地 commit 後遺失回應；重啟續送仍只有 4 列，outbox 全部 sql_committed、payload digest 不變、baseline 不重複前移 |
| 本地閉窗交易失敗 | SQLite／plain PG17、PG18 | 注入 baseline 寫入失敗，整筆 closure 回滾；樣本仍未 consumed，重試後只前進一次 |
| 缺表重建 | SQLite／plain PG17、PG18／Timescale | 重建已確認結構，generation=2、歷史可能遺失；共享群組 pending payload／receipt 保留，續送不重複；累積 baseline 不因表消失而歸零 |
| cutover／drain | 同上 | 舊寬表歷史保留，非對齊邊界窗帶 partial／truncated；舊 pending 送原表、邊界樣本屬新長表；新累積型態先 baseline、下一窗精確 delta |
| schema operation 回應遺失 | SQLite／plain PG17、PG18 | 查回原 operation，不補跑未知 DDL；本項不據此宣稱 Timescale 同路徑覆蓋 |

**decimal 證據層級**：目前沒有 Decimal PLC Tag 型別；此處證明計量引擎、編碼與三種 SQL 儲存的精確性，不宣稱 collector 支援 Decimal Tag。

## Timescale catalog 與跨 chunk

- 唯一時間維度為 bucket_start；catalog interval 以整數微秒比對 **604800000000**，不是四捨五入到秒。
- 複合主鍵 (bucket_start, record_id)，存在 source/time index；compression=false、目的地綁定 jobs=0。
- 獨立跨 chunk 案例經 encoder → outbox → sender 寫入 4 筆相隔超過 7 天的記錄，catalog chunks≥4 且 COUNT(DISTINCT tableoid)=4。
- 排序讀回 9007199254740993 至 9007199254740996；重送後資料與 receipts 仍各 4 列，直接重複 INSERT 被複合主鍵拒絕。
- 此跨 chunk 案例採已形成的 typed records，不冒稱為實體設備採集驗證。

## 清理與未驗證範圍

- 各 fixture 的 schema／SQLite 檔案由測試清理。兩階段 catalog 核對 gw_* schema 均為 0。
- 主代理建立的 PG18 容器已刪除；原有 gw-ts-6 已恢復 running（17.11／2.31.0）。未接觸正式資料庫或 PLC。
- NOT RUN：瀏覽器操作／截圖及完整交付檢查（5.2）、Windows、現場設備；MySQL／MariaDB／SQL Server 是另外的後續變更。
