# 草案檢查與交接

## 狀態

本案是 **文件草案／產品實作未開始／正式 CLI 檢查待執行**。只新增本 change 內的文件，沒有修改 `openspec/specs/`、產品程式、測試程式或資料庫。

本次建立的文件共 8 份：`.openspec.yaml`、`proposal.md`、`design.md`、`tasks.md`、`evidence.md`、本檔，以及兩份差異規格。tasks.md 有 14 項未完成工作；數字是本次文件計數，不代表工期或測試覆蓋率。

## 本次已執行的檢查

初步結果：115 個簡單檢查均通過；文件包含 13 個需求、53 個情境與 14 項未完成待辦。這些是文件結構的實際計數，不是產品測試通過數。

在提交前對實際文件執行獨立的簡單檢查：必要文件存在、差異規格標題與 WHEN／THEN 情境、MODIFIED 名稱對照已讀主規格、每個需求有待辦對照、待辦編號唯一且依賴無循環、無已完成的實作勾選、無空白錯誤、文件路徑只限此 change。這只是格式及一致性初查，不是 Spectra analyzer 或官方驗證。

本次 local `git diff --cached --check` 僅針對新增提案文件的暫存工作目錄執行；完整倉庫未成功 clone，不能據此宣稱全倉庫檢查通過。GitHub 提交後另核對實際 commit 的變更清單與分支來源。

## 無法在本次環境執行

- `spectra` 與 `openspec` 指令均未安裝；網路名稱解析失敗，未安裝替代工具。
- 未執行 `spectra new change`／`instructions`／`analyze`／`validate`／`park`。文件依現有 OpenSpec spec-driven 目錄、主規格與已讀專案規範手工起草，並非宣稱完成 Spectra skill 的全套流程。
- 未執行 OpenSpec 官方 validate，未將簡單自製檢查標記為官方通過。
- 沒有開放或驗證產品功能，因此未跑前後端產品測試、完整建置、瀏覽器、真資料庫與硬體測試；這些保持在待辦中。

## 接手順序

1. 取得此分支，重新查 main 並比較新變更；遇到同功能更新，先更新本案來源、設計與任務，不直接套舊修補包。
2. 使用專案已採用的 Spectra CLI 讀取此 change 的 status 與 instructions，核對當前模板／規則。工具版本不同時先確認指令說明，不安裝或改寫專案流程來假造通過。
3. 執行下列正式文件檢查，修正問題後再記錄真實輸出；只有通過且 owner 進一步要求實作，才執行 tasks.md。

```bash
spectra status --change fix-studio-v2-database-workflow --json
spectra analyze fix-studio-v2-database-workflow --json
spectra validate fix-studio-v2-database-workflow
```

本次保留於 `openspec/changes/fix-studio-v2-database-workflow/`，未冒充已 park 或已 archive；有工具的環境完成提案工作流時，按當前專案規範處理 park，記錄真實位置和狀態。

## 產品實作後的驗收基準

依當時的 AGENTS.md 執行，至少包含：

```bash
cd frontend && npm run lint && npm test -- --run && npm run build
npm run test:e2e
cd ..
go test ./...
go vet ./...
golangci-lint run ./...
make check-lines
git diff --check
```

真資料庫驗證需用可丟棄環境；核對所選目標、真表、測試列內容、回讀、清理、斷線後重試及正式資料未受影響。390／768／1440 CSS px 與中文輸入測試是計畫條件，尚無驗收截圖。

如後續修改 Studio inventory，必須同步使用 `go run ./cmd/studio_inventory_changelog` 更新其 changelog；本次未碰該目錄，沒有虛構 changelog 紀錄。
