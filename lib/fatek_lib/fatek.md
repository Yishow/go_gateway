# FATEK FBs 通訊協定實作規範 (Serial & Over TCP)

本文件根據 FATEK 官方通訊協定手冊 (附錄一) 整理，旨在指導 `go-gateway` 專案中 FATEK 驅動的開發。

## 1. 協定概述

FATEK PLC 採用 **ASCII** 格式進行通訊。無論是透過串列埠 (RS-232/485) 或乙太網 (TCP/IP)，其應用層的訊框格式皆相同。

*   **架構**: Master-Slave (主從式)。`go-gateway` 作為 Master 主動發起請求，PLC 作為 Slave 回應。
*   **編碼**: ASCII 字元 (除了 STX 和 ETX 為控制碼)。
*   **傳輸層**:
    *   **Serial**: 標準 COM Port 通訊 (預設 9600, 7, E, 1 或自訂)。
    *   **Over TCP**: 將完整的 ASCII 訊框 (包含 STX, ETX, Checksum) 直接作為 TCP Payload 發送。預設 Port 通常為 **500**。

---

## 2. 訊框格式 (Frame Structure)

所有通訊 (命令與回應) 皆遵循以下格式：

```text
|  STX  | Station | Command |   Body    | Checksum |  ETX  |
| 1 Byte| 2 Bytes | 2 Bytes | Var Bytes | 2 Bytes  | 1 Byte|
```

| 欄位 | 長度 | 說明 | 範例 (Hex) |
| :--- | :--- | :--- | :--- |
| **STX** | 1 | 起始字元，固定為 `0x02` | `02` |
| **Station** | 2 | 站號 (00~FE)，ASCII Hex | 1號站 -> `"01"` (`30 31`) |
| **Command** | 2 | 指令碼，ASCII Hex | Cmd 40 -> `"40"` (`34 30`) |
| **Body** | N | 數據內容 (視指令而定) | - |
| **Checksum** | 2 | 縱向冗餘校驗 (LRC)，ASCII Hex | `0xC7` -> `"C7"` (`43 37`) |
| **ETX** | 1 | 結束字元，固定為 `0x03` | `03` |

### 2.1 Checksum (LRC) 計算演算法

**關鍵規則**: Checksum 計算範圍包含 **STX**、Station、Command 到 Body 的最後一個字元 (不包含 Checksum 本身與 ETX)。

1.  將範圍內所有 Byte 的數值 (Hex) 相加。
2.  取總和的最低 1 Byte (Modulo 256)。
3.  將該 Byte 轉換為 2 個大寫 ASCII Hex 字元。

---

## 3. 效能優化 (Optimization)

為了支援高頻率的數據採集，驅動程式實作了以下優化：

### 3.1 Buffer Pool (sync.Pool)
使用 `sync.Pool` 重用 `bytes.Buffer` 物件，大幅減少封包建構時的記憶體分配 (Memory Allocation) 與 GC 壓力。

*   **API**: `GetBuffer()` / `PutBuffer()`
*   **應用**: `BuildFrameToBuffer` 直接寫入重用的緩衝區，避免產生暫時性的 `[]byte`。

### 3.2 Zero-Allocation Request
在 `FatekClient.execute` 流程中，從封包建構到發送至 TCP Socket，全程使用指標傳遞，不產生額外的數據拷貝。

### 3.3 Benchmark 結果
```text
BenchmarkBuildFrame_Alloc-22    141.6 ns/op    4 B/op    2 allocs/op
BenchmarkBuildFrame_Pool-22     138.7 ns/op    2 B/op    1 allocs/op
```
使用 Pool 後，單次封包建構的記憶體分配次數減半。

---

## 4. 核心指令詳解

以下列出 `go-gateway` 需優先實作的關鍵指令。

### 4.1 讀取連續單點狀態 (Cmd 44)
用於批量讀取 X, Y, M, S 等狀態。

*   **Request**:
    *   Command: `"44"`
    *   Body: `[數量N (2 Hex)]` + `[起始位址 (5 Char)]`
    *   *N 範圍: 01~FF (1~255)*
*   **Response**:
    *   Code: `"0"` (成功)
    *   Body: `[狀態0]` `[狀態1]` ... `[狀態N-1]` (每個狀態 1 Char: '0' 或 '1')

**範例**: 讀取 X0~X2 (共3點)
*   Req Body: `"03"` + `"X0000"`
*   Full Req: `STX` `"014403X0000"` `LRC` `ETX`

### 4.2 寫入連續單點狀態 (Cmd 45)
用於批量寫入 X, Y, M, S 等狀態。

*   **Request**:
    *   Command: `"45"`
    *   Body: `[數量N (2 Hex)]` + `[起始位址 (5 Char)]` + `[狀態數據...]`
    *   *狀態數據: N 個 '0' 或 '1'*
*   **Response**:
    *   Code: `"0"`
    *   Body: (無)

### 4.3 讀取連續暫存器 (Cmd 46)
用於批量讀取 R, D, RT, RC 等數值。

*   **Request**:
    *   Command: `"46"`
    *   Body: `[數量N (2 Hex)]` + `[起始位址 (6 or 7 Char)]`
    *   *N 範圍: 01~40 (最多 64 個 Word)*
*   **Response**:
    *   Code: `"0"`
    *   Body: `[數據0]` `[數據1]` ... (每個數據為 4 Chars Hex String)

**範例**: 讀取 D0~D1 (共2個)
*   Req Body: `"02"` + `"D00000"`
*   Resp Body: `"1234"` `"ABCD"` (表示 D0=0x1234, D1=0xABCD)

### 4.4 寫入連續暫存器 (Cmd 47)
用於批量寫入 R, D 等數值。

*   **Request**:
    *   Command: `"47"`
    *   Body: `[數量N (2 Hex)]` + `[起始位址 (6 or 7 Char)]` + `[數據...]`
    *   *數據: 每個 Word 為 4 Chars Hex String*
*   **Response**:
    *   Code: `"0"`

### 4.5 混合讀取 (Cmd 48) - *進階優化用*
可一次讀取不連續的點或暫存器。

*   **Request**:
    *   Command: `"48"`
    *   Body: `[數量N (2 Hex)]` + `[元件1]` + `[元件2]` ...
*   **Response**:
    *   Code: `"0"`
    *   Body: `[數據1]` `[數據2]` ...
    *   *數據格式: 單點回 '0'/'1' (1 Char)，暫存器回 Hex (4 Chars)*

---

## 5. 錯誤處理 (Error Handling)

當 PLC 回應的 Command 欄位後的第一個字元**不是** `'0'` (0x30) 時，代表發生錯誤。該字元即為錯誤碼 (Error Code)。

**錯誤回應格式**:
```text
| STX | Station | Command | ErrorCode | Checksum | ETX |
```
*注意: 錯誤時沒有 Body，ErrorCode 取代了原本 Code '0' 的位置。*

| 錯誤碼 | 說明 | 處理建議 |
| :--- | :--- | :--- |
| **2** | 非法數值 | 檢查發送的數據是否包含非 Hex 字元 |
| **4** | 非法格式/指令 | 檢查指令碼是否支援，LRC 是否正確 |
| **A** | 非法位址 | 檢查讀寫的位址是否超出 PLC 範圍 (Boundary) |

---

## 6. TCP 實作注意事項

1.  **無 Header**: 與 Modbus TCP 不同，FATEK TCP 通常**不增加**額外的 6-byte Header (如 MBAP)。它只是將 Serial 的 ASCII Bytes 透過 Socket 透明傳輸。
2.  **分包處理**: TCP 是串流協定。讀取回應時，必須持續 Read 直到讀到 `ETX` (0x03) 為止，才能視為一個完整的訊框。
3.  **Timeout**: 建議設定 Read Deadline (例如 1~2秒)，避免因斷線或 PLC 無回應導致永久阻塞。
4.  **併發**: FATEK PLC 的處理能力有限，建議對同一個 IP:Port 的連線使用 **Mutex** 鎖，確保「一問一答」順序執行，避免 Pipeline 請求導致 PLC 錯亂。

---

## 7. 測試案例 (Loopback Test)

使用 **Cmd 4E** (Loopback) 驗證通訊與 LRC 算法是否正確。

*   **Request**: `STX` + `"01"` + `"4E"` + `"AABB"` + `LRC` + `ETX`
*   **Expected Response**: `STX` + `"01"` + `"4E"` + `"AABB"` + `LRC` + `ETX`
*   若回應內容與請求完全一致，則底層通訊正常。