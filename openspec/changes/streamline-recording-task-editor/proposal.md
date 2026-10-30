## Summary

將新記錄入口整合為同頁的「選來源、怎麼記錄、存到哪裡」，需要時就地設定設備與讀取點位。重用現有 autosave、mapping、profile、schema review 與 scoped start，不新增另一套採集或設定引擎。

## Motivation

首頁清楚後，仍需解決跳到不同步驟補設備／點位、重新找原記錄與設定一堆非必要欄位的負擔。點數、暫存器跨度、讀取週期與寫入週期必須分清楚；簡化不能丟失既有協議能力或安全確認。

## Proposed Solution

- 新增 simple-recording-editor 能力：同頁漸進展開來源／記錄／目的地，設備與點位以置中 modal就地編輯，完成後回到原草稿與焦點。
- 修改 guided-recording-workflow 的 Four-step intent-led setup：保留原四步驟進階入口，新增同一 V2 的簡易呈現，不改其安全／資料語意。
- 已有設備及點位直接選取，不要求重建；原始值記錄不強迫設定衍生語意。跨設備與特殊 layout 依真實既有契約給進階入口，不以 demo 宣稱所有組合已支援。
- 保存、檢查、schema impact 確認與 scoped start 區分；修改使相關 readiness/consent 失效，未決結果查狀態，不重送。

## Impact

- Affected specs: simple-recording-editor (new); guided-recording-workflow (modified)
- Affected code:
  - Modified: frontend/src/features/datalink/workbench-v2/steps/step1/ConnectionConfigForm.tsx; frontend/src/features/datalink/workbench-v2/steps/step2/RuleEditor.tsx; frontend/src/features/datalink/workbench-v2/steps/step4/GuidedRecordingPreparation.tsx; frontend/src/pages/datalink/workbench-v2/useStudioV2AutosaveState.ts
  - New: frontend/src/features/datalink/workbench-v2/recording-editor/; frontend/tests/unit/workbench-v2/recording-editor.test.tsx
  - Removed: none
- 依賴 simplify-recording-home-guidance；下批 simplify-recording-lifecycle-recovery 接續已保存記錄的操作。既有完成 fix-step4-recording-authoring 與 guide-step4-recording-workflow 不重做。

本次更新納入 `DESIGN.md` 語意 tokens、system/light/dark 主題、RWD、簡潔文案與字型階層；新彈出設定為 viewport 寬高各75%的置中 modal。已選手機近滿版／小確認框較小例外，分類依 DESIGN.md；舊 demo 不構成此次確認。

新增同工作階段快速新增：同協議帶入上次成功保存的非敏感連線參數，顯示提醒與清除鍵；只在目前workspace/page memory，不新增後端或長期偏好，不帶設備身份／名稱／點位／啟用與安全證據。

已選context help具體化為每個參數控制旁可點擊／觸控／鍵盤操作的圓圈ⓘ；精準欄位說明使用輕量popover，不隱藏必填／錯誤／風險，並提供逐欄coverage map。

已設定目的地／資料表從workspace write-group配置選取，依connector/database/schema/table與新鮮ownership/schema gate判斷；不掃DB。設備常用參數常顯、點位compact逐列；每點字序UI依另立前置契約，不能僅補下拉冒稱runtime支援。

每點字序由獨立前置 preserve-per-point-word-order 提供 durable intent／effective order、rule更新保留override與production decoder契約；本案僅消費其API／capability，完成之前不宣稱每點功能已可用。

本輪v6回饋：不用「每規則單一計算方式」等內部用語，以目的／數字例子解釋實際記錄方式與每點role行為；讀取範圍明確保存成功後關閉modal，失敗保留；編輯記錄初始只顯示已選來源，新增／更換透過有界搜尋選擇器，不列出全部設備與點位。v5未核准，修訂後重新owner checkpoint，不正式apply。

新增明確需求：同一記錄每個member可選最後值／平均／增加量，各member自己的對齊窗口、到期同列；依獨立前置 support-per-member-recording-calculation 完成後接API，不存global Point、不拆成多record冒充完成。


本輪已確認逐member interval與UTC epoch對齊：5秒00/05…55、15秒00/15/30/45、10分鐘00/10…50；wide採due邊界聯集，:15/:30/:45的60秒delta留NULL＋not_due，:00同列各自前15秒mean／前60秒delta；long只新增到期member列。不carry、不補0；required-not-due不skip，到期失敗／無baseline／gap/reset分別標原因。顯示各member的[start,end)、partial coverage與next due，SQL commit時間另列；first partial、nullable schema、安全gate與durable驗收依support-per-member-recording-calculation。此語意確認不是新版demo或apply核准。

採集次數預覽只消費 integrate-protocol-aware-acquisition-batching 的同planner、scoped/revision-bound read-only API；計畫估算與實際協議attempt／retry分開，不能用demo或選取點數宣稱效能改善。此依賴只限預覽／批次能力，不阻擋獨立member時間窗決策。

新增記錄無法新增點位的回饋補入既有有界來源選擇器：空設備可就地新增讀取範圍，新設備／範圍保存後回到可發現新點的選取畫面，不能被已選限定或舊搜尋隱藏；同一modal步驟返回保留記錄草稿、member設定與selection。這是缺口補充，不宣稱已確認使用者的完整重現路徑，也不代表v6 owner checkpoint或正式apply已核准。

本輪選點流程修訂：主頁只管理已加入點與逐member方法／間隔，唯一「選擇點位」開加法picker；同一搜尋／設備／點位類型工具列控制同一結果表，固定footer的待加入數與「取消／加入 N 個點位」一致。pending membership不直接修改record draft，既有點標已加入，主頁移除可復原且重加保留意圖。來源建立收斂為單一次入口、同modal內容切換，目標設備必明確，設備與範圍分別保存；取消不回滾已保存來源。沿原tasks補足驗收，不新增change或任務；正式apply仍等修訂固定demo的owner確認。
