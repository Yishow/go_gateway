## Context



今日 **Source Rule** 只承載協議讀取所需的 `data_type`（對應 Point），`sourcerule` 同步建立 Tag 時強制 **Tag.data_type = Point.data_type**，且 Mapping **`transform_pipeline` 為空**。



數值轉換發生在 **`runtime/ingestor`** 對每條 Mapping 呼叫 **`mapping.ExecutePipeline`**；步驟定義於 **`schema.TransformStep`**（**`cast`**、**`scale`** 等）。



**讀取解碼**（多 word 的 float／int 如何拼成位元組）依連線與協議適配器，常涉及 **`data_format`**（`lib/hsllogic` 字節序）。此層與映射 **`cast`** 互補：**字節序錯誤**會導致「同一地址讀到錯的 float」，**無法**靠把 Tag 改成 `float64` 或事後 `cast` 修正。



## Goals / Non-Goals



**Goals:**



- **目標資料型態**：規則層新增可選 **目標型別**；預設等於 Point 讀取型別；當二者不同時自動產出合法 **`cast`** 步驟，滿足 **「讀對 uint16，Tag 想用 float64 存」** 等語意（在**已正確解碼**的前提下）。

- **縮放／單位**：規則層可選 **scale／offset**（或單一線性 `scale` 步驟參數），同步時寫入 Mapping 預設 **`scale`** 步驟，滿足 **工程單位換算** 而不僅型別轉換。

- **連線 `data_format`**：**修正** 連線設定 schema 與實作，使 **Modbus TCP／RTU／UDP** 等凡需要多暫存器解碼者，與 **MC 3E** 一致地可選 **ABCD／BADC／CDAB／DCBA**；實際解碼路徑與 `internal/datalink/connector/adapters`、`lib/hsllogic` 對齊。

- **參數分層說明**：在實作與文件區分 **連線層**（傳輸、逾時、slave、字節序）、**點位層**（address、`data_type`）、**映射層**（管線）之職責。



**Non-Goals:**



- 完整視覺化任意公式編輯器（可沿用既有 `formula` 步驟與後續 UI）。

- 在單一變更內支援所有協議的位元組序細節若某協議尚未接 hsllogic（需適配器先行）。



## Decisions



1. **目標型別欄位**：Source Rule 持久化 **`target_data_type`**（可空 = 與 Point 相同）。



2. **縮放欄位（可選）**：`scale_multiplier`／`scale_offset`（或合併為 `scale` 步驟參數結構）存於規則；同步 Mapping 時若非預設值則插入 **`scale`** 步驟；**與 `cast` 順序**：先 **`cast`** 至 Tag 語意型別，再 **`scale`**（若驗證不允許則調整並寫入測試）。



3. **Tag／Point**：`resolveRuleTag` 允許 Tag 使用 **目標型別**；Point 維持 **協議讀取型別**。



4. **`data_format`（連線）**  

   - **決策**：在 **`connector` 之 `protocolInfos` ConfigSchema** 中，為 **Modbus TCP／RTU／UDP**（及任何使用相同多暫存器解碼之適配器）新增 **`data_format`** 欄位，預設值與現行解碼預設一致；**MC 3E** 維持現有欄位。  

   - **實作**：Modbus adapter 讀取 `connection_config` 中的 `data_format` 並傳入 hsllogic／讀取路徑（與既有 MC3E 模式對齊）。  

   - **MQTT** 等非暫存器塊協議：不強制 `data_format`，或標為 N/A。



5. **相容矩陣**：`cast`／`scale` 組合須通過 **`ValidateTransformPipeline`**；連線儲存時校驗 `data_format` 枚舉。



## 資料抓取相關參數一覽（分層）



| 層級 | 主要參數 | 寫入位置 / UI 入口 | 職責邊界 | 本變更對齊 |
| --- | --- | --- | --- | --- |
| **裝置連線** | `host`、`port`、`serial`、`baud`、`parity`、`slave_id`、`timeout`、`persistent_connection`、`keep_alive`、**`data_format`** | `devices.connection_config`、Workbench Step 1 裝置編輯器、`connector` `ConfigSchema` | 決定如何連線，以及如何把多 register 原始值解碼成正確數值；**不負責** Tag 型別或工程單位換算 | Modbus TCP／RTU／UDP 與 MC 3E 皆可設定 `data_format`；空值仍沿用既有解碼預設 |
| **點位 Point** | **`address`**、**`data_type`**、polling 群組 | `points`、Source Rule 同步產生的 Point 定義 | 決定從哪裡讀、以哪種協議語意讀、需要幾個 register；若 Point `data_type` 不對，後續 `cast`／`scale` 無法補救錯誤解碼 | `data_type` 持續代表**協議讀取型別**，不因 Tag 目標型別而改寫 |
| **來源規則（規劃）** | `start`、`count`、`naming`、**`target_data_type`**、**`scale_multiplier`**、**`scale_offset`** | `source_rules`、Workbench Source Rule 表單 | 描述批次建立 Point／Tag／Mapping 的規劃意圖；自身不執行轉換，但會在同步時決定 Tag 型別與預設管線 | 新增 `target_data_type` 與 scale 欄位；空值表示沿用 Point 型別 / 不插入 scale |
| **映射 Mapping** | `transform_pipeline`（`cast`、`scale`、`formula`…） | `mappings.transform_pipeline`、runtime `ExecutePipeline` | 對**已正確解碼**的 Point 值做後續轉換；固定順序為 **`cast -> scale`** | Source Rule 同步依規則產生預設 `cast` / `scale` 步驟，並經 `ValidateTransformPipeline` 驗證 |

**補充：**

- `Tag.data_type` 是規則同步的**結果**，不是獨立的抓取參數層；當規則宣告 `target_data_type` 時，Tag 跟隨目標型別，而 Point 保持協議讀取型別。
- 排查數值異常時，優先順序為：**連線 `data_format`** → **Point `data_type` / address** → **規則的目標型別與縮放宣告** → **Mapping 管線**。



**注意**：「float 顯示 100.123 卻變成 0.488…」類問題，優先查 **連線 `data_format`** 與 **Point `data_type`** 是否與裝置一致，再查映射 **scale**。



## Risks / Trade-offs



- **[Risk] Modbus adapter 尚未傳遞 `data_format`** → **Mitigation**：後端單元測試多字組 float 在四種序下與裝置手冊對照。  

- **[Risk] cast 與 scale 順序錯誤** → **Mitigation**：以固定順序與 golden tests 鎖定。  

- **[Risk] 規則欄位膨脹** → **Mitigation**：縮放為可選折疊進階區塊。



## Migration Plan



1. DB：規則表新增 `target_data_type`、可選 scale 欄位；連線 JSON 新增鍵預設可選。  

2. 既有裝置：無 `data_format` 時使用與今日相同之解碼預設。  

3. 回滾：清空新欄位或還原 migration。



## Open Questions



- 規則更新時與**手動管線**衝突的處理策略（仍待產品定案）。  

- Modbus **預設 `data_format`** 與市場多數裝置（ABCD vs CDAB）是否需在 UI 標示建議。



## 支援的型別轉換矩陣



### Cast 轉換規則



系統支援以下資料型別之間的轉換：



| 來源型別 → 目標型別 | bool | int16 | uint16 | int32 | uint32 | int64 | uint64 | float32 | float64 | string |

|---------------------|------|-------|--------|-------|--------|-------|--------|---------|---------|--------|

| **bool** | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |

| **int16** | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |

| **uint16** | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |

| **int32** | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |

| **uint32** | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |

| **int64** | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |

| **uint64** | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |

| **float32** | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |

| **float64** | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |

| **string** | ⚠ | ⚠ | ⚠ | ⚠ | ⚠ | ⚠ | ⚠ | ⚠ | ⚠ | ✓ |



**說明：**

- ✓：直接支援轉換

- ⚠：string 轉數值時需為有效數字格式，否則可能產生非預期結果



### 常見使用案例



| 使用案例 | 來源型別 | 目標型別 | 說明 |

|----------|----------|----------|------|

| 暫存器原始值轉工程值 | uint16 | float64 | PLC 暫存器存整數，Tag 需浮點運算 |

| 縮放後的整數值 | int16 | float64 | 含正負號的原始值 |

| 大數值轉換 | uint32 | float64 | 雙暫存器無號整數 |

| 累計值轉換 | int64 | float64 | 64 位元累計器 |



## 轉換管線執行順序



### 固定順序規則



當 Source Rule 同時宣告 `target_data_type` 與 `scale` 參數時，系統依以下順序建構 Mapping 的 `transform_pipeline`：



```

1. cast（型別轉換）→ 先將原始值轉為目標型別

2. scale（縮放）→ 對已轉型的值進行線性縮放

```



**公式：** `result = (cast(raw_value) * scale_multiplier) + scale_offset`



### 管線範例



**案例：uint16 → float64 並縮放**



```json

{

  "transform_pipeline": [

    {

      "type": "cast",

      "order": 0,

      "params": { "target_type": "float64" }

    },

    {

      "type": "scale",

      "order": 1,

      "params": { "scale": 0.1, "offset": 0 }

    }

  ]

}

```



### 驗證規則



1. **型別相容性**：`cast` 步驟的 `target_type` 必須為有效的 DataType 枚舉值

2. **數值運算**：`scale` 步驟僅對數值型別有效（非 bool/string）

3. **管線完整性**：所有步驟須通過 `ValidateTransformPipeline` 驗證

4. **冪等性**：相同參數的規則重複同步應產生相同的管線



### 衝突處理策略



當規則同步發現 Mapping 已有手動建立的管線時：



1. **空管線**：直接插入規則衍生的 cast/scale 步驟

2. **僅 cast**：若目標型別相同，保留既有；若不同，標記衝突需人工確認

3. **僅 scale**：若參數相同，保留既有；若不同，標記衝突需人工確認

4. **完整管線**：不自動覆寫，需人工介入決定
