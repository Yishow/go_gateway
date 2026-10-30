## 1. 已保存記錄的操作

- [ ] 1.1 在共同owner demo確認gate及 streamline-recording-task-editor 完成後，以 Red-Green-Refactor 為 Recording actions follow matching server availability 建立測試；鎖定「操作選單依後端事實而非狀態猜測」及stale scope禁止命令。
- [ ] 1.2 [after: 1.1] 整合停止與移除的影響確認及安全選單，保留資料表／歷史、source引用保護；用simple-recording-lifecycle及既有writeGroupLifecycle測試驗證不誤刪。
- [ ] 1.3 [after: 1.2] 實作 Resume revalidates and copying creates only a new draft 與「恢復需重新檢查且複製只產生草稿」的恢復流程；測試停止後以fresh evidence確認start，結果先waiting而非success。
- [ ] 1.4 [after: 1.3] 加入可表示配置的複製草稿與不支援原因；測試原group不變、無舊ID/operation/readiness/receipt/proof、managed ensure去重不假稱副本及不自動start。

## 2. 表格變更與復原

- [ ] 2.1 [after: 1.4] 實作 Table changes use reviewed forward-only successors 與「前向表格變更沿用既有安全階段」的影響摘要；測試新表接續、舊歷史保留、wide/long可用及受阻，無原地破壞轉換。
- [ ] 2.2 [after: 2.1] 串接既有schema/review/start-or-cutover呈現；測試first-start與cutover不同、boundary顯示、輸入改動後consent失效及pending不能當effective。
- [ ] 2.3 [after: 2.2] 實作 Unknown operations recover without replacement commands 與「異常修復保留未決結果與歷史」；測試關閉重開、Escape、Back/Forward、double click、切scope及unknown先查原operation；fresh匹配proof滿足原guard後保留explicit retryExact/repreview正向測試，無proof維持禁止重送負向測試。

## 3. 結果及修復指引

- [ ] 3.1 [after: 2.3] 實作 Results and repair paths are actionable without invented values 與「細節按需展示並保留進階交接」；測試規則/receipt摘要、saved/effective差異、無raw值不造假與advanced-only唯讀。
- [ ] 3.2 [after: 3.1] 對partial source、stale、DB failure、unknown加入帶scope修復/日誌入口，保留安全文案與返回位置；用locator/runtime/log相關測試驗證正確對象、舊成功不當目前正常。

## 4. 驗證及交接

- [ ] 4.1 [after: 3.2] 執行本批lifecycle/successor/recovery/locator相關unit及integration、lint/typecheck與line gate；對browser中斷/返回流程記實際結果，mock與jsdom不作SQL/視覺驗收。
- [ ] 4.2 [after: 4.1] 獨立review操作scope、複製authority清除、前向歷史保留與unknown recovery，修復本批問題並重跑選測；將結果交給adopt-simple-recording-entry的整體驗收。

設計驗收補充：既有操作相關任務須同時驗證 DESIGN.md 的light/dark modal/portal/status/focus/reduced-motion；75%指modal本身的viewport寬高，手機與小確認框依 DESIGN.md 已選尺寸分類。

相關modal／導航任務依「已選 modal 與返回引導契約」回歸一般／手機／小確認分類、固定操作區、contexthelp風險可見性及同scope返回狀態，不新增未選功能。
