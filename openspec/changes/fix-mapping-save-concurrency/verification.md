# 映射並行保存驗證紀錄

本次範圍：同規則保存與來源 mutation 的協調、原子連結替換、前端規則排隊與安全失敗／明確重試。基準 `d9ba336ecd575ec7f71950677c4ad64af1552ecb`。修正、完整檢查與兩套實際隔離 UI 均已通過，最後審查發現bootstrap readiness背景刷新退步，已補最小修正／六項真實觀察者回歸；變更版本完整前端／mapping UI與Basic／recovery均通過。最後 Git HEAD／clean 狀態於交接時釘選。

## 已重現與修正的依據

- 真 SQL pool 1 的八筆並行保存原先出現 409／500，受控連結 DELETE 窗口可造成另一列 404；同 fixture sequential control 全部成功。
- 新回歸檢查 pool 1／4、每池四次八筆並行保存，以及完整 point／tag／mapping 身份、pipeline、accepted hash。pool 4 使用 production WAL 與 busy_timeout15000 設定，沒有重送 SQL 來掩蓋錯誤。
- 原子替換使用同一 SQL transaction，插入唯一鍵失敗或第二筆故障時，reader 看完整舊集合。Memory 單一 lock swap 並複製資料。
- 後段 PUT 失敗完整還原上一個已接受狀態；直接修改過的 mapping／tag 使用 exact-field CAS 拒絕補償覆蓋。新的 POST 可保留唯一 owned disabled unconfirmed draft，再以同身份明確重試。
- SQL Update 漏寫既有 tag_id 的問題會使改 Tag Key 保存失敗；修正既有欄位，沒有 schema／request DTO 變更。
- DELETE 先更新連結，mapping 刪除失敗可恢復連結；已刪除／已保存而只有 orphan tag cleanup 失敗，區分完成結果與清理警告。

正式 production confirmation 測試的舊預期要求失敗後留下 disabled pending pipeline。新的補償會恢復完整 prior accepted pipeline，因此改為檢查 prior pipeline／tag／hash／status／enabled 均恢復、derived sync 不採用未確認變更，以及移除故障後明確重試成功；沒有移除「失敗不得發布 pending pipeline」契約。

## 規格與測試對照

| 規格情境 | 實作／可執行測試 | 本次結果 |
| --- | --- | --- |
| Eight same-rule rows save concurrently | rule_mutation_guard、WorkspaceMappingEightConcurrentSQLiteSaves | focused／完整 Go PASS |
| Atomic replacement fails | repository_links_atomic、SQLReplaceLinksFailureRollsBackWholeSetAndReadersNeverSeeDeleteGap、SQLReplaceLinksUniqueFailurePreservesExistingLinks | focused／完整 Go PASS |
| Source mutation races save | WorkspaceMappingSaveCoordinatesSourceMutations、LateDirectManualEditPreservesExactCASAndAcceptedHash、ProductionWorkspaceConfirmationOnlySuccessfulOwnedSave | focused／完整 Go PASS |
| Another rule progresses | WorkspaceMappingOtherRuleSaveAndDeleteProgressWhileRuleHeld、前端 deferred global read | focused／完整 Go PASS |
| Repeated batch and rapid edits | mapping-autosave-queue、mapping-autosave-queue.lifecycle、隔離 mapping-save-concurrency UI | 實際 mapping UI PASS |
| Failure and explicit retry | mapping-save-error-ui、failed PUT／DELETE／POST SQL 回歸與隔離 fault→Retry UI | 實際 mapping UI PASS |
| Default numeric pipeline remains exact | mapping-autosave-queue.numeric、ProductionDefaultUIUint64ToSQL、ProductionDefaultUIScaledInt16ToSQL | focused／完整 Go PASS |

案例值：八列 int16 40001–40008；uint64 `9007199254740993`；int16 `243` × `0.5` + `10` = float64 `131.5`。全部情境均有跨呼叫端或失敗分支契約，沒有以表面樣式理由排除測試。

## 審查修正

1. owner resolver 不得等待不相關 rule 的 guard；使用實際 owner 再在鎖內重讀，且保留原有 ownership mismatch 行為。
2. 舊請求拒絕不得把失敗資訊套到已替換的新列；以 generation／身份保護成功與失敗兩條路徑。
3. 每筆 mutation 與批次 refresh 不得等待會被其他 rule 擋住的全工作區查詢；同時避免其舊 snapshot 覆蓋新確認值。
4. accepted PUT／POST 已完成但 orphan cleanup 失敗，不得回報保存失敗／引導重試已完成的操作。
5. 舊查詢不可漏掉查詢開始後新增成功的映射；補上 POST acknowledgement merge，同時確認下一次 authoritative absence 可以移除該 ID，避免永久保留。
6. 背景bootstrap readiness不能僅標記stale而沒有refetch；恢復原生active query invalidation，但不等待GET。六項真實observer回歸先RED再GREEN，涵蓋create/update/delete、blocked GET仍讓其他rule下一列前進、取消舊response，以及metadata故障不改變已接受save結果。
7. guard registry 測試等候真正 release 完成後才檢查清理，避免測試自身競態。

以上修正已有 focused 回歸及跨 owner 重審；後端與前端沒有未解審查項目。完整 frontend gates 與實際 mapping UI 已通過；Basic／recovery 已重跑通過。

## 先前0ba／9dd版本完整檢查與 UI 證據

最終後端檢查已通過：`GOMAXPROCS=2 go test -p 1 -parallel 1 ./...`（49 個含測試套件、7 個無測試套件）、`GOMAXPROCS=2 go vet -p 1 ./...`、`GOMAXPROCS=2 golangci-lint run --concurrency 1 ./...`；均退出0。首次 lint 28項問題按既有規則修正，沒有 exemptions。

前端 `NODE_OPTIONS=--max-old-space-size=2048 npm run lint`、`npm test -- --run --maxWorkers=1 --no-file-parallelism`、`npm run build` 均退出0；完整198檔／1269測試通過。新 focused 排隊／safe UI／numeric suite13檔／88測試通過。

隔離實際 mapping UI：`mapping-concurrency-20261005b.json` PASS。8個POST新增、66個PUT成功、32個preview POST成功，只有明確注入的1個PUT500；無409／404／額外500。六次SQL snapshot每次8個唯一point/tag/mapping，三次batch enable與bulk float64/0.5/10保持身份。快速修改最終值、131.50顯示（數值131.5）、保留草稿、不自動重送、故障移除後明確Retry成功均通過。已實際看過error-state與Saved 截圖，page errors0。browser closed，兩個owned程序正常退出0／無KILL，namespace已刪。

程式提交 `0ba22461397e443d1150d8553151a653a6320179`；只有驗收格式修正的提交 `9dd92a34a7425764d96b65153d3ab54285913c47`。UI binary在9dd92a34建置；完整source manifest雜湊 `d5eb24320e8052fe02fc42d6e7a4bbaea01f8209a36d3de6158f34a5233a6d4f`，記在成功JSON。Go build metadata的modified=true來自尚未提交的本任務驗證紀錄／證據；產品code沒有dirty差異。最後證據提交前後會核對同一source內容，不冒稱binary的VCS revision已包含後續純證據提交。

首次 UI runner按顯示文字比對131.5，而實際UI正確顯示131.50＋Live preview，造成假失敗。保留a.json／a-error.png／a-context.json，改為數值比對並保留顯示原文，再用新namespace b重跑通過；未覆寫首次結果。

首次完整frontend有6項失敗：locale檢查只找舊namespace、Step3文案mock不識別新namespace，及DB fixture人工ID／float64 width1與40002不符合derivePoints。修正為正式生成ID／width4／第二點40005，以及valid DELETE runtime metadata；三個DB測試檔與原始hydration/table-order/independent-target斷言沒有修改。locale測試仍要求所有canonical安全code，新增configured resource與getFixedT雙語驗證。再重跑198檔全通過。

Basic／recovery：`mapping-concurrency-20261005-basic.json` PASS，使用相同 gateway／simulator SHA-256。兩台各8量測點（register寬度依型別），八列來源型別預設保持（其中兩列uint64），包含uint64 9007199254740993與18446744073709551615、string／bool／float。各有3個連續60秒row／bucket／outbox／receipt獨立讀回；A215→B187→A215實際SSE切換沒有外設備tag。缺表修好後Basic Retry送達一次，poison row明確確認Skip→operator_skipped，不假稱送達，frozen destination／payload身份保持，後續bucket能送達。停用drain後正常重啟，停用組SQL列數不增加，另一啟用組仍前進。390／768／1440無文件水平溢出，page errors0；已看過Basic／restart截圖。browserclosed，2sim＋2gateway（含restart前後）正常退出0／無KILL，namespace removed。

所有 commands、退出碼、source／binary SHA、UI結果與cleanup摘要見 `docs/plans/studio-v2-flow-completion/evidence-f/mapping-concurrency-20261005-checks.json`。完整與focused RED／GREEN原始logs位於同目錄的 `mapping-concurrency-20261005/`；沒有把 sandbox bind 阻擋或共享source mutation等待列為行為GREEN。原有migration／legacy take-over／unknown safety回歸隨完整Go套件重跑；實際UI不宣稱正式PostgreSQL／permission／unknown重送驗收。

OpenSpec與最後行數檢查見checks摘要；所有7個scenario有對應可執行回歸，補上5個具體Examples，沒有增加需求／範圍。

首次四套件非提升權限試跑遇到 `127.0.0.1` bind sandbox 阻擋；此項是環境阻擋，不算行為 RED。首次授權完整 Go 檢查僅 production confirmation 舊狀態預期失敗，保留原結果，再重驗修正後 source。

外部 Claude reviewer 曾在 shared worktree 將 SQL transaction executor 改成 db executor 做 mutation test，導致其自身與兩個本任務舊 binary 測試等待；只恢復本任務 source，不停止外部程序。確認其預定恢復副本與修正檔 SHA-256一致，待外部程序自然結束，重新通過 focused 與完整 Go gates。此期間的 stalled logs 不列為 GREEN 證據。

## 最後背景刷新修正版本

修正提交 `2c6919da2ff2d1df1fb7ae757035b81a71f45199`，完整前端199檔／1275測試、lint、build通過；focused19檔／109測試及TypeScript／scoped lint通過。六項真實bootstrap observer regression在舊mark-only行為全部RED，再GREEN。Go source與已通過全量gates的0ba22461逐路徑diff完全相同，沒有重跑相同版本檢查。

新建置source manifest SHA-256 `d7ac58695889fa2ebad349fb06caeda075dd9b30b82c4a7d6b2ff026eb49b57c`。新namespace `mapping-concurrency-20261005c.json` PASS：8 POST201、65 PUT200、32 preview200、只有owned注入1 PUT500，無409／404／額外500。六次SQL8列identity、latest draft、131.50／數值131.5、明確Retry、safe UI皆通過；已檢視兩張截圖，browserclosed／2owned程序exit0／無KILL／namespace absent。舊a／b／Basic結果完整保留，不覆寫。新 `mapping-concurrency-20261005-basic2.json` PASS且binary SHA與c相同：兩台各八量測點／三個60秒row、default exact uint64、A215／B187／A215、缺表Retry→sql_committed、明確Skip→operator_skipped與後續row、停用重啟前後6→6且另一組前進。390／768／1440無文件overflow，page errors0，七張最新Basic／runtime／recovery／restart截圖均已檢視。browserclosed，四個owned程序正常exit0／無KILL，namespace removed。

## 邊界與交接

- guard 僅協調同一 gateway process 的同一 sourcerule.Service；不冒稱跨程序鎖，持久 source revision／pipeline CAS 繼續有效。
- 單列 mapping enabled=false 目前由來源 rule lifecycle 決定，本次保持既有政策；全部啟用／批次轉換有本次驗收，不能宣稱單列停用持久化政策已修好。WriteGroup 停用／重啟另以原有驗收重跑。
- 不操作原本3222／5173的服務、使用者既有 DB／設備。最後只讀程序／port檢查，先前PID76744／76814已不在，3222／5173沒有listener；本任務沒有停止它們，不推測原因或冒稱仍在執行。只用模擬 PLC 與 run-owned SQLite；沒有真 PLC、Windows、LAN、長時間 soak 或正式 PostgreSQL 驗收。
- 實際runtime切換截圖仍出現live-stream-degraded與projection-stale提示，同時選定設備八列值正確、fresh且沒有外設備tag。此批驗收覆蓋設備身份與數值，不宣稱整個dashboard健康提示已清乾淨；原因未在此批證明。
- service.go 與 service_test.go 為歷史超長檔，本次分別減少20與9行；規則 guard／原子替換已抽出。後續若需增長，先把 CRUD／rollback 分拆，不能依靠本次例外擴張。
- Spectra CLI 此版本在 shared-Git worktree 解析原專案 root；artifact validate／analyze 與實際 source scope 分開。程式審查／最後身份以 worktree Git 差異、source content manifest 與 exact HEAD 核對；不把原樹 CLI snapshot 宣稱為新程式審查。task checkbox 在已驗證後以實際工作樹更新，沒有有效 worktree task baseline 時不偽造 source attribution。
- 只本地提交；不 push、merge、archive、deploy。原工作樹保持原 branch／HEAD，最終只移除本任務暫存的未追蹤 change artifacts。隔離worktree於最後同一HEAD detach，釋出原repo的 `fix/mapping-save-concurrency` 分支供使用者切換；不自動更新現場。
