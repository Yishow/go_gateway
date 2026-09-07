## Context

以 main@4762f566 為基準；參考 point-catalog、source-rule-runtime、source-rule-target-datatype。
既有規則擁有 revision 與下游身份，且禁止未對齊模型就單獨擴張 parser。
本 change 是讀取與記錄中間的共用資料意義層，不是新協議或 PLC 控制功能。

## Goals / Non-Goals

目標：同一設備混合格式與用途可正確解碼；來源身份、精度與品質可供長期記錄。
不做：資料庫 delivery worker、用量演算法、任意腳本、PLC 寫入、全自動猜測現場接線。

## Decisions

1. 以 measurement_id 關聯既有 point_id/tag_id；equipment_id 表示實際設備／回路，不能只用 IP 當身份。
   Gateway 經同一網路端點讀到多台邏輯設備時，仍保留站號與來源範圍。
2. semantic_kind 採共用契約八種；rate 是瞬時速率而非累積量；delta 必須有事件或區間身份。
   counter 預設單調累積；signed_counter 必須明確確認可下降的淨值。
   text/state/event 不提供數值平均；bitmask 保留原值並可衍生具名 bit。
3. Source rule 增加 versioned `layout_mode=homogeneous|mixed` 與 mixed item definitions。
   item 包含 stable item_id、offset/address、width、read_type、word/byte order、scale/offset、
   target_type、measurement binding；逐項有效設定只計算一次。
   homogeneous 完全走原行為；mixed 未覆寫欄可繼承，但預覽必須顯示最終有效值及來源。
4. 位址規畫先保留 span，再交既有協議讀取優化器；同一相容請求可讀一段，再按 item 解碼。
   不能保證整段 PLC 掃描原子性；跨 request 保留各自 acquisition/time。
   32 位項目占兩個 16 位位置；64 位與字串依格式保留完整長度，不允許邊界截斷。
   bit slices 可顯式共用一個 word，其餘 span 重疊預設拒絕。
5. 先解碼，再應用既有有效換算／cast 管線；阻擋重複 scale、溢位、NaN/Inf、不合法 enum。
   uint64/精確小數使用 decimal string 傳遞；顯示四捨五入不可改變保存或累積計算值。
6. 名稱與畫面排序調整不換身份。改位址、倍率、單位意義或換表需新 series_epoch；
   讓使用者預覽哪些歷史可比、哪些不能續算。拆分舊 span 不自動把舊 counter checkpoint 接給新點。
7. 範本記錄型號／版次／角色／格式／占用長度與建議用途，不附真實密碼。
   批次套用按設備產生 preview diff；各台讀值、接線角色、CT/PT 倍率和單位需確認。
   存在 PLC 端已換算工程量時，不可再自動加一次倍率。

## Risks / Trade-offs

mixed 讀取可能增加請求數，先保正確性再合併；無同步證據不宣稱各項同時量測。
既有大量規則需加法式 migration；禁止為簡化新 UI 改掉其原始值或預設用途。
用量資料精度不等同感測器物理精度；保留來源解析度，不誇大高位數意義。

## Migration Plan

新增定義與關聯表／欄，既有來源標示 legacy/unclassified，維持原寫入行為。
新方案啟用前逐項確認用途；回退停用新功能，保留新增歷史與定義，不做 destructive down。
測試原有 uniform rules 與 Local Modbus span/ownership；API、runtime、UI 一起通過後才開 mixed 選項。

## Open Questions

現場位址表、字序、單位、回捲上限由使用者在設定時提供，不阻擋實作固定契約。
沒有資料時走「尚未確認」，不以範例做生產預設；沒有其他阻擋本提案的產品決策。
