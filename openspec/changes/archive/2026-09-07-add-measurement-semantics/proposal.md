# 為混合感測器建立逐項資料意義

## Why

MC D0～D7 可能同時放溫度、壓力、流量、累積數字與狀態；一段相同通訊範圍不代表一種記錄方式。
現有 Point/Mapping 提供型別、倍率與單位，但記錄用途與累積世代需要明確契約。
沒有這一層，資料庫只能「把數字塞進去」，無法可靠產生曲線或用量。

## What Changes

- 新增逐項 measurement 定義、用途、穩定身份、品質與樣本 envelope。
- 支援一段範圍包含不同占用長度與格式；保留原有均一型別規則。
- 對每個項目定義有效解碼、換算與用途，不把整段設定重複套倍率。
- 建立三相電表與混合感測器的版本化範本契約；未知位址與單位需確認。
- 既有 Point/Tag/Mapping、rule revision、Local Modbus ownership 保留，不另造平行來源。

## Capabilities

### New Capabilities
- `measurement-semantics`: 用途、品質、精度、身份與世代。
- `mixed-source-layout`: 逐項位址占用、解碼與範本。

### Modified Capabilities
- 無。既有均一規則保持不變；新增 mixed 模式的逐項覆寫優先於該模式的預設值。

## Impact

涉及 `internal/datalink/point`、`sourcerule`、`mapping`、`schema/migrations`、協議讀取計畫，
以及前端 types/services/hooks。新 measurement 模組與 API 為擬新增，需與既有模型關聯。
先實作此 change；2～6 依賴其穩定 sample 與身份契約。共用規格見 `docs/plans/telemetry-recording/contracts.md`。
