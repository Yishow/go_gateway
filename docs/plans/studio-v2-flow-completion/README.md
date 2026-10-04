# Studio V2：流程收斂與既有缺陷修復

## 交付範圍

`bf463e38` 的原始交付為 **OpenSpec 起草＋一次 review**；當時六案共 50 項 tasks 保持未完成，沒有 archive 或產品實作。先前受阻的提交嘗試記錄於 review，並非產品驗證結果；原始交付記錄保留。

2026-10-04 已依後續授權逐案完成 A–E 的 apply、verify、review、archive 與精準提交；F 的 fresh UI／SQL／恢復驗收矩陣已執行完畢，待 verify／review／archive 與提交。以 `6676b2b0` 為程式基準，既有修復不重做；必要修補與原始證據見 [review-repairs.md](review-repairs.md)，各案實際結果見下表。沒有 push 或部署。

| 已交付案 | 提交 | 實際驗證 |
| --- | --- | --- |
| A | `eb30a21f0e7aff7345ea04292bc8d286fe9c66db` | [生命週期](A-lifecycle-verification.md) |
| B | `c26c8a52a1712d922015ba2d3822263c76a5dd7d` | [交付期限](B-deadline-verification.md) |
| C | `b560926c88ce4b477f53ec73e20e5d721964ba81` | [Entity 列](C-entity-verification.md) |
| D | `38c023e403efefc7601dbf8628d2fd8e8627b86f` | [受管理儲存](D-managed-verification.md) |
| E | `f8882bb6264302027b9590e6ca1510dcd05dac83` | [四步設定](E-setup-verification.md) |
| F | 尚未提交 | [空白目的地驗收](F-first-recording-verification.md) |

唯一產品目標：在既有 `/studio/v2`，快速設定設備與資料，透過確認式準備寫入 SQLite／PostgreSQL，看到可靠的真實記錄。不是再建立一個平台。

## 基準與證據等級

- 起點：`d16f1294276e9569dc2dae54575ae89303a54b42`，包含該 commit 的規劃脈絡；其後產品變更以差異及最新 source 檢視。
- 本輪讀取 main：`d2c22ef625903264ec6d74f6fd4bf76c5a48f916`，與上一輪分析相同；提交前須再查 main，若前進就重新對齊，不 force push。
- 該 main 的 tracked `openspec/changes/` 僅有 archive。未連線開發機器的 untracked/parked/WIP 不可視為已檢查或擅自移動。
- 上輪 findings 是 source review；deadline 與 until 有隔離重現，不等於本輪專案整合通過。各 owner 必須先在正式接線重現，再修正並驗收；若新證據推翻問題，記錄差異而不是硬改。
- 既有六案的 [最終驗收](../studio-v2-write-groups/final-verification.md) 與 [技術合約](../../technical/studio-v2-write-groups.md) 保留，不能把新提案算成它們已驗過的能力。

## 六案與順序

| 次序 | Change | 唯一責任 | 前置 |
| --- | --- | --- | --- |
| A | [fix-write-group-runtime-lifecycle](../../../openspec/changes/archive/2026-10-04-fix-write-group-runtime-lifecycle/proposal.md) | 草稿不停舊版、cutoff/race、離線冷啟動及歷史 journal 排空 | 無 |
| B | [fix-write-group-delivery-deadlines](../../../openspec/changes/archive/2026-10-04-fix-write-group-delivery-deadlines/proposal.md) | 遠端完成後才起算本地結果保存期限 | 無 |
| C | [fix-write-group-entity-row-layout](../../../openspec/changes/archive/2026-10-04-fix-write-group-entity-row-layout/proposal.md) | 同群組多 entity 的欄位、readiness、編碼一致 | 無 |
| D | [complete-write-group-managed-storage](../../../openspec/changes/archive/2026-10-04-complete-write-group-managed-storage/proposal.md) | canonical group 受控建表、必要身分/時間/品質、既有 receipt | A、B；C 是多 entity 正式驗收前置 |
| E | [streamline-studio-v2-recording-setup](../../../openspec/changes/archive/2026-10-04-streamline-studio-v2-recording-setup/proposal.md) | 四步基本操作、少輸入、開始記錄協調、交付事實 | D，包含其前置 |
| F | [validate-studio-v2-first-recording](../../../openspec/changes/validate-studio-v2-first-recording/proposal.md) | 真 UI 空白目的地與完整故障組合驗收 | E，包含其前置 |

先完成 A/B 的必要採集與交付保護，接 D/E 的單設備 SQLite 垂直流程；F 的初步 UI→首列觀察隨此流程進行，不等整批 50 tasks 全完成。C 可獨立修復，與 D 同檔修改按 ownership 整合；多 entity、PostgreSQL 與完整故障矩陣仍在正式 release/archive 前完成。D 不重做 A 的恢復描述，E 不重做 D 的 schema engine，F 不擴充產品功能。每案自己的回歸、文件與 focused tests 隨該案完成。

早期回饋的最小見證：全新 workspace、單設備、空白 SQLite，真 UI 完成四步→schema 預覽與明確確認→開始→獨立 SQL 首列。它回答基本流程是否可用，不能宣稱 A-F 完成或代替正式驗收。先看使用者真正需要輸入哪些資料、哪裡停住、首桶等待是否可理解，再沿同一主線補 PostgreSQL 與相容情境；不另建「快速模式」。

## 上輪問題到唯一 owner

| 問題 | Owner | 必須看到的結果 |
| --- | --- | --- |
| R1 儲存草稿影響執行中版本 | A | 未 Apply 跨兩桶及重啟仍用舊 applied membership |
| R2 settlement 提早到期 | B | 慢成功／失敗都得到完整的後置本地保存時間 |
| R3 離線重啟無法接住新樣本 | A | 有 verified 本地描述的版本可繼續 ACK，交付等恢復 |
| R4 停用／切版後歷史 journal 無人排空 | A | Disable/Delete/supersede 後重啟只排空原版，不再接新資料 |
| R5 同群組多 entity 被錯擋或錯誤 skipped | C | 合法兩 entity 實際各自落庫，不用兩群組替代 |
| R6 until 並行讀寫 | A | 真正受影響 package 的 race 回歸通過 |
| U1 空白資料庫沒有完整 UI 建表路徑 | D | 無人工 SQL、無 advanced plan，確認後才建 owned 表 |
| U2 四個畫面卻要求太多技術決策 | E | 基本映射／群組由系統保存，revision 不需使用者操作 |
| U3 snapshot 等待與缺值行為不清楚 | E | 直接顯示週期/規則與等待事實，不偷改桶或補值 |
| U4 SQL 缺可靠時間與來源資訊 | D | 新基本 managed 表保留原時間、record/group/device 與品質 |
| 整體是否真的簡單而可靠 | F | 空目標 UI→SQL、三桶、三次計時與故障矩陣 |

## 不可擴大範圍

保留 Go/Gin、React/TypeScript、既有四步、單一 WriteGroup 設定來源、exact codec、journal/outbox、不可變版本、destination ownership、CAS/readiness/token 與 Share 獨立保護。SQLite/PostgreSQL 是本批目的地；只用目前已支援的設備採集能力。

不做：V3、完整重寫、換框架、新協議、設備自動探索、新 CSV 匯入器、MQTT 打通、MySQL/MSSQL 新支援、報表/保留政策/聚合/用量、跨設備同步或 join、新通用 job engine、第二套 queue/ledger/token、效能平台、全域 dashboard 改版或正式環境操作。不得順手整理全部舊文件或刪除相容資料。

每台設備一個 managed group 只適用 **新基本設定**。既有 custom、多 entity、跨設備設定不自動拆分、轉換或抹除；C 仍必須修好目前合法的多 entity。原始資料型別、單位與位址不得猜測，未知值不得補 0／good。

受控建表只做明確確認的新 owned 結構或已驗證 no-op；不改使用者既有表、建 PG 帳號/資料庫或搬資料。目的 SQLite 與 gateway 設定/journal DB 分離。custom all-good 表不被強迫新增 metadata/receipt。試寫與清理仍是診斷，不是開始採集的偽成功替代品。

若驗收發現新需求：先確認是否為上述結果不可或缺；範圍內 bug 回唯一 owner，新增能力另行取得授權，不塞進目前 tasks。必要 additive migration 與薄 API 是既定流程的實作細節，不授權一般化擴充。

## 使用者流程與完成標準

沿用連設備 → 選點 → 確認資料/必要 Point-to-Tag → 目的地與開始記錄。基本預設每設備一組、目前 60 秒 snapshot；週期明示可改，缺值遵守既有政策。Step 4 只有必要選擇，DDL 有獨立明確確認，後端協調開始動作但不偽裝跨資源原子交易。

必須從全新 workspace 與沒有採集表的目標開始，全部 setup mutation 經 UI；API/SQL 只作觀察，不能預先 CREATE 目標表或 INSERT 假讀值。SQLite 缺檔用 stat、已有檔用 read-only URI，不讓 fixture open 偷建目的檔；每次使用獨立 fresh file/schema。SQLite／PostgreSQL 各驗證七種既有型別、exact uint64（含 2^53+1、2^63、最大值）、時間、品質、身分、至少三個正式桶；受控首次使用各跑三次，從首開 setup 到首列 SELECT 的每次結果以 300 秒為門檻，並拆分 setup、完整桶等待與交付時間。自動化計時與操作者理解/操作證據分別記錄。前置環境與限制要記錄，此數字不是實測或現場保證。

跨案必驗 R1-R6、延遲補送原時間、重送/失去回覆/stale/部分啟動、custom 與 Share-only。沿用現有容量、poison、fencing、unknown、cleanup 回歸。顯示已採集／本地已接受／SQL committed，實際讀回才叫 verified。任一必驗情境 blocked/NOT RUN 就不能宣稱整批實作完成。

## Source anchors

以下均為本輪基準的 repo-relative 原始碼定位，不是已執行結果：

| 主題 | 來源 |
| --- | --- |
| 正式接線 | `cmd/test_ui/service_wiring.go` |
| 草稿/執行/恢復 | `internal/datalink/workspace/write_group_service.go`、`write_group_lifecycle.go`；`internal/datalink/grouppipeline/reconcile.go` |
| 時間邊界/ledger | `internal/datalink/runtime/group_boundary.go` |
| 交付期限 | `internal/datalink/groupdelivery/sender.go` |
| entity 列 | `internal/datalink/snapshot/assembler.go`；`internal/datalink/dbtarget/group_row_layout.go` |
| UI 建表缺口 | `frontend/src/features/datalink/workbench-v2/steps/step4/SchemaSetupSection.tsx`、`writeGroup/GroupColumnProposal.tsx` |
| 舊設定耦合 | 同目錄 `Step4Database.tsx`；`writeGroup/GroupEditor.tsx` |
| 預設與 metadata | `frontend/src/features/datalink/workbench-v2/state/writeGroup/draft.ts`；`internal/datalink/dbtarget/group_row_insert.go` |
| 已有 schema/operation | `internal/api/router_studio_v2_recording_routes.go`；`internal/datalink/recordingplan/` |
| 既有驗收起點 | `scripts/tests/f_device_to_sql/run.mjs` |

Delta 沿用七個既有 capability，沒有新平行能力：`studio-v2-write-groups`、`runtime-write-group-delivery`、`write-group-row-semantics`、`recording-database-setup`、`guided-recording-workflow`、`datalink-workbench-v2-step4-database`、`studio-v2-device-to-sql-acceptance`。同 capability 的跨案 changes 各修改不同 requirement，不互相取代；既有未改 requirement/scenario 保留。

## 本輪檢查與實作前 gate

原始起草 review 結果及當時限制見 [review.md](review.md)；該日 CLI 未能執行的紀錄保留。後續本機修補與實際驗證見 [review-repairs.md](review-repairs.md)，不改寫歷史結果。

原始實作前 gate 的要求為逐案 validate 再依上表順序操作；目前實際使用 repo-local Spectra skills 與真實 CLI。未執行的 CLI、全套 source gate、真 DB 或現場項目，不能靠 tasks 勾選、舊報告或模擬輸出補成 PASS。各案新結果由上表的實際驗證與 Git 歷史判定，原始起草交付不作產品完成證據。
