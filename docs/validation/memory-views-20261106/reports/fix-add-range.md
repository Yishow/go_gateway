> 歷史報告閱覽副本：正規化本機絕對路徑與附件連結；保留當時結論與未驗項目。原始位元組及完整附件見 `local-only/original-evidence.tar.gz`，不要將歷史 PASS 當成新版重驗。

# 新增點位被空資料庫版本阻擋：修正驗證

結果：PASS。修正位於 `feat/memory-views-validation-pending-20261105`，基準 HEAD `7a01e7d5` 之上的未提交變更。

## 原因與最小修正

尚未設定資料庫時，write-groups 與 locator 的 database setup revision 合法地回傳空字串。來源設備的穩定讀取檢查誤把它視為缺失，導致 MC 與 Modbus 新增範圍拿不到設備設定、保存按鈕停用。

僅允許明確的空字串通過版本格式檢查；缺失、空白、工作區不符、前後版本不一致、設備清單不一致、過期請求仍拒絕。Modbus Share 來源版本的保護與群組保存 ACK/CAS 保護不變。

## 修改檔案

- `frontend/src/pages/datalink/workbench-v2/useLocatorDeviceSourceLease.ts`：修正兩處空版本判斷及註解。
- `frontend/tests/unit/workbench-v2/recording-home-source-entry.test.tsx`：MC／Modbus 未設定資料庫仍可準備新增範圍。
- `frontend/tests/unit/workbench-v2/recording-editor-source-lease.test.tsx`：穩定空版本、讀取中建立資料庫、缺失／空白版本的回歸案例。

規格脈絡：[simple-recording-home](../../../../openspec/specs/simple-recording-home/spec.md)。沒有改動規格語意或建立新工作流程。

## 實際證據

- RED：未修改產品程式時，MC／Modbus 兩案例均在保存按鈕應啟用的斷言失敗，見 `red3.log`。
- GREEN：最終四組直接回歸 52／52；後續測試資料缺失值更正後，版本保護 17／17。
- 前端完整測試：643 個檔案、7,381 個測試 PASS。最終測試輔助資料更正後另外重跑受影響測試。
- `npm run lint` PASS；最後修改的檔案 eslint PASS；`npm run build` PASS；最後測試資料更正後 `npx tsc --noEmit` PASS。
- `git diff --check`、`bash scripts/check_file_lines.sh` PASS（3 個檔案）。
- 真實嵌入式 Gateway + 專用 Chromium + 可丟棄 SQLite 副本，沒有攔截或偽造 API。資料庫設定版本在保存前後均為空字串。
- MC D2、Modbus 40003 各新增 1 個 int16 範圍：source-rule 與 mapping 均收到 201 ACK，canonical GET 讀回相同設定，重新載入網頁後格子仍為 used，名稱與位址吻合。原有範圍資料逐筆相同。
- 瀏覽器 JavaScript 錯誤：0。

| 實測 | MC | Modbus |
| --- | --- | --- |
| 保存前 | [畫面](../images/fix-add-range-20261106-7a01e7d5/mc-before-save.png) | [畫面](../images/fix-add-range-20261106-7a01e7d5/modbus-before-save.png) |
| 保存後 | [畫面](../images/fix-add-range-20261106-7a01e7d5/mc-saved.png) | [畫面](../images/fix-add-range-20261106-7a01e7d5/modbus-saved.png) |
| 重載後 | [畫面](../images/fix-add-range-20261106-7a01e7d5/mc-reloaded.png) | [畫面](../images/fix-add-range-20261106-7a01e7d5/modbus-reloaded.png) |

詳細 request／ACK／讀回：`browser-results.json`；重現腳本：`verify-add.cjs`；最終差異與雜湊：`fix.patch`、`manifest.json`。

驗證腳本初版曾使用錯誤的 API 名稱欄位、假定新增完成會切換至另一個表單，以及將名稱誤當為格子可見文字；已依真實 API 與畫面修正。這些腳本失敗不列為產品 RED。最後完整執行已通過。

## 限制與交接

本輪沒有啟動 HSL 模擬器，未重驗真實協議讀值／即時更新／斷線復原；畫面的讀值失敗不列為本輪 PASS。Windows、真 PLC、正式資料庫與部署為 NOT RUN。

只寫入本輪隔離副本，未修改使用者工作資料庫；未 commit／push／merge。嵌入式舊執行檔必須重新建置及重啟才會載入修正。此報告不代表先前五種視圖的視覺差異已修正。
