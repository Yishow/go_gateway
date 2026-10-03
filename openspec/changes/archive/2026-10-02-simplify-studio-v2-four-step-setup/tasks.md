前置：A/B/C/D驗證後順序實作，不平行修改共用Step4/state/locale檔。

所有項目是尚未執行的產品實作。每個行為項先建立具名測試並確認修正前失敗（RED），再實作使其通過（GREEN），最後refactor並重跑；高複雜度狀態機／並行／交易不可跳過此順序。下列測試名是預定驗收標籤，不宣稱已存在。每項保持單次session可驗證；若發現過大，先拆出可獨立驗收的小批並更新依賴，不以未驗收的半成品勾選。

## 1. 前三步一致

- [x] 1.1 以CapabilityAndReloadReadiness測試統一backend能力及persisted device validity的initial/reload validity顯示與activation gate；未打通MQTT顯示不可用原因，connection identity edit使probe失效。OfflineDraftNavigation測試未probe/設備offline仍能Save/Next/Back草稿、只阻擋相關activation，reload保持unverified而非已完成。
- [x] 1.2 [after: 1.1] 以DeviceScopedAddressConflict修改sourceRule衝突key與Step2呈現；兩設備40001合法，同設備normalized area/address/span重疊正確定位，取消/返回/reload不假完成。
- [x] 1.3 [after: 1.2] 以BasicPointTagToGroup簡化Step3名/type/scale與真值preview，保留point/tag/mapping stable IDs；不強制physical measurement語意，保留IME/Enter/blur及dirty drafts保護。

## 2. 單一群組Step4

- [x] 2.1 [after: 1.3] 以WriteGroupEditorManagedAndCustom切換Step4為群組editor；persisted members/target/revisions來自同authority，managed不要求另一套manual db targets，legacy讀取不成第二editor；呈現rename/edit member/disable/delete及backlog影響，tombstone後可查operation/交付結果。
- [x] 2.2 [after: 2.1] 以EmptyErrorMismatchPlanRecovery涵蓋無device/connector/table/group、load failure與所有plan不匹配；提供create/explicit repair/retry而非只有disabled選單，不猜第一筆。
- [x] 2.3 [after: 2.2] 接手舊2.3，以RealMetadataAndReviewableAssignment驗證保留相容配對、無表/forbidden/failed與新表建議，合法shared-column有row identity，跨group/unsafe upsert和不足欄位不自動wrap。
- [x] 2.4 [after: 2.3] 以RevisionSafeGroupActivationAndTestWrite接D的preview/confirm/status；stale回覆不更新新scope，readonly阻擋全部mutation，double-click保留same operation，Share-only仍獨立通過原gate。

## 3. 真實完成與可用性

- [x] 3.1 [after: 2.4] 以GroupCompletionTruth修改CommitSummary/DeliveryTruthStrip，依server readiness與applied revisions；saved/running/buffered/SQL committed/verified/cleanup分開，mixed device/target failure可修復。
- [x] 3.2 [after: 3.1] 接手舊4.1，以GuidedGroupKeyboardAndViewport驗證搜尋/問題篩選及explicit bulk scope/count，390/768/1440 CSS px獨立table scroll、focus/label/error repair可鍵盤操作，新文案完整en/zh-TW。
- [x] 3.3 [after: 3.2] 執行AGENTS前後端相關完整檢查與focused交互回歸，保存實測截圖；確認route identity未變，受影響Studio inventory若更新同步changelog；將actual UI留給F production驗收。
