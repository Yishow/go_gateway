## Context

V2 的 Step 1~4 autosave change 會把資料變成後端 truth，但使用者另外定了重要邏輯：設備第一次啟動後，之後合法 autosave 應直接套用到正在執行的設備，而不是再次等待 Step 4。這個 change 專門定義「保存成功」到「runtime 套用成功」之間的契約。

## Goals / Non-Goals

**Goals:**

- 已啟動設備在合法 autosave 成功後直接套用新設定
- autosave 回應能回報 runtime apply 結果
- 未啟動設備 autosave 仍只保存，不偷啟動

**Non-Goals:**

- 不處理 invalid edit 導致的停用，那是另一個 change
- 不改 Step 4 第一次啟動的選取規則
- 不改 runtime 工作區切換 UI

## Decisions

### Running device update application

只要設備已經進入 `running` 狀態，任何會影響該設備 runtime 行為的合法 autosave 成功後，都應觸發重新套用，而不是再等使用者按額外按鈕。

### Autosave response exposes runtime apply state

前端需要知道「這次 save 只是保存」還是「已保存且已套用到執行中設備」。因此 autosave response 需帶出 apply 結果。

### Not-running devices stay dormant

尚未第一次啟動的設備，即使 autosave 成功，也只能更新 persisted config，不得偷偷變成 running。

## Implementation Contract

- Behavior:
  - 已啟動設備在合法 autosave 成功後，runtime 直接改用新設定
  - autosave response 可明確區分 `not_running`, `applied`, `apply_failed`
  - 未啟動設備 autosave 成功後仍保持未啟動
- Interface / data shape:
  - Step 1~4 autosave success payload 新增 `runtime_apply_status`
  - `runtime_apply_status` 值至少包含：
    - `not_running`
    - `applied`
    - `apply_failed`
  - `apply_failed` 時需附 `message`
- Failure modes:
  - 若 runtime apply 失敗，保存成功與套用失敗要分開表達
  - apply 失敗不得回報成單純 save success
- Acceptance criteria:
  - backend tests 覆蓋 running device autosave 直接套用與 not-running device 不自動啟動
  - frontend tests 覆蓋 `runtime_apply_status` 在 UI state 中可被識別
- Scope boundaries:
  - In scope: running device 的 valid autosave re-apply
  - Out of scope: invalid edit stop、first activation flow、workspace runtime list UI
- Not complete if:
  - 已在跑設備 autosave 成功後還要手動再按 Step 4
  - 前端無法知道這次 save 是否真的已套用到 runtime
  - 未啟動設備 autosave 會被偷偷啟動

## Risks / Trade-offs

- [Risk] save success 與 apply success 混在一起會讓故障難排查 → Mitigation: 用獨立 `runtime_apply_status` 表達，不把 apply 失敗藏掉
- [Risk] 某些欄位變更只需熱更新，某些欄位可能需要完整 reload → Mitigation: contract 只規定「已套用」，實作可選擇 hot reload 或 controlled restart，但回應語意要一致
