## Why

目前 frontend-test-suite-stability 將 release acceptance 定義為未設定 worker 的 Vitest default，但同一份 187 個檔案、1000 個測試的 baseline 在 default no-flags NormalFull 出現 991 pass 與 9 個動態 5 秒 timeout；這使 canonical requirement 與可重現的 repository 行為衝突。已有證據顯示平行的 maxWorkers=6 能在同一 inventory 連續通過，且套用後的 repository default 連續三次通過，因此需要以明確 gate 授權一個 bounded parallel default cap，而不是把 diagnosis-only 分支冒充 release acceptance。

## What Changes

- 修改既有 frontend-test-suite-stability capability：在同一 inventory 的 default baseline 失敗、至少兩次 explicit maxWorkers=6 full pass、config contract RED/GREEN 與三次無 CLI worker flags 的 full pass 全部成立後，允許 repository default 設定 maxWorkers=6。
- 保持 release gate 仍以無 CLI worker flags 的完整 Vitest inventory 為 acceptance mode；maxWorkers=6 是 repository default 的平行上限，不是 single-worker 或 reduced-worker diagnosis branch。
- 在 frontend/vite.config.ts 設定 Vitest maxWorkers=6，並在 frontend/tests/unit/utils/viteConfig.test.ts 以 TDD contract test 鎖定該 default。
- 保留既有 Vitest test timeout、hook timeout、isolate、pool、minWorkers、fileParallelism、skip、retry、coverage、dependency 與測試 inventory 行為；不得以提高 timeout、停用平行、排除測試或重試掩蓋問題。
- 在低核心主機上將 6 視為上限而非必須建立的 worker 數；diagnostic runner 的 DefaultGroup、SingleWorkerGroup 與 OneVariable 分支維持原本 diagnosis-only／blocked 語意。

## Capabilities

### New Capabilities

（無）

### Modified Capabilities

- frontend-test-suite-stability：將證據限定的 repository default maxWorkers=6 納入 release acceptance 契約，並保留 default no-CLI full-suite、diagnosis-only 分支與 fail-closed 邊界。

## Impact

- Affected specs: frontend-test-suite-stability
- Affected code:
  - Modified: frontend/vite.config.ts
  - Modified: frontend/tests/unit/utils/viteConfig.test.ts
  - New: 無
  - Removed: 無
- Dependencies: 無
- Scope boundary: 不修改 SourceStep mocks、QueryClient、product frontend code、legacy /studio、diagnostic scripts、dependency manifests 或其他 capability。
