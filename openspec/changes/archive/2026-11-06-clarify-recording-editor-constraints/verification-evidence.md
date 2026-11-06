# 記錄編輯器限制說明：驗收紀錄

日期：2026-11-06。範圍：任務 4.1；本輪只補驗收及交接，未修改產品程式或測試。

## 本輪結果

| 結果 | 指令／範圍 | 實際證據 |
| --- | --- | --- |
| PASS | `cd frontend && npm test -- --run` | 633 檔、7277 項全數通過，102.31 秒。此次沒有重現前次 `recording-shared-configuration-navigation.test.tsx` 失敗；未修改該測試或產品程式，也未判定歷史失敗根因。 |
| PASS | `cd frontend && npx tsc --noEmit` | 退出碼 0。 |
| PASS | `cd frontend && npx eslint <相關檔案>` | 47 個 source／測試檔，清單見 [eslint-files.txt](evidence/2026-11-06/eslint-files.txt)。 |
| PASS | `bash scripts/check_file_lines.sh` | 行數檢查通過；既有超長檔依原有規則未增行。 |
| PASS | `git diff --check` | 無空白錯誤。 |
| PASS | dev server 截圖與實際檢視 | 新建各點位、已儲存長表、舊式、瞬時值，1440／768／390 各一張，共 12 張；另補 3 張互動／區段截圖。 |
| NOT RUN | 完整 build、完整 E2E、真實後端／資料庫寫入、PLC、Windows、現場驗收 | 不在本輪 4.1 驗收範圍，不能以本次測試或模擬畫面推論通過。 |

全套原始輸出見 [vitest.log](evidence/2026-11-06/vitest.log)。本輪結果取代 tasks.md 前次覆核對 4.1 的未完成判斷；前次紀錄仍保留為歷史證據。

## 截圖環境與限制

- macOS 本機；以 hc 管理 Vite dev server，只綁定 loopback。`VITE_BUNDLED_DEV=false npm run dev -- --host 127.0.0.1 --port 3223 --strictPort`。
- 使用 `agent-browser` 0.41.2、獨立工作階段與 Playwright Chromium headless shell 1248；先執行 `npx playwright install chromium` 確認專用瀏覽器可用。
- 既有 `frontend/tests/e2e/fixtures/recording-constraints-host.html?scenario=<new|long|legacy|instant>` 渲染真實 `GroupEditor` 元件，但身分、型錄、readiness、lifecycle、delivery 全部為模擬資料；不是 production route 或真實資料庫驗收。
- 兩個 Axios client 均由既有 fixture adapter 攔截，拒絕所有寫入與未知讀取；本輪未點保存、套用、檢查結構、試寫或移除。
- [browser-results.json](evidence/2026-11-06/browser-results.json) 的 12 次截圖紀錄皆只有已處理的模擬 GET，page errors 為空，整頁 `scrollWidth` 等於 viewport 寬度。390px 點位表使用既有內部水平捲動，不代表每個表格欄位同時可見。
- 初次開啟時，Vite 預先打包模式曾將 fixture URL 導至主頁，並對未啟動的 loopback backend 產生讀取連線失敗；切換上述逐模組模式後才擷取本次證據。未修改持久設定或啟動後端。
- 本輪沒有新產品行為，未新增測試或宣稱新的「修改前先失敗」證據。

## 畫面

| 情境 | 1440 | 768 | 390 |
| --- | --- | --- | --- |
| 新建各點位（尚未選點） | [畫面](evidence/2026-11-06/new-1440.png) | [畫面](evidence/2026-11-06/new-768.png) | [畫面](evidence/2026-11-06/new-390.png) |
| 已儲存長表 | [畫面](evidence/2026-11-06/long-1440.png) | [畫面](evidence/2026-11-06/long-768.png) | [畫面](evidence/2026-11-06/long-390.png) |
| 舊式寬表 | [畫面](evidence/2026-11-06/legacy-1440.png) | [畫面](evidence/2026-11-06/legacy-768.png) | [畫面](evidence/2026-11-06/legacy-390.png) |
| 瞬時值 | [畫面](evidence/2026-11-06/instant-1440.png) | [畫面](evidence/2026-11-06/instant-768.png) | [畫面](evidence/2026-11-06/instant-390.png) |

額外實際操作：

- 瞬時值選「算術平均」後仍維持最新值，先顯示量測用途確認；截圖後按「維持最新值」取消，未保存。[確認畫面](evidence/2026-11-06/instant-390-mean-consent.png)
- 長表輸入搜尋 `temperature` 並勾「只看有問題的列」，再從平均值分組按「到來源表格」：搜尋變空、篩選取消，焦點落在 `group-member-calculation-pt-1` 的 `SELECT`。[聚焦畫面](evidence/2026-11-06/long-390-mean-focus.png)
- 長表「目的地與確認」顯示共用表鎖定原因、複製入口及停用的欄位。[目的地畫面](evidence/2026-11-06/long-390-destination.png)

## 規格情境與既有測試

以下是本輪閱讀的斷言對照；執行證據為上述目前版本的全套通過，不等於真實 API 寫入或現場驗收。測試路徑以 `frontend/tests/unit/workbench-v2/` 為基準。

| 規格 Scenario | 分類 | 測試及可觀察斷言 |
| --- | --- | --- |
| New per-point recording chooses its table | 有測試 | `recording-mixed-destination-lock.test.tsx`：保存前可選／輸入 schema、table，受管固定、他人擁有的表停用。 |
| Saved long-table recording explains the lock | 有測試 | 上述鎖定測試及 `recording-destination-lock-notice.test.tsx`：已儲存維持鎖定、長表／設定檔原因不同、複製呼叫既有動作。 |
| Grouped per-point methods | 有測試 | `recording-method-groups.test.tsx`：3 個 latest／15 秒及 1 個 mean／60 秒的分組文字，並捲動聚焦第一個匹配控制項。 |
| Legacy recording is summarized, not converted | 有測試 | `recording-legacy-method.test.tsx`：60 秒舊式摘要、建立相同點位與連接器的新記錄，不更新原記錄。 |
| Switching to mean | 有測試 | `recording-instantaneous-method.test.tsx`：latest／15 秒改為 mean／30 秒，確認前不變、更版保存 payload、重新檢查結構提示；另涵蓋取消、拒絕、越界與累積值無平均選項。 |
| Unsaved managed recording | 有測試 | `recording-hidden-and-pending-areas.test.tsx`：未保存受管記錄提示保存後才能檢查結構；非受管不顯示。 |
| Unreachable destination | 有測試 | 同上：不讀取時無 loading；失敗顯示連線錯誤，重試後清除錯誤並顯示資料表。 |
| Schema confirmation required | 有測試 | `recording-readiness-plain-language.test.tsx`：編輯器及 runtime 的五種代碼、未知代碼、診斷文字與搜尋 autofill。 |

8 個 Scenario 皆有測試；delta spec 無 `##### Example:` 區塊，分組表格的兩列數值已對應斷言。無情境排除或未覆蓋缺口。單一畫面的寬度與間距依 TDD test-scope 準則不新增單元斷言，使用上述實際截圖檢視。

## 交接

本輪新增本檔與 `evidence/2026-11-06/`，並更新 tasks.md 的 4.1 狀態及本輪結果；既有產品 WIP、其他 change 與 index 保留。未 commit、archive、push、merge 或部署。

CLI 以順序編號 `9` 對應 4.1；`task start` 回報 `baseline_created: false`。`task done --file` 拒絕把 spec 目錄內文件登記為 source／test；本輪沒有產品檔案變更，因此只更新 checkbox，不以整個 dirty tree 補登歸屬。

後續執行 `$spectra-verify clarify-recording-editor-constraints` 及 `$spectra-review clarify-recording-editor-constraints`，再決定 archive。本紀錄只完成 apply 驗收，不替代獨立 verify／review。
