## Why

四步畫面仍把多次保存、Apply、revision 與不同設定模型交給使用者協調。只縮排畫面無法解決流程中斷；本案在既有 Studio V2 內，讓系統處理可推導設定與版本接線，使用者只確認真正必要的選擇。

## What Changes

- 保留四步及真實讀值；基本 Point-to-Tag 沿用已有生成/保存能力，不再要求重複建立相同映射。
- 新基本記錄預設每台設備一個 canonical managed 群組，已有群組不自動拆分或覆寫。
- 主畫面保留設備、所選資料、目的地、記錄週期與開始記錄；revision、dedupe、自訂欄位等移到進階。
- 開始記錄由一個狹窄後端協調入口串起既有儲存/驗證/Apply/activation；schema 仍是事先獨立明確確認，不偷偷 DDL。
- 完成區分已採集、本地已接受、SQL 已提交；試寫成功或 runtime running 不能代替持續記錄。

## Capabilities

### New Capabilities

無；沿用現有 capability，不建立第二套採集或交付模型。

### Modified Capabilities

- `guided-recording-workflow`
- `datalink-workbench-v2-step4-database`

## Impact

限定 `frontend/src/features/datalink/workbench-v2/`、相關 hydration/autosave/services/hooks/types/locales、既有 workspace activation 與 operation ledger 的薄協調；必要 runtime handoff 只承接該 workspace/group/device identity 與交付狀態，不改建 dashboard。

## Scope and Dependencies

不新增協議、CSV 匯入器、設備自動偵測、報表、聚合、V3、UI framework 或全域 job engine。保留 custom/advanced 編輯與離線草稿；不重寫 `/test`、`/gateway/*`，不改 Share gates。單設備基本開始記錄依賴 A/B/D 的正確性與安全契約；C 的多 entity 矩陣是多 entity 功能與正式 release/archive gate，不阻擋單設備 SQLite 的早期垂直回饋。

前置：`complete-write-group-managed-storage`。
本案只起草；產品 tasks 全部未完成。共同邊界與來源見 [總覽](../../../docs/plans/studio-v2-flow-completion/README.md)。
