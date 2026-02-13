# p2-optimize-packet-aggregation 設計文件

## 1. 聚合演算法 (Aggregation Algorithm)

系統將在啟動時（或配置變更時）執行靜態分析，將所有的 `PollingPoint` 分組：

1.  **分組 (Grouping)**: 根據 `DeviceID`, `Protocol`, `MemoryArea` (如 D, M, X, Y) 進行分組。
2.  **排序 (Sorting)**: 在組內根據 `Address` 偏移量進行排序。
3.  **合併 (Merging)**:
    - 遍歷排序後的 Tag。
    - 若 `Current.Address` + `Wait` < `Next.Address`，且合併後的總長度 < `MaxPDU`，則合併為一個 `Block`。
    - 允許「空洞 (Holes)」：為了合併 D100 和 D105，系統會連同 D101-D104 一起讀取（犧牲少許頻寬換取 RTT）。

## 2. 元件架構

### 2.1 RequestOptimizer (新增)

負責執行上述演算法，產出 `OptimizedBlock` 列表。

- 輸入：`[]RawPoint`
- 輸出：`[]ReadRequest` (包含多個 Point 的引用)

### 2.2 ResponseDispatcher (新增)

收到聚合回應後，負責將數據「分發」回原始的 Point 物件。

- 解析 `RawBytes`。
- 根據 Offset 切割並轉換為 Int/Float/Bool。

## 3. 數據流

```mermaid
graph LR
    Config[Tag Config] --> Optimizer
    Optimizer -->|合併請求| Scheduler
    Scheduler -->|Read Block| Protocol[PLC]
    Protocol -->|Raw Bytes| Dispatcher
    Dispatcher -->|Slice & Decode| Datalink[Tag Values]
```

## 4. 權衡取捨

| 策略               | 說明               | 決策                                    |
| ------------------ | ------------------ | --------------------------------------- |
| 動態聚合 (Runtime) | 每次 Loop 重新計算 | 效能差，**不採用**                      |
| 靜態聚合 (Startup) | 啟動時計算一次     | 效能好，**採用**                        |
| 激進合併 (大空洞)  | 讀取大量無用數據   | 設定 `MaxHoleSize` 限制 (例如 10 words) |

## 5. 協議適配

不同協議需提供 `PDU Constraints`：

- **Modbus TCP**: Max 125 Registers (Read Holding Registers).
- **Fatek**: Max 64 Registers (視型號而定).
- **MC Protocol**: 可變長度。
