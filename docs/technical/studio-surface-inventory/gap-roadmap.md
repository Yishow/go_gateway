# Gap Register And Change Roadmap

## 目的

這份文件記錄目前最重要的前後端整合缺口，並把它們轉成可執行的 change 候選。  
不是所有缺口都要立刻做，但應該先有排序，避免下一輪規劃時又重新盤點一次。

## 目前優先級原則

1. `/studio`：立即刪除 dedicated route 與 proven legacy-only graph members；不再建立 fallback
2. `/studio/v2`：唯一產品 setup 入口、重點施作
3. `/studio/runtime`：作為 V2 的 post-setup observer 一起補齊
4. `/test`：維持工程工具入口
5. `/gateway/*`：只保留記錄與 prototype 說明

## P0: 現在就該優先處理

| Change 建議名 | 問題 | 主要影響面 | 為什麼優先 |
| --- | --- | --- | --- |
| `retire-legacy-studio-and-polish-v2` | `/studio` dedicated surface 已獲 owner 授權立即刪除，且刪除後必須與 generic unknown route 等價 | router、legacy import/chunk/test、inventory、release evidence | 先移除已不使用的 legacy surface，避免持續被誤認為產品契約 |
| `integrate-studio-v2-multi-device-commit` | V2 commit 仍是前端模擬，且還不能一次送出多台設備到後端 | `/studio/v2`、後端 datalink contract | 不補這個，V2 仍只是草稿 shell |
| `adapt-runtime-handoff-for-multi-device` | runtime 目前主要以單台 `device_id` handoff 為前提 | `/studio/v2`、`/studio/runtime` | 不補這個，多台 commit 成功後無法穩定交接 |
| `normalize-runtime-lifecycle-status-contract` | runtime snapshot / SSE 缺正式 lifecycle semantics | `/studio/runtime`、未來 `/studio/v2` | 前端目前還要靠錯誤字串判斷狀態 |

## P1: 緊接著要處理

| Change 建議名 | 問題 | 主要影響面 | 說明 |
| --- | --- | --- | --- |
| `refactor-studio-mainline-surface-boundaries` | `/studio` 歷史 inventory 曾建議整理 page boundary | pre-delete historical record | 已被立即刪除決策 supersede，不應重新建立 legacy surface |
| `promote-source-rule-apply-contracts-in-studio` | tags / outputs 的 apply endpoints 曾由 legacy 主線使用 | backend capability / V2 contract | 只有在保留 surface 明確需要時另立 change，不以 `/studio` 為 owner |
| `document-surface-api-inventory` | 這次新增的維護文件需要正式維護機制 | docs / review 流程 | 否則文件很快 drift |

## P2: 視方向決定是否要做

| Change 建議名 | 問題 | 主要影響面 | 說明 |
| --- | --- | --- | --- |
| `align-gateway-pages-with-datalink-product-contract` | `/gateway/*` 是否要轉正仍未定義 | `/gateway/*` | 目前先不推進，只保留選項 |
| `promote-readiness-first-device-flow` | `/devices/:id/readiness` 已存在但 V2 主線未全面採用 | `/studio/v2` | 可讓使用者看到更明確的 staged readiness |
| `add-runtime-e2e-handoff-coverage` | setup -> runtime 的整條路徑仍缺完整 e2e | 前後端整合測試 | 補產品信心 |

## 建議缺口類型

| 類型 | 定義 | 例子 |
| --- | --- | --- |
| `missing-api` | 後端沒有頁面需要的 API | V2 正式多台設備 commit contract |
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

### 2. `integrate-studio-v2-multi-device-commit`

- 目標: V2 commit 後能一次送出多台設備，並拿到逐台 persisted 結果
- 後端工作:
  - 提供正式多台設備 commit endpoint 或 orchestration contract
- 前端工作:
  - Step 4 改用真實 commit response
  - 不再使用模擬 commit logs 當 truth
- 驗證:
  - `setup -> commit -> runtime handoff`

### 3. `adapt-runtime-handoff-for-multi-device`

- 目標: 讓 runtime 能接住單台或多台 commit 結果
- 前端工作:
  - 單台成功時可直接帶 `device_id`
  - 多台成功時需有正式 handoff context，不再假設只有單一 `device_id`
- 驗證:
  - 單台成功可直達 focused runtime
  - 多台成功可穩定落到可觀察的 runtime 入口

### 4. `normalize-runtime-lifecycle-status-contract`

- 目標: 將 runtime snapshot / stream 語意正式化
- 後端工作:
  - 對 `device not found`、`starting`、`not running yet` 提供明確語意
- 前端工作:
  - 移除字串比對
  - route state 改用正式 status enum
- 驗證:
  - missing device、starting、live、degraded 都能穩定呈現

### 5. `refactor-studio-mainline-surface-boundaries`

- 狀態: **superseded by immediate legacy deletion**
- 不再以 `/studio` 為新功能 owner；若保留的 V2/runtime surface 需要相同 capability，另以 typed contract 與現行 route 為範圍提出 change。

### 6. `align-gateway-pages-with-datalink-product-contract`

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
