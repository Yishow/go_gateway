# Fatek PLC 測試工具

此工具用於測試 Fatek PLC 的 TCP 和 Serial 通訊功能，支援讀取、寫入和隨機讀取操作。

## 編譯

```bash
go build -o fatek_test.exe ./cmd/fatek_test
```

或在 Linux/Mac 上：

```bash
go build -o fatek_test ./cmd/fatek_test
```

## 使用方法

> **重要提示**：在 Windows PowerShell 中，IP 位址和包含逗號的參數值必須用引號包裹，例如 `-host="192.168.1.5"` 或 `-random="X:0,D:0,R:10"`。在 CMD 或 Linux/Mac 中，引號是可選的。

### TCP 模式

#### 讀取離散狀態 (X, Y, M, S, T, C)

```bash
# 使用預設埠號 500（Windows PowerShell 需用引號）
fatek_test.exe -mode=tcp -host="192.168.1.5" -action=read -symbol=X -addr=0 -count=10

# 指定自訂埠號
fatek_test.exe -mode=tcp -host="192.168.1.5" -port=502 -action=read -symbol=X -addr=0 -count=10
```

#### 讀取暫存器 (R, D, RT, RC, F, DR, DD, DF, WX, WY, WM, WS, WT, WC)

```bash
# 使用預設埠號 500（Windows PowerShell 需用引號）
fatek_test.exe -mode=tcp -host="192.168.1.5" -action=read -symbol=D -addr=0 -count=10

# 指定自訂埠號
fatek_test.exe -mode=tcp -host="192.168.1.5" -port=502 -action=read -symbol=D -addr=0 -count=10
```

#### 寫入離散狀態

```bash
# 使用預設埠號 500（Windows PowerShell 需用引號）
fatek_test.exe -mode=tcp -host="192.168.1.5" -action=write -symbol=Y -addr=0 -values=true,false,true

# 指定自訂埠號
fatek_test.exe -mode=tcp -host="192.168.1.5" -port=502 -action=write -symbol=Y -addr=0 -values=true,false,true
```

#### 寫入暫存器

```bash
# 使用預設埠號 500（Windows PowerShell 需用引號）
fatek_test.exe -mode=tcp -host="192.168.1.5" -action=write -symbol=D -addr=0 -values=100,200,300

# 指定自訂埠號
fatek_test.exe -mode=tcp -host="192.168.1.5" -port=502 -action=write -symbol=D -addr=0 -values=100,200,300
```

#### 隨機讀取（混合讀取多個不同組件）

```bash
# 使用預設埠號 500（Windows PowerShell 需用引號）
fatek_test.exe -mode=tcp -host="192.168.1.5" -action=random -random="X:0,D:0,R:10"

# 指定自訂埠號
fatek_test.exe -mode=tcp -host="192.168.1.5" -port=502 -action=random -random="X:0,D:0,R:10"
```

### Serial 模式

#### 讀取離散狀態

```bash
fatek_test.exe -mode=serial -serial=COM3 -action=read -symbol=X -addr=0 -count=5
```

#### 讀取暫存器

```bash
fatek_test.exe -mode=serial -serial=COM3 -action=read -symbol=D -addr=0 -count=10
```

#### 寫入離散狀態

```bash
fatek_test.exe -mode=serial -serial=COM3 -action=write -symbol=Y -addr=0 -values=true,false,true
```

#### 寫入暫存器

```bash
fatek_test.exe -mode=serial -serial=COM3 -action=write -symbol=D -addr=0 -values=100,200,300
```

#### 隨機讀取

```bash
fatek_test.exe -mode=serial -serial=COM3 -action=random -random="X:0,D:0,R:10"
```

### 參數說明

#### 必填參數

- `-mode`: 連接模式，`tcp` 或 `serial`
- `-action`: 操作類型，`read`、`write` 或 `random`

#### TCP 模式參數

- `-host`: PLC IP 位址（必填）
- `-port`: TCP 埠號（預設 500，可選）
  - 範例：`-port=502` 指定使用埠號 502
- `-station`: PLC 站號（預設 1）
- `-timeout`: 逾時時間（預設 2s）

#### Serial 模式參數

- `-serial`: 串列埠名稱，例如 `COM3`（Windows）或 `/dev/ttyUSB0`（Linux）（必填）
- `-baud`: 波特率（預設 9600）
- `-databits`: 資料位元數（預設 7）
- `-stopbits`: 停止位元數（預設 1）
- `-parity`: 同位檢查，`N`/`E`/`O`（預設 E）
- `-station`: PLC 站號（預設 1）
- `-timeout`: 逾時時間（預設 2s）

#### 讀取操作參數

- `-symbol`: 組件符號（X, Y, M, S, T, C, R, D, RT, RC, F, DR, DD, DF, WX, WY, WM, WS, WT, WC 等）
- `-addr`: 起始位址
- `-count`: 讀取數量

#### 寫入操作參數

- `-symbol`: 組件符號
- `-addr`: 起始位址
- `-values`: 寫入值，逗號分隔
  - 離散狀態：`true,false,true`
  - 暫存器：`100,200,300`

#### 隨機讀取參數

- `-random`: 讀取項目列表，格式為 `SYMBOL1:ADDR1,SYMBOL2:ADDR2`
  - 例如：`"X:0,D:0,R:10"`

### 符號與使用限制說明 🔧

- 符號 **不區分大小寫**（例如 `d` 等同 `D`）。
- 新增支援符號說明：
  - `DD`, `DF`：32-bit 類型（Double/檔案暫存器）。
  - `WX`, `WY`, `WM`, `WS`, `WT`, `WC`：16-bit（以 word 方式存取離散位元）。
- 常見限制：
  - 讀取離散狀態（Cmd 44）最大 `count` = **255**。
  - 讀取/寫入暫存器（Cmd 46/47）對 16-bit 最多 **64** 個；對 32-bit 最多 **32** 個。
  - 隨機讀取（Cmd 48）最大項目數 **64**。
- 若符號或位址超出範圍，命令會回傳錯誤並帶有說明（請檢查回應錯誤碼或日誌）。

### 額外範例

- 讀取 32-bit `DD`（兩個 32-bit 值）

```bash
fatek_test.exe -mode=tcp -host="192.168.1.5" -action=read -symbol=DD -addr=0 -count=2
```

- 以 16-bit word 方式讀取離散 `WX`（5 個 word）

```bash
fatek_test.exe -mode=tcp -host="192.168.1.5" -action=read -symbol=WX -addr=0 -count=5
```

## 範例輸出

### 讀取操作

```
建立 TCP 客戶端: 192.168.1.5:500 (站號: 1, 逾時: 2s)
正在連線...
連線成功！
讀取暫存器: D0-9
結果: [100, 200, 300, 400, 500, 600, 700, 800, 900, 1000]
```

### 寫入操作

```
建立 TCP 客戶端: 192.168.1.5:500 (站號: 1, 逾時: 2s)
正在連線...
連線成功！
寫入暫存器: D0-2 = [100, 200, 300]
寫入成功！
```

### 隨機讀取操作

```
建立 TCP 客戶端: 192.168.1.5:500 (站號: 1, 逾時: 2s)
正在連線...
連線成功！
隨機讀取 3 個項目:
  - X0
  - D0
  - R10

結果:
  X0 = true
  D0 = 100
  R10 = 200
```
