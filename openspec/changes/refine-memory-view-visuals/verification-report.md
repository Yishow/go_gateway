# D 實作對照驗證

## 範圍與結論

基準 HEAD：5047ab54f94e2bc185cc3f5c33e092ebc1a4b458。以官方 touched tracking、完整目前差異和明確授權的 recovery 回歸測試為範圍。DESIGN.md 開始前的 45 行增補由呼叫端指定屬 D；保留完整內容，只更新執行順序及實作狀態。逐檔歸屬不是逐 hunk 歷史證明。

官方 list 目前為 8/16 tasks：雲端非 browser 工作完成，八項 browser/closure 保留未完成。

四項 requirement 均有實作或驗收 harness；八個 scenario 的 assertions 與執行狀態分開記錄。原 spec 沒有 Example 區段，因此具體 Example 值追蹤為 0 項，不虛構例子。analyze 有七項缺少具體 Example 的 Suggestion，沒有 Critical/Warning；validate 通過。

## Requirement / scenario 對照

| Requirement | 實作 | 有測試的 scenarios | 本次執行 |
| --- | --- | --- | --- |
| Memory views expose distinct inspection structures | RecordingHomeMemoryPanels.tsx:26–152、View.tsx:214–226；CSS grid/章節/摘要/inspector 分區 | A user chooses each inspection mode：visual-navigation/map/chapters/detail 單元；Dense map or long chapter list：256 格、精確 slider、40 range 與最後一項定位單元，以及新 E2E | 單元見 non-browser-evidence；E2E not-run |
| Typed cells preserve physical address meaning and safety | Cells.tsx:25–50；Panels.tsx:74–152，共用 canonical map/values 不變 | A multi-register value crosses a boundary：continuity/visual-cells/values；A cell action is requested：recovery/scope/keyboard/visual-detail；unknown/conflict/partial fail-closed assertions 保留 | 單元見 non-browser-evidence；真實 geometry not-run |
| Layout changes retain one shared interaction context | View.tsx:68、97–124 與既有唯一 parent owner | Repeated layout changes and pointer inspection：hover/observer-states/values/layouts；Small viewport and dismissal：既有 sizing/navigation、production-modal view suite 和新 visuals low-height/touch | 單元見 non-browser-evidence；browser focus/history/dirty/unknown barriers not-run，最後 routed A/B harness 聯合補強 |
| Visual acceptance uses current original-reference evidence | reference-manifest.json；新 recording-memory-visuals.spec.ts:84–231；audit helper | Current visual and interaction evidence is reviewed：30 基礎 captures 和 9 邊界 cases，SHA/viewport/fixture metadata；Browser execution is blocked：接受流程保持未完成，沒有偽造替代證據 | 原圖五張直接檢視且完整 SHA/尺寸吻合；只 discovery/types/lint，擷取/對照/互動 not-run |

前三項行為 scenario 已有測試。最後環境受阻 scenario 是流程契約，依 test-scope 不跨產品使用點，排除產品單元斷言；以未完成 checklist 與報告保留其真實狀態，不算產品 coverage gap。局部 padding/radius 等不補孤立單元測試，但不得省略明確要求的最終視覺驗收。

## Critical：最後 Mac 驗收尚未執行

- tasks.md 3.1–3.5：五版型 × light/dark × 1440/768/390 的 30 張目前產品畫面及逐圖 accepted / intentional difference / defect 判定。
- tasks.md 3.6：low-height、安全區域、完整長值/名稱、unknown/conflict、四格跨列/頁、實際 sizing/對比。
- tasks.md 3.7：keyboard/touch/hover、切換、Close/Cancel/Escape、Back/Forward、dirty/unknown-operation barriers 的真 browser 行為。
- tasks.md 4.2：完成上述證據後重跑 verify/review、修正重驗，才 archive。

這八項 incomplete tasks 是 verify 的 Critical，不代表目前 source review 發現八個 bug。按授權留待四案全部雲端實作和 harness 完成後合併在 Mac 執行；不得勾選或 archive。3.8 非 browser QA 結果另記，不能以 artifacts status isComplete 當作產品已驗收。

## Coherence 與限制

DESIGN.md 是唯一 geometry/token authority；原圖只管構圖，Quiet mineral 配色為明確有意差異，不冒稱逐 pixel 相同。保留 memory-only 85vw/85dvh exception、初始兩列密度、單一 grid/parent owner、實體位址、unknown/MC proof 與 B/C scope/ACK fences。

對比 helper 量測文字及部分控制/focus；最後仍需人工/像素補 native slider thumb/track、hatch 與全部必要狀態，文字至少 4.5、控制/focus 至少 3。isolated div-dialog 不可當作 production/embedded focus/history 證明。

範圍快照、完整 bytes/hash、command/log/exit 與最終 scope check 隨提交計畫提供。無 browser、Mac、archive、stage、commit、push、後端操作或權限提升。
