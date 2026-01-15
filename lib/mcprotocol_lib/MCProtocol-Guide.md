# MC Protocol 完整實現指南
# Complete MC Protocol Implementation Guide

**版本 / Version: 3.0**  
**支持 / Supports: 3E/4E Frame, Binary Encoding, Auto-batching, Error Handling, Logging**

---

## 📋 目錄 / Table of Contents

1. [快速開始 / Quick Start](#快速開始--quick-start)
2. [協議概述 / Protocol Overview](#協議概述--protocol-overview)
3. [幀結構 / Frame Structure](#幀結構--frame-structure)
4. [API 文件 / API Documentation](#api-文件--api-documentation)
5. [設備類型 / Device Types](#設備類型--device-types)
6. [錯誤處理 / Error Handling](#錯誤處理--error-handling)
7. [批次管理 / Batch Management](#批次管理--batch-management)
8. [範例代碼 / Code Examples](#範例代碼--code-examples)

---

## 快速開始 / Quick Start

### 安裝 / Installation

```bash
# 無需外部依賴 / No external dependencies required
python3 mcprotocol_complete.py
```

### 基本使用 / Basic Usage

```python
from mcprotocol_complete import MCProtocol, MCConfig, DeviceType, FrameType

# 配置 / Configure
config = MCConfig(
    host='192.168.0.2',
    port=2000,
    frame_type=FrameType.FRAME_3E
)

# 連接 / Connect
with MCProtocol(config) as plc:
    # 讀取 / Read
    data = plc.read_batch(DeviceType.D, 0, 100)
    
    # 寫入 / Write
    plc.write_batch(DeviceType.M, 0, [1, 0, 1])
    
    # 讀取位元 / Read bits
    bits = plc.read_batch(DeviceType.M, 0, 16)
```

---

## 協議概述 / Protocol Overview

### MC Protocol 基礎 / MC Protocol Basics

**MC Protocol** (MELSEC 通訊協議) 用於與三菱 PLC 進行 Ethernet 通訊

**特性 / Features:**
- ✅ **3E Frame**: 標準 SLMP 幀格式 / Standard SLMP frame format
- ✅ **4E Frame**: 帶序列號的增強幀 / Enhanced frame with serial number
- ✅ **Binary Encoding**: 高效率二進制編碼 / Efficient binary encoding
- ✅ **Auto-batching**: 自動分包超大請求 / Automatically split large requests
- ✅ **Error Handling**: 完整的異常和錯誤碼處理 / Complete exception and error code handling
- ✅ **Checksum**: 校驗和驗證 / Checksum verification
- ✅ **Logging**: 詳細的通訊日誌 / Detailed communication logging

### 通訊流程 / Communication Flow

```
PC/外設                          PLC
 │                              │
 ├─── 編碼請求 ─── TCP ────────→ │ 接收請求
 │    (Request Frame)            │ (Parse Request)
 │                               │
 │                      執行命令 │
 │                      (Execute)│
 │                               │
 │ ← TCP ─── 編碼回應 ──────────┤ 發送回應
 │  (Response Frame)              │ (Send Response)
 │
 └─ 解析回應、驗證校驗和
   (Parse Response, Verify Checksum)
```

---

## 幀結構 / Frame Structure

### 3E Frame (Binary Mode)

**請求 / Request Frame:**

```
┌─────────────────────────────────────────────────────────────┐
│ Header (2B) │ Access Route (7B) │ Command │ Data │ Checksum │
├─────────────┼──────────────────┼─────────┼──────┼──────────┤
│  50 00      │ 00 FF 03FF 00    │  0401   │ ...  │    XX    │
│ SubHeader   │ Nw PC IOMod Stat │ Cmd Sub │      │  (opt)   │
└─────────────────────────────────────────────────────────────┘
```

**回應 / Response Frame:**

```
┌──────────────────────────────┬─────────────┬──────────┐
│ Header (2B) │ Access (7B) │ End Code (2B) │ Data │ CS │
├─────────────┼─────────────┼───────────────┼──────┼────┤
│  D0 00      │ 00 FF ... │    00 00       │ ...  │ XX │
│ SubHeader   │ Route      │    (NORMAL)    │      │    │
└──────────────────────────────┴─────────────┴──────┴────┘
```

### 4E Frame (Binary Mode)

**請求 / Request Frame:**

```
┌────────────────────────────────────────────────────┐
│ Header (4B) │ Access (7B) │ Command │ Data │ CS   │
├──────────────┼─────────────┼────────┼─────┼──────┤
│ 54 00 + SN  │ 00 FF ...  │  0401  │ ... │  XX  │
│ Sub + Serial │ Route       │ Cmd Sub│     │      │
└────────────────────────────────────────────────────┘
```

---

## API 文件 / API Documentation

### MCConfig 類 / MCConfig Class

```python
@dataclass
class MCConfig:
    host: str                                       # PLC IP
    port: int = 2000                              # TCP Port
    frame_type: FrameType = FrameType.FRAME_3E   # 3E or 4E
    plc_model: PLCModel = PLCModel.MELSEC_Q      # PLC Model
    encoding: str = "binary"                      # binary/ascii
    network_number: int = 0x00                    # Network No.
    pc_number: int = 0xFF                         # PC No.
    timeout: float = 5.0                          # Socket timeout
    serial_number: int = 0x0000                   # For 4E frame
    sum_check: bool = True                        # Enable checksum
```

### MCProtocol 類 / MCProtocol Class

#### 連接管理 / Connection Management

```python
# 連接 / Connect
plc = MCProtocol(config)
plc.connect()  # → bool

# 斷開 / Disconnect
plc.disconnect()  # → bool

# Context manager
with MCProtocol(config) as plc:
    # 自動連接/斷開 / Auto connect/disconnect
    ...
```

#### 讀取操作 / Read Operations

```python
# 批次讀取 (自動分包) / Batch read (auto-batching)
values = plc.read_batch(
    device_type=DeviceType.D,     # 設備類型
    start_address=0,               # 起始位址
    count=100                       # 讀取點數
)  # → List[int]

# 隨機讀取 (混合設備) / Random read (mixed devices)
results = plc.read_random({
    "temperature": (DeviceType.D, 0),
    "pressure": (DeviceType.D, 1),
    "status": (DeviceType.M, 10)
})  # → Dict[str, int]

# 隨機讀取 (單一設備類型) / Random read (single device type)
values_dict = plc.read_random_batch(
    device_type=DeviceType.D,
    addresses=[0, 5, 10, 50, 100]
)  # → Dict[int, int]
```

#### 寫入操作 / Write Operations

```python
# 批次寫入 (自動分包) / Batch write (auto-batching)
success = plc.write_batch(
    device_type=DeviceType.M,
    start_address=0,
    values=[1, 0, 1, 0, 1]
)  # → bool
```

#### 工具方法 / Utility Methods

```python
# 讀取 PLC 型號 / Read PLC model
model = plc.get_plc_model()  # → str

# 迴路測試 / Loopback test
success = plc.device_loopback_test()  # → bool
```

### MCProtocolSafe 類 / MCProtocolSafe Class

具有自動重試的安全客戶端

```python
plc = MCProtocolSafe(config, max_retries=3)

# 自動重試的讀取 / Read with auto-retry
data = plc.read_batch_safe(DeviceType.D, 0, 100)

# 自動重試的寫入 / Write with auto-retry
success = plc.write_batch_safe(DeviceType.M, 0, [1, 0, 1])
```

---

## 設備類型 / Device Types

### 支援的設備 / Supported Devices

| 設備 | 代碼 | 位寬 | 說明 |
|------|------|------|------|
| **D** | 0xA8 | 16-bit | 數據暫存器 / Data register |
| **M** | 0x90 | 1-bit | 內部繼電器 / Internal relay |
| **X** | 0x9C | 1-bit | 輸入 / Input |
| **Y** | 0xA8 | 1-bit | 輸出 / Output |
| **S** | 0x98 | 1-bit | 狀態繼電器 / State relay |
| **T** | 0xC0 | 16-bit | 計時器 / Timer |
| **C** | 0xC4 | 16-bit | 計數器 / Counter |
| **F** | 0xAF | 1-bit | 文件暫存器位 / File register bit |
| **B** | 0xA0 | 1-bit | 連結繼電器 / Link relay |
| **W** | 0xB4 | 16-bit | 連結暫存器 / Link register |
| **R** | 0xAF | 16-bit | 文件暫存器 / File register |
| **TN** | 0xC1 | 16-bit | 計時器當前值 / Timer current value |
| **CN** | 0xC5 | 16-bit | 計數器當前值 / Counter current value |
| **TS** | 0xC7 | 1-bit | 計時器狀態 / Timer status |
| **CS** | 0xC8 | 1-bit | 計數器狀態 / Counter status |

### 使用設備類型 / Using Device Types

```python
# 字設備 (16-bit) / Word devices
data = plc.read_batch(DeviceType.D, 100, 10)  # D100-D109

# 位設備 (1-bit) / Bit devices
bits = plc.read_batch(DeviceType.M, 0, 16)    # M0-M15
```

---

## 錯誤處理 / Error Handling

### 異常類型 / Exception Types

```python
MCProtocolException              # 基礎異常 / Base exception
├── ConnectionException          # 連線錯誤 / Connection error
├── TimeoutException            # 超時 / Timeout
├── ChecksumException           # 校驗和錯誤 / Checksum error
├── EndCodeException            # 結束碼錯誤 / End code error
├── BatchSizeException          # 批次大小超限 / Batch size exceeded
└── DeviceException             # 設備錯誤 / Device error
```

### 結束碼 / End Codes

| 碼 | 值 | 說明 |
|----|-------|------|
| NORMAL | 0x0000 | 正常完成 / Normal completion |
| DEVICE_ERROR | 0xC051 | 設備錯誤 / Device error |
| BLOCK_ERROR | 0xC052 | 區塊錯誤 / Block error |
| CHECKSUM_ERROR | 0xC056 | 校驗和錯誤 / Checksum error |
| FORMAT_ERROR | 0xC057 | 格式錯誤 / Format error |
| UNKNOWN_ERROR | 0xC0FF | 未知錯誤 / Unknown error |

### 錯誤處理範例 / Error Handling Examples

```python
try:
    with MCProtocol(config) as plc:
        data = plc.read_batch(DeviceType.D, 0, 100)
        
except ConnectionException:
    print("無法連接 PLC")  # Cannot connect to PLC
    
except TimeoutException:
    print("通訊超時")  # Communication timeout
    
except EndCodeException as e:
    print(f"PLC 錯誤: 0x{e.end_code:04X}")  # PLC error
    
except ChecksumException:
    print("資料校驗失敗")  # Checksum verification failed
    
except MCProtocolException as e:
    print(f"MC Protocol 錯誤: {e}")
```

---

## 批次管理 / Batch Management

### 自動分包 / Auto-batching

當請求超過最大點數限制時，自動分割成多個請求

**最大點數限制 / Max Points per Request:**
- **字設備 (Word)**: 256 個點 / 256 points
- **位設備 (Bit)**: 256 個點 / 256 points

### 分包範例 / Batching Example

```python
# 單一請求會自動分成多個 / Single request auto-splits
data = plc.read_batch(DeviceType.D, 0, 500)
# 會執行:
# - Request 1: D0-D255 (256 points)
# - Request 2: D256-D500 (244 points)

logger.info("Read 500 points from D0 - auto-batching enabled")
```

### 手動批次控制 / Manual Batch Control

```python
from mcprotocol_complete import BatchManager

# 分析批次 / Analyze batches
batches = BatchManager.split_requests(
    DeviceType.D,
    addresses=[0, 100, 200, 300, 400],
    is_write=False
)
# [(0, 401), ...]  # 起始位址, 點數 / start_address, count
```

---

## 範例代碼 / Code Examples

### 1. 基本讀寫 / Basic Read/Write

```python
from mcprotocol_complete import MCProtocol, MCConfig, DeviceType, FrameType

config = MCConfig(
    host='192.168.0.2',
    port=2000,
    frame_type=FrameType.FRAME_3E
)

with MCProtocol(config) as plc:
    # 讀取數據暫存器 D0-D9 / Read data registers D0-D9
    data = plc.read_batch(DeviceType.D, 0, 10)
    print(f"D0-D9: {data}")
    
    # 寫入 M0-M4 / Write to M0-M4
    plc.write_batch(DeviceType.M, 0, [1, 0, 1, 1, 0])
    print("Written to M0-M4")
    
    # 讀取位元 / Read bits
    bits = plc.read_batch(DeviceType.M, 0, 16)
    print(f"M0-M15: {bits}")
```

### 2. 4E Frame 帶序列號 / 4E Frame with Serial Number

```python
config = MCConfig(
    host='192.168.0.2',
    frame_type=FrameType.FRAME_4E,
    serial_number=0x1234  # 序列號 / Serial number
)

with MCProtocol(config) as plc:
    data = plc.read_batch(DeviceType.D, 0, 100)
```

### 3. 自動重試 / Auto-retry

```python
from mcprotocol_complete import MCProtocolSafe

config = MCConfig(
    host='192.168.0.2',
    timeout=5.0
)

# 最多重試 3 次 / Max 3 retries
plc = MCProtocolSafe(config, max_retries=3)

try:
    data = plc.read_batch_safe(DeviceType.D, 0, 100)
except Exception as e:
    print(f"Failed after 3 retries: {e}")
```

### 4. 混合設備隨機讀取 / Random Read Mixed Devices

```python
results = plc.read_random({
    "temperature": (DeviceType.D, 0),
    "pressure": (DeviceType.D, 1),
    "pump_status": (DeviceType.M, 10),
    "valve_position": (DeviceType.D, 100)
})

print(f"Temperature: {results['temperature']}")
print(f"Pump Status: {results['pump_status']}")
```

### 5. 大量資料讀取 (自動分包) / Large Data Read (Auto-batching)

```python
# 讀取 1000 個點，自動分成 4 個請求
# Read 1000 points, auto-split into 4 requests
data = plc.read_batch(DeviceType.D, 0, 1000)
print(f"Read {len(data)} values")

# 寫入 500 個值，自動分成 2 個請求
# Write 500 values, auto-split into 2 requests
plc.write_batch(DeviceType.M, 0, [i % 2 for i in range(500)])
```

### 6. 錯誤處理和日誌 / Error Handling and Logging

```python
import logging
from mcprotocol_complete import MCProtocol, MCProtocolException

# 設置日誌級別 / Set logging level
logging.getLogger('mcprotocol').setLevel(logging.DEBUG)

try:
    with MCProtocol(config) as plc:
        # 檢查 PLC 連接 / Check PLC connection
        if plc.device_loopback_test():
            print("✓ PLC Connected")
        else:
            print("✗ PLC Not responding")
            
        # 讀取 PLC 型號 / Read PLC model
        model = plc.get_plc_model()
        print(f"PLC Model: {model}")
        
except MCProtocolException as e:
    print(f"Error: {e}")
```

### 7. 數據打包和解析 / Data Packing and Parsing

```python
from mcprotocol_complete import DataParser

# 字資料打包 / Pack word data
word_values = [100, 200, 300, 400]
packed = DataParser.pack_word_data(word_values)

# 位資料打包 / Pack bit data (16 bits per word)
bit_values = [1, 0, 1, 1, 0, 1, 0, 1, 1, 1, 0, 0, 1, 0, 1, 1]
packed_bits = DataParser.pack_bit_data(bit_values)

# 解析回應 / Parse response
word_response = bytes([0x64, 0x00, 0xC8, 0x00, 0x2C, 0x01])
parsed = DataParser.parse_word_response(word_response)
print(f"Parsed words: {parsed}")

# 解析位回應 / Parse bit response
bit_response = bytes([0xFF, 0x00, 0x55, 0xAA])
parsed_bits = DataParser.parse_bit_response(bit_response, 32)
print(f"Parsed bits: {parsed_bits}")
```

### 8. 監測 PLC 狀態 / Monitor PLC Status

```python
import time

config = MCConfig(host='192.168.0.2')

with MCProtocol(config) as plc:
    while True:
        try:
            # 監測運行位 / Monitor run bit
            run_bits = plc.read_batch(DeviceType.M, 0, 1)
            
            # 監測計數器 / Monitor counters
            counter = plc.read_batch(DeviceType.CN, 0, 1)
            
            # 監測計時器 / Monitor timer
            timer = plc.read_batch(DeviceType.TN, 0, 1)
            
            print(f"Run: {run_bits[0]}, Counter: {counter[0]}, Timer: {timer[0]}")
            
            time.sleep(1)
            
        except Exception as e:
            print(f"Monitor error: {e}")
            break
```

---

## 性能優化建議 / Performance Optimization

### 1. 使用 Binary 編碼 / Use Binary Encoding

```python
# ✓ 推薦 / Recommended (50% 更小的資料)
config = MCConfig(encoding='binary')  # 預設 / Default

# ✗ 避免 / Avoid (2 倍的資料大小)
config = MCConfig(encoding='ascii')
```

### 2. 合併請求 / Merge Requests

```python
# ✗ 不推薦 / Not recommended (3 個請求)
v1 = plc.read_batch(DeviceType.D, 0, 10)
v2 = plc.read_batch(DeviceType.D, 50, 10)
v3 = plc.read_batch(DeviceType.D, 100, 10)

# ✓ 推薦 / Recommended (1 個請求)
v = plc.read_batch(DeviceType.D, 0, 110)
```

### 3. 使用 Context Manager / Use Context Manager

```python
# ✓ 推薦 / Recommended (自動連接/斷開)
with MCProtocol(config) as plc:
    data = plc.read_batch(DeviceType.D, 0, 100)

# ✗ 手動管理 / Manual management
plc = MCProtocol(config)
plc.connect()
data = plc.read_batch(DeviceType.D, 0, 100)
plc.disconnect()
```

### 4. 批次延遲 / Batch Delay

```python
# 自動分包時批次間延遲 / Delay between batches
time.sleep(0.1)

# 可在批次讀寫時設置 / Can be set in batch operations
```

---

## 生產環境檢查清單 / Production Environment Checklist

- [ ] **Network**: PLC 與 PC 在同一網段 / Same subnet
- [ ] **Firewall**: Port 2000 已開放 / Port 2000 opened
- [ ] **PLC Setup**: 
  - [ ] Ethernet 模組已配置 / Ethernet module configured
  - [ ] MC Protocol 已啟用 / MC Protocol enabled
  - [ ] Network Number & PC Number 已設置 / Set correctly
- [ ] **Code**:
  - [ ] Error handling 完整 / Complete
  - [ ] Logging 已啟用 / Enabled
  - [ ] Timeout 值適當 / Appropriate value
  - [ ] Auto-retry 已配置 / Configured
- [ ] **Testing**:
  - [ ] 迴路測試通過 / Loopback test passed
  - [ ] PLC 型號讀取成功 / Model read successfully
  - [ ] 小數據量讀寫測試 / Small data test passed
  - [ ] 大數據量自動分包測試 / Auto-batching test passed

---

## 常見問題 / FAQ

### Q1: 無法連接到 PLC
### Q1: Cannot connect to PLC

**解決方案 / Solutions:**
1. 檢查 IP 和 Port / Check IP and Port
2. 檢查網路連接 / Check network connectivity
3. 檢查 PLC Ethernet 模組配置 / Check PLC Ethernet module setup
4. 檢查防火牆 / Check firewall

### Q2: 校驗和錯誤
### Q2: Checksum error

**解決方案 / Solutions:**
1. 確認 sum_check=True / Ensure sum_check=True
2. 檢查網路穩定性 / Check network stability
3. 嘗試禁用校驗和測試 / Try disabling checksum to test

### Q3: 超時錯誤
### Q3: Timeout error

**解決方案 / Solutions:**
1. 增加 timeout 值 / Increase timeout value
2. 檢查 PLC 負載 / Check PLC load
3. 檢查網路延遲 / Check network latency

### Q4: 批次大小超限
### Q4: Batch size exceeded

**解決方案 / Solutions:**
- 庫會自動分包 / Library automatically batches
- 檢查日誌確認分包 / Check logs for batching

---

## 版本歷史 / Version History

### v3.0 (Current)
- ✅ 完整 3E/4E Frame 支持 / Full 3E/4E Frame support
- ✅ Binary 編碼 / Binary encoding
- ✅ 自動分包 / Auto-batching
- ✅ 位/字資料解析與打包 / Bit/word data parsing and packing
- ✅ 結束碼與狀態碼解析 / End code and status code parsing
- ✅ 異常類別與錯誤處理 / Exception classes and error handling
- ✅ 日誌記錄 / Logging
- ✅ 自動重試 / Auto-retry

### v2.0
- 基本 3E Frame 支持
- ASCII 編碼

### v1.0
- 初始版本

---

## 許可 / License

本庫可自由使用於商業和個人項目

---

## 支持 / Support

詳見代碼註釋和日誌輸出

**Log file**: `mcprotocol.log`

---

**Made with ❤️ for MELSEC PLC Communication**
