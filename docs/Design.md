# Legacy `/studio` Workbench UI 設計原型（歷史 pre-delete evidence）

> **文件狀態（2026-08-24）**：本文件只保留 dedicated `/studio` 移除前的設計證據，不是 `/studio/v2` 的現行設計。`WorkbenchFrame`、`WorkbenchInspectorPanel`、`TagBindingStudio` 與 `frontend/src/pages/datalink/workbench/**` 均已隨 legacy surface 移除；目前 contract 為 `/studio/v2` setup、`/studio/runtime` observer、保留 `/test` 與 `/gateway/*`，而 `/studio` 依 generic unknown-route policy 處理。詳見 [退場紀錄](./releases/retire-legacy-studio-and-polish-v2.md)。

> 更新重點（歷史）：**主工作區（`WorkbenchFrame` / `main`）**與**右側檢查面板（`WorkbenchInspectorPanel`）**的現代化層次與可掃讀性；與 Step 3 Tag 候選列卡片語彙對齊。
> 設計取向依 **B2B／工業儀表板**：信任感、低炫光、資訊分塊（Miller’s Law／chunking），**刻意避免**大面積毛玻璃與「SaaS 紫／網格英雄區」等泛用 AI 版型。

---

## 1. 設計約束與受眾

| 項目 | 決策 |
|------|------|
| 受眾 | 單人工程／維運操作，長時間盯屏 |
| 品牌 | 延續既有 **slate 深色 + cyan 重點色**（閘道產品語意） |
| 技術 | Tailwind CSS、無新增設計套件 |
| 無障礙 | 維持語意化 `aside`／`main`、對比以現有色階為基礎微調 |

---

## 2. 視覺語彙（tokens 概念）

### 2.0 統一外殼 `WB_SHELL_SURFACE`（第二輪收斂）

**歷史檔案（已刪除）**：`frontend/src/pages/datalink/workbench/workbenchShellTokens.ts`

五區共用同一底層，避免各區各自漸層／重陰影／多層 ring 疊加造成的視覺疲勞與「數位炫光」感（呼應 frontend-design： restraint、solid hierarchy）。

```text
rounded-2xl border border-slate-800/70 bg-slate-900/95 shadow-sm shadow-black/30
```

**套用元件（歷史 legacy，已刪除）**：`WorkbenchContextBar`、`WorkbenchStepRail`、`WorkbenchFrame` 的 `main`、`WorkbenchInspectorPanel` 的 `aside`、`WorkbenchBottomSummaryBar`。
**外框網格（歷史 legacy，已刪除）**：`WorkbenchFrame` 根節點 `gap-3`、`p-3`，讓區塊之間有穩定呼吸（8pt 網格）。

### 2.1 主工作區 `main`（`workbench-primary-work-area`）

- **殼層**：`WB_SHELL_SURFACE` + **`p-6`**。
- **意圖**：與頂／左／右／底**同色階實心面**；角色差異改由**內容與內距**表達，而非另一套漸層。

### 2.2 檢查面板 `aside`（`workbench-inspector-panel`）

- **殼層**：`WB_SHELL_SURFACE` + **`p-5`**。
- **標頭**：`border-b border-slate-800/55 pb-4` 分隔標題與內容；內層 `flex-1 min-h-0` 維持捲動與空狀態置中行為。

### 2.3 檢查面板內卡片 `INSPECTOR_CARD_CLASS`

共用於裝置／來源規則／Tag／輸出追蹤等區塊：

```text
rounded-xl border border-slate-800/65 bg-slate-950/70 p-4 shadow-sm shadow-black/25
```

- **意圖**：比外殼**略深一階**即可形成 chunking；**不再**加 `ring-white`，避免與外殼邊線雙重描邊發糊。

### 2.4 步驟導覽（Step Rail，水平於 Context Bar）

- **位置**：`WorkbenchStepRail` 已**嵌入** `WorkbenchContextBar` 左側，不再佔用獨立左欄；主版面網格為 **`main | inspector`** 兩欄（中列）。
- **版面**：「步驟」標籤 + **cyan 豎條**（`sm+`）+ **橫向**四顆步驟按鈕（編號、截斷標題、readiness 點）；`lg` 時與裝置資訊區之間有 **直向分隔線**。
- **使用中步驟**：`border-cyan-500/40` + `bg-slate-800/90`；編號塊 `ring-cyan-400/25`。
- **測試**：`context-bar-step-summary` 使用 **`data-active-step`** 標示目前步驟鍵（`device` | `source` | `tag` | `output`）。

### 2.5 區塊 eyebrow `INSPECTOR_SECTION_EYEBROW_CLASS`

```text
text-[10px] font-semibold uppercase tracking-[0.16em] text-slate-500
```

- **意圖**：統一「區塊標籤」階層；英文仍可用 uppercase，中文環境下字重與字距仍提供節奏感。

---

## 3. Tag 檢查面板（Step 3）專項

與左側 **Tag 候選列**（歷史 `TagBindingStudio`，已刪除）對齊：

| 元素 | 設計 |
|------|------|
| 主標題（Tag key） | `text-base font-semibold tracking-tight text-slate-50` |
| 綁定狀態 chip | `rounded-md`、邊界與內高光與候選列 `getStatusToneClass` 家族一致 |
| 明細 `dl` | 外層 `rounded-lg` + `divide-y` + 弱 `bg-slate-900/25` + `ring`，**原始／轉換值**使用 `font-mono tabular-nums` |
| 位址 | `font-mono` 利於對齊與複製 |
| 警示（衝突／已連結） | `rounded-xl`、略提高邊界與 ring，**不**用過強純色塊 |

---

## 4. 與 frontend-design 技能對齊摘要

- **Hick**：側欄維持單一焦點區塊，細項收斂在卡片內。
- **Fitts**：主區 padding 與按鈕既有最小高度延續；檢查面板空狀態置中於剩餘高度。
- **Miller**：卡片 + `dl` 分區降低工作記憶負載。
- **Von Restorff**：cyan／emerald／rose 僅用於狀態與危險操作，中性面積佔多數（約 60-30-10）。
- **反模式迴避**：無 mesh gradient、無紫色主軸、無整頁強毛玻璃。

---

## 5. 歷史實作對照（程式位置，檔案已刪除）

| 歷史區域 | 歷史檔案（已刪除） |
|------|------|
| 外殼 token | `frontend/src/pages/datalink/workbench/workbenchShellTokens.ts` |
| 網格與主工作區 | `frontend/src/pages/datalink/workbench/WorkbenchFrame.tsx` |
| 頂欄／步驟軌／底欄 | `WorkbenchContextBar.tsx`、`WorkbenchStepRail.tsx`、`WorkbenchBottomSummaryBar.tsx` |
| 檢查面板殼層、卡片常數、Tag 明細 | `frontend/src/pages/datalink/workbench/WorkbenchInspectorPanel.tsx` |
| Tag 候選列卡片（左欄） | `frontend/src/pages/datalink/workbench/TagBindingStudio.tsx` |

---

## 6. 變更紀錄

| 日期 | 摘要 |
|------|------|
| 2026-03-22 | 初版：`main`／`aside` 層次與陰影、檢查面板卡片與 Tag 明細 `dl` 現代化；本文件記錄原型決策。 |
| 2026-03-22 | 第二輪：**`WB_SHELL_SURFACE` 統一五區外殼**；主區／側欄改實心底 + 輕陰影；Step Rail 使用中態去 glow；`INSPECTOR_CARD` 去 ring；Context 主按鈕 `min-h-9`。 |
