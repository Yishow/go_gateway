## 1. 移除與保留

- [x] 1.1 刪除真 PLC 與長負載測試入口，保留兩個排程 pilot；核對測試清單並跑 pilot。
- [x] 1.2 不再註冊 SQLite 不支援組合，保留其餘 entity SQL 測試並跑 SQLite 分支。
- [x] 1.3 三個 crash/kill helper 不再獨立註冊，父測試仍可啟動子程序；跑實際 kill 與重啟驗證。
- [x] 1.4 SQL Server 關閉能力分支不產生 skip，仍斷言拒絕原因；跑 native gate 測試。

## 2. 驗證與交付

- [x] 2.1 確認本批只有測試與文件變動；執行 Go 完整測試、lint、vet 與行數檢查。
- [x] 2.2 完成 Spectra verify/review 與可封存提交範圍核對，不將刪除項目當作通過。
