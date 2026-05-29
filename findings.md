# 發現與決策

## 需求
- 依序實作 `datalink-workbench-v2-shell` 等 6 個變更。
- 另開 Git branch 來進行開發。

## 研究發現
- `spectra list` 顯示所有 6 個 workbench 變更皆已在 `in-progress` 狀態，總任務數量如下：
  - `datalink-workbench-v2-shell`: 34 任務
  - `datalink-workbench-v2-step1-device`: 24 任務
  - `datalink-workbench-v2-step2-rule`: 24 任務
  - `datalink-workbench-v2-step3-mapping`: 22 任務
  - `datalink-workbench-v2-step4-database`: 25 任務
  - `datalink-workbench-v2-settings`: 23 任務

## 技術決策
| 決策 | 理由 |
|------|------|
| 新建分支 `feature/datalink-workbench-v2` | 提供乾淨的開發環境 |
| surface / API 維護文件採多檔案拆分 | `/studio` 主線、`/studio/v2`、runtime、gateway、test 的成熟度差異很大，若混在單一文件會難維護 |

## 遇到的問題
| 問題 | 解決方案 |
|------|---------|
| 點位高頻 stale 事件狂刷 `setLogs` 造成嚴重重繪卡頓與日誌爆量 | 實作 `pointStaleStatesRef` 僅在點位 stale 狀態「改變」的瞬間才記錄 transition log (並新增 recovered 復原狀態日誌)。 |
| 新日誌 prepend 到最上方，與診斷面板的「滾動置底」UX 機制相反 | 將新日誌修改為 append 到陣列最後，使時間排序由上到下遞增，完美咬合 auto-scroll-to-bottom。 |

## 新研究發現

- `/studio` 是目前唯一真正深度綁定 datalink persisted APIs 的主產品 surface。
- `/studio/v2` 雖然已具備導向型流程與 runtime handoff，但 commit 仍主要是前端 reducer / 模擬日誌，不是正式 backend lifecycle。
- `/studio/runtime` 已有 snapshot + SSE，但 lifecycle 語意仍不足，前端部分狀態需要字串判斷。
- `/gateway/*` 目前主要依賴 feature flag 與 `/test/connect`，應視為 prototype / connect-only，而不是 datalink persisted flow。
- router 預設入口目前仍是 `/studio`，與未來希望服務預設採 `/studio/v2` 的方向不一致。
- 使用者已明確決定：`/studio` 這一輪先暫停；`/studio/v2` 與 `/studio/runtime` 是近期重點；`/test` 先暫停但需要完整記錄；文件需要反映這個優先順序，而不是只記技術成熟度。
- 第一版 `index.html` 看不到資料的根因是內嵌 JS template string 內含未跳脫反引號內容，導致整段 script parse 失敗、render 完全沒啟動；已改為 `index.html` + `inventory.css` + `inventory.js` 三檔。
- repo 目前沒有既存的 `studio-surface-inventory` 變更履歷機制；若只靠 md/html 本身，無法避免文件被更新後難以追溯。
- `go.mod` 已內建 `modernc.org/sqlite`，可直接在 repo 內提供小型 Go CLI 來管理 inventory changelog，不必引入新外部依賴。
- `cmd/studio_inventory_changelog` 現已落地，並已初始化 `docs/technical/studio-surface-inventory/changelog.sqlite`；後續 inventory 文件更新應以這個 CLI 寫入 changelog，而不是手動維護另一份文字紀錄。
- 若要讓接手 AI 快速上手，關鍵不是再增加長文，而是提供固定第一入口與 machine-readable 摘要，讓它先讀 `START_HERE.md` + `context.json`，再按需展開長文件。
- 要進一步接近「新對話直接接上」，還需要一份最新狀態快照；否則 AI 雖知道去哪裡讀，仍要自己整理最近一次的決策與下一步。

## 資源
- [useRuntimeStream.ts](file:///Users/yishow/prj/go_gateway/frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts)

## 視覺/瀏覽器發現
- 新的即時診斷日誌面板 (RealtimeLogsPanel.tsx) 配合 auto-scroll 效果在點位狀態轉變時非常流暢且易讀，在點位等待時呈現呼吸波動 Skeleton 骨架屏顯得富有質感。

---
*每執行2次查看/瀏覽器/搜尋操作後更新此檔案*
*防止視覺資訊遺失*
