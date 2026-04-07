## 1. Schema 與契約

- [x] 1.1 在 `source_rules`（或專案實際表名）新增可空 `target_data_type`（或對應 snake_case）欄位與 migration；預設 NULL 表示與 Point 讀取型別相同
- [x] 1.2 新增可空 **scale** 相關欄位（倍率／偏移，結構對齊 `TransformParamsScale`）與 migration
- [x] 1.3 更新 Go `schema`、repository、API DTO（create/update/get list）與 Swagger／OpenAPI 註解
- [x] 1.4 定義並文件化「支援的 cast 對」矩陣與 **cast → scale** 管線順序（重用 `mapping` 驗證）

## 2. 規則同步與 Tag／Mapping

- [x] 2.1 調整 `sourcerule`：`resolveRuleTag` 等路徑在宣告目標型別時允許 Tag 使用目標型別、Point 維持協議型別
- [x] 2.2 建立／更新 Mapping 時：依規則插入 **`cast`** 與／或 **`scale`** 步驟；空管線與手動管線衝突策略依 design
- [x] 2.3 補齊單元／整合測試：`uint16`→`float64`、僅 scale、ingestor 執行結果

## 3. 前端 Workbench（規則）

- [x] 3.1 更新 TypeScript 型別與 `useSourceRules` 以傳遞 `target_data_type` 與 scale 欄位
- [x] 3.2 在 `SourceCanvasSection` 規則表單新增「目標資料型態」與可選「縮放／偏移」、i18n
- [x] 3.3 補 Vitest／Workbench 相關測試

## 4. 連線 `data_format`（跨協議）

- [x] 4.1 在 `internal/datalink/connector/registry.go` 為 **Modbus TCP／RTU／UDP** 的 `ConfigSchema` 增加 **`data_format`**（enum 與 MC3E 一致），預設值與現行解碼預設對齊
- [x] 4.2 Modbus 適配器讀取 `connection_config.data_format` 並傳入 `lib/hsllogic`（或等價）解碼路徑；補多字組 float／int32 單元測試
- [ ] 4.3 前端裝置連線表單（動態 schema）驗證 Modbus 與 MC3E 皆能選擇並儲存 `data_format`

## 5. 文件與驗收

- [ ] 5.1 在 `docs/` 或 `design.md` 延伸處整理「連線／Point／規則／映射」參數分層表（與本變更 proposal 一致）
- [ ] 5.2 手動驗收：① `data_format` 變更後 float 讀值正確；② `uint16`→`float64` cast；③ scale 單獨與合併
- [ ] 5.3 更新本 `tasks.md` 勾選進度
