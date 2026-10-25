## Context

本案是後續功能提案，不是前案重新驗收。前置 add-circuit-recording-profiles 的產品已提交且 tasks 為 15/15，但尚未 archive；其 delta specs 仍是交接依據。本案不改動前案或其完成紀錄。專案版本與接線以 AGENTS.md、go.mod、frontend/package.json、internal/api/router_studio_v2_recording_start.go 及目前原始碼為準；CLI context 中過時的 Go／Viper 敘述不作為改動理由。

已讀碼確認的現況：
- BasicRecordingPanel 的 hasLegacyBasicWide 對整份 liveGroups 判斷，任何舊群組都能讓其他設備回到舊面板。BasicRecordingMembers 預設展開完整清單；Step4Database 把群組編輯與 CommitSummary 的全域啟用並排於 Advanced。
- useCircuitProfileStart 發送 scoped recording-start；RecordingStartService.runRecordingStart 已逐組檢查、Apply，再 ActivateScope。GroupLifecycleBar 另有 Apply，而工作區 ActivateEligible 啟用所有 eligible 設備，不會替草稿群組完成 Apply。
- snapshot.finalizeProfile 與 Selector 現為最新合格樣本；既有 wide finalize 有 required/skip_row 規則。newLongGroupBoundary 只接 long profile，沒有 production 平均閉窗。
- GroupRowLayout／LongRowLayout 將輸出綁定來源型別；整數平均 1 與 2 必須產生 1.5，不能塞回整數 slot。已接受樣本沿既有 journal、CommitClosure、outbox／receipt 可重播；不能用另一個未接線 float mean helper 宣稱已有平均。
- 現有 snapshot Config／Member 沒有每窗 expected-count 契約，缺樣本也不能直接推論為斷線。使用者已選有效樣本算術平均，不採時間加權。

畫面截圖是問題證據，不是像素設計稿；沒有 .dc.html 或需快取的視覺交付來源。本案沿用現有 Studio 元件／樣式，不另做全站視覺設計。

## Goals / Non-Goals

**Goals:**
- 以「整台設備快照（寬表）」與「指定量測項目（長表）」表達用途，布局與瞬時／累積語意分離。
- 新建瞬時值預設 15 秒算術平均，保留 1–300 秒範圍；累積保留 1–600 分鐘與 10 分鐘預設。這些為需求設計條件，不是效能量測。
- 應用同一份 canonical 設定，提供設備方塊、缺值說明、收合明細、單一 scoped 開始及可查證的記錄狀態。
- 保護所有舊版本、正式來源身分、schema 確認、Share-only、耐久佇列、精確 counter 與歷史目的地。

**Non-Goals:**
- 不實作 CMS 報表引擎、查詢／匯出新 API、費率、告警、需量分析、圖表、跨日用量分攤或新 Runtime 頁；電錶／水／瓦斯只是本案混合來源與原值／增量的驗收資料。
- 不做時間加權、插值、補零、跨窗 carry、預測值、每筆永久保存、新採集器、改 polling／scheduler／協議、新健康監控或新 counter 算法。
- 不新增 MySQL／MariaDB／SQL Server 支援，不 unpark 其提案；不把寬表做成 Timescale、不安裝 extension、不改現有 7 日 chunk、retention、compression、分表或通用效能調校。
- 不自動合併電錶資料表、跨設備重編群組、不重寫歷史、不自動雙向轉表、不重做 cutover 引擎；需要布局交接時只擴充既有 reviewed/CAS cutover 的必要目標驗證。
- 不修改 Step 1–3、Tag Key／顯示名稱與既有 ⓘ 修正、launcher、toolchain、Windows tray、未關聯的 baseline failures；不新增依賴／框架、全域狀態或第二個 journal／operation ledger。
- 不因 reviewer 的建議自行增加功能、測試平台、效能矩陣或模組重構；不自動 commit、archive、push 或部署。

## Decisions

### 兩種用途共用語意與 canonical 設定

選擇「整台設備快照」建立新受管 wide profile；「指定量測項目」使用既有 managed_long_v1。新 wide 使用獨立 layout_version managed_wide_profile_v1，避免拿既有 managed_wide_v1 的整數欄與 latest 語意接平均結果。wide 的一列代表所選設備／迴路的一個 profile 時窗；long 的一列代表一來源的一個時窗。各 kind 的 interval 仍為唯一權威，不能在連線或 Advanced 多設一份秒數。

兩種布局均可建立瞬時與累積 kind；canonical key 沿用 workspace＋circuit 或 unassigned device scope＋kind，layout 不增加第二個同 kind owner。Ensure 對不同 intent 仍衝突，禁止改 circuit ID／刪除再建來繞過。建立後布局不可直接修改；跨布局只走明確準備與既有 reviewed/CAS 交接。unassigned 的既有限制保留，cumulative 仍須真實 circuit／measurement 綁定。

替代方案：只留長表不能服務固定多欄電錶快照；只新增寬表／長表名稱卻維持全域 legacy 分支仍會混淆；每設備新建一套表管理或 metric registry 超出範圍，均不採用。

### 明選平均與結果型別契約

**Supersedes**: add-circuit-recording-profiles / 瞬時選樣與累積精確計算

增加 persisted recording_profile.sample_reduction = latest | mean；缺值一律解讀 latest（向後相容），新建瞬時預設 mean。cumulative 不接受 mean。每個瞬時 member 的 value_role 明示為 gauge、counter_raw 或 state；numeric gauge 用 mean，已確認 counter_raw 及 bool/text state 用最新合格值。未確認 numeric 不由名稱／單位猜 counter；操作者選 gauge mean 是明確意圖，基本 raw Tag 不強迫建立計量定義；選 counter_raw／cumulative 才沿用 ConfirmCounter。

平均僅包含實際時窗內、身分／版本相符且依既有品質／freshness 可用的唯一 SampleID。bad、invalid、stale、晚到已關窗資料與重複投遞不加到分母；相同 ID 異內容仍衝突。使用 [start,end) 的 UTC／實際截斷界線，不跨窗加值。

exact int／uint／decimal 以精確和與整數 count 計算；結果型別獨立於 source_value_type，輸出 decimal。非終止小數採 18 位小數、half-even（本案設計條件），同列 provenance 保留精確 sum、count、輸出 scale／rounding，不能宣稱無限小數完全精確。float 來源依穩定 observed_at＋SampleID 順序算 float64 平均並標 basis=float64，不宣稱還原精度；禁止把 exact 大整數先轉 float64。

mean 的凍結 runtime layout 同時記 source_value_type 與 result_value_type。long 使用既有 value_decimal／value_float slot，source type 移入明確 metadata，不能靠改 Tag 型別滿足檢查；wide profile 對 mean 生成相容的數值／exact TEXT 欄，int 不截斷為整數。自訂／既有不相容欄位在準備前阻擋，不自動 ALTER 或猜轉型。

替代方案：用 latest 冒充平均、整數平均截斷、套時間權重或把累積原值平均化，均拒絕。

### 缺資料證據不杜撰完整率

每窗的 aggregation metadata 為 version、method、source_value_type、result_value_type、basis、valid_count、rejected_count、first_observed_at、last_observed_at、coverage = complete | partial | unknown | missing、reasons。僅 exact mean 另有 sum／scale／rounding。樣本數為去重後觀察數，不等於應採集次數；不輸出杜撰的 expected_count 或百分比。

現有 grouppipeline.sourceEvidence 僅投影 long，故必須在既有 delivery API 的 source_evidence 擴充平均 aggregation 的 method/count/coverage 與新 wide profile 的唯讀 member 投影。從凍結 payload／applied layout 與匹配 receipt 取證，維持 bounded candidate 與無原始數值的回應；不得查整段歷史、新增報表 endpoint 或拿草稿目的地當已提交證據。無效 proof／metadata 明確 unavailable，不能默默當 complete。

沒有有效值時 value=NULL、coverage=missing、原因依已接受證據為 no-samples／bad／stale／invalid；沒有斷線證據就不能叫斷線。有可用平均但含已知失敗、缺口或截斷時 coverage=partial；有有效值但無法證明整窗覆蓋時 coverage=unknown，仍可寫平均且明示未知。complete 只能由既有已保存證據證明，不能因沒有錯誤事件就成立；本案不新增推估採集排程或健康監測系統。這是保守且可驗收的界線。

新 mean wide 預設 optional＋partial，目的地必須可保存 NULL 與 per-member metadata；整窗全缺仍保存一列 NULL／原因。若操作者在詳細設定明確設 required，required 缺值仍按政策略過整列並記原因。long 各來源獨立，不顯示整列缺值選項。既有 wide／custom 的 required、skip_row、no_data 與不支援 partial 時的拒絕完全保留。

### 兩種布局重用累積與耐久閉窗

cumulative 原值、delta、baseline、精確數字、reset／signed／rollover／gap span total 全部重用 aggregation.CloseCounterWindow。wide profile 的來源 value 欄保存原值，既有 provenance 欄新增 per-member counter 物件保存 typed delta／區間／狀態／policy；long 保留現有 delta_* 欄。不得用 JSON Number 傳 exact 值，不額外產生一套報表／counter 表。完整電錶的瞬時 wide 中 counter_raw 只記最新原值；同來源另選 cumulative kind 才計算該 kind 的增量，不能以平均結果當基準。

平均的 accepted samples、counts 與結果使用既有 journal 回放，closure 與 consumed checkpoint／outbox 同交易。sender 只重送凍結 payload，不重算平均或 counter；SQL ACK 遺失仍一個 effect。沒有擴充後 runtime proof 的新群組不能 Apply。沿用既有 backlog／backpressure 上限，不要求新效能指標或無界暫存。

### 設備方塊與單一 scoped 開始

Step 4 使用設備名稱小方塊單選切換，無設備下拉；以 persisted device ID 為權威，名稱可改、同位址設備不混。每個 workspace＋device＋用途／circuit 的草稿獨立；裝置切換、reload、舊非同步回覆均不得污染其他 scope。選取設備只是篩選，不啟停、不勾選所有來源、不 Apply。跨設備 custom 群組保留 Advanced 檢視，不自動拆分。

主流程共用目的地卡 → 設備方塊／用途 → 型態與週期／來源 → schema preview／confirm → 套用並開始記錄。準備明確 DDL 確認仍與開始不同。工作區總覽沒有第二個啟用按鈕；Advanced 仍儲存相同草稿與提供明確停用／刪除；本案受管記錄的生效動作回同一個 scoped review/start，而非另發 ActivateEligible。既有 cross-device custom 的工程 Apply／生命週期不重寫，也不把它偷偷包進單設備開始範圍。既有全域啟用 API 留給相容呼叫端，不移除 endpoint。Share-only 由同一主入口走既有零群組分支及 Share barriers。

開始請求凍結目前設備與所選 groups；不把「查看全部狀態」或選設備當作「啟用全部」。不新增多設備批次啟用功能。來源清單預設 summary＋數量，展開只檢視目前設備，超出容器則內部捲動；異常計數與修復入口不因收合消失。一般 DOM 使用顯示名稱／位址，技術 key／revision 留診斷。

### 有界交接與舊資料遷移

展開 → 移轉 → 收斂：先新增 nullable 設定／解碼及不開放 Apply 的新 layout；再接平均、typed encoders、兩種布局與既有 counter／durability；最後接兩條用途 UI、設備方塊、單一啟動並移除主線的重複啟用入口。每批可建置、有 focused tests。

舊 sample_reduction 缺值仍 latest；讀取、切換畫面或 reload 不更新設定。既有 latest long 明確改 mean 需新 revision、preview／confirm 及 Apply；舊 wide 改新 wide profile 需既有 reviewed cutover 延伸到該 target，絕不覆寫舊表，歷史與 pending 沿原目的地排空。舊寬表只影響它自己的設備／scope，不再封鎖其他設備選新長表。布局不可直接改，已在記錄的同 kind/source owner 仍需明確交接。

回復時停止新 intake 並排空 frozen outbox；恢復已保存舊版本需同一套明確 review，不刪表／歷史／ledger。尚不理解新 layout 的二進位必須安全拒絕該新群組，不能將它當 legacy wide 啟動；不能承諾舊二進位讀新 payload 的向下相容。

## Implementation Contract

| 可驗收行為／介面 | 失敗／限制 | 指定證據 |
| --- | --- | --- |
| sample_reduction、value_role、source/result type 存取與版本凍結 | unknown enum／mean+counter／不相容欄位先拒絕；缺欄讀 latest | TestRecordingReductionCompatibility、TestMeanResultTypeProof |
| [00:00:00,00:00:15) 的 10、20、30 平均 20；1、2 平均 1.5；邊界樣本歸下窗 | bad／重送不加分母；全缺 NULL；不填舊值 | TestWindowArithmeticMean、TestMeanWindowEvidence |
| exact 9007199254740992 與 9007199254740993 的平均 9007199254740992.5；0、0、1 為 0.333333333333333333 | 精確 sum/count 可查；float basis 另標；原始 Tag 型別不改 | TestMeanExactProjection、SQL readback |
| wide／long 的每來源平均及 counter 結果邏輯一致 | required／custom 舊規則仍成立；wide counter delta 在 provenance | TestWideProfileRecording、TestWideCounterParity |
| journal／closure／SQL ACK 故障後相同平均、count、counter 與 effect | 未提交不 consume；已凍結不重算；目的地離線仍有 pending | TestAverageClosureRecovery、TestAverageCommitAckLost |
| device A/B 方塊切換、重名同位址、草稿與異步scope | 不呼叫 create/update/apply/activate；沒有設備 select | device-recording-tiles.test.tsx、瀏覽器真實鍵盤 |
| 單一開始只作用目前設備／已審核群組，工作區只總覽 | 相同request復原、stale拒絕、失敗不影響B、Share gates不減 | recording-workflow-unified.test.tsx、TestUnifiedRecordingStartScope |
| 舊寬表60秒、custom、latest long與新平均共存 | 不能改歷史、默默換表或改成員／必填 | TestRecordingWorkflowCoexistence、切換前後SQL／版本證據 |

需求數字、例值與邊界時間皆為驗收設計條件。驗證 viewport 1024×720 與 1600×1000 為既有測試規格的選定條件，不宣稱所有尺寸；需要的元件尺寸沿用既有設計 token，不新增像素設計目標。

禁止外擴的執行規則：
1. 每個修改與測試必須對應本表、spec scenario 或本案造成的回歸，交接需指出對應項。
2. 必要依賴才可新增檔案／nullable migration／typed metadata；不得因「比較完整、較安全、未來可能需要」建立新服務或橫向整理。
3. 若需新增本規格受管目的地表與 additive 設定 migration 以外的持久化資料表、改採集或斷線事件系統、變更既有來源身分、加入新DB／報表、改歷史、超出上述輸出契約，先列最小方案與原驗收失敗證據，向使用者取得授權；只暫停依賴該決定的部分。
4. reviewer、測試產生的範圍外發現只簡報；既有 baseline 失敗先比對來源，不藉此接成任務。代理不得自行增加任務或驗收平台。
5. 每項估計 30–60 分鐘，這是待實作起迄／focused checks 驗證的估算，不是承諾。超出一個 context 可完成的切片時，只拆同一驗收的剩餘工作並保留界線，不能追加需求。
6. 所列需求與必要檢查通過即停止；未執行明列 NOT RUN；無授權不 commit／archive／push／部署。本輪僅規劃，所有產品任務保持未勾。

## Risks / Trade-offs

- [沒有可靠應採集次數／連線事件] → 標 coverage unknown／missing，不新增監測子系統以追求漂亮百分比；不以缺樣本直接宣稱斷線。
- [平均造成型別及資料意義改變] → source/result 分離、exact sum/count＋明確 rounding、僅新版本；既有整數與 counter 不轉 float。
- [新寬表的 counter 查詢需讀 JSON metadata] → 沿用 provenance 以限制範圍，文件示例說明原值與 delta 位置；另加固定 delta 欄或報表層需新授權。
- [兩種布局導致 UI 再分裂] → 只分用途與對應規則，共用同一開始、目的地、狀態與 Advanced 權威。
- [舊 specs／前案 newest 與全域啟用描述衝突] → 本案四份 delta 明確保留 legacy latest、明選 mean 例外及全域 API 相容；不直接重寫已完成 change 或 base specs。
- [效能與資料量未知] → 不宣稱長表更快或負載承諾；只驗現有有界佇列與 deterministic closure，獨立容量規劃不屬本案。
