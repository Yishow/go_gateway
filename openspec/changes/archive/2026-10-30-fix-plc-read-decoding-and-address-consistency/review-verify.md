# PLC 修復最終 Review／Verify

本案五項初次 Warning 已修復；最終品質審查沒有 active Critical、Warning 或 Suggestion。32 個 Scenario、5 個 Example／33 列均有測試對應，沒有 test-scope exclusion。實測狀態與範圍另記於 [validation.md](validation.md)，未把靜態對照視為測試通過。

## Review

使用 Spectra 的 change scope；基準為 `76a853d71fda23ee3a1d3cf9bf801e45f85f24f2`，涵蓋本次全部 66 個 source／test 檔案。前端與後端分別唯讀檢查 correctness、efficiency、simplification/reuse、convention；前端追加修復後針對異動內容重新審查。主代理回讀 source／diff、測試與真實 browser 證據，保留最終裁決。

- R1：新 mapping／database target 草稿不再因沒有 server record 被預設值取代。
- R2：cold MC locator 使用既存 mapping 的 point ID／address。
- R3：新規則預填遵循 MC radix，舊規則以實際已占用地址計算。
- R4：拒絕 24-bit 範圍外地址及展開溢位，並保留正常 error UI。
- R5：mutation 前 validation 不執行 restore writes；mutation 後的 validation failure 仍回復。
- Browser 額外確認的同案缺口：全字母 hex suffix 不誤當 area；compatible evidence 不掩蓋目前 invalid draft。

無效編輯收到 backend 400、顯示可行動錯誤並阻擋 Continue，不構成成功寫入；未新增「所有 invalid input 一律零 HTTP」政策。legacy evidence conflict／checking／unavailable 的既有 mutation gate 保持有效。

## Verify

| Requirement | 實作 | 測試對照／執行 |
| --- | --- | --- |
| Scalar point reads preserve protocol count units | connector count helper、collector periodic／manual、device read-now | 4 Scenario；Go passed-current |
| Multiword numeric decoding follows all configured byte orders | hsllogic byte transform／register converter | 3 Scenario、2 Example／8 列；Go passed-current |
| FATEK typed register decoding preserves complete values | FATEK adapter typed conversion／error propagation | 3 Scenario；Go passed-current |
| Explicit raw Modbus addresses preserve requested function | Modbus shared address parser | 5 Scenario、1 Example／8 列；Go passed-current |
| FATEK discrete payloads are validated before success | FATEK client Cmd44／adapter／三個 scalar caller | 3 Scenario；Go passed-current |
| Protocol-aware point address derivation | sourceRule projection／backend planner | 3 Scenario、1 Example／7 列；Go／TS passed-current |
| Protocol-radix source address validation and expansion | frontend addressParser／backend MC plan | 3 Scenario、1 Example／10 列；Go／TS passed-current |
| MC persisted source identities survive radix correction | MC guard／runtime reconcile validation boundary | 4 Scenario；Go passed-current |
| Persisted MC points remain authoritative in Step 2 | scoped evidence／projection／autosave／locator | 4 Scenario；TS passed-current |

完整 scenario／example assertion 對照見 [scenario-coverage.md](scenario-coverage.md)。spec-driven 四類 artifact 均完成；21 個 task 的 checkbox 與 Spectra list 另以最終交付檢查核對。

## 實測與限制

- Go：8 packages、835 個含 subtest 的 PASS 項目，API 2 packages／49 項，均 0 fail／skip。
- Frontend：29 files／209 tests；29 個 TS 檔 ESLint、TypeScript、Vite／Go embed build 通過。
- 受影響 packages vet／lint、MC boundary／reconcile focused race、diff／行數與 Spectra validate 檢查通過。
- Browser：最終 binary＋loopback MC simulator＋可丟棄 SQLite，實際驗證 X10／W11 預填、Word label、兩種 overflow 的 error／blocked UI；[去敏結果](evidence/followup-browser-verification.json)。
- 真 PLC、正式資料庫、Windows 與現場驗收：NOT RUN。

範圍只包含原五類 PLC 修復及其必要 caller／回歸；四個 runtime-dashboard WIP 排除。沒有 push、部署、歷史回填或使用者資料變更。封存與提交仍依完整內容計畫確認後執行。
