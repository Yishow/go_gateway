## 1. 必填問題與既有功能

- [x] 1.1 保留既有驗證順序並合併必填問題：以 Modbus／FATEK mixed missing+empty 的 RED 重現 Aggregated device connection validation 漏報，再最小修復；既有單類、nil、有效設定、range／enum、條件欄位以 focused Go tests 與 mutation 檢查回歸。
- [x] 1.2 README 與計畫只描述可從 wiring/crash test 證明的 accepted journal／unknown 恢復，現場成功保留為歷史來源、不含私有位址或執行識別；以 source review、敏感資訊檢查及文件連結檢查驗證。
- [x] 1.3 [after: 1.1, 1.2] 完整 Go test/vet/lint、review、analyze/validate、line/diff gate 驗證同一最終 source，記錄真實結果與限制；無 Critical 後完成封存準備與 preview；tasks 依真實結果完成後，由外層授權執行 archive／提交。
