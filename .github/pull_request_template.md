## Summary
- 

## Scope
- [ ] Quick Setup route/entry impact reviewed
- [ ] Expert Workbench route/entry impact reviewed
- [ ] Feature flag behavior (`ENABLE_GATEWAY_DUAL_ENTRY`) verified

## Validation
- [ ] Unit tests added/updated
- [ ] Build passed (`npm run build` in `frontend/`)
- [ ] Flag OFF regression checked (must stay on legacy `/datalink` flow)
- [ ] File line limit check passed (`bash scripts/check_file_lines.sh`)
- [ ] If any file exceeds 300 lines, rationale and split plan documented

## Risks / Rollback
- [ ] Risks documented
- [ ] Rollback path documented (set `ENABLE_GATEWAY_DUAL_ENTRY=false`)
