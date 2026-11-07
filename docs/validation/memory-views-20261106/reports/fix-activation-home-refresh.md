> 歷史報告閱覽副本：正規化本機絕對路徑與附件連結；保留當時結論與未驗項目。原始位元組及完整附件見 `local-only/original-evidence.tar.gz`，不要將歷史 PASS 當成新版重驗。

# 已保存範圍啟用 400 與首頁重掛修正

結果：本輪修正與驗證 PASS，變更尚未提交。基準 HEAD `7a01e7d5`，分支 `feat/memory-views-validation-pending-20261105`。

## 問題與修正

使用者的來源規則已以停用狀態保存，但設備沒有成功的連線／讀取測試。直接切換啟用再儲存，會被後端 activation readiness 保護拒絕，回傳 400。已在全新 MC 測試設備的真實 UI 重現相同情境與回應，見 `enable-red.json`。

保留後端保護。已保存範圍由停用改成啟用時，前端先查就緒狀態；未就緒會說明原因，提供「測試連線並儲存」。只有使用者明確點擊才執行連線測試。測試成功後再核對就緒狀態、工作區、來源修訂與連線設定，才交給原有儲存佇列；失敗保留未送出的編輯，不送啟用 PUT。

保存連線會觸發背景重新查詢；首頁原先把 `isFetching` 當成沒有有效工作區，暫時卸載再重掛整個首頁。修正為同一個已核對的工作區／版本在重查期間繼續掛載，保留篩選狀態。缺失資料、讀取失敗、跨工作區或版本變動仍依原有判斷處理。

## 本輪修改檔案

- `frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeSavedRangeEditor.tsx`
- `frontend/src/features/datalink/workbench-v2/recording-home/RecordingHomeRoute.tsx`
- `frontend/src/i18n/locales/en/recording-home.json`
- `frontend/src/i18n/locales/zh-TW/recording-home.json`
- `frontend/tests/unit/workbench-v2/recording-home-source-modal.test.tsx`
- `frontend/tests/unit/workbench-v2/recording-home-route.test.tsx`

前一輪空版本修正的 3 個檔案逐一雜湊比對，保持完全相同；目前工作樹共 9 個變更檔案。

規格脈絡：[simple-recording-home](../../../../openspec/specs/simple-recording-home/spec.md)。本輪為既有流程的獨立修錯，沒有修改規格或勾選其他 change 的任務。

## 驗證結果

| 驗證 | 結果 |
| --- | --- |
| 修正前真實 MC 啟用 PUT | RED：400，明確回報 activation readiness 未通過 |
| 修正前啟用引導回歸 | RED：沒有引導且錯送更新 |
| 修正前首頁保持掛載回歸 | RED：同一工作區重查會移除原首頁 DOM |
| 前端完整測試 | PASS：643 檔案、7,389 測試 |
| 最後資料存在保護調整後的直接回歸 | PASS：5 檔案、72 測試 |
| 完整 lint／最終變更檔案 eslint | PASS |
| 型別檢查與前端正式建置 | PASS：`npm run build` |
| 差異空白與行數檢查 | PASS：9 檔；2 份 locale 376 行、測試檔 408 行為警告，沒有超過 500 行 |
| 真實 MC＋HSL | PASS：先顯示就緒引導且沒有無效 PUT，明確測試後保存成功，重載仍啟用 |
| 真實 Modbus＋HSL | PASS：同上 |
| 離線設備 | PASS：真實連線測試失敗，保留編輯與停用的保存值，沒有 PUT |
| 保存連線後首頁 DOM | PASS：沒有替換，沒有瀏覽器文件導航 |

所有實測使用專用 Chromium、嵌入式 Gateway、可丟棄 SQLite 副本與 localhost HSL；沒有攔截／偽造 API，也沒有連真 PLC。

MC／Modbus 的成功流程記錄於 `fixed-browser-results.json`。同一腳本後續要開離線設備時，誤以首頁僅列 4 筆的摘要能找到所有設備而中止；離線案例改由完整設備表開啟，獨立完成於 `offline-browser-results.json`。初版另曾把後端正規化的大寫前綴誤當小寫；這些腳本問題不列為產品 RED。

## 等待時間的實際觀察

| 隔離環境單次量測 | 修正前 | 修正後 |
| --- | --- | --- |
| 保存連線到點位表單 | 117 ms | 120 ms |
| 設備建立 HTTP 請求 | 15 ms | 8 ms |
| 原首頁被移除 | 是 | 否 |
| 瀏覽器文件導航 | 0 次 | 0 次 |

這是單次本機模擬器／可丟棄設定量測，不是效能基準。使用者提到的 3～5 秒未在此環境重現；已確認並消除的是首頁不必要的重掛，不能宣稱所有延遲已解決。真正連線／讀取測試仍可能等到設備逾時。

## 截圖與證據

- [MC 測試後保存](../images/fix-range-activation-home-refresh-20261106-7a01e7d5/mc-enabled-saved.png)、[MC 重載](../images/fix-range-activation-home-refresh-20261106-7a01e7d5/mc-enabled-reloaded.png)
- [Modbus 測試後保存](../images/fix-range-activation-home-refresh-20261106-7a01e7d5/modbus-enabled-saved.png)、[Modbus 重載](../images/fix-range-activation-home-refresh-20261106-7a01e7d5/modbus-enabled-reloaded.png)
- [離線測試失敗時的清楚提示與操作](../images/fix-range-activation-home-refresh-20261106-7a01e7d5/offline-failed-test.png)
- `create-before.json`、`create-after.json`：請求耗時、首頁 DOM 與文件導航量測。
- `manifest.json`、`fix.patch`：最終檔案雜湊、執行檔雜湊及差異；已比對目前程式與建置來源相同。
- `editfix-*.log`：執行的回歸、lint、建置與瀏覽器證據。

## 限制與交接

未修改使用者的工作資料庫，未 commit／push／merge／部署。測試服務與模擬器於驗證後關閉。開發服務重新整理頁面即可重新載入表單；舊嵌入式執行檔需重新建置及重啟。

本輪沒有重跑五種記憶體布局的完整點對點視覺／數值矩陣，也沒有做 Windows、真 PLC、正式資料庫或現場驗收。
