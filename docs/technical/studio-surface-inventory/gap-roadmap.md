# Gap Register And Change Roadmap

## 目的

這份文件記錄目前最重要的前後端整合缺口，並把它們轉成可執行的 change 候選。  
不是所有缺口都要立刻做，但應該先有排序，避免下一輪規劃時又重新盤點一次。

## 目前優先級原則

1. `/studio`：保留完整版定位，但這一輪先暫停
2. `/studio/v2`：預設入口、重點施作
3. `/studio/runtime`：作為 V2 的 post-setup observer 一起補齊
4. `/test`：先暫停
5. `/gateway/*`：只保留記錄與 prototype 說明

## P0: 現在就該優先處理

| Change 建議名 | 問題 | 主要影響面 | 為什麼優先 |
| --- | --- | --- | --- |
| `make-studio-v2-default-entry` | 使用者希望預設入口是 `/studio/v2`，但 router 目前仍導向 `/studio` | router、入口 IA | 產品入口與使用者預期不一致 |
| `integrate-studio-v2-commit-with-runtime-lifecycle` | V2 commit 仍是前端模擬，不能提供正式 persisted handoff | `/studio/v2`、`/studio/runtime` | 不補這個，V2 無法成為真正產品主線 |
| `normalize-runtime-lifecycle-status-contract` | runtime snapshot / SSE 缺正式 lifecycle semantics | `/studio/runtime`、未來 `/studio/v2` | 前端目前還要靠錯誤字串判斷狀態 |

## P1: 緊接著要處理

| Change 建議名 | 問題 | 主要影響面 | 說明 |
| --- | --- | --- | --- |
| `refactor-studio-mainline-surface-boundaries` | `/studio` 功能完整但過於複雜，步驟責任交纏 | `/studio` | 重要，但依目前決策先暫停 |
| `promote-source-rule-apply-contracts-in-studio` | tags / outputs 的 apply endpoints 存在，但主線仍以細粒度 CRUD 組裝 | `/studio` | 可降低使用者操作複雜度 |
| `document-surface-api-inventory` | 這次新增的維護文件需要正式維護機制 | docs / review 流程 | 否則文件很快 drift |

## P2: 視方向決定是否要做

| Change 建議名 | 問題 | 主要影響面 | 說明 |
| --- | --- | --- | --- |
| `align-gateway-pages-with-datalink-product-contract` | `/gateway/*` 是否要轉正仍未定義 | `/gateway/*` | 目前先不推進，只保留選項 |
| `promote-readiness-first-device-flow` | `/devices/:id/readiness` 已存在但主線未全面採用 | `/studio`, `/studio/v2` | 可讓使用者看到更明確的 staged readiness |
| `add-runtime-e2e-handoff-coverage` | setup -> runtime 的整條路徑仍缺完整 e2e | 前後端整合測試 | 補產品信心 |

## 建議缺口類型

| 類型 | 定義 | 例子 |
| --- | --- | --- |
| `missing-api` | 後端沒有頁面需要的 API | V2 正式 commit contract |
| `exists-not-wired` | API 存在但前端沒用 | `source-rules/:id/tags/apply` |
| `semantics-gap` | API 可呼叫但語意不足 | runtime `device not found` / `starting` |
| `frontend-local-only` | 行為仍只存在 reducer / mock | `/studio/v2` commit flow |
| `connect-only` | 只做測試連線，不會持久化 | `/gateway/*` submit |

## 我額外補的維護重點

這些不是新功能，但很容易在後續規劃時漏掉：

1. **頁面受眾**  
   同樣叫 setup，不代表都是產品主線。文件必須持續標示 operator、一般使用者、工程工具的區別。

2. **source of truth**  
   每個頁面要清楚標示資料是來自 backend persisted state、runtime snapshot、SSE，還是 local reducer draft。

3. **handoff contract**  
   特別是 `setup -> runtime`、`entry -> quick/expert`、`studio mainline -> output review` 這些頁面切換，不能只記 route，還要記參數與前置條件。

4. **feature flags**  
   `/gateway/*` 受 `ENABLE_GATEWAY_DUAL_ENTRY` 控制，這層條件若不寫進文件，後續容易誤判為失效 route。

5. **verification path**  
   不只記 API，還要記完整操作鏈。例如：  
   `create device -> activate -> source rule -> points -> tags -> output -> runtime snapshot -> SSE`

6. **字串判斷風險**  
   任何前端若靠錯誤訊息字串判斷狀態，都應視為契約缺口，而不是「前端已處理完成」。

## change 拆分建議

### 1. `make-studio-v2-default-entry`

- 目標: 將產品預設入口從 `/studio` 轉為 `/studio/v2`
- 前端工作:
  - router default / wildcard redirect 調整
  - legacy redirect 對齊
  - 文案與入口 CTA 對齊
- 驗證:
  - `/`、未知 route、legacy route 都能落到正確入口

### 2. `integrate-studio-v2-commit-with-runtime-lifecycle`

- 目標: V2 commit 後能拿到 persisted IDs 與 runtime handoff context
- 後端工作:
  - 提供正式 commit endpoint 或 orchestration contract
- 前端工作:
  - Step 4 改用真實 commit response
  - 不再使用模擬 commit logs 當 truth
- 驗證:
  - `setup -> commit -> /studio/runtime?device_id=...`

### 3. `normalize-runtime-lifecycle-status-contract`

- 目標: 將 runtime snapshot / stream 語意正式化
- 後端工作:
  - 對 `device not found`、`starting`、`not running yet` 提供明確語意
- 前端工作:
  - 移除字串比對
  - route state 改用正式 status enum
- 驗證:
  - missing device、starting、live、degraded 都能穩定呈現

### 4. `refactor-studio-mainline-surface-boundaries`

- 目標: 降低 `/studio` 複雜度，但不破壞正式能力
- 狀態: **先暫停**
- 前端工作:
  - 先重畫 step boundaries
  - 把 runtime、review、output 邊界整理清楚
- 驗證:
  - 主線任務完成率不下降
  - deep-link 與 shared context 不破

### 5. `align-gateway-pages-with-datalink-product-contract`

- 目標: 決定 `/gateway/*` 長期命運
- 狀態: **先暫停**
- 選項:
  - 保持 prototype / connect-only
  - 轉正並接 datalink persisted APIs
- 驗證:
  - 文件、命名與行為一致，不再混淆

## 文件維護 checklist

每次新增或修改頁面 / API / route 時，至少問一次：

1. 這個頁面是給誰用？
2. 這裡的 source of truth 是什麼？
3. 對應哪個 backend contract？
4. 是 `wired` 還是只是 `exists-not-wired`？
5. 是否新增了新的 handoff / search param / lifecycle state？
6. 是否要補一條 roadmap change？
