# 儲存、檢查及資料表確認交接

歸屬 streamline-recording-task-editor 3.1；來源 scoped Start 與成功記錄窗口屬後續 3.2。本項沒有新增後端、操作 ledger、持久化 store 或平行請求 owner。

## 單一下一步與原安全邊界

- GroupEditor 的原 handleSave 保持唯一保存權威：fetchLatest／CAS、busy、readback、unknown／late scope 保護不變。新 preparation adapter 只投影下一步。
- 原 GroupSchemaPanel 的 guided commands／schemaOperationRecovery 處理 preview、明確 confirm／Apply、未知操作查詢。保存不自動 preview，preview 不等執行，schema confirmation 不 Start。
- 直接呈現服務回傳的完整 tables（含 receipt table）、database／schema／revision、columns／indexes 與原 SQL details；沒有服務未提供的 layout／index 就不編造。
- 來源／layout／destination 及 native consent generation 異動退役舊 preview／confirm；晚到回覆與舊 captured command 不借新 scope。已知 preview failure 的 A→B→A 明確 retry 保留，未知 DDL 僅查原 operation。
- 原 test-write recovery 的當前 matching references 在 command invocation 再讀，以阻止 render 後新保留的 cleanup reference 被舊 callback 忽略。已知 reference 優先；不可讀但沒有已知 reference 不被新增成必須先成功 test-write 的基本流程門檻。沒有 fake test-write 成功或新 store。
- 原 schema observation unresolved 且有 operation ID 才註冊同一 pending navigation guard；沒有 discard callback，不清 ledger。純 preview 沒有已提交 DDL identity，不誤稱已送出變更。與既有草稿合併同一次離開確認。

## 精確測試及審查

- 新動作／impact／consent／late scope 有真 RED → GREEN；原 preview-failure retry 的單一回歸保留原失敗紀錄，修為獨立 consent fence，而非擴張 owner identity 清掉重試入口。
- 新 owner 整合 12 案例、原 retry 18 案例及最後三個直接 caller 檔案 45 案例通過。Live-reference race 修前 preview／Apply 各誤送一次，修後 15 案例通過，含不可讀空 store 非新必過前置。
- lint／typecheck／locale 舊 tokens／差異檢查通過；GroupEditor 473、GroupSchemaPanel 409 行。選測分次閉合，不宣稱單次全專案通過。
- 獨立唯讀審查檢查精確 baseline／八檔 hash、原 request owner、consent／retry／unknown／pending guard；最終與各窄增量均無剩餘具體 blocker。

## 真 Browser 與修復迭代

QA27 從新 Advanced GroupEditor 明選 managed（原預設 custom），建立合法 VA／VB generic 記錄。新 footer 首次 Save 201 後沒有 preview／DDL；明確 Check 200 顯完整主表／receipt 影響。移除 VB 即退役舊確認，再明確 Save 同 group 新 revision、fresh Check 成功。

QA27 真 Light 發現新增 impact／actions 文字直接覆在舊深色父表面，不可讀，因此沒有 Confirm。只在兩個新元件補定義齊全的 gw surface／text／border／spacing，新增欄位不用低對比單一顏色表意。QA28 Light 可讀，但 390px 的 Confirm 固定單行溢出新卡片，仍沒有 Confirm；僅局部 Button 換行、自高及 max-width 修正。原 FAIL 圖與收據均保留，不推給後續 3.3。

v4 精確 style RED → GREEN、相關 lint 通過；type 未變沿用先前 typecheck。QA29 與 root 實際看 Light 桌面／390px：impact／no-Start 文字可讀，Confirm 完整換行留在卡內，焦點環可見。只有讀清真 impact、fresh preview 確認自有新目標兩個 CREATE 後，才明確 Confirm。

## Schema-only 實際結果

- 續用明確授權的 fixture27 settings／target pair；先有一致備份，其他 28 份舊 DB 不變。沒有複製設定卻誤寫另一個舊 target。
- fresh preview49／200 → schema-apply61／200，operation succeeded、executed_statements=2。
- 本機 SQLite 真建立 run27_preparation_values 與 gw_effect_receipts，含各自 PK 自動索引；兩表各 0 rows。這是建表證據，不是資料寫入成功。
- group 同 ID／revision 仍 draft，applied_revision 空；僅 schema revision／digest 等核驗欄位更新。沒有 source Start、group Apply 或 test write。
- 最後 Back／Step1／Step4／重開 editor 成功，沒有 React #185 或 uncaught exception。自有服務 21:20:00.714 UTC 停止；source362 hash 相符，pair 備份 integrity 通過。

## 必須保留的限制與後續

Submitted unknown／late ACK／scope／live recovery race 是單元／整合證據，沒有在 Browser 偽造 lost ACK。Generic 真 preview 沒有 explicit index／layout metadata，完整 index-list 顯示另有單元覆蓋。

最後正常回 Step1 自動 acquisition-preview GET78 使用 Apply 前的 workspace head 而收到 409，畫面顯 stale／Refresh；這不是刻意負向測試，也不是本項 schema-preview POST 失敗。未操作 Refresh，恢復能力尚未驗證，已交原採集預覽／4.1 終驗診斷，不宣稱預期可恢復或已修。Schema-only 真成功不被此混淆。

原整頁 Light 混色、文字換行及 sticky 遮擋仍在 3.3／4.1；本項僅新 preparation 元件可讀性閉合。3.2 仍須真正 scoped source Start、獨立 mixed 與成功 mean／delta 窗口 SQL/readback。

## 證據識別

- 八檔 frozen-v4 manifest：2dbde5bc726e737df683a2720f3dc314ab6ed6ba7bd6431130f47505b4590afd。
- QA29 更正後報告 SHA256：7510067698ddab4b7eeb96ab54ef4b1a343caf81f9319e36028bbf5194bec2ae；原版本另留歷史，不能引用舊版本冒稱已含 GET78 澄清。
- 證據組：editor31-preparation-evidence、gateway-native-27／28／29；保留所有紅／綠結果、增量 patch、network／SQL、真實截圖與清理紀錄。沒有跑全套或同步 main。
