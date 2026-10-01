# C 方案 W4 上線檢核清單（Dual Entry）

> 範圍：`/gateway/entry`、`/gateway/quick-setup`、`/gateway/expert-workbench`
> 原則：不破壞既有 API 契約，隨時可回退至 legacy (`/datalink`)。

## 1. 發佈前檢核

- [ ] Feature Flag `ENABLE_GATEWAY_DUAL_ENTRY=true` 已在目標環境配置
- [ ] `test:gateway:gate` 全數通過（unit + e2e + build）
- [ ] Quick / Expert 提交流程可呼叫 `/api/v1/test/connect`
- [ ] `X-UI-Version` header 已驗證：
  - [ ] `/gateway/quick-setup` 送出時帶 `dual_quick`
  - [ ] `/gateway/expert-workbench` 送出時帶 `dual_expert`
- [ ] Flag OFF 時 `/datalink/devices/new` 可正確回退到 `/datalink`

## 2. 灰度觀測（建議）

- [ ] 先 10% 流量灰度
- [ ] 觀測指標：
  - [ ] connect 成功率（Quick vs Expert）
  - [ ] 轉換漏斗（Entry -> Quick/Expert -> Submit）
  - [ ] 前端錯誤率（JSON parse / submit error）
- [ ] 若異常上升（錯誤率/失敗率）超門檻，立即回退 Flag

## 3. 回滾策略

### 即時回滾（無停機）
1. 將 `ENABLE_GATEWAY_DUAL_ENTRY=false`
2. 驗證 `/datalink/devices/new` 回到 legacy 路徑
3. 保留新頁面程式碼，待問題修復後再啟用

### 回滾後確認
- [ ] 新增流程回到 `/datalink`
- [ ] 既有 API 路徑與契約行為不變
- [ ] 監控告警恢復正常

## 4. 指令

```bash
cd frontend
npm run test:gateway:gate
```

---

## 5. 當前基準（2026-03-02）

- Unit: PASS
- E2E: PASS（3 條主流程）
- Build: PASS
- 最新相關 commits:
  - `252b4d6` W3 Expert Workbench 最小可操作版
  - `1df64cf` Quick/Expert 串接 Connect API 提交流程
  - `40f9ad6` C 方案三條 Playwright E2E 驗收
