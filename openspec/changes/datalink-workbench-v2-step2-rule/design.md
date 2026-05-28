## Context

Step 2 是整個 v2 工作台中互動最密集、邏輯最複雜的一頁。原型 `step2-rule.jsx`（676 行）涵蓋多規則 tab、規則編輯器、點位記憶體網格、Modbus Share 配置、批次工具列、合併點位表、衝突偵測等多個互動領域。直接 TypeScript 化單一檔案會超過 repo 500 行硬上限，需要拆檔；同時 derive 邏輯與互動 state 必須抽離，否則 component 內 useEffect / useMemo 會難以維護。

點位網格是 Dark Industrial Telemetry 風格的視覺核心：玩家熟悉的 Modbus 「記憶體格子」隱喻，搭配 24px 格線底紋與顏色語意（規則色 / amber skipped / 紅 ring 衝突 / 藍 ring 多選）。這個元件決定使用者在 Step 2 的整體體感，必須移植到位且測試覆蓋。

## Goals / Non-Goals

**Goals:**

- 1:1 視覺與互動復刻原型 Step 2，含網格三種點擊修飾鍵語意。
- 把 derive 邏輯（`derivePoints` / `deriveAllPoints` / `computeShareLayout` / `detectAddressConflicts` / `fnFromAddr`）抽到純函式 module，便於單元測試。
- 在 selector 層處理跨規則合併與衝突計算，避免每次 render 重算 O(rules × count) 邏輯。
- 為後續 Step 3 提供 `state.points` 衍生來源；Step 3 將直接消費。

**Non-Goals:**

- 不接後端 `POST /source-rules`。
- 不處理規則排序拖拉與跨裝置複製。
- 不變更全域 `state.settings.modbus_share`（只讀取）。
- 不支援超過 64 個 register 的單條規則（與原型一致）。
- 不在 Step 2 渲染 Tag 與 mapping 細節（屬於 Step 3）。

## Decisions

### 拆檔策略：以「資料區塊 / 互動行為」拆 10 個元件 + 2 個 state module

原型 676 行拆為 10 個 component 檔 + 2 個 state module：

- `Step2Rule.tsx`（~220 行）：容器，state derivation 與 onContinue gate。
- `RuleTabRail.tsx`（~150 行）：上層 rule tab 列。
- `RuleEditor.tsx`（~200 行）：左 col 5 容器，組合 RangeSummary + ScaleSection + ShareSection。
- `RangeSummary.tsx`（~80 行）：範圍摘要 chip + Modbus function code chip。
- `ScaleSection.tsx`（~100 行）：線性轉換 details。
- `ShareSection.tsx`（~150 行）：Modbus Share details。
- `PointGrid.tsx`（~200 行）：右 col 7 網格本體，含三種點擊修飾鍵 handler。
- `PointGridToolbar.tsx`（~120 行）：批次工具列。
- `MergedPointTable.tsx`（~200 行）：下層合併點位表。
- `index.ts` barrel。

2 個 state module：
- `state/sourceRule.ts`（~150 行）：純函式 derive helpers。
- `state/selectors.ts`（~80 行）：基於 reducer state 的 selector hook（`useAllPoints` / `useShareLayout` / `useConflictAddrs`）。

**Alternatives considered**：(A) 把網格與工具列合併在同一檔（~320 行）— 超過 300 行警戒；(B) 不抽 selectors，每次 component 用 `useMemo` 重算 — 跨多個元件重複計算，且測試 derive 邏輯時需要 mount React。

### 衍生 state：純函式 + selector hook

`state/sourceRule.ts` 包含全部純函式：

```ts
export function derivePoints(rule, deviceId, skippedSet): Point[];
export function deriveAllPoints(rules, fallbackDeviceId): Point[];
export function computeShareLayout(rules, baseRegister): Record<string, ShareLayout | null>;
export function detectAddressConflicts(allPoints): Set<string>;
export function fnFromAddr(addr): 'coil' | 'discrete_input' | 'input_register' | 'holding_register';
export function dataTypeWidth(type): number;
export function formatAddr(n): string;
```

`state/selectors.ts` 提供 hooks：

```ts
export function useAllPoints(): Point[];
export function useShareLayout(): Record<string, ShareLayout | null>;
export function useConflictAddrs(): Set<string>;
export function useRulePoints(ruleId: string): Point[];
```

每個 hook 用 `useMemo` cache 結果，依賴 `state.rules` 與 `state.devices[0]?.id` 與 `state.settings.modbus_share.base_register`。

**Alternatives considered**：(A) 把 derive 結果存在 reducer state（如原型 `setState({ points })`）— 增加 state 反向依賴，容易讓 reducer 變胖且需要在 reducer 內呼叫 derive；(B) 用 Reselect — 引入新依賴，且 selectors 規模小不需要 reselect 級別優化。

### 點位網格互動：點擊 + Shift + Ctrl/⌘ 三模式

`PointGrid` 內部 state：

```ts
const [gridSelection, setGridSelection] = useState<Set<string>>(new Set());
const [lastClickedIdx, setLastClickedIdx] = useState<number | null>(null);
```

點擊單一 cell 的 onClick handler：

```ts
function onCellClick(e: MouseEvent, address: string, idx: number) {
  if (e.shiftKey && lastClickedIdx != null) {
    // 範圍切換 skip
    const [lo, hi] = [Math.min(lastClickedIdx, idx), Math.max(lastClickedIdx, idx)];
    const set = new Set(rule.skipped_addresses);
    const baseSkipped = set.has(address);
    for (let j = lo; j <= hi; j++) {
      const a = rulePoints[j].address;
      if (baseSkipped) set.delete(a);
      else set.add(a);
    }
    dispatch({ type: 'updateRuleSkipped', ruleId: rule.id, skipped: [...set] });
  } else if (e.ctrlKey || e.metaKey) {
    // 多選 set toggle
    const next = new Set(gridSelection);
    next.has(address) ? next.delete(address) : next.add(address);
    setGridSelection(next);
  } else {
    // 一般 toggle skip
    dispatch({ type: 'toggleRuleSkippedAddress', ruleId: rule.id, address });
  }
  setLastClickedIdx(idx);
}
```

切換 rule 時 `gridSelection` 與 `lastClickedIdx` 重置（用 useEffect 依賴 `selectedRuleId`）。

**Alternatives considered**：(A) 用 react-dnd 處理拖拉選取 — 過度工程，鍵盤修飾鍵足以滿足；(B) 把 `gridSelection` 也放進全域 state — 只有 PointGrid 用，沒必要全域。

### Modbus Share 位址布局：自動接續 + 手動覆寫

`computeShareLayout(rules, baseRegister)` 行為：

- 用 `cursor = baseRegister` 起點。
- 依 `rules` 順序遍歷，每條 rule：
  - 若 `share_enabled === false`，layout[rule.id] = null，跳過。
  - 否則 `stride = rule.share_stride ?? dataTypeWidth(rule.data_type)`，`enabledCount = rule.count - rule.skipped_addresses.length`。
  - 若 `rule.share_start_register != null`，`start = rule.share_start_register`，layout 標記 `auto: false`。
  - 否則 `start = cursor`，layout 標記 `auto: true`。
  - `end = start + enabledCount * stride`。
  - `cursor = max(cursor, start + enabledCount * stride)`。

cursor 用 `max` 不是直接賦值，是為了讓「手動指定的規則 share 起點 < cursor」這種情況不會把後續自動分配往回拉。

**Alternatives considered**：(A) 強制手動 / 自動二選一全域 — 失去原型混合配置的彈性；(B) 在 reducer 觸發即時計算並寫入 state — 與 derive 邏輯散落兩處。

### 衝突偵測：在 `state.points` 層去重

```ts
function detectAddressConflicts(allPoints: Point[]): Set<string> {
  const usage: Record<string, string[]> = {};
  for (const p of allPoints) {
    if (p.skipped) continue;
    (usage[p.address] ??= []).push(p.rule_id);
  }
  return new Set(Object.entries(usage).filter(([, arr]) => arr.length > 1).map(([k]) => k));
}
```

衝突只在「跨規則」之間判定；同一規則內無法產生重複 address（derive 順序遞增）。`allPoints` 由 `deriveAllPoints(state.rules)` 來。

### 線性轉換預設值繼承到 Step 3

`Rule.scale_multiplier`、`Rule.scale_offset` 在 Step 2 設定後，由 `derivePoints` 寫入每個 Point 的 `_rule_scale` 與 `_rule_offset`（內部 underscore 標示「來自規則的預設值」）。Step 3 的 mapping 初始化會優先讀 `point._rule_scale` 作為初值，但允許個別覆寫。

**Alternatives considered**：(A) 在 Step 3 init 時直接讀 `state.rules` 反查 — 增加耦合；(B) 把 scale 放在 mapping 層而非 rule 層 — 失去「同一規則的點位預設共享相同轉換」便利。

### 改變參數時清空 skipped_addresses

修改 `start_address`、`count`、`data_type` 任一時，自動清空 `skipped_addresses` 陣列。原因：這三個參數任一改變後，舊的 skipped 位址可能已不屬於新範圍，保留會造成「skipped 位址不在網格內」的死狀態。

`updateRule` reducer 收到 patch 含上述三欄之一時，會主動把 `skipped_addresses` 重設為 `[]`。

## Implementation Contract

#### Behavior

- 進入 `current === 2`：顯示一張預設 rule（`Holding Registers` / `40001` / `int16` / `count=8`），右側網格顯示 8 個 cell（40001–40008）。
- 點「+ 新增規則」：tab 尾端新增 `規則 N`，初始 `count=4`，`naming_prefix=BLOCK{N}_`，所屬 device 預設取當前 selected rule 的 device 或 `devices[0]`。
- inline 改名：rule card 文字即時更新。
- 切到不同 rule：右側網格與左側編輯器內容換成該 rule；`gridSelection` 與 `lastClickedIdx` 重置。
- 改 `start_address` / `count` / `data_type` 任一：`skipped_addresses` 清空，網格 cell 數量與寬度同步刷新；範圍摘要 chip 更新；Modbus function code chip 依新位址第一位數字重新推斷。
- 改 `scale_multiplier` / `scale_offset`：rule 層預設值更新，但已存在的 `state.mappings` 不受影響（mapping 由 Step 3 init）。
- 點網格 cell：toggle 該位址的 skip 狀態；amber 與規則色互換。
- Shift+click：對 `lastClickedIdx → i` 範圍批次切換 skip；基準是「當前 cell 是否已 skip」。
- Ctrl/⌘+click：把 cell address 加入或移出 `gridSelection`；藍 ring 標示。
- 點批次工具列「全部啟用」：`updateRuleSkipped(ruleId, [])`。
- 點「全部略過」：把該規則所有 address 加入 skipped。
- 點「反轉啟用」：把 `skipped_addresses` 與「未 skipped」位址對調。
- 點「略過選取」/「啟用選取」：把 `gridSelection` 內位址加入/移出 skipped，並清空 `gridSelection`。
- 啟用 share 時，網格 cell 右下角浮出 mono `→{shareAddr}`；衝突 cell 額外渲染紅 ring 與右上 ! 徽章。
- 合併點位表反映所有啟用 rule 的衍生 point，跨 rule 共用相同位址時兩列都標 `衝突` 紅字。
- 表尾「繼續到映射」按鈕在 `totalEnabled === 0 || conflictAddrs.size > 0` 時 disabled；點擊呼叫 `onContinue`。
- 刪除規則：tab 至少保留一條；刪除後若被刪是當前 selected，自動切到第一個剩下的。

#### Interface / Data Shape

- `Rule` 型別（state/types.ts 既有）：
  ```ts
  type Rule = {
    id: string;
    device_id: string | null;
    name: string;
    start_address: string;
    count: number;          // 1-64
    data_type: 'bool' | 'int16' | 'int32' | 'int64' | 'uint16' | 'uint32' | 'uint64' | 'float32' | 'float64' | 'string';
    naming_prefix: string;
    enabled: boolean;
    scale_multiplier: number;
    scale_offset: number;
    data_format: '' | 'ABCD' | 'BADC' | 'CDAB' | 'DCBA';
    skipped_addresses: string[];
    share_enabled: boolean;
    share_start_register: number | null;
    share_stride: number | null;
  };
  ```
- `Point` 型別：
  ```ts
  type Point = {
    id: string;
    device_id: string;
    rule_id: string;
    rule_name: string;
    name: string;
    address: string;
    data_type: string;
    function: 'coil' | 'discrete_input' | 'input_register' | 'holding_register';
    width: number;
    enabled: boolean;
    skipped: boolean;
    display?: string;
    unit?: string;
    tag_key_suggest?: string;
    _rule_scale: number;
    _rule_offset: number;
  };
  ```
- `ShareLayout` 型別：`{ start: number; stride: number; end: number; auto: boolean } | null`
- Reducer rule actions：
  - `{ type: 'addRule'; rule: Rule }`
  - `{ type: 'removeRule'; ruleId: string }`
  - `{ type: 'updateRule'; ruleId: string; patch: Partial<Rule> }`（patch 含 start_address/count/data_type 任一時自動 reset skipped_addresses）
  - `{ type: 'renameRule'; ruleId: string; name: string }`
  - `{ type: 'toggleRuleEnabled'; ruleId: string }`
  - `{ type: 'updateRuleSkipped'; ruleId: string; skipped: string[] }`
  - `{ type: 'toggleRuleSkippedAddress'; ruleId: string; address: string }`
  - `{ type: 'toggleRuleShareEnabled'; ruleId: string }`
  - `{ type: 'updateRuleShareStart'; ruleId: string; start: number | null }`
  - `{ type: 'updateRuleShareStride'; ruleId: string; stride: number | null }`
- Selectors：
  - `useAllPoints(): Point[]`
  - `useShareLayout(): Record<string, ShareLayout | null>`
  - `useConflictAddrs(): Set<string>`
  - `useRulePoints(ruleId: string): Point[]`

#### Failure Modes

- 切換 rule 後 `gridSelection` 重置但 user 已習慣「跨 rule 選取」：不支援；明確以 useEffect 重置避免錯亂。
- `start_address` 輸入非數字字串：用 `parseInt(start.replace(/[^0-9]/g, ''), 10) || 40001` 回退；UI 仍接受任意字串作為標籤但內部以數字運算。
- `count` 超過 64：reducer 用 `Math.min(64, Math.max(1, count))` 夾住；UI input 設 `min=1 max=64`。
- 衝突偵測：跨規則同址同時啟用時，合併表雙列標紅 `衝突` chip，繼續按鈕 disabled。
- Modbus Share 全域 disabled：規則層 `share_enabled` 仍可設 true，但網格 cell 不顯示 share 位址、details summary 顯示「全域未啟用」chip。

#### Acceptance Criteria

- `frontend/tests/workbench-v2/sourceRule.test.ts` 覆蓋：
  - `derivePoints(rule, deviceId, skippedSet)`：基本 8 個 point + 部分 skipped。
  - `deriveAllPoints(rules)`：兩條 rule 合併。
  - `computeShareLayout(rules, 40001)`：自動接續 / 部分手動 / disabled rule 的混合場景。
  - `detectAddressConflicts(allPoints)`：兩條規則同址、單條規則無衝突。
  - `fnFromAddr('40001') === 'holding_register'`；其他 prefix 對應。
  - `dataTypeWidth` 對 10 種 data type 的回傳值。
- `frontend/tests/workbench-v2/step2-rule.test.tsx` 覆蓋：
  - 預設 1 條 rule + 8 個網格 cell。
  - 新增第 2 條 rule：tab 與合併表反映 2 條。
  - 改 count 從 8 到 4：網格 cell 數變 4、`skipped_addresses` 清空。
  - 切到 mc_3e device 並改 start_address：function code chip 更新。
- `frontend/tests/workbench-v2/step2-grid.test.tsx` 覆蓋：
  - 單擊 toggle skip：cell tone 切換到 amber。
  - Shift+click 範圍：lastClickedIdx=0、shift+click idx=3 → 4 個 cell skip 切換。
  - Ctrl+click 多選：3 個 cell ring blue、清除選取後 ring 消失。
  - 批次按鈕「全部略過」：rule.skipped_addresses 含所有位址。
- `frontend/tests/workbench-v2/step2-share.test.tsx` 覆蓋：
  - share 全域 enabled + rule share_enabled：網格 cell 右下浮出 mono → 標籤。
  - 手動設定 share_start_register=50001：cell 標籤從 50001 起算。
  - 點「改回自動分配」：share_start_register=null、stride=null。
  - 兩條規則 share 開啟、第二條起點未指定：第二條 cell 標籤從第一條 end 起算。
- `cd frontend && npm run lint && npm run typecheck && npm run test -- --run workbench-v2 && npm run build` 全綠。
- 手動 smoke：完整跑「新增第二條規則、改 count、Shift+click skip 範圍、Ctrl+click 多選、批次略過選取、開啟 share、手動改 share 起點、看合併表衝突高亮」。

#### Scope Boundaries

**In scope:**

- Step 2 完整靜態 UI（10 個元件檔）。
- `state/sourceRule.ts` 7 個純函式 + `state/selectors.ts` 4 個 hook。
- `useWorkbenchV2State` 新增 10 個 rule actions。
- `state/types.ts` 補完 Rule / Point / ShareLayout 型別。
- shell 換掉 Step2RulePlaceholder；shell spec 縮窄 placeholder。
- i18n `workbench-v2` 加入 Step 2 microcopy。
- 上述項目的 unit test。

**Out of scope:**

- 後端 API、SSE 推送、即時點位值（commit 後才有）。
- 規則排序、拖拉、複製。
- 跨裝置複製規則。
- 規則匯入 / 匯出。
- E2E 測試。
- Step 3/4/Settings 行為。
- 全域 Modbus Share 設定（屬於 Settings change）。
- 失敗 fallback fail tone（rule 編輯本身不會失敗，無需 fail UI）。

## Risks / Trade-offs

- [拆 10 個元件增加 import 開銷與 cross-file 跳轉成本] → 用 `index.ts` barrel 集中 export；元件之間 props drilling 限制在 2 層以內。
- [衝突偵測是 O(allPoints) 但每次 render 重算] → `useConflictAddrs` 用 `useMemo` cache，依賴 `state.rules` 與 `state.devices[0]?.id`；rules 變更頻率低。
- [Shift+click 範圍切換的「基準」語意可能讓使用者困惑（基準是 last clicked cell 的 skipped 狀態而非 current cell）] → 與原型一致；在 i18n hint 明寫「Shift+click：以上一格為基準切換範圍」。
- [share_start_register 留空 vs `null` 語意] → reducer 把空字串轉 null；input 顯示 placeholder `(自動)`；不允許負數。
- [`updateRule` 內隱含 reset skipped 行為可能讓「精準微調 data_type」場景丟失 skip 設定] → 在 i18n hint 明寫「變更位址 / 數量 / 型態會清空略過清單」。

## Migration Plan

1. 在 `state/types.ts` 補完 Rule / Point / ShareLayout 型別。
2. 建立 `state/sourceRule.ts`，TDD：先寫 `sourceRule.test.ts` 對 7 個純函式各 2–3 個 case。
3. 建立 `state/selectors.ts`，補對應 hook 測試。
4. 擴充 reducer 加入 10 個 rule actions，TDD：reducer test 覆蓋 `updateRule` 內隱含 reset 行為與 `toggleRuleSkippedAddress` 邊界。
5. 從子元件實作起：`RangeSummary` → `ScaleSection` → `ShareSection` → `RuleEditor`。
6. `PointGridToolbar` → `PointGrid`（含三種點擊修飾鍵）。
7. `RuleTabRail` → `MergedPointTable`。
8. 組裝 `Step2Rule.tsx`，覆寫 shell `current === 2` 分支。
9. 補 i18n。
10. 跑 lint / typecheck / test / build / check-lines。
11. 手動 smoke 全互動。

**Rollback strategy**：不刪除 shell change 的 placeholder；rollback 把 shell 的 `current === 2` 改回 placeholder、移除 step2 目錄、reducer actions 變 dead code。

## Open Questions

- 是否要支援「跨規則同步 skip 同位址」操作（例如把規則 A 與 B 都有的 `40001` 一鍵全部 skip）？傾向 **不**，使用者應該調整其中一條規則的範圍而非 skip；衝突反而是 UX 提示「請調整」訊號。
- 「裝置色 vs 規則色」在合併表上同時顯示兩個顏色點，是否會讓 UI 過於擁擠？傾向 **保留**，多裝置情境下「設備來源」與「規則來源」是兩個維度資訊；在單裝置情境下裝置欄會 hidden（只看到一個顏色點）。
