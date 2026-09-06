# 感測資料記錄方案：整體實作入口

狀態：六份 OpenSpec 提案，尚未實作。日期：2026-09-07（Asia/Taipei）。
基準：`Yishow/go_gateway` 的 `main@4762f566ea28a1dbb4bbe313b156462d4af52517`。
工作分支：`spec/telemetry-recording-plans-20260907`。不修改 main，不自動啟動設備或寫入外部資料庫。

交付狀態：完整 45 份提案與配套文件隨本次文件 commit 提交至上述分支，未合併 main。
先前草稿包記載的提交阻擋屬上一輪交付狀態；本輪僅同步文件，不代表產品已實作或通過實機驗收。

## 要解決的事情

使用者不應先懂資料表、Tag、Row group，才能留下可用的資料。
主流程保留四步：連接設備 → 選擇資料 → 確認數值與用途 → 設定記錄並啟動。
先確認每個數字代表什麼，再選「留曲線／算用量／記狀態／留事件／批次紀錄／只看最新」，最後才選資料庫。
同一台設備、同一段讀取範圍，可以包含不同用途、單位、格式與記錄頻率。

## 六份 change 與依賴

| 順序 | Change | 交付邊界 | 前置 |
| --- | --- | --- | --- |
| 1 | [add-measurement-semantics](../../../openspec/changes/add-measurement-semantics/proposal.md) | 混合位址、逐項格式、用途、身份、品質與範本契約 | 無 |
| 2 | [add-recording-plans](../../../openspec/changes/add-recording-plans/proposal.md) | 記錄方案、資料表、型別保存、觸發與安全建表 API | 1 |
| 3 | [add-usage-and-aggregation](../../../openspec/changes/add-usage-and-aggregation/proposal.md) | 用量、時間統計、狀態時長、衍生值、補算 | 1、2 |
| 4 | [add-durable-recording-delivery](../../../openspec/changes/add-durable-recording-delivery/proposal.md) | 本機耐久暫存、補送、去重、版本切換與容量保護 | 1、2；整合驗收需 3 |
| 5 | [redesign-studio-recording-flow](../../../openspec/changes/redesign-studio-recording-flow/proposal.md) | 四步白話流程、批次套用、真實欄位、驗證與修復 | 1～4 |
| 6 | [add-recording-history-reports](../../../openspec/changes/add-recording-history-reports/proposal.md) | 趨勢、用量、事件、CSV 與全流程驗收 | 1～5 |

3、4 可在 2 的介面凍結後平行實作；5、6 可先做不接正式入口的測試元件。
任何草稿、模擬資料或未完成 adapter 均不得提前顯示為可正式啟動。
每份 change 有 proposal、design、tasks、delta specs；實作任務全部保持未勾選。
本文件是導覽；行為由各 change 的 delta specs 定義，共用欄位見 [contracts.md](contracts.md)。

## 使用者案例

### 多台 Modbus TCP 三相電表

每台各自對應 3 個電流、3 個電壓、1 個總功率、1 個累積電量。
前七項留明細與趨勢；kWh 留原讀值，另算區間用量。
確認電壓是相電壓或線電壓，電量是進電、出電或淨值；不靠單位猜測。
同型號可套用範本，仍逐台確認讀取位置、倍率、單位與實測值。
不把總表和分表自動相加；不把三相電壓或電流自動加成「總量」。

### MC Protocol 混合範圍

以下是測試示意，不是現場 PLC 位址表；一個 D 暫存器不必然等於一個完整感測值。

| 占用範圍 | 項目 | 示意格式與換算 | 用途 |
| --- | --- | --- | --- |
| D0 | 溫度 | int16 × 0.1 °C | 曲線、平均、採樣最高最低 |
| D1 | 壓力 | uint16 × 0.01 bar | 曲線、平均、採樣最高最低 |
| D2 | 瞬時流量 | uint16 × 0.1 L/min | 曲線；經確認可另估流量積分 |
| D3～D4 | 累積流量 | uint32 × 0.001 m³，明確指定字序 | 區間用量；不得再建立 D4 獨立點 |
| D5 | 運轉狀態 | uint16，0 停止／1 運轉／2 故障 | 變化紀錄、有效觀察時長 |
| D6 | 警報位元 | uint16 bitmask | 警報出現與解除，保留原始位元值 |
| D7 | 累積良品數 | uint16，是否回捲需明確設定 | 區間產量；不能自動猜回捲 |

八個 D 暫存器在此例形成七個主要項目；D6 可再衍生具名警報位元。
相容的範圍可合併通訊讀取，但逐項解碼、品質、用途與身份不可合併掉。
MC、Modbus、FATEK 共用記錄模型；不新增不實宣稱的協議能力。

## 完成標準

使用者能不寫 SQL 完成設備到記錄，且看到真實已存紀錄。
瞬時值不覆寫歷史；累積值不相加；缺資料不補零；重送不重複算用量。
讀取頻率、明細記錄、摘要頻率、批次送出與保存期限各自明確。
報表保留時區、缺口、估算與計算版本，不能用假完整度掩飾資料遺失。
既有設定與 Local Modbus 轉發保留，不重建 `/studio`，不要求只用轉發的人設定資料庫。

## 查閱與驗收

- [contracts.md](contracts.md)：跨 change 的資料、API、時間與身份契約。
- [acceptance.md](acceptance.md)：數值範例、異常情境、平台與資料庫驗收矩陣。
- [baseline.md](baseline.md)：主分支證據、規格衝突與資料來源。
- [validation.md](validation.md)：本次提案檢查與尚未執行項目。
