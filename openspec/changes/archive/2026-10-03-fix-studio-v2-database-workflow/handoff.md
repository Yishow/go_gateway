# 2026-10-01：未完成工作移交

原基準main `4c00af81fb51b3239b4e0d0d9c6b83680c808b36` 的tasks有5項完成、5項未完成。此輪僅整理提案文件，沒有執行產品修正，也沒有把任何未完成項勾選。

## 唯一接手者

| 原任務／需求 | 新owner與任務 | 原需求如何保留 |
| --- | --- | --- |
| 2.3 real metadata／assignment／shared-column驗收 | `simplify-studio-v2-four-step-setup` 2.3 | 完整移轉Tag-to-column auto-assignment及Column conflict detection delta |
| 3.3 真實idempotent test-write | `implement-confirmed-write-group-test-write` 1.1–2.1、3.1–3.2 | 完整移轉Truthful explicit test writes，擴充canonical group及production sender |
| 3.4 readback／cleanup | 同上 2.2–2.3 | 原七個scenario保留，另補typed payload／owned cleanup／restart |
| 4.1 引導與可用性 | `simplify-studio-v2-four-step-setup` 2.1–3.2 | 完整移轉Guided database setup and accessible controls；Connector configuration form轉交新群組editor維護 |
| 4.3 整合與browser／DB驗收 | `validate-studio-v2-device-to-sql` 1.1–3.2 | 原partial failure、reload、disposable DB、full checks與evidence要求保留，增加production SQL鏈 |

D/E/F前置為A `unify-studio-v2-write-group-contract` → B `enforce-write-group-sample-semantics` → C `wire-durable-write-group-delivery`。完整順序 [見總覽](../../../../docs/plans/studio-v2-write-groups/README.md)。

## 保留與不變的證據

- 已完成1.2、1.3、2.1、2.2、2.4的checkbox和原任務文字完整保留，對應validation及source沒有在此輪重新認證
- 剩餘delta保留Protected row editing、Fresh schema results、Credential and path safety、Persisted recording membership、Consistent database setup persistence
- Persisted recording membership只約束選用原進階recording plan的情境；新basic WriteGroup不需fabricated measurement，也不改寫既有advanced plan的語意
- 移交需求從本案可執行delta移除，避免兩案修改同一requirement。未完成任務原始版本由上述基準Git保存，不以刪除checkbox表示完成
- schema confirmation、workspace/settings revisions、readiness token、scope ownership與Share protection全保留；schema setup前置已落地，不建立第二套ledger

## 如何使用舊design與proposal

原design的Decisions與A/B/C/D Implementation Contract記錄舊分階段計畫。2026-10-01起，未完D階段的實作owner與順序以本文件及新案design為準；已完成B/C及歷史證據仍有效，不能重新套用舊順序造成雙重實作。

本案留未勾選5.1 closeout gate，須在新案D/E/F和前置全部有實際驗收後才核對；不在文件交付時archive。新案尚未實作時，本案仍未完成。

## 2026-10-03 實作後核對

A–E 已依序完成並歸檔；F 已有真 UI／SQL／故障與 Linux container 證據，最後 reviewer 發現的 harness cleanup 與 late group save 草稿問題已修正，受影響 full runs 繼續核對。原五項完成文字保留；移交的五項以新 owner 實作與證據判定，5.1 在最終結果齊全前保持未完。逐項對照及 production migration／rollback 未執行邊界見 [handoff-closeout.md](../../../../docs/plans/studio-v2-write-groups/handoff-closeout.md)。


## 2026-10-03 完成移交

最後 review／verify、完整基準與受影響重跑皆通過；A–F 全部需求依核定範圍完成，5.1 現可按實際證據完成。原五項完成文字未變，五项移交結果與未驗界線見 [final-verification.md](../../../../docs/plans/studio-v2-write-groups/final-verification.md)。archive由Spectra preview/core transaction處理，沒有部署。
