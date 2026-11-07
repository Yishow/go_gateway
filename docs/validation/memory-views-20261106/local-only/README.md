# 本機封存與資料庫

此目錄內的資料庫、壓縮檔與原始 manifest 已由 `.gitignore` 排除；只有本說明與 ignore 規則可進 Git。它們保留在本機 docs 內，刪除舊 worktree／evidence 不會刪到這裡；若要跨機器保存，須另外備份這個目錄。

- `original-evidence.tar.gz`：原 evidence 除 4 個舊編譯執行檔外的全部檔案；每個成員已與原始 SHA-256 核對。包含成功／失敗測試、完整原始報告、JSON、SQL／SSE／TCP 與截圖、歷史腳本，以及當時的測試 DB。
- `original-evidence-manifest.json`：完整原始成員及排除執行檔清單，含 SHA-256 與大小。
- `worktree-datalink.sqlite`：舊 worktree 的獨立一致性快照，包含 WAL 中已提交資料；`PRAGMA integrity_check` 為 ok。
- `evidence-gateway-test.sqlite`：原 evidence 測試 DB 的同類快照，完整性檢查為 ok。
- `database-snapshots.json`：兩份快照雜湊與檢查結果。不需要原 DB 的 WAL／SHM 就能讀取；不要把舊 WAL 複製到快照旁。

唯讀查看快照可使用 SQLite URI 的 `mode=ro`；不要讓 Gateway 直接把保留證據當作工作 DB。原始封存如需展開，請另建空目錄後解壓，保留本檔不變。沒有將快照覆蓋回 main 使用中的資料庫，也沒有執行資料 migration。
