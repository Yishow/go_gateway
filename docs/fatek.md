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

**Go 實作範例**:
```go
func CalculateLRC(data []byte) (string, error) {
    // data 應包含 STX (0x02) 開始，直到 Body 結束
    var sum byte = 0
    for _, b := range data {
        sum += b
    }
    // 轉為大寫 Hex 字串，例如 0xC7 -> "C7"
    return fmt.Sprintf("%02X", sum), nil
}
```

**驗證案例 (Cmd 40)**:
*   數據: `STX` + `"01"` + `"40"`
*   Hex: `02` + `30` + `31` + `34` + `30`
*   Sum: `0xC7`
*   Checksum String: `"C7"`

---

## 3. 元件位址編碼 (Addressing)

FATEK 協定對不同類型的元件有固定的字串長度格式。

| 元件類型 | 代號 | 範圍 | 格式長度 | 範例字串 | 說明 |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Discrete (單點)** | X | 0000~9999 | 5 | `"X0000"` | 輸入接點 |
| | Y | 0000~9999 | 5 | `"Y0000"` | 輸出繼電器 |
| | M | 0000~9999 | 5 | `"M0000"` | 內部繼電器 |
| | S | 0000~9999 | 5 | `"S0000"` | 步進繼電器 |
| | T | 0000~9999 | 5 | `"T0000"` | 計時器 (狀態) |
| | C | 0000~9999 | 5 | `"C0000"` | 計數器 (狀態) |
| **Register (16-bit)** | R | 00000~4167 | 6 | `"R00000"`| 資料暫存器 |
| | D | 00000~4999 | 6 | `"D00000"`| 資料暫存器 |
| | RT | 0000~9999 | 6 | `"RT0000"`| 計時器 (數值) |
| | RC | 0000~9999 | 6 | `"RC0000"`| 計數器 (數值) |
| **Register (32-bit)** | DR | 00000~65534| 7 | `"D00000"`| (D暫存器雙字) |
| | DW | | 7 | `"DWX000"`| (離散組雙字) |

*注意: 暫存器位址字串需補零至固定長度。例如 R10 必須寫為 `"R00010"` (6碼)。*

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
