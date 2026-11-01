## Context

基準6f4a5fc：sourcerule.Service.Create 將 plannedAddresses 展開成 Points；collector.pollDevicePoints 在設備鎖內逐點 pollPointWithAcquisition，conn.Read 的 Count 是單一 typed scalar 需要量。MC adapter 只用 BatchReadWord/Bit，FATEK 用 ReadRegisters/Status；低階 RandomRead／ReadRandom 可用但沒接線。optimizer.BlockMerger 與 Dispatcher 存在，production 無呼叫；前者只按 area 分組且預設 MaxHoleSize=10，後者單一 ByteOrder，不能原樣接入。

collection-optimizer 舊 spec 的 max_hole_size=10 不代表可安全跨不可讀區；本案修改該 requirement。已掃 active、parked、archived，逐點字序由 parked preserve-per-point-word-order 唯一負責；recording editor/member calculation 只消費，不重作 acquisition。

## Goals / Non-Goals

**Goals:** 在既有 collector 內減少冗餘協議 round trips，完整保留 Point 值、格式、品質、身份與時間；安全公平調度、同源預覽與可重現量測。

**Non-Goals:** 不改記錄窗口／SQL commit 時間、不根據 recording selector 縮減或增加來源集合、不承諾 PLC 原子快照、不新增 driver／第二採集引擎、不碰真 PLC／正式 DB、不做任意位址探測或自動宣告設備可讀範圍。

## Decisions

### 唯一不可變計畫與能力界線

擴展既有 optimizer 為純函式 planner，輸入 immutable device/station/connection revision、protocol capability、effective Point snapshots、polling group/cadence 與 readable coverage revision；輸出 deterministic sorted requests、每點 wire slices/decoder policy、plan fingerprint、限制／fallback 原因。無法解析地址或不支持型別是明確 point error，不准像舊 Merger 靜默略過。

僅同設備、station/unit、連線版本、protocol/frame、function/access area、到期 polling group 合批；不同 cadence 不因這次同時到期就合併排程權威。跨 area random 只在能力明確支持時允許。capability 限制取協議、client、裝置配置三者最嚴，未知 random 能力預設 scalar/安全 contiguous；不得發測試讀來猜能力。MC client 目前 word random 1–255 items、FATEK 64 items 是程式上限，不等於所有 PLC 支持；只配置已支持命令才用，限制需對 simulator fixture 驗證。每 scalar words/native components 必須完整放同 request；超限 scalar 拒絕，不拼跨時刻半個值。MC hexadecimal address、FATEK D/DR width、Modbus normalized offsets 全用各協議既有 parser，不用 decimal regex 猜。

**Byte budget 澄清（屬上述唯一不可變計畫與能力界線）**

2026-10-30 澄清上述「byte限制」：immutable capability snapshot 分別承載 Protocol／Client／Device 的 `{RequestBytes, ResponseBytes}` budgets；request與response各自採正值中最嚴限制，0只代表未指定，不能代表已驗證成本為0。計算涵蓋受支持的完整application frame（Modbus MBAP或RTU unit＋CRC、MC3E binary header＋end code、FATEK ASCII framing＋checksum），不含IP／TCP／Ethernet開銷，不是通用payload byte數。成本必須依實際command/frame serializer與response fixture求得；未知成本明列unknown，遇正值budget則以 `unsupported-frame-cost` 拒絕該候選，不宣稱符合限制。random成本在其策略交付時另以真實serializer驗證，不能套用contiguous公式。

來源查核：`internal/protocol/modbus/frame.go` 的 BuildTCPFrame／BuildRTUFrame／BuildReadRequest／ParseReadResponse，`internal/protocol/mcprotocol/frame.go` 的 RequestFrame.BuildPacket與 `internal/protocol/mcprotocol/client.go` 的BatchReadWord／Bit，以及 `internal/protocol/fatek/frame.go` 的BuildFrameToBuffer／ParseResponse與 `internal/protocol/fatek/client.go` 的ReadRegisters／Status。測試須比對實際序列化長度與正常回應fixture，包括bit packing、native component width、FormatAddress長度、恰好符合與超過一byte的邊界；以下為目前支援路徑的code-derived測試條件，不是通用協議最大值或實測效能承諾：

| Contiguous／scalar frame | Request bytes | Response bytes（正常讀取） |
|---|---|---|
| Modbus TCP | 12 | 9＋payload；word為2N、bit為ceil(N/8) |
| Modbus RTU | 8 | 5＋同上payload |
| MC3E binary | 21 | 11＋payload；word為2N、bit為ceil(N/2) |
| FATEK ASCII | 8＋2＋len(FormatAddress) | word為9＋N×(native bits/4)，bit為9＋N |

item limit也須clamp到實際client：目前Modbus為125 words／2000 bits，MC BatchRead為125 words／2000 bits（不是MC協議通用上限），FATEK為64個16-bit／32個32-bit components或255 bits；裝置設定只能再縮小，不能放寬client實際能力。限制與成本模型參與不可變plan/capability identity；serializer或能力版本更新後須重驗fixture，不沿用未證明的bytes估算。

### 安全覆蓋與去重解碼

MC/FATEK 針對稀疏 scalar 展開完整 native units，用現有 random client 一次讀 D0/D10/D20（已宣告可支持、int16 的設計 fixture）；只支持 batch contiguous 的組合採合法區塊，否則 scalar。Modbus 合併僅當 union 完全由有效 scalar spans 覆蓋，或 gap 每個 cell 都在同版本、明確 operator/config declared readable coverage 內且非 forbidden。SourceRule start/count 或曾經讀成功不構成 gap 授權；skipped/disabled 未被其他 active point 覆蓋的 cells 預設 excluded。max_hole_size 只是已許可範圍內成本限制，不能授權。無 coverage 的0/10/20三點仍三次；已確認0–20整段可讀且限制允許可一次讀21words，僅發射原三點。

同輪同wire address/span重疊可去重 raw read，但 fanout 必須保留每個 canonical Point 的各自 typed decoding、quality、format 與sample identity；不得按名稱或record合併。相同 raw words 的 ABCD/CDAB 點值不同是正確結果。每 Point 使用 preserve-per-point-word-order 的 immutable resolver，不使用舊 Dispatcher 的全域 ByteOrder。不能改共享 connection converter。去重僅同一 plan/scan，不快取跨 scan 舊值冒充新取樣。

### 調度、重試與背壓有界

3.2 交接澄清（2026-10-31）：以下只精確定義原「調度、重試與背壓有界」的容量、coalescing 與 loss 歸屬，不改已完成 3.1、task IDs、來源權威或建立另一排程引擎。

保留 device lock 與各設備並行，在既有 scheduler 加 bounded fair ready queue。每 device 最多一個 running scan；每已配置 polling group 最多一個 coalesced pending container，內含每 device 最多一份 latest due snapshot。Global pending capacity 計未 dispatch 的 `(device, group)` work entries；同 tuple replacement 不增加容量，新 tuple 超量則拒絕並計數，不踢除／覆寫另一 group 已 admitted 的 Point set。Group 部分 device 已 dispatch 後，新 tick 只更新同 tuple 的 pending 或在容量允許時新增；running scan 的 captured set/deadline 不被改寫。不同 ready devices round-robin，同 device 的不同 groups 用穩定先到／輪替順序避免飢餓。『保留多 group latest due』以成功 admitted entries 為界，飽和例外必須可觀察。

每份 immutable snapshot 保留自己的 captured time、Point set、ticker/config/install identity 與固定 scan deadline；coalescing 是新 due snapshot 替代尚未執行的舊 due，需計 coalesced/overrun，不可抹掉missed history或延長 running deadline。Global active concurrency 計取得既有 device gate 後的連線／讀取 active phase；等待 gate 或 retry backoff 不佔 active permit。Backoff 暫放 permit但保留 device scan/gate 所有權，恢復時按原 context/deadline 公平重新取得，不能越過已等候其他設備。每 configured device 最多一個等待／執行 scan goroutine，另加固定 dispatcher／既有 group producer；不得按 tick 或 pending entry 生無界 goroutine。這是依配置設備數的明確上界，不宣稱與設備數無關的固定 goroutine 上限。

沿既有 scheduler retry wait 與 wirebudget.Wait 的可選 context cooperation 接點暫放／重取 active capacity；不建立第二retry owner，不改budget、delay、decoder或actual-send admission，不把阻塞 acquire 放入短且非阻塞的 admission/send observer。取得permit後仍必須通過原 captured-install/config/stop send fence；cancel/stop在gate/backoff/permit wait/read均可中斷，transport解除等待不得關閉新版連線。既有ConnectionManager per-key single-flight/短global鎖/late generation guard已實作，僅保留回歸，不重造。此階段限定scheduled容量管理；manual同device原gate保留，manual共用planner/executor與容量契約在3.3完成，不提前宣稱。

Global active concurrency與pending capacity用既有SchedulerConfig入口提供正整數設定；未指定者維持constructor/default相容，明確無效值須可辨識地拒絕，不以新增產品設定頁補此工作。原配置下的defaults是設計值，不宣稱已量測最佳效能。Due coalesced、capacity-rejected/skipped、overrun與pending depth分開診斷，拒絕不製造讀值或刷新last-good。測試用既有可控seam或Go synctest，不把目前只用於completion的Clock冒充ticker clock，不建通用clock框架。

結果buffer loss也唯一歸3.2，與due admission loss分開計數；4.1只消費／呈現，不再建第二套counter。只處理既有emitValue實際丟掉的CollectedValue，使用該值已捕捉的InstallationScope/PointRevision/PointID，不在drop後lookup新revision，不改SampleEnvelope/StableSampleID，不把wire Good偽稱下游已接受，也不將drop當通訊失敗。保留累積總loss及有界近期scope/revision/point明細；近期最多256組（設計上限，不是量測值），超量淘汰舊明細並明示已淘汰的loss彙總，不能假稱完整歷史或靜默丟總數。未具完整identity的legacy值明示unattributed，不借當前identity。只新增native唯讀診斷供驗證與4.1，不在3.2新增public路由／診斷平台或泛化其他channel。

每 request 用既有 DefaultRetryCount/Delay 上限與read deadline；統一嘗試預算，adapter自癒重送也計入，禁止batch×scalar×adapter巢狀放大。僅明確 unsupported-command 拒絕允許同scan一次降級至已驗證 scalar/contiguous plan，且計入總預算；timeout/CRC/短回應不猜unsupported，不盲目fanout重讀。下次scan對同connection capability revision保留降級狀態，重連／capability修改使其失效。總預算定義為該scan initial plan request數加fallback plan request數（僅真的被授權觸發時）乘以(1+retry count)；每request仍個別限額且總deadline不延長。

### 結果、時間與版本不冒充一致快照

wire response完整驗證後依point span解碼；不可判定對齊的短／壞checksum整包受影響點 Bad，其他已成功request不受牽連。若協議可明確定位個別item錯誤，只該item受影響；禁止client缺項map零值成Good。每point每scan最多一個final emission，重試中間結果不重覆入measurement；整個scalar解碼失敗只該點Bad，絕不使用舊值補Good。

同device/poll scan保留 acquisitionID，新增request/batch identity供追查，fanout各point仍各自sampleidentity。無trusted device time時同成功request各point的ObservedAt/ReceivedAt採同一gateway read completion，跨request各自completion；重試用最終成功completion。不得把tick起點、預計時間或SQL時間當採集時間。同包也不宣稱PLC scan atomic。

計畫鍵納入 source membership、address、width/type、function、effective format/source revisions、station/connection/capability、cadence、coverage；更新原子安裝新snapshot，使preview/ready token與affected source proofs失效。舊inflight response不改綁新point metadata：依既有ingestor舊revision eligibility接收為舊歷史或拒絕並計數，不能刷新新revision freshness。停用／移除立即阻止其後新request；已發出不能撤回的回應按舊身份處理。

### 同planner預覽與量測證明

新增scoped read-only acquisition-plan preview，輸入workspace/device scope與expected revision（不接受任意host或地址）；回傳plan_revision、capability_revision、acquisition_point_count、unique_wire_units、planned_requests、request kinds/spans、safe coverage/fallback reasons、poll interval、estimate=true。零PLC I/O、不變更設定；stale revision回409，不把未知能力顯示成random=1。UI各group label為「每輪計畫協議讀取次數（不含故障重試）」；actual attempted/succeeded/failed/retry/fallback/canceled、cycle latency、overrun和last-good age另列，同scope/revision標示，不稱TCP封包數或SQL次數。

量測使用同機fake transports／simulators、相同seed/data/point set/cadence/limits、baseline逐點與新planner兩模式；baseline只測試hook，不建立產品第二引擎。必要fixtures：稀疏三點MC/FATEK 3→1；Modbus未授權gap 3→3，明確coverage 3→1；21連續int16 21→1（限額允許）；125/126 limit切分、多word邊界、重疊去重不同格式、invalid address、forbidden gap、unsupported command、timeout、短回應、CRC、cancel、rapid revision swap。這些是設計斷言，不是已量測性能。

負載矩陣固定1/10/100設備×每設備3/128點，1/20/100ms模擬RTT，至少100完整scan；另一台timeout與全部timeout、queue飽和。報告request count、每device p50/p95 cycle latency、last-good age、missed/overrun、峰值queue/goroutine與配置，不用總平均掩蓋starvation。無錯fixture值/quality逐點等價；除有意同包completion外時間單調與來源正確。合批fixture無fault時requests嚴格符合planner；同條件controlled RTT的cycle duration不劣於baseline，否則不啟用該策略並查明成本。公平性用fake clock斷言有available capacity時持續ready設備每輪ready-queue rotation獲服務，全部slot阻塞時由deadline釋放，慢設備不佔住其他設備slot；queue永不超配置，stop後無新send或late-good污染。CPU／memory實測列數值，不捏造速度百分比。

補充獨立source review後的必要收斂（屬上述決策，不新增引擎）：

- optimizer.union end使用max(existingEnd,newEnd)，修正舊merger在D0 int64＋D1 int16會縮短block的問題；invalid parser不得靜默略過。不同record訂閱同Point原已共享，不重報為新節省；真正去重是不同Point重疊wire units。
- 策略先union exact overlap，再對完整合法group比較all-contiguous、random packing與有界hybrid候選，以request count、wire bytes依序最小、地址與ID穩定tie-break。不得先固化dense subruns使D0/D1/D10變兩次；dense21同一次時選更少bytes，不做指數搜尋或自動調參。
- pollGroup不再wg.Wait全設備才接受下輪due；producer只送有界intent，每group保留自己的due generation／points，不能覆寫另一group。ConnectionManager.GetOrCreate目前持全域mu進行dial且按device ID重用，改每key single-flight、全域鎖僅publish；protocol/config/generation必須相符，late dial不覆蓋新版。shared physical serial endpoint別名無可證明序列化則拒絕，不新建serial subsystem。
- MC/FATEK random須完整驗長度、status、item encoding；FATEK空items、非0/1離散值、trailing garbage、missing map key均失敗。Modbus bit batch取得完整bit vector，不可拿現adapter first bool當fanout。
- PollNowContext走同planner/executor且僅exact requested subset，不建scheduled duplicate；保存scheduler_observation.go的原installation lease，禁止回應時才lookup新binding。emitValue現有drop-oldest改為可歸因scope/revision/point loss counters，drop不刷新downstream last-good，不能宣稱無損。
- 量測warm/cold分開；planner-only另測100/1000/10000 Point CPU／allocations。新增慢device dial、nested overlap、256 MC units、65 FATEK units、serial aliases、完整Modbus bit scatter、manual cancel fixtures。所有數量為設計測試條件，非現有測量。

FATEK 原廠契約複核（2026-10-31，task 2.5 補償先前純 planner 假設）：Cmd48 的 N 是 component 數，但協議另限 64 words；bit／16-bit 各計 1W、32-bit 計 2W，configured component cap 與 bytes 上限仍獨立取最嚴。20 個 32-bit 加 24 個 bit／16-bit 為 N=44、64W；33 個不同 DR wire 地址的 int32 不可因 component 數小於 64 就單包送出。不得把 N 改成 weighted word 數。既有 NativeUnit 保持 component 語意，唯一地址步進契約由 scalar expansion、contiguous 相容判斷、random 去重、frame cost、scatter 共用：DR/DD 為 2 個底層 word 地址、W* 為 16 個 bit 地址、DW* 為 32 個 bit 地址，普通 word／離散點為 1。目前支援的 WX/WY/WM/WS、DWX/DWY/DWM/DWS 起始 bit 地址須為 8 倍數而非整個 component stride 的倍數，完整 scalar footprint 不得越界。合法的 DR1 不得當成 DR0 連續命令的第二回應（其實是 DR2），也不能為避開錯誤全面拒絕 odd address。修正計畫 model identity，舊 preview 失效；MC/Modbus 既有座標不改。原 completed task 描述與證據保留，task 2.5 的 native width／正式 adapter 驗收承擔此次修正，未經正式 wire fixture 與獨立複驗不交 collector。來源為 [FATEK 原廠 Appendix 1 文件](https://www.esea.cz/support/fatek/FBs_Manual/Manual_2/Appendix1.pdf) 的地址表、Cmd48 限額與 Cmd47/48 完整範例，以及 [FBe Advanced Chapter 14](https://www.esea.cz/support/fatek/FBe_Manual/Advanced/Chapter_14.pdf) 印刷 14-17 的 DR0/DR2、DD0/DD2 組成圖；fixture 必須解析實際地址及數值，不以舊 serializer 自我比對代替原廠規格。

能力宣告只沿用既有 `Device.ConnectionConfig` 的 optional `acquisition` 物件；不新增 table、record selector authority 或第二份設定來源。格式固定 `version:1`、`profile`，可附 `limits:{native_units,request_bytes,response_bytes}` 與 `coverage:{revision,readable,forbidden,excluded}`。缺省宣告維持 unknown 並僅採安全 contiguous/scalar；`safe-v1` 明確禁止 random。非缺省的 null、未知欄位/profile/version、非正整數 limits 與不合法 coverage 均在既有 device validation/persist/install 邊界拒絕。coverage 的區域與半開實體地址範圍重用 optimizer 既有 canonical parser/bounds，content 本身納入 fingerprint，opaque revision 不足以授權讀取 gap。`native_units` 是每 request 的 native component 上限，單向 request/response bytes 分別取更嚴正值；缺省不是 0-byte 已證成本。

初始經官方文件核證的可選 family profiles 僅 `fatek-fbs-ascii-v1` 與 `mitsubishi-fx5-cpu-3e-binary-v1`；後者限 FX5 CPU 內建 Ethernet 3E binary，不能泛稱所有 MC/Q/L/R 或 Ethernet modules。client 函式存在、protocol subtype 或成功連線都不是裝置 random 支持證明；即使 FATEK subtype 明確，舊 config 沒宣告仍 unknown。單一純 resolver 同供 install 與後續 preview，registry schema 從相同宣告定義輸出。FX5 CPU 0403/0000 binary 的 192 access-point protocol cap 只以 `RandomProtocolLimit` 額外限制 random，不套給 contiguous/bit，也不能以 client255 冒稱裝置上限。既有 `ProtocolLimit` 若配置仍為所有相關候選共同更嚴上限，新欄不能放寬它。FX5 CPU 0401 word960/bit3584 與現有 client word125/bit2000 各自取更嚴；FATEK 48 的 component count 與 64-word weighted budget 獨立保留。

每次 scheduled due 在一個既有 scheduler 短鎖內捕捉 device config、source Points、原 group/cadence、immutable ReadProfile、PointRevision 與原 InstallationLease/Scope。pure plan fingerprint 只含 revision-bound 設定、source membership、capability/coverage content 與原 decoder policy，不含捕捉時間、AcquisitionID、requestID、observer 開關或 ephemeral connection token；lease/scope 仍保留於原 binding/admission，不能充當新 SourceEpoch。record selectors 不參與 source-set 推導。既有 ingestion 先凍結 raw→mapped binding、revision CAS 與 publication guards 均不更改。

scan 固定 deadline 在同一 snapshot 建立時決定：取既有 request timeout ×（initial request 數＋預計一次安全 fallback 最多 request 數）的 checked 上界，尊重較早 caller deadline。預留時間不是 fallback 權限；只有完整且明確 unsupported-command response 才能同 scan 一次安全 replan，且僅重讀 unresolved Points。每 request 仍是原 timeout 與 `1+DefaultRetryCount` actual-send budget，adapter reconnect 與 scheduler retry 共用之；backoff、排隊與所有 fallback 不延長 scan deadline。MC C059 及嚴格 FATEK command48/status9 可授降級；FATEK4、timeout、CRC/短/錯站/錯命令/尾隨payload 不可。

actual send 在既有 `wirebudget.Take` 最前檢查 context 和捕捉的 install/config/stop；成功檢查為 admission 線性化點，後發 disable 無法撤回已 admitted send，但每次後續 retry/fallback 都須重檢。callback 不得做 SQL/provider/scheduler callback 或逆轉 mapping→pointMeta→scheduler 鎖序。短 callback 後仍使用原 budget CAS，再以無鎖 token getter stamp 實際 final wire attempt。

連線 token 僅 manager 成功 publication generation 加 adapter 成功 connect/reconnect publication generation，failed/late completion 不可清掉較新 downgrade。同 token/capability revision 的已證 unsupported 在下次 scan 保留；新成功 token 或 capability revision 不沿用。無有效 token 可完成同 scan 一次安全 fallback，但不可永久 cache downgrade。requestID 與 plan fingerprint 僅附 final `CollectedValue`/acquisition observation；不擴充 durable `SampleEnvelope` 或改 `StableSampleID`。同 request 用真 final completion，跨 request 不冒稱同一 PLC scan snapshot，每 Point 每 scan 最多 emit 一次。

Primary profile evidence: Mitsubishi FX5 Ethernet manual JY997D56001K printed p57 / PDF p58 (CPU/module tables distinct), https://www.mitsubishielectric.com/dl/fa/document/manual/plcf/jy997d56001/jy997d56001k.pdf#page=58 ; FATEK FBs manual Appendix1 https://www.esea.cz/support/fatek/FBs_Manual/Manual_2/Appendix1.pdf .


## Implementation Contract

In scope：optimizer、collector、既有adapters與production wiring、scoped preview/diagnostics及其consumer契約；Out of scope：所有Non-Goals與其他change的UI布局。預覽屬本案backend authority；editor owner只接此結果，不自己算封包。Implementation先完成preserve-per-point-word-order；planner純函式可先TDD，但接線需等待其decoder/fingerprint契約。member timing並不依batch才正確：本案為其「優化／請求數」展示前置，不阻塞獨立窗口語意決策。

每垂直task先失敗fixture→最小實作→重構，每批可build。集成須穿cmd/test_ui scheduler→adapter→measurement path，不以helper存在當完成；相關Go tests/vet/race、API/frontend consumer測試與最終相關側完整gate。純提案僅analyze/validate/diff/line gates，不執行PLC或性能宣稱。

## Risks / Trade-offs

- 批次擴大失敗blast radius：僅合法範圍、嚴格驗包、不盲目fallback，保留每點品質。
- 過多合併增加多餘bytes：同planner列extra units，先安全再成本；不以最大合併為唯一目標。
- 共享格式／過期回應污染：immutable decoder與revision隔離、race與swap fixtures。
- 大矩陣耗時：4.2依「負載矩陣分批驗收」保留完整條件與樣本，按具名case／有證據的續跑區段執行，不靠漏測兌現完成；本次僅調整驗收執行分批，不新增產品子change。

### 手動讀取保留回傳與生命週期

3.3 沿原 PollNowContext 的明確 caller 意圖：可在尚未 Start 或 scheduled Stop 後讀已安裝且非 recoveryPending 的 requested Points；不自動擴到全部來源或建立 scheduled due。Unknown／recoveryPending／空集合沿原回傳邊界，重複 canonical Point ID 最多一份 final 結果，不按地址合併不同 Point。每 manual/device 保留一次 AcquisitionID（它進入既有 StableSampleID/SampleEnvelope，不可誤當純 transient）；一個 coherent device snapshot 包含全部 requested Points 及各 Point 原 group/cadence，沿既有 immutable planner 的 group 隔離與唯一 request index，不造 fake group 或晚取 device/config 拼舊 Point。已知資料不足不得杜撰 cadence；若與既有 before-Start caller 契約衝突，先保原可用契約並明列證據，不以啟用來源補洞。

共用 executor 透過明確 final delivery policy 保 scheduled emitValue 與 manual return-only；manual 的既有 caller delivery／observer 語意保留，不同時再送 scheduled ValueChannel，known requested Points 在 caller cancel 仍回其原 Bad 結果、不額外重複。非 bounded PLC 的既有 custom protocol 保單次 scheduler-level scalar invocation，共用既有 count/read/result helper；adapter 原內部 retry、request budget/deadline 不因此變成只允許一次 wire attempt。MC/FATEK/Modbus 走同 shared planned executor，不為假 adapter 放行 unsafe scalar bypass。

3.3 完成 manual capacity parity 時，active permit pool 改為 constructor 建立的 Scheduler-lifetime 物件；Start 只重建本輪 due queue，不在 manual 活動中重置或拆分全域容量，也不新增 Start-busy 行為。這精確化3.2階段僅scheduled時的per-run permit實作；3.2歷史交接不當成最終manual契約。配置值仍constructor凍結；明確無效容量的manual fail closed，不永久等待。Stop／WaitStopped 繼續只取消並join原scheduled workers，manual依原caller context；跨Start/Stop/restart仍共用同device gate、fair active pool、原backoff cooperation與actual-send fence。舊scheduled workers先join，不能回返改新queue；manual合法release屬同lifetime pool。Manual caller原每request/device fanout不冒稱受scheduled per-configured-device goroutine上界涵蓋，不新增內部無界queue。

此為原3.3 migration的契約交接，沒有新public route、decoder、source/operation authority或持久ID；4.1預覽及4.2性能驗收仍獨立。


### 同源預覽的單位、head與實測契約

以下精確化原4.1／editor1.4單位與不可用邊界，未新增task、SourceEpoch、第二planner、SQL source-only規劃或診斷平台；同一次原生metadata更新只由一個owner計數。

#### 同一 installed snapshot 與分組估算

預覽只使用可證明 current 的 installed native acquisition snapshot，沿既有 capture/build helpers 及唯一 optimizer planner。沒有 current installation、無法證明原 polling group/cadence、或 context 不可證明時回明確 unavailable，不以空計畫或 count=0 假成功；不增加 SQL source-only planner 分支，不為預覽安裝來源或啟動 scheduler。已 installed 但尚無完成 scan 可以同時回傳有效 estimate 與 actual unavailable；因此不以執行過 PLC 讀取作估算前置。

device 回傳 groups[]，每 group 有自己的 cadence、plan fingerprint、planned request count；不同 group 不合成同一輪。每組 acquisition_point_count 是 captured acquisition Point membership（包括被 planner 明確拒絕的 Point），point_errors 另列，不能靜默消失。unique_wire_units 為有效 scalar physical-address spans 在同一完整 compatibility key/area 內的 union，排除單為 bridge gap 多讀的 cells；按 physical address unit 分列 bit 與 16-bit word cells。不得將 PointSlice.NativeBits（request component width）貼到 span length 當 cell width：FATEK DR 的 component 是 32 bits，但 physical address cell 是 16 bits；WM/DWM 的 physical address cell 是 bit。API 只消費 pure planner 的同源只讀 summary，不自行解析/換算/union；例如 DR0:uint32 的 span [0,2) 是2個16-bit地址單位，不能乘component的32；DR0+DR1 span union [0,3) 為3個16-bit地址單位；WM8:uint16 為[8,24)的16個bit地址單位，DWM8:uint64為[8,72)的64個bit地址單位。僅在既有canonical area/compatibility namespace內union，不自行推定跨D/DR/R等別名關係。不同 function/area/station/connection/group/cadence 不互相 dedup。這是 wire source geometry，不是 recording members 或 fetched request length。

planned_requests 表示該 group 當前已證明 context 下之初始執行策略的協議 request 數，不含故障重試與本 scan 未觸發的 fallback requests；零點/零 request 僅在完整有效 snapshot 證明該事實時才可為零。無法建立有效 snapshot 不能產生數值。

#### 唯讀 GET、head 與一致性

擬新增 GET 使用現有 `/api/v1/datalink/studio-v2/workspace/` prefix；建議 suffix `acquisition-plan-preview`（尚不存在，最終由 API owner 定名）。只接受 workspace_id、device_id、expected_workspace_revision 與可選 expected_preview_revision；拒絕任意 endpoint/address/member selector。

expected_workspace_revision 必須出現；空字串只依既有 workspace legacy-empty revision 語意驗證，不能等同缺少欄位。expected_preview_revision 首讀可省略，回傳同次 coherent snapshot 的 preview_revision 與 estimate，不先發另一個無保護的 head request；後續提供的 revision 不相符回 409。提供空 preview revision 視為無效輸入，不等同省略。query 重複值與不支援欄位由既有嚴格 parsing 慣例拒絕。

preview_revision 是既有 immutable plan/config/Point revision、installation lease/scope、group/cadence、capability/coverage 及所採 connection/downgrade context evidence 的 deterministic opaque composite，不是新的 SourceEpoch、持久 counter 或 authority。plan_revision 仍是 canonical planner fingerprint；兩者不能混稱。actual observation 的更新不改 preview_revision。

捕捉前驗 workspace/device ownership，建置後重驗同一 captured identity、workspace revision、installed slots/group/capability/connection context，變動回 409，不重抓新版後宣稱原 requested revision 成功。first read 無 expected_preview_revision 也需前後同一 snapshot。無法觀察 context（不是證實 changed）為 unavailable；若使用者提供舊 head 而目前已失去其 installation/context，舊 head 不再有效，回 stale 409。foreign scope 按既有 typed foreign error 拒絕，不能先洩露計畫或 metrics。

#### warm downgrade 與不探測

只讀 existing exact manager entry/token 與既有 downgrade cache。確認沒有 existing connection 時按 declared capability 建初始 plan，標 basis=declared_no_connection；有匹配的 existing entry/token 時，依同 capability revision 的 cache 決定 declared 或 safe initial plan。存在 entry 但 config/generation/token/context 無法證明，回 unavailable；不能當成不存在後樂觀顯示 random。無 dial、GetOrCreate、probe、reconnect、read、write 或設定變更。

preview evidence 明確區分 declared plan fingerprint 與 effective initial plan fingerprint：warm safe plan 是本輪初始策略，不是本輪剛發生的 fallback。connection replacement/cache change 使 context-bound head 失效；不能改寫 Point/sample 已有 fingerprint 或 AcquisitionID。

#### 最近完成 scheduled scan 的 actual 單位

只保留每 configured (device, polling group) 最近完成 scheduled scan 的 bounded native evidence，不建立歷史 registry。manual scan 不寫入此 view，也不能以 same device 或 subset point 結果充當 scheduled scan。保存 captured scope/install/Point/group/capability/plan/context proof，API 只有 proof 與 preview 完全匹配才取用；scope 證明不足或 revision 不符均 unavailable。遲到的舊 scan 不得取代新版 tuple 的 evidence。停用/移除配置時清除對應 tuple，restart 不借前次 lifetime evidence。

完成指 scan execution 路徑終止（包含失敗或取消），不是保證成功採集；不包含尚未 dispatch 的 rejected/coalesced due intents。每次記錄保留 AcquisitionID、declared plan、effective initial plan 與各 phase fingerprint 關聯；fallback phase 的 sample fingerprint 不改綁 initial。

actual 欄位採下列明確單位，保留原 attempted/succeeded/failed/retry/fallback/canceled 所需資訊，不宣稱每 wire success/failure：

| 欄位 | 單位與來源 |
| --- | --- |
| wire_attempts（attempted） | 原 wirebudget send observer 成功 admission/budget 後、即將 transport I/O 的 send 次數，含 scheduler/adapter 重送；不是 TCP packets |
| retry_wire_attempts（retry） | 各 logical request 第一次 actual send 之後的 actual sends；adapter self-retry 也算 |
| fallback_wire_attempts（fallback） | phase=fallback 的 actual sends；warm safe initial 不算本輪 fallback |
| logical_requests.succeeded | 一個已進入 readScheduledRequest 的 plan request 在完整驗證後 err=nil 且全部結果 Good 的最終結果，一 request 計一次 |
| logical_requests.failed | 已進 request 層而非 canceled 的最終失敗；明確 unsupported-command 造成 fallback 的原 request 也記一次 failed |
| logical_requests.canceled | 已進 request 層、最終為 cancellation/deadline interruption 的 request；不等於剩餘未執行 requests 數 |
| logical_requests.status | not_started（未進 request 層）或 observed；不能把 gate/dial/planning failure 的 Point Bad fanout 當 request failed |
| cycle_latency_ms | captured due 到 execution 終止的 monotonic elapsed，含 queue/gate/capacity/backoff/connection/read；不是 SQL latency |
| overrun | 同 scan elapsed > captured positive cadence 的布林；無 cadence 證據不可計算 |
| last_good_age | 僅能使用 matching captured Point/scan 的真正 Good completion 證據；無 Good 證據為 unavailable，不借其他 revision 的最後成功 |

retry_wire_attempts 與 fallback_wire_attempts 是 wire_attempts 的可重疊子集：fallback request 的第二次 send 同時屬兩者；不可相加當 total。wire_attempts 不等於 logical succeeded+failed+canceled，也不等於 Point count。plan-invalid Point 尚未形成 request 不列 logical failed；fallback replanning 所拒絕的 Point 同理。

真實已完成 scan 可以有 measured wire_attempts=0（例如 gate/dial 前失敗），但 logical_requests.status=not_started、三種 logical outcomes 不提供數值；仍可顯示真 latency。沒有 matching completed scan 時 actual 整體 unavailable，不能全補零。若已進 request 層而未送任何 wire，該 request 可以有實際 canceled/failed 結果與 measured wire=0，不與 not_started 混淆。顯示 true zero 必須有對應 observation。

取消發生於 request 完成後不得追溯改寫已確定 outcome；分類取該 logical request 的最終 validated error/result，不讀稍晚的 scan ctx 去改判。為避免 early-return 漏數，outcome 先記錄再處理 executor 的 canceled exit / fallback transition。

#### 已有 diagnostics 最小投影

不把 DueDiagnostics 全域 Coalesced/Rejected/Overruns 等數值灌入 scoped actual；本版 overrun 由同一最近 matching scan 的 elapsed/cadence 投影，不新增全域 counter。既有 output-loss counter 仍唯一 owner，僅投影能按 captured installation scope/Point revision/PointID 證明匹配的 retained detail，標 partial/history-incomplete；沒有 retained row 不代表零損失。EvictedDetailLoss/Unattributed/Total 全域值不得暴露為 scoped 總數；不得以 output loss 當 wire 或 SQL failure。最多 256 retained series 是既有設計上限，非量測值。


### 負載矩陣分批驗收

**Supersedes**: integrate-protocol-aware-acquisition-batching / 同planner預覽與量測證明（僅澄清長矩陣執行分批，並取代 Risks 的超時即新增子change排程規則；原量測門檻與矩陣不減）。

4.2保持原一項完整驗收；1/10/100設備、3/128點、1/20/100ms、至少100完整scan及cold/warm、fault/pressure條件全部保留。先短pilot實測時間，再依具名case分批安排，避免長測阻塞其他已授權介面驗證；每批保留版本、seed、完整參數、樣本數及結果，只有全部要求完成才勾4.2。超長case可分可續跑區段，但必須保存同一實驗的連線／warm狀態或明確重新分類cold，不能拼接不同條件冒充同次分布。任何仍未測／失敗case均明列，不把分批完成稱整體通過，不新增功能或第二引擎。

每device分位數必須列實際樣本數；100設備的一次sample不可冒稱100筆每device分布。Cold首scan單列啟動與首讀條件，不能用單樣本宣稱可靠cold p95；warm分布與setup/probe流量分開。同canonical設定／limits的singleton-planner逐Point正式adapter讀取可作test-only scalar基準，whole-set planner為paired策略；不得改能力／limits來偽造scalar，也不得將多次scalar塞入原optimized單request預算。正式scheduled負載、queue／fairness／drop仍須沿現有production path另證；adapter-only計時不冒稱包含scheduled排隊的cycle latency。既有無變code的相關通過證據可沿用，race耗時不能作效能比較。
