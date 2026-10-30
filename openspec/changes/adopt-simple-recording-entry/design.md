## Context

AppRoutes 的 root/wildcard經GuidedWorkbenchEntryRedirect到 /studio/v2 並保留search/hash；legacyRoutes有step/target/section相容意圖。現行明確entity locator、runtime focus、Modbus與settings等交接不可被新首頁攔截。已核對現存不一致：buildLocalModbusCompatRedirect生成step=output&target=modbus，但parseStudioEntityLocator把非readiness的step當locator意圖而拒絕；此批需先用RED回歸證明，再在presentation routing範圍修復，不宣稱該舊連結目前已正常。前三案新增同route可選呈現並保留原configuration authority；本案只決定何時採用為預設。

## Goals / Non-Goals

**Goals:** 完整操作驗收與owner頁面確認後，新使用者直接到簡單入口、既有明確意圖保持可用；交付實際可證明的完成範圍及限制。

**Non-Goals:** 不新建V3或恢復/studio、不刪原進階配置能力、不擴真PLC/SCADA/正式DB/外網部署權限、不新增測試環境、不將demo當正式產品或自動發布。沒有視覺來源直接指定本案route行為，因此本案沒有Design Source；樣式承接前三案快取及owner確認版本。

## Decisions

### 預設入口採用與明確意圖優先

無明確意圖的 /studio/v2採記錄首頁；view=advanced保留原四步驟呈現，view=recordings明確開首頁。query先做raw encoding/duplicate驗證；只要有entry或workspace/device/rule/mapping/group等真正locator欄位，一律嚴格locator解析，錯誤不得降級。沒有真正locator欄位時，才容許已列舉的legacy step（device/source/tag/output）、target（database/modbus）、section/modal由既有allowlist處理；不得放任任意值。runtime readiness focus沿原有限定step與scope處理。既有entity locator、runtime readiness focus、純合法step/target/section/modal等有效相容意圖先解析，進入對應既有surface，不能被default搶走。root/wildcard仍用generic導向，不新增/studio特殊handler。不提供任意external return URL；return target沿現有受控route。

### 驗收矩陣與分層證據

四change的產品任務全完成才算功能完成。每批先相關測試；最終批次跑前端npm run lint、npm test -- --run、npm run build，並檢查built assets同步到cmd/test_ui/static的現有嵌入式交付流程。若只改前端不廣跑所有SQL engine；修改Go時才依docs/testing.md跑涉及側完整go test/vet/lint。

本環境已具Go1.25.5、相同lock的node_modules；不安裝新服務。已知browser啟動/loopback限制尚未解除，不能把歷史browser結果或waiver移到此批。Browser keyboard/Back/Forward/dialog/focus/responsive必須有本批實證；無法執行列NOT RUN並在必要斷點請owner看畫面，這不會自動取代全部browser回歸。設備/SQLmock、JSdom、真SQL、真PLC分開記錄。

### 使用者斷點與失敗回退

先把已審查規格的demo補齊，提供固定版本連結、情境操作、差異與待決定項讓owner確認。記錄owner選擇後再apply，不把沉默當同意。最終預設切換前，未過的必要驗收必須解決或取得本批明確接受的限制；未獲接受維持opt-in，不假稱完整驗收。回退只改UI入口選擇，不撤回已接受的資料庫/採集設定或外部操作。

### 唯一交接清單與最終批次覆核

release記錄列每個scenario的測試名稱/命令、PASS/FAIL/NOT RUN及證據層級、owner demo確認版本與時間、剩餘風險。不要使用歷史測試數替代此次結果。最終source變更之後重驗相關gate；涉及route/build整批變更時完整檢查必須對最終batch，不拿先前綠燈代表。每change獨立commit需使用spectra-commit及當時授權；本案不授權push或部署。

## Implementation Contract

共同 apply 前置：四案規格完成審查、補齊 demo，並已取得 owner 對本次固定版本關鍵頁面／路徑的確認。此 gate 在第一案產品實作前執行，不等到第四案；demo確認不代表正式瀏覽器或實機驗收通過。

- 依賴順序為home→editor→lifecycle→adoption；各change單獨可build並保留舊入口。前三案沒完成或owner斷點未確認，本案不得切換預設。
- 驗收路徑：空白連設備/點位/目的地/檢查/明確開始/等待receipt；已有來源跳過重設；草稿保存與reload；修改後再檢查；停止/恢復/複製/移除；寬長/目的地forward transition；partial/stale/DBfail/unknown修復；advanced-only配置與舊深連結。
- 特別驗證錯誤載入不呈空白、歷史成功不當現況、欄位變更後舊consent不能用、dialog關閉不取消已送出命令、double click/late response/切scope不重送、普通DOM無secret/raw exceptions。
- 前端瀏覽器journey用現有Playwright fixture pattern，不新增測試服務；真SQL僅用既有可丟棄同機測試目標且依變更範圍。沒有現場授權不接PLC或正式DB。
- 以source測試及build驗證新default候選；任何必要gate未過則最終交付維持原default，報告阻礙，不用默默跳過換取完成。針對同一最終batch記錄結果，不重複無意義全套。
- 不需編輯studio inventory即可完成本案release紀錄；若實作確實更新inventory，必須依AGENTS同步changelog.sqlite，不手改資料庫。

## Risks / Trade-offs

- [入口切換破壞舊連結] → explicit intent優先和route回歸，view=advanced持續可用。
- [測試環境缺口被誤稱通過] → matrix明列NOT RUN、必要gate停住預設切換及owner明確確認。
- [回退被誤認資料rollback] → 僅UI呈現回退，既有DB歷史與operation不可撤回。
- [新舊介面長期分叉] → 共用service/state，原進階保留能力不另演進第二產品；所有修改須兩入口回歸。

## 2026-10-30 需求更新（尚待版面選項確認）

根目錄 `DESIGN.md` 為唯一人讀設計規範；依 Taste app-screens 與 designtoken.md 的 primitive → semantic → component 分層，沿用既有 ThemeProvider 與字型，不另建主題系統。舊 demo 的淺色色值僅是歷史快取證據，不再是實作目標。

新簡易介面的彈出設定改為置中 modal，modal 本身寬度為目前 browser viewport 的 75%、高度為 75%（75vw × 75dvh），內容內捲；主頁仍正常 RWD，不是頁面縮為75%。已選手機近滿版與小型確認較小例外，尺寸分類依根 DESIGN.md；一般設定仍75% viewport。固定header/footer、contexthelp與同scope返回保留搜尋／篩選／捲動納入；選項8已選為dev/test專用樣式驗收頁，不放主選單，3/5/6/7未選，不新增tasks。文案簡潔、標籤保留意義，繁中／英文與數值／單位不能因版面省略。

共同 owner gate 必須用這次修訂後固定 demo 版本；過去 demo 檢視不視為確認。現處需求收集階段，未 apply、未做產品驗收；確定的協議保真與主題要求先納入，未選建議不擴 scope。

### 新版設計與能力驗收

acceptance matrix逐protocol列必要／進階欄位、既存設定round-trip與真validator，MQTT disabled；不把五類source可達當實機PASS。逐light/dark/system列OS live change、legacy偏好、raw server key缺失、late hydration、保存失敗、跨頁與portal、focus/contrast/reduced-motion、1440/768/390寬度與75% modal。手機／小dialog採已選例外與 DESIGN.md 分類，選項8已選為隱藏dev/test驗收入口；必要視覺gate無法執行記NOT RUN。Taste純文字規範採用不等於scripts或browser通過。

### 已選 modal 與返回引導契約

依 DESIGN.md 尺碼分類：一般popup75vw×75dvh；mobile<640px或short viewport<480px近滿版；小型短確認較小autoheight。第二主畫面editor維持RWD，不是popup。modal固定header/footer、body內捲；criticalrisk/validation/consent不能藏contexthelp。返回搜尋/篩選/捲動僅同scope、entity與view，revision變更重驗identity，不能復活舊selection/availability/receipt。dirtyguard、focus與pending操作語意不變。選項8已選為dev/test專用樣式驗收頁，不放主選單，未選候選不入實作；此需求更新不是owner對新版demo的視覺核准。

## v4 owner 未驗收

v4配色、密度與部分行為被使用者否決；既有mock測試不是UI確認。本輪先修demo及規格，修訂固定版本重新確認後才可apply；已以DESIGN.md Quiet mineral候選替換舊blue/slate；新版仍待owner視覺驗收，不當核准品牌。

### 混合member計算前置與整體驗收

新增前置support-per-member-recording-calculation（獨立於source字序）；editor consume完成才進adoption。驗收mean/mean/delta同record各member獨立window與真SQL一列22/12/8、firstbaseline requiredskip/optionalNULL、gapspan原因、同source不同record隔離、legacy不變、method更改freshschema/proof及replay不重複。新mixed需canonical managed wide／long production gates完成；wide未到期NULL＋not_due、long只到期member列，其他未支持shape受阻原因可見。不以mock/helper替代production closure/readback；新版owner確認仍必要。


本輪已確認逐member interval與UTC epoch對齊：5秒00/05…55、15秒00/15/30/45、10分鐘00/10…50；wide採due邊界聯集，:15/:30/:45的60秒delta留NULL＋not_due，:00同列各自前15秒mean／前60秒delta；long只新增到期member列。不carry、不補0；required-not-due不skip，到期失敗／無baseline／gap/reset分別標原因。顯示各member的[start,end)、partial coverage與next due，SQL commit時間另列；first partial、nullable schema、安全gate與durable驗收依support-per-member-recording-calculation。此語意確認不是新版demo或apply核准。

### 同源採集次數預覽交接

integrate-protocol-aware-acquisition-batching 唯一負責planner、preview API/shared types與實際計數來源；本案僅接既有consumer呈現／驗收，不另推導合批、不新增採集引擎。該案production gate完成後才顯正式「每輪計畫協議讀取次數（不含故障重試）」與estimate標示；顯示acquisition point count、unique wire units、planned requests及安全coverage/fallback原因，不能用recording已選點數當整設備採集集合或自行猜random=1。API預覽零PLC I/O，沿workspace/device scope與expected revision，stale409／載入失敗清除可用性並提示重新取得，不展示過期成功。

actual attempted/succeeded/failed/retry/fallback/canceled及cycle latency使用另列同scope/revision的真實觀測；缺乏實測就顯尚無實測，不補0、不將估算稱TCP封包數、SQL寫入數或已證明加速。demo只可明標模擬契約，無須也不能假造真實性能。獨立member recording window仍由support-per-member-recording-calculation負責；更改selector不改採集plan，批次預覽依賴不升為member排程的硬性前置。最終adoption必須分開列預覽契約測試與實際transport量測證據，保留新版demo owner與apply gate。

本輪「新版設計與能力驗收」另覆核選點實際旅程：主頁唯一入口／逐member設定→加法picker→同表搜尋設備類型篩選→計數相符的加入／取消；既有點locked、remove/undo/re-add保意圖、非法輸入文字保留、空設備續建、無第一台默認、同框setup與獨立設備／範圍保存、不回滾已存來源、同scope位置與pending恢復。這些是原editor／adoption tasks內的驗收細化，不增change／任務；新版owner checkpoint與必要browser證據依舊不得由mock取代。
