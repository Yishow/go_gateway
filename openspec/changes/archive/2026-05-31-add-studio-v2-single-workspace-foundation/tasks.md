## 1. Workspace 模型與持久化

- [x] 1.1 依照 `### Singleton workspace model` 交付 `Singleton studio v2 workspace`：建立唯一 workspace 的資料模型、儲存與重啟後重讀能力，並以後端 integration test 驗證首次建立後可在服務重啟後讀回同一個 workspace。
- [x] 1.2 依照 `### Legacy isolation` 交付 `Workspace excludes legacy /studio data`：讓 workspace 初始化時不自動匯入 `/studio` legacy devices 或 rules，並以內容檢查與 integration test 驗證 bootstrap 回傳的資料只屬於 V2。

## 2. Bootstrap API

- [x] 2.1 交付 `Studio v2 workspace bootstrap API`：提供 `GET /api/v1/datalink/studio-v2/workspace` 讀取即自動建立的契約，並以 handler / router test 驗證 response 至少包含 `id`、`kind`、`status`、`ordered_device_ids`。
- [x] 2.2 依照 `### Auto-bootstrap on read` 固定錯誤路徑與可行動訊息，不讓 bootstrap 失敗時回傳半成品 workspace，並以 API error path test 驗證。

## 3. 前端啟動

- [x] 3.1 交付 `First-open workspace bootstrap`：讓 `/studio/v2` 頁面啟動時先讀 workspace API，不再只靠本地預設 state，並以前端 boot test 驗證首次開頁會自動拿到 workspace。
- [x] 3.2 依照 `### Workspace metadata contract` 把 workspace metadata 收斂成前端 hook / service 可重用契約，並以型別檢查與 hook test 驗證欄位完整性。
- [ ] 3.3 交付 `Empty workspace renders true empty state`：當 singleton workspace 尚無 devices / rules 時，`/studio/v2` 不得靠 seeded local draft device/rule 撐頁，並以前端 route boot / empty-state test 驗證真實空狀態。
