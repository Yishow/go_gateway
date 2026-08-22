## Context

目前 9 個 frontend/tests/unit 測試 monolith 已超過專案硬上限 500 行，且最大的 DatalinkWorkbenchSourceStep.test.tsx 達 1902 行。這些檔案把多個 UI surface、狀態轉移、mock 與 fixture 放在同一個 suite；拆分時最大的風險不是檔案搬移，而是重複 setup、遺失測試探索、改變測試名稱或改變共享狀態。

這個 change 只處理 test-only 結構。原始測試的可觀測契約是 Vitest 能在 frontend/tests/unit 探索所有分割檔，測試名稱與斷言結果保持不變，且每個產出的 TS/TSX 檔案不超過 300 行。Workbench full-suite timeout 是既有基線，必須被量測與記錄，但不是本 change 的隱含修復目標。

## Goals / Non-Goals

**Goals:**

- 將 9 個指定 monolith 依可辨識的行為群組拆成 53 個 test-only harness 或測試檔，並移除原始 monolith。
- harness 只集中共用 render、router、mock、fixture、factory 與 cleanup；行為斷言留在 discoverable 的 .test.tsx 檔。
- 保留原始 test title、assertion、mock response、fixture data、路由與測試 order 不變；只做必要的 import 與 harness 介接。
- 每個產出的 TS/TSX 檔案維持不超過 300 行，且所有新 discoverable 測試仍位於 frontend/tests/unit。
- 先對每個原始檔建立 characterization 證據，再對 moved tests 執行 focused regression，最後執行 frontend quality gates 與 repo line gate。

**Non-Goals:**

- 修改 frontend/src、Go、API、runtime、database、Vite/Vitest 設定或 production bundle。
- 修改測試語意、斷言、測試名稱、fixture 值、mock 契約或增加平行入口。
- 新增 .line-limit-ignore、提高測試 timeout、跳過測試，或保證既有 full-suite timeout 因拆分消失。
- 將測試搬出 frontend/tests/unit，或把共用 harness 提升成 production/shared runtime module。

## Decisions

### Decision: 以行為邊界拆分 monolith

每個原始檔先以既有 describe 與 test title 的語意分組，再建立小型 discoverable test file。Foundation 依 legacy route、operator flow、device list、device editor、device inspector 分組；Output 依 Modbus surface/bindings/canvas、database target/schema/columns、inspector、cross-step handoff 分組；Source 依 planning、inspector、rule actions、canvas overlays、batch selection、inline editing、conflict queue、datatype/span、visual span 分組。這比固定等量切割更能保留 setup 與測試責任，也讓失敗訊息對應單一 operator behavior。

替代方案是只按原始檔案每 N 行切開，會把 describe、fixture 與狀態轉移截斷，增加重複 setup 與隱藏耦合，因此不採用。

### Decision: test-only typed harness 與行為檔分離

每個來源群組提供一個同目錄的 .testHarness.tsx，集中 typed render helper、router wrapper、mock registration、fixture factory 與 teardown。harness 檔名不含 .test.，不會被 Vitest 當成測試入口；discoverable 行為檔保留 .test.tsx。若某群組的 setup 不需要跨檔共享，仍可留在該行為檔，禁止為了抽象而新增不必要的跨群組 helper。

替代方案是把 setup 複製到每個新測試檔，會讓 mock drift 與 cleanup 不一致，因此不採用。替代方案是建立 frontend/src 共用測試 runtime，會擴大 production 邊界，因此不採用。

### Decision: 保留測試探索與契約

所有 discoverable 檔案維持在原始目錄下的 frontend/tests/unit，測試命名與 describe 階層可透過必要的最小外層調整保留辨識性，但不得刪除、重寫或合併原有 test title。每個原始 suite 的 characterization title set 必須與 moved suite 的 title set 一致；mock、fixture、路由與 API request shape 必須由同一測試情境驗證。

替代方案是改寫成新的參數化測試或共用 describe factory，可能改變 failure granularity 與測試名稱，因此不採用。

### Decision: characterization-first 與分層驗證

每個來源群組在 split 前先執行原始檔的 focused test，記錄 line count、test title set、pass/fail 與已知 timeout；split 後執行對應 destination glob，核對 title set、line count 與 focused result。全部群組完成後才執行完整 frontend test、lint、build 與 repo line gate。Workbench full-suite timeout 必須以實際執行時間與原始錯誤獨立記錄，不以 timeout 增加掩蓋。

替代方案是先大幅搬移再以完整 suite 猜測回歸，難以定位遺失測試或 setup drift，因此不採用。

## Implementation Contract

### Observable behavior

- Vitest 在 frontend/tests/unit 下能探索每個新 .test.tsx 檔，且不探索 .testHarness.tsx 為獨立測試。
- 9 個原始 monolith 不再存在；每個來源群組的測試 title、assertion、mock、fixture 與路由契約在 destination files 中可被驗證。
- 每個產出的 TS/TSX 檔案（含 test-only harness）行數不超過 300；不得以 ignore 規則豁免。
- product source、runtime behavior、API payload、Vite/Vitest 設定與測試 timeout policy 不變。

### Exact split ownership

以下 ownership 以原始 test title 與群組責任為準；實作時不得把一個 title 複製到多個 destination，也不得留下未歸屬 title。

#### Foundation source

Source: frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchFoundation.test.tsx

- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchFoundation.testHarness.tsx：共用 router、render、device fixture、mock 與 inspector helpers。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchFoundation.legacy-routes.test.tsx：legacy local-modbus、nested workbench、/datalink、/test redirect 與 /test shell。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchFoundation.operator-flow.test.tsx：/datalink/workbench redirect、deep-link query、navigation unlock、compact toolbar、hero removal、main operator flow、normal /studio flow、source engineering tools。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchFoundation.device-list.test.tsx：compact device row endpoint/protocol display、device select、master-detail、context bar、create action、refresh selection、description clearing。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchFoundation.device-editor.test.tsx：inline create form、localized protocol options、inline editor、draft config test、clone defaults、duplicate clone validation。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchFoundation.device-inspector.test.tsx：zero timestamp hiding、selected details、pending test action、latest result、recent timeline、connect/probe phases、polite live region。

#### Output source

Source: frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.test.tsx

- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.testHarness.tsx：共用 selected device/tag、output fixture、router、mock API、render helpers。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.modbus-surface.test.tsx：empty state、Local Modbus start/port、tag chips、supporting panels、unified selection surface。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.modbus-bindings.test.tsx：Modbus register bind、tag/register synchronization、register conflict blocking。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.modbus-canvas.test.tsx：register map canvas ranges/conflicts、multi-word overlap、high registers、unbind、empty-slot bind、auto-map strategy/sequential behavior、dry-run validation。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.database-targets.test.tsx：database direct column bind、mapping removal、selected-tag display alignment。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.database-schema.test.tsx：connector switch stale scope、write-mode/timestamp synchronization、manual write-mode preservation、connector field collapse、schema supporting-panel placement。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.database-schema-columns.test.tsx：schema column type badges、primary-key badge、grouped row planner、required unmapped column highlighting。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.output-inspector.test.tsx：source-to-tag-to-output trace、partial/ready readiness reasons、database trace path、inspector load-failure live region。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchOutputStep.cross-step-handoff.test.tsx：Step 3 focused tag preselection and first-candidate fallback when focusedTagIds do not match。

#### Shell source

Source: frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchShellUi.test.tsx

- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchShellUi.testHarness.tsx：共用 shell render、step/device fixtures、route state 與 MUI query helpers。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchShellUi.frame.test.tsx：shell regions、step rail embedding、desktop sizing、legacy ActionDock/HeaderBar absence。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchShellUi.context.test.tsx：four-step rail、active aria-selected、device/no-device context、primary action、source command rail/workspace sizing、tag-review recovery routing。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchShellUi.inspector.test.tsx：step-dependent inspector heading 與 device-step empty state。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchShellUi.summary.test.tsx：selected-device summary metrics、four-step readiness、output target badge、rich readiness data attributes、blocked tag/output status、active-step emphasis。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchShellUi.switching.test.tsx：step rail click 後將正確內容接入 PrimaryWorkArea。

#### Source source

Source: frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.test.tsx

- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.testHarness.tsx：共用 selected device、source rules、planner fixture、router、API mock、canvas queries。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.planning.test.tsx：device gate、planner/added-rules tabs、last start/protocol defaults、count validation、rule layer/gap cells、persisted rule load、runtime counters、per-device draft isolation。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.inspector.test.tsx：unmanaged points、rule inspector、span inspector、selected cell aria-pressed、edited-rule inspector retargeting。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.rule-actions.test.tsx：rule apply/update persistence、add/delete workflow、active rule workspace、source template save/reapply、navigation persistence、orphan point deletion。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.canvas-overlays.test.tsx：fixed 16-bit lattice、primary/supporting workspace layout、canvas scroll dominance、plan/live/link overlays、merged spans/gap cells。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.batch-selection.test.tsx：batch persistence with safe/conflicting spans、selection toolbar、create point、skip planned address、safe-span filtering。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.inline-editing.test.tsx：inline start/count/type editing、cancel without mutation、immediate canvas update、non-planned cell toolbar absence。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.conflict-queue.test.tsx：actionable conflict queue、visible inline edit entry、continuation-only conflicts、skip-span conflict isolation。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.datatype-span.test.tsx：grouped backend-aligned data type selector、protect-plan wording、root-cell snapping、selected logical span aria state。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchSourceStep.visual-span.test.tsx：merged 32-bit gridColumn rendering、continuation hiding、stale skipped-span cleanup after geometry edit、safe persistence beside conflicts、overlap root resolution。

#### Tag source

Source: frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchTagStep.test.tsx

- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchTagStep.testHarness.tsx：共用 candidate/tag fixture、selected rows、router、mutation mocks、render helpers。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchTagStep.review-overview.test.tsx：empty state、go-to-source、preview key/prefix、source metadata、review metrics、adaptive guidance、manual correction grouping。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchTagStep.review-selection.test.tsx：existing-tag selectors、selection-gated batch actions、delete confirmation、row/inspector conflict detail、candidate keyword/status filtering、master overview counts。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchTagStep.tag-binding.test.tsx：standalone tag create/delete、selected-point existing-tag binding、preview conflict block。
- frontend/tests/unit/pages/datalink/workbench/DatalinkWorkbenchTagStep.tag-results.test.tsx：partial bind failure、batch diff preview/skipped/toBind rows、created/linked/skipped/failed summary、single/batch unbind cancel/confirm、collapsed template panel。

#### Mui source

Source: frontend/tests/unit/pages/datalink/workbench/MuiOutputIncidentDesk.reopen.test.tsx

- frontend/tests/unit/pages/datalink/workbench/MuiOutputIncidentDesk.reopen.testHarness.tsx：共用 incident desk bootstrap、ready output flow、render/query helpers。
- frontend/tests/unit/pages/datalink/workbench/MuiOutputIncidentDesk.reopen.command.test.tsx：command surfaces、Local Modbus dry-run、Database desk connector editing。
- frontend/tests/unit/pages/datalink/workbench/MuiOutputIncidentDesk.reopen.recovery.test.tsx：blocker-first recovery copy、calmer summary/handoff layout、noRule handoff copy。

#### Database autosave source

Source: frontend/tests/unit/workbench-v2/database-autosave-page.test.tsx

- frontend/tests/unit/workbench-v2/database-autosave-page.testHarness.tsx：共用 database config/target/row-group fixture、hydration、request mock 與 render helpers。
- frontend/tests/unit/workbench-v2/database-autosave-page.hydration.test.tsx：persisted connector/targets 對應 point rows 的 hydration。
- frontend/tests/unit/workbench-v2/database-autosave-page.validation.test.tsx：valid connector save 與 invalid local edit 保留。
- frontend/tests/unit/workbench-v2/database-autosave-page.row-groups.test.tsx：single target failure isolation、row-group config ordering、stale row-group cleanup/table change、connector save completion gating。
- frontend/tests/unit/workbench-v2/database-autosave-page.bulk-toggle.test.tsx：每個 persisted row 的 bulk database target enable autosave。

#### Device autosave source

Source: frontend/tests/unit/workbench-v2/device-autosave-page.test.tsx

- frontend/tests/unit/workbench-v2/device-autosave-page.testHarness.tsx：共用 device fixture、persisted/running state、draft、save request mock 與 render helpers。
- frontend/tests/unit/workbench-v2/device-autosave-page.hydration.test.tsx：persisted device order 與 running truth hydration。
- frontend/tests/unit/workbench-v2/device-autosave-page.draft-validation.test.tsx：valid local save、invalid local values、draft recovery marker、invalid running device unavailable state。
- frontend/tests/unit/workbench-v2/device-autosave-page.failure-recovery.test.tsx：per-device save failure isolation、apply_failed visibility、apply_failed 不被誤判為 unrecovered draft loss。

#### Mapping autosave source

Source: frontend/tests/unit/workbench-v2/mapping-autosave-page.test.tsx

- frontend/tests/unit/workbench-v2/mapping-autosave-page.testHarness.tsx：共用 mapping rule/address row fixture、hydration、mutation mock 與 render helpers。
- frontend/tests/unit/workbench-v2/mapping-autosave-page.hydration.test.tsx：persisted mapping rows 對應同一 rule/address rows 的 initialization。
- frontend/tests/unit/workbench-v2/mapping-autosave-page.validation.test.tsx：valid draft row save 與 invalid local edit 保留。
- frontend/tests/unit/workbench-v2/mapping-autosave-page.reconciliation.test.tsx：per-row failure isolation、stale mapping_id recreation、save-error owning rule suppression。
- frontend/tests/unit/workbench-v2/mapping-autosave-page.bulk-toggle.test.tsx：每個 persisted row 的 bulk mapping enable toggle autosave。

### File and discovery contract

- Each listed source monolith SHALL be removed after its destination title set is verified.
- Every destination path above SHALL remain below 300 physical lines, including harness files.
- Only files ending in .test.ts or .test.tsx contain test cases; .testHarness.tsx contains no top-level test registration.
- All destination files remain under frontend/tests/unit and are included by the existing Vitest discovery pattern without configuration changes.
- No test title may appear in more than one destination. A title-set comparison SHALL cover all nine source suites.

### Failure modes and boundaries

- If characterization and moved title sets differ, the split is incomplete and SHALL stop before deleting the source monolith.
- If focused tests fail only after a move, the implementer SHALL restore the original local behavior through import/harness wiring, not alter assertions or increase timeout.
- If any destination exceeds 300 lines, the group SHALL be split further by behavior; .line-limit-ignore SHALL remain unchanged.
- Existing full-suite timeout, environment failure, or unrelated test failure SHALL be reported with exact command and output as baseline evidence; it SHALL NOT be reclassified as a split defect without a focused reproduction.
- Scope is limited to the 9 source test files, their 53 destination test/harness files, and the new capability artifacts. Production source, config, dependency manifests, timeout settings, and unrelated tests are out of scope.

### Acceptance criteria

- Before each group split, the original focused test and line count are recorded; after each split, destination focused tests pass or preserve the exact pre-existing failure and title-set comparison is equal.
- All destination TypeScript/TSX files are at most 300 lines, original 9 monoliths are absent, and make check-lines or its repo-supported equivalent reports no new line-limit violation.
- frontend lint, frontend test, frontend build, and the repository line gate are run; full-suite timeout evidence is recorded separately.
- Final diff review reports only the approved test/harness paths and OpenSpec artifacts, with no production source, timeout, ignore-list, or test-configuration changes.

## Risks / Trade-offs

- [Risk] Shared harness extraction can hide per-file assumptions or leak mutable mock state → Mitigation: keep harness typed and minimal, reset mocks per test, and compare focused title/result evidence before source removal.
- [Risk] Test discovery can silently omit a moved file → Mitigation: use destination glob collection and title-set comparison, plus the existing frontend test command.
- [Risk] Splitting can increase import/setup overhead and expose a pre-existing full-suite timeout → Mitigation: measure before and after, report exact timeout baseline, and do not apply blanket timeout changes.
- [Risk] Behavior grouping can still leave one destination above 300 lines → Mitigation: check physical line count after each group and subdivide within the same ownership boundary.

## Migration Plan

1. Create characterization records for each source group, including line count and test title set.
2. Implement the nine independent source-group splits in dependency order, keeping originals until destination evidence matches.
3. Run moved focused tests and line checks after each group; then run title-set/discovery verification for all 53 destinations.
4. Run frontend lint, full test, build, and make check-lines; record known timeout or unrelated baseline failures.
5. Perform final diff/audit review. Rollback is deleting destination-only changes and restoring the original monolith from version control; no production migration or data rollback is required.

## Open Questions

- 無。full-suite timeout 的後續處理需另立 change，不能在本 change 中擴大範圍。
