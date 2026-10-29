## 1. 共用 scope 與引導模型

- [x] 1.1 依「分期準備與啟用門檻」與「相依責任與漸進替換」，在獨占新 recordingWorkflowView.ts 固定既有 group/revision 的純 caller contract；以隔離型別／contract tests 證明 saved/applied/unknown 不混用、缺 lifecycle evidence 回 unavailable、沒有第二份 readiness authority。可先於前案驗收準備，但不掛載或啟用新 live transition。
- [x] 1.2 [after: 1.1] 僅在新 recordingWorkflowView／RecordingNextAction 及專屬 tests 以測試先行實現「單一 scope 與下一步」及 `One selected scope and next safe primary action` 的 pure derived model，覆蓋空來源、offline draft、missing semantics、schema conflict、unknown operation；先觀察失敗再加 next-action，focused tests 通過後重構，不改 server state machine。
- [x] 1.3 [after: 1.2] 實現「標準 scope 選擇介面」及 `Canonical setup selection revalidates ownership`；先測 typed request、foreign/missing/deleted ID、historical revision hint 無 mutation 授權，再實作只接受 read-only ownership/revision evidence provider 的獨立 studioSelection adapter，無證據回 unavailable；component tests 證明選擇不 save/apply/start，真 query/route 掛載留 2.2 後；URL parsing/resolver 唯一移交後續 runtime-rules change。

- [x] 1.4 [after: 1.3] 補償已完成 1.1／1.2 的 caller contract 與 pure next-action：依 design「安全生效證據交接」改消費 runtime 唯一安全 DTO，不要求 effective revision 等於 saved 的 latest requested applied_revision，不把 GroupFacts cast／補值重建成 WriteGroup；以 RED→GREEN contract/next-action tests 驗證 saved C、scheduled B、as-of effective A、boundary 前後、stale/missing、ended/tombstone retained snapshot 不授權 intake，保留 unknown-operation 優先與所有原 guards。先與 runtime 1.2 的契約補償交接；本任務不掛載、不宣稱真 reader 完成。

## 2. 摘要與漸進式流程

- [x] 2.1 [after: 1.4] 僅在未掛載 RecordingRuleSummary 與專屬 tests 完成「摘要先行與進階詳情」及 `Rule summaries distinguish saved and effective configuration`，新增只讀摘要列來源、layout/profile/reduction、週期與單位、缺值規則及目的表；先測 draft/applied 不同及 stale evidence，再實作並通過 summary tests。
- [x] 2.2 [after: 2.1] 確認 freeze→cutover→lifecycle→authoring 最終驗收、runtime 1.3–1.4 安全 read projection／redaction／effective provider 及共用 service/hook 真實驗收交接後（不等待 runtime dashboard／locator／最終驗收），並完成 shared Step4／locale 檔案排他移交後，將 destination/source/mode 區逐批接入共用 scope 與摘要，保留可展開型別、欄位、binding、counter policy、epoch 與 revisions；用 component tests 確認所有既有選項能力不縮減、非計數 basic 不多要求量測語意，每批可 typecheck/build。
- [x] 2.3 [after: 2.2] 將 preparation/start 區改為一個主要下一步，但保留 save、DDL review/confirm、start、SQL committed/readback/cleanup 的不同階段；先測沒有 confirmation 不發 DDL/start，實作後 integration tests 驗證阻擋理由與 repair focus。 依「Owner-local 操作觀察與唯讀導覽分界」補正已準備 caller：區分 unresolved／observed owner-local none 或 resolved／unavailable，不假造全域 absence 或 lifecycle review_scoped_start authority；O1／無 ID retained request 優先，unreadable／invalid 不當 idle。NAV 只 open/focus，preview 是既有 owner POST command；review scoped start 只開重用摘要的 presentation-only scope-review／explicit-confirm 區，bridge 不傳 raw request/draft/token/error 或 operation DTO；schema 非空 opaque ID 仍 retained、不新增格式解析，schema/start 最終確認走同一既有 handler 並於每次 invocation 重驗 scope/epoch/recovery/proof。先測上述分支、offline NAV／network gate、無 consent／stale callback 無 mutation；必要時僅 additive readonly inspection 既有兩個 helper，不改 persistence format／retry／backend authority，保留 standalone／Advanced 行為。

## 3. 記錄管理與導航保護

- [x] 3.1 [after: 2.3] 實現「已設定記錄與工作區狀態」及 `Recording management is discoverable without a second activation path`，卡片重用 lifecycle action availability，Advanced 更名；先測沒有新增 workspace start、legacy/custom 不被轉換，再替換重複 presentation 並通過 tests。
- [x] 3.2 [after: 3.1] 交付 `Navigation preserves drafts and safety boundaries`，覆蓋 Basic/Advanced、Back/Cancel、切設備、空設備、reload 與未決 operation；先建立不發 mutation／不丟草稿測試，再修正焦點與來源捲動，維持相依 change 的 operation recovery。
- [x] 3.3 [after: 3.2] 將完成摘要與 runtime handoff 綁到實際 activated IDs，部分成功明列未完成範圍；以 integration/router tests 證明不擴大為全部設備、重複 click 不重啟、返回仍看到本 scope，完成後才移除重複操作區。

## 4. 最終使用者驗收

- [x] 4.1 [after: 3.3] 用正式 browser tests 與截圖驗收新設定、重開 wide/long/custom、schema/lookup failure、unknown operation、partial start、typed scope selection；覆蓋 keyboard Tab、repair focus、drawer focus restore 與 390/768/1280 CSS px 視窗寬（設計驗收條件），主要動作及錯誤修復可操作且無頁面水平溢出，普通非計數設定不需繞 Advanced 或另一個 start，逐一核對 spec scenarios。 依「現有雲端瀏覽器驗收的條件式交付例外」驗證 `Browser verification exceptions remain explicit`：保留正式 cases／discovery／E2E 專用 typecheck／lint／scenario 對照；僅當既有 cloud blockers 未變且實際執行仍不可用，實際 browser／截圖標 NOT RUN／SKIP／未驗證（不可標 PASS），其餘原驗收照跑，不新增環境或繞過安全限制。
- [x] 4.2 [after: 4.1] 執行相關 frontend suites、typecheck、lint/build、git diff --check、make check-lines，檢查 en／zh-TW 的 scope/階段文案一致；交接實際 PASS／FAIL／NOT RUN，不以 UI 綠色狀態代替 SQL／現場驗收。

## 待辦對應與執行界線（非完成紀錄）

- 原 1.1 的「前案全完成後才開工」替換為 1.1 的 isolated caller contract；原 backend／lifecycle／authoring 整合核對移到 2.2 的掛載 gate，不取消。
- 原 1.2、1.3、2.1 保留編號與完整行為；準備限新檔，1.3 的真 query/canonical receiver 接線保留在 2.2–3.2。其餘 2.2–4.2 原編號、順序及最終驗收保留。
- labels、fixtures、未掛載 adapter 不代表真 UI 已完成；所有產品 spec scenarios 仍須 4.1–4.2 證明。已完成 1.1／1.2 的描述及 provenance 原樣保留；新增 1.4 為 pending compensation，不回寫歷史或勾選產品完成，既有 [P] 不變。

- 本次插入 1.4：CLI checkbox indices 1–3（文件 1.1–1.3）不變；原 CLI 4–11 → 5–12，原文件 2.1–4.2 編號不變。新 CLI 4 為 1.4；所有 after 指向現行文件編號。
