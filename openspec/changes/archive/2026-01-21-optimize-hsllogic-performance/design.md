# optimize-hsllogic-performance 設計文件

## 1. 架構概述

此設計旨在將 HslCommunication 的核心效能模式引入 `go-gateway` 專案，主要涉及三個層面：

```
┌─────────────────────────────────────────────────────────┐
│                    hsllogic Package                     │
│  ┌───────────────┐  ┌───────────────┐  ┌─────────────┐  │
│  │ AddressCache  │  │ PacketLogger  │  │ ByteTransform│ │
│  │  (新增)       │  │  (優化)       │  │  (現有)     │  │
│  └───────────────┘  └───────────────┘  └─────────────┘  │
└─────────────────────────────────────────────────────────┘
                          ▲
                          │
┌─────────────────────────────────────────────────────────┐
│                Protocol Connectors                      │
│  ┌───────────────┐  ┌───────────────┐  ┌─────────────┐  │
│  │ PersistentConn│  │ ConnectionPool│  │ Modbus/MC3E │  │
│  │  Interface    │  │   (新增)      │  │ Fatek...    │  │
│  │   (新增)      │  │               │  │             │  │
│  └───────────────┘  └───────────────┘  └─────────────┘  │
└─────────────────────────────────────────────────────────┘
```

## 2. 核心變更

### 2.1 PersistentConnection 介面（新增）

參考 HSL 的 `NetworkDoubleBase`，定義一個可切換長短連接模式的介面：

```go
// PersistentConnection 長連接支援介面
type PersistentConnection interface {
    // SetPersistentConnection 設定是否使用長連接模式
    SetPersistentConnection(enabled bool)
    // Connect 建立連接（長連接模式下只在首次調用）
    Connect() error
    // Close 關閉連接
    Close() error
    // IsConnected 檢查連接狀態
    IsConnected() bool
}
```

**影響範圍**：`internal/protocol/modbus`, `internal/protocol/fatek`, `internal/protocol/mc3e`

### 2.2 PacketLogger 物件池（優化）

利用 `sync.Pool` 減少 `PacketLog` 結構的分配：

```go
var packetLogPool = sync.Pool{
    New: func() interface{} {
        return &PacketLog{}
    },
}

func (pl *PacketLogger) log(...) {
    entry := packetLogPool.Get().(*PacketLog)
    defer packetLogPool.Put(entry)
    // ... 使用 entry ...
}
```

**注意**：需確保 `Put` 前清空 `entry` 的敏感欄位。

### 2.3 AddressCache（新增）

快取解析後的位址結構，避免熱迴圈中重複解析：

```go
type AddressCache struct {
    cache sync.Map // map[string]*ParsedAddress
}

func (c *AddressCache) Get(rawAddr string) (*ParsedAddress, error) {
    if cached, ok := c.cache.Load(rawAddr); ok {
        return cached.(*ParsedAddress), nil
    }
    addr, err := ParseAddress(rawAddr)
    if err != nil { return nil, err }
    c.cache.Store(rawAddr, addr)
    return addr, nil
}
```

### 2.4 ConnectionPool（新增）

參考 HSL 的 `ConnectPool<T>`，為多 Goroutine 存取同一設備提供連接池：

```go
type ConnectionPool struct {
    factory   func() (Connector, error)
    maxConns  int
    pool      chan Connector
}
```

## 3. 權衡取捨

| 方案       | 優點     | 缺點                       | 決策                     |
| ---------- | -------- | -------------------------- | ------------------------ |
| 全面長連接 | 極低延遲 | 需管理 KeepAlive, 資源佔用 | **採用**，作為預設       |
| sync.Pool  | 減少 GC  | 狀態管理複雜               | **採用**，限於 PacketLog |
| 連接池     | 高併發   | 增加複雜度                 | **延後**，Phase 2        |

## 4. 驗證策略

- **單元測試**：`PacketLogger` 的池化行為。
- **基準測試**：比較優化前後的分配次數與延遲。
- **整合測試**：長連接模式下的斷線重連。
