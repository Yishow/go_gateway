## Context

八個舊測試項來自先前跳過清單，不代表八個產品缺陷。使用者要求移除全部項目，需避免刪除父測試依賴的子程序程式。

## Goals / Non-Goals

範圍內是測試註冊與相依 helper 整理。範圍外是產品 API、資料庫驗證、安全 gate、設備通訊與 UI；不連真 PLC。

## Decisions

刪除真 PLC 檔及長負載入口，保留短 pilot。SQLite 不支援組合在註冊前排除。子程序程式改普通 helper，透過已有父測試選入並在父測試最前方分派，子程序不再獨立註冊或空跑成功。SQL Server gate 使用只提供 Helper/Skip 的介面，單元測試以 recorder 核對是否拒絕與拒絕原因，正式 fixture 仍傳入 testing.T。

## Implementation Contract

八個舊身份不再產生測試事件，不用直接 return 偽裝通過。三組父測試仍執行實際終止與恢復；沒有產品層可觀察行為變更。未配置的外部資料庫測試仍如實跳過。

## Risks / Trade-offs

移除負載與現場驗收入口降低該兩種自動測試覆蓋，依使用者指示保留在 Git 歷史而非宣稱已驗證。子程序分派必須在父測試開啟其他資源之前，避免遞迴；以 kill 回歸驗證。

## Migration Plan

單一測試整理提交；若要恢復已刪除的現場與長負載入口，可從父提交取得，不改歷史。Mac 手測沿原產品建置流程。

## Open Questions

無。
