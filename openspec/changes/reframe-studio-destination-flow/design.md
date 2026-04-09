# Design: Reframe Studio Destination Flow

## Overview

This change establishes a new operator model:

- `/studio` is the product workflow
- `/test` is the engineering console
- `Destination` replaces the old Output framing
- `Database` and `Share / Publish` are parallel destination families

## Route model

### Studio routes

- `/studio`
- `/studio?step=destination`
- `/studio/database`
- `/studio/share`
- `/studio/share/local-modbus`
- `/studio/share/mqtt`

### Test route

- `/test`

## Domain object model

### Point

- internal acquisition unit
- source/runtime/diagnostics object
- not the main operator-facing product object

### Tag

- primary semantic data object
- used by Destination planning

### Tag Group

- semantic grouping inside Tag step

### Delivery Group

- destination planning object
- system-suggested, operator-adjusted
- global base object with destination-specific projections

## Workspace model

### Tag Workspace

- left: candidate/tag groups
- center: semantic editor
- right: traceability/preview/readiness

### Database Workspace

- left: delivery groups
- center: row planner
- right: schema/apply preview/validation

### Local Modbus Workspace

- left: delivery groups
- center: register block planner
- right: register map/conflict/apply

### MQTT Workspace

- left: delivery groups
- center: payload/topic planner
- right: preview/validation/apply

## Readiness model

Step readiness remains:

- device
- source
- tag
- destination

`Destination` is complete when at least one valid destination path is configured and apply-ready. It does not require every destination family to be complete.

## `/test` positioning

`/test` remains independent and uses the same visual language as `/studio`, but with a denser, more tool-oriented interaction model.

## Backend integration

The redesigned frontend should continue to use the current backend contract through:

- existing `services/*.ts`
- existing TanStack Query hooks
- new view-model adapters

Direct screen-level raw fetch integration is explicitly discouraged.
