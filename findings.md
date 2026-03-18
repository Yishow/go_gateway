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
