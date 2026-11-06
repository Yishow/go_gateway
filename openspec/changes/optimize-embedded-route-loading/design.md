## Context

基準 main 0f54740e2df083ff0e6695897bfa4963039d5fda，2026-11-06 再以遠端 main 核對。App eager 載入 DatalinkWorkbenchV2Page；該檔直接載入 advanced shell、autosave 與來源編輯。SetupStaticFiles 對 assets 設長快取，現有函式沒有 gzip negotiation。量測數字見 proposal，均非真 browser 時間。

相關既有決策 reliable-embedded-frontend-delivery / Decision: Use route-level lazy imports behind one deterministic Suspense boundary，理由是避免空白頁與競爭載入狀態；本案保留單一可辨識 loading/error 契約。其固定 logical vendor grouping 與完整同步要求仍保留，不能以目前 Vite 未配置就宣稱已取消。退休 /studio 由 legacy-studio-retirement 及 embedded spec 的 Final removal preserves non-target embedded delivery 優先，不重建退休路由。

## Goals / Non-Goals

**Goals:** 減少首頁未使用程式與傳輸；保留第一畫面、導航與可重現交付。
**Non-Goals:** 不改設定／保存／採集／SQL／logs lifecycle；不做全站重寫、service worker、CDN、HTTP API compression、font 換型或 deployment。實作僅限本案已授權範圍。

## Decisions

### 輕量首頁與按需編輯

先新增 AdvancedWorkbenchV2Page 模組，保留輕量 frame 與路由解析在既有入口；再把 advanced/editor 的既有 wiring 原樣遷移至 lazy boundary。最後移除首頁 eager import 對該模組的依賴。不得拆散 mutation owners、dirty navigation barrier 或因 loading remount 清除草稿。替代的整頁直接 lazy 成 generic spinner 拒絕，因會失去目前 reading frame。

### 靜態傳輸與版本快取

只為 embed 中具 content-hashed identity 的 JS/CSS 使用標準庫 gzip，啟動時各壓縮一次並存固定資產集合，避免每次 request 壓縮。生成內容不寫 repository；新增記憶體成本需與 bytes 一起量測。未產生更小 representation 時保留 identity。Accept-Encoding 遵守 q=0、wildcard 與 identity；gzip;q=0 不得送 gzip。Range request 在 identity 可接受時使用 identity；若 Range 與 identity;q=0 同時出現，本案選擇 406，不回被拒絕的 identity 或偷偷忽略 range。HEAD 與 GET 的選定 representation headers 相同且無 body；Vary: Accept-Encoding、MIME、Content-Length、ETag 必須對應 representation。沒有可接受 representation 回 406。

HTML 設 Cache-Control: no-cache，hashed assets 設 public,max-age=31536000,immutable；非 hashed 檔不得取得 immutable。304 正確繼承快取／Vary；API、SSE 不進此處壓縮。資產缺失 404，不能回 SPA HTML。只採 gzip，不增加 Brotli dependency。保留安全 error boundary，不自動無限重載舊 HTML。

### 分層量測與交付

傳輸 byte、未壓縮 JS byte、請求數、browser 首次可讀與可互動分開記錄。候選設計門檻是相同 build/fixture 的首頁 eager JS bytes 與 gzip transfer bytes 均小於基準，且 warm return 不重抓 unchanged immutable 資產；不承諾未測的固定百分比或毫秒。實作前先固定相同 toolchain/build/input，五次冷／暖跑記錄 median/p95 及原始值；五次是驗證樣本條件，不是統計可靠性保證。

### 最終 Mac 瀏覽器驗收

**Supersedes**: optimize-embedded-route-loading / 分層量測與交付

保留分層量測與交付的fixture、bytes、請求數與五次cold/warm原始值、median/p95條件，僅替換執行順序與browser地點。2026-11-06使用者明確要求四案先完成雲端非browser實作，最後一起在Mac驗收。先逐案verify/review並修復程式問題；尚未執行的browser與最終結案保持未完成，不以雲端DOM／HTTP結果替代，不提前archive。原「不使用Mac」限制由這次明確授權取代，但只涵蓋全部實作完成後的browser驗證；雲端仍不得擴權，真PLC／正式DB／部署未新增授權。

## Implementation Contract

- 保留 /studio/v2、/studio/runtime、/studio/logs、/test、/gateway/* 與現有 redirect/query/hash/Back/Forward；/studio 只走 generic unknown-route policy。
- Static asset graph 必須包含所有 dynamic imports 與 shared logical dependency groups；不因 source split 移除 runtime logs 或工程工具。
- 兩個獨立驗證面：app-routing-lazy-load.test.tsx 證明 frame/deep-link/chunk failure；internal/web/embed_delivery_test.go 覆蓋 GET/HEAD/Range/q=0/304/MIME/404/HTML cache，並檢查不改 API/SSE。
- frontend production build 後以現有同步路徑建立 embedded binary，tests/shell/embedded-frontend-delivery.sh 驗證 stale removal、placeholder、圖完整性；不提交生成 bundles。
- 真 browser 以 frontend/tests/e2e/embedded-frontend-delivery.spec.ts 驗證 cold/warm、deep-link、chunk failure、版本切換。雲端既有 Chromium socket EPERM 保留為歷史 blocked evidence，不是通過；依最終 Mac 瀏覽器驗收決策，四案非browser實作完成後才在Mac統一執行。
- In scope 是 proposal named paths 與相應測試。Out of scope 是 runtime、DB、廣泛依賴更新；add-windows-tray-runtime-logs 保留自己的路由及 host lifecycle 所有權，碰同檔時先比對其進度。
- 每一遷移批次可 build；完成後 affected frontend/full build、affected Go tests、diff/line gate。跨平台與 browser 未跑必須分列。

## Risks / Trade-offs

- [拆分造成額外等待或草稿 remount] → stable owner 與 dirty navigation regression，先 add 再 migrate。
- [gzip 佔 RAM] → 固定 embedded 集合、不接受任意 cache keys、量測額外 bytes，若超過原始資產大小的設計上限則只使用 identity 並回報不達改善門檻。
- [HTML/asset 版本錯配] → 可重驗 HTML、完整同步、不可變 hash、失敗明示；不以刪除測試解決。
