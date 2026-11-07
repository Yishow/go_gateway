## Context

原始完整測試 JSON 記錄 524 個 skip；其中只有三個是無條件停用的 handler 測試，其餘包含外部資料庫、耗時測量、硬體或環境條件。SQL migration 已以 UNIQUE 限制輪詢群組名稱。

## Goals / Non-Goals

**Goals:** 恢復三個 handler 測試，對齊記憶體與 SQL 名稱完整性，讓一般測試不會自動探測真 PLC。

**Non-Goals:** 不新增名稱正規化、不更改 HTTP error status、不移除正常 fixture/平台 skip、不連真 PLC、不改 UI。

## Decisions

### 使用真實 loopback Modbus peer

設備測試使用既有 virtual Modbus server 的動態埠與獨立 ConnectionManager，不 mock 掉啟用的 probe。測試結束先關連線再停 peer。

### 名稱檢查與寫入持有同一把鎖

MemoryRepository 在 Create/Update 的現有 mutex 臨界區內掃描同名群組；Update 排除自身 ID。保持 SQL 已有的大小寫敏感字串比較，不在 Service 做有競態的先查後寫。

### 硬體驗收須明確選入

GATEWAY_REAL_DEVICE_INTEGRATION 必須恰為 1 才能進入原有實體設備測試；未設定時在任何 Dial 之前 skip。此開關不是取代真 PLC 個別操作授權。

### 區分過期 CAS 與損毀 lineage

資料庫 proof 測試原以 revision+1 模擬 corruption，但此值在服務既有契約屬 stale CAS、應先回 409。保留真正 corruption 的 500 斷言，以刪除 reserved generation 建立損毀 fixture，另加入 stale_lineage 回歸。不改 production CAS 判斷次序。

## Implementation Contract

- MemoryRepository.Create/Update 簽名不變；同名錯誤回傳時不修改 map 或原物件。API 沿用既有 500 error envelope，狀態改善不屬本次範圍。
- 兩個設備 handler 測試要求實際 create → activate → persisted active，與 activate → disable → persisted disabled。
- 三個舊 skip 測試搬至小型專用測試檔，避免再增加既有超長檔。
- 驗收包含先失敗的名稱測試、修復後受影響 package、race 選測、lint/vet 與 524 skip 分類。

## Risks / Trade-offs

- 外部 DB fixture 未執行仍須獨立列為 NOT RUN，不能以此修復宣稱全資料庫驗證完成。
- 記憶體名稱掃描為 O(n)，與其現有儲存定位相符；不新增索引與同步狀態。
