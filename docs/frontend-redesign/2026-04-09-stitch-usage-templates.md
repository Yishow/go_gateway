# 2026-04-09 Stitch Usage Templates

## 用途

這份文件提供 Google Stitch 實際貼 prompt 時的共用模板。

兩份 prompt 文件都應搭配這份一起用：

- `2026-04-09-google-stitch-prompts.md`
- `2026-04-09-test-console-stitch-prompts.md`

## 第一輪用

用途：先建立產品或頁面的整體方向，不要求沿用既有生成畫面。

```text
I am redesigning an existing industrial data gateway frontend.
Please help me generate one screen direction at a time.
Do not redesign the whole product in one response.
Start with the overall direction first.
```

## 後續收斂用

用途：保持既有方向，只修改同一張頁面的局部內容。

```text
Keep the accepted direction.
Now focus only on this screen or this one refinement request.
Do not redesign the whole page again.
```

## 後續固定頁面用

用途：固定同一張頁面骨架，避免每輪都變成不同頁。

```text
Keep the accepted direction and preserve the current screen structure.
Do not redesign the whole page again.
Only change the requested area.

Preserve:
- [fill the regions that must stay stable]
- [fill the interactions that must stay stable]
```

## 使用規則

- 每次先貼模板，再貼對應 prompt body
- 第一輪用 `第一輪用`
- 後續一般微調用 `後續收斂用`
- 要固定同一張頁面骨架時，用 `後續固定頁面用`

## 建議組裝方式

建議實際貼到 Stitch 時，用這種順序：

1. 先貼共用模板
2. 再貼單一 prompt body
3. 若你已經有一版畫面，再補一小段限制句

限制句可直接補在最後，例如：

```text
Do not add extra pages.
Do not turn this into a dashboard.
Keep the current hierarchy.
```

## 可直接貼用的組裝範例

### 範例 A：第一輪開新頁

用途：第一次產 `/studio` 或 `/test` 的新方向。

```text
I am redesigning an existing industrial data gateway frontend.
Please help me generate one screen direction at a time.
Do not redesign the whole product in one response.
Start with the overall direction first.

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

### 範例 B：第二輪局部微調

用途：方向接受，但只要調同一張頁面的某個區塊。

```text
Keep the accepted direction.
Now focus only on this screen or this one refinement request.
Do not redesign the whole page again.

Design the Device step for /studio. This screen is where operators define and validate data source devices.

Key capabilities:
- create, edit, clone devices
- separate connect and probe diagnostics
- show protocol traits, address model, and capability context
- expose device readiness for the next step

The screen should feel operational and trustworthy, not like a generic form page.
```

### 範例 C：固定頁面骨架後再調整

用途：你已經有一版可接受的 layout，不想每輪都被重做。

```text
Keep the accepted direction and preserve the current screen structure.
Do not redesign the whole page again.
Only change the requested area.

Preserve:
- top context bar
- left step rail
- main work surface layout
- right diagnostics panel

Keep the same /test structure, but make the diagnostics rail feel more authoritative and more useful during active debugging.

Increase the perceived importance of:
- packet stream
- log stream
- selected packet detail
- display mode controls

Do not make it overpower the main workspace entirely.
```

## 何時該改 prompt，何時該在成品上微調

- 如果整體方向不對：回到第一輪 prompt，改高層描述。
- 如果結構大致對，但畫面細節不喜歡：用 `後續收斂用` 微調。
- 如果同一張頁面每輪都跑版：改用 `後續固定頁面用`。
- 如果只是密度、視覺語氣、單一模組位置不順：不要重下完整第一輪 prompt。

## 常見失敗與修正

- 失敗：一次貼太多頁面需求。
  修正：一次只做一張頁面或一個 refinement。
- 失敗：第二輪又把整個產品重講一次。
  修正：第二輪只保留同一張頁面的限制與局部目標。
- 失敗：Stitch 每次都長出完全不同頁。
  修正：加入 `preserve the current screen structure`，並列出要保留的區塊。
- 失敗：畫面像 generic SaaS dashboard。
  修正：補強 domain language，例如 control-room、engineering console、register planning、delivery group。
