## Why

四步外觀仍存在，但MQTT選項與位址parser能力不一致、跨設備同位址被誤報、重載完成判定與即時probe不一致；Step 4還要求使用者跨多個模型補設定。需要在既有V2路由內讓四步和同一份persisted readiness一致，並給清楚的修復路徑。

## What Changes

- 保留連線→採集點→Point-to-Tag→群組寫資料庫四步；前三步簡化驗證和反饋，不引入V3或新的framework。
- 協議選項依setup+parser+collector能力交集開放；MQTT未打通時明示不可用，不在本案新寫協議。位址衝突以device＋normalized address/area範圍判定。
- Step 1初次／reload共用persisted validation；offline仍可保存與導航draft，只有相關activation受阻；Step 3基本mapping可直接供write-group使用，進階measurement僅於已選語意功能要求。
- Step 4只編輯write-groups；managed/custom同authority，保留真metadata、schema確認、草稿保護、readonly、CAS與Share獨立性。所有舊方案不匹配時仍可新建或明確修復。
- CommitSummary依後端群組readiness與applied revisions判定，不以另一套enabledTargetCount當managed路線必備。完成畫面不把running/buffered當delivered。

## Capabilities

### New Capabilities

無。

### Modified Capabilities

- `guided-recording-workflow`: 四步基礎流程與可選進階語意分開
- `datalink-workbench-v2-step4-database`: 接手舊active的欄位配對、衝突及引導規格並統一群組authority
- `studio-v2-workspace-readiness`: 持久ready來源、capability與群組revision gate一致

## Impact

frontend/src/features/datalink/workbench-v2/、frontend/src/pages/datalink/workbench-v2/hydratedProgress.ts、frontend/src/utils/addressParser.ts、services/types/hooks、en/zh-TW locales、internal/datalink/workspace/service_readiness*。前置：前四案皆驗證；API語意先定後改UI，避免平行修改同檔。

本次僅起草文件，沒有產品實作。來源、依賴與移交見 [總覽](../../../docs/plans/studio-v2-write-groups/README.md) 及 [現況證據](../../../docs/plans/studio-v2-write-groups/evidence.md)。
