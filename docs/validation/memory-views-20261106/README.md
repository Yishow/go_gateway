# 記憶體五視圖與完整記錄驗證交接

此目錄集中保留 2026-11-05～06 的原圖、驗證結論、修復紀錄與本機資料，供移除舊 worktree／evidence 目錄前使用。整理沒有重新執行產品驗收，也沒有修改產品程式或提交 Git。

## 先看這些

| 用途 | 入口 |
| --- | --- |
| 五張原圖、明暗／寬度切換、逐項版面判定 | [互動對照](visual/comparison.html) |
| 五視圖原圖差異與 30 組實際畫面結論 | [原圖重驗報告](reports/visual-revalidation.md) |
| UI 新增設備到 SQLite／PostgreSQL 的實證 | [完整記錄報告](reports/recording-e2e.md) |
| 五視圖數值、斷線與新值復原 | [數值與恢復](reports/memory-values-recovery.md) |
| mixed 試寫 422 的設計限制 | [能力邊界](reports/testwrite-limit.md) |
| 新增範圍的空 DB 版本修復 | [修復報告](reports/fix-add-range.md) |
| 已存範圍啟用、首頁背景重掛與延遲量測 | [修復報告](reports/fix-activation-home-refresh.md) |
| 可重用的 E2E skill 與踩坑筆記 | [Skill 快照](skill/SKILL.md)、[經驗](skill/references/lessons.md) |
| 本機資料庫與完整原始附件 | [本機封存說明](local-only/README.md) |

## 結論與限制

- 實作已合入 main，提交 `0f54740e2df083ff0e6695897bfa4963039d5fda`，含原 4 個 feature 提交及來源修復。
- 實際 UI 建立 MC／Modbus 六個 int16 點，以新草稿預設 latest／15 秒寫入 SQLite 與 PostgreSQL；各三個完整區間對上來源、版本、SQL 與雙 receipt。五視圖與斷線後新值復原另有實證。
- mixed 試寫 execution／readback／cleanup 不支援，不能稱試寫成功。PG 停用後已驗有界穩定，ACK 即刻停止 intake 尚未證實。
- 原圖版面仍有已列出的結構差異；數值與互動 PASS 不能當成外觀一致。手機暗色、單獨 stale、完整型別／多 entity／三次 fresh run 依各報告分別列已驗與未驗。
- 一條命令的整合 runner 尚未 fresh 重跑；保留的是實際逐段 UI 證據，不以語法檢查代替 E2E。
- 舊報告中的「未提交」等措辭是當時狀態；此整理保留歷史結論，沒有改成新的 PASS。

## 保留與排除

- `visual/`：五張原圖、30 張真實 HSL 畫面、36 張明示 SAMPLE 圖與互動式對照；圖片按位元組比對。
- `images/`：使用者回報貼圖、最後五視圖與故障／恢復畫面，以及兩次來源修復的截圖。
- `reports/`：可直接閱讀的報告；只正規化私人路徑與附件連結，完整原檔未修改地保存在本機封存。
- `skill/`：已安裝 skill 的逐位元組快照；全域安裝仍獨立存在，不受移除舊 worktree 影響。這份 docs 快照不會自動安裝成 repo skill。
- `local-only/`：完整原始 evidence 封存、兩份包含已提交 WAL 資料的獨立 SQLite 快照、原始清單與 SHA-256。它們可能包含本機設定或歷史私人路徑，已由該目錄的 `.gitignore` 排除，不能當成已上 Git 的備份。
- 排除原 evidence 根目錄的 4 個舊 `test-ui*` 執行檔；它們是可重建輸出，保留其名稱、大小與雜湊。完整證據、原始報告、失敗嘗試、重現腳本和圖片仍保留。
- 舊 worktree 的 node_modules、bin、dist、embed 靜態輸出與測試快取不搬移；程式碼已在 main。SQLite 用一致性快照取代需成套配對的 DB／WAL／SHM。

## 移除舊目錄前

1. 確認 [保留清單](RETENTION.json) 與 [核對結果](VERIFICATION.md)。本機封存需連同此目錄另外備份，單獨 Git clone 不會帶回 ignored 備份。
2. 舊目錄的 `start.sh`、Air、Gateway 與 Vite 已依指示停止；8080／5173 與資料庫檔案占用均已核對。不要再從舊目錄啟動服務後直接刪除。
3. 使用 Git worktree 管理方式移除舊 `go_gateway-memory-validation`；evidence 原目錄再由使用者自行移除。本次沒有刪除這兩個來源。
4. 原始附件與腳本仍記錄歷史埠及路徑；需要重跑時先換成新 owned run，不能直接對備份或正式資料執行。

相關規格：[device-to-SQL](../../../openspec/specs/studio-v2-device-to-sql-acceptance/spec.md)、[記錄生命週期](../../../openspec/specs/recording-lifecycle-management/spec.md)、[記憶體視圖 change](../../../openspec/changes/add-switchable-memory-views/proposal.md)。
