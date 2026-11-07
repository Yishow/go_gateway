# 已確認的踩坑與經驗

每次重跑核對現況；歷史症狀不代表目前仍有該缺陷，舊測試數、埠與私人路徑不是固定條件。

## 2026-11-06：合法空 DB 版本擋住新增來源

- 症狀：MC／Modbus 顯示沒有設備連線設定、無法核對版本，填好範圍仍不能保存。
- 根因：來源 lease 把尚未設定資料庫時合法的 `workspace_revision: ''` 當成缺失，未交回已存在設備。
- 重跑：空工作區 UI 新增第一個範圍，及 memory 空格新增；真正保存、GET、F5。兩協議都做。
- 保護：undefined、純空白、跨工作區、讀取中版本改變仍拒絕。DB setup 與 Modbus Share revision 是不同 authority，不可互換。
- 歷史證據目錄名稱：`fix-add-range-20261106-7a01e7d5`；搜尋當下驗證目錄，不依賴固定絕對路徑。

## 2026-11-06：新增通過，不代表已保存範圍能啟用

- 症狀：已保存的離線範圍啟用後 PUT 400，只有資料無效。
- 根因：設備沒成功 connect／probe，後端正確擋 activation；前端缺 readiness 查詢與行動說明。
- 重跑：未測試設備先保存停用範圍，再編輯／啟用；未就緒不送 PUT，明確測試後才保存。離線埠測失敗保留草稿。
- 不刪 readiness、不假 ready、不用 API success 替身補過。
- 歷史證據：`fix-range-activation-home-refresh-20261106-7a01e7d5`。

## 2026-11-06：背景重查造成首頁重掛

- 原 `isFetching` 判斷暫時卸載 scoped home，看似整頁重載；React 重掛不等於文件導航。
- 重跑要保留原 DOM／filters，用 MutationObserver 看是否被移除，另計主 frame navigation。
- 隔離單次保存到表單約 117／120 ms，導航均 0；只證明消除重掛，沒重現使用者 3～5 秒，不能宣稱所有延遲已修。

## 2026-11-06：腳本假設造成誤判

- raw locator 不接受隨意加 `modbus_share`；它屬 unscoped workspace，fixture 應照真 DTO。
- 規則欄位是 `naming_prefix`，不是假設的 `name`；後端會正規化大小寫，核對契約及 ACK。
- 新增保存可能留在原表單的已保存狀態，不會出現既有編輯的 `home-saved-range-form`；驗真實狀態與 ID。
- 格子名稱可能在 aria-label／title，不是可見文字；不要拿錯文字斷言當產品失敗。
- 首頁設備摘要當時只列 4 筆；第 5 台到完整設備表找，不重建設備。
- DOM 存在／`toBeVisible` 不保證在 viewport；截圖／操作前確認提示與控制實際可見、未被覆蓋。
- Testing Library `getByRole` 沒有 Playwright 的 `exact` option；型別或 harness 錯誤不能算產品 RED。
- 設備表單先選協議再填名稱／連線欄位，且值已相同就不要重選。先填名稱再觸發協議選擇，會遇到正常的未送出變更確認；後續欄位 timeout 是腳本操作順序，不是設備保存 bug。
- 一輪 agent-browser batch fill 沒有取代原預設文字，曾保存錯誤的 PG connector 欄值；改成逐欄 Playwright `fill` 並立即核對 `inputValue`、保存後再讀 canonical GET。修正只涉及該自有測試 connector，不當成產品缺陷。

## 2026-11-06：真實數值與原圖是兩條驗收軸

- 五布局與 MC／Modbus 值變動／斷線復原曾通過；章節位置、地圖摘要及詳情區塊仍與原圖有差異。
- 0、負值、float32、float64、uint16 邊界、int32 負值由種值與 wire readback 核對；MC wire word bytes 與 simulator order 名稱不可直接畫等號。
- SAMPLE 密集格子只能驗版面，原圖假數值／計數也不是資料真相。
- 模擬器 register API 可能直接呼叫 `server.Read`，不能只憑名稱稱獨立 wire read；本輪另用 TCP socket 保存 MC 3E／Modbus request/response 才取得網路協議證據。來源碼與 binary 對應未證時也要明列。
- 恢復驗證實際出現 error／SSE `quality=bad`／格子 `!`，恢復種入 MC 321、Modbus 654 後，TCP、Gateway GET 與五視圖都對上新值。未單獨觸發 stale，不能將 error case 改稱 stale PASS；頂部 `Live` 代表串流狀態，不能單獨證明設備成功讀值。
- 歷史證據：`revalidation-20261106-7a01e7d5` 的 REPORT 與 reference checklist，不直接繼承成新版本 PASS。

## 歷史 first-recording 驗收的限制

- 2026-10-04 一輪 PostgreSQL fresh multi-device setup 曾有 `Tag binding missing`／9 blockers。這是舊次失敗，不能宣稱目前仍有該 bug。
- code、spec、runner 可能已變；重查 Point／Tag／mapping，不沿用舊 binary 或 archive 勾選作完成證據。
- API 預建 Step 1–3 再驗 SQL，只能代表後半段，不是 UI 新增設備到 SQL 全程。
- SQLite 不能替代 PostgreSQL／SQL Server／Windows／現場；正式 spec 的精確 uint64、多 entity 與 fresh 次數不能被單一 happy path 取代。

## 2026-11-06：skill 自身也需驗證

- 官方 `quick_validate.py` 在系統 Python 缺 PyYAML，這是工具環境錯誤。hc Terra 使用 owned 暫存 venv 安裝驗證依賴，未改全域 Python 或產品依賴。
- 原證據工具只看 Git diff／未追蹤內容，會漏掉 `assume-unchanged` 隱藏的追蹤檔變更、腳本執行權限及目錄 symlink 目標變更。Terra 先重現反例再補快照與測試，主代理另跑 11 項 CLI 回歸及官方檢查。
- 新快照不是原子讀取，也不證明 ignored／外部依賴或 binary／SPA 等同來源；仍需建置 manifest。舊格式不能直接繼承新版 PASS，見 evidence reference。
- 新的長命令名稱須符合當前 hc 長度限制；先核對 CLI，不把命名失敗當產品問題。

## 2026-11-06：完整 fresh 流程新增發現的 MC 準備卡住

- 症狀：新保存的停用 MC 範圍，明確連線／probe 成功且 readiness ready，仍停在準備中，沒有送出啟用 PUT；Modbus 同路徑可完成。
- 根因：`useRecordingRangeSubmission.prepareSaved` 在排入規則更新前清掉剛核對的 MC proof，原本的 queue gate 因缺少相同 scoped revision 的證據拒絕寫入，promise 因而 pending。
- 修正：保留重新核對的 `after` proof 到 queue ACK，不移除 identity／revision／MC ownership gate。hc Terra 在真實頁面呼叫鏈補 RED；主代理另跑 strict ACK／lease；新 binary 的 fresh MC／Modbus 建立、準備、編輯與重載已完成。
- 與先前 400 不同：這次是新增表單的 `Check connection and prepare points`，不是既有範圍標題上的啟用開關。兩條路都要驗；不能以前一條通過推論另一條。
- 來源準備成功不等於已開始採集；當下 UI 明示不啟動 collection，真實數值可在記錄啟動後回驗。此條只記來源階段，不表示完整 UI→SQL 已通過。

## 2026-11-06：mixed 試寫能力限制與正常記錄分開

- 預設簡化 UI 建立 mixed、latest／15 秒、受管寬表。明確 schema apply 成功後，試寫 preview 422 `test-ownership-unsupported`；因保存的 row policy 沒有現行試寫 codec 要求的 entity key。
- PG 初次 start review 曾因就緒證據改變而停用確認，依 UI 關閉／重載後重新 review 才成功；這是保留中的 stale scope 保護，不改掉 guard 或偽造 ready。
- hc Terra 唯讀追到 `grouptestwrite.buildPlan` 與正式 mixed runtime 的不同 codec。規格要求無法證明 operation ownership 時拒絕，未找到本案 mixed 必須試寫成功的契約；不為清單變綠新增能力或移除 guard。
- 簡化 mixed UI 沒有可填的 entity key；改 long 或硬補欄位不是已證解法。進階非 mixed 路徑不同，未實測就不能當本案成功。
- 正常 SQLite 記錄另已觀察到三完整 15 秒 bucket、六點與雙 receipt；試寫執行、讀回及 cleanup 仍未執行。對正常記錄與試寫分別報告，不互相代替。
- 群組 receipt 要查 `wg_delivery_*`，不可混用 legacy `gw_delivery_receipts`；`sql_committed` 也不可誤斷言為 `committed`。這些是 harness 假設錯誤，不是產品交付失敗。
- UI 停止會更新 saved revision，不能拿停用後的 current saved revision 強制等同歷史記錄版本。hc Terra 重現父代理 verifier 的這個取時錯誤後，改核對 `wg_runtime_versions`、`write_group_versions`、當時 catalog、outbox 與雙 receipt；主代理直接唯讀重跑通過。啟動截斷窗也不算三個完整 bucket。
- 停用 ACK 後仍可能觀察到後續 SQL commit；要分開 sample observed／received、runtime accepted、outbox created 與 committed 時間。某輪 PG 停用後最後一列 due 晚於 ACK，數值的來源 observed 卻早於 ACK；缺 accepted timestamp 時，不能逕稱全部是 pre-ACK backlog，也不能只憑後續筆數穩定宣稱即時停止 intake。列明已驗的穩定觀察期間與未證範圍。

## 追加格式

只記日期／版本、症狀、分類（產品／腳本／環境）、已證根因、最小重現、修正重驗及證據位置。保留失敗嘗試，不抄 secrets 或未證猜測。
