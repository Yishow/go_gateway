## Context

前面的 changes 會讓 Step 1~4 逐步建立 workspace 真實資料，Step 4 也只負責第一次啟動。接下來 `/studio/runtime` 不能再只假設自己接一個單台 `device_id` 就夠了，必須承接 workspace 內所有已成功存進後端的設備，並把可用與不可用設備一起展示。

## Goals / Non-Goals

**Goals:**

- `/studio/runtime` 載入 singleton workspace 內的設備集合
- 預設選擇 V2 順序最前面的可用設備
- 頁內可切換其它設備
- 沒有可用設備時保留 runtime 空狀態
- unavailable device 留在清單中但標示為不可用

**Non-Goals:**

- 不做 fleet-wide dashboard
- 不做 multi-workspace 切換
- 不改 `/studio/v2` 的 autosave 或 activation 契約

## Decisions

### Runtime page is workspace-scoped first

runtime page 先以 workspace 取設備集合，再決定目前選到哪台設備。`device_id` query 可以保留當 override，但不是唯一 source of truth。

### Default device follows v2 order

預設設備使用 V2 workspace 排列順序中的第一個可用設備，而不是依 runtime 事件先後或隨機第一筆資料。

### Unavailable devices remain visible

即使設備目前不可用，也要留在 runtime 清單中並標示 unavailable，這樣使用者才知道是哪台設備出了問題，而不是誤以為設備消失。

### Empty runtime state stays on page

若一台可用設備都沒有，runtime 頁保留空狀態與返回 V2 的修正指引，不強制 redirect。

## Implementation Contract

- Behavior:
  - `/studio/runtime` 先載入 singleton workspace device set
  - 預設選擇 V2 順序最前面的可用設備
  - 使用者可在頁內切換其它設備
  - unavailable device 保留在清單中但不可視為可觀察主設備
  - 若沒有可用設備，runtime 頁顯示空狀態並引導回 V2
- Interface / data shape:
  - `GET /api/v1/datalink/studio-v2/workspace/runtime-context`
  - payload 至少包含：
    - `workspace_id`
    - `devices[]` ordered by V2 order
    - each device: `device_id`, `name`, `running`, `availability_status`, `availability_reason`
    - `default_device_id`
- Failure modes:
  - 若 query `device_id` 指向 unavailable device，頁面仍可顯示該設備不可用狀態，但預設主觀察區不得假裝它可用
  - 若 workspace 內沒有任何 device，回傳空集合與空狀態 metadata，不 redirect
- Acceptance criteria:
  - frontend runtime tests 覆蓋 default selection、device switching、empty state、unavailable row visibility
  - backend tests 覆蓋 runtime context ordering 與 unavailable device payload
- Scope boundaries:
  - In scope: workspace runtime context、default selection、switching、empty state
  - Out of scope: fleet dashboard、多 workspace、autosave internals
- Not complete if:
  - runtime 還只吃單台 `device_id` 才能工作
  - unavailable device 直接從清單消失
  - 沒有可用設備時被硬轉回 `/studio/v2`

## Risks / Trade-offs

- [Risk] workspace context 與既有 single-device runtime hooks 並存一段時間 → Mitigation: 保留 `device_id` query 作為 override，但以 workspace context 為主
- [Risk] unavailable device 若還能被當主觀察區，會混淆使用者 → Mitigation: default selection 永遠選第一台可用設備，unavailable 只留在切換清單與狀態區
