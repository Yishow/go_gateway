## 1. 來源與就地設定

- [ ] 1.1 在 simplify-recording-home-guidance 完成及 owner demo 斷點確認後，以 Red-Green-Refactor 為 Progressive editing preserves source capabilities 建立 editor/modal 測試；鎖定「就地補設備與點位但共用原資料權威」及沒有 fake identity。
- [ ] 1.2 [after: 1.1] 實作 Source editing begins with selected items and bounded discovery 與「編輯先看已選來源再按需搜尋新增」、「加法選點與來源保存分層」，接唯一選點入口／同表搜尋設備類型filter／獨立pendingAdd／固定計數與Cancel/Add N；測既有點locked、N=0、跨頁隱藏pending、取消不改draft或非法輸入、主頁remove/undo/同ID重加保留method與interval、失效ID阻擋、單一active modal與無第一台默認；接來源選取與就地設備設定，保存置中 modal後回到原名稱、來源、目的地及焦點；執行 recording-editor 與既有 device/autosave 相關測試，驗證取消只捨棄未送出buffer且零隱性啟動；依「空來源建立後回到可選點位」回歸零選取新記錄、可選既有空設備並新增範圍、同modal view stack往返／取消保留草稿與member method/interval及selection；實作 Closing an editor does not roll back submitted source changes，測試已成功source保存後捨棄recording不rollback，以及pending/unknown查證不冒稱取消。
- [ ] 1.3 [after: 1.2] 實作 Successful read-range saving closes only its matching modal 與「範圍保存成功才關閉並回到原位置」，接就地 RuleEditor，保留 datatype/address/word-order/scale/unit 驗證與現有採集排程摘要／進階入口（不新增per-rule poll）；測成功關modal/失敗或unknown保輸入/pending防重送/late不關新modal/新buffer generation，pending後新編輯不被舊success清掉，以及new-device→range-save與既有空設備range-save返回bounded discovery、明示清除遮擋新點的query／類型filter並定位新點且不自動全選、保留原member設定，以及設備保存後取消範圍仍保設備、取消picker仍保已存來源、取消setup還原query/filter/page/scroll/focus/pending，以及10 float32=20 word、跨設備同名與同位址分離、persisted MC compatibility 受阻提示。

- [ ] 1.4 [after: 1.3] 實作 Protocol and parameter parity is explicit before implementation，依「既有協議與參數保真矩陣」投影5個可達protocol及MQTT disabled；測試RTU完整serial、FATEK mode、MC aliases/radix、server baud限制、單欄更新保留隱藏值與原validator，逐欄標示可編輯/進階/unsupported；在batching前置完成後依「同源採集次數預覽交接」實作 Acquisition request preview consumes the canonical planner，測同scope/revision估算與actual分列、stale409／無實測不補0、零PLC I/O預覽，不自行推導合批。

- [ ] 1.5 [after: 1.4] 實作 Recent connection defaults are session-scoped and explicitly clearable 與「同協議最近連線參數只作新表單預填」；測再新增預選上一成功protocol/subtype但不覆蓋已輸入form、同協議success後預填、FATEK子mode隔離、whitelist排secret/identity/proof、成功request snapshot而非當下form、cancel/failed/apply_failed/unknown/late不更新、clear後晚回應不復活、清除回defaults不改existing、dirty取消不清、reload/離頁/scopechange清cache及明確save不autoactivate。

- [ ] 1.6 [after: 1.5] 在 preserve-per-point-word-order 完成其API／capability gate後，實作 Point order controls consume the canonical override contract 與「每點字序依正式前置能力呈現」；逐點呈現inherit/override/legacy/effective，測同PLC混合ABCD/CDAB、range更新不覆override、resetinherit、explicit等於default仍保override、單word不新增override亦不丟legacybyteswap、API readback mode/configured/effective/source/applicable/revisions、datatype/protocol受限與stale revision，UI不自建decoder或假稱保存成功。

## 2. 草稿與記錄設定

- [ ] 2.1 [after: 1.6] 實作 Basic and advanced recording semantics remain faithful 與「同頁漸進引導與能力保真」，在 support-per-member-recording-calculation 完成後依「同一記錄每個成員選擇計算結果」實作 Per-record member controls consume mixed calculation intent，並實作 Recording choices explain outcomes without inventing restrictions 與「以使用目的說明記錄結果」，基本 raw 值不強迫衍生語意，特殊配置只讀摘要可交接進階；測試白話label/例子及gauge mean但counter_raw/state latest、first baseline NULL與signed/rollover不能一律判missing、同組mean/mean/delta各member自己的對齊窗口、到期同列APIintent、同source不同record隔離、原始counter/legacy契約不默轉、requiredskip/optionalNULL與cross-device不可默轉。
- [ ] 2.2 [after: 2.1] 實作 Draft navigation preserves or explicitly discards unsaved intent 與「草稿返回與確認失效契約」；以 dirty-navigation 測試覆蓋 Close/Escape/Back/Forward/切scope/reload/autosave失敗，驗證保留與捨棄語意。
- [ ] 2.3 [after: 2.2] 實作 Four-step intent-led setup 的同V2新舊呈現共用配置，保留原進階入口、routine mapping、offline draft及既有role約束；用原 guided workflow 與新 editor 整合測試確認導航不轉換配置。

- [ ] 2.4 [after: 2.3] 實作 Configured destinations offer safe table reuse 與「從已設定目的地安全選表」；測不同connector/database同名不合併、兼容長表可選、寬表受owner限制導原編輯、無候選/失敗/過期狀態，選取零寫入且fresh schema/ownership/revision gate仍必要。

## 3. 檢查並開始

- [ ] 3.1 [after: 2.4] 實作 One next action retains every safety boundary 與「重用檢查及明確開始」的save/check/schema impact流程；測試輸入變動清除舊consent、preview不是成功、schema確認不自動start。
- [ ] 3.2 [after: 3.1] 串接既有 scoped start，開始後僅waiting-first-write；測試double click、過期response、unknown查狀態與不啟動整workspace。
- [ ] 3.3 [after: 3.2] 依 DESIGN.md 補齊雙主題、繁中/英文簡潔指引、字型階層、常用連線參數常顯、點位compact rows／低頻細節才折疊、依「緊湊表單與短確認分類」實作 Compact forms and short confirmations preserve useful density，普通75%modal與dirty/stop/remove/discard短確認autoheight、鍵盤及RWD；實作 Parameter help is accessible and source-grounded 與「每個參數都有可操作的圓圈說明」，逐欄coverage map對照；測click/touch/keyboard開關、focus/Escape/outside/no-submit、risk inline及source準確。執行component/a11y檢查並記錄瀏覽器視覺項實際PASS或NOT RUN。

## 4. 驗證與交接

- [ ] 4.1 [after: 3.3] 跑本批editor/device/rule/mapping/autosave/preparation/start相關測試、lint/typecheck、git diff --check與make check-lines；記錄實際結果，不把mock或jsdom當真實PLC/SQL或視覺驗收。
- [ ] 4.2 [after: 4.1] 獨立review資料權威、草稿與confirmation失效、advanced交接，修復本批問題並重跑受影響測試；列出 lifecycle後續依賴，不重做已完成引擎。

相關modal／導航任務依「已選 modal 與返回引導契約」回歸一般／手機／小確認分類、固定操作區、contexthelp風險可見性及同scope返回狀態，不新增未選功能。


本輪已確認逐member interval與UTC epoch對齊：5秒00/05…55、15秒00/15/30/45、10分鐘00/10…50；wide採due邊界聯集，:15/:30/:45的60秒delta留NULL＋not_due，:00同列各自前15秒mean／前60秒delta；long只新增到期member列。不carry、不補0；required-not-due不skip，到期失敗／無baseline／gap/reset分別標原因。顯示各member的[start,end)、partial coverage與next due，SQL commit時間另列；first partial、nullable schema、安全gate與durable驗收依support-per-member-recording-calculation。此語意確認不是新版demo或apply核准。
