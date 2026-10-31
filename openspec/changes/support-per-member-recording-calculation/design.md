## Context

基準6f4a5fc。已核對WriteGroupMember無計算欄位、WriteGroupRecordingProfile僅group reduction，runtime/group_boundary_wide_profile及long依group Kind決定cumulative；snapshot/assembler_counter要求全group counter或無counter。mean profile雖可gauge平均＋counter_raw/state最新值，但counter_raw不做delta。group_row_layout_wide_counter_test確認舊wide欄位存raw、delta在provenance，不能改下拉就說能mean/mean/delta同列。

## Goals / Non-Goals

**Goals:** 同record group的各member選latest/mean/delta，各自使用對齊member時間窗，wide同一列各欄是真正選定結果；不確定delta有NULL原因／required skip與原子baseline進展。提供editor consumer契約。

**Non-Goals:** 不把計算寫入global Point、不新增採集／任意公式／觸發引擎、不新增DB dialect或protocol、不遷移歷史SQL值、不自動拆成多record、不放寬目的表／欄位ownership與legacy canonical create-once，不直接改四案UI。新mixed shape涵蓋既有managed wide及long布局；long一列一個到期member，不假裝一列包含全部members。native-monthly/legacy custom保原功能，未具新shape證明的組合明確受阻。

## Decisions

### Member意圖與新mixed形狀獨立於來源

新增WriteGroupMember.calculation=latest|mean|delta，僅在顯式versioned mixed profile使用；API讀回configured method、effective method、source/result type、applicability、member identity及group/source/mapping/measurement revisions。canonical group member/mapping是authority，Point／Tag只提供來源與型別，沒有global method。新member初始化latest；群組預設只是建立新member的UI捷徑，既有explicit不被批次默覆。相同source可在不同record使用不同method與獨立baseline。

新增profile.kind=mixed與result shape member-calculation-v1，不借舊instantaneous/cumulative假裝支援；member method必需，unknown/omitted混合member拒絕。既有profile不因reload/upgrade增欄或改method。既有profile/sample_reduction與role投影只供legacy摘要：mean gauge平均、counter_raw/state最新，cumulative維持原raw column＋delta provenance。明確改用mixed先建立獨立mixed record的candidate/revision與fresh安全證據，不in-place重寫legacy身份或已寫入資料；同canonical identity內的相容變更才沿既有successor/schema流程。

### 每member對齊時間窗與結果

每member保存自己的interval/unit與alignment basis，明確支援15秒mean、60秒delta及10分鐘等需求；原「mixed單一1–300秒上限／同窗」已被新需求取代，不當新契約上限。沿既有snapshot selector的UTC epoch anchored half-open windows，啟用時刻不能改phase；15秒對齊00/15/30/45、5秒00/05…55、600秒於使用者UTC+8顯示分鐘00/10…50。member interval為正整數second/minute，正規化為精確duration並檢查overflow；必須支援5s/15s/60s/600s驗收值。其餘可選值以backend明示且經資源限制驗證的bounded capability為準，不臆定新的通用最大值、不用UI任意unit繞過validator；既有legacy bounds不變。此處不新增cron／日曆／DST引擎，顯示alignment basis與實際next-boundary，不能宣稱任何timezone都自動有相同local minute phase。latest用既有eligible latest，mean用既有等權arithmetic mean且numeric gauge才可選；bool/text只latest。delta須已確認counter binding/policy，不能按name/unit猜counter。新mixed不以sample_reduction決定所有member；若同時送衝突舊group reduction則拒絕，不忽略。

派送重用snapshot selectors、measurement.Mean、aggregation.CloseCounterWindow，只有delta members有counterbinding；mean/latest不能偷偷產生baseline。scope/CAS、型別、scale、source-boundproof均保留。完整wide result欄存每member結果：Amean、Bmean、Cdelta。delta raw、from/to、sample IDs、counter status/policy與coverage存既有可擴充provenance；不能以raw覆蓋delta欄。exact integer/decimal結果沿既有typed SQL encoding，不經float64或JSONnumber損精度。

### 缺值與baseline進展明確

每member required沿既有配置。missing/invalid/stale、無有效mean、delta首窗無baseline／不確定reset／無法解釋下降皆沒有可用結果。optional寫NULL＋具體原因，僅fresh schema證明nullable與quality/provenance支援才允許partial；任何到期required沒有結果則wide整個row skipped（long只跳過該member列）；未到期required不造成skip，不製造0、不退回raw、不拆成兩筆mean row。當首窗因required delta跳過row，baseline與skip closure仍原子保存，下一窗才能算，不能永遠首窗。

counter的signed、rollover、reset證據沿既有確認policy與可靠tail規則；不發明所有負差invalid。missing window不advancebaseline，後續可能計算跨窗span total；新mixed同窗delta欄只接受屬於該member當窗的可靠增量，CounterGap/Span結果寫NULL或required skip並在provenance保留真實from/to/span與reason，不把跨窗total冒稱當窗用量。此是新resultshape明確政策，不改legacy span記錄語意。

### 版本與durable closure不能分離

member method、resultshape/type、window、required、counterbinding/epoch/policy納入managed digest、candidate、source/mapping associations、readiness、frozen runtime descriptor、schema proof與CAS。method改變使舊readiness/consent失效，舊snapshot/sample proof不能被新method重解。counterbaseline以workspace/group revision/member/measurement binding隔離；新method或binding revision透過既有受審migration/handoff，無可證明相容baseline則新baseline，不跨record借用。

單一closure transaction保存所有member計算證據、counter checkpoints、bucket進度與row/skip outbox；提交失敗不得只advancebaseline。重啟replay沿immutable runtimeversion與stableeffect id，不重算已接受舊output成新method；SQL response lost使用既有receipt查證，unknown不重建另一row。copy新record不複製baseline/readiness/proof。

### 已確認不同週期輸出與邊界

採UTC epoch對齊，member僅當邊界T是其interval整數倍才到期，計算半開窗[T-interval,T)；恰在T的觀測屬下一窗，delta baseline沿既有可靠前窗tail而非強行用T的值。於T加既有allowed lateness後才可closure，實際SQL commit可更晚。wide以各member到期邊界聯集產一列；未到期欄NULL＋not_due，不carry舊delta、不補0、不advance其baseline。15秒mean＋60秒delta在:15/:30/:45僅mean有值，:00使用各自前15秒／前60秒結果同列。long每個到期member各一列，未到期不新增列；到期失敗沿required skip／optional NULL及具體reason，不混成not_due。fresh schema必須證明wide非同頻欄位可nullable並保留not_due品質，即使該member在到期時為required。

每個到期結果保留member identity、method、interval、aligned window start/end、actual coverage、quality/reason及provenance；not_due不偽造已計算window或sample，保留其設定與next due。wide列的T是共同到期標記，不代表全部欄同長window。first partial沿snapshot.windowBounds及markTruncation：aligned identity不改，實際起點max(window start, StartsAt)，可用latest/mean保留boundary-truncated與partial coverage，不捏造啟用前樣本；無baseline delta仍不可用，required skip／optional NULL，baseline與skip原子保存。不把partial mean冒稱完整15秒，無可靠單窗delta不可冒稱60秒用量。

### 儲存擴展與草稿能力分期

**Supersedes**: support-per-member-recording-calculation / 分期擴展與consumer交接（僅澄清1.2與1.3交付界線，後續runtime安全gate不變）

1.2只擴展nullable calculation／result_shape及內部repository mixed snapshot寫入／scoped readback；legacy CRUD/CAS與foreign identity保護維持。正常mixed create/update在此階段仍拒絕，沒有member cadence前不得用0、任意預設或假group interval通過normalizeWriteGroupProfile。1.3完成member cadence、適用性與record-scoped identity後才接通正常mixed草稿CRUD；草稿可儲存不是可Apply／啟動，runtime、encoder、schema/readiness/proof與capability仍依後續任務逐項驗收。拒絕用新平行draft store或提前打開runtime來繞過依賴。

已核對029的write_groups.profile_kind CHECK只允許legacy kinds，新增兩欄不能使mixed snapshot roundtrip；1.2安全擴展該CHECK保留舊值與約束。採專用sql.Conn保存原foreign_keys設定，交易前只在該連線暫停FK，交易中取寫鎖、確認已知schema、建立新暫存表、以明列完整欄位copy、驗完整內容與row count、替換舊表並還原indexes／triggers及相關views，foreign_key_check全庫無違反才commit；失敗rollback並恢復原FK。不得先rename舊表造成child references漂移，不使用writable_schema、不修改歷史migration、不轉換legacy intent。未知schema／未識別constraint fail closed，不靠任意SQL字串替換猜結構。FK恢復失敗必須報錯並丟棄該連線，不將FK-off連線歸還pool。

1.2不為未開放的mixed canonical create提前重建profile_keys。1.3接record identity時再同步其CHECK與047 lineage key validator，保留既有PK／UNIQUE／FK、lineage genesis／head／reservation與create-once身份。047已初始化但缺損的corruption不得由repair猜測或重建；nullable缺欄repair只補結構、不回填method。測migration前後legacy完整列／JSON、members、versions、migration maps、lineage及revision不變，覆FK原ON／OFF、重跑／重啟、併發startup和copy／replace／index／FK-check故障rollback；不以row count相等代替資料完整性。

### 來源共享與記錄身份隔離

同source可在同一workspace、同DB/schema的獨立mixed records使用不同method／cadence，不得把需求縮為只允許跨DB。共享的是來源引用，不是目的表／欄位寫入權、canonical key或lineage。legacy circuit/profile create-once、migration source claim與既有destination ownership契約保持原樣；不可直接刪除checkProfileSourceOwnership、放寬所有profile guard或靠不同名稱繞過。

1.3建立明確且持久的mixed record-scoped canonical identity：同record重送沿既有idempotency／CAS找到同一lineage，另一record有獨立identity／member/baseline；不得再以source device或circuit/profile組合作為所有mixed records唯一身份。優先在既有canonical key／selection_scope與lineage機制內擴展可區分的record scope，與legacy circuit／unassigned角色保持不相交；不要新增第二套權威。只在已驗證mixed record scope上允許來源重用，仍驗workspace／source revisions與managed target ownership。具體key shape與CHECK須一起覆蓋persist／reload／restart，不能只在記憶體helper隔離。

正例為同DB/schema、不同獨立managed tables的record A latest/5s與record B confirmed delta/600s引用同source，互改不影響對方；反例為第二record企圖奪取第一record的同target table／columns，仍拒絕且不改owner。shared long table僅沿既有明確storage scope／member row identity授權，不因來源可共享而擴大任意共寫權。上述草稿配置正例不代替後續runtime baseline／SQL隔離驗收。

### 分期擴展與consumer交接

expand新增optional persistence字段／版本讀取，legacy不變；migrate明確mixed候選走fresh安全gate；contract只在runtime/encoder/readiness/proof全部支援時開capability。不刪舊profile。新shape需要schema preview確認result欄型別與nullable，已有raw wide table改用delta結果語意不能默認相容；legacy到mixed須另建獨立record；wide raw轉result須新owned target，long依下述既有shared-storage邊界，歷史留存。

**Supersedes**: support-per-member-recording-calculation / Member意圖與新mixed形狀獨立於來源（僅取代「必要時走既有successor/schema流程」對legacy轉mixed的歧義，不減少fresh proof與歷史保護驗收）。已核對circuitProfileKeyOf／mixedRecordProfileKey：legacy circuit或unassigned＋instantaneous/cumulative身份與mixed record＋mixed＋RecordID不相交；recreatedProfileIdentityInTx要求candidate key等於genesis。因此legacy原地Update及同lineage successor改mixed均須typed拒絕，不鬆canonical create-once、ownership或新增cross-canonical遷移引擎／scope alias。使用既有獨立mixed建立流程，以fresh RecordID／lineage、明確member method／cadence重新取得schema/result-type/nullable/quality/provenance及readiness證明；wide raw轉result必須新owned target，long可沿既有shared-long ownership／row partition重用同表，但必須fresh mixed result-shape／nullable／provenance／schema proof全部成立且不重解舊rows。不符合則走既有新storage scope／新目標流程或明確blocked，不任意換tablename、新增alias或放寬ownership guard；舊proof不可沿用、baseline不可複製。舊record、歷史raw、legacy create-once原樣，不自動停止或刪除舊record。這是操作員明確另建／複製definition後改用新記錄，不是同record原地轉換成功，也不是把一個mixed需求拆成多record；新mixed仍須同列mean/mean/delta。editor／lifecycle consumer交接須明示另建與獨立啟停，不能把typed拒絕包裝為已轉換；本案不新增自動handover或copy endpoint，4.2仍驗真production鏈路。

editor每已選member顯示「存最後讀值／平均／增加量」與例子，不讓操作員以為與ABCD字序是同一設定；同source不同record不互改。此案不實作頁面，streamline-recording-task-editor依本案完成後接API；adopt-simple-recording-entry加同row mean/mean/delta驗收。正式apply授權以docs/verification/simple-recording/implementation-decision-2026-10-30.md為準；demo/helper不代替runtime或SQL驗收。

## Implementation Contract

- 本案獨立於preserve-per-point-word-order的decoder能力；兩者都會影響source/runtimeproof時不得平行修改重疊檔，整合重驗revision與samplebinding。
- 與integrate-protocol-aware-acquisition-batching僅共享觀測輸入契約：batch大小／協議分段不得更改observed_at、sample identity、quality或source revision，不以read batch重定記錄window；不新增硬性先後依賴，重疊runtime/proof檔需序列整合重驗。
- 實作選測workspace/snapshot/aggregation/runtime/dbtarget/groupdelivery與scoped API，每個垂直批次可build；完整功能最後跑相關側完整tests/vet/lint與前端總驗收。
- 同機可丟棄SQLite為必要local durable＋SQLreadback證據；既有其他SQL dialect僅在受影響capability範圍測。fake transport，不連PLC／正式DB、不新增環境。
- 真驗收fixture：同bucket A=[20,22,24]→22、B=[10,14]→12、C可靠baseline1000→1008→8，SQL一列三欄；另驗firstbaseline、requiredskip/optionalNULL、gapspan、policy/reset/rollover、CAS、crash/replay、lostreceipt與legacy無變。

## Risks / Trade-offs

- 混合欄位被舊raw encoder誤寫 → versioned resultshape與真SQLreadback斷言。
- 首窗skip卡死baseline → skipclosure同tx持久化checkpoint並restart測試。
- group defaults覆蓋明確member → 只新member初始化，immutable intent與CAS。
- 相同來源不同record共享baseline → group/member/binding/revision隔離＋雙recordfixture。
- 以15個task硬塞大改 → 每task是單一可驗證垂直契約；實作若超出一context或發現新backend外延，停下修規格，不省略gate。

## 最新需求取代共同窗口，儲存語意已確認

**Supersedes**: support-per-member-recording-calculation / 單一時間窗與每member結果

2026-10-30使用者要求同record內15秒平均＋1分鐘delta，以及多record的5秒／10分鐘整齊對齊。共同window限制與mixed 300秒上限撤回；原pending task1.3/2.1改為各member cadence，未標完成。interval合法值由上述backend bounded capability與驗收值約束，不沿用任意demo範圍。

**Supersedes**: 前版必要storage決策blocker由2026-10-30 16:50 UTC使用者「這個決策依你說的避免記錯意思下去做」解除，確認寬表未到期留空、長表不新增未到期列、不補0、不沿用增加量，到期失敗另記原因且每筆保留自身區間。該時點只確認計算／儲存語意；其後正式apply授權已記錄於docs/verification/simple-recording/implementation-decision-2026-10-30.md，未追認任何未完成測試。

**已確定時間事實:** selector.go按UTC UnixNano/interval定位半開窗；Plan等window end+allowedlateness後閉窗；現pipeline/worker tick與SQL重試造成實際commit時間延遲。畫面分「資料窗口／到期時間」與「實際SQL寫入時間」，不能保證00整秒落盤；delay/restart不得把補寫時間當資料時間。採集ticker與記錄對齊是不同層，不改採集protocol週期。
