# polling-group-name-integrity Specification

## Purpose

補齊既有 Polling Groups API 的名稱完整性：記憶體儲存庫與 SQL UNIQUE 約束一致，拒絕重複名稱且不破壞原資料。

## Requirements

### Requirement: 輪詢群組名稱唯一性保持一致

系統 SHALL 在記憶體與 SQL 輪詢群組儲存庫維持相同的名稱唯一性，新增或更新發生同名衝突時回傳錯誤，並保留原有資料。

#### Scenario: 拒絕建立相同名稱

- **WHEN** 已有名稱為 group-a 的輪詢群組，呼叫端以另一 ID 建立相同名稱
- **THEN** 儲存庫回傳錯誤，原群組保持不變，列表中仍只有一個 group-a

##### Example: 相同名稱不會新增第二筆

- **GIVEN** first 的名稱為 same-name，interval 為 1000
- **WHEN** duplicate 以 same-name、interval 2000 建立
- **THEN** 回傳錯誤，記憶體儲存庫仍只有 first 一筆

#### Scenario: 拒絕改成其他群組的名稱

- **WHEN** 群組 B 嘗試改成群組 A 已使用的名稱
- **THEN** 儲存庫回傳錯誤，群組 B 的原名稱與其他欄位保持不變

##### Example: 衝突不改動其他欄位

- **GIVEN** first 的名稱為 same-name，second 的名稱為 other-name、interval 為 1000
- **WHEN** second 嘗試改成 same-name、interval 2000
- **THEN** 回傳錯誤，second 仍為 other-name、interval 1000

#### Scenario: 允許保留自己的名稱

- **WHEN** 群組保留自己的名稱並更新其他欄位
- **THEN** 更新不會被當成名稱衝突

##### Example: 自身名稱不阻擋更新

- **GIVEN** first 的名稱為 same-name、interval 為 1000
- **WHEN** first 保留 same-name，將 interval 改為 3000
- **THEN** 更新成功且持久值為 3000

<!-- @trace
source: restore-skipped-handler-tests
updated: 2026-11-07
code:
  - internal/datalink/sourcerule/acquisition_source_mixed.go
  - internal/datalink/collector/scheduler_auto_interval.go
  - internal/datalink/migrator_device_acquisition.go
  - internal/datalink/dbtarget/writer_group_flush.go
  - internal/datalink/sourcerule/point_write_journal.go
  - internal/datalink/optimizer/acquisition_builder.go
  - internal/datalink/sourcerule/mixed_rule_plan.go
  - internal/datalink/connector/adapters/mc3e_planned_scatter.go
  - internal/datalink/dbtarget/writer_recovery.go
  - internal/datalink/workspace/recording_home_sources.go
  - internal/datalink/anomaly/sql_repo.go
  - internal/protocol/mcprotocol/random_read.go
  - internal/datalink/runtime/ingestor.go
  - internal/datalink/collector/scheduler_config_methods.go
  - internal/datalink/point/format_applicability_update.go
  - internal/datalink/optimizer/acquisition_random.go
  - internal/api/handlers/studio_v2_workspace_audit.go
  - internal/datalink/runtimeview/contract_effective_validation.go
  - internal/datalink/pollinggroup/memory_repo.go
  - internal/datalink/dbtarget/writer_snapshot.go
  - internal/datalink/runtimeview/contract_mixed.go
  - docs/testing/skip-audit-20261107.md
tests:
  - internal/datalink/dbtarget/destination_open_test.go
  - internal/datalink/dbtarget/writer_skip_reason_test.go
  - internal/datalink/dbtarget/writer_snapshot_test.go
  - internal/api/handlers/device_activation_fixture_test.go
  - internal/api/router_studio_v2_successor_activation_postgres_proof_test.go
  - internal/datalink/dbtarget/writer_worker_test.go
  - internal/datalink/dbtarget/writer_decouple_timeout_regression_test.go
  - internal/datalink/anomaly/sql_repo_empty_list_test.go
  - internal/datalink/dbtarget/writer_recovery_lock_test.go
  - internal/datalink/connector/adapters/mc3e_planned_matrix_test.go
  - cmd/test_ui/target_writer_throughput_test.go
  - internal/api/handlers/studio_v2_workspace_source_rules_mixed_test.go
  - internal/datalink/dbtarget/writer_recovery_test.go
  - internal/protocol/mcprotocol/random_read_bit_lint_regression_test.go
  - internal/api/router_studio_v2_successor_sqlserver_proof_test.go
  - internal/datalink/dbtarget/writer_recovery_postgres_test.go
  - internal/api/handlers/device_acquisition_handler_test.go
  - internal/datalink/anomaly/sql_repo_rollback_test.go
  - internal/datalink/real_device_integration_test.go
  - internal/datalink/dbtarget/service_writer_test.go
  - cmd/test_ui/target_writer_lifecycle_test.go
  - internal/datalink/dbtarget/delivery_outcome_test.go
  - internal/api/router_studio_v2_successor_mysql_proof_test.go
  - internal/datalink/pollinggroup/memory_repo_name_test.go
  - internal/api/handlers/mixed_item_identity_regression_test.go
  - internal/datalink/dbtarget/writer_decouple_group_regression_test.go
  - internal/datalink/dbtarget/writer_interval_test.go
  - internal/api/handlers/device_handler_extended_test.go
  - internal/api/handlers/polling_group_handler_test.go
  - cmd/test_ui/bug_a_activation_test.go
  - internal/datalink/anomaly/sql_repo_retention_regression_test.go
  - internal/datalink/dbtarget/tooling_service_test.go
  - internal/api/handlers/device_acquisition_profile_regression_test.go
  - internal/api/handlers/polling_group_duplicate_test.go
-->