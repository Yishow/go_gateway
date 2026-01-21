# p3-enable-virtual-device-mode 任務清單 (TDD Edition)

## 階段 1：虛擬記憶體核心 (Virtual Memory Core)

- [ ] **1.1** [RED] 定義 `MemoryBank` 行為測試。
  - 建立 `internal/virtual/memory/bank_test.go`。
  - **Case A (Byte讀寫)**: WriteByte(Offset 0, 0xFF) -> ReadByte(0) == 0xFF。
  - **Case B (Word讀寫)**: WriteWord(Offset 0, 0x1234) -> ReadByte(0)==0x12, ReadByte(1)==0x34 (BigEndian)。
  - **Case C (越界)**: WriteByte(Offset 99999) -> Error。
- [ ] **1.2** [GREEN] 實作 `MemoryBank`。
  - 使用 `[]byte` slice 模擬記憶體。
  - 實作 BigEndian/LittleEndian 處理。
  - 通過 1.1 測試。

- [ ] **1.3** [RED/GREEN] 實作 `AddressMapper`。
  - 測試：Map("40001") -> Offset 0。
  - 測試：Map("D100") -> Offset 200 (假設 D=2 bytes)。
  - 實作邏輯確保映射正確。

## 階段 2：Modbus TCP Server 實作

- [ ] **2.1** [RED] 建立 Server 集成測試。
  - 建立 `protocol/server/modbus/integration_test.go`。
  - 啟動 Server (空實作) 監聽隨機 Port。
  - 使用 Go 的 Modbus Client (`simonvetter/modbus` 或其他庫) 連接。
  - 嘗試寫入 40001=12345，讀取 40001，預期 Timeout 或 Error。

- [ ] **2.2** [GREEN] 實作 Server 請求處理。
  - 實作 Socket Listener。
  - 解析 Modbus Frame。
  - 串接 `MemoryBank`。
  - 通過 2.1 測試 (Client 成功讀回 12345)。

## 階段 3：模擬與橋接功能

- [ ] **3.1** [RED/GREEN] 實作 `SimulationEngine`。
  - 測試：設定規則 "D100 AutoIncrement 1s"。
  - 執行 2秒，斷言 `MemoryBank` 中 D100 的值增加。
- [ ] **3.2** [TDD] 實作 Datalink 寫入虛擬設備。
  - 測試：`WriteToTarget("Virtual:D100", 999)`。
  - 斷言：虛擬記憶體 D100 變為 999。

## 階段 4：API (API First)

- [ ] **4.1** [RED] 定義 API 契約測試。
  - `GET /api/virtual/memory/dump`。
  - 期望回傳 JSON 格式的記憶體快照。

- [ ] **4.2** [GREEN] 實作 API Handler。
