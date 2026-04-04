## Why



來源規則建立器目前只讓使用者選「協議讀取」用的 Point `data_type`，自動建立的 Tag 與 Point 必須同型別，且 Mapping 的 `transform_pipeline` 為空。使用者若已預期將暫存器數值轉成工程上常用的另一種型別（例如 **讀對 `uint16` 後 Tag 想用 `float64` 存**），或需要 **縮放／單位換算**（例如 `scale`／offset），只能事後在 Tag／Mapping 手動調整，與「規劃當下」的心智模型脫節。



另有一類問題與**型別 cast 無關**：多暫存器浮點或 32-bit 數值**解碼字節序**錯誤會導致數值完全錯亂（例如 float 顯示與實值無關）。目前 **`data_format`（ABCD／BADC／CDAB／DCBA）** 在連線設定中的可見性與適用範圍不一致（例如 MC 3E 有、Modbus 家族於 UI schema 未必露出），**不同協議連線**都應在**資料抓取解碼**所需時可設定，並與 `lib/hsllogic` 等解碼路徑一致。



## What Changes



- **規則／映射語意**

  - 在 **Studio 來源步驟／規則建立器**中，於選擇 Point 資料型態時，一併提供 **目標／預想資料型態**（預設與 Point 相同；可選與協議型別不同），涵蓋 **「讀對 uint16，Tag 用 float64」** 等情境。

  - 可選提供 **縮放／單位** 相關之規劃欄位（例如 scale、offset），於同步 Mapping 時預設插入 **`scale`** 步驟（與既有 `TransformParamsScale` 一致），並經 **`runtime.ingestor` → `mapping.ExecutePipeline`** 執行；與 **cast** 的順序（先 cast 再 scale 或依驗證規則）在 design 中固定。

- **連線解碼**

  - **修正並補齊 `data_format`**：凡涉及多暫存器／浮點解碼之協議，**裝置連線設定** SHALL 暴露 `data_format`（或等價命名）供選擇；**Modbus TCP／RTU／UDP** 等與 **MC 3E** 等連線類型皆依實際解碼需求對齊（見 `protocol-connectors` delta）。

- **文件化「抓取時所需參數」清單**：在 design 中列出裝置連線層（逾時、slave、字節序、協議特有欄位）與點位層（address、`data_type`）與映射層（管線）之分工，避免與「目標型別／縮放」混淆。



## Capabilities



### New Capabilities



- `source-rule-target-datatype`：規則層「協議讀取型別」「目標 Tag 型別」「可選縮放／偏移」與 Mapping 預設管線（`cast`／`scale`）之契約與驗證。



### Modified Capabilities



- `source-rule-runtime`：來源規則模型與同步至 Point／Tag／Mapping 時，納入目標型別與預設轉換管線（含可選 scale）之產出規則。

- `mapping-pipeline`：由來源規則衍生之預設管線（`cast`、`scale`）與手動覆寫／合併策略。

- `datalink-workbench-desktop`：規則建立器 UI（目標型別、可選縮放；連線表單露出 `data_format` 等）。

- `protocol-connectors`：連線 `config_schema` 與執行路徑對 **字節序／`data_format`** 的暴露與行為一致性（跨 Modbus 與其他適用協議）。



## Impact



- **後端**：`internal/datalink/sourcerule`、`mapping`、`connector/registry`（及各 adapter 是否讀取 `data_format`）、`schema` 遷移。

- **前端**：Workbench 裝置連線表單（動態 schema）、規則表單、型別與 i18n。

- **執行路徑**：仍單一路徑 **`ingestor` → `ExecutePipeline`**；**位元組序**在**讀取解碼**階段處理，**cast／scale**在映射管線處理。

- **相容性**：新欄位可選；既有專案預設行為不變。


