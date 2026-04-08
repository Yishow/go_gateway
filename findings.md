# Findings

## 2026-04-09 Phase 5 shared foundation 新發現
- Phase 5 的起點不是先畫新 shell，而是 **先鎖 ownership boundary**；如果沒有 shared compare contract，後續 `v2` / `v1` / `v3` 很容易把 shell 擴張成各自的 version-only workflow center。
- 目前 `phase5.scenarioFocus` 舊值 `readiness / global-blockers / repair-hop / return-to-flow` 太鬆，會讓 shell 看起來像能直接接管修復流程；對齊成 `readiness-summary / active-blocker / diagnostics-refresh / return-to-mainline` 後，shared contract 才真正符合 OpenSpec 的四個 shell owned surfaces。
- Phase 5 shared acceptance 必須明講 **shell 只做 summary / refresh / return，不做 step-local edit / mutation / validation**；不把這條寫死，variant UI 很容易在 compare 前就先漂成不同 ownership model。
- `v2` / `v1` / `v3` 在 Phase 5 的 archetype 差異應該建立在 summary / alert / blocker framing 上，而不是把 step-local workflow 拉上 shell；這讓後續 baseline 與三版 compare 能聚焦在 diagnostics 讀法，而不是誰偷接管更多操作。

## 2026-04-09 Reopened Phase 4 compare gate 新發現
- **`v2` 仍是 Output phase 的最佳 canonical owner**：它把 readiness、dry-run、apply、blocker diagnosis 都提升到 incident-desk command surface，對 Phase 4 這種「先判讀 target 狀態，再決定 publish / repair 動作」的節奏最完整。
- **`v1` 是最佳 long-session operator console**：calm rail 能把 rule revision、attention count、return path 與 blocker framing 壓成最低認知負擔；但實際 mutation 仍在 shared board，所以關鍵操作速度沒有超過 `v2`。
- **`v3` 的 telemetry density 與 quick-action dock 最強**：top KPI + bottom dock 讓 expert operator 一眼讀到全局並直接觸發真動作；代價是 wrapper/global-style/action-anchor orchestration touched 面最大，維護風險明顯高於 `v1` / baseline / `v2`。
- baseline 雖然 shared board 最直接、raw functionality 也最完整，但 compare 真正看的是 archetype-level reframing；baseline 仍讓 operator 自己從 shared page 解讀 target diagnosis，因此不適合作為 rollout owner。
- 四版在同一條 real API 路徑上的 raw latency 沒有 user-visible 差異；Phase 4 真正拉開差距的是：
  - 第一屏是否先講清楚「現在是哪個 target、哪種狀態、下一步該做什麼」
  - 從 landing 到 meaningful action / blocker 的 click path 長度
- Phase 4 正式 recommendation：
  - 維持 `v2` canonical owner
  - 優先吸收 `v1` 的 calm blocker / revision framing
  - 次優先吸收 `v3` 的 telemetry strip 與 bottom operations dock
  - 不更換 canonical owner，除非後續要回 OpenSpec 提 amendment

## 2026-04-08 Phase 4 baseline Current Studio 新發現
- `4.1.baseline` 一開始看起來像是 baseline Output 沒有 linked tags，可是真正的 shared blocker 是 **`/api/v1/datalink/mappings` 因 nullable lifecycle 欄位直接掃進 `string` 而回 500**；修正 `mapping/sql_repo.go` 的 NULL scan 後，real API 才重新露出 linked tag candidates，baseline `/studio` 才能做真正的 Output snapshot。
- 真實裝置 `UI 4.3 Modbus TCP` 現在可在同一個 `/studio?step=output` 裡同時切換 Local Modbus 與 Database，且共用同一組 tag chips（`TAG_40002` / `TAG_40001`）；這證明 baseline 已經具備跨 target family 的 shared selection surface。
- baseline 最明顯的語意衝突是：**global tag chips 已能直接操作 linked tag，但 rule-scoped review strip 同時又顯示 `來源規則 REVISION —`、`0 筆候選項` 與「請先聚焦來源規則」**。對 operator 來說，這會形成「上方已可動作、下方卻說目前沒有 candidate」的雙重訊號。
- Local Modbus focused state 的 traceability 不差：切換 `TAG_40001` 後，右側 inspector 會同步顯示來源位址、TAG key、Local Modbus / Database mapping 狀態；但 blocker narrative 仍分散在 strip、workspace controls 與 page-local status message，沒有 incident-style 單點診斷面。
- Local Modbus blocker 在真實操作下可直接再現：按下 `啟動 Server` 後，頁內會明確回報 `port 5020 is already in use, please stop the conflicting process or choose another port`。診斷內容夠具體，但仍停留在 page-local message，而不是更高位階的 blocker ownership。
- Database overview 雖然成功共用同一組 tag chips 與 selected-tag preview，但在 connector 尚未存在時就先露出 table / write-mode / write-row preview 骨架；這讓 operator 在尚未完成 connector readiness 前，就被迫看見過早的 mapping form。
- Database blocker 也能以真實操作直接重現：展開 connector editor 後直接 `儲存 Connector`，頁內會清楚回報 `請先輸入 connector 名稱`。copy 清楚，但仍屬 page-local validation，沒有被提升成更清晰的 target-diagnosis surface。
- 本輪正式 evidence（`agent-browser`，session artifacts）：
  - `baseline-output-modbus-overview.png`
  - `baseline-output-modbus-focused.png`
  - `baseline-output-modbus-blocker.png`
  - `baseline-output-database-overview.png`
  - `baseline-output-database-focused.png`
  - `baseline-output-database-blocker.png`

## 2026-04-08 Phase 4 shared foundation 新發現
- `4.1` 的缺口和前兩個 reopened phase 一樣，首先不是 UI，而是 **shared compare contract 不存在**；在 `main` 先補 `workbenchOutputCompareContract.ts`，才能把 Output phase 的 compare gate 從 OpenSpec 文字落到前端主契約。
- Output phase 的 shared scenario focus 必須明確從舊的 `mapping-visibility / apply-output` 收斂成 spec 指定的 `readiness / dry-run / apply / blocker-diagnosis`；否則後面 baseline / v1 / v2 / v3 很容易又各自定義不同的 Output 主線。
- 這輪 contract 也把一個重要 spec 約束寫死：**Local Modbus 與 Database 兩個 target families 都必須被同一份 compare contract 覆蓋**，避免 Phase 4 只拿單一 target family 做假 compare。

## 2026-04-08 Reopened Phase 3 compare gate 新發現
- 四個 compare 輸入（baseline / `v2` / `v1` / `v3`）在同一條 critical path 上的 raw latency 幾乎沒有差異：以 local preview + 同一個 real device `UI 4.3 Modbus TCP` 測得，`/studio -> 選 device -> Tag Review` 約都在 `2.9s` 同級；真正決定優劣的是資訊排序與 action zoning。
- **`v2` 仍是 Tag phase 的最佳 canonical owner**：它把 review readiness、retry / recovery、以及 Output handoff narrative 都放在 incident-desk command surface 的首屏，既能保留 shared diff/apply flow，又最符合「先判斷、再決策、再交接」的操作節奏。
- **`v1` 是最佳 long-session control room**：左 rail 讓 active rule、counts、return path 與 calm pacing 最清楚，操作心理負擔最低；但主要 command 仍散落在 rail / workboard / right-side tools 之間，所以 system completeness 與 action clarity 還是略輸 `v2`。
- **`v3` 的首屏資訊密度與產品辨識度最強**：cockpit banner、telemetry tiles、bottom handoff dock 讓 operator 一眼讀到全局；代價是 layout override 與 visual orchestration touched 面最大，shared 結構若再變，v3 最容易需要跟著補修。
- baseline 雖然共享 flow 完整，但 compare 真正看的是 archetype-level reframing；baseline 仍主要仰賴操作者自行從 dense shared page 推理 decision hierarchy，因此不適合作為後續 rollout owner。
- Phase 3 compare 正式結論：**維持 `v2` 作為 Tag canonical owner**；後續若要吸收亮點，優先回收 `v1` 的 calm rail pacing 與 `v3` 的 telemetry density，而不是更換 canonical line。

## 2026-04-08 Reopened Phase 3 v3 ClickHouse Data Cockpit 新發現
- `v3` Tag 最重要的不是再做一張 neon 皮，而是把 `TagBindingStudio` 改成 **top telemetry -> full-width review board -> bottom operations dock** 的 cockpit 讀法；這樣才能和 `v1` 的 left rail、`v2` 的 incident desk 真正拉開 archetype 差異。
- `CockpitTagGlobalStyles` 的核心價值在 orchestration，而不是 fork domain：它只重排 shared `TagBindingStudio` 的 section/aside hierarchy，讓 board 成為 primary instrument、exception/master surface 下沉成 operations dock，因此 shared review/apply/recovery contract 仍維持一致。
- `agent-browser` 實際驗收證明 v3 的 first-screen density 已經成立：真實裝置 `UI 4.3 Modbus TCP` 進入 Tag step 後，會先讀到 cockpit banner、telemetry tiles 與 bottom handoff dock，再進入 review board。
- 目前真實資料下 handoff 仍 blocked（linked tags = 0），但 v3 已把 blocker narrative 提升成 cockpit dock，而不是只在 context bar 留一顆 disabled CTA；這讓 compare 可以公平評估 v3 的 blocker legibility。
- 本輪正式 evidence：
  - `phase3-v3-tag-overview.png`
  - `phase3-v3-tag-focused.png`
  - `phase3-v3-tag-handoff-blocked.png`

## 2026-04-08 Reopened Phase 3 v1 Linear Control Room 新發現
- `v1` Tag 這輪最安全、也最符合 spec 的做法，是把 archetype 差異放在 **zoning / pacing**：左側 `TagRailPanel` 承接 active rule、candidate count、point blocker 與 shortest next action，中央仍保留 shared `TagBindingStudio`，外層 inspector 完全不 fork。
- 這種 layout-only control-room shell 對 shared Tag flow 的干擾極低；official `woe-v1-radix` targeted Tag suite `36/36` 全綠，證明 review queue、exception tools、tag master 與既有 apply/recovery contract 都還在原位運作。
- `agent-browser` 實際驗收顯示 v1 的首屏讀法已和 baseline / `v2` 拉開：操作者進入 Tag step 後，會先看到 calm rail 的 rule review / revision / point count framing，再進入 shared workboard，而不是直接被 dense candidate controls 淹沒。
- 真實裝置 `UI 4.3 Modbus TCP` 在目前資料下仍是 handoff blocked（linked tags = 0），但 v1 已把這個狀態收斂成 control-room framing，而不是只留一個漂浮在 shared page 的 disabled CTA。
- 本輪正式 evidence：
  - `phase3-v1-tag-overview.png`
  - `phase3-v1-tag-focused.png`
  - `phase3-v1-tag-handoff-blocked.png`

## 2026-04-08 Reopened Phase 3 v2 Tag incident desk 新發現
- `v2` 這輪真正成立的關鍵不是把 `TagBindingStudio` 換掉，而是把它降成 shared workboard，再由 `MuiTagIncidentDesk` 把 blocker / summary / handoff 提到首屏；一旦用 early-return 把 workboard 整塊藏掉，既有 Tag step integration contract 會立刻整片失效。
- `tagReady` 不能只靠候選列存在就成立；`useWorkbenchSummary()` 的 ready gate 仍看 linked tag/mapping，所以 incident-desk handoff 驗證必須用真實 linked state，不能只用「有 candidates」當 proxy。
- official `woe-v2-mui` 一度殘留了不可信 subagent 寫下的 composition 版本與 dead incident locale；這輪已收斂回 rescue 驗證過的 shell，並清掉過時 test / locale 殘留，避免 compare 時再次混入假 branch state。
- `agent-browser` 實際驗收證明 v2 首屏已經具備明確的 incident-desk 節奏：上方先給 review readiness / Output handoff narrative，下方仍保留 shared review queue、exception tools、tag master 與 batch diff preview。
- browser 驗收還抓到一個真實 i18n 漏洞：handoff blocker 一開始落回英文 fallback；補上 `workbench.tag.results.outputBlocked` 後，正式 evidence 才可算數。
- 本輪正式 evidence：
  - `phase3-v2-tag-overview.png`
  - `phase3-v2-tag-handoff-blocked.png`
  - `phase3-v2-tag-diff-preview.png`

## 2026-04-08 Reopened Phase 3 baseline Tag evidence 新發現
- baseline Tag 在同一個 `/studio` flow 下已能直接打開 review queue，但 **decision surfaces 分散**：候選列本身有 row-level actions，batch action 又在另一塊 checkbox toolbar，tag library 再是第三塊；對大量 review 來說，Hick's Law 成本偏高。
- baseline 的 handoff 問題很明顯：`Review outputs` 只顯示 disabled CTA，沒有把「為什麼還不能去 Output」提升成明確 blocker narrative；這會迫使操作者回頭自己猜是 selection、mapping 還是 apply 狀態造成的阻塞。
- baseline focused review 雖然能看見 candidate-level 資訊，但 diff / review focus 與 handoff 關係沒有被重新 framing，整體仍比較像 shared data page，而不是一個清楚的 Tag operation language。

## 2026-04-08 Reopened Phase 3 shared foundation 新發現
- Tag phase 的 shared contract 與 Source 類似，也適合先用獨立 compare contract 檔把 matrix 鎖住，再由 `workbenchExperimentContract.ts` 統一 re-export；這樣 variant branches 之後可以直接依同一份 phase3 scenario / acceptance / archetype / critical task 做實作與 compare。
- `3.1` 最重要的不是先做某個 variant UI，而是先把 `review candidates -> diff preview -> choose/apply -> failure/retry/recovery -> handoff Output` 這條 shared operator path 固定下來，否則後面三版很容易又各自漂出不同 Tag flow。
- Phase 3 design 對 archetype 差異已足夠明確：`v2` 是 incident-desk review command deck，`v1` 是 control-room batch review skeleton，`v3` 則只能保留 cockpit board / summary density，不能把 shared apply / recovery contract 取代掉。

## 2026-04-08 Reopened Phase 2 compare gate 新發現
- 這輪 compare 以同一台真實 device `UI 4.3 Modbus TCP` 的 `overview -> focused -> handoff` 路徑做 controller-side 判讀，四個輸入都維持同一個 `/studio` route、同一組 backend API 與 shared domain flow。
- **`v2` 仍是最佳 canonical owner**：它把 blocker / recovery / handoff 全都放在 incident-desk command surface 裡，對 Source 這種「先判斷、再編修、再交接」的操作節奏最順；compare 不足以推翻 OpenSpec 既定 owner。
- **`v1` 最適合長時間操作與低干擾巡檢**：左欄 rail 把 framing 做得最清楚，操作心理負擔最低；但因主要 command 仍在 canvas / outer inspector 間分配，所以在 Source phase 的 decisive command clarity 還是略輸 `v2`。
- **`v3` 的 first-screen density 最強，但 maintenance risk 也最高**：top KPI strip + cockpit shell 讓資訊一眼很多，但它仰賴額外 wrapper 與 global style orchestration 去重排 shared `SourceCanvasSection`，後續若 shared 結構再變，v3 最容易需要跟著補修。
- baseline 雖然功能沒壞、關鍵路徑也短，但 reopened compare 真正要拉開的是 archetype-level framing；baseline 在這一點上仍主要依賴操作者對 shared workbench 的既有記憶，因此不適合作為繼續 rollout 的 owner。
- 結論：Phase 2 compare 可以正式判定三個 variant 在非 Device phase **仍 clearly distinct**，且 `v2` 無需換線。

## 2026-04-08 Reopened Phase 2 v3 ClickHouse Data Cockpit completion 新發現
- 實際 code 與 browser 驗收顯示，v3 並不是只有 banner：`SourceCanvasSection` 已經把 ClickHouse cockpit 語言深度接到 Source step，包含 top banner、dense rule layer、coverage readout 與 cockpit-styled canvas/ledger framing。
- v3 與 v2 的差異點不在 shared domain flow，而在 **首屏讀法**：v2 是 command-deck / handoff-first，v3 則把 KPI / rule layer / coverage readout 提到 top-of-fold，讓操作員先看 cockpit instruments，再進入細部編修。
- controller-side `agent-browser` 驗收已確認 v3 overview / focused / handoff 三態都能在同一組真實 API 下成立，證據為：
  - `phase2-v3-source-overview.png`
  - `phase2-v3-source-focused.png`
  - `phase2-v3-source-handoff.png`
- 就 reopened `2.1.v3` 的 minimum-obligation compare 目標來說，目前 official `woe-v3-antd` branch state 已足夠作為 compare 輸入，不需要再做第二條 v3 rescue 線。

## 2026-04-08 Reopened Phase 2 v1 Linear Control Room completion 新發現
- v1 `Source` 這輪終於不再只是 header wrapper：`SourceControlRoom` 改成真正的 inner three-zone skeleton，左欄 `SourceRailPanel` 讓規則導覽、blocker 與 handoff 從 canvas 中抽離，中央仍保留 shared `SourceCanvasSection`，右側則沿用 `WorkbenchFrame inspector` slot 承接現有 inspector banner。
- 這個做法把 archetype-level 差異放在 **action zoning / work framing**，而不是去 fork shared Source domain logic；因此 `SourceCanvasSection` 不需要再被 v1 branch 深改，shared API 與 state ownership 也維持一致。
- `SourceRailPanel` 目前是 read-only navigator，而不是第二套 command surface；這讓 v1 能維持 long-session friendly 的 calm pacing，同時不破壞 canonical `v2` already-defined flow ownership。
- controller-side `agent-browser` evidence 已保存：
  - `phase2-v1-source-overview.png`
  - `phase2-v1-source-focused.png`
  - `phase2-v1-source-handoff.png`

## 2026-04-08 Reopened Phase 2 v2 rescue completion 新發現
- official `v2` 雖然已有 `Sentry Incident Desk` 外觀，但 reopened acceptance 真正的缺口在三個 shared-contract 層面：`phase2.scenarioFocus` 漂移、`workbenchSourceCompareContract.ts` 沒有 consumer、以及 stale preview recovery 沒有可見 UI surface。
- rescue 修正的關鍵不是再做一層 MUI wrapper，而是把 `workbenchExperimentContract.ts` 改成 live-link Source compare contract，並讓 stale template auto-upgrade 真的經過 `source-template-warning` surfaced state。
- 補做 review-driven regression 後確認：mount-time `templateRecoveryWarning` 也不能永久壓過後續 template apply warning；因此最終顯示優先序必須是 `templateWarning ?? templateRecoveryWarning`，否則 legacy template 使用者會看不到 capability mismatch 提示。
- `agent-browser` 的有效 v2 evidence 必須來自正確指向 backend `8080` 的 preview；若 4175 上跑的是未帶 `VITE_API_BASE_URL=http://127.0.0.1:8080/api/v1` 的 dev server，UI create/select device 會落到 local 404，這種 run 不能算正式 acceptance evidence。
- 實際 handoff 驗證最穩定的 real device 是 `UI 4.3 Modbus TCP`：它已存在 `4 pts`，所以能直接驗證 `Source ready -> Tag Review` 的真正可行路徑；本輪 v2 三態 evidence 已保存為：
  - `phase2-v2-rescue-source-overview.png`
  - `phase2-v2-rescue-source-focused.png`
  - `phase2-v2-rescue-source-handoff.png`

## 2026-04-08 Reopened Phase 2 baseline evidence 新發現
- baseline `Source` overview 仍從 device-pick 起手；即使直接 deep-link 到 `?step=source`，第一個最明顯的動作還是先決定裝置，而不是進入某種 archetype-specific planning workspace。
- 選到 `Browser Smoke PLC` 後，focused state 立即把 `Rule builder`、`Plan / Live / Link`、`Review tags` CTA 與 dense grid 一次鋪開；操作是完整的，但 hierarchy 仍然偏向 shared workbench，而不是某個鮮明的 phase-specific product language。
- 點 `40004` 的 focused state 會把 cell 選取壓在同一張 dense grid 裡，沒有額外的 work/readout scaffold 幫操作員理解「現在正在處理哪個 Source 任務」，這正是 reopened compare 要求重新拉開的地方。
- `Review tags` handoff 目前可正確導到 `http://127.0.0.1:4173/studio?step=tag`，但交接語意仍主要靠 step switch 本身，而不是由 `Source` surface 內部提供更強的 handoff framing。
- 本輪 baseline evidence 已保存到 session files：
  - `phase2-reopen-baseline-source-overview.png`
  - `phase2-reopen-baseline-source-focused.png`
  - `phase2-reopen-baseline-source-handoff.png`

## 2026-04-08 Reopened Phase 2 shared contract 新發現
- `workbenchExperimentContract.ts` 原本只有 Phase 0/1 matrix 與 Device 1R contract；這輪補上正式的 Source reopened shared contract export，避免 reopened `2.1` 只存在於 OpenSpec 文字、沒有進入前端主 contract。
- 新增 `workbenchSourceCompareContract.ts` 後，`Source` phase 現在有四個正式 shared artifacts：
  - `WORKBENCH_SOURCE_COMPARE_SCENARIOS`
  - `WORKBENCH_SOURCE_COMPARE_ACCEPTANCE`
  - `WORKBENCH_SOURCE_COMPARE_ARCHETYPES`
  - `WORKBENCH_SOURCE_COMPARE_CRITICAL_TASK`
- `WORKBENCH_EXPERIMENT_PHASES` 的 `phase2.scenarioFocus` 也已對齊 reopened gate：`create-rule / apply-template / plan-live-link / stale-preview-recovery / handoff-tag`。
- 這代表 reopened `2.1 共用基礎` 不再只是規格口號，而是已被主程式與單元測試正式鎖住。

## 2026-04-08 Phase 2–5 reopen governance 新發現
- 使用者實際看過 `v1` / `v2` / `v3` live preview 後，直接判定「除了 Device，其他都像同一版」；這個回饋與目前程式碼現況相符：非 Device phase 雖然有 branch-specific 檔案修改，但 operator-facing 差異仍然不夠強。
- 使用者已明確選擇：**reopen Phase 2–5**，並且 **保留 winner-led owner model**（`v2` canonical、`v1` full-flow high-polish、`v3` minimum-obligation）。
- 新治理的核心不是改 shared logic，而是要求所有 non-Device phases 都必須重做 operator surface：action placement / ordering、primary work surface、preview / summary framing、visual language 都要真正分化。
- OpenSpec amendment 已完成且 `openspec validate --changes workbench-ux-operator-efficiency` 通過；`openspec instructions apply` 目前回到 `24/69`，下一個正式 pending task 是 `2.1 共用基礎`。

## 2026-04-08 Phase 3 execution-governance 新發現
- 目前 `workbench-ux-operator-efficiency` 的 `proposal.md`、`design.md`、`tasks.md` 都一致採用 **winner-led rollout**：Phase 2–5 固定節奏是 `shared -> baseline -> v2 -> v1 -> v3 -> compare`，且 `v2` 是 canonical owner。
- 你最新要求改成 `baseline -> v1 -> v2 -> v3 -> compare`，且要求「Spec is Law」；因此目前不是單純執行順序偏好，而是 **使用者指令與現行 OpenSpec 工件衝突**，不能直接假裝兩者都成立。
- `openspec instructions apply --change workbench-ux-operator-efficiency --json` 目前顯示進度為 `33/68`；已完成到 `3.1.v2`，剩餘最近的 pending tasks 是 `3.1.v1`、`3.1.v3`、`3.1.compare`。
- 正式 worktree 已存在且可直接沿用：baseline `woe-base-current-ui`、`v1` `woe-v1-radix`、`v2` `woe-v2-mui`、`v3` `woe-v3-antd`；`.worktrees/` 也已被 git ignore。

## 2026-04-08 Phase 3 v2 Tag 新發現
- `v2` 已把 Tag 首屏改成 Sentry Incident Desk 語言：command deck 先給 active rule / readiness / Output handoff，再把 review surface 與 batch board 收進同一個 command-center 節奏，首屏產品辨識度比 baseline 明顯高。
- canonical `v2` 這輪已吸收 shared contract：refresh failure 會明確 surfaced feedback + retry，不再讓錯誤只停在 raw candidate error 區塊。
- 真正有效的 browser evidence 來自修正後的 `VITE_API_PROXY_TARGET=http://127.0.0.1:8080`；先前誤指向 `3333` 的 run 已作廢，不納入 acceptance。
- `woe-v2-mui` 最終 head 為 `ec0e0c1`，`3.1.v2` 已在 branch tasks 打勾，並保存 `phase3-v2-tag-{overview,focused,handoff}.png`。

## 2026-04-08 Phase 3 baseline evidence 新發現
- `main` 已補齊 Tag shared acceptance / scenario matrix：Phase 3 compare focus 改成 `review-candidates / diff-preview / apply-decision / failure-retry-recovery / handoff-output`，`SourceRuleTagReviewSurface` 也補上 refresh failure feedback + retry CTA。
- `workbench.tag.reviewSurface` 原本缺整段 locale；這輪已把 rule-scoped review copy、actions、feedback、stale state 一次補齊，避免 baseline Tag review surface 在真實 UI 出現 raw translation keys。
- baseline Tag overview 雖然已經有 rule-scoped review queue，但 rename / skip / override controls 會在首屏一次展開多列，review surface 與 batch board 彼此競爭注意力，decision hierarchy 偏吵。
- baseline diff preview 必須先手動 `Select bindable only` 才會浮現；preview scaffold 不是主動引導，而更像藏在 batch board 後面的 secondary readout。
- baseline handoff 到 Output 仍卡在 blocker：即使完成真實 bind，Step 4 仍只給出 `Finish tag review first so Step 4 has linked tags to publish.`，操作員很難立即理解還缺哪個 review / apply state。
- 本輪 agent-browser evidence 已保存到 session files：`phase3-baseline-tag-overview.png`、`phase3-baseline-tag-diff-preview.png`、`phase3-baseline-tag-handoff-blocked.png`。

## 2026-04-08 Phase 2 Source compare 新發現
- baseline 共享 flow 完整，但首屏仍直接把 device-picker 與 planner/canvas 疊在同一張工作檯，operator 必須自己理解 handoff；local compare flow 到 Source planner ready 約 `0.99s`。
- v1 `Linear Control Room` 的長時間可讀性最好，`SourceControlRoom` + runtime/tag-review surface 讓 focused / handoff 狀態更 calm，但需要額外 scoped CSS layer 與多筆收尾 commit；local compare flow 約 `1.02s`。
- v2 `Sentry Incident Desk` 在不改 shared flow 的前提下，把 overview / focused / handoff 三態收進同一個 command deck，首屏資訊密度、引導順序與 canonical 行為契約最平衡。
- v3 `ClickHouse Data Cockpit` 的產品語言辨識度與資訊密度最強，但 neon cockpit 噪音最高、shared surface touched 面最大，維護風險高於 v1 / v2。
- Compare 結論：維持 `v2` 作為 Source phase canonical owner；後續如需加值，優先收割 `v1` 的 calm control-room pacing 與 `v3` 的 density cues，不需要回頭修 OpenSpec owner。

## 2026-04-08 Phase 2 baseline evidence 新發現
- baseline Source overview 仍以 device-picker + rule-builder 雙區塊起手；選到 `Browser Smoke PLC` 後，首屏立即暴露大量欄位與稀疏 canvas，重複輸入負擔高。
- baseline focused state 的 `Review tags` 仍是單一前進動作；即使操作員已聚焦到 `40004`，切到 Tag step 後仍直接落到整個 device 範圍的候選清單，handoff 焦點沒有被保住。
- baseline `Link` mode 雖然能把 cell 文案切成 downstream 狀態（例如 `Unbound` / `No downstream state`），但 mode switch 只是在同一張密集 canvas 上改字，沒有額外的 compare / preview scaffold 幫操作員理解差異。
- 本輪 agent-browser evidence 已保存到 session files：`phase2-baseline-source-overview.png`、`phase2-baseline-source-focused.png`、`phase2-baseline-source-selected.png`、`phase2-baseline-source-handoff.png`、`phase2-baseline-source-link-mode.png`。

## 2026-04-08 Source -> Tag boundary shared-contract 新發現
- Tag step 原本的 batch diff preview 只看 selected device 全部 points，沒有跟著 active source rule / candidate snapshot 收斂 scope，因此會把 unrelated device-scope candidates 帶進同一個 diff。
- 這輪已把 active-rule selection 抽成 shared helper，Tag board 透過 candidate snapshot 收斂 point scope；因此 focused-rule handoff、rule 切換 stale preview invalidation、candidate snapshot 變更 stale preview invalidation，現在都由同一個 shared contract 驅動。
- 為了不吞掉既有的 review surface，Tag board 的整頁 empty gate 仍維持 `points.length === 0`；當 active rule 的 candidate set 為空時，只清空 board candidates，不會把 review surface 一起隱藏。

## 2026-04-08 Phase 2 shared contract TDD 新發現
- Source step 的第一個真正 shared-contract 缺口不是 planner helper，而是 `useSourceRulesQuery` 還在 loading / error 時，畫面仍直接進入一般 planner/canvas。
- 這輪已用 TDD 補上 `loading / loadFailed / retry` 三個 surfaced state，並新增 `SourceCanvasStatusState.tsx` 承接狀態 UI，避免繼續放大 `SourceCanvasSection.tsx` 這種歷史超長檔。
- Vitest 只收 `frontend/tests/unit/**`，因此新增 Source status-state 測試時要補 wrapper `frontend/tests/unit/pages/datalink/workbench-source-status-state.test.tsx` 去匯入 `src/.../__tests__/SourceCanvasStatusState.test.tsx`。
- 為符合超長檔規則，`SourceCanvasSection.tsx` 最終比 HEAD 少 1 行、`DatalinkWorkbenchSourceStep.test.tsx` 維持與 HEAD 同行數；新行數增量由新小檔承接。

## 2026-04-08 Phase 2 shared contract reconnaissance
- Source shared contract 目前主要落在三個地方：`frontend/src/pages/datalink/workbench/__tests__/DatalinkWorkbenchSourceStep.test.tsx`、`frontend/src/features/datalink/sourcePlannerContract.ts`、`frontend/src/pages/datalink/workbench/SourceCanvasSection.tsx`。
- `sourcePlannerContract.test.ts` 目前只鎖住 naming prefix、count parsing、template round-trip 等 planner helper contract，尚未覆蓋這輪 spec 新增的 `diff preview scope`、stale preview invalidation、Source -> Tag handoff 邊界。
- `SourceCanvasSection.tsx` 已明確持有 `viewMode`（`plan` / `live` / `link`）與 template load/save state，代表這輪 shared contract 更可能要先在 Source step integration test 上補 failing test，而不是只加 helper unit test。

## 2026-04-08 Workbench UX OpenSpec correction 新發現
- `proposal.md`、`design.md`、`tasks.md` 已對齊 approved winner-led rollout spec：`v2` canonical、`v1` full-flow high-polish、`v3` minimum-obligation compare track。
- `tasks.md` 的 Phase 2–5 順序已從舊的 `baseline -> v1 -> v2 -> v3 -> compare` 改為 `shared -> baseline -> v2 -> v1 -> v3 -> compare`，並補上 baseline 定義、Phase 4 雙 target families、Phase 5 shell ownership 邊界。
- `openspec validate --changes workbench-ux-operator-efficiency` 已通過；這代表下一個真正可以開始的實作項是 Phase 2 shared contract，而不是再回頭修 spec artifact。

## 2026-04-07 Workbench UX Phase 2–5 rollout planning 新發現
- `docs/superpowers/specs/2026-04-07-workbench-phase-2-to-5-rollout-design.md` 已完成 reviewer loop，且兩個 reviewer 都明確轉為 `✅ Approved`；代表後續可進入 implementation planning，但仍要先修正 OpenSpec change files。
- Phase 2–5 的新硬規則已固定：`v2` 是 canonical behavior track、`v1` 是 full-flow high-polish track、`v3` 是 minimum-obligation comparison track；compare 不能直接默默改寫 canonical owner，若要改 owner 必須先回到 OpenSpec amendment。
- Phase 2–5 的 baseline 已被明確定義為：`main` branch 上、該 phase shared acceptance commit 完成後、任何 variant branch UI work 開始前的 `/studio` frozen control snapshot。
- Phase 4 Output 不再能抽象地當成單一表面處理；規劃與驗收都必須同時覆蓋 `Local Modbus register binding` 與 `Database schema/column binding` 兩個 target families。
- 目前 `openspec/changes/workbench-ux-operator-efficiency/tasks.md` 仍停留在舊的三等份 Phase 2–5 實驗矩陣；真正開始 Phase 2 前，必須先把 proposal/design/tasks 一起改成 winner-led rollout model。

## 2026-04-07 Workbench UX Phase 0 啟動新發現
- `openspec status/apply` 顯示目前 change 為唯一 active change，schema 為 `spec-driven`，進度 `4/60`；可直接從 Phase 0 開始，不需再猜測使用哪個 change。
- 目前 repo 已有舊的 `.worktrees/workbench-v1|v2|v3`，但它們仍停在舊基線 commit，且命名不符合本次 OpenSpec matrix；最安全做法是保留不動，另外建立 `woe-*` 正式 worktree。
- `frontend/src/App.tsx` 已明確把 `/studio` 定為唯一主產品 route，`/datalink/*` 收斂為 redirect；這代表本 change 的 route guard 主要是「不可破壞既有 contract」，不是新增路由。
- `frontend/package.json` 目前只有 Radix 相關依賴，尚未引入 MUI / Ant Design；因此 v2、v3 Phase 0 需要由各自 branch 明確新增依賴與 theme mapping，不可假設 repo 已預裝。
- 現有 `frontend/src/styles/tokens.ts` 與 `designSystem.ts` 偏向目前 dark dashboard token；本輪若要做公平實驗，必須另建一層 kit-agnostic semantic token contract，而不是直接把現有 blue/slate token 當成共享實驗主題。
- `frontend-design` 的 UX / color guidance 與本案目標一致：工業操作員屬 B2B / efficiency-first 情境，應優先控制 cognitive load、資訊分塊、低疲勞 dark surface 與清楚的 severity / diagnostics，而不是追求炫光視效。
- `openspec validate` 對 validate 子命令使用 `--changes`，不是 `--change`；apply/status 與 validate 的 CLI 旗標不一致，後續 automation 需注意。

## 2026-04-07 Workbench UX Phase 0 compare 新發現
- 四套表面在 compare 當下都能以同一個 `/studio` route 啟動，且 `?step=source` deep-link 也都可正常落到來源步驟；這代表 route contract 與 route-state sync 在四版共同起跑線上成立。
- 透過四個 frontend dev server 的 `/api/v1/datalink/devices` proxy 回應做 hash，比對結果完全一致；因此本輪 compare 的 backend 條件確實是同一組真實 API，而不是版本分流。
- `v1` 的 shadcn/Radix foundation 最保守，幾乎維持 baseline 的資訊結構，只在 shell 語彙與 token 掛載上建立版本身份；優點是風險低，缺點是第一眼辨識度與 step-state 語意提升有限。
- `v2` 的 MUI foundation 在 Phase 0 就已經把 step rail、context bar、summary bar 與 theme shell 做出最明確的版本身份；a11y snapshot 中只有 v2 明確呈現 selected tab 語意，這對後續 compare 的「邏輯清晰度」很有利。
- `v3` 的 Ant Design foundation 版本身份也很清楚，且 Layout/Steps 能快速形成工作台感；但目前從 a11y snapshot 看，步驟選取狀態不像 v2 那麼直接，且 App-level ConfigProvider 帶來的全域影響面比 v1 大。
- 以 Phase 0 這個「共同起跑線」目標來看，推薦版本先給 `v2 / MUI`：不是因為它已經是最終 winner，而是它在不改 API/route 的前提下，最成功地建立了可比較、可辨識、可延續到後續 phase 的 shell 基礎。
- 風險排序（僅 Phase 0）：`baseline` 最低、`v1` 次低、`v3` 中高、`v2` 最高；v2 的清晰度最好，但也是目前 diff 面與維護面最大的版本，後續 phase 需持續觀察是否會放大成本。

## 2026-04-07 Workbench UX Phase 1 Device shared criteria 新發現
- Device phase 若不先鎖定 `create / edit / clone / connect / probe / diagnostics` 六個動作，很容易讓各版本只優化自己擅長的一段，失去 compare 的可比性。
- baseline 的 `agent-browser` 實際探查顯示：
  - Device list 首屏有明確 `Create device`
  - 已選設備 detail 可見 `Edit device`、`Clone from selected`、`Test connection`
  - `Edit device` 之後可在編輯器中看到 `Test current draft settings`
  - 執行 draft diagnostics 後，頁面文字明確出現 `Connect phase` 與 `Protocol probe`
- 換言之，baseline 並不是沒有 probe，而是 **probe state 被包在 diagnostics 結果裡，而不是作為首屏就能辨識的一級動作**。
- 這代表 Phase 1 compare 的一個核心問題不是單純視覺好不好看，而是：**版本能否把 connect / probe 分段診斷語意做得更清楚，同時不破壞現有 create/edit/clone 路徑。**
- `workbench-device-step-editor.test.tsx` 與 `workbench-device-form-model.test.ts` 已經把多協議欄位、draft preservation 與 clone / edit 基本約束固定下來；Phase 1 variants 應該重用這些現有 contract，而不是另外發明新流程。

## 2026-04-07 Workbench UX Phase 1 v2 / MUI 新發現
- v2 把 Device monolith 拆成 6 個 MUI 元件後，connect / probe / detail / editor 的責任邊界明顯清楚許多，這對 compare 中的「邏輯清晰度」與「維護風險」都很有利。
- `MuiDeviceDiagnosticsPanel` 把 `connect` 與 `probe` 變成一級 stage，這正面回應了 baseline 只有在 diagnostics 結果裡才看得到 probe state 的問題。
- v2 的真實瀏覽器驗證顯示：Device 相關 API 正常打到 backend 8080，而 `/mappings` 的 500 是 pre-existing backend issue，不應誤記到 v2 Device surface 頭上。
- v2 保留 create / edit / clone 主線，同時讓 mode（CREATE / EDIT）在 editor 內變成可視語意，這可能會在「關鍵操作時間」與「操作順暢度」上勝過 baseline。

## 2026-04-07 Workbench UX Phase 1 v1 / shadcn-Radix 新發現
- v1 也把 connect / probe diagnostics 升成明確分層，不再像 baseline 那樣只在診斷結果裡被動顯示 probe state。
- v1 的操作時間非常接近 baseline：
  - list -> detail：`459ms`
  - detail -> editor：`397ms`
  - editor -> diagnostics：`381ms`
- 這代表 v1 的改良方向偏向 **保守演進**：它把 diagnostics 語意拉高，但沒有像 v2 那樣大幅重塑整個 Device shell。
- 風險面上，v1 雖然拆成多個子元件，但主協調檔 `WorkbenchDeviceStep.tsx` 仍有 462 行；若 compare 重視維護面，這會明顯輸給 v2。

## 2026-04-07 Workbench UX Phase 1 v3 / Ant Design 新發現
- v3 把 Device shell 的資訊密度拉得比 v1 / v2 都高，特別是在 list/detail/card/drawer 的結構語意上更接近「控制台」感。
- v3 也明確呈現 connect / probe 兩階段 diagnostics，但互動時間比 baseline / v1 / v2 都慢一些：
  - list -> detail：`486ms`
  - detail -> editor：`467ms`
  - editor -> diagnostics：`462ms`
- v3 的主要風險不是功能缺失，而是維護面：
  - `DeviceFormDrawer.tsx` 422 行，超過 300 行軟上限
  - foundation tests 28 failures 屬於 Phase 0 Ant Design shell 既有問題，會拖累版本整體穩定性評價
- 因此 v3 在 compare 中很可能會同時拿到：
  - 較強的首屏資訊密度 / workbench identity
  - 較高的實作與維護風險

## 2026-04-07 Workbench UX Phase 1 compare 結論
- **推薦版本：v2 / MUI**
- 排名理由：
  - `v2`：在 Device 場景下，connect / probe 兩階段 diagnostics 的語意最清楚，editor mode 也最一眼可辨識；雖 detail -> editor 比 baseline / v1 稍慢，但整體認知負擔最低，且主要檔案都壓在 300 行內。
  - `v1`：是最穩健的低風險備案。它保留 baseline 節奏，又把 diagnostics 語意前移；但主協調檔仍有 462 行，維護面不如 v2。
  - `v3`：資訊密度與控制台感最強，但互動時間最慢，且 `DeviceFormDrawer.tsx` 422 行、foundation tests 28 failures（Phase 0 shell 既有問題）都會放大後續維護成本。
  - `baseline`：速度與風險都最保守，但 Device probe 語意仍然過於隱藏，不適合作為 Phase 1 最終推薦。
- Device phase 關鍵操作時間（ms）：
  - baseline：`458 / 381 / 377`
  - v1：`459 / 397 / 381`
  - v2：`502 / 493 / 386`
  - v3：`486 / 467 / 462`
  - 順序分別為：`list -> detail / detail -> editor / editor -> diagnostics`

## 2026-04-07 Workbench UX 使用者直接 review 新發現
- 使用者明確判定上一輪 Phase 1 Device round **不合格**，核心原因不是單一 bug，而是設計方向錯了：
  - 三版太像彼此，也太像 baseline
  - UI kit 被當成「換皮」而不是不同 interaction model
  - `design-md` 的骨架 / 資訊區 / 數據視覺語言沒有真正落地
- 因此後續 Phase 1 不應再沿用「同一個骨架換 kit」策略，而要改為：
  - `v1` = Linear Control Room
  - `v2` = Sentry Incident Desk
  - `v3` = ClickHouse Data Cockpit
- 本輪最重要的驗收標準不再只是 connect / probe 有沒有變清楚，而是：**三版是否一眼就能看出不同 interaction model 與不同視覺語言。**

## 2026-04-07 Workbench UX Phase 1R shared token contract 新發現
- 原始 Phase 0 token contract 雖然足夠支撐「同一骨架換 kit」的公平比較，但不足以支撐 Phase 1R 的三 archetypes reboot：
  - `Linear` 需要近黑 + 單一 indigo-violet + 超細白邊界
  - `Sentry` 需要 warm purple-black + Rubik/uppercase + lime/coral emphasis
  - `ClickHouse` 需要 pure black + neon yellow-green + heavy KPI display
- v2 在重做途中開始長出 branch-local `sentryVisualTokens.ts`，這不是單純命名問題，而是明確訊號：**shared contract 太窄，版本只好自己補 style contract。**
- 依使用者規則，遇到 token conflict 時主代理不能讓版本偷偷完成，而必須先回到 shared contract；因此本輪先暫停 `phase1r-v1/v2/v3` 收尾，重開 `phase1r-shared`。
- 修補方向不是為三版各開一份 token file，而是在同一份 `workbench-experiment-tokens.ts` 補上：
  - `archetype.linear / sentry / clickhouse`
  - shared typography cues（510 / 700 / 900 權重與 label tracking）
  - shared treatment semantics（glass / inset / neon glow / ambient purple）
- 這樣三版仍共用**同一份 shared semantic contract**，但已不再被單一灰黑 palette 綁死。

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

## 2026-04-06 Phase -1 API 實作新發現
- SourceRule 的 `database_outputs` candidate 在現行模型下，若沒有已存在的 DB mapping scope（connector/table/column），API 層無法安全自動補建 mapping；因此 `apply` 需回傳 per-item `schema_missing`，不能默默成功。
- Local Modbus output snapshot 目前仍是 `deferred`（payload 可為空），`/local-modbus/apply` 在此階段應回傳 per-item `skipped` + `deferred`，避免前端誤判成後端故障。
- DB tooling 三支 API 若以「直接套用既有 mapping repo + inspectTables」為核心，可在不引入 mock 的前提下快速提供：
  - schema generate（preview/execute）
  - mapping dry-run（candidate-level blocked reason）
  - write history（由 writer 寫入事件聚合）
- 在 sandbox 內，`internal/api/handlers` 全量測試會因既有測試需 bind TCP 埠而失敗；本輪採 targeted tests 驗證新增 API 契約，並保留此環境限制說明。
