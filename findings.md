# Findings

## 2026-04-09 Studio 主線重定義

- 使用者已明確重新定義產品主目標：**從 PLC 取資料並存入 Database**。
- 現行 `Output = Local Modbus + Database` 的 unified step 模型不再符合產品心智。
- `Local Modbus` 與 `MQTT` 應視為 **Share / Publish** 能力，而不是主線必經 step。
- 新 IA 方向已確認為：將分享能力抽成 `/studio/share` 獨立子路由。
- 後續新補充推翻了「Database 是唯一主線終點」這個假設：
  - 有些情境只是要收集資料、整理順序，然後直接 share 給其他終端
  - 因此 `Database` 不是唯一終點，`Share / Publish` 也可能是直接目的地
  - 主線更合理的抽象應該是：`Device -> Source -> Tag -> Destination`
- `Destination` 已確定為新的主線最後一步命名；它比 `Database` 更寬，也比 `Publish` 更中性。
- `Destination` 的最佳 IA 目前已確定為先進 `Destination Hub` 再分流，而不是直接跳進某個子頁。
- `Destination Hub` 的價值在於：
  - 明確告訴操作員目前是在做「資料去向選擇」
  - 讓 `Database` 與 `Share / Publish` 形成平行目的地，而不是互為前後依賴
  - 非常適合後續做 Stitch / Pencil / skill 模板化
- 主線完成條件已進一步確定：
  - 只要在 `Destination` 內完成任一合法目的地，就算完成主線
  - 不需要同時完成 `Database` 與 `Share / Publish`
  - 這代表 `Destination` step 的任務是「讓資料有有效去向」，而不是「逼使用者配置所有去向」
- `Destination Hub` 的最佳互動不是完全中立，也不是自動導流，而是：
  - 依目前 tag / data shape / operator intent 顯示「建議目的地」
  - 同時保留所有目的地的明確入口
  - 不進行自動跳轉，避免把 Destination hub 變回隱性分支
- `Share Hub` 目前最合理的邊界是「只做能力入口」，不承擔：
  - 全域分享 KPI dashboard
  - 跨協定統一設定表
  - 大型監控總覽
- 這能避免 share 區再次長成一張萬用 output dashboard，讓 `Local Modbus` / `MQTT` 維持清楚的獨立工作區。
- 使用者已明確補充：雖然 `Database`、`Local Modbus`、`MQTT` 可以共用同一組上游資料上下文，但 **每個目的地都必須提供自己的批次規劃便利操作**，不能退化成逐筆綁定 UI。
- `Database` 的代表性便利操作：
  - 將多個 point/tag 規劃為同一筆資料列的多個欄位
  - 例如電力儀表的 `point1 -> tag_A1 ... point10 -> tag_KW`，需要快速寫入同一筆 row
- `Local Modbus` 的代表性便利操作：
  - 將一組 point/tag 快速規劃成同一區段的連續暫存器
  - 避免逐筆指定 register
- 因此新 IA 除了 hub 與 route，還必須把「destination-specific bulk planner」視為一級需求。
- `Delivery Group` 已確認為新的共用抽象：
  - 它代表一組應被一起規劃、一起送往目的地的資料
  - 對 `Database` 來說是一筆 row 的欄位集合
  - 對 `Local Modbus` 來說是一段連續 register block
  - 對 `MQTT` 來說是一個 topic / payload message
- `Delivery Group` 的價值是把 grouping logic 從各目的地抽出，避免三套 planner 各自重做同一件事。
- `Delivery Group` 的狀態模型已進一步確認：
  - 建立方式：`system-suggested, operator-adjusted`
  - 作用域：全域共用主物件
  - 各目的地只保存自己的 projection / override，而不是各自維護一套獨立 group
- `Delivery Group` 已進一步被確認為 destination UX 的主列表單位：
  - 進到 `Database`、`Local Modbus`、`MQTT` 時，主視角應該先看 group
  - 不是先看單一 tag
  - 這能避免 destination planner 退化回逐筆綁定 UI
- `Database Workspace` 的最佳骨架目前已確定為三區式：
  - 左側 `Delivery Groups`
  - 中央 `Row Planner`
  - 右側 `Schema / Apply Preview / Validation`
- 這個骨架的核心價值是把「群組選擇」、「row 規劃」、「schema/apply 驗證」拆開，避免所有操作擠進一張大 grid。
- `Database Row Planner` 的主操作模型已確定為 `Template-assisted row composition`：
  - 不是空白起稿
  - 也不是單純一格一格對 column
  - 而是由系統先提出 row / column 組成建議，再由操作員快速修正
- 第一版必備便利功能已鎖定：
  - `Auto column suggestion`
  - `Rename all with pattern`
  - `Column order drag sort`
  - `Drop / keep toggle`
  - `Data type hint / coercion warning`
  - `Row preview`
- `Local Modbus Workspace` 也已確定採三區式骨架，但中區改為 `Register Block Planner`。
- `Local Modbus` 的主操作模型已確定為 `Block-first register composition`：
  - 先把一組資料規劃成一段 block
  - 再由系統協助排成連續 registers
  - 不再以單一 tag / 單一 register 為主操作心智
- 第一版必備便利功能已鎖定：
  - `Sequential auto-pack`
  - `Gap-aware auto-pack`
  - `Width-aware placement`
  - `Base register quick set`
  - `Contiguous block preview`
  - `Conflict heatmap`
  - `Dry run before apply`
- `MQTT Workspace` 也已確定採三區式骨架，但中區改為 `Payload / Topic Planner`。
- `MQTT` 的主操作模型已確定為 `Template-assisted message composition`：
  - 先決定 topic / message family
  - 系統先提出 payload 結構建議
  - 操作者再快速修正 key、結構與 metadata
- 第一版必備便利功能已鎖定：
  - `Topic pattern generator`
  - `Payload schema suggestion`
  - `Flat vs nested payload switch`
  - `Rename keys with pattern`
  - `Metadata include toggle`
  - `Sample message preview`
  - `Dry run / publish validation`

## 為什麼不維持舊模型

- 舊模型把「主資料落地」與「附屬資料分享」混成同一個完成條件。
- 它也會讓使用者誤以為一定要先進 Database 才能走分享路徑。
- 對後續 AI 產版來說，舊模型會持續輸出錯誤資訊架構，因為 prompt 會把 share 當成 mainline。

## 目前最佳 IA 判斷

- 不建議把 `Local Modbus` / `MQTT` 再綁回舊的 unified Output step。
- 不建議把分享能力完全切出 `/studio`，否則會斷掉 tag、資料整理與目的地選擇的上下文。
- 目前最平衡的方向是：
  - `/studio` 主線完成於一個更抽象的 `Destination / Deliver` 階段
  - `/studio/share` 承接 `Local Modbus` / `MQTT`
  - `Database` 與 `Share / Publish` 應成為平行目的地，而不是先後依賴關係

## 與現行規格的衝突

- 現行 OpenSpec 仍要求在 `/studio` 的同一個 Output step 中統一處理 `Local Modbus` 與 `Database`。
- 若採新模型，後續需要進行 spec amendment，而不只是換一版 UI layout。
- 新的 OpenSpec change 不應只被理解成「讓 AI 產出的新前端接回現有 backend」：
  - 它首先是產品 IA / domain flow 的正式重定義
  - 其次才是定義這個新前端如何對接既有 backend contract
- 對 AI 產版流程本身，也需要一套明確的迭代規則：
  - 結構 / IA / workflow 問題應優先回改 prompt
  - screen direction 正確但局部不順時，才保留產出做 refinement
  - 只有當 screen map 與 page ownership 穩定後，才應進入前端實作與 backend 對接
- 如果這三層混做，最常見的結果是：
  - prompt 還在改，卻已經開始接 API
  - layout 尚未定稿，就先 patch 實作細節
  - 最後 UI 與 backend 契約一起反覆重工
- `/test` 也必須被納入這次 redesign package，但角色是：
  - 工程 debug / 現場工具頁
  - 不併入 `/studio` 主產品主線
  - 仍需要獨立的 UI/UX 盤點、API inventory 與 integration guidance
- `/test` 的產品定位已確認為 `Field Engineer Debug Console`
- `/test` 與 `/studio` 的視覺關係也已確認：
  - 共用同一套品牌 / 設計語言
  - 但 `/test` 的操作語氣要更工具化、更高密度、更直接
  - 不需要沿用 `/studio` 的產品引導節奏
- 進一步分析目前 `/test` 頁面後，可確認它的真正問題不是功能不足，而是：
  - 主角是很多 card，不是 current session
  - command、monitor、scan、debug 混在同一層
  - `polling`、`monitor`、`RTU polling` 三種持續讀值工具心智重疊
- `/test` 下一版最合理的 IA 目前已確定為：
  - 單一路由 `/test`
  - `Session Header`
  - `Setup Rail`
  - `Main Workspace`
  - `Diagnostics Rail`
- `Main Workspace` 建議切為四種工作區：
  - `Command Lab`
  - `Live Monitor`
  - `Scan`
  - `RTU Polling Jobs`
- `Diagnostics Rail` 應升級為固定右側診斷柱，而不只是可收折的 debug card。
- skill 結構目前最合理的方式已確認為雙層：
  - 通用版 `frontend-redesign-integration-review`
  - 專案版 `studio-frontend-redesign-review`
- 通用版只保留可跨 repo 重用的方法與檢查清單。
- 專案版才承載本 repo 的 `/studio`、`/test`、datalink domain 與現有 integration 規則。

## 接下來需要定義的關鍵問題

- `/studio/share` 從哪裡進入最自然
- `Local Modbus` 與 `MQTT` 在 share 內是 tab、兩頁，還是先 overview 再 drill-down
- 主線 readiness 與 share readiness 要如何切開顯示

## 2026-04-09 Point 角色重定義

- `Point` 已確認不再是主產品主角，而是 **internal acquisition unit**。
- 在新主線裡：
  - `Source` 負責規劃怎麼從設備讀資料
  - `Point` 是該規劃落地後的 machine-facing 採集單位
  - `Tag` 是 operator-facing 的資料語意單位
  - `Destination` 與 `Delivery Group` 應以 `Tag` 為主，不應回頭以 `Point` 作主要操作對象
- 這代表 UI 原則必須改成：
  - `Point` 主要出現在診斷、追蹤、例外修正
  - 不應回到 manual point-first 的主流程
- `Point` 在 UI 的保留位置已初步確定為 4 類：
  - `Source step`
  - `Inspector / Diagnostics`
  - `Traceability view`
  - `Exception tools`
- `Point` 不應再以獨立主工作台或 destination 主列表的形式出現。

## 2026-04-09 Tag 角色重定義

- `Tag` 已確認為新主線中的 **primary semantic data object**。
- 在新主線裡：
  - `Source` 產生 point / acquisition result
  - `Tag` 將 point 整理成可理解、可交付的語意資料
  - `Destination` 再以 tag 作為主要輸入，組成 `Delivery Group`
- 因此 `Tag` 將成為：
  - operator-facing 的主要資料物件
  - `Delivery Group` 的穩定輸入單位
  - 命名、型別、資料語意、業務意圖的主要承載層
- `Tag step` 的主操作模型也已初步確認為 `Semantic Refinement Board`：
  - 不是只做 candidate review
  - 而是將來源資料整理成可交付的語意資料物件
  - 其結果將直接供 `Destination` 與 `Delivery Group` 使用
- `Tag Workspace` 的最佳骨架目前已確定為三區式：
  - 左側 `Candidate / Tag Groups`
  - 中央 `Semantic Editor`
  - 右側 `Traceability / Preview / Readiness`
- 第一版必備便利功能已鎖定：
  - `Auto naming suggestion`
  - `Rename with pattern`
  - `Data type confirm / convert hint`
  - `Unit / meaning annotation`
  - `Keep / drop toggle`
  - `Merge / split candidate`
  - `Key tag / critical tag marker`
  - `Destination-readiness preview`
- `Tag Groups` 與 `Delivery Group` 的關係已確定採：
  - `Tag Group -> Delivery Group suggestions`
  - `Tag Group` 側重語意分組
  - `Delivery Group` 側重目的地投影與交付方式
  - 前者可作為後者的建議來源，但不強制一對一綁定
