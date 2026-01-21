# Persistent Connection Mode

本文件說明協議連接器（Protocol Connectors）的長連接（Persistent Connection）與短連接（Short Connection）模式，以及如何在應用程式中使用它們。

## 概述

為了優化高頻率輪詢（Polling）場景下的效能，系統在協議連接器層級引入了「雙模式」連線策略：

1.  **長連接模式 (Persistent Mode)**：預設模式。連線建立後會持續保持，供後續請求重複使用。這消除了每次請求都要進行 TCP 三向交握（Handshake）的開銷，顯著降低延遲並減少設備負擔。
2.  **短連接模式 (Short Connection Mode)**：適用於低頻率或間歇性通訊。每次請求都會建立新的連線，並在操作完成後立即關閉。這有助於釋放網路資源，避免長期佔用設備連線數。

## 支援的協議

目前以下協議連接器支援雙模式切換：

*   **Modbus TCP** (`modbus_tcp`)
*   **Modbus RTU** (`modbus_rtu`)
*   **Modbus UDP** (`modbus_udp`)
*   **FATEK FBs** (`fatek`)
*   **Mitsubishi MC 3E** (`mc3e`)

## 介面定義

所有支援雙模式的連接器都實作了 `connector.PersistentConnection` 介面：

```go
type PersistentConnection interface {
    Protocol

    // SetPersistentConnection 設定是否使用長連接模式
    // enabled: true 啟用長連接，false 使用短連接（每次操作後斷線）
    SetPersistentConnection(enabled bool)

    // IsPersistentMode 檢查當前是否為長連接模式
    IsPersistentMode() bool

    // Disconnect 顯式斷線（僅在長連接模式下有意義）
    // 短連接模式下此方法等同於 Close()
    Disconnect() error

    // Reconnect 重新連線（用於斷線恢復）
    Reconnect(ctx context.Context) error
}
```

## 使用行為

### 預設行為

當透過 `connector.Get()` 或 `connector.New*Connector()` 建立連接器實例，並呼叫 `Connect()` 時，**系統預設啟用長連接模式**。

這是因為在工業自動化場景中，持續監控是常態，長連接能提供最佳效能。

### 切換模式

開發者可以透過型別斷言（Type Assertion）將 `Protocol` 介面轉換為 `PersistentConnection` 介面，並動態切換模式。

#### 範例：切換至短連接模式

```go
conn, err := connector.Get("modbus_tcp", config)
if err != nil {
    return err
}
defer conn.Close()

// 檢查是否支援長連接介面
if pc, ok := conn.(connector.PersistentConnection); ok {
    // 切換為短連接模式
    // 之後的每次 Read/Write 都會自動執行 Connect -> Operation -> Close
    pc.SetPersistentConnection(false)
}

// 執行讀取（內部自動處理連線生命週期）
result, err := conn.Read(ctx, request)
```

#### 範例：顯式斷線與重連

在長連接模式下，有時需要強制重置連線（例如設備狀態異常時）：

```go
if pc, ok := conn.(connector.PersistentConnection); ok {
    // 強制斷線並重新建立連線
    if err := pc.Reconnect(ctx); err != nil {
        log.Printf("Reconnect failed: %v", err)
    }
}
```

## 錯誤處理與自動恢復

在長連接模式下，連接器內建了基本的斷線偵測與自動重連機制：

1.  當 `Read` 或 `Write` 操作因網路錯誤失敗時，連接器會標記連線為「不健康」。
2.  下一次操作嘗試時，系統會檢測到連線已中斷，並嘗試自動執行 `Connect` 重新建立連線。
3.  若重連成功，則繼續執行操作；若失敗，則返回錯誤。

這種「惰性重連（Lazy Reconnect）」策略確保了系統在短暫網路波動後的自我修復能力。

## 效能考量

*   **高頻輪詢 (< 1s)**：強烈建議使用 **長連接模式**。TCP 握手可能佔用 10-100ms，在 100ms 的輪詢週期中佔比過高。
*   **低頻監控 (> 10s)**：可考慮使用 **短連接模式**，特別是當設備同時服務多個上位機系統，且連線數受限時。
