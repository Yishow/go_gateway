## 1. 正式 Sender 的期限回歸

- [ ] 1.1 補 RED：遠端合法耗時大於 settlement budget、但小於 attempt timeout 時，成功、暫時失敗與永久拒絕仍應各自正確保存結果；使用 barrier 避免脆弱 sleep。
- [ ] 1.2 在遠端返回後才建立獨立有界 settlement context，覆蓋所有分支並讓 1.1 GREEN；同步 sender 的生命週期註解。
- [ ] 1.3 補 caller 在遠端 commit 後取消及 settlement 真實故障測試，確認前者可記帳、後者保留恢復證據且不盲目重送。

## 2. Lease 與關閉保護

- [ ] 2.1 驗證 lease 涵蓋 attempt 加 settlement，失效 owner 仍被 fenced，shutdown 不超出既有界線；必要修正只限期限銜接。
- [ ] 2.2 在 disposable SQLite／PostgreSQL 重跑慢寫、commit 回覆遺失與 local receipt 失敗，逐一比對 outbox／SQL／effect key；更新期限與限制文件。

## 3. 整合驗證

- [ ] 3.1 執行 sender／worker focused tests 及 repository 後端最低檢查，保留命令與未執行項目；不得把隔離範例當成專案整合通過。
- [ ] 3.2 review 確認未改預設重試、去重或狀態語意，修復本案發現後重跑對應測試。
