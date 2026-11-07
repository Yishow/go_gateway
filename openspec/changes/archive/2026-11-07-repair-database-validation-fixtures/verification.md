# 測試 fixture 驗證與審查

## 範圍與結果

本案僅包含 proposal 的四個測試檔案；沒有產品程式或delta spec。基準c7045b9e，scope attribution為approximated_base，旁邊的schema-confirm產品修復由另一change負責。

5/5任務完成。validate通過（no-spec的No delta specs found提示符合選定schema）、analyze 0 findings，diff與行數檢查通過。

## 實際測試

- PG17.11/Timescale2.22.1 recovery/stop/remove/hypertable相關選測：33PASS/0FAIL/0SKIP。
- 原跨分鐘案例count=3：15PASS/0FAIL/0SKIP。清理後extension count=1，服務已停止。
- SQL Server successor完整回歸：97PASS/0FAIL/0SKIP，覆蓋native cleanup和runtime時間邊界；這個結果取自新增產品負向回歸前的完整綠燈紀錄，四個fixture檔案內容與當前一致。
- 最後API全套、full lint/vet/build由整合發布gate執行，本報告不提前聲稱完成。

## Review

已逐檔依正確性、效率、簡化重用、專案慣例核對，沒有具體阻擋缺陷。獨立整合review也核對fresh pool、glmTargetPartition及時間比較。

- Timescale重建使用全新pool並明確close，保留成功斷言。
- fixture同批Apply只使用同一clock，defer WithClock(nil)符合service還原UTC wall clock的既有契約。
- NativeMonthly僅補漏的證據，不放寬產品schema/ownership檢查。
- RFC3339Nano解析+time.Equal保留奈秒精度，只接受同一瞬間。

## 限制

本案不修schema-confirm離線重放產品行為，不代表全524個skip均已補測。完整18組長負載矩陣、真PLC與合理的helper/dialect skip仍依總驗證報告如實標示。
