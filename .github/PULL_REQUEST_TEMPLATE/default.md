## 變更摘要
- 

## 風險評估
- [ ] 低風險（文件/註解/非行為變更）
- [ ] 中風險（重構/測試調整）
- [ ] 高風險（流程邏輯/API 契約變更）

## 必要檢查（請勾選）
- [ ] 已閱讀並符合本專案規範
- [ ] 已附上測試命令與結果
- [ ] 已確認無非預期 UI 外觀變更
- [ ] 已確認無非預期 API 契約變更

### Frontend 變更時（必填）
- [ ] `cd frontend && npm test -- --run`
- [ ] `cd frontend && npm run lint`
- [ ] `cd frontend && npm run build`

### Backend 變更時（必填）
- [ ] `go test ./...`
- [ ] `go vet ./...`

## 驗收重點
- 

## 回滾方案
- 
