# 採集預覽交接（2026-10-31 UTC）

## 已完成的邊界

- 唯一來源是已安裝的 native snapshot 與原 planner；未安裝／無法證明的情境回 unavailable，不以 SQL 配置或記錄選取點數冒充實際採集集合。
- GET 嚴格綁 workspace/device、DB setup revision 與可選 preview head；前後重驗，foreign/stale 安全拒絕。查詢不安裝、不連線、不探測、不讀寫 PLC。
- 各原 polling group 分列估算與實測。實體地址單位、計畫請求、wire attempts、logical outcomes、retry/fallback 與部分 output loss 不混算；沒有證據不補零。
- Scoped recording-start 已沿原完整 SQL graph 與實際安裝 slots 發布 geometry；不借 full Apply，不移除／重啟未選 slots，不把未知舊 slots 補成已授權。SQL acceptance 前不發布新 mapping／observer／projection 接受證據，原 raw wire 品質層級不改。
- capability revision 依 native 固定格式獨立解析；通用 scope identity 不放寬。真 API wire fixture 穿過 service、hook、panel 的測試防止手製 fixture 掩蓋契約。

## 已執行驗證

- Native/API 精確 normal、race、vet、官方 new-diff lint 通過；原 source owner／scoped／foreign／reentry／取消與 SQL acceptance 邊界涵蓋。
- Protocol 五類欄位的局部驗證沿既有 editor protocol handoff；本輪沒有假稱實體串口／PLC 現場驗收。
- 真 wire service＋consumer 23 cases 通過；完整候選 strict TypeScript、Vite、靜態嵌入及 Go build 通過。
- Run17 從未 Start 的隔離 DB 沿原 UI 設新 SQLite 目標、schema consent、recording-start 成功，Gateway 未重啟後取得 available；DOM 與 GET 一致顯示同 group 的 2 個 float32 points、4 個 16-bit 地址單位、planned 1、actual wire 1、logical success 1，retry/fallback 為真實 0。
- Refresh 產生新 GET；診斷／預覽與 footer 沒有覆蓋。最後頁面導航、console 與 SQL 檢查完成，沒有 React #185 或未處理例外；外部 font 拒連另列為環境限制。
- Run17 target 有實際 SQL rows／receipts；這只證明該次完整流程，不把採集次數當 SQL 寫入次數。

## 修復與保留的失敗證據

- Run10 的預覽布局覆蓋已由既有 SectionCard/Button 與內在高度 grid 修正，經真 Refresh/截圖重驗。
- Run13 重現 DB 逐字輸入 React #185／漏字；根因是未變來源被每個 DB 字元重建並重做草稿 hydration。CircuitProfileBasicPanel 的 memo 改為四個實際 source slices，沒有改 CAS、generation 或延遲工具輸入。Run14 同速原 87 字路徑及長名稱／路徑在 DOM、PUT、ACK、SQL、reload 完全一致。
- Run14 正式 scoped Start 後 preview unavailable 暴露 publication 漏接；Run15 API available 但畫面 unavailable 暴露 capability revision 解析不符。兩項皆有修前證據與修後精確回歸；Run17 閉合正式 UI 路徑，不用 full bootstrap 遮蔽 scoped 問題。

## 尚未宣稱

4.2 完整負載矩陣、全專案終驗與實體硬體驗收尚未完成。此交接不代表全部七個 changes 完成或已發布 GitHub main。
