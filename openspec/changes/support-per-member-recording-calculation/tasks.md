## 1. 持久化與相容基礎

- [x] 1.1 共同新版demo owner確認後，以Red-Green-Refactor實作 Recording members own calculation intent 與「Member意圖與新mixed形狀獨立於來源」的typed member method/profile/resultshape；測unknown/omitted拒絕、同source不同record不互改、groupdefault不覆explicit與legacy未改。
- [x] 1.2 [after: 1.1] 以Red-Green-Refactor交付「儲存擴展與草稿能力分期」的nullable member/resultshape與內部mixed snapshot寫入／scoped readback，實作 Legacy profiles keep their original meaning 與 Persisted reduction and independent value roles；安全擴展write_groups CHECK，測可丟棄SQLite migration/reload/restart、rollback／repair重入、legacy完整列與既有CRUD/CAS不變、foreign identity拒絕。正常mixed mutation／activation保持受阻，不虛構group interval。
- [x] 1.3 [after: 1.2] 實作 Member windows align independently of activation 與「每member對齊時間窗與結果」的member cadence驗證（使用backend bounded capability，不新增任意上限）與member適用性；測15秒／5秒整秒phase、60秒／10分鐘整分phase、非邊界啟用與restart不漂移、fraction/unknownunit、bool/text mean、未確認counter delta與衝突group reduction失敗；接通正常mixed草稿scoped CRUD／CAS及「來源共享與記錄身份隔離」，以record-scoped canonical identity支援同DB不同table的同source獨立method／cadence，測重送同record不重建、foreign scope／stale拒絕、同target衝突及legacy create-once不變；runtime／proof gates未完成仍不得activation。

## 2. 窗口派送與正確結果

- [x] 2.1 [after: 1.3] 實作 Different member cadences preserve sparse due semantics，擴展既有snapshot mixed mode，只對delta members建立counterbinding並使用各member對齊window，依「已確認不同週期輸出與邊界」實作wide聯集／long due-only、[start,end)與first partial，測exact-boundary sample歸下一窗、required-not-due不skip及長表到期失敗原因；測兩gauge mean、一counter delta各自window及latest/mean不產生baseline，不另建reducer engine。
- [x] 2.2 [after: 2.1] 實作 Counter policy and member state remain isolated，重用CloseCounterWindow；測firstbaseline、signed/rollover/reset、same-source不同group baseline隔離、method/binding revision變更不能借舊checkpoint。
- [ ] 2.3 [after: 2.2] 實作 Missing results obey explicit completeness policy 與「缺值與baseline進展明確」；測optional NULL+reason、required整row skip、gapspan不進單窗delta欄、不補raw/0，以及缺值mean coverage保留。
- [ ] 2.4 [after: 2.3] 實作 Mixed result columns are not legacy raw columns 的typed managed-wide及long encoder；測mean/mean/delta結果欄22/12/8、exact大數與decimal、provenance raw/from/to分離，legacy raw欄不改，wide not_due=NULL且長表未到期零列、到期required僅skip所屬輸出，unsupported layout明確受阻。

## 3. 安全gate與durable整合

- [ ] 3.1 [after: 2.4] 實作 Calculation changes invalidate configuration proofs 與「版本與durable closure不能分離」的candidate/digest/resulttype/schema preview/readiness；測method改變過期proof拒絕及nullable partial能力需fresh確認，legacy轉mixed走既有安全transition。
- [ ] 3.2 [after: 3.1] 串接production runtime boundary與counterproof resolver，frozen descriptor完整包含membermethod；測stale/foreign observation、member變更、source/mapping/measurement revision與exactscale保護，舊sample不能新method重解。
- [ ] 3.3 [after: 3.2] 實作 Closure and counter progress are atomic and replayable；以ledger故障注入驗證row或skip、checkpoint、bucket/outbox同tx，required首窗skip仍可restart後下一窗delta，不部分advance。
- [ ] 3.4 [after: 3.3] 完成immutable replay/receipt復原測試，覆crash前後、重試payload衝突、SQL已commit回覆遺失、copy不帶baseline；驗證stable effect不重複寫或重算用量，unknown查證不新建row。

## 4. 生產證據與交接

- [ ] 4.1 [after: 3.4] 實作 Production proof precedes editor support 與「分期擴展與consumer交接」的API/types/capability輸出；測legacy入口仍可用、simple editor讀回member計算/型別/revision而非global Point，不先開未完成能力。
- [ ] 4.2 [after: 4.1] 使用既有同機可丟棄SQLite執行真scoped save→readiness/schema確認→runtime→durable closure→SQL write/readback，驗證共同到期22/12/8一列、15秒／60秒不同窗的wide not_due與long due-only、exact邊界／first partial、NULL/skip/receipt；fake transport，不連PLC/正式DB、不新增環境。
- [ ] 4.3 [after: 4.2] 獨立review並執行相關workspace/snapshot/aggregation/runtime/dbtarget/groupdelivery/API測試與Go lint/vet、diff/line gate；本案完整契約完成後交給editor/adoption，整功能最終才完整總測，明列PASS/FAIL/NOT RUN與未過必要gate，不以helper或demo代替真SQL。

需求變更註：原共同window／1–300秒限制被新每member cadence取代；storage策略已明確確認，未勾選項目仍待完成；正式實作授權已依docs/verification/simple-recording/implementation-decision-2026-10-30.md取得，不代替產品驗收。1.2限定storage expand；原正常mixed CRUD驗收移至1.3待cadence與record identity完整後交付，工作總數與完成歷史不變。
