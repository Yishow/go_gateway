## 1. 頁面確認與整合驗收

- [x] 1.1 在前三change規格審查後、任何產品apply前完成 Default adoption follows explicit acceptance gates 的「使用者斷點與失敗回退」確認清單；提供此次修訂固定demo版本、協議參數／75%modal與RWD／亮暗主題情境及待決定例外，記錄owner明確回覆，不把沉默或舊waiver當確認。
- [ ] 1.2 [after: 1.1] 在home/editor/lifecycle實作完成後依「新版設計與能力驗收」實作 Revised design and protocol parity precede acceptance，並建立 Acceptance records distinguish evidence layers 的「驗收矩陣與分層證據」；為每個既定scenario及協議欄位parity、system/liveOS/override/persistence/portal/contrast/reduced-motion對應測試與真實層級，依「混合member計算前置與整體驗收」納入 Mixed member acceptance preserves one-row and legacy contracts，並依「同源採集次數預覽交接」驗收 Adoption distinguishes planned and actual acquisition counts，測估算／actual分層、stale409、零I/O及unsupported原因，實際性能證據由batching owner提供，不用demo冒充；內容review確認無漏列或以mock冒充SQL/PLC。
- [ ] 1.3 [after: 1.2] 以現有Playwright pattern補齊首次建立、已有來源、草稿、正常、異常與修改/停恢/複製/表格變更journeys；驗證Close/Escape/Back/Forward/double click/scope切換、焦點及窄畫面；依新版設計與能力驗收覆蓋唯一Select Points、filter即時改同表／計數、加法pending與Add N、取消保非法輸入、remove/undo/re-add、空設備續建、無默選首設備、分步保存後取消不回滾來源及同scope返回位置；不能執行者明列NOT RUN與影響。

## 2. 預設採用與相容

- [ ] 2.1 [after: 1.3] 以前置gate與owner確認為條件實作「預設入口採用與明確意圖優先」；以Red-Green-Refactor測試 Explicit existing intents survive default adoption，鎖定locator、runtime focus、step/target/section/modal、search/hash及invalid locator fail-closed；先RED證明現行step=output&target=modbus被誤判locator，再修bounded legacy辨識，測invalid locator混legacy不得降級。
- [ ] 2.2 [after: 2.1] 實作 Entry rollback preserves configuration and operations，保留view=advanced、runtime/logs/test及原default回退；測試切入口零mutation、pending operation與既有identity不變，不恢復/studio。

## 3. 最終驗證與交付

- [ ] 3.1 [after: 2.2] 對最終功能批次執行npm run lint、npm test -- --run、npm run build及現有embedded static交付檢查；有Go改動才執行涉及側完整test/vet/lint，記實際PASS/FAIL/NOT RUN，不新增環境或接真PLC/正式DB。
- [ ] 3.2 [after: 3.1] 獨立review安全與完整生命週期、修本批問題並重跑受影響gate；必要驗收未過且未獲本批明確接受則保持opt-in，驗證報告與實際default一致。
- [ ] 3.3 [after: 3.2] 完成「唯一交接清單與最終批次覆核」的docs/releases/simple-recording-entry.md，列四change完成狀態、owner確認版本、測試命令/證據/風險及rollback邊界；跑diff/line gate與Spectra verify/review，未授權不push或部署。

相關modal／導航任務依「已選 modal 與返回引導契約」回歸一般／手機／小確認分類、固定操作區、contexthelp風險可見性及同scope返回狀態，不新增未選功能。


本輪已確認逐member interval與UTC epoch對齊：5秒00/05…55、15秒00/15/30/45、10分鐘00/10…50；wide採due邊界聯集，:15/:30/:45的60秒delta留NULL＋not_due，:00同列各自前15秒mean／前60秒delta；long只新增到期member列。不carry、不補0；required-not-due不skip，到期失敗／無baseline／gap/reset分別標原因。顯示各member的[start,end)、partial coverage與next due，SQL commit時間另列；first partial、nullable schema、安全gate與durable驗收依support-per-member-recording-calculation。此語意確認不是新版demo或apply核准。
