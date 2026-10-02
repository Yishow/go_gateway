# Studio V2：從採集到資料庫的收斂計畫

2026-10-02 已授權依 **A→F 順序實作並驗證六案**；各 change 的 tasks 與實作驗證紀錄表示當前進度。保留 Go／React，在 `/studio/v2` 原地改善；不建立 V3、不恢復 `/studio`，也不部署。

## 凍結基本需求

使用者可以依序完成「連線 → 採集點位 → Point-to-Tag → 選取 Tags 分組寫資料庫」。分組的 destination、columns、週期、品質規則、row identity 及 revision 只由一份 WriteGroup authority 決定。managed／custom 是同一群組的 storage strategy。既有 Local Modbus 獨立輸出與 readiness 安全保護保留。

基本寫入不要求建立報表、retention政策、window aggregation 或區間電量。使用者明確選擇這些進階能力時才要求對應 measurement semantics；不得依名稱或單位猜測累積量／瞬時量。

## 實作順序與唯一負責範圍

| 順序 | Change | 唯一負責交付 | 必要前置 |
| --- | --- | --- | --- |
| A | [unify-studio-v2-write-group-contract](../../../openspec/changes/archive/2026-10-02-unify-studio-v2-write-group-contract/proposal.md) | 群組模型、相容adapter、原子保存、authority／切換介面 | 現有安全契約 |
| B | [enforce-write-group-sample-semantics](../../../openspec/changes/archive/2026-10-02-enforce-write-group-sample-semantics/proposal.md) | 原始時間/型別/品質、snapshot closure、row identity | A驗證 |
| C | [wire-durable-write-group-delivery](../../../openspec/changes/archive/2026-10-02-wire-durable-write-group-delivery/proposal.md) | production intake/journal/outbox/sender/restart及交付truth | B驗證 |
| D | [implement-confirmed-write-group-test-write](../../../openspec/changes/archive/2026-10-02-implement-confirmed-write-group-test-write/proposal.md) | confirm/write/readback/cleanup，接手舊3.3/3.4 | C驗證 |
| E | [simplify-studio-v2-four-step-setup](../../../openspec/changes/simplify-studio-v2-four-step-setup/proposal.md) | capability、前三步、群組Step4/readiness，接手舊2.3/4.1 | D驗證 |
| F | [validate-studio-v2-device-to-sql](../../../openspec/changes/validate-studio-v2-device-to-sql/proposal.md) | production UI→實際SQL及故障矩陣，接手舊4.3 | E驗證 |

字母只是本表順序，change名稱才是交接識別。依賴圖為 A→B→C→D→E→F，沒有互相等待。各案可先研究，實作 gate 必須核對前案tasks及測試證據；CLI artifact ready 不等於跨案已解鎖。共用文件及程式檔依此串行，避免多套writer或多份delta競爭同一需求。

## 與原有工作如何銜接

- [舊資料庫流程 change](../../../openspec/changes/fix-studio-v2-database-workflow/handoff.md) 的五項已完成任務與證據原封保留；五項未完需求改由D/E/F完成，沒有勾成完成
- `fix-studio-v2-database-result-truthfulness`、`implement-studio-v2-verified-schema-setup` 已歸檔；沿用其實際schema metadata/preview/confirm/operation保護，不重做第二套ledger
- 9月的 measurement、recording、aggregation、delivery與history specs是既有契約；[現況證據](evidence.md)區分component/API與production wiring。這批只補基本grouped-writing缺口，不藉此宣稱那些能力全數端到端完成
- canonical specs 不在本次直接改寫；新需求留delta，後續實作和verify/review完成才考慮archive

## 交接資料

- [Source基準與缺口](evidence.md)
- [跨案資料契約](contracts.md)
- [預定驗收矩陣](acceptance.md)
- [本次文件檢查結果](validation.md)
- [A 的設定／移轉／回復說明](../../technical/studio-v2-write-groups.md)
- [A 的實作驗證紀錄](../../../openspec/changes/archive/2026-10-02-unify-studio-v2-write-group-contract/validation.md)

原提案建立時，新增 tasks 全未勾選，測試名稱是待實作驗收標籤；現在只按實際實作與驗證結果勾選，不因 artifacts 齊全表示產品完成。提案起草日期2026-10-01、基準main `4c00af81fb51b3239b4e0d0d9c6b83680c808b36`；本輪實作起始 SHA 與實際證據記於各 change 的 validation。
