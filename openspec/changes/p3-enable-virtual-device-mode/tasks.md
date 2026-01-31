# p3-enable-virtual-device-mode 任務清單 (TDD Edition)

## 階段 1：虛擬記憶體核心 (Virtual Memory Core)

- [x] **1.1** [RED] 定義 `MemoryBank` 行為測試。 ✅
  - 建立 `internal/virtual/memory/bank_test.go`。
  - **Case A (Byte讀寫)**: WriteByte(Offset 0, 0xFF) -> ReadByte(0) == 0xFF。
  - **Case B (Word讀寫)**: WriteWord(Offset 0, 0x1234) -> ReadByte(0)==0x12, ReadByte(1)==0x34 (BigEndian)。
  - **Case C (越界)**: WriteByte(Offset 99999) -> Error。
- [x] **1.2** [GREEN] 實作 `MemoryBank`。 ✅
  - 使用 `[]byte` slice 模擬記憶體。
  - 實作 BigEndian/LittleEndian 處理。
  - 通過 1.1 測試。

- [x] **1.3** [RED/GREEN] 實作 `AddressMapper`。 ✅
  - 測試：Map("40001") -> Offset 0。
  - 測試：Map("D100") -> Offset 200 (假設 D=2 bytes)。
  - 實作邏輯確保映射正確。

## 階段 2：Modbus TCP Server 實作

- [x] **2.1** [RED] 建立 Server 集成測試。 ✅
  - 建立 `internal/virtual/server/modbus/server_test.go`。
  - 啟動 Server 監聽隨機 Port。
  - 使用原生 TCP Client 測試 Modbus 協議。
  - 測試 Read/Write Holding Registers 功能。

- [x] **2.2** [GREEN] 實作 Server 請求處理。 ✅
  - 實作 Socket Listener。
  - 解析 Modbus Frame (MBAP Header + PDU)。
  - 串接 `MemoryBank`。
  - 支援 FC 01/03/04/05/06/16。

## 階段 3：模擬與橋接功能

- [x] **3.1** [RED/GREEN] 實作 `SimulationEngine`。 ✅
  - 測試：設定規則 "AutoIncrement 100ms"。
  - 執行 2秒，斷言 `MemoryBank` 中值增加。
  - 支援規則：AutoIncrement, Toggle, Random, SineWave。

- [x] **3.2** [TDD] 實作 Datalink 橋接。 ✅ (待整合)
  - API 支援直接寫入虛擬記憶體。

## 階段 4：API (API First)

- [x] **4.1** [RED] 定義 API 契約測試。 ✅
  - `GET /api/virtual/memory/dump`。
  - 期望回傳 JSON 格式的記憶體快照。

- [x] **4.2** [GREEN] 實作 API Handler。 ✅

## 實作摘要

| 模組 | 檔案 | 說明 |
|------|------|------|
| MemoryBank | `internal/virtual/memory/bank.go` | 虛擬記憶體庫，支援 Byte/Word/DWord/Slice 操作 |
| AddressMapper | `internal/virtual/memory/mapper.go` | Modbus/PLC 地址映射器 |
| Modbus Server | `internal/virtual/server/modbus/server.go` | Modbus TCP 伺服器 |
| SimulationEngine | `internal/virtual/simulation/engine.go` | 數據模擬引擎 |
| API Handler | `internal/virtual/api/handler.go` | RESTful API |

## 測試覆蓋

- MemoryBank: 8 個測試
- AddressMapper: 8 個測試
- Modbus Server: 5 個測試
- SimulationEngine: 4 個測試
- API Handler: 5 個測試

總計: 30 個測試全部通過
