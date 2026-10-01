# 2026-04-09 AI Layout Iteration Playbook

## 目的

回答一個很實際的問題：

AI 產出的版面如果不喜歡，下一步應該改 `prompt`，還是沿用產出慢慢修？

建議不要混著做。

先判斷問題層級，再選對迭代手段，速度才會快，也比較不會把對接節奏打亂。

## 我建議的三層迭代法

### Layer 1: 先改 prompts

適用於方向錯了。

典型訊號：

- 資訊架構不對
- 主次層級不對
- workflow 節奏不對
- `/studio` 與 `/test` 氣質混掉
- `Tag`、`Destination`、`Delivery Group` 的主角位置不對
- AI 產出的畫面又回到 tag-first 或 table-first 舊心智

這一層不要急著修成品。

因為你現在不滿意的不是一個按鈕，而是整個 screen 的設計前提。

### Layer 2: 再修生成結果

適用於方向已對，但局部不順。

典型訊號：

- 區塊順序差一點
- 某個 panel 太大或太小
- label 不夠準
- 某一組操作太分散
- 視覺密度不夠穩
- `/test` 還不夠像 debug console

這一層可以沿用現有產出做 refinement。

也就是：

- 保留 screen map
- 保留主 layout
- 只改一個 screen 或一組局部問題

### Layer 3: 最後才做實作修補

適用於設計方向已接受，只剩落地問題。

典型訊號：

- layout 已可接受，但元件間距要修
- 某些互動要接現有 hooks
- 某些 panel 需要換成真資料
- loading / empty / error / dry-run 狀態要補
- 設計沒問題，但型別與契約還沒接上

這一層才進入前端實作。

不然很容易出現：

- prompt 還沒定，就先硬改 code
- layout 還沒穩，就先綁 API
- 最後 UI 和 backend contract 一起被反覆重做

## 實際判斷公式

你可以用這個簡單判斷：

- 如果你想說的是「整個頁面不對」：改 prompts
- 如果你想說的是「這頁差不多對，但幾個地方不順」：修生成結果
- 如果你想說的是「畫面我接受了，只差接資料與互動」：進實作

## 對這個專案最重要的判斷點

### 一定要回去改 prompt 的問題

- `Destination` 沒先進 hub
- `Database` 與 `Share / Publish` 沒被畫成平行去向
- `Local Modbus` 被綁成 Database 後處理
- `Point` 又回到主工作流中心
- `Delivery Group` 沒成為 destination planner 的主列表單位
- `/test` 被做成精簡版 `/studio`

### 可以直接修生成結果的問題

- `Tag` 右側 readiness panel 太弱
- `Database` 的 row preview 太小
- `Local Modbus` 的 conflict 區塊不夠醒目
- `MQTT` topic preview 與 payload preview 排序不理想
- `/test` 的 monitor、packets、logs 密度不夠

### 才該進入實作層的問題

- 畫面要接 `frontend/src/services/datalink.ts`
- 畫面要接 `frontend/src/services/api.ts`
- 新 screen 要接 TanStack Query hooks
- 需要加 adapter 做 view-model 轉換
- mutation、dry-run、retry、error handling 要真正可用

## 建議迭代節奏

### 第 1 輪

只修方向。

目標：

- 主線對不對
- screen map 對不對
- `/studio` 與 `/test` 的分流對不對

### 第 2 輪

只修各 screen 骨架。

目標：

- 左中右工作區是否成立
- 主操作物件是否正確
- 便利操作是否有被體現

### 第 3 輪

只修局部操作體驗。

目標：

- panel 順序
- 資訊密度
- microcopy
- 視覺強弱

### 第 4 輪

才開始接回 backend。

目標：

- route shell
- adapter
- hook / service reuse
- mutation wiring
- verify

## 我對你的實際建議

我建議你把 AI 產版分成兩段：

### 段 A：先把 prompts 調到你喜歡的方向

先不要急著接 code。

因為只要方向還會改，後面每接一次 API 都是在重工。

### 段 B：當 screen map 穩定後，再接回現有 backend

這時才做：

- route 對接
- adapter 設計
- services / hooks reuse
- 測試與驗證

## 最後判準

不是「畫面看起來不錯」就算完成。

真正完成要同時成立：

- 你喜歡 screen direction
- 資訊架構正確
- 新 UI 能接回現有 backend contract
- 主要互動已驗證可用

因此答案不是二選一。

正確順序是：

1. 大問題先改 prompts
2. 小問題再修生成結果
3. 方向穩了才進實作
