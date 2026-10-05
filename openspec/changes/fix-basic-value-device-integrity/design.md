## Context
基準 d41638b 已接 canonical WriteGroup、journal、outbox、receipt。此批修四步路徑中的值／身份錯誤，不重建 runtime。

## Goals / Non-Goals
**Goals:** 預設 mapping 保留來源型別、checked exact conversions、單設備顯示。
**Non-Goals:** 不新 protocol／衍生量／報表／event 多 tag 模型；不重建前端，不更新正式設備與資料庫；不自動改既存 float64 映射。

## Decisions
### 保留型別並拒絕失真
buildDefaultMapping 的 target_type 跟隨 point.data_type。executeCast 以 typed integer／十進位字串處理整數；有限 float 僅在精確且範圍合法時轉換。整數轉 float 亦拒絕 precision loss；不以飽和、截斷或零代替失敗。float32 的既有近似語意只允許明確浮點轉換，禁止非有限值。
### 依選取設備與 point 身份配對
LivePointsTable 接 selectedDeviceId；過濾 mapping 與 event，point_id 精確配對；若保留地址 fallback，必須同 device 且無歧義。metadata 不使用全域地址 key。切換時舊 SSE 不顯示為新設備。

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
