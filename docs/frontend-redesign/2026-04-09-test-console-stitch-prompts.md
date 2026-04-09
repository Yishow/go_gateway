# 2026-04-09 Test Console Stitch Prompts

## 使用原則

這組 prompts 專門給 `/test`。

建議順序：

1. 先產 shell
2. 再產 session header
3. 再產各 workspace
4. 最後修 diagnostics 與密度

不要一次把整個 `/test` 所有要求塞進同一個 prompt。

## Prompt T0: Test Console Foundation

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

## Prompt T1: Test Console Shell

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

## Prompt T2: Session Header

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

## Prompt T3: Setup Rail

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

## Prompt T4: Command Lab

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

## Prompt T5: Live Monitor Workspace

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

## Prompt T6: Scan Workspace

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

## Prompt T7: RTU Polling Jobs Workspace

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

## Prompt T8: Diagnostics Rail

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

## Prompt T9: Packet Detail Expansion

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

## Prompt T10: Density Refinement

```text
Keep the same /test information architecture, but increase information density and reduce empty decorative spacing.

Make the interface feel more like an industrial engineer console.

Do not change the screen structure.
Do not turn it into a cluttered wall of controls.
```

## Prompt T11: Session-Centered Refinement

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

## Prompt T12: Diagnostics-First Refinement

```text
Keep the same /test structure, but make the diagnostics rail feel more authoritative and more useful during active debugging.

Increase the perceived importance of:
- packet stream
- log stream
- selected packet detail
- display mode controls

Do not make it overpower the main workspace entirely.
```

## Prompt T13: Setup-to-Action Rhythm

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
