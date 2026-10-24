## Why

目前 Studio 基本記錄以設備分組、每個點位各占一欄，無法直接表達同一迴路的瞬時值與累積值採不同週期記錄。操作人員需要簡單選擇迴路與記錄型態，累積值同時保留原始讀值及可信的區間差值，並讓新增點位與時間分區管理不再依賴擴充寬表欄位。

## What Changes

- 在既有 Studio V2 基本流程加入迴路與資料選擇，可分別啟用瞬時值、累積值；同一迴路可同時記錄功率與累積電量，支援批次套用並保留個別覆寫。
- 瞬時值可設 1–300 秒、預設 15 秒；累積值可設 1–600 分鐘、預設 10 分鐘。以上為使用者指定契約，記錄週期不更動設備採集週期。
- 累積記錄同時保存原始值與區間差值；明確呈現首次基準、重置、回捲、缺值與換表，避免零值假象及重啟後重複計量。
- 新基本記錄採固定欄位的 typed long table，沿用 canonical WriteGroup、版本確認、journal、outbox 與 SQL receipt；既有寬表／自訂群組維持原行為，轉換只影響明確確認後的新資料。
- PostgreSQL 支援普通長表及可選 TimescaleDB hypertable；SQLite 保持同一邏輯模型。TimescaleDB 能力先查證，未安裝時不偷偷安裝或靜默降級；不宣稱長表必然比較快。
- 長表首次建置沿用 schema 預覽與明確確認；已確認且允許自動重建的受管表若消失，依凍結布局恢復並重試耐久資料，不能改造同名異質表或偽稱歷史資料可恢復。
- 基本設定、進階群組、工作區啟動及實際 SQL 寫入狀態各有清楚用途與修復入口，不用 runtime running 代替寫入證據。

## Capabilities

### New Capabilities

- `circuit-recording-profiles`: 迴路來源選擇、雙記錄型態、週期邊界、原始累積值與差值，以及可恢復的計算基準。
- `long-table-group-recording`: canonical 群組的 typed long table、可選 hypertable、受管表恢復、耐久交付與前向切換。

### Modified Capabilities

- `guided-recording-workflow`: 新基本記錄改由迴路型態驅動長表群組，保留既有寬表編輯；更新首次準備與重新載入契約，說明進階群組及工作區責任。

## Impact

- 設定／API：internal/datalink/workspace/、internal/datalink/measurement/、internal/datalink/schema/migrations/、internal/api/handlers/ 與 frontend/src/types/、frontend/src/services/。
- 執行／儲存：internal/datalink/runtime/、internal/datalink/aggregation/、internal/datalink/grouppipeline/、internal/datalink/groupdelivery/、internal/datalink/recordingplan/ 及 cmd/test_ui/service_wiring.go；擴充既有路徑，不新增第二套排程或佇列。
- 畫面：frontend/src/features/datalink/workbench-v2/、對應 hooks、en／zh-TW locales 與 frontend/tests/；維持四步 setup 與現有 focused monitor 的交付證據契約。
- 延續 measurement-aggregation、recording-plan-storage、write-group-row-semantics 與 runtime-write-group-delivery 的既有品質、精度及原子性要求；新長表是獨立 layout，不改寫 timeseries-storage 舊月分區路徑。
- 唯一負責 change 為 add-circuit-recording-profiles。已封存的 streamline-studio-v2-recording-setup、enforce-write-group-sample-semantics、add-usage-and-aggregation 提供既有邊界；其完成證據不代表本案已落地。現行 add-windows-tray-runtime-logs 與工作樹內 launcher、mapping、legacy writer 修復不屬本案。
- 提案階段僅建立本規劃；後續已依使用者授權完成實作與模擬來源／可丟棄 SQLite、PostgreSQL、TimescaleDB 驗收。未連正式設備或資料庫、不遷移歷史資料；產品不自動安裝 extension。
