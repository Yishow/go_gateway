## 1. 首頁證據與導引

- [x] 1.1 在共同owner demo確認gate完成後，以 Red-Green-Refactor 建立 Home guidance projects canonical scoped evidence 純投影；驗證 opening 零 mutation、舊 scope 回覆不覆蓋，執行 recording-home projection 單元測試。
- [x] 1.2 [after: 1.1] 實作 Guidance states coexist with deterministic priority 的「首頁優先序與同時存在的狀態」，涵蓋空白、有設備無點、有點無規則、草稿、正常、異常與混合；以相同 resolver fixtures 回歸測試驗證不遮蔽其他項目。
- [ ] 1.3 [after: 1.2] 串接首頁唯讀 hooks 與 safe error state；驗證 loading/error/unprovable-scope 不偽裝空白，切換 workspace 後過期回覆不啟用命令。

## 2. 同路由與結果呈現

- [ ] 2.1 [after: 1.3] 實作 Optional home preserves existing Studio navigation 與「同路由可選入口與既有識別保持」，保留 view、search/hash、locator 和 runtime focus 優先權；用 AppRoutes、legacyRoutes、studioEntityLocator 相關測試驗證valid locator+單一view、invalid locator+view、duplicate view及malformed encoding，invalid意圖不得降級首頁。
- [ ] 2.2 [after: 2.1] 實作 Large recording collections remain searchable and bounded 與「大量設備與規則可定位而非無限卡片」，首頁以compact管理表、keyword/filter/stable sort/25或50分頁與設備最多4列摘要及同頁管理呈現，不另建 workspace authority；以 recording-home component tests 驗證大量fixture的100設備/300規則設計fixtures的每頁bounded rows（非backend效能保證）、跨頁定位/同scope返回、查無結果/empty分開、revision移除後page clamp、可達草稿/正常/異常與導航零寫入。
- [ ] 2.3 [after: 2.2] 實作 Healthy summaries distinguish history from present evidence 與「結果證據而非狀態綠燈」；測試已保存、等待首次寫入、SQL receipt 與舊成功後斷線分開顯示，不製造資料列數值。
- [ ] 2.4 [after: 2.3] 依「Quiet mineral 取代被否決的 v4 配色」與 DESIGN.md 完成雙主題語意 tokens、繁中與英文精簡文案及字型階層、鍵盤焦點/返回焦點、長名稱與窄畫面排版；以 accessibility/component tests 及可用瀏覽器截圖驗證，無法執行的視覺項明列 NOT RUN。

- [ ] 2.5 [after: 2.4] 實作 Theme preference resolves without inventing persisted intent，依「系統主題與單一語意 tokens」與 UI / API / Diagnostics options 連接現有provider及settings adapter；測試raw key有無、legacy相容、OS live change、override、回system、失敗fallback、late hydration與backend context切換／workspace不重設，theme-only updateKey只寫theme且modbus_share/default_interval零write。
- [ ] 2.6 [after: 2.5] 實作 Semantic design tokens govern affected surfaces、Design tokens scope 與 New configuration popups are centered viewport-sized modals，依 DESIGN.md 對新 surfaces／重用modal／portal導入單一CSS semantic tokens與TS aliases；檢查light/dark/focus/disabled/status對比、reduced-motion與viewport75% modal，不改未觸及舊頁。

- [ ] 2.7 [after: 2.6] 依「已選 modal 與返回引導契約」實作 Classified modals and contextual return state preserve safety；測一般75%／手機近滿版／小確認尺寸、固定header/footer與body捲動、help不藏risk、同scope返回搜尋/篩選/scroll/focus及scope/revision變更不復活舊項目。

- [ ] 2.8 [after: 2.7] 實作 Development style preview is isolated from operator navigation 與「開發驗收樣式頁不成為第三主流程」；以actualcomponents/tokens展示狀態，驗證dev/test入口可達、主選單無項目、production build不bundle且零資料寫入/啟用/credential。

## 3. 驗證與交接

- [ ] 3.1 [after: 2.8] 執行本批首頁/route/locator/runtime-return 相關測試、受影響檔 lint、型別檢查、git diff --check 與 make check-lines；記錄命令和實際 PASS/FAIL/NOT RUN，不提前跑整個產品全套。
- [ ] 3.2 [after: 3.1] 獨立 review 首頁零 mutation、scope、異常及既有路徑回歸，修復本批問題後重跑相關測試；交接 editor/lifecycle 尚未完成範圍，不宣稱整功能驗收完成。

需求變更註：原 pending 2.4 的「完成淺色樣式」由雙主題／tokens取代，不算完成；原側欄pending工作由同語意的置中modal取代。沒有已完成工作被重標。
