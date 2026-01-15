# Go 程式調試指南

本指南介紹多種調試 Go 程式的方法，從簡單到進階。

## 1. 使用 fmt.Printf 進行簡單調試

最簡單的調試方法，適合快速查看變數值。

```go
package main

import "fmt"

func main() {
    value := 42
    fmt.Printf("DEBUG: value = %d\n", value)
    fmt.Printf("DEBUG: value (hex) = %x\n", value)
    fmt.Printf("DEBUG: value (binary) = %b\n", value)
}
```

### 進階：條件編譯調試

創建一個調試輔助文件：

```go
// +build debug

package main

var Debug = true

func debugPrint(format string, args ...interface{}) {
    if Debug {
        fmt.Printf("[DEBUG] "+format+"\n", args...)
    }
}
```

使用時：
```go
//go:build debug
// +build debug

debugPrint("變數值: %v", myVar)
```

編譯時：
```bash
# 調試版本
go build -tags debug

# 發布版本
go build
```

## 2. 使用 log 包進行結構化日誌

```go
import (
    "log"
    "os"
)

func main() {
    // 設置日誌格式
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
    
    // 設置輸出到文件
    f, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    log.SetOutput(f)
    
    log.Println("程式開始執行")
    log.Printf("變數值: %v", myVar)
}
```

## 3. 使用 Delve (dlv) 調試器

Delve 是 Go 的官方調試器，功能強大。

### 安裝 Delve

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

### 基本使用

```bash
# 啟動調試
dlv debug ./cmd/fatek_test

# 設置斷點
(dlv) break main.main
(dlv) break internal/protocol/fatek/client.go:38

# 運行
(dlv) continue
(dlv) c

# 單步執行
(dlv) next
(dlv) n

# 進入函數
(dlv) step
(dlv) s

# 查看變數
(dlv) print variableName
(dlv) p variableName

# 查看局部變數
(dlv) locals

# 查看堆棧
(dlv) stack
(dlv) bt

# 查看所有斷點
(dlv) breakpoints

# 刪除斷點
(dlv) clear <breakpoint-id>

# 退出
(dlv) exit
```

### 調試運行中的程式

```bash
# 附加到運行中的進程
dlv attach <pid>

# 或者讓程式等待調試器連接
dlv exec ./fatek_test.exe -- -mode=tcp -host="127.0.0.1" -action=read
```

### 調試範例

```bash
# 1. 啟動調試
dlv debug ./cmd/fatek_test

# 2. 設置斷點在 ReadRandom 函數
(dlv) break internal/protocol/fatek/client.go:211

# 3. 運行到斷點
(dlv) continue

# 4. 查看參數
(dlv) args

# 5. 查看局部變數
(dlv) locals

# 6. 單步執行
(dlv) next

# 7. 查看特定變數
(dlv) print body
(dlv) print items
```

## 4. 使用 VS Code 調試

VS Code 有優秀的 Go 調試支持。

### 安裝擴展

1. 安裝 "Go" 擴展（由 Go Team 提供）
2. 安裝 "Delve"（調試器）

### 創建調試配置

在項目根目錄創建 `.vscode/launch.json`：

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug Fatek Test",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/cmd/fatek_test",
            "args": [
                "-mode=tcp",
                "-host=127.0.0.1",
                "-port=2000",
                "-action=random",
                "-random=X:0,D:0,R:10"
            ],
            "env": {},
            "showLog": true
        },
        {
            "name": "Debug Fatek Test (Serial)",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/cmd/fatek_test",
            "args": [
                "-mode=serial",
                "-serial=COM3",
                "-action=read",
                "-symbol=X",
                "-addr=0",
                "-count=10"
            ]
        },
        {
            "name": "Attach to Process",
            "type": "go",
            "request": "attach",
            "mode": "local",
            "processId": 0
        },
        {
            "name": "Debug Tests",
            "type": "go",
            "request": "launch",
            "mode": "test",
            "program": "${workspaceFolder}",
            "args": [
                "-test.v"
            ]
        }
    ]
}
```

### 使用步驟

1. 在代碼中設置斷點（點擊行號左側）
2. 按 `F5` 開始調試
3. 使用調試工具欄：
   - `F5` - 繼續
   - `F10` - 單步跳過
   - `F11` - 單步進入
   - `Shift+F11` - 單步跳出
   - `Shift+F5` - 停止調試

## 5. 使用條件斷點和日誌點

### VS Code 條件斷點

1. 設置斷點
2. 右鍵點擊斷點
3. 選擇 "Edit Breakpoint"
4. 輸入條件，例如：`len(items) > 3`

### 日誌點（Logpoints）

在 VS Code 中：
1. 設置斷點
2. 右鍵選擇 "Add Logpoint"
3. 輸入日誌表達式，例如：`Items: {len(items)}, Body: {body}`

## 6. 調試網路程式

對於網路程式（如我們的 Fatek 客戶端），可以：

### 使用 Wireshark 抓包

```bash
# 過濾 TCP 流量
tcp.port == 2000
```

### 添加調試輸出

在 `transport.go` 中添加：

```go
func (t *TCPTransport) SendReceive(data []byte) ([]byte, error) {
    // 調試輸出
    if os.Getenv("DEBUG") == "1" {
        fmt.Printf("[DEBUG] Sending: %q\n", string(data))
        fmt.Printf("[DEBUG] Sending (hex): %x\n", data)
    }
    
    // ... 原有代碼 ...
    
    if os.Getenv("DEBUG") == "1" {
        fmt.Printf("[DEBUG] Received: %q\n", string(response))
        fmt.Printf("[DEBUG] Received (hex): %x\n", response)
    }
    
    return response, nil
}
```

使用：
```bash
$env:DEBUG="1"; .\fatek_test.exe -mode=tcp -host="127.0.0.1" -action=read
```

## 7. 調試測試

### 調試單個測試

```bash
dlv test ./internal/protocol/fatek -- -test.run TestReadRegisters
```

### 在 VS Code 中調試測試

1. 打開測試文件
2. 點擊測試函數上方的 "debug test"
3. 或使用調試配置中的 "Debug Tests"

## 8. 性能調試

### 使用 pprof

```go
import (
    _ "net/http/pprof"
    "net/http"
)

func main() {
    // 啟動 pprof 服務器
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    
    // ... 你的代碼 ...
}
```

查看性能數據：
```bash
# CPU 性能分析
go tool pprof http://localhost:6060/debug/pprof/profile

# 內存分析
go tool pprof http://localhost:6060/debug/pprof/heap

# 查看網頁界面
go tool pprof -http=:8080 http://localhost:6060/debug/pprof/profile
```

## 9. 實用調試技巧

### 打印結構體

```go
import "fmt"

// 使用 %+v 打印結構體字段名
fmt.Printf("%+v\n", myStruct)

// 使用 %#v 打印 Go 語法表示
fmt.Printf("%#v\n", myStruct)
```

### 打印字節數組

```go
fmt.Printf("Bytes: %x\n", data)        // 十六進制
fmt.Printf("Bytes: %q\n", data)        // 可打印字符
fmt.Printf("Bytes: %+q\n", data)      // 轉義字符
```

### 使用 recover 捕獲 panic

```go
func safeFunction() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Panic recovered: %v\n", r)
            debug.PrintStack()
        }
    }()
    
    // 可能 panic 的代碼
}
```

### 使用 runtime.Caller 追蹤調用棧

```go
import "runtime"

func debugTrace() {
    pc, file, line, ok := runtime.Caller(1)
    if ok {
        fmt.Printf("Called from %s:%d (function: %s)\n", 
            file, line, runtime.FuncForPC(pc).Name())
    }
}
```

## 10. 推薦的調試工作流

1. **開發階段**：使用 `fmt.Printf` 快速調試
2. **複雜問題**：使用 Delve 或 VS Code 斷點調試
3. **網路問題**：使用 Wireshark + 調試輸出
4. **性能問題**：使用 pprof
5. **生產環境**：使用結構化日誌（log 包）

## 11. 針對本項目的調試建議

### 調試 Fatek 協議

```go
// 在 client.go 的 execute 函數中添加
func (c *FatekClient) execute(cmd, body string) (string, error) {
    if os.Getenv("FATEK_DEBUG") == "1" {
        fmt.Printf("[FATEK] Command: %s, Body: %q\n", cmd, body)
    }
    // ... 原有代碼 ...
}
```

### 調試地址格式化

```go
// 在 address.go 中添加
func FormatAddress(comp ComponentType, addr int) string {
    result := fmt.Sprintf(format, comp.Name, addr)
    if os.Getenv("FATEK_DEBUG") == "1" {
        fmt.Printf("[FATEK] FormatAddress: %s, %d -> %q\n", 
            comp.Name, addr, result)
    }
    return result
}
```

### 調試測試工具

```bash
# 使用 Delve 調試
dlv debug ./cmd/fatek_test -- -mode=tcp -host="127.0.0.1" -action=read

# 或設置環境變數
$env:FATEK_DEBUG="1"
.\fatek_test.exe -mode=tcp -host="127.0.0.1" -action=read
```

## 參考資源

- [Delve 官方文檔](https://github.com/go-delve/delve)
- [VS Code Go 調試指南](https://github.com/golang/vscode-go/wiki/debugging)
- [Go 官方調試指南](https://go.dev/doc/diagnostics)
