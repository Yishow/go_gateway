## 1. Backend truth contract cleanup

- [x] 1.1 Deliver Runtime backend expresses degraded and empty states explicitly by extending snapshot and stream contracts with truthful empty or degraded markers, and verify it with runtime handler or stream tests that cover empty and unavailable states.

## 2. Frontend truthful rendering

- [x] 2.1 Deliver Runtime dashboard renders only backend-backed data and Runtime dashboard does not invent unsupported runtime values under Dashboard truthfulness beats visual continuity by removing synthetic panel fallbacks, and verify it with runtime dashboard page/state tests that assert empty and degraded rendering instead of guessed values.
- [x] 2.2 Deliver Workspace runtime view stays aligned to the selected device truth and Runtime dashboard keeps selected device context truthful under Selected device context never silently falls back by refusing cross-device fallback summaries, and verify it with route and dashboard state tests for missing selected-device data.

## 3. Snapshot versus stream behavior

- [x] 3.1 Deliver Snapshot, stream, and degraded fallback stay semantically distinct and Post-setup runtime route shows truthful empty and degraded states under Snapshot and stream failures stay distinct by preserving truthful snapshot state during stream failure without pretending live freshness, and verify it with runtime stream tests plus frontend degraded-live-state assertions.
