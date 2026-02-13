# OpenSpec 草案總結：Add Device Onboarding and Readiness Check

## 📋 草案概覽

**Change ID**: `add-device-onboarding`  
**狀態**: 草案（待審核）  
**建立日期**: 2026-01-21  
**基於**: 深度架構分析報告

## 🎯 目標

解決當前系統在設備配置流程中的三個關鍵問題：

1. **缺少使用者引導** → 實作 Onboarding Wizard
2. **缺少狀態驗證** → 實作 Readiness Check
3. **缺少即時反饋** → 實作 Status Dashboard
4. **Sidebar 導航不完整** → 新增 Points 頁面並優化導航結構

## 📁 文件結構

```
openspec/changes/add-device-onboarding/
├── proposal.md          # 提案說明（Why/What/Impact）
├── design.md            # 技術設計決策
├── tasks.md             # 實施任務清單（含 TDD 工作流程）
└── specs/
    ├── device-registry/
    │   └── spec.md      # 設備註冊規格變更
    ├── datalink-api/
    │   └── spec.md      # API 規格變更
    └── datalink-ui/
        └── spec.md      # UI 規格變更
```

## 🔄 TDD 工作流程整合

所有實施任務都遵循 **RED-GREEN-REFACTOR** 循環：

### 範例：Readiness Check Service

```markdown
## 3. Readiness Check Service（TDD）

- [ ] **RED**: 撰寫 Readiness Check 單元測試
  - [ ] 測試設備狀態檢查
  - [ ] 測試點位存在性檢查
  - [ ] ...（更多測試案例）
  - [ ] 執行測試確認失敗（RED）

- [ ] **GREEN**: 實作 Readiness Check 服務
  - [ ] 實作各項檢查邏輯
  - [ ] 執行測試確認通過（GREEN）

- [ ] **REFACTOR**: 優化程式碼
  - [ ] 提取重複邏輯
  - [ ] 優化查詢效能
  - [ ] 確保測試仍通過
```

## 📊 規格變更摘要

### 1. Device Registry (`device-registry/spec.md`)

**新增需求**:
- ✅ Device readiness check
- ✅ Device collection statistics
- ✅ Device onboarding wizard
- ✅ Device status dashboard

**修改需求**:
- 🔄 Device registry（擴充收集統計）

### 2. Datalink API (`datalink-api/spec.md`)

**新增需求**:
- ✅ Device readiness API (`GET /devices/:id/readiness`)
- ✅ Device collection statistics API (`GET /devices/:id/stats`)
- ✅ Device status stream API (`GET /devices/stream`)

### 3. Datalink UI (`datalink-ui/spec.md`)

**新增需求**:
- ✅ Device onboarding wizard 組件
- ✅ Device status dashboard 頁面
- ✅ Device readiness indicator
- ✅ Sidebar navigation improvements
- ✅ Points management page

## 🛠️ 實施任務（13 個主要任務）

1. **資料模型擴充**（TDD）
2. **資料庫遷移**（TDD）
3. **Readiness Check Service**（TDD）
4. **Readiness Check API**（TDD）
5. **收集統計更新機制**（TDD）
6. **Onboarding Wizard 前端組件**（TDD）
7. **設備狀態儀表板**（TDD）
8. **整合測試**（TDD）
9. **文件更新**
10. **驗證與部署**
11. **Sidebar 導航改進**（TDD）
12. **Points 頁面實作**（TDD）
13. **路由配置更新**

## ✅ 符合 OpenSpec 規範

- ✅ 使用 `## ADDED|MODIFIED Requirements` 格式
- ✅ 每個需求至少包含一個 `#### Scenario:`
- ✅ Scenario 使用 `WHEN/THEN` 格式
- ✅ 包含 `proposal.md`, `tasks.md`, `design.md`
- ✅ 規格變更檔案位於 `specs/` 目錄

## 🧪 TDD 整合特點

1. **每個任務都包含 TDD 循環**
   - RED: 撰寫失敗測試
   - GREEN: 實作最小程式碼通過測試
   - REFACTOR: 優化程式碼

2. **測試優先級明確**
   - 單元測試優先
   - 整合測試其次
   - 端到端測試最後

3. **測試覆蓋範圍**
   - 成功場景
   - 失敗場景
   - 邊界條件
   - 錯誤處理

## 📈 預期效果

### 使用者體驗提升
- ✅ 清晰的引導流程，減少配置錯誤
- ✅ 即時狀態反饋，快速發現問題
- ✅ 配置完成度顯示，明確缺少的項目
- ✅ 完整的導航結構，所有核心功能可直接訪問
- ✅ 統一的 UI/UX 設計，提升操作效率

### 系統可靠性提升
- ✅ 配置驗證機制，減少運行時錯誤
- ✅ 即時監控，及時發現異常
- ✅ 統計資訊，支援問題診斷

### 開發效率提升
- ✅ TDD 工作流程，確保程式碼品質
- ✅ 清晰的任務分解，易於追蹤進度
- ✅ 完整的規格定義，減少溝通成本

## 🚀 下一步

1. **審核提案**
   - 檢視 `proposal.md` 確認範圍和影響
   - 檢視 `design.md` 確認技術決策
   - 檢視 `specs/` 確認規格完整性

2. **驗證規格**
   ```bash
   openspec validate add-device-onboarding --strict --no-interactive
   ```

3. **開始實施**
   - 按照 `tasks.md` 順序執行
   - 嚴格遵循 TDD 工作流程
   - 每個任務完成後更新狀態

## 📝 注意事項

1. **TDD 強制性**
   - 所有高複雜度邏輯必須遵循 TDD
   - 測試必須先失敗（RED）才能開始實作

2. **向下相容性**
   - 資料庫遷移必須向下相容
   - 現有 API 不應被破壞

3. **效能考量**
   - Readiness Check 需要優化查詢
   - SSE 推送需要控制頻率

4. **測試覆蓋率**
   - 目標：>80% 程式碼覆蓋率
   - 關鍵路徑：100% 覆蓋率

## 📚 相關文檔

- [架構分析報告](../docs/ARCHITECTURE_ANALYSIS.md)
- [新增設備流程指南](../docs/DEVICE_ONBOARDING_FLOW.md)
- [分析總結](../docs/ANALYSIS_SUMMARY.md)
- [OpenSpec 規範](../openspec/AGENTS.md)
- [TDD 工作流程](../../.agent/skills/tdd-workflow/SKILL.md)

---

**草案狀態**: ✅ 已完成  
**待審核**: 是  
**建議優先級**: P0（立即實施）
