# Studio v2 資料庫規格：發布狀態

基準 main：`1a0311c8e8db9c62fe4388f0701ba38afe552ff7`

分支：`spec/studio-v2-database-hardening`

本次只新增 OpenSpec 草稿，沒有實作、套用上一輪 patch、修改主規格、執行資料庫變更或部署。

## 2026-09-15 本機 review 更新

目前本機 `main` 基準為 `569da9f98e52b2635638d85250a4af16601d4364`；此狀態已包含以下三個 active change，未查詢遠端 main 是否有更新。下面的「初稿發布紀錄」保留當時發布結果與限制，不是目前的實作前置清單。

| 階段 | 已存在的 change／待辦 | 負責內容 |
| --- | --- | --- |
| A | `fix-studio-v2-database-result-truthfulness` 1.1–1.8 | 安全封鎖、能力／錯誤、真實進度與完成卡 |
| B | `fix-studio-v2-database-workflow` 1.2、1.3、2.1、2.2、2.4 | 編輯／readonly、已存連線及成員歸屬、版本與原子儲存 |
| C | `implement-studio-v2-verified-schema-setup` 3.1–3.10 | 真實欄位、持久預覽及操作、確認建表、入口防旁路 |
| D | `fix-studio-v2-database-workflow` 2.3、3.3、3.4、4.1、4.3 | 配對、真實試寫／讀回／清理、引導介面與整合驗收 |

**Supersedes**：原 C1→C2→C3→C4→C5 發布安排由以上 A→B→C→D 取代。未發布的 C2/C4/C5 沒有當作已完成或已讀來源；其對應工作以 repo 中 workflow 的現行規格為準。依賴是階段與待辦，不是要求 workflow 整案先完成；每階段切換須檢查前置完成及測試證據。共享程式檔串行修改，主規格的每條需求只有一份 active delta 負責。

本次只修文件，沒有產品實作、資料庫操作、部署或 Git commit。正式文件檢查與限制見 [review checkpoint](../../../openspec/changes/fix-studio-v2-database-workflow/validation.md)。

## 初稿發布紀錄（歷史）

## 本分支包含

| 代號 | Change | 草稿內容 |
| --- | --- | --- |
| C1 | [fix-studio-v2-database-result-truthfulness](../../../openspec/changes/fix-studio-v2-database-result-truthfulness/proposal.md) | 停止未執行卻回報成功；結果、能力、錯誤與完成卡的判定 |
| C3 | [implement-studio-v2-verified-schema-setup](../../../openspec/changes/implement-studio-v2-verified-schema-setup/proposal.md) | 真實欄位查詢、版本綁定預覽、明確確認後建表、部分完成與未知結果 |

每案均含 proposal.md、design.md、tasks.md、.openspec.yaml 及 specs/*/spec.md。兩案共 8 條 requirements、31 個 scenarios、18 個未完成 implementation tasks；加上本文件共 12 個新增檔案。

## 已起草但未發布到 GitHub

以下三案的 GitHub create_tree 寫入被平台安全檢查阻擋，未改用其他方式繞過。本次對話提供的完整草稿附件保留這三案；它們不在本分支，不得當作已提交。

- C2：`fix-studio-v2-database-scope-and-persistence`
- C4：`implement-studio-v2-verified-test-writes`
- C5：`improve-studio-v2-database-guided-experience`

完整草稿包合計五案、24 條 requirements、97 個 scenarios、47 個待執行 tasks，另有來源／驗收對照與文件檢查紀錄。這些是文件計數，不是測試通過數。

## 相依與完成界線

預定整合順序為 C1 → C2 → C3 → C4 → C5。**C3 依賴尚未發布的 C2，不能略過前置規格直接開始。** C1 可單獨審閱；完成安全封鎖不等於真實建表或試寫完成。

本機文件副本已檢查結構、情境格式、task 依賴、重複 requirement 與空白；已發布兩案的 Git tree 識別碼與檢查過的本機副本一致。環境沒有 OpenSpec/Spectra CLI，因此沒有官方 analyze/validate/park 的成功證據。所有產品測試、真資料庫及瀏覽器驗收都未執行。

這是部分發布的規格草稿，**不是五案皆已在 GitHub、官方驗證通過或 ready-to-apply**。接手時先核對最新 main、補齊尚未發布的規格與官方檢查，再取得實作授權。
