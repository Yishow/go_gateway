# Findings

## 核心結論
- 原始需求始終沒有改變：datalink UI 要回到單純主線，而不是讓使用者在 SmartDashboard、Tag、Local Modbus、資料庫之間切頁與切心智模型。
- 最適合的實作路徑仍是 **混合式過渡**：新 workbench 承接主線，舊頁只做 fallback / compat。
- 真正該重用的是 domain 與 hooks / services / types，不是舊 UI 外觀本身。

## 仍有效的重要發現

### 架構層
- `SmartDashboardPage.tsx` 是高耦合 orchestration god component；不適合作為後續主線的長期基地。
- `device -> point -> tag -> output` 這條資料模型仍成立，尤其是 Step 4 目前仍綁在 Tag 上，不是直接綁 Point。
- `runtime` 與 `database target` phase 2 契約已完成，不是當前主阻塞。

### UI / UX 層
- Step 1 最怕的是「列表、篩選、編輯表單一起搶主位」；master-detail + inline editor 比 modal 更接近使用者回饋。
- Step 2 的主角必須一直是格狀 lattice，不是工具列牆；能降權的動作都應降到 supporting controls。
- Step 2 point create 若被全域 conflict 阻擋，使用者會直接感覺成「UI 不能用」；per-span safety 比 all-or-nothing gate 更合理。
- Step 3 的額外 ceremony 很容易讓使用者質疑「規則都已命名，為什麼還要再綁一次」；因此 Step 3 只能保留真正必要的 lifecycle / exception handling，不該再做成厚重中間層。
- Step 4 的核心訴求已被驗證：使用者真正要的是「先選好 tag，再直接點輸出表面綁定」，而不是再看一塊大型候選清單。

### 測試與驗證層
- 真桌面 smoke 很重要，因為 deep-link route-lock 這種問題不容易只靠靜態 code reading 看出來。
- `cmd/test_ui` 走 embedded static，驗最新前端時不能直接相信它送出的資產；要用 Vite dev server + `VITE_API_PROXY_TARGET` 做新 UI 驗證。
- 規劃檔若不定期收斂，很快會從「可續作記憶體」退化成「歷史堆積」。

## 目前未解 / 待處理議題
- round 2 polish 已技術收尾，剩下的是產品決策：
  - 是否把 `/datalink/workbench` 升格為 datalink 主入口。
  - 何時把最新前端重新嵌入 `cmd/test_ui/static`。
- Step 3/4 雖已完成需求對齊，但仍值得透過實機驗收再確認是否還有局部操作阻力。

## 2026-03-18 OpenSpec 前置分析新發現
- 使用者的重點已從單純 UI polish 轉成 **資料模型與流程語意**：
  - Step1 需要把 `connect` 與 `probe` 分流顯示與判定。
  - Step2 的 rule 不只是暫存規劃，而是要進 DB、可啟停、重啟還原。
  - Step3 傾向改成「建立規則後自動建立 Tag + Mapping，再由 UI 做覆核」。
  - Database 頁這輪要連 `SQLite + PostgreSQL` connector scope 一起定。
- 本輪選項分析曾出現 agent 狀態漂移：
  - SQL 顯示 `step1/database/step4` analysis 仍為 `in_progress`
  - 但 background agents 已不存在
  - 後續需要以新一輪 agent 補跑分析，再做總整合

## 2026-03-18 各 step 可優化處（目前已回收）
- Step1：
  - 最值得優先的是 **connect / probe 分流** + **允許 connect 成功但 probe 失敗時先存設備、但禁止啟動規則/採集**
  - 純粹只加 timeout 已不足，核心其實是狀態機與診斷粒度
- Step4：
  - 最值得優先的是 **每個 output target 分離 selection 狀態** + **用單一 source of truth 管 register/tag/binding state**
  - 問題核心是 drift，不是缺更多按鈕
- Database：
  - 最值得優先的是 **Connector / Schema / Mapping 分層**
  - 後續 drift check、versioning、智慧預設都應建立在分層之上

## 2026-03-19 SourceRule / Step 2 新發現
- Step 2 若把所有 existing point 都顯示成同一種 `used` 狀態，操作員無法分辨：
  - 這是 persisted rule 已落地的 span
  - 還是只有 point、沒有 rule 的 unmanaged legacy 狀態
- 因此 Step 2 狀態至少要拆成：
  - `planned`：rule 已存在但 point 尚未落地
  - `used`：rule-backed point 已落地
  - `unmanaged`：只有 point、尚未納入 persisted rule
  - `conflict`：rule / point / merge semantics 不一致
- `Create rule points` 與 `Create selected points` 的語意必須分開：
  - `Create rule points` = persist source rules（必要時帶 `skipped_addresses`）
  - `Create selected points` = manual / unmanaged exception path
- runtime live values 與 restart restoration 只有在 Step 2 直接吃 persisted rule + derived point 狀態時，grid / inspector 才不會和 backend lifecycle drift。

## round 2 已確認有效的收斂方向
- Step 1：editor 進中央區，不再用 modal。
- Step 2：
  - 32/64-bit merge 必須視覺正確。
  - secondary tools 要降權，但不代表要藏到找不到。
  - 刪 rule 與已建立 point 的一致性必須被處理。
- Step 3：bound item 必須可以取消 / 刪除。
- Step 4：binding action 要回到 output surface 本身，而不是先去操作大型 tag 候選區。

## 精簡歷史歸檔

### 2026-03-15
- 完成 datalink UI 深度分析。
- 確認新 workbench 主線：來源設定 -> 可視化 -> Tag -> 輸出。

### 2026-03-16
- 完成桌面 workbench shell / source / tag / output 重整基礎。
- quiet desktop / scan-first 成為後續 polish 核心原則。
- database target / runtime 契約落地。

### 2026-03-17
- 完成 master-detail polish、route-lock regression fix、1920 desktop smoke。
- 使用者再追加 round 2 回饋，焦點轉到 Step 1/2/3/4 的真實操作阻力。
- round 2 已完成：
  - Step 1 inline editor
  - Step 2 workflow consistency
  - Step 3 unbind/lifecycle 降噪
  - Step 4 direct surface binding
  - regression / review / build 驗證
