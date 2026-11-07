# 整理核對結果

本檔只記錄搬移與保留完整性，沒有新增產品 PASS 或改寫原驗收結論。

| 核對 | 結果 |
| --- | --- |
| 原 evidence 保留 | PASS：1,426 個原始檔封存；逐成員 SHA-256 對上原始清單 |
| 可重建輸出排除 | 僅排除原 evidence 根目錄 4 個舊 test-ui 執行檔，大小與 SHA 仍在保留清單 |
| 複製／閱覽副本 | PASS：129 份圖片、報告與 skill 檔案，逐檔核對；正規化副本另記原始及新雜湊 |
| 原圖／圖片可達性 | PASS：五原圖及全部 60 種 layout／寬度／主題／資料來源組合檔案存在；另保留額外 SAMPLE 圖 |
| 對照判定資料 | PASS：HTML 內嵌原圖、criteria 及 labels 與原始版完全相同；只調整報告連結 |
| SQLite | PASS：兩個獨立 backup 快照，已納入 committed WAL，兩者 integrity_check 均為 ok |
| 原始證據／資料庫的 Git 忽略 | PASS：10 個本機封存／資料檔均被 local-only/.gitignore 排除 |
| 公開閱覽副本 | PASS：沒有原本的使用者絕對路徑；原始未修改內容只留在 ignored 封存中 |
| 舊開發服務 | PASS：已停止舊 start.sh／Air／Gateway／Vite，8080／5173 無監聽，舊 DB 無程序占用 |
| 來源資料 | 沒有移除兩個來源目錄；沒有把備份覆寫回 main 的工作 DB |
| Git 範圍 | 僅新增本 docs 目錄，index 空，沒有 commit／push |
| diff／行數／連結檢查 | PASS；沒有改動 line-limit ignore 或放寬規則 |
| UI、SQL、PLC、效能重新驗收 | NOT_RUN：本次是資料整理，歷史驗證與剩餘限制見各報告 |

原始封存及各快照的 SHA、大小見 [RETENTION.json](RETENTION.json)。詳細逐原檔 manifest 與閱覽副本來源對應保留在 local-only 內，跟隨本機備份，不放 Git。
