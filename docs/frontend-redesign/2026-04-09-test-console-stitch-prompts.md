# 2026-04-09 Test Console Stitch Prompts

## 使用原則

這組 prompts 專門給 `/test`。

建議順序：

1. 先產 shell
2. 再產 session header
3. 再產各 workspace
4. 最後修 diagnostics 與密度

不要一次把整個 `/test` 所有要求塞進同一個 prompt。

## 前置說明模板

共用模板見 [2026-04-09-stitch-usage-templates.md](/Users/yishow/prj/go_gateway/docs/frontend-redesign/2026-04-09-stitch-usage-templates.md)

本文件建議：

- 第一輪先用 `第一輪用`
- 後續一般微調用 `後續收斂用`
- 要固定 `/test` 的 shell 骨架時，用 `後續固定頁面用`
- `Preserve` 建議先填：top session header / left setup rail / main workspace / right diagnostics rail

## 第一輪建議順序

`Prompt T0 -> T1 -> T2 -> T4 -> T8 -> T11`

這一輪先驗證 `session-first`、`command vs diagnostics`、`工程 console 氣質`。

## 後續建議順序

第一輪方向穩定後，再依需求補：`T3 / T5 / T6 / T7 / T9 / T10 / T12 / T13`

## 每輪你要接受的是什麼

- `T0`：只接受 `/test` 的角色與操作語氣。
- `T1`：只接受整體 console shell。
- `T2`：只接受 session header 是否足夠像工程儀表列。
- `T4`：只接受 command lab 是否順手。
- `T8`：只接受 diagnostics rail 是否有權重。
- `T11-T13`：只修操作節奏，不重做整張頁。

建議你每輪只看三件事：

1. 現在是不是 session-first
2. setup / action / diagnostics 是否分工清楚
3. 整體是否更像現場 debug tool，而不是管理後台

## 實際貼法建議

### 第一輪貼法

先貼共用模板檔裡的 `第一輪用`，再貼 `Prompt T0` 或 `Prompt T1`。

### 第二輪貼法

如果 console 方向已經對，只是某個 workspace 不夠順，就改貼：

1. `後續收斂用`
2. 單一 `/test` prompt body
3. 補 1 到 3 句限制句

### 固定 shell 貼法

如果 Stitch 每次都把 `/test` shell 改掉，就改貼：

1. `後續固定頁面用`
2. 目標 prompt body
3. `Preserve` 裡列出固定區塊

建議 `/test` 常用 Preserve：

- top session header
- left setup rail
- main command workspace
- right diagnostics rail
- current route remains single-page console

## Prompt T0: Test Console Foundation Body

用途：先定 `/test` 的角色與操作語氣。

```text
Design /test as a field engineer debug console for an industrial data gateway.

This is not part of the /studio workflow.

The page should feel like a session-centered engineering console, not a set of generic dashboard cards.

Use the same product family as /studio, but make /test denser, more operational, and more immediate.

Prioritize:
- current connection session awareness
- fast command execution
- live diagnostics
- protocol-level debugging
```

## Prompt T1: Test Console Shell Body

用途：建立 `/test` 的整體 shell。

```text
Create the main shell for /test.

The shell should include:
- a top session header
- a left setup rail
- a main workspace area
- a right diagnostics rail

This should remain a single-route console, not a multi-step wizard.

The page should feel like a real engineering tool used on-site in a factory.
```

## Prompt T2: Session Header Body

用途：強化 current session 的主角地位。

```text
Design a top session header for /test.

This header should make the current test session obvious at a glance.

Include:
- current profile
- protocol family
- connection mode
- endpoint summary
- connected or disconnected state
- connection id
- monitor state
- latest packet or error activity
- quick connect and disconnect actions

This should feel more like an instrument status bar than a marketing hero area.
```

## Prompt T3: Setup Rail Body

用途：補左側 setup rail。

```text
Design the left setup rail for /test.

This is where engineers prepare and switch testing context.

Include:
- profile manager
- protocol selector
- connection mode selector
- connection configuration
- saved preset or reusable config area

The layout should support frequent switching between profiles and protocols without feeling like a long form page.
```

建議接受條件：

- profile / protocol / endpoint 設定切換很快
- 不像長表單
- 可重用 preset 很明顯
- 不會搶走主工作區焦點

## Prompt T4: Command Lab Body

用途：建立 `Command Lab`。

```text
Design the main Command Lab workspace for /test.

This workspace is for immediate protocol commands and quick verification.

Include:
- single read or write mode
- batch command mode
- command polling mode
- result panel
- recent command history or recent result list

Make this feel like a controlled command bench, not like a CRUD form.
```

建議接受條件：

- single read/write、batch、polling 有清楚分區
- result panel 與 recent history 靠近主操作
- 看起來像命令台，不像設定頁

## Prompt T5: Live Monitor Workspace Body

用途：補長時間觀察用的 `Live Monitor`。

```text
Design the Live Monitor workspace for /test.

This workspace is for longer-running observation after a connection is already established.

Include:
- monitor item configuration
- monitor interval
- chart area
- latest values list
- monitor event log

The screen should help engineers understand signal changes over time without looking like a business analytics dashboard.
```

建議接受條件：

- 監看設定與趨勢檢視彼此連動
- 看得出這是 debug monitor，不是 BI chart
- 最新值與事件紀錄足夠靠近

## Prompt T6: Scan Workspace Body

用途：補掃描與發現工具頁。

```text
Design the Scan workspace for /test.

This workspace is used to discover devices, stations, or targets in the field.

Include:
- scan configuration
- station range or IP range setup
- scan speed controls
- progress area
- scan statistics
- result table or result list
- search and filter tools

This should feel like a discovery utility, separate from direct command execution.
```

建議接受條件：

- scan setup、進度、結果關係清楚
- discovery utility 感很強
- 不會跟 command lab 混在一起

## Prompt T7: RTU Polling Jobs Workspace Body

用途：補進階 `RTU Polling Jobs` 工具頁。

```text
Design the RTU Polling Jobs workspace for /test.

This is a specialized advanced tool, not the default command mode.

Include:
- polling job list
- per-job operation setup
- enable or disable controls
- interval summary
- live job status

The layout should make it obvious that this is persistent job-oriented polling, not the same as quick command polling.
```

建議接受條件：

- job-oriented 與 quick polling 明確區隔
- 每個 job 的狀態很容易掃描
- 不會誤以為這是一般命令工具

## Prompt T8: Diagnostics Rail Body

用途：建立右側固定 diagnostics rail。

```text
Design the right-side diagnostics rail for /test.

This rail should stay visible while engineers work in the main workspace.

Include:
- packet stream
- system logs
- selected packet detail
- protocol display mode toggle such as hex, ascii, parsed
- quick clear and export actions

This should feel like a serious diagnostics surface, not a decorative sidebar.
```

## Prompt T9: Packet Detail Expansion Body

用途：細修右側 packet detail 展開互動。

```text
Keep the existing /test shell, but redesign the packet detail interaction in the diagnostics rail.

When a packet is selected, show:
- timestamp
- direction
- protocol
- payload length
- connection id
- copy hex action
- optional analysis panel

Make the expanded detail feel precise and engineer-friendly.
```

## Prompt T10: Density Refinement Body

用途：不改架構，只提高資訊密度。

```text
Keep the same /test information architecture, but increase information density and reduce empty decorative spacing.

Make the interface feel more like an industrial engineer console.

Do not change the screen structure.
Do not turn it into a cluttered wall of controls.
```

## Prompt T11: Session-Centered Refinement Body

用途：把 current session 再拉成頁面主角。

```text
Keep the same /test structure, but make the current connection session feel like the main character of the page.

The operator should always know:
- what profile is active
- what endpoint is being tested
- whether the session is connected
- whether live monitoring is running
- whether errors or packets are currently active

Do not change the route structure.
```

## Prompt T12: Diagnostics-First Refinement Body

用途：把 diagnostics rail 做得更有權重。

```text
Keep the same /test structure, but make the diagnostics rail feel more authoritative and more useful during active debugging.

Increase the perceived importance of:
- packet stream
- log stream
- selected packet detail
- display mode controls

Do not make it overpower the main workspace entirely.
```

## Prompt T13: Setup-to-Action Rhythm Body

用途：微調 setup 到 action 的節奏。

```text
Keep the same /test structure, but improve the rhythm from setup to action.

The page should feel like:
1. choose or confirm profile
2. confirm protocol and connection
3. run commands or monitoring
4. inspect diagnostics

Do not turn this into a wizard.
Keep it as a fast engineer console.
```

## 常用 follow-up 限制句

這幾句可以直接貼在任一 `/test` body 後面：

```text
Do not redesign the whole /test console again.
Keep the single-route console structure.
Preserve the top session header and right diagnostics rail.
Do not turn this into a wizard.
Do not make it look like a generic admin dashboard.
```

## 我建議先跑的最小組合

如果你要先快速出第一輪，我建議先跑：

1. `Prompt T0`
2. `Prompt T1`
3. `Prompt T2`
4. `Prompt T4`
5. `Prompt T8`
6. `Prompt T11`

這樣可以先驗證最重要的結構：

- session-first 是否成立
- command 與 diagnostics 是否分明
- `/test` 是否真的更像工程 debug console

## 後續收斂提醒

如果你已經有一版 shell，不想讓 Stitch 每輪都重算整張頁面，就先貼共用模板檔裡的 `後續固定頁面用`，再貼你要調整的 prompt body。

## 建議你在 Stitch 的實際操作法

1. 先跑 `T0`，只看角色是否對。
2. 再跑 `T1`，確認 shell 結構。
3. 接著跑 `T2`、`T4`、`T8`，確認 session、command、diagnostics 三大主角。
4. 若核心成立，再補 `T3`、`T5`、`T6`、`T7`。
5. 最後才用 `T9-T13` 做精修。

我的建議：

- 如果你覺得 `/test` 不像現場工具，回去改 `T0`。
- 如果只是 command 不順，直接改 `T4`，不要重做整頁。
- 如果 diagnostics 很弱，先用 `T8` 或 `T12`，不要把 setup 區一起重設計。
