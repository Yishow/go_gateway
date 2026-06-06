## 1. Persisted snapshot recovery

- [ ] 1.1 Deliver Workspace bootstrap hydrates persisted setup state before rendering steps by defining a deterministic bootstrap composition over workspace-owned devices, source rules, mappings, database config, and database targets, and verify it with backend workspace/bootstrap tests plus frontend workspace boot tests.
- [ ] 1.2 Deliver Workspace bootstrap returns a complete persisted setup snapshot under Persisted snapshot is server truth by ensuring persisted identities survive restart and hydrate in one consistent snapshot, and verify it with service restart tests or handler tests that assert stable IDs across reload.

## 2. Step-by-step reload restore behavior

- [ ] 2.1 Deliver Step 1 reload restores persisted device definitions and save truth and Step 2 reload restores persisted source-rule definitions and save truth by hydrating the Step 1 and Step 2 editors from persisted records instead of local defaults, and verify it with device autosave/rule autosave reload-focused unit tests plus a manual refresh assertion.
- [ ] 2.2 Deliver Step 3 reload restores persisted mapping rows against current point identities, Step 4 reload restores persisted database setup and save truth, and Hydration preserves cross-step identities under Hydration runs in dependency order by reconciling mappings and database targets only after current point identities exist, and verify it with mapping/database autosave page tests that reload after persisted saves plus a point-identity restore assertion.

## 3. Draft loss and restored shell truth

- [ ] 3.1 Deliver Reload makes unrecovered drafts explicit under Unrecovered drafts surface as explicit loss by surfacing when a hard refresh discards unsaved local edits, and verify it with a frontend unit test that reloads after save_state is not saved plus a manual assertion that the persisted value is restored with a draft-loss warning.
- [ ] 3.2 Deliver Shell status reflects restored activation truth after reload under Activation summary is persisted separately from UI animation by restoring committed/runtime truth for the shell top bar from persisted activation state instead of reducer defaults, and verify it with shell or step4 commit tests plus a manual refresh check after activation.
