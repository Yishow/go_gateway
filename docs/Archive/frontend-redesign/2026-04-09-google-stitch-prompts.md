# 2026-04-09 Google Stitch Prompt Package

本文件依 Google Stitch prompt guide 拆成：

- 先用 plain language 高層次定方向
- 一次一個 screen
- 一次一種修改
- 先產出 layout，再逐步 refinement

## 使用方式

1. 先貼「前置說明模板」
2. 再貼 `Prompt 0`
3. 後續每一輪改用 follow-up 模板，再貼對應 screen prompt
4. 每個 screen 只做少量迭代
5. 不要把所有頁面需求一次塞進同一個 prompt

## 前置說明模板

共用模板見 [2026-04-09-stitch-usage-templates.md](/Users/yishow/prj/go_gateway/docs/frontend-redesign/2026-04-09-stitch-usage-templates.md)

本文件使用方式：

- 第一輪先用 `第一輪用`
- 後續一般微調用 `後續收斂用`
- 要固定同一張頁面骨架時，用 `後續固定頁面用`
- 下方各段都是 `prompt body`

## 第一輪建議順序

`Prompt 0 -> 1 -> 2 -> 3 -> 4 -> 5`

先確認 `/studio` 主線與 shell，不要一開始就深入 destination 細頁。

## 後續建議順序

主線方向穩定後，再跑：`Prompt 6 -> 7 -> 8 -> 9 -> refinement prompts`

## 每輪你要接受的是什麼

- `Prompt 0`：只接受產品家族方向與整體語氣。
- `Prompt 1`：只接受 `/studio` shell 骨架。
- `Prompt 2-5`：只接受各步驟主工作面的 layout 與主次層級。
- `Prompt 6-9`：只接受 destination 細頁的操作模型。
- `refinement prompts`：只修正單一問題，不重開新方向。

建議你每一輪只判斷三件事：

1. 主操作物件有沒有變對
2. 操作順序有沒有順
3. 是否還像這個產品，不像 generic dashboard

## 實際貼法建議

### 第一輪貼法

先貼共用模板檔裡的 `第一輪用`，再貼 `Prompt 0` 或 `Prompt 1`。

### 第二輪貼法

如果你已經接受產品方向，但只要調整某一頁，就改貼：

1. `後續收斂用`
2. 單一 screen 的 prompt body
3. 最後補 1 到 3 句限制句

### 固定骨架貼法

如果 Stitch 每次都把 `/studio` 變成不同架構，就改貼：

1. `後續固定頁面用`
2. 單一 screen 的 prompt body
3. `Preserve` 裡列出不能變的區塊

建議 `/studio` 常用 Preserve：

- top context bar
- left 4-step rail
- center working canvas
- right inspector / readiness panel
- bottom summary strip

## Prompt 0: Product Foundation Body

用途：先定整體產品氣質。

```text
Design a B2B industrial data gateway web application with two major areas:
1. /studio = the main product workflow
2. /test = a field engineer debug console

The visual language should feel like one product family, but the two areas should have different operating tones.

/studio is workflow-driven, structured, and product-oriented.
/test is more tool-like, denser, and optimized for fast diagnostics.

Use a dark industrial control-room style with strong information hierarchy, monospaced operational data where appropriate, compact status chips, and clear distinction between workflow steps and engineering diagnostics.
```

## Prompt 1: Studio Shell Body

用途：建立 `/studio` 的整體 shell 骨架。

```text
Create the main /studio shell for an industrial data platform. The primary workflow is:
Device -> Source -> Tag -> Destination

This shell should be desktop-first and include:
- a top context bar
- a 4-step navigation rail
- a main working surface
- an inspector / diagnostics panel
- a bottom summary area

The shell should support clear workflow progression and not look like a generic dashboard.
```

建議接受條件：

- 看得出這是 workflow shell，不是 dashboard
- 四步驟層級清楚
- 主工作面與 inspector 有明顯分工
- 後續可以容納 Device / Source / Tag / Destination 四頁

## Prompt 2: Device Screen Body

用途：只做 `Device` screen。

```text
Design the Device step for /studio. This screen is where operators define and validate data source devices.

Key capabilities:
- create, edit, clone devices
- separate connect and probe diagnostics
- show protocol traits, address model, and capability context
- expose device readiness for the next step

The screen should feel operational and trustworthy, not like a generic form page.
```

建議接受條件：

- connect 與 probe 是分開的
- device 定義不是長表單到底
- readiness 明確對接到下一步
- 協議與能力資訊不是藏在次層

## Prompt 3: Source Screen Body

用途：只做 `Source` screen。

```text
Design the Source step for /studio. This step is used to plan how data is read from industrial devices.

Key capabilities:
- source rule planning
- address lattice / register planning canvas
- live values and link view
- source coverage awareness

The screen should make spatial address planning understandable, not feel like a spreadsheet.
```

建議接受條件：

- 主角是 Source Rule 與格狀規劃，不是 list page
- register / address coverage 很容易掃描
- live values 與 planning 有關聯
- 不會被做成一般 table CRUD

## Prompt 4: Tag Screen Body

用途：只做 `Tag` screen。

```text
Design the Tag step as a Semantic Refinement Board.

This is not just a review list. It is where raw acquisition results are refined into semantic data objects.

The layout should have:
- left: candidate or tag groups
- center: semantic editor
- right: traceability, preview, and readiness

Essential actions:
- auto naming suggestions
- rename with pattern
- data type confirmation hints
- unit and meaning annotation
- keep or drop
- merge or split
- mark key or critical tags
- preview destination readiness
```

建議接受條件：

- review-first 成立，不是把手動 mapping 當主流程
- semantic editor 是主工作面
- traceability 與 preview 能幫你快速判斷是否保留
- 有整理與語意補全的感覺，不是單純 rename list

## Prompt 5: Destination Hub Body

用途：先做 `Destination Hub`。

```text
Design the Destination step as a Destination Hub.

This is the final workflow step where users choose where the refined data should go.

Important rules:
- the step label remains Destination
- users should first land on a hub, not jump directly into one destination
- destination options are:
  - Database
  - Share / Publish
- the hub can show a recommendation based on current data shape, but must never auto-redirect

The screen should make it obvious that Database and Share are parallel destination options.
```

建議接受條件：

- 先進 hub，再進能力頁
- Database 與 Share / Publish 並列
- 不會自動導向某一種 output
- 能看出這是 workflow 最後一步，不是另一套產品首頁

## Prompt 6: Database Workspace Body

用途：做 `Database` workspace。

```text
Design the Database workspace for /studio.

The primary object is not a single tag, but a Delivery Group.

The layout should have:
- left: Delivery Groups
- center: Row Planner
- right: Schema / Apply Preview / Validation

The main interaction model is template-assisted row composition.

Key actions:
- auto column suggestion
- rename all with a pattern
- drag-sort columns
- drop or keep fields
- show data type or coercion warnings
- show final row preview

This should support cases like grouping 10 power-meter tags into one database row.
```

建議接受條件：

- 主角是 Delivery Group，不是單一 tag
- row composition 很直觀
- schema preview 與 validation 足夠靠近
- 支援一筆資料列對應多個量測點

## Prompt 7: Share Hub Body

用途：做 `Share / Publish` 的 capability hub。

```text
Design the Share / Publish hub under /studio.

This is only a capability-entry page, not a dashboard.

It should present two clear destination cards:
- Local Modbus
- MQTT

Use the same visual language as /studio, but emphasize that these are delivery targets for already refined data.
```

建議接受條件：

- 這只是 capability hub，不是 metrics dashboard
- Local Modbus 與 MQTT 是平行入口
- 強調分享 / 發佈，不與 DB 儲存混在一起

## Prompt 8: Local Modbus Workspace Body

用途：做 `Local Modbus` workspace。

```text
Design the Local Modbus workspace.

The primary object is a Delivery Group.

Layout:
- left: Delivery Groups
- center: Register Block Planner
- right: Register Map Preview / Conflict / Apply

The main interaction model is block-first register composition, not per-tag register entry.

Key actions:
- sequential auto-pack
- gap-aware auto-pack
- width-aware placement
- base register quick set
- contiguous block preview
- conflict heatmap
- dry run before apply
```

建議接受條件：

- block-first 成立，不是 per-tag register form
- 連續區塊與衝突風險易讀
- base register 與 auto-pack 很順手
- 能看出它是 share destination，不是 source planning

## Prompt 9: MQTT Workspace Body

用途：做 `MQTT` workspace。

```text
Design the MQTT workspace.

The primary object is a Delivery Group.

Layout:
- left: Delivery Groups
- center: Payload / Topic Planner
- right: Topic Preview / Payload Preview / Validation / Apply

The main interaction model is template-assisted message composition.

Key actions:
- topic pattern generator
- payload schema suggestion
- flat vs nested payload switch
- rename keys with pattern
- metadata include toggle
- sample message preview
- dry run or publish validation
```

建議接受條件：

- topic 與 payload composition 是主角
- Delivery Group 還是主要操作單位
- preview 區能快速看最終 message
- 不會退化成一般 broker 設定表單

## Prompt 10: Test Console Body

用途：快速產出 `/test` 的高層方向；更細的 `/test` 請看另一份文件。

```text
Design /test as a Field Engineer Debug Console for factory and on-site debugging.

This is not part of the /studio workflow.

Key capabilities:
- protocol selection
- connection config
- connect, read, write, batch operations
- monitor and packet/log debugging
- profile switching
- compact tool cards

Use the same product family visual language as /studio, but make this screen denser, more tool-like, and more immediate.
```

更細的 `/test` page-level prompts 另見 [2026-04-09-test-console-stitch-prompts.md](/Users/yishow/prj/go_gateway/docs/frontend-redesign/2026-04-09-test-console-stitch-prompts.md)。

## Refinement prompts

## 常用 follow-up 限制句

這幾句可以直接貼在任一 body 後面：

```text
Do not redesign the whole /studio shell again.
Keep the current step order.
Keep Destination as the final workflow step.
Do not turn this into a generic analytics dashboard.
Preserve the accepted information hierarchy.
```

### Refine Studio consistency body

用途：主線結構已接受，但產品家族感還不夠一致時使用。

```text
Keep the same overall layout, but make the visual language across Device, Source, Tag, and Destination feel more like one product family.
Do not change the information architecture.
```

### Refine Test density body

用途：`/test` 方向正確，但密度還不夠像工程 console 時使用。

```text
Keep the same /test structure, but increase the information density and make it feel more like an engineering console and less like a product workflow page.
Do not make it look messy.
```

### Refine Delivery Group framing body

用途：destination pages 已成形，但 `Delivery Group` 還沒成為主操作單位時使用。

```text
Keep the existing destination layouts, but make Delivery Groups feel like the primary planning unit across Database, Local Modbus, and MQTT.
Do not revert to tag-first visual hierarchy.
```

## 建議你在 Stitch 的實際操作法

1. 先跑 `Prompt 0`，只選方向，不急著看細節。
2. 接著跑 `Prompt 1`，確認 shell。
3. 再逐頁跑 `Prompt 2` 到 `Prompt 5`。
4. 如果 hub 與主線已接受，再跑 `Prompt 6` 到 `Prompt 9`。
5. 若某一頁只是不夠順，不要重跑第一輪，改用 refinement。

我的建議：

- 如果你對整個產品調性不滿意，回去改 `Prompt 0`。
- 如果你只是不喜歡某一頁布局，改該頁 prompt，不要重講全產品。
- 如果你只是不想讓版面一直變，先換成 `後續固定頁面用`。
