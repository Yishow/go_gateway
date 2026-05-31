## 1. Workspace runtime context

- [x] 1.1 交付 `Runtime page is workspace-scoped first`：讓 `/studio/runtime` 先讀 singleton workspace device set，而不是只能靠單台 `device_id`，並以 backend runtime context test 與前端 route boot test 驗證。
- [x] 1.2 交付 `Default device follows v2 order`：讓 runtime 預設選擇 V2 排列最前面的可用設備，並以前端 runtime test 驗證 default selection 不受事件先後影響。

## 2. 切換與空狀態

- [x] 2.1 交付 `Unavailable devices remain visible`：讓 unavailable device 保留在 runtime 切換清單中並帶出原因，並以前端 UI test 驗證設備不會消失。
- [x] 2.2 交付 `Empty runtime state stays on page`：在沒有任何可用設備時保留 runtime 空狀態與返回 V2 指引，並以前端 runtime test 驗證不會自動 redirect。

## 3. API 契約

- [x] 3.1 交付 `Studio v2 workspace runtime context API`：提供 `GET /api/v1/datalink/studio-v2/workspace/runtime-context` 契約，並以 router / handler test 驗證 ordered device payload 與 `default_device_id`。
- [x] 3.2 依照 `### Runtime page is workspace-scoped first` 保留 `device_id` query 作為 optional override，而不是唯一 truth，並以前端 route test 驗證 query 存在與不存在兩種開頁方式。
