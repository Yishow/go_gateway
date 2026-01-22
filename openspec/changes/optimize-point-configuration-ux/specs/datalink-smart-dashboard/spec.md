# Spec: Datalink Smart Dashboard UI

> **Version**: 1.0  
> **Status**: Draft  
> **Related**: `device-registry`, `datalink-api`

## Overview

智慧 Dashboard 是資料點配置的一站式工作台，採用類 IDE 三欄佈局，將設備管理、點位配置、標籤映射整合於單一頁面。

---

## ADDED Requirements

### REQ-DASH-001: 三欄式佈局

Dashboard 頁面應採用三欄式佈局：設備樹狀導覽（左）、動態內容區（中）、快速操作面板（右）。

#### Scenario: 頁面載入

- **Given** 使用者開啟 SmartDashboard 頁面
- **When** 頁面完成載入
- **Then** 應顯示三欄佈局
- **And** 左側顯示設備樹狀導覽（240px）
- **And** 中間顯示設備選擇提示或記憶體格子
- **And** 右側顯示快速操作面板（280px）

#### Scenario: 響應式收合

- **Given** 視窗寬度小於 1280px
- **When** 頁面調整佈局
- **Then** 左側導覽應自動收合至圖標模式（64px）
- **And** 右側面板應改為浮動按鈕

---

### REQ-DASH-002: 設備樹狀導覽

樹狀導覽應顯示設備層級結構，支援展開/收合子項目。

#### Scenario: 設備列表顯示

- **Given** 系統有 3 個設備
- **When** 樹狀導覽載入
- **Then** 應顯示 3 個設備節點
- **And** 每個節點顯示設備名稱與圖標
- **And** 啟用設備顯示綠色指示點

#### Scenario: 展開設備子項目

- **Given** 設備 "PLC-001" 有 5 個點位和 3 個映射
- **When** 使用者點擊展開圖標
- **Then** 應顯示 "Points (5)" 子節點
- **And** 應顯示 "Mappings (3)" 子節點

#### Scenario: 選中設備

- **Given** 使用者點擊設備節點
- **When** 選中事件觸發
- **Then** 該節點應顯示選中樣式
- **And** 中間區域應載入該設備的記憶體格子

#### Scenario: 拖曳排序設備

- **Given** 使用者拖曳設備節點 "PLC-002"
- **When** 放置到 "PLC-001" 上方
- **Then** 設備順序應更新為 PLC-002, PLC-001, PLC-003
- **And** 應呼叫 `PATCH /api/v1/datalink/devices/reorder`
- **And** 順序應持久化到資料庫

#### Scenario: 拖曳視覺反饋

- **Given** 使用者開始拖曳設備節點
- **When** 拖曳進行中
- **Then** 被拖曳節點應顯示半透明效果
- **And** 目標位置應顯示藍色插入線
- **And** 拖曳手柄圖標應可見

---

### REQ-DASH-003: 記憶體格子視覺化

中間區域應顯示記憶體格子，視覺化 PLC 暫存器使用狀況。

#### Scenario: 格子渲染

- **Given** 設備協議為 MC3E，中心位址為 D100
- **When** 記憶體格子載入
- **Then** 應顯示 100 個格子（D050-D149）
- **And** 已使用格子顯示綠色
- **And** 可用格子顯示灰色
- **And** 每個格子顯示位址標籤

#### Scenario: 格子 Tooltip

- **Given** 格子 D100 已綁定點位 "Temp_Sensor"
- **When** 使用者懸停該格子
- **Then** 應顯示 Tooltip
- **And** Tooltip 包含點位名稱、最後讀取值、時間

#### Scenario: 單擊選取

- **Given** 使用者點擊可用格子 D105
- **When** 點擊事件觸發
- **Then** 該格子應顯示選中樣式
- **And** 側邊面板應開啟「新增點位」表單

#### Scenario: 範圍選取

- **Given** 使用者已選取 D100
- **When** 使用者 Shift+點擊 D109
- **Then** D100-D109 共 10 個格子應全部選中
- **And** 快速操作面板應顯示「選取 10 個位址」

---

### REQ-DASH-004: 快速操作面板

右側面板應提供常用操作按鈕與即時狀態監控。

#### Scenario: 批量建立按鈕

- **Given** 使用者選取 5 個格子
- **When** 快速操作面板顯示
- **Then** 「批量建立」按鈕應啟用
- **And** 按鈕應顯示「批量建立 (5)」

#### Scenario: 即時狀態顯示

- **Given** 設備 PLC-001 已連線
- **When** 快速操作面板載入
- **Then** 應顯示連線狀態指示燈（綠色）
- **And** 應顯示最後同步時間

---

### REQ-DASH-005: 側邊滑出面板

點選格子或操作按鈕時，應從右側滑出詳細設定面板。

#### Scenario: 開啟面板

- **Given** 使用者點擊「批量建立」按鈕
- **When** 面板開啟
- **Then** 面板應從右側滑入（300ms）
- **And** 背景應顯示半透明遮罩

#### Scenario: 關閉面板

- **Given** 面板已開啟
- **When** 使用者按 ESC 鍵
- **Then** 面板應滑出關閉
- **And** 遮罩應淡出

---

### REQ-DASH-006: 批量點位建立

批量建立器應支援一次建立多個連續位址的點位。

#### Scenario: 命名模板

- **Given** 起始位址 D100，數量 5，模板 "Sensor\_{index}"
- **When** 預覽生成
- **Then** 應顯示：Sensor_0, Sensor_1, Sensor_2, Sensor_3, Sensor_4
- **And** 對應位址：D100, D101, D102, D103, D104

#### Scenario: 批量建立成功

- **Given** 使用者填寫完整表單
- **When** 點擊「建立」按鈕
- **Then** 應呼叫 `POST /api/v1/datalink/points/batch`
- **And** 成功後顯示成功訊息
- **And** 記憶體格子應更新顯示新建立的點位

---

## ADDED API Requirements

### REQ-API-BATCH-001: 批量建立點位 API

#### Scenario: 成功批量建立

- **Given** 請求包含有效的 device_id 和 points 陣列
- **When** POST `/api/v1/datalink/points/batch`
- **Then** 應返回 201 Created
- **And** 回應包含 created_count 和 points 陣列

#### Scenario: 部分失敗回滾

- **Given** 請求包含 10 個點位，第 5 個位址衝突
- **When** POST `/api/v1/datalink/points/batch`
- **Then** 應返回 400 Bad Request
- **And** 已建立的 4 個點位應回滾
- **And** 回應包含錯誤詳情

---

### REQ-API-REORDER-001: 設備排序 API

#### Scenario: 成功重新排序

- **Given** 請求包含有效的 device_ids 陣列
- **When** PATCH `/api/v1/datalink/devices/reorder`
- **Then** 應返回 200 OK
- **And** 設備順序應按照 device_ids 陣列更新

---

## ADDED UI Consistency Requirements

### REQ-UI-001: 全系統設計 Token 統一

所有頁面（包含 TestPage）應使用統一的設計 Token。

#### Scenario: 背景色一致性

- **Given** 使用者開啟 SmartDashboard 頁面
- **When** 切換到 TestPage
- **Then** 背景色應相同（#0F172A）
- **And** 卡片樣式應一致
- **And** 按鈕樣式應一致

#### Scenario: 圖標一致性

- **Given** 所有頁面
- **When** 顯示圖標
- **Then** 應使用 Lucide Icons
- **And** 圖標大小應統一（20x20 或 24x24）
- **And** 不應使用 Emoji 作為圖標
