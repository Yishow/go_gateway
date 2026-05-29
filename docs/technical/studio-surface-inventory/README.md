# Studio Surface Inventory

## 目的

這組文件用來維護 `go_gateway` 前端頁面與後端 API 的對照關係，重點不是列出所有 route，而是回答以下問題：

1. 這個頁面是給誰用的？
2. 使用者在頁面上能做什麼？
3. 前端真正呼叫了哪些 API？
4. 哪些 API 只是存在，但尚未接到頁面？
5. 哪些行為還停留在前端 draft / mock / connect-only？
6. 後續拆 change 或規劃時，應先補哪一段契約？

## 目前產品決策

1. `/studio`
   - 定位：完整版主線，功能應可全面覆蓋需求面
   - 決策：先暫停，不列為這一輪重點施作
2. `/studio/v2`
   - 定位：面向使用者的簡化版入口
   - 特性：簡單、有指引、容易觀察
   - 決策：列為網站服務預設入口與這一輪重點施作
3. `/test`
   - 定位：工程測試工具
   - 決策：先暫停，後續再完善
4. `/gateway/*`
   - 定位：實驗性 / prototype surfaces
   - 決策：保留記錄，但不列為目前主產品優先線

## 閱讀順序

### 接手 AI / 維護者第一入口

1. [START_HERE.md](./START_HERE.md)
2. [context.json](./context.json)
3. [CURRENT_STATE.md](./CURRENT_STATE.md)
4. [總覽 HTML](./index.html)

### 需要展開時再讀

1. [Studio 主線 `/studio`](./studio-mainline.md)
2. [Studio V2 與 Runtime `/studio/v2`、`/studio/runtime`](./studio-v2-runtime.md)
3. [Test 工具 `/test`](./test-tooling.md)
4. [Gateway 原型 `/gateway/*`](./gateway-experiments.md)
5. [後端 API 登錄表](./backend-api-registry.md)
6. [缺口與 change roadmap](./gap-roadmap.md)

## Surface 總表

| Surface | Route | 主要受眾 | 目前定位 | 主要資料真相 | 整合成熟度 | 目前策略 |
| --- | --- | --- | --- | --- | --- |
| Studio Mainline | `/studio` | 進階操作者、維運人員 | 正式主產品入口，但仍過於複雜 | `datalink` CRUD / review / output APIs | 高，但需要重整 | 暫停 |
| Studio V2 | `/studio/v2` | 一般使用者 | 簡化導向型 setup flow | 目前多數仍是前端 reducer state | 低到中 | 重點施作 |
| Runtime Dashboard | `/studio/runtime` | setup 後操作者 | post-setup focused monitor | `runtime/status` + `runtime/stream` | 中 | 追隨 V2 一起施作 |
| Test Tool | `/test` | 工程測試使用者 | 工程工具，不是產品主線 | `/test/*`、`/debug/*` | 中 | 暫停 |
| Gateway Entry / Quick / Expert | `/gateway/*` | 實驗性雙入口 | feature-flag 保護下的探索面 | `settings` flag + `/test/*` connect-only | 低 | 僅記錄 |

## 目前最重要的產品事實

- 使用者期待網站服務預設入口是 `/studio/v2`，但目前 router 預設仍是 `/studio`。
- `/studio` 已經深度綁定正式 datalink APIs，是現在最完整的產品主線，但這一輪先暫停。
- `/studio/v2` 已具備導覽、引導與 post-setup handoff，但 commit 仍不是正式 backend-persisted lifecycle，因此雖然是重點線，仍有關鍵契約缺口。
- `/studio/runtime` 已經是可用的 snapshot + SSE 監看頁，應作為 V2 setup 後的 focused observer，而不是獨立 fleet dashboard。
- `/test` 是工程工具，後續可持續擴充，但不應綁住這一輪產品主線。
- `/gateway/*` 目前不應被誤認為 datalink 主產品流程，因為它們主要依賴 `/test/*` 的 connect-only API，而不是 datalink persisted APIs。

## 維護原則

新增或修改頁面時，至少同步更新以下內容：

1. 這個頁面的受眾與定位是否改變。
2. 使用者新增了哪些前端動作。
3. 對應的 query / mutation / SSE / local draft 狀態是什麼。
4. 是否新增 route params、search params 或 handoff contract。
5. 是否有 API 已存在但尚未 wiring，或前端開始依賴字串判斷等脆弱語意。
6. 驗證路徑是否改變，例如 `create -> activate -> runtime -> stream`。

## Onboarding 規則

- 接手 AI 不應一開始就掃完整包 md/html。
- 先讀 [START_HERE.md](./START_HERE.md)、[context.json](./context.json)、[CURRENT_STATE.md](./CURRENT_STATE.md)，只在資訊不足時再往下展開。
- 若任務只涉及 `/studio/v2` 或 `/studio/runtime`，預設先看 `studio-v2-runtime.md`，不要先掃 `/studio` 主線文件。

## Changelog 機制

- `studio-surface-inventory` 的維護紀錄使用 SQLite，檔案路徑是 `docs/technical/studio-surface-inventory/changelog.sqlite`。
- 任何對這個目錄下 `md`、`html`、`js`、`css` 的修改，都必須同步新增一筆 changelog；不要只改文件。
- changelog 寫入工具：

```bash
go run ./cmd/studio_inventory_changelog init
go run ./cmd/studio_inventory_changelog add \
  -surface "studio-v2-runtime" \
  -summary "補齊 runtime handoff 文件與預設入口策略" \
  -files "docs/technical/studio-surface-inventory/README.md,docs/technical/studio-surface-inventory/studio-v2-runtime.md" \
  -reason "讓後續規劃可追溯 inventory 變更"
go run ./cmd/studio_inventory_changelog list -limit 20
```

- 建議欄位用法：
  - `surface`：例如 `overview`、`studio-mainline`、`studio-v2-runtime`、`test-tooling`、`api-registry`
  - `summary`：一句話描述這次改了什麼
  - `files`：逗號分隔的實際修改檔案
  - `reason`：為什麼要改，不要只寫「update」
  - `notes`：可選，用來補充限制、後續待辦或驗證方式

## 狀態標記

| 標記 | 意義 |
| --- | --- |
| `wired` | 前端已穩定使用，且契約語意足夠 |
| `exists-not-wired` | 後端 API 存在，但前端尚未使用 |
| `frontend-local-only` | 行為仍停留在前端 reducer / local draft |
| `connect-only` | 僅做測試連線，不持久化進 datalink domain |
| `semantics-gap` | API 可呼叫，但產品語意不足 |
| `needs-change` | 建議拆成正式 change 處理 |
