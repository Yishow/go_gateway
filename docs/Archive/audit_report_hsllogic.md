# HSLLogic 與 MC3E Adapter 程式碼審查報告

**日期**: 2026-01-21
**審查範圍**: `lib/hsllogic` 核心套件、`internal/datalink/connector/adapters/mc3e.go`、及相關測試檔案。

## 總結 (Executive Summary)

整體架構設計清晰，成功實現了仿 HSL Communication 的泛型與設計模式。然而，在 **陣列讀取 (Batch Read)** 與 **特定設備地址解析** 上發現了關鍵的邏輯錯誤，這將導致資料讀取不完整或解析錯誤。建議優先修復這些功能性 Bug，並移除冗餘代碼。

## 1. 嚴重缺陷 (Critical Bugs)

### 1.1 MC3E 陣列讀取邏輯錯誤 (Array Read Logic Failure)

**位置**: `internal/datalink/connector/adapters/mc3e.go` (Read 方法)
**嚴重性**: 🔴 高 (High)

**問題描述**:
目前 `Read` 方法無法正確處理陣列（多筆連續資料）的讀取。

1. **暫存器數量計算錯誤**: 當 `req.Count > 1` 時，程式碼直接設定 `wordCount = count`。對於 32 位元 (如 Int32, Float32) 或 64 位元資料類型，需要的暫存器數量應為 `count * 2` 或 `count * 4`。目前的邏輯會導致讀取的暫存器數量不足 (Data truncation)。
2. **數值轉換錯誤**: `c.dataConverter.RegistersToValue` 僅設計用於轉換「單一數值」。當傳入陣列資料的暫存器時，它只會轉換並返回第一個元素，導致陣列其餘部分遺失。

**建議修復**:

- 修正 `wordCount` 計算公式：`wordCount = req.Count * hsllogic.RegisterCountForDataType(...)`。
- 在 `hsllogic/data_types.go` 中新增 `RegistersToValues` 方法，支援將暫存器陣列轉換為 `[]interface{}` 或特定類型的 Slice (如 `[]int32`)。
- 更新 `mc3e.go` 以使用新的陣列轉換方法。

### 1.2 三菱地址解析 - Hex 設備處理錯誤 (Incorrect Hex Parsing)

**位置**: `lib/hsllogic/address_parser.go` (`parseMitsubishiAddress`)
**嚴重性**: 🔴 高 (High)

**問題描述**:
三菱 PLC 的 `W` (Link Register) 與 `B` (Link Relay) 等設備通常使用 **十六進位 (Hex)** 編號 (例如 `W10` 代表暫存器 16)。
目前的解析邏輯優先嘗試 `Atoi` (十進位)，解析失敗才嘗試 Hex。這會導致 `W10` 被錯誤解析為十進位的 `10`，而非正確的 `16`。`X` 和 `Y` 正確使用了八進位，但 W/B 的 Hex 邏輯有誤。

**建議修復**:

- 針對 `W`, `B`, `SB`, `SW` 等 Hex 編址設備，應強制使用 `ParseInt(..., 16, 32)` 進行解析，而非先嘗試十進位。

## 2. 架構與邏輯問題 (Logic & Architecture Issues)

### 2.1 冗餘的地址解析回退 (Redundant Fallback)

**位置**: `internal/datalink/connector/adapters/mc3e.go`
**嚴重性**: 🟡 中 (Medium)

**問題描述**:
`MC3EConnector` 內部保留了 `parseMC3EAddress` 函數，並在 `hsllogic.ParseAddress` 失敗時作為 fallback。這違反了 "Single Source of Truth" 原則。若 `hsllogic` 有 Bug，應直接在 `hsllogic` 修復，而非在 adapter 中維護另一套邏輯。且兩者的解析邏輯均存在上述 Hex 解析問題。

**建議修復**:

- 移除 `mc3e.go` 中的 `parseMC3EAddress` 及相關 fallback 邏輯。
- 確保 `hsllogic.ParseAddress` 覆蓋所有舊有邏輯的需求。

### 2.2 `ToByteOffset` 的定義模糊

**位置**: `lib/hsllogic/address_parser.go`
**嚴重性**: 🟡 中 (Medium)

**問題描述**:
`ToByteOffset` 對非 Siemens 協議一律回傳 `Offset * 2`。這假設了所有 Offset 都是以 Word (16-bit) 為單位。然而，Modbus Coil/Status 或三菱的 Bit 設備 (M, X, Y) 是位元地址。如果此方法被用於計算位元組偏移量 (例如在 Buffer 中的位置)，對於位元設備來說結果可能不準確 (視上下文而定)。目前程式碼中似乎未大量使用此方法，但留著容易造成誤用。

**建議修復**:

- 明確定義該方法的用途，或針對位元設備拋出錯誤/返回正確的位元組索引 (例如 `Offset / 8`)。

### 2.3 缺少 `intSliceToBytes` 實作

**位置**: `internal/datalink/connector/adapters/mc3e.go`
**嚴重性**: 🟡 中 (Medium)

**問題描述**:
`Read` 方法中呼叫了 `intSliceToBytes`，但在該檔案中未見其定義。如果未在同 package 其他檔案定義，這將導致編譯錯誤。

## 3. 優化建議 (Optimization Suggestions)

### 3.1 擴充 `ByteTransform` 支援陣列

建議在 `ByteTransform` 中加入批量轉換方法，例如：

- `TransByteToInt32s(buffer []byte, index int, count int) []int32`
  這能大幅簡化 Adapter 中處理陣列讀取的程式碼，並提升效能 (減少重複的 function call)。

### 3.2 PacketLogger 緩衝區策略

目前 `LogToFile` 在 channel buffer 滿時會丟棄舊日誌 (Drop oldest)。對於高可靠性需求的場景，建議增加配置選項，允許選擇 "Block" (阻塞等待寫入) 或 "Drop" (丟包)，並在丟包發生時輸出警告日誌。

## 4. 程式碼品質 (Code Quality)

- **註解完整性**: 優秀。所有 public method 都有繁中 JSDoc 風格註解。
- **測試覆蓋**: 核心邏輯 (`address_parser`, `bytetransform`) 測試案例豐富，但建議補上 `W` (Hex) 設備的解析測試，以驗證上述 Bug 修復。

## 接下來的行動 (Next Steps)

1.  **實作修正**: 修改 `hsllogic/address_parser.go` 修復 Hex 解析問題。
2.  **功能增強**: 在 `hsllogic/data_types.go` 增加 `RegistersToValues`。
3.  **Adapter 重構**: 修正 `mc3e.go` 的陣列讀取邏輯並移除 fallback 解析器。
