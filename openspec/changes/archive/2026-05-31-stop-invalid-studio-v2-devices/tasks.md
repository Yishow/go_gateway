## 1. 停止規則

- [x] 1.1 交付 `Invalid edit stops running device`：讓 running device 在 invalid edit 出現時立即停止，並以 backend runtime test 驗證狀態會從 running 轉成 stopped / unavailable。
- [x] 1.2 依照 `### No silent last-good fallback` 固定系統不再默默沿用舊的 running 狀態，並以 service test 驗證 invalidation 後不會繼續報告可用。

## 2. 可用性狀態

- [x] 2.1 交付 `Studio v2 device availability status` 與 `Availability state is explicit`：在 workspace / runtime device payload 中加入 `availability_status` 與 `availability_reason`，並以 handler test 與前端 type test 驗證欄位完整。
- [x] 2.2 交付 `studio-v2-device-validity` 的 UI state 呈現：讓 unavailable device 仍保留在集合中但可被辨識，並以前端 state test 驗證 unavailable 不會被直接移除。
