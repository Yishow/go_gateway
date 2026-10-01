<!-- SPECTRA:START v1.3.0 -->

# Spectra Instructions

This project uses Spectra for Spec-Driven Development(SDD). Specs live in `openspec/specs/`, change proposals in `openspec/changes/`.

## Skills

Each `/spectra-*` skill carries its own trigger description; these are the groups:

- Shape and plan → `/spectra-discuss`, `/spectra-propose`
- Continue tasks for an identified change → `/spectra-apply`
- Update requirements or plans for an identified change → `/spectra-ingest`
- Quality gate → `/spectra-verify`, `/spectra-review`, `/spectra-analyze`, `/spectra-audit`, `/spectra-drift`, `/spectra-debug`
- Finish → `/spectra-archive`, `/spectra-commit`

Explicit skill invocation takes precedence. Apply existing authorization within its unchanged scope.

## Workflow

discuss? → propose → apply ⇄ ingest → verify / review → archive

- `discuss` is optional — skip if requirements are clear
- Requirements change mid-work? Plan mode → `ingest` → resume `apply`

## Parked Changes

Changes can be parked（暫存）— temporarily moved out of `openspec/changes/`. Parked changes won't appear in `spectra list` but can be found with `spectra list --parked`. To restore: `spectra unpark <name>`. The `/spectra-apply` and `/spectra-ingest` skills disclose parking and restore when the named operation is already explicitly requested; respect a known refusal, otherwise ask for missing authorization.

<!-- SPECTRA:END -->

# Claude 執行補充

共通規範的單一來源是 [AGENTS.md](AGENTS.md)。本檔與它一起讀，重疊或衝突以 AGENTS.md 為準，不重複建置、安全、命名與測試清單。

## 開始工作

1. 確認任務是研究、proposal、實作或文件整理；只做獲授權的層級。
2. 查 git 狀態與來源基準；保留其他人的修改。
3. 明確指定 Spectra 時讀 `.agents/skills/` 對應 skill。先核對 CLI 及既有 change，再執行該工作流，不用舊機器記憶或自製格式替代。
4. Studio 任務依 AGENTS.md 先讀 inventory 入口，再開實際 route、handler、service wiring 與測試。不要由 archived tasks 的勾選推論現在 production 已接線。

## 調查與交接方式

- 沿使用者操作追到 persisted identity，再追 runtime target writer；UI、API 與 SQL 證據要能對上同一 workspace／revision。
- 遇到 competing model 時，先寫明誰是 authority、哪些只是 projection／相容 adapter，避免新增一套平行狀態。
- 既有 schema operation、readiness 與 Modbus Share 安全契約是保留邊界；不能用簡化 UI 為理由移除。
- 測試與 source 不一致時保留事實及風險，不能為了讓完成報告好看而改掉預期。
- 每次完成回報列出：修改檔案、已執行／未執行的驗證、風險、後續建議。proposal 的完成只表示規格可交接，不表示產品能力已完成。

## 維護本檔

共通政策寫 AGENTS.md；只有 Claude 特定的操作脈絡才放這裡。不要增加不存在的工具、私人備份路徑或與 README 重複的產品介紹。
