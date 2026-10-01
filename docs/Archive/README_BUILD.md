# 構建說明

## 快速開始

### 一鍵構建（推薦）

```powershell
# 基本構建（清理 + 測試 + 構建）
.\build.ps1

# 清理舊文件並構建
.\build.ps1 -Clean

# 構建多平台版本
.\build.ps1 -Cross

# 跳過測試（僅構建）
.\build.ps1 -Test:$false

# 跳過 Lint 檢查
.\build.ps1 -Lint:$false

# 詳細輸出
.\build.ps1 -Verbose

# 指定版本號
.\build.ps1 -Version "v1.0.0"

# 指定輸出目錄
.\build.ps1 -Output "dist"
```

## 參數說明

- `-Clean`: 清理舊構建文件和輸出目錄
- `-Test`: 運行測試（默認：true）
- `-Lint`: 運行 Lint 檢查（默認：true）
- `-Cross`: 構建多平台版本（Windows/Linux/Mac，amd64/arm64）
- `-Output`: 輸出目錄（默認：bin）
- `-Version`: 版本號（默認：自動從 git 檢測）
- `-Verbose`: 顯示詳細輸出

## 構建流程

1. 檢查 Go 環境
2. 清理舊文件（如果指定 -Clean）
3. 運行測試（如果指定 -Test）
4. 運行 Lint 檢查（如果指定 -Lint）
5. 構建 Windows 版本
6. 構建多平台版本（如果指定 -Cross）

## 輸出文件

構建完成後，文件將輸出到 `bin` 目錄：

- `fatek_test.exe` - Windows 版本
- `fatek_test-<os>-<arch>` - 多平台版本（如果使用 -Cross）

## 手動構建

如果不想使用構建腳本，也可以手動構建：

```powershell
# 構建 Windows 版本
go build -o bin\fatek_test.exe ./cmd/fatek_test

# 構建 Linux 版本
$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o bin\fatek_test-linux-amd64 ./cmd/fatek_test

# 構建 Mac 版本
$env:GOOS="darwin"
$env:GOARCH="amd64"
go build -o bin\fatek_test-darwin-amd64 ./cmd/fatek_test
```

## 故障排除

### 權限問題

如果遇到執行策略錯誤，請運行：

```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### Go 環境問題

確保已安裝 Go 1.25.5 或更高版本：

```powershell
go version
```

### 依賴問題

確保所有依賴已下載：

```powershell
go mod download
go mod tidy
```
