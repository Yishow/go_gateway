# `/gateway/*` Experimental Surfaces

## 定位

- Routes:
  - `/gateway/entry`
  - `/gateway/quick-setup`
  - `/gateway/expert-workbench`
- 狀態: feature-flag 保護下的原型 / 實驗入口
- 目前策略: **僅記錄，不列為近期主產品重點**

## 核心事實

1. 這組頁面目前不是 datalink persisted 主產品流程。
2. 它們會讀 `ENABLE_GATEWAY_DUAL_ENTRY` feature flag。
3. Quick / Expert 最後的 submit 目前主要走 `/api/v1/test/connect`，屬於 connect-only。

## 對應 API

| 行為 | API | 狀態 |
| --- | --- | --- |
| dual-entry feature flag | `GET /api/v1/datalink/settings` | `wired` |
| quick submit | `POST /api/v1/test/connect` | `connect-only` |
| expert submit | `POST /api/v1/test/connect` | `connect-only` |

## 維護提醒

1. 不要把 `/gateway/*` 的 submit 成功誤判為「已完成 datalink setup」。
2. 若未來要轉正，應先重新定義它與 `/studio/v2` 的關係。
3. 若只是 prototype，則它應繼續被文件標記為 experimental，而不是 product default path。
