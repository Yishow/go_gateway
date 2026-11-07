## 1. 契約與測試先行

- [x] 1.1 補能重現平板長格子與非全螢幕 modal 的回歸測試，驗證 Selected memory facts remain readable without expanding address rows 及 New configuration popups are centered viewport-sized modals；先記錄舊版失敗，涵蓋長描述、平板form/confirm與桌面不變。
- [x] 1.2 [after: 1.1] 更新 DESIGN.md 的手機/平板全螢幕分類與格外選取資訊，對照五張原圖和本案delta做文件review，保留token authority與使用者最新記錄點位表變更。

## 2. 響應式與選取資訊

- [ ] 2.1 [after: 1.1] 實作全螢幕分類共用但限定真正 modal，讓 Radix、來源host、Dialog、ConfirmDialog、SlidePanel 在手機/平板填滿viewport，保留非modal popover與桌面尺寸；依 Classified modals and contextual return state preserve safety 驗證各primitive邊界、safe area與內容捲動。
- [ ] 2.2 [after: 1.1] 實作選取資訊脫離格子並維持穩定空間，以及可見值訂閱保留canonical地址；保留完整aria-label、同源值與unknown保護，單元測試與選取前後量測確認列高穩定、continuation正確訂閱。
- [ ] 2.3 [after: 2.1, 2.2] 整理五種記憶體視圖的工具列、格子與footer層級，於手機/平板和明暗主題比對原圖構圖，驗證不溢出、不縮小觸控目標且所有layout可切換。

## 3. 實測與交付

- [ ] 3.1 [after: 1.2, 2.3] 用 agent-browser 與現有E2E在390×844、768×1024、1024×768、1366×1024 touch、1440×900 desktop驗證五視圖、各modal類別、旋轉、連續選取、Escape/Close與dirty返回保護，保存當前畫面和結果，不把blocked標通過。
- [ ] 3.2 [after: 3.1] 對最終前端執行完整lint/test/build、diff/line gate及獨立review；修復本案問題後重驗受影響範圍，紀錄驗證限制，完成Spectra verify/review與可封存交接。
