# 設計：沒有證據，就不顯示完成

## Context
提案時已確認 `StudioV2WorkspaceRecordingPlansHandler.SchemaApply` 只驗證非空 token 後回覆 applied，`TestWrite` 只驗證 plan_id 後回覆 written_verified。一般 `StudioV2WorkspaceDatabaseHandler.GenerateSchema` 則有呼叫真正服務，兩條路徑不可混淆。這是程式行為證據，不是現場資料遺失證據；本次修正與驗證結果見 [validation.md](validation.md)。

`CommitProgress` 現為 logs/status 顯示，`CommitSuccessCard` 已用 activation response；固定十段計時是主規格的過時要求，不宣稱目前仍有該計時器。空 results 的正常採集與儲存文案須有後端證據。

## Goals / Non-Goals
In scope：記錄方案未實作入口的安全封鎖、能力與錯誤格式、Step 4 完成判定、API 文件與回歸測試。
Out of scope：真正資料庫操作、migration、連線權限制度重寫、PLC 協議、整個 runtime 重構。

## Decisions
### 1. 先修後端，再修畫面
未實作的合法要求回 HTTP 501，`success:false`，不得帶 `applied:true`、`written_verified` 或 fabricated record_id。格式錯誤維持既有驗證拒絕；規格不將所有錯誤改成 501。不存在或不屬於工作區的資源，在未實作階段也不能取得成功；`fix-studio-v2-database-workflow` 2.2 完成後由共用 scope 檢查提供一致的 404。

固定 codes 為 `RECORDING_SCHEMA_NOT_IMPLEMENTED`、`RECORDING_TEST_WRITE_NOT_IMPLEMENTED`。錯誤帶 code、可翻譯安全 message、retryable=false、action=`wait_for_supported_operation`、不含帳密的 request_id。機器判斷依 code，不能比對文案。

未實作錯誤 envelope 的 code 只允許上述兩個 `RECORDING_*_NOT_IMPLEMENTED`，action 只允許 `wait_for_supported_operation`；不得放寬任意原始 code 或 action。測試須在 normalize／translation 後 assert code、action、retryable=false 與完整 request_id，raw diagnostics 不得進 DOM。

### 2. 能力以「目前可執行的操作」為準
`supported` 不可代表所有操作已實作。全域能力 mask 僅取 adapter kind 能力與已接線 operation 的交集；本案兩個未接線操作固定為 false。`supports_managed_schema` 與 `supports_test_writes` 要同時滿足 adapter 能力及已接線服務。只有試寫未實作時，不得禁用已有的安全查詢功能。connector scope 仍由 workflow／C3 負責，不在本案擴充 scope API。畫面先說明原因，舊版畫面直接呼叫也由後端拒絕。

### 3. 成功狀態與網路結果分離
只有伺服器確認的動作能成為完成項目。計時器只能控制顯示，不產生建表、排程啟動或第一筆寫入的成功證據。HTTP 200 本身不足以算完成，必要欄位缺失要視為回覆格式錯誤。未知列舉值顯示「尚未確認」，不可落入成功分支。

逾時、斷線或使用者離開頁面只代表沒有收到確定答案，顯示「結果尚未確認」。未實作入口的 501 是確定未執行；真實操作的未知結果則交 `implement-studio-v2-verified-schema-setup` 與 `fix-studio-v2-database-workflow` 的 operation_id 查詢，不得換新編號自動重送。

### 4. 完成卡與 runtime 交接
設定已存、設定已套用、設備採集中、資料已交付及已讀回核對是不同證據。部分設備成功時保留成功設備的導航，同時保留失敗摘要，不能宣告全部成功。單一確定設備延續 `/studio/runtime?device_id=...`，不能確定單一設備時延續 `/studio/runtime`，不得改成無作用按鈕。第一筆寫入時間不得只由 write_interval_seconds 推算成承諾。

Remount recovery 讀取既有後端 workspace/runtime 狀態，不重送 activation。若沒有本次 workspace revision 的完成證據，狀態保持未確認；operation_id 只在既有回覆提供時保留，不新建 activation ledger 或重寫 runtime。

## Implementation Contract

1.1–1.3 在真正 Gin router 驗證：缺 body／欄位保留驗證拒絕；合法的未接線要求各回固定 501、完整安全錯誤、零 target mutation。能力與 handler 可用性一致，已有 custom-table 查詢／建表不被一併停用。

1.4–1.6 驗證 501 不重試、未知或缺欄回覆不成功、斷線保留未確認狀態；完成判定保留 workspace_revision、settings_revision、readiness_token 與 Modbus Share guards。部分設備、原有設備仍運作與重載不能證明本次設定成功，runtime 導航依已確認設備保留。

1.7–1.8 以具名回歸及 AGENTS.md 的前後端基準驗證，產生 release／rollback 紀錄。範圍限本案安全封鎖與結果判定，不能把 501 封鎖當真實建表／試寫或實機驗收。

## Contract and Compatibility
此案改變的是假成功回覆，不更改 request URL。舊客戶端收到 501 會進入失敗流程，這是刻意的安全相容性中斷。Release note 須明示；不可為舊畫面保留假的 200。上線後 rollback 只能退回同樣拒絕未實作操作的版本。

## Risks and Mitigations
若能力列舉與實際 handler 不一致，按鈕可能可見卻不可用：加 router-to-capability 合約測試。若只修 UI，直接 API 仍假成功：測試必須經真正 Gin router。若把全部 database output 關掉，會傷到現有功能：保留一般建表與 Local Modbus-only 的回歸。

## Verification and Rollout
先寫會在原碼失敗的測試，再修 handler 與 UI，再重整重複程式。驗收包括任意非空 token、不存在方案、重複點擊、未知成功值、錯誤缺欄、部分啟用、一般建表與 Modbus-only 未受影響。數量與效能不得引用上一輪 patch 的結果。

2026-09-15 已對本次授權後的最終產品版本執行測試、審查與規格對照；結果、未執行項目、已解除的 Go lint 基線阻擋及安全 rollback 均記錄於 [validation.md](validation.md)。未取得本次 revision 完成證據的 recovery 維持未確認；既有 runtime-context 的 running 可能來自已存設定，不能據此宣稱本次採集或資料庫交付完成。
