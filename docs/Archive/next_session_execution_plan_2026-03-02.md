# go_gateway 下一步執行計畫（重開會話版）

更新時間：2026-03-02 02:10（Asia/Taipei）
適用分支：`feat/datalink-smartdashboard-batch5-i18n`

---

## 0) 已完成基線（重開會話先確認）
- SmartDashboard 重構 Batch 1~4 已合併
- 前端/後端最小 CI 已建立
- 手動 gate 流程已建立
- 前端入口：`http://localhost:5173`

---

## 1) 本輪目標（不改 API 契約、不大改視覺）
1. 統一關鍵流程四態體驗：`Loading / Empty / Error / Success`
2. 優化 Commit / Retry / Rollback / Workflow 的文案與回饋節奏
3. 清理 SmartDashboard 硬編碼字串，改用 i18n key（`zh-TW` / `en` 同步）
4. 維持測試與 CI 全綠

### 1.1) 模型分工規則（本輪固定）
- **Codex**：負責規劃、非前端程式、驗證執行、交付封裝（文件、checklist、回報、風險控管）。
- **Gemini**：負責 **全部前端實作**（`frontend/src` 與前端語系檔調整）。
- 共同限制：不改 API 契約；所有變更需可回滾（小步提交、可分段還原）。

---

## 2) 重開會話後「直接執行」步驟（照順序）

### Step A：環境與分支確認
```bash
cd /home/yishow/github/go_gateway
git fetch --all --prune
git switch feat/datalink-smartdashboard-batch5-i18n
git status --short --branch
```

### Step A.1：長任務啟動模式（固定）
- 一般任務：Copilot CLI（`gpt-5.3-codex`）
- 長任務：Terminal 互動式 Copilot + `/fleet` 平行拆分（前提：當前終端支援）
- 前端任務：Gemini（`gemini-3-pro-preview`）
- 背景啟動入口：
```bash
/home/yishow/.openclaw/agents/code-pro/agent/workspace/scripts/run_longtask.sh   --name batch5-copilot   -- "cd /home/yishow/github/go_gateway && ./scripts/copilot_watchdog.sh --workdir /home/yishow/github/go_gateway --prompt-file /home/yishow/github/go_gateway/docs/templates/codex_batch5_orchestrator.prompt.md --model gpt-5.3-codex --idle-minutes 8 --max-restarts 2 --poll-seconds 20 --log-file /tmp/go_gateway_codex_batch5.log"
```
- 若當前命令列無法使用 `/fleet`，改採多 session 平行執行並維持同等進度回報節奏。

### Step B：可見體驗優化（先做）
- 範圍：`frontend/src` 內 SmartDashboard 相關元件（CommitPanel、WorkflowModal、Workspace）
- 交付標準：
  - 四態 UI 呈現一致（含按鈕可用狀態與提示文案）
  - Error 狀態附「可操作下一步」
  - 成功 / 部分成功狀態可被明確辨識

### Step C：i18n 收斂（第二階段）
- 將新增/既有硬編碼字串改為 i18n key
- `zh-TW`、`en` 兩語系同步補齊
- 測試改用穩定 selector（避免直接綁中文文案）

### Step D：驗收與交付（最後）
```bash
cd /home/yishow/github/go_gateway/frontend
npm test -- --run
npm run lint
npm run build
```

交付時固定輸出：
- 變更檔案清單
- 測試結果（test/lint/build）
- 風險與未完成項
- 下一步建議

---

## 3) 重開會話可直接貼上的任務指令（給代理）

```text
請接續執行 go_gateway SmartDashboard Batch5：
1) 先完成四態體驗統一（Loading/Empty/Error/Success）與 Commit/Workflow 回饋優化
2) 再完成 i18n 硬編碼清理（zh-TW/en 同步）
3) 不改 API 契約，不做大幅視覺 redesign
4) 每 5 分鐘固定回報：完成/進行中/風險/下一步
5) 完成後回傳：commit、測試結果、變更檔案、PR 連結
```

---

## 4) 風險與守門
- PR 變更量可能偏大：務必分批 commit（UI 行為 vs i18n/測試）
- i18n 改寫可能導致測試波動：先改 selector 再改文案
- 目前無 branch protection：需嚴守 manual gate（測試全綠才可推進）

---

## 5) 完成定義（DoD）
- [ ] 四態體驗一致且可見
- [ ] SmartDashboard 新增/既有硬編碼字串已 key 化
- [ ] `frontend` 的 test/lint/build 全綠
- [ ] 交付回報含 commit/測試/檔案/PR
- [ ] 文件同步更新（本檔 + 必要 PR 說明）
