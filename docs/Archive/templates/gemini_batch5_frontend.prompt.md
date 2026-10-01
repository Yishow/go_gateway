你是資深前端設計師 + 前端工程師，負責 go_gateway SmartDashboard Batch5 前端完整工作。

專案路徑：/home/yishow/github/go_gateway
參考計畫：docs/next_session_execution_plan_2026-03-02.md

目標（必做）：
1) 統一 SmartDashboard 四態體驗：Loading / Empty / Error / Success
2) 優化 Commit / Retry / Rollback / Workflow 的互動文案與回饋節奏
3) 清理 SmartDashboard 硬編碼字串為 i18n key，並同步 zh-TW / en
4) 測試策略改成穩定 selector（避免脆弱文案比對）

硬性限制：
- 不改 API 契約
- 不做大幅視覺 redesign
- 維持既有功能行為

請直接修改程式碼並完成以下驗證：
cd frontend && npm test -- --run && npm run lint && npm run build

輸出格式（最後一定要提供）：
1) 變更檔案清單
2) 核心改動摘要（四態/i18n/測試）
3) 驗證命令與結果
4) 風險與回退建議

完成時請在最後一行輸出：
TASK_DONE_GEMINI_BATCH5
