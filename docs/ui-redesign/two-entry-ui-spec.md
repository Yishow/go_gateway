# Two-Entry UI 重構規格書（SPEC）

版本：v1.0
日期：2026-03-02
決策：1A（快速入口預設）/ 2B（Feature Flag 灰度）/ 3C（雙入口同時推進）

---

## 1. 目標與範圍

### 1.1 目標
- 在不破壞既有功能與 API 契約前提下，重構 SmartDashboard 操作流程為雙入口：
  1) Quick Setup（預設）
  2) Expert Workbench（進階）
- 降低新手完成首個設定任務的操作成本。
- 保留進階使用者完整控制能力。

### 1.2 非目標
- 不變更後端 API contract。
- 不重寫資料模型。
- 不一次性移除舊流程（需透過旗標灰度）。

---

## 2. 角色與主要場景

### 2.1 角色
- 操作員（新手/日常）：偏好任務導向、快速完成。
- 工程師（進階）：偏好全局控制、批次與進階設定。

### 2.2 核心場景
- S1：建立新映射（設備→點位→Tag→驗證→提交）。
- S2：調整既有映射與 Tag（含衝突提示）。
- S3：檢視與處理提交/回滾/重試。

---

## 3. 資訊架構（IA）

- `/datalink`：入口頁（預設）
  - 入口卡 A：Quick Setup（主 CTA）
  - 入口卡 B：Expert Workbench
  - 側欄：最近任務/草稿/最近提交狀態
- `/datalink/quick-setup/*`：精靈式流程
- `/datalink/workbench/*`：進階工作台

---

## 4. 流程規格

### 4.1 Quick Setup（預設）
步驟固定：
1) 選設備
2) 選點位
3) 綁定 Tag（既有/新建）
4) 驗證（衝突/完整性）
5) Commit（成功/部分成功/失敗）

每步都要有：
- Loading / Empty / Error / Success 四態
- 下一步 CTA
- 可返回上一步且不丟草稿

### 4.2 Expert Workbench
- 保留既有 SmartDashboard 進階能力（批次、全域編輯、進階檢查）。
- 可隨時切回 Quick Setup。
- 切換入口不丟失草稿。

---

## 5. 相容與不破壞策略

1) 單一資料來源（draft/commit state 共用）
2) Adapter 層轉換：UI payload -> 現有 API 格式
3) 保留舊路由相容導向（301/前端轉址）
4) 嚴禁破壞性 schema 變更

---

## 6. 互動與文案標準

### 四態標準
- Loading：明確進度文案 + 控制可否操作
- Empty：原因 + 建議下一步
- Error：錯誤原因 + 可執行修復按鈕（Retry / 返回 / 切換入口）
- Success：成功或部分成功需明確區分

### 關鍵操作
- Commit / Retry / Rollback / Workflow：
  - 統一按鈕狀態
  - 統一事件回饋節奏
  - i18n key 必填（zh-TW / en）

---

## 7. Feature Flag（2B）

旗標：`feature.datalink.twoEntryUI`

- Off：維持舊路徑
- On（灰度）：指定使用者/環境可見雙入口
- 全開：雙入口成為預設

觀測指標：
- 任務完成率
- 首次成功配置時間
- 失敗重試次數
- 回滾率

---

## 8. 測試規格

- 單元測試：流程狀態機、資料轉換、guardrail
- 整合測試：Quick/Workbench 切換不丟草稿、提交流程一致
- E2E：入口到提交全流程、錯誤回復路徑

固定驗證：
```bash
cd frontend
npm test -- --run
npm run lint
npm run build
```

---

## 9. 驗收標準（DoD）

- [ ] 雙入口可用且功能對齊
- [ ] 不破壞既有 API 契約
- [ ] 入口切換不丟草稿
- [ ] 四態與關鍵操作回饋一致
- [ ] i18n（zh-TW/en）完整
- [ ] test/lint/build 全綠
