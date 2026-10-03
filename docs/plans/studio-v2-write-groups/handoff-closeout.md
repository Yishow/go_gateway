# 舊資料庫流程移交核對

2026-10-03；對象為 `fix-studio-v2-database-workflow` 5.1 與 F 3.2。最終執行結果以 F validation 與 evidence-f 當次 witness 為準；本檔不以文件完成代替產品驗收。

## 唯一責任與原完成紀錄

原五項完成任務 1.2／1.3／2.1／2.2／2.4 的完整文字與 main `4c00af81fb51b3239b4e0d0d9c6b83680c808b36` 比對相同；沒有改写完成歷史。原未完五項由下列 owner 實作，沒有於舊案平行重做。

| 原任務 | 完成 owner | 可核對實作與驗收 |
| --- | --- | --- |
| 2.3 真 metadata／配對／共用欄位 | E 2.3 | `GroupMemberTable`、`columns.ts`、`GroupColumnProposal`；writeGroupMembers／writeGroupProposal 測試涵蓋真實查詢狀態、建議需確認、distinct entity 合法共享與碰撞拒絕 |
| 3.3 防重複試寫 | D 1.1–2.1、3.1–3.2 | grouptestwrite preview／atomic claim／同一 operation ledger；F 真 UI confirm、retained 200／GET、真三分鐘 lease 後 cleanup-only |
| 3.4 型別化 readback／owned cleanup | D 2.2–2.3 | production GroupRowLayout／InsertGroupRow；F SQLite 精確 uint64／neighbor 不變與 PG no-SELECT／no-DELETE；DELETE 後 SELECT 失敗回 unknown，不能宣稱 cleaned |
| 4.1 引導與可用性 | E 2.1–3.2 | 既有四步的 readiness／空清單與讀取失敗／群組 Save/Apply／readonly／scope bulk；E validation 保存 390／768／1440 CSS px 的實際檢查 |
| 4.3 完整整合／browser／DB | F 1.1–3.2 | embedded production UI→兩個真 Modbus reads→SQL；SQLite、PostgreSQL、quality、recovery、capacity／CAS、worker overlap、test-write、Linux container 各自保存 witness |

A/B/C 為 D/E/F 的前置；各自 tasks 與 validation 已歸檔。A–E 共 45 項完成，F 的九項完成狀態以其最終 tasks 為準。

## 保留需求的現行核對

| 保留需求 | 現行實作／測試 | 邊界 |
| --- | --- | --- |
| Protected row editing and single submission | `GroupEditor` explicit Save、local draft、CAS conflict、mounted response guard；writeGroupSection 的 double-click／held draft 與 writeGroupLifecycle 的 late Save A→B 回歸 | E 已移除舊 TargetMappingTable 的逐列 Enter／blur autosave。現行群組欄位沒有鍵盤／blur mutation handler，IME 不觸發 SQL 或 Save；不把已退役 UI 的歷史測試說成當次執行 |
| Fresh schema results and comprehensive readonly controls | `SchemaSetupSection` scope signature／operation 狀態，Step4Database readonly；step4-database、step4-database-controls 與 writeGroupLifecycle 的過期 revision／running／lost response／readonly tests | 忽略遲到回覆只保護畫面，不表示後端操作取消；已開始試寫以同一 operation 查詢 |
| Credential and path safety | stored connector resolver／identity revision；dbtarget connector_identity/service_postgres 與 workspace_database_identity tests | dialect 由 saved connector 決定；password 原 bytes、clear／replace／masked password 與 foreign connector 拒絕維持；fixture witnesses 不保存 bearer token |
| Persisted recording membership and explicit plan selection | `validateRecordingPlanMembership`、workspace-scoped recordingplan CRUD；router membership／members／target resolver tests | 基本群組 UI 不呼叫 advanced plan 建立、選第一筆或猜 measurement ID；advanced API 保留 persisted measurement／device 與 safe404，沒有宣稱 advanced reporting／aggregation 端到端完成 |
| Consistent database setup persistence | workspace database local transaction／expected_setup_revision；AtomicSetupSave 四項 actual SQLite rollback／first-save／stale tests | 外部 DDL／SQL write／readback／cleanup 與本地設定交易獨立，不宣稱跨 DB 原子還原 |
| schema／readiness／Modbus Share | recording schema preview/apply operation tests、workspace schema gate、router_modbus_share／startup hydration／readiness-token tests | 未確認 schema mutation 仍拒絕；既有 Share fail-closed、workspace/settings revisions 保留，未部署或操作真設備 |

## 移轉與回復

- A 的 024 migration 為 additive；single mapping／row-group review 使用 current revisions/digest 與 stable ID map，重跑不覆寫 canonical edits；recording plan migration 維持明確 blocked，不猜轉換。
- A 的 migration／legacy-write tests 與 C 的 025 restart／retained ownership tests在目前完整 Go suite 核對；D 的 026 operation 欄位與既有 schema ledger 同庫保留。
- F 不對正式設定資料做 migration。所有外部 DDL、fault、permission role 與 proxy 只限 run-owned disposable resources；SQLite DB/log 為診斷明列 retained，PG exact-run resources 清除後才發布結果。
- 回退 binary 先停新 intake，保留 immutable revisions、journal／outbox／receipt 與外部 committed effects；不能以舊備份覆蓋新 accepted data，不能盲目補 INSERT。none unknown 需人工核對，endpoint revision 改變後舊 backlog blocked、不轉送新 endpoint。
- 正式備份／restore、部署、Windows／native ARM／native Linux hardware／embedded browser／LAN／真 PLC／SCADA／production DB／長時間 soak 均 NOT RUN。Linux amd64 container 在 arm64 host 的實跑不升格為這些驗收。

## 規格與審查界線

舊案 design 四個原未完主題由本表 D/E/F evidence 完成移交；分析工具的「未出現在本案 tasks」提示不表示缺少唯一 owner。將此移交對照加在 5.1 附註，保留舊 Decisions 與 task 原文。

舊 delta 僅保留上述五條需求；與 E 的群組設定及 D 的確認试寫義務相容。archive 由 Spectra preview／core transaction 一次套用 delta，沒有手動改寫 canonical spec。最終 review／verify 與實跑結果已完成，逐項 scenarios、最後修正、來源身分及未驗邊界見 [final-verification.md](final-verification.md)。
