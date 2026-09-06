## Context

/studio/runtime 的產品角色仍是完成設定後看同一工作區／記錄方案。
新頁籤不是另開報表產品，先讓用戶能看見剛設定的電表與混合感測器紀錄。

## Goals / Non-Goals

目標：明細、曲線、用量與事件數字一致，可回查、有缺口標記與安全匯出。
不做：任意 SQL 報表設計器、排程寄信、財務帳單、AI 預測或未取樣尖峰推估。

## Decisions

1. 每 plan 明確選 query-capable primary destination，API 回傳來源 identity 與 delivered 範圍。
   禁把多個鏡像目的地資料相加。缺 SELECT 權限時顯示不可查歷史；可顯示本機暫存狀態，
   但必須標明 local_pending，而不是冒充目的地已保存的歷史。
2. samples/intervals/events/export 查詢均帶 workspace、plan、measurement/equipment、
   time range、timezone、resolution、quality filter 與 revision cutoff。
   先校驗 ownership，再做 parameterized bounded query；有限筆數、cursor 與取消／逾時。
   pagination/export 使用一致 cutoff，不能前一頁舊版、後一頁新版導致重複或漏項。
3. 趨勢：快速看最新與短期明細，長期使用已存摘要；回應實際 resolution、資料保留起點與來源。
   mean/min/max 與有效 coverage 同時可查，空窗顯示斷線；不補零、不無條件連線、不开假即時動畫。
   超出保存範圍顯示已過期；raw 已不存在時不能號稱每5秒曲線。
4. 電表將三相 A、三相 V、kW 分成有單位的圖組；phase_role 清楚，電壓接線不明時提示。
   kWh 原累積值曲線與區間用量柱形／表格分開；每日用量可點開查看起訖原值、缺口與算法版次。
   MC 例則各顯溫度、壓力、流量、累積用量、運轉狀態、警報與產量，不平均文字或累積值。
5. 日／時總量依3的結果；complete=false 時完整總量空值，known_subtotal另列。
   signed quantity 保留負號。跨設備合計要求兼容單位、同時間邊界與非重疊計量角色，
   缺一台時明確顯示幾台有資料，不把缺台算成零。
6. state/event 顯示初始狀態、發現時間、持續區間、unknown、清除時間、alarm bit與batch identity。
   由輪詢推得的運轉時長和 transition timing 顯示採樣限制，不誇大為精確事件時間。
7. 狀態列分設備採集、記錄本機落地、各目的地送達與查回驗證；
   點「修正」回 /studio/v2 的對應設備、項目／plan與欄位，而非從第一步重做。
   連線失敗保留最後成功結果但標過期，需明確刷新，不把舊值當新值。
8. CSV 匯出保留 timezone、observed_at、unit、quality、estimated、completeness、
   measurement identity、plan/calculation revision；精確數值用不變字串。
   textual cells 防 spreadsheet formula injection，保留合法負數數值；分頁串流、限時與範圍限制。
   不匯出憑證、DSN、SQL exception 或跨 workspace 資料；取消匯出釋放資源。
9. 性能驗收使用可重播負載與可量測門檻，不能只說「高效」。
   最小 fixture：100設備×8項每5秒（160樣本/秒），查24小時每分鐘摘要與最長保留範圍。
   在明確記錄的參考硬體上，目標24小時查詢p95≤2秒、無無界記憶體／goroutine成長；
   ARM限資源時允許調整已公告上限，但要測量、提示與拒絕超限，而非偷偷掉資料。
   這是待驗收目標，不是現有性能聲稱。外部批次大小與資料庫版本都記入測試證據。

## Risks / Trade-offs

報表真實顯示未知会出现空白，這是必要資訊而不是UI故障。
SQLite 與遠端資料庫查询能力不同，adapter 顯示有效保留／可查範圍；不可隱式切資料源。
長期報表精度取决於保存的摘要與證據，匯出必須帶粒度。

## Migration Plan

新增 runtime tabs與API，不改既有 status/stream 基本用途。
僅對具備 plan/history capability 的方案顯示功能；既有 latest-only 資料清楚標無歷史。
回退不刪已存資料；保留API與安全停止記錄能力，前後版本相容需測。

## Open Questions

部署硬體與最終性能上限在實作基準測試時記錄；案例、功能與正確性驗收不延後。
