# 一鍵啟動腳本使用說明

## 概述

`start.ps1` 是一個 PowerShell 腳本，用於自動執行 Go Gateway 項目的完整開發工作流程，包括：

1. **一鍵啟動專案** - 構建並啟動後端 API 服務
2. **golangci-lint 靜態分析** - 全面的代碼質量檢查
3. **構建可執行文件** - 編譯後端服務
4. **代碼質量檢查** - 格式檢查和未使用代碼檢查
5. **單元測試** - 運行所有測試並檢查覆蓋率

## 環境配置

### 首次使用

1. **複製環境變數範例文件**
   ```powershell
   Copy-Item .env.sample .env
   ```

2. **編輯配置**（可選）
   根據需要修改 `.env` 文件中的配置值，詳見 [環境變數配置說明](./ENV_CONFIG.md)

### 基本使用

```powershell
# 顯示主選單
.\start.ps1

# 一鍵啟動（構建 + 啟動服務）
.\start.ps1 -QuickStart
```

### 進階選項

```powershell
# 跳過 lint 檢查
.\一鍵啟動.ps1 -SkipLint

# 跳過測試
.\一鍵啟動.ps1 -SkipTest

# 跳過構建
.\一鍵啟動.ps1 -SkipBuild

# 顯示詳細覆蓋率報告（生成 HTML）
.\一鍵啟動.ps1 -Coverage

# 詳細輸出模式
.\一鍵啟動.ps1 -Verbose

# 組合使用
.\一鍵啟動.ps1 -SkipTest -Coverage
```

## 功能詳解

### 1. golangci-lint 靜態分析

自動檢查並安裝 `golangci-lint`（如果未安裝），然後執行：

```bash
golangci-lint run ./...
```

**檢查內容包括：**
- 錯誤處理（errcheck）
- 代碼簡化建議（gosimple）
- 可疑代碼結構（govet）
- 無效賦值（ineffassign）
- 靜態分析（staticcheck）
- 未使用代碼（unused）
- 安全性檢查（gosec）
- 代碼風格（revive, gocritic）
- 拼寫檢查（misspell）
- 以及更多...

**配置文件：** `.golangci.yml`

### 2. 單元測試

執行所有單元測試並顯示覆蓋率：

```bash
go test -cover ./internal/protocol/...
```

**覆蓋率目標：** 85%+

**當前覆蓋率：**
- Fatek: ~62%
- MC Protocol: ~57%
- Modbus: ~48%

### 3. 構建可執行文件

構建以下工具：

- `bin/fatek_test.exe` - Fatek PLC 測試工具
- `bin/test_all.exe` - 完整協議測試工具

### 4. 代碼質量檢查

- 檢查代碼格式（`go fmt`）
- 檢查未使用的代碼（可選）

## 輸出說明

腳本使用顏色編碼的輸出：

- ✅ **綠色** - 成功
- ❌ **紅色** - 錯誤
- ⚠️ **黃色** - 警告
- ℹ️ **青色** - 信息

## 退出碼

- `0` - 所有檢查通過
- `1` - 發現問題（lint 錯誤、測試失敗等）

## 依賴要求

### 必需

- Go 1.25.5+
- PowerShell 5.1+

### 可選（會自動安裝）

- golangci-lint - 如果未安裝，腳本會嘗試自動安裝

## 故障排除

### golangci-lint 安裝失敗

如果自動安裝失敗，手動安裝：

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

或訪問：https://golangci-lint.run/usage/install/

### 測試失敗

檢查測試輸出，確保：
1. 所有依賴已安裝（`go mod download`）
2. 測試環境正確配置
3. 沒有硬體依賴的測試被跳過（RTU 測試需要實際串列埠）

### 構建失敗

檢查：
1. Go 版本是否符合要求
2. 所有依賴是否正確
3. 是否有編譯錯誤

## 示例輸出

```
============================================
   Go Gateway 一鍵啟動腳本
============================================

ℹ️  檢查工具依賴...
✅ golangci-lint 已安裝

[1/4] 執行 golangci-lint 靜態分析...
✅ golangci-lint 檢查通過，未發現問題

[2/4] 執行單元測試...

測試覆蓋率：
ok  	go-gateway/internal/protocol/fatek	coverage: 61.9% of statements
ok  	go-gateway/internal/protocol/mcprotocol	coverage: 57.2% of statements
ok  	go-gateway/internal/protocol/modbus	coverage: 48.2% of statements

✅ 所有測試通過

[3/4] 構建可執行文件...
✅ fatek_test 構建成功: bin/fatek_test.exe
✅ test_all 構建成功: bin/test_all.exe

[4/4] 代碼質量檢查...
✅ 代碼格式正確

============================================
✅ 所有檢查完成！
============================================
```

## 與 CI/CD 集成

腳本設計為可在 CI/CD 管道中使用：

```yaml
# GitHub Actions 示例
- name: Run Quality Checks
  run: |
    pwsh -File ./一鍵啟動.ps1 -SkipBuild
```

## 相關文件

- `.golangci.yml` - golangci-lint 配置文件
- `build.ps1` - 構建腳本（包含跨平台構建）
- `go.mod` - Go 模組定義

## 更新日誌

- **v1.0** - 初始版本
  - 支持 golangci-lint 靜態分析
  - 支持單元測試和覆蓋率檢查
  - 支持自動構建
  - 支持代碼質量檢查
