後端重構完整 Checklist 報告

建立時間: 2026-02-13T18:32:11.636Z

一、目標與現況

- 目標：所有 Go 生產碼檔案 <= 400 行
- 目標：降低高耦合（handler/service/adapter）與重複邏輯
- 現況：>400 行檔案共 35（生產碼 21、測試 14）

二、重構優先順序（P1→P3）

P1（先做，高影響）

- internal/api/handlers/test.go (1152)
  - 拆成：test_connection_handler.go、test_read_handler.go、test_write_handler.go、test_debug_handler.go
- internal/datalink/connector/adapters/modbus.go (923)
  - 拆成：modbus_connect.go、modbus_read.go、modbus_write.go、modbus_parse.go
- internal/datalink/mapping/service.go (844)
  - 拆成：mapping_validation.go、mapping_transform.go、mapping_execute.go
- internal/datalink/collector/scheduler.go (700)
  - 拆成：scheduler_loop.go、scheduler_dispatch.go、scheduler_retry.go

P2（第二批，穩定性/可維護）

- internal/datalink/schema/models.go (615) → schema model 分檔（device/tag/point/mapping）
- internal/datalink/device/service.go (567) → CRUD / state / validation 分離
- internal/datalink/tag/service.go (537) → query / lifecycle / binding 分離
- internal/datalink/point/service.go (520) → read/write / transform / history 分離
- internal/datalink/api/router.go (501) → route registry 分模組註冊

P3（第三批，技術債）

- internal/datalink/connector/adapters/mqtt.go (484)
- internal/datalink/storage/timeseries.go (472)
- internal/datalink/point/sql_repo.go (448)
- internal/datalink/connector/adapters/fatek.go (433)
- internal/datalink/connector/manager.go (426)
- internal/protocol/modbus/transport.go (424)
- internal/virtual/server/modbus/server.go (416)
- lib/hsllogic/*（多檔）

三、共通重構規範（每批都要勾）

- 單檔行數檢查（Go 生產碼全部 <= 400）
- 不改變對外 API 行為（路由、request/response 相容）
- 每批完成後跑：go test ./...
- 每批完成後跑：pnpm lint && pnpm exec tsc --noEmit（前端不回歸）
- 新增/調整最小必要測試，避免重構破壞流程

四、驗收標準（Definition of Done）

- 生產碼無檔案超過 400 行
- 全量測試通過（Go + Frontend）
- 無功能回歸、核心流程（device/tag/point/mapping/polling）可正常操作
- 每批有獨立 commit（可回滾）

五、推薦優先順序（實作建議）

1. 統一 API 回應封裝（Response Envelope）→ 建立 internal/api/handlers/response.go
2. 抽取 SQL Repository 共用基底（Generic BaseRepository）→ internal/datalink/repository/base.go
3. 協議 Adapter 初始化流程共用（BaseProtocolAdapter）→ internal/datalink/connector/adapters/base.go
4. 拆分大型 handler/service（按 P1 逐檔分割）
5. 清理 TODO 與補測試

六、估計步驟（無時間框）

- Step A: 建立 Response Envelope，更新 3 個最常用 handler 做範例，跑測試
- Step B: 建立 BaseRepository，重構 device repo 為樣板，跑測試
- Step C: 抽取 BaseProtocolAdapter，修正 modbus & fatek 初始化流程，跑測試
- Step D: 依 P1 列表逐檔拆分，確保每步都有測試覆蓋

七、風險與注意事項

- 拆檔時需小心包內的未導出 (unexported) 函數/變數，拆分前先重構為明確可導出的接口。
- 每次改動務必小步提交並執行 CI 測試。
- 測試套件中若有過度依賴檔案內部非公開狀態的測試，需先以黑盒方式重寫測試。

---

若要開始實作，請選：
- "從統一 API 回應封裝開始（推薦）"
- "從抽取 SQL Repository 開始"
- "從協議 Adapter 初始化流程開始"
- "先建立更詳細的實作計畫與變更集（逐檔）"

（若要我把此檔案另外 commit，回覆「commit」。）