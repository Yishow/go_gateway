## Why

實際 collector 逐 Point 讀取；MC／FATEK 低階 random-read 與舊 optimizer 並未接上 production。使用者要求將 D0、D10、D20 的合批、請求數預覽納入 change，並加深優化；不能把協議能力或舊規格誤報為現有性能。

## What Changes

- 在既有 collector 接入單一 protocol-aware planner：MC／FATEK 能力允許時 random read，Modbus 依已證明可讀覆蓋合併連續區間，禁止盲目跨 gap。
- 同輪同設備去重 wire spans，再依各 Point 身分與有效字序分派；保持採集集合、cadence、品質與真實觀測時間。
- 可取消、有界重試與公平派送，慢設備不阻塞其他設備；immutable plan/revision 防止舊回應套新設定。
- 同 planner 提供不發 PLC 請求的預覽，區分計畫請求數、實際嘗試／重試及非 TCP 封包數。
- 用可重現 simulator 比較 baseline 與 optimized 請求數、cycle latency、freshness、值正確性，公布量測條件，不承諾未量測速度。

## Capabilities

### New Capabilities

無；不建立第二採集引擎。

### Modified Capabilities

- `collection-optimizer`: 替換僅啟動／max-hole 的舊合併契約，加入安全覆蓋、random、去重、版本計畫與同源預覽。
- `collection-scheduler`: 合批結果品質／時間、bounded retry、公平取消與版本隔離。

## Impact

主要影響 internal/datalink/optimizer、collector、connector/adapters、runtime 與 cmd/test_ui/service_wiring.go；重用 internal/protocol/mcprotocol/client.go、internal/protocol/fatek/client.go、internal/protocol/modbus/client.go。預覽接既有 internal/api/handlers 與 frontend/src，不新開採集入口。

preserve-per-point-word-order 是 decoder contract 前置；本案保留其 legacy/effective format 與 invalidation，不修改該案。streamline-recording-task-editor、support-per-member-recording-calculation、adopt-simple-recording-entry 只消費採集事實，另由 owner 加依賴，不由本案改其 artifacts。舊 collection-optimizer 規格存在但 runtime 尚未接線；本案以明確 MODIFIED delta 修正，不能宣稱已完成。僅提案，不授權 apply、PLC／正式 DB、部署。
