## 1. Backend Red→Green
- [ ] 1.1 建立 fixture database 與期望查詢結果，覆蓋 samples/intervals/events、權限、分頁 cutoff與大整數。
- [ ] 1.2 實作各 adapter 有界索引查詢、primary destination、retention與capability回應。
- [ ] 1.3 實作用量/evidence drilldown、缺口與estimation/coverage欄位，直接使用3的結果。
- [ ] 1.4 實作串流CSV、formula injection防護、取消/逾時、跨workspace拒絕與精度往返測試。

## 2. Frontend TDD
- [ ] 2.1 先測電表三相曲線、kW/kWh區別與MC混合感測器頁籤，再實作runtime UI。
- [ ] 2.2 加入缺口不補零、過期樣本、實際粒度、缺台總量、signed負值與state unknown測試。
- [ ] 2.3 實作方案／設備／項目／時間切換、查原值與返回精確修正位置。
- [ ] 2.4 整合CSV與安全錯誤，確保圖表、表格、CSV相同cutoff的數字一致。

## 3. Refactor、全組驗收與發行
- [ ] 3.1 將 acceptance.md 全部案例轉為後端／前端／E2E自動測試；逐案留下證據。
- [ ] 3.2 三種DB實測斷線、重啟、回應遺失、補送、更正、資料表權限與quota限制。
- [ ] 3.3 執行100×8×5秒負載、固定硬體查詢基準與長時間資源觀察，量測不足不得聲稱達標。
- [ ] 3.4 Windows/Linux/ARM與真PLC按矩陣驗收；未實測項目明確列pending。
- [ ] 3.5 全組go test/vet/lint、frontend lint/test/build、Playwright、route/Share回歸與六份strict spec驗證。
- [ ] 3.6 更新部署、retention、備份／恢復／回退文件與使用者操作說明，所有完整性案例通過再開正式功能。
