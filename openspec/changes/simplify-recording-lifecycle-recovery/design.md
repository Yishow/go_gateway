## Context

已讀 SavedRecordingActions 的 server availability/revision gate、RecordingChangeDialog 的 choose/impact/reserved/schema/activation/unknown流程、recordingWorkflowView 的 scoped intent。此批改變使用者看到的動作與指引，沿用既有 engine。demo 的手動成功/失敗切換是測試fixture，不得接到正式應用。

## Design Source

2026-10-30 demo v2 原稿快取 .spectra/design-cache/simplify-recording-lifecycle-recovery/gateway-simple-ui-v2.dc.html；SHA256 d6987f1c8749db5d0ff9898544e3b2d505e98a99774b63d936a75daab2ffa306。原稿實測正常chip #eaf3ec/#34735c、警告 #f5efde/#94722a、異常 #f9eae5/#a15e45、未知 #edeef0/#6c7180；按鈕圓角8px、內距10px 17px，卡片圓角12px；原稿無命名圖示組，正式使用專案有標籤圖示。

原稿缺恢復、複製、前向表格變更與精確修復，這些由本規格補齊demo再於owner斷點確認。不得沿用原稿直接取代record即代表完成的模擬語意。

## Goals / Non-Goals

**Goals:** 從建立後到修改、停用、恢復、移除與異常修復可完成既有支援流程，明確顯示目前設定與實際執行結果。

**Non-Goals:** 不新增複製後端API、condition trigger engine、DB dialect、歷史資料遷移/刪除、原地寬長互轉、任意SQL瀏覽器、背景自動重送；不變更已有 cutover/ownership/CAS/receipt 行為。

## Decisions

### 操作選單依後端事實而非狀態猜測

主操作隨目前可用性與下一步變化；edit/stop/remove/recreate/transition仍以匹配scope/revision的 lifecycle actions為準。不適用項目顯示原因及下一步，不能憑灰色或綠色推斷。停止與移除需明確確認對象與影響；移除只改設定，保留資料表與歷史。

### 恢復需重新檢查且複製只產生草稿

恢復重用 review-scoped-start與最新證據，不能把本地status改回running。複製從可表示的設定產生未啟用新草稿，新的identity、operation、readiness、receipt與proof不得取自舊group。destination/table需明確檢查，不自動共用寫入責任。managed ensure若只能回原canonical角色，顯示不可建立重複角色與進階入口；不得製造副本成功回報。未支援simple editor的配置不做有損複製。

### 前向表格變更沿用既有安全階段

在可用transition/recreate動作中預先顯示原布局/目的地、目標布局/新表與來源影響，再依正式successor契約review、reserve、schema plan、必要confirm、start-or-cutover。預設新表接續不原地改歷史；具體新表名與唯一性由既有服務決定，使用者可確認但不得越過保護。顯示切換邊界與未完成階段；首次啟動與已開始記錄切換不可混淆。簡化頁面層次但不把多個權限命令合成隱性一鍵。

### 異常修復保留未決結果與歷史

partial source、stale、DB write failure、unknown各自有安全原因與scope導航。unknown優先查原operation，沒有fresh proof不用新命令替代。保留RecordingChangeController.retryExact與canRepreviewAfterStatus既有合法復原：fresh匹配status證實not_submitted／preview_not_submitted等原契約條件後，才顯示明確的exact retry或repreview；不一律鎖死為只能查詢，也不自行放寬server proof。關閉dialog不是取消已送出操作；回首頁顯示pending/unknown提示，可resume查同一operation。舊成功時間保持歷史標示，不轉成現在正常。

### 細節按需展示並保留進階交接

記錄詳情先顯示來源設備/點位、讀取與記錄規則、目的地/表格、最近SQL receipt及目前狀態。saved/effective差異才突出，不重複堆兩份相同摘要。runtime/logs入口帶現有locator；一般delivery沒有raw值時不展示假資料列，示例/試寫readback與正式receipt分開。unsupported配置保持唯讀摘要與進階入口。

## Implementation Contract

共同 apply 前置：四案規格完成審查、補齊 demo，並已取得 owner 對本次固定版本關鍵頁面／路徑的確認。此 gate 在第一案產品實作前執行，不等到第四案；demo確認不代表正式瀏覽器或實機驗收通過。

- 前置 streamline-recording-task-editor；只從現有 lifecycle/actions 與 recordingWorkflowView 產生操作。frontend新增presentation adapter，不新增後端truth或資料migration。
- 受影響 source/device/range 刪除或identity改動，透過既有依賴檢查顯示引用記錄與修復入口，不能讓新editor繞過原有保護；基本頁無法完成時交接既有進階流程。
- 變更editor的dirty guard與operation recovery共用；重複點擊、關閉重開、Back/Forward、切scope及過期回覆不能造成重送或丟失unknown狀態。
- 測試覆蓋停止→重驗證→恢復、移除保留歷史、複製清除能力/證據、managed去重、wide→long及long→wide的可用/受阻、首次啟動/切換、schema確認後輸入異動、pending/unknown恢復。
- 複製安全驗證包括來源group ID不作新group ID、無舊operation/readiness/receipt/proof、原group完全未變、提交新draft時仍有既有ownership/CAS檢查。
- mock HTTP/SSE測試只驗證UI契約，SQL和PLC須獨立證據。此批只跑相關lifecycle/receipt/locator/runtime/log tests與lint/typecheck；如未改後端不廣跑無關DB或真設備。

## Risks / Trade-offs

- [複製創造重複writer] → 只草稿，不帶權限證據；提交仍由現有canonical ownership判斷。
- [表格互轉被誤認歷史轉檔] → 新表接續、邊界與舊歷史保留明示；能力未允許就受阻。
- [關閉dialog被誤認取消命令] → 已送出狀態留存及resume，同一operation查詢。
- [錯誤細節洩露endpoint或憑證] → 保留typed safe copy和按需locator，不把raw exception送正常DOM。

## 2026-10-30 需求更新（尚待版面選項確認）

根目錄 `DESIGN.md` 為唯一人讀設計規範；依 Taste app-screens 與 designtoken.md 的 primitive → semantic → component 分層，沿用既有 ThemeProvider 與字型，不另建主題系統。舊 demo 的淺色色值僅是歷史快取證據，不再是實作目標。

新簡易介面的彈出設定改為置中 modal，modal 本身寬度為目前 browser viewport 的 75%、高度為 75%（75vw × 75dvh），內容內捲；主頁仍正常 RWD，不是頁面縮為75%。已選手機近滿版與小型確認較小例外，尺寸分類依根 DESIGN.md；一般設定仍75% viewport。固定header/footer、contexthelp與同scope返回保留搜尋／篩選／捲動納入；選項8已選為dev/test專用樣式驗收頁，不放主選單，3/5/6/7未選，不新增tasks。文案簡潔、標籤保留意義，繁中／英文與數值／單位不能因版面省略。

共同 owner gate 必須用這次修訂後固定 demo 版本；過去 demo 檢視不視為確認。現處需求收集階段，未 apply、未做產品驗收；確定的協議保真與主題要求先納入，未選建議不擴 scope。

### 已選 modal 與返回引導契約

依 DESIGN.md 尺碼分類：一般popup75vw×75dvh；mobile<640px或short viewport<480px近滿版；小型短確認較小autoheight。第二主畫面editor維持RWD，不是popup。modal固定header/footer、body內捲；criticalrisk/validation/consent不能藏contexthelp。返回搜尋/篩選/捲動僅同scope、entity與view，revision變更重驗identity，不能復活舊selection/availability/receipt。dirtyguard、focus與pending操作語意不變。選項8已選為dev/test專用樣式驗收頁，不放主選單，未選候選不入實作；此需求更新不是owner對新版demo的視覺核准。

## v4 owner 未驗收

v4配色、密度與部分行為被使用者否決；既有mock測試不是UI確認。本輪先修demo及規格，修訂固定版本重新確認後才可apply；已以DESIGN.md Quiet mineral候選替換舊blue/slate；新版仍待owner視覺驗收，不當核准品牌。
