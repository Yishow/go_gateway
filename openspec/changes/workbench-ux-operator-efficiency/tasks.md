# Tasks: Workbench UX Operator Efficiency Experiment

> 對應 proposal: `openspec/changes/workbench-ux-operator-efficiency/proposal.md`
> 對應 design: `openspec/changes/workbench-ux-operator-efficiency/design.md`

## 共同規則

- baseline 與三個新版本都必須使用**同一個 `/studio` route**；版本隔離只靠 worktree / branch / port。
- baseline 與三個新版本都必須使用**同一組真實 API 契約**；禁止 mock-only 流程。
- 三個新版本都必須先共用**同一組 design-token 主題語意**，再映射到各自 UI kit。
- Phase 2–5 的 shared logic 保持不變：同一個 domain flow、同一組 blocker / retry / recovery semantics、同一組 state ownership。
- Phase 2–5 真的要重做的是 operator surface：action placement / ordering、primary work surface、preview / summary framing、視覺語言。
- 只有換 kit、加 scoped CSS、或只調整密度，不足以算完成；如果 review 仍認為「除了 Device，其他都像同一版」，該 phase 直接視為未完成。
- Phase 2–5 角色固定為：
  - `v2` = primary functional track
  - `v1` = full-flow high-polish track
  - `v3` = necessary-consistency comparison track
- compare 可以建議吸收 `v1` / `v3` 亮點或停下來修 spec，但**不可默默改寫 `v2` canonical owner**；若要換 owner，必須先回到 OpenSpec amendment。
- Phase 2–5 的 baseline 定義固定為：`main` branch 上、該 phase shared acceptance commit 完成後、任何 variant UI work 開始前的 `/studio` frozen snapshot。
- Phase 2–5 每個 phase 都必須完成 `共用基礎 -> baseline -> v2 -> v1 -> v3 -> compare`，才能前進。
- 每個 shared task 都必須先在 `main` 上定義 branch-neutral acceptance / scenario matrix，再開始 variant work。
- Phase 4 Output 一律同時覆蓋：
  - `Local Modbus register binding`
  - `Database schema/column binding`
- Phase 5 shell 只擁有：
  - readiness summary
  - active blocker summary
  - diagnostics refresh status
  - shortest return-to-mainline action
- step-local edit / mutation / validation 仍屬各步驟表面，不可被 shell 吃掉。
- 每個 `compare` 子任務都必須輸出：
  - 操作順暢度
  - 邏輯清晰度
  - 對系統的完整性
  - 首屏資訊密度
  - 關鍵操作時間
  - 實作 / 維護風險
  - 推薦版本與理由

---

## Phase -1：共享 API 基線（已完成）

> 目標：先完成 baseline + 三個新版本共用的真實 API，避免 UI 實驗階段各自定義契約。

- [x] -1.1 補齊 OpenSpec API 契約（`specs/datalink-api/spec.md`）
  - SourceRule output apply：database / local-modbus
  - DB schema generate（dry_run）、mapping dry-run、write history

- [x] -1.2 後端路由與 handler 實作
  - `POST /api/v1/datalink/source-rules/:id/database-outputs/apply`
  - `POST /api/v1/datalink/source-rules/:id/local-modbus/apply`
  - `POST /api/v1/datalink/db-targets/connectors/:id/schema/generate`
  - `POST /api/v1/datalink/db-targets/connectors/:id/mappings/dry-run`
  - `GET /api/v1/datalink/db-targets/connectors/:id/write-history`

- [x] -1.3 補後端測試（handler + service）
  - revision mismatch / conflict / connector invalid / schema missing
  - per-item partial success response contract

- [x] -1.4 將 Phase -1 作為 baseline + 三個新版本的共同起點

---

## Phase 0：baseline、tokens 與工作樹隔離

> 目標：建立四套表面的共同比較框架與 shared token 主題。

- [x] 0.1 共用基礎
  - 固定 baseline / v1 / v2 / v3 的版本矩陣
  - 固定 branch / worktree / port 命名
  - 固定 compare 輸出格式與量測 rubric
  - 固定共同驗收場景（Device / Source / Tag / Output）
  - 明確記錄「主線不受影響、共用 `/studio`、禁止 mock」
  - 明確記錄 phase 執行順序：`baseline -> v1 -> v2 -> v3 -> compare`

- [x] 0.2 共用 tokens
  - 從 `/Users/yishow/prj/awesome-design-md/design-md/` 萃取主題元素
  - 以 `linear.app + sentry + clickhouse` 作為目前選定的 shared token 參考混合
  - 區分 `preserve / adapt / reject` 三種抽取結果
  - 建立共享 design-token 語意層
  - 定義三個 UI kit 的 token mapping 邊界

- [x] 0.3 命名與執行順序落地
  - `baseline`：`woe-base-current-ui` / `.worktrees/woe-base-current-ui` / `4173`
  - `v1`：`woe-v1-radix` / `.worktrees/woe-v1-radix` / `4174`
  - `v2`：`woe-v2-mui` / `.worktrees/woe-v2-mui` / `4175`
  - `v3`：`woe-v3-antd` / `.worktrees/woe-v3-antd` / `4176`
  - 記錄 shared backend 使用策略與 compare 順序
  - 補出可直接照做的 worktree / backend / frontend 啟動操作稿

- [x] 0.4 語意 token 命名落地
  - 至少明確 `surface.*`、`text.*`、`accent.*`、`status.*`、`border.*`、`focus.*`、`density.*`、`data.*`
  - 確保三個 UI kit 共用相同 token family，而不是各自發明不同命名

- [x] 0.2.baseline 建立 baseline worktree
  - 從共享 API 基線切出 current UI baseline
  - 凍結 baseline 的 `/studio` 行為供後續比較
  - 收集 baseline 的初始截圖與量測方式

- [x] 0.2.v1 建立 `v1` worktree（shadcn/Radix）
  - 從共享 API 基線切出 `v1`
  - 在該 worktree 內讓 `/studio` 指向 shadcn/Radix 版本
  - 使用 shared design tokens

- [x] 0.2.v2 建立 `v2` worktree（MUI）
  - 從共享 API 基線切出 `v2`
  - 在該 worktree 內讓 `/studio` 指向 MUI 版本
  - 使用 shared design tokens

- [x] 0.2.v3 建立 `v3` worktree（Ant Design）
  - 從共享 API 基線切出 `v3`
  - 在該 worktree 內讓 `/studio` 指向 Ant Design 版本
  - 使用 shared design tokens

- [x] 0.2.compare 驗證四套表面共同起跑線
  - baseline 與三個新版本皆可在各自 port 開啟 `/studio`
  - baseline 與三個新版本皆連到同一組真實 API
  - 三個新版本皆使用相同語意 token
  - 四套表面皆可載入相同資料

---

## Phase 1：Device 對照實驗

> 目標：比較設備建立、編輯、clone、connect / probe 的最佳表面。

- [x] 1.1 共用基礎
  - 固定 Device 場景與驗收標準
  - 固定需要比較的動作：create、edit、clone、connect、probe、錯誤診斷

- [x] 1.1.baseline Current Studio
  - 量測目前 UI 在 Device 場景的操作順暢度、邏輯清晰度與主線承接完整性
  - 記錄 baseline 的優勢與缺點

- [x] 1.1.v1 shadcn/Radix
  - 以 shared tokens + shadcn/Radix primitives 設計 Device surface
  - 可依 kit 特性微調表單、結果區與提示區邏輯

- [x] 1.1.v2 MUI
  - 以 shared tokens + MUI components 設計 Device surface
  - 可依 kit 特性微調表單、結果區與提示區邏輯

- [x] 1.1.v3 Ant Design
  - 以 shared tokens + Ant Design components 設計 Device surface
  - 可依 kit 特性微調表單、結果區與提示區邏輯

- [x] 1.1.compare Device 比較與推薦
  - 輸出 baseline + 三個新版本在 Device 場景的比較表
  - 給出本 phase 推薦版本與理由
  - 已被直接 review 判定為「三版差異不足且過度接近 baseline」，以下重開 Device round

---

## Phase 1R：Device 三 archetypes 重做實驗

> 目標：把 Device 從「同一表面換 kit」提升為三種真正不同的工作台 archetype。

- [x] 1R.1 共用基礎
  - 把 `v1` / `v2` / `v3` 明確定義為三種不同 archetype
  - 固定三版各自的 design source mapping
  - 固定 Device 共通 domain model：selection / task / diagnostics / system
  - 明確規定 connect / probe 必須是第一公民，create / edit / clone 必須有明確 mode
  - 若 shared token contract 不足以表達三版語言，必須先擴充同一份 shared semantic contract；不得新增 version-only token file

- [x] 1R.1.v1 Linear Control Room
  - 以 `linear.app` 為骨架，重做 v1 Device interaction skeleton
  - 形成左欄 device rail / 中央 task canvas / 右欄 live inspector
  - 視覺上必須一眼區別於 baseline 與其他版本

- [x] 1R.1.v2 Sentry Incident Desk
  - 以 `sentry` 為骨架，重做 v2 Device interaction skeleton
  - 形成任務導向 command center，讓 diagnostics 成為主敘事
  - 視覺上必須一眼區別於 baseline 與其他版本

- [x] 1R.1.v3 ClickHouse Data Cockpit
  - 以 `clickhouse` 為骨架，重做 v3 Device interaction skeleton
  - 形成高密度 cockpit，首屏就呈現 KPI / health / recent tests / diagnostics
  - 視覺上必須一眼區別於 baseline 與其他版本

- [x] 1R.1.compare Device 重做比較與推薦
  - 重新輸出 baseline + 三個重做版本在 Device 場景的比較表
  - compare 必須明確說明三版在 interaction model 與視覺語言上的差異
  - 推薦 `v2 / Sentry Incident Desk` 作為本輪最佳平衡；`v1 / Linear Control Room` 保留作為低風險 fallback，`v3 / ClickHouse Data Cockpit` 保留高密度 cockpit 亮點供後續吸收
  - 給出新的 Device phase 推薦版本與理由

---

## Phase 2：Source winner-led rollout（reopened）

> 目標：在 **邏輯不變** 的前提下，重新讓 `Source` 形成三個真正不同的 archetype-level 操作表面，而不是 shared flow 外面包不同皮。

- [x] 2.1 共用基礎
  - 在 `main` 上更新 Source shared acceptance / scenario matrix
  - 固定 Source 共同動作：rule create/edit、template apply、plan/live/link switch、handoff to Tag
  - 明確定義 `diff preview scope`、stale preview invalidation、failure / retry / recovery
  - 明確定義什麼叫做 `Source` phase 的 archetype-level 差異（不是只換 shell / CSS）

- [x] 2.1.baseline Current Studio
  - 以 `main` branch baseline snapshot 收集 overview / focused / handoff evidence
  - 記錄 baseline 在重複輸入、模式切換與 stale preview 上的痛點

- [x] 2.1.v2 Sentry Incident Desk（primary functional track）
  - 完成 canonical Source flow 與 blocker / retry / recovery surfaced state
  - 以 incident-desk command surface 重做 `Source` 的 primary work / preview / handoff arrangement
  - 以 `v2` 作為後續 `v1` / `v3` 的行為基準，但不可只停在既有 wrapper-level 差異

- [x] 2.1.v1 Linear Control Room（full-flow high-polish track）
  - 保持與 `v2` 相同的 Source domain flow
  - 以 clearly different control-room editing skeleton 重做 `Source` surface
  - 操作感受必須一眼區別於 `v2` / `v3`

- [x] 2.1.v3 ClickHouse Data Cockpit（minimum-obligation track）
  - 只實作 compare 所需的 Source cockpit surface
  - 保留 overview / focused / handoff / blocker / retry evidence，不擴張成第二條主線
  - 但 primary work/readout surface 必須仍然是 clearly different cockpit form，不可退化成 shared page + dense skin

- [x] 2.1.compare Source 比較與推薦
  - 比較 baseline + `v2` + `v1` + `v3`
  - 若 compare 認為 canonical owner 必須改變，先停下來修 OpenSpec，不可默默換線
  - 若 review 仍認為三版除了 Device 幾乎一樣，本 phase 不得關閉

---

## Phase 3：Tag winner-led rollout（reopened）

> 目標：在 **邏輯不變** 的前提下，重新讓 `Tag` 形成三個真正不同的 archetype-level review/apply 操作表面。

- [x] 3.1 共用基礎
  - 在 `main` 上更新 Tag shared acceptance / scenario matrix
  - 固定共同動作：review candidates、preview diff、choose action、apply、failure / retry / recovery
  - 明確定義 `Tag` phase 的 archetype-level 差異與 fail gate

- [x] 3.1.baseline Current Studio
  - 以 baseline snapshot 收集 review queue、diff preview、handoff to Output evidence
  - 記錄 baseline 在 batch decision 與失敗回饋上的痛點

- [x] 3.1.v2 Sentry Incident Desk（primary functional track）
  - 完成 canonical Tag review/apply flow
  - 讓 blocker、retry、partial-failure surfaced state 在 `v2` 先穩定
  - 以 clearly different incident-desk review command surface 重做主要 decision flow

- [x] 3.1.v1 Linear Control Room（full-flow high-polish track）
  - 保持與 `v2` 相同的 Tag domain flow
  - 以 clearly different control-room review skeleton 提升大批量 decision 的 legibility

- [x] 3.1.v3 ClickHouse Data Cockpit（minimum-obligation track）
  - 只保留 compare 所需的 board / summary density
  - 不得改變 shared apply / recovery contract
  - 但 board / summary / apply surface 必須仍然 clearly different，而不是 shared Tag page 加密度樣式

- [x] 3.1.compare Tag 比較與推薦
  - 比較 baseline + `v2` + `v1` + `v3`
  - compare 只能推薦吸收亮點或回補 spec，不可直接改 canonical owner
  - 若 review 仍認為三版除了 Device 幾乎一樣，本 phase 不得關閉

---

## Phase 4：Output winner-led rollout

> 目標：以 `v2` 建立 canonical Output flow，且同時覆蓋 Local Modbus 與 Database 兩個 target families。

- [x] 4.1 共用基礎
  - 在 `main` 上更新 Output shared acceptance / scenario matrix
  - 固定共同動作：readiness、dry-run、apply、blocker diagnosis
  - 明確要求 Local Modbus 與 Database 兩個 target families 都必須被驗收
  - 明確定義 `Output` phase 的三 archetype operation language

- [x] 4.1.baseline Current Studio
  - 以 baseline snapshot 收集兩個 target families 的 overview / focused / blocker evidence
  - 記錄 baseline 在狀態可見性與阻塞理解上的痛點

- [x] 4.1.v2 Sentry Incident Desk（primary functional track）
  - 完成 Local Modbus + Database 的 canonical Output flow
  - 先讓 readiness、dry-run、apply、failure / retry / recovery 在 `v2` 正確
  - incident-desk target diagnosis / apply 操作面必須 clearly different
  - ✅ branch commit：`a73fff4`
  - ✅ 已完成 Local Modbus / Database 真實 browser evidence 與 branch docs 封帳

- [x] 4.1.v1 Linear Control Room（full-flow high-polish track）
  - 保持與 `v2` 相同的 Output domain flow與 target coverage
  - 以 clearly different control-room operator-console 語言表達 blocker 與狀態轉移
  - ✅ branch commit：`e714be8`
  - ✅ 已完成 Local Modbus / Database 真實 browser evidence 與 branch docs 封帳

- [x] 4.1.v3 ClickHouse Data Cockpit（minimum-obligation track）
  - 只保留 compare 所需的 dense mapping/state surface
  - 仍必須同時覆蓋 Local Modbus + Database 與 shared apply / blocker contract
  - 但不能只是 shared output page 加 density surface
  - ✅ branch commit：`5f357df`
  - ✅ 已完成 Local Modbus / Database 真實 browser evidence 與 branch docs 封帳

- [x] 4.1.compare Output 比較與推薦
  - 比較 baseline + `v2` + `v1` + `v3`
  - compare 必須明確評估兩個 target families，不可只看其中一個
  - 若 review 仍認為三版除了 Device 幾乎一樣，本 phase 不得關閉
  - ✅ compare 輸入：
    - baseline `4173`
    - `v2` `a73fff4` / `4175`
    - `v1` `e714be8` / `4174`
    - `v3` `5f357df` / `4176`
  - ✅ compare 結論：
    - **推薦版本：`v2 Sentry Incident Desk`**
    - `v1` 最適合吸收 calm rail / blocker framing
    - `v3` 最適合吸收 telemetry strip / bottom operations dock
    - compare 不足以推翻 Phase 4 的 canonical owner，維持 `v2`

---

## Phase 5：Cross-step shell / diagnostics winner-led rollout

> 目標：以 `v2` 建立 canonical shell / diagnostics flow，同時守住 shell ownership 邊界。

- [x] 5.1 共用基礎
  - 在 `main` 上更新 shell shared acceptance / scenario matrix
  - 固定 shell 只擁有：readiness summary、active blocker summary、diagnostics refresh status、shortest return action
  - 明確禁止 shell 吃掉 step-local edit / mutation / validation
  - 明確定義 shell / diagnostics 的三 archetype 語言與 fail gate
  - ✅ 新增 `workbenchShellCompareContract.ts`
  - ✅ `workbenchExperimentContract.ts` 已正式 re-export shell shared compare contract
  - ✅ `phase5.scenarioFocus` 已從舊的 `readiness / global-blockers / repair-hop / return-to-flow` 對齊為 `readiness-summary / active-blocker / diagnostics-refresh / return-to-mainline`

- [x] 5.1.baseline Current Studio
  - 以 baseline snapshot 收集 readiness、blocker、return-to-mainline evidence
  - 記錄 baseline 在跨步驟診斷上的強弱項
  - ✅ baseline `4173` + real device `UI 4.3 Modbus TCP`
  - ✅ 正式 evidence：
    - `baseline-shell-readiness.png`
    - `baseline-shell-return-to-mainline.png`
    - `baseline-shell-blocker.png`

- [x] 5.1.v2 Sentry Incident Desk（primary functional track）
  - 完成 canonical cross-step shell / diagnostics flow
  - 先讓 global blocker、refresh、retry、return action 在 `v2` 正確
  - incident-desk shell / diagnostics 操作面必須 clearly different
  - ✅ branch / commit / preview：`woe-v2-mui` / `3408f63` / `4175`
  - ✅ 正式 evidence：
    - `v2-shell-overview.png`
    - `v2-shell-refresh-success.png`
    - `v2-shell-return-to-output.png`

- [x] 5.1.v1 Linear Control Room（full-flow high-polish track）
  - 保持與 `v2` 相同的 shell ownership model
  - 以 clearly different、低干擾的 control-room shell 語言表達同一組能力
  - ✅ branch / commit / preview：`woe-v1-radix` / `0579eef` / `4174`
  - ✅ 正式 evidence：
    - `v1-shell-overview.png`
    - `v1-shell-refresh-success.png`
    - `v1-shell-return-to-output.png`

- [x] 5.1.v3 ClickHouse Data Cockpit（minimum-obligation track）
  - 只保留 compare 所需的 cockpit summary / alert surface
  - 不得把 shell 擴成 version-only workflow center
  - 但 summary / alert surface 必須仍然 clearly different，而不是 shared shell 加 cockpit 裝飾
  - ✅ branch / commit / preview：`woe-v3-antd` / `6f9bc5c` / `4176`
  - ✅ 正式 evidence：
    - `v3-shell-overview.png`
    - `v3-shell-refresh-success.png`
    - `v3-shell-return-to-output.png`

- [x] 5.1.compare Shell / Diagnostics 比較與推薦
  - 比較 baseline + `v2` + `v1` + `v3`
  - compare 必須確認 shell 沒有改寫 step-local ownership
  - 若 review 仍認為三版除了 Device 幾乎一樣，本 phase 不得關閉
  - ✅ baseline control：`b446045` / `4173`
  - ✅ compare rubric：
    - **操作順暢度**
      - baseline：最低；需先進 step-local surface 才能理解真正 blocker
      - `v2`：最佳；blocker / refresh / return action 集中在同一條 incident strip
      - `v1`：次佳；同樣可直接回主線，但 calmer spacing 讓掃描稍慢於 `v2`
      - `v3`：中上；return action 明確，但高密度 cockpit framing 提高首屏解析負擔
    - **邏輯清晰度**
      - baseline：最低；shell 只給 generic CTA，沒有 shell-level blocker summary
      - `v2`：最佳；四個 shell-owned surfaces 的責任邊界最清楚
      - `v1`：高；ownership 一樣乾淨，但需要多看一眼才會感覺到 urgency
      - `v3`：中；alert/telemetry identity 強，但較容易把注意力帶向 cockpit 語氣
    - **對系統的完整性**
      - baseline：部分；readiness 可讀，但 blocker / refresh / return contract 不完整
      - `v2`：最佳；是唯一把 canonical shell loop 做成最完整、最直接的版本
      - `v1`：高；shared shell loop 完整成立，且沒有侵入 step-local workflow
      - `v3`：達標；minimum-obligation shell loop 成立，但刻意不擴成 full workflow center
    - **首屏資訊密度**
      - baseline：最低
      - `v1`：中
      - `v2`：中高，平衡最好
      - `v3`：最高
    - **關鍵操作時間**
      - baseline：最慢；需要先進各 step 才能完成 blocker triage
      - `v2`：最短；讀取 shell 後可直接 refresh 或 return
      - `v1`：接近 `v2`；主要差距是 calmer framing 帶來的額外掃描時間
      - `v3`：接近 `v2`，但首屏 telemetry 較多，初次判讀成本最高
    - **實作 / 維護風險**
      - baseline：最低，但功能不足
      - `v1`：中低；新增 surface 小且 ownership 穩定
      - `v2`：中；surface 最完整，但 shell semantics 也最強，後續需守 canonical owner 邊界
      - `v3`：最高；視覺密度與 cockpit identity 最強，也最容易放大維護面與噪音
  - ✅ 推薦版本：維持 `v2 Sentry Incident Desk` 為 Phase 5 canonical owner
  - ✅ 推薦理由：
    - `v2` 在不吃掉 step-local ownership 的前提下，把 readiness、blocker、refresh、return 這條 shell loop 做到最短、最清楚
    - `v1` 是最好的長時間操作備案，建議回收 calmer blocker framing / lower-interference pacing
    - `v3` 的 alert line 在首屏資訊密度與版本辨識度上最強，建議只回收 alert / telemetry cues，不延用整條 cockpit intensity

---

## Phase 6：最終驗收與推薦

> 目標：以同一個 end-to-end 主線重跑 baseline、`v2`、`v1`、`v3`，輸出正式推薦與保留亮點。

- [x] 6.1 共用基礎
  - 固定最終 E2E 場景：
    - 使用真實 API + 真實裝置 `UI 4.3 Modbus TCP`
    - 依同一條 `/studio` 主線重跑 `Device -> Source -> Tag -> Output -> Shell`
    - `Source` 必須進入可操作的 focused rule / preview state，並保留 handoff to `Tag`
    - `Tag` 必須進入可操作的 review/apply state，並保留 handoff to `Output`
    - `Output` 必須同時覆蓋 Local Modbus 與 Database 兩個 target families
    - `Shell` 必須再次驗證 readiness summary、active blocker summary、diagnostics refresh、shortest return action
  - 固定最終 compare 輸出格式與證據欄位：
    - 每個版本至少保留 `overview`、`focused`、`next-step handoff` 三組 browser evidence
    - compare 必須逐一輸出：操作順暢度、邏輯清晰度、對系統的完整性、首屏資訊密度、關鍵操作時間、實作 / 維護風險、推薦版本與理由、非 Device phase distinctness 結論
    - 最終 compare 必須額外寫出：保留亮點、淘汰理由、後續 follow-up 建議
  - 固定最終「三版本在非 Device phase 仍 clearly distinct」驗收題：
    - reviewer 必須能在 `Source`、`Tag`、`Output`、`Shell` 四段，透過首屏與第一個主操作辨識 `v2` / `v1` / `v3` 是不同產品語言
    - 若差異只剩 component kit、spacing、顏色或 chrome，而主操作骨架與資訊節奏仍幾乎一樣，Phase 6 不得關閉

- [x] 6.1.baseline 整體驗收
  - 完成 baseline 端對端操作、截圖與記錄
  - ✅ preview / route / device：`4173` / `/studio` / `UI 4.3 Modbus TCP`
  - ✅ 正式 evidence：
    - `baseline-final-overview.png`
    - `baseline-final-source-focused.png`
    - `baseline-final-tag-focused.png`
    - `baseline-final-output-local-modbus.png`
    - `baseline-final-output-database.png`
  - ✅ baseline 記錄：
    - 同一條主線可從 `Device -> Source -> Tag -> Output` 穩定重跑，且 Output 兩個 target families 都可見
    - `Shell` 仍然只有 generic readiness / CTA framing，沒有 final compare 所需的 blocker / refresh narrative
    - 作為 control 足夠，但在非 Device phase 幾乎沒有 archetype differentiation，可作為最終 compare 的下限參考

- [ ] 6.1.v2 整體驗收（primary functional track）
  - 完成 `v2` 端對端操作、截圖與記錄

- [ ] 6.1.v1 整體驗收（full-flow high-polish track）
  - 完成 `v1` 端對端操作、截圖與記錄

- [ ] 6.1.v3 整體驗收（minimum-obligation track）
  - 完成 `v3` 端對端操作、截圖與記錄

- [ ] 6.1.compare 最終比較與推薦版本
  - 彙整所有 phase compare 結果
  - 明確指出推薦版本、保留亮點、淘汰理由與後續 follow-up 建議
  - 明確回答三版本是否在整條 `/studio` 主線上都仍是 recognizably different products

---

## 驗收條件

- [ ] baseline 與 `v2` / `v1` / `v3` 全程共用同一組真實 API，無 mock 分流
- [ ] baseline 與 `v2` / `v1` / `v3` 全程維持同一個 `/studio` route
- [ ] baseline 依 Phase 2–5 baseline 定義凍結，不隨 variant work 漂移
- [ ] 三個新版本全程使用同一組 shared design tokens
- [ ] 每個 shared task 都先在 `main` 上定義 branch-neutral acceptance / scenario matrix
- [ ] compare 子任務完成前，不得跳到下一個 phase
- [ ] 每個 compare 結果都包含操作順暢度、邏輯清晰度、對系統的完整性、首屏資訊密度、關鍵操作時間、實作 / 維護風險
- [ ] 每個 compare 結果都明確回答三版本在非 Device phase 是否仍 clearly distinct
- [ ] Phase 4 compare 同時覆蓋 Local Modbus 與 Database 兩個 target families
- [ ] 若 compare 挑戰 `v2` canonical owner，必須先回到 OpenSpec amendment，不可直接換線
- [ ] 最終輸出包含推薦版本與明確理由

---

## 後續擴充規則（v4+）

- [ ] 若新增 `v4+`，必須先補 proposal / design / tasks 矩陣，再開始該版本實作
- [ ] 若新增 `v4+`，必須沿用 shared design tokens 與同一組 `/studio` + API 契約
- [ ] 若新增 `v4+`，必須補齊對應的 `baseline -> vN -> compare` 證據欄位，才能納入正式比較
- [ ] 若新增 `v4+`，必須先分配新的 branch / worktree / port，並重新 `openspec validate`
