## Context

目前 Studio V2 的 bootstrap 依賴多個 query 與多個 autosave hooks 拼裝本地狀態，但真正的 persisted truth 分散在 workspace、devices、source rules、mappings、database config、database targets 等不同來源。這種混合模式在單次編輯流程可用，但一旦 hard refresh、服務重啟、部分 autosave 尚未完成，畫面就會退回本地預設或只恢復一部分資料，讓操作員無法判斷哪一份狀態才是真相。

這個 change 的目的不是新增離線草稿，而是先把「已保存資料必須可恢復」與「未保存草稿遺失必須被明示」建立成系統契約。只有先把 reload recovery 做對，後續 readiness、runtime、database delivery 才有可信的起點。

## Goals / Non-Goals

**Goals:**

- 讓 Studio V2 在 reload 後先以 persisted setup state 完成 hydrate，再顯示 Step 1 到 Step 4。
- 讓 Step 1 到 Step 4 的 save truth 反映後端真相，而不是本地預設或殘留 reducer 狀態。
- 讓 unrecovered draft loss 成為顯式事件，而不是沉默覆蓋。
- 讓 shell committed/runtime 指示器可在 reload 後恢復為已持久化 truth。

**Non-Goals:**

- 不在這個 change 內加入跨瀏覽器同步或離線 localStorage 草稿保存。
- 不在這個 change 內重新設計 activation readiness、runtime refresh 或 database write path。
- 不把每一個 Step 的所有 UI 細節都改成獨立持久化，只處理與 persisted truth/reload recovery 直接相關的契約。

## Decisions

### Persisted snapshot is server truth

Studio V2 reload 後的第一份真相來自已持久化後端資料，而不是前端上一次 reducer 記憶體。任何已持久化欄位在 reload 後都必須以後端值還原，不能讓 factory defaults 覆蓋 persisted state。

替代方案是保留現在的 reducer 預設，等 query 回來後再局部覆蓋，但這會在 bootstrap 過程中短暫顯示假資料，且會掩蓋哪些欄位其實根本沒有 hydrate 成功。

### Hydration runs in dependency order

reload 恢復不能平行地把所有 slice 直接灌回 reducer，因為 Step 3 mappings 與 Step 4 targets 依賴 Step 1 devices、Step 2 rules 與當前 point identity。設計上必須維持 devices -> rules/points -> mappings -> database 的依賴順序，並在相依資料尚未就緒前阻止畫面誤顯示空白或預設值。

### Unrecovered drafts surface as explicit loss

對於尚未到達 persisted state 的本地草稿，本 change 不承諾 reload 後恢復；但系統必須在 reload 後明示這份草稿已遺失，而不是讓使用者誤以為 autosave 成功。這個決策能把 persisted truth 與 local draft 的責任邊界切清楚。

### Activation summary is persisted separately from UI animation

Step 4 的 commit animation、逐條 log 動畫、按鈕 loading 是前端短生命週期狀態，不需要完整持久化；但 activation 是否曾成功、目前 workspace 是否處於 committed/running truth，必須有可恢復來源，讓 shell top bar 與 reload 後的 landing state 不會回到 idle 假象。

## Implementation Contract

- Behavior:
  - reload /studio/v2 時，Step 1 到 Step 4 顯示的欄位與列表必須以 persisted workspace state 為準。
  - 若某筆編輯在 reload 前仍未達到 persisted state，reload 後畫面回到最後 persisted value，並顯示 unrecovered draft loss。
  - shell committed/runtime indicator 在 reload 後必須反映已持久化 activation truth，而不是本地預設 idle。
- Interface / data shape:
  - Studio V2 bootstrap 組合 devices、source rules、mappings、database config、database targets 的 hydrate 流程。
  - 前端狀態需區分 persisted truth、saving、save-error、draft-invalid、discarded draft loss 等顯式狀態。
  - activation truth 需能讓 TopBar、WorkbenchV2Shell、Step 4 commit summary 在 reload 後判定 workspace 已提交或仍未提交。
- Failure modes:
  - 若任何 persisted slice hydrate 失敗，頁面進入 bootstrap error/degraded state，而不是靜默回退到 default setup。
  - 若 local draft 未恢復，系統只可回到 persisted value，且必須提示草稿遺失。
- Acceptance criteria:
  - reload 後 Step 1 到 Step 4 已保存欄位完整恢復。
  - 未保存草稿在 reload 後明示遺失。
  - 已 activation 的 workspace reload 後，shell committed/runtime indicator 不回到預設 idle。
  - 受影響前端單測與 bootstrap/reload 測試覆蓋這些情境。
- Scope boundaries:
  - In scope: persisted state restore、draft loss signaling、shell activation truth restore。
  - Out of scope: offline drafts、跨瀏覽器同步、runtime refresh 策略、database write reliability。

## Risks / Trade-offs

- [Hydration 順序更嚴格] → 需要更明確的 bootstrap loading/degraded state，初次載入感知可能變慢，但能換取較少的假資料閃爍。
- [顯式顯示草稿遺失] → 使用者會更常看到「沒存成功」的現實，但這比沉默錯誤安全。
- [activation truth 持久化來源可能需要新欄位或新投影] → 先用最小可恢復模型，避免把整套 commit animation 歷史一起持久化。

## Migration Plan

1. 先補齊 bootstrap / hydrate 契約與前端測試。
2. 再補 draft loss signaling 與 shell activation truth restore。
3. 最後用 manual reload scenario 驗證 service restart、browser refresh、partial autosave 三條路徑。

## Open Questions

- selected tab / current step 是否也要跨 reload 恢復，還是只保證資料內容恢復？
- unrecovered draft loss 要以 toast、inline banner、還是 summary rail 提示為主？
