## Context

Step4Database 仍有 targets/row_groups 的建表前置，新 group editor 又有獨立草稿、readiness 與 Apply。現有四步與進階元件可重用，不需要重寫頁面。來源與目的見 proposal.md。

## Goals / Non-Goals

減少必要輸入與技術狀態，而不是移除一致性。此案只整合 A-D 的既有能力；本身不負責新增 schema engine、重試模式或自動推測設備型別。

## Decisions

1. 保留 Step 1 連線、Step 2 點位、Step 3 名稱/型別與 Point-to-Tag 確認、Step 4 目的地及開始記錄。Step 3 基本映射由現有 source-rule/tag/mapping service 建立真 persisted identity；不製造假 IDs 或重複 Tag。既有批次選取/模板照用，不新增匯入方式。
2. 新基本設定每台設備一個 managed group，create-once key 為 (workspace_id, device_id, canonical managed role)，由 backend 的持久唯一性與同一 local transaction 保證，不用 display name／列順序或 browser memory。相同意圖雙擊／reload 重用原 ID；同設備換 connector/table 是明確 draft/CAS 變更，經重新 preview/confirm/Apply 後才切換，舊 accepted payload 不移轉。既有 advanced 群組不被當成可轉換的基本群組；相同設備需第二目的地走既有進階明確建組。使用持久來源身分與已保存 draft 做 create-once/reuse，重新開頁或雙擊不得多建群組。既有 custom、多 entity 或跨設備群組仍讀取原設定並走進階，不自動拆分、切換 storage strategy 或搶寫入 ownership。
3. 主畫面直接顯示記錄週期、snapshot 語意、缺值政策與下一個可記錄區間。沿用目前 60 秒新草稿預設與 UTC 完整桶規則，不為了快出首列而偷寫半桶；使用者可明確調整既有 interval。提供可行動的「等候資料／缺點位／等待補送」，而非假倒數成功。
4. 移除基本路徑對 legacy targets/row_groups 的必填依賴；canonical group 是設定真相。revision/digest/token 留在 service 與 diagnostics，非一般必填。自訂表真 metadata、欄位修復與既有試寫工具保留進階。
5. 既有 ActivateEligible 會讀整個 workspace；本案必須提供真正 scoped readiness/activation，不能只過濾回覆。選設備 A 不啟動 B，B 未完成的設定不阻擋 A，已有健康 B 不被停止或重投影；共享 Share 保護仍依既有契約驗證。開始記錄為狹窄的應用協調動作，不是通用工作流。由 backend 對同一選取意圖保存必要本地狀態、取得當前 revisions/readiness、Apply 指定群組及呼叫既有 activation。只接已完成 D 的明確 schema confirmation；缺結構回到準備階段，不自動 DDL。不由 browser 自行拼湊多個最新版本。
6. 沿用既有 operation ledger 記錄 start 的 action、workspace/group 範圍、request identity、意圖 digest、進度與安全結果；與 schema/test_write action 分開，不共用錯誤種類的 token。重送與失去第一個回覆用相同 request identity 回同一 operation，查明已完成步驟再接續；不得建立第二套 queue/job repository。
7. 外部結構已建、群組已 Apply 或部分設備已啟動時，不宣稱整體 rollback。重新整理能查進度，恢復只作用原選取 scope，使用者修改意圖後必須重新驗證。保存草稿不能觸發啟動；重試不停止其他健康設備。
8. 使用既有 delivery query/SSE/polling 呈現三段事實及最後 SQL commit 的群組/record 證據；只有實際讀回才稱 verified。完成 handoff 仍到 `/studio/runtime` 的 focused monitor，不新增儀表板或全量監控面板。

## Implementation Contract

- 行為：保留四步，開始只作用使用者選取的 device/group；未完成鄰近群組不阻擋，健康鄰近設備持續採集。
- 介面：start operation 固定保存 workspace/group/device set、原意圖 digest、connector/source/schema revisions、Apply 結果與 activation 進度；request identity 在同 scope 重送只能指向同 operation。
- 恢復：操作產生的新 applied revision/進度是該 operation 的已記錄結果；外部重新修改 intent/group/connector/schema 時 stale/revalidate，不以新 scope 繼續舊進度。不把全 workspace ActivateEligible 包裝成 scoped 成功。
- 唯一性：基本 managed role 由持久 (workspace_id, device_id, role) 身分重用；目的地改變不多建預設組、不移轉既有 backlog、不自動改 advanced 組。
- 失敗：save/readiness/Apply/activation 的階段與部分結果可查；未確認 schema 回準備區，不偷偷 DDL，也不宣稱跨資源 rollback。
- 驗收：A ready/B incomplete、A start/B healthy、雙擊/reload/lost response/restart、換目的地與同名設備測試；真 UI→首筆 SQL 與等待/queued/committed 分開。
- 範圍：既有設定與操作接線；不新增 workflow/job/queue/token/ledger、協議、模式或 dashboard。

## Risks / Trade-offs

- 把進階欄位藏起來會失去修復入口：保留「詳細／進階」及確切 blocked action，不隱藏失敗。
- 一次啟動不是跨資源原子交易：同 operation 保存進度，明示部分結果與恢復方式。
- 既有複雜設定：只優化新基本路徑，不自動轉換或刪除相容資料。

## Migration Plan

沿用 persisted workspace/group identities。只拆除基本畫面對舊模型的耦合，不刪除仍被使用的 legacy API 或保存資料；需要的 ledger schema 只 additive。rollback 不回復舊備份覆蓋新 accepted data，也不把已建表／已提交效果當可撤銷。

## Verification

先寫互動與 coordinator 回歸：未設定、離線草稿、保存失敗、stale revision、double click、reload/lost response、部分啟動、advanced 保留與 Share-only。確認 390/768/1440 寬度、鍵盤、en/zh-TW。最後真 UI 的空白目的地驗收由 F 負責。
