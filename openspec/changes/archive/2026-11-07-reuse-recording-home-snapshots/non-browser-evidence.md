# 非瀏覽器實作證據

## 1.1 Canonical parity fixture

共享 testdata/recording_home/source_projection.json 47 個 literal cases，先以 legacy frontend transport 驗證。新增 fixture test 加既有來源測試 63 tests pass。移除 point-device ownership guard 的 mutation 初次存活，補同 workspace 不同 device 負案例後確認 1 test failed，再立即恢復 production source 並重跑 green。Server 對照同 fixture 由 1.3 完成；此處不宣稱 server 已驗。

## 1.2 Existing-workspace accessor reuse

讀取 internal/datalink/workspace/read_snapshot.go 證實現有 Service.ReadSnapshot 已符合只讀 contract，故 ingest 將「新增」改成重用，沒有新建 authority。

執行 go test ./internal/datalink/workspace -run '^TestReadSnapshot' -count=1 通過。暫時把 s.repo.Get(ctx) 改成 s.getOrCreate(ctx)，TestReadSnapshotAbsentAndFailedReadNeverSave 如預期失敗（不存在 workspace 卻被建立）；立即還原檔案，再跑相同 suite 通過。既有測試涵蓋不存在、讀取錯誤、零 saves、取消、nil service、detached copy、兩次 observation 間刪除不重建。Production diff 為空。

## 驗收邊界

本次尚未執行 browser，所有四項非瀏覽器實作完成後才集中上 Mac。不得將上述測試當成 browser 或 field 驗收。

## 1.3 Compact projection

ReadRecordingHomeSources 以唯讀 reader 讀兩輪，完整 allowlisted source fingerprint 包含 source revision；兩輪各自驗 identity／membership，未知不回 partial arrays/digest。SHA-256 僅 hash canonical presentation（workspace、devices、candidates、excluded），不作 command proof。47 shared fixtures 與 cancellation／alias mutation／nil／privacy／revision guards 通過。新測試先以可執行 stub 產生 behavioral RED，再完成實作；presentation-only digest 與 inconsistent-reader sentinel 另有獨立 RED/GREEN。

實際執行 go test ./internal/datalink/workspace -count=1：PASS，172.312s。Focused 最終測試亦 PASS，證據 task1.3-final-green.log；完整 workspace evidence task1.3-workspace-suite.log。這些為雲端 Go 測試，未跑 browser。

## 1.4 Readonly source adapter

新 handler 僅注入 readonly interfaces；正式 mapping GET 的 recovery/ReplaceLinks 不在呼叫路徑。真實 memory service graph、寫入 spies 與前後記錄比對證明 repeated GET 不變更 workspace/link/entity。400／404／503、nil／typed nil、完整 unknown、敏感欄位排除、malformed skips、duplicate／dangling／foreign identities、legacy direct binding、disabled／skipped 與 source變動皆有測試。未配置的 link 只有兩個refs皆nil且該point無mapping時是合法空結果；repair-needed狀態不修復，也不跳過成部分available。

新增 Links allowlist 到完整source fingerprint，捕捉沒有candidate的link身分變動；digest仍僅presentation。新測試先捕捉「unconfigured link ID changed仍available」失敗，再修復通過。Go/TS共享fixture仍通過。

執行 focused handler/workspace TestRecordingHomeSources、既有 TestStudioV2WorkspaceMappingsHandler regression、go vet ./internal/api/handlers ./internal/datalink/workspace 全通過；git diff --check通過。證據 task1.4-{red,link-red,green,affected,vet}.log。

合成2,000-link in-memory初步量測：12,010次service read、48.001921ms、340,498response bytes。此非原20-device SQLite fixture，不作DB延遲改善證據；後續須單列實際endpoint成本。memory device repository的不存在錯誤未typed，安全回503；SQL sql.ErrNoRows可辨識為unknown，不解析error字串。

## Canonical label parity follow-up

共享fixtures新增BOM(U+FEFF)與NEL(U+0085)兩個獨立literal case，總數49。既有TS .trim()先通過，Go strings.TrimSpace兩案失敗後，以僅作用於display label的ECMAScript whitespace helper修復。Go workspace／handler TestRecordingHomeSources通過，TS parity＋existing source共65 tests通過，diff whitespace通過。證據 trim-parity-{red,ts-baseline,go-green,ts-green}.log。

## 1.5 Production wiring and generated Swagger

Route在既有recording routes註冊點掛入原六個service，nil依賴仍可達且回503；未建立第二套service graph。cmd/test_ui以wireGatewayServices＋migrated temporary SQLite的真實graph驗證available candidate、repeated GET total_changes不增加、dangling link保持unknown且不repair。Route／production wiring／Swagger先取得404／missing-path behavioral RED，再全部green。api、handlers、workspace、cmd/test_ui四package focused tests通過。

Swagger僅由make gen-docs重產；官方github.com/swaggo/swag module v1.16.6（binary內嵌--version字串為v1.16.4），binary SHA256 7c6aad055ddfaf90ea9810a64f9993a89d610f98d59ee4ea79f98a87bcf8cce1。task1.5-{route-red,production-red,swagger-red,final-green,gen-docs,swag-provenance,lines}.log保留命令與輸出。未執行browser。

## 原始20-device／2000-point SQLite fixture量測

從既有synthetic資料庫read-only backup到獨立可丟棄copy；20個draft device皆127.0.0.1，20個disabled rules、2000 points、2000 unbound links、0 tags／mappings。各flow獨立copy，未啟動runtime collection；量測是production-service handler的in-process HTTP，不是browser或歷史loopbackHTTP。

十次各自量測：compact單GET median196.771ms、2336bytes；legacy source十GET median417.249ms、2394466bytes；legacy home十四GET median582.428ms、3057450bytes。歷史十四GET為3058890bytes；不能把摘要單GET當成整頁只剩一個request，也不宣稱整頁browser timing已驗。

以實際legacy回應餵既有TS adapter精確比對：20 devices、0 candidates、0 excluded完全相同。此fixture是原baseline而非通用容量承諾；49個共同literal fixtures另覆蓋非空候選與負案例。Compact十次均0SQLite changes且全部table fingerprints不變。Counted handler實測4086次service read，其中4000mapping.List、40device.GetByID、40rules.ListLinks，workspace／rules.ListByDeviceIDs／points.List各2；這不是SQLstatement數。Counted handler median207.832ms，mapping list耗時median92.759ms，server仍有明確成本，不把HTTP減量宣稱零DB成本。

Legacy source十GET每次80 SQLite changes、十四GET每次160；只有devices.updated_at變更，投影欄位不變。保留此發現，不將舊讀取說成零mutation。證據performance-summary.json、performance-sqlite-*.{json,log}、performance-canonical-parity.{json,log}、performance-legacy-device-changes.json與raw responses、commands、ephemeral harness原文。暫時repo測試檔已移除，無常駐server。

## Backend lint

官方golangci-lint v2.6.2（Go1.25.5）對api、handlers、workspace、cmd/test_ui執行：35 findings中25個B變更問題已修復。剩餘10個baseline findings的檔案逐一與ca372b31920917a5370df948a9f306a4f64b09f7比對byte-identical，未擴張修復範圍；同四package --new-from-rev查核0issues。四package focused tests（含49fixtures、trim/link、Swagger、wiring）重跑通過。未宣稱unfiltered lint全綠。詳backend-lint-{unfiltered,final-unfiltered,new-from-base,baseline-attribution,fixes-tests,lines}.log。

## 全批次 backend gate（parent 執行時的 backend snapshot）

實際 go test -p4 -json ./...：FAIL，422.607s，14762 test/subtest pass、523 skip、1 fail；59 packages pass、6 no-test skip、1 fail，無cached結果。唯一失敗 recordingplan/TestStorageRestoreIsDoneOnceWhenMeetingConcurrently 已在隔離20次重跑重現，local dependency closure 21 package variants／912 source+embed files均與原main0f54740byte-identical；不是B變更。未修復範圍外concurrency flow。Opt-in external DB／acceptance skips不是通過。

go vet -p4 ./...：PASS，6.962s。執行期間3222 backend source/config fingerprints不變。

全域 golangci-lint run --allow-serial-runners ./...：FAIL，41 issues；39位於29個與原main0f54740byte-identical的repo檔案，2位於已安裝npm flatted3.4.4的Go helper（非repo authored code）。沒有本次A/B changed issue file。未改baseline或node_modules；B四package new-from-base lint仍0issues。

上述證據在同批次final-backend/{summary.json,results.json,tests.log,vet.log,recordingplan-dependency-identity.json,recordingplan-isolated-20.log,lint.log,lint-baseline-attribution.json}。不宣稱全backend gate全綠；已知baseline／dependency失敗與B focused通過分開。

## 2.1 Client observations and context fencing

首頁已改單一summary讀取；每visit保留epoch／AbortSignal，query key含QueryClient／backend identity與完整API base URL partition，再附workspace／revision／epoch。每partition只保留最新八個workspace觀察，unknown／error不提升freshness。返回立即可呈現observed＋refreshing；錯誤保留stale警示，authority仍要current visit success。

Task2.1先以service 3個behavioral RED、route migration 5個behavioral RED落實；初次14files／189tests、TypeScript、scopedESLint、diff／linegate通過。後續補duplicate point claim／blank identity檢查（10個behavioral RED，修復後3files／50tests PASS）。Captured API base URL固定於request；presentation cleanup亦用captured context，backend／QueryClient變更會重建scoped presentation。

React concurrent render的兩項負案例先失敗：未commit的新A visit會毒化舊A refresh，未commit的第九workspace會evict舊scope。修復將entry建立／eviction／epoch registration移至query execution，render僅read，原epoch／abort／generation／entry-identity fences保留；4files／52tests與TypeScript／ESLint通過。證據task2.1-{service-red,route-red,focused,identity-red,identity-green,concurrent-red,concurrent-green,tsc,eslint,lines}.log。後續全前端gate另列，browser仍未執行。

## Review 修復：不將未知 skip 資料當成空集合

獨立review確認 mapped rule 的 persisted skipped_addresses=null 被解為nil並normalize成[]，會把legacy Home不能呈現的資料升成available。四個behavioral RED後修復：explicit null／non-array回unknown；mapped empty string回unknown；真正未配置、無mapping的empty string保留available；正常[]不變。十個persisted-source cases檢查repeated GET與零repair／零mutation。新增檔案為studio_v2_recording_home_sources_skips_test.go，修改read.go及cost_test.go（cost fixture明確使用[]）。獨立review已重新檢查並確認Warning resolved。

此三檔在上述full-backend run之後變更，因此full-suite結果對它們是prior snapshot，不冒稱當前全套已重跑。當前affected驗證實際命令：
- go test ./internal/api ./internal/api/handlers ./internal/datalink/workspace ./cmd/test_ui -run 'TestRecordingHomeSources|TestReadSnapshot|TestProductionRecordingHomeSources' -count=1：PASS。
- go vet ./internal/api ./internal/api/handlers ./internal/datalink/workspace ./cmd/test_ui：PASS。
- golangci-lint run --timeout=10m --new-from-rev 0f404446 ./internal/api ./internal/api/handlers ./internal/datalink/workspace ./cmd/test_ui：PASS，0issues。

使用Go1.25.5／golangci-lint2.6.2，證據backend-skips-{red,final-green,final-vet,final-lint}.log；未改其他backend modules／DTO／Swagger／49共同fixtures。先前成對效能時間為此窄guard修復前量測；final guarded fixture的bytes／canonical／zero-write另做確認。

## 2.2 Mutation owners and bounded demand

成功source ACK以captured QueryClient／backend／完整baseURL／workspace精確失效，並獨立啟動active observer refresh，沒有新awaited readback／save barrier。Inactive Home不發request，返回才重新觀察。Optimistic mutation在onMutate入口、cancelQueries await前捕捉scope，原rollback／onSettled行為保留；locator mode從既有owner傳入workspace，不猜可能殘留的bootstrap。Device／rule／point／mapping與影響tag label的owners、direct wizard mapping均有實際hook測試。

Home保持mount於source modal後方；實測普通close原先沒有revalidation。現在成功mutation owner安排單次background summary，移除onCreated跨context廣泛refresh。測試保持summary pending而mutateAsync仍依既有barrier完成、modal返回更新名稱、late ACK不影響新的scope。1000group測試保留25個detail demand上限與全部request count；僅將該大型fixture的測試timeout調至15s，沒有降低fixture或產品延遲承諾。

Final affected selection116files／1275tests PASS；manifest-scoped ESLint、full TypeScript、diff／line gate PASS。證據task2.2-{workspace-red,workspace-green,owner-red,owner-green,background-red,background-green,final-focused,tsc}.log與36-path manifest。

## 2.3 Contract relocation and final frontend gate

GroupEditor仍有明確locator-scoped discovery caller，所以依已ingest且經parent確認的scope，把原adapter原樣搬至recording-editor/recordingSourceDiscovery.ts／readRecordingSourceDiscovery；old home adapter路徑刪除，首頁只用compactservice。Transport／projector／guard function bodies在symbol rename後byte-identical；49共同fixture與editor caller／locator回歸仍有效，不改editorauthority／await／readback，也不刪full source APIs。

Relocation focused8files／152tests PASS。首次full frontend run暴露五個舊minimal QueryClient／Axios mocks缺新增context依賴，以additive fixture補齊，所有既有assertions不變，初次失敗log保留。最終實際執行npm run lint、npm test -- --run、npm run build皆exit0；655files／7567tests全部PASS。Build保留既有large-chunk warning，未宣稱零warning。Final diff／line gates PASS。證據frontend-{lint,test,build}.log、frontend-exit-statuses.txt、frontend-test.initial.log、task2.3-relocation-verification.txt與17-path manifest（含rename兩端）。最終embedded static sync／batch build和Mac browser仍由整批後續處理。

Final guarded-handler再查原20／2000fixture：available且exact2336bytes，與先前compactresponse byte-identical、與legacyTS canonical projection精確相同；0SQLite changes且all-table fingerprints不變。此次單trial378.953ms時fullfrontend正在並行，不能取代先前成對比較。證據performance-final-guard{.log,.json,-result.json,-compact_route-00.json,-commands.txt}。暫時harness已移除。
