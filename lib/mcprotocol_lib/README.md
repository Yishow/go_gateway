# MC Protocol Library (Python)

這是一個輕量級、無需依賴外部套件的 Mitsubishi MELSEC MC Protocol (3E Binary Frame) Python 實作。
專為工業自動化應用設計，提供穩定、易用的 API 來讀寫三菱 PLC。

## 特色
- **純 Python 實作**: 僅使用標準庫 (`socket`, `struct`)，無需 `pip install`。
- **完整支援 3E Binary**: 支援最常用的乙太網通訊幀格式。
- **隨機讀取支援**: 支援單次請求讀取多個不連續地址 (`random_read_word`)。
- **連線檢測**: 提供 `check_connection` 方法即時確認連線狀態。
- **類型安全**: 提供完整的 Type Hinting。
- **中文文件**: 程式碼與文件皆使用繁體中文。
- **自動重連**: 內建斷線偵測與重連機制。

## 安裝
將 `mcprotocol_lib` 資料夾複製到您的專案中即可使用。

## 快速開始

```python
from mcprotocol_lib.mc_protocol import MCProtocol

# 1. 初始化連線
# host: PLC IP
# port: MC Protocol Port (需在 PLC 參數中設定，通常為 5000 或 5002)
plc = MCProtocol(host='192.168.1.5', port=5000)

try:
    plc.connect()
    
    # 檢查連線狀態
    if plc.check_connection():
        print("PLC 連線正常")

    # 2. 寫入數據
    # 寫入 D100, D101, D102 為 123, 456, 789
    plc.write_d(100, [123, 456, 789])
    
    # 寫入 M0 為 ON
    plc.write_m(0, 1)

    # 3. 讀取數據 (批量)
    # 讀取 D100 開始的 5 個 Word
    d_values = plc.read_d(100, 5)
    print(f"D100-D104: {d_values}")
    
    # 讀取 M0 開始的 10 個 Bit
    m_values = plc.read_m(0, 10)
    print(f"M0-M9: {m_values}")

    # 4. 隨機讀取 (讀取不連續地址)
    # 同時讀取 D100, D200, W50
    random_values = plc.random_read_word([('D', 100), ('D', 200), ('W', 50)])
    print(f"Random Read Values: {random_values}")

except Exception as e:
    print(f"發生錯誤: {e}")
finally:
    plc.close()
```

## API 參考

### `MCProtocol(host, port, network_no=0, pc_no=255, timeout=2.0)`
建構子。
- `network_no`: 網路編號 (預設 0)
- `pc_no`: PC 編號 (預設 0xFF)

### 連線管理
- `connect()`: 建立連線。
- `close()`: 關閉連線。
- `check_connection() -> bool`: 檢查目前連線是否有效 (發送輕量級請求)。

### 讀寫方法

#### Word (字組) 存取
- `batch_read_word(device_type, head_device, count)`: 通用批量讀取
- `batch_write_word(device_type, head_device, values)`: 通用批量寫入
- `random_read_word(word_devices)`: 隨機讀取不連續字組
    - `word_devices`: Tuple 列表，例如 `[('D', 100), ('W', 50)]`
- `read_d(address, count)`: 讀取 D 暫存器 (捷徑)
- `write_d(address, values)`: 寫入 D 暫存器 (捷徑)

#### Bit (位元) 存取
- `batch_read_bit(device_type, head_device, count)`: 通用批量讀取
- `batch_write_bit(device_type, head_device, values)`: 通用批量寫入
- `read_m(address, count)`: 讀取 M 繼電器 (捷徑)
- `write_m(address, values)`: 寫入 M 繼電器 (捷徑)

### 支援的軟元件 (Device Code)

| 代號 | 類型 | 說明 |
|------|------|------|
| D | Word | 資料暫存器 |
| W | Word | 連結暫存器 |
| R | Word | 檔案暫存器 |
| M | Bit | 內部繼電器 |
| X | Bit | 輸入 |
| Y | Bit | 輸出 |
| L | Bit | 鎖存繼電器 |
| B | Bit | 連結繼電器 |
| F | Bit | 報警器 |
| TN | Word | 計時器(目前值) |
| CN | Word | 計數器(目前值) |
| ... | ... | ... |

(完整列表請參考程式碼中的 `DEVICE_MAP`)

## 錯誤處理

- `MCConnectionError`: 連線失敗或通訊中斷。
- `MCResponseError`: PLC 回傳錯誤代碼 (例如地址超出範圍)。
- `MCProtocolError`: 封包格式錯誤或其他協定層級錯誤。