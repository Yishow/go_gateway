# p3-enable-virtual-device-mode 設計文件

## 1. 架構概述

```
┌─────────────────┐       ┌─────────────────────────┐       ┌─────────────────┐
│  Real PLC (A)   │ <──── │      go-gateway         │ <──── │   SCADA / HMI   │
└─────────────────┘       │ ┌─────────────────────┐ │       └─────────────────┘
                          │ │  Virtual Memory Map │ │          (Modbus Client)
┌─────────────────┐       │ │  [D0...D1000]       │ │
│  Real PLC (B)   │ <──── │ └──────────▲──────────┘ │
└─────────────────┘       │            │            │
      (Collection)        │            │(Serve)     │
                          │   ┌────────▼────────┐   │
                          │   │ ModbusTcpServer │   │
                          │   └─────────────────┘   │
                          └─────────────────────────┘
```

## 2. 核心組件

### 2.1 VirtualMemory (虛擬記憶體)

一個執行緒安全的記憶體儲存結構。

- `Read(address, length)`
- `Write(address, data)`
- 支援訂閱機制（當數值改變時觸發事件，用於 Datalink 反向寫入）。

### 2.2 ProtocolServer (服務器介面)

參考 HSL 的 Server 實作。

- `Start(port)`
- `Stop()`
- `HandleRequest(conn)`: 解析傳入的 Modbus/MC3E 封包，操作 `VirtualMemory` 並回傳。

### 2.3 SimulationManager

負責管理虛擬設備的生命週期與數據生成（如自動生成正弦波、隨機數）。

## 3. 數據流場景

### 場景 A：純模擬 (Simulation)

1. 啟動 `ModbusTcpServer` 監聽 502 port。
2. 配置 `SimulationRule`：每秒將 D100 + 1。
3. 外部 Client 連接 502 讀取 D100，看到數值跳動。

### 場景 B：協議轉換 (Gateway/Proxy)

1. `CollectionEngine` 從真實 Fatek PLC 讀取 D200。
2. 將數據寫入 `VirtualMemory` 的 D100 位置。
3. 外部 SCADA 透過 Modbus TCP 讀取 Gateway 的 D100，間接取得 Fatek 數據。

## 4. 權衡取捨

| 選擇       | 方案                                               | 決策                                                         |
| ---------- | -------------------------------------------------- | ------------------------------------------------------------ |
| 記憶體實作 | `map[string]interface{}` (靈活) vs `[]byte` (高效) | 採用 **`[]byte`** 配合 Hsl 地址映射算法，以模擬真實 PLC 行為 |
| 協議支援   | 支援所有 vs 僅 Modbus                              | Phase 1 **僅支援 Modbus TCP Server** (最通用)                |
