## Context

## 目前決策／Supersedes

本段依 `docs/verification/simple-recording/owner-checkpoint-2026-10-30.md` 與 `implementation-decision-2026-10-30.md`，取代本文較早「待 owner 確認／未 apply／owner gate 停止」的當前狀態解讀。v2、v4、v5 原型與待確認敘述保留為歷史；正式實作參考已固定為 Site v9（commit `c7698bf137dc0d966e2be19c0ba6b06110da02d4`）。2026-10-30 的 owner 指示及已記錄判斷已解除是否開始正式實作的待決狀態，無需再次詢問布局。

本段不宣稱 owner 已逐頁驗收正式產品，不變更現行 DESIGN.md 的 tokens、功能範圍、task 依賴、checkbox 或完成證據。正式 browser／SQL／協議驗收及 adoption 預設切換 gate 均維持；未完成必要驗收不得自行切換預設。Taste app-screen 文字規範與第三方 scripts 未執行的證據邊界不變。

AppRoutes 的 root/wildcard經GuidedWorkbenchEntryRedirect到 /studio/v2 並保留search/hash；legacyRoutes有step/target/section相容意圖。現行明確entity locator、runtime focus、Modbus與settings等交接不可被新首頁攔截。已核對現存不一致：buildLocalModbusCompatRedirect生成step=output&target=modbus，但parseStudioEntityLocator把非readiness的step當locator意圖而拒絕；此批需先用RED回歸證明，再在presentation routing範圍修復，不宣稱該舊連結目前已正常。前三案新增同route可選呈現並保留原configuration authority；本案只決定何時採用為預設。

## Goals / Non-Goals

**Goals:** 完整操作驗收與owner頁面確認後，新使用者直接到簡單入口、既有明確意圖保持可用；交付實際可證明的完成範圍及限制。

**Non-Goals:** 不新建V3或恢復/studio、不刪原進階配置能力、不擴真PLC/SCADA/正式DB/外網部署權限、不新增測試環境、不將demo當正式產品或自動發布。沒有視覺來源直接指定本案route行為，因此本案沒有Design Source；樣式承接前三案快取及owner確認版本。

## Decisions

### 預設入口採用與明確意圖優先

無明確意圖的 /studio/v2採記錄首頁；view=advanced保留原四步驟呈現，view=recordings明確開首頁。query先做raw encoding/duplicate驗證；只要有entry或workspace/device/rule/mapping/group等真正locator欄位，一律嚴格locator解析，錯誤不得降級。沒有真正locator欄位時，才容許已列舉的legacy step（device/source/tag/output）、target（database/modbus）、section/modal由既有allowlist處理；不得放任任意值。runtime readiness focus沿原有限定step與scope處理。既有entity locator、runtime readiness focus、純合法step/target/section/modal等有效相容意圖先解析，進入對應既有surface，不能被default搶走。root/wildcard仍用generic導向，不新增/studio特殊handler。不提供任意external return URL；return target沿現有受控route。

### Demo 強化先接簡易呈現，預設另依所選主流程驗收

**Supersedes**: adopt-simple-recording-entry / 驗收矩陣與分層證據

依2026-11-01使用者最新「以討論清單及最後決議做為驗收點，強化demo」，本輪落實「精簡原Step4」及「最終完成版就是demo的強化版」，不再將完整四案／七案、Settings全頁或剩餘大型效能矩陣全部完成當成簡易host接線前置。既有home/editor可用契約先整合；未完成lifecycle、採集優化及其原owner任務保留pending，不因本批入口交付假done或刪需求。討論清單80–87不是全部已選；UTC epoch等既有實作選擇不改寫為使用者逐項確認。

原已列的bounded legacy辨識修復隨1.2 routing foundation先做並保留RED；2.1預設採用沿用該證據及相關回歸，避免1.3先要求相容驗收、2.1才修的順序死結，不增加新的相容模式。

本決議只取代舊「全案完成才可接線」順序，不取消安全保護、真Browser或所選主流程的必要SQL證據。1.2負責同一Step4的簡易呈現；1.3負責總覽新增／修改／返回接線及受影響主流程驗收。預設仍待「已選demo主旅程＋其必要生命週期安全」完整：新增、修改、停止／恢復、複製、移除與已選表格／目的地前向變更逐項真驗，必要安全缺陷閉合及2.1 gate；未過維持opt-in。原1.1完成描述及所有完成證據不變。未受影響證據可明列原版本沿用，不能冒充本批重跑；本批影響的行為須重驗。

### 驗收矩陣與分層證據

依「Demo 強化先接簡易呈現，預設另依所選主流程驗收」分開接線與預設；未完成原owner任務繼續pending，不把本批完成寫成全案完成。每批先相關測試；最終批次跑前端npm run lint、npm test -- --run、npm run build，並檢查built assets同步到cmd/test_ui/static的現有嵌入式交付流程。若只改前端不廣跑所有SQL engine；修改Go時才依docs/testing.md跑涉及側完整go test/vet/lint。

環境與Browser可用性以本批實際檢查為準，不新增測試服務；不能把歷史browser限制、結果或waiver移到此批。Browser keyboard/Back/Forward/dialog/focus/responsive必須有本批實證；無法執行列NOT RUN並在必要斷點請owner看畫面，這不會自動取代全部browser回歸。設備/SQLmock、JSdom、真SQL、真PLC分開記錄。

### 使用者斷點與失敗回退

先把已審查規格的demo補齊，提供固定版本連結、情境操作、差異與待決定項讓owner確認。記錄owner選擇後再apply，不把沉默當同意。最終預設切換前，未過的必要驗收必須解決或取得本批明確接受的限制；未獲接受維持opt-in，不假稱完整驗收。回退只改UI入口選擇，不撤回已接受的資料庫/採集設定或外部操作。

### 唯一交接清單與最終批次覆核

release記錄列每個scenario的測試名稱/命令、PASS/FAIL/NOT RUN及證據層級、owner demo確認版本與時間、剩餘風險。不要使用歷史測試數替代此次結果。最終source變更之後重驗相關gate；涉及route/build整批變更時完整檢查必須對最終batch，不拿先前綠燈代表。每change獨立commit需使用spectra-commit及當時授權；本案不授權push或部署。

## Implementation Contract

共同 apply 前置：既有固定Site v9（commit `c7698bf137dc0d966e2be19c0ba6b06110da02d4`）owner checkpoint與最新強化demo指示為本輪依據；不再新增布局核准迴圈。接線不是預設切換或正式驗收完成。

- **In scope**：重用原Step4及既有home/editor，增加同route簡易呈現與總覽新增／修改／返回；保留進階及深連結。簡易頁可先掛載，不等待未用到的完整lifecycle或效能矩陣。
- **Authority / interface**：`WriteGroupSection`既有new/edit selection及`GroupEditor`的`group`、`workspaceId`、`workspaceRevision`、`fetchLatest`、`onSaved`契約繼續負責保存與readback；`DeviceSourcePort`負責設備／範圍保存；`GroupEditorPreparation`與`SavedRecordingStartReview`／`LocatorSavedRecordingStartReview`維持確認及scoped Start。新增host props僅傳呈現、new/edit意圖及受guard的返回，不新增record DTO、保存queue、API或activation owner；同時僅掛載一套owner，不能把新舊owner以hidden雙掛。
- **Failure / lifecycle**：沿`RecordingDraftNavigation`、locator freshness與原pending/unknown owner；載入失敗不可變空白成功。lifecycle僅呈現匹配workspace/group revisions的`WriteGroupLifecycle.actions`之`allowed`、`reason_code`、`requires_review`；允許的動作走原owner且可操作；真正不支援要說原因／導原進階，不以status猜命令、不為demo補造resume/copy API。尚未接妥的已選必要動作仍是pending，不得冒稱unsupported或藏掉來提前切default；必要限制只能依本批明確接受記錄處理。阻擋所選主流程的安全問題仍須原owner修復後才能驗收。
- **Out of scope**：不另造editor、schema／計算／採集／傳送engine；不翻修Settings全頁、不擴協議或DB支援、不以接線為由新增效能矩陣、監測掛鉤、搬移歷史或重做全部既有頁面。採集優化需求不取消，未完成項留原owner pending。

本批只有下列七個可觀察驗收點；子情境用既有scenario／owner證據對照，不另增使用者清單：

1. **兩主畫面**：總覽的新增及既有記錄修改直接進入v9方向的單頁設定，返回回到總覽；一般使用者不用先展開Advanced／WriteGroups或重走四頁。無來源／未存目的地也能先開本地草稿並就地補齊；這不解除保存及Start的原prerequisite。
2. **同頁與來源返回**：名稱、已選來源、逐點算法／週期、目的地及必要下一步同頁；唯一加法選點與來源modal沿原元件，取消保留原草稿、已保存設備／範圍不回滾，同scope搜尋／篩選／位置保留。
3. **真正保存與明確開始**：save/reload對回同workspace/group/revision與逐點設定；schema影響／確認／Start沿原owner，保存不啟動、變更使舊consent失效、只啟動所選scope，開始後先等第一筆而非假稱SQL成功。
4. **生命週期與失敗真實**：修改、停止／恢復、複製、移除、已選表格／目的地變更逐項驗證；lifecycle依server availability，允許的走原owner可操作，真正不支援明示原因及原advanced交接，未完成必要動作保留pending；stale、partial、DBfail、pending/unknown能看原因及原修復入口，未知結果查原operation、不盲目重送或使停用writer復活。
5. **結果可核對**：總覽與原runtime/logs可查來源／記錄規則、saved/effective差異、待送與最近真實寫入證據；採集中、SQL committed與readback verified分開，無actual時不補0或宣稱效能。
6. **保留原能力**：advanced-only能力、明確locator、runtime focus及合法legacy step/target/section/modal仍到原對象；invalid locator fail-closed；切呈現零mutation，不改資料身分或歷史。
7. **已選demo主旅程＋其必要生命週期安全通過才預設**：以本批真Browser完成首次來源建立→選點／逐點設定→目的地→保存／確認→明確Start→真實結果，以及既有記錄修改／reload／返回、第4點生命週期及必要中斷回歸；所選mixed寬／長路徑沿原due/NULL與SQL證據。核對v9資訊層級、DESIGN.md主題／modal／RWD。必要項FAIL／NOT RUN不切預設，不能以DOM/mock或全部任務勾選代替。

- 開發只跑當下受影響測試；完整功能批次才跑既定前端lint/test/build及embedded static檢查。Browser互動沿使用者已選工具，既有Playwright fixture可作回歸但不是對真Browser旅程的替代；真SQL限本輪既有可丟棄目標，未授權不接真PLC／正式DB。
- 不需編輯studio inventory即可完成release紀錄；若實作確實更新inventory，仍依AGENTS同步changelog.sqlite。

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
