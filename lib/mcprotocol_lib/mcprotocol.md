# Mitsubishi MC Protocol (3E Binary) Implementation Plan

## 1. 協定概述 (Protocol Overview)

本模組將實作 Mitsubishi MELSEC Communication Protocol (MC Protocol) 的 **3E Frame (Binary Mode)**。
此協定廣泛用於與 Mitsubishi Q/L/iQ-R/iQ-F 系列 PLC 進行乙太網通訊。

*   **傳輸層**: TCP/IP
*   **編碼**: Binary (Little Endian)
*   **Frame**: 3E Frame (QnA 相容 3E 幀)

## 2. 功能範疇 (Scope)

參考 Python 原型 (`lib/mcprotocol_lib/mc_protocol.py`)，Go 版本將完整支援以下功能：

### 2.1 核心指令 (Commands)
| 功能 | 指令碼 (Command) | 子指令 (Subcmd) | 說明 |
| :--- | :--- | :--- | :--- |
| **Batch Read Word** | `0401` | `0000` | 批量讀取字組 (D, W, R...) |
| **Batch Write Word** | `1401` | `0000` | 批量寫入字組 |
| **Batch Read Bit** | `0401` | `0001` | 批量讀取位元 (M, X, Y...) |
| **Batch Write Bit** | `1401` | `0001` | 批量寫入位元 |
| **Random Read** | `0403` | `0000` | 隨機讀取字組 (支援多個不連續位址) |

### 2.2 支援元件 (Devices)
| 元件 | 代碼 (Hex) | 類型 | 說明 |
| :--- | :--- | :--- | :--- |
| **D** | `A8` | Word | 資料暫存器 |
| **W** | `B4` | Word | 連結暫存器 |
| **R** | `AF` | Word | 檔案暫存器 |
| **M** | `90` | Bit | 內部繼電器 |
| **X** | `9C` | Bit | 輸入 |
| **Y** | `9D` | Bit | 輸出 |
| **L** | `92` | Bit | 鎖存繼電器 |
| **B** | `A0` | Bit | 連結繼電器 |
| **F** | `93` | Bit | 報警器 |
| **TN**| `C2` | Word | 計時器(當前值) |
| **CN**| `C5` | Word | 計數器(當前值) |
| **TS**| `C1` | Bit | 計時器(觸點) |
| **CS**| `C4` | Bit | 計數器(觸點) |
| **SB**| `A1` | Bit | 特殊連結繼電器 |
| **SW**| `B5` | Word | 特殊連結暫存器 |

## 3. 架構設計 (Go Structure)

目錄: `internal/protocol/mcprotocol/`

### 3.1 檔案結構
*   `const.go`: 定義 Device Code, Command Code 常數。
*   `frame.go`: 處理 3E Frame 的 Header 封裝與解包、位元壓縮/解壓縮邏輯。
*   `client.go`: 主要 `MCClient` 結構，提供高階 API (`ReadD`, `WriteM` 等)。
*   `transport.go`: TCP 連線管理、超時重連機制 (參考 Fatek 的穩健設計)。

### 3.2 關鍵邏輯
1.  **位元壓縮 (Bit Packing)**:
    *   3E Binary 協定在讀寫 Bit 時，採用 **2 points / byte** 的壓縮方式。
    *   **Write**: 0x10 (ON/OFF), 0x11 (ON/ON), 0x00 (OFF/OFF)。
    *   **Read**: 回傳長度為 `(Count + 1) / 2` Bytes。需正確解壓 High/Low Nibble。
2.  **封包結構**:
    *   **Request**: `SubHeader(5000)` + `Net` + `PC` + `IO` + `Station` + `Len` + `Timer` + `Cmd` + `SubCmd` + `Data`
    *   **Response**: `SubHeader(D000)` + `Net` + `PC` + `IO` + `Station` + `Len` + `EndCode` + `Data`
3.  **錯誤處理**:
    *   檢查 Response Header 的 `EndCode` (非 0 代表 PLC 錯誤)。
    *   錯誤碼封裝為 `MCError` 提供詳細訊息。

## 4. 實作計畫

1.  建立基礎 `const.go` 定義所有代碼。
2.  實作 `transport.go` (可複用/參考 Fatek 的 `bufio` 讀取邏輯)。
3.  實作 `frame.go` 處理 Header 與 Bit Packing。
4.  實作 `client.go` 包含所有 Batch/Random 讀寫方法。
5.  撰寫 `mcprotocol_test.go` 進行單元測試 (Mock Server)。

---
*註：針對 Fatek 分析結果，確認專案採用的 Fatek ASCII 與提供的 Python Binary Lib 不同，但符合專案原始規劃。MC Protocol 則將完全依照 Python Lib 的規格 (3E Binary) 進行 Go 移植。*
