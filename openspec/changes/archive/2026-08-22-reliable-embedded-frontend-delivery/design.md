## Context

`go-gateway` packages the React frontend inside the Go binary through `cmd/test_ui/static`. The current frontend direction uses route-level dynamic imports and Vite vendor chunks, so one HTML entry now depends on a graph of hashed JavaScript, CSS, and route assets. A build that deletes the static directory or copies only a subset of `frontend/dist` can remove the tracked directory required by `//go:embed static`, leave stale assets from an earlier build, or ship an embedded binary whose route navigation requests a missing chunk.

The staged implementation is not yet accepted by this change. This design records the contract needed to implement and verify it without changing the established route product decisions: `/studio/v2` remains the default user-facing path, `/studio` remains the full-workbench fallback, `/studio/runtime` remains the focused monitor, `/test` remains an engineering tool, and `/gateway/*` remains experimental.

## Goals / Non-Goals

**Goals:**

- Make route-level lazy loading deterministic, observable, and compatible with the existing route table and redirects.
- Make Vite output a reproducible entry/lazy/vendor asset graph whose relative references remain valid after embedding.
- Make every supported build and start path synchronize the complete `frontend/dist` tree into `cmd/test_ui/static`, remove stale generated entries, and preserve the tracked placeholder.
- Keep fresh-clone compilation possible before any frontend build has produced generated assets.
- Provide focused checks for source build output, synchronization completeness, embedded binary loading, and lazy-load navigation without a chunk 404.
- Make the POSIX launcher fail closed when the frontend source is missing, without invoking the backend build in the same run.
- Make the exact Playwright smoke command start and await its webServer on both Windows and POSIX hosts.
- Keep historical `start.ps1` and `start.sh` line counts at or below their pre-change baselines after repairs.

**Non-Goals:**

- No route redirect or product-flow changes, including no change to `/studio/v2`, `/studio`, `/studio/runtime`, `/test`, or `/gateway/*` semantics.
- No Workbench autosave, mapping, live-values, API, database, connector, or runtime behavior changes.
- No Air entrypoint or lock-algorithm work, pure test-directory migration, lint/type cleanup, or font/browser lockfile update.
- No generated asset commit beyond the tracked placeholder; generated files remain ignored and are recreated by the build or start flow.

## Decisions

### Decision: Use route-level lazy imports behind one deterministic Suspense boundary

Heavy route/page modules are loaded with `React.lazy` and dynamic `import()` at the route boundary. Named exports are adapted to the lazy default-export shape without changing their public route components. The route tree is rendered inside one Suspense boundary with a stable, non-random loading surface that occupies the page and clearly states that the route is loading. Existing redirect elements and path-to-component mappings remain unchanged.

This choice reduces the initial entry payload while preserving the current navigation contract. Per-component boundaries would create multiple competing loading states and make smoke assertions less deterministic; no boundary would expose a blank page while a chunk is pending. The existing error boundary remains the failure surface for an import rejection, so a missing chunk is surfaced as an application error rather than silently treated as a successful route.

Supporting context: `frontend/src/App.tsx`.

### Decision: Keep vendor chunk groups explicit and stable in Vite output

The Vite build uses a `manualChunks` policy for the established dependency families: React/runtime, Material UI and Emotion, charting, router/i18n/query application dependencies, and a residual vendor group. The policy returns stable logical names for matching `node_modules` paths and leaves application modules to Vite's normal route chunks. Hashed filenames remain enabled so a new build cannot accidentally reuse stale content.

Explicit groups are preferred over one monolithic vendor file because they bound cache invalidation and make the generated graph inspectable. A fully automatic split is rejected because dependency upgrades can reshuffle chunk boundaries without a reviewable policy, while manually naming every individual package is too brittle. The policy does not introduce a new runtime base URL or alter route paths.

Supporting context: `frontend/vite.config.ts`.

### Decision: Synchronize generated assets by invariant, not by directory replacement

Every supported synchronization path follows the same invariant:

1. Ensure `cmd/test_ui/static` exists.
2. Remove every generated direct child under that directory except the tracked `embed-placeholder.txt`.
3. Copy the complete contents of `frontend/dist`, including nested directories and all asset types, into the static directory.
4. Fail visibly when the frontend build or source directory is unavailable; do not continue to backend compilation while claiming a successful frontend sync.

The PowerShell paths are `scripts/build.ps1` and the `Build-Frontend` function in `start.ps1`; the POSIX path is `build_frontend` in `start.sh`. The scripts intentionally share the invariant rather than introducing a cross-platform helper that would add another deployment dependency. The placeholder is the only tracked child allowed to survive cleanup, so stale hashed assets cannot mask an incomplete copy and a fresh clone still provides a match for `//go:embed static`.

Supporting context: `.gitignore`, `cmd/test_ui/static/embed-placeholder.txt`, `scripts/build.ps1`, `start.ps1`, and `start.sh`.

### Decision: Verify the asset graph at source, embedded, and browser boundaries

Verification is layered so a successful frontend build cannot be mistaken for embedded runtime acceptance:

- The frontend build check confirms `frontend/dist/index.html` and every referenced static entry exists.
- The synchronization check runs the supported script in a disposable workspace, compares the generated tree with `frontend/dist`, confirms the placeholder survives, and places a deliberately stale generated file to prove cleanup.
- The fresh-clone compile check starts from a checkout containing the tracked placeholder but no generated assets and runs the `cmd/test_ui` Go compilation target.
- The Playwright route smoke in `frontend/tests/e2e/embedded-frontend-delivery.spec.ts` starts the embedded binary, opens `/studio/v2`, `/studio`, `/studio/runtime`, `/test`, and one `/gateway/*` entry, waits for each lazy route to settle, and fails on any asset request returning 404. The assertions check the route's existing identity or redirect, not a new product label.

This layered approach is preferred over inspecting only filename counts because a complete file count can still hide broken relative references or missing dynamic imports. The focused tests are placed in `frontend/tests/unit/app-routing-lazy-load.test.tsx`, `tests/shell/embedded-frontend-delivery.sh`, and `frontend/tests/e2e/embedded-frontend-delivery.spec.ts`; they must assert observable behavior rather than merely assert that a file was edited.

### Decision: Fail closed when the frontend source is missing

The POSIX `build_frontend` path checks that the frontend source directory and the inputs needed to build it exist before it reports success. When that source is unavailable, it returns nonzero and the enclosing build path stops before invoking the backend compiler. The shell fixture adds a `missing-frontend-source` case that removes the source directory, runs the same build-to-backend sequence with a stub compiler marker, and asserts both the failure status and the absence of that marker.

This explicit guard is required because treating a missing frontend as an optional development mode is unsafe for embedded delivery: the backend binary could otherwise be rebuilt while retaining stale or absent static assets. The guard is limited to the build/synchronization path and does not change the normal Vite development-server path when the frontend source is present.

Supporting context: `start.sh` and `tests/shell/embedded-frontend-delivery.sh`.

### Decision: Make Playwright webServer startup cross-platform

`frontend/playwright.config.ts` owns a platform-neutral webServer command and environment setup. The command `npm run test:e2e -- tests/e2e/embedded-frontend-delivery.spec.ts`, launched from the `frontend` working directory, SHALL work without a POSIX-only inline environment assignment or shell chaining. Playwright SHALL wait for the configured URL before the route smoke begins on Windows and POSIX hosts.

Using Playwright's `webServer` environment configuration is preferred over embedding `VITE_DEV_PORT=...` in the command because the latter is not interpreted as an environment assignment by the Windows command shell. A second shell wrapper is rejected because it introduces another platform-specific quoting and process-lifecycle boundary.

Supporting context: `frontend/playwright.config.ts`, `frontend/package.json`, and `frontend/tests/e2e/embedded-frontend-delivery.spec.ts`.

### Decision: Enforce historical launcher line-count ceilings

The final gate compares the current line counts of `start.ps1` and `start.sh` with their pre-change baselines and requires each repaired file to be no longer than its baseline. If a synchronization or source-guard repair increases either historical file, the implementation must refactor within the same file to remove the increase while preserving the behavior contract. The gate runs the repository line-limit checker plus an explicit baseline comparison so a pre-existing over-500-file warning is not mistaken for permission to grow the file.

Supporting context: `start.ps1`, `start.sh`, `scripts/check_file_lines.sh`, and the final content review task.

## Implementation Contract

### Observable behavior

- A fresh checkout can compile `cmd/test_ui` before `frontend/dist` exists because `cmd/test_ui/static/embed-placeholder.txt` keeps the embedded directory matchable.
- A successful frontend build followed by any supported synchronization path produces a static tree containing the complete `frontend/dist` contents. The tracked placeholder remains, and generated files from an earlier build that are absent from the new dist are removed.
- The generated HTML and JavaScript references resolve relative to the embedded static root. Navigating to the existing core routes loads the corresponding lazy chunk and vendor chunks without a 404.
- While a lazy route is pending, the same deterministic fallback is rendered. If a chunk import rejects, the existing error boundary exposes an error state; the route is not presented as successfully loaded.
- If the POSIX frontend source directory or required source inputs are absent, `start.sh` returns nonzero and does not invoke the backend build in that invocation.
- Route redirects and product semantics remain the same: `/` and generic datalink landing continue toward `/studio/v2`, `/studio` remains available as fallback, `/studio/runtime` remains focused on the selected runtime monitor, `/test` remains independent, and `/gateway/*` stays experimental.

### Interface and data shape

- `frontend/src/App.tsx` keeps the existing route path and element contract; only the module loading mechanism changes. Lazy wrappers resolve the same component props and route parameters as the current eager imports.
- `frontend/vite.config.ts` emits relative, hashed asset references and stable vendor chunk labels. No new API, database schema, environment variable, or external runtime service is introduced.
- `frontend/playwright.config.ts` exposes a `webServer` command and environment setup that the exact `npm run test:e2e -- tests/e2e/embedded-frontend-delivery.spec.ts` invocation can start from `frontend` on Windows and POSIX, and its readiness URL is awaited before tests run.
- `cmd/test_ui/static` is the Go embed root. `embed-placeholder.txt` is a tracked text marker; all other children produced by a frontend build are generated deployment assets and remain ignored by Git.
- The synchronization commands continue to be invoked by the existing `scripts/build.ps1`, `start.ps1`, and `start.sh` entrypoints; their success is represented by a zero exit code, and a missing build output or copy failure is represented by a nonzero exit code with an actionable message.

### Failure modes

- If dependency installation or `frontend/dist` generation fails, the owning script stops and does not invoke a backend build as if assets were current.
- If static synchronization cannot create, clean, or copy the target tree, the command returns nonzero and preserves the error in its existing console/log channel.
- If a stale generated file is found after synchronization, the focused asset test fails; the expected state contains only current dist output plus the placeholder.
- If a route chunk or vendor chunk is missing, browser smoke observes a failed network request and the lazy import error boundary; the change does not add a silent fallback that hides the missing asset.
- If the frontend source directory is missing on the POSIX build path, the command surfaces the missing-source error, returns nonzero, and does not continue to a backend build.
- If the Playwright webServer cannot be started by the exact smoke command on either supported host family, the E2E command fails before route assertions instead of silently using an unavailable server.
- If a fresh clone lacks the placeholder, the Go embed compile check fails and identifies the missing static match; the implementation must not solve this by committing generated bundles.
- If either historical launcher exceeds its pre-change line-count baseline, the final gate fails and the file remains unaccepted until the repair is shortened without changing the synchronization contract.

### Acceptance criteria

- The implementation passes the focused route lazy-load unit test in `frontend/tests/unit/app-routing-lazy-load.test.tsx`.
- The frontend build command succeeds and produces an entry HTML plus all lazy/vendor assets required by its import graph.
- `tests/shell/embedded-frontend-delivery.sh` passes for the PowerShell and POSIX synchronization paths, including stale-file removal, placeholder preservation, full tree parity, and nonzero failure behavior for a missing dist directory.
- A disposable fresh-clone check passes `go build ./cmd/test_ui` (or the repository's equivalent `cmd/test_ui` compile command) before any generated frontend assets exist.
- The Playwright spec `frontend/tests/e2e/embedded-frontend-delivery.spec.ts` runs with `cd frontend; npm run test:e2e`, reaches the five route families listed above, observes no chunk or vendor 404, and passes the missing-asset negative case; existing redirect and route identity assertions remain green.
- The exact command `cd frontend; npm run test:e2e -- tests/e2e/embedded-frontend-delivery.spec.ts` starts Playwright's configured webServer and reaches its readiness URL on both Windows and POSIX.
- `start.sh`'s `missing-frontend-source` shell fixture exits nonzero and proves that the backend build marker is absent after the source directory is removed.
- The final line-count gate proves that current `start.ps1` and `start.sh` counts are each less than or equal to their pre-change baselines, and `scripts/check_file_lines.sh` reports no new hard-limit violation.
- The diff is limited to the files named by the proposal and does not include generated bundles, lockfile drift, unrelated Workbench changes, or product-route rewrites.

### Scope boundaries

In scope are lazy loading in `frontend/src/App.tsx`, Vite chunk policy in `frontend/vite.config.ts`, ignore/placeholder governance, asset synchronization in the three named build/start paths, cross-platform webServer startup in `frontend/playwright.config.ts`, the historical line-count gate for `start.ps1` and `start.sh`, and the three focused verification surfaces `frontend/tests/unit/app-routing-lazy-load.test.tsx`, `tests/shell/embedded-frontend-delivery.sh`, and `frontend/tests/e2e/embedded-frontend-delivery.spec.ts`. Out of scope are all route meaning and redirects, Workbench and runtime data behavior, Air entrypoint/lock changes, generated asset commits, and unrelated staged files. No database, API, protocol, or deployment-control contract changes are permitted.

## Risks / Trade-offs

- [Risk] A route chunk can be generated but omitted from the copied static tree. → [Mitigation] Compare the complete dist tree with the synchronized tree and run browser smoke that fails on every asset 404.
- [Risk] A stale chunk can make a local run appear healthy while a clean deployment fails. → [Mitigation] Delete generated static children before copying and assert stale sentinel removal in the shell test.
- [Risk] A slow or failed import can produce an indistinguishable blank screen. → [Mitigation] Use one deterministic Suspense fallback and retain the existing error boundary for rejected imports.
- [Risk] Vendor grouping can increase cache misses when dependency families change. → [Mitigation] Keep groups explicit and reviewable, with hashed filenames and no package-by-package overfitting.
- [Risk] PowerShell and POSIX semantics can drift. → [Mitigation] Keep the same four-step synchronization invariant and exercise both paths in the focused shell verification.

## Migration Plan

1. Apply the lazy route, chunk policy, placeholder/ignore, and synchronization changes within the named files.
2. Run the focused frontend and shell checks, then build the embedded binary from a disposable fresh checkout.
3. Run `cd frontend; npm run test:e2e` against the newly built embedded binary using `frontend/tests/e2e/embedded-frontend-delivery.spec.ts` and record asset request results.
4. Release the binary using the existing packaging path; no database or persistent-data migration is required.
5. Roll back by reverting the change commit and rebuilding the binary; the previous eager imports and synchronization behavior require no data conversion.

## Open Questions

No unresolved design decision remains for this proposal. Browser smoke uses the exact command `cd frontend; npm run test:e2e -- tests/e2e/embedded-frontend-delivery.spec.ts`; `frontend/playwright.config.ts` supplies the cross-platform environment and webServer readiness behavior while preserving the route, network, and failure assertions stated in the implementation contract.
