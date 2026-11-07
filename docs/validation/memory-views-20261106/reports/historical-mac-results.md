> 歷史報告閱覽副本：正規化本機絕對路徑與附件連結；保留當時結論與未驗項目。原始位元組及完整附件見 `local-only/original-evidence.tar.gz`，不要將歷史 PASS 當成新版重驗。

# Mac branch validation, 2026-11-05

Tested detached bfceaa99c32aff49302f39b21d62e66e012c0450 plus verified cloud patches. No local source implementation, commit or push. All changes are in $WORKTREE; original $REPO was not edited. Original main status counts remain147 tracked modifications/100 untracked, not a bytewise user-file audit.

## Passed
- Final full Playwright memory suite118/118, including desktop/tablet density20/20,3 widths×2 themes, navigation/recovery/hover/physical geometry/keyboard/320px layout.
- Standard test-managed Vite starts fixtures correctly using candidate config patch; no external VITE_BUNDLED_DEV workaround.
- First focus patch repeat stability36/36. Browser event trace proved stationary-pointer scroll caused mouseenter overriding keyboard focus; mousemove fix verified.
- Final Mac npm build and Go embedded Gateway build. Large chunk warning only.
- Final production bundle on3223: Modbus30 layouts+6 Close/reopen; MC10 layouts; no page errors, overflow or input clipping. Desktop/tablet all layouts show at least2 full rows initially; tablet rows3/4/2/4/3. Mobile unchanged.
- Real MC and Modbus exact typed values0,-42,12.5,123456.125,65535,-123456. Independent raw TCP readback12/12. Real live update, transport pause→error within3s, resume→fresh37.5→restored12.5,108 SSE events. Both collectors6/6healthy at end.
- All5 original reference image hashes verified locally; visual comparison follows agreed hierarchy/existing themes, not pixel identity.

## Isolation and remaining state
- Gateway final binary:127.0.0.1:3223, PID7994. DisposableDB gateway-test.sqlite retained.
- Isolated HSL PID61612: MC15100,FATEK2100,Modbus1503,management18081,all loopback,empty initial environment/trial mode. No original data used.
- Original HSL PID7396 on15000/2000/1502 was never stopped/seeded/modified by these tests.
- Vite4173 stopped; all dedicated test browsers closed. No Codex CLI or separate Mac Codex chat launched. Plan quota accounting was not independently observable.
- Final source changes10 files, hashes in final-source-sha256.txt. Cloud reported matching key source hashes. Cloud commit/sync decision pending; do not claim HEAD already advanced.

## Evidence
- final-118-report/ and final-118-results/: final fixture run.
- live-hsl/RESULTS.md and live-hsl/production-final/: final real embedded-bundle results.
- live-hsl/both-protocol-results.json and wire-readback.json: authoritative live interruption and packet proof.
- qa-status.json: machine-readable status.
- references/: original images kept outside Git. Screenshots never uploaded or externally shared.

Excluded trials are preserved but not treated as passing evidence: initial Vite bundled-dev fixture routing failures; first patched run serving stale transformed module; HSL initial1000 second timeout fixture; density RED/partial fixes before final green.

## Additional locale coverage
Final production Traditional Chinese matrix passed30/30 plus6/6 Close/reopen. Real Chinese labels verified, all typed values/hover/keyboard passed, desktop/tablet minimum2fullrows in every layout, no overflow or input clipping. See live-hsl/production-zh-TW/. Cloud canonical commit8be1219f56b2e68a8c3a6881f25deef28fd3c1f9 created with matching validated blobs; publication/sync pending.

## Published commit and final synchronization
The isolated Mac worktree now tracks origin on feat/memory-views-validation-pending-20261105, HEAD 1dea0af76b8f3e10dfbf4931aa1dfce5adb4fbfe, tree 3f731f01be5d84d4110a4ea25ba9bb3f34322682, clean status. All 10 candidate blobs matched before staging; staged tree exactly matched published tree, then non-force git switch created the tracking branch. No reset, forced operation or patch discard. Original main remained a9539351 and status counts147/100.
The isolated binary was rebuilt from this clean HEAD: vcs.revision matches and vcs.modified=false. Gateway PID18916 serves3223; original HSL7396 and isolated HSL61612 remain running. Post-restart runtime collectors are both running6/6healthy, exact values retained. The earlier local cloud commit8be1219f was metadata-equivalent; published1dea0af is authoritative.

## Documentation sync and targeted browser supplement
Clean branch HEAD is now 2d437234a9ea76a20f254805cc8bdc486ccc31c6, tree 71e9b1039beaa25dc3bfec6b8c510d655976e160. This fast-forward changes tasks.md only; all validated product/test hashes and running product code remain 1dea0af7. No service restart or product-suite rerun was needed for documentation.
New targeted real-browser coverage passed 30 four-register cross-row cases and 30 cross-page cases across five layouts, three widths and both themes. Values at 41009/41016/41025/41032 cover actual 10/8/5/4-column boundaries; 41249–41252 crosses the 41001–41250 page. Preserved physical addresses/focus, bounded shared spans, continuation owner identity and no false start key. No page errors or navigation-induced Gateway mutations.
Five layout-specific real HSL pushes, 701.125 through 705.125, all reached the UI. Captured 38 value events, one device-scoped point-list request, and maximum one concurrent SSE stream. Three sequential streams reflect changed visible point sets; direct React hook lifecycle was not instrumented/asserted.
All five new seeds restored; only the new wrap collector disabled, configuration retained. Original 12 values unchanged and both original test collectors healthy 6/6. Evidence: four-register-browser/wrap-results.json, live-layout-results.json and final-state.json. No repository source change or upload.
