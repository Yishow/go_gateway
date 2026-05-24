## Context

repo 目前同時擁有完整工作台 `/studio`、快速設定 `/gateway/quick-setup` 與進階設定 `/gateway/expert-workbench`。這些能力本身不是問題，問題在於產品缺少一個前置入口，讓第一次進來的操作員先做路徑選擇：是要用最短路徑快速完成設定，還是直接進入既有完整工作台或進階模式。

前一版提案把焦點放在收斂入口，但使用者已明確指出目前系統功能不能遺失。這次 ingest 後，change 的正確方向是 additive UI：新增入口頁，保留既有流程、既有 deep link、既有工作台，不以功能整併為成功條件。

## Goals / Non-Goals

**Goals:**

- 新增一個清楚的產品入口頁，幫助操作員快速開始完整設定。
- 讓入口頁把 `快速完整設定` 作為主要 CTA，優先導向現有 quick setup 流程。
- 保留既有 `/studio`、`/gateway/quick-setup`、`/gateway/expert-workbench`、`/test` 的能力與可達性。
- 確保入口頁只是前置導流，不把既有工作台降級為相容模式或隱藏功能。
- 補齊路由、copy 與測試，證明新增入口沒有造成舊功能回歸。

**Non-Goals:**

- 不收斂既有路由家族成單一路由。
- 不移除或退役 quick setup、expert workbench、`/studio`。
- 不重做 workbench provider、Source/Tag/Output 業務邏輯或 backend contract。
- 不把 `/test` 納入入口頁主線。

## Decisions

### Preserve existing product routes and shells
- Decision: `/studio`、`/gateway/quick-setup`、`/gateway/expert-workbench` 維持原本產品能力與直接可達性；新增入口頁不會取代它們。
- Rationale: 使用者已明確要求功能不得遺失，最小改動就是新增前置導流，不碰既有核心頁面定位。
- Alternative considered:
  - 將現有路由全部收斂到 `/studio`：與新需求衝突，不採用。

### Reuse `/gateway/entry` as the additive setup front door
- Decision: 直接使用現有 `/gateway/entry` 路由與頁面作為新增入口頁的承載位置，而不是再創一個新路由家族。
- Rationale: repo 已有對應 skeleton，重用現有頁面是最小改動，也能保留既有 feature flag 與導流接點。
- Alternative considered:
  - 新增全新入口路由：會增加維護面與導流複雜度，不採用。

### Quick-complete CTA composes existing setup flows
- Decision: 入口頁的主要 CTA 導向現有 `/gateway/quick-setup`，並用 copy 說明它能快速完成完整設定；次要 CTA 保留前往 `/studio` 與 `/gateway/expert-workbench`。
- Rationale: 這符合「快速完整設定」需求，同時不重寫 quick setup 已有的實作。
- Alternative considered:
  - 在入口頁內直接重新實作 quick setup 表單：重複既有邏輯，不採用。

### Feature-flag behavior only gates the new entry page, not existing workflows
- Decision: 若保留 `ENABLE_GATEWAY_DUAL_ENTRY`，它只影響是否顯示或導流到入口頁，不得讓 `/studio`、quick setup、expert workbench 失效。
- Rationale: 旗標可以控制新入口 rollout，但不能拿來封鎖既有功能。
- Alternative considered:
  - 使用旗標關閉整個 gateway 路由族：有功能遺失風險，不採用。

### Verification protects against regression of existing routes
- Decision: 驗證優先檢查三件事：入口頁是否存在、快速 CTA 是否能進入既有 quick setup、既有 `/studio` 與 expert route 是否仍可直接開啟。
- Rationale: 這次 change 的風險不是新功能做不出來，而是新增入口時誤傷既有路徑。
- Alternative considered:
  - 只測入口頁畫面：不足以證明沒有回歸，不採用。

## Implementation Contract

**Behavior**

- 使用者開啟 `/gateway/entry` 時，系統 SHALL 顯示一個產品入口頁，清楚區分 `快速完整設定`、`進入工作台`、`進階設定` 等路徑。
- 使用者點擊 `快速完整設定` 後，系統 SHALL 進入現有 `/gateway/quick-setup` 流程，而不是新的重複表單。
- 使用者仍 SHALL 能直接開啟 `/studio` 與 `/gateway/expert-workbench`，不需要先經過入口頁。
- `/datalink/devices/new` 或其他既有建立裝置導流若使用入口頁，導入後 SHALL 仍能到達原有 quick / expert / studio 路徑。
- `/test` SHALL 維持原樣，不受入口頁影響。

**Interface / data shape**

- 入口頁路由沿用 `/gateway/entry`。
- 入口頁至少提供三個可辨識行動：導向 `/gateway/quick-setup`、導向 `/studio`、導向 `/gateway/expert-workbench`。
- 若使用 feature flag，旗標關閉時的 fallback SHALL 指向既有可用流程，不得落到不存在或被封鎖的頁面。

**Failure modes**

- 若入口頁旗標或設定不存在，系統 SHALL 保持既有 `/studio`、quick setup、expert workbench 直接可用。
- 若入口頁無法取得旗標狀態，系統 SHALL 顯示可恢復的 loading / fallback，而不是讓既有路由失效。
- 若使用者跳過入口頁直接進入既有頁面，系統 SHALL 正常工作，不依賴入口頁先寫入任何前置 state。

**Acceptance criteria**

- 路由測試可證明 `/gateway/entry` 可開啟入口頁，而 `/studio`、`/gateway/quick-setup`、`/gateway/expert-workbench` 仍可直接到達。
- 入口頁測試可證明 `快速完整設定` CTA 會導向既有 quick setup 路徑。
- 回歸測試可證明既有 expert 與 workbench 路由不會因入口頁新增而被 redirect、隱藏或封鎖。
- 文案檢查可證明入口頁將 quick setup 描述為快速完成完整設定，而不是暗示其他功能已被移除。

**Scope boundaries**

- In scope: 入口頁本身、入口頁導流、feature flag fallback、相關前端文案與測試。
- Out of scope: 收斂既有功能頁、重寫 quick setup 內容、重做 workbench shell、後端 API 調整。

## Risks / Trade-offs

- [Risk] 入口頁文案若過度偏向 quick setup，可能讓使用者誤以為 expert 或 `/studio` 被弱化  
  → Mitigation: 在同頁保留明確次要 CTA，清楚說明既有完整功能仍可直接進入。

- [Risk] feature flag 守門若寫錯，可能封鎖既有 gateway 路由  
  → Mitigation: 將旗標責任限制為入口頁 rollout，不影響既有 quick / expert route 的直接可用性。

- [Risk] 新入口頁與現有 entry skeleton 的測試基準不同，容易造成測試漂移  
  → Mitigation: 測試集中驗證可達性與 CTA 導流，不依賴脆弱的全頁細節快照。

## Migration Plan

1. 重新定義 `/gateway/entry` 的產品定位與 copy，使其成為新增入口頁。
2. 調整 route guard / feature flag 行為，確保它只影響入口頁 rollout，不影響既有 direct routes。
3. 補齊入口頁 CTA 與導流測試，驗證 quick setup、expert workbench、`/studio` 仍可直接開啟。
4. 執行前端回歸測試，確認新增入口頁沒有造成既有功能遺失。

## Open Questions

- 無。新需求已明確把範圍限制在 additive entry page，不需要再擴張到功能整併。
