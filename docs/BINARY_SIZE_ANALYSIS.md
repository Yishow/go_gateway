# .exe 程式肥大分析報告

## 📊 現況分析

### 檔案大小

- **當前 gateway.exe**: 28.89 MB
- **優化後（使用 `-s -w`）**: 20.37 MB
- **可減少**: 8.52 MB（約 29.5%）

### 前端資源大小

嵌入的前端靜態檔案總計約 **320 KB**：

- `index-B2F_Dyh8.js`: 305.91 KB（約 313 KB）
- `index-CEapu_ni.css`: 12.96 KB（約 13 KB）
- `index.html`: 0.44 KB
- `vite.svg`: 1.46 KB

**結論**：前端資源僅佔總大小的 **1.1%**，不是主要問題。

---

## 🔍 肥大原因分析

### 1. **除錯資訊與符號表（主要因素）**

**影響**：約 8.52 MB（29.5%）

當前編譯命令：

```powershell
go build -ldflags "-H=windowsgui" -o bin/gateway.exe ./cmd/test_ui
```

**問題**：

- 未使用 `-s`（strip symbol table）移除符號表
- 未使用 `-w`（omit DWARF symbol table）移除除錯資訊

**解決方案**：

```powershell
go build -ldflags "-H=windowsgui -s -w" -o bin/gateway.exe ./cmd/test_ui
```

### 2. **大型依賴庫**

#### 主要依賴分析

| 依賴                 | 類型     | 說明          | 影響                     |
| -------------------- | -------- | ------------- | ------------------------ |
| `quic-go/quic-go`    | 間接依賴 | QUIC 協議實作 | ⚠️ 大型（可能 5-10 MB）  |
| `getlantern/systray` | 直接依賴 | 系統托盤功能  | ⚠️ 中型（可能 2-3 MB）   |
| `gin-gonic/gin`      | 直接依賴 | HTTP 框架     | ✅ 中型（約 1-2 MB）     |
| `go.bug.st/serial`   | 直接依賴 | 串列埠通訊    | ✅ 小型（約 0.5 MB）     |
| `bytedance/sonic`    | 間接依賴 | JSON 編解碼器 | ✅ 小型（約 0.5 MB）     |
| `golang.org/x/*`     | 間接依賴 | 標準擴展庫    | ✅ 小型（總計約 1-2 MB） |

**特別注意**：

- `quic-go/quic-go` 是透過 `gin-gonic/gin` 的間接依賴引入的，但可能未被實際使用
- `getlantern/systray` 可能包含 CGO 依賴（雖然當前 CGO_ENABLED=0）

### 3. **Go 運行時（Runtime）**

**影響**：約 2-3 MB（正常範圍）

Go 編譯器會將運行時庫打包進二進位檔，這是正常且必要的。

### 4. **前端資源嵌入**

**影響**：約 320 KB（1.1%）

使用 `//go:embed static` 嵌入的前端資源，佔比很小。

---

## 💡 優化建議

### 優先級 1：立即優化（可減少約 30%）

#### 1.1 移除除錯資訊

修改 `start.ps1` 和 `build.ps1` 中的編譯命令：

```powershell
# 修改前
go build -ldflags "-H=windowsgui" -o $outputPath $target.Path

# 修改後
go build -ldflags "-H=windowsgui -s -w" -o $outputPath $target.Path
```

**效果**：可減少約 8.5 MB（29.5%）

#### 1.2 啟用編譯優化

```powershell
go build -ldflags "-H=windowsgui -s -w" -trimpath -o $outputPath $target.Path
```

**說明**：

- `-s`: 移除符號表
- `-w`: 移除 DWARF 除錯資訊
- `-trimpath`: 移除檔案路徑資訊（提升可移植性）

### 優先級 2：依賴優化（可減少約 10-20%）

#### 2.1 檢查未使用的依賴

```bash
# 檢查未使用的依賴
go mod tidy
go mod verify

# 使用工具分析依賴大小
go install github.com/nikolaydubina/go-binsize-treemap/cmd/go-binsize-treemap@latest
go-binsize-treemap bin/gateway.exe
```

#### 2.2 移除未使用的 QUIC 依賴

如果未使用 HTTP/3（QUIC），可以考慮：

- 檢查 `gin-gonic/gin` 是否真的需要 QUIC 支援
- 或使用更輕量的 HTTP 框架（如 `net/http` + `gorilla/mux`）

#### 2.3 考慮替換 systray

如果系統托盤功能不是核心需求，可以考慮：

- 使用更輕量的替代方案
- 或移除托盤功能，改用瀏覽器通知

### 優先級 3：前端優化（可減少約 1-2%）

#### 3.1 前端資源壓縮

- 確保 Vite 建置時啟用壓縮
- 考慮使用 Brotli 或 Gzip 壓縮（需在運行時解壓）

#### 3.2 移除未使用的前端依賴

檢查 `web/test-ui/package.json`，移除未使用的依賴：

- `recharts`（如果未使用圖表功能）
- 其他未使用的 React 套件

### 優先級 4：進階優化（可減少約 5-10%）

#### 4.1 使用 UPX 壓縮（可選）

```powershell
# 安裝 UPX
# 下載：https://upx.github.io/

# 壓縮（可能影響啟動速度）
upx --best bin/gateway.exe
```

**注意**：UPX 壓縮可能被防毒軟體誤報，且可能影響啟動速度。

#### 4.2 分離前端資源（進階）

如果前端資源很大，可以考慮：

- 不嵌入前端，改為從外部目錄讀取
- 或使用 CDN 載入前端資源

---

## 📈 預期優化效果

| 優化項目              | 預期減少    | 累計大小     |
| --------------------- | ----------- | ------------ |
| 當前大小              | -           | 28.89 MB     |
| 移除除錯資訊（-s -w） | -8.52 MB    | 20.37 MB     |
| 依賴優化              | -2~4 MB     | 16~18 MB     |
| 前端優化              | -0.1~0.3 MB | 15.9~17.9 MB |
| UPX 壓縮（可選）      | -5~7 MB     | 10~12 MB     |

**最終預期**：約 **15-18 MB**（不使用 UPX）或 **10-12 MB**（使用 UPX）

---

## 🛠️ 實作步驟

### 步驟 1：修改編譯腳本

修改 `start.ps1` 和 `build.ps1`，加入 `-s -w` 標誌。

### 步驟 2：驗證優化效果

```powershell
# 重新編譯
.\start.ps1 -SkipLint -SkipTest -SkipQuality

# 檢查檔案大小
Get-ChildItem bin/*.exe | Select-Object Name, @{Name="Size(MB)";Expression={[math]::Round($_.Length/1MB, 2)}}
```

### 步驟 3：測試功能

確保優化後的程式功能正常：

- 系統托盤功能
- HTTP 服務
- 前端頁面載入
- 協議通訊功能

### 步驟 4：依賴分析（可選）

```bash
# 分析依賴大小
go install github.com/nikolaydubina/go-binsize-treemap/cmd/go-binsize-treemap@latest
go-binsize-treemap bin/gateway.exe > dependency_analysis.html
```

---

## 📝 注意事項

1. **除錯資訊移除後**：

   - 無法使用 `go tool pprof` 進行效能分析
   - 無法使用 `go tool trace` 進行追蹤
   - 錯誤堆疊資訊會較簡略

2. **建議保留兩個版本**：

   - 開發版本：保留除錯資訊（用於開發和除錯）
   - 發布版本：移除除錯資訊（用於分發）

3. **CGO 依賴**：
   - 當前 `CGO_ENABLED=0`，這是好的（純 Go 編譯）
   - 如果未來需要使用 CGO，檔案大小會進一步增加

---

## 🔗 參考資源

- [Go Build Flags](https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies)
- [Go Binary Size Optimization](https://go.dev/doc/tutorial/compile-install)
- [UPX - Ultimate Packer for eXecutables](https://upx.github.io/)
