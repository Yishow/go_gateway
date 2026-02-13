## 1. Routing and Container Refactor

- [x] 1.1 移除 `/datalink` 外層 `DatalinkLayout` 包裝，讓 `SmartDashboard` 成為首頁容器
- [x] 1.2 新增 section intent 解析（query 或 route state）以承接 legacy 導向
- [x] 1.3 調整 `/datalink/devices`、`/datalink/settings`、`/datalink/points`、`/datalink/mappings`、`/datalink/wizard` 導向規則

## 2. Dashboard Top Controls Integration

- [ ] 2.1 實作「分頁導覽 + 已選設備條」同列水平布局（桌面）
- [ ] 2.2 補齊平板/手機斷點行為（改為上下堆疊）
- [ ] 2.3 以 context bar 顯示設備名稱、協議、狀態與切換入口
- [ ] 2.4 context bar 顯示最近切換時間
- [ ] 2.5 設備狀態色階（active/offline/read-only）與對應文案
- [ ] 2.6 搜尋列與上方控制列採同一視覺網格對齊
- [ ] 2.7 無設備空狀態僅提供「選擇設備」與「新增設備」兩個主要行動
- [ ] 2.8 右側 Quick Actions 依設備狀態動態排序

## 3. Device Switching Drawer and Guardrails

- [ ] 3.1 將 `DeviceTreeNav` 改為 on-demand drawer 模式
- [ ] 3.2 完成切換設備成功流程（更新 context + 重新載入 workspace 資料）
- [ ] 3.3 完成異常流程：設備離線、切換失敗保留舊設備、權限不足 read-only
- [ ] 3.4 完成未儲存變更確認彈窗（儲存後切換 / 放棄 / 取消）
- [ ] 3.5 未儲存變更時在切換設備入口顯示紅點提示
- [ ] 3.6 確認彈窗預設焦點為 `cancel`
- [ ] 3.7 切換請求進行中鎖定切換按鈕並顯示 loading
- [ ] 3.8 切換失敗 toast 提供 `重試` 與 `查看詳情`

## 4. Device Modal Flow (Dashboard-first)

- [ ] 4.1 將 `/datalink/devices/new` 轉為 dashboard 內 modal 流程（非獨立頁）
- [ ] 4.2 新增設備 modal stepper（identity -> connection -> validation）
- [ ] 4.3 建立成功後提供「立即切換到新設備」

## 5. Local Modbus 5020 Dedicated Workbench

- [ ] 5.1 新增 `/datalink/local-modbus` 獨立頁（完整功能頁）
- [ ] 5.2 實作本機 `5020` server lifecycle（status/start/stop/bind 狀態）
- [ ] 5.3 實作本機 memory grid 映射編輯、批次操作、import/export
- [ ] 5.4 實作 `5020` 目標位址衝突掃描與阻擋機制
- [ ] 5.5 實作 test write/read 與 sync-from-mappings
- [ ] 5.6 在規則上明確區分「設備來源隔離」與「本機 `5020` 目標衝突」

## 6. Verification and Regression Coverage

- [ ] 6.1 新增/調整路由測試，驗證 legacy 路徑導向與 section intent
- [ ] 6.2 新增/調整 `SmartDashboard` 互動測試（context bar、drawer、切換流程）
- [ ] 6.3 新增/調整異常流程測試（offline、switch failure、unsaved changes、read-only）
- [ ] 6.4 增加 E2E 主線：切換成功、切換失敗、未儲存攔截、read-only
- [ ] 6.5 新增 `local-modbus` 獨立頁流程測試與衝突測試
- [ ] 6.6 執行 `cd frontend; npm run test` 與 `cd frontend; npm run lint`，附上結果
