## Context

recordingHomeSources.ts 的 readRecordingHomeSources 連做兩輪 workspace/devices/rules/mappings/points 讀取，依 fingerprint 比對，再由 projectSources 檢查 identity。RecordingHomeRoute 每次 visit 增加 epoch、gcTime=0。這是有意的保真，不是無條件可刪的重複請求。

simple-recording-home / Home guidance projects canonical scoped evidence 與 simplify-recording-home-guidance 的唯讀原則繼續有效。查到的 guide-step4-recording-workflow / Owner-local 操作觀察與唯讀導覽分界，理由是 owner-local observed 不等於全域 readiness；本案也不把摘要提升成 command proof。

## Goals / Non-Goals

**Goals:** 首頁只傳需要的來源摘要；返回保留 observed 畫面；來源變動／錯誤不能誤報 fresh。
**Non-Goals:** 不新建 persisted source state、全域 revision、原子快照保證、runtime health index、write API、全 point catalog 分頁工程或記憶體 live hook 改寫。

## Decisions

### 唯讀來源觀察與完整性

新增 GET /api/v1/datalink/studio-v2/recording-home/sources?workspace_id=<id>。handler 注入既有 Workspace/Device/SourceRule/Mapping/Point read interfaces，先讀已存在 singleton，不呼叫 GetOrCreate／repair／reconcile，不因 GET 建立資料。已核對 internal/datalink/workspace/read_snapshot.go 的 Service.ReadSnapshot：只委派 Repository.Get，保留 nil／取消與 detached copy；直接重用，不新增 accessor 或 authority。

回應維持現有 success/data envelope；data 為 {workspace_id, source_digest, state, devices, candidates, excluded}。state=available 才附完整 canonical 投影；unknown 附安全 reason code 而非部分成功資料；503 為 read unavailable、404 為未知 scope、400 為缺失／無效 scope。source_digest 只對可呈現欄位的確定性 canonical serialization 計 SHA-256，不是版本、readiness token 或 mutation 前提。

在 server 對同一 workspace 讀取兩次完整投影並比較；兩次都有 workspace membership、duplicate id、rule→device、mapping→rule/point/device/address 一致性檢查，任何差異或不完整均 unknown。投影不含 connection_config、credentials、values、raw backend errors。devices 僅 id/name/protocol/availability；candidates 與 excluded 精確對齊既有 buildGroupCandidates 結果，透過共同 fixtures 做 Go/TS parity。這是搬移既有 observation 驗證，不宣稱跨 repository atomic。避免只比較 Share revision 或只延長 staleTime。

### 快取呈現與請求隔離

client cache 沿既有 WeakMap 的 QueryClient/backend context 分區，包含當前 API origin；不跨 QueryClient 或 backend 複用，即使 workspace ID 相同。每個 context 最多保留既有 8 個 workspace，每個只保留最新一份摘要／digest，沿現有最舊 scope eviction；context dispose 後可回收，不逐 digest 累積歷史。這些是沿用既有容量與新增保留條件，不是容量量測。workspace identity＋摘要 source_digest 只識別該 context 的已驗證 presentation；每次 visit 仍有獨立 epoch/abort fence，不能把前次 pending promise 當此次回應。同 scope 返回立即顯示 observed cache 和 refreshing 文案，再啟動一次摘要請求；未完成 revalidation 不能標 fresh、healthy 或啟用依賴來源 freshness 的命令。失敗時可保留 last observed 但明示 stale/error，retry 唯讀。

scope 變更立即停止展示其他 workspace cache；A→B→A 的舊 epoch 回覆不可覆寫新的 A。所有 device/rule/point/mapping mutations 使相關摘要 invalidated；未知影響範圍則 invalidate 當前 workspace，而非猜未變。只有成功的新觀察可更新 source_digest。filter/page/scroll/focus 還原不會跳過現有 scoped locator revalidation。

### 漸進替換與有界首頁

Expand：新增 endpoint 與 parity tests，現有首頁不改。Migrate：首頁來源改讀摘要，保留可見頁的 bounded detail demand，不同時跑兩套 transport。Contract：新舊 fixtures 與 race cases 通過後移除首頁雙輪 HTTP adapter 的 ownership／路徑。已核對 GroupEditor 的 locator-scoped point discovery 仍呼叫 readRecordingHomeSources；將既有 adapter 原樣搬到 recording-editor/recordingSourceDiscovery.ts，僅更新該 editor 的 import／名稱與 tests，不把 editor discovery 改成 compact presentation evidence。完整設備／點位 APIs、editor await／locator／readback 語意保留；summary type 移至 shared types。共同 fixtures 繼續對 live editor discovery 與 Go projector 驗證 parity。

設計驗收門檻：一次首頁來源 revalidation 只送一個摘要 GET，不再發十個原始來源 GET；冷啟其他 workspace/group/detail requests 不冒稱零。相同 20-device/2,000-point fixture 的摘要 bytes 小於原 3.06 MB 重播，且候選與 excluded 完全相同；數字是既有合成 baseline，測試 fixture 不是產品容量上限。server 仍有讀取成本，必須量測並單列，不把 HTTP 減量宣稱 DB 延遲已減少。

### 雲端實作與最終 Mac 瀏覽器驗收

**Supersedes**: reuse-recording-home-snapshots / 不新建 persisted source state、全域 revision、原子快照保證、runtime health index、write API、全 point catalog 分頁工程或記憶體 live hook 改寫。只草擬，未實作。

依最新授權由草擬進入雲端實作，保留全部上述產品 Non-Goals。亦取代原 tasks 3.1 的「雲端 browser／不使用 Mac」限制：所有四項 change 的非瀏覽器實作完成後，才集中在 Mac 驗證 browser；不擴權、不在受限雲端重試 Chromium。非瀏覽器品質檢查可先完成，browser 與最終驗收保持 pending，archive 依最新指示保留至最終 Mac 驗收；已通過非瀏覽器品質檢查的 change 可先做明確範圍 commit，再繼續下一項實作，不因待驗 browser 停工。

## Implementation Contract

- Endpoint 必須在 internal/api/router.go 以既有 service graph 接線，正式 cmd/test_ui entrypoint 可達；nil read dependency 回 unavailable，不回示範資料。
- 重用既有 Service.ReadSnapshot；新增聚合可位於 internal/datalink/workspace/recording_home_sources.go；不得透過 handler 呼叫 mutating bootstrap 來取得 scope。既有 listWorkspaceRuleMappings 會 recover 並 ReplaceLinks，不能重用該 mutating GET helper；聚合讀原 rule/link/mapping/tag/point authority，repair-needed identity 回 unknown，不能修復或略過成部分 available。
- 保留 useRecordingHomeEvidence/useRecordingHomeFacts/cachedHomeReads 的 owner/epoch/group revision 與 visible demand 界線，不把未觀察群組當健康。
- 測試 workspace switch、late response、source mutation between passes、QueryClient/backend 同 ID 隔離與第九 scope eviction、duplicate/foreign/missing member、empty legitimate workspace、read failure、cancel/retry、filter/page/scroll/focus、bounded counts、navigation 零 mutation。
- Go focused tests 用新 handler tests 與 API router wiring test；frontend tests 以 recording-home-source-summary、recording-home-sources、route/evidence/facts/collection 為選測範圍。API annotations 變更需 make gen-docs 產生 Swagger，不手寫 generated JSON。
- In scope 是首頁來源唯讀服務與 client observation。Out of scope 是全域 source cache、command contracts、protocol/SQL/MC authority、memory renderer。若安全 parity 無法證明，保留舊 transport 並標未達門檻，不降級 unknown guard。
- 實作順序在 optimize-embedded-route-loading 的 shared page split 後；shorten-recording-save-feedback 整合本案 exact invalidation keys。

## Risks / Trade-offs

- [兩次相同投影不是 atomic snapshot] → 明確 observed 語意，所有 commands 仍讀原 scope/revision/readiness proof。
- [雙語前後端投影漂移] → 同 fixtures 對照 canonical candidates/excluded，未知欄位不直接丟棄造成 false ready。
- [快取藏住錯誤] → freshness 分離、失敗保留警示、不得把 stale 映射成健康。
