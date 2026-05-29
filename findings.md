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

## 遇到的問題
| 問題 | 解決方案 |
|------|---------|
| 點位高頻 stale 事件狂刷 `setLogs` 造成嚴重重繪卡頓與日誌爆量 | 實作 `pointStaleStatesRef` 僅在點位 stale 狀態「改變」的瞬間才記錄 transition log (並新增 recovered 復原狀態日誌)。 |
| 新日誌 prepend 到最上方，與診斷面板的「滾動置底」UX 機制相反 | 將新日誌修改為 append 到陣列最後，使時間排序由上到下遞增，完美咬合 auto-scroll-to-bottom。 |

## 資源
- [useRuntimeStream.ts](file:///Users/yishow/prj/go_gateway/frontend/src/features/datalink/runtime-dashboard/useRuntimeStream.ts)

## 視覺/瀏覽器發現
- 新的即時診斷日誌面板 (RealtimeLogsPanel.tsx) 配合 auto-scroll 效果在點位狀態轉變時非常流暢且易讀，在點位等待時呈現呼吸波動 Skeleton 骨架屏顯得富有質感。

---
*每執行2次查看/瀏覽器/搜尋操作後更新此檔案*
*防止視覺資訊遺失*
