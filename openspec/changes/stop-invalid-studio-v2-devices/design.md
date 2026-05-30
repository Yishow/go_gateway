## Context

V2 autosave 讓資料持續進到後端後，系統還需要清楚定義「什麼時候一台設備變得不可用」。使用者的決策不是保留舊設定繼續跑，而是只要目前畫面已改成不合法，就把這台設備立刻停掉並標成不可用。

## Goals / Non-Goals

**Goals:**

- running device 遇到 invalid edit 時立即停止
- 系統能回報 `available` / `unavailable` 與原因
- runtime 後續可以顯示這台設備還在，但現在有錯

**Non-Goals:**

- 不在這個 change 內決定 runtime 的切換 UI 排版
- 不處理 first activation
- 不用舊的 last-good-running 行為當 fallback

## Decisions

### Invalid edit stops running device

一旦設備目前編輯值變成 invalid，就不允許它繼續被視為可用中的 running device。這是使用者明確選擇的嚴格模式。

### Availability state is explicit

前後端都需要有清楚的 availability state，而不是靠字串湊推論。至少要能表示 `available`、`unavailable` 與 `reason`。

### No silent last-good fallback

不能對使用者說目前設備已經改壞了，卻又默默讓舊的 runtime 繼續當作正常可用。這會讓 `/studio/v2` 與 runtime 的 truth 分裂。

## Implementation Contract

- Behavior:
  - running device 遇到 invalid edit 時，runtime 立即停止該設備
  - 該設備狀態轉成 unavailable，並帶出原因
  - 裝置仍保留在 workspace device set 中，不會消失
- Interface / data shape:
  - workspace / runtime device status payload 至少包含：
    - `availability_status = "available" | "unavailable"`
    - `availability_reason`
    - `running`
- Failure modes:
  - invalidation 處理失敗時，不得把設備假裝成仍可用
  - 前端不得只靠 generic error string 判斷 unavailable
- Acceptance criteria:
  - backend runtime tests 覆蓋 running -> unavailable -> stopped
  - frontend state tests 覆蓋 unavailable 狀態與原因可被正確顯示
- Scope boundaries:
  - In scope: invalid edit stop rule、availability state
  - Out of scope: runtime 切換 UI、first activation
- Not complete if:
  - 改成 invalid 後設備仍繼續跑
  - 前後端沒有明確 availability 狀態
  - unavailable device 直接從集合消失

## Risks / Trade-offs

- [Risk] 嚴格 stop 規則可能讓使用者在編輯途中就中斷收集 → Mitigation: 這是明確產品決策，UI 需清楚提示此變更會使設備不可用
- [Risk] availability 與 runtime health/breaker state 混在一起 → Mitigation: availability 是配置合法性，health 是執行品質，兩者分欄位表達
