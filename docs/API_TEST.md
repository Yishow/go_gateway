# 測試 API 文件

本文件描述後端測試 API（`/api/v1/test`）與相關端點的請求格式、欄位與範例。

## 基礎路徑

- 預設：`/api/v1`

## 協議與操作類型

### 協議代碼（`protocol`）

- `modbus_tcp`
- `modbus_udp`
- `modbus_rtu`
- `fatek_tcp`
- `fatek_serial`
- `mc_tcp`
- `mc_serial`

### 讀取操作（`operation` for `/test/read`）

**Modbus**

- `read_coils`
- `read_discrete_inputs`
- `read_holding_registers`
- `read_input_registers`

**Fatek**

- `read_status`
- `read_registers`

**MC Protocol**

- `batch_read_word`
- `batch_read_bit`

### 寫入操作（`operation` for `/test/write`）

**Modbus**

- `write_single_coil`
- `write_single_register`
- `write_multiple_coils`
- `write_multiple_registers`

**Fatek**

- `write_status`
- `write_registers`

**MC Protocol**

- `batch_write_word`
- `batch_write_bit`

### 補充限制

- **MC Protocol**：`device` 名稱大小寫不敏感（例如 `d` 等同 `D`）。
- **Fatek**：`symbol` 支援 `DR/DD/DF` 與 `WX/WY/WM/WS/WT/WC` 等擴充符號。

## API 端點

### 1) 建立連線

`POST /api/v1/test/connect`

**Request Body**

```json
{
  "protocol": "fatek_tcp",
  "config": {
    "host": "192.168.1.5",
    "port": 500,
    "station": 1,
    "timeout": 2000
  }
}
```

**Response**

```json
{
  "connection_id": "conn_...",
  "status": "connected"
}
```

---

### 2) 斷開連線

`POST /api/v1/test/disconnect?connection_id=...`

**Response**

```json
{
  "status": "disconnected"
}
```

---

### 3) 取得連線狀態

`GET /api/v1/test/status?connection_id=...`

**Response**

```json
{
  "id": "conn_...",
  "protocol": "fatek_tcp",
  "config": {
    "host": "192.168.1.5",
    "port": 500,
    "station": 1,
    "timeout": 2000
  },
  "connected": true,
  "created_at": "2026-01-18T00:00:00Z"
}
```

---

### 4) 讀取

`POST /api/v1/test/read`

**Modbus 範例**

```json
{
  "connection_id": "conn_...",
  "operation": "read_holding_registers",
  "address": 0,
  "count": 10
}
```

**Fatek 範例**

```json
{
  "connection_id": "conn_...",
  "operation": "read_registers",
  "symbol": "DD",
  "address": 0,
  "count": 10
}
```

**MC Protocol 範例（device 大小寫不敏感）**

```json
{
  "connection_id": "conn_...",
  "operation": "batch_read_word",
  "device": "d",
  "address": 0,
  "count": 10
}
```

**Response**

```json
{
  "values": [100, 200, 300],
  "count": 3
}
```

---

### 5) 寫入

`POST /api/v1/test/write`

**Modbus 範例（多線圈）**

```json
{
  "connection_id": "conn_...",
  "operation": "write_multiple_coils",
  "address": 0,
  "values": [true, false, true]
}
```

**Fatek 範例**

```json
{
  "connection_id": "conn_...",
  "operation": "write_registers",
  "symbol": "WX",
  "address": 0,
  "values": [100, 200, 300]
}
```

**MC Protocol 範例**

```json
{
  "connection_id": "conn_...",
  "operation": "batch_write_bit",
  "device": "M",
  "address": 0,
  "values": [true, false, true]
}
```

**Response**

```json
{
  "status": "success"
}
```

---

### 6) 批量測試

`POST /api/v1/test/batch`

```json
{
  "connection_id": "conn_...",
  "operations": [
    {
      "type": "read",
      "read_request": {
        "operation": "read_holding_registers",
        "address": 0,
        "count": 10
      }
    },
    {
      "type": "write",
      "write_request": {
        "operation": "write_multiple_registers",
        "address": 0,
        "values": [1, 2, 3]
      }
    }
  ]
}
```

**Response**

```json
{
  "results": [{ "success": true, "data": [1, 2, 3] }, { "success": true }]
}
```

---

### 7) 監控模式

`POST /api/v1/test/monitor/start`

```json
{
  "connection_id": "conn_...",
  "interval": 500,
  "items": [
    {
      "operation": "read_holding_registers",
      "address": 0,
      "count": 10
    }
  ]
}
```

`POST /api/v1/test/monitor/stop`

```json
{
  "connection_id": "conn_..."
}
```

---

## 連線配置（`config`）

### Modbus TCP/UDP

```json
{
  "host": "127.0.0.1",
  "port": 502,
  "unitID": 1,
  "timeout": 2000
}
```

### Modbus RTU

```json
{
  "port": "COM4",
  "baudRate": 9600,
  "dataBits": 8,
  "stopBits": 1,
  "parity": "N",
  "unitID": 1,
  "timeout": 2000
}
```

### Fatek TCP

```json
{
  "host": "127.0.0.1",
  "port": 500,
  "station": 1,
  "timeout": 2000
}
```

### Fatek Serial

```json
{
  "port": "COM4",
  "baudRate": 9600,
  "dataBits": 7,
  "stopBits": 2,
  "parity": "E",
  "station": 1,
  "timeout": 2000
}
```

### MC Protocol TCP

```json
{
  "host": "127.0.0.1",
  "port": 5000,
  "timeout": 2000
}
```

### MC Protocol Serial

```json
{
  "port": "COM4",
  "baudRate": 9600,
  "dataBits": 7,
  "stopBits": 2,
  "parity": "E",
  "timeout": 2000
}
```
