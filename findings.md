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

## 2026-03-19 4.1 新發現
- `source_rule_links.tag_id / mapping_id` 不能只當預留欄位；一旦 rule-driven flow 改成自動建 Tag/Mapping，它們就必須成為 rule-derived relationship 的 authoritative anchor，不然 Step 3/4 仍會 drift。
- strict `1 Point : 1 Tag` 不能只靠前端 candidate filtering；至少要在 `mapping.Service.Create` 做 service-level gate，不然 manual path 或 race condition 仍會灌出重複關聯。
- rule delete / shrink 如果只靠 `points -> mappings ON DELETE CASCADE`，會留下 orphan tags；auto-generated tag 必須帶 rule-managed metadata，後續 cleanup 才能安全判斷哪些 tag 可以跟著移除。
- `tag` / `mapping` 的 not-found 判斷若只靠錯誤字串比對，很容易在 rollback / cleanup path 漂移；這類 lifecycle-sensitive domain 最好直接用 sentinel error + `errors.Is`。
- source-rule mutation 只 invalidate `points` 不夠；一旦 backend 自動建立 Tag/Mapping，前端 cache 也必須同步 invalidates `tags` / `mappings`，不然 Step 3 review surface 會短暫顯示舊狀態。

## 2026-03-19 Step1 / Step2 bug trace
- Step 1 的 `Test connection` 目前從 `WorkbenchInspectorPanel` 直接呼叫 `useTestConnectionMutation(deviceId)`，只測 **已存檔的 selectedDevice**，不會吃 `WorkbenchDeviceStep` inline editor 裡尚未儲存的 draft config。
- Step 2 的 capability context 也只看 `selectedDevice.connection_config`；如果使用者剛在 Step 1 改了 protocol / host / address-related config 但還沒 save，Step 2 仍會沿用舊設備上下文。
- `SourceCanvasSection` 的起始位址 state 目前直接 `useState('40001')`，沒有依協議切換預設基準；這會讓 FATEK / MC3E 之類的裝置看起來仍像 Modbus 規劃。
- 「尚未規畫前已有被規畫的點位」的高機率來源有兩種：
  1. 使用者其實仍停留在舊的 selectedDevice context（草稿未存，Step2 仍看舊設備）；
  2. 該 device 已有 persisted source rules，Step2 會依設計載回它們，但目前 UI 對「這是既有 persisted rule，不是你剛新增的草稿」說明還不夠強。
- 2026-03-19 實測 `192.168.31.62`：
  - `ping 192.168.31.62` 成功，代表 ICMP reachability 正常。
  - 直接從目前執行環境用 Python `socket.connect(('192.168.31.62', 502))` 會得到 `OSError: [Errno 65] No route to host`，與 workbench 回報一致。
  - `127.0.0.1:502` / `localhost:502` 在目前機器上則是 `Connection refused`，表示此刻本機沒有服務在 502 上 listening。
  - 因此至少目前這個錯誤不是前端捏造；更像目標主機 / 防火牆 / port bind 問題，或使用者所測的「本機可連」不是同一個 IP/port 組合。
- 2026-03-19 Step 1 診斷修正已落地：
  - `WorkbenchDeviceStep` 新增 draft-aware 測試入口，inline editor 可直接呼叫 `/datalink/devices/test-draft` 測目前草稿設定，不必先 save。
  - draft test 後端刻意不走 `ConnectionManager.GetOrCreate()` 快取，而是用一次性 `connector.Get()` + direct probe read，避免沿用舊的 saved-device 連線狀態。
  - `WorkbenchInspectorPanel` 與 `WorkbenchDeviceStep` 都新增 backend-host hint，明確說明 TCP dial / probe 是從 backend 所在主機發起；saved test 與 draft test 的設定來源也因此被分開講清楚。
- 2026-03-19 Step 2 實作決策已落地：
  - `startAddress` 不再只是 `SourceCanvasSection` 本地 state；每台設備的最後規劃起點會記在 `WorkbenchSourcePlanningState.plannerStartAddressByDeviceId`。
  - 切設備時，Step 2 會先恢復該設備上次的起點；若沒有記憶值，則透過 `getDefaultPlannerStartAddress(protocol)` fallback 到協議預設（目前 Modbus=`40001`、FATEK=`D0`、MC3E=`D0`）。
  - `clearSourcePlanningState()` 現在只清規則/選取，不會把 per-device 起點記憶一併抹掉。
- 2026-03-19 Step 2 主畫面已開始 Tag-first 化：
  - `AddressCanvas` / `AddressLedger` 主標題改為 `tagDisplayName -> tagKey -> point.name -> generic label`，讓主畫面優先講 Tag，但未綁 Tag 的既有點位仍可辨識。
  - source rule 卡片會顯示 `既有規則 / 草稿規則` badge，降低 persisted state 被誤認成新規劃的風險。
  - 主畫面文案（summary / actions / canvas / ledger）已回拉到中性 `rule/source` 語氣，避免在 UI 端過度宣稱「已直接建立 Tag」，但 Point 仍只保留在 inspector/debug 細節，符合使用者選擇的 `B`。
  - `addressParser.offset()` 已補 protocol-aware lower bound：Modbus 維持從 `1` 起算，FATEK / MC3E 改為允許 `D0`，避免規劃器把 `D0` 錯誤偏移成 `D1`。
  - review 第二輪已確認上述三項修正後沒有新的實質問題。

## 2026-03-19 Step 4 Local Modbus 新發現
- `LocalModbusBoard` 原本雖然看起來是 HR 介面，但實際上整個 surface 都直接暴露 internal 0-based register：
  - canvas slot label 直接顯示 `HR{i}`
  - register input 預設 `0`
  - inspector trace 顯示 `HR${modbusMapping.register}`
  - conflict / message / dry-run 也全都直接印 internal register
- 只把 label 改成 1-based 不夠；tag chips、inspector、dry-run、conflict 文案若沒一起改，會出現同畫面混用 `HR0` / `HR1` 的語意漂移。
- 直接移除 64-slot cap 也不安全：若使用者輸入極高位址（例如 UI 允許的 `65536`），canvas 會一次 render 六萬多個 button，造成瀏覽器卡死。
- 這輪最後採用的安全方案是：
  - backend / stored mapping 維持 internal 0-based，不動既有 API 契約
  - workbench surface 全面改成 1-based 顯示與輸入，再於 bind 時做轉換
  - canvas 改成 bounded viewport，並以目前選取 / 輸入 register 作為 anchor，因此 `HR200` 可見，但高位址也不會炸 DOM

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
