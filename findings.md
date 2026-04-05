# Findings

## 2026-04-06 未提交變更 code review 新發現
- `scripts/check_file_lines.sh` 在本地 fallback 模式原本只看 staged 或 unstaged diff，未涵蓋 untracked 新檔；這會導致開發者在 `git add` 前先跑 `make check-lines` 時漏檢新建立的大檔案。
- 修正方式：
  - fallback 模式改為合併 `git diff --cached`、`git diff`、`git ls-files --others --exclude-standard`
  - 再 `sort -u` 去重後統一檢查
- 修正後驗證：
  - `bash scripts/check_file_lines.sh`
  - `make check-lines`
  - 兩者皆能正確涵蓋本次變更檔案並通過門檻檢查（warning 僅反映 >300 行檔案，不阻擋）。

## 2026-04-06 檔案行數規範強制落地新發現
- repo 內已存在多個 >500 行歷史檔案（前端 workbench 頁、部分 backend service、測試檔、大型 docs/lock 檔）；若直接做「全倉 hard fail」，會立即阻擋幾乎所有實務修改。
- 可行的強制策略是「只檢查本次變更檔案」並加入 legacy guard：
  - 新檔或本次修改後 >500 行：阻擋
  - 歷史 >500 行檔案若本次修改後行數不增加：允許通過（並要求後續逐步縮減）
- lock files、build artifacts、嵌入式靜態資產與外部匯入 docs 必須有 ignore 清單，否則行數規範會被非程式碼類型檔案干擾。
- 只靠文件宣告不足以形成約束，需同時落地：
  - `scripts/check_file_lines.sh`（單一規則實作）
  - CI workflow（伺服器端強制）
  - pre-commit hook（本地提早阻擋）
  - PR 模板（人類審查時補充 >300 行理由與拆分計畫）
- 「不能只讀一份規範文件」必須在 `AGENTS.md` 與 `CLAUDE.md` 雙向聲明，才能避免 agent 僅讀自身專屬文件（例如只讀 `CLAUDE.md`）而遺漏共通規範。

## 2026-04-06 AGENTS / CLAUDE / README 規範整併新發現
- `README.md` 已完成完整化，但 `AGENTS.md` 與 `CLAUDE.md` 仍偏向「綜合敘述」，缺少獨立、可快速查閱的章節（尤其是 `Error Handling Pattern` 與 `禁止事項`）。
- `AGENTS.md` 雖有「程式風格與命名慣例」，但命名規則與錯誤處理規則混在描述中，不利於 code review 或 onboarding 時快速對照。
- `CLAUDE.md` 原本以工作流程與脈絡為主，對「提交前測試要求」與「禁止事項」缺少明確條列，容易產生執行邊界模糊。
- `Makefile` 實際支援 `gatev11`、`points-precheck-down`、`points-migrate-down`、`longtask-smoke`，若規範文件未列出，會讓維運與驗證流程被低估。
- 三份文件一致化的關鍵不是內容完全重複，而是：
  - `AGENTS.md` 作為共通規範主體（完整、可操作）
  - `CLAUDE.md` 補 agent 視角下的落地邊界
  - `README.md` 提供入口層摘要，且欄位名稱能對齊前兩者

## 2026-04-06 README 規範整併新發現
- 現有 `README.md` 過於精簡，只覆蓋產品入口說明，缺少實際開發/維運會依賴的規範資訊（命令、樣式、測試、安全、禁止事項、OpenSpec 流程）。
- `AGENTS.md` 與 `CLAUDE.md` 在規範內容上已高度對齊，但 `README.md` 尚未承接這份對齊成果，導致新成員無法從入口文件一次建立正確心智模型。
- 實際可執行命令需要以 `Makefile`、`frontend/package.json`、`scripts/build.ps1` 為準，而不是只列最常見命令：
  - Makefile 另含 `gate*`、`points-*`、`longtask-smoke` 等工作流命令
  - 前端另含 `test:gateway:*` 任務
- 產品入口與路由收斂已在實作層明確落地（`frontend/src/App.tsx`）：
  - `/studio` 為主線
  - `/test` 為測試工具入口
  - `/datalink/*` 屬 compat redirect 收斂
- 靜態資源供應與單一可執行檔模型在程式碼層清楚可驗證（`cmd/test_ui/main.go` + `internal/web/embed.go`），README 應明確寫出 embed 與 SPA 路由處理機制，避免誤解部署型態。
- `.github/instructions/go.instructions.md` 明確要求 error handling pattern（`%w` 包裝、`errors.Is/As`、錯誤訊息風格），原 README 缺漏，這是開發一致性風險點。
- `openspec/project.md` 仍有部分歷史技術棧描述（如 viper/gorm）與目前 `go.mod` 不完全一致，README 應以 repo 現況（`go.mod` 與實際 import）為主，避免引用歷史描述造成偏差。

## 2026-03-23 AGENTS / CLAUDE 文件對齊新發現
- `AGENTS.md` 原本已涵蓋結構、測試、UI 主線與文件工作流，但缺少獨立的「安全考量」區塊，無法完整承接 repo 對輸入驗證、secret 管理、參數化查詢與 `gosec` 的要求。
- `CLAUDE.md` 與 `AGENTS.md` 原本在規範優先順序上存在描述差異；本輪已收斂為同一套規則：`AGENTS.md -> Agent 專屬文件 -> .github/instructions/`。
- repo 的建置與驗證命令不只 `make build` / `go test` / `npm run test`，還包含 `make gen-docs`、`go vet ./...`、前端 `test:e2e` / `test:gateway:*` 與 Makefile 裡的 gate / migration scripts，文件應視情況明確列出。
- `/studio`、`/test`、`SQLite + PostgreSQL` 與 `cmd/test_ui/static` embed 流程，已是 repo 現況的一部分；若文件只描述舊的高層概念，容易和實際操作脫節。
- 本次文件對齊屬 active docs 更新，不涉及其他使用者正在修改的程式碼面。

## 2026-03-20 `/test` 精簡改造新發現
- `frontend/src/App.tsx` 目前仍把 `/test`、`/templates`、`/history`、`/compare`、`/analyzer` 一起掛在舊 `Layout` 下，代表「只保留測試頁」至少會涉及 legacy 測試工具 routes 清理。
- `frontend/src/components/Layout.tsx` 的側邊欄與 `max-w-7xl mx-auto` 是 `/test` 現在看起來像多頁後台、且在 1920 螢幕下內容偏窄的直接來源。
- `TestPage.tsx` 本身已經是完整的單頁工具集合，若改成 page-owned shell，大多數精簡需求不需要動它的核心操作流程，只需處理容器與入口層。
- 目前「移除用不到的頁面」最明顯的候選是 `/templates`、`/history`、`/compare`、`/analyzer`，但是否連同對應 page 檔與測試一起刪除，仍需先向使用者鎖定範圍。
- 使用者已確認：
  - 清理範圍不只 `/test`，還包含舊 datalink 頁。
  - datalink 主入口應改為 `/datalink` 直接進 `/datalink/workbench`。
  - 最終主線命名不保留 `datalink`，改採 `/studio`。
- `frontend/src/features/datalink/legacyRoutes.ts` 已經提供多條 compat redirect helper，代表這輪更適合採「入口切換 + legacy redirect 收斂」而不是只刪檔不處理舊連結。
- `SmartDashboard.tsx` 目前只是 `SmartDashboardPage.tsx` 的薄 wrapper；若要盤掉舊 datalink 頁，實際待處理的重量級 legacy surface 仍是 `SmartDashboardPage.tsx` 與 `LocalModbusWorkbenchPage.tsx`。
- `LocalModbusWorkbenchPage.test.tsx` 仍在測舊 page 行為，若移除 legacy page，測試也要一起轉向 redirect 或 compat contract。
- 這輪最安全的落地方式不是暴力刪掉整個 SmartDashboard implementation tree，而是：
  - 對外主入口先改成 `/studio`
  - 舊 `datalink` 路由改走 redirect
  - 明確無用的 test-tool pages 與舊 local modbus page 再實體刪除
- `buildWorkbenchRedirect()` 改成直接產生 `/studio` 後，原本從 SmartDashboard 內部導向 workbench 的操作也會自然落到新主線，不需要額外再補一層 route glue。
- 2026-03-20 續查 `openspec/` 時發現：
  - `openspec/specs/` 底下其實有多份歷史 `TBD - created by archiving change ...` Purpose placeholder。
  - 但和這輪 archive 直接新增/變動強相關、且最適合立即收尾的是：
    - `openspec/specs/database-target-workbench/spec.md`
    - `openspec/specs/source-rule-runtime/spec.md`
  - 因此這輪 follow-up 採 **最小收尾**：只補這兩份的 Purpose，不順手擴大整理整個 openspec 舊債。
- 2026-03-20 下一輪 legacy cleanup 盤查結果：
  - `SmartDashboard.tsx` 只是 `SmartDashboardPage.tsx` wrapper。
  - `SmartDashboardPage.tsx` 仍直接依賴 `frontend/src/pages/datalink/smart-dashboard/` 整個子樹。
  - 已確認的 SmartDashboard page-level tests 至少包含：
    - `SmartDashboard.interaction.test.tsx`
    - `SmartDashboardGridOverlaysSection.test.tsx`
    - `useSmartDashboardWorkspaceContentState.test.tsx`
    - `useSmartDashboardWorkspaceState.test.ts`
    - `useSmartDashboardCommitFlow.test.ts`
    - `useSmartDashboardPanelsState.test.ts`
  - `frontend/tests/integration/ui/smart-dashboard-regression.test.tsx` 仍直接 import `@/pages/datalink/SmartDashboard`，若 repo 層完整移除 legacy UI，這支 integration test 也應一起移除。
  - `frontend/src/styles/dashboard.ts` 初步搜尋無任何引用，傾向視為 orphan 一併刪除。
- 這輪 docs 邊界已由使用者鎖定為：**只更新 active docs，保留 historical docs / archived specs**。
- `useSmartDashboardShortcuts` 僅剩 SmartDashboard 舊頁與其專屬測試使用；在刪除 legacy page 後，保留它只會留下無主 API，因此應連同 export 與測試一起收掉。
- `frontend/FILE_CLASSIFICATION.md` 與 `frontend/tests/README.md` 屬於 active docs，若不一起更新，repo 說明會與實際檔案狀態衝突。

## 2026-04-06 staged code review 新發現
- `internal/datalink/db.go` 將內嵌 SQLite DSN 切到 WAL 後，repo root 會額外產生 `datalink.db-wal` / `datalink.db-shm`；若 `.gitignore` 不同步補上，每次啟動 `cmd/test_ui` 都會污染工作樹。
- 本輪 staged diff 在 `frontend/src/utils/addressParser.ts` 新增 MQTT topic-based 位址解析與 `sensor/data` 預設起點，但現有 source planner / canvas / source rule backend 仍是 sequential-address 模型：
  - frontend `buildPlannedPointAddresses()` 與 `buildAddressCanvasItems()` 會用 offset / expand 推導連續位址
  - backend `sourcerule.Service` 會依規則建立衍生 points / links
  - backend point address validation 仍不接受 MQTT topic（例如 `/`）格式
- 結論：這不是單點 parser bug，而是 staged 變更提前暴露了尚未真正打通的 MQTT source planner contract；最安全的修補是先回收 topic-based planner support，而不是讓 UI 接受 topic 後在 runtime / point create 才失敗。

## 2026-04-06 uncommitted review 新發現（本輪）
- `frontend/src/pages/datalink/workbench/sourceCanvasModel.ts` 的 `formatSourceValue()` 雖然新增了 object / array / JSON payload 支援，但 numeric string 走到 `string` 分支時會直接回傳原字串：
  - `point.last_value` 在前端型別上是 `unknown`，實際上常以字串形式出現
  - `SourceCanvasSection` 的 value format toolbar 允許 `decimal / hex / binary / float`
  - 結果是同一筆數值若以字串傳入，切到 `hex` / `binary` / `float` 會失去格式化，形成 UI regression
- 最小且正確的修補不是回退新 formatter，而是在字串分支先做 `Number(trimmed)` 正規化；這樣：
  - 純數字字串可延續既有格式模式
  - JSON 物件內的 `value: "255"` 也能透過遞迴套用相同邏輯
  - 非數值字串仍保持原樣，不會誤傷一般文字 payload

## 核心結論
- 原始需求始終沒有改變：datalink UI 要回到單純主線，而不是讓使用者在 SmartDashboard、Tag、Local Modbus、資料庫之間切頁與切心智模型。
- 最適合的實作路徑仍是 **混合式過渡**：新 workbench 承接主線，舊頁只做 fallback / compat。
- 真正該重用的是 domain 與 hooks / services / types，不是舊 UI 外觀本身。

## 2026-03-20 剩餘 canonical OpenSpec staged 變更盤點
- 目前尚未提交的 5 份 canonical spec 不是純格式調整，而是補入一整組與這輪 workbench / runtime 收斂相符的 requirement：
  - `datalink-workbench-desktop`
    - 補 Step 1 `connect` / `probe` 分段診斷
    - 補 Step 2 device capability / persisted rule state 語意
    - 補 Step 3 review-first / exception-handling requirement
  - `local-modbus-memory-workbench`
    - 補 per-target isolated selection state
    - 補單一 authoritative binding state model
    - 補 bind / unbind inline feedback requirement
  - `point-catalog`
    - 補 rule-derived point 是 primary runtime asset
    - 補 unmanaged / legacy point 仍需可辨識
    - 補 point collection 跟隨 rule lifecycle 啟停
  - `protocol-connectors`
    - 補 connector test 的 connect-stage / probe-stage 分段診斷
    - 補 protocol-specific probe configuration requirement
  - `tag-dictionary`
    - 把 Tag/Point cardinality 收斂成 `1 Point : 1 Tag`
    - 補 source rule 自動建立 Tag + Mapping requirement
- 結論：這 5 份不是單純補 `Purpose`，而是把最近一輪 workbench / runtime / mapping 決策正式寫回 canonical spec。

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

## 2026-03-19 Step 2 / Step 3 新發現
- Step 2 的 draft rule state 若不帶 `deviceId`，只靠 `selectedDeviceId` 切畫面，切設備時很容易把 device-local 草稿誤當成全域草稿清掉。
- `clearSourcePlanningState()` 若同時負責清 selection 與清 rules，很容易在「切設備」和「重置目前畫布」兩種語意之間漂移；這次證實兩者必須拆開，至少 rules 要能跨 device 切換保留。
- Step 3 現有單筆 unbind 雖然存在，但只要使用者一次選到多筆已綁定 row，主 CTA 就會停在 blocked/conflicts，形成實際上的 rebind dead-end。
- 對這種 dead-end，解法不是放寬 `bind` gate，而是補一條 selection-aware escape hatch：讓使用者能直接對目前 selection 做 batch unbind，再回到正常 bind 流。
- unbind 成功摘要若仍沿用 `created / linked` 語彙，會讓 Step 3 lifecycle 語意變混亂；因此結果摘要至少要能區分 `bind` 與 `unbind` 兩種 action mode。
- Step 3 若沒有一個 visible tag master surface，使用者會被迫在「綁定流程」中順手做資料管理，結果就是找不到全域 Tag、也無法先整理舊 Tag 再回來綁定。
- 這次驗證後比較安全的 `Tag master` 邊界是：
  - 允許快速建立 standalone Tag
  - 允許刪除未使用 Tag
  - 對已綁定 Tag 先顯示 `使用中` 並停用刪除，避免在 review flow 直接拆壞既有 mapping
- batch unbind 這種多步 mutation 不能只依賴 React Query 單一 mutation 的 `isPending`；若沒有本地 in-flight guard，按鈕會在迴圈間短暫重新啟用，造成重入。

## 2026-03-19 Step 4 / Database 新發現
- Step 4 的 target-isolated state 不能只做「各 target 各自記住 selected tag」；若保留 Step 3 → Step 4 的 `focusedTagIds` 交接，handoff 必須高於既有 selection，否則回到 Step 4 時會看起來像焦點沒有接上。
- 但 handoff 只該在 `focusedTagIds` 改變時覆寫 selection；一旦進到 Step 4 內部，`modbus` / `database` 仍必須能各自保留後續手動改選，否則會重新回到 cross-target drift。
- `DatabaseTargetBoard` 原本真正的 drift 來源不是單一 bug，而是 `selectedConnectorId / tableKey / columnName / writeMode / timestampColumn` 五段 state 分散在多個 effect 裡互相修正，connector 切換時很容易短暫殘留上一個 schema 的 table/column。
- 把 database output state 收斂成明確的 `DatabaseOutputScope`（connector / table / column / writeMode / timestamp）後，才比較容易保證：
  - connector 切換時 downstream scope 一起 reset
  - table 切換時 column / timestamp 會跟著重新正規化
  - `upsert` / `insert` 切換不會留下失效的 timestamp column
- 這次也證實：Database flow 若要做 `Connector / Schema / Mapping` 分層，最小安全做法不是先重做 UI，而是先把 scope model 顯性化，再讓 UI 反映該 scope。

## 2026-03-20 Step 3 / Output 收尾新發現
- Step 3 若要真正符合「建立規則後自動建立 Tag + Mapping」，不能只留原功能再換標題；必須把主畫面資訊排序改成：
  - 先看 generated / needs review
  - 再把 create / existing / unbind 放進 exception handling
  這樣操作員才不會誤解為「還要再手動綁一次才算完成」。
- `5.3` 這類 output feedback 類需求，不一定需要再重寫 UI；若現有行為已正確，補上 bind / unbind / delete mapping 的 inline feedback regression，反而是更安全的完成方式。
- `6.2` 的「migration and UX coverage」可以由三條 seam 組成：
  - Step 2 unmanaged legacy point 顯示與 inspector 說明
  - Step 4 per-target output selection drift contract
  - Database layered connector/schema/mapping scope regression
  三者一起成立，才足以證明新 flow 不會把 legacy state、output state、database scope 混在一起。

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
