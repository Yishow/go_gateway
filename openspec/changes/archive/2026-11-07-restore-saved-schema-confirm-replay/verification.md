# 已保存schema-confirm重放驗證

## 範圍

基準c7045b9e，僅handler studio_v2_workspace_write_group_successor_schema.go與router_studio_v2_successor_sqlserver_token_test.go。其他四個fixture檔另案封存。scope為approximated_base，精確兩個source檔案內容已對照實測指紋。

## 規格核對

Exact persisted successor schema confirmation replay 的兩個scenario與具體例子皆有測試：plain/native_monthly成功重放，以及兩模式各5個digest/operation/connector_revision/lineage_revision/alias_token拒絕案例。全部斷言零probe、local/external狀態不變。

實作proofSaved分支檢查canonical Token/Action後，使用既有ConfirmSuccessorSchema。該service仍經successorReservationInTx與VerifiedSchemaProofInTx驗證predecessor、group/lineage/connector CAS、reservation、source digest及完整ledger/token；不以caller聲稱成功作為權威。非proofSaved與首次scope/apply路徑維持原樣。

Completeness：3/3任務；Correctness：1 requirement、2 scenarios皆已映射，沒有未覆蓋例子；Coherence：未新增平行proof authority，沒有放寬第一次確認或activation。

## 本輪實測

- 原真SQLServer離線精確重放先失敗422；修復後完整successor回歸107PASS/0FAIL/0SKIP，package完成時間76.710秒。
- 其中10個新增負向案例全部通過，包含plain/native_monthly各5種錯誤。兩個精確offline replay回200並保留相同operation_id。
- 兩個source SHA256分別為a57bf9c4c956887c4a0037ce2ad57ad54b900c641ffb15ece5ae08a55645ccf9、a8b9c55578781e13ffb7076ef905ab897682f138c93d6fb5e0c0379e918b7f17，與實測終版一致。
- 最後6檔版本的全庫lint 0 issues、vet及產品build exit 0。
- validate/analyze通過且0 findings；diff與行數檢查通過。最後API全套由整合發布gate執行，不將未完成全套提前記成pass。

## Review

已逐diff與相關service依正確性、效率、簡化重用、慣例審查；獨立reviewer再次追到既有持久proof/CAS驗證，未發現具體新增缺陷。alias token/action檢查避免同operation的別名token繞入精確重放。

## 限制

成功重放代表既有已確認結果，不代表目前目的端在線；新activation仍需live readiness。真SQLServer結果不直接宣稱其他引擎已驗證，本次整合API gate與外部DB覆蓋矩陣另行記錄。
