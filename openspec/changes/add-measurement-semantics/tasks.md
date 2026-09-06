## 1. 契約與測試先行
- [ ] 1.1 以共用契約凍結 MeasurementDefinition/SampleEnvelope 與 mixed layout JSON，建立 fixtures。
- [ ] 1.2 Red：新增 int16/uint16/uint32/uint64、字序、邊界重疊、bit slice、倍率只套一次的單元測試。
- [ ] 1.3 Red：補 measurement rename、source replace、epoch 與 JSON 大整數往返測試。

## 2. 後端完整落地
- [ ] 2.1 建立加法式 schema migration/repository/ownership 驗證與保留 legacy 的遷移測試。
- [ ] 2.2 Green：實作逐項有效格式解析、span 保留與相容讀取合併，重用既有協議 parser。
- [ ] 2.3 Green：串接 Point/Tag/Mapping 候選、revision 與有效換算；拒絕部分更新造成的下游不一致。
- [ ] 2.4 實作 sample identity、時間來源、品質與型別編碼；單項失敗保留其他項目。
- [ ] 2.5 實作範本版次、批次預覽與接受；未知現場資訊不可直接生效。
- [ ] 2.6 Refactor：拆開 layout、semantic、identity 與 serialization 模組，維持既有 Go 閘門。

## 3. API、前端契約與驗收
- [ ] 3.1 新增 workspace-scoped measurement API、typed error、revision CAS 與 Swagger；UI types/services/hooks 同步。
- [ ] 3.2 加入三相電表與 MC 混合 fixture 的 protocol→point→sample 整合測試。
- [ ] 3.3 確認 legacy uniform、rule enable/disable、Local Modbus geometry/probe gate 無回歸。
- [ ] 3.4 執行 go test/vet/golangci-lint；前端相關契約測試、lint/build；記錄未能執行的平台。
- [ ] 3.5 執行本 change OpenSpec strict 驗證、更新實作證據與相依契約；所有反例通過才開功能。
