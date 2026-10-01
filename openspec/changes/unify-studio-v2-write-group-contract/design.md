## Context

已讀proposal與現況證據。Step4的membership、plan selection和enabledTargetCount各自使用不同資料。既有workspace DatabaseRowGroup及target refs需保留；本案只建立單一authority與相容讀寫，尚不啟動新delivery。

現有定位：internal/datalink/workspace/service_database_row_groups.go、internal/api/handlers/studio_v2_workspace_mappings_handler.go、frontend/src/features/datalink/workbench-v2/steps/step4/recordingPlanMembership.ts。

詳細來源見 [evidence](../../../docs/plans/studio-v2-write-groups/evidence.md)，共用欄位與政策見 [跨案契約](../../../docs/plans/studio-v2-write-groups/contracts.md)。這是擬議設計，不是完成報告。

## Goals / Non-Goals

**Goals:** 基本Tag寫入有一份可保存、可重載、可驗證的WriteGroup；existing配置無損相容；新舊不雙寫。

**Non-Goals:** 重寫framework、改採集協議、報表/aggregation/retention產品、外部DDL或現場操作。

## Decisions

1. **單一authority與optional semantics。** 採跨案契約A的欄位。WriteGroup保存persisted Tag與來源point/device，managed/custom僅選layout策略；measurement ID可選。拒絕繼續用recording plan、row group和db target各持一份可獨立修改的設定，也不為基本raw Tag自動編造semantic_kind。既有measurement/plan高階模式仍保留，但非基本readiness前置。
2. **API與CAS。** 擴充現有workspace API，以 `/studio-v2/workspace/write-groups` 提供list/create及 `/:id` get/update/delete；所有mutation要求expected workspace/group（create除外）/connector revisions。一次本地transaction保存group、相容projection及workspace revision；第二段故障全回滾。unknown/foreign資源同安全404，stale409，invalid422。新增API不自行建表、不啟用設備。契約及Swagger在實作批次成套更新。
3. **只讓一份設定可寫。** group啟用前先保存draft；activation在既有workspace/settings/readiness barrier下原子切換applied revision與writer ownership。舊API修改能無損轉成canonical patch則同transaction/CAS更新；不能對應的payload回409及修復動作，不能默默成功寫到過時projection。新delivery接線完成前保持舊runtime路徑，不提早宣稱durable。
4. **分批移轉。** 第一批simple single-point/custom mapping；第二批row-group含shared columns/unique metadata；第三批可精確表達為basic snapshot的recording plan。每批dry-run列出原ID、新ID、scope、row layout、identity與差異；review/apply要同revision。若shared-column將落到同一row、缺unique key、來源模糊、advanced stream或混合目標不能等價則blocked，原配置保持原狀。
5. **保存原意與rollback。** migration provenance及ID map持久化，重跑相同來源revision只得到同一group。每批比較before/after persisted payload與legacy輸出意圖，不把column碰撞一概禁止或一概放行。舊read API從canonical projection回覆已移轉資料。等同一scope的所有consumers已轉接、backlog有明確歸屬且回歸通過後才撤舊write authority，不能先刪表再搬資料。

6. **基本lifecycle。** rename保留group ID及既有record/receipt identity，不重啟已套用writer；所有edit仍走CAS。增刪member、改column/policy/destination保存新draft，未Apply不改running revision；成功Apply在下一個bucket邊界切換，舊bucket用原snapshot收尾。disable停止新intake，已ACK backlog按原snapshot繼續。delete是soft-delete/tombstone並停止新intake；保留不可變payload、revision、ownership、operation/receipt及backlog查詢，worker仍可完成既有delivery；不得重用已刪ID。無法保證backlog歸屬與查詢的舊格式先拒絕delete並提供disable，不能cascade刪accepted記錄。實體purge不屬本案。

### Migration Plan

expand schema及adapter → 三批reviewed migration → C接新delivery並以activation切換 → E移除平行editor → F驗收後contract過時write路徑。各批維持build/tests。rollback切回仍相容舊reader，只可停用新intake；已接受backlog與外部效果保留，不能靠版本回退刪除或重送。

### Open Questions

無需使用者先補答的產品選擇。實作前的migration fixture必須證明legacy row layout可等價映射；不符合者已明定blocked/manual review，不推測。

## Risks / Trade-offs

[已有舊client持續寫入] → adapter統一CAS並揭露conflict，不維護雙authority。

[現有shared-column規則語意不足] → 保留原資料與legacy路徑，以preview對照同row是否衝突，未證明等價不遷移。

[一份群組模型變成新巨型framework] → 僅基本typed snapshots與已驗證兩種DB，advanced streams不在本案。
