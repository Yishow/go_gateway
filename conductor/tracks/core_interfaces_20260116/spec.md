# Track Specification: 建立核心通訊與資料庫介面

## 1. 概述
本軌道的目標是實作 `go-gateway` 專案中最重要的抽象層：通訊協議介面 (`Protocol`) 與資料庫驅動介面 (`Driver`)。這些介面將定義系統如何與外部工業設備通訊以及如何與不同的資料庫系統互動。

## 2. 技術細節
### 2.1 通訊協議介面 (internal/protocol/protocol.go)
*   **介面名稱**: `Protocol`
*   **關鍵方法**:
    *   `Connect(ctx context.Context) error`
    *   `Disconnect() error`
    *   `Read(ctx context.Context, req *ReadRequest) (*ReadResult, error)`
    *   `Write(ctx context.Context, req *WriteRequest) error`
*   **數據結構**:
    *   `ReadRequest`, `ReadResult`, `WriteRequest`
    *   `DataType` (Coil, Register, etc.)
    *   `ValueType` (Bool, Int, Float, etc.)

### 2.2 資料庫驅動介面 (internal/database/driver.go)
*   **介面名稱**: `Driver`
*   **關鍵方法**:
    *   `Connect(ctx context.Context) error`
    *   `BatchInsert(ctx context.Context, table string, columns []string, rows [][]any) error`
    *   `Upsert(ctx context.Context, table string, columns []string, values []any, keyColumns []string) error`
    *   `GetTables(ctx context.Context) ([]string, error)`

## 3. 驗證標準
*   介面定義必須符合 `WBS.md` 與 `plan.md` 中的規格。
*   必須包含單元測試，確保數據結構的序列化與枚舉定義正確。
*   程式碼必須通過 `golangci-lint` (如果已配置) 與類型檢查。
