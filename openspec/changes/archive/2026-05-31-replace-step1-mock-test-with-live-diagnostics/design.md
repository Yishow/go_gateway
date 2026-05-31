## Context

Step 1 已經有 autosave、availability、runtime apply 等後端契約，但 `執行測試` 仍停在前端本地動畫。使用者已明確要求 `/studio/v2` 全站對接，因此 Step 1 不能再用 mock test 冒充真實診斷。

## Goals / Non-Goals

**Goals:**

- `執行測試` 對目前 draft 或 persisted device 發出真實 backend request
- connect / probe 的成功或失敗來自 backend diagnostics，而不是本地合成
- Step 1 continue gate 不再接受 fake success
- 失敗時保留可行動錯誤與階段資訊

**Non-Goals:**

- 不重做 Step 1 autosave 契約
- 不改寫 Step 2-4 autosave 或 first activation
- 不處理 settings backend wiring

## Decisions

### Backend-produced diagnostics only

Step 1 的測試狀態必須由 backend 回傳驅動。前端可以做 loading / rendering state，但不可自行用 timer 或隨機值推導「成功」。

### Draft-first diagnostics contract

V2 的 Step 1 編輯器允許尚未完全保存的當前表單值，因此 diagnostics 契約必須能接受 current draft config，而不是只接受已存在的 device id。若既有 `test-draft` / readiness API 足夠，優先重用；若不足，補齊而不是另開平行 mock flow。

### Failure must remain visible

測試失敗時，UI 應保留 connect/probe 的階段落點與訊息，讓使用者知道是連線層失敗還是協議層失敗，而不是只回 generic error。

## Implementation Contract

- Behavior:
  - 點擊 `執行測試` 會對 backend 發出真實 diagnostics request
  - connect/probe stage 的結果由 backend response 決定
  - diagnostics success 才能讓該 device 進入 tested state
  - diagnostics failure 不得讓 device.status 變成 `tested`
- Interface / data shape:
  - 以既有 `POST /api/v1/datalink/devices/test-draft` 或等價真實 endpoint 為主
  - 若需要 persisted readiness 補充，可重用既有 readiness contract，但前端不得因此回退到 mock
  - response 必須足夠讓 UI 區分 connect / probe outcome 與訊息
- Failure modes:
  - backend failure 時不得合成 success stage
  - draft validation 不足時不得偷偷改以假資料跑通
  - 前端不得保留 operator-visible 的 setTimeout 成功動畫
- Acceptance criteria:
  - backend handler / service tests 覆蓋 success、connect failure、probe failure
  - frontend Step 1 tests 覆蓋真 request 發送、成功 gate、失敗不通過 gate
- Scope boundaries:
  - In scope: Step 1 live diagnostics、stage rendering、continue gate
  - Out of scope: Step 2-4、settings、runtime dashboard rewrite
- Not complete if:
  - Step 1 仍透過 timer 或 Math.random 合成 success
  - failure path 仍把 device 標成 tested
  - 使用者在沒有真實 diagnostics evidence 時仍可被引導到下一步

## Risks / Trade-offs

- [Risk] 既有 `test-draft` response 對 V2 的 stage rendering 不夠細 → Mitigation: 在既有 endpoint 上補足 V2 所需欄位，不另造平行 mock contract
- [Risk] draft 與 persisted state 並存時容易混淆 diagnostics 對象 → Mitigation: request payload 一律以當前編輯值為準，UI 明確綁定當前 tab
