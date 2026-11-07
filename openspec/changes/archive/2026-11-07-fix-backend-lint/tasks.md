## 1. 清除後端靜態檢查問題

- [x] 1.1 修正錯誤處理、型別與 bit 轉換邊界，維持空列表和既有讀寫語義，新增相關回歸測試並執行受影響 package 測試。
- [x] 1.2 完成等價程式整理與測試 helper context 傳遞，以 golangci-lint run ./... 零問題確認，不變更檢查規則。

## 2. 驗證與交付

- [x] 2.1 [after: 1.1, 1.2] 執行完整 go test、go vet 和 lint，分列跳過與未執行範圍；完成 spectra-verify、spectra-review 並修復發現的問題。
- [x] 2.2 [after: 2.1] 驗證 artifacts、檔案行數和 diff，準備限本變更的封存與提交清單，核對遠端 main 為可安全快轉的祖先；發佈結果另記交付報告。
