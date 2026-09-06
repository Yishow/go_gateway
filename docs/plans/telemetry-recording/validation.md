# 本次提案驗證

範圍：提案與設計文件；不是產品實作驗收。

## 已執行

- 核對main基準與建立獨立分支；本次檔案全部位於openspec/changes與docs/plans/telemetry-recording。
- 本機離線檢查：6份change、12份delta spec、59條requirement、80個scenario、85項未勾選實作任務。
- 確认proposal/design/tasks必要章節、ADDED/MODIFIED標題、規範關鍵字、WHEN/THEN情境。
- 確認每份proposal宣告的capability與實際spec檔一致，7個modified requirement各只有一個擁有者。
- 檢查相依圖無循環、相對Markdown連結存在、所有實作任務維持未完成。
- 12個Decimal/identity/容量數值fixture通過；這只驗證文件數字，不代表後端已實作。
- 在只包含本次新增檔案的本機驗證倉庫執行git diff --check，無空白錯誤。

## GitHub 文件同步（2026-09-07）

- 更新前再次核對 main 與目標分支，兩者皆為 `4762f566ea28a1dbb4bbe313b156462d4af52517`。
- 將六份 change 的36份文件，以及9份共用規畫、契約和驗收文件，一併提交至 `spec/telemetry-recording-plans-20260907`。
- 同步 README、progress、task_plan 與本文件的交付說明；不修改 main、canonical specs 或產品程式。
- 本輪重新檢查45份文件清單、必要artifact、12份delta spec、59條requirement、80個scenario、85項未勾選實作任務、相對連結與12個數值fixture。
- 在僅包含本次新增文件的本機驗證倉庫重跑 `git diff --cached --check`。
- 遠端交付以包含本文件的 commit 與目標分支為準；先前未完成提交的敘述只適用上一輪草稿包。

## 未執行與原因

- 官方OpenSpec CLI：起草時環境未安裝；當時嘗試npm取得套件，registry.npmjs.org DNS回EAI_AGAIN。
  本輪僅同步文件，沒有重跑官方CLI。
  因此沒有宣稱執行過openspec new/status/instructions或官方strict驗證。
  本次依repo openspec-propose skill與spec-driven格式建立所有artifact；離線檢查不取代官方CLI。
- Go/前端/資料庫/PLC/平台/性能測試：本次只起草，沒有修改runtime程式，這些仍是各change的實作任務。
- 不將任何未執行的測試標記通過；所有產品測試與性能目標見acceptance.md。

## 接手後的官方檢查

在有repo完整checkout與OpenSpec CLI的環境執行：

```bash
for change in \
  add-measurement-semantics \
  add-recording-plans \
  add-usage-and-aggregation \
  add-durable-recording-delivery \
  redesign-studio-recording-flow \
  add-recording-history-reports
do
  openspec status --change "$change" --json
  openspec validate "$change" --strict
done
git diff --check
```

先補官方驗證再依README依賴實作；每項任務以實際測試結果更新，不能因artifact齊全就把實作勾完。
