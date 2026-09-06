## Context

演算法消費帶身份、品質與觀測時間的 SampleEnvelope，不依賴 writer flush tick 或前端表格。
共用來源：contracts v1；counter/state checkpoint 的保存與結果落地必須由一個本機 transaction 完成。

## Goals / Non-Goals

目標：曲線與計量保留原始證據，明確區分實測、已知小計、估算與未知。
不做：電費／稅務計算、需量費率判定、保護電驛、快速脈衝保證捕獲、未經確認的資料補齊。

## Decisions

1. 視窗 `[start,end)`，UTC 存時間，依 plan timezone 切日；台灣預設 Asia/Taipei。
   事件時刻與收到時刻分開，設備時間不可信時採 Gateway 回應時間並標來源。
   禁止用 flush 時間分桶；clock jump、同時刻不同值與來源倒退記 uncertain，不能擅自猜順序。
2. Gauge/rate：保留有效樣本的 count、min/max 及 timestamp、first/last 與時間加權平均。
   預設以有效值保持至下一個樣本，但最多保持 max_hold（初始建議 2×採樣間隔，需在 plan 明列）。
   `mean = Σ(value_i × covered_seconds_i) / valid_duration`；有效段裁切到視窗邊界。
   count=0 時 mean/min/max 為 null，不是零；品質變壞立即中斷有效段。
   max_hold 之後是 unknown，不無限延用。summary 記錄方法與 coverage，min/max 稱「採樣最高最低」，
   不能宣稱捕捉到採樣間隔之間的尖峰。較大視窗以加權分子／有效時間合併，不平均各小窗平均值。
3. Counter：按同一 measurement+series_epoch 的有效樣本排序；首筆建立基準，後筆差分。
   相等為合法零用量；不同設備、倍率、計量世代不得相減。保存 decimal 差值，不提前四捨五入。
   負差、不合理速率或超範圍先隔離可疑樣本、標 discontinuity，不使用 abs 或自動歸零。
4. 重置預設需確認或可信 reset signal；保存事件前後值與時間，開新 epoch／segment。
   無法知道重置間漏掉多少，就保留不完整度，不把重置後的第一個數字冒充整段用量。
   回捲僅在明確 modulus M 與可驗證至多一次回捲時計 `(M-prev)+curr`；
   如65530→4且M=65536，差值10。沒有上限、最大合理速率或不能排除多次回捲則不自動補算。
5. signed_counter：允許淨值下降，差值可負；必須明確標識 net/import/export 角色。
   delta：來源必須提供不重疊區間與穩定 sequence/event identity，同一筆重送只算一次。
   同值、不同有效 event id 是兩筆；相同寄存器在每次 poll 仍顯示同值不代表有新用量。
6. 長缺口兩端若可信且同 epoch，保留整段差值與 span；沒有內部邊界證據不能分到每分鐘。
   若已知歸零可能发生且無連續性證據，連整段差值也標 uncertain。
   日界線需要有效邊界讀值；缺邊界顯示日量未知及已知小計。
   預設不插值；使用者啟用線性分配時必留 estimated=true、方法、缺口長度，且不覆寫原實測結果。
7. Rate integration：只對 quantity/unit 可換算的速率開啟，預設相鄰有效採樣梯形積分。
   kW×seconds/3600→kWh；L/min×seconds/60→L。
   超過 max_gap 的段落不積分；結果是估算，不偽装 meter counter。
   同一設備已有權威累積量時，不把積分量再加到累積差值；需選 primary 或用作比較。
8. State：首筆記 initial snapshot，不虛構轉態；確認的值改變才記 transition。
   運轉時間按有效觀察段計；缺口變 unknown，不能把上一狀態無限延長。
   除非事件來源提供精確時間，輪詢轉態只能反映發現時間，報表須揭露採樣精度。
   bitmask 支援每位元 transition；debounce 可設定但必須揭露可能略去短事件。
   text／批號只留值與變化，不平均；held level 不當成每次新事件。
9. 衍生值只用型別化運算白名單（+、-、×、÷、顯式單位換算與量測引用），
   禁 eval、任意 Go/JS/SQL。驗證單位、除零、缺值、DAG 循環與來源時間偏差；
   輸出保留公式版次與全部 evidence。跨設備加總需要明確非重疊計量範圍，阻擋總表+分表重複計量。
10. 每 series 單序列 checkpoint；有限亂序 grace 及 correction_horizon 明列在 plan。
    寬限內排序後結算；已有結果遇晚到樣本，在證據仍保留時重算受影響區間，
    產生同 result_id 的更高 calculation_revision，與 outbox 同 transaction 提交。
    舊 revision 重送不能覆蓋新結果；不把 late sample 當 counter reset。
    證據已過期、時間衝突或來源身份不明時拒絕自動重算，保留原因及人工修正事件。
11. 壓縮後的摘要要保留可合併統計量與邊界資料；不足以精確重算則明示不可重算。
    檢測重置所需原值、counter checkpoint、receipt 與最終結果之保留期限分開。
    計量誤差以來源解析度揭露，不承諾財務級電費結算。

## Risks / Trade-offs

時間加權採保持模型、速率積分採梯形模型，都是明確近似，不等於連續物理量測。
完整用量與已知小計分開會讓某些圖有空白，但比錯誤補零可靠。
有限亂序寬限增加最終報表延遲，必須顯示 provisional/final 和最後更新。

## Migration Plan

新引擎只作用於明確啟用的 plan；現有 raw 記錄不自動重解。
回退保留 checkpoint、epochs、已存結果與版本；重新上線從 durable checkpoint 接續。
資料定義變更不靜默重算歷史，需預覽受影響時間與證據範圍。

## Open Questions

現場最大合理流率、回捲上限、reset signal 是可設定資料；缺少時採保守阻擋，非開發待定。
