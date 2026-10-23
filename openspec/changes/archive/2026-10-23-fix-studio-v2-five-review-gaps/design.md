## Context

既有四步與group pipeline已可寫SQLite；review的乾淨單點D3 uint32因Count二次寬度得到type-mismatch。範本真API回confirmed=false；snapshot required仍優先於partial。

## Goals / Non-Goals

Goals: 修復五項既有可觀察契約，最小修改，保留既有能力。
Non-Goals: 新平台、協定重寫、Windows native驗收、正式PLC/DB、部署。

## Decisions

### 單點協定數量與直接批次分離

沿既有connector建立單點Count helper；MC用1，其他協定維持RegisterCountForDataType。不改MC直接批次的乘寬度語意。

### 既有UI依真實證據與明確意圖顯示

確認按鈕才送confirmed=true；設備/範本切換作廢preview與過期回應。既有安全error工具呈現失敗。Header只陳述running/availability/stream，不自行組合SQL成功。範圍分別表示起址跨度與完整占用；required保留，只揭露並連既有group editor。

## Implementation Contract

1. R1 單點periodic/manual/read-now傳MC Count=1；直接Read Count=2仍回兩值。測32/64-bit scalar、Modbus/FATEK原契約與MC三桶SQL。
2. R2 TemplateApplyPreview僅明確確認送confirmed=true；未確認backend仍拒絕；切設備/範本與stale async結果不得誤套。query/preview/apply錯誤安全可見可重試，明確操作離線不暫存至重連自動送出；reload保持展開，測斷線/409/500及真API。
3. R3 FocusedDeviceHeader不得由available或running推論SQL成功；stream狀態沿既有來源。測等待、Share-only、停止、失敗、unavailable。
4. R4 RangeSummary終點涵蓋最後項寬度，stride與width分別核實，使用既有addressParser處理進位。測MC/Modbus、32/64-bit、單word/bit，不改點位。
5. R5 Basic清楚說required仍阻擋整列、顯示目前必填並連既有member required控制；不暗改群組；既有受管unsigned精確TEXT成員不誤報不相容，自訂表檢查保持。更新後照常重新DDL預覽/確認，不繞版本保護；相容no-op的statements正規化為空陣列，先於token digest與儲存，前端嚴格parser不放寬。驗required skip、optional partial NULL/原因、恢復後寫入。
6. 所有錯誤/版本/DDL/unknown/Retry/Skip保護不變；lint/typecheck/tests/build與review才可完成；Windows/PG未跑不得填PASS。

## Risks / Trade-offs

Count單位改錯影響其他協定 → 入口真request測試加直接batch保護。Async preview跨scope → scope fence忽略舊回覆。停用仍排空accepted backlog → 觀察穩定桶而非立即行數。測試只用隔離loopback/temp SQLite，現場與Windows限制明列。
