# Go Gateway 專案架構與流程深度分析

## 目錄
1. [系統架構圖](#系統架構圖)
2. [新增設備完整流程](#新增設備完整流程)
3. [資料流向圖](#資料流向圖)
4. [流程缺失分析](#流程缺失分析)
5. [改進建議](#改進建議)

---

## 系統架構圖

### 整體架構

```mermaid
graph TB
    subgraph "前端層 (Frontend)"
        UI[React UI<br/>TypeScript]
        Pages[頁面組件<br/>DevicesPage, PointsPage, TagsPage, MappingsPage]
        Components[UI 組件<br/>DeviceForm, MappingWizard]
    end

    subgraph "API 層 (API Layer)"
        Router[Gin Router<br/>RESTful API]
        Handlers[HTTP Handlers<br/>Device/Point/Tag/Mapping]
    end

    subgraph "業務邏輯層 (Business Logic)"
        DeviceService[Device Service<br/>設備管理]
        PointService[Point Service<br/>點位管理]
        TagService[Tag Service<br/>標籤管理]
        MappingService[Mapping Service<br/>映射管理]
        Scheduler[Collection Scheduler<br/>資料收集排程器]
    end

    subgraph "連接層 (Connection Layer)"
        ConnMgr[Connection Manager<br/>連線池管理]
        ProtocolAdapters[協議適配器<br/>Modbus/FATEK/MC3E/MQTT]
    end

    subgraph "資料層 (Data Layer)"
        Repositories[Repository Pattern<br/>SQL Repository]
        DB[(資料庫<br/>SQLite/PostgreSQL)]
    end

    subgraph "儲存層 (Storage Layer)"
        TimeSeriesWriter[TimeSeries Writer<br/>時序資料寫入器]
        BatchWriter[Batch Writer<br/>批次寫入緩衝]
    end

    subgraph "外部系統 (External Systems)"
        PLC[PLC 設備<br/>Modbus/FATEK/MC3E]
        MQTTBroker[MQTT Broker]
        ExternalDB[外部資料庫<br/>CMS/其他系統]
    end

    UI --> Router
    Pages --> Router
    Components --> Router
    
    Router --> Handlers
    Handlers --> DeviceService
    Handlers --> PointService
    Handlers --> TagService
    Handlers --> MappingService
    
    DeviceService --> Repositories
    PointService --> Repositories
    TagService --> Repositories
    MappingService --> Repositories
    Repositories --> DB
    
    Scheduler --> ConnMgr
    Scheduler --> Repositories
    ConnMgr --> ProtocolAdapters
    ProtocolAdapters --> PLC
    ProtocolAdapters --> MQTTBroker
    
    Scheduler --> TimeSeriesWriter
    TimeSeriesWriter --> BatchWriter
    BatchWriter --> DB
    BatchWriter --> ExternalDB
```

### 核心模組關係

```mermaid
graph LR
    subgraph "設備生命週期"
        Device[Device<br/>設備]
        DeviceStatus[Status<br/>draft/active/disabled]
    end

    subgraph "資料收集管線"
        Point[Point<br/>點位]
        PollingGroup[Polling Group<br/>輪詢群組]
        Mapping[Mapping<br/>映射]
        Tag[Tag<br/>標籤]
    end

    subgraph "資料流"
        CollectedValue[Collected Value<br/>收集值]
        TransformedValue[Transformed Value<br/>轉換值]
        TimeSeriesRecord[TimeSeries Record<br/>時序記錄]
    end

    Device -->|1:N| Point
    Point -->|N:1| PollingGroup
    Point -->|1:1| Mapping
    Mapping -->|N:1| Tag
    
    Point -->|產生| CollectedValue
    CollectedValue -->|轉換| TransformedValue
    TransformedValue -->|儲存| TimeSeriesRecord
```

---

## 新增設備完整流程

### 完整工作流程圖

```mermaid
flowchart TD
    Start([開始新增設備]) --> CreateDevice[1. 建立設備<br/>POST /devices]
    CreateDevice --> DeviceDraft[設備狀態: draft<br/>未啟用]
    
    DeviceDraft --> TestConnection{2. 測試連線<br/>POST /devices/:id/test}
    TestConnection -->|失敗| FixConfig[修正連線配置]
    FixConfig --> TestConnection
    TestConnection -->|成功| UpdateTestResult[更新測試結果<br/>last_test_success = true]
    
    UpdateTestResult --> ActivateDevice{3. 啟用設備<br/>POST /devices/:id/activate}
    ActivateDevice -->|失敗| CheckError[檢查錯誤]
    CheckError --> FixConfig
    ActivateDevice -->|成功| DeviceActive[設備狀態: active<br/>可參與資料收集]
    
    DeviceActive --> CreatePoints{4. 建立點位<br/>POST /points}
    CreatePoints -->|未建立| PointWarning[⚠️ 警告: 設備已啟用但無點位<br/>不會收集任何資料]
    CreatePoints -->|已建立| PointCreated[點位建立完成]
    
    PointCreated --> AssignPollingGroup{5. 指派輪詢群組<br/>point.polling_group_id}
    AssignPollingGroup -->|未指派| GroupWarning[⚠️ 警告: 點位未指派群組<br/>不會被排程收集]
    AssignPollingGroup -->|已指派| GroupAssigned[輪詢群組已指派]
    
    GroupAssigned --> CreateTag{6. 建立標籤<br/>POST /tags}
    CreateTag -->|未建立| TagWarning[⚠️ 警告: 無標籤無法建立映射]
    CreateTag -->|已建立| TagCreated[標籤建立完成]
    
    TagCreated --> CreateMapping{7. 建立映射<br/>POST /mappings<br/>或使用 MappingWizard}
    CreateMapping -->|未建立| MappingWarning[⚠️ 警告: 無映射無法儲存資料]
    CreateMapping -->|已建立| MappingCreated[映射建立完成]
    
    MappingCreated --> CheckScheduler{8. 檢查排程器狀態}
    CheckScheduler -->|未啟動| StartScheduler[啟動 Scheduler<br/>scheduler.Start]
    CheckScheduler -->|已啟動| SchedulerRunning[排程器運行中]
    
    StartScheduler --> SchedulerRunning
    SchedulerRunning --> DataCollection[9. 開始資料收集<br/>依輪詢群組間隔執行]
    
    DataCollection --> TransformData[10. 轉換資料<br/>執行 Transform Pipeline]
    TransformData --> StoreData[11. 儲存時序資料<br/>寫入 timeseries 表]
    
    StoreData --> Complete([✅ 流程完成<br/>資料開始流入])
    
    PointWarning --> Complete
    GroupWarning --> Complete
    TagWarning --> Complete
    MappingWarning --> Complete
    
    style DeviceDraft fill:#fff3cd
    style DeviceActive fill:#d4edda
    style PointWarning fill:#f8d7da
    style GroupWarning fill:#f8d7da
    style TagWarning fill:#f8d7da
    style MappingWarning fill:#f8d7da
    style Complete fill:#d4edda
```

### 狀態轉換圖

```mermaid
stateDiagram-v2
    [*] --> Draft: 建立設備
    
    Draft --> Testing: 測試連線
    Testing --> Draft: 測試失敗
    Testing --> Active: 測試成功 + 啟用
    
    Active --> Disabled: 停用設備
    Disabled --> Active: 重新啟用
    
    Active --> Collecting: 有啟用的點位 + 群組 + 映射
    Collecting --> Active: 移除點位/群組/映射
    
    note right of Draft
        設備已建立但未啟用
        無法參與資料收集
    end note
    
    note right of Active
        設備已啟用
        可參與資料收集
        但需要完整配置才能實際收集
    end note
    
    note right of Collecting
        設備正在收集資料
        需要滿足條件:
        1. 狀態 = active
        2. 至少一個 enabled 點位
        3. 點位指派到 enabled 群組
        4. 點位有對應的 enabled 映射
        5. Scheduler 正在運行
    end note
```

---

## 資料流向圖

### 資料收集與處理流程

```mermaid
sequenceDiagram
    participant Scheduler as Collection Scheduler
    participant ConnMgr as Connection Manager
    participant Protocol as Protocol Adapter
    participant PLC as PLC Device
    participant Point as Point Repository
    participant Mapping as Mapping Service
    participant Transform as Transform Pipeline
    participant Storage as TimeSeries Storage
    participant DB as Database

    Note over Scheduler: 輪詢群組 Ticker 觸發
    Scheduler->>Point: 查詢群組中的點位
    Point-->>Scheduler: 返回點位列表
    
    loop 每個點位
        Scheduler->>ConnMgr: GetOrCreate(device_id)
        ConnMgr->>Protocol: 建立/取得連線
        Protocol->>PLC: 讀取點位 (Read Request)
        PLC-->>Protocol: 返回原始資料 (Raw Bytes)
        Protocol-->>ConnMgr: 解析後的數值
        ConnMgr-->>Scheduler: CollectedValue
        
        Scheduler->>Mapping: 查詢點位的映射
        Mapping-->>Scheduler: 返回映射配置
        
        alt 有映射
            Scheduler->>Transform: 執行轉換管線
            Transform->>Transform: 執行各轉換步驟<br/>(decode, cast, scale, formula...)
            Transform-->>Scheduler: TransformedValue
            
            Scheduler->>Storage: WriteBatch(TimeSeriesRecord)
            Storage->>DB: 批次寫入 timeseries 表
        else 無映射
            Note over Scheduler: 資料被丟棄<br/>不儲存
        end
    end
```

### 映射轉換流程

```mermaid
flowchart LR
    RawValue[原始值<br/>Raw Bytes] --> Decode[解碼步驟<br/>decode]
    Decode --> Cast[型別轉換<br/>cast]
    Cast --> Scale[縮放轉換<br/>scale + offset]
    Scale --> Lookup[查詢表<br/>lookup_table]
    Lookup --> Formula[公式計算<br/>formula expression]
    Formula --> Conditional[條件規則<br/>conditional]
    Conditional --> FinalValue[最終值<br/>Transformed Value]
    
    FinalValue --> Store[儲存到<br/>TimeSeries]
    
    style RawValue fill:#fff3cd
    style FinalValue fill:#d4edda
    style Store fill:#cfe2ff
```

---

## 流程缺失分析

### 1. 缺少明確的引導流程

**問題描述：**
- 新增設備後，系統沒有明確指引下一步該做什麼
- 使用者不知道需要完成哪些步驟才能開始收集資料
- 缺少「快速開始」或「設定精靈」功能

**影響：**
- 使用者可能啟用設備後就停止，不知道還需要建立點位、標籤、映射
- 導致設備處於 active 狀態但實際上沒有收集任何資料

**證據：**
- `DeviceForm` 只有基本的設備建立功能
- 沒有「下一步建議」或「完成度檢查」
- `DevicesPage` 沒有顯示設備的「配置完成度」

### 2. 缺少狀態驗證與檢查

**問題描述：**
- 設備啟用時只檢查連線，不檢查是否有點位、群組、映射
- 沒有「就緒檢查」(Readiness Check) 機制
- 缺少配置完整性驗證

**影響：**
- 設備可能處於 active 狀態但無法實際收集資料
- 使用者不知道缺少哪些配置

**證據：**
```go
// internal/datalink/device/service.go:214
func (s *Service) Activate(ctx context.Context, id string) error {
    // 只檢查連線，不檢查點位、群組、映射
    if err := s.TestConnection(ctx, id); err != nil {
        return fmt.Errorf("連線測試失敗，無法啟用設備: %w", err)
    }
    // ...
}
```

### 3. 缺少自動化建議

**問題描述：**
- 沒有根據設備類型自動建議點位配置
- 沒有根據協議類型提供預設配置模板
- 缺少「一鍵設定」功能

**影響：**
- 使用者需要手動輸入所有配置，容易出錯
- 新手使用者不知道如何開始

### 4. 缺少即時狀態反饋

**問題描述：**
- 設備啟用後，沒有即時顯示「是否正在收集資料」
- 沒有「最後收集時間」或「收集統計」資訊
- 缺少資料流狀態監控

**影響：**
- 使用者無法確認設備是否正常運作
- 無法快速發現配置問題

**證據：**
- `Device` 模型沒有 `last_collected_at` 欄位
- `DevicesPage` 沒有顯示收集狀態

### 5. 缺少錯誤處理與恢復機制

**問題描述：**
- 資料收集失敗時，沒有自動重試機制（雖然 Scheduler 有重試，但缺少更高層的恢復）
- 連線中斷後，缺少自動恢復策略
- 沒有錯誤通知機制

**影響：**
- 長時間的連線問題可能導致資料遺失
- 使用者無法及時發現問題

### 6. 缺少批次操作支援

**問題描述：**
- 無法批次建立多個點位
- 無法批次建立映射
- 缺少「從模板匯入」功能

**影響：**
- 大量設備配置時效率低下
- 重複性工作增加

### 7. 缺少配置驗證與預覽

**問題描述：**
- 建立映射時缺少「測試映射」功能（雖然有 preview，但缺少在建立前的驗證）
- 沒有「配置影響分析」（例如：修改群組間隔會影響哪些點位）

**影響：**
- 配置錯誤可能導致資料收集異常
- 無法預先評估配置變更的影響

---

## 改進建議

### 優先級 P0（立即實施）

#### 1. 新增設備配置完成度檢查

**實施方案：**
- 在 `Device` 模型中新增 `readiness_status` 欄位
- 實作 `DeviceService.CheckReadiness()` 方法，檢查：
  - ✅ 設備狀態 = active
  - ✅ 至少一個 enabled 點位
  - ✅ 點位指派到 enabled 群組
  - ✅ 點位有對應的 enabled 映射
  - ✅ Scheduler 正在運行

**API 端點：**
```go
GET /api/v1/datalink/devices/:id/readiness
Response: {
    "ready": bool,
    "checks": {
        "device_active": bool,
        "has_points": bool,
        "has_polling_groups": bool,
        "has_mappings": bool,
        "scheduler_running": bool
    },
    "missing": []string,  // 缺少的配置項目
    "suggestions": []string  // 建議的下一步操作
}
```

**前端顯示：**
- 在 `DeviceCard` 中顯示完成度進度條
- 顯示缺少的配置項目清單
- 提供「快速完成配置」按鈕

#### 2. 新增設備後引導流程

**實施方案：**
- 建立 `DeviceOnboardingWizard` 組件
- 步驟：
  1. 建立設備 ✅
  2. 測試連線 → 啟用設備
  3. 建立點位（可批次）
  4. 指派輪詢群組
  5. 建立標籤（可批次）
  6. 建立映射（使用現有 MappingWizard）
  7. 確認開始收集

**UI 流程：**
```typescript
// 新增設備成功後，自動開啟引導流程
const handleCreateSuccess = (device: Device) => {
    showOnboardingWizard(device.id);
};
```

#### 3. 設備狀態儀表板

**實施方案：**
- 新增 `DeviceStatusDashboard` 頁面
- 顯示：
  - 設備總數、啟用數、收集中數
  - 每個設備的收集狀態（最後收集時間、收集次數、錯誤次數）
  - 配置完成度統計

**資料模型擴充：**
```sql
ALTER TABLE devices ADD COLUMN last_collected_at TIMESTAMPTZ;
ALTER TABLE devices ADD COLUMN collection_count BIGINT DEFAULT 0;
ALTER TABLE devices ADD COLUMN error_count BIGINT DEFAULT 0;
```

### 優先級 P1（短期實施）

#### 4. 配置模板與批次匯入

**實施方案：**
- 建立「設備配置模板」功能
- 支援 JSON/YAML 格式的批次匯入
- 提供常見設備類型的預設模板

**API 端點：**
```go
POST /api/v1/datalink/devices/:id/apply-template
POST /api/v1/datalink/devices/batch-import
GET /api/v1/datalink/templates
```

#### 5. 即時資料流監控

**實施方案：**
- 使用 WebSocket 或 SSE 推送即時收集狀態
- 顯示每個設備的即時資料流
- 提供「資料流預覽」功能

**API 端點：**
```go
GET /api/v1/datalink/devices/:id/stream  // SSE
GET /api/v1/datalink/devices/:id/stats  // 統計資訊
```

#### 6. 配置驗證與影響分析

**實施方案：**
- 建立 `ConfigValidator` 服務
- 在修改配置前進行驗證
- 顯示配置變更的影響範圍

**API 端點：**
```go
POST /api/v1/datalink/config/validate
POST /api/v1/datalink/config/impact-analysis
```

### 優先級 P2（中期實施）

#### 7. 自動化配置建議

**實施方案：**
- 根據設備協議類型自動建議點位配置
- 使用機器學習或規則引擎提供配置建議
- 提供「一鍵完成配置」功能

#### 8. 錯誤恢復與通知機制

**實施方案：**
- 實作更完善的錯誤恢復策略
- 新增錯誤通知系統（Email/Webhook）
- 提供錯誤日誌查詢與分析

#### 9. 配置版本管理

**實施方案：**
- 實作配置版本控制
- 支援配置回滾
- 提供配置變更歷史記錄

---

## 實施路線圖

### 階段一：基礎改進（1-2 週）
1. ✅ 實作設備配置完成度檢查
2. ✅ 新增設備後引導流程
3. ✅ 設備狀態儀表板

### 階段二：增強功能（2-3 週）
4. ✅ 配置模板與批次匯入
5. ✅ 即時資料流監控
6. ✅ 配置驗證與影響分析

### 階段三：進階功能（3-4 週）
7. ✅ 自動化配置建議
8. ✅ 錯誤恢復與通知機制
9. ✅ 配置版本管理

---

## 總結

### 當前狀態
- ✅ 核心功能完整：設備、點位、標籤、映射、收集都已實作
- ⚠️ 缺少使用者引導：新增設備後不知道下一步
- ⚠️ 缺少狀態驗證：無法確認配置是否完整
- ⚠️ 缺少即時反饋：無法監控資料收集狀態

### 關鍵改進點
1. **引導流程**：新增設備後自動引導完成所有必要配置
2. **狀態檢查**：實作配置完成度檢查，明確顯示缺少的項目
3. **即時監控**：提供資料收集狀態的即時反饋
4. **批次操作**：支援批次建立點位、標籤、映射，提升效率

### 預期效果
- 📈 使用者體驗提升：明確的引導流程，減少配置錯誤
- 📈 系統可靠性提升：配置驗證與狀態檢查，減少運行時錯誤
- 📈 運維效率提升：即時監控與批次操作，減少人工干預
