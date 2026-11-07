## Context

真 Timescale/SQL Server 補測揭露四個測試 fixture 問題，產品邏輯不在本案變更範圍。

## Goals / Non-Goals

**Goals:** 測試精確保留資料庫證據、避免受真實時鐘邊界影響，且清理能恢復原測試環境。

**Non-Goals:** 不改產品流程、精度、權限與 fail-closed 行為，不以重跑掩蓋不穩定。

## Decisions

### Timescale 清理使用新連線

DROP extension 後以新 pgx pool 執行 CREATE，避開已載入 extension 的 backend，結束時關閉新 pool並保留成功斷言。

### 同批 Apply 固定時鐘

群組確認後擷取單一時間，套用同批群組；defer 還原原有預設 wall clock，保留所有 EffectiveAt 相同斷言。

### NativeMonthly 證據完整傳遞

測試 converter 使用既有 glmTargetPartition 轉換，不以缺失證據要求產品放行 native schema。

### 時間以瞬間而非字串比較

RFC3339Nano 解析後使用 time.Equal，保留奈秒精度與解析失敗斷言，只接受相同瞬間。

## Implementation Contract

只修改四個測試檔。驗收是真 Timescale recovery/cleanup 與 SQL Server native/runtime 回歸；原先失敗不可只改為 skip。最終 API test/lint/vet由整合工作執行。

## Risks / Trade-offs

不能以此 fixture 修復推論產品 schema-confirm 離線重放已修好；該產品缺陷由另一 change 負責。
