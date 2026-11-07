## Purpose

補齊既有 Polling Groups API 的名稱完整性：記憶體儲存庫與 SQL UNIQUE 約束一致，拒絕重複名稱且不破壞原資料。

## ADDED Requirements

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
