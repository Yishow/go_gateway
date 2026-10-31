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

### Explicit recovery on the existing reconcile route (4.2 / editor 1.6)

`POST /datalink/modbus-share/reconcile` accepts **only** a nested
`format_recovery` intent for this flow. It is mutually exclusive with ordinary
reconcile fields, desired mappings, Restore flags or client-supplied formats.
The nested request carries workspace/device/point/rule IDs, the original Save
request ID, **saved-result** point and source/Share workspace revisions, the
original rule/connection revisions, and the settings revision captured when the
operator starts recovery. It never uses database-setup revision or silently
rebases to latest. A previously observed settings revision can be supplied too;
changes reject rather than being accepted implicitly.

A successful response contains `outcome: recovered|replayed`,
`runtime_state: installed|not_started`, `share_applied` and hydration evidence.
This means source configuration was restored, not that a Good measurement exists.
Draft/Disabled devices remain inactive (`not_started`, no invented lease).
Disabled Share remains disabled (`share_applied: false`); no listener or recording
is started. Normal activation still requires the fresh readiness token and its
existing guards. The first accepted recovery issues a token; replay preserves it
and never reinstalls the point. Editor wording should distinguish
「設定已恢復，設備尚未啟用」from「設定已重新套用；等待新的採集結果」.

SQL unknown/foreign dirty, changed scope/defaults, failed hydration or a newer
installation reject the original recovery. Keep the original operation ID and
show a safe refusal; do not retry Save or recreate points. The original receipt's
`runtime_apply_status: pending` records the historical invalidation, not a live
health flag; only this recovery response/current runtime can describe recovery.
The backend implementation is under 4.2 acceptance review; UI integration remains
editor 1.6 and must not be advertised until both gates pass.

For network/503/malformed responses, retain the same operation and request ID.
Use `pointFormatAPI.inspect(operation)` to inspect its durable receipt without
replaying callbacks. A missing receipt (404) does not prove the write failed;
remain unknown, do not auto-resubmit or generate a replacement ID. Foreign scope,
revision conflicts and unsupported intent have bounded, translated safe messages;
raw backend/Axios diagnostics are not rendered.

The `point-format` English and Traditional Chinese namespace supplies mode,
source, policy, pending and error text. Runtime health, managed downgrade and full
production E2E acceptance remain separate gates; this API does not advertise them.
