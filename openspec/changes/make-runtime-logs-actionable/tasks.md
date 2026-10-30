## 1. 擴充共用事件契約

- [x] 1.1 在 expose-runtime-recording-rules 的 1.1 精確 RuleEvidenceView／EntityLocator／stage／action 契約完成並審定後，依「事件只引用共用真相」建立 domain-event refs 與 producer matrix；逐項記錄 authority、emit boundary、identity、evidence，通過契約 review 及缺 identity 回 unavailable 的 parser 測試，不新增 delivery authority。
- [x] 1.2 [after: 1.1] 限獨占新 domain_projection 檔，以 Red-Green-Refactor 完成 Event projection preserves redaction and local access 與「安全白名單與本機門禁」：新增實際 fixed code/typed scalar 白名單函式，測 raw error、secret、endpoint、nested payload、超長字段被拒絕或 suppressed；既有 sink 對未知 code 仍 suppression，不提前修改 project/broker。此為安全模組準備；memory/file/SSE/DOM 與 IPv4/IPv6 local、LAN/proxy/origin 的完整整合驗收保留在 2.1、3.3、4.1。
- [x] 1.3 [after: 1.2] 限新 domain_coalescer 與注入非阻塞 admission/clock 的專屬測試，以 Red-Green-Refactor 完成 Repeated events remain bounded and nonblocking 與「有界合併與隔離」：首次／heartbeat-window summary／恢復、first/last／字串 occurrences 與 drop 分離；測同 key 重複、revision 切換、高 cardinality、count/byte cap、取消／關閉並執行 diagnostics race tests。

## 2. 接入真實執行結果

- [x] 2.1 [after: 1.3] 在 runtime change 最終交接、backend/lifecycle 前置及 parent 逐檔 owner 移交後，先接 projector/coalescer 至既有安全 admission 並驗證 memory/file/SSE redaction、caps/drop accounting、local guard 不回歸，再完成 Production events carry attributable evidence 的 source apply/reject producer，僅在已確定結果發出 identity／revision；以 source-rule reconciliation 的成功、拒絕、unknown scope 測試驗證，不把 HTTP 200 當 runtime 已套用。
- [x] 2.2 [after: 2.1] 接入 collector failure/recovery 的安全事件並保持原採集結果；以 simulator 斷線／恢復、不同設備同時錯誤與 logging admission full 測試確認 identity 正確且採集不被 event sink 阻塞。
- [x] 2.3 [after: 2.2] 接入 mapping failure/recovery 的安全事件，引用 point/tag/mapping revision 而非 value/raw transform；用轉換錯誤、恢復、過期 mapping 與 secret params 測試，確認無每 sample spam。
- [x] 2.4 [after: 2.3] 接入 delivery retry/blocked/quarantined/unknown/recovered 與 SQL first-confirmed/成功摘要，僅從 durable transition／匹配 receipt 發出；以可丟棄 DB 測 queued、ack uncertainty、receipt mismatch、重啟舊 receipt 不當新 commit 及 logging failure 不改交付結果；驗證 counter-replay-proof-*／recording-ownership-*／recording-handoff-proof-conflict／recording-cutover-finalized-overlap 消費原 reason，不修復資料或重送 finalized output。
- [x] 2.5 [after: 2.4] 接入 group intake applied/paused/retiring 與 quota refusal，維持 manage-recording-lifecycle-safely 的 mutation ownership；測 disable 後仍有 backlog、revision 邊界、quota 拒收與恢復，事件不能宣稱所有舊資料已送達。
- [x] 2.6 [after: 2.5] 接入 storage repair／cutover 的已確定 operation 結果，沿 ledger ID 對應原 group/connector revision；測 repair blocked、結果不確定、cutover 成功及歷史 schema recreation 不代表資料恢復，確保不從 log action 重送 DDL。

## 3. 可行動的瀏覽與定位

- [x] 3.1 [after: 2.6] 完成 Actionable filters retain event and transport identity 與「可行動篩選與歷史定位」的 typed workspace/device/rule/group/stage/actionable filters；測 snapshot/replay/live 同條件、cursor 字串精度、filter A→B 遲到取消、gap 後 unknown 不當 healthy，維持舊 level/source/q 行為。
- [x] 3.2 [after: 3.1] 完成 Events navigate to exact settings without mutations：事件呈現問題／影響／時間／次數／action，透過共用 resolver 開啟原 revision 規則或 Step 4；先寫 UI 測試再實作，涵蓋 renamed/deleted/foreign identity、resolver unavailable、無任何 retry/skip/activation/DDL request。
- [x] 3.3 [after: 3.2] 完成新舊 event renderer 相容與 en／zh-TW 可達操作，區分 connection/process/delivery 證據；測 hostile text、鍵盤展開、403 停 retry、pause/resume/clear、history gap 與遠端說明，確認不建立 process-log proxy。

## 4. 收斂與驗證

- [x] 4.1 [after: 3.3] 真 producer 已覆蓋後移除同一結果的重複泛用事件，保留 startup／HTTP／shutdown 舊事件；執行 cmd/test_ui、diagnostics、runtime、grouppipeline、groupdelivery 與 API 受影響測試及 race、frontend runtime-logs/service tests／lint／typecheck／build，確認每個 spec scenario 有具名測試且沒有 log endpoint 自我放大。
- [x] 4.2 [after: 4.1] 以 simulator／可丟棄 DB 跑 embedded browser 的錯誤→恢復→證據定位與 offline/reconnect/gap 驗收，保存截圖及 production identity 證據；依 docs/testing.md 完成最終回歸、make check-lines、git diff --check、Spectra analyze/validate，明列未測 Windows／現場，不改既有 add-windows-tray-runtime-logs 完成狀態。 依「現有雲端瀏覽器驗收的條件式交付例外」驗證 `Browser verification exceptions remain explicit`：保留正式 cases／discovery／E2E 專用 typecheck／lint／scenario 對照；僅當既有 cloud blockers 未變且實際執行仍不可用，實際 browser／截圖標 NOT RUN／SKIP／未驗證（不可標 PASS），其餘原驗收照跑，不新增環境或繞過安全限制。

## 待辦對應與執行界線（非完成紀錄）

- 原 1.1 僅將 whole-change 前置拆為 runtime 1.1 契約前置與 2.1 的 production 前置；producer matrix、identity/authority/evidence 要求不變。
- 原 1.2 的 fixed-code projection 留 1.2，新 2.1 接回真 admission／memory/file/SSE 與 guard 整合，DOM 保留 3.3/4.1；原 1.3 的 coalescer 留 1.3，真 writer 非阻塞與 drop accounting 再由 2.1–2.6/4.1 驗證。不是把 fixture-only 測試當成產品完成。
- 2.1–4.2 保留原編號、producer 範圍及全部最終驗收，僅在 2.1 增加準備成果接線核對；14 項 pending，無完成 provenance 或 [P] 改寫。
