# go_gateway 下一步執行計畫（重開會話版）

更新時間：2026-03-02 01:34 (Asia/Taipei)

## 0) 目前已完成基線
- SmartDashboard 重構 Batch 1~4 已合併（PR #2）
- 前端最小 CI 已合併（PR #3）
- 後端最小 CI 已合併（PR #4）
- 替代保護流程（PR template + manual gate）已建立（PR #5 待確認）
- 前端入口確認：`http://localhost:5173`

---

## 1) 下一步總目標（本輪）
在不改 API 契約、不大改視覺的前提下，做出「使用者可明顯感知」的體驗優化：

1. 統一四種狀態體驗：`Loading / Empty / Error / Success`
2. 優化關鍵操作文案與回饋節奏（Commit / Retry / Rollback / Workflow）
3. 確保 i18n key 化落地，不再新增硬編碼字串
4. 維持測試與 CI 綠燈

---

## 2) 建議執行順序（重開會話後直接照做）

### Phase A：可見體驗優化（優先）
- 範圍：CommitPanel、WorkflowModal、Workspace 主要狀態
- 交付：
  - 針對四態建立一致 UI 呈現規則
  - 錯誤訊息改為「可操作建議」
  - 成功/部分成功狀態明確區分

### Phase B：i18n 收斂
- 範圍：SmartDashboard 相關元件
- 交付：
  - 硬編碼字串改為 i18n key
  - `zh-TW` / `en` 文案同步
  - 測試改用穩定查找策略（避免脆弱文案比對）

### Phase C：驗收與交付
- 跑測試 + lint + build
- 更新 PR 說明與驗收清單
- 輸出「變更檔案 / 測試結果 / 風險 / 下一步」

---

## 3) 技術與流程守門（固定）

### 必守原則
- 不改 API 契約
- 不做大幅視覺 redesign
- 先測試再重構
- 單批次小步提交（結構重構與測試分開 commit）

### 固定驗證命令
```bash
cd frontend
npm test -- --run
npm run lint
npm run build
```

### 進度回報規範（固定格式）
每 5 分鐘主動回報：
1. 完成
2. 進行中
3. 風險
4. 下一步

---

## 4) 重開會話後第一條指令模板（可直接貼）

```text
請接續執行 go_gateway SmartDashboard 下一步：
1) 先做可見體驗優化（Loading/Empty/Error/Success 四態統一 + Commit/Workflow 文案與回饋節奏）
2) 再做 i18n 硬編碼清理
3) 不改 API 契約與主要視覺
4) 全程使用 Copilot gpt-5.3-codex + /fleet
5) 每 5 分鐘固定回報（完成/進行中/風險/下一步）
6) 完成後回傳 commit、測試結果、變更檔案、PR 連結
```

---

## 5) 風險清單（先看）
- PR 累積變更量較大：需分批 review，避免漏看
- i18n 改寫可能影響測試 selector：需先調整測試查找策略
- 無 branch protection（方案限制）：需嚴格遵循 manual gate

---

## 6) 完成定義（DoD）
- [ ] 關鍵流程四態體驗一致且可見
- [ ] SmartDashboard 不再新增硬編碼字串
- [ ] frontend test/lint/build 全綠
- [ ] PR 描述附驗收重點與風險
- [ ] 交付回報包含 commit/測試/檔案/PR
