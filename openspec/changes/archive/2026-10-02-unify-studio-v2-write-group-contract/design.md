## Context

已讀proposal與現況證據。Step4的membership、plan selection和enabledTargetCount各自使用不同資料。既有workspace DatabaseRowGroup及target refs需保留；本案只建立單一authority與相容讀寫，尚不啟動新delivery。

現有定位：internal/datalink/workspace/service_database_row_groups.go、internal/api/handlers/studio_v2_workspace_mappings_handler.go、frontend/src/features/datalink/workbench-v2/steps/step4/recordingPlanMembership.ts。

詳細來源見 [evidence](../../../docs/plans/studio-v2-write-groups/evidence.md)，共用欄位與政策見 [跨案契約](../../../docs/plans/studio-v2-write-groups/contracts.md)。這是擬議設計，不是完成報告。

## Goals / Non-Goals

**Goals:** 基本Tag寫入有一份可保存、可重載、可驗證的WriteGroup；existing配置保留原資料並以明確審閱相容；新舊不雙寫。

**Non-Goals:** 重寫framework、改採集協議、報表/aggregation/retention產品、外部DDL或現場操作。

## Decisions

1. **單一authority與optional semantics。** 採跨案契約A的欄位。WriteGroup保存persisted Tag與來源point/device，managed/custom僅選layout策略；measurement ID可選。拒絕繼續用recording plan、row group和db target各持一份可獨立修改的設定，也不為基本raw Tag自動編造semantic_kind。既有measurement/plan高階模式仍保留，但非基本readiness前置。
2. **API與CAS。** 擴充現有workspace API，以 `/studio-v2/workspace/write-groups` 提供list/create及 `/:id` get/update/delete；所有mutation要求expected workspace/group（create除外）/connector revisions。一次本地transaction保存group、相容projection及workspace revision；第二段故障全回滾。unknown/foreign資源同安全404，stale409，invalid422。新增API不自行建表、不啟用設備。契約及Swagger在實作批次成套更新。
3. **只讓一份設定可寫。** group啟用前先保存draft；activation在既有workspace/settings/readiness barrier下原子切換applied revision與writer ownership。舊API修改能無損轉成canonical patch則同transaction/CAS更新；不能對應的payload回409及修復動作，不能默默成功寫到過時projection。新delivery接線完成前保持舊runtime路徑，不提早宣稱durable。
4. **分批移轉。** 第一批simple single-point/custom mapping經確認轉為snapshot草稿；第二批row-group含shared columns/unique metadata；第三批 recording plan 完整預覽及等價阻擋；目前不可證明等價的 plan 保留原資料，基本寫入由 canonical WriteGroup 建立。每批dry-run列出原ID、新ID、scope、row layout、identity與差異；review/apply要同revision。若shared-column將落到同一row、缺unique key、來源模糊、advanced stream、未支援upsert或混合目標不能表達則blocked，原配置保持原狀。
5. **保存原意與rollback。** migration provenance及ID map持久化，重跑相同來源revision只得到同一group。每批比較before/after persisted payload與legacy輸出意圖，不把column碰撞一概禁止或一概放行。舊read API從canonical projection回覆已移轉資料。等同一scope的所有consumers已轉接、backlog有明確歸屬且回歸通過後才撤舊write authority，不能先刪表再搬資料。

6. **基本lifecycle。** rename保留group ID及既有record/receipt identity，不重啟已套用writer；所有edit仍走CAS。增刪member、改column/policy/destination保存新draft，未Apply不改running revision；成功Apply在下一個bucket邊界切換，舊bucket用原snapshot收尾。disable停止新intake，已ACK backlog按原snapshot繼續。delete是soft-delete/tombstone並停止新intake；保留不可變payload、revision、ownership、operation/receipt及backlog查詢，worker仍可完成既有delivery；不得重用已刪ID。無法保證backlog歸屬與查詢的舊格式先拒絕delete並提供disable，不能cascade刪accepted記錄。實體purge不屬本案。

### Reviewed single-mapping snapshot conversion

**Supersedes**: unify-studio-v2-write-group-contract / 分批移轉的第一批等價映射前提

已核對legacy無GroupKey的writer：每次採集直接執行SQL，write interval並未限流，insert不寫timestamp column；這與B的週期snapshot不等價。使用者於2026-10-02確認採建議：preview明列逐筆→snapshot、bucket選值／late／missing政策與timestamp差異；review必須顯式確認差異，才保存canonical draft，不自動Apply。

替代方案是全部blocked直到另增every-sample模式；使用者已選擇本次snapshot轉換，因此不新增every-sample產品模式。來源不明、grouped布局、advanced與未支援upsert仍blocked，不能以確認繞過。

### Blocked recording-plan migration

**Supersedes**: unify-studio-v2-write-group-contract / 分批移轉的第三批可等價 basic recording plans 前提

使用者於2026-10-02確認：現有 recording plan 缺少 canonical column／row identity，raw_history／on_change／triggered snapshot 與週期 snapshot 無法證明等價。保留等價要求；本批提供完整原 intent、來源版本、blocked 原因與修復提示，不自動保存群組，不改 plan 的 revision/status/applied identity；新基本寫入沿 canonical WriteGroup 的 Create 路徑完成。

未採用「猜欄位或將每筆記錄改為 snapshot」及另增舊 every-sample runtime；這些會改變已保存意圖或擴大本案。沒有正向等價 fixture 時，不把無法移轉者勾成已移轉。

### LegacyWriteAdapterConflict

本案既有設計允許 canonical transaction 或 actionable409；A3.1 採後者。global target CRUD 與 Studio target/row-group API 在來源已由 canonical 群組擁有時拒絕修改，提供 `open_write_groups`；不另實作缺 revision envelope 的無損 patch 翻譯。未移轉 scope 保留原 CRUD 與 runtime，C 前不切 writer。

### WriterOwnershipActivationBarrier

A3.2 準備 `WriteGroupService.Apply` 既有 local transaction 內的可注入 owner transition；實際 production consumer 與 activation route 由 C 接線。沿用 `UpdateDatabaseSetup` 的 mutex/CAS/commit/rollback，不新增平行 owner table、runtime 或 worker。沒有安裝 barrier 時，domain Apply 仍只是既有 bucket/version 排程，不能把 applied revision 說成已切 production writer。

Share 的 workspace revision、settings revision、readiness token 維持現有 activation preflight；它們不是 `Record.DatabaseSetupRevision`。不在 local transaction 巢狀呼叫會另開 transaction／外部 inspector 的 readiness、settings CAS 或 Share revision store。

## Implementation Contract

- 範圍內：canonical群組的transaction/CAS、readiness、lifecycle、三批reviewed migration與單一writer ownership；範圍外：新協議、外部DDL、正式DB／PLC及報表。
- A2.1 preview只讀本地設定，回傳workspace/connector revisions、adapter version、review digest、來源revision、BeforeIntent與draft candidate；差異必含every-sample→snapshot，不能宣稱等價。
- `ReviewSingleMappingMigration` 接收workspace ID、expected workspace/connector revisions、preview digest、source IDs及明確`confirm_snapshot_conversion=true`。交易內重算preview，任何scope／來源／設定／digest變動均409；blocked來源422，未確認422，foreign/unknown安全404。
- 同一交易保存draft群組、projection、持久化migration ID map及workspace revision；故障全回滾。ID map以workspace/source kind/source ID定位stable group，保存來源revision／原始intent／review digest／adapter version；相同來源revision重跑回同group且不覆寫後續canonical編輯，不重複增workspace revision。來源revision變更要求重新preview，不能重用舊確認。
- review保留原legacy mapping與enabled狀態，applied revision空白；舊writer持續到既有activation barrier下有效Apply。舊read對已移轉資料使用canonical相容projection；舊write不得靜默另成authority（A3.1）。
- 以`LegacySingleMappingMigration`驗證明確確認、stale digest/CAS、來源變動、foreign、blocked、重跑stable IDs／不覆寫、交易故障、legacy讀取一致及review未切writer。A3.2驗證Apply才切owner；C以前不啟用新runtime。

- A2.2 `PreviewRowGroupMigration`／`ReviewRowGroupMigration` 使用同一revision/digest/explicit confirmation envelope，`source_ids`為workspace保存的row-group IDs。每項保留完整row-group與target/source intent；原target mapping與enabled狀態保持原狀，review只保存draft。
- canonical member新增optional `entity_key`，保存legacy target的非空GroupKey，僅作buffer row partition；不同entity可共用target column，同entity/column碰撞blocked。空entity維持原group-scope row；不得由point ID猜SQL business key值。row policy保存`group_key_columns`／`unique_key_columns`規劃metadata；共享column缺unique metadata、缺GroupKey、混合interval／timestamp／target／未支援upsert及模糊來源均blocked。
- migration provenance保存`legacy_row_group_id`及`target_mapping_points`（original target mapping ID → point ID），SourceIDs保留原target IDs；durable ID map仍以workspace/source kind/原row-group ID定位新stable group。workspace projection保留原row-group ID、member與key metadata，不能另增相同legacy布局或遺失refs。canonical edit後舊read按provenance對應member；無法單一表示則actionable409，不回退舊authority。
- row-group preview明列arrival-order→max observed time/tie sample ID、UTC year-one→Unix epoch bucket對齊、late與missing政策；timestamp intent只保留原資料，不猜新SQL binding。一般source rule的已套用signature是provenance，不能單因存在而拒絕；有blocking／pending signature／不明transform則blocked。以LegacyRowGroupMigration驗證shared column分列、同row碰撞／缺key blocked、scope／CAS／來源digest、idempotence與rollback、canonical read投影及沒有writer啟動。C前不交付新row writer。


- A2.3 `PreviewRecordingPlanMigration`／`ReviewRecordingPlanMigration` 使用同一 workspace/connector revision、review digest 與 source IDs envelope；source IDs 是 persisted plan IDs。preview 以一個 local transaction 讀取完整 plan、被引用 measurement 的 workspace/source 身分與 saved connector，回傳 `before_recording_plan_intent`、可重算 source revision、明確 issues／repair action，沒有 candidate。unsupported stream、多 destination、missing/stale source 或缺 column／row identity 均 blocked；foreign/unknown plan 或 foreign source 安全404。
- 本批不創造等價 candidate。review 先驗 revisions/digest；stale409，完整 current request 仍422及`open_write_groups`修復提示，不能用 confirmation 繞過。原 plan、measurement、target、workspace revision、group/map/projection、writer 與外部DB均不改；preview/review 不建表、不 probe、不自動 Apply。完整 before intent 保留於原 persisted plan 與返回的安全預覽；blocked 只描述移轉資格，不覆寫原 plan Status。
- `BasicPlanMigration` 以 single raw_history/every_sample fixture 驗證 blocked／完整原 intent／無候選／原資料可重載；另驗所有 advanced/on_change/sampled/batch/latest、多target、missing/stale/foreign measurement、source/destination變動 digest、review refusal與無任何保存副作用。SDK 僅傳 plan IDs 與同 revision envelope，bounded parser 保留原 plan/source 與修復提示，拒絕被偽造的 needs_review/candidate；不新增 UI editor（E負責）。

- A3.1 global MappingService 的 Create/Update/Delete 與 Studio target 保存都走同一 legacy-write guard；已移轉 mapping ID（single 與 row-group SourceIDs，含 disabled/tombstone）一律拒409。Create 沒有 source ID，依 normalized Tag/connector 判斷：canonical member/current destination，或該 group 保留的原 target mapping Tag/connector，均屬已擁有 scope，不能藉 table/column 或 destination 變更建立第二份 legacy authority。其它 unmigrated scope 保留原 CRUD。
- legacy target 變更先做本地 ownership preflight，再進既有 external table validation；正式寫入在同一 local SQL transaction 重查 guard。global runner 與 canonical mutation 共用 workspace service mutex，Studio 沿既有 setup CAS/mutex/transaction；任何 ownership 變更或拒絕均先於 legacy INSERT/UPDATE/DELETE，不能只有 handler 外的獨立 snapshot。未移轉 global 變更不自行建立 workspace/group 或新增 projection。
- Studio row-group replacement 的 nil 表示未要求替換；explicit `[]` 須檢查 current 的 owned groups。owned group 的移除、layout/member/key/scope 變更拒409；等值保留的 group 與未移轉 scope 可走原 transaction。guard 在 connector/mapping persistence 與 ReplaceDatabaseRowGroups 前；拒絕時 connector/workspace/ref/group/map 全不變，canonical Save/Review 的合法 projection 不受此 legacy guard 阻擋。
- `ErrWriteGroupLegacyWriteConflict` 對外回 `WRITE_GROUP_LEGACY_WRITE_CONFLICT`、HTTP409、`open_write_groups`，安全訊息不帶 credential/DSN/SQL。以 `LegacyWriteAdapterConflict`、owned row-group explicit-empty、unmigrated CRUD、同 transaction ownership 變更／rollback 與 canonical reload 驗證。範圍不含 connector lifecycle、外部 DDL、owner activation 或新 runtime（A3.2/C 負責）。

- A3.2 新增 `WriterOwnershipActivationRequest`，包含 workspace/group IDs、expected database setup/group/connector revisions、previous/new applied revision 與 UTC effective_at；`WriterOwnershipActivationBarrier.ActivateInTx(context.Context, *sql.Tx, request)` 僅接受同一 local transaction，實作者不得 commit、巢狀開 transaction、啟停 runtime 或修改 Share。以 `WithWriterOwnershipActivationBarrier` 注入，沒有 production 注入／Apply route 時保持舊 writer；這個 preparation seam 不表示 durable runtime 已接線。
- `Apply` 先沿既有 readiness preflight，再於 `UpdateDatabaseSetup` 重驗 local CAS/live connector/source，完成新的 immutable version/applied projection 後、local commit 前呼叫 barrier。barrier 可驗 owner/consumer 並寫同一 transaction 的 owner projection；barrier拒絕、projection/workspace save/commit 失敗時，version/applied/owner/workspace 全回滾，caller保留可識別原始 error。版本過期或來源無效時不呼叫 barrier。metadata-only rename 保留 applied identity，不再呼叫 owner transition；下一 UTC cutover與pending semantic version規則保留。
- `WriterOwnershipActivationBarrier` 驗證成功transition在同tx可讀新的immutable version與精確effective_at、stale setup/group/connector/source與未ready不呼叫、barrier錯誤與owner寫入後workspace故障全回滾、rename不重置、原legacy enabled與accepted payload/revisions/receipts不變。另實跑既有Share hydration/settings/token/revision與backlog lifecycle regression；其通過證據與local setup CAS分開。scope 不含新的HTTP Apply端點、production owner consumer、sender、DDL或真正runtime rollback（C負責）。

### Migration Plan

expand schema及adapter → 三批reviewed migration → C接新delivery並以activation切換 → E移除平行editor → F驗收後contract過時write路徑。各批維持build/tests。rollback切回仍相容舊reader，只可停用新intake；已接受backlog與外部效果保留，不能靠版本回退刪除或重送。

### Open Questions

使用者已確認 single-point 的逐筆寫入可經明確審閱轉為週期snapshot，以及現有 recording plans 無法證明等價時保留原 plan／提供 blocked 修復提示；row-group 以已驗證 fixture 保存 row partition，新基本寫入使用 canonical WriteGroup。

## Risks / Trade-offs

[已有舊client持續寫入] → adapter統一CAS並揭露conflict，不維護雙authority。

[現有shared-column規則語意不足] → 保留原資料與legacy路徑，以preview對照同row是否衝突，未證明等價不遷移。

[一份群組模型變成新巨型framework] → 僅基本typed snapshots與已驗證兩種DB，advanced streams不在本案。
