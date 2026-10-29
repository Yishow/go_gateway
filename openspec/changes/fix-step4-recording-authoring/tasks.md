## 1. 角色與必要資訊

- [x] 1.1 依「首個可執行任務封包：1.1」先以 repo-relative 正式測試重現 numeric state 缺選項，再 RED→GREEN→Refactor 實現「角色與計數名稱前置提示」中的 `Numeric state is an explicit latest-value role`：numeric 可明選 state/latest，gauge/mean 與 counter_raw confirmation 不混淆，saved role 不自動改；以 authored value_role=state、非 gauge 的 payload 及 saved-role disabled／不變測試通過後重構，不在前端假算 1→2；最終真 backend latest 2、非 1.5 與 saved-role review 留 4.1–4.2，現在不解鎖 saved-role mutation。
- [x] 1.2 [after: 1.1] 先以 repo-relative 正式測試重現 instant-only 缺名稱提示，再 RED→GREEN→Refactor 實現 `Counter-dependent names have visible prerequisites`：instant-only counter 缺名稱先顯示 inline 必填與 disabled 理由，一般非計數瞬時名稱仍選填；先測兩分支與 stable circuit ID，再實作並通過 component tests。

## 2. 選擇範圍與能力

- [x] 2.1 [after: 1.2] 先以 repo-relative 正式測試重現新累積全選，再 RED→GREEN→Refactor 完成「累積來源明確選擇與查詢重試」的 `New cumulative selection is explicit`：新 cumulative_keys 為空且引導 picker，既存選擇與關閉重開不丟失，不連帶選 disabled profile；先測混合 numeric/text 及 saved scope，再實作、通過並重構。
- [x] 2.2 [after: 2.1] 先以 repo-relative 正式測試重現 instant-only lookup failure 缺 retry，再 RED→GREEN→Refactor 交付 `Every selected counter exposes measurement-read recovery`，讓 instant-only counter_raw lookup failure 也顯示 error/retry；測試 retry 保留草稿、無 save/start、晚回覆不寫入新 scope，修正後 focused Vitest 通過。
- [x] 2.3 [after: 2.2] 交付「布局能力與不支援設定」與 `Destination layout choices are capability bound`，將同 connector revision 的能力接入 UI 及 builder；先測 MySQL、MariaDB、SQL Server wide 不可送、stale/unknown capability 不放行、舊值可讀不自動轉換，再實作並跑 component/service tests。

## 3. 試寫復原

- [x] 3.1 [after: 2.3] 等待 design「本案檔案分工與開工條件」的 backend／lifecycle 交接後，以失敗測試固定「試寫 operation 與 editor 分離」及 `Unresolved test writes remain independently recoverable`，覆蓋 pending/running、lost response、新 preview、scope/revision change、reload；先確認舊 state reset 可重現，再新增 non-secret scoped operation references，通過後重構。
- [x] 3.2 [after: 3.1] 由 parent 另派且取得精確檔案 ownership 的 backend worker 確認並補齊 workspace-scoped server ledger recovery query，讓 local-reference 遺失仍找回未決 operation；用 API ownership tests 證明 foreign workspace 不可讀、safe error 不含 endpoint/credential，前端 integration tests 證明讀取失敗不被當成未寫入。
- [x] 3.3 [after: 3.2] 讓新 preview 對同 group 未決 write/cleanup 維持阻擋，write/readback/cleanup 分開顯示；先測 cleanup failed/unknown 與 supported same-operation repair 不重插資料，再實作並通過 mutation call-count／operation-ID integration tests。

## 4. 驗收與交接

- [x] 4.1 [after: 3.3] 在 freeze→cutover→lifecycle 前置實際通過且共享檔案移交後，依「已保存角色鎖定與支援的 lifecycle 導覽」整合新舊 Basic/Advanced 與既有 lifecycle host 支援的 review／導覽；saved role 維持鎖定並明示角色變更不支援，不將布局／目的地 review 當角色編輯、不新增 ensure/update fallback；測 saved role／frozen binding 不變，依 design 的 TestMeanLongRecording SQLite fixture 驗證新 authored numeric state 1→2 得 latest 2、gauge mean 1.5（不需真設備／新 DB 環境），執行相關 frontend targeted suites、typecheck、lint/build；若新增 API 則執行 handlers/router tests 並重產 Swagger，確認 capability／recovery shape 與 UI 使用一致。
- [x] 4.2 [after: 4.1] 用 browser tests 與截圖覆蓋 numeric state、缺名稱、累積 picker、instant-only retry、每種目的地布局 gate、unknown test-write reload/Back/切設備；核對後端 frozen-binding 前置與已保存語意未被導覽改寫，各 spec scenario 記錄實際結果。 依「現有雲端瀏覽器驗收的條件式交付例外」驗證 `Browser verification exceptions remain explicit`：保留正式 cases／discovery／E2E 專用 typecheck／lint／scenario 對照；僅當既有 cloud blockers 未變且實際執行仍不可用，實際 browser／截圖標 NOT RUN／SKIP／未驗證（不可標 PASS），其餘原驗收照跑，不新增環境或繞過安全限制。
- [x] 4.3 [after: 4.2] 將復原指引與已知不支援組合寫入適用操作文件，執行 git diff --check、make check-lines 與 spec scenario-to-test 對照審查；只記錄實際通過項目，未執行 DB/現場驗收保持未完成。

## 待辦對應與執行界線（非完成紀錄）

- 原 1.1 的四項 RED 證據拆回新 1.1（numeric state）、1.2（名稱）、2.1（新累積選擇）、2.2（lookup retry）；每個切片自行 RED→GREEN，不留下其他已啟用失敗測試。原 1.2 → 新 1.1，原 1.3 → 新 1.2；原 1.2 的 saved-role lifecycle 整合保留在 4.1–4.2。
- 原 2.1–4.3 各項保留原編號及行為；只補上準備／接線門檻並修正 2.1 的 after。沒有取消或完成任何產品要求；原始任務全部未完成，沒有既有 [P] 或完成 provenance 需改寫。
- 1.1–2.3 可在既有契約上循序先行；3.1 以後須通過 design 的依賴及檔案交接。fixture／純前端通過不代替真 server listing、saved-role review、SQL 或 browser 最終驗收。
