## 1. Route consolidation

- [ ] 1.1 Redirect legacy `/datalink/workbench` entry points into `/studio`
- [ ] 1.2 Preserve step or selection context on redirect where feasible
- [ ] 1.3 Remove shell-level assumptions that `/datalink/workbench` remains a separate primary workflow

## 2. Guided workflow cleanup

- [ ] 2.1 Update guided workflow content so `/studio` is the sole primary product path
- [ ] 2.2 Demote drag-drop point/tag/mapping tooling to secondary engineering surfaces
- [ ] 2.3 Remove remaining manual-first product language from workbench shell and navigation

## 3. Validation

- [ ] 3.1 Add tests for legacy route redirects into `/studio`
- [ ] 3.2 Add tests proving the normal operator workflow completes inside `/studio`
- [ ] 3.3 Add tests that secondary engineering tools do not replace `/studio` as the primary workflow route
