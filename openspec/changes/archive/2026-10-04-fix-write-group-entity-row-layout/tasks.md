## 1. 欄位範圍與 readiness

- [x] 1.1 補 RED：同群組兩 entity 共用 value 欄位合法，同 entity 重複欄位非法；比對 frontend、readiness、Apply 的相同判定。
- [x] 1.2 修正 entity-scoped layout 與共用驗證，讓 1.1 GREEN；同步說明列身分，不改舊 group/member key。

## 2. 實際資料列與錯誤保存

- [x] 2.1 補 RED：同群組 A/B 同欄及不同欄各自 EncodeRow；A missing 而 B good、partial 與 no_data 不互相污染。
- [x] 2.2 修正每列 member 選取與錯誤分類，讓 2.1 GREEN；結構性 layout 錯誤保留 accepted input／checkpoint，不冒充正常 skipped。
- [x] 2.3 在真 SQLite／PostgreSQL 驗證同群組兩 entity、record identity 與每列內容，並覆蓋重啟後 frozen layout；更新相容性與限制文件。

## 3. 整合驗證

- [x] 3.1 執行受影響前後端套件與 repository 最低檢查，確認既有單 entity／custom 群組沒有回歸。
- [x] 3.2 review 對照所有 scenarios；合法多 entity 必須修好並實測，不以一律拒絕或拆成兩群組當作完成。
