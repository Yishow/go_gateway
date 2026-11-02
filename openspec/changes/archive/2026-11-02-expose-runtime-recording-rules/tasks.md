## 1. 擴充安全規則讀取

- [x] 1.1 依「共用安全讀取模型」與 design「契約準備與 production 讀取分界」，先逐欄審定 RuleEvidenceView／EntityLocator／stage/action 的 Go/TypeScript 契約與來源，不新增 authority；再於指定新 runtimeview/type/parser 檔實作未掛載 parser，RED→GREEN 測 missing/unknown、unsafe fields、cross revision 與 nullable evidence，契約 review、Go contract tests、frontend parser tests／typecheck 通過才交給 logs。不得改 production reader、以 fixture 填 active 值或宣稱 operator API 已 redacted。

- [x] 1.2 [after: 1.1] 隔離補償已完成 1.1 契約：依「共用安全讀取模型」及 design「安全生效證據交接」區分 saved、latest requested applied 與 as-of effective revision／effective_at，更新同一 Go/TS DTO/parser，交接 guide 1.4 並通知 logs；不掛載 production。以 RED→GREEN contract/parser 及受影響 diagnostics consumer tests 驗證 missing/unknown、scheduled B/effective A、mixed revision、ended snapshot 非 intake permission、相容降級與 redaction；將既有 DeliveryView.CounterReplayQuarantines 的 workspace/group/revision、pending_samples、typed reason/action 納入安全契約；測 accepted-input quarantine 不混為 finalized SQL output、missing/unknown 與跨 scope 拒絕。維持既有 stage/reason/action 語意，不自行擴張 producer payload。此任務可在 backend 交接前隔離完成，不宣稱真 reader 或全 API 已安全化。
- [x] 1.3 [after: 1.2] 在 freeze→cutover→lifecycle 最終證據及逐檔移交後，完成第一批 production 安全 projection/redaction：同批移除／遮蔽所有 legacy/new operator API endpoint／credential／private-path，保留舊 UI 缺值相容。先補全舊新 API、remote-view 與 parser redaction 失敗測試，再實作並驗證 typed unavailable、missing/unknown/mixed revision 與所有 unsafe fixture 不外洩；此 gate 通過後才可接 1.4 真 effective provider，不等待 guided 最終驗收，不放寬 local-only guard。
- [x] 1.4 [after: 1.3] 在 freeze→cutover→lifecycle 交接後，以 Red-Green-Refactor 前移原 2.1 的 applied-config-provider：將 ResolveAppliedAt(server as-of) 與已驗證 frozen runtime layout 投影為同一安全 RuleEvidenceView，提供共用 read-only service/hook 供 guided 2.2 與後續 dashboard 消費；不新增平行 DTO／authority，不接 locator／URL。API/service/hook tests 覆蓋 saved C、requested B 未到 boundary 仍 effective A、boundary 後 B、foreign/deleted、讀取中切版、missing proof、retained ended snapshot 與 intake 分離；驗證 redaction、同 workspace/revision/as-of consistency、缺證據 unavailable，不以 editable GET 或 fixture 代替實況。
- [x] 1.5 [after: 1.4] 等待 guided 4.1–4.2 最終驗收及檔案交接後，以 Red-Green-Refactor 完成 Scoped device rule inspection 的單設備來源讀取與詳情：Point/PollingGroup 取得有效輪詢來源，顯示地址／解碼／轉換／Tag，effective 證據不足回 unavailable；用 scale、byte order、poll/write 不同週期與 foreign workspace fixtures 驗證。
- [x] 1.6 [after: 1.5] 將來源規則卡遷移到完整搜尋／分頁列表及詳情，只保留舊 caller 的安全欄位相容；測六規則十 mappings、預設五十／上限一百、cursor revision stale、A→B 遲到結果及請求數界線，確認不再只有 +N。

## 2. 對照生效規則與交付

- [x] 2.1 [after: 1.6] 以 Red-Green-Refactor 完成 Revision-bound recording rule inspection 的 dashboard 呈現及「生效版本與交付證據」：重用 1.4 唯一 provider／共用 service/hook，不另建讀取 authority；測 saved C、scheduled B、effective A 的目的地／layout／interval／reduction／boundary 對照、讀取中 revision 切換及 frozen proof unavailable。1.4 已負責原本本項的 applied-config-provider，不延後到 dashboard 才提供 guide。
- [x] 2.2 [after: 2.1] 完成 Delivery evidence remains authoritative 的群組 stage／backlog 卡，重用 DeliveryForRevision；測 collecting／queued 不當 SQL committed、disabled 仍有歷史 backlog、舊 connector revision 不被新草稿重新標記及共用群組 totals 標示；納入 counter-replay-proof-*／recording-ownership-*／recording-handoff-proof-conflict／recording-cutover-finalized-overlap 的 typed reason fixtures，finalized 與 replay quarantine 分開且只提供指引。
- [x] 2.3 [after: 2.2] 接上 source evidence／receipt 詳情並區分 test-write、runtime commit、readback、cleanup；測 mean partial、delta unknown、legacy/custom 無 proof、receipt mismatch 與同群組不同設備，不把缺資料補零或把群組成功分配給每個 Tag。
- [x] 2.4 [after: 2.3] 加入 evidence-backed timing 顯示及 nullable fallback；測 window／committed／oldest pending 的真實時間、沒有 next_due／next_retry／duration 時顯示 unavailable，確認 UI 不從 interval 合成執行倒數。

## 3. 定位與監看呈現

- [x] 3.1 [after: 2.4] 完成 Exact safe return to settings 與「精準定位與設備範圍」：resolver 綁定 workspace/device/rule/mapping/group identity，串接 guide-step4-recording-workflow 的 canonical scope-selection interface；測 Step 1–4、Back/Forward、deleted/stale/foreign locator，不以同名替代且不發 mutation。
- [x] 3.2 [after: 3.1] 落實 Safe focused monitoring and truthful freshness 的安全 projection，確認首批 redaction gate 持續覆蓋所有 legacy/new operator contract 與 DOM，保留安全 table alias；用 secret/endpoint/private-path/driver-error fixtures 對所有舊新 API 與 remote-view 做負測試，驗證無 log proxy、無 local-only guard 放寬，不將秘密移除延後到此批。
- [x] 3.3 [after: 3.2] 完成「Focused monitor 與降級」資訊順序、transport-event 標籤與 LIVE freshness；先寫 UI 測試再調整畫面，驗證規則／交付先可見、歷史報表仍可達、SSE 故障保留帶時間的 stale snapshot、重連恢復及 en／zh-TW 鍵盤操作。

## 4. 收斂與驗證

- [x] 4.1 [after: 3.3] 所有新讀取 caller 遷移後移除過時摘要與重複的安全 projection，unsafe 欄位已於 1.3 的第一個 production 批次移除且不能恢復；保留 backward-compatible parser 降級；執行受影響 Go handler／workspace／grouppipeline 測試及 frontend runtime-dashboard／service／routing 測試、lint、型別與 build，確認每個 spec scenario 有具名測試。
- [x] 4.2 [after: 4.1] 以 simulator／可丟棄 DB 跑 embedded browser 的規則→設定→runtime、draft/applied 差異、離線 backlog 與 redaction 驗收並保留截圖；依 docs/testing.md 執行最終受影響回歸、make check-lines、git diff --check、Spectra analyze/validate，記錄 NOT RUN 的現場與平台界線；若改 Studio inventory 同步 changelog.sqlite。 依「現有雲端瀏覽器驗收的條件式交付例外」驗證 `Browser verification exceptions remain explicit`：保留正式 cases／discovery／E2E 專用 typecheck／lint／scenario 對照；僅當既有 cloud blockers 未變且實際執行仍不可用，實際 browser／截圖標 NOT RUN／SKIP／未驗證（不可標 PASS），其餘原驗收照跑，不新增環境或繞過安全限制。

## 待辦對應與執行界線（非完成紀錄）

- 原 1.1 拆成新 1.1（精確契約／隔離 parser 準備）與 1.2（第一批 production API projection／全舊新欄位 redaction／相容與實際契約測試）。原 1.2 → 1.3；原 1.3 → 1.4；2.1 的 after 改為 1.4，其餘編號不變。
- 新 1.1 的歷史準備證據保留；1.2 補償契約並通知 logs，1.3–1.4 production 只等 backend→cutover→lifecycle 與逐檔移交，1.5 起 UI／locator／最終驗收仍等 authoring、guided 整合交接。不能先架上不安全 API，安全移除仍在第一個 production 批次完成。
- 原所有 reader、source proof、locator、timing、redaction、browser 與最終 acceptance 保留；已完成 1.1 原描述與 provenance 保留，1.2 是 pending compensation；其餘未執行任務不勾選，既有 [P] 不變。

- 本次順序修正取代上列歷史 pending 對應：1.2 僅契約補償，新增 1.3 全舊新 production redaction、1.4 provider（承接原 2.1 provider）；原 1.3 → 1.5、1.4 → 1.6，原 2.1 留 dashboard 並 after 1.6。CLI indices 1–2 不變；新 CLI 3–4 為 1.3–1.4，原 CLI 3–13 → 5–15。文件 2.1 以後編號不變。
