# Point format consumer contract

This contract supports the existing point form and the recording editor. It does
not add a second format resolver or change the editor layout.

## Editing-start evidence

Use `pointFormatAPI.read({ workspace_id, device_id, point_id, rule_id? })` once
when opening format editing. GET `/api/datalink/points/:id/format` is read-only,
including when the Share revision row is missing. An absent revision row returns
`workspace_revision_origin: initial` with the existing `rev-1` authority; it does
not initialize a row. An explicitly empty `point_revision` is observed legacy
evidence, not permission to omit the field.

`workspace_revision` is the **source/Share revision**, never
`DatabaseSetupRevision`. `connection_revision` is an opaque digest, not credentials.
The response contains no connection configuration or volatile point values.
`editable` combines the read-only SQL snapshot with the same live Share scope
check as Save, using the snapshot's original revisions. Missing, pending or stale
live authority returns `editable: false` while still returning the configuration
resolution. Neither GET nor this check initializes Share state or writes receipts.
This is permission to attempt an edit, not proof that a live sample was decoded.

Render the server's `resolution`: configured mode, effective format, effective
source/policy, applicability, native compatibility and versions. Empty effective
format with native compatibility is meaningful FATEK component assembly, not an
implicit ABCD format. Legacy MC ignores the point's compatibility `data_format`;
legacy FATEK retains component assembly. Do not derive these policies in JavaScript.

## Save, reset and ordinary changes

`createPointFormatOperation(initial, intent, requestId, patch?)` captures a detached,
frozen request. Explicit override uses `{ format_mode: 'override',
configured_format: 'ABCD' | 'BADC' | 'CDAB' | 'DCBA' }`; reset uses
`{ format_mode: 'inherit', configured_format: '' }`. Save uses the original point,
rule and source revisions. Never fetch latest revisions during Save or silently
rebase after a conflict. The backend validates protocol, address and data type.

The actual advanced consumer is `PointForm`, not a component named PointEditor.
It already omits `data_format` and intent fields on ordinary edits; a regression
test protects that behavior. The new format service is ready for recording-editor
integration; this change does not claim a new format UI has been wired.

## Unknown outcomes and recovery

`pointFormatAPI.save(operation)` makes one request. A successful response means
configuration persisted, **not live decoding applied**. `runtime_apply_status`
remains `pending`. Explain:「設定已儲存，執行中的採集仍待重新套用設備設定。」
Do not advise recreating points or repeatedly pressing Save.

A successful invalidation receipt releases only that operation's owned durable
pending marker. It does **not** mark live hydration or acquisition ready. Foreign
`dirty_unknown`, failed callbacks and uncertain commits remain fail-closed.

### Recovery integration still required (4.2 / editor 1.6)

The existing startup path can rebuild canonical source projections after a true
service restart; integration tests exercise the same startup helper used by main.
This proves recovery is possible, but restarting the entire service after every
format edit is **not** the accepted user workflow. The existing activation and
ordinary Share reconcile endpoints require a ready token and cannot currently
serve as pending-format recovery. Do not label them a working recovery action.

The proposed narrow extension is an explicit recovery intent on the existing
Share reconcile endpoint, bound to the original operation receipt and captured
source revisions. It must rebuild from server-authoritative sources, preserve
unknown/foreign dirty and late-operation fences, and issue a fresh readiness token
without starting recording. Activation must still check that fresh token. This
extension and its editor action remain pending safety review and implementation;
clients must not invent it, bypass activation gates, or silently rebase revisions.

For network/503/malformed responses, retain the same operation and request ID.
Use `pointFormatAPI.inspect(operation)` to inspect its durable receipt without
replaying callbacks. A missing receipt (404) does not prove the write failed;
remain unknown, do not auto-resubmit or generate a replacement ID. Foreign scope,
revision conflicts and unsupported intent have bounded, translated safe messages;
raw backend/Axios diagnostics are not rendered.

The `point-format` English and Traditional Chinese namespace supplies mode,
source, policy, pending and error text. Runtime health, managed downgrade and full
production E2E acceptance remain separate gates; this API does not advertise them.
