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

### 第一輪用

```text
I am redesigning an existing industrial data gateway frontend.
Please help me generate one screen direction at a time.
Do not redesign the whole product in one response.
Start with the overall product direction first.
```

### 後續 screen iteration 用

```text
Keep the accepted overall product direction.
Now focus only on this screen or this one refinement request.
Do not redesign the whole product again.
```

## 使用規則

- 下方各段 `Prompt` 都是 prompt body
- 不建議單獨裸貼
- 建議每次都先貼前置說明模板，再貼 prompt body

## Prompt 0: Product Foundation Body

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

## Prompt 2: Device Screen Body

```text
Design the Device step for /studio. This screen is where operators define and validate data source devices.

Key capabilities:
- create, edit, clone devices
- separate connect and probe diagnostics
- show protocol traits, address model, and capability context
- expose device readiness for the next step

The screen should feel operational and trustworthy, not like a generic form page.
```

## Prompt 3: Source Screen Body

```text
Design the Source step for /studio. This step is used to plan how data is read from industrial devices.

Key capabilities:
- source rule planning
- address lattice / register planning canvas
- live values and link view
- source coverage awareness

The screen should make spatial address planning understandable, not feel like a spreadsheet.
```

## Prompt 4: Tag Screen Body

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

## Prompt 5: Destination Hub Body

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

## Prompt 6: Database Workspace Body

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

## Prompt 7: Share Hub Body

```text
Design the Share / Publish hub under /studio.

This is only a capability-entry page, not a dashboard.

It should present two clear destination cards:
- Local Modbus
- MQTT

Use the same visual language as /studio, but emphasize that these are delivery targets for already refined data.
```

## Prompt 8: Local Modbus Workspace Body

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

## Prompt 9: MQTT Workspace Body

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

## Prompt 10: Test Console Body

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

### Refine Studio consistency body

```text
Keep the same overall layout, but make the visual language across Device, Source, Tag, and Destination feel more like one product family.
Do not change the information architecture.
```

### Refine Test density body

```text
Keep the same /test structure, but increase the information density and make it feel more like an engineering console and less like a product workflow page.
Do not make it look messy.
```

### Refine Delivery Group framing body

```text
Keep the existing destination layouts, but make Delivery Groups feel like the primary planning unit across Database, Local Modbus, and MQTT.
Do not revert to tag-first visual hierarchy.
```
