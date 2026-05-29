## Context

當前 Datalink Workbench V2 提供了完整的設備配置、規則格狀編輯、Tag 映射與資料庫綁定流。然而，系統缺乏一個直觀、實時的運作監控儀表板。用戶在配置完成並提交後，無法直接在系統內監看實時的採集吞吐量、輪詢延遲、點位即時轉換值，以及資料庫寫入佇列是否積壓。

為了填補這一空白，本設計提出在 `/studio/dashboard` 路由下實作 Runtime Dashboard，作為運作狀態的主控面板，提供實時數據監控與故障排障介面。

## Goals / Non-Goals

**Goals:**

- 實作單頁的實時運行儀表板 `/studio/dashboard`，整合於 Workbench V2 的 Shell 佈局中。
- 採用 SSE (Server-Sent Events) 單向推送技術，讓前端能以低延遲接收後端發送的即時數據。
- 實作五大數據區塊：採集績效指標 (Metrics)、設備連線網格 (Device Grid)、點位實時數值表 (Live Values)、寫入佇列積壓 (Queue Backlog) 與實時診斷日誌 (Diagnostic Logs)。
- 建立連線異常或積壓過高時的視覺警告機制。

**Non-Goals:**

- 實作歷史趨勢圖或複雜的 Recharts 折線圖（本階段只做實時值展示，歷史報表不在本次範疇內）。
- 提供下發控制 PLC 暫存器寫入的功能（此儀表板僅限於唯讀監控，不包含寫入操作）。

## Decisions

### Decision: 使用 SSE (Server-Sent Events) 進行實時數據推送

- **Rationale**: 相比於 WebSocket，SSE 更加輕量化，且專為「伺服器單向推送，前端僅讀取」的場景設計，非常符合儀表板即時更新的需求。此外，後端 Go/Gin 本身對 SSE 的支援非常簡潔，能重用現有的 HTTP 連線管理與憑證，開發與維運成本低。
- **Alternatives considered**:
  - *WebSocket*: 提供雙向通訊，但協定較為複雜，且我們不需要從前端高頻發送控制命令，使用 WebSockets 會增加不必要的連線管理開銷。
  - *HTTP Polling (輪詢)*: 前端定時發送 GET 請求拉取狀態。這會對嵌入式閘道器（如 Raspberry Pi）的 CPU 與網絡頻寬造成持續的輪詢開銷，在大點位數量時尤為嚴重。

### Decision: 前端監控元件模組化與狀態隔離

- **Rationale**: 為了保持前端檔案行數嚴格符合 `≤ 300` 行的專案規範，我們拒絕將所有區塊寫在單一元件中。`RuntimeDashboard` 將作為 SSE 連線管理與數據分流的 Container，而指標、設備、點位表格、佇列、日誌等則拆分為獨立的 Presentation components。
- **Alternatives considered**:
  - *單一長檔案實作*: 將所有 UI 邏輯與 SSE 連線寫在一起。雖然開發快速，但檔案極易突破 500 行硬上限，且無法針對單一區塊進行單元測試，可維護性極差。

## Implementation Contract

- **Behavior**: 使用者點擊導覽列中的「監控面板」後，系統載入 `/studio/dashboard`。頁面將展示五大區塊並實時跳動數值。當設備斷線，對應卡片轉為紅色並顯示錯誤；當佇列積壓 > 500 筆，佇列指標轉為黃色並閃爍。
- **Interface / Data Shape**:
  - 後端 SSE 端點 `/api/v1/datalink/runtime/stream` 預期推送以下 Event 類型：
    - `event: metrics` -> `{"points_per_sec": number, "success_rate": number, "avg_latency_ms": number}`
    - `event: devices` -> `[{"id": string, "name": string, "status": "online"|"offline"|"testing", "latency_ms": number}]`
    - `event: live_values` -> `{"point_id": string, "raw": number|string, "transformed": number|string, "timestamp": string}`
    - `event: queue_status` -> `{"db_backlog": number, "modbus_backlog": number}`
    - `event: logs` -> `{"timestamp": string, "level": "error"|"warn", "message": string}`
- **Failure modes**: 當 SSE 連線中斷時，Dashboard 頂部會顯示黃色連線中斷 Banner，並於 5 秒後自動重連。重連期間畫面呈半透明遮罩，保留最後已知數據防止畫面閃爍。
- **Acceptance criteria**:
  - 新增單元測試 `runtime-dashboard.test.tsx`。
  - 斷言接收到 `metrics` 事件時，畫面上即時值（如吞吐量、延遲）能正確更新。
  - 斷言當 `db_backlog` 超過 500 時，畫面上會出現琥珀色警告（amber warning）與對應 Alert 訊息。
- **Scope boundaries**: 僅限於本端即時監控數據的展現，不包括控制命令的發送，也不包括歷史報表（例如 Chart/Graph 曲線圖，本階段只做數值文字顯示與基本指標卡片）。

## Risks / Trade-offs

- `[Risk]` SSE 連線因瀏覽器後台標籤頁凍結而斷線，重新喚醒時數據脫節 -> `[Mitigation]` 實作 Page Visibility API 監聽，當頁面由 hidden 轉為 visible 時，強制觸發 SSE 重連，刷新所有狀態。
- `[Risk]` 高頻點位更新導致 DOM 頻繁重繪 (Repaint/Reflow)，造成網頁卡頓 -> `[Mitigation]` 對實時點位表格的更新加上 Throttle（限制每 100ms 最多重繪一次），確保低配硬體渲染流暢。
