# 測試與交付檢核清單（SmartDashboard Batch5）

日期：2026-03-01  
分支：`feat/datalink-smartdashboard-batch5-i18n`

## 1) 範圍與限制
- Codex：只負責規劃、非前端程式、驗證與交付封裝。
- Gemini：負責全部前端實作（`frontend/src` 與語系檔）。
- 不改 API 契約；所有變更需可回滾。

## 2) test / lint / build 結果欄位
| 類別 | 命令 | 結果（PASS/FAIL/SKIP） | 摘要 | 證據時間（UTC） |
| --- | --- | --- | --- | --- |
| Backend Test | `go test ./...` | FAIL | 執行環境缺少 `go`（`bash: go：指令找不到`） | 2026-03-01T18:29Z |
| Backend Lint | `golangci-lint run ./...` | FAIL | 執行環境缺少 `golangci-lint`（`bash: golangci-lint：指令找不到`） | 2026-03-01T18:29Z |
| Frontend Test | `cd frontend && npm test -- --run` | FAIL | Vitest：`1 failed | 171 passed (172)`，失敗案例 `MemoryGrid.test.tsx > should call onCellClick with point info` | 2026-03-01T18:29Z |
| Frontend Lint | `cd frontend && npm run lint` | PASS | ESLint 完成，退出碼 0 | 2026-03-01T18:29Z |
| Frontend Build | `cd frontend && npm run build` | FAIL | `tsc && vite build` 失敗：共 20 個 TS 錯誤（3 檔） | 2026-03-01T18:29Z |

## 3) 風險欄位
| 風險項目 | 影響 | 緩解策略 | 目前狀態 |
| --- | --- | --- | --- |
| 前端批次改動量大 | 測試與 review 成本上升 | UI 行為與 i18n 拆 commit、逐項驗證 | 追蹤中 |
| i18n key 調整 | 既有測試可能波動 | 優先改穩定 selector，再收斂文案 | 追蹤中 |
| 驗證環境工具缺失（go/golangci-lint） | 後端驗證無法在本機完成 | 於 CI 或補齊工具後重跑 `go test`/`golangci-lint` | 需處理 |
| 前端基線存在測試/型別錯誤 | Gemini 實作完成前無法達成全綠 | 先修復基線失敗，再進行 Batch5 前端改動驗收 | 需處理 |

## 4) PR 交付欄位
| 欄位 | 內容 |
| --- | --- |
| PR 連結 | 待 Gemini 前端實作完成後填寫 |
| 主要 commits | 待本輪提交後填寫 |
| 變更檔案清單 | `docs/next_session_execution_plan_2026-03-02.md`、`docs/test_delivery_checklist.md` |
| API 契約變更 | 無（固定） |
| 回滾方式 | 以小步 commit 回滾；必要時按 commit 範圍還原 |

## 5) 完成回報模板（複製使用）
```text
[Batch5 交付回報]
- Commit：
  - <sha1> <message>
  - <sha2> <message>
- 測試：
  - go test ./... => <PASS/FAIL + 摘要>
  - golangci-lint run ./... => <PASS/FAIL + 摘要>
  - cd frontend && npm test -- --run => <PASS/FAIL + 摘要>
  - cd frontend && npm run lint => <PASS/FAIL + 摘要>
  - cd frontend && npm run build => <PASS/FAIL + 摘要>
- 變更檔案：
  - <file1>
  - <file2>
- 風險：
  - <risk1 + mitigation>
- 下一步：
  - <next action>
```

## 6) 歷史備註（舊任務）
- 舊版 Datalink Gate v1~v1.2 交付內容請參考 git 歷史版本。
