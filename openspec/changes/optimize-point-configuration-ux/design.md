# Design: 智慧 Dashboard 模式 - 技術設計

## Context

採用 **Option C: 智慧 Dashboard 模式**，將 9 個頁面整合為類 IDE 的一站式工作台。本文件詳述各元件的技術設計與介面定義。

---

## Goals / Non-Goals

### Goals

- ✅ 實現三欄式佈局（設備樹 | 內容區 | 操作面板）
- ✅ 記憶體格子視覺化，支援批量選取
- ✅ 側邊滑出面板，無需頁面跳轉
- ✅ 響應式設計，支援平板與桌面
- ✅ TDD 開發，先測試後實作
- ✅ **全系統 UI 一致性**（包含 TestPage）
- ✅ **設備樹拖曳排序功能**

### Non-Goals

- ❌ 不支援手機版（工控軟體主要在桌面使用）
- ❌ Phase 1 不實作離散位址選取
- ❌ 不修改後端資料模型

---

## 設計系統定義

### 顏色 Token

```typescript
// frontend/src/styles/tokens.ts

export const colors = {
  // 主題色
  primary: {
    50: "#EFF6FF",
    100: "#DBEAFE",
    500: "#3B82F6",
    600: "#2563EB",
    700: "#1D4ED8",
  },

  // 背景色
  background: {
    primary: "#0F172A", // Slate-900
    secondary: "#1E293B", // Slate-800
    tertiary: "#334155", // Slate-700
  },

  // 文字色
  text: {
    primary: "#F1F5F9", // Slate-100
    secondary: "#94A3B8", // Slate-400
    muted: "#64748B", // Slate-500
  },

  // 邊框色
  border: {
    default: "#334155", // Slate-700
    hover: "#475569", // Slate-600
  },

  // 狀態色
  status: {
    success: "#22C55E", // Green-500
    warning: "#EAB308", // Yellow-500
    error: "#EF4444", // Red-500
    info: "#3B82F6", // Blue-500
  },

  // 格子狀態色
  cell: {
    used: "#22C55E", // Green-500
    available: "#334155", // Slate-700
    selected: "#3B82F6", // Blue-500
    pending: "#EAB308", // Yellow-500
  },
} as const;
```

### 間距 Token

```typescript
export const spacing = {
  xs: "4px",
  sm: "8px",
  md: "16px",
  lg: "24px",
  xl: "32px",
  "2xl": "48px",
} as const;
```

### 佈局 Token

```typescript
export const layout = {
  treeNav: {
    width: "240px",
    collapsedWidth: "64px",
  },
  quickActions: {
    width: "280px",
  },
  slidePanel: {
    width: "400px",
  },
  memoryGrid: {
    cellWidth: "64px",
    cellHeight: "48px",
    gap: "4px",
  },
} as const;
```

---

## Component Design

### 1. SmartDashboard (主頁面)

```typescript
// frontend/src/pages/datalink/SmartDashboard.tsx

interface SmartDashboardProps {
  /** 無屬性，從 Context 獲取狀態 */
}

interface SmartDashboardState {
  /** 當前選中的設備 ID */
  selectedDeviceId: string | null;
  /** 當前選中的點位 ID（用於側邊面板） */
  selectedPointId: string | null;
  /** 側邊面板是否開啟 */
  isPanelOpen: boolean;
  /** 面板類型 */
  panelType: "point-detail" | "batch-create" | "quick-mapping" | null;
  /** 樹狀導覽是否收合 */
  isTreeCollapsed: boolean;
}
```

**佈局結構**：

```
┌─────────────────────────────────────────────────────────────────┐
│ Header (h-14)                                                   │
├────────────┬────────────────────────────────┬───────────────────┤
│            │                                │                   │
│ TreeNav    │ ContentArea                    │ QuickActions      │
│ (w-60)     │ (flex-1)                       │ (w-70)            │
│            │                                │                   │
│            │ ┌────────────────────────────┐ │                   │
│            │ │ MemoryGrid                 │ │                   │
│            │ │                            │ │                   │
│            │ └────────────────────────────┘ │                   │
│            │                                │                   │
│            │ ┌────────────────────────────┐ │                   │
│            │ │ PointTable (optional)      │ │                   │
│            │ └────────────────────────────┘ │                   │
│            │                                │                   │
└────────────┴────────────────────────────────┴───────────────────┘
                                                        │
                                              SlidePanel (overlay)
```

---

### 2. DeviceTreeNav (設備樹狀導覽)

```typescript
// frontend/src/components/datalink/DeviceTreeNav.tsx

interface DeviceTreeNavProps {
  /** 設備列表 */
  devices: Device[];
  /** 當前選中的設備 ID */
  selectedDeviceId: string | null;
  /** 選中設備事件 */
  onSelectDevice: (deviceId: string) => void;
  /** 是否收合模式 */
  isCollapsed: boolean;
  /** 切換收合事件 */
  onToggleCollapse: () => void;
}

interface TreeNode {
  id: string;
  label: string;
  icon: React.ReactNode;
  type: "device" | "points" | "mappings";
  count?: number;
  children?: TreeNode[];
  isExpanded?: boolean;
}
```

**視覺規格**：

| 狀態   | 背景            | 邊框                       | 文字           |
| ------ | --------------- | -------------------------- | -------------- |
| 預設   | transparent     | none                       | text-slate-300 |
| Hover  | bg-slate-700/50 | none                       | text-slate-100 |
| 選中   | bg-slate-700    | border-l-2 border-blue-500 | text-white     |
| 收合   | w-16            | none                       | 只顯示圖標     |
| 拖曳中 | bg-blue-500/20  | ring-2 ring-blue-500       | text-white     |

**拖曳排序支援（@dnd-kit/sortable）**：

```typescript
// 擴展 DeviceTreeNavProps
interface DeviceTreeNavProps {
  /** 設備列表 */
  devices: Device[];
  /** 當前選中的設備 ID */
  selectedDeviceId: string | null;
  /** 選中設備事件 */
  onSelectDevice: (deviceId: string) => void;
  /** 是否收合模式 */
  isCollapsed: boolean;
  /** 切換收合事件 */
  onToggleCollapse: () => void;
  /** 拖曳排序完成事件 */
  onReorder: (deviceIds: string[]) => void;
}

// 拖曳項目包裝器
interface SortableDeviceItemProps {
  device: Device;
  isSelected: boolean;
  onSelect: () => void;
}
```

**拖曳互動規格**：

| 操作          | 行為                       |
| ------------- | -------------------------- |
| 長按/拖曳手柄 | 開始拖曳，顯示半透明預覽   |
| 拖曳中        | 顯示插入位置指示線         |
| 放下          | 更新順序，呼叫 `onReorder` |
| ESC 鍵        | 取消拖曳                   |
| 鍵盤 Space    | 開始/結束拖曳（無障礙）    |

**後端 API 支援**：

```
PATCH /api/v1/datalink/devices/reorder

Request:
{
  "device_ids": ["uuid-1", "uuid-2", "uuid-3"]  // 新順序
}

Response:
{
  "success": true
}
```

---

### 3. MemoryGrid (記憶體格子視覺化)

```typescript
// frontend/src/components/datalink/MemoryGrid.tsx

interface MemoryGridProps {
  /** 設備 ID */
  deviceId: string;
  /** 協議類型 */
  protocol: ProtocolType;
  /** 中心位址 */
  centerAddress: string;
  /** 顯示範圍（上下各多少格） */
  range?: number;
  /** 已存在的點位 */
  existingPoints: Point[];
  /** 選中狀態（多選） */
  selectedAddresses: string[];
  /** 選取事件 */
  onSelect: (addresses: string[]) => void;
  /** 點選詳情事件 */
  onCellClick: (address: string, point?: Point) => void;
}

interface GridCell {
  address: string;
  displayAddress: string;
  status: "used" | "available" | "selected" | "pending";
  point?: Point;
  row: number;
  col: number;
}
```

**互動規格**：

| 操作            | 行為                                 |
| --------------- | ------------------------------------ |
| 單擊空格子      | 選中該格子，開啟側邊面板（新增點位） |
| 單擊已使用格子  | 開啟側邊面板（編輯點位）             |
| Shift + 單擊    | 範圍選取                             |
| Ctrl/Cmd + 單擊 | 切換選取狀態                         |
| 拖曳            | 框選多個格子                         |
| 懸停            | 顯示 Tooltip（點位名稱、值、時間）   |

**Tooltip 內容**：

```typescript
interface CellTooltip {
  address: string;
  pointName?: string;
  lastValue?: string;
  lastReadAt?: string;
  dataType?: DataType;
}
```

---

### 4. QuickActions (快速操作面板)

```typescript
// frontend/src/components/datalink/QuickActions.tsx

interface QuickActionsProps {
  /** 當前設備 */
  device: Device | null;
  /** 選中的格子數量 */
  selectedCount: number;
  /** 批量建立事件 */
  onBatchCreate: () => void;
  /** 快速映射事件 */
  onQuickMapping: () => void;
  /** 測試連線事件 */
  onTestConnection: () => void;
}

interface LiveStatus {
  cpuUsage: number;
  memoryUsage: number;
  connectionStatus: "connected" | "disconnected" | "connecting";
  lastSync: Date;
}
```

---

### 5. SlidePanel (側邊滑出面板)

```typescript
// frontend/src/components/datalink/SlidePanel.tsx

interface SlidePanelProps {
  /** 是否開啟 */
  isOpen: boolean;
  /** 標題 */
  title: string;
  /** 關閉事件 */
  onClose: () => void;
  /** 內容 */
  children: React.ReactNode;
  /** 寬度 */
  width?: "sm" | "md" | "lg"; // 320px | 400px | 560px
}
```

**動畫規格**：

```css
/* 進入動畫 */
.slide-enter {
  transform: translateX(100%);
}
.slide-enter-active {
  transform: translateX(0);
  transition: transform 300ms ease-out;
}

/* 離開動畫 */
.slide-exit {
  transform: translateX(0);
}
.slide-exit-active {
  transform: translateX(100%);
  transition: transform 200ms ease-in;
}
```

---

### 6. BatchPointCreator (批量建立器)

```typescript
// frontend/src/components/datalink/BatchPointCreator.tsx

interface BatchPointCreatorProps {
  /** 設備 ID */
  deviceId: string;
  /** 協議類型 */
  protocol: ProtocolType;
  /** 預選的位址（從格子選取帶入） */
  preselectedAddresses?: string[];
  /** 建立完成事件 */
  onCreated: (points: Point[]) => void;
  /** 取消事件 */
  onCancel: () => void;
}

interface BatchPointConfig {
  startAddress: string;
  quantity: number;
  nameTemplate: string;
  dataType: DataType;
  pollingGroupId: string;
  enabled: boolean;
  description?: string;
}

interface NameTemplateVars {
  index: number; // 0, 1, 2, ...
  address: string; // D100, D101, ...
  num: number; // 100, 101, ...
}
```

**命名模板支援**：

| 模板             | 範例輸出                     |
| ---------------- | ---------------------------- |
| `Sensor_{index}` | Sensor_0, Sensor_1, Sensor_2 |
| `{address}_Temp` | D100_Temp, D101_Temp         |
| `Motor_{num}`    | Motor_100, Motor_101         |

---

### 7. AddressParser (位址解析器)

```typescript
// frontend/src/utils/addressParser.ts

interface AddressParser {
  /** 解析單一位址 */
  parse(address: string): ParsedAddress;

  /** 展開位址範圍 */
  expand(startAddress: string, count: number): string[];

  /** 驗證位址格式 */
  validate(address: string): ValidationResult;

  /** 計算位址偏移 */
  offset(address: string, delta: number): string;
}

interface ParsedAddress {
  protocol: ProtocolType;
  area: string; // 'D', 'M', 'HR', 'IR', etc.
  startNumber: number;
  raw: string;
}

interface ValidationResult {
  valid: boolean;
  error?: string;
}
```

**各協議位址格式**：

| 協議   | 格式             | 範例            | 正則表達式           |
| ------ | ---------------- | --------------- | -------------------- |
| Modbus | `{FC}{Address}`  | 40001, 30100    | `/^[0134]\d{4,5}$/`  |
| FATEK  | `{Area}{Number}` | D0100, R0, M100 | `/^[DRMXYWTS]\d+$/i` |
| MC3E   | `{Area}{Number}` | D100, W0, M0    | `/^[DWMXYB]\d+$/i`   |

---

## 狀態管理

### Context 結構

```typescript
// frontend/src/contexts/SmartDashboardContext.tsx

interface SmartDashboardContextValue {
  // 設備狀態
  selectedDeviceId: string | null;
  selectDevice: (id: string) => void;

  // 格子選取狀態
  selectedAddresses: string[];
  selectAddresses: (addresses: string[]) => void;
  clearSelection: () => void;

  // 側邊面板狀態
  panelState: PanelState;
  openPanel: (type: PanelType, data?: any) => void;
  closePanel: () => void;

  // 設備資料（React Query）
  devices: Device[];
  points: Point[];
  isLoading: boolean;
}
```

---

## API 擴充

### 批量建立 API

```
POST /api/v1/datalink/points/batch
```

**Request**:

```json
{
  "device_id": "uuid",
  "polling_group_id": "uuid",
  "data_type": "int16",
  "enabled": true,
  "points": [
    { "name": "Sensor_0", "address": "D100" },
    { "name": "Sensor_1", "address": "D101" }
  ]
}
```

**Response**:

```json
{
  "success": true,
  "data": {
    "created_count": 10,
    "points": [...]
  }
}
```

---

## Verification Plan

### 單元測試（TDD）

#### 1. AddressParser 測試

```typescript
// frontend/src/utils/__tests__/addressParser.test.ts

describe("AddressParser", () => {
  describe("Modbus", () => {
    it("should parse holding register address", () => {
      expect(parser.parse("40001")).toEqual({
        protocol: "modbus_tcp",
        area: "HR",
        startNumber: 1,
        raw: "40001",
      });
    });

    it("should expand address range", () => {
      expect(parser.expand("40001", 3)).toEqual(["40001", "40002", "40003"]);
    });
  });

  describe("FATEK", () => {
    it("should parse D register address", () => {
      expect(parser.parse("D0100")).toEqual({
        protocol: "fatek_fbs",
        area: "D",
        startNumber: 100,
        raw: "D0100",
      });
    });
  });

  describe("MC3E", () => {
    it("should parse D register address", () => {
      expect(parser.parse("D100")).toEqual({
        protocol: "mc_3e",
        area: "D",
        startNumber: 100,
        raw: "D100",
      });
    });
  });
});
```

#### 2. MemoryGrid 組件測試

```typescript
// frontend/src/components/datalink/__tests__/MemoryGrid.test.tsx

describe('MemoryGrid', () => {
  it('should render correct number of cells', () => {
    render(<MemoryGrid centerAddress="D100" range={50} ... />);
    expect(screen.getAllByTestId('grid-cell')).toHaveLength(100);
  });

  it('should highlight used cells', () => {
    const points = [{ address: 'D100', name: 'Test' }];
    render(<MemoryGrid existingPoints={points} ... />);
    expect(screen.getByTestId('cell-D100')).toHaveClass('cell-used');
  });

  it('should call onCellClick when clicked', () => {
    const handleClick = jest.fn();
    render(<MemoryGrid onCellClick={handleClick} ... />);
    fireEvent.click(screen.getByTestId('cell-D100'));
    expect(handleClick).toHaveBeenCalledWith('D100', undefined);
  });

  it('should support range selection with Shift+Click', () => {
    const handleSelect = jest.fn();
    render(<MemoryGrid onSelect={handleSelect} ... />);
    fireEvent.click(screen.getByTestId('cell-D100'));
    fireEvent.click(screen.getByTestId('cell-D105'), { shiftKey: true });
    expect(handleSelect).toHaveBeenCalledWith(['D100', 'D101', 'D102', 'D103', 'D104', 'D105']);
  });
});
```

#### 3. BatchPointCreator 測試

```typescript
// frontend/src/components/datalink/__tests__/BatchPointCreator.test.tsx

describe('BatchPointCreator', () => {
  it('should generate names from template', () => {
    render(<BatchPointCreator deviceId="..." ... />);

    fireEvent.change(screen.getByLabelText('起始位址'), { target: { value: 'D100' } });
    fireEvent.change(screen.getByLabelText('數量'), { target: { value: '3' } });
    fireEvent.change(screen.getByLabelText('命名模板'), { target: { value: 'Sensor_{index}' } });

    expect(screen.getByText('Sensor_0')).toBeInTheDocument();
    expect(screen.getByText('Sensor_1')).toBeInTheDocument();
    expect(screen.getByText('Sensor_2')).toBeInTheDocument();
  });

  it('should call API on submit', async () => {
    const mockOnCreated = jest.fn();
    render(<BatchPointCreator onCreated={mockOnCreated} ... />);

    // Fill form...
    fireEvent.click(screen.getByText('建立'));

    await waitFor(() => {
      expect(mockOnCreated).toHaveBeenCalled();
    });
  });
});
```

### 手動驗證

1. 啟動開發伺服器
2. 開啟 SmartDashboard 頁面
3. 選取設備 → 確認記憶體格子載入
4. 點選格子 → 確認側邊面板開啟
5. 批量選取 → 確認批量建立流程
6. 驗證響應式佈局（1024px, 1440px, 1920px）
