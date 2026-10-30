## 1. 保留既有資料並建立可追溯意圖

- [x] 1.1 以 TDD RED-GREEN-REFACTOR 交付 Preserve legacy decoding semantics／舊版語意以 legacy 保留：在現有 adapter fixture 鎖定 Modbus fallback、MC3E connection、FATEK component assembly 與舊16-bit行為，再新增相容 resolver；驗證升級／restart／無關編輯前後值不變。
- [x] 1.2 [after: 1.1] 交付 Point format intent and effective order／點位意圖是唯一 authority：additive configuration migration、Point mode/revision、inherit/override/reset 與讀回有效來源；先寫模式相等但意圖不同的失敗測試，再實作，驗證 disposable DB roundtrip／migration repeat，省略欄位不變且不猜舊intent。
- [x] 1.3 [after: 1.2] 交付 Applicability and safe decoding failure：以 protocol/address width/type 的 capability matrix 驗證新32/64-bit格式、拒絕新的bool/16-bit/string override、保留legacy及新inherit的driver既有16-bit byte語意，鎖定64-bit全word反轉；以 resolver boundary tests 驗證 unknown格式、缺word及不適用不會QualityGood零值。

## 2. 安全編輯與規則生命週期

- [x] 2.1 [after: 1.3] 交付 Scope and revision guarded updates／原子更新與安全失敗 的 scoped mutation 及舊API相容入口；以 handler/service tests 先重現foreign/stale與重播問題，再實作ownership、CAS、request_id同payload重播與異payload409，驗證revision只前進一次及typed errors。
- [x] 2.2 [after: 2.1] 交付 Preserve overrides across rule revisions／規則衍生以 identity 保留且版本化：先重現 Service.Update 抹除override，再使rename/count/skip/default/enable/rematerialize保留既有point意圖，新point inherit，skip/unskip與range slide不依index搬移，同址新identity不繼承舊override；以 sourcerule persistence/restart tests 驗證全部案例。
- [x] 2.3 [after: 2.2] 交付不相容type變更拒絕、rollback及readiness失效的一致性；internal materializer/rollback不被當operator intent；以 fault injection tests 驗證Point／rule／candidate任何一步失敗不留半套格式，uncertain時受影響採集fail closed，舊response不覆寫新revision且Share gate不弱化。

## 3. 逐協議正式讀取接線

- [ ] 3.1 [after: 2.3] 交付 Protocol decoder parity／共用正規化後逐協議解碼 的 Modbus TCP/RTU/UDP 分支：scheduled/manual傳resolved format，legacy分支不變；以既有 adapter fake transport 加同連線A/B交錯與六種numeric type×四格式fixture驗證，不修改共享converter。
- [ ] 3.2 [after: 3.1] 交付 MC3E 的逐request format decoder，保留legacy connection政策及既有count語意；先寫override被忽略的RED fixture，再實作GREEN並REFACTOR，以MC3E production Read測試確認兩點不同格式、32/64-bit與缺word失敗。
- [ ] 3.3 [after: 3.2] 交付 FATEK 的逐point格式與16/32-bit component正規化，保留legacy assembler；先寫忽略format的RED fixture再實作，以D/DR fake transport檢查全部numeric type/order、request count、不同點交錯與短response，禁止把它誤稱connection-default解碼。
- [ ] 3.4 [after: 3.3] 交付 Version decoding identities and observation proofs：SourceRevision／AcquisitionConfigFingerprint版本化並含mode/effective policy/default來源，install/bind/poll snapshot同步；驗證MC相同舊DataFormat但legacy→override的fingerprint必改、舊proof/late sample拒絕，並交付resolver的不可變collector snapshot／manual路徑與restart一致性；以collector和adapter相關race／integration測試驗證併發兩點不同格式無污染、revision切換不重新解碼歷史量測。

## 4. 使用端契約與整體驗收

- [ ] 4.1 [after: 3.4] 交付 Production evidence and consumer handoff 的共用API schema、frontend types/services及en/zh-TW安全錯誤契約，供 streamline-recording-task-editor 顯示mode/effective source/applicability；以API contract與前端service tests確認unsupported不可假存成功，annotations變更經既有generator產生Swagger，不在本案改布局。
- [ ] 4.2 [after: 4.1] 經 cmd/test_ui 真實wiring測試 rule save→point reload→collector→adapter→ingestor 與reset/restart/rollback，完成全protocol/type/format acceptance matrix；使用既有fake transport/disposable設定DB，列出實跑與NOT RUN，不以helper測試替代production證據。
- [ ] 4.3 [after: 4.2] 完成expand→migrate→contract與downgrade安全演練：新override存在時以受管理回退程序／相容檢查阻止不相容downgrade（不新增部署engine，不聲稱攔任意舊binary），驗證設定備份還原不改歷史SQL；本案交接執行受影響測試與必要integration/acceptance、git diff --check、make check-lines、Spectra verify/review並交接測試矩陣；完整前後端全套檢查留到七案整批實作完成、最終交付前統一執行，不在單案交接提前總測，未做真PLC／正式DB驗收須明列。

需求澄清（2026-10-30）：依 owner「做的期間只跑相關測試，直到完成才跑全部」，上述全套時點指七案整批最終交付，取代可能被誤讀為本案結束即總測的舊措辭；各 task 必要focused integration、production wiring與驗收矩陣不減免，task數與依賴不變。
