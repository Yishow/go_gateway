# PR Manual Gate（替代保護流程）

> 適用情境：目前 repo 因方案限制無法啟用 Branch Protection。

## 目標
在無分支保護情況下，仍維持「先檢查、再合併」的最低品質門檻。

## Gate 流程
1. PR 作者填寫 PR template（必填檢查項）。
2. Reviewer 至少 1 人確認：
   - 變更範圍與風險等級一致
   - 測試命令與結果可重現
3. 合併前由維護者執行：
   - `gh pr checks <pr-number> --repo Yishow/go_gateway`
   - 若無 checks，改人工核對 template 勾選與測試證據
4. 合併後留存：
   - 於 PR 留言附上最終驗收摘要（完成/風險/下一步）

## 維護者快捷指令
```bash
# 查看 PR 狀態
-gh pr view <pr-number> --repo Yishow/go_gateway

# 查看 PR checks
-gh pr checks <pr-number> --repo Yishow/go_gateway
```

## 備註
- 未來若升級方案支援 Branch Protection，應改採正式保護規則。
