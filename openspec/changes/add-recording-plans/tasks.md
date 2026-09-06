## 1. 契約與 Red 測試
- [ ] 1.1 凍結 plan/stream/destination、schema preview token 與 repository transaction 介面，供 3、4 共用。
- [ ] 1.2 Red：測試多設備多用途、每筆/變更/降採樣、心跳、品質變化與傳輸 batch 不丟樣本。
- [ ] 1.3 Red：測試 typed columns、uint64/decimal precision、NULL/0、late snapshot、batch trigger 去重與缺項。

## 2. 儲存與模式
- [ ] 2.1 新增安全 migration、plan repository、stream compiler 與 definition ledger，保留 legacy mapping。
- [ ] 2.2 Green：實作 raw/event/state/batch/latest 模式，以及供 aggregation 使用的 pre-decimation input。
- [ ] 2.3 實作 managed schema、索引、穩定 record key 與 receipt transaction 介面；進階 wide 輸出不另算值。
- [ ] 2.4 實作保留期限、計算證據期限與容量估算驗證；不讓設定默默縮短重算證據。

## 3. 真實資料庫準備
- [ ] 3.1 Red→Green：逐 kind 測試 capability、連線與真實 schema introspection；不得用 sample schema fallback。
- [ ] 3.2 實作 schema diff/token/confirm、碰撞與權限驗證，不在 runtime 自動建庫。
- [ ] 3.3 實作 credential reference、身份變更失效、SQLite 路徑邊界、SQL identifier allowlist/quoting。
- [ ] 3.4 實作 preview/test-write/readback 三段結果，去重及 test 資料排除，覆蓋有副作用的既有表。
- [ ] 3.5 更新新 API、types/services/hooks、Swagger 與 backend typed errors，不改既有 endpoint 的舊請求語意。

## 4. 完整驗收
- [ ] 4.1 SQLite/PostgreSQL/MySQL integration 測試相同契約；不支援者在 UI/後端同時阻擋。
- [ ] 4.2 Refactor 與既有單點/grouped/Local Modbus 回歸；所有測試、lint/build、migration 演練留證。
- [ ] 4.3 執行 strict spec 驗證並更新契約；新 managed activation 保持關閉，直到 4 的耐久交付驗收通過。
