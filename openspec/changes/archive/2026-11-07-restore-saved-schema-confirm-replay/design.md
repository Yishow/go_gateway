## Context

真SQLServer既有 TestSQLServerSuccessorSchemaProofReplayOffline 在plain/native_monthly中重現422：handler識別proofSaved後仍走h.scope觸發capability probe。service已有完全持久化的精確重放分支。

## Goals / Non-Goals

**Goals:** 已保存proof的完全相同schema-confirm可離線重放，錯誤identity/current CAS仍必須拒絕。

**Non-Goals:** 不變更第一次確認、activation、runtime或來源採集，不新增proof authority，不放寬過期scope或alias token。

## Decisions

### 已保存proof沿用既有確認權威

proofSaved時先核對ledger operation的Token/Action與提供token相同，再呼叫既有ConfirmSuccessorSchema。該service以readSuccessorSchemaSnapshotChecked與verifySuccessorSchemaProof核對持久化完整證據，不做fresh probe；HTTP handler不另建驗證邏輯。

### 非精確重放保持拒絕

增加digest、operation、connector_revision、lineage_revision與alias_token負向案例。每個案例斷言拒絕、零probe、本地及外部狀態不變。

## Implementation Contract

- proofSaved成立時必須滿足成功operation及saved schema operation/digest一致；canonical token/action不一致回ErrSchemaOperationMismatch。
- proofSaved不成立仍走原scope/preview token驗證；首次schema execution路徑不變。
- schema confirm的成功只能保存/讀回schema proof，不更新head、啟用intake或代替activation。
- 回歸是原有真SQLServer plain/native_monthly離線成功與兩模式各5個負向案例；其他dialect由既有service保護與相關API suite驗證，不從SQLServer結果推論全部引擎。

## Risks / Trade-offs

離線重放不是重新確認目前目的端健康，只重放已完成且與當前本地scope一致的持久結果；新activation仍須原有live readiness。外部狀態讀取只用於測試前後對照，產品重放本身不可呼叫外部probe。
