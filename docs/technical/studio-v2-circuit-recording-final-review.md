# 迴路記錄主 change 收尾審查

對應 `add-circuit-recording-profiles`，2026-10-25；僅收尾本 change。
MySQL／MariaDB 與 SQL Server 兩批保持 parked，未實作、未驗收。

## 規格與程式審查

- verify：15/15 任務、17 個 requirements、49 個 scenarios；逐列對照見[情境覆蓋](evidence/circuit-recording-scenario-coverage.md)。
- review：主代理核對 workspace、pipeline、production wiring、API 與關鍵前端；唯讀代理補查 counter／delivery／schema 與前端，發現均由主代理回讀原碼及規格判斷。
- analyze：Coverage、Consistency、Gaps、Localization 結構檢查 Clean；37 個 Suggestion 為未另外設 `##### Example` 區塊，不是缺少對應測試或阻擋性矛盾。
- 語意核對將 proposal／design／tasks 中「本輪只提案／尚未執行」標為提案階段歷史，避免與目前完成狀態混淆；沒有改寫 requirement 語意。
- 早期測試含實作後補寫及 mutation 驗證，不把它們回寫成事前失敗歷程。

## 確認並修正的兩個缺口

1. 已套用寬表保存新草稿後，`status=draft`，原 applied revision 仍記錄；ownership 查詢原本只讀 ready/running，普通長表 Apply 可能繞過 `cutover_required`。查詢現在包含有 applied revision 的 draft，仍排除 disabled/deleted。
2. 長表 INSERT 原本省略 owner 欄而依賴 DEFAULT；布局證明不要求 default，缺少 default 的相同布局會在寫入時失敗。長列現在明確寫入既有 owner marker 的 `gateway` 值，未更動 DDL、migration、receipt 或舊寬表 writer。

重現與回歸：

- `TestProfileApplyRefusesAnAppliedWideGroupWithANewerDraft`：原程式 SQLite 回 200，修後 SQLite／PostgreSQL 都回 409／`cutover_required`；無新增長表 version、無 cutover／revision-end，舊草稿與 applied 身分保持，舊 writer 真實送達 SQL。
- `TestLongRowWritesItsOwnerWithoutAColumnDefault`：真 SQLite／PostgreSQL 表只移除長表 owner 的 default，保留 receipt default；停用修正時兩庫皆因 owner NOT NULL 失敗，修後寫入並讀回 owner=`gateway`、raw=`9007199254740993`。
- 首版 owner 測試誤將 receipt default 一併移除，屬 fixture 錯誤；修正 fixture 後重新確認 owner 失敗原因，未把該錯誤當成另一個產品缺陷。
- raw 的 MaxAge 選樣與完整 pre-decimation reset 證據不同，沒有因窗內較舊讀值而抹去 reset；首窗維持 raw＋NULL delta＋baseline 的既定契約。

## 修正後實際驗證

環境：Go 1.27.3（go.mod 1.25.5）、modernc SQLite、PostgreSQL 17.11／TimescaleDB 2.31.0；沿用既有 loopback 可丟棄 fixture。

| 範圍 | 結果 |
|---|---|
| dbtarget 全套 | 324 PASS 路徑，0 SKIP |
| workspace 全套 | 361 PASS 路徑，0 SKIP |
| runtime 全套 | 158 PASS 路徑，0 SKIP |
| grouppipeline 全套 | 100 PASS 路徑，0 SKIP |
| groupdelivery 全套 | 214 PASS 路徑；1 個 crash child helper SKIP |
| API 受影響整合 | 88 PASS 路徑，0 SKIP |
| 合計 | 1,245 PASS 路徑，0 FAIL；路徑數含父層，不是葉測試數 |
| go vet | dbtarget／workspace／runtime／grouppipeline／groupdelivery／api 通過 |
| go build ./... | 通過 |
| golangci-lint | dbtarget／workspace 無新問題；api 仍只有範圍外 router_logger.go:109 gocritic |
| git diff --check／行數檢查 | 工作樹通過；保留歷史長檔警告 |

API 選測包含實際被選中的：

```text
TestCircuitRecordingIntegration
TestCircuitProfilesPartialPolicyUsesTheLongLayout
TestProfileApplyRefusesAnAppliedWideGroupWithANewerDraft
TestPartialSourceCutover / TestProfileCutoverFailure
TestCumulativeLongGroupRolloverFromSavedPolicy
TestLongGroupAppliedByTheServiceRunsInTheProductionPipeline
TestCumulativeLongGroupCountsFromASavedMeasurement
TestLongStorageCumulativeRecovery
TestLongStorageMissingRecoveryKeepsFrozenRows
TestLongStorageMissingRecovery
TestCutoverSourceEvidenceReportsTheActualWindow
```

上述 API 測試包含 SQLite、PostgreSQL plain 與實際選中的 Timescale 子案例；未以安裝 extension 推定所有案例皆跑 Timescale。
完成後 `gw_*` 測試 schema 為 0；既有 fixture 容器維持 running，未操作正式資料庫。

## 證據與限制

- 修正前的完整結果仍保留：[前端 260 檔／1,812 項及 Go 11 套件／2,476 PASS 路徑](evidence/circuit-recording-test-results-20261024.json)；本輪 Go 修正後的受影響重測以上表為準，不冒稱重新跑過全部套件。
- 本輪沒有修改前端，重用已通過的 lint／build／1,812 測試與[真實瀏覽器證據](studio-v2-circuit-recording-browser-verification.md)，未重新宣稱一輪瀏覽器操作。
- [DB 矩陣](studio-v2-circuit-recording-db-verification.md)的 PG18 無 extension／Timescale 實測與 skip 配對保持獨立；本輪沒有新起 PG18 fixture。
- 已知 legacy recording-plan PostgreSQL `verification_unavailable`、storage partition 基線及 router_logger lint 未轉成此次修復範圍。
- 非正式設備、Windows、LAN 或現場驗收；未 archive、push 或部署。提交另依完整檔案與差異計畫確認，保留其他 WIP 與暫存內容。
