# Git Worktree 工作流程指南

> 本文檔詳細說明如何在 Cursor IDE 中正確使用 Git Worktree 進行多任務並行開發。

## 目錄

- [概念介紹](#概念介紹)
- [流程圖](#流程圖)
- [完整工作流程](#完整工作流程)
- [詳細步驟說明](#詳細步驟說明)
- [最佳實踐](#最佳實踐)
- [常見問題](#常見問題)
- [命令速查表](#命令速查表)

---

## 概念介紹

### 什麼是 Git Worktree？

Git Worktree 允許你在同一個 Git 倉庫中同時檢出多個分支到不同的目錄。這意味著你可以：

- **並行開發多個功能**：不需要 stash 或 commit 未完成的工作
- **同時比較不同分支**：在不同視窗中查看不同版本的代碼
- **獨立的工作環境**：每個 worktree 有自己的工作目錄和暫存區

### 為什麼在 Cursor 中使用 Worktree？

| 優點 | 說明 |
|------|------|
| **多任務並行** | 可以同時處理多個功能或 bug 修復 |
| **AI 輔助隔離** | 每個 Cursor 視窗的 AI 上下文獨立，不會混淆 |
| **快速切換** | 不需要 checkout，直接切換視窗即可 |
| **安全性** | 每個任務在獨立分支，降低衝突風險 |

---

## 流程圖

### 整體工作流程

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           Git Worktree 工作流程                              │
└─────────────────────────────────────────────────────────────────────────────┘

                              ┌──────────────┐
                              │   開始任務   │
                              └──────┬───────┘
                                     │
                                     ▼
                    ┌────────────────────────────────┐
                    │  Step 1: 在主專案創建 Worktree  │
                    │  git worktree add -b <branch>  │
                    └────────────────┬───────────────┘
                                     │
                                     ▼
                    ┌────────────────────────────────┐
                    │  Step 2: 在 Cursor 開啟資料夾  │
                    │  File → Open Folder → 選擇路徑 │
                    └────────────────┬───────────────┘
                                     │
                                     ▼
                    ┌────────────────────────────────┐
                    │  Step 3: 在 Worktree 中開發    │
                    │  編輯、測試、commit            │
                    └────────────────┬───────────────┘
                                     │
                                     ▼
                         ┌───────────────────┐
                         │  開發完成？        │
                         └─────────┬─────────┘
                                   │
                      ┌────────────┴────────────┐
                      │ 否                      │ 是
                      ▼                         ▼
            ┌─────────────────┐    ┌─────────────────────────┐
            │ 繼續開發        │    │ Step 4: Push 分支       │
            │ 返回 Step 3     │    │ git push -u origin HEAD │
            └─────────────────┘    └───────────┬─────────────┘
                                               │
                                               ▼
                                  ┌─────────────────────────┐
                                  │ Step 5: 創建 PR 或合併  │
                                  │ gh pr create / merge    │
                                  └───────────┬─────────────┘
                                               │
                                               ▼
                                  ┌─────────────────────────┐
                                  │ Step 6: 關閉 Cursor 視窗│
                                  └───────────┬─────────────┘
                                               │
                                               ▼
                                  ┌─────────────────────────┐
                                  │ Step 7: 移除 Worktree   │
                                  │ git worktree remove     │
                                  └───────────┬─────────────┘
                                               │
                                               ▼
                                        ┌──────────┐
                                        │   完成   │
                                        └──────────┘
```

### 多 Worktree 並行開發示意圖

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              主專案 (main)                                   │
│                         C:\AIProject\go_gateway                             │
│                                                                             │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                        .git (共享的 Git 目錄)                        │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
                                     │
          ┌──────────────────────────┼──────────────────────────┐
          │                          │                          │
          ▼                          ▼                          ▼
┌─────────────────────┐  ┌─────────────────────┐  ┌─────────────────────┐
│   Worktree 1        │  │   Worktree 2        │  │   Worktree 3        │
│   feature/scanner   │  │   feature/monitor   │  │   bugfix/timeout    │
│                     │  │                     │  │                     │
│ 路徑:               │  │ 路徑:               │  │ 路徑:               │
│ ~/.cursor/worktrees │  │ ~/.cursor/worktrees │  │ ~/.cursor/worktrees │
│ /go_gateway/scanner │  │ /go_gateway/monitor │  │ /go_gateway/timeout │
│                     │  │                     │  │                     │
│ ┌─────────────────┐ │  │ ┌─────────────────┐ │  │ ┌─────────────────┐ │
│ │ Cursor 視窗 1   │ │  │ │ Cursor 視窗 2   │ │  │ │ Cursor 視窗 3   │ │
│ │ 獨立 AI 上下文  │ │  │ │ 獨立 AI 上下文  │ │  │ │ 獨立 AI 上下文  │ │
│ └─────────────────┘ │  │ └─────────────────┘ │  │ └─────────────────┘ │
└─────────────────────┘  └─────────────────────┘  └─────────────────────┘
          │                          │                          │
          │                          │                          │
          └──────────────────────────┼──────────────────────────┘
                                     │
                                     ▼
                          ┌─────────────────────┐
                          │   合併回 main 分支   │
                          │   git merge <branch> │
                          └─────────────────────┘
```

### 分支生命週期

```
時間軸 ──────────────────────────────────────────────────────────────────────►

main     ●────────●────────●────────●────────●────────●────────●────────●
         │        │                          ▲        ▲        ▲
         │        │                          │        │        │
         │        └──────────────────────────┼────────┼────────┘
         │                                   │        │
feature/ │        ●────●────●────●────●──────┘        │
scanner  └────────┘    │    │    │    │               │
                       │    │    │    │               │
                       │    │    │    └── commit 4    │
                       │    │    └─────── commit 3    │
                       │    └──────────── commit 2    │
                       └───────────────── commit 1    │
                                                      │
feature/          ●────●────●────●────────────────────┘
monitor  ─────────┘    │    │    │
                       │    │    └── commit 3
                       │    └─────── commit 2
                       └──────────── commit 1

         ▲             ▲                   ▲          ▲
         │             │                   │          │
    創建分支      開發中...            合併 scanner  合併 monitor
    創建 worktree                     移除 worktree  移除 worktree
```

---

## 完整工作流程

### 階段一：準備工作

```
┌─────────────────────────────────────────────────────────────────┐
│                         準備階段                                 │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  1. 確保主專案是最新狀態                                         │
│     $ cd C:\AIProject\go_gateway                                │
│     $ git checkout main                                         │
│     $ git pull origin main                                      │
│                                                                 │
│  2. 規劃 worktree 存放位置                                       │
│     建議: ~/.cursor/worktrees/<project>/<feature>               │
│     例如: ~/.cursor/worktrees/go_gateway/scanner                │
│                                                                 │
│  3. 確認沒有未提交的變更                                         │
│     $ git status                                                │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 階段二：創建 Worktree

```
┌─────────────────────────────────────────────────────────────────┐
│                      創建 Worktree 階段                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  方式一：創建新分支並建立 worktree（推薦）                        │
│  ─────────────────────────────────────────                      │
│  $ git worktree add <路徑> -b <新分支名>                         │
│                                                                 │
│  範例：                                                          │
│  $ git worktree add ~/.cursor/worktrees/go_gateway/scanner \    │
│        -b feature/device-scanner                                │
│                                                                 │
│  方式二：使用現有分支                                            │
│  ─────────────────────────                                      │
│  $ git worktree add <路徑> <現有分支名>                          │
│                                                                 │
│  範例：                                                          │
│  $ git worktree add ~/.cursor/worktrees/go_gateway/hotfix \     │
│        hotfix/urgent-bug                                        │
│                                                                 │
│  ⚠️ 注意：避免使用 detached HEAD 模式！                          │
│     錯誤示範：git worktree add <路徑>  ← 沒有指定分支            │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 階段三：在 Cursor 中開啟

```
┌─────────────────────────────────────────────────────────────────┐
│                    在 Cursor 開啟 Worktree                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  方法一：使用選單                                                │
│  ─────────────────                                              │
│  1. 開啟新的 Cursor 視窗 (Ctrl+Shift+N)                         │
│  2. File → Open Folder                                          │
│  3. 導航到 worktree 路徑                                         │
│  4. 選擇資料夾並開啟                                             │
│                                                                 │
│  方法二：使用命令列                                              │
│  ─────────────────                                              │
│  $ cursor ~/.cursor/worktrees/go_gateway/scanner                │
│                                                                 │
│  方法三：使用 PowerShell 腳本（批量開啟）                        │
│  ─────────────────────────────────────                          │
│  $ Get-ChildItem ~/.cursor/worktrees/go_gateway | ForEach {     │
│      cursor $_.FullName                                         │
│  }                                                              │
│                                                                 │
│  ✅ 每個 worktree 應該在獨立的 Cursor 視窗中開啟                 │
│  ✅ 這樣 AI 助手的上下文會是獨立的，不會混淆                     │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 階段四：開發工作

```
┌─────────────────────────────────────────────────────────────────┐
│                         開發階段                                 │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  在 Worktree 中進行正常開發：                                    │
│                                                                 │
│  1. 編輯代碼                                                     │
│     - 使用 Cursor AI 輔助開發                                    │
│     - 每個視窗的 AI 上下文獨立                                   │
│                                                                 │
│  2. 測試變更                                                     │
│     $ npm run test                                              │
│     $ go test ./...                                             │
│                                                                 │
│  3. 提交變更                                                     │
│     $ git add .                                                 │
│     $ git commit -m "feat: add device scanner component"        │
│                                                                 │
│  4. 定期同步主分支（避免衝突）                                   │
│     $ git fetch origin main                                     │
│     $ git rebase origin/main                                    │
│     # 或者                                                       │
│     $ git merge origin/main                                     │
│                                                                 │
│  💡 提示：                                                       │
│     - 小步提交，每個 commit 做一件事                             │
│     - 寫清楚的 commit message                                   │
│     - 定期 push 到遠端備份                                       │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 階段五：合併與清理

```
┌─────────────────────────────────────────────────────────────────┐
│                       合併與清理階段                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Step 1: Push 分支到遠端                                         │
│  ─────────────────────────                                      │
│  $ git push -u origin feature/device-scanner                    │
│                                                                 │
│  Step 2: 創建 Pull Request（推薦）                               │
│  ─────────────────────────────────                              │
│  $ gh pr create --title "feat: 新增設備掃描功能" \              │
│        --body "## 變更內容\n- 新增掃描組件\n- 支援所有協議"      │
│                                                                 │
│  或者直接合併（小型變更）                                        │
│  ─────────────────────────                                      │
│  # 在主專案目錄                                                  │
│  $ cd C:\AIProject\go_gateway                                   │
│  $ git checkout main                                            │
│  $ git merge feature/device-scanner                             │
│  $ git push origin main                                         │
│                                                                 │
│  Step 3: 關閉 Cursor 視窗                                        │
│  ─────────────────────────                                      │
│  ⚠️ 重要：先關閉 Cursor 視窗，再移除 worktree                    │
│                                                                 │
│  Step 4: 移除 Worktree                                           │
│  ─────────────────────                                          │
│  $ cd C:\AIProject\go_gateway                                   │
│  $ git worktree remove ~/.cursor/worktrees/go_gateway/scanner   │
│                                                                 │
│  Step 5: 刪除遠端分支（可選）                                    │
│  ─────────────────────────                                      │
│  $ git push origin --delete feature/device-scanner              │
│                                                                 │
│  Step 6: 清理本地分支（可選）                                    │
│  ─────────────────────────                                      │
│  $ git branch -d feature/device-scanner                         │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 詳細步驟說明

### 創建 Worktree 的完整命令解析

```bash
git worktree add <path> -b <branch-name> [<start-point>]
```

| 參數 | 說明 | 範例 |
|------|------|------|
| `<path>` | worktree 的存放路徑 | `~/.cursor/worktrees/go_gateway/scanner` |
| `-b <branch-name>` | 創建新分支的名稱 | `-b feature/device-scanner` |
| `[<start-point>]` | 分支的起始點（可選） | `main` 或 `origin/main` 或 commit hash |

### 範例：創建 4 個並行 Worktree

```bash
# 進入主專案目錄
cd C:\AIProject\go_gateway

# 確保 main 是最新的
git checkout main
git pull origin main

# 創建 4 個 worktree，每個對應不同的功能
git worktree add ~/.cursor/worktrees/go_gateway/scanner -b feature/device-scanner
git worktree add ~/.cursor/worktrees/go_gateway/monitor -b feature/realtime-monitor
git worktree add ~/.cursor/worktrees/go_gateway/auth -b feature/authentication
git worktree add ~/.cursor/worktrees/go_gateway/bugfix -b bugfix/connection-timeout

# 查看所有 worktree
git worktree list
```

輸出範例：
```
C:/AIProject/go_gateway                                    abc1234 [main]
C:/Users/Yishow/.cursor/worktrees/go_gateway/scanner       abc1234 [feature/device-scanner]
C:/Users/Yishow/.cursor/worktrees/go_gateway/monitor       abc1234 [feature/realtime-monitor]
C:/Users/Yishow/.cursor/worktrees/go_gateway/auth          abc1234 [feature/authentication]
C:/Users/Yishow/.cursor/worktrees/go_gateway/bugfix        abc1234 [bugfix/connection-timeout]
```

### 在 Cursor 中開啟所有 Worktree

```powershell
# PowerShell 腳本：批量開啟所有 worktree
$worktrees = git worktree list --porcelain | Select-String "worktree" | ForEach-Object {
    $_.Line -replace "worktree ", ""
}

foreach ($wt in $worktrees) {
    if ($wt -ne (Get-Location).Path) {
        Start-Process cursor -ArgumentList $wt
    }
}
```

---

## 最佳實踐

### ✅ 應該做的事

| 實踐 | 說明 |
|------|------|
| **使用有意義的分支名** | `feature/device-scanner` 而非 `wt1` |
| **每個 worktree 一個功能** | 保持專注，避免混淆 |
| **定期同步 main** | 減少合併衝突 |
| **小步提交** | 每個 commit 做一件事 |
| **先關視窗再刪 worktree** | 避免檔案鎖定問題 |
| **使用 PR 進行合併** | 便於 code review |

### ❌ 不應該做的事

| 避免 | 原因 |
|------|------|
| **使用 detached HEAD** | commit 可能會丟失 |
| **在主專案編輯 worktree 檔案** | 會造成混亂 |
| **同時 checkout 同一分支** | Git 不允許 |
| **忘記 push 就刪除 worktree** | 變更會丟失 |
| **長時間不同步 main** | 會累積大量衝突 |

### 分支命名規範

```
feature/    ← 新功能
bugfix/     ← Bug 修復
hotfix/     ← 緊急修復
refactor/   ← 重構
docs/       ← 文檔更新
test/       ← 測試相關
```

---

## 常見問題

### Q1: Worktree 顯示 "detached HEAD" 怎麼辦？

```bash
# 在 worktree 目錄中
cd ~/.cursor/worktrees/go_gateway/scanner

# 創建並切換到新分支
git checkout -b feature/device-scanner

# 或者切換到現有分支
git checkout feature/device-scanner
```

### Q2: 如何查看所有 worktree？

```bash
git worktree list
```

### Q3: 移除 worktree 時出錯？

```bash
# 強制移除（確保已保存所有變更）
git worktree remove <path> --force

# 清理無效的 worktree 引用
git worktree prune
```

### Q4: 如何在 worktree 之間共享 node_modules？

```bash
# 方法一：使用 pnpm（推薦）
pnpm install  # 自動使用全局 store

# 方法二：使用符號連結
# 在 worktree 中
mklink /D node_modules C:\AIProject\go_gateway\node_modules
```

### Q5: 如何處理合併衝突？

```bash
# 在 worktree 中
git fetch origin main
git rebase origin/main

# 如果有衝突
# 1. 編輯衝突檔案
# 2. git add <resolved-files>
# 3. git rebase --continue
```

---

## 命令速查表

### 創建與管理

| 命令 | 說明 |
|------|------|
| `git worktree add <path> -b <branch>` | 創建新分支的 worktree |
| `git worktree add <path> <branch>` | 使用現有分支創建 worktree |
| `git worktree list` | 列出所有 worktree |
| `git worktree remove <path>` | 移除 worktree |
| `git worktree remove <path> --force` | 強制移除 worktree |
| `git worktree prune` | 清理無效的 worktree 引用 |

### 分支操作

| 命令 | 說明 |
|------|------|
| `git checkout -b <branch>` | 創建並切換分支 |
| `git push -u origin <branch>` | Push 並設置上游分支 |
| `git fetch origin main` | 獲取最新的 main |
| `git rebase origin/main` | 將變更 rebase 到最新 main |
| `git merge origin/main` | 合併最新的 main |

### Cursor 相關

| 命令 | 說明 |
|------|------|
| `cursor <path>` | 在 Cursor 中開啟資料夾 |
| `Ctrl+Shift+N` | 開啟新 Cursor 視窗 |
| `File → Open Folder` | 開啟資料夾 |

---

## 自動化腳本

### create-worktree.ps1

```powershell
# 創建 worktree 的 PowerShell 腳本
param(
    [Parameter(Mandatory=$true)]
    [string]$FeatureName,
    
    [string]$BranchPrefix = "feature"
)

$ProjectRoot = "C:\AIProject\go_gateway"
$WorktreeBase = "$env:USERPROFILE\.cursor\worktrees\go_gateway"
$BranchName = "$BranchPrefix/$FeatureName"
$WorktreePath = "$WorktreeBase\$FeatureName"

# 進入主專案
Set-Location $ProjectRoot

# 更新 main
git checkout main
git pull origin main

# 創建 worktree
git worktree add $WorktreePath -b $BranchName

# 在 Cursor 中開啟
cursor $WorktreePath

Write-Host "✅ Worktree 已創建: $WorktreePath"
Write-Host "✅ 分支: $BranchName"
```

### cleanup-worktree.ps1

```powershell
# 清理 worktree 的 PowerShell 腳本
param(
    [Parameter(Mandatory=$true)]
    [string]$FeatureName,
    
    [switch]$DeleteBranch
)

$ProjectRoot = "C:\AIProject\go_gateway"
$WorktreeBase = "$env:USERPROFILE\.cursor\worktrees\go_gateway"
$WorktreePath = "$WorktreeBase\$FeatureName"

# 進入主專案
Set-Location $ProjectRoot

# 移除 worktree
git worktree remove $WorktreePath --force

# 清理
git worktree prune

if ($DeleteBranch) {
    # 獲取分支名
    $branches = git branch --list "*$FeatureName*"
    foreach ($branch in $branches) {
        $branchName = $branch.Trim().TrimStart("* ")
        git branch -d $branchName
        Write-Host "✅ 已刪除分支: $branchName"
    }
}

Write-Host "✅ Worktree 已清理: $WorktreePath"
```

---

## 總結

使用 Git Worktree 配合 Cursor 可以大幅提升多任務開發效率：

1. **創建** → 使用 `git worktree add -b` 創建帶分支的 worktree
2. **開啟** → 在獨立的 Cursor 視窗中開啟每個 worktree
3. **開發** → 正常開發，定期 commit 和同步 main
4. **合併** → 通過 PR 或直接 merge 合併變更
5. **清理** → 先關閉 Cursor，再移除 worktree

遵循這個流程，可以安全、高效地進行並行開發！
