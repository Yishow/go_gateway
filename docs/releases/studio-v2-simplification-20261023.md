# Studio V2 新手操作簡化驗證（2026-10-23）

基準：`33f6a1d10218aa0e577dd51aedfbad076c96c8bf`。在原有 Studio V2 實作，沒有建立另一套前端。原工作樹的 12 個未提交／新增檔案未納入本次提交。

## 完成的行為

1. 點格單擊、Shift 區段及 Ctrl／⌘ 多選只改選取；明確「採集」勾選與既有選取後批次操作才改採集。
2. 映射批次操作明示來源，逐列勾選對象，先呈現欄位前後差異，再確認。全部啟用列仍可明確勾選；未勾選與停用列不被修改。來源／對象值改變會撤銷舊預覽。
3. Runtime 無 failure 只表示未偵測到錯誤，不推論全入庫，不補造 `database_delivery` 成功階段。SQL 成功沿用基本記錄的群組、applied revision、connector revision 與 committed effect 證據。
4. Step 4 先呈現 DB 目的地；尚未儲存時不呈現記錄項目／頻率／開始。已有連線以摘要與修改入口呈現。收合只隱藏，保留元件與密碼／設定狀態；Share-only 保留。
5. 基本建表一次準備並自動取得唯讀預覽，DDL 仍必須明確確認原 token／operation。回覆遺失與 unknown 保留原 ID，重整、兩張同群組卡片均不能重送。已確認 rollback 的 failed 可修好後重新預覽。prepare_schema 指向基本建表卡。
6. 基本受管群組提供新增／移除項目差異確認。既有欄位、required、max_age、量測歸屬、row/write policy 保留；自訂、entity、特殊去重／來源配置仍由進階頁管理。workspace／group／connector CAS 與後端來源版本驗證保留。
7. 映射預設五欄：來源／位址、名稱、單位、目前結果、啟用。進階展開保留 Tag Key、所有型別、換算、中間值與 payload；切換不重建資料。
8. 四步為「連接設備 → 選要讀的資料 → 命名與換算 → 儲存並開始記錄」。讀取範圍有可修改的名稱建議，型別持續可見；4 個 float32 明示占 8 words，bit 區使用 bits。
9. Retry 是主要動作；「不再補送此筆」需展開影響、填原因與確認。保留原 payload、state revision、decision_id 與逐筆稽核；unknown 沒有 Retry／Skip。

## 實測發現及修復

新增點位的真實瀏覽器測試發現原受管 schema 僅允許完全相同的欄位集合。這使前端雖能保存新增項目，既有 DDL 預覽卻拒絕。此次在同一 token／確認／ledger 路徑補上 `add_columns`：

- 必須驗證原 group owner marker 與所有結構欄位，僅新增後端產生的 `v_<24 hex>` 可空值資料欄位。
- 不變更既有型別、required 或主鍵，不 DROP／重建表；移除成員後保留可空值歷史欄位與歷史列。
- receipt 表仍採原完整嚴格契約；確認前重新檢查目標、精確比對 DDL，目標改變即拒絕舊 token。
- SQLite 真實 DDL／SQL 驗證通過；PostgreSQL 產生精確 `NUMERIC(20,0)` 的新增欄位由單元契約驗證，未宣稱有本輪 PostgreSQL 瀏覽器驗收。

整體自行 review 另修復：查回成功操作後自動預覽蓋掉紀錄、同群組另一張建表卡重送、來源編輯與批次對象混淆、確認後來源變更、保存成員時丟失尚未修改的基本頻率。沒有執行獨立 agent review。

## 實際驗證

環境：Mac Studio、macOS arm64，Go 1.27.3（repo 目標仍為 1.25.5）、repo lockfile 的 Node／React／Vite 與 Playwright Chromium。以下是本次執行結果：

| 檢查 | 結果 |
| --- | --- |
| `cd frontend && npm run lint` | PASS |
| `cd frontend && npm test -- --run` | 210 檔、1331 項 PASS |
| `cd frontend && npm run build`（含 tsc） | PASS |
| `cd frontend && npm run test:e2e` | 14 項 PASS |
| `node --test scripts/tests/f_device_to_sql/*.test.mjs` | 62 項 PASS |
| `go test -p 1 ./...` | PASS |
| `go vet ./...` | PASS |
| `golangci-lint run ./...` | 0 issues |
| `git diff --check`、`make check-lines` | PASS；歷史超長檔不增行 |

完整 Go 回歸後只有 NUMERIC 字串常數抽取，另重跑 recordingplan tests 與完整 lint；不影響行為。最後新增的前端 rollback 修復已含於完整 1331 項回歸及最終 release 瀏覽器流程。

真瀏覽器使用 production embedded UI、loopback Modbus simulator 與本次新建且可丟棄的 SQLite；HTTP 設定由 UI 操作，SQL 查詢獨立唯讀。故障注入只在自有 SQLite 建立拒收 trigger，不以直接 INSERT 代替採集。

- 雙設備同位址、每台九個混合型別點位、兩個 Basic 群組，各三個連續 60 秒 bucket；SQL、receipt、provenance 比對通過。uint64 包含 `9007199254740993`、`9223372036854775808`、`18446744073709551615`。
- 最終 release 單設備：點格只選取、批次差異取消、基本／進階保值；真 DDL 執行後丟棄回覆，重整只查原 operation，apply POST 恰為一次。
- signed int64 SQL 查回 `9223372036854775807` 與 `-9223372036854775808`，無 JS number 捨入。
- 返回新增點位，窄範圍確認由 2 項增至 3 項，原欄位／policy 不變；確認新增欄位後 SQL 出現新值 `777`，原歷史列仍存在。
- 真實拒收列修好後 Retry 寫入 1 筆；Skip 須展開、原因與確認，目的地該筆為 0。兩者 payload identity 不變，各有唯一 decision_id 稽核。
- 雙設備在 390／768／1440 寬度無 document 水平溢出。自有 browser、gateway、simulator 均正常結束，暫存 DB namespace 已清理。
- 連點、保存中切頁、取消、過期 scope、unknown 禁重送與 Share-only 等其餘邊界由本次完整單元／整合回歸覆蓋；不是每項都做了獨立真人操作。

## 證據

[摘要與 SHA256](../plans/studio-v2-flow-completion/evidence-f/simplify-summary-20261023.json) 包含 binary／原始 witness SHA、SQL 結果與清理證明。[來源檔案指紋](../plans/studio-v2-flow-completion/evidence-f/simplify-source-inputs-20261023.json) 對應本次程式輸入。

- [最終 release 真瀏覽器／SQL／Retry／Skip witness](../plans/studio-v2-flow-completion/evidence-f/simplify-release-20261023.json)
- [完整雙設備原始 witness（gzip JSON）](../plans/studio-v2-flow-completion/evidence-f/simplify-dual-final-20261023.json.gz)
- [五欄映射畫面](../plans/studio-v2-flow-completion/evidence-f/simplify-release-20261023-mapping.png)
- [記錄與寫入證據畫面](../plans/studio-v2-flow-completion/evidence-f/simplify-release-20261023.png)
- [雙設備手機寬度](../plans/studio-v2-flow-completion/evidence-f/simplify-dual-final-20261023-390.png)

原始雙設備 witness 的 binary 早於純常數抽取及已確認 rollback 之 UI 重新預覽修復；最終 release 流程使用包含這兩項的 release binary。前幾輪失敗的原因均已記錄於上方並重驗，未把失敗當成完成。

未執行真人新手驗收、現場 PLC、正式 DB、Windows／Linux 或部署；沒有改系統名稱或重啟。Vite 仍有既有 chunk size／Browserslist 提示，未以額外依賴升級擴大本次範圍。GitHub 該 commit 的 CI 狀態在 push 後另行查證，不由本機 PASS 推論。
