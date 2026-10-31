## Context

基準 6f4a5fc；已核對 RuleEditor、RangeSummary、GuidedRecordingPreparation、DatalinkWorkbenchV2Page 與 useStudioV2AutosaveState。RuleEditor 的資料項目數不是 word 數，persisted MC 點位另有 compatibility 證據；新呈現不得替換協議驗證器。AppRoutes 與首頁 change 的 view 共用同一 V2 workspace scope。

## Design Source

2026-10-30 HTML demo v2 原樣快取 .spectra/design-cache/streamline-recording-task-editor/gateway-simple-ui-v2.dc.html；SHA256 d6987f1c8749db5d0ff9898544e3b2d505e98a99774b63d936a75daab2ffa306。實測畫布 #f6f5f0、卡片 #fffefa、文字 #263b3b、主色 #27766f；輸入圓角 7px、內距 10px 12px；卡片圓角 12px、內距 24px；段落間距 22px。原稿無命名 icon library，正式沿用 lucide-react 有標籤圖示。數字是來源 CSS 實測，不作效能保證。

原 demo 只驗證互動構想，不是正式 capability map。MC D-only、Modbus 40001-only、固定模擬值、所有點必填語意與單一全組 reduction 不能套成正式限制。補齊 demo 後需 owner 頁面斷點確認，再更新這份設計來源差異；沒有確認不開始產品 apply。

## Goals / Non-Goals

**Goals:** 完成來源→方式→目的地→檢查→必要影響確認→明確開始的單頁任務；已有設定能直接用，缺項就地補。

**Non-Goals:** 不新建 protocol/trigger/DB/durable draft engine，不更改量測算法、PLC地址規則或 Modbus Share 契約；不把所有跨設備群組轉成 basic profile；不建通用 SQL history viewer，不收集 demo 密碼。

## Decisions

### 就地補設備與點位但共用原資料權威

記錄 editor 保持未提交表單意圖；置中 modal重用 ConnectionConfigForm、RuleEditor 的欄位與validator，先保留本地 pending buffer，按明確保存才 dispatch 同一既有 autosave action/queue；不建立新的持久化引擎。取消只捨棄尚未送出的 buffer。新增完成後等待可證明的 persisted source/mapping identity，再供明確勾選及加入；不能先製造 fake IDs。顯示資料項目數、解碼型態及實際占用暫存器跨度，讀取週期與記錄週期分開。既有 word order、scale、單位、地址空間與型態能力原樣保留於詳細設定。

### 同頁漸進引導與能力保真

基本流程優先原始 typed Tag 最新值及既有預設；瞬時平均或累積用量由明確目的選擇進入已存在 profile/semantics。advanced/raw/cross-device group 若基本模型不能表示，顯示唯讀摘要與帶 scope 的進階入口；不默默拆分、覆寫或套 demo 規則。資料表長／寬與量測語意分開。預覽欄位和示例值清楚區分真讀值與示意，不把預覽當 SQL 已寫入。

### 草稿返回與確認失效契約

重用現有 autosave/settlement；已保存草稿與尚在頁面未保存的欄位有區別。只有後端可恢復的內容可標示可重新載入；未提交表單離開時提示保留/捨棄，不假稱跨 reload 持久。關閉、Escape、Back、Forward、切設備／group 與回首頁均經同一 dirty guard；保留所有已填欄位、點選與原位置，不只記錄名稱。已送出的非同步操作關閉視窗仍存在，回來查狀態。保存後應等 settlement；失敗或 unknown 明示而非稱取消。已成功保存的設備/點位不因關閉記錄editor或捨棄尚未保存的記錄form而回滾；UI在離開前明示哪些來源變更已保存。

### 重用檢查及明確開始

顯示 save、檢查、schema impact 與 start 的實際階段，同一主 CTA 不等於合併後端安全 gate。schema review 的 table/columns/layout/revision 依服務 evidence；修改任何相關輸入清除過期 confirmation，再重新檢查。Start 使用既有 scoped ownership/readiness/CAS，不要求 workspace-wide 啟用；started 只轉 waiting，不轉 success。

## Implementation Contract

共同 apply 前置：四案規格完成審查、補齊 demo，並已取得 owner 對本次固定版本關鍵頁面／路徑的確認。此 gate 在第一案產品實作前執行，不等到第四案；demo確認不代表正式瀏覽器或實機驗收通過。

- 完整 editor 整合前置仍為 simplify-recording-home-guidance 完成且 demo/owner 斷點已記錄；僅 tasks 1.1/1.2 的未掛載 foundation 可依下列分期先行。新 recording-editor 接現有 hooks/service，source identity 採正式 workspace/device/rule/point/tag ID，不採名稱或瀏覽器陣列 index。
- 分期依據：owner 正式實作決定見 docs/verification/simple-recording/implementation-decision-2026-10-30.md；2026-10-30 18:49 UTC 訊息 Sentinel_77b0a35bf0c48191aa14aff6c65cfe3a 明確要求「7個change」「多個agent平行實作，不可重覆」。本分期只取代原「home 全案完成後才可開始任何 editor 工作」的順序，不更動功能範圍、task IDs 或終驗 gate。
- Foundation 前置為上述固定版本 owner acceptance 已記錄及 home tasks 1.1–2.2 完成。允許以 Red-Green-Refactor 在獨立 recording-editor 新檔建立純 UI state、picker/table/footer 內容及相關測試：同表 filters、獨立 pendingAdd、既有 ID locked、Cancel/Add N、失效 ID 阻擋、remove/undo/re-add 保留原 member 意圖與非法文字、單一 active view stack 及無第一設備默選。1.2 的此局部可在 1.1 對應 foundation 測試契約完成後先行；不是解除其完整任務的 [after: 1.1]。
- Foundation 只接受既有 canonical projection/type 的 props 與明確事件 callbacks，不掛 production route、不新增 API/service mutation owner、不建立第二份 durable draft 或共享 modal primitive。候選資料的 device/point/tag ID 不等於 source/rule/mapping revision 證明；未知語意不得由名稱或 datatype 推猜，未知能力不可假成功。保存分步／不回滾與 late generation 可先測狀態契約，但不能以測試 callback 當作已完成正式來源保存。
- 整合階段才接既有 device/rule action queue 與 matching settlement/readback，保留已保存來源、pending/unknown 查證及原 validator；不得將全域 barrier 清空當作本次保存成功。共享 tokens、ThemeProvider 與 modal shell 由 home 唯一 owner 提供，foundation 不修改其共享檔；完整掛載及 1.1/1.2 勾選仍等 home 完成、整合相關測試及原任務驗收。先行部分只記局部證據，不宣稱任務或 browser 完成。
- Tasks 1.2/1.3 驗收分工澄清：1.2交付唯一picker/member draft、同modal view stack、設備保存及已保存設備不回滾、零選取／既有空設備進range buffer與取消完整返回；不要求1.3的真range persistence先完成。1.3承接此buffer，集中驗收matching submission、rule與mapping/point/tag讀回、新點定位、成功才關閉、失敗／unknown保輸入、防重送／late generation，以及已存range後取消picker／捨棄recording不回滾。這只消除重複造成的循環驗收，不刪使用者情境、不改task IDs／數量／checkbox／[after:]與既有home/backend/browser前置。
- 上述matching是既有queue的提交snapshot、scope/session/generation觀察與canonical readback，不新增device/rule API未提供的expectedRevision或operation ledger。讀回目前設定吻合不等於證明某未知request的因果結局；無法查明仍保留unknown與輸入，不盲重送或宣稱取消／成功。
- Batching 預覽、逐點字序及逐 member 計算仍分別等待 tasks 1.4、1.6、2.1 的正式 backend/capability 前置；source 保存、schema consent、scoped start、browser 修復重驗及 adoption gate 均不放寬。開發只跑相關測試，七案完成才按既有決定跑全套。
- 接口：editor 接收明確 create/edit 意圖及 scoped group identity；modal 回傳 persisted identity 或取消結果，未保存 connection 不能冒充 ready。既有設備連線變更依現行服務使 probe 過期。
- 任何無法由既有 API 表達的操作顯示能力原因與進階交接，不新增後端能力以填滿 demo；如必要業務流程確有 gap，回報 owner 修規格後才實作。
- 點數驗收使用設計測試例：10 float32 項需20個16-bit word；正式協議支援集仍由既有 validator決定。重複名稱與不同設備相同位址不合併。
- 草稿選測涵蓋置中 modal取消／保存、autosave失敗、reload、返回、stale response、異動後 confirmation invalidation、double click；格式合法不代表設備已連線。
- 有關 source/rule/mapping 的既有 related tests 必須一起跑；本批 lint/typecheck與line gate，不廣跑無關 DB engine。browser 視覺/歷史導航證據缺少時記 NOT RUN，不把 jsdom 代替它。

## Risks / Trade-offs

- [置中 modal產生第二份 autosave 狀態] → 共享同一 action/settlement，UI form 與 persisted projection 明確區分。
- [為追求簡單縮減正式協議] → 詳細欄位可展開，capability tests 覆蓋原 validator，不搬 demo 限制。
- [保存就誤啟動或一次按鈕繞過 schema consent] → 每個 gate 保留、輸入變更使相關 evidence 失效。
- [既有複雜記錄無法用基本模型] → 可讀摘要與精確進階定位，不自動轉換。

## 2026-10-30 需求更新（尚待版面選項確認）

根目錄 `DESIGN.md` 為唯一人讀設計規範；依 Taste app-screens 與 designtoken.md 的 primitive → semantic → component 分層，沿用既有 ThemeProvider 與字型，不另建主題系統。舊 demo 的淺色色值僅是歷史快取證據，不再是實作目標。

新簡易介面的彈出設定改為置中 modal，modal 本身寬度為目前 browser viewport 的 75%、高度為 75%（75vw × 75dvh），內容內捲；主頁仍正常 RWD，不是頁面縮為75%。已選手機近滿版與小型確認較小例外，尺寸分類依根 DESIGN.md；一般設定仍75% viewport。固定header/footer、contexthelp與同scope返回保留搜尋／篩選／捲動納入；選項8已選為dev/test專用樣式驗收頁，不放主選單，3/5/6/7未選，不新增tasks。文案簡潔、標籤保留意義，繁中／英文與數值／單位不能因版面省略。

共同 owner gate 必須用這次修訂後固定 demo 版本；過去 demo 檢視不視為確認。現處需求收集階段，未 apply、未做產品驗收；確定的協議保真與主題要求先納入，未選建議不擴 scope。

### 既有協議與參數保真矩陣

來源：state/protocols.ts、protocolSupport.ts、ConnectionConfigForm、studioV2DeviceAutosave.toBackendConnectionConfig、addressParser、internal/datalink/connector/registry.go、schema_device_models.go、sourcerule/protocol_address_planner.go；可達不等於此次硬體測試通過。

| 協議 | 必要／基本欄位 | 進階保留與驗證 | 簡易呈現 |
| --- | --- | --- | --- |
| Modbus TCP/UDP | host/port/slave_id，timeout | data_format；0/1/3/4五至六位位址，非僅40001 | 基本可編輯、進階展開 |
| Modbus RTU | serial_port/baud_rate/slave_id/parity | data_bits7/8、stop_bits1/2、timeout、data_format；port屬backend主機 | 全欄位可編輯；baud依server enum，不以舊UI2400/4800擴允許值 |
| FATEK FBs | mode/station_no；TCP host/port或serial_port/baud_rate/parity | serial data_bits/stop_bits/timeout，現存serial不能默轉TCP；X/Y/M/S/T/C/D/R | 新modal顯示mode分支；既有值round-trip |
| MC 3E | host/port/station/timeout | network_no/pc_no/io_no/station_no/data_format與aliases；D/M十進位、X/Y八進位、W/B十六進位；persisted compatibility fail closed | 原值優先，未設port提示可改預設；不聲稱6000或5000普遍正確 |
| MQTT | picker disabled | schema存在但V2 register planner不支援 | 明示unsupported，不收真密碼、不假稱可記錄 |

欄位只有「簡易可編輯／進階可達且值保留／明確unsupported」三類，不可消失。每類必需與advanced以真正server validation為準；新增表單只是投影既有schema，不新增driver或planner。protocol切換僅合法draft並提示重設，不改persisted identity。更新一欄須保留未呈現原始欄位與alias正規化。

Rule datatype保留bool/int16/uint16/int32/uint32/int64/uint64/float32/float64/string，按protocol/area相容性驗證；count1–64是項目數。保留naming_prefix/enabled/skipped_addresses/scale_multiplier/scale_offset/data_format與Share的share_enabled/start_register/stride及ownership gate。unit/scale/offset/target_type在mapping層，不誤寫成Rule欄位。

正式Rule/create/update沒有每範圍poll；demo原per-block poll撤成目前排程來源唯讀摘要或既有settings/PollingGroup進階交接，不能新增per-rule scheduling。採集排程與SQL寫入窗分別呈現。demo缺參數或僅模擬驗證應標明，不能等產品寫code才發現丟能力。

### 已選 modal 與返回引導契約

依 DESIGN.md 尺碼分類：一般popup75vw×75dvh；mobile<640px或short viewport<480px近滿版；小型短確認較小autoheight。第二主畫面editor維持RWD，不是popup。modal固定header/footer、body內捲；criticalrisk/validation/consent不能藏contexthelp。返回搜尋/篩選/捲動僅同scope、entity與view，revision變更重驗identity，不能復活舊selection/availability/receipt。dirtyguard、focus與pending操作語意不變。選項8已選為dev/test專用樣式驗收頁，不放主選單，未選候選不入實作；此需求更新不是owner對新版demo的視覺核准。

### 同協議最近連線參數只作新表單預填

使用者新增需求是短時間連續建立設備；本案以同workspace、目前掛載頁面的memory session界定，不臆造5或30分鐘TTL。連續流程的「再新增設備」可預選上一筆成功保存的協議／子型（仍可切換）；只在初始化新form且沒有明確protocol或已輸入欄位時預填，late save不得覆蓋已開form的選擇或文字。新設備modal開啟時可讀該協議最近一次成功persisted save的allowlisted connection config；必須顯示「已帶入上次設定」及「請確認IP／串列埠，避免重複連到同台」，可修改後再明確保存。不是複製設備／probe／點位；沒有新identity前不假稱建立，不autoactivate。

只收目前協議schema明確允許的非敏感連線參數：TCP/UDP host、port、slave_id、timeout、data_format；RTU serial_port、baud_rate、slave_id、parity、data_bits、stop_bits、timeout、data_format；FATEK mode、station_no及該mode適用TCP或serial參數；MC host、port、network_no、pc_no、io_no、station_no、timeout、data_format，沿既有alias正規化。未知server keys及任何credential/token/password不入cache；device name/ID、rule/point/mapping、enabled、readiness/receipt/operation/ownership proof一律不帶。MQTT目前disabled不建立cache能力。

cache key為workspace identity＋protocol＋必要subtype（FATEK TCP/serial）；跨協議不混用，mode切換不得留不相容參數。僅matching scope/session epoch、protocol/submode與request generation的成功request snapshot或server canonical response經whitelist可更新；不能watch save_state=saved後讀當下form，因useStudioV2AutosaveState可能保留尚未保存的latestDevice.config。runtime_apply_status=apply_failed亦不更新。clear遞增cache generation，晚到舊response不得復活已清cache。僅已確定保存成功的回覆可更新，cancel、failed、pending/unknown與late舊scope回覆不更新；為避免其他編輯意外覆蓋近期值，此cache只由簡易新設備流程成功保存更新，不由背景載入既有設備或編輯舊設備灌入。

「清除帶入值」刪除目前protocol的所有subtype cache，將目前新表單connection config重設為所選protocol/mode的正式預設（無預設留空），清除帶入提示；不更動新表單名稱、不刪除／修改任何既有設備。FATEK仍維持使用者當下選定mode但清掉另一mode不相容欄位。已被使用者修改的未保存連線內容要經既有dirty確認後才重設；取消清除則原form/cache保持。清除不自動保存／probe／啟用。

重新載入、離開/卸載整個Studio頁面或workspace切換清空；同頁關閉新增modal不清除上一筆已成功的cache，但取消此次輸入不更新cache。使用者若只關modal再新增仍屬同page session。不寫localStorage/sessionStorage、不新增設定API、後端或長期偏好。記錄editor草稿與source保存的既有安全分界不變。

### 每個參數都有可操作的圓圈說明

沿根 DESIGN.md，input/select/checkbox等設定參數旁ⓘbutton（type=button）以click/touch/Enter/Space開輕量popover，不是75%設定modal。用途、unit/range、例子、何時修改四要素依需要提供，值域取schema/validator與protocol source，不臆造通用安全值。label／必填／validation／disabled reason／criticalrisk／consent始終inline。

Escape或點外可關，鍵盤關閉回trigger focus，點外不搶新目標焦點；不submit、不mutate、不清草稿，同light/dark tokens且窄viewport不clip。所有實際呈現參數需coverage map（canonicalfield→內容source→demo位置→formal差異）；未呈現正式欄位只能精確標parity／進階handoff，不冒稱完成。既有協議matrix為能力authority，help不放寬validator。

## v4 owner 未驗收

v4配色、密度與部分行為被使用者否決；既有mock測試不是UI確認。本輪先修demo及規格，修訂固定版本重新確認後才可apply；已以DESIGN.md Quiet mineral候選替換舊blue/slate；新版仍待owner視覺驗收，不當核准品牌。

### 從已設定目的地安全選表

來源GET /studio-v2/workspace/write-groups、useWriteGroupsQuery及WriteGroupDestination，候選identity是workspace＋connector_id＋database＋table_schema＋table_name，包含storage_strategy/layout與顯示用目前revision。用已有配置整理而非live DB inventory，不列未配置外部table，也不索取真密碼。名稱相同不同DB/connector不可合併。

相同scope/layout的managed長表依既有write_group_long_storage scope可以共用，table名稱由canonical scope導出，不允許任意覆寫；寬表owner是workspace+groupID，因此新group選已歸屬寬表只能顯示原因／導原記錄編輯，不暗放寬跨groupwriter。無法重用候選仍可見具體原因及安全下一步。選候選只填目的地intent，不複製舊group identity/schema proof/readiness/digest能力；必須fresh connector/schema/revision/ownership review，變動使consent失效。

### 緊湊表單與短確認分類

**Supersedes**: streamline-recording-task-editor / 同頁漸進引導與能力保真

保留漸進任務與advanced能力交接，但設備常用連線參數常顯，不為減少畫面而折疊；點位compact rows而非大卡，name/address/readonly datatype/perpoint wordorder/semantic常顯，位址跨度與每點有效字序容易比較；generic dialog不設尺寸，以data-size明確分類，form75%、confirm auto/fit-content及auto grid rows。低頻細節可折疊，criticalrisk/validation不藏。首頁與editor視覺先修demo，再重新owner checkpoint。

「尚未套用的修改」與stop/remove/discard按語意共用short-confirm分類，auto/fit-content且清native inset stretching，不因dialog id不同套75%regular。限制max-height與body scroll，header/footer不留下巨量空白；實際layout需browser證據，mock DOM不是驗收。

### 每點字序依正式前置能力呈現

前置change `preserve-per-point-word-order`、capability `per-point-word-order` 為authority；canonical Point是唯一format intent authority，ruleLink只作ownership，不建立第二份link override。本案不自行新增持久化或decoder。Point format intent and effective order、Preserve overrides across rule revisions、Protocol decoder parity、Preserve legacy decoding semantics、Scope and revision guarded updates、Version decoding identities and observation proofs 由該案負責。既有Point.data_format並不足：source-rule目前更新覆蓋generated points，MC/FATEK原runtime不採request format，不能下拉就稱完成。

API mode精確為legacy/inherit/override，legacy-preserved僅人讀label；新點inherit，explicit值即使等於default仍是override。每點行投影API readback的mode/configured/effective/source/applicable/revisions，而非只看storedformat猜effective；使用者可設定單點、reset回明確inherit，調整range default不得覆已有individualoverride。datatype／protocol適用選項來自該前置capability matrix，不假設所有多word型別都是ABCD/CDAB四byte；單word/bit不提供新的multiword override；隱藏selector不得強制ABCD、清data_format或破壞既有driver byte-swap policy。舊MC/FATEK維持driver legacy政策，不能把既存Point.data_format當有效override。legacy preserved不可偷偷稱inherit；preview已知rawfixture與正式decode證據分層，無PLC授權不連設備。

保存使用前置scoped/revision API；stale revision、point已刪／重建、unsupported格式均明確受阻，不用local field假成功。前置案未完成時demo只能標示目標互動／mock，正式editor task不可跳過依賴。

## v6 易懂記錄方式與來源選取

### 以使用目的說明記錄結果

**Supersedes**: streamline-recording-task-editor / 同頁漸進引導與能力保真

只取代不易懂／過度概括的呈現，不改引擎：主問題「這組資料要怎麼記？」；選項依正式能力顯示「保留讀到的值」、「把這段時間的量測值取平均」、「計算這段時間增加多少」。示意20/22/24→最新24或平均22；累積電表1000→1008→增加8。例子明示示範，不當SQL receipt。

已核對measurement/mean.go、snapshot/mean_test.go、runtime/group_boundary_long.go longMemberShape、aggregation/profile_counter.go及workspace write-group reduction contracts。平均是各有效樣本等權算術平均，不是跨點平均或時間加權；缺失/拒絕樣本與coverage沿既有evidence，不補0。profile有group reduction不代表每個member都同算法：mean的gauge會平均，counter_raw/state仍latest，counter_raw是原累積讀值不是用量；不能寫「同組所有點都平均」或「後端不支援任何混合」，這是現有legacy行為描述；新明確逐record-member需求由下述獨立前置提供，不把legacy限制當新功能限制。

用量依已確認counter policy及既有可靠baseline計算，首窗無baseline為NULL不是0；不確定連續性、reset或無法解釋下降顯示原因，已確認signed/rollover政策可合法解釋變化，不能一律負差視缺值。跨缺窗是具真實跨度的total，不冒充固定單窗。每點摘要顯示其角色與實際結果用途；UI不能表示的現有profile保留摘要＋進階交接，不重寫／降級。原saved reduction原樣保留，legacy選項只投影真正API能力，新mixed在前置完成後才開啟。

### 範圍保存成功才關閉並回到原位置

範圍modal明確Save捕捉submitted buffer generation，等待相同scope/request/modal session及generation的persisted成功與必要identity；若pending期間有新編輯，舊request成功只能settle已送部分，保留新dirty欄位且不關modal、不清新dirty，明示仍有未保存修改。所有識別仍匹配的成功才清本次buffer dirty、關modal，刷新來源並回原記錄草稿／選擇位置與焦點，不再跳「尚未套用」確認。驗證失敗、save failure或unknown仍開著保留所有輸入與安全原因；pending防重送。較舊回覆不得關掉新開modal或新scope。關閉不代表開始採集／SQL成功，已保存source也不因稍後捨棄recording而rollback。

### 編輯先看已選來源再按需搜尋新增

主頁僅顯示已加入點及每member method／interval，唯一「選擇點位」入口開75%置中modal，不再以新增／更換直接修改membership。搜尋、設備、點位類型（量測值／累積值／狀態值）為同一工具列，立即共同過濾同一結果表／符合數，query/filter變更回第一頁；大量列表25/50分頁與canonical-ID穩定排序，短列表不硬加分頁。既有點標「已加入」不可取消或重複加入；pending跨query/filter/page保留。主頁移除與復原不刪來源，復原或同ID重加保留原member方法／間隔及未通過驗證的文字。

選取由canonical workspace/device/rule/point/tag identity維持；跨查詢不可見不等於刪選。source實際刪除或revision失效保留可辨識的失效項與修復理由，阻擋未驗證加入／保存，不以同名或同址新ID取代。切workspace隔離state，同scope返回恢復query/filter/page/scroll/focus。搜尋只使用既有已載入projection；loading/error/no-results/empty分開。100設備/300規則與單設備多頁是測試fixture，不作runtime效能宣稱。

### 加法選點與來源保存分層

**Supersedes**: streamline-recording-task-editor / 編輯先看已選來源再按需搜尋新增

取代原決策中「新增／更換來源」及可在原頁混排的呈現選項；保留其canonical identity、有界發現及返回契約，以本節加法modal為唯一互動方式。沿現有tasks 1.2/1.3修訂驗收，沒有取消已完成工作、沒有新增任務。

狀態分為record draft members、picker pendingAdd、source form buffer與已送出operation；皆沿現有authority而非新增後端draft engine。開picker不改record draft；「加入 N 個點位」重驗scope/revision後只合併pending新ID，只初始化真正新成員。既有或同session已移除後重加成員保留原method／interval意圖，不覆蓋非法輸入文字或行內錯誤。pending有失效項時清楚阻擋並供處理，不能靜默忽略後宣稱全數加入。主頁remove／undo記憶僅目前record draft scope，不跨workspace或冒稱reload可恢復。

固定footer顯「待加入 N 點」（另列目前結果外pending數）、取消及唯一primary「加入 N 個點位」，N=0禁用；列表摘要為符合X點／Y已加入，數值皆由對應canonical集合推導。取消／Close／Escape丟棄本次pending且回入口焦點，不觸發整份記錄捨棄；背景點擊不默默關閉。source未送出buffer仍用dirty guard；短確認替換active內容而非疊dialog，拒絕捨棄返回原form。主頁remove可復原，picker不能移除既有成員。

只有「找不到點位？新增來源」次入口與對應empty-state替代動作，主picker沒有第二個目標設備selector或逐列設定。單設備filter進setup直接標唯讀設備；全部設備時setup明確選既有設備／新增設備，不預設第一台。query/filter只改結果，不能偷偷換建立目標。新設備按「儲存設備並設定點位」才保存；下一步顯「設備已儲存，尚未建立點位」，範圍按「儲存點位並返回選取」才保存，旁說來源保留、加入此記錄稍後選。兩次保存各用既有settlement／revision，非原子交易；取消後一步不刪已保存設備，取消picker不刪已保存來源。unknown/pending查狀態、不重送，不把未決叫取消。

視覺依根DESIGN.md的使用者任務層級：同高工具列、identity最大欄、結果優先、次入口降階、固定footer結合計數與加入；點位類型與資料型別分開。此輪不換palette、不增新引擎或宣稱prototype已接後端；正式產品apply仍待新版owner checkpoint。

### 空來源建立後回到可選點位

補充「編輯先看已選來源再按需搜尋新增」：新增記錄的零選取不是死路；單一次要來源流程提供新增設備，已載入但尚無點位的既有設備仍可篩選並提供在該設備「新增讀取範圍」，不要求先有point才能進入範圍設定。無匹配結果與確實無點位分開，載入失敗不能偽裝空資料。

設備→範圍→來源選擇以同一modal view stack切換，不疊modal。進入前保留recording名稱、目的地、所有未送出欄位、canonical draft members、picker pendingAdd與每member的method/interval及未通過驗證文字等設定，以及chooser query/filter/page/scroll/focus。取消設備／範圍只捨棄該步驟未送出buffer，原草稿與selection不變，依dirty guard返回原chooser位置。

匹配成功的新增範圍保存（含已有多頁點位的設備）沿既有settlement關閉其設定步驟，返回bounded discovery而非selected-only空畫面；定位該device／range並清除會遮住新點的舊query／點位類型filter並明示原因、重設至新點所在頁，使已確認identity的新點可直接勾選，不自動勾整台／整個range。這是成功建立後的明確定位例外，不改一般編輯既有範圍或取消時恢復原query/filter/page的契約；原selection與member設定不被重建或預設覆蓋。新增設備成功但尚無range時，提供立即新增範圍與可返回的空設備狀態，不假裝已有點。

驗收依Implementation Contract分工：1.2涵蓋零選取新記錄、既有空設備進range buffer、已選member有自訂設定時往返／取消；1.3涵蓋既有空設備及new-device→range-save、正式新點讀回、舊query／點位類型遮擋後定位，以及已保存range後捨棄記錄仍保留source。失敗／unknown／late generation及不關窗／查證亦由1.3驗收，不要求1.2預先完成真range-save。v6尚未取得owner確認；此補充不新增change／引擎／任務、不開始產品apply。

v5目前未核准；以上修正現在完成demo再給owner確認，不留正式版才修，不改根DESIGN視覺方向。

### 同一記錄每個成員選擇計算結果

**Supersedes**: streamline-recording-task-editor / 以使用目的說明記錄結果

保留白話label/例子、真實legacy說明與未知值保護；取代單一group method作新記錄唯一選項的呈現。前置 `support-per-member-recording-calculation` 完成後，已選來源每列提供「存最後讀值／取平均／此時間區間增加量」，method屬recording member/mapping，不屬來源Point；ABCD/CDAB是來源解碼，計算是此記錄要存的結果，兩欄不同。相同source在另一record可有不同method。新member預設latest，群組捷徑只初始化新增member、不覆explicit。

新mixed是明確versioned shape：同group各member自己的interval（second/minute正整數、依backend bounded capability）及UTC epoch對齊到期邊界，每member獨立結果，managedwide同一列可為22/12/8。新mixed的delta欄是增加量，legacy cumulative欄仍raw＋provenance，不默轉。每列顯configured/effective/resulttype/applicable及缺失原因；mean僅numeric gauge，delta需confirmedcounter，unsupported組合disabled且給原因。不提供global Point method，不改semanticdefinition來硬過validator。

required delta首窗無baseline→wide整列skip／long該member列skip但baseline原子進展，optional無結果→schema允許才NULL+reason；gapspan不冒充此窗增加量。改method清舊readiness/consent，fresh revision/schema/ownership；legacy變mixed若需新表沿既有successor。預覽是示範，不代表已runtime可用；前置未完成只明示mock目標，不正式啟用。


本輪已確認逐member interval與UTC epoch對齊：5秒00/05…55、15秒00/15/30/45、10分鐘00/10…50；wide採due邊界聯集，:15/:30/:45的60秒delta留NULL＋not_due，:00同列各自前15秒mean／前60秒delta；long只新增到期member列。不carry、不補0；required-not-due不skip，到期失敗／無baseline／gap/reset分別標原因。顯示各member的[start,end)、partial coverage與next due，SQL commit時間另列；first partial、nullable schema、安全gate與durable驗收依support-per-member-recording-calculation。此語意確認不是新版demo或apply核准。

### 同源採集次數預覽交接

integrate-protocol-aware-acquisition-batching 唯一負責planner、preview API/shared types與實際計數來源；本案僅接既有consumer呈現／驗收，不另推導合批、不新增採集引擎。該案production gate完成後才顯正式「每輪計畫協議讀取次數（不含故障重試）」與estimate標示；顯示acquisition point count、unique wire units、planned requests及安全coverage/fallback原因，不能用recording已選點數當整設備採集集合或自行猜random=1。API預覽零PLC I/O，沿workspace/device scope與expected revision，stale409／載入失敗清除可用性並提示重新取得，不展示過期成功。

actual attempted/succeeded/failed/retry/fallback/canceled及cycle latency使用另列同scope/revision的真實觀測；缺乏實測就顯尚無實測，不補0、不將估算稱TCP封包數、SQL寫入數或已證明加速。demo只可明標模擬契約，無須也不能假造真實性能。獨立member recording window仍由support-per-member-recording-calculation負責；更改selector不改採集plan，批次預覽依賴不升為member排程的硬性前置。最終adoption必須分開列預覽契約測試與實際transport量測證據，保留新版demo owner與apply gate。
