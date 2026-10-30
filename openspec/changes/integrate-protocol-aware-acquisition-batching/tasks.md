## 1. 純計畫與兼容性擴展

每task採Red→Green→Refactor且維持build；時間預估每項不超一小時／單一context，實作超界先拆獨立可驗證子change。preserve-per-point-word-order完成後才接production decoder；可先做純planner fixtures。僅勾選項目已完成，其餘待實作。

- [x] 1.1 交付Static request optimization的「唯一不可變計畫與能力界線」：在既有optimizer新增revision-bound輸入／結果，不改runtime；測invalid-address顯式錯誤、device/station/function/cadence隔離與 deterministic plan。
- [x] 1.2 [after: 1.1] 交付PDU awareness與「安全覆蓋與去重解碼」raw span union／safe contiguous planner；先以D0 int64＋D1 int16重現shrink，再驗125/126、多word不切半、禁止gap／明確coverage、radix與native width。
- [x] 1.3 [after: 1.2] 交付Capability-aware random acquisition的planner及bounded candidate比較；fixture驗MC255/256、FATEK64/65 native units、D0/D1/D10一次、dense21選較少bytes及unknown capability安全退回。

## 2. 有界傳輸與逐點結果

- [x] 2.1 交付Concurrency control的connection generation／短鎖建立；在connector以慢dial第二device進度、同key single-flight、config改變／late dial及shared-serial alias拒絕fixtures證明不借舊connection。
- [x] 2.2 [after: 2.1] 交付Retry and backoff之MC實際wire attempt預算與cancel deadline；fixture涵蓋scheduler/adapter/transport重送、unsupported與timeout分流，確認取消後不新send、不關新版connection。
- [ ] 2.3 [after: 2.2] 將相同Retry and backoff契約套FATEK／Modbus現有transport，不新增引擎；fixture驗backoff取消、deadline、serial解除等待與最終Bad，所有attempt計入統一預算。
- [ ] 2.4 [after: 1.3, 2.3] 交付MC raw random/contiguous adapter執行與Identity-preserving raw deduplication；以正式adapter fixture驗完整scalar、重疊多Point及ABCD/CDAB各自解碼，空／短／壞回應不產生Good。
- [ ] 2.5 [after: 2.4] 交付FATEK random native scatter與嚴格response驗證；正式adapter fixture驗D/DR/bit widths、空items、短／多餘字元／非0或1bit、缺map key與逐點格式，不默填零。
- [ ] 2.6 [after: 2.5] 交付Modbus safe word/bit batch scatter；fixture驗全部bits不是first bool、不同function隔離、授權／禁止gap及整包錯誤，保留現有TCP/RTU/UDP scalar語意。

## 3. 既有排程遷移

- [ ] 3.1 [after: 2.6] 在既有collector遷移Scheduled collection為planned reads，落實「結果、時間與版本不冒充一致快照」與Quality and timestamp；production wiring fixture驗scan/request identity、same-batch completion、每點一次emission、舊lease/revision不刷新新設定。
- [ ] 3.2 [after: 3.1] 交付「調度、重試與背壓有界」：due producer脫離pollGroup全設備等待、bounded per-group coalescing／fair queue與drop counters；fake clock／race驗慢device不擋同group快device、多group不互覆、queue飽和可觀測及stop零新send。
- [ ] 3.3 [after: 3.2] PollNowContext使用同planner/executor，保留exact requested subset／cancel並移除重複推導；manual與scheduled交錯fixture驗device lock、未要求點無新讀值、一次final結果與source-set不受record selectors影響。

## 4. 預覽與優化驗收

- [ ] 4.1 [after: 3.3] 交付Revision-bound request preview及「同planner預覽與量測證明」scoped API／shared types，前端consumer僅呈現同源estimate與actual分類；API測零PLC I/O、foreign scope拒絕、stale409、預覽execution parity，frontend測不稱封包／SQL次數。
- [ ] 4.2 [after: 4.1] 完成Reproducible optimization evidence：設計矩陣baseline／optimized、warm/cold、planner100/1000/10000點與fault/pressure量測，報request/bytes/p50/p95/freshness/alloc/drop及值品質等價；不達控制RTT不退化門檻的策略不啟用。重跑相關Go test/vet/race、API/frontend受影響gate與cmd/test_ui端到端fake transport，交接實測／未測與逐點字序依賴結果，不把提案當驗收。
