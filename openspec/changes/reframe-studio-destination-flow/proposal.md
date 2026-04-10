# Proposal: Reframe Studio Destination Flow

## Why

The current frontend still models `/studio` Output as one unified workspace where Database and Local Modbus are treated as one final step. That no longer matches the intended product flow.

The new product direction is:

- `/studio` should represent the main product workflow
- `/test` should remain a separate engineering console
- `Point` should no longer be the primary operator object
- `Tag` should become the primary semantic data object
- `Destination` should replace the old Output concept
- `Database` and `Share / Publish` should be parallel destination options

This change also needs to define how a redesigned AI-generated frontend will reconnect to the current backend API contract.

## What Changes

### `/studio`

Reframe the main workflow as:

`Device -> Source -> Tag -> Destination`

Replace the current unified Output step with:

- `Destination Hub`
- `Database Workspace`
- `Share Hub`
- `Local Modbus Workspace`
- `MQTT Workspace`

### `/test`

Keep `/test` as a first-class engineering tool route and formalize it as:

`Field Engineer Debug Console`

### Data model reframing

- `Point` = internal acquisition unit
- `Tag` = primary semantic data object
- `Tag Group` = upstream semantic grouping
- `Delivery Group` = downstream destination grouping

### Integration strategy

The redesign should preserve the current backend contract first, then adapt the new frontend through services, query hooks, and view-model adapters before proposing backend changes.

## Non-Goals

- Do not merge `/test` into `/studio`
- Do not require Database and Share to complete together
- Do not revert to point-first or mapping-first operator flow
- Do not redesign backend APIs first unless the frontend adapter layer proves insufficient
