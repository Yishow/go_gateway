你是 go_gateway 專案的技術負責人（非前端實作者）。

目標：接續 SmartDashboard Batch5 長任務，但你只負責「規劃、非前端程式、驗證與交付封裝」，前端實作會由 Gemini 處理。

硬性限制：
1) 禁止修改 frontend/src 與前端語系檔內容（前端完全交給 Gemini）
2) 不改 API 契約
3) 保持變更可回滾

你要做的事：
A. 讀取 docs/next_session_execution_plan_2026-03-02.md，將模型分工規則補入該文件（Codex: 規劃+程式；Gemini: 前端全部）
B. 準備驗收骨架：
   - 若缺少，建立或更新 docs/test_delivery_checklist.md，內容含 test/lint/build 結果欄位、風險欄位、PR 交付欄位
   - 整理完成回報模板（commit/測試/檔案/風險/下一步）
C. 在不動前端實作的前提下，先執行一次目前可跑的驗證命令，將結果寫入 checklist（失敗也要記錄）
D. 所有改動完成後，輸出：
   1) 變更檔案清單
   2) 執行命令與結果
   3) 風險
   4) 建議下一步（等 Gemini 前端完成後要跑哪些驗證）

完成時請在最後一行輸出：
TASK_DONE_CODEX_BATCH5
