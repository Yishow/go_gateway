# Gemini 前端設計流程（go_gateway）

## 目的
讓 Gemini 穩定扮演前端設計師，產出可直接派工的 UI/UX 規格。

## 流程
1. 先蒐集上下文（相關頁面、元件、測試、現況問題）
2. 呼叫 Gemini（固定 gemini-3-pro-preview；失敗直接報錯，不做 fallback）
3. 收斂為三份輸出：
   - PRD（目標與範圍）
   - UI Spec（元件/互動/狀態）
   - Ticket Backlog（優先級與驗收）
4. 寫入 docs 並開 PR

## 命令範本
```bash
# 固定使用 3-pro-preview
gemini -m gemini-3-pro-preview -p "你是資深前端產品設計師，請針對 SmartDashboard ..."

# 若不可用：直接報錯（不降級）
# 請回報 ModelNotFound/權限/配額錯誤並停止流程
```

## 驗收門檻
- 設計建議需可驗證（含明確驗收條件）
- 不得只給抽象建議，需含實作路徑
- PR 必附測試與風險說明
