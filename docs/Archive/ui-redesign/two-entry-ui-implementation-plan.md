# Two-Entry UI 實作計畫（PLAN）

版本：v1.0
日期：2026-03-02
對應規格：`docs/ui-redesign/two-entry-ui-spec.md`

---

## 1. 實作策略

- 策略：雙入口同時推進（3C），但以可回滾小步提交。
- 原則：
  - 不動 API 契約
  - 每階段可獨立驗收
  - 任何階段失敗可快速回退

---

## 2. 工作分工

- Copilot/Codex：規劃、非前端程式、驗證、交付封裝
- Gemini：前端 UI/UX 與前端實作（含預覽圖）

---

## 3. 里程碑（W1~W4）

### W1：基建與入口頁
- 建立 `/datalink` 入口頁
- 建立 Feature Flag 開關與守門
- 接上最近任務/草稿摘要區

交付：
- 入口頁可開關
- 舊路由仍可用

### W2：Quick Setup 流程
- 完成 5 步精靈（設備/點位/Tag/驗證/Commit）
- 補齊四態與錯誤回復
- 接上 i18n key

交付：
- 可從頭到尾完成一次提交

### W3：Expert Workbench 對接
- 整併現有進階能力至 workbench 路由
- 實作 Quick <-> Workbench 切換且草稿不丟失

交付：
- 進階流程可完整執行

### W4：灰度與收斂
- 打開灰度（小流量）
- 補齊監控與回滾腳本
- 修復回歸問題，準備全開

交付：
- 灰度報告 + 全開建議

---

## 4. Commit 規則

- Commit A：入口頁與路由（不含進階流程）
- Commit B：Quick Setup 核心流程
- Commit C：Workbench 對接
- Commit D：i18n / 測試 / 文件

---

## 5. 驗證節奏

每個里程碑都跑：
```bash
cd frontend
npm test -- --run
npm run lint
npm run build
```

長任務啟動前跑：
```bash
make longtask-smoke
```

---

## 6. 風險與應對

- 風險：雙入口狀態不一致
  - 應對：單一 store + adapter + 切換整合測試
- 風險：灰度期間行為分歧
  - 應對：旗標紀錄 + 事件追蹤
- 風險：文案/i18n 回歸
  - 應對：key 規範 + 快照測試

---

## 7. 上線條件

- W1~W4 驗收全過
- 灰度數據穩定（完成率與錯誤率達標）
- 回滾演練完成
