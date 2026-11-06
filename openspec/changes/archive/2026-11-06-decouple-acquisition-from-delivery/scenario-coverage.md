# 本次 apply 的情境測試對照

此表區分測試內容與執行結果；執行指令、環境及限制見 [本次交接](apply-completion-evidence.md)。這是 apply 要求的覆蓋檢查，不代替完整 `spectra-verify` 或 `spectra-review`。

主代理依實際斷言、正式元件、當次本機 wire／API／瀏覽器結果判定：24 個 Scenario 中 23 個有測試覆蓋、1 個部分覆蓋；2 個 Example 均有對應具體值。沒有依 test-scope 排除的整個 Scenario／Example。單一畫面的尺寸、置中與間距以實際截圖檢視，未增加 CSS 尺寸斷言。

## database-output-delivery

| Scenario／Example | 分類 | 測試及具體值 |
| --- | --- | --- |
| Unreachable connector does not slow ingest | covered | `internal/datalink/runtime/ingestor_database_decoupling_test.go`：`TestHandleCollectedValue_UnreachableDatabaseConnectorDoesNotDelayValueEvents`，32 筆各 <50ms、32 個事件；`writer_decouple_timeout_regression_test.go` 驗證卡住時仍回報失敗。當次 loopback 另實測 32 點採集吞吐。 |
| Only the latest value per interval is written | covered | `internal/datalink/dbtarget/writer_snapshot_test.go`：`TestWriter_WritesOnlyLatestValuePerInterval`，15 秒內 10／11／12 → 一列 12，無新值不增加列。 |
| Example: interval snapshots | covered | 上項直接覆蓋 10／11／12；`writer_backoff_test.go` 的 `TestWriter_RecordsSkipWhenIntervalHasNoNewValue` 覆蓋 none → 無新增列，history 含 `no_new_value`。 |
| Failed connector backs off | covered | `writer_backoff_test.go`：`TestWriter_BackoffRetriesLatestSnapshotWithInjectedClock`，1／2／4／8／16／32／60／60 秒、逐次 failed history、首次 failed outcome、最後 success 與只寫 12。當次真 API 另觀察失敗期間零新增列及恢復寫入。 |

32 個事件的每個 ID／值未在同一個阻塞測試逐項比較；單點事件內容由既有 stream 測試保護。退避測試也未逐次比較 outcome 的時間戳。這些是斷言粒度的限制；原任務明列的逐筆耗時、事件數、逐次失敗歷史及退避結果已有測試，未另加驗收條件。

## acquisition-anomaly-log

| Scenario | 分類 | 測試及具體值 |
| --- | --- | --- |
| Stale point opens and closes an event | covered | `internal/datalink/anomaly/detector_test.go`：`TestStalePointOpensAndClosesOneEvent`，D5、1000ms、10 秒門檻、12 秒無讀取後開立，成功讀取時結束。 |
| Frozen detection is opt-in | covered | `TestFrozenDetectionIsOptIn`，值 0 一小時不開事件，600 秒門檻啟用後開 frozen。 |
| Disconnect is recorded once | covered | `TestDisconnectIsRecordedOnce`，重複失敗 30 秒後恢復，只有一筆及正確起訖。 |
| Filtered listing | covered | `internal/api/handlers/acquisition_anomaly_handler_test.go`：`TestAcquisitionAnomalyHandlerFiltersByStateAndKind`，active／stale、設備隔離及 total。 |
| Retention keeps active events | covered | `internal/datalink/anomaly/sql_repo_retention_regression_test.go`：`TestAnomalyRetentionDefaultPrunesOldestEndedAtFiveThousand`，5000 筆含 3 active，新事件後移除最舊 ended，原 active ID 保留。 |
| Active events exceed the retention target | covered | `TestAnomalyRetentionOverFiveThousandActivePrunesOnEnd`，5001 active 全保留，結束一筆後剩 5000 active；新增 HTTP 回歸另驗證超額時查詢與摘要。 |
| Warning state with count | covered | `frontend/tests/unit/workbench-v2/recording-home-acquisition-anomalies.test.tsx`：兩個首頁位置均斷言「採集異常 2」、公開 rendered selector `data-anomaly-state=active`、緊接讀取點位。主代理另核對該 selector 的 warning tokens 及當次瀏覽器警告呈現；自動測試未比較 computed color。 |
| Read failure is not an empty state | covered | 同檔 `shows an actionable error with retry instead of an empty list, then recovers`，錯誤、重試、無空列表假象與恢復；當次 mock API 瀏覽器另實際操作。 |

## device-acquisition-interval

| Scenario | 分類 | 測試及具體值 |
| --- | --- | --- |
| Interval change affects only one device | covered | `internal/datalink/device/service_acquisition_test.go`：`TestService_Update_AcquisitionIntervalReassignsOnlyTargetDevicePointsAndUpdatesScheduler`，A 1000→500ms、B 保持 1000ms。真 API 另測 B 1000→3000ms、A 保持 1000ms。 |
| Upgrade keeps frequency | covered | `internal/datalink/migrator_device_acquisition_test.go`：`TestDeviceAcquisitionBackfillKeepsIntervalsAndIsIdempotent`，原 1000ms 保持、點位轉移、重跑不重複。 |
| Upgrade combines mixed intervals at the shortest interval | covered | `TestDeviceAcquisitionBackfillUsesShortestEnabledInterval`，A 的 500／2000ms 統一 500ms、B 2000ms；停用／未指派與停用點位保留、重跑冪等。 |
| Slow device stretches the interval | covered | `internal/datalink/collector/scheduler_auto_interval_test.go`：`TestSchedulerAutoIntervalStretchesAndShrinksBack`，1250ms→1500ms、連續 5 輪後回 1000ms。 |
| Protocol-specific options | covered | `frontend/tests/unit/workbench-v2/recording-acquisition-settings.test.tsx`，Modbus TCP／UDP／RTU 無 random、MC 四種選項、FATEK 三種選項與實際 payload。 |

## random-acquisition-range

| Scenario／Example | 分類 | 測試及具體值 |
| --- | --- | --- |
| Mixed word, double word and bit in one request | covered | `internal/datalink/optimizer/acquisition_random_contract_test.go`、`internal/protocol/mcprotocol/client_planned_exact_test.go`、`internal/datalink/connector/adapters/mc3e_planned_peer_test.go`，D39 int16／D100 int32／M5 bool；0403/0000 兩 word、一 dword、M0 對齊與 bit 5 解碼。 |
| Long timer is excluded | covered | `acquisition_random_long_timer_regression_test.go` 的 `TestAcquisitionMCLongTimerRejectedAlongsideRandomWords` 與 `client_word_read_test.go` 的 `TestBatchReadWordRejectsLongTimerCounterWithoutSending`，明確錯誤、零 LTN 請求；random 入口另有既有拒絕測試。 |
| Example: long timer beside random items | covered | `TestAcquisitionMCLongTimerRejectedAlongsideRandomWords` 使用 D10／D20／LTN0 的原值，錯誤歸屬 LTN0，其餘有效點沿用既有規劃方式。 |
| Limit cannot be raised | gap（部分覆蓋） | `internal/datalink/connector/acquisition_profile_random_limit_test.go`、handler 與前端設定測試含 120 接受／300 拒絕且原宣告不變／FATEK 48 接受。實際請求拆分測試使用上限 2，尚無以字面值 120 驗證請求拆分的測試。 |
| Mixed rule round trip | covered | `internal/api/handlers/studio_v2_workspace_source_rules_mixed_test.go`：`TestWorkspaceSourceRuleMixedRoundTrip`，Q/iQ-R、D39／D100／M5 儲存與讀回、三個 point/link。 |
| Overlap is rejected | covered | 同檔 `TestWorkspaceSourceRuleMixedOverlapIsFieldError`，D3 uint32 與 D4 int16，欄位錯誤列出兩項，資料不變。 |
| Safe profile keeps excluding mixed rules | covered | `internal/datalink/sourcerule/acquisition_source_mixed_test.go` profile matrix、`internal/datalink/workspace/runtime_acquisition_mixed_test.go`，safe-v1 回 `source-layout-unsupported`。 |
| Warning at 90 percent | covered | `frontend/tests/unit/scatteredRangeForm.test.tsx`，172 不警告、173／192 警告，RuleEditor 有 warning DOM，超額仍可儲存。 |
| One range replaces several | covered | `cmd/test_ui/acquisition_scheduled_mixed_test.go`：`TestProductionScheduledAcquisitionMixedRuleWireRequests`，D0–D31 加 D39／D40／D100／D200／D1000／D2000，斷言 ≤2 requests、38 個同輪 Good samples 與六個分散值；當次 wire log 實際是一個 random、38 units。 |

唯一部分覆蓋項目不阻擋 apply 完成。若另行加強，可補有效上限 120 的精確計畫拆分測試；本輪未因此擴充功能或新增完成條件。真 PLC、Windows 與正式資料庫不是這份測試覆蓋表的替代證據，也未宣稱已驗收。
