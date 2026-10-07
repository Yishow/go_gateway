## Context
基準 d41638b 已接 canonical WriteGroup、journal、outbox、receipt。此批修四步路徑中的值／身份錯誤，不重建 runtime。

## Goals / Non-Goals
**Goals:** 預設 mapping 保留來源型別、checked exact conversions、單設備顯示。
**Non-Goals:** 不新 protocol／衍生量／報表／event 多 tag 模型；不重建前端，不更新正式設備與資料庫；不自動改既存 float64 映射。

## Decisions
### 保留型別並拒絕失真
buildDefaultMapping 在 neutral scale=1／offset=0 時 target_type 跟隨 point.data_type。executeCast 以 typed integer／十進位字串處理整數；有限 float 僅在精確且範圍合法時轉換。整數轉 float 亦拒絕 precision loss；不以飽和、截斷或零代替失敗。float32 的既有近似語意只允許明確浮點轉換，禁止非有限值。
### 保留合法縮放與舊步驟順序
**Supersedes**: fix-basic-value-device-integrity / 保留型別並拒絕失真
上述保留來源型別限 neutral；使用者明確 non-neutral numeric scale 預設 float64。新 workspace mapping 先 scale 再 checked cast 到宣告 target（包含同 source type 的整數 target）；scale 輸入大整數不可精確轉 float64、非法參數、非有限結果都拒絕。既存 pipeline 的已發布步驟順序保持原樣，不 migration；合法 int16×0.5+10 保持 float64，宣告 tag 型別與值一致。
### 依選取設備與 point 身份配對
LivePointsTable 接 selectedDeviceId；過濾 mapping 與 event，point_id 精確配對；移除 address fallback；metadata 不使用全域地址 key。切換時舊 SSE 不顯示為新設備。

### 保留合法縮放與舊步驟順序的確認邊界
只在操作員 workspace mapping Save 全部成功時，以 source revision／link／exact pipeline／tag type CAS 保存已確認 signature。source edit 或未確認 pipeline edit 繼續 out_of_sync，不更新 accepted signature；同 rule 其他 mapping 不被順便接受，先同步其他 links、最後 CAS 當前 mapping。兩次衍生同步／重啟仍保持確認。來源變更後需既有 candidate 明確 Reapply，409 不猜測新來源。
### 保留型別並拒絕失真的即時顯示
既有 ValueEvent JSON 將超過 JavaScript safe integer 的 int64／uint64 輸出十進位字串，native sample 與 SQL 型別不改；不新 protocol。pipeline 失敗的第一 binding event quality=bad、transformed=null，保留 raw 診斷，不發布假 good。float underflow 至零拒絕。

## Implementation Contract
1. 預設 UI（包含 source uint64、無 scale）產生 uint64 target，workspace mapping 不插 float cast。production 鏈 journal／outbox／SQLite readback 對 9007199254740993 字元精確一致。
2. executeCast 對非法字串、負到 unsigned、越界、分數到整數、非有限數與不精確轉換回 error；pipeline quality／sample 機制阻止錯值當 good 寫入。需重用既有 exact codec／helpers 可行部分。
3. 同 workspace A.D0=215、B.D0=187：切 A 只出 A 名稱／值／型別單位，切 B 不暫借 215；缺 matching event 呈等待。
4. 先寫 focused Go／Vitest 失敗回歸，再實作與 review/fix/rereview。測試串行、Go -p 1、GOMAXPROCS=2，Vitest 單 worker。全量／build／browser 先向主對話协调；只 simulator 與 owned disposable DB。
5. scope 限 proposal 列出 modules 及直接 caller／測試。執行 git diff --check、line gate、affected tests；repo required gates 與 UI 實測未執行時明列限制；最後 aggregate exact HEAD review。

## Risks / Trade-offs
- 嚴格 cast 可能讓舊非法資料失敗 → 邊界 matrix 與 pipeline quality 回歸；不默默修補值。
- 舊 SSE event 缺 device identity → 拒絕不明身份，測切換 waiting 狀態。
- 16GiB Mac 與同機 Qycms → 不啟 Docker／大型 build／browser／全量 tests，重驗證先協調。
