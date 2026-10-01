# Copilot CLI 任務提示詞模板（可重用）

> 使用方式：把 `{{...}}` 區塊替換後，存成實際 prompt 檔，再丟給 `scripts/copilot_watchdog.sh`。

你正在 `{{WORKDIR}}` 工作，分支：`{{BRANCH}}`。

## 目標
{{GOAL}}

## 交付要求（必達）
{{DELIVERABLES}}

## 限制
- 優先使用現有檔案與既有架構，避免無必要重構。
- 不要只做靜態分析，必須實際執行驗證命令。
- 若主機無 `go`，請改用 Docker `golang:1.25` 執行。

## 驗證命令（照順序執行）
{{VALIDATION_COMMANDS}}

## 回報格式（固定）
1. changed files list
2. key diff summary
3. commands run
4. test/benchmark summary（含是否 PASS）
5. remaining risks

## 完成事件（必做）
完成後請執行：
`openclaw system event --text "{{DONE_EVENT_TEXT}}" --mode now`
